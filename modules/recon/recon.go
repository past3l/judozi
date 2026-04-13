//go:build linux

package recon

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/judozi/judozi/modules/recon/techniques"
	"github.com/judozi/judozi/pkg/module"
	"github.com/judozi/judozi/pkg/ui"
)

type ReconModule struct{}

func New() module.Module { return &ReconModule{} }

func (r *ReconModule) Name() string        { return "recon" }
func (r *ReconModule) Description() string { return "Advanced System Reconnaissance & Enumeration" }
func (r *ReconModule) Category() string    { return "recon" }

func (r *ReconModule) Run(args []string) error {
	ui.ModuleHeader(r.Name(), r.Category(), r.Description())

	fs := flag.NewFlagSet("recon", flag.ContinueOnError)
	listFlag := fs.Bool("list", false, "List all techniques")
	outDir   := fs.String("dir", "recon", "Directory to save output files")
	noSave   := fs.Bool("nosave", false, "Don't save files, print to stdout")
	if err := fs.Parse(args); err != nil {
		return err
	}

	techs := techniques.GetAll()

	if *listFlag {
		listTechniques(techs)
		return nil
	}

	timestamp := time.Now().Format("20060102_150405")
	remaining := fs.Args()

	if !*noSave {
		if err := os.MkdirAll(*outDir, 0700); err != nil {
			return fmt.Errorf("cannot create output dir %s: %w", *outDir, err)
		}
	}

	if len(remaining) > 0 {
		for _, arg := range remaining {
			idx := 0
			fmt.Sscanf(arg, "%d", &idx)
			if idx < 1 || idx > len(techs) {
				ui.Error("Invalid selection: %s", arg)
				listTechniques(techs)
				return fmt.Errorf("invalid selection")
			}
			runOne(techs[idx-1], *outDir, timestamp, *noSave)
		}
		return nil
	}

	return runParallel(techs, *outDir, timestamp, *noSave)
}

func runParallel(techs []techniques.Technique, outDir, timestamp string, noSave bool) error {
	start := time.Now()

	if !noSave {
		absDir, _ := filepath.Abs(outDir)
		ui.Info("Running %s%d%s techniques in parallel → saving to: %s%s%s",
			ui.Bold, len(techs), ui.Reset, ui.Cyan, absDir, ui.Reset)
	} else {
		ui.Info("Running %s%d%s techniques in parallel...", ui.Bold, len(techs), ui.Reset)
	}
	fmt.Println()

	type techResult struct {
		name    string
		file    string
		content string
	}

	results := make(chan techResult, len(techs))
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, t := range techs {
		wg.Add(1)
		go func(t techniques.Technique) {
			defer wg.Done()
			content := t.Run()

			file := ""
			if !noSave {
				fname := techFileName(t.Name()) + "_" + timestamp + ".txt"
				fpath := filepath.Join(outDir, fname)
				header := fmt.Sprintf("# judozi recon — %s\n# Generated: %s\n\n", t.Name(), timestamp)
				_ = os.WriteFile(fpath, []byte(header+content), 0600)
				file = fpath
			}

			mu.Lock()
			if !noSave {
				fmt.Printf("  %s[✓]%s %-34s → %s%s%s\n",
					ui.Green, ui.Reset, t.Name(), ui.Cyan, file, ui.Reset)
			} else {
				fmt.Printf("\n%s══ %s ══%s\n%s", ui.Cyan+ui.Bold, t.Name(), ui.Reset, content)
			}
			mu.Unlock()

			results <- techResult{name: t.Name(), file: file, content: content}
		}(t)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for range results {
	}

	elapsed := time.Since(start)
	fmt.Println()
	ui.Success("Completed %d techniques in %s%.1fs%s", len(techs), ui.Bold, elapsed.Seconds(), ui.Reset)
	if !noSave {
		absDir, _ := filepath.Abs(outDir)
		fmt.Printf("  %sFiles saved in: %s%s\n\n", ui.Dim, absDir, ui.Reset)
	}
	return nil
}

func runOne(t techniques.Technique, outDir, timestamp string, noSave bool) {
	start := time.Now()
	ui.Info("Running: %s%s%s", ui.Yellow+ui.Bold, t.Name(), ui.Reset)
	content := t.Run()
	elapsed := time.Since(start)

	if noSave {
		fmt.Println(content)
	} else {
		fname := techFileName(t.Name()) + "_" + timestamp + ".txt"
		fpath := filepath.Join(outDir, fname)
		header := fmt.Sprintf("# judozi recon — %s\n# Generated: %s\n\n", t.Name(), timestamp)
		_ = os.WriteFile(fpath, []byte(header+content), 0600)
		ui.Success("Saved → %s%s%s (%.1fs)", ui.Cyan, fpath, ui.Reset, elapsed.Seconds())
	}
}

func listTechniques(techs []techniques.Technique) {
	fmt.Println()
	fmt.Printf("  %s%-4s %-34s %-10s %s%s\n",
		ui.Bold+ui.Cyan, "ID", "TECHNIQUE", "STEALTH", "DESCRIPTION", ui.Reset)
	fmt.Printf("  %s%s%s\n", ui.Dim, strings.Repeat("─", 95), ui.Reset)

	for i, t := range techs {
		stealth := stealthColor(t.StealthLevel())
		fmt.Printf("  %s%-4d%s %-34s %-20s %s\n",
			ui.Yellow, i+1, ui.Reset, t.Name(), stealth, t.Description())
	}
	fmt.Println()
	fmt.Println(ui.Dim + "  Use: recon [id...]     → run specific technique(s)" + ui.Reset)
	fmt.Println(ui.Dim + "       recon             → run all in parallel" + ui.Reset)
	fmt.Println(ui.Dim + "       recon -dir /tmp   → save files to /tmp" + ui.Reset)
	fmt.Println(ui.Dim + "       recon -nosave      → print to stdout only" + ui.Reset)
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

// techFileName converts technique name to a safe filename prefix
func techFileName(name string) string {
	r := strings.NewReplacer(
		" ", "_", "/", "_", "&", "", "-", "_",
		"(", "", ")", "", ",", "", ".", "",
	)
	result := r.Replace(strings.ToLower(name))
	for strings.Contains(result, "__") {
		result = strings.ReplaceAll(result, "__", "_")
	}
	return strings.Trim(result, "_")
}
