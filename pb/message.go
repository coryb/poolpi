package pb

import (
	"bytes"

	"github.com/logrusorgru/aurora/v3"
)

func (m *MessageEvent) Plain() string {
	return plain(m.GetMessage())
}

func (m *MessageEvent) Fancy() string {
	return fancy(m.GetMessage())
}

func (m *MessageUpdateEvent) Plain() string {
	return plain(m.GetMessage().GetMessage())
}

func (m *MessageUpdateEvent) Fancy() string {
	return fancy(m.GetMessage().GetMessage())
}

func plain(b []byte) (text string) {
	joinText := joinSeparator(b)
	width := lineWidth(b)
	for i := 0; i < len(b); i += width {
		line := bytes.TrimSpace(b[i : i+width])
		if len(line) == 0 {
			continue
		}
		if i > 0 {
			text += joinText
		}
		for _, r := range line {
			text += string(r & 0x7f)
		}
	}
	return text
}

func fancy(b []byte) (text string) {
	joinText := joinSeparator(b)
	width := lineWidth(b)
	for i := 0; i < len(b); i += width {
		line := bytes.TrimSpace(b[i : i+width])
		if len(line) == 0 {
			continue
		}
		if i > 0 {
			text += joinText
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

func lineWidth(b []byte) int {
	if len(b)%16 == 0 {
		return 16
	}
	return 20
}

func joinSeparator(b []byte) string {
	// Most lines are some Name on the top line then a Value on the 2nd line,
	// so we want `Name: Value` but menu names are the edge case where we want
	// `Settings Menu`
	if bytes.Contains(b, []byte("Menu")) {
		return " "
	}
	return ": "
}
