package techniques

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type SystemInfo struct{}

func NewSystemInfo() Technique { return &SystemInfo{} }
func (s *SystemInfo) Name() string         { return "System Information" }
func (s *SystemInfo) Description() string  { return "OS, kernel, hostname, hardware, CPU, memory" }
func (s *SystemInfo) StealthLevel() string { return "passive" }

func (s *SystemInfo) Run() string {
	var b strings.Builder

	b.WriteString("[+] HOSTNAME & OS\n")
	b.WriteString(cmd("uname -a"))
	b.WriteString(cmdFile("/etc/os-release"))
	b.WriteString(cmdFile("/etc/issue"))
	b.WriteString("\n")

	b.WriteString("[+] CPU / ARCHITECTURE\n")
	b.WriteString(cmd("lscpu 2>/dev/null | head -30"))
	b.WriteString(cmd("cat /proc/cpuinfo 2>/dev/null | grep 'model name' | uniq"))
	b.WriteString("\n")

	b.WriteString("[+] MEMORY\n")
	b.WriteString(cmd("free -h"))
	b.WriteString("\n")

	b.WriteString("[+] DISK\n")
	b.WriteString(cmd("df -h 2>/dev/null"))
	b.WriteString("\n")

	b.WriteString("[+] UPTIME & BOOT TIME\n")
	b.WriteString(cmd("uptime"))
	b.WriteString(cmd("who -b 2>/dev/null"))
	b.WriteString(cmd("last reboot 2>/dev/null | head -5"))
	b.WriteString("\n")

	b.WriteString("[+] ENVIRONMENT\n")
	b.WriteString(cmd("env 2>/dev/null | sort"))
	b.WriteString("\n")

	b.WriteString("[+] PATH ANALYSIS\n")
	pathVal := os.Getenv("PATH")
	for _, p := range strings.Split(pathVal, ":") {
		info := fmt.Sprintf("  %-40s", p)
		fi, err := os.Stat(p)
		if err != nil {
			info += " [NOT FOUND]"
		} else {
			info += fmt.Sprintf(" [%s]", fi.Mode())
			// check writable
			if fi.Mode().Perm()&0002 != 0 {
				info += " *** WORLD WRITABLE ***"
			}
		}
		b.WriteString(info + "\n")
	}
	b.WriteString("\n")

	b.WriteString("[+] KERNEL MODULES\n")
	b.WriteString(cmd("lsmod 2>/dev/null | head -30"))
	b.WriteString("\n")

	b.WriteString("[+] DMESG (last 20 lines)\n")
	b.WriteString(cmd("dmesg 2>/dev/null | tail -20"))
	b.WriteString("\n")

	b.WriteString("[+] INSTALLED TOOLS\n")
	tools := []string{
		"gcc", "cc", "g++", "make", "python", "python3", "perl", "ruby", "php",
		"wget", "curl", "nc", "ncat", "nmap", "tcpdump", "strace", "ltrace",
		"gdb", "socat", "sqlmap", "docker", "kubectl", "nfs", "mount",
	}
	for _, tool := range tools {
		path, err := exec.LookPath(tool)
		if err == nil {
			b.WriteString(fmt.Sprintf("  %-15s %s\n", tool, path))
		}
	}

	return b.String()
}
