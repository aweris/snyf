package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/samber/lo"
	"golang.org/x/mod/modfile"
	"golang.org/x/mod/module"
)

// GoModule represents a module with dependencies
type GoModule struct {
	Path      string           `json:"path"`
	GoVersion string           `json:"go_version"`
	Module    string           `json:"module"`
	Require   []module.Version `json:"require,omitempty"`
	Replace   []Replace        `json:"replace,omitempty"`
	DependsOn []string         `json:"depends_on,omitempty"`
	UsedBy    []string         `json:"used_by,omitempty"`
	Gomod     *modfile.File    `json:"-"`
}

// A Replace is a single replace statement.
type Replace struct {
	Old module.Version
	New module.Version
}

// ModuleGraph stores module relationships
type ModuleGraph struct {
	Modules map[string]*GoModule `json:"modules"`
}

func main() {
	// flags
	var (
		source   string
		extended bool
	)

	flag.StringVar(&source, "source", ".", "source directory to scan for go.mod files")
	flag.BoolVar(&extended, "extended", false, "include require and replace information for each module. !!WARNING!! This can be very verbose.")

	flag.Parse()

	modules := make(map[string]*GoModule)

	// Scan for go.mod files
	err := filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Name() == "go.mod" {
			mod, err := parseGoMod(path)
			if err != nil {
				fmt.Printf("Error processing %s: %v\n", path, err)
				return nil
			}

			relPath, err := filepath.Rel(source, path)
			if err != nil {
				fmt.Printf("Error converting path %s to relative: %v\n", path, err)
				relPath = path // Fallback to absolute path if conversion fails
			}

			// Store module information
			gomod := &GoModule{
				Path:      filepath.Dir(relPath),
				GoVersion: mod.Go.Version,
				Module:    mod.Module.Mod.Path,
				DependsOn: []string{},
				UsedBy:    []string{},
				Gomod:     mod,
			}

			if extended {
				gomod.Require = lo.Map(mod.Require, func(r *modfile.Require, _ int) module.Version { return r.Mod })
				gomod.Replace = lo.Map(mod.Replace, func(r *modfile.Replace, _ int) Replace { return Replace{Old: r.Old, New: r.New} })
			}

			modules[mod.Module.Mod.Path] = gomod
		}

		return nil
	})

	if err != nil {
		fmt.Println("Error scanning directories:", err)
		return
	}

	// find depends_on and used_by relationships
	for _, modData := range modules {
		for _, req := range modData.Gomod.Require {
			// Only add dependencies if they exist in our repo
			if dep, exists := modules[req.Mod.Path]; exists {
				modData.DependsOn = append(modData.DependsOn, req.Mod.Path)
				dep.UsedBy = append(dep.UsedBy, modData.Module)
			}
		}
	}

	// Output JSON
	graph := ModuleGraph{Modules: modules}
	output, err := json.MarshalIndent(graph, "", "  ")
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	fmt.Println(string(output))
}

// returns parsed go.mod file
func parseGoMod(path string) (*modfile.File, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return modfile.Parse("go.mod", []byte(content), nil)
}
