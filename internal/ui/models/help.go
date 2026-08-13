package models

import (
	"strings"
)

func HelpFromKeys(keys ...string) string {
	transformed := make([]string, len(keys))
	for i, key := range keys {
		transformed[i] = strings.ReplaceAll(key, "ctrl+", "^")
	}
	return strings.Join(transformed, "/")
}
