package lib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
)

// Package represents a single package.json file
type Package struct {
	RawBin             interface{} `json:"bin"`
	RawDependencies    interface{} `json:"dependencies"`
	RawDevDependencies interface{} `json:"devDependencies"`
	Name               string      `json:"name"`
	Resolved           string      `json:"_resolved"`
	Scripts            Scripts     `json:"scripts"`
	Version            string      `json:"version"`
	Main               string      `json:"main"`

	Bin             map[string]string `json:"-"`
	DevDependencies map[string]string `json:"-"`
	Dependencies    map[string]string `json:"-"`
	Filepath        string            `json:"-"`
	Dir             string            `json:"-"`
}

// Scripts represents a scripts property from a package.json file
type Scripts struct {
	PostInstall string `json:"postinstall"`
	Install     string `json:"install"`
}

// NewPackage creates a new package.json struct
func NewPackage() *Package {
	return &Package{}
}

// LoadPackageFromDir loads a package file from a directory
// assumes file is named package.json
func LoadPackageFromDir(dir string) (*Package, error) {
	return LoadPackage(path.Join(dir, "package.json"))
}

// LoadPackage loads a package from a file
func LoadPackage(filePath string) (*Package, error) {
	pkg := NewPackage()
	err := pkg.Load(filePath)
	if err != nil {
		return nil, err
	}
	return pkg, nil
}

// Load loads data from file
func (p *Package) Load(filePath string) error {
	jsonstr, err := ioutil.ReadFile(filePath)

	if err != nil {
		return err
	}

	err = json.Unmarshal(jsonstr, p)
	if err != nil {
		GetContext().Debug("Failed to parse %s", filePath)
		GetContext().Debug("\n%s\n", jsonstr)
		return err
	}

	// This is needed because npm is too forgiving in how you specify the
	// bin value. Ideally this should be specified as:
	// {
	//     "bin": {
	//         "name": "./path/to/bin/file"
	//     }
	// }
	//
	// However, this is also valid:
	// {
	//     "bin": "./path/to/bin/file"
	// }
	//
	// When this syntax is used, the bin name is set to the package name
	p.Bin = make(map[string]string)
	binStr, ok := p.RawBin.(string)
	if ok {
		p.Bin[p.Name] = binStr
	}
	binMap, ok := p.RawBin.(map[string]interface{})
	if ok {
		for k, v := range binMap {
			sv, ok := v.(string)
			if ok {
				p.Bin[k] = sv
			}
		}
	}

	// This is needed because some people think it is entertaining to specify
	// dependencies in their package.json as an array instead of as a map
	p.Dependencies = make(map[string]string)
	deps, ok := p.RawDependencies.(map[string]interface{})
	if ok {
		for k, v := range deps {
			sv, ok := v.(string)
			if ok {
				p.Dependencies[k] = sv
			}
		}
	}

	// This is needed because some people think it is entertaining to specify
	// devDependencies in their package.json as an array instead of as a map
	p.DevDependencies = make(map[string]string)
	devDeps, ok := p.RawDevDependencies.(map[string]interface{})
	if ok {
		for k, v := range devDeps {
			sv, ok := v.(string)
			if ok {
				p.DevDependencies[k] = sv
			}
		}
	}

	p.Filepath = filePath
	p.Dir = path.Dir(filePath)
	return nil
}

// Commit saves to disk
func (p *Package) Commit() error {
	bytes, err := json.MarshalIndent(p, "", "    ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(p.Filepath, bytes, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

// RunScript an npm script from package.json
func (p *Package) RunScript(script string) error {
	var src string

	ctx := GetContext()

	switch script {
	case "install":
		src = p.Scripts.Install
	case "postinstall":
		src = p.Scripts.PostInstall
	default:
		return fmt.Errorf("Do not know how to npm run-script %s", script)
	}

	if src != "" {
		ctx.Info("Running %s script from %s: `%s`", script, p.Filepath, src)
		ctx.Debug(src)

		parts := strings.Split(src, " ")
		cmd := exec.Command(parts[0], parts[1:]...)
		cmd.Dir = p.Dir
		out, err := cmd.CombinedOutput()

		if err != nil {
			ctx.Info("ERROR running `%s`", src)
			ctx.Info(string(out))
			return err
		}
	}

	return nil
}

// LinkBin creates symlink for {bin: '..'} entries
func (p *Package) LinkBin(relBin string, relModule string) error {
	if len(p.Bin) == 0 {
		return nil
	}

	if path.IsAbs(relBin) {
		return fmt.Errorf(
			"Please provide relative path to node_modules/bin. Got: %s",
			relBin)
	}

	target := path.Join(p.Dir, relBin)
	err := os.MkdirAll(target, os.ModePerm)
	if err != nil {
		return err
	}

	for binName, relPath := range p.Bin {
		symlinkFile := path.Join(target, binName)
		targetFile := path.Join(relModule, relPath)

		if IsFile(symlinkFile) {
			err := os.RemoveAll(symlinkFile)
			if err != nil {
				return err
			}
		}

		err := os.Symlink(targetFile, symlinkFile)
		if err != nil {
			GetContext().Debug("Failed to symlink %s => %s", symlinkFile, targetFile)
			return err
		}
		GetContext().Debug("Created symlink %s => %s", symlinkFile, targetFile)

		err = os.Chmod(path.Join(p.Dir, relPath), os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}
