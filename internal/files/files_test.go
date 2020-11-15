package files_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/goreleaser/nfpm"
	"github.com/goreleaser/nfpm/internal/files"
)

func TestListFilesToCopy(t *testing.T) {
	info := &nfpm.Info{
		Overridables: nfpm.Overridables{
			Files: []*files.FileToCopy{
				{"../../testdata/scripts/*", "/test","", 0},
				{"../../testdata/whatever.conf", "/whatever","config", 0},
			},
		},
	}

	regularFiles, err := files.ExpandFiles(info.Files, info.DisableGlobbing)
	require.NoError(t, err)

	// all the input files described in the config in sorted order by source path
	require.Equal(t, []*files.FileToCopy{
		{"../../testdata/scripts/postinstall.sh", "/test/postinstall.sh", "", 0},
		{"../../testdata/scripts/postremove.sh", "/test/postremove.sh", "", 0},
		{"../../testdata/scripts/preinstall.sh", "/test/preinstall.sh", "", 0},
		{"../../testdata/scripts/preremove.sh", "/test/preremove.sh", "", 0},
	}, regularFiles)

	require.Equal(t, []*files.FileToCopy{
		{"../../testdata/scripts/*", "/test","", 0},
		{"../../testdata/whatever.conf", "/whatever","config", 0},
	}, regularFiles)
}

func TestListFilesToCopyWithAndWithoutGlobbing(t *testing.T) {
	_, err := files.Expand(map[string]string{
		"../../testdata/{file}*": "/test/{file}[",
	}, false)
	assert.EqualError(t, err, "glob failed: ../../testdata/{file}*: no matching files")

	mappedFiles, err := files.Expand(map[string]string{
		"../../testdata/{file}[": "/test/{file}[",
	}, true)
	require.NoError(t, err)
	assert.Equal(t, []*files.FileToCopy{
		{"../../testdata/{file}[", "/test/{file}[", "", 0},
	}, mappedFiles)
}
