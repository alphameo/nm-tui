package tabview

import tea "charm.land/bubbletea/v2"

type SizedModel interface {
	tea.Model
	Resize(width, height int)
	Width() int
	Height() int
}
