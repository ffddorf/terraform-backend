package scaffold_test

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/liamg/memoryfs"
	"github.com/nimbolus/terraform-backend/pkg/scaffold"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const expectedBackendConfig = `terraform {
  backend "http" {
    address        = "https://ffddorf-terraform-backend.fly.dev/state/terraform-backend/default"
    lock_address   = "https://ffddorf-terraform-backend.fly.dev/state/terraform-backend/default"
    unlock_address = "https://ffddorf-terraform-backend.fly.dev/state/terraform-backend/default"
    username       = "github_pat"
  }
}
`

func compareFiles(t *testing.T, fs1, fs2 fs.FS, name string) {
	file1, err := fs.ReadFile(fs1, name)
	require.NoError(t, err)

	file2, err := fs.ReadFile(fs2, name)
	require.NoError(t, err)

	assert.Equal(t, string(file1), string(file2), "in file: %s", name)
}

type confirmer struct{}

var confirmation = []byte{'y', '\n'}

func (c *confirmer) Read(dst []byte) (int, error) {
	return copy(dst, confirmation), nil
}

func TestScaffolding(t *testing.T) {
	nativeFS := os.DirFS("files")

	tests := map[string]struct {
		stdin  io.Reader
		assert func(*testing.T, *memoryfs.FS)
	}{
		"empty": {
			assert: func(t *testing.T, memfs *memoryfs.FS) {
				backendOut, err := memfs.ReadFile("backend.tf")
				require.NoError(t, err)
				assert.Equal(t, expectedBackendConfig, string(backendOut))

				subFS, err := memfs.Sub(".github/workflows")
				require.NoError(t, err)
				compareFiles(t, nativeFS, subFS, "tf-preview.yaml")
				compareFiles(t, nativeFS, subFS, "tf-run.yaml")
			},
		},
		"update": {
			stdin: &confirmer{},
			assert: func(t *testing.T, memfs *memoryfs.FS) {
				expectedFS, err := memfs.Sub("expected")
				require.NoError(t, err)
				compareFiles(t, expectedFS, memfs, "backend.tf")
				compareFiles(t, expectedFS, memfs, ".github/workflows/tf-run.yaml")

				subFS, err := memfs.Sub(".github/workflows")
				require.NoError(t, err)
				compareFiles(t, nativeFS, subFS, "tf-preview.yaml")
			},
		},
	}

	for name, params := range tests {
		t.Run(name, func(t *testing.T) {
			memfs := memoryfs.New()

			s, err := os.Stat("testdata/" + name)
			if err == nil {
				require.True(t, s.IsDir(), "testdata for test name needs to be a directory")
				testFS := os.DirFS("testdata/" + name)
				err := fs.WalkDir(testFS, ".", func(path string, d fs.DirEntry, err error) error {
					if err != nil {
						return err
					}
					if !d.Type().IsRegular() {
						return nil
					}

					contents, err := fs.ReadFile(testFS, path)
					if err != nil {
						return err
					}
					if err := memfs.MkdirAll(filepath.Dir(path), 0755); err != nil {
						return err
					}
					return memfs.WriteFile(path, contents, d.Type().Perm())
				})
				require.NoError(t, err)
			}

			var stdin io.Reader = os.Stdin
			if params.stdin != nil {
				stdin = params.stdin
			}

			cmd := scaffold.NewCommand(memfs, stdin)
			require.NoError(t, cmd.Execute())

			params.assert(t, memfs)
		})
	}
}
