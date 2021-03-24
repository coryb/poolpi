package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncodeKey(t *testing.T) {
	s := NewSystem(nil)

	s.encodeKey(KEY_AUX3)
	got := formatBytes(<-s.queue)
	expected := "[10 02] [00 03] [00 08 00 00 00 08 00 00] [00 25] [10 03]"
	require.Equal(t, expected, got)

	s.encodeKey(KEY_MINUS)
	got = formatBytes(<-s.queue)
	expected = "[10 02] [00 03] [00 00 01 00 00 00 01 00] [00 17] [10 03]"
	require.Equal(t, expected, got)

	s.encodeKey(KEY_PLUS)
	got = formatBytes(<-s.queue)
	expected = "[10 02] [00 03] [00 00 02 00 00 00 02 00] [00 19] [10 03]"
	require.Equal(t, expected, got)
}
