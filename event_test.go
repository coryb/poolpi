package poolpi

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEscDLE(t *testing.T) {
	e := NewEvent([]byte{0xa, 0xb, FrameDLE})
	require.Equal(t, "[10 02] [0a 0b] [10 00] [00 37] [10 03]", e.Format())

	e = NewEvent([]byte{0xa, 0xb, FrameDLE, FrameDLE})
	require.Equal(t, "[10 02] [0a 0b] [10 00 10 00] [00 47] [10 03]", e.Format())

	e = NewEvent([]byte{0xa, 0xb, FrameDLE, 0x0, FrameDLE})
	require.Equal(t, "[10 02] [0a 0b] [10 00 00 10 00] [00 47] [10 03]", e.Format())
}
