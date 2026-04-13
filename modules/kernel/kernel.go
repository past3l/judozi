//go:build linux

package kernel

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/judozi/judozi/modules/kernel/pkg/exploit"
	"github.com/judozi/judozi/modules/kernel/pkg/kernel"
	"github.com/judozi/judozi/modules/kernel/pkg/output"
	"github.com/judozi/judozi/modules/kernel/pkg/vulndb"
	"github.com/judozi/judozi/pkg/module"
	"github.com/judozi/judozi/pkg/ui"
)

// KernelModule implements the Module interface for kernel privilege escalation
type KernelModule struct{}

func init() {
	// Auto-register this module
	// Will be called from main
}

// New creates a new kernel module instance
func New() module.Module {
	return &KernelModule{}
}

func (k *KernelModule) Name() string {
	return "kernel"
}

func (k *KernelModule) Description() string {
	return "Automated Linux Kernel Privilege Escalation"
}

func (k *KernelModule) Category() string {
	return "privesc"
}

func (k *KernelModule) Run(args []string) error {
	ui.ModuleHeader(k.Name(), k.Category(), k.Description())
	
	fs := flag.NewFlagSet("kernel", flag.ExitOnError)
	mirror := fs.String("mirror", "", "Base URL or local path for exploit source files (optional)")
	list := fs.Bool("list", false, "List all exploits in database")
	cve := fs.String("cve", "", "Try a specific CVE (e.g. CVE-2022-0847)")
	auto := fs.Bool("auto", false, "Automatically try all matching exploits")
	noCleanup := fs.Bool("no-cleanup", false, "Keep temp files after execution")
	
	if err := fs.Parse(args); err != nil {
		return err
	}

	output.Banner()

	db, err := vulndb.Load()
	if err != nil {
		return fmt.Errorf("failed to load exploit database: %v", err)
	}

	if *list {
		listExploits(db)
		return nil
	}

	if exploit.IsRoot() {
		output.Warn("Already running as root!")
		return nil
	}

	kv, err := kernel.GetVersion()
	if err != nil {
		return fmt.Errorf("kernel version detection failed: %v", err)
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
			return fmt.Errorf("CVE %s not found in database", *cve)
		}
		matches = []vulndb.Exploit{*e}
	} else {
		matches = db.Search(kv, arch)
	}

	if len(matches) == 0 {
		output.Error("No matching exploits for kernel %s (%s)", kv.String(), arch)
		return nil
	}

	output.Success("Found %s%d%s potential exploit(s):\n", output.Bold+output.Green, len(matches), output.Reset)
	printExploitTable(matches)

	selected := matches
	if !*auto {
		selected = interactiveSelect(matches)
		if len(selected) == 0 {
			return nil
		}
	}

	runExploits(selected, *mirror, *noCleanup)
	return nil
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
			output.Cyan, strings.Join(e.Tags, " "), output.Reset)
	}
	fmt.Printf("\n  %sTotal: %d exploits%s\n\n", output.Bold+output.Yellow, len(db.Exploits), output.Reset)
}

func printExploitTable(exploits []vulndb.Exploit) {
	for i, e := range exploits {
		fmt.Printf("  %s[%d]%s %s%s%s\n", output.Bold+output.Cyan, i+1, output.Reset, output.Bold, e.Name, output.Reset)
		fmt.Printf("      %s%s%s\n", output.Magenta, e.ID, output.Reset)
		fmt.Printf("      %s\n", e.Description)
		fmt.Printf("      Kernel: %s%s%s → %s%s%s  |  Tags: %s%s%s\n\n",
			output.Green, e.MinKernel, output.Reset,
			output.Red, e.MaxKernel, output.Reset,
			output.Cyan, strings.Join(e.Tags, ", "), output.Reset)
	}
}

func interactiveSelect(matches []vulndb.Exploit) []vulndb.Exploit {
	fmt.Println(output.Dim + "  ╔════════════════════════════════════════════════════════════════════╗")
	fmt.Println("  ║ " + output.Reset + "Enter exploit number (1-" + fmt.Sprint(len(matches)) + "), 'all' to try all, or 'q' to quit" + output.Dim + " ║")
	fmt.Println("  ╚════════════════════════════════════════════════════════════════════╝" + output.Reset)
	fmt.Print("  > ")

	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return nil
	}

	input := strings.TrimSpace(scanner.Text())
	if input == "q" || input == "quit" {
		return nil
	}
	if input == "all" {
		return matches
	}

	idx, err := strconv.Atoi(input)
	if err != nil || idx < 1 || idx > len(matches) {
		output.Error("Invalid selection")
		return nil
	}

	return []vulndb.Exploit{matches[idx-1]}
}

func runExploits(exploits []vulndb.Exploit, mirror string, noCleanup bool) {
	for i, e := range exploits {
		fmt.Println()
		output.Info("Attempting exploit %s[%d/%d]%s: %s%s%s (%s)",
			output.Bold, i+1, len(exploits), output.Reset,
			output.Yellow+output.Bold, e.Name, output.Reset, e.ID)

		if err := exploit.Run(e, mirror, noCleanup); err != nil {
			output.Error("Exploit failed: %v", err)
			continue
		}

		if exploit.IsRoot() {
			output.Success("Privilege escalation successful!")
			output.RootBanner()
			os.Exit(0)
		}
	}

	output.Error("All exploits failed")
}
