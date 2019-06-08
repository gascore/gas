package new

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/gascore/gas/gasx/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	var cmdCompile = &cobra.Command{
		Use:   "new [template] [project name]",
		Short: "Create new gas project",
		Long: `Create new gas project.

Support three templates:

1. default - default template (only gas and gas-web)
2. router  - default template with gas-router
3. full    - template with all *gas technologies*: gas-router, gas-store
`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("new requires at least 1 argument (project name)")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			var pName string // project name
			var pType string // project type

			if len(args) == 1 { // if only one argument will create standard project with name = argument
				pName = args[0]
				pType = "default"
			} else {
				pType = args[0]
				pName = args[1]
			}

			extraFiles := make(map[string]string)
			switch pType {
			case "router":
				extraFiles["routes.go"] = routesGo
				extraFiles["components/about.gas"] = routerAboutGas
				extraFiles["components/hello.gas"] = routerHelloGas
				extraFiles["main.gas"] = routerMainGas
			case "full":
				extraFiles["store/store.go"] = fullStoreGo
				extraFiles["routes.go"] = routesGo
				extraFiles["components/about.gas"] = fullAboutGas
				extraFiles["components/hello.gas"] = fullHelloGas
				extraFiles["main.gas"] = fullMainGas
			case "default":
				extraFiles["main.gas"] = defaultMainGas
				extraFiles["components/hello.gas"] = defaultHelloGas
			}

			currentDir, err := os.Getwd()
			utils.Must(err)

			i := strings.Index(currentDir, "/go/src/")
			if i > 0 {
				currentDir = currentDir[i+len("/go/src/"):] + "/" + pName
			} else {
				currentDir = "your_project_path"
			}

			utils.Must(InitFiles(currentDir, pName, map[string]string{
				"static/index.html":      indexHtml,
				"static/index.gojs.html": indexGoJsHtml,
				"clear.sh":               clearSh,
				"config.json":            configJSON,
				".gitignore":             gitIgnore,
			}))
			utils.Must(InitFiles(currentDir, pName, extraFiles))

			red := color.New(color.FgRed).SprintFunc()
			gray := color.New(color.FgHiBlue).SprintFunc()

			fmt.Println("> cd ./" + pName)

			if currentDir == "your_project_path" {
				fmt.Println(fmt.Sprintf(`> find ./ -type f \( -iname \*.gas -o -iname \*.go \) -exec sed -i 's/your_project_path/%s/g' {} +`, red("<YOUR_PROJECT_PATH_HERE>")))
				fmt.Println(gray("(for dynamic imports use YOUR_PROJECT_PATH_HERE=\".\", and \"..\" in components/*.gas)"))
			}

			fmt.Println("> gasx run")
		},
	}

	return cmdCompile
}

func InitFiles(currentDir, namesPrefix string, files map[string]string) error {
	haveDir := make(map[string]bool)
	for fileName, fileBody := range files {
		fileName = namesPrefix + "/" + fileName
		dirName := filepath.Dir(fileName)
		if !haveDir[dirName] {
			if _, err := os.Stat(dirName); err != nil {
				err := os.MkdirAll(dirName, os.ModePerm)
				if err != nil {
					panic(err)
				}
			}

			haveDir[dirName] = true
		}

		file, err := os.Create(fileName)
		if err != nil {
			return err
		}

		fileBody = strings.Replace(fileBody, "your_project_path", currentDir, -1)

		_, err = file.WriteString(fileBody)
		if err != nil {
			return err
		}
	}

	return nil
}
