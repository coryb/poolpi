package poolpi

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
	"unicode"

	"github.com/coryb/poolpi/pb"
	"github.com/logrusorgru/aurora/v3"
)

type EventType uint16

const (
	EventReady       EventType = 0x0101
	EventLEDs        EventType = 0x0102
	EventMessage     EventType = 0x0103
	EventLongDisplay EventType = 0x040a
	EventPumpRequest EventType = 0x0c01
	EventPumpStatus  EventType = 0x000c
	// EventLocalKey    EventType = 0x0002 // key press on main console
	EventRemoteKey EventType = 0x0003 // key press on wired remote
	// EventWirelessKey EventType = 0x0083 // key press on wireless remote
)

func NewEventType(b []byte) EventType {
	return EventType(binary.BigEndian.Uint16(b[:2]))
}

func (et EventType) ToBytes() []byte {
	seq := make([]byte, 2)
	binary.BigEndian.PutUint16(seq, uint16(et))
	return seq
}

type Event struct {
	Type EventType
	Data []byte
}

func NewEvent(b []byte) Event {
	return Event{
		Type: NewEventType(b[:2]),
		Data: b[2:],
	}
}

func (e Event) Format() string {
	b := e.ToBytes()
	l := len(b)
	// Expected format: DLE+STX CMD[2] DATA[...] CRC[2] DLE+ETX

	// special handling for message, last bit in data is "flags" also they use
	// the high bit to indicate "blink".  For logging if not graphic char
	// (excluding space) just print the hex value.
	if e.Type == EventMessage {
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
		return fmt.Sprintf("[% x] [% x] [%s] [% x] [% x] [% x]",
			b[:2],
			b[2:4],
			strings.Join(data, " "),
			b[l-5:l-4],
			b[l-4:l-2],
			b[l-2:],
		)
	}
	return fmt.Sprintf("[% x] [% x] [% x] [% x] [% x]",
		b[:2],
		b[2:4],
		b[4:l-4],
		b[l-4:l-2],
		b[l-2:],
	)
}

func (e Event) Summary() string {
	pbEvent := e.ToPB()
	switch ev := pbEvent.Event.(type) {
	case *pb.Event_State:
		active := []string{}
		blink := func(s string, should bool) string {
			if !should {
				return s
			}
			return aurora.SlowBlink(s).String()
		}
		if ind := ev.State.GetHeater1(); ind.GetActive() {
			active = append(active, blink("Heater1", ind.GetCaution()))
		}
		if ind := ev.State.GetValve3(); ind.GetActive() {
			active = append(active, blink("Valve3", ind.GetCaution()))
		}
		if ind := ev.State.GetCheckSystem(); ind.GetActive() {
			active = append(active, blink("CheckSystem", ind.GetCaution()))
		}
		if ind := ev.State.GetPool(); ind.GetActive() {
			active = append(active, blink("Pool", ind.GetCaution()))
		}
		if ind := ev.State.GetSpa(); ind.GetActive() {
			active = append(active, blink("Spa", ind.GetCaution()))
		}
		if ind := ev.State.GetFilter(); ind.GetActive() {
			active = append(active, blink("Filter", ind.GetCaution()))
		}
		if ind := ev.State.GetLights(); ind.GetActive() {
			active = append(active, blink("Lights", ind.GetCaution()))
		}
		if ind := ev.State.GetAux1(); ind.GetActive() {
			active = append(active, blink("Aux1", ind.GetCaution()))
		}
		if ind := ev.State.GetAux2(); ind.GetActive() {
			active = append(active, blink("Aux2", ind.GetCaution()))
		}
		if ind := ev.State.GetService(); ind.GetActive() {
			active = append(active, blink("Service", ind.GetCaution()))
		}
		if ind := ev.State.GetAux3(); ind.GetActive() {
			active = append(active, blink("Aux3", ind.GetCaution()))
		}
		if ind := ev.State.GetAux4(); ind.GetActive() {
			active = append(active, blink("Aux4", ind.GetCaution()))
		}
		if ind := ev.State.GetAux5(); ind.GetActive() {
			active = append(active, blink("Aux5", ind.GetCaution()))
		}
		if ind := ev.State.GetAux6(); ind.GetActive() {
			active = append(active, blink("Aux6", ind.GetCaution()))
		}
		if ind := ev.State.GetValve4(); ind.GetActive() {
			active = append(active, blink("Valve4", ind.GetCaution()))
		}
		if ind := ev.State.GetSpillover(); ind.GetActive() {
			active = append(active, blink("Spillover", ind.GetCaution()))
		}
		if ind := ev.State.GetSystemOff(); ind.GetActive() {
			active = append(active, blink("SystemOff", ind.GetCaution()))
		}
		if ind := ev.State.GetAux7(); ind.GetActive() {
			active = append(active, blink("Aux7", ind.GetCaution()))
		}
		if ind := ev.State.GetAux8(); ind.GetActive() {
			active = append(active, blink("Aux8", ind.GetCaution()))
		}
		if ind := ev.State.GetAux9(); ind.GetActive() {
			active = append(active, blink("Aux9", ind.GetCaution()))
		}
		if ind := ev.State.GetAux10(); ind.GetActive() {
			active = append(active, blink("Aux10", ind.GetCaution()))
		}
		if ind := ev.State.GetAux11(); ind.GetActive() {
			active = append(active, blink("Aux11", ind.GetCaution()))
		}
		if ind := ev.State.GetAux12(); ind.GetActive() {
			active = append(active, blink("Aux12", ind.GetCaution()))
		}
		if ind := ev.State.GetAux13(); ind.GetActive() {
			active = append(active, blink("Aux13", ind.GetCaution()))
		}
		if ind := ev.State.GetAux14(); ind.GetActive() {
			active = append(active, blink("Aux14", ind.GetCaution()))
		}
		if ind := ev.State.GetSuperChlorinate(); ind.GetActive() {
			active = append(active, blink("SuperChlorinate", ind.GetCaution()))
		}
		return fmt.Sprintf("State: %v", active)
	case *pb.Event_Message:
		return fmt.Sprintf("Message: %s", ev.Message.Fancy())
	case *pb.Event_PumpRequest:
		return fmt.Sprintf("Pump speed request: %d%%", ev.PumpRequest.SpeedPercent)
	case *pb.Event_PumpStatus:
		return fmt.Sprintf("Pump speed status: %d%% %dW",
			ev.PumpStatus.SpeedPercent,
			ev.PumpStatus.PowerWatts,
		)
	}
	return e.Format()
}

