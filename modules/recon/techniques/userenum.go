package techniques

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type UserEnum struct{}

func NewUserEnum() Technique        { return &UserEnum{} }
func (u *UserEnum) Name() string         { return "User & Group Enumeration" }
func (u *UserEnum) Description() string  { return "Users, groups, sudoers, passwd, shadow status, login history" }
func (u *UserEnum) StealthLevel() string { return "passive" }

func (u *UserEnum) Run() string {
	var b strings.Builder

	b.WriteString("[+] CURRENT USER\n")
	b.WriteString(cmd("id"))
	b.WriteString(cmd("whoami"))
	b.WriteString(cmd("groups"))
	b.WriteString("\n")

	b.WriteString("[+] ALL USERS (with shells)\n")
	parsePasswd(&b)
	b.WriteString("\n")

	b.WriteString("[+] USERS WITH UID >= 1000\n")
	b.WriteString(cmd("awk -F: '$3>=1000{print $1,$3,$4,$6,$7}' /etc/passwd 2>/dev/null"))
	b.WriteString("\n")

	b.WriteString("[+] GROUPS\n")
	b.WriteString(cmd("cat /etc/group 2>/dev/null"))
	b.WriteString("\n")

	b.WriteString("[+] SUDOERS / SUDO ACCESS\n")
	b.WriteString(cmd("sudo -l 2>/dev/null"))
	b.WriteString(cmdFile("/etc/sudoers"))
	b.WriteString(cmd("cat /etc/sudoers.d/* 2>/dev/null"))
	b.WriteString("\n")

	b.WriteString("[+] SHADOW FILE STATUS\n")
	fi, err := os.Stat("/etc/shadow")
	if err == nil {
		b.WriteString(fmt.Sprintf("  /etc/shadow exists: %s\n", fi.Mode()))
		if fi.Mode().Perm()&0004 != 0 {
			b.WriteString("  *** SHADOW FILE IS WORLD-READABLE! ***\n")
			b.WriteString(cmd("cat /etc/shadow 2>/dev/null"))
		}
	} else {
		b.WriteString("  /etc/shadow: not accessible\n")
	}
	b.WriteString("\n")

	b.WriteString("[+] PASSWORD POLICIES\n")
	b.WriteString(cmdFile("/etc/login.defs"))
	b.WriteString(cmdFile("/etc/pam.d/common-password"))
	b.WriteString("\n")

	b.WriteString("[+] LOGGED IN USERS\n")
	b.WriteString(cmd("w 2>/dev/null"))
	b.WriteString(cmd("who 2>/dev/null"))
	b.WriteString("\n")

	b.WriteString("[+] LAST LOGINS\n")
	b.WriteString(cmd("last -n 20 2>/dev/null"))
	b.WriteString("\n")

	b.WriteString("[+] LASTLOG\n")
	b.WriteString(cmd("lastlog 2>/dev/null | grep -v 'Never'"))
	b.WriteString("\n")

	b.WriteString("[+] HOME DIRECTORIES\n")
	b.WriteString(cmd("ls -la /home/ 2>/dev/null"))
	b.WriteString(cmd("ls -la /root/ 2>/dev/null"))
	b.WriteString("\n")

	b.WriteString("[+] USER BASH HISTORIES\n")
	b.WriteString(cmd("cat /root/.bash_history 2>/dev/null"))
	b.WriteString(cmd("cat ~/.bash_history 2>/dev/null"))
	b.WriteString(cmd("find /home -name '.bash_history' -readable 2>/dev/null -exec echo '--- {}' \\; -exec cat {} \\;"))
	b.WriteString(cmd("find /home -name '.zsh_history' -readable 2>/dev/null -exec echo '--- {}' \\; -exec cat {} \\;"))
	b.WriteString("\n")

	b.WriteString("[+] INTERESTING .rc / DOTFILES\n")
	b.WriteString(cmd("find /home -maxdepth 2 -name '.*' -readable 2>/dev/null | head -40"))
	b.WriteString(cmd("ls -la ~/.* 2>/dev/null"))
	b.WriteString("\n")

	return b.String()
}

func parsePasswd(b *strings.Builder) {
	f, err := os.Open("/etc/passwd")
	if err != nil {
		b.WriteString("  [cannot read /etc/passwd]\n")
		return
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, ":")
		if len(fields) < 7 {
			continue
		}
		shell := fields[6]
		// filter for real shells
		if strings.Contains(shell, "sh") || strings.Contains(shell, "bash") || strings.Contains(shell, "fish") || strings.Contains(shell, "zsh") {
			b.WriteString(fmt.Sprintf("  %-20s uid=%-6s gid=%-6s home=%-30s shell=%s\n",
				fields[0], fields[2], fields[3], fields[5], shell))
		}
	}
}
