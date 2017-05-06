package lilspec

import (
	"sync"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/onsi/gomega"
)

type TCtx interface {
	Errorf(format string, args ...interface{})
	Run(string, func(*testing.T)) bool
	Fatalf(format string, args ...interface{})
}

// S provides BDD Spec features
type S struct {
	bfns []func(s *S)
	TCtx
	once     sync.Once
	mockCtrl *gomock.Controller
}

// T returns a base spec given t
func T(t TCtx) *S {
	return &S{TCtx: t}
}

// Expect allows for calling assertions on the subject s
func (s *S) Expect(it interface{}, e ...interface{}) gomega.GomegaAssertion {
	gomega.RegisterTestingT(s.TCtx)
	return gomega.Expect(it, e...)
}

func (s *S) mock() *gomock.Controller {
	s.once.Do(func() {
		s.mockCtrl = gomock.NewController(s.TCtx)
	})

	return s.mockCtrl
}

func (s *S) specBlock(desc string, descb func(*S)) {
	for _, bfn := range s.bfns {
		bfn(s)
	}
	s.TCtx.Run(desc, func(t *testing.T) {
		descb(T(t))
	})
	if s.mockCtrl != nil {
		s.mockCtrl.Finish()
	}
}

// Describe begins a sub test with the given description.
func (s *S) Describe(desc string, block func(*S)) {
	s.specBlock(desc, block)
}

// Context begins a sub test with the given description.
func (s *S) Context(desc string, block func(*S)) {
	s.specBlock(desc, block)
}

// It begins a sub test with the given description.
func (s *S) It(desc string, block func(*S)) {
	s.specBlock(desc, block)
}

func (s *S) BeforeEach(bfn func(s *S)) {
	s.bfns = append(s.bfns, bfn)
}
