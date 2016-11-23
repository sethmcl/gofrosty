package test

import (
	"github.com/sethmcl/gofrosty/lib/npm"
	"github.com/sethmcl/gofrosty/lib/testutil"
	"github.com/sethmcl/gofrosty/vendor/semver"
	"testing"
)

// TestCreateCache tests creating a new FrostyCache instance
func TestCreateCache(t *testing.T) {
	test := testutil.New(t)

	dir, cleanup := test.TempDir()
	defer cleanup()
	cleanup()

	cache, err := npm.NewFrostyCache(dir)
	test.Assert(err, nil)
	test.AssertDir(cache.RootDir)
}

// TestCacheIndex
func TestCacheIndex(t *testing.T) {
	test := testutil.New(t)

	dir, cleanup := test.TempDir()
	defer cleanup()

	index, err := npm.NewCacheIndex(dir)
	test.Assert(err, nil)

	err = index.Add("my-module", "1.0.0", "/foo/bar")
	test.Assert(err, nil)

	err = index.Add("my-module", "~1.0.x", "/foo/bar/biz")
	test.Assert(err, nil)

	val, err := index.Get("my-module", "1.0.0")
	test.Assert(err, nil)
	test.Assert(val, "/foo/bar")

	val, err = index.Get("my-module", "~1.0.x")
	test.Assert(err, nil)
	test.Assert(val, "/foo/bar/biz")

	err = index.Commit()
	test.Assert(err, nil)

	index2, err := npm.NewCacheIndex(dir)
	test.Assert(err, nil)

	val, err = index2.Get("my-module", "~1.0.x")
	test.Assert(err, nil)
	test.Assert(val, "/foo/bar/biz")

	rr, err := semver.ParseRange(">=1.0.0 <1.1.0")
	test.Assert(err, nil)

	vers := []string{
		"*",
		"1.0.0",
		"1.0.100",
		"1.0.1+alpha",
		"1.1.0",
		"1.9999.0",
	}

	for _, ver := range vers {
		v, _ := semver.Parse(ver)
		if rr(v) {
			test.Logf("version %s satisfies range %s", ver, ">=1.0.0 <1.1.0")
		} else {
			test.Logf("version %s DOES NOT satisfy range %s", ver, ">=1.0.0 <1.1.0")
		}
	}
}
