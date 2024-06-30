package scaffold

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/nimbolus/terraform-backend/pkg/git"
	"github.com/nimbolus/terraform-backend/pkg/tfcontext"
	"github.com/zclconf/go-cty/cty"
)

func writeBackendConfig(ctx context.Context, dir string) (reterr error) {
	var file *hclwrite.File
	var outFile io.WriteCloser
	var backendBlock *hclwrite.Block

	_, filename, err := tfcontext.FindBackendBlock(dir)
	if err == nil {
		relPath, _ := filepath.Rel(dir, filename)
		ok, err := promptYesNo(ctx, fmt.Sprintf("There is an existing backend config at %s. Do you want to update it?", relPath))
		if err != nil {
			return err
		}
		if !ok {
			fmt.Println("Skipping backend update")
			return nil
		}

		b, err := os.ReadFile(filename)
		if err != nil {
			return err
		}

		var diags hcl.Diagnostics
		file, diags = hclwrite.ParseConfig(b, filename, hcl.Pos{})
		if len(diags) > 0 {
			return errors.Join(diags)
		}
		var tfBlock *hclwrite.Block
		for _, block := range file.Body().Blocks() {
			if block.Type() != "terraform" {
				continue
			}
			tfBlock = block
			for _, innerBlock := range block.Body().Blocks() {
				if innerBlock.Type() == "backend" {
					backendBlock = innerBlock
				}
			}
		}
		if backendBlock == nil {
			return errors.New("backend block not found anymore")
		}
		if backendBlock.Labels()[0] != "http" {
			tfBlock.Body().RemoveBlock(backendBlock)
			backendBlock = tfBlock.Body().AppendNewBlock("backend", nil)
		}

		outFile, err = os.Create(filename)
		if err != nil {
			return err
		}
		defer func() {
			if reterr != nil {
				// restore original content
				_, _ = outFile.Write(b)
			}
			_ = outFile.Close()
		}()
	} else {
		file = hclwrite.NewEmptyFile()
		tfBlock := file.Body().AppendNewBlock("terraform", nil)
		backendBlock = tfBlock.Body().AppendNewBlock("backend", nil)
		filename = filepath.Join(dir, "backend.tf")
		outFile, err = os.Create(filename)
		if err != nil {
			return err
		}
		defer outFile.Close()
	}

	origin, err := git.RepoOrigin()
	if err != nil {
		return fmt.Errorf("failed to read repo origin: %w", err)
	}
	segments := strings.Split(origin.Path, "/")
	if len(segments) < 2 {
		return fmt.Errorf("invalid repo path: %s", origin.Path)
	}
	repo := segments[1]

	backendURL, err := url.Parse(backendAddress)
	if err != nil {
		return err
	}
	backendURL.Path = filepath.Join(backendURL.Path, "state", repo, "default")
	address := backendURL.String()

	backendBlock.SetLabels([]string{"http"})
	backendBody := backendBlock.Body()
	backendAttributes := []string{"address", "lock_address", "unlock_address", "username"}
	for name := range backendBody.Attributes() {
		if slices.Contains(backendAttributes, name) {
			continue
		}
		backendBody.RemoveAttribute(name)
	}
	backendBody.SetAttributeValue("address", cty.StringVal(address))
	backendBody.SetAttributeValue("lock_address", cty.StringVal(address))
	backendBody.SetAttributeValue("unlock_address", cty.StringVal(address))
	backendBody.SetAttributeValue("username", cty.StringVal("github_pat"))

	if _, err := file.WriteTo(outFile); err != nil {
		return err
	}

	relPath, _ := filepath.Rel(dir, filename)
	fmt.Printf("Wrote backend config to: %s\n", relPath)
	return nil
}
