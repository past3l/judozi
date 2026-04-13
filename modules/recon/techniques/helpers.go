package techniques

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"strings"
	"time"
)

// cmd runs a shell command with a 10-second timeout
func cmd(c string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, "sh", "-c", c).Output()
	if err != nil || len(bytes.TrimSpace(out)) == 0 {
		return ""
	}
	return strings.TrimRight(string(out), "\n") + "\n"
}

// cmdFile reads a file and returns its contents with a header label
func cmdFile(path string) string {
	data, err := os.ReadFile(path)
	if err != nil || len(bytes.TrimSpace(data)) == 0 {
		return ""
	}
	return "  [" + path + "]\n" + string(data) + "\n"
}
