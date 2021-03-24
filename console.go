package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"flag"
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

var unknown bool
var withReady bool

func init() {
	flag.BoolVar(&unknown, "unknown", false, "only display unknown event")
	flag.BoolVar(&withReady, "ready", false, "include READY messages in binary log")
}

// ASCII constants for serial transmissions
const FRAME_DLE = 0x10 // "data link escape"
const FRAME_ESC = 0x00 // escape
const FRAME_STX = 0x02 // start
const FRAME_ETX = 0x03 // end

func main() {
	flag.Parse()
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
	c := &serial.Config{Name: "/dev/ttyS0", Baud: 19200, StopBits: serial.Stop2}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}

	sys := NewSystem(s)

	// hack to simulate moving through menus for testing
	// go func() {
	// 	for {
	// 		prompt := sys.keyUntil(KEY_MENU, "Menu")
	// 		if messagePlain(prompt) == "Settings Menu" {
	// 			break
	// 		}
	// 	}
	// 	log.Printf("Got Settings Menu")

	// 	sys.keyUntil(KEY_RIGHT, "Spa")
	// 	prompt := sys.keyUntil(KEY_RIGHT, "Pool")
	// 	log.Printf("Pool Prompt: %s", messageFancy(prompt))
	// }()

	// TODO make system func
	read := func() byte {
		buf := make([]byte, 1, 1)
		for {
			n, err := s.Read(buf)
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
		switch b := read(); b {
		case FRAME_DLE:
			rawData = append(rawData, b)
			switch b := read(); b {
			case FRAME_ESC:
				// escape sequence 0x10 0x00 => 0x10
				rawData = append(rawData, b)
				data = append(data, 0x10)
			case FRAME_STX:
				rawData = append(rawData, b)
				// start
				collectData = true
				continue
			case FRAME_ETX:
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
				eventType := binary.BigEndian.Uint16(data[:2])
				eventData := data[2:]
				if EventType(eventType) != EVENT_READY || withReady {
					rawLog.Print(formatBytes(rawData))
				}
				sys.event(EventType(eventType), eventData)
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
type KeyType uint16

const (
	EVENT_READY        EventType = 0x0101
	EVENT_LEDS                   = 0x0102
	EVENT_MSG                    = 0x0103
	EVENT_LONG_DISPLAY           = 0x040a
	EVENT_PUMP_REQUEST           = 0x0c01
	EVENT_PUMP_STATUS            = 0x000c
	EVENT_LOCAL_KEY              = 0x0002 // key press on main console
	EVENT_REMOTE_KEY             = 0x0003 // key press on wired remote
	EVENT_WIRELESS_KEY           = 0x0083 // key press on wireless remote

	KEY_NONE     KeyType = 0x0000
	KEY_RIGHT            = 0x0001
	KEY_MENU             = 0x0002
	KEY_LEFT             = 0x0004
	KEY_SERVICE          = 0x0008
	KEY_MINUS            = 0x0010
	KEY_PLUS             = 0x0020
	KEY_POOL_SPA         = 0x0040
	KEY_FILTER           = 0x0080
	KEY_LIGHTS           = 0x0100
	KEY_AUX1             = 0x0200
	KEY_AUX2             = 0x0400
	KEY_AUX3             = 0x0800
	KEY_AUX4             = 0x1000
	KEY_AUX5             = 0x2000
	KEY_AUX6             = 0x4000
	KEY_AUX7             = 0x8000
)

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

func decodeKey(key KeyType) string {
	switch key {
	case KEY_NONE:
		return "NONE"
	case KEY_RIGHT:
		return "RIGHT"
	case KEY_MENU:
		return "MENU"
	case KEY_LEFT:
		return "LEFT"
	case KEY_SERVICE:
		return "SERVICE"
	case KEY_MINUS:
		return "MINUS"
	case KEY_PLUS:
		return "PLUS"
	case KEY_POOL_SPA:
		return "POOL_SPA"
	case KEY_FILTER:
		return "FILTER"
	case KEY_LIGHTS:
		return "LIGHTS"
	case KEY_AUX1:
		return "AUX1"
	case KEY_AUX2:
		return "AUX2"
	case KEY_AUX3:
		return "AUX3"
	case KEY_AUX4:
		return "AUX4"
	case KEY_AUX5:
		return "AUX5"
	case KEY_AUX6:
		return "AUX6"
	case KEY_AUX7:
		return "AUX7"
	}
	return "UNKNOWN"
}

type system struct {
	s *serial.Port
	// currentMenu string
	state       map[string]bool
	displayText []byte
	queue       chan []byte
	watchers    sync.Map
}

func NewSystem(s *serial.Port) *system {
	return &system{
		s:     s,
		queue: make(chan []byte, 100),
	}
}

func formatBytes(b []byte) string {
	l := len(b)
	// Expected format: DLE+STX CMD[2] DATA[...] CRC[2] DLE+ETX

	// special handling for message
	if bytes.Equal(b[2:4], []byte{0x01, 0x03}) {
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

func (s *system) display(m []byte) {
	s.displayText = m
	s.watchers.Range(func(_, value interface{}) bool {
		if f, ok := value.(func([]byte)); ok {
			f(m)
		}
		return true
	})
}

func (s *system) keyUntil(key KeyType, expected string) (prompt []byte) {
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
		timeout := time.After(500 * time.Millisecond)
		select {
		case <-done:
			break loop
		case <-timeout:
		}
	}
	return prompt
}

func (s *system) encodeKey(key KeyType) {
	msg := []byte{
		FRAME_DLE, FRAME_STX,
		0x0, 0x0, // type
		0x0, 0x0, 0x0, 0x0, // key
		0x0, 0x0, 0x0, 0x0, // repeat
		0x0, 0x0, //crc
		FRAME_DLE, FRAME_ETX,
	}
	binary.BigEndian.PutUint16(msg[2:], EVENT_REMOTE_KEY)
	// copy(msg[4:8], key)
	// copy(msg[8:12], key)
	binary.BigEndian.PutUint32(msg[4:], uint32(key))
	binary.BigEndian.PutUint32(msg[8:], uint32(key))
	var crc uint16
	for _, c := range msg[:12] {
		crc += uint16(c)
	}
	binary.BigEndian.PutUint16(msg[12:], crc)
	s.queue <- msg
}

func messageWidth(data []byte) int {
	// data is 16/32 or 20/40 + 1 char at end for "display flags"
	if (len(data)-1)%20 == 0 {
		return 20
	}
	return 16
}

func messagePlain(data []byte) (text string) {
	width := messageWidth(data)
	for i := 0; i < len(data)-1; i += width {
		line := bytes.TrimSpace(data[i : i+width])
		if len(line) == 0 {
			continue
		}
		if i > 0 {
			text += ": "
		}
		for _, r := range line {
			text += string(r & 0x7f)
		}
	}
	return text
}

func messageFancy(data []byte) (text string) {
	width := messageWidth(data)
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
				text += string(r & 0x7f)
			}
		}
	}
	return text
}

func (s *system) event(typ EventType, data []byte) {
	maybeLog := log.Printf
	if unknown {
		maybeLog = func(_ string, _ ...interface{}) {}
	}
	switch typ {
	case EVENT_READY:
		// TODO send any outgoing messages now
		select {
		case m := <-s.queue:
			maybeLog("OUTPUT: %x", m)
			s.s.Write(m)
		default:
		}
	case EVENT_LONG_DISPLAY:
	case EVENT_MSG:
		s.display(data[0 : len(data)-1])
		maybeLog("Message: %s", messageFancy(data))
	case EVENT_LEDS:
		// 1st 4 bytes are for light indicators
		// 2nd 4 bytes are for blink indicators
		// maybeLog("Lights: %08b Blink: %08b", data[:4], data[4:])
		leds := decodeLeds(data[:4])
		blinking := decodeLeds(data[4:])
		maybeLog("LEDs: %v  Blinking: %v", leds, blinking)
	case EVENT_PUMP_REQUEST:
		speed := binary.BigEndian.Uint16(data)
		maybeLog("Pump speed request: %d%%", speed)
	case EVENT_PUMP_STATUS:
		// hex := []string{}
		// for _, b := range data {
		// 	hex = append(hex, fmt.Sprintf("%X", b))
		// }
		// maybeLog("PumpStatus %04x Data: %v", typ, hex)
		speed := data[2]
		power := ((((int(data[3]) & 0xf0) >> 4) * 1000) +
			((int(data[3]) & 0x0f) * 100) +
			(((int(data[4]) & 0xf0) >> 4) * 10) +
			(int(data[4]) & 0x0f))
		maybeLog("Pump speed status: %d%% %dW", speed, power)
	case EVENT_REMOTE_KEY:
		// Button Presse events are frequently duplicated
		// we will get [KEY KEY] on the original press
		// and [KEY NONE] on subsequent 100ms triggers where the key is still pressed
		key := binary.BigEndian.Uint32(data[:4])
		upper := binary.BigEndian.Uint32(data[4:])
		if KeyType(upper) == KEY_NONE {
			// only report original keypress
			return
		}
		maybeLog("Button: %s", decodeKey(KeyType(key)))
	default:
		hex := []string{}
		for _, b := range data {
			hex = append(hex, fmt.Sprintf("%X", b))
		}
		log.Printf("Event %04x Data: %v", typ, hex)
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
