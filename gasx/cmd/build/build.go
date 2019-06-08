// Big part of code copied from github.com/swissChili/go-wasm
package build

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gascore/gas/gasx/cmd/compile"
	cfg "github.com/gascore/gas/gasx/cmd/config"
	"github.com/gascore/gas/gasx/cmd/css"
	"github.com/gascore/gas/gasx/cmd/lock"
	"github.com/gascore/gas/gasx/cmd/pm"
	"github.com/gascore/gas/gasx/utils"
	copyPkg "github.com/otiai10/copy"
	"github.com/spf13/cobra"
)

func Return() *cobra.Command {
	var cmdBuild = &cobra.Command{
		Use:   "build",
		Short: "gas projects builder",
		Run: func(cmd *cobra.Command, args []string) {
			allConfig, err := cfg.ParseConfig()
			utils.Must(err)

			utils.Must(Body(allConfig))
		},
	}

	return cmdBuild
}

func Body(allConfig *cfg.Config) error {
	gasLock, buildExternal, err := lock.ParseGasLock(allConfig)
	if err != nil {
		return err
	}

	if _, err := os.Stat("/dist"); os.IsNotExist(err) {
		err = os.MkdirAll("dist", os.ModePerm)
		if err != nil {
			return err
		}
	} else {
		err = os.RemoveAll("dist")
		if err != nil {
			return err
		}
	}

	if _, err := os.Stat("/css"); os.IsNotExist(err) {
		err = os.MkdirAll("css", os.ModePerm)
		if err != nil {
			return err
		}
	}

	info("Templates building")

	err = compile.Body(allConfig, gasLock, buildExternal)
	if err != nil {
		return err
	}

	err = css.ACSSBody(allConfig, gasLock, buildExternal)
	if err != nil {
		return err
	}

	info("Clearing ./dist")

	err = copyPkg.Copy("static", "dist")
	if err != nil {
		return err
	}

	info("Importing dependencies")
	if len(allConfig.Build.FilesDependencies) != 0 {
		for _, dep := range allConfig.Build.FilesDependencies {
			if _, err := os.Stat(dep.Path); os.IsExist(err) {
				continue
			}

			err = RunCommand(fmt.Sprintf("cp %s %s", dep.Src, dep.Path))
			if err != nil {
				return err
			}
		}
	}

	info("Compiling code")

	switch allConfig.Build.Platform {
	case "gopherjs":
		err = RunCommand(fmt.Sprintf("gopherjs build -o ./dist/index.js"))
		if err != nil {
			return err
		}
	case "wasm":
		err = RunCommand(fmt.Sprintf("GOOS=js GOARCH=wasm go build -o dist/main.wasm"))
		if err != nil {
			return err
		}
	case "tinygo":
		return errors.New("comming soon(?)")
	default:
		return errors.New("invalid target platform")
	}

	info("Compiling css")

	err = compileCSS(allConfig.Build)
	if err != nil {
		return err
	}

	info("Bundling web_modules")

	err = pm.BuildBody(allConfig)
	if err != nil {
		return err
	}

	if allConfig.Build.Platform == "wasm" {
		err = ioutil.WriteFile("dist/index.js", []byte(execScript+wasmLoader), 0644)
		if err != nil {
			return err
		}
	}

	err = gasLock.Save()
	if err != nil {
		return err
	}

	err = allConfig.Save()
	if err != nil {
		return err
	}

	info("Building finished")
	return nil
}

var info = utils.Info

func RunCommand(command string) error {
	cmd := exec.Command("sh", "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func compileCSS(cfg cfg.Build) error {
	return filepath.Walk("css", func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			f := strings.Split(path, "/")
			file := f[len(f)-1]
			e := strings.Split(file, ".")
			extension := e[len(e)-1]
			command := ""
			comp := ""

			if extension == "sass" {
				comp = cfg.Sass
			} else if extension == "scss" {
				comp = cfg.Sass
			} else if extension == "less" {
				comp = cfg.Less
			}

			if comp != "" {
				file = file[:len(file)-5]
				command = strings.Replace(comp, "INPUT", path, -1)
				command = strings.Replace(command, "OUTPUT", "dist/"+file+".css", -1)
			}

			err = RunCommand(command)

			return err
		}
		return nil
	})
}
