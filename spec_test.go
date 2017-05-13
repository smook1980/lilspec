package lilspec_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	. "github.com/smook1980/lilspec"
)

func TestSpec(t *testing.T) {
	t.Run("T(t)", func(t *testing.T) {
		var ctrl *gomock.Controller
		var mockT *MockTestContext
		var subject B

		beforeFn := func(t *testing.T) {
			ctrl = gomock.NewController(t)
			mockT = NewMockTestContext(ctrl)
			subject = T(mockT)
		}

		t.Run("Describe creates a subtest", func(t *testing.T) {
			beforeFn(t)
			desc := "describe block description"
			mockT.EXPECT().Run(desc, gomock.Any()).Return(true)
			subject.Describe(desc, func(b B) {})

			ctrl.Finish()
		})

		t.Run("Expect calls Errorf in the correct context", func(t *testing.T) {
			beforeFn(t)
			mockSubT := NewMockTestContext(ctrl)
			mockT.EXPECT().
				Run(gomock.Any(), gomock.Any()).
				Do(func(_ string, subt func(TestContext)) {
					subt(mockSubT)
				})

			mockSubT.EXPECT().Errorf(gomock.Any(), gomock.Any(), gomock.Any())

			subject.It("blah", func(s S) {
				s.Expect(true).To(BeFalse())
			})

			ctrl.Finish()
		})
	})

	T(t).It("works", func(s S) {
		s.Expect(false).To(BeFalse())
	})

	t.Run("s.Describe", func(it *testing.T) {
		T(it).Describe("Describe", func(b B) {
			b.It("has assertions", func(s S) {
				s.Expect(true).To(BeTrue())
				s.Expect(false).ToNot(BeTrue())
			})
		})
	})

	t.Run("s.Context", func(t *testing.T) {
		T(t).Context("when a context block is given", func(b B) {
			b.It("has a nested scope", func(s S) {
				s.Expect(true).To(BeTrue())
			})
		})
	})

	t.Run("spec blocks", func(t *testing.T) {
		T(t).Describe("blocks", func(b B) {
			b.Context("when nested", func(b B) {
				b.It("works", func(s S) {
					s.Expect(true).To(BeTrue())
				})
			})
		})
	})

	t.Run("BeforeEach", func(t *testing.T) {
		x := 1

		T(t).Describe("BeforeEach", func(b B) {
			b.BeforeEach(func(s S) {
				x = 0
			})

			b.BeforeEach(func(s S) {
				x = x + 1
			})

			b.It("runs before each block", func(s S) {
				s.Expect(x).To(Equal(1))
			})

			b.Context("when nested scopes", func(b B) {
				b.BeforeEach(func(s S) {
					x = x + 5
				})

				b.It("runs each BeforeEach", func(s S) {
					s.Expect(x).To(Equal(6))
				})
			})

			b.It("runs before each block", func(s S) {
				s.Expect(x).To(Equal(1))
			})
		})
	})
}
