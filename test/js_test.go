package test

import (
	"github.com/sethmcl/gofrosty/lib/js"
	"github.com/sethmcl/gofrosty/lib/test"
	"testing"
)

func TestRequire(t *testing.T) {
	filepath := test.DataPath("gofrostyjs/gofrosty_simple.js")
	js.Reset()
	ident, err := js.Require(filepath)
	test.Assert(err, nil, t)
	test.Assert(ident, "__module1__", t)
}
