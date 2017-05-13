package lilspec_test

import (
	"testing"

	. "github.com/smook1980/lilspec"
	"github.com/spf13/afero"
)

func TestMatchers(t *testing.T) {
	T(t).Describe("BeARegularFile", func(b B) {
		fs := afero.NewOsFs()
		b.It("matches an existing file", func(s S) {
			s.Expect("matchers_test.go").To(BeAFile(fs))
		})
	})
}
