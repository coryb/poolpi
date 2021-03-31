package poolpi

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
	"unicode"

	"github.com/coryb/poolpi/pb"
)

// Update Events have variable content with Maybe EventsLEDs and/or maybe EventMessage
// DLE STX 04a0 83 00 STX EventsLEDs ETX EventMessage CRC DLE ETX
// [10 02] [04 0a] [83 00 03 20 20 20 20 46 69 6c 74 65 72 20 53 70 65 65 64 20 20 20 20 20 20 20 20 20 20 20 20 4f 66 66 20 20 20 20 20 20 20 20 20 00] [09 58] [10 03]
// [04 0a] [83 0 2 8 10 0 0 0 0 0 0] [00]

// is 00 04 from heater????

// [10 02] [04 07] [] [00 1d] [10 03]
// [10 02] [00 04] [05] [00 1b] [10 03]
// [10 02] [04 a0] [] [00 b6] [10 03]
// [10 02] [00 04] [30 33 30 30 20] [00 f9] [10 03]
// [10 02] [04 ae] [] [00 c4] [10 03]
// [10 02] [00 04] [6f 4a] [00 cf] [10 03]

// [10 02] [00 0c] [00 00 5a 16 49] [00 d7] [10 03] => 90% + 1649W
// [10 02] [00 0c] [00 00 32 03 00] [00 53] [10 03] => 50% + 0300W

// Pump speed request: 255%

type EventType uint16

