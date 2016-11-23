package lib

import (
	"fmt"
)

// InstalledModule keeps track of an installed npm module
type InstalledModule struct {
	CacheKeys   []string
	InstallDirs []string
}

// NewInstalledModule creates a new installed module object
func NewInstalledModule(cacheKey string, installDir string) *InstalledModule {
	return &InstalledModule{
		CacheKeys:   []string{cacheKey},
		InstallDirs: []string{installDir},
	}
}

// AddCacheKey adds a new cache key
func (im *InstalledModule) AddCacheKey(cacheKey string) {
	if !im.contains(im.CacheKeys, cacheKey) {
		im.CacheKeys = append(im.CacheKeys, cacheKey)
	}
}

// AddInstallDir adds a new install dir
func (im *InstalledModule) AddInstallDir(dir string) {
	if !im.contains(im.InstallDirs, dir) {
		im.InstallDirs = append(im.InstallDirs, dir)
	}
}

// contains check if string is present in specified string slice check if cache key is present in package
func (im *InstalledModule) contains(haystack []string, needle string) bool {
	for _, item := range haystack {
		if item == needle {
			return true
		}
	}
	return false
}

// InstallContext keeps track of npm modules which are installed, so that they can be added to the cache after
// the install process completes.
type InstallContext struct {
	Modules map[string]*InstalledModule
}

// NewInstallContext Create new install context
func NewInstallContext() *InstallContext {
	return &InstallContext{
		Modules: make(map[string]*InstalledModule, 0),
	}
}

// Add module which has been installed
func (ictx *InstallContext) Add(name string, explicitSemver string, cacheKey string, installDir string) {
	id := ictx.GenID(name, explicitSemver)
	m, ok := ictx.Modules[id]
	if !ok {
		ictx.Modules[id] = NewInstalledModule(cacheKey, installDir)
		return
	}

	m.AddCacheKey(cacheKey)
	m.AddInstallDir(installDir)
}

// Get a module which has been installed
func (ictx *InstallContext) Get(name string, explicitSemver string) (*InstalledModule, error) {
	id := ictx.GenID(name, explicitSemver)
	m, ok := ictx.Modules[id]
	if !ok {
		return nil, fmt.Errorf("%s has not been installed", id)
	}
	return m, nil
}

// GenID Return module id, which is name@version
func (ictx *InstallContext) GenID(name string, explicitSemver string) string {
	return fmt.Sprintf("%s@%s", name, explicitSemver)
}
