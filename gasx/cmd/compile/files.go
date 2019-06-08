package compile

import (
	"go/build"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	cfg "github.com/gascore/gas/gasx/cmd/config"
	"github.com/gascore/gas/gasx/utils"
	"github.com/visualfc/fastmod"
)

var (
	reImport  = regexp.MustCompile(`(?m)(.?)import(|\s)*?"(?P<import>((.|\n|\r)*?))"(.|\n|\r)?`)
	reImports = regexp.MustCompile(`(?m)(.?)import(|\s)*?\((?P<name>((.|\n|\r)*?))\)(.|\n|\r)?`)
	rePaths   = regexp.MustCompile(`.*?"(?P<import>(.*?))".*?`)

	goPath = (func() string {
		goPath := os.Getenv("GOPATH")
		if goPath == "" {
			goPath = build.Default.GOPATH
		}
		return goPath
	})()
)

func GetGasFiles(currentDir string, buildExternal bool, config *cfg.Config) ([]File, error) {
	already := make(map[string]bool)
	files, err := getGasFilesBody(currentDir, buildExternal, already, false, config)
	if err != nil {
		return nil, err
	}

	if !config.GoModSupport || !buildExternal {
		return files, nil
	}

	pkg, err := fastmod.LoadPackage(currentDir, &build.Default)
	if err != nil {
		return files, nil
	}

	for _, nodeValue := range pkg.NodeMap {
		modDir := nodeValue.ModDir()
		newFiles, err := getGasFilesBody(modDir, false, already, true, config)
		if err != nil {
			return nil, err
		}

		files = append(files, newFiles...)
	}

	return files, nil
}

func getGasFilesBody(root string, buildExternal bool, already map[string]bool, isExternal bool, allConfig *cfg.Config) ([]File, error) {
	var files []File
	if _, ok := already[root]; ok {
		return files, nil
	}
	already[root] = true

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		var forImports []byte

		if strings.HasSuffix(path, ".gas") {
			files = append(files, File{Path: path, IsExternal: isExternal})

			if buildExternal && !allConfig.GoModSupport {
				file, err := ioutil.ReadFile(path)
				if err != nil {
					return nil
				}

				forImports = file
			}
		} else if !allConfig.GoModSupport && buildExternal &&
			(strings.HasSuffix(path, ".go") || strings.HasSuffix(path, ".gos")) &&
			!strings.HasSuffix(path, allConfig.Compile.FilesSuffix+".go") &&
			!strings.HasSuffix(path, allConfig.Compile.ExternalFilesSuffix+".go") {

			file, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			forImports = file
		}

		if len(forImports) == 0 || allConfig.GoModSupport {
			return nil
		}

		aloneImport := reImport.Find(forImports)
		if len(aloneImport) != 0 {
			aloneImportS := string(aloneImport)
			aloneImportS = strings.Replace(aloneImportS, " ", "", -1)
			aloneImportS = aloneImportS[len(`import"`):]
			aloneImportS = aloneImportS[:len(aloneImportS)-2]

			fromDepth, err := goDepth(aloneImportS, goPath, root, already, allConfig)
			if err != nil {
				return err
			}

			files = append(files, fromDepth...)
		}

		rawImports := reImports.Find(forImports)
		if len(rawImports) == 0 {
			return nil
		}

		for _, item := range rePaths.FindAllString(string(rawImports), -1) {
			item = item[strings.Index(item, `"`):]
			item = item[:len(item)-strings.Index(item, `"`)-1]
			item = item[1:]

			fromDepth, err := goDepth(item, goPath, root, already, allConfig)
			if err != nil {
				return err
			}

			files = append(files, fromDepth...)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return files, err
}

func goDepth(item string, goPath string, root string, already map[string]bool, config *cfg.Config) ([]File, error) {
	newRoot := goPath + "/src/" + item

	if strings.HasPrefix(newRoot, root) {
		return []File{}, nil
	}

	if !utils.Exists(newRoot) {
		return []File{}, nil
	}

	newItems, err := getGasFilesBody(newRoot, true, already, true, config)
	if err != nil {
		return nil, err
	}

	return newItems, nil
}
