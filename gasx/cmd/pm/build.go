package pm

import (
	"io/ioutil"
	"os"
	"strings"

	cfg "github.com/gascore/gas/gasx/cmd/config"
	"github.com/gascore/gas/gasx/utils"
	"github.com/spf13/cobra"
)

func Build() *cobra.Command {
	var cmdPMGet = &cobra.Command{
		Use:   "build",
		Short: "Build all dependencies to two file (js and css)",
		Run: func(cmd *cobra.Command, args []string) {
			allConfig, err := cfg.ParseConfig()
			utils.Must(err)

			utils.Must(BuildBody(allConfig))
		},
	}

	return cmdPMGet
}

func BuildBody(allConfig *cfg.Config) error {
	var outJS, outCSS string
	for _, dep := range allConfig.Deps.Deps {
		file, err := ioutil.ReadFile(packagesFolder + "/" + dep.Name + dep.DefaultFile)
		if err != nil {
			return err
		}

		if strings.HasSuffix(dep.DefaultFile, ".js") {
			outJS = outJS + "\n\n" + string(file)
		} else if strings.HasSuffix(dep.DefaultFile, ".css") {
			outCSS = outCSS + "\n\n" + string(file)
		} else {
			panic("invalid package default file type")
		}
	}

	err := saveToFile(outJS, "js", allConfig.Deps.BuildJSOut)
	if err != nil {
		return err
	}

	err = saveToFile(outCSS, "css", allConfig.Deps.BuildCSSOut)
	if err != nil {
		return err
	}

	return nil
}

func saveToFile(out string, t string, fileName string) error {
	if fileName == "" {
		fileName = "deps." + t
	}

	if utils.Exists(fileName) {
		err := os.Remove(fileName)
		if err != nil {
			return err
		}
	}

	return ioutil.WriteFile(fileName, []byte(out), os.ModePerm)
}
