package floating

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

type Anchor int

const (
	Begin Anchor = iota
	Center
	End
)

func resolvePos(fgW, fgH, bgW, bgH int, XAnch, YAnch Anchor, xOffset, yOffset int) (int, int) {
	var xPos, yPos int
	switch XAnch {
	case Begin:
		xPos = 0
	case Center:
		xPos = (bgW - fgW) / 2
	case End:
		xPos = bgW - fgW
	}
	switch YAnch {
	case Begin:
		yPos = 0
	case Center:
		yPos = (bgH - fgH) / 2
	case End:
		yPos = bgH - fgH
	}

	return xPos + xOffset, yPos + yOffset
}

func Compose(fg, bg string, xAnchor, yAnchor Anchor, xOffset, yOffset int) string {
	fgW, fgH := lipgloss.Size(fg)
	bgW, bgH := lipgloss.Size(bg)
	fgXmin, fgYmin := resolvePos(fgW, fgH, bgW, bgH, xAnchor, yAnchor, xOffset, yOffset)
	fgXmax := fgXmin + fgW
	fgYmax := fgYmin + fgH

	if (fgW >= bgW && fgH >= bgH) || fgXmin >= bgW || fgYmin >= bgH || fgXmax < 0 || fgYmax < 0 {
		return bg
	}

	fgLines := lines(fg)
	bgLines := lines(bg)

	var sb strings.Builder

	var fgInd int
	if fgYmin < 0 {
		fgInd -= fgYmin
	}

	for bgY, bgLine := range bgLines {
		if bgY > 0 {
			sb.WriteByte('\n')
		}
		if bgY < fgYmin || bgY >= fgYmax {
			sb.WriteString(bgLine)
			continue
		}

		if fgXmin > 0 {
			left := ansi.Truncate(bgLine, fgXmin, "")
			sb.WriteString(left)
		}

		fgLine := fgLines[fgInd]
		fgInd++
		if fgXmin < 0 {
			fgLine = ansi.TruncateLeft(fgLine, -fgXmin, "")
		}

		if fgXmax <= bgW {
			sb.WriteString(fgLine)
		} else {
			sb.WriteString(ansi.Truncate(fgLine, fgW-fgXmax+bgW, ""))
			continue
		}

		right := ansi.TruncateLeft(bgLine, fgXmax, "")
		sb.WriteString(right)
	}
	return sb.String()
}

func lines(s string) []string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	return strings.Split(s, "\n")
}
