package js

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dop251/goja"
)

// CompiledHelpers holds pre-compiled global helper code.
type CompiledHelpers struct {
	program *goja.Program
}

// CompileHelpers compiles global helper code from inline JS and files.
// Returns nil if no helpers are configured.
func CompileHelpers(jsCode string, jsFiles []string, configDir string) (*CompiledHelpers, error) {
	var parts []string

	if jsCode != "" {
		parts = append(parts, jsCode)
	}

	for _, file := range jsFiles {
		path := file
		if !filepath.IsAbs(path) {
			path = filepath.Join(configDir, path)
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to read helper file %s: %w", file, err)
		}
		parts = append(parts, string(content))
	}

	if len(parts) == 0 {
		return nil, nil
	}

	combined := strings.Join(parts, "\n;\n")
	program, err := goja.Compile("helpers", combined, true)
	if err != nil {
		return nil, fmt.Errorf("failed to compile helpers: %w", err)
	}

	return &CompiledHelpers{program: program}, nil
}

// InjectInto executes helpers in the VM, adding them to globalThis.
func (h *CompiledHelpers) InjectInto(vm *goja.Runtime) error {
	if h == nil || h.program == nil {
		return nil
	}
	_, err := vm.RunProgram(h.program)
	return err
}
