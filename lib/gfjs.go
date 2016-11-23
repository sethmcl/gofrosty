package lib

import (
	js "github.com/sethmcl/gofrosty/vendor/ottoman"
)

// GoFrostyJS represents a parsed gofrosty.js config file
type GoFrostyJS struct {
	JSModule   string
	PathOnDisk string
}

// New returns *gfjs.GoFrostyJS
func LoadGoFrostyJS(filename string) (*GoFrostyJS, error) {
	jsModuleIdent, err := js.Require(filename)
	if err != nil {
		return nil, err
	}

	g := &GoFrostyJS{
		JSModule:   jsModuleIdent,
		PathOnDisk: filename,
	}

	return g, nil
}

// ParseDependency reads dependency and modifies transformed module
func (g *GoFrostyJS) ParseDependency(dep interface{}, module interface{}) error {
	var err error
	jsDep := "__frostyDependency__"
	jsMod := "__frostyModule__"

	err = js.SetGlobal(jsDep, dep)
	if err != nil {
		return err
	}

	err = js.SetGlobal(jsMod, module)
	if err != nil {
		return err
	}

	_, err = js.RunSrc("%s.parseDependency(%s, %s);", g.JSModule, jsDep, jsMod)
	if err != nil {
		return err
	}

	return nil
}
