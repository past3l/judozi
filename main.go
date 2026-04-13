//go:build linux

package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/judozi/judozi/pkg/exploit"
	"github.com/judozi/judozi/pkg/kernel"
	"github.com/judozi/judozi/pkg/output"
	"github.com/judozi/judozi/pkg/vulndb"
)

func main() {
	mirror := flag.String("mirror", "", "Base URL or local path for exploit source files (optional, exploits have default URLs)")
	list := flag.Bool("list", false, "List all exploits in database")
	cve := flag.String("cve", "", "Try a specific CVE (e.g. CVE-2022-0847)")
	auto := flag.Bool("auto", false, "Automatically try all matching exploits")
	noCleanup := flag.Bool("no-cleanup", false, "Keep temp files after execution")
	flag.Parse()

	output.Banner()

	db, err := vulndb.Load()
	if err != nil {
		output.Fatal("Failed to load exploit database: %v", err)
	}

	if *list {
		listExploits(db)
		return
	}

	if exploit.IsRoot() {
		output.Warn("Already running as root!")
		return
	}

	kv, err := kernel.GetVersion()
	if err != nil {
		output.Fatal("Kernel version detection failed: %v", err)
	}
	output.Info("Kernel: %s%s%s (%s%s%s)", output.Bold+output.Yellow, kv.String(), output.Reset, output.Dim, kv.Raw, output.Reset)

	arch := kernel.GetArch()
	output.Info("Arch:   %s%s%s", output.Bold+output.Cyan, arch, output.Reset)

	if kernel.IsContainer() {
		output.Warn("Running inside a container — exploits target the host kernel")
	}

	var matches []vulndb.Exploit
	if *cve != "" {
		e := db.GetByID(*cve)
		if e == nil {
			output.Fatal("CVE %s not found in database", *cve)
		}
		matches = []vulndb.Exploit{*e}
	} else {
		matches = db.Search(kv, arch)
	}

	if len(matches) == 0 {
		output.Error("No matching exploits for kernel %s (%s)", kv.String(), arch)
		return
	}

	output.Success("Found %s%d%s potential exploit(s):\n", output.Bold+output.Green, len(matches), output.Reset)
	printExploitTable(matches)

	selected := matches
	if !*auto {
		selected = interactiveSelect(matches)
		if len(selected) == 0 {
			return
		}
	}

	runExploits(selected, *mirror, *noCleanup)
}

func listExploits(db *vulndb.Database) {
	fmt.Printf("\n  %s%-17s %-28s %-10s %-10s %s%s\n",
		output.Bold+output.Cyan, "CVE", "NAME", "MIN", "MAX", "TAGS", output.Reset)
	fmt.Printf("  %s%s%s\n", output.Dim, strings.Repeat("─", 90), output.Reset)
	for _, e := range db.Exploits {
		fmt.Printf("  %s%-17s%s %s%-28s%s %s%-10s%s %s%-10s%s %s%s%s\n",
			output.Magenta, e.ID, output.Reset,
			output.Bold, e.Name, output.Reset,
			output.Green, e.MinKernel, output.Reset,
			output.Red, e.MaxKernel, output.Reset,
			output.Cyan, strings.Join(e.Tags, ", "), output.Reset)
	}
	fmt.Printf("\n  %sTotal: %d exploits%s\n\n", output.Bold+output.Yellow, len(db.Exploits), output.Reset)
}

func printExploitTable(exploits []vulndb.Exploit) {
	for i, e := range exploits {
		fmt.Printf("  %s[%d]%s %s%-28s%s %s%s%s\n",
			output.Yellow+output.Bold, i+1, output.Reset,
			output.Bold+output.Cyan, e.Name, output.Reset, 
			output.Magenta, e.ID, output.Reset)
		fmt.Printf("      %s%s%s\n", output.Dim, e.Description, output.Reset)
		fmt.Printf("      %sKernel:%s %s%s%s → %s%s%s  %s|%s  %sTags:%s %s%s%s\n\n",
			output.Blue, output.Reset,
			output.Green, e.MinKernel, output.Reset,
			output.Red, e.MaxKernel, output.Reset,
			output.Dim, output.Reset,
			output.Yellow, output.Reset,
			output.Cyan, strings.Join(e.Tags, ", "), output.Reset)
	}
}

func interactiveSelect(matches []vulndb.Exploit) []vulndb.Exploit {
	fmt.Printf("  %s╔════════════════════════════════════════════════════════════════════╗%s\n", 
		output.Yellow+output.Bold, output.Reset)
	fmt.Printf("  %s║%s Enter exploit number (%s1-%d%s), '%sall%s' to try all, or '%sq%s' to quit %s║%s\n", 
		output.Yellow+output.Bold, output.Reset,
		output.Bold, len(matches), output.Reset,
		output.Green+output.Bold, output.Reset,
		output.Red+output.Bold, output.Reset,
		output.Yellow+output.Bold, output.Reset)
	fmt.Printf("  %s╚════════════════════════════════════════════════════════════════════╝%s\n  %s>%s ", 
		output.Yellow+output.Bold, output.Reset,
		output.Cyan+output.Bold, output.Reset)

	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return nil
	}
	input := strings.TrimSpace(scanner.Text())

	switch strings.ToLower(input) {
	case "q", "quit", "exit", "":
		return nil
	case "all", "a":
		return matches
	default:
		n, err := strconv.Atoi(input)
		if err != nil || n < 1 || n > len(matches) {
			output.Error("Invalid selection: %s", input)
			return nil
		}
		return []vulndb.Exploit{matches[n-1]}
	}
}

func runExploits(exploits []vulndb.Exploit, mirror string, noCleanup bool) {
	workDir, err := os.MkdirTemp("", "kexploit-*")
	if err != nil {
		output.Fatal("Failed to create temp directory: %v", err)
	}

	runner := exploit.NewRunner(mirror, workDir)
	if !noCleanup {
		defer runner.Cleanup()
	} else {
		output.Info("Work directory: %s (no-cleanup mode)", workDir)
	}

	for i, e := range exploits {
		if len(exploits) > 1 {
			fmt.Printf("\n%s╔═══════════════════════════════════════════════════════════════════╗%s\n",
				output.Cyan+output.Bold, output.Reset)
			fmt.Printf("%s║%s [%s%d%s/%s%d%s] %s%s%s (%s%s%s) %s║%s\n",
				output.Cyan+output.Bold, output.Reset,
				output.Yellow+output.Bold, i+1, output.Reset,
				output.Yellow+output.Bold, len(exploits), output.Reset,
				output.Bold+output.White, e.Name, output.Reset,
				output.Magenta, e.ID, output.Reset,
				output.Cyan+output.Bold, output.Reset)
			fmt.Printf("%s╚═══════════════════════════════════════════════════════════════════╝%s\n",
				output.Cyan+output.Bold, output.Reset)
		} else {
			fmt.Printf("\n%s── %s (%s) %s\n",
				output.Cyan, e.Name, e.ID, output.Reset)
		}

		if err := runner.Run(e); err != nil {
			output.Error("%s failed: %v", e.Name, err)
			continue
		}

		output.Success("Exploit %s%s%s completed", output.Bold+output.Green, e.Name, output.Reset)
		if exploit.IsRoot() {
			output.RootBanner()
			return
		}
	}

	if len(exploits) > 1 {
		output.Warn("All exploits attempted — check if any spawned a root shell")
	}
}
