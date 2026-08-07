package compositor

import (
	"strings"
	"testing"
)

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

func TestCompose(t *testing.T) {
	tests := []struct {
		name             string
		fg, bg           string
		xAnchor, yAnchor Anchor
		xOffset, yOffset int
		want             string
	}{
		{
			"fg-covers-bg-returns-bg",
			"AAAA", "BBBB",
			Begin, Begin, 0, 0,
			"BBBB",
		},
		{
			"fg-offscreen-right-leaves-bg",
			"AA", "BBB\nBBB\nBBB",
			Begin, Begin, 10, 0,
			"BBB\nBBB\nBBB",
		},
		{
			"fg-offscreen-below-leaves-bg",
			"AA", "BBB\nBBB\nBBB",
			Begin, Begin, 0, 10,
			"BBB\nBBB\nBBB",
		},
		{
			"overwrite-corner",
			"XY", "BBC\nBBC\nBBB",
			Begin, Begin, 0, 0,
			"XYC\nBBC\nBBB",
		},
		{
			"offset-right-preserves-left-column",
			"XY", "BBC\nBBC\nBBC",
			Begin, Begin, 1, 0,
			"BXY\nBBC\nBBC",
		},
		{
			"center-placement",
			"AB", "####\n####\n####\n####",
			Center, Center, 0, 0,
			"####\n#AB#\n####\n####",
		},
		{
			"multi-row-fg",
			"A\nB", "CC\nCC",
			Begin, Begin, 0, 0,
			"AC\nBC",
		},
		{
			"fg-smaller-than-bg-stays-intact",
			"XY", "abcdef",
			Begin, Begin, 0, 0,
			"XYcdef",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Compose(tt.fg, tt.bg, tt.xAnchor, tt.yAnchor, tt.xOffset, tt.yOffset)
			if got != tt.want {
				t.Errorf("Compose() =\n%q\nwant\n%q", got, tt.want)
			}
		})
	}
}

func TestComposeTrimNewlines(t *testing.T) {
	got := Compose("XY", "A", Begin, Begin, 0, 0)
	if strings.Contains(got, "\n") {
		t.Errorf("Compose() introduced trailing newline: %q", got)
	}
}
