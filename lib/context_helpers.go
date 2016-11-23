package lib

import (
	"os"
	"path"
)

// GetFrostyHome returns absolute path to default frosty home directory
func GetFrostyHome() string {
	home := os.Getenv("HOME")
	frostyHome := os.Getenv("FROSTY_HOME")

	if frostyHome != "" {
		return frostyHome
	}

	return path.Join(home, ".gofrosty")
}

// GetNpmAuthToken returns NPM_TOKEN environment variable
func GetNpmAuthToken() string {
	return os.Getenv("NPM_TOKEN")
}

// GetNpmRegistryURL returns NPM_REGISTRY_URL environment variable
// if this is not set, then returns default "https://registry.npmjs.org" value
func GetNpmRegistryURL() string {
	url := os.Getenv("NPM_REGISTRY_URL")

	if url == "" {
		url = "https://registry.npmjs.org"
	}

	return url
}
