package techniques

import (
	"fmt"
	"os"
	"strings"
)

type DockerEnum struct{}

func NewDockerEnum() Technique      { return &DockerEnum{} }
func (d *DockerEnum) Name() string         { return "Container & Docker Recon" }
func (d *DockerEnum) Description() string  { return "Detect containers, docker socket, k8s, namespaces, cgroups" }
func (d *DockerEnum) StealthLevel() string { return "passive" }

func (d *DockerEnum) Run() string {
	var b strings.Builder

	b.WriteString("[+] CONTAINER DETECTION\n")
	// Check /.dockerenv
	if _, err := os.Stat("/.dockerenv"); err == nil {
		b.WriteString("  *** RUNNING INSIDE DOCKER CONTAINER (/.dockerenv exists) ***\n")
	}
	// Check cgroup
	data, _ := os.ReadFile("/proc/1/cgroup")
	if strings.Contains(string(data), "docker") {
		b.WriteString("  *** DOCKER CGROUP DETECTED in /proc/1/cgroup ***\n")
	}
	if strings.Contains(string(data), "kubepods") {
		b.WriteString("  *** KUBERNETES POD DETECTED in /proc/1/cgroup ***\n")
	}
	if strings.Contains(string(data), "lxc") {
		b.WriteString("  *** LXC CONTAINER DETECTED ***\n")
	}
	b.WriteString(cmd("cat /proc/1/cgroup 2>/dev/null"))
	b.WriteString("\n")

	b.WriteString("[+] DOCKER SOCKET\n")
	for _, sock := range []string{"/var/run/docker.sock", "/run/docker.sock"} {
		fi, err := os.Stat(sock)
		if err == nil {
			b.WriteString(fmt.Sprintf("  *** DOCKER SOCKET EXISTS: %s [%s] ***\n", sock, fi.Mode()))
			// Check if we can read it
			f, err2 := os.Open(sock)
			if err2 == nil {
				f.Close()
				b.WriteString("  *** DOCKER SOCKET IS READABLE — potential escape! ***\n")
			}
		}
	}
	b.WriteString("\n")

	b.WriteString("[+] PRIVILEGED / CAPABILITIES\n")
	b.WriteString(cmd("cat /proc/self/status 2>/dev/null | grep -E 'Cap|Uid|Gid|Groups'"))
	b.WriteString(cmd("capsh --print 2>/dev/null"))
	b.WriteString(cmd("cat /proc/self/attr/current 2>/dev/null"))
	b.WriteString("\n")

	b.WriteString("[+] NAMESPACES\n")
	b.WriteString(cmd("ls -la /proc/self/ns/ 2>/dev/null"))
	b.WriteString(cmd("lsns 2>/dev/null"))
	b.WriteString("\n")

	b.WriteString("[+] CGROUP VERSION & WRITABILITY\n")
	b.WriteString(cmd("stat -fc %T /sys/fs/cgroup/ 2>/dev/null"))
	b.WriteString(cmd("cat /proc/1/cgroup 2>/dev/null"))
	b.WriteString(cmd("find /sys/fs/cgroup -name 'release_agent' -writable 2>/dev/null"))
	b.WriteString(cmd("find /sys/fs/cgroup -name 'notify_on_release' 2>/dev/null | head -10"))
	b.WriteString("\n")

	b.WriteString("[+] MOUNTED FILESYSTEMS\n")
	b.WriteString(cmd("mount 2>/dev/null || cat /proc/mounts 2>/dev/null"))
	b.WriteString(cmd("cat /proc/1/mountinfo 2>/dev/null | head -30"))
	b.WriteString("\n")

	b.WriteString("[+] INTERESTING HOST MOUNTS\n")
	b.WriteString(cmd("cat /proc/mounts 2>/dev/null | grep -E '/etc|/proc|/sys|/dev|host|overlay'"))
	b.WriteString("\n")

	b.WriteString("[+] KUBERNETES ENV & TOKENS\n")
	b.WriteString(cmd("env | grep -E 'KUBERNETES|K8S|KUBE'"))
	b.WriteString(cmd("cat /var/run/secrets/kubernetes.io/serviceaccount/token 2>/dev/null"))
	b.WriteString(cmd("cat /var/run/secrets/kubernetes.io/serviceaccount/ca.crt 2>/dev/null | openssl x509 -text -noout 2>/dev/null | head -20"))
	b.WriteString("\n")

	b.WriteString("[+] METADATA SERVICE CHECK (cloud)\n")
	b.WriteString(cmd("curl -s -m 2 http://169.254.169.254/latest/meta-data/ 2>/dev/null | head -20"))
	b.WriteString(cmd("curl -s -m 2 -H 'Metadata-Flavor: Google' http://169.254.169.254/computeMetadata/v1/ 2>/dev/null | head -20"))
	b.WriteString(cmd("curl -s -m 2 -H 'Metadata: true' 'http://169.254.169.254/metadata/instance?api-version=2021-02-01' 2>/dev/null | head -20"))
	b.WriteString("\n")

	b.WriteString("[+] DOCKER CONTAINERS (if docker available)\n")
	b.WriteString(cmd("docker ps -a 2>/dev/null"))
	b.WriteString(cmd("docker images 2>/dev/null"))
	b.WriteString(cmd("docker network ls 2>/dev/null"))
	b.WriteString("\n")

	return b.String()
}

