package testutil

import (
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"runtime/debug"
	"testing"
)

// TestUtil provides testing related helper functions
type TestUtil struct {
	t *testing.T
}

// New create new TestUtil instance
func New(t *testing.T) *TestUtil {
	return &TestUtil{
		t: t,
	}
}

// DataPath return path to data file under /test/data
func (tut *TestUtil) DataPath(parts ...string) string {
	wd, _ := os.Getwd()
	parts = append([]string{wd, "data"}, parts...)
	return path.Join(parts...)
}

// TempDir create temporary directory for use in tests
func (tut *TestUtil) TempDir() (string, func()) {
	dir, err := ioutil.TempDir("", "go-testutil")
	tut.Assert(err, nil)
	return dir, func() {
		os.RemoveAll(dir)
	}
}

// Log logs test message
func (tut *TestUtil) Log(args ...interface{}) {
	tut.t.Log(args...)
}

// Logf logs formatted test message
func (tut *TestUtil) Logf(format string, args ...interface{}) {
	tut.t.Logf(format, args...)
}

// IsFile returns true if file exists
func (tut *TestUtil) IsFile(file string) bool {
	f, err := os.Open(file)
	if err != nil {
		return false
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		return false
	}

	return !stat.IsDir()
}

// IsDir returns true if directory exists
func (tut *TestUtil) IsDir(dir string) bool {
	f, err := os.Open(dir)
	if err != nil {
		return false
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		return false
	}

	return stat.IsDir()
}

// Assert asserts a value
func (tut *TestUtil) Assert(actual interface{}, expected interface{}, args ...interface{}) {
	if actual != expected {
		loc := "(" + tut.getCallerStack() + ")"
		tut.t.Log(loc, "\nACTUAL   = ", actual, "\nEXPECTED = ", expected, args)
		tut.t.Fail()
	}
}

// AssertDir asserts that a path is a directory that exists on disk
func (tut *TestUtil) AssertDir(dir string) {
	if !tut.IsDir(dir) {
		loc := "(" + tut.getCallerStack() + ")"
		tut.t.Log(loc, "\n", dir, "is not a directory")
		tut.t.Fail()
	}
}

// AssertFile asserts that a path is a file that exists on disk
func (tut *TestUtil) AssertFile(file string) {
	if !tut.IsFile(file) {
		loc := "(" + tut.getCallerStack() + ")"
		tut.t.Log(loc, "\n", file, "is not a file")
		tut.t.Fail()
	}
}

// getStackTrace returns stack trace
func (tut *TestUtil) getCallerStack() string {
	raw := debug.Stack()
	re, err := regexp.Compile("[^/]*\\.go:\\d+")
	if err != nil {
		return ""
	}

	results := re.FindAll([]byte(raw), -1)

	if len(results) < 3 {
		return ""
	}

	return string(results[2])
}
