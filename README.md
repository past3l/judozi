# KExploit - Kernel Exploit Suggester & Auto-Exploiter

Zero-dependency Linux kernel exploit suggester and automated exploitation tool written in Go.

## Features

- **Zero dependencies** - Single static binary (~5MB), works on any Linux system
- **23 kernel exploits** - From kernel 2.6.22 to 6.7.1
- **Auto-detect** - Detects kernel version and suggests matching exploits
- **Auto-download** - Downloads exploit source from public repos (exploit-db, GitHub)
- **Auto-compile** - Compiles on-the-fly using system gcc
- **Auto-execute** - Runs exploit to attempt privilege escalation
- **Educational** - Always downloads and compiles source code, never uses pre-built binaries

## Usage

```bash
# List all exploits
./kexploit -list

# Auto-detect kernel and show matching exploits
./kexploit

# Try a specific CVE
./kexploit -cve CVE-2022-0847

# Try all matching exploits automatically
./kexploit -auto

# Use custom mirror/local path
./kexploit -mirror ./my-exploits
```

## Requirements on Target System

- `gcc` (Alpine: `apk add gcc musl-dev`, Ubuntu: `apt install gcc`)
- `bash` or `sh` (for shell-based exploits)

## Supported Exploits

| CVE | Name | Kernel Range | Tags |
|-----|------|--------------|------|
| CVE-2024-1086 | nf_tables UAF | 3.15.0 - 6.7.1 | nftables, uaf |
| CVE-2023-2640 | GameOver(lay) | 5.15.0 - 6.2.0 | overlayfs, ubuntu |
| CVE-2022-0847 | DirtyPipe | 5.8.0 - 5.16.10 | pipe, write-anywhere |
| CVE-2022-0185 | fsconfig Heap Overflow | 5.1.0 - 5.16.1 | fsconfig, userns |
| CVE-2021-3493 | OverlayFS Ubuntu | 3.13.0 - 5.10.99 | overlayfs, ubuntu |
| CVE-2021-22555 | Netfilter xt_compat | 2.6.19 - 5.12.0 | netfilter, oob-write |
| CVE-2019-13272 | ptrace_link | 4.10.0 - 5.1.17 | ptrace, credentials |
| CVE-2017-16995 | eBPF Arbitrary R/W | 4.4.0 - 4.14.7 | ebpf, sign-extension |
| CVE-2016-5195 | DirtyCow | 2.6.22 - 4.8.2 | cow, race-condition |
| ... | ... | ... | 23 total exploits |

## Build

```bash
# Static binary, zero dependencies
CGO_ENABLED=0 go build -ldflags="-s -w" -o kexploit .

# Cross-compile for different architectures
GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-s -w" -o kexploit-arm64 .
```

## Security & Legal

**Educational purposes only.** Only use on systems you own or have explicit permission to test.

This tool:
- ✅ Always downloads source code from public repositories
- ✅ Compiles exploits locally using gcc
- ✅ Transparent - you can inspect what it downloads
- ❌ Never downloads pre-compiled binaries
- ❌ Does not hide its activity

## How It Works

1. Reads `/proc/version` to detect kernel version
2. Searches embedded exploit database for matches
3. Downloads exploit source from public GitHub repos
4. Compiles using local gcc
5. Executes compiled binary

## License

Educational use only. Exploit code belongs to original authors (linked in references).
