package lib

import (
	"fmt"
	"io/ioutil"
	"path"
	"regexp"
)

// CacheIndex manages local cache on disk
type CacheIndex struct {
	RootDir string
	Modules map[string]*ModuleIndex
}

// LoadCacheIndex creates a new cache index object
func LoadCacheIndex(dir string) (*CacheIndex, error) {
	err := Mkdirs([]string{dir})
	if err != nil {
		return nil, err
	}

	c := &CacheIndex{
		RootDir: dir,
		Modules: make(map[string]*ModuleIndex),
	}

	c.Load()
	return c, nil
}

// Commit saves index to disk
func (c *CacheIndex) Commit() error {
	for _, module := range c.Modules {
		err := module.Commit()
		if err != nil {
			return err
		}
	}
	return nil
}

// Load loads indices from files on disk
func (c *CacheIndex) Load() error {
	if !IsDir(c.RootDir) {
		return fmt.Errorf("Cache index directory (%s) does not exist", c.RootDir)
	}

	list, err := ioutil.ReadDir(c.RootDir)
	if err != nil {
		return err
	}

	suffixRe, err := regexp.Compile("\\.json$")
	if err != nil {
		return err
	}

	for _, file := range list {
		if !file.IsDir() {
			name := suffixRe.ReplaceAllString(file.Name(), "")
			c.Modules[name] = NewModuleIndex(name, c.RootDir)
			c.Modules[name].Load()
		}

		if file.IsDir() {
			scopedName := file.Name()
			scopedList, err := ioutil.ReadDir(path.Join(c.RootDir, scopedName))
			if err != nil {
				return err
			}

			for _, scopedFile := range scopedList {
				if !scopedFile.IsDir() {
					name := fmt.Sprintf("%s/%s", scopedName, scopedFile.Name())
					name = suffixRe.ReplaceAllString(name, "")
					c.Modules[name] = NewModuleIndex(name, path.Join(c.RootDir, scopedName))
					c.Modules[name].Load()
				}
			}
		}
	}

	return nil
}

// Add module to index
func (c *CacheIndex) Add(name string, version string, dir string) error {
	m, ok := c.Modules[name]
	if !ok {
		m = NewModuleIndex(name, c.RootDir)
		c.Modules[name] = m
	}

	m.Add(version, dir)
	return nil
}

// Get lookup module in index
func (c *CacheIndex) Get(name string, version string) (string, error) {
	cacheMiss := fmt.Errorf("Cache does not contain entry for %s@%s", name, version)
	m, ok := c.Modules[name]
	if !ok {
		return "", cacheMiss
	}

	entry := m.Get(version)
	if entry == "" {
		return "", cacheMiss
	}

	return entry, nil
}
