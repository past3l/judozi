package techniques

import "strings"

// Helper functions used by techniques

func splitLines(s string) []string {
	return strings.Split(s, "\n")
}

func joinLines(lines []string) string {
	return strings.Join(lines, "\n")
}

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
