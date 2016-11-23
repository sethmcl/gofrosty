package test

import (
	"github.com/sethmcl/gofrosty/lib/gfjs"
	"github.com/sethmcl/gofrosty/lib/module"
	"github.com/sethmcl/gofrosty/lib/npm"
	"github.com/sethmcl/gofrosty/lib/test"
	"testing"
)

func TestParseDependency(t *testing.T) {
	g, err := gfjs.New(test.DataPath("gofrostyjs/gofrosty_simple.js"))
	test.Assert(err, nil, t)

	mockDep := &npm.Dependency{}
	mockModule := &module.Module{}

	err = g.ParseDependency(mockDep, mockModule)
	test.Assert(err, nil, t)
	test.Assert(mockModule.Name, "FooBar", t)
}
