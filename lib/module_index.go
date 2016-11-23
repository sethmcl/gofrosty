package lib

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

// ModuleIndex cache index for a single module (unique by name)
type ModuleIndex struct {
	Filepath string            `json:"file"`
	Index    map[string]string `json:"index"`
}

// NewModuleIndex returns a new instance of a *ModuleIndex
func NewModuleIndex(name string, root string) *ModuleIndex {
	parts := strings.Split(name, "/")
	if len(parts) > 1 {
		name = parts[1] + ".json"
	} else {
		name = name + ".json"
	}

	return &ModuleIndex{
		Filepath: path.Join(root, name),
		Index:    make(map[string]string),
	}
}

// Add adds a new cache entry
func (m *ModuleIndex) Add(key string, dir string) {
	m.Index[key] = dir
}

// Get retrieves an existing cache entry
func (m *ModuleIndex) Get(key string) string {
	return m.Index[key]
}

// Delete removes an existing cache entry
func (m *ModuleIndex) Delete(key string) {
	delete(m.Index, key)
}

// Commit saves to disk
func (m *ModuleIndex) Commit() error {
	bytes, err := json.MarshalIndent(m, "", "    ")
	if err != nil {
		return err
	}

	// ensure parent directory exists
	err = os.MkdirAll(path.Dir(m.Filepath), os.ModePerm)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(m.Filepath, bytes, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

// Load module index file from disk
func (m *ModuleIndex) Load() error {
	bytes, err := ioutil.ReadFile(m.Filepath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, m)
	if err != nil {
		return err
	}

	return nil
}
