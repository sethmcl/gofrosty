package lib

import (
	"encoding/json"
	"io/ioutil"
	"path"
)

// Dependency represents a dependency specification inside an npm-shrinkwrap.json file
type Dependency struct {
	Version        string                 `json:"version"`
	From           string                 `json:"from"`
	Resolved       string                 `json:"resolved"`
	Dependencies   map[string]*Dependency `json:"dependencies"`
	Name           string
	ShrinkwrapDir  string
	ShrinkwrapPath string
}

// Shrinkwrap represents a single npm-shrinkwrap.json file
type Shrinkwrap struct {
	Path         string
	Dir          string
	Name         string                 `json:"name"`
	Version      string                 `json:"version"`
	Dependencies map[string]*Dependency `json:"dependencies"`
}

// LoadShrinkwrapFile returns Shrinkwrap object
func LoadShrinkwrapFile(filePath string) (*Shrinkwrap, error) {
	jsonstr, err := ioutil.ReadFile(filePath)

	if err != nil {
		return nil, err
	}

	shrink := &Shrinkwrap{Path: filePath, Dir: path.Dir(filePath)}
	root := &Dependency{}
	err = json.Unmarshal(jsonstr, &root)
	if err != nil {
		return nil, err
	}

	shrink.Name = root.Name
	shrink.Version = root.Version
	shrink.Dependencies = root.Dependencies
	return shrink, nil
}

func addDeps(
	deps map[string]*Dependency,
	result *[]*Dependency,
	seendeps map[string]struct{},
	sp string) {
	for depName, dep := range deps {
		dep.Name = depName
		dep.ShrinkwrapPath = sp
		dep.ShrinkwrapDir = path.Dir(sp)
		key := depName + "@" + dep.Version
		_, ok := seendeps[key]
		if !ok {
			*result = append(*result, dep)
			seendeps[key] = struct{}{}
		}
		addDeps(dep.Dependencies, result, seendeps, sp)
	}
}

// FlattenDeps returns a flattened list of unique dependencies
func (n *Shrinkwrap) FlattenDeps() []*Dependency {
	var result []*Dependency
	seendeps := make(map[string]struct{}, 0)
	addDeps(n.Dependencies, &result, seendeps, n.Path)
	return result
}
