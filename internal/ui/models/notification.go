package models

import (
	"time"

	"charm.land/lipgloss/v2"
)

type Notification struct {
	message   string
	active    bool
	title     string
	closeTime time.Duration
	style     *lipgloss.Style
}
