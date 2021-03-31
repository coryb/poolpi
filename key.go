package poolpi

import (
	"encoding/binary"

	"github.com/coryb/poolpi/pb"
)

type Key uint32

const (
	KeyNone    Key = 0x00000
	KeyRight   Key = 0x00001
	KeyMenu    Key = 0x00002
	KeyLeft    Key = 0x00004
	KeyService Key = 0x00008
	KeyMinus   Key = 0x00010
	KeyPlus    Key = 0x00020
	KeyPoolSpa Key = 0x00040
	KeyFilter  Key = 0x00080
	KeyLights  Key = 0x00100
	KeyAux1    Key = 0x00200
	KeyAux2    Key = 0x00400
	KeyAux3    Key = 0x00800
	KeyAux4    Key = 0x01000
	KeyAux5    Key = 0x02000
	KeyAux6    Key = 0x04000
	KeyAux7    Key = 0x08000
	KeyValve3  Key = 0x10000
	KeyValve4  Key = 0x20000
	KeyHeater  Key = 0x40000
)

func NewKey(b []byte) Key {
	return Key(binary.LittleEndian.Uint32(b[:4]))
}

func (kt Key) ToBytes() []byte {
	seq := make([]byte, 4)
	binary.LittleEndian.PutUint32(seq, uint32(kt))
	return seq
}

func (kt Key) ToEvent() Event {
	data := kt.ToBytes()
	return Event{
		Type: EventRemoteKey,
		// repeat key bytes twice for protocol
		Data: append(data, data...),
	}
}

func KeyFromPB(key pb.Key) Key {
	switch key {
	case pb.Key_Right:
		return KeyRight
	case pb.Key_Menu:
		return KeyMenu
	case pb.Key_Left:
		return KeyLeft
	case pb.Key_Service:
		return KeyService
	case pb.Key_Minus:
		return KeyMinus
	case pb.Key_Plus:
		return KeyPlus
	case pb.Key_PoolSpa:
		return KeyPoolSpa
	case pb.Key_Filter:
		return KeyFilter
	case pb.Key_Lights:
		return KeyLights
	case pb.Key_Aux1:
		return KeyAux1
	case pb.Key_Aux2:
		return KeyAux2
	case pb.Key_Aux3:
		return KeyAux3
	case pb.Key_Aux4:
		return KeyAux4
	case pb.Key_Aux5:
		return KeyAux5
	case pb.Key_Aux6:
		return KeyAux6
	case pb.Key_Aux7:
		return KeyAux7
	case pb.Key_Valve3:
		return KeyValve3
	case pb.Key_Valve4:
		return KeyValve4
	case pb.Key_Heater:
		return KeyHeater
	}
	return KeyNone
}

func (kt Key) ToPB() pb.Key {
	switch kt {
	case KeyRight:
		return pb.Key_Right
	case KeyMenu:
		return pb.Key_Menu
	case KeyLeft:
		return pb.Key_Left
	case KeyService:
		return pb.Key_Service
	case KeyMinus:
		return pb.Key_Minus
	case KeyPlus:
		return pb.Key_Plus
	case KeyPoolSpa:
		return pb.Key_PoolSpa
	case KeyFilter:
		return pb.Key_Filter
	case KeyLights:
		return pb.Key_Lights
	case KeyAux1:
		return pb.Key_Aux1
	case KeyAux2:
		return pb.Key_Aux2
	case KeyAux3:
		return pb.Key_Aux3
	case KeyAux4:
		return pb.Key_Aux4
	case KeyAux5:
		return pb.Key_Aux5
	case KeyAux6:
		return pb.Key_Aux6
	case KeyAux7:
		return pb.Key_Aux7
	case KeyValve3:
		return pb.Key_Valve3
	case KeyValve4:
		return pb.Key_Valve4
	case KeyHeater:
		return pb.Key_Heater
	}
	return pb.Key_None
}

func (kt Key) String() string {
	return pb.Key_name[int32(kt.ToPB())]
}
