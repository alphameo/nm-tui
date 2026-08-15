package styles

import (
	"bytes"
	"fmt"
	"image/color"
	"os"
	"time"

	"github.com/charmbracelet/x/ansi"
	"github.com/charmbracelet/x/term"
)

// queryTerminalForegroundColor sends an OSC 10 query to the controlling
// terminal and reads back the reply, returning the color as a "#rrggbb" hex
// string. It fails gracefully when the process has no terminal or the terminal
// does not answer within the deadline.
func queryTerminalForegroundColor() (string, error) {
	tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	if err != nil {
		return "", err
	}
	defer func() { _ = tty.Close() }()

	// In canonical mode read() waits for a newline, but the OSC reply is
	// terminated by BEL/ST. Switch to raw mode so the reply is readable as it
	// arrives, then restore the previous state before bubbletea takes over.
	state, err := term.MakeRaw(tty.Fd())
	if err != nil {
		return "", err
	}
	defer func() { _ = term.Restore(tty.Fd(), state) }()

	if _, err = tty.WriteString(ansi.RequestForegroundColor); err != nil {
		return "", err
	}

	if err = tty.SetReadDeadline(time.Now().Add(150 * time.Millisecond)); err != nil {
		return "", err
	}

	buf := make([]byte, 0, 64)
	tmp := make([]byte, 64)
	for {
		var n int
		n, err = tty.Read(tmp)
		if n > 0 {
			buf = append(buf, tmp[:n]...)
			if bytes.Contains(buf, []byte{ansi.BEL}) || bytes.Contains(buf, []byte{ansi.ESC, '\\'}) {
				break
			}
		}
		if err != nil {
			return "", err
		}
	}

	hex, err := parseOSCColor(buf)
	if err != nil {
		return "", err
	}
	return hex, nil
}

// parseOSCColor extracts the color payload from an OSC 10 response such as
// "\x1b]10;rgb:dddd/eeee/ffff\x07" or "\x1b]10;#ffffff\x07".
func parseOSCColor(buf []byte) (string, error) {
	_, after, ok := bytes.Cut(buf, []byte("10;"))
	if !ok {
		return "", fmt.Errorf("no OSC 10 payload in %q", string(buf))
	}
	payload := after

	if i := bytes.IndexByte(payload, ansi.BEL); i >= 0 {
		payload = payload[:i]
	}
	if i := bytes.Index(payload, []byte{ansi.ESC, '\\'}); i >= 0 {
		payload = payload[:i]
	}

	c := ansi.XParseColor(string(payload))
	if c == nil {
		return "", fmt.Errorf("unparseable color %q", string(payload))
	}
	return colorHex(c), nil
}

// colorHex renders a [color.Color] as a "#rrggbb" string.
func colorHex(c color.Color) string {
	if c == nil {
		return ""
	}
	r, g, b, _ := c.RGBA()
	return fmt.Sprintf("#%02x%02x%02x", r>>8, g>>8, b>>8)
}
