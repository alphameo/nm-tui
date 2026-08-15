package styles

import (
	"testing"
)

func TestParseOSCColor(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []byte
		want string
		err  bool
	}{
		{name: "rgb 4-digit bel", in: []byte("\x1b]10;rgb:dddd/eeee/ffff\x07"), want: "#ddeeff"},
		{name: "rgb 2-digit bel", in: []byte("\x1b]10;rgb:ff/aa/00\x07"), want: "#ffaa00"},
		{name: "rgba", in: []byte("\x1b]10;rgba:1122/3344/5566/7788\x07"), want: "#113355"},
		{name: "hex bel", in: []byte("\x1b]10;#865fff\x07"), want: "#865fff"},
		{name: "hex st", in: []byte("\x1b]10;#ffffff\x1b\\"), want: "#ffffff"},
		{name: "extra trailing bytes", in: []byte("\x1b]10;rgb:0102/0304/0506\x07extra"), want: "#010305"},
		{name: "no payload", in: []byte("\x1b]11;#ffffff\x07"), err: true},
		{name: "color name unsupported", in: []byte("\x1b]10;default\x07"), err: true},
		{name: "empty", in: nil, err: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := parseOSCColor(tt.in)
			if tt.err {
				if err == nil {
					t.Errorf("parseOSCColor(%q) = %q, want error", tt.in, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("parseOSCColor(%q) error: %v", tt.in, err)
			}
			if got != tt.want {
				t.Errorf("parseOSCColor(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}
