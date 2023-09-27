package goload

import (
	"fmt"
	"io"
	"math"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type UI struct {
	output io.Writer
}

func NewUI(output io.Writer) *UI {
	return &UI{
		output: output,
	}
}

func (ui *UI) PrintStartMessage() {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("#48d597"))

	fmt.Fprintf(
		ui.output,
		"Starting a %s Test\n",
		style.Render("GoLoad"),
	)
	fmt.Fprintln(ui.output)
}

func (ui *UI) PrintAbortMessage() {
	fmt.Fprintln(ui.output)
	fmt.Fprintln(ui.output, "Interrupt received.")
	fmt.Fprintln(ui.output, "Waiting for existing requests to finish...")
	fmt.Fprintln(ui.output)
}

func (ui *UI) ReportInitialRPM(rpm int32) {
	fmt.Fprintf(ui.output, "Starting load test with initial %d RPM\n", rpm)
}

func (ui *UI) ReportIncreaseInRPM(rpm int32) {
	fmt.Fprintf(ui.output, "Increasing RPM to %d\n", rpm)
}

func (ui *UI) ReportDecreaseInRPM(rpm int32) {
	fmt.Fprintf(ui.output, "Decreasing RPM to %d\n", rpm)
}

func (ui *UI) ReportResults(results *LoadTestResults) {
	fmt.Fprintln(ui.output)

	rows := generateResultRows(results)

	addTotalRequestsColumn(rows, results)
	addFailedRequestsColumn(rows, results)
	addAverageResponseTimeColumn(rows, results)

	for _, row := range rows {
		fmt.Fprintln(
			ui.output,
			strings.Join(row, " "),
		)
	}
}

func generateResultRows(results *LoadTestResults) [][]string {
	maxLen := 0
	for _, endpoint := range results.Iter() {
		if maxLen < len(endpoint.Name) {
			maxLen = len(endpoint.Name)
		}
	}

	columns := [][]string{}
	for _, endpoint := range results.Iter() {
		columns = append(
			columns,
			[]string{
				fmt.Sprintf(
					"%s%s",
					endpoint.Name,
					lipgloss.
						NewStyle().
						Foreground(lipgloss.Color("#8C8FA3")).
						Render(
							fmt.Sprintf(
								"%s:",
								strings.Repeat(".", maxLen-len(endpoint.Name)+3),
							),
						),
				),
			},
		)
	}

	return columns
}

func addTotalRequestsColumn(columns [][]string, results *LoadTestResults) {
	items := []string{}
	for _, endpoint := range results.Iter() {
		items = append(
			items,
			fmt.Sprintf("total=%d", endpoint.GetTotalRequests()),
		)
	}

	maxLen := 0
	for _, item := range items {
		if maxLen < len(item) {
			maxLen = len(item)
		}
	}

	for index, item := range items {
		columns[index] = append(
			columns[index],
			fmt.Sprintf("%s%s", item, strings.Repeat(" ", maxLen-len(item))),
		)
	}
}

func addFailedRequestsColumn(columns [][]string, results *LoadTestResults) {
	items := []string{}
	for _, endpoint := range results.Iter() {
		items = append(
			items,
			fmt.Sprintf("failed=%d", endpoint.GetTotalFailedRequests()),
		)
	}

	maxLen := 0
	for _, item := range items {
		if maxLen < len(item) {
			maxLen = len(item)
		}
	}

	for index, item := range items {
		columns[index] = append(
			columns[index],
			fmt.Sprintf("%s%s", item, strings.Repeat(" ", maxLen-len(item))),
		)
	}
}

func addAverageResponseTimeColumn(columns [][]string, results *LoadTestResults) {
	items := []string{}
	for _, endpoint := range results.Iter() {
		averageDuration := endpoint.GetAverageDuration()
		if math.IsNaN(averageDuration) {
			items = append(items, "")
			continue
		}

		items = append(
			items,
			fmt.Sprintf("avg=%.2fms", endpoint.GetAverageDuration()),
		)
	}

	maxLen := 0
	for _, item := range items {
		if maxLen < len(item) {
			maxLen = len(item)
		}
	}

	if maxLen == 0 {
		return
	}

	for index, item := range items {
		columns[index] = append(
			columns[index],
			fmt.Sprintf("%s%s", item, strings.Repeat(" ", maxLen-len(item))),
		)
	}
}
