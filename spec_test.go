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

	t.Run("spec blocks", func(t *testing.T) {
		T(t).Context("when a context block is given", func(s *S) {
			s.Expect(true).To(BeTrue())
			s.It("has a nested scope", func(s *S) {
				s.Expect(true).To(BeTrue())
			})
		})
	})

	t.Run("BeforeEach", func(t *testing.T) {
		T(t).Describe("BeforeEach", func(s *S) {
			x := 1
			s.BeforeEach(func(s *S) {
				x = x + 1
			})

			s.It("runs before each block", func(s *S) {
				s.Expect(x).To(Equal(2))
				y := x

				s.BeforeEach(func(s *S) {
					y = y + 1
				})

				s.It("test block", func(s *S) {
					s.Expect(y).To(Equal(3))
				})
			})

			s.It("runs before each block", func(s *S) {
				s.Expect(x).To(Equal(3))
			})
		})
	})
}
