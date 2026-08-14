package compositor_test

import (
	"strings"
	"testing"

	"github.com/alphameo/nm-tui/internal/ui/tools/compositor"
)

func TestCompose(t *testing.T) {
	tests := []struct {
		name             string
		fg, bg           string
		xAnchor, yAnchor compositor.Anchor
		xOffset, yOffset int
		want             string
	}{
		{
			"fg-covers-bg-returns-bg",
			"AAAA", "BBBB",
			compositor.Begin, compositor.Begin, 0, 0,
			"BBBB",
		},
		{
			"fg-offscreen-right-leaves-bg",
			"AA", "BBB\nBBB\nBBB",
			compositor.Begin, compositor.Begin, 10, 0,
			"BBB\nBBB\nBBB",
		},
		{
			"fg-offscreen-below-leaves-bg",
			"AA", "BBB\nBBB\nBBB",
			compositor.Begin, compositor.Begin, 0, 10,
			"BBB\nBBB\nBBB",
		},
		{
			"overwrite-corner",
			"XY", "BBC\nBBC\nBBB",
			compositor.Begin, compositor.Begin, 0, 0,
			"XYC\nBBC\nBBB",
		},
		{
			"offset-right-preserves-left-column",
			"XY", "BBC\nBBC\nBBC",
			compositor.Begin, compositor.Begin, 1, 0,
			"BXY\nBBC\nBBC",
		},
		{
			"center-placement",
			"AB", "####\n####\n####\n####",
			compositor.Center, compositor.Center, 0, 0,
			"####\n#AB#\n####\n####",
		},
		{
			"multi-row-fg",
			"A\nB", "CC\nCC",
			compositor.Begin, compositor.Begin, 0, 0,
			"AC\nBC",
		},
		{
			"fg-smaller-than-bg-stays-intact",
			"XY", "abcdef",
			compositor.Begin, compositor.Begin, 0, 0,
			"XYcdef",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := compositor.Compose(tt.fg, tt.bg, tt.xAnchor, tt.yAnchor, tt.xOffset, tt.yOffset)
			if got != tt.want {
				t.Errorf("Compose() =\n%q\nwant\n%q", got, tt.want)
			}
		})
	}
}

func TestComposeTrimNewlines(t *testing.T) {
	got := compositor.Compose("XY", "A", compositor.Begin, compositor.Begin, 0, 0)
	if strings.Contains(got, "\n") {
		t.Errorf("Compose() introduced trailing newline: %q", got)
	}
}
