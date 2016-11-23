package test

import (
	"github.com/sethmcl/gofrosty/lib/npm"
	"github.com/sethmcl/gofrosty/lib/testutil"
	"testing"
)

func TestPackage(t *testing.T) {
	test := testutil.New(t)

	pkgFile := test.DataPath("package.json")
	pkg := npm.NewPackage()
	err := pkg.Load(pkgFile)
	test.Assert(err, nil)
	test.Assert(pkg.Name, "test-module")
	test.Assert(pkg.Version, "1.0.0")
	test.Assert(pkg.Dependencies["small-uuid"], "^1.0.1")
	test.Assert(pkg.Scripts.PostInstall, "make post-install")
	test.Assert(pkg.Scripts.Install, "make install")
}
