package techniques

import (
	"fmt"
	"os"
	"strings"
	"syscall"
)

type SUIDSGIDEnum struct{}

func NewSUIDSGID() Technique         { return &SUIDSGIDEnum{} }
func (s *SUIDSGIDEnum) Name() string         { return "SUID/SGID Binaries" }
func (s *SUIDSGIDEnum) Description() string  { return "All SUID/SGID binaries, capabilities, writable paths" }
func (s *SUIDSGIDEnum) StealthLevel() string { return "low" }

func (s *SUIDSGIDEnum) Run() string {
	var b strings.Builder

	b.WriteString("[+] SUID BINARIES\n")
	suidOut := cmd("find / -perm -4000 -type f -not -path '/proc/*' -not -path '/sys/*' -not -path '/dev/*' -not -path '/run/*' 2>/dev/null")
	b.WriteString(suidOut)
	b.WriteString("\n")

	b.WriteString("[+] SGID BINARIES\n")
	b.WriteString(cmd("find / -perm -2000 -type f -not -path '/proc/*' -not -path '/sys/*' -not -path '/dev/*' -not -path '/run/*' 2>/dev/null | head -40"))
	b.WriteString("\n")

	b.WriteString("[+] FILE CAPABILITIES\n")
	b.WriteString(cmd("getcap -r / 2>/dev/null"))
	b.WriteString("\n")

	b.WriteString("[+] INTERESTING SUID BINS (GTFOBins candidates)\n")
	gtfoBins := []string{
		"bash", "dash", "sh", "ksh", "csh", "tcsh", "zsh",
		"find", "vim", "vi", "nmap", "awk", "gawk", "nawk", "man",
		"less", "more", "perl", "python", "python3", "ruby", "lua",
		"php", "cc", "gcc", "g++", "tar", "zip", "unzip", "ar",
		"cpio", "rsync", "cp", "mv", "scp", "dd", "tee", "cat",
		"head", "tail", "cut", "base64", "base32",
		"openssl", "curl", "wget", "git", "env", "sudo", "doas",
		"pkexec", "newgrp", "passwd", "chsh", "chfn",
		"strace", "ltrace", "gdb", "screen", "tmux",
		"docker", "podman", "lxc", "kubectl",
		"mount", "umount", "nsenter", "unshare", "ionice", "nice",
		"chroot", "systemctl", "journalctl", "service",
	}
	for _, line := range strings.Split(suidOut, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		for _, g := range gtfoBins {
			if strings.HasSuffix(line, "/"+g) || line == g {
				b.WriteString(fmt.Sprintf("  *** GTFOBIN: %s\n", line))
				break
			}
		}
	}
	b.WriteString("\n")

	b.WriteString("[+] WRITABLE DIRECTORIES IN PATH\n")
	pathVal := os.Getenv("PATH")
	for _, p := range strings.Split(pathVal, ":") {
		fi, err := os.Stat(p)
		if err != nil {
			continue
		}
		if fi.Mode().Perm()&0002 != 0 {
			b.WriteString(fmt.Sprintf("  *** WORLD-WRITABLE PATH: %s\n", p))
		} else if isWritable(p) {
			b.WriteString(fmt.Sprintf("  *** WRITABLE BY CURRENT USER: %s\n", p))
		}
	}
	b.WriteString("\n")

	b.WriteString("[+] SUDO NOPASSWD\n")
	b.WriteString(cmd("sudo -l 2>/dev/null | grep -i 'nopasswd'"))
	b.WriteString("\n")

	b.WriteString("[+] POLKIT / PKEXEC VERSION\n")
	b.WriteString(cmd("pkexec --version 2>/dev/null || dpkg -l policykit-1 2>/dev/null | tail -1 || rpm -q polkit 2>/dev/null"))
	b.WriteString("\n")

	b.WriteString("[+] WRITABLE SENSITIVE FILES\n")
	for _, sf := range []string{"/etc/passwd", "/etc/shadow", "/etc/sudoers"} {
		if isWritable(sf) {
			b.WriteString(fmt.Sprintf("  *** %s IS WRITABLE!\n", sf))
		} else {
			b.WriteString(fmt.Sprintf("  %s — not writable\n", sf))
		}
	}
	b.WriteString("\n")

	return b.String()
}

func isWritable(path string) bool {
	return syscall.Access(path, syscall.O_WRONLY) == nil
}

