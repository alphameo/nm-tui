package tabview_test

import (
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
	"github.com/alphameo/nm-tui/internal/ui/models/tabview"
)

func testTabStyles() (lipgloss.Style, lipgloss.Style) {
	active := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		Foreground(lipgloss.Color("11")).
		Background(lipgloss.Color("60"))
	inactive := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		Foreground(lipgloss.Color("7")).
		Background(lipgloss.Color("8"))
	return active, inactive
}

func TestRenderTabBar(t *testing.T) {
	active, inactive := testTabStyles()

	tests := []struct {
		name      string
		titles    []string
		fullWidth int
		active    int
	}{
		{"single-tab", []string{"wifi"}, 10, 0},
		{"three-tabs-even", []string{"a", "bb", "ccc"}, 30, 1},
		{"three-tabs-remainder", []string{"a", "bb", "ccc"}, 32, 2},
		{"active-out-of-range", []string{"a", "b"}, 20, 5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tabview.RenderTabBar(tt.titles, active, inactive, tt.fullWidth, tt.active)
			if got == "" {
				t.Fatalf("RenderTabBar() returned empty string")
			}
			for _, title := range tt.titles {
				if !strings.Contains(got, title) {
					t.Errorf("RenderTabBar() output missing title %q: %q", title, got)
				}
			}
		})
	}
}

func TestRenderTabBarEdgeBorders(t *testing.T) {
	active, inactive := testTabStyles()

	gotActiveFirst := tabview.RenderTabBar([]string{"one", "two"}, active, inactive, 20, 0)
	gotActiveLast := tabview.RenderTabBar([]string{"one", "two"}, active, inactive, 20, 1)

	checks := []struct {
		name   string
		output string
		want   string
	}{
		{"active-first", gotActiveFirst, "┌"},
		{"active-last", gotActiveLast, "┐"},
	}
	for _, c := range checks {
		t.Run(c.name, func(t *testing.T) {
			if !strings.Contains(c.output, c.want) {
				t.Errorf("%s output %q missing %q", c.name, c.output, c.want)
			}
		})
	}

	if gotActiveFirst == gotActiveLast {
		t.Errorf("active-first and active-last renderings are identical: %q", gotActiveFirst)
	}
}
