package lib

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path"
)

// InstallCmdInit initialized application Context
func InstallCmdInit(args []string) error {
	ctx := GetContext()

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	installCmd := flag.NewFlagSet("install", flag.ExitOnError)
	cwdFlag := installCmd.String("C", cwd, "Set working directory")
	verboseFlag := installCmd.Bool("verbose", false, "Show verbose log output")
	forceFlag := installCmd.Bool("force", false,
		"Force install -- continue even if some packages fail to install")
	frostyHomeFlag := installCmd.String("frosty-home", GetFrostyHome(),
		"Location of frosty home directory")
	configFlag := installCmd.String("config", "", "Path to gofrosty.js configuration")
	shrinkwrapFlag := installCmd.Bool("shrinkwrap", false, "Force usage of npm-shrinkwrap.json")
	packageFlag := installCmd.Bool("package", false, "Force usage of package.json")

	installCmd.Parse(args)

	ctx.Cwd = *cwdFlag
	ctx.Verbose = *verboseFlag
	ctx.Force = *forceFlag
	ctx.FrostyHome = *frostyHomeFlag
	ctx.GoFrostyJSPath = *configFlag

	ctx.Cwd = ResolvePath(*cwdFlag, cwd)
	ctx.NodeModulesDir = path.Join(ctx.Cwd, "node_modules")
	ctx.FrostyHome = ResolvePath(*frostyHomeFlag, cwd)

	if *configFlag != "" {
		ctx.GoFrostyJSPath = ResolvePath(*configFlag, cwd)
	}

	if *shrinkwrapFlag && *packageFlag {
		return errors.New(
			"Please specify either --shrinkwrap or --package flag, but not both")
	}

	if *shrinkwrapFlag {
		ctx.UseShrinkwrap = true
		ctx.UsePackage = false
	}

	if *packageFlag {
		ctx.UseShrinkwrap = false
		ctx.UsePackage = true
	}

	candidateShrinkwrapPath := path.Join(ctx.Cwd, "npm-shrinkwrap.json")
	candidatePackagePath := path.Join(ctx.Cwd, "package.json")

	if ctx.UseShrinkwrap {
		ctx.ShrinkwrapPath = candidateShrinkwrapPath
	}

	if ctx.UsePackage {
		ctx.PackagePath = candidatePackagePath
	}

	if !ctx.UseShrinkwrap && !ctx.UsePackage {
		if IsFile(candidateShrinkwrapPath) {
			ctx.UseShrinkwrap = true
			ctx.ShrinkwrapPath = candidateShrinkwrapPath
			ctx.UsePackage = false
		} else {
			ctx.UsePackage = true
			ctx.PackagePath = candidatePackagePath
			ctx.UseShrinkwrap = false
		}
	}

	LoadGoFrostyJSFile(ctx)

	// initialize the local cache
	cache, err := LoadCache(path.Join(ctx.FrostyHome, "cache"))
	if err != nil {
		return err
	}
	ctx.Cache = cache

	return nil
}

// InstallCmdRun runs the install command
func InstallCmdRun(args []string) error {
	ctx := GetContext()

	err := InstallCmdInit(args)
	if err != nil {
		return err
	}

	// Print run-time context values, for debugging
	ctx.Debug(ctx.String())

	if ctx.UseShrinkwrap {
		return installFromShrinkwrapJSON(ctx)
	}

	if ctx.UsePackage {
		return installFromPackageJSON(ctx)
	}

	return nil
}

// Install npm modules from npm-shrinkwrap.json file
func installFromShrinkwrapJSON(ctx *Context) error {
	ctx.Debug("Installing modules from %s...", ctx.ShrinkwrapPath)

	// // Parse npm-shrinkwrap.json file
	// shrinkwrap, err := LoadShrinkwrapFile(ctx.ShrinkwrapPath)
	// if err != nil {
	// 	return err
	// }

	// // Map dependencies from npm-shrinkwrap.json to frosty module objects
	// modules, err := module.TransformDependencies(shrinkwrap.FlattenDeps())
	// if err != nil {
	// 	return err
	// }

	// // Add each module to the local cache, downloading from network as needed.
	// err = populateCache(cache, modules)
	// if err != nil {
	// 	return err
	// }

	// // Materialize npm-shrinkwrap.json dependency tree in ./node_modules
	// // from modules stored  in local cache
	// err = materialize(shrinkwrap)
	// if err != nil {
	// 	return err
	// }

	return nil
}

// Install npm modules from package.json file
func installFromPackageJSON(ctx *Context) error {
	ctx.Debug("Installing modules from %s...", ctx.PackagePath)
	ictx := NewInstallContext()
	pkg, err := LoadPackage(ctx.PackagePath)
	if err != nil {
		return err
	}

	err = installDepMap(pkg.Dependencies, ctx.NodeModulesDir, ictx)
	if err != nil {
		return err
	}

	err = ctx.Cache.Index.Commit()
	if err != nil {
		return err
	}

	return nil
}

func installDepMap(deps map[string]string, nodeModulesDir string, ictx *InstallContext) error {
	for name, version := range deps {
		installDir := path.Join(nodeModulesDir, name)
		err := installDep(name, version, installDir, ictx)
		if err != nil {
			return err
		}
	}
	return nil
}

