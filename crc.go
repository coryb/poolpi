package poolpi

import (
	"bytes"
	"encoding/binary"
)

type CRC uint16

func NewCRC(b []byte) CRC {
	return CRC(binary.BigEndian.Uint16(b))
}

func ComputeCRC(b []byte) CRC {
	start := 0
	if bytes.Equal(b[:2], []byte{FrameDLE, FrameSTX}) {
		start = 2
	}
	var crc uint16 = 0x10 + 0x02
	for _, c := range b[start:] {
		crc += uint16(c)
	}
	return CRC(crc)
}

func (crc CRC) ToBytes() []byte {
	out := make([]byte, 2)
	binary.BigEndian.PutUint16(out, uint16(crc))
	return out
}
