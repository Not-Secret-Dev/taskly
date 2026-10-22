package ui

import (
	"bufio"
	"fmt"
	"strings"
)

// Clear clears the terminal.
func Clear() { fmt.Print("\033[H\033[2J") }

// Success prints a success message.
func Success(msg string) { fmt.Printf("\n  %s✔ %s%s\n", green, msg, reset) }

// Failure prints an error message.
func Failure(msg string) { fmt.Printf("\n  %s✖ %s%s\n", red, msg, reset) }

// Info prints an informational message.
func Info(msg string) { fmt.Printf("\n  %sℹ %s%s\n", cyan, msg, reset) }

// FileCreated reports that a new data file was written.
func FileCreated(name string) {
	fmt.Printf("%sCreated %s and wrote headers.%s\n", dim, name, reset)
}

// Pause waits for the user to press Enter.
func Pause(scanner *bufio.Scanner) {
	fmt.Printf("\n  %sPress Enter to continue…%s", gray, reset)
	scanner.Scan()
}

// Prompt prints label and returns the trimmed line the user typed.
func Prompt(scanner *bufio.Scanner, label string) string {
	fmt.Printf("  %s▸%s %s", magenta, reset, label)
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}

// Goodbye prints the farewell message.
func Goodbye() {
	fmt.Printf("\n  %s✦ Goodbye! Stay productive. ✦%s\n\n", cyan, reset)
}
