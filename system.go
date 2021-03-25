package poolpi

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/logrusorgru/aurora/v3"
	"github.com/rs/xid"
	"github.com/tarm/serial"
)

const FrameDLE = 0x10 // "data link escape"
const FrameESC = 0x00 // escape
const FrameSTX = 0x02 // start
const FrameETX = 0x03 // end

func (s *System) read() byte {
	buf := make([]byte, 1)
	for {
		n, err := s.s.Read(buf)
		if err != nil && !errors.Is(err, io.EOF) {
			log.Fatalf("Got Error: %s", err)
		}
		if n <= 0 {
			log.Printf("No Bytes")
			continue
		}
		return buf[0]
	}
}

func (s *System) Loop() {
	f, err := os.Create("binary.log")
	if err != nil {
		log.Fatalf("Failed to open binary.log: %s", err)
	}
	defer f.Close()
	rawLog := log.New(f, "", log.Ldate|log.Ltime|log.Lmicroseconds)

	data := []byte{}
	rawData := []byte{}
	var crc uint16
	collectData := false
	for {
		if !collectData {
			data = data[:0]
		}
		switch b := s.read(); b {
		case FrameDLE:
			rawData = append(rawData, b)
			switch b := s.read(); b {
			case FrameESC:
				// escape sequence 0x10 0x00 => 0x10
				rawData = append(rawData, b)
				data = append(data, 0x10)
			case FrameSTX:
				rawData = append(rawData, b)
				// start
				collectData = true
				continue
			case FrameETX:
				rawData = append(rawData, b)
				// end
				collectData = false
				if len(data) < 2 {
					rawData = rawData[:0]
					continue
				}
				crc = binary.BigEndian.Uint16(data[len(data)-2:])
				data = data[:len(data)-2]
				var chk uint16 = 0x10 + 0x02
				for _, c := range data {
					chk += uint16(c)
				}
				if crc != chk {
					log.Printf("WARNING: Invalid CRC %d/%d: %v", crc, chk, data)
					rawData = rawData[:0]
					continue
				}
				eventType := newEventType(data[:2])
				eventData := data[2:]
				if eventType != EventReady || s.WithReady {
					rawLog.Print(formatBytes(rawData))
				}
				s.event(eventType, eventData)
				rawData = rawData[:0]
			default:
				log.Printf("Unknown Sequence: 0x10%x", b)
				collectData = false
			}
		default:
			rawData = append(rawData, b)
			if collectData {
				data = append(data, b)
			}
		}
	}
}

type EventType uint16
type KeyType uint32

const (
	EventReady       EventType = 0x0101
	EventLEDs        EventType = 0x0102
	EventMsg         EventType = 0x0103
	EventLongDisplay EventType = 0x040a
	EventPumpRequest EventType = 0x0c01
	EventPumpStatus  EventType = 0x000c
	// EventLocalKey    EventType = 0x0002 // key press on main console
	EventRemoteKey EventType = 0x0003 // key press on wired remote
	// EventWirelessKey EventType = 0x0083 // key press on wireless remote

	KeyNone    KeyType = 0x00000
	KeyRight   KeyType = 0x00001
	KeyMenu    KeyType = 0x00002
	KeyLeft    KeyType = 0x00004
	KeyService KeyType = 0x00008
	KeyMinus   KeyType = 0x00010
	KeyPlus    KeyType = 0x00020
	KeyPoolSpa KeyType = 0x00040
	KeyFilter  KeyType = 0x00080
	KeyLights  KeyType = 0x00100
	KeyAux1    KeyType = 0x00200
	KeyAux2    KeyType = 0x00400
	KeyAux3    KeyType = 0x00800
	KeyAux4    KeyType = 0x01000
	KeyAux5    KeyType = 0x02000
	KeyAux6    KeyType = 0x04000
	KeyAux7    KeyType = 0x08000
	KeyValve3  KeyType = 0x10000
	KeyValve4  KeyType = 0x20000
	KeyHeater  KeyType = 0x40000
)

func newEventType(b []byte) EventType {
	return EventType(binary.BigEndian.Uint16(b[:2]))
}

func (et EventType) ToBytes() []byte {
	seq := make([]byte, 2)
	binary.BigEndian.PutUint16(seq, uint16(et))
	return seq
}

func newKeyType(b []byte) KeyType {
	return KeyType(binary.LittleEndian.Uint32(b[:4]))
}

func (kt KeyType) ToBytes() []byte {
	seq := make([]byte, 4)
	binary.LittleEndian.PutUint32(seq, uint32(kt))
	return seq
}

