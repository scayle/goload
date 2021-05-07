package goload

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

type UI struct{}

func NewUI() *UI {
	return &UI{}
}

func (ui *UI) PrintStartMessage() {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("#48d597"))

	fmt.Printf(
		"Starting a %s Load Test\n",
		style.Render("GoLoad"),
	)
	fmt.Println()
}

func (ui *UI) PrintAbortMessage() {
	fmt.Println()
	fmt.Println("Interrupt received.")
	fmt.Println("Please wait for GoLoad to exit or data loss may occur.")
	fmt.Println("Waiting for existing requests to finish...")
	fmt.Println()
}

func (ui *UI) ReportInitialRPM(rpm int32) {
	fmt.Printf("Starting load test with initial %d RPM\n", rpm)
}

func (ui *UI) ReportIncreaseInRPM(rpm int32) {
	fmt.Printf("Increasing RPM to %d\n", rpm)
}

func (ui *UI) ReportDecreaseInRPM(rpm int32) {
	fmt.Printf("Decreasing RPM to %d\n", rpm)
}
