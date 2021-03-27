package pb

import (
	"bytes"

	"github.com/logrusorgru/aurora/v3"
)

func (m *MessageEvent) lineWidth() int {
	if m == nil || len(m.Message)%16 == 0 {
		return 16
	}
	return 20
}

func (m *MessageEvent) joinSeparator() string {
	// Most lines are some Name on the top line then a Value on the 2nd line,
	// so we want `Name: Value` but menu names are the edge case where we want
	// `Settings Menu`
	if bytes.Contains(m.Message, []byte("Menu")) {
		return " "
	}
	return ": "
}

func (m *MessageEvent) Plain() (text string) {
	if m == nil {
		return ""
	}

	joinText := m.joinSeparator()
	width := m.lineWidth()
	for i := 0; i < len(m.Message); i += width {
		line := bytes.TrimSpace(m.Message[i : i+width])
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

func (m *MessageEvent) Fancy() (text string) {
	if m == nil {
		return ""
	}

	joinText := m.joinSeparator()
	width := m.lineWidth()
	for i := 0; i < len(m.Message); i += width {
		line := bytes.TrimSpace(m.Message[i : i+width])
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
