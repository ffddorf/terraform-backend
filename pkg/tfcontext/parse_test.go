package tfcontext_test

import (
	"os"
	"testing"

	"github.com/nimbolus/terraform-backend/pkg/tfcontext"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFindBackend(t *testing.T) {
	dir := os.DirFS("./testdata")
	be, err := tfcontext.FindBackend(dir)
	require.NoError(t, err)

	assert.Equal(t, "https://dummy-backend.example.com/state", be.Address)
	assert.Equal(t, "my_user", be.Username)
	assert.Empty(t, be.Password)
}