func (kt KeyType) String() string {
	switch kt {
	case KeyNone:
		return "NONE"
	case KeyRight:
		return "RIGHT"
	case KeyMenu:
		return "MENU"
	case KeyLeft:
		return "LEFT"
	case KeyService:
		return "SERVICE"
	case KeyMinus:
		return "MINUS"
	case KeyPlus:
		return "PLUS"
	case KeyPoolSpa:
		return "POOL_SPA"
	case KeyFilter:
		return "FILTER"
	case KeyLights:
		return "LIGHTS"
	case KeyAux1:
		return "AUX1"
	case KeyAux2:
		return "AUX2"
	case KeyAux3:
		return "AUX3"
	case KeyAux4:
		return "AUX4"
	case KeyAux5:
		return "AUX5"
	case KeyAux6:
		return "AUX6"
	case KeyAux7:
		return "AUX7"
	case KeyValve3:
		return "VALVE3"
	case KeyValve4:
		return "VALVE4"
	case KeyHeater:
		return "HEATER"
	}
	return "UNKNOWN"
}

func decodeLeds(data []byte) []string {
	leds := []string{}
	bitmasks := [4][8]string{{
		"Heater1", "Valve3", "Check System", "Pool", "Spa", "Filter", "Lights", "Aux1",
	}, {
		"Aux2", "Service", "Aux3", "Aux4", "Aux5", "Aux6", "Valve4/Heater2", "Spillover",
	}, {
		"System off", "Aux7", "Aux8", "Aux9", "Aux10", "Aux11", "Aux12", "Aux13",
	}, {
		"Aux14", "Super Chlorinate",
	}}
	for byteIx, bitmask := range bitmasks {
		for bitIx, label := range bitmask {
			if data[byteIx]&(0b1<<bitIx) > 0 {
				leds = append(leds, label)
			}
		}
	}
	return leds
}

type System struct {
	s *serial.Port
	// currentMenu string
	// state       map[string]bool
	displayText []byte
	queue       chan []byte
	watchers    sync.Map

	// FIXME move logging control out
	Unknown   bool
	WithReady bool
}

func NewSystem(s *serial.Port) *System {
	return &System{
		s:     s,
		queue: make(chan []byte, 100),
	}
}

func formatBytes(b []byte) string {
	l := len(b)
	// Expected format: DLE+STX CMD[2] DATA[...] CRC[2] DLE+ETX

	// special handling for message, last bit in data is "flags" also they use
	// the high bit to indicate "blink".  For logging if not graphic char
	// (excluding space) just print the hex value.
	eventType := newEventType(b[2:4])
	if eventType == EventMsg {
		data := []string{}
		for _, c := range b[4 : l-5] {
			highbit := c & 0x80
			if unicode.In(rune(c&0x7f), unicode.L, unicode.M, unicode.N, unicode.P) {
				s := string(c & 0x7f)
				if highbit > 0 {
					s = "*" + s
				}
				data = append(data, s)
			} else {
				data = append(data, fmt.Sprintf("%0x", c))
			}
		}
		return fmt.Sprintf("[% x] [% x] [%s] [% x] [% x] [% x]", b[:2], b[2:4], strings.Join(data, " "), b[l-5:l-4], b[l-4:l-2], b[l-2:])
	}
	return fmt.Sprintf("[% x] [% x] [% x] [% x] [% x]", b[:2], b[2:4], b[4:l-4], b[l-4:l-2], b[l-2:])
}

func (s *System) display(m []byte) {
	s.displayText = m
	s.watchers.Range(func(_, value interface{}) bool {
		if f, ok := value.(func([]byte)); ok {
			f(m)
		}
		return true
	})
}

func (s *System) keyUntil(key KeyType, expected string) (prompt []byte) {
	id := xid.New().String()
	done := make(chan struct{})
	once := sync.Once{}
	s.watchers.Store(id, func(message []byte) {
		if strings.Contains(string(message), expected) {
			once.Do(func() {
				prompt = message
				close(done)
			})
		}
	})
	defer s.watchers.Delete(id)
loop:
	for {
		s.encodeKey(key)
		timeout := time.After(1000 * time.Millisecond)
		select {
		case <-done:
			break loop
		case <-timeout:
		}
	}
	return prompt
}

func escDLE(in []byte) []byte {
	for i, b := range in {
		if b == FrameDLE {
			in = append(in[:i+1], append([]byte{FrameESC}, in[i+1:]...)...)
		}
	}
	return in
}

func crc(b []byte) []byte {
	var crc uint16
	for _, c := range b {
		crc += uint16(c)
	}
	out := make([]byte, 2)
	binary.BigEndian.PutUint16(out, crc)
	return out
}

func (s *System) encodeKey(key KeyType) {
	msg := []byte{FrameDLE, FrameSTX}
	msg = append(msg, EventRemoteKey.ToBytes()...)
	msg = append(msg, escDLE(key.ToBytes())...)
	msg = append(msg, escDLE(key.ToBytes())...)
	msg = append(msg, crc(msg)...)
	msg = append(msg, FrameDLE, FrameETX)
	s.queue <- msg
}

