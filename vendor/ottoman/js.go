package js

import (
	"fmt"
	"github.com/robertkrimen/otto"
	"io/ioutil"
)

var instance *otto.Otto
var jsModuleID = 0

// GetInstance returns otto.Otto JS engine instance
func GetInstance() *otto.Otto {
	if instance == nil {
		instance = otto.New()
	}
	return instance
}

// Reset force new instance of JS engine
func Reset() {
	instance = nil
	jsModuleID = 0
}

func nextJSModuleID() int {
	jsModuleID = jsModuleID + 1
	return jsModuleID
}

// Require load a JavaScript module
func Require(filepath string) (string, error) {
	src, err := ioutil.ReadFile(filepath)
	if err != nil {
		return "", err
	}

	engine := GetInstance()
	wrappedSrc := fmt.Sprintf("(function(module) {%s;return module.exports;})({});", src)
	result, err := engine.Run(wrappedSrc)
	if err != nil {
		return "", err
	}

	ident := fmt.Sprintf("__module%d__", nextJSModuleID())
	err = engine.Set(ident, result)
	if err != nil {
		return "", err
	}

	return ident, nil
}

// RunSrc executes a source string in the JS engine
func RunSrc(src string, fmtVals ...interface{}) (*otto.Value, error) {
	engine := GetInstance()
	fmtSrc := fmt.Sprintf(src, fmtVals...)
	wrappedSrc := fmt.Sprintf("(function() {%s})();", fmtSrc)
	result, err := engine.Run(wrappedSrc)
	return &result, err
}

// RunSrcGetStr executes a source string in the JS engine, expects string return value
func RunSrcGetStr(src string, fmtVals ...interface{}) (string, error) {
	result, err := RunSrc(src, fmtVals...)
	if err != nil {
		return "", err
	}

	resultStr, err := result.ToString()
	if err != nil {
		return "", err
	}
	return resultStr, err
}

// RunFile executes a file in the JS engine
func RunFile(filepath string) (*otto.Value, error) {
	src, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	return RunSrc(string(src[:]))
}

// SetGlobal sets a global value
func SetGlobal(name string, value interface{}) error {
	engine := GetInstance()
	return engine.Set(name, value)
}
