package scaffold

import (
	"context"
	"embed"
	"fmt"
	"path/filepath"

	"dario.cat/mergo"
	"github.com/nimbolus/terraform-backend/pkg/fs"
	"gopkg.in/yaml.v3"
)

var (
	//go:embed files
	assets       embed.FS
	filesToWrite = []string{
		"tf-preview.yaml",
		"tf-run.yaml",
	}
)

func writeGithubActionsWorkflows(ctx context.Context, dir fs.FS) error {
	if err := dir.MkdirAll(".github/workflows", 0755); err != nil {
		return err
	}

	for _, filename := range filesToWrite {
		outFilename := filepath.Join(".github", "workflows", filename)

		_, err := dir.Stat(outFilename)
		fileExists := err == nil
		if fileExists {
			ok, err := promptYesNo(ctx, fmt.Sprintf("Workflow at %s already exist. Do you want to replace it? (This is experimental and might not deal well with your edits.)", outFilename))
			if err != nil {
				return err
			}
			if !ok {
				fmt.Printf("Skipping update of %s\n", outFilename)
				continue
			}
		}

		srcFile, err := assets.Open(filepath.Join("files", filename))
		if err != nil {
			return err
		}
		defer srcFile.Close()

		var config yaml.Node
		if err := yaml.NewDecoder(srcFile).Decode(&config); err != nil {
			return err
		}

		if fileExists {
			oldFile, err := dir.Open(outFilename)
			if err != nil {
				return err
			}
			defer oldFile.Close()

			var oldConfig yaml.Node
			if err := yaml.NewDecoder(oldFile).Decode(&oldConfig); err != nil {
				return err
			}

			if err := mergo.Merge(&oldConfig, &config, mergo.WithSliceDeepCopy); err != nil {
				return err
			}
			config = oldConfig
		}

		f, err := fs.Create(dir, outFilename)
		if err != nil {
			return err
		}
		defer f.Close()

		enc := yaml.NewEncoder(f)
		enc.SetIndent(2)
		if err := enc.Encode(&config); err != nil {
			return err
		}
		if err := enc.Close(); err != nil {
			return err
		}
		fmt.Printf("Wrote workflow to: %s\n", outFilename)
	}
	return nil
}
