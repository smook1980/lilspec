package lilspec

import (
	"testing"

	"github.com/onsi/gomega"
)

type TCtx interface {
	Errorf(format string, args ...interface{})
	Run(string, func(*testing.T)) bool
}

// S provides BDD Spec features
type S struct {
	TCtx
}

// S returns a base spec given t
func T(t TCtx) *S {
	return &S{TCtx: t}
}

// Expect allows for calling assertions on the subject s
func (s *S) Expect(it interface{}, e ...interface{}) gomega.GomegaAssertion {
	gomega.RegisterTestingT(s.TCtx)
	return gomega.Expect(it, e...)
}

// Describe begins a sub test with the given description.
func (s *S) Describe(desc string, descb func(*S)) {
	s.TCtx.Run(desc, func(t *testing.T) {
		descb(T(t))
	})
}

// Describe begins a sub test with the given description.
func (s *S) When(desc string, descb func(*S)) {
	s.TCtx.Run(desc, func(t *testing.T) {
		descb(T(t))
	})
}