func installDep(name string, version string, installDir string, ictx *InstallContext) error {
	if IsDir(installDir) {
		return nil
	}

	ctx := GetContext()
	ctx.Debug("Installing %s@%s to %s", name, version, installDir)

	if IsDir(installDir) {
		err := os.RemoveAll(installDir)
		if err != nil {
			return err
		}
	}

	os.MkdirAll(installDir, os.ModePerm)
	cacheDir, cacheMiss := ctx.Cache.GetPath(name, version)

	var pkg *Package
	var err error

	if cacheMiss != nil {
		ctx.Debug("CACHE MISS %s@%s --- %s", name, version, cacheMiss.Error())
		return cacheMiss
		// pkg, err = installDepFromNetwork(name, version, installDir, ictx)
		// if err != nil {
		// 	return err
		// }
	} else {
		ctx.Debug("CACHE HIT %s@%s [%s]", name, version, cacheDir)
		pkg, err = installDepFromCacheDir(cacheDir, installDir)
		if err != nil {
			return err
		}
	}

	err = installDepMap(pkg.Dependencies, path.Join(installDir, "node_modules"), ictx)
	if err != nil {
		return err
	}

	// once deps are installed, we can run install and postinstall scripts from package.json
	// only need to do this if dep was not installed from cache, because once a module is in
	// the cache, these scripts have already been run
	if cacheMiss != nil {
		err = pkg.RunScript("install")
		if err != nil {
			return err
		}

		err = pkg.RunScript("postinstall")
		if err != nil {
			return err
		}

		err := ctx.Cache.Add(name, version, pkg)
		if err != nil {
			return err
		}
	}

	// symlink bin files from package.json
	err = pkg.LinkBin(path.Join("..", ".bin"), path.Join("..", path.Base(installDir)))
	if err != nil {
		return err
	}

	return nil
}

func installDepFromCacheDir(cacheDir, installDir string) (*Package, error) {
	err := CopyDir(cacheDir, installDir)
	if err != nil {
		return nil, err
	}

	pkg, err := LoadPackageFromDir(installDir)
	if err != nil {
		return nil, err
	}

	return pkg, nil
}

func installDepFromNetwork(name, version, installDir string, ictx *InstallContext) (*Package, error) {
	res, err := downloadDep(name, version)
	if err != nil {
		return nil, err
	}

	// untar file
	err = ExtractTar(res.Body, installDir, 1)
	if err != nil {
		return nil, err
	}
	res.Body.Close()

	// update package.json to reflect where code was downloaded from
	pkg, err := LoadPackageFromDir(installDir)
	if err != nil {
		return nil, err
	}

	pkg.Resolved = res.Request.URL.String()
	err = pkg.Commit()
	if err != nil {
		return nil, err
	}

	ictx.Add(name, pkg.Version, version, installDir)

	return pkg, nil
}

func downloadDep(name string, version string) (*http.Response, error) {
	// At this point, we know the module is not in the cache, so we need to
	// grab it from the network. The version string can be one of the following:
	//
	//    - A URL to a TAR file (ex: http://foo.com/my-tar.tgz)
	//    - An explicit semver string (ex: 1.0.0)
	//    - A semver range string (ex: 1.0.x)
	//    - A URL to a GIT repo (not supported)
	//
	// If the URL is a git repo, we return an error (not supported).
	//
	// If the URL is an explicit semver version, we generate a TAR file URL,
	// based on the configured npm registry host.
	//
	// If the URL is a semver range, we query the registry to get the newest version
	// which satisfies the semver. We then generate the TAR file URL for that version,
	// based on the configured npm registry host.
	//
	// If version is a TAR file URL, or once we have converted it to a TAR file URL in one
	// of the previous steps, we download it from the network, and install it.

	ctx := GetContext()

	if IsGitURL(version) {
		return nil, fmt.Errorf("Git dependencies are not supported [%s]", version)
	}

	downloadURL, err := ctx.NpmRegistry.GetTarURL(name, version)
	if err != nil {
		return nil, err
	}

	return ctx.NpmRegistry.Get(downloadURL)
}

// func materialize(shrinkwrap *Shrinkwrap) error {
// nodeModulesDir := path.Join(shrinkwrap.Dir, "node_modules")
// err := os.MkdirAll(nodeModulesDir, os.ModePerm)
// if err != nil {
// 	return err
// }

// err = materializeDeps(nodeModulesDir, shrinkwrap.Dependencies)
// if err != nil {
// 	return err
// }

// 	return nil
// }

// func materializeDeps(nodeModulesDir string, deps map[string]*Dependency) error {
// for _, dep := range deps {
// 	module, err := module.TransformDependency(dep)
// 	if err != nil {
// 		return err
// 	}

// 	materializeModule(nodeModulesDir, module)
// 	subNodeModulesDir := path.Join(nodeModulesDir, dep.Name, "node_modules")
// 	materializeDeps(subNodeModulesDir, dep.Dependencies)
// }

// 	return nil
// }

// func materializeModule(nodeModulesDir string, module *module.Module) error {
// ctx := GetContext()
// if !ctx.Cache.Contains(module) {
// 	return fmt.Errorf(
// 		"Cannot materialize %s@%s [not in cache]",
// 		module.Name,
// 		module.Version)
// }

// source := module.PathOnDisk
// dest := path.Join(nodeModulesDir, module.Name)

// if module.PathOnDisk == "" {
// 	source = cache.ModulePackageDir(module)
// }

// err := CopyDir(source, dest)
// if err != nil {
// 	return err
// }

// 	return nil
// }

// func populateCache(cache *module.Cache, modules []*module.Module) error {
// ctx := GetContext()

// for _, module := range modules {
// 	ctx.Debug(module.String())
// 	if !cache.Contains(module) {
// 		err := cache.Add(module)
// 		if err != nil {
// 			return err
// 		}
// 	} else {
// 		ctx.Debug("Cache already contains %s@%s, continuing.",
// 			module.Name, module.Version)
// 	}
// }

// 	return nil
// }
