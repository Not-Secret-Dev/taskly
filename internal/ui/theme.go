// Package ui renders the terminal interface: styling primitives, prompts,
// the menu, the stats bar, and the task table.
package ui

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

const (
	reset   = "\033[0m"
	bold    = "\033[1m"
	dim     = "\033[2m"
	strike  = "\033[9m"
	red     = "\033[31m"
	green   = "\033[32m"
	yellow  = "\033[33m"
	blue    = "\033[34m"
	magenta = "\033[35m"
	cyan    = "\033[36m"
	gray    = "\033[90m"
)

var ansiRe = regexp.MustCompile(`\x1b\[[0-9;]*m`)

// visibleLen counts printable runes, ignoring ANSI escape codes.
func visibleLen(s string) int {
	return utf8.RuneCountInString(ansiRe.ReplaceAllString(s, ""))
}

// pad right-pads s to w visible characters (ANSI-safe).
func pad(s string, w int) string {
	n := w - visibleLen(s)
	if n < 0 {
		n = 0
	}
	return s + strings.Repeat(" ", n)
}

// center centers s within w visible characters (ANSI-safe).
func center(s string, w int) string {
	n := w - visibleLen(s)
	if n <= 0 {
		return s
	}
	left := n / 2
	return strings.Repeat(" ", left) + s + strings.Repeat(" ", n-left)
}

// truncate shortens s to w runes, ending with an ellipsis when cut.
func truncate(s string, w int) string {
	r := []rune(s)
	if len(r) <= w {
		return s
	}
	return string(r[:w-1]) + "…"
}