type EnvSecrets struct{}

func NewEnvSecrets() Technique       { return &EnvSecrets{} }
func (e *EnvSecrets) Name() string         { return "Environment & Secret Variables" }
func (e *EnvSecrets) Description() string  { return "env vars, /proc env, secrets in processes memory" }
func (e *EnvSecrets) StealthLevel() string { return "passive" }

func (e *EnvSecrets) Run() string {
	var b strings.Builder

	b.WriteString("[+] CURRENT PROCESS ENVIRONMENT\n")
	b.WriteString(cmd("env 2>/dev/null | sort"))
	b.WriteString("\n")

	b.WriteString("[+] SECRET KEYWORDS IN ENVIRONMENT\n")
	b.WriteString(cmd("env 2>/dev/null | grep -Ei 'pass|secret|token|key|api|auth|cred|db_|database|aws|azure|gcp|google|slack|github|gitlab'"))
	b.WriteString("\n")

	b.WriteString("[+] /proc ENVIRONMENT OF RUNNING PROCESSES\n")
	b.WriteString(readProcEnvs())
	b.WriteString("\n")

	b.WriteString("[+] PROC 1 ENVIRONMENT\n")
	b.WriteString(readEnvFile("/proc/1/environ"))
	b.WriteString("\n")

	b.WriteString("[+] PROC cmdline (all)\n")
	b.WriteString(cmd("cat /proc/*/cmdline 2>/dev/null | tr '\\0' ' ' | tr '\\n' '\\n' | grep -Ei 'pass|secret|token|key' | head -20"))
	b.WriteString("\n")

	return b.String()
}

func readEnvFile(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Sprintf("  [cannot read %s]\n", path)
	}
	lines := strings.Split(strings.ReplaceAll(string(data), "\x00", "\n"), "\n")
	var b strings.Builder
	keywords := []string{"pass", "secret", "token", "key", "api", "auth", "cred", "db_", "aws", "azure", "gcp"}
	for _, line := range lines {
		lower := strings.ToLower(line)
		for _, kw := range keywords {
			if strings.Contains(lower, kw) {
				b.WriteString("  " + line + "\n")
				break
			}
		}
	}
	return b.String()
}

func readProcEnvs() string {
	var b strings.Builder
	entries, _ := os.ReadDir("/proc")
	count := 0
	for _, e := range entries {
		if !e.IsDir() || count > 20 {
			break
		}
		var pid int
		if _, err := fmt.Sscanf(e.Name(), "%d", &pid); err != nil {
			continue
		}
		result := readEnvFile(fmt.Sprintf("/proc/%d/environ", pid))
		if result != "" {
			b.WriteString(fmt.Sprintf("  --- pid %d ---\n%s", pid, result))
			count++
		}
	}
	return b.String()
}
