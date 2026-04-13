package techniques

import (
	"strings"
)

type CronEnum struct{}

func NewCronEnum() Technique        { return &CronEnum{} }
func (c *CronEnum) Name() string         { return "Cron & Scheduled Tasks" }
func (c *CronEnum) Description() string  { return "All cron jobs, at jobs, anacron, systemd timers" }
func (c *CronEnum) StealthLevel() string { return "passive" }

func (c *CronEnum) Run() string {
	var b strings.Builder

	b.WriteString("[+] SYSTEM-WIDE CRONTAB\n")
	b.WriteString(cmd("cat /etc/crontab 2>/dev/null"))
	b.WriteString("\n")

	b.WriteString("[+] CRON.D DIRECTORY\n")
	b.WriteString(cmd("ls -la /etc/cron.d/ 2>/dev/null"))
	b.WriteString(cmd("cat /etc/cron.d/* 2>/dev/null"))
	b.WriteString("\n")

	b.WriteString("[+] CRON HOURLY/DAILY/WEEKLY/MONTHLY\n")
	b.WriteString(cmd("ls -la /etc/cron.hourly/ /etc/cron.daily/ /etc/cron.weekly/ /etc/cron.monthly/ 2>/dev/null"))
	b.WriteString("\n")

	b.WriteString("[+] USER CRONTABS\n")
	b.WriteString(cmd("crontab -l 2>/dev/null"))
	b.WriteString(cmd("cat /var/spool/cron/crontabs/* 2>/dev/null"))
	b.WriteString(cmd("cat /var/spool/cron/* 2>/dev/null"))
	b.WriteString(cmd("ls -la /var/spool/cron/crontabs/ 2>/dev/null"))
	b.WriteString("\n")

	b.WriteString("[+] ANACRON\n")
	b.WriteString(cmd("cat /etc/anacrontab 2>/dev/null"))
	b.WriteString("\n")

	b.WriteString("[+] SYSTEMD TIMERS\n")
	b.WriteString(cmd("systemctl list-timers --all 2>/dev/null"))
	b.WriteString(cmd("find /etc/systemd /lib/systemd /usr/lib/systemd -name '*.timer' 2>/dev/null | xargs grep -l '' | head -20"))
	b.WriteString("\n")

	b.WriteString("[+] CRON LOG\n")
	b.WriteString(cmd("cat /var/log/cron 2>/dev/null | tail -30"))
	b.WriteString(cmd("cat /var/log/syslog 2>/dev/null | grep CRON | tail -20"))
	b.WriteString(cmd("journalctl -u cron 2>/dev/null | tail -20"))
	b.WriteString("\n")

	b.WriteString("[+] WRITABLE CRON SCRIPTS\n")
	b.WriteString(cmd("find /etc/cron* /var/spool/cron -writable 2>/dev/null"))
	b.WriteString(cmd("for f in $(find /etc/cron* /var/spool/cron -type f 2>/dev/null); do ls -la \"$f\"; done"))
	b.WriteString("\n")

	b.WriteString("[+] AT JOBS\n")
	b.WriteString(cmd("atq 2>/dev/null"))
	b.WriteString(cmd("ls -la /var/spool/at/ 2>/dev/null"))
	b.WriteString("\n")

	return b.String()
}

type ServiceEnum struct{}

func NewServiceEnum() Technique       { return &ServiceEnum{} }
func (s *ServiceEnum) Name() string         { return "Service Enumeration" }
func (s *ServiceEnum) Description() string  { return "Systemd services, init.d, running daemons, version info" }
func (s *ServiceEnum) StealthLevel() string { return "passive" }

func (s *ServiceEnum) Run() string {
	var b strings.Builder

	b.WriteString("[+] SYSTEMD SERVICES\n")
	b.WriteString(cmd("systemctl list-units --type=service --state=running 2>/dev/null"))
	b.WriteString("\n")

	b.WriteString("[+] ALL SYSTEMD UNITS\n")
	b.WriteString(cmd("systemctl list-units --all 2>/dev/null | head -60"))
	b.WriteString("\n")

	b.WriteString("[+] FAILED SERVICES\n")
	b.WriteString(cmd("systemctl list-units --state=failed 2>/dev/null"))
	b.WriteString("\n")

	b.WriteString("[+] ENABLED SERVICES\n")
	b.WriteString(cmd("systemctl list-unit-files --state=enabled 2>/dev/null"))
	b.WriteString("\n")

	b.WriteString("[+] INIT.D SCRIPTS\n")
	b.WriteString(cmd("ls -la /etc/init.d/ 2>/dev/null"))
	b.WriteString(cmd("find /etc/init.d/ -writable 2>/dev/null"))
	b.WriteString("\n")

	b.WriteString("[+] SERVICE FILES (interesting ones)\n")
	b.WriteString(cmd("find /etc/systemd /lib/systemd /usr/lib/systemd -name '*.service' 2>/dev/null | head -40"))
	b.WriteString("\n")

	b.WriteString("[+] WRITABLE SYSTEMD SERVICE FILES\n")
	b.WriteString(cmd("find /etc/systemd /lib/systemd /usr/lib/systemd -name '*.service' -writable 2>/dev/null"))
	b.WriteString("\n")

	b.WriteString("[+] APPLICATIONS / VERSIONS\n")
	apps := []string{
		"apache2", "nginx", "mysql", "mysqld", "postgres", "mongod",
		"redis-server", "memcached", "tomcat", "jetty", "node", "npm",
		"docker", "containerd", "kubelet", "sshd", "vsftpd", "proftpd",
		"named", "bind9", "postfix", "sendmail", "exim",
	}
	for _, app := range apps {
		out := cmdSilent(app + " -v 2>&1 | head -1 || " + app + " --version 2>&1 | head -1")
		if out != "" && out != "\n" {
			b.WriteString("  " + app + ": " + out)
		}
	}
	b.WriteString("\n")

	b.WriteString("[+] PACKAGE MANAGER VERSIONS\n")
	b.WriteString(cmd("dpkg -l 2>/dev/null | grep -E '^ii' | wc -l"))
	b.WriteString(cmd("rpm -qa 2>/dev/null | wc -l"))
	b.WriteString(cmd("dpkg -l 2>/dev/null | grep -E 'lib|sudo|ssh|curl|wget|openssl' | head -30"))
	b.WriteString("\n")

	return b.String()
}

func cmdSilent(c string) string {
	return cmd(c)
}
