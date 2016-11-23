package lib

import (
	"fmt"
	"io"
	"math"
	"os"
	"path"
	"regexp"
	"strings"
)

// IsFile returns true if file exists
func IsFile(file string) bool {
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
func IsDir(dir string) bool {
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

// IsTarURL returns true if url string points to a tar file
func IsTarURL(url string) bool {
	match, _ := regexp.MatchString("^http.*[\\.tar\\.gz|\\.tgz]$", url)
	return match
}

// IsGitURL returns true if this module's URL is a git repository
func IsGitURL(url string) bool {
	match, _ := regexp.MatchString("^git\\+", url)
	return match
}

// IsFileURL returns true if this module's URL is a local path on disk
func IsFileURL(url string) bool {
	match, _ := regexp.MatchString("^file:", url)
	return match
}

// CopyFile copies a file
func CopyFile(source string, dest string) error {
	sourceFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		sourceInfo, err := os.Stat(source)
		if err != nil {
			err = os.Chmod(dest, sourceInfo.Mode())
			return err
		}
	}

	return nil
}

// CopyDir copies a directory, recursively
func CopyDir(source string, dest string) error {
	return CopyDirBlacklist(source, dest, []string{})
}

// CopyDirBlacklist copies a directory, recursively, skipping files
// with name matching entry in blacklist
func CopyDirBlacklist(source string, dest string, blacklist []string) error {
	sourceInfo, err := os.Stat(source)
	if err != nil {
		return err
	}

	if IsDir(dest) {
		err := os.RemoveAll(dest)
		if err != nil {
			return err
		}
	}

	err = os.MkdirAll(dest, sourceInfo.Mode())
	if err != nil {
		return err
	}

	dir, err := os.Open(source)
	if err != nil {
		return err
	}

	objects, err := dir.Readdir(-1)
	if err != nil {
		return err
	}

	for _, obj := range objects {
		sourceFilePointer := path.Join(source, obj.Name())
		destFilePointer := path.Join(dest, obj.Name())

		if len(blacklist) > 0 {
			if StringSliceContains(blacklist, path.Base(obj.Name())) {
				continue
			}
		}

		if obj.IsDir() {
			err = CopyDir(sourceFilePointer, destFilePointer)
			if err != nil {
				return err
			}
		} else {
			err = CopyFile(sourceFilePointer, destFilePointer)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Mkdirs creates 1..N directories. Returns error if any fail.
func Mkdirs(dirs []string) error {
	for _, dir := range dirs {
		if !IsDir(dir) {
			err := os.MkdirAll(dir, os.ModePerm)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// ResolvePath resolves relative path to absolute path
func ResolvePath(filepath string, cwd string) string {
	if path.IsAbs(filepath) {
		return filepath
	}

	return path.Join(cwd, filepath)
}

// TruncatePath removes levels in a filepath. Path must be absolute.
// Example:
//   TruncatePath("/foo/bar/biz", 0) => "/foo/bar/biz")
//   TruncatePath("/foo/bar/biz", 1) => "/bar/biz")
//   TruncatePath("foo/bar/biz", 1) => "bar/biz")
//   TruncatePath("/foo/bar/biz", -1) => "/foo/bar")
//   TruncatePath("/foo/bar/biz", -2) => "/foo")
//   TruncatePath("/foo/bar/biz", 3) => "")
//   TruncatePath("/foo/bar/biz", 4) => ERROR!)
func TruncatePath(filepath string, count int) (string, error) {
	prepend := ""
	if path.IsAbs(filepath) {
		if count > 0 {
			count++
		}
		prepend = string(os.PathSeparator)
	}

	if count == 0 {
		return filepath, nil
	}

	absCount := int(math.Abs(float64(count)))

	parts := strings.Split(filepath, string(os.PathSeparator))
	if len(parts) == absCount {
		return "", nil
	}

	if len(parts) < absCount {
		return "", fmt.Errorf("Cannot truncate %d parts from path %s", count, filepath)
	}

	if count < 0 {
		return prepend +
			strings.Join(parts[1:len(parts)+count], string(os.PathSeparator)), nil
	}

	return prepend + strings.Join(parts[count:], string(os.PathSeparator)), nil
}
