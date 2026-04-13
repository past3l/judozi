//go:build linux

package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/judozi/judozi/modules/kernel"
	"github.com/judozi/judozi/modules/persistence"
	"github.com/judozi/judozi/pkg/module"
	"github.com/judozi/judozi/pkg/shell"
	"github.com/judozi/judozi/pkg/ui"
)

func main() {
	registry := module.NewRegistry()
	
	registry.Register(kernel.New())
	registry.Register(persistence.New())
	
	ui.ShowBanner()
	
	if len(os.Args) > 1 {
		moduleName := os.Args[1]
		
		if moduleName == "-h" || moduleName == "--help" {
			showHelp(registry)
			return
		}
		if moduleName == "-l" || moduleName == "--list" {
			listModules(registry)
			return
		}
		if moduleName == "-s" || moduleName == "--shell" {
			sh := shell.NewShell(registry)
			sh.Start()
			return
		}
		
		mod, err := registry.Get(moduleName)
		if err != nil {
			ui.Error("Module not found: %s", moduleName)
			fmt.Println()
			listModules(registry)
			os.Exit(1)
		}
		
		if err := mod.Run(os.Args[2:]); err != nil {
			ui.Error("Module execution failed: %v", err)
			os.Exit(1)
		}
		
		if os.Getuid() == 0 {
			fmt.Println()
			ui.Info("Entering interactive shell...")
			sh := shell.NewShell(registry)
			sh.Start()
		}
		
		return
	}
	
	selectAndRun(registry)
}

func showHelp(registry *module.Registry) {
	fmt.Println(ui.Bold + ui.Cyan + "USAGE:" + ui.Reset)
	fmt.Println("  judozi [module] [options]")
	fmt.Println("  judozi --shell              # Start interactive shell")
	fmt.Println()
	fmt.Println(ui.Bold + ui.Cyan + "GLOBAL FLAGS:" + ui.Reset)
	fmt.Println("  -l, --list     List all available modules")
	fmt.Println("  -s, --shell    Start interactive shell")
	fmt.Println("  -h, --help     Show this help message")
	fmt.Println()
	fmt.Println(ui.Bold + ui.Cyan + "EXAMPLES:" + ui.Reset)
	fmt.Println("  judozi                    # Interactive module selection + shell")
	fmt.Println("  judozi --shell            # Direct to interactive shell")
	fmt.Println("  judozi kernel             # Run kernel module + shell")
	fmt.Println("  judozi kernel -list       # List all kernel exploits")
	fmt.Println("  judozi kernel -auto       # Auto-exploit mode")
	fmt.Println("  judozi persistence        # Run persistence module")
	fmt.Println("  judozi persistence -list  # List persistence techniques")
	fmt.Println()
	listModules(registry)
}

func listModules(registry *module.Registry) {
	modules := registry.List()
	
	fmt.Println(ui.Bold + ui.Green + "AVAILABLE MODULES:" + ui.Reset)
	fmt.Println()
	
	categories := make(map[string][]module.Module)
	for _, mod := range modules {
		cat := mod.Category()
		categories[cat] = append(categories[cat], mod)
	}
	
	for cat, mods := range categories {
		fmt.Printf("  %s%s%s\n", ui.Bold+ui.Magenta, strings.ToUpper(cat), ui.Reset)
		fmt.Printf("  %s%s%s\n", ui.Dim, strings.Repeat("─", 78), ui.Reset)
		
		for _, mod := range mods {
			recommendation := ""
			if mod.Category() == "persistence" && os.Getuid() == 0 {
				recommendation = ui.Yellow + ui.Bold + " [RECOMMENDED]" + ui.Reset
			} else if mod.Category() == "persistence" && os.Getuid() != 0 {
				recommendation = ui.Dim + " [Requires root]" + ui.Reset
			}
			
			fmt.Printf("  %s%-15s%s  %s%s\n",
				ui.Yellow+ui.Bold, mod.Name(), ui.Reset,
				mod.Description(), recommendation)
		}
		fmt.Println()
	}
}

func selectAndRun(registry *module.Registry) {
	modules := registry.List()
	
	if len(modules) == 0 {
		ui.Error("No modules available")
		return
	}
	
	fmt.Println(ui.Bold + ui.Cyan + "SELECT MODULE:" + ui.Reset)
	fmt.Println()
	
	for i, mod := range modules {
		fmt.Printf("  %s[%d]%s %s%s%s\n",
			ui.Bold+ui.Cyan, i+1, ui.Reset,
			ui.Yellow+ui.Bold, mod.Name(), ui.Reset)
		fmt.Printf("      %s%s%s - %s\n",
			ui.Magenta, mod.Category(), ui.Reset,
			mod.Description())
		fmt.Println()
	}
	
	fmt.Println(ui.Dim + "  ╔════════════════════════════════════════════════════════════════════╗")
	fmt.Println("  ║ " + ui.Reset + "Enter module number (1-" + fmt.Sprint(len(modules)) + ") or 'q' to quit" + ui.Dim + "                        ║")
	fmt.Println("  ╚════════════════════════════════════════════════════════════════════╝" + ui.Reset)
	fmt.Print("  > ")
	
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return
	}
	
	input := strings.TrimSpace(scanner.Text())
	if input == "q" || input == "quit" {
		return
	}
	
	idx, err := strconv.Atoi(input)
	if err != nil || idx < 1 || idx > len(modules) {
		ui.Error("Invalid selection")
		return
	}
	
	selectedModule := modules[idx-1]
	
	fmt.Println()
	ui.Info("Starting module: %s%s%s", ui.Yellow+ui.Bold, selectedModule.Name(), ui.Reset)
	fmt.Println()
	
	if err := selectedModule.Run([]string{}); err != nil {
		ui.Error("Module execution failed: %v", err)
		os.Exit(1)
	}
	
	fmt.Println()
	ui.Info("Entering interactive shell. Type 'help' for commands.")
	fmt.Println()
	
	sh := shell.NewShell(registry)
	sh.Start()
}
