//go:build linux

package persistence

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/judozi/judozi/modules/persistence/techniques"
	"github.com/judozi/judozi/pkg/module"
	"github.com/judozi/judozi/pkg/ui"
)

// PersistenceModule implements the Module interface for persistence techniques
type PersistenceModule struct{}

func New() module.Module {
	return &PersistenceModule{}
}

func (p *PersistenceModule) Name() string {
	return "persistence"
}

func (p *PersistenceModule) Description() string {
	return "Advanced Persistence Mechanisms for Root Access"
}

func (p *PersistenceModule) Category() string {
	return "persistence"
}

func (p *PersistenceModule) Run(args []string) error {
	// Check if root
	if os.Getuid() != 0 {
		ui.Error("This module requires root privileges")
		ui.Warn("Run the 'kernel' module first to gain root access")
		return fmt.Errorf("root required")
	}
	
	ui.ModuleHeader(p.Name(), p.Category(), p.Description())
	
	fs := flag.NewFlagSet("persistence", flag.ExitOnError)
	list := fs.Bool("list", false, "List all persistence techniques")
	remove := fs.Bool("remove", false, "Remove persistence (cleanup)")
	
	if err := fs.Parse(args); err != nil {
		return err
	}
	
	// Get all available techniques
	techs := techniques.GetAll()
	
	if *list {
		listTechniques(techs)
		return nil
	}
	
	if *remove {
		return removePersistence(techs)
	}
	
	// Check if specific technique number provided
	remaining := fs.Args()
	if len(remaining) > 0 {
		idx, err := strconv.Atoi(remaining[0])
		if err != nil || idx < 1 || idx > len(techs) {
			ui.Error("Invalid technique number")
			listTechniques(techs)
			return fmt.Errorf("invalid selection")
		}
		
		return runTechnique(techs[idx-1])
	}
	
	// Interactive selection
	return interactiveSelect(techs)
}

func listTechniques(techs []techniques.Technique) {
	fmt.Println()
	fmt.Printf("  %s%-4s %-25s %-12s %-10s %s%s\n",
		ui.Bold+ui.Cyan, "ID", "TECHNIQUE", "STEALTH", "REBOOT", "DESCRIPTION", ui.Reset)
	fmt.Printf("  %s%s%s\n", ui.Dim, strings.Repeat("─", 100), ui.Reset)
	
	for i, tech := range techs {
		stealthLevel := getStealthIndicator(tech.StealthLevel())
		rebootSafe := ui.Green + "✓" + ui.Reset
		if !tech.SurvivesReboot() {
			rebootSafe = ui.Red + "✗" + ui.Reset
		}
		
		fmt.Printf("  %s%-4d%s %-25s %-12s %-10s %s\n",
			ui.Yellow, i+1, ui.Reset,
			tech.Name(),
			stealthLevel,
			rebootSafe,
			tech.Description())
	}
	
	fmt.Println()
	fmt.Printf("  %sStealth Levels:%s %sLOW%s | %sMEDIUM%s | %sHIGH%s\n",
		ui.Bold, ui.Reset,
		ui.Red, ui.Reset,
		ui.Yellow, ui.Reset,
		ui.Green, ui.Reset)
	fmt.Printf("  %sReboot Safe:%s %s✓%s Survives reboot | %s✗%s Memory only\n\n",
		ui.Bold, ui.Reset,
		ui.Green, ui.Reset,
		ui.Red, ui.Reset)
}

func getStealthIndicator(level string) string {
	switch strings.ToLower(level) {
	case "low":
		return ui.Red + "LOW    " + ui.Reset
	case "medium":
		return ui.Yellow + "MEDIUM " + ui.Reset
	case "high":
		return ui.Green + "HIGH   " + ui.Reset
	default:
		return ui.Dim + "UNKNOWN" + ui.Reset
	}
}

func interactiveSelect(techs []techniques.Technique) error {
	listTechniques(techs)
	
	fmt.Println(ui.Dim + "  ╔════════════════════════════════════════════════════════════════════╗")
	fmt.Println("  ║ " + ui.Reset + "Enter technique number (1-" + fmt.Sprint(len(techs)) + "), 'all', or 'q' to quit" + ui.Dim + "         ║")
	fmt.Println("  ╚════════════════════════════════════════════════════════════════════╝" + ui.Reset)
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
		return runAllTechniques(techs)
	}
	
	idx, err := strconv.Atoi(input)
	if err != nil || idx < 1 || idx > len(techs) {
		ui.Error("Invalid selection")
		return nil
	}
	
	return runTechnique(techs[idx-1])
}

func runTechnique(tech techniques.Technique) error {
	fmt.Println()
	ui.Info("Deploying: %s%s%s", ui.Yellow+ui.Bold, tech.Name(), ui.Reset)
	ui.Info("Stealth: %s | Reboot Safe: %v", tech.StealthLevel(), tech.SurvivesReboot())
	fmt.Println()
	
	if err := tech.Install(); err != nil {
		ui.Error("Installation failed: %v", err)
		return err
	}
	
	ui.Success("Persistence technique deployed successfully!")
	
	if details := tech.Details(); details != "" {
		fmt.Println()
		fmt.Println(ui.Cyan + "Details:" + ui.Reset)
		fmt.Println(details)
	}
	
	return nil
}

func runAllTechniques(techs []techniques.Technique) error {
	ui.Warn("Installing ALL persistence techniques...")
	fmt.Println()
	
	success := 0
	failed := 0
	
	for i, tech := range techs {
		ui.Info("[%d/%d] Installing: %s", i+1, len(techs), tech.Name())
		
		if err := tech.Install(); err != nil {
			ui.Error("Failed: %v", err)
			failed++
		} else {
			ui.Success("Installed: %s", tech.Name())
			success++
		}
		fmt.Println()
	}
	
	fmt.Println(ui.Bold + "Summary:" + ui.Reset)
	ui.Success("Successful: %d", success)
	if failed > 0 {
		ui.Error("Failed: %d", failed)
	}
	
	return nil
}

func removePersistence(techs []techniques.Technique) error {
	ui.Warn("Removing all persistence mechanisms...")
	fmt.Println()
	
	for _, tech := range techs {
		ui.Info("Removing: %s", tech.Name())
		if err := tech.Remove(); err != nil {
			ui.Error("Failed to remove %s: %v", tech.Name(), err)
		} else {
			ui.Success("Removed: %s", tech.Name())
		}
	}
	
	ui.Success("Cleanup complete")
	return nil
}
