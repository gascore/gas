package cmd

import (
	"fmt"
	"os"

	"github.com/gascore/gas/gasx/cmd/build"
	"github.com/gascore/gas/gasx/cmd/compile"
	"github.com/gascore/gas/gasx/cmd/css"
	commandNew "github.com/gascore/gas/gasx/cmd/new"
	"github.com/gascore/gas/gasx/cmd/pm"
	"github.com/gascore/gas/gasx/cmd/serve"
	"github.com/gascore/gas/gasx/cmd/watcher"
	"github.com/gascore/gas/gasx/utils"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "gasx",
	Short:   "CLI for gas apps",
	Long:    `CLI for gas apps. Built with https://github.com/spf13/cobra`,
	Version: "0.0...1",
	Run: func(cmd *cobra.Command, args []string) {
		utils.Must(cmd.Help())
	},
}

func Execute() {
	cmdCss := css.Return()
	cmdPM := pm.Return()

	rootCmd.AddCommand(build.Return())
	rootCmd.AddCommand(compile.Return())
	rootCmd.AddCommand(commandNew.New())
	rootCmd.AddCommand(serve.Serve())
	rootCmd.AddCommand(watcher.Watch())
	rootCmd.AddCommand(watcher.RunAlias())
	rootCmd.AddCommand(cmdCss)
	rootCmd.AddCommand(cmdPM)

	cmdCss.AddCommand(css.Minify())
	cmdCss.AddCommand(css.ACSS())

	cmdPM.AddCommand(pm.Get())
	cmdPM.AddCommand(pm.Build())
	cmdPM.AddCommand(pm.Info())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
