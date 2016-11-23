package test

import (
	"github.com/sethmcl/gofrosty/lib"
	"github.com/sethmcl/gofrosty/lib/testutil"
	"testing"
)

func TestInstallContext(t *testing.T) {
	test := testutil.New(t)

	ictx := lib.NewInstallContext()
	ictx.Add("foo", "1.1.5", "~1.1.0", "/module/foo")
	m, err := ictx.Get("foo", "1.1.5")
	test.Assert(err, nil)
	test.Assert(len(m.CacheKeys), 1)
	test.Assert(m.CacheKeys[0], "~1.1.0")
	test.Assert(len(m.InstallDirs), 1)
	test.Assert(m.InstallDirs[0], "/module/foo")

	ictx.Add("foo", "1.1.5", "^1.0.0", "/module/foo/1")
	ictx.Add("foo", "1.1.5", "^1.1.0", "/module/foo/1")
	test.Assert(len(m.CacheKeys), 3)
	test.Assert(len(m.InstallDirs), 2)
}
