package lilspec

//go:generate mockgen -package mocks -destination ./mocks/fs_mock.go github.com/spf13/afero Fs
//go:generate mockgen -package mocks -destination ./mocks/file_mock.go github.com/spf13/afero File
//go:generate mockgen -package lilspec_test -destination ./mock_test_context_test.go github.com/smook1980/lilspec TestContext

import (
	"fmt"
	"sync"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/onsi/gomega"
)

// TestContext test runner interface
type TestContext interface {
	Errorf(format string, args ...interface{})
	Run(string, func(TestContext)) bool
	Fatalf(format string, args ...interface{})
}

// B is a block providing context for S
type B interface {
	Describe(string, func(B))
	Context(string, func(B))
	BeforeEach(func(S))
	It(string, func(S))
}

// S is a spec
type S interface {
	Expect(it interface{}, e ...interface{}) gomega.GomegaAssertion
	Mock() *gomock.Controller
}

type specContext struct {
	bfns []func(s S)
	TestContext
	parent   *specContext
	once     sync.Once
	mockCtrl *gomock.Controller
}

type tAdapter struct {
	*testing.T
}

func (ta *tAdapter) Run(desc string, tfn func(TestContext)) bool {
	atfn := func(t *testing.T) {
		tfn(&tAdapter{T: t})
	}

	return ta.T.Run(desc, atfn)
}

// T returns a base spec given t
func T(t interface{}) B {
	var tctx TestContext

	switch ctx := t.(type) {
	case TestContext:
		tctx = ctx
	case *testing.T:
		tctx = &tAdapter{T: ctx}
	default:
		panic(fmt.Errorf("T(t) called when t was not *testing.T or TestContext.  Was %+v.", t))
	}

	return &specContext{TestContext: tctx}
}

func (s *specContext) subTest(t TestContext) *specContext {
	return &specContext{TestContext: t, parent: s}
}

// Expect allows for calling assertions on the subject s
func (s *specContext) Expect(it interface{}, e ...interface{}) gomega.GomegaAssertion {
	return gomega.ExpectWithOffset(1, it, e...)
}

// Mock returns the gomock Controller for the current scope
func (s *specContext) Mock() *gomock.Controller {
	s.once.Do(func() {
		s.mockCtrl = gomock.NewController(s)
	})

	return s.mockCtrl
}

func (s *specContext) assertMocks() {
	if s.mockCtrl == nil {
		return
	}

	ctrl := s.mockCtrl
	s.mockCtrl = nil
	s.once = sync.Once{}
	ctrl.Finish()
}

func (s *specContext) specBlock(desc string, descb func(B)) {
	s.TestContext.Run(desc, func(t TestContext) {
		descb(s.subTest(t))
	})
}

func (s *specContext) spec(desc string, spec func(S)) {
	s.Run(desc, func(t TestContext) {
		s = s.subTest(t)
		gomega.RegisterTestingT(s)
		for _, bfn := range s.beforeFns() {
			bfn(s)
		}
		spec(s)
		s.assertMocks()
	})
}

func (s *specContext) beforeFns() []func(S) {
	if s.parent != nil {
		return append(s.parent.beforeFns(), s.bfns...)
	}

	return s.bfns
}

// Describe begins a sub test with the given description.
func (s *specContext) Describe(desc string, block func(B)) {
	s.specBlock(desc, block)
}

// Context begins a sub test with the given description.
func (s *specContext) Context(desc string, block func(B)) {
	s.specBlock(desc, block)
}

// It begins a sub test with the given description.
func (s *specContext) It(desc string, block func(S)) {
	s.spec(desc, block)
}

// BeforeEach specifies a block to run before each spec block
func (s *specContext) BeforeEach(bfn func(S)) {
	s.bfns = append(s.bfns, bfn)
}
