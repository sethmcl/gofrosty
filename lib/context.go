package lib

import (
	"fmt"
	"log"
	"os"
	"path"
	"runtime/debug"
)

var instance *Context
var stdout *log.Logger

// Context represents an instance of running frosty
type Context struct {
	Cache          *Cache
	Cwd            string
	Force          bool
	FrostyHome     string
	GoFrostyJSPath string
	GoFrostyJS     *GoFrostyJS
	NodeModulesDir string
	NpmAuthToken   string
	NpmRegistry    *NpmRegistryClient
	PackagePath    string
	ShrinkwrapPath string
	UsePackage     bool
	UseShrinkwrap  bool
	Verbose        bool
}

// GetContext returns Context singleton
func GetContext() *Context {
	if instance == nil {
		instance = &Context{}
		instance.NpmAuthToken = GetNpmAuthToken()
		instance.NpmRegistry = NewNpmRegistryClient(
			GetNpmRegistryURL(),
			instance.NpmAuthToken)

		stdout = log.New(os.Stdout, "", log.Ldate)
	}

	return instance
}

// String returns string representation of the context object
func (c *Context) String() string {
	return fmt.Sprintf(`<<Context>>
  Cache           = %s
  Cwd             = %s
  Force           = %t
  FrostyHome      = %s
  GoFrostyJSPath  = %s
  NodeModulesDir  = %s
  NpmAuthToken    = %s
  NpmRegistry     = %s
  PackagePath     = %s
  ShrinkwrapPath  = %s
  UsePackage      = %t
  UseShrinkwrap   = %t
  Verbose         = %t
`,
		c.Cache.RootDir,
		c.Cwd,
		c.Force,
		c.FrostyHome,
		c.GoFrostyJSPath,
		c.NodeModulesDir,
		c.NpmAuthToken,
		c.NpmRegistry.RootURL,
		c.PackagePath,
		c.ShrinkwrapPath,
		c.UsePackage,
		c.UseShrinkwrap,
		c.Verbose)
}

// Info prints log message
func (c *Context) Info(args ...interface{}) {
	str := fmt.Sprintf(args[0].(string), args[1:]...)
	stdout.Println(str)
}

// Debug prints debug log message
func (c *Context) Debug(args ...interface{}) {
	if c.Verbose == true {
		str := fmt.Sprintf(args[0].(string), args[1:]...)
		stdout.Println(str)
	}
}

// DumpStack prints stack to stdout
func (c *Context) DumpStack() {
	if c.Verbose == true {
		stdout.Println(string(debug.Stack()))
	}
}

// LoadGoFrostyJSFile parse config file
func LoadGoFrostyJSFile(ctx *Context) error {
	candidates := []string{
		ctx.GoFrostyJSPath,
		path.Join(ctx.Cwd, "gofrosty.js"),
		path.Join(ctx.FrostyHome, "gofrosty.js"),
	}

	for _, candidate := range candidates {
		if IsFile(candidate) {
			gfj, err := LoadGoFrostyJS(candidate)
			if err != nil {
				ctx.Debug("Failed to load %s; Error: %s", candidate, err)
				continue
			} else {
				ctx.GoFrostyJS = gfj
				ctx.Debug("Loaded configuration from %s", candidate)
				break
			}
		}
	}

	return nil
}
