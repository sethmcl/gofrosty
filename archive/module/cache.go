package module

// import (
// 	"errors"
// 	"fmt"
// 	"github.com/sethmcl/gofrosty/lib/context"
// 	"github.com/sethmcl/gofrosty/lib/npm"
// 	"github.com/sethmcl/gofrosty/lib/tar"
// 	"github.com/sethmcl/gofrosty/lib/util"
// 	"net/http"
// 	"os"
// 	"os/exec"
// 	"path"
// 	"strings"
// )

// // Cache provides interaction with the frosty cache
// type Cache struct {
// 	Dir string
// }

// // Contains returns true if module is in cache
// func (c *Cache) Contains(module *Module) bool {
// 	if module.Cacheable == false {
// 		return false
// 	}

// 	moduleCacheDir := c.ModuleDir(module)

// 	file, err := os.Open(moduleCacheDir)
// 	if err != nil {
// 		return false
// 	}
// 	defer file.Close()

// 	stat, err := file.Stat()
// 	if err != nil {
// 		return false
// 	}

// 	return stat.IsDir()
// }

// // Add adds an npm module to the cache
// func (c *Cache) Add(module *Module) error {
// 	// Do not cache local modules
// 	if module.Type.IsFile() {
// 		return nil
// 	}

// 	if len(module.DownloadURLs) == 0 {
// 		return fmt.Errorf("No DownloadURL found for module %s", module.Name)
// 	}

// 	if module.Type.IsTar() {
// 		err := c.addTarModule(module)
// 		if err != nil {
// 			return err
// 		}
// 		return c.installModule(module)
// 	}

// 	if module.Type.IsGit() {
// 		return errors.New("Git NPM dependencies are not supported")
// 	}

// 	return fmt.Errorf("Do not know how to add module of type %s", module.Type.Value)
// }

// // Remove removes a module from the cache
// func (c *Cache) Remove(module *Module) error {
// 	dir := c.ModuleDir(module)
// 	return os.RemoveAll(dir)
// }

// func (c *Cache) addTarModule(module *Module) error {
// 	var (
// 		res         *http.Response
// 		err         error
// 		selectedURL string
// 	)

// 	ctx := context.GetInstance()

// 	for _, url := range module.DownloadURLs {
// 		ctx.Info("GET %s", url)
// 		res, err = downloadFile(url)
// 		if err == nil {
// 			defer res.Body.Close()
// 			selectedURL = url
// 			break
// 		} else {
// 			return err
// 		}
// 	}

// 	if res == nil {
// 		return fmt.Errorf("Unable to download module %s@%s", module.Name, module.Version)
// 	}

// 	if res.StatusCode >= 400 && res.StatusCode < 600 {
// 		return fmt.Errorf(
// 			"ERROR: [HTTP %d] Unable to download %s",
// 			res.StatusCode,
// 			selectedURL)
// 	}

// 	ctx.Debug("Downloaded %s", selectedURL)
// 	cachePath := c.ModuleDir(module)
// 	err = tar.Extract(res.Body, cachePath)
// 	if err != nil {
// 		ctx.Info("Unable to extract %s", selectedURL)
// 		return err
// 	}
// 	module.PathOnDisk = path.Join(cachePath, "package")

// 	return nil
// }

// // ModuleDir returns absolute path to module in cache, on disk
// func (c *Cache) ModuleDir(module *Module) string {
// 	return path.Join(c.Dir, module.Name, module.Version)
// }

// // ModulePackageDir returns absolute path to module in cache, on disk (package)
// func (c *Cache) ModulePackageDir(module *Module) string {
// 	return path.Join(c.ModuleDir(module), "package")
// }

// func (c *Cache) runNpmScript(module *Module, script string) error {
// 	var src string
// 	ctx := context.GetInstance()
// 	pkg := npm.NewPackage()
// 	err := pkg.Load(path.Join(module.PathOnDisk, "package.json"))
// 	if err != nil {
// 		return err
// 	}

// 	switch script {
// 	case "install":
// 		src = pkg.Scripts.Install
// 	case "postinstall":
// 		src = pkg.Scripts.PostInstall
// 	default:
// 		return fmt.Errorf("Do not know how to npm run-script %s", script)
// 	}

// 	if src != "" {
// 		ctx.Info("[npm run-script %s] %s@%s", script, module.Name, module.Version)
// 		ctx.Debug(src)

// 		parts := strings.Split(src, " ")
// 		cmd := exec.Command(parts[0], parts...)
// 		cmd.Dir = module.PathOnDisk
// 		out, err := cmd.CombinedOutput()

// 		if err != nil {
// 			ctx.Info("ERROR running `%s`", src)
// 			ctx.Info(string(out))

// 			// Since install failed, we should remove module from cache
// 			// c.Remove(module)

// 			return err
// 		}
// 	}

// 	return nil
// }

// func (c *Cache) installModule(module *Module) error {
// 	ctx := context.GetInstance()
// 	err := c.runNpmScript(module, "install")
// 	if err != nil {
// 		ctx.Info(err.Error())
// 		// return err
// 	}

// 	err = c.runNpmScript(module, "postinstall")
// 	if err != nil {
// 		ctx.Info(err.Error())
// 		// return err
// 	}
// 	// ctx := context.GetInstance()
// 	// pkg, err := npm.LoadPackageFile(path.Join(module.PathOnDisk, "package.json"))
// 	// if err != nil {
// 	// 	return err
// 	// }

// 	// // run `npm install`
// 	// if pkg.Scripts.Install != "" {
// 	// 	parts := strings.Split(pkg.Scripts.Install, " ")
// 	// 	ctx.Debug("[npm install] %s@%s", module.Name, module.Version)
// 	// 	ctx.Debug(strings.Join(parts, " "))

// 	// 	cmd := exec.Command(parts[0], parts...)
// 	// 	cmd.Dir = module.PathOnDisk
// 	// 	out, err := cmd.Output()

// 	// 	if err != nil {
// 	// 		ctx.Info("=== ERROR ===")
// 	// 		ctx.Info(string(out))

// 	// 		// Since install failed, we should remove module from cache
// 	// 		// c.Remove(module)

// 	// 		return err
// 	// 	}
// 	// }

// 	// // run `npm postinstall`
// 	// if pkg.Scripts.PostInstall != "" {
// 	// 	parts := strings.Split(pkg.Scripts.PostInstall, " ")
// 	// 	ctx.Debug("[npm postinstall] %s@%s", module.Name, module.Version)
// 	// 	ctx.Debug(pkg.Scripts.PostInstall)

// 	// 	cmd := exec.Command(parts[0], parts[1:]...)
// 	// 	out, err := cmd.Output()
// 	// 	if err != nil {
// 	// 		ctx.Info("ERROR running `%s`", pkg.Scripts.PostInstall)
// 	// 		ctx.Info(string(out))

// 	// 		// Since install failed, we should remove module from cache
// 	// 		// c.Remove(module)

// 	// 		return err
// 	// 	}
// 	// }

// 	return nil
// }

// func downloadFile(url string) (*http.Response, error) {
// 	ctx := context.GetInstance()
// 	client := &http.Client{}
// 	req, err := http.NewRequest("GET", url, nil)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if util.IsNpmRegistryURL(url) {
// 		ctx.Debug("Attaching NPM token (%s) to request", ctx.NpmAuthToken)
// 		auth := fmt.Sprintf("Bearer %s", ctx.NpmAuthToken)
// 		req.Header.Add("authorization", auth)
// 	}
// 	return client.Do(req)
// }
