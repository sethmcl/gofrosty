package lib

import (
	"os"
	"path"
)

// var instance *Cache

// Cache manages local cache on disk
type Cache struct {
	RootDir    string
	ModulesDir string
	IndexDir   string
	Index      *CacheIndex
}

// LoadCache creates a new instance of a Cache from a given cache directory
func LoadCache(dir string) (*Cache, error) {
	cache := &Cache{
		RootDir:    dir,
		ModulesDir: path.Join(dir, "modules"),
		IndexDir:   path.Join(dir, "index"),
	}

	err := Mkdirs([]string{
		cache.RootDir,
		cache.ModulesDir,
		cache.IndexDir})
	if err != nil {
		return nil, err
	}

	index, err := LoadCacheIndex(cache.IndexDir)
	if err != nil {
		return nil, err
	}

	cache.Index = index
	return cache, nil
}

// GetPath get path of module in cache based on name and version string
func (c *Cache) GetPath(name string, version string) (string, error) {
	dir, err := c.Index.Get(name, version)
	if err != nil {
		return "", err
	}
	return dir, nil
}

// Contains check if cache contains module name/version pair
func (c *Cache) Contains(name string, version string) bool {
	_, err := c.GetPath(name, version)
	if err != nil {
		return false
	}
	return true
}

// Add module to cache
func (c *Cache) Add(name, version string, pkg *Package) error {
	if c.Contains(name, version) {
		return nil
	}

	cacheDir := path.Join(c.ModulesDir, name, pkg.Version)
	if !IsDir(cacheDir) {
		// Directory exists but is not in index.
		err := os.MkdirAll(cacheDir, os.ModePerm)
		if err != nil {
			return err
		}

		// Copy module dir to cache. Exclude ./node_modules.
		err = CopyDirBlacklist(pkg.Dir, cacheDir, []string{"node_modules"})
		if err != nil {
			return err
		}
	}

	GetContext().Debug("CACHED %s@%s => %s", name, version, cacheDir)
	return c.Index.Add(name, version, cacheDir)
}