const (
	EventReady       EventType = 0x0101
	EventLEDs        EventType = 0x0102
	EventMessage     EventType = 0x0103
	EventUpdate      EventType = 0x040a
	EventPumpRequest EventType = 0x0c01
	EventPumpStatus  EventType = 0x000c
	EventLocalKey    EventType = 0x0002 // key press on main console
	EventRemoteKey   EventType = 0x0003 // key press on wired remote
	EventWirelessKey EventType = 0x0083 // key press on wireless remote
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

type BinaryFormat string

const (
	FormatHex    BinaryFormat = "[% x]"
	FormatBinary BinaryFormat = "%08b"
)

func (e Event) Format(f BinaryFormat) string {
	b := e.ToBytes()
	l := len(b)
	// Expected format: DLE+STX CMD[2] DATA[...] CRC[2] DLE+ETX

	if e.Type == EventUpdate {
		d := e.Data
		if bytes.HasPrefix(d, []byte{0x83, 0x0, 0x3}) ||
			(bytes.HasPrefix(d, []byte{0x83, 0x0, 0x2}) && len(d) > 12) {
			text := ""
			if l > (40 + 1 + 2 + 2) {
				text = dataToTextIsh(b[l-(40+1+2+2) : l-5])
			}
			return fmt.Sprintf(
				strings.ReplaceAll("[% x] [% x] F [%s] F [% x] [% x]", "F", string(f)),
				b[:2],
				b[2:4],
				b[4:l-(40+1+2+2)],
				text,
				b[l-5:l-4],
				b[l-4:l-2],
				b[l-2:],
			)
		}
	}
	if e.Type == EventMessage {
		text := dataToTextIsh(b[4 : l-5])
		return fmt.Sprintf(
			strings.ReplaceAll("[% x] [% x] [%s] F [% x] [% x]", "F", string(f)),
			b[:2],
			b[2:4],
			text,
			b[l-5:l-4],
			b[l-4:l-2],
			b[l-2:],
		)
	}
	return fmt.Sprintf(
		strings.ReplaceAll("[% x] [% x] F [% x] [% x]", "F", string(f)),
		b[:2],
		b[2:4],
		b[4:l-4],
		b[l-4:l-2],
		b[l-2:],
	)
}

// dataToTextIsh will attempt to decode special ascii characters in the data
// for one of the message events.  For messages the last bit in data is "flags"
// also they use the high bit to indicate "blink".  For logging if not graphic
// char (excluding space) just print the hex value.
func dataToTextIsh(b []byte) string {
	data := []string{}
	for _, c := range b {
		highbit := c & 0x80
		if unicode.In(rune(c&0x7f), unicode.L, unicode.M, unicode.N, unicode.P) {
			s := string(c & 0x7f)
			if highbit > 0 {
				s = "*" + s
			}
			data = append(data, s)
		} else if unicode.In(rune(c&0x7f), unicode.Space) {
			data = append(data, ".")
		} else {
			data = append(data, fmt.Sprintf("%0x", c))
		}
	}
	return strings.Join(data, " ")
}

func (e Event) Summary() string {
	if s := e.ToPB().Summary(); s != "" {
		return s
	}
	return e.Format(FormatHex)
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
				State: toStateEvent(e.Data[:4], e.Data[4:8]),
			},
		}
	case EventMessage:
		l := len(e.Data)
		return &pb.Event{
			Event: &pb.Event_Message{
				Message: &pb.MessageEvent{
					Message: e.Data[:l-1],
					Flags:   uint32(e.Data[l-1]),
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
					RawData:      e.Data,
				},
			},
		}
	case EventUpdate:
		l := len(e.Data)
		// one of:
		// 1) [83 0 2 EventLEDs 3 EventMessage]
		// 2) [83 0 3 EventMessage]
		// 3) [83 0 2 EventLEDs]
		if bytes.HasPrefix(e.Data, []byte{0x83, 0x0, 0x2}) {
			if l == 3+8 {
				// This is 3, just state
				return &pb.Event{
					Event: &pb.Event_StateUpdate{
						StateUpdate: &pb.StateUpdateEvent{
							State: toStateEvent(e.Data[3:7], e.Data[7:11]),
						},
					},
				}
			} else {
				// This is 1, state + message
				return &pb.Event{
					Event: &pb.Event_CurrentState{
						CurrentState: &pb.CurrentStateEvent{
							State: toStateEvent(e.Data[3:7], e.Data[7:11]),
							Message: &pb.MessageEvent{
								Message: e.Data[12 : l-1],
								Flags:   uint32(e.Data[l-1]),
							},
						},
					},
				}

			}
		}
		// This is 2, just message
		return &pb.Event{
			Event: &pb.Event_MessageUpdate{
				MessageUpdate: &pb.MessageUpdateEvent{
					Message: &pb.MessageEvent{
						Message: e.Data[3 : l-1],
						Flags:   uint32(e.Data[l-1]),
					},
				},
			},
		}
	case EventRemoteKey:
		return &pb.Event{
			Event: &pb.Event_Key{
				Key: &pb.KeyEvent{
					Key:    NewKey(e.Data).ToPB(),
					Source: pb.KeySource_Remote,
				},
			},
		}
	case EventLocalKey:
		return &pb.Event{
			Event: &pb.Event_Key{
				Key: &pb.KeyEvent{
					Key:    NewKey(e.Data).ToPB(),
					Source: pb.KeySource_Local,
				},
			},
		}
	case EventWirelessKey:
		l := len(e.Data)
		return &pb.Event{
			Event: &pb.Event_Key{
				Key: &pb.KeyEvent{
					Key:    NewKey(e.Data[1 : l-1]).ToPB(),
					Source: pb.KeySource_Wireless,
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

func EventFromPB(e *pb.Event) Event {
	switch ev := e.Event.(type) {
	case *pb.Event_State:
		return fromStateEvent(ev.State)
	case *pb.Event_Message:
		return Event{
			Type: EventMessage,
			Data: append(ev.Message.Message, byte(ev.Message.Flags)),
		}
	case *pb.Event_PumpRequest:
		data := make([]byte, 2)
		binary.BigEndian.PutUint16(data, uint16(ev.PumpRequest.SpeedPercent))
		return Event{
			Type: EventPumpRequest,
			Data: data,
		}
	case *pb.Event_PumpStatus:
		return Event{
			Type: EventPumpStatus,
			Data: ev.PumpStatus.RawData,
		}
	case *pb.Event_MessageUpdate:
		return Event{
			Type: EventUpdate,
			Data: append(
				append([]byte{0x83, 0x0, 0x3}, ev.MessageUpdate.Message.Message...),
				byte(ev.MessageUpdate.Message.Flags),
			),
		}
	case *pb.Event_StateUpdate:
		tmp := fromStateEvent(ev.StateUpdate.State)
		return Event{
			Type: EventUpdate,
			Data: append([]byte{0x83, 0x0, 0x2}, tmp.Data...),
		}
	case *pb.Event_CurrentState:
		stateTmp := fromStateEvent(ev.CurrentState.State)
		return Event{
			Type: EventUpdate,
			Data: append(
				append(
					append(
						append(
							append([]byte{0x83, 0x0, 0x2}, stateTmp.Data...),
							0x03),
					),
					ev.CurrentState.Message.Message...,
				),
				byte(ev.CurrentState.Message.Flags),
			),
		}
	case *pb.Event_Key:
		switch ev.Key.Source {
		case pb.KeySource_Local:
			return Event{
				Type: EventLocalKey,
				Data: KeyFromPB(ev.Key.Key).ToBytes(),
			}
		case pb.KeySource_Wireless:
			return Event{
				Type: EventWirelessKey,
				Data: append(append([]byte{0x1}, KeyFromPB(ev.Key.Key).ToBytes()...), 0x0),
			}
		default:
			return Event{
				Type: EventRemoteKey,
				Data: KeyFromPB(ev.Key.Key).ToBytes(),
			}
		}
	case *pb.Event_Unknown:
		return Event{
			Type: NewEventType(ev.Unknown.Type),
			Data: ev.Unknown.Data,
		}
	default:
		panic("Invalid proto event")
	}
}

func toStateEvent(active []byte, blinking []byte) *pb.StateEvent {
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

func fromStateEvent(state *pb.StateEvent) Event {
	e := Event{
		Type: EventLEDs,
		Data: make([]byte, 8),
	}
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
			if indicator == nil {
				continue
			}
			if (*indicator).GetActive() {
				e.Data[byteIx] |= 0b1 << bitIx
			}
			if (*indicator).GetCaution() {
				e.Data[byteIx+4] |= 0b1 << bitIx
			}
		}
	}
	return e
}
