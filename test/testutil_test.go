package test

import (
	"github.com/sethmcl/gofrosty/lib/test"
	"os"
	// "os/exec"
	"path"
	"testing"
)

func TestDataPath(t *testing.T) {
	wd, _ := os.Getwd()
	expected := path.Join(wd, "data", "foo", "bar.json")
	actual := test.DataPath("foo", "bar.json")

	if actual != expected {
		t.Error("Actual:", actual, ", Expected:", expected)
	}
}

func TestTempDir(t *testing.T) {
	dir, cleanup := test.TempDir()
	test.AssertDir(dir)
	cleanup()
	test.Assert(test.IsDir(dir), false)
}

// func TestTest(t *testing.T) {
// 	cmd := exec.Command("env")
// 	cmd.Dir = "/tmp"
// 	out, err := cmd.Output()
// 	test.Assert(err, nil, t)
// 	t.Log(string(out))
// }
