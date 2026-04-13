package techniques

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

type NetworkInfo struct{}

func NewNetworkInfo() Technique    { return &NetworkInfo{} }
func (n *NetworkInfo) Name() string         { return "Network Enumeration" }
func (n *NetworkInfo) Description() string  { return "Interfaces, routes, ports, connections, ARP, DNS, proxies" }
func (n *NetworkInfo) StealthLevel() string { return "passive" }

func (n *NetworkInfo) Run() string {
	var b strings.Builder

	b.WriteString("[+] NETWORK INTERFACES\n")
	ifaces, _ := net.Interfaces()
	for _, iface := range ifaces {
		addrs, _ := iface.Addrs()
		b.WriteString(fmt.Sprintf("  %-20s flags=%v\n", iface.Name, iface.Flags))
		for _, a := range addrs {
			b.WriteString(fmt.Sprintf("    %s\n", a))
		}
	}
	b.WriteString(cmd("ip addr 2>/dev/null || ifconfig 2>/dev/null"))
	b.WriteString("\n")

	b.WriteString("[+] ROUTING TABLE\n")
	b.WriteString(cmd("ip route 2>/dev/null || route -n 2>/dev/null"))
	b.WriteString(cmd("ip -6 route 2>/dev/null"))
	b.WriteString("\n")

	b.WriteString("[+] ARP TABLE / NEIGHBORS\n")
	b.WriteString(cmd("arp -an 2>/dev/null"))
	b.WriteString(cmd("ip neigh 2>/dev/null"))
	b.WriteString("\n")

	b.WriteString("[+] LISTENING PORTS\n")
	b.WriteString(cmd("ss -tlnp 2>/dev/null || netstat -tlnp 2>/dev/null"))
	b.WriteString("\n")

	b.WriteString("[+] ALL CONNECTIONS\n")
	b.WriteString(cmd("ss -anp 2>/dev/null || netstat -anp 2>/dev/null"))
	b.WriteString("\n")

	b.WriteString("[+] UDP PORTS\n")
	b.WriteString(cmd("ss -ulnp 2>/dev/null || netstat -ulnp 2>/dev/null"))
	b.WriteString("\n")

	b.WriteString("[+] DNS CONFIGURATION\n")
	b.WriteString(cmdFile("/etc/resolv.conf"))
	b.WriteString(cmdFile("/etc/nsswitch.conf"))
	b.WriteString("\n")

	b.WriteString("[+] HOSTS FILE\n")
	b.WriteString(cmdFile("/etc/hosts"))
	b.WriteString("\n")

	b.WriteString("[+] PROXY SETTINGS\n")
	for _, env := range []string{"http_proxy", "https_proxy", "HTTP_PROXY", "HTTPS_PROXY", "no_proxy", "ALL_PROXY"} {
		v := os.Getenv(env)
		if v != "" {
			b.WriteString(fmt.Sprintf("  %s = %s\n", env, v))
		}
	}
	b.WriteString(cmdFile("/etc/environment"))
	b.WriteString("\n")

	b.WriteString("[+] FIREWALL RULES\n")
	b.WriteString(cmd("iptables -L -n -v 2>/dev/null"))
	b.WriteString(cmd("ip6tables -L -n -v 2>/dev/null"))
	b.WriteString(cmd("nft list ruleset 2>/dev/null"))
	b.WriteString(cmd("ufw status verbose 2>/dev/null"))
	b.WriteString("\n")

	b.WriteString("[+] TCPDUMP INTERFACES\n")
	b.WriteString(cmd("tcpdump --list-interfaces 2>/dev/null"))
	b.WriteString("\n")

	b.WriteString("[+] NETWORK NAMESPACES\n")
	b.WriteString(cmd("ip netns list 2>/dev/null"))
	b.WriteString("\n")

	b.WriteString("[+] NFS MOUNTS\n")
	b.WriteString(cmd("cat /etc/exports 2>/dev/null"))
	b.WriteString(cmd("showmount -e 2>/dev/null"))
	b.WriteString("\n")

	b.WriteString("[+] /proc/net/tcp PARSED (IPv4 listening)\n")
	b.WriteString(parseProcNetTCP())

	return b.String()
}

func parseProcNetTCP() string {
	var b strings.Builder
	f, err := os.Open("/proc/net/tcp")
	if err != nil {
		return "  [unavailable]\n"
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	scanner.Scan() // skip header
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 4 || fields[3] != "0A" {
			continue // only LISTEN state (0A)
		}
		raw := fields[1]
		parts := strings.Split(raw, ":")
		if len(parts) != 2 {
			continue
		}
		// decode little-endian hex ip
		portHex := parts[1]
		port := 0
		fmt.Sscanf(portHex, "%X", &port)
		b.WriteString(fmt.Sprintf("  0.0.0.0:%d\n", port))
	}
	return b.String()
}
