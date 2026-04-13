package shell

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/judozi/judozi/pkg/module"
	"github.com/judozi/judozi/pkg/ui"
)

// Shell represents an interactive shell session
type Shell struct {
	registry *module.Registry
	history  []string
	running  bool
}

// NewShell creates a new interactive shell
func NewShell(registry *module.Registry) *Shell {
	return &Shell{
		registry: registry,
		history:  make([]string, 0),
		running:  true,
	}
}

// Start begins the interactive shell session
func (s *Shell) Start() {
	ui.Info("Entering interactive shell. Type %shelp%s for available commands.", ui.Yellow+ui.Bold, ui.Reset)
	fmt.Println()
	
	scanner := bufio.NewScanner(os.Stdin)
	
	for s.running {
		s.printPrompt()
		
		if !scanner.Scan() {
			break
		}
		
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		
		s.history = append(s.history, line)
		
		if err := s.execute(line); err != nil {
			ui.Error("%v", err)
		}
		fmt.Println()
	}
}

func (s *Shell) printPrompt() {
	hostname, _ := os.Hostname()
	user := "user"
	if os.Getuid() == 0 {
		user = ui.Red + ui.Bold + "root" + ui.Reset
	} else {
		user = ui.Green + user + ui.Reset
	}
	
	cwd, _ := os.Getwd()
	if cwd == "" {
		cwd = "~"
	}
	
	fmt.Printf("%s@%s%s:%s%s%s$ ", user, ui.Cyan, hostname, ui.Blue, cwd, ui.Reset)
}

func (s *Shell) execute(line string) error {
	parts := strings.Fields(line)
	if len(parts) == 0 {
		return nil
	}
	
	cmd := parts[0]
	args := parts[1:]
	
	// Check for shell built-in commands
	switch cmd {
	case "help", "?":
		return s.showHelp()
	case "exit", "quit":
		s.running = false
		ui.Success("Goodbye!")
		return nil
	case "clear", "cls":
		fmt.Print("\033[H\033[2J")
		return nil
	case "history":
		return s.showHistory()
	case "modules", "list":
		return s.listModules()
	case "use", "module", "run":
		if len(args) == 0 {
			return fmt.Errorf("usage: %s <module_name> [args...]", cmd)
		}
		return s.runModule(args[0], args[1:])
	case "background", "bg":
		ui.Warn("Background execution not yet implemented")
		return nil
	default:
		// Try to execute as system command
		return s.executeSystemCommand(parts)
	}
}

func (s *Shell) showHelp() error {
	fmt.Println(ui.Bold + ui.Cyan + "JUDOZI INTERACTIVE SHELL" + ui.Reset)
	fmt.Println()
	fmt.Println(ui.Bold + "Framework Commands:" + ui.Reset)
	fmt.Println("  " + ui.Yellow + "help, ?" + ui.Reset + "              Show this help message")
	fmt.Println("  " + ui.Yellow + "modules, list" + ui.Reset + "        List all available modules")
	fmt.Println("  " + ui.Yellow + "use <module> [args]" + ui.Reset + "  Run a specific module")
	fmt.Println("  " + ui.Yellow + "run <module> [args]" + ui.Reset + "  Alias for 'use'")
	fmt.Println("  " + ui.Yellow + "module <name> [args]" + ui.Reset + " Alias for 'use'")
	fmt.Println("  " + ui.Yellow + "history" + ui.Reset + "             Show command history")
	fmt.Println("  " + ui.Yellow + "clear, cls" + ui.Reset + "          Clear screen")
	fmt.Println("  " + ui.Yellow + "exit, quit" + ui.Reset + "          Exit shell")
	fmt.Println()
	fmt.Println(ui.Bold + "System Commands:" + ui.Reset)
	fmt.Println("  Any other command will be executed as a system command")
	fmt.Println("  Examples: " + ui.Dim + "ls, pwd, whoami, ps aux, netstat -tulpn" + ui.Reset)
	fmt.Println()
	fmt.Println(ui.Bold + "Examples:" + ui.Reset)
	fmt.Println("  " + ui.Green + "use kernel -list" + ui.Reset + "       # List kernel exploits")
	fmt.Println("  " + ui.Green + "run persistence -list" + ui.Reset + "  # List persistence methods")
	fmt.Println("  " + ui.Green + "use persistence 1" + ui.Reset + "      # Run specific persistence technique")
	fmt.Println("  " + ui.Green + "id" + ui.Reset + "                      # Check current user")
	fmt.Println("  " + ui.Green + "ps aux | grep root" + ui.Reset + "     # List root processes")
	return nil
}

func (s *Shell) showHistory() error {
	if len(s.history) == 0 {
		ui.Info("No command history yet")
		return nil
	}
	
	fmt.Println(ui.Bold + ui.Cyan + "Command History:" + ui.Reset)
	for i, cmd := range s.history {
		fmt.Printf("  %s%3d%s  %s\n", ui.Dim, i+1, ui.Reset, cmd)
	}
	return nil
}

func (s *Shell) listModules() error {
	modules := s.registry.List()
	
	if len(modules) == 0 {
		ui.Warn("No modules available")
		return nil
	}
	
	fmt.Println(ui.Bold + ui.Green + "AVAILABLE MODULES:" + ui.Reset)
	fmt.Println()
	
	// Group by category
	categories := make(map[string][]module.Module)
	for _, mod := range modules {
		cat := mod.Category()
		categories[cat] = append(categories[cat], mod)
	}
	
	for cat, mods := range categories {
		fmt.Printf("  %s%s%s\n", ui.Bold+ui.Magenta, strings.ToUpper(cat), ui.Reset)
		fmt.Printf("  %s%s%s\n", ui.Dim, strings.Repeat("─", 78), ui.Reset)
		
		for _, mod := range mods {
			// Show recommendation badge
			recommendation := ""
			if mod.Category() == "persistence" && os.Getuid() == 0 {
				recommendation = ui.Yellow + " [RECOMMENDED]" + ui.Reset
			}
			
			fmt.Printf("  %s%-15s%s  %s%s\n",
				ui.Yellow+ui.Bold, mod.Name(), ui.Reset,
				mod.Description(), recommendation)
		}
		fmt.Println()
	}
	
	return nil
}

func (s *Shell) runModule(name string, args []string) error {
	mod, err := s.registry.Get(name)
	if err != nil {
		return fmt.Errorf("module '%s' not found. Use 'modules' to list available modules", name)
	}
	
	ui.Info("Running module: %s%s%s", ui.Yellow+ui.Bold, name, ui.Reset)
	fmt.Println()
	
	if err := mod.Run(args); err != nil {
		return fmt.Errorf("module execution failed: %v", err)
	}
	
	return nil
}

func (s *Shell) executeSystemCommand(parts []string) error {
	// Handle pipes and redirections by using sh -c
	cmdLine := strings.Join(parts, " ")
	
	// Check if command contains shell operators
	if strings.ContainsAny(cmdLine, "|&><;") {
		cmd := exec.Command("sh", "-c", cmdLine)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}
	
	// Simple command execution
	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return fmt.Errorf("command exited with code %d", exitErr.ExitCode())
		}
		return err
	}
	
	return nil
}
