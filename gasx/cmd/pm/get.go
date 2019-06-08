package pm

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
	"time"

	cfg "github.com/gascore/gas/gasx/cmd/config"
	"github.com/gascore/gas/gasx/utils"
	"github.com/spf13/cobra"
	"github.com/valyala/fasthttp"
)

func Get() *cobra.Command {
	var version, defaultFile string

	var cmdPMGet = &cobra.Command{
		Use:   "get [package name]",
		Short: "Get new package",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				panic("no package name")
			}
			packageName := args[0]

			allConfig, err := cfg.ParseConfig()
			utils.Must(err)

			versions, err := getVersions(packageName)
			utils.Must(err)

			if version != "" && !inArray(version, versions.Versions) {
				panic("invalid package version")
			} else {
				version = versions.Tags.Latest
			}

			/* Get info about Files */
			packageFiles, err := getPackageFiles(packageName, version)
			utils.Must(err)

			if defaultFile == "" {
				if packageFiles.Default == "" {
					panic("no default file in package")
				}
				defaultFile = packageFiles.Default
			}

			/* Get default file */
			dep := cfg.Dep{
				Name:        packageName,
				Version:     version,
				DefaultFile: defaultFile,
			}

			fileBody, err := getFile(dep)
			utils.Must(err)

			allConfig.Deps.Deps = append(allConfig.Deps.Deps, dep)
			utils.Must(savePackage(dep, fileBody))
			utils.Must(allConfig.Save())
		},
	}

	cmdPMGet.Flags().StringVarP(&version, "version", "v", "", "package version")
	cmdPMGet.Flags().StringVarP(&defaultFile, "default", "d", "", "required file in package")

	return cmdPMGet
}

type Versions struct {
	Tags struct {
		Latest string `json:"latest"`
		Alpha  string `json:"alpha"`
		Beta   string `json:"beta"`
		Next   string `json:"next"`
	} `json:"tags"`
	Versions []string `json:"versions"`
}

type PackageFiles struct {
	Default string  `json:"default"`
	Files   []Files `json:"files"`
}

type Files struct {
	Type  string    `json:"type"`
	Name  string    `json:"name"`
	Files []Files   `json:"files,omitempty"`
	Hash  string    `json:"hash,omitempty"`
	Time  time.Time `json:"time,omitempty"`
	Size  int       `json:"size,omitempty"`
}

func inArray(a string, arr []string) bool {
	for _, el := range arr {
		if a == el {
			return true
		}
	}
	return false
}

func savePackage(dep cfg.Dep, fileBody []byte) error {
	if !utils.Exists(packagesFolder) {
		err := os.Mkdir(packagesFolder, os.ModePerm)
		if err != nil {
			return err
		}
	}

	depFolder := packagesFolder + "/" + dep.Name
	if utils.Exists(depFolder) {
		err := utils.RemoveContents(depFolder)
		if err != nil {
			return err
		}
	}

	err := os.MkdirAll(depFolder+"/"+path.Dir(dep.DefaultFile), os.ModePerm)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(depFolder+"/"+dep.DefaultFile, fileBody, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func getFile(dep cfg.Dep) ([]byte, error) {
	client := fasthttp.Client{}
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(getFileURL + dep.Name + "@" + dep.Version + dep.DefaultFile)
	res := fasthttp.AcquireResponse()

	err := client.Do(req, res)
	if err != nil {
		return nil, err
	}

	if res.StatusCode() != 200 {
		return nil, errors.New("invalid package name")
	}

	return res.Body(), nil
}
