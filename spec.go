package lilspec

//go:generate mockgen -package mocks -destination ./mocks/fs_mock.go github.com/spf13/afero Fs
//go:generate mockgen -package mocks -destination ./mocks/file_mock.go github.com/spf13/afero File

import (
	"sync"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/onsi/gomega"
)

type TestContext interface {
	Errorf(format string, args ...interface{})
	Run(string, func(*testing.T)) bool
	Fatalf(format string, args ...interface{})
}

// S provides BDD Spec features
type S struct {
	bfns []func(s *S)
	TestContext
	once     sync.Once
	mockCtrl *gomock.Controller
}

// T returns a base spec given t
func T(t TestContext) *S {
	return &S{TestContext: t}
}

// Expect allows for calling assertions on the subject s
func (s *S) Expect(it interface{}, e ...interface{}) gomega.GomegaAssertion {
	gomega.RegisterTestingT(s.TestContext)
	return gomega.Expect(it, e...)
}

func (s *S) Mock() *gomock.Controller {
	s.once.Do(func() {
		s.mockCtrl = gomock.NewController(s.TestContext)
	})

	return s.mockCtrl
}

func (s *S) specBlock(desc string, descb func(*S)) {
	for _, bfn := range s.bfns {
		bfn(s)
	}
	s.TestContext.Run(desc, func(t *testing.T) {
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

// BeforeEach specifies a block to run before each spec block
func (s *S) BeforeEach(bfn func(s *S)) {
	s.bfns = append(s.bfns, bfn)
}
