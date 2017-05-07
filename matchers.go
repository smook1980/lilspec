package lilspec

import (
	"github.com/onsi/gomega/types"
	"github.com/smook1980/lilspec/matchers"
	"github.com/spf13/afero"
)

// BeAFile succeeds iff a file exists and is a regular file.
// Actual must be a string representing the abs path to the file being checked.
func BeAFile(fs afero.Fs) types.GomegaMatcher {
	return &matchers.BeARegularFileMatcher{FS: fs}
}
