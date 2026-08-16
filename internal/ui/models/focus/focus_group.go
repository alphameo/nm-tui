// Package focus provides helper structures for management of focusable elements.
package focus

import tea "charm.land/bubbletea/v2"

type Focusable interface {
	Focused() bool
	Focus() tea.Cmd
	Blur()
}

type Group struct {
	focuses  []Focusable
	focusIdx int
}

func NewGroup(focuses []Focusable) *Group {
	return &Group{focuses: focuses, focusIdx: 0}
}

// SetFocusIdx sets focus on last if [idx] >= length, and on first if [idx] < 0.
func (f *Group) SetFocusIdx(idx int) tea.Cmd {
	if len(f.focuses) == 0 {
		return nil
	}
	f.focuses[f.focusIdx].Blur()
	if idx >= len(f.focuses) {
		idx = len(f.focuses) - 1
	} else if idx < 0 {
		idx = 0
	}
	f.focusIdx = idx
	return f.focuses[f.focusIdx].Focus()
}

func (f *Group) FocusPrevCmd() tea.Cmd {
	if f.focusIdx <= 0 {
		return nil
	}
	f.focuses[f.focusIdx].Blur()
	f.focusIdx--
	return f.focuses[f.focusIdx].Focus()
}

func (f *Group) FocusCyclePrevCmd() tea.Cmd {
	f.focuses[f.focusIdx].Blur()
	f.focusIdx = (f.focusIdx + len(f.focuses) - 1) % len(f.focuses)
	return f.focuses[f.focusIdx].Focus()
}

func (f *Group) FocusNextCmd() tea.Cmd {
	if f.focusIdx >= len(f.focuses)-1 {
		return nil
	}
	f.focuses[f.focusIdx].Blur()
	f.focusIdx++
	return f.focuses[f.focusIdx].Focus()
}

func (f *Group) FocusCycleNextCmd() tea.Cmd {
	f.focuses[f.focusIdx].Blur()
	f.focusIdx = (f.focusIdx + len(f.focuses) + 1) % len(f.focuses)
	return f.focuses[f.focusIdx].Focus()
}

func (f *Group) FocusIdx() int {
	return f.focusIdx
}

func (f *Group) FocusCurrent() tea.Cmd {
	return f.focuses[f.focusIdx].Focus()
}
