package css

import (
	"bytes"
	"os"
	"strings"

	"github.com/gascore/gas/gasx/utils"
	"github.com/spf13/cobra"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
)

func Minify() *cobra.Command {
	var inputName string
	var outputName string

	var cmdCssMinify = &cobra.Command{
		Use:   "minify",
		Short: "Minify css files",
		Long:  `Minify css files by tdewolff/minify(https://github.com/tdewolff/minify)`,
		Run: func(cmd *cobra.Command, args []string) {
			m := minify.New()
			m.AddFunc("text/css", css.Minify)

			file, err := os.Open(inputName)
			utils.Must(err)

			buf := bytes.NewBufferString("")
			utils.Must(m.Minify("text/css", buf, file))

			minified := buf.String()

			if inputName != "main.css" && outputName == "main.min.css" {
				outputName = strings.TrimSuffix(inputName, ".css") + ".min.css"
			}

			minifiedFile, err := os.Create(outputName)
			utils.Must(err)

			_, err = minifiedFile.WriteString(minified)
			utils.Must(err)
		},
	}

	cmdCssMinify.Flags().StringVarP(&inputName, "input", "i", "main.css", "input file name (main.css by default)")
	cmdCssMinify.Flags().StringVarP(&outputName, "output", "o", "main.min.css", "output file name (main.min.css by default)")

	return cmdCssMinify
}
