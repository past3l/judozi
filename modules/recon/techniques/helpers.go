package techniques

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
)

// cmd runs a shell command and returns trimmed output
func cmd(c string) string {
	out, err := exec.Command("sh", "-c", c).Output()
	if err != nil || len(bytes.TrimSpace(out)) == 0 {
		return ""
	}
	return strings.TrimRight(string(out), "\n") + "\n"
}

// cmdFile reads a file and returns its contents (with header)
func cmdFile(path string) string {
	data, err := os.ReadFile(path)
	if err != nil || len(bytes.TrimSpace(data)) == 0 {
		return ""
	}
	return "  [" + path + "]\n" + string(data) + "\n"
}
