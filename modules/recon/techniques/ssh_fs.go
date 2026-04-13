package techniques

import (
	"fmt"
	"os"
	"strings"
)

type SSHEnum struct{}

func NewSSHEnum() Technique       { return &SSHEnum{} }
func (s *SSHEnum) Name() string         { return "SSH & Lateral Movement" }
func (s *SSHEnum) Description() string  { return "SSH keys, known_hosts, authorized_keys, agent sockets" }
func (s *SSHEnum) StealthLevel() string { return "passive" }

func (s *SSHEnum) Run() string {
	var b strings.Builder

	b.WriteString("[+] SSH KEYS\n")
	b.WriteString(cmd("find / \\( -name 'id_rsa' -o -name 'id_ecdsa' -o -name 'id_ed25519' -o -name 'id_dsa' \\) 2>/dev/null | head -20"))
	b.WriteString(cmd("find / \\( -name 'id_rsa.pub' -o -name 'id_ecdsa.pub' -o -name 'id_ed25519.pub' \\) 2>/dev/null | head -20"))
	b.WriteString("\n")

	b.WriteString("[+] AUTHORIZED_KEYS\n")
	b.WriteString(cmd("find / -name 'authorized_keys' 2>/dev/null | xargs cat 2>/dev/null"))
	b.WriteString("\n")

	b.WriteString("[+] KNOWN_HOSTS\n")
	b.WriteString(cmd("find / -name 'known_hosts' 2>/dev/null | xargs cat 2>/dev/null"))
	b.WriteString("\n")

	b.WriteString("[+] SSH AGENT SOCKETS\n")
	b.WriteString(cmd("find /tmp -name 'agent.*' -o -name 'ssh-*' 2>/dev/null | head -10"))
	v := os.Getenv("SSH_AUTH_SOCK")
	if v != "" {
		b.WriteString(fmt.Sprintf("  *** SSH_AUTH_SOCK=%s (agent hijack possible if readable)\n", v))
	}
	b.WriteString("\n")

	b.WriteString("[+] SSH CONFIG FILES\n")
	b.WriteString(cmdFile("/etc/ssh/sshd_config"))
	b.WriteString(cmdFile("/etc/ssh/ssh_config"))
	b.WriteString(cmd("find /home -name 'config' -path '*/.ssh/*' 2>/dev/null | xargs cat 2>/dev/null"))
	cmdFile("/root/.ssh/config")
	b.WriteString("\n")

	b.WriteString("[+] HOSTS FROM BASH HISTORY\n")
	b.WriteString(cmd("grep -E '^ssh ' /root/.bash_history /home/*/.bash_history 2>/dev/null | sort -u | head -30"))
	b.WriteString("\n")

	b.WriteString("[+] NETRC FILES (cleartext credentials)\n")
	b.WriteString(cmd("find / -name '.netrc' -readable 2>/dev/null | xargs cat 2>/dev/null"))
	b.WriteString("\n")

	b.WriteString("[+] /etc/hosts.equiv\n")
	b.WriteString(cmdFile("/etc/hosts.equiv"))
	b.WriteString(cmd("cat /root/.rhosts 2>/dev/null"))
	b.WriteString(cmd("find /home -name '.rhosts' 2>/dev/null | xargs cat 2>/dev/null"))
	b.WriteString("\n")

	return b.String()
}

type FSWritable struct{}

func NewFSWritable() Technique      { return &FSWritable{} }
func (f *FSWritable) Name() string         { return "Filesystem & Sensitive Files" }
func (f *FSWritable) Description() string  { return "World-writable dirs, sensitive files, interesting paths" }
func (f *FSWritable) StealthLevel() string { return "low" }

func (f *FSWritable) Run() string {
	var b strings.Builder

	b.WriteString("[+] WORLD-WRITABLE DIRECTORIES\n")
	b.WriteString(cmd("find / -not -path '/proc/*' -not -path '/sys/*' -not -path '/dev/*' -not -path '/run/*' -perm -0002 -type d 2>/dev/null | head -30"))
	b.WriteString("\n")

	b.WriteString("[+] WORLD-WRITABLE FILES (non-symlink)\n")
	b.WriteString(cmd("find / -not -path '/proc/*' -not -path '/sys/*' -not -path '/dev/*' -not -path '/run/*' -perm -0002 -type f 2>/dev/null | head -30"))
	b.WriteString("\n")

	b.WriteString("[+] SENSITIVE FILES CHECK\n")
	sensitiveFiles := []string{
		"/etc/shadow", "/etc/shadow-", "/etc/sudoers",
		"/etc/gshadow", "/etc/gshadow-",
		"/root/.bash_history", "/root/.ssh/id_rsa",
		"/var/lib/mysql/mysql/user.frm",
		"/var/lib/postgresql/",
		"/etc/ssl/private/",
	}
	for _, sf := range sensitiveFiles {
		fi, err := os.Stat(sf)
		if err == nil {
			readable := ""
			f2, err2 := os.Open(sf)
			if err2 == nil {
				f2.Close()
				readable = " *** READABLE ***"
			}
			b.WriteString(fmt.Sprintf("  %-50s [%s]%s\n", sf, fi.Mode(), readable))
		}
	}
	b.WriteString("\n")

	b.WriteString("[+] BACKUP FILES\n")
	b.WriteString(cmd("find /home /root /var /opt /etc /srv -maxdepth 8 \\( -name '*.bak' -o -name '*.backup' -o -name '*.old' -o -name '*.orig' -o -name '*.save' \\) -readable 2>/dev/null | head -20"))
	b.WriteString("\n")

	b.WriteString("[+] INTERESTING FILES IN /var\n")
	b.WriteString(cmd("find /var -readable -type f 2>/dev/null | grep -Ev '^/var/(lib|cache|log)/(dpkg|apt)' | head -40"))
	b.WriteString("\n")

	b.WriteString("[+] DATABASES\n")
	b.WriteString(cmd("find / \\( -name '*.db' -o -name '*.sqlite' -o -name '*.sqlite3' \\) -readable 2>/dev/null | grep -v '/sys/' | head -20"))
	b.WriteString("\n")

	b.WriteString("[+] LOG FILES (interesting)\n")
	b.WriteString(cmd("ls -la /var/log/ 2>/dev/null"))
	b.WriteString(cmd("grep -Ei 'password|failed|error|root|sudo' /var/log/auth.log 2>/dev/null | tail -30"))
	b.WriteString(cmd("grep -Ei 'password|failed|error|root|sudo' /var/log/secure 2>/dev/null | tail -30"))
	b.WriteString("\n")

	b.WriteString("[+] /opt, /srv, /app CONTENTS\n")
	b.WriteString(cmd("find /opt /srv /app /data /mnt -maxdepth 4 -type f -readable 2>/dev/null | head -40"))
	b.WriteString("\n")

	return b.String()
}
