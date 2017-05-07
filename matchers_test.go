package lilspec_test

import (
	"testing"

	. "github.com/smook1980/lilspec"
	"github.com/spf13/afero"
)

func TestMatchers(t *testing.T) {
	T(t).Describe("BeARegularFile", func(s *S) {
		fs := afero.NewOsFs()
		s.Expect("matchers_test.go").To(BeAFile(fs))
	})
}
