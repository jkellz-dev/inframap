package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/spf13/afero"

	"github.com/cycloidio/inframap/errcode"
)

// resourceBody is the subset of a Terraform resource that inframap needs.
type resourceBody struct {
	Config   hcl.Body
	DefRange hcl.Range
}

// resourceSet replaces terraform's configs.Module. inframap only ever read
// ManagedResources from it, so the old full-language loader (pinned to TF 0.15)
// is unnecessary and rejected modern blocks (import, removed, check, preconditions).
type resourceSet struct {
	ManagedResources map[string]*resourceBody
}

// topLevelSchema only names "resource"; PartialContent ignores every other block
// type, so new Terraform language features never cause a parse failure.
var topLevelSchema = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{{Type: "resource", LabelNames: []string{"type", "name"}}},
}

func loadResources(fs afero.Fs, path string) (*resourceSet, error) {
	files, err := tfFiles(fs, path)
	if err != nil {
		return nil, err
	}
	parser := hclparse.NewParser()
	rs := &resourceSet{ManagedResources: map[string]*resourceBody{}}
	for _, f := range files {
		src, err := afero.ReadFile(fs, f)
		if err != nil {
			return nil, fmt.Errorf("reading %s: %w", f, err)
		}
		file, diags := parser.ParseHCL(src, f)
		if diags.HasErrors() {
			return nil, fmt.Errorf("parsing %s: %s", f, diags.Error())
		}
		content, _, diags := file.Body.PartialContent(topLevelSchema)
		if diags.HasErrors() {
			return nil, fmt.Errorf("decoding %s: %s", f, diags.Error())
		}
		for _, b := range content.Blocks {
			key := b.Labels[0] + "." + b.Labels[1]
			if prev, ok := rs.ManagedResources[key]; ok {
				return nil, fmt.Errorf("%s: duplicate resource %q, already declared at %s", b.DefRange, key, prev.DefRange)
			}
			rs.ManagedResources[key] = &resourceBody{Config: b.Body, DefRange: b.DefRange}
		}
	}
	return rs, nil
}

// tfFiles returns the primary .tf files for path. Override files are skipped:
// Terraform merges them into the base resource, which already carries the
// references inframap needs, and loading them separately would look like a
// duplicate. JSON configuration is rejected as before.
func tfFiles(fs afero.Fs, path string) ([]string, error) {
	info, err := fs.Stat(path)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		if strings.HasSuffix(path, ".json") {
			return nil, errcode.ErrGenerateFromJSON
		}
		return []string{path}, nil
	}
	entries, err := afero.ReadDir(fs, path)
	if err != nil {
		return nil, err
	}
	var out []string
	for _, e := range entries {
		name := e.Name()
		switch {
		case e.IsDir():
		case strings.HasSuffix(name, ".tf.json"):
			return nil, errcode.ErrGenerateFromJSON
		case !strings.HasSuffix(name, ".tf"), isOverrideFile(name):
		default:
			out = append(out, filepath.Join(path, name))
		}
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("no .tf files in %s: %w", path, os.ErrNotExist)
	}
	return out, nil
}

func isOverrideFile(name string) bool {
	base := strings.TrimSuffix(name, ".tf")
	return base == "override" || strings.HasSuffix(base, "_override")
}
