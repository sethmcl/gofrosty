package test

import (
	"github.com/sethmcl/gofrosty/lib/module"
	"github.com/sethmcl/gofrosty/lib/test"
	"io/ioutil"
	"log"
	"os"
	"path"
	"testing"
)

func CreateCache() *module.Cache {
	tempDir, _ := ioutil.TempDir("", "gofrosty_test")
	return &module.Cache{
		Dir: tempDir,
	}
}

func TestContains(t *testing.T) {
	cache := CreateCache()
	module := &module.Module{
		Name:      "foo",
		Version:   "1.0.0",
		URL:       "http://foo.com/t.tgz",
		Cacheable: true,
	}
	test.Assert(cache.Contains(module), false, t)

	err := os.MkdirAll(path.Join(cache.Dir, module.Name, module.Version), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	test.Assert(cache.Contains(module), true, t)
}
