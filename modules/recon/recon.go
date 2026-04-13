//go:build linux

package recon

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/judozi/judozi/modules/recon/techniques"
	"github.com/judozi/judozi/pkg/module"
	"github.com/judozi/judozi/pkg/ui"
)

type ReconModule struct{}

func New() module.Module {
	return &ReconModule{}
}

func (r *ReconModule) Name() string        { return "recon" }
func (r *ReconModule) Description() string { return "Advanced System Reconnaissance & Enumeration" }
func (r *ReconModule) Category() string    { return "recon" }

func (r *ReconModule) Run(args []string) error {
	ui.ModuleHeader(r.Name(), r.Category(), r.Description())

	fs := flag.NewFlagSet("recon", flag.ExitOnError)
	list := fs.Bool("list", false, "List all recon techniques")
	all  := fs.Bool("all", false, "Run all recon techniques")
	out  := fs.String("out", "", "Save output to file")
	if err := fs.Parse(args); err != nil {
		return err
	}

	techs := techniques.GetAll()

	if *list {
		listTechniques(techs)
		return nil
	}

	var output *os.File
	if *out != "" {
		f, err := os.Create(*out)
		if err != nil {
			return fmt.Errorf("cannot create output file: %w", err)
		}
		defer f.Close()
		output = f
		ui.Success("Output will also be saved to: %s%s%s", ui.Green, *out, ui.Reset)
		fmt.Println()
	}

	remaining := fs.Args()

	if *all || len(remaining) == 0 {
		return runAll(techs, output)
	}

	// Run specific technique by number
	if len(remaining) > 0 {
		for _, arg := range remaining {
			idx := 0
			fmt.Sscanf(arg, "%d", &idx)
			if idx < 1 || idx > len(techs) {
				ui.Error("Invalid technique number: %s", arg)
				listTechniques(techs)
				return fmt.Errorf("invalid selection")
			}
			runTechnique(techs[idx-1], output)
		}
		return nil
	}

	return runInteractive(techs, output)
}

func listTechniques(techs []techniques.Technique) {
	fmt.Println()
	fmt.Printf("  %s%-4s %-30s %-10s %s%s\n",
		ui.Bold+ui.Cyan, "ID", "TECHNIQUE", "STEALTH", "DESCRIPTION", ui.Reset)
	fmt.Printf("  %s%s%s\n", ui.Dim, strings.Repeat("─", 95), ui.Reset)

	for i, t := range techs {
		stealth := stealthColor(t.StealthLevel())
		fmt.Printf("  %s%-4d%s %-30s %-20s %s\n",
			ui.Yellow, i+1, ui.Reset,
			t.Name(), stealth, t.Description())
	}
	fmt.Println()
	fmt.Println(ui.Dim + "  Use: judozi recon <id> [id...] | judozi recon -all | judozi recon -out report.txt" + ui.Reset)
	fmt.Println()
}

func stealthColor(level string) string {
	switch strings.ToLower(level) {
	case "passive":
		return ui.Green + "PASSIVE" + ui.Reset
	case "low":
		return ui.Yellow + "LOW    " + ui.Reset
	case "medium":
		return ui.Red + "MEDIUM " + ui.Reset
	default:
		return ui.Dim + level + ui.Reset
	}
}

func runAll(techs []techniques.Technique, out *os.File) error {
	ui.Info("Running %d recon techniques...", len(techs))
	fmt.Println()
	for _, t := range techs {
		runTechnique(t, out)
	}
	return nil
}

func runTechnique(t techniques.Technique, out *os.File) {
	header := fmt.Sprintf("\n%s══════════════════════════════════════════════════════════════════════════════════%s\n  %s%s%s  |  %s\n%s══════════════════════════════════════════════════════════════════════════════════%s\n",
		ui.Cyan+ui.Bold, ui.Reset,
		ui.Yellow+ui.Bold, t.Name(), ui.Reset, t.Description(),
		ui.Cyan+ui.Bold, ui.Reset)

	fmt.Print(header)
	if out != nil {
		fmt.Fprintf(out, "\n## %s\n%s\n\n", t.Name(), t.Description())
	}

	result := t.Run()
	fmt.Println(result)
	if out != nil {
		fmt.Fprintln(out, result)
	}
}

func runInteractive(techs []techniques.Technique, out *os.File) error {
	listTechniques(techs)

	fmt.Println(ui.Dim + "  ╔════════════════════════════════════════════════════════════════════╗")
	fmt.Println("  ║ " + ui.Reset + "Enter number(s), 'all', or 'q' to quit" + ui.Dim + "                          ║")
	fmt.Println("  ╚════════════════════════════════════════════════════════════════════╝" + ui.Reset)
	fmt.Print("  > ")

	var input string
	fmt.Scanln(&input)

	if input == "q" || input == "quit" {
		return nil
	}
	if input == "all" {
		return runAll(techs, out)
	}

	idx := 0
	fmt.Sscanf(input, "%d", &idx)
	if idx < 1 || idx > len(techs) {
		ui.Error("Invalid selection")
		return nil
	}

	runTechnique(techs[idx-1], out)
	return nil
}
