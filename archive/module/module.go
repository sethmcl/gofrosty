package module

// import (
// 	"errors"
// 	"fmt"
// 	"github.com/sethmcl/gofrosty/lib/context"
// 	"github.com/sethmcl/gofrosty/lib/npm"
// 	"github.com/sethmcl/gofrosty/lib/util"
// 	"path"
// 	"strings"
// )

// // This is the name of the function which can be implemented in gofrosty.js
// // to parse dependencies.
// const parseDependencyFnName string = "parseDependency"

// // Module represents a single npm module
// type Module struct {
// 	Name         string
// 	Version      string
// 	URL          string
// 	DownloadURLs []string
// 	Type         *Type
// 	PathOnDisk   string
// 	Cacheable    bool
// }

// // ToString returns string representation of object
// func (m *Module) String() string {
// 	return fmt.Sprintf(`<<Module>>
//   Name         = %s
//   Version      = %s
//   URL          = %s
//   DownloadURLs = %s
//   Type         = %s
//   Cacheable    = %t
//   PathOnDisk   = %s
// `, m.Name, m.Version, m.URL, m.DownloadURLs, m.Type, m.Cacheable, m.PathOnDisk)
// }

// // AppendDownloadURL add download url
// func (m *Module) AppendDownloadURL(url string) {
// 	m.DownloadURLs = append(m.DownloadURLs, url)
// }

// // PrependDownloadURL prepend download url
// func (m *Module) PrependDownloadURL(url string) {
// 	m.DownloadURLs = append([]string{url}, m.DownloadURLs...)
// }

// // SetDownloadURLs replace download urls
// func (m *Module) SetDownloadURLs(urls []string) {
// 	m.DownloadURLs = urls
// }

// // TransformDependency convert shrinkwrap dependency to an npm module object
// func TransformDependency(dep *npm.Dependency) (*Module, error) {
// 	module := &Module{
// 		Name:      dep.Name,
// 		Version:   dep.Version,
// 		URL:       dep.Resolved,
// 		Cacheable: true,
// 		Type:      &Type{},
// 	}

// 	// If there is not a URL in the npm-shrinkwrap.json for a module,
// 	// then try to generate it based on the known format of npm registry
// 	// urls
// 	if module.URL == "" {
// 		module.URL = util.GenerateNpmRegistryURL(dep.Name, dep.Version)
// 	}

// 	if util.IsFileURL(module.URL) {
// 		relativePath := dep.Resolved[len("file:"):]
// 		module.Cacheable = false
// 		module.PathOnDisk = path.Join(dep.ShrinkwrapDir, relativePath)
// 		module.Type.SetFile()
// 	}

// 	if util.IsTarURL(module.URL) {
// 		parts := strings.Split(module.URL, "/")
// 		module.Type.SetTar()
// 		module.AppendDownloadURL(module.URL)

// 		if !util.IsNpmRegistryURL(module.URL) {
// 			module.Version = parts[len(parts)-1]
// 		}
// 	}

// 	if util.IsGitURL(module.URL) {
// 		return nil, errors.New("Git dependencies are not supported.")
// 	}

// 	ctx := context.GetInstance()
// 	if ctx.GoFrostyJS != nil {
// 		err := ctx.GoFrostyJS.ParseDependency(dep, module)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}

// 	return module, nil
// }

// // TransformDependencies convert slice of shrinkwrap dependencies to a slice of npm modules
// func TransformDependencies(deps []*npm.Dependency) ([]*Module, error) {
// 	var modules []*Module
// 	for _, dep := range deps {
// 		module, err := TransformDependency(dep)
// 		if err != nil {
// 			return nil, err
// 		}
// 		modules = append(modules, module)
// 	}
// 	return modules, nil
// }
