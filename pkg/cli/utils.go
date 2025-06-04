package cli

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

// PrintSuccess prints a success message
func PrintSuccess(message string) {
	fmt.Printf("✅ %s\n", message)
}

// PrintWarning prints a warning message
func PrintWarning(message string) {
	fmt.Printf("⚠️  %s\n", message)
}

// PrintError prints an error message
func PrintError(message string) {
	fmt.Printf("❌ %s\n", message)
}

// PrintInfo prints an info message
func PrintInfo(message string) {
	fmt.Printf("💡 %s\n", message)
}

// PrintHeader prints a formatted header
func PrintHeader(title string) {
	fmt.Printf("\n%s\n", title)
	fmt.Println(strings.Repeat("=", len(title)))
}

// PrintSubHeader prints a formatted sub-header
func PrintSubHeader(title string) {
	fmt.Printf("\n%s\n", title)
	fmt.Println(strings.Repeat("-", len(title)))
}

// CreateTable creates a new tabwriter for formatted table output
func CreateTable() *tabwriter.Writer {
	return tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
}

// PrintTable prints data in table format
func PrintTable(headers []string, rows [][]string) {
	w := CreateTable()
	defer w.Flush()

	// Print headers
	fmt.Fprintln(w, strings.Join(headers, "\t"))

	// Print rows
	for _, row := range rows {
		fmt.Fprintln(w, strings.Join(row, "\t"))
	}
}

// FormatAge formats duration strings to be more readable
func FormatAge(age string) string {
	return age
}

// TruncateString truncates a string to specified length with ellipsis
func TruncateString(s string, length int) string {
	if len(s) <= length {
		return s
	}
	return s[:length-3] + "..."
}

// FormatStatus formats status strings
func FormatStatus(status string) string {
	return status
}

// PromptUser prompts the user for input and returns the response
func PromptUser(prompt string) string {
	fmt.Print(prompt)
	var input string
	fmt.Scanln(&input)
	return strings.TrimSpace(input)
}

// FormatJSON formats a string as JSON (placeholder - in real implementation would use proper JSON formatting)
func FormatJSON(content string) string {
	// In a real implementation, this would parse YAML and convert to JSON
	// For now, just return the content with a JSON comment
	return fmt.Sprintf("# JSON format not yet implemented\n%s", content)
}

// PrintBanner prints the application banner
func PrintBanner() {
	banner := `
 ╔═══════════════════════════════════════╗
 ║        🤖 Crossplane AI CLI           ║
 ║   Intelligent Infrastructure as Code  ║
 ╚═══════════════════════════════════════╝
`
	fmt.Println(banner)
}
