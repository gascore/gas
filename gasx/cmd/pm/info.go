package pm

import (
	"encoding/json"
	"fmt"

	"github.com/gascore/gas/gasx/utils"
	"github.com/go-yaml/yaml"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/valyala/fasthttp"
)

func Info() *cobra.Command {
	var version, outputType string

	var cmdPMGet = &cobra.Command{
		Use:   "info [package name]",
		Short: "Print information about package",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				panic("no package name")
			}
			packageName := args[0]

			if outputType != "json" && outputType != "yaml" {
				panic("invalid output type")
			}

			versions, err := getVersions(packageName)
			utils.Must(err)

			if version != "" && !inArray(version, versions.Versions) {
				panic("invalid package version")
			} else {
				version = versions.Tags.Latest
			}

			packageFiles, err := getPackageFiles(packageName, version)
			utils.Must(err)

			info := packageInfo{
				Versions:      versions,
				PackageFiles:  packageFiles,
				LatestVersion: versions.Tags.Latest,
				DefaultFile:   packageFiles.Default,
			}

			var out []byte
			switch outputType {
			case "yaml":
				out, err = yaml.Marshal(info)
				utils.Must(err)
			case "json":
				outUgly, err := json.Marshal(info)
				utils.Must(err)

				out, err = utils.PrettyPrint(outUgly)
				utils.Must(err)
			}

			fmt.Println(string(out))
		},
	}

	cmdPMGet.Flags().StringVarP(&outputType, "type", "t", "json", "output format type")
	cmdPMGet.Flags().StringVarP(&version, "version", "v", "", "package version")

	return cmdPMGet
}

type packageInfo struct {
	Versions     Versions     `json:"versions" yaml:"versions"`
	PackageFiles PackageFiles `json:"files" yaml:"files"`

	LatestVersion string `json:"latest_version" yaml:"latest_version"`
	DefaultFile   string `json:"default_file" yaml:"default_file"`
}

// getVersions Get info about package versions
func getVersions(packageName string) (Versions, error) {
	client := fasthttp.Client{}
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(getInfoURL + packageName)
	res := fasthttp.AcquireResponse()
	err := client.Do(req, res)
	if err != nil {
		return Versions{}, err
	}

	if res.StatusCode() != 200 {
		return Versions{}, errors.New("invalid package name")
	}

	var versions Versions
	err = json.Unmarshal(res.Body(), &versions)
	if err != nil {
		return Versions{}, nil
	}

	return versions, nil
}

// getPackageFiles Get info about package files
func getPackageFiles(packageName, version string) (PackageFiles, error) {
	client := fasthttp.Client{}
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(getInfoURL + packageName + "@" + version)
	res := fasthttp.AcquireResponse()
	err := client.Do(req, res)
	if err != nil {
		return PackageFiles{}, err
	}

	if res.StatusCode() != 200 {
		return PackageFiles{}, errors.New("invalid package name")
	}

	var packageFiles PackageFiles
	err = json.Unmarshal(res.Body(), &packageFiles)
	if err != nil {
		return PackageFiles{}, nil
	}

	return packageFiles, nil
}
