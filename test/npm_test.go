package test

import (
	"github.com/sethmcl/gofrosty/lib/npm"
	"github.com/sethmcl/gofrosty/lib/test"
	"path"
	"testing"
)

func LoadPackageJSON() *npm.Package {
	file := test.DataPath("package.json")
	pkg, err := npm.LoadPackageFile(file)
	if err != nil {
		panic(err)
	}
	return pkg
}

func LoadShrinkwrapJSON() *npm.Shrinkwrap {
	file := test.DataPath("npm-shrinkwrap-2.11.3.json")
	shrink, err := npm.LoadShrinkwrapFile(file)
	if err != nil {
		panic(err)
	}
	return shrink
}

func TestName(t *testing.T) {
	shrink := LoadShrinkwrapJSON()
	actual := shrink.Name
	expected := "project"
	test.Assert(actual, expected, t)
}

func TestPath(t *testing.T) {
	shrink := LoadShrinkwrapJSON()
	actual := shrink.Path
	expected := test.DataPath("npm-shrinkwrap-2.11.3.json")
	test.Assert(actual, expected, t)
}

func TestDir(t *testing.T) {
	shrink := LoadShrinkwrapJSON()
	actual := shrink.Dir
	expected := path.Dir(test.DataPath("npm-shrinkwrap-2.11.3.json"))
	test.Assert(actual, expected, t)
}

func TestDependencies(t *testing.T) {
	shrink := LoadShrinkwrapJSON()

	matches := 0
	for k := range shrink.Dependencies {
		switch k {
		case "lib":
			matches++
		case "reference-node-module":
			matches++
		default:
			matches--
		}
	}
	test.Assert(matches, 2, t)

	lib := shrink.Dependencies["lib"]
	test.Assert(lib.Version, "1.0.0", t)
	test.Assert(lib.From, "../lib", t)
	test.Assert(lib.Resolved, "file:../lib", t)

	smallUUID := lib.Dependencies["small-uuid"]
	test.Assert(smallUUID.Version, "1.0.1", t)
	test.Assert(smallUUID.From, "small-uuid@>=1.0.1 <2.0.0", t)
	test.Assert(
		smallUUID.Resolved,
		"https://registry.npmjs.org/small-uuid/-/small-uuid-1.0.1.tgz",
		t,
	)

	matches = 0
	for k := range smallUUID.Dependencies {
		switch k {
		case "node-uuid":
			matches++
		default:
			matches--
		}
	}
	test.Assert(matches, 1, t)
}

func TestFlattenDeps(t *testing.T) {
	shrink := LoadShrinkwrapJSON()
	deps := shrink.FlattenDeps()
	test.Assert(len(deps), 4, t)

	matches := 0
	for _, dep := range deps {
		switch dep.Name {
		case "node-uuid":
			matches++
		case "small-uuid":
			matches++
		case "reference-node-module":
			matches++
		case "lib":
			matches++
		default:
			matches--
		}
	}
	test.Assert(matches, 4, t)

	test.Assert(deps[0].Name, "lib", t)
	test.Assert(deps[0].ShrinkwrapPath, shrink.Path, t)
	test.Assert(deps[0].ShrinkwrapDir, shrink.Dir, t)
}

func TestLoadPackageFile(t *testing.T) {
	pkg := LoadPackageJSON()
	test.Assert(pkg.Scripts.PostInstall, "make post-install", t)
	test.Assert(pkg.Scripts.Install, "make install", t)
}
