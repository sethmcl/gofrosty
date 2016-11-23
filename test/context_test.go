package test

import (
	"github.com/sethmcl/gofrosty/lib/context"
	"github.com/sethmcl/gofrosty/lib/test"
	"os"
	"path"
	"testing"
)

func TestGetFrostyHome(t *testing.T) {
	var (
		actual   string
		expected string
	)

	actual = context.GetFrostyHome()
	expected = path.Join(os.Getenv("HOME"), ".gofrosty")
	test.Assert(actual, expected, t)

	mock := "/foo/bar/.home"
	os.Setenv("FROSTY_HOME", mock)
	actual = context.GetFrostyHome()
	expected = mock
	test.Assert(actual, expected, t)
}
