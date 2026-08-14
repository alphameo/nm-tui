package compositor

import "testing"

func TestResolvePos(t *testing.T) {
	tests := []struct {
		name               string
		fgW, fgH, bgW, bgH int
		XAnch, YAnch       Anchor
		xOffset, yOffset   int
		wantX, wantY       int
	}{
		{"begin-begin", 3, 3, 10, 10, Begin, Begin, 0, 0, 0, 0},
		{"center-center", 2, 2, 10, 10, Center, Center, 0, 0, 4, 4},
		{"end-end", 3, 3, 10, 10, End, End, 0, 0, 7, 7},
		{"begin-center", 3, 3, 10, 10, Begin, Center, 0, 0, 0, 3},
		{"center-begin", 3, 3, 10, 10, Center, Begin, 0, 0, 3, 0},
		{"center-with-offset", 10, 10, 20, 20, Center, Center, 2, -3, 7, 2},
		{"odd-size-center-floor", 9, 9, 20, 20, Center, Center, 0, 0, 5, 5},
		{"negative-offsets", 3, 3, 10, 10, Begin, Begin, -5, -5, -5, -5},
		{"larger-than-bg-covers", 20, 20, 10, 10, Begin, Begin, 0, 0, 0, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x, y := resolvePos(tt.fgW, tt.fgH, tt.bgW, tt.bgH, tt.XAnch, tt.YAnch, tt.xOffset, tt.yOffset)
			if x != tt.wantX || y != tt.wantY {
				t.Errorf("resolvePos(%d,%d,%d,%d,%d,%d,%d,%d) = (%d,%d), want (%d,%d)",
					tt.fgW, tt.fgH, tt.bgW, tt.bgH, tt.XAnch, tt.YAnch, tt.xOffset, tt.yOffset,
					x, y, tt.wantX, tt.wantY)
			}
		})
	}
}