func messageLineWidth(data []byte) int {
	// data is 16/32 or 20/40 + 1 char at end for "display flags"
	if (len(data)-1)%20 == 0 {
		return 20
	}
	return 16
}

func messagePlain(data []byte) (text string) { //nolint:deadcode
	width := messageLineWidth(data)
	for i := 0; i < len(data)-1; i += width {
		line := bytes.TrimSpace(data[i : i+width])
		if len(line) == 0 {
			continue
		}
		if i > 0 {
			text += " "
		}
		for _, r := range line {
			text += string(r & 0x7f)
		}
	}
	return text
}

func messageFancy(data []byte) (text string) {
	width := messageLineWidth(data)
	for i := 0; i < len(data)-1; i += width {
		line := bytes.TrimSpace(data[i : i+width])
		if len(line) == 0 {
			continue
		}
		if i > 0 {
			text += ": "
		}
		for _, r := range line {
			if r&0x80 > 0 {
				text += aurora.SlowBlink(string(r & 0x7f)).String()
			} else {
				if r&0x7f == '_' {
					text += "°"
				} else {
					text += string(r & 0x7f)
				}
			}
		}
	}
	return text
}

func (s *System) event(typ EventType, data []byte) {
	maybeLog := log.Printf
	if s.Unknown {
		maybeLog = func(_ string, _ ...interface{}) {}
	}
	switch typ {
	case EventReady:
		select {
		case m := <-s.queue:
			eventKey := newEventType(m[2:4])
			if eventKey == EventRemoteKey {
				maybeLog("|--> Button: %s", newKeyType(m[4:8]))
			} else {
				maybeLog("|--> %s", formatBytes(m))
			}
			s.s.Write(m)
		default:
		}
	case EventLongDisplay:
		// maybeLog("|<-- Long: %s", formatBytes(data))
	case EventMsg:
		s.display(data)
		maybeLog("|<-- Message: %s", messageFancy(data))
	case EventLEDs:
		// 1st 4 bytes are for light indicators
		// 2nd 4 bytes are for blink indicators
		leds := decodeLeds(data[:4])
		blinking := decodeLeds(data[4:])
		maybeLog("|<-- LEDs: %v  Blinking: %v", leds, blinking)
	case EventPumpRequest:
		speed := binary.BigEndian.Uint16(data)
		maybeLog("|<-- Pump speed request: %d%%", speed)
	case EventPumpStatus:
		speed := data[2]
		power := ((((int(data[3]) & 0xf0) >> 4) * 1000) +
			((int(data[3]) & 0x0f) * 100) +
			(((int(data[4]) & 0xf0) >> 4) * 10) +
			(int(data[4]) & 0x0f))
		maybeLog("|<-- Pump speed status: %d%% %dW", speed, power)
	case EventRemoteKey:
		// Button Press events are frequently repeated we will get [KEY KEY] on
		// the original press and [KEY NONE] on subsequent 100ms triggers where
		// the key is still pressed
		key := newKeyType(data[:4])
		upper := newKeyType(data[4:])
		if upper == KeyNone {
			// only report original keypress
			return
		}
		maybeLog("|<- Button: %s", key)
	default:
		log.Printf("|<-- [% x] [% x]", typ.ToBytes(), data)
	}

	// Jandy Valve?? Spa/Pool
	// Waterfall Pump?
	// Ozone??
	// Heater??
	// Lights??
	// 2021/03/22 05:18:28 Event 0407 Data: []
	// 2021/03/22 05:18:28 Event 0004 Data: [2]
	// 2021/03/22 05:18:28 Event 04a0 Data: []
	// 2021/03/22 05:18:28 Event 0004 Data: [30 33 30 30 20]
	// 2021/03/22 05:18:28 Event 04ae Data: []
	// 2021/03/22 05:18:28 Event 0004 Data: [6f 4a]

	// 2021/03/22 05:56:06 Event 0083 Data: [1 0 0 0 0 0 0 0 0 AB] // remote plugged in
	// 2021/03/22 05:56:28 Event 0083 Data: [1 0 0 0 0 0 0 0 0 FF] // remote unplugged
	// 2021/03/22 09:02:28 Event 0083 Data: [1 0 0 0 0 0 0 0 0 F0] // remote startup?
	// 2021/03/22 09:02:47 Event 0083 Data: [1 0 0 0 0 0 0 0 0 EF]
	// 2021/03/22 09:03:06 Event 0083 Data: [1 0 0 0 0 0 0 0 0 EE]
	// 2021/03/22 09:03:52 Event 0083 Data: [1 0 0 0 0 0 0 0 0 ED]
	// 2021/03/22 09:04:30 Event 0083 Data: [1 0 0 0 0 0 0 0 0 EC]
}
