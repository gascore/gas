package pm

import (
	cfg "github.com/gascore/gas/gasx/cmd/config"
	"github.com/gascore/gas/gasx/utils"
	"github.com/spf13/cobra"
)

const packagesFolder = "web_modules"
const getInfoURL = "https://data.jsdelivr.com/v1/package/npm/"
const getFileURL = "https://cdn.jsdelivr.net/npm/"

func Return() *cobra.Command {
	var isReset bool

	var cmdPMMain = &cobra.Command{
		Use:   "pm [command]",
		Short: "Update all dependencies",
		Run: func(cmd *cobra.Command, args []string) {
			allConfig, err := cfg.ParseConfig()
			utils.Must(err)

			for _, dep := range allConfig.Deps.Deps {
				// if not reset and package already downloaded just skip
				if !isReset && utils.Exists(packagesFolder+"/"+dep.Name) {
					continue
				}

				if len(dep.DefaultFile) != 0 {
					fileBody, err := getFile(dep, dep.DefaultFile)
					utils.Must(err)
					utils.Must(savePackage(dep, fileBody))
				}

				for _, file := range dep.RequiredFiles {
					fileBody, err := getFile(dep, file)
					utils.Must(err)
					utils.Must(savePackage(dep, fileBody))
				}
			}
		},
	}

	cmdPMMain.Flags().BoolVarP(&isReset, "reset", "r", false, "delete and download all packages")

	return cmdPMMain
}
