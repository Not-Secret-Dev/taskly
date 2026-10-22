package ui

import (
	"fmt"
	"strings"

	"todo-cli/internal/domain"
)

const bannerWidth = 38

// Banner prints the application header.
func Banner() {
	w := bannerWidth
	fmt.Println(cyan + "╭" + strings.Repeat("─", w) + "╮" + reset)
	fmt.Println(cyan + "│" + reset + center(bold+"✦  T O D O   M A N A G E R  ✦"+reset, w) + cyan + "│" + reset)
	fmt.Println(cyan + "╰" + strings.Repeat("─", w) + "╯" + reset)
}

// Menu prints the available actions.
func Menu() {
	fmt.Println()
	fmt.Printf("  %s1%s  📋  List tasks\n", bold+cyan, reset)
	fmt.Printf("  %s2%s  ➕  Add task\n", bold+cyan, reset)
	fmt.Printf("  %s3%s  🔁  Toggle complete\n", bold+cyan, reset)
	fmt.Printf("  %s4%s  🗑   Delete task\n", bold+cyan, reset)
	fmt.Printf("  %s5%s  👋  Exit\n", bold+cyan, reset)
	fmt.Println()
}

// ChoosePrompt prints the menu input prompt.
func ChoosePrompt() {
	fmt.Printf("  %sChoose an option (1-5):%s ", dim, reset)
}

// StatsBar prints the progress bar and counters.
func StatsBar(s domain.Stats) {
	fmt.Printf("  %s %s%3d%%%s  %s%d/%d done%s",
		progressBar(s.Done, s.Total, 20), bold, s.Percent(), reset, dim, s.Done, s.Total, reset)
	if s.Overdue > 0 {
		fmt.Printf("  %s⚠ %d overdue%s", red, s.Overdue, reset)
	}
	fmt.Println()
}

func progressBar(done, total, width int) string {
	if total == 0 {
		return gray + strings.Repeat("░", width) + reset
	}
	filled := done * width / total
	return green + strings.Repeat("█", filled) + reset +
		gray + strings.Repeat("░", width-filled) + reset
}
