package test

import (
	"github.com/sethmcl/gofrosty/lib"
	"github.com/sethmcl/gofrosty/lib/testutil"
	"testing"
)

func TestListModuleVersions(t *testing.T) {
	test := testutil.New(t)
	reg := lib.NewNpmRegistryClient("https://registry.npmjs.org", "")
	_, err := reg.ListModuleVersions("webpack")
	test.Assert(err, nil)
	// test.Log(versions)
}

func TestGetTarURL(t *testing.T) {
	test := testutil.New(t)
	reg := lib.NewNpmRegistryClient("https://registry.npmjs.org", "")
	version := "1.0.0"
	table := map[string]string{
		"foo":     "https://registry.npmjs.org/foo/-/foo-1.0.0.tgz",
		"biz-buz": "https://registry.npmjs.org/biz-buz/-/biz-buz-1.0.0.tgz",
		"@me/bar": "https://registry.npmjs.org/@me/bar/-/bar-1.0.0.tgz",
	}

	for input, expected := range table {
		actual := reg.GetTarURL(input, version)
		test.Assert(actual, expected, input)
	}
}
