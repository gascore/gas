package css

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/gascore/gas/gasx/cmd/compile"
	cfg "github.com/gascore/gas/gasx/cmd/config"
	"github.com/gascore/gas/gasx/cmd/lock"
	"github.com/gascore/gas/gasx/utils"
	"github.com/spf13/cobra"
	"golang.org/x/net/html"
)

func ACSS() *cobra.Command {
	var cmdCssMinify = &cobra.Command{
		Use:   "acss",
		Short: "Generate a.css for elements classes styles",
		Long:  "Generate a.css for elements classes styles. \nFor more information see https://github.com/gascore/gas/gasx/blobl/master/acss.md",
		Run: func(cmd *cobra.Command, args []string) {
			allConfig, err := cfg.ParseConfig()
			utils.Must(err)

			gasLock, buildExternal, err := lock.ParseGasLock(allConfig)
			utils.Must(err)

			utils.Must(ACSSBody(allConfig, gasLock, buildExternal))

			utils.Must(gasLock.Save())
		},
	}

	return cmdCssMinify
}

func ACSSBody(allConfig *cfg.Config, gasLock *lock.Lock, buildExternal bool) error {
	config := allConfig.ACSS

	if config.BreakPoints == nil {
		config.BreakPoints = make(map[string]string)
	}
	if config.Custom == nil {
		config.Custom = make(map[string]string)
	}

	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	templates, err := compile.GetGasFiles(currentDir, buildExternal)
	if err != nil {
		return err
	}

	var generated = make(map[string]bool)
	var out string
	externalClasses := make(map[string]string)

	for _, fileInfo := range templates {
		file, err := os.Open(fileInfo.Path)
		if err != nil {
			return err
		}

		nodes, err := compile.Parse(file)
		if err != nil {
			return err
		}

		template, _, err := compile.ParseTemplate(filepath.Dir(fileInfo.Path), nodes)
		if err != nil {
			return err
		}

		classes, err := getClasses(template)
		if err != nil {
			return err
		}

		for _, class := range classes {
			if generated[class] {
				continue
			}

			generated[class] = true

			if inArray(class, config.Exceptions) {
				continue
			}

			classOut := buildClass(class, config)
			if classOut == "" {
				continue
			}

			out = out + classOut

			if buildExternal && strings.Index(fileInfo.Path, currentDir) == -1 {
				externalClasses[class] = classOut
			}
		}
	}

	if buildExternal { // save for ignoring
		gasLock.AtomicStyles = externalClasses
	} else if gasLock.AtomicStyles != nil {
		for _, classOut := range gasLock.AtomicStyles {
			out = out + classOut
		}
	}

	err = ioutil.WriteFile(config.Out, []byte(out), os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func buildClass(class string, config cfg.ACSS) string {
	rawClass := class

	var pOperator string
	if strings.Contains(class, ":") {
		pIndex := strings.Index(class, ":")
		pOperator = class[1+pIndex : 2+pIndex]
		class = class[:pIndex] + class[2+pIndex:]

		switch pOperator {
		case "a":
			pOperator = "active"
			break
		case "c":
			pOperator = "checked"
			break
		case "f":
			pOperator = "focus"
			break
		case "h":
			pOperator = "hover"
			break
		case "d":
			pOperator = "disabled"
			break
		}
	}

	var breakPoint string
	if strings.Contains(class, "--") {
		breakPoint = config.BreakPoints[class[strings.Index(class, "--")+len("--"):]]
		class = class[:strings.Index(class, "--")]
	}

	styleValue := GenerateStyleForClass(class, config.Custom)
	if len(styleValue) == 0 {
		return ""
	}

	styleOut := "{\n" + styleValue + "\n}\n"
	styleClass := "." + ClearClass(rawClass)
	if len(pOperator) == 0 {
		styleOut = styleClass + styleOut
	} else {
		styleOut = styleClass + ":" + pOperator + styleOut
	}

	if len(breakPoint) != 0 {
		styleOut = breakPoint + " {\n" + styleOut + "}\n"
	}

	return styleOut
}

func getClasses(t *html.Node) ([]string, error) {
	var classes []string
	var generated = make(map[string]bool)
	for c := t.FirstChild; c != nil; c = c.NextSibling {
		for _, attr := range c.Attr {
			if attr.Key != "class" {
				continue
			}

			for _, class := range strings.Split(attr.Val, " ") {
				if generated[class] {
					continue
				}

				classes = append(classes, class)

				generated[class] = true
			}
		}

		cClasses, err := getClasses(c)
		if err != nil {
			return nil, err
		}

		classes = append(classes, cClasses...)
	}

	return classes, nil
}

func ClearClass(class string) string {
	class = strings.Replace(class, "(", "\\(", -1)
	class = strings.Replace(class, ")", "\\)", -1)
	class = strings.Replace(class, ",", "\\,", -1)
	class = strings.Replace(class, ":", "\\:", -1)
	class = strings.Replace(class, "#", "\\#", -1)
	class = strings.Replace(class, ".", "\\.", -1)
	class = strings.Replace(class, "%", "\\%", -1)
	class = strings.Replace(class, "!", "\\!", -1)

	return class
}

func inArray(el string, arr []string) bool {
	for _, x := range arr {
		if x == el {
			return true
		}
	}

	return false
}
