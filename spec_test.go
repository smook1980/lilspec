package lilspec_test

import (
	"testing"

	. "github.com/onsi/gomega"
	. "github.com/smook1980/lilspec"
)

func TestSpec(t *testing.T) {
	T(t).Expect(false).To(BeFalse())

	t.Run("s.Describe", func(it *testing.T) {
		T(it).Describe("Myself", func(s *S) {
			if s == nil {
				t.Error("Expected block to have been passed *S, got nil.")
			}

			s.Expect(true).To(BeTrue())
			s.Expect(false).ToNot(BeTrue())
		})
	})

	t.Run("s.When", func(it *testing.T) {
		T(it).When("when blocks", func(s *S) {
			s.Expect(true).To(BeTrue())
			s.When("nested", func(s *S) {
				s.Expect(true).To(BeTrue())
			})
		})
	})
}