func (s *SUIDSGIDEnum) Run() string {
	var b strings.Builder

	b.WriteString("[+] SUID BINARIES\n")
	suidBins := findSUID()
	for _, bin := range suidBins {
		b.WriteString(bin)
	}
	if len(suidBins) == 0 {
		b.WriteString("  [none found]\n")
	}
	b.WriteString("\n")

	b.WriteString("[+] SGID BINARIES\n")
	b.WriteString(cmd("find / -perm -2000 -type f 2>/dev/null | head -60"))
	b.WriteString("\n")

	b.WriteString("[+] FILE CAPABILITIES\n")
	b.WriteString(cmd("getcap -r / 2>/dev/null"))
	b.WriteString("\n")

	b.WriteString("[+] INTERESTING SUID BINS (GTFOBins candidates)\n")
	gtfoBins := []string{
		"bash", "dash", "sh", "ksh", "csh", "tcsh", "zsh",
		"find", "vim", "vi", "nmap", "awk", "gawk", "nawk", "man",
		"less", "more", "perl", "python", "python3", "ruby", "lua",
		"php", "cc", "gcc", "g++", "tar", "zip", "unzip", "ar",
		"cpio", "rsync", "cp", "mv", "scp", "dd", "tee", "cat",
		"head", "tail", "cut", "base64", "base32",
		"openssl", "curl", "wget", "git", "env", "sudo", "doas",
		"pkexec", "newgrp", "passwd", "chsh", "chfn",
		"strace", "ltrace", "gdb", "screen", "tmux",
		"docker", "podman", "lxc", "kubectl",
		"mount", "umount", "nsenter", "unshare", "ionice", "nice",
		"chroot", "systemctl", "journalctl", "service",
	}
	for _, bin := range suidBins {
		for _, g := range gtfoBins {
			if strings.Contains(bin, "/"+g) || strings.HasSuffix(strings.TrimSpace(bin), g) {
				b.WriteString(fmt.Sprintf("  *** GTFOBIN FOUND: %s\n", strings.TrimSpace(bin)))
			}
		}
	}
	b.WriteString("\n")

	b.WriteString("[+] WRITABLE DIRECTORIES IN PATH\n")
	pathVal := os.Getenv("PATH")
	for _, p := range strings.Split(pathVal, ":") {
		fi, err := os.Stat(p)
		if err != nil {
			continue
		}
		if fi.Mode().Perm()&0002 != 0 {
			b.WriteString(fmt.Sprintf("  *** WORLD-WRITABLE PATH: %s\n", p))
		}
		if isWritable(p) {
			b.WriteString(fmt.Sprintf("  *** WRITABLE PATH (current user): %s\n", p))
		}
	}
	b.WriteString("\n")

	b.WriteString("[+] SUDO WITHOUT PASSWORD (NOPASSWD)\n")
	b.WriteString(cmd("sudo -l 2>/dev/null | grep -i 'nopasswd'"))
	b.WriteString("\n")

	b.WriteString("[+] POLKIT / PKEXEC VERSION\n")
	b.WriteString(cmd("pkexec --version 2>/dev/null"))
	b.WriteString(cmd("dpkg -l policykit-1 2>/dev/null | tail -1"))
	b.WriteString(cmd("rpm -q polkit 2>/dev/null"))
	b.WriteString("\n")

	b.WriteString("[+] WRITABLE /etc/passwd\n")
	if isWritable("/etc/passwd") {
		b.WriteString("  *** /etc/passwd IS WRITABLE! (classic privesc vector)\n")
	} else {
		b.WriteString("  /etc/passwd not writable\n")
	}

	b.WriteString("[+] WRITABLE /etc/shadow\n")
	if isWritable("/etc/shadow") {
		b.WriteString("  *** /etc/shadow IS WRITABLE!\n")
	} else {
		b.WriteString("  /etc/shadow not writable\n")
	}
	b.WriteString("\n")

	return b.String()
}

func findSUID() []string {
	var result []string
	searchPaths := []string{"/", "/usr", "/bin", "/sbin", "/usr/bin", "/usr/sbin", "/usr/local"}
	visited := map[string]bool{}
	for _, root := range searchPaths {
		_ = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return nil
			}
			if visited[path] {
				return nil
			}
			visited[path] = true
			info, err := d.Info()
			if err != nil {
				return nil
			}
			if info.Mode()&os.ModeSetuid != 0 {
				stat := info.Sys().(*syscall.Stat_t)
				result = append(result, fmt.Sprintf("  %-60s [uid=%d mode=%s]\n",
					path, stat.Uid, info.Mode()))
			}
			return nil
		})
	}
	return result
}

func isWritable(path string) bool {
	return syscall.Access(path, syscall.O_WRONLY) == nil
}
