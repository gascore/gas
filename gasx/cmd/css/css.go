package css

import (
	"fmt"
	"github.com/spf13/cobra"
)

func Return() *cobra.Command {
	var cmdCssMain = &cobra.Command{
		Use:   "css [command]",
		Short: "Collection of small commands helps working with styles",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Choose your command")
		},
	}

	return cmdCssMain
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
