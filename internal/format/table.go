package format

import (
	"fmt"
	"strings"
)

func Table(headers []string, rows [][]string) string {
	if len(rows) == 0 {
		return "No data to display."
	}

	// Calculate column widths
	widths := make([]int, len(headers))
	for i, h := range headers {
		widths[i] = len(h)
	}
	for _, row := range rows {
		for i, cell := range row {
			if i < len(widths) && len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
	}

	var sb strings.Builder

	// Header
	for i, h := range headers {
		if i > 0 {
			sb.WriteString("  ")
		}
		sb.WriteString(fmt.Sprintf("%-*s", widths[i], h))
	}
	sb.WriteString("\n")

	// Separator
	for i, w := range widths {
		if i > 0 {
			sb.WriteString("  ")
		}
		sb.WriteString(strings.Repeat("─", w))
	}
	sb.WriteString("\n")

	// Rows
	for _, row := range rows {
		for i, cell := range row {
			if i >= len(widths) {
				break
			}
			if i > 0 {
				sb.WriteString("  ")
			}
			sb.WriteString(fmt.Sprintf("%-*s", widths[i], cell))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

func Gold(amount int) string {
	if amount >= 0 {
		return fmt.Sprintf("%d gold", amount)
	}
	return fmt.Sprintf("%d gold", amount)
}

func Trophies(amount int) string {
	return fmt.Sprintf("%d 🏆", amount)
}
