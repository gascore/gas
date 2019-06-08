package serve

import (
	"fmt"
	"net/http"
	"os"

	"github.com/fatih/color"
	cfg "github.com/gascore/gas/gasx/cmd/config"
	"github.com/gascore/gas/gasx/utils"
	"github.com/spf13/cobra"
)

func Serve() *cobra.Command {
	var cmdServe = &cobra.Command{
		Use:   "serve",
		Short: "Serve current directory",
		Run: func(cmd *cobra.Command, args []string) {
			allConfig, err := cfg.ParseConfig()
			utils.Must(err)

			utils.Must(Body(allConfig.Serve))
		},
	}

	return cmdServe
}

func Body(config cfg.Serve) error {
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	utils.Info(fmt.Sprintf("Server starting on port: %s", color.GreenString(config.Port[1:])))

	err = http.ListenAndServe(config.Port, http.FileServer(http.Dir(currentDir+"/"+config.Dir)))
	if err != nil {
		return err
	}

	return nil
}