func escDLE(in []byte) []byte {
	count := bytes.Count(in, []byte{FrameDLE})
	if count == 0 {
		return in
	}
	out := make([]byte, len(in)+count)
	offset := 0
	for i, b := range in {
		out[i+offset] = b
		if b == FrameDLE {
			offset++
			out[i+offset] = FrameESC
		}
	}
	return out
}

func (e Event) ToBytes() []byte {
	msg := []byte{FrameDLE, FrameSTX}
	msg = append(msg, e.Type.ToBytes()...)
	msg = append(msg, escDLE(e.Data)...)
	msg = append(msg, ComputeCRC(msg).ToBytes()...)
	msg = append(msg, FrameDLE, FrameETX)
	return msg
}

func (e Event) ToPB() *pb.Event {
	switch e.Type {
	case EventLEDs:
		return &pb.Event{
			Event: &pb.Event_State{
				State: e.toStateEvent(),
			},
		}
	case EventMessage:
		return &pb.Event{
			Event: &pb.Event_Message{
				Message: &pb.MessageEvent{
					Message: e.Data[:len(e.Data)-1],
					Flags:   uint32(e.Data[len(e.Data)-1]),
				},
			},
		}
	case EventPumpRequest:
		speed := binary.BigEndian.Uint16(e.Data)
		return &pb.Event{
			Event: &pb.Event_PumpRequest{
				PumpRequest: &pb.PumpRequestEvent{
					SpeedPercent: uint32(speed),
				},
			},
		}
	case EventPumpStatus:
		speed := e.Data[2]
		power := ((((int(e.Data[3]) & 0xf0) >> 4) * 1000) +
			((int(e.Data[3]) & 0x0f) * 100) +
			(((int(e.Data[4]) & 0xf0) >> 4) * 10) +
			(int(e.Data[4]) & 0x0f))
		return &pb.Event{
			Event: &pb.Event_PumpStatus{
				PumpStatus: &pb.PumpStatusEvent{
					SpeedPercent: uint32(speed),
					PowerWatts:   uint32(power),
				},
			},
		}
	default:
		return &pb.Event{
			Event: &pb.Event_Unknown{
				Unknown: &pb.UnknownEvent{
					Type: e.Type.ToBytes(),
					Data: e.Data,
				},
			},
		}
	}
}

func (e Event) toStateEvent() *pb.StateEvent {
	if e.Type != EventLEDs {
		return nil
	}

	active := e.Data[:4]
	blinking := e.Data[4:8]

	state := pb.StateEvent{}
	bitmasks := [4][8]**pb.Indicator{{
		&state.Heater1, &state.Valve3, &state.CheckSystem, &state.Pool, &state.Spa, &state.Filter, &state.Lights, &state.Aux1,
	}, {
		&state.Aux2, &state.Service, &state.Aux3, &state.Aux4, &state.Aux5, &state.Aux6, &state.Valve4, &state.Spillover,
	}, {
		&state.SystemOff, &state.Aux7, &state.Aux8, &state.Aux9, &state.Aux10, &state.Aux11, &state.Aux12, &state.Aux13,
	}, {
		&state.Aux14, &state.SuperChlorinate,
	}}
	for byteIx, bitmask := range bitmasks {
		for bitIx, indicator := range bitmask {
			if active[byteIx]&(0b1<<bitIx) > 0 || blinking[byteIx]&(0b1<<bitIx) > 0 {
				*indicator = &pb.Indicator{
					Active:  active[byteIx]&(0b1<<bitIx) > 0,
					Caution: blinking[byteIx]&(0b1<<bitIx) > 0,
				}
			}
		}
	}
	return &state
}
