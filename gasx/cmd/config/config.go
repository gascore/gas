package config

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/gascore/gas/gasx/utils"
)

type ACSS struct {
	BreakPoints map[string]string `json:"breakPoints"`
	Exceptions  []string          `json:"exceptions"` // ignoring classes
	Custom      map[string]string `json:"custom"`     // custom values for styles
	Out         string            `json:"out"`        // output file name (a.css by default)"
}

type Build struct {
	Platform string `json:"platform"` // gopherjs or wasm

	Sass string `json:"sass"`
	Scss string `json:"scss"`
	Less string `json:"less"`

	FilesDependencies []FileDep `json:"files_deps"`
}

type FileDep struct {
	Path string
	Src  string
}

type Compile struct {
	SupportStyles bool   `json:"supportStyles"` // extract <style> to one file (main.css by default)
	StylesOut     string `json:"stylesOut"`     // custom file name for SupportStyles output

	FilesSuffix         string `json:"suffix"`          // {filename}{suffix}.go
	ExternalFilesSuffix string `json:"external_suffix"` // in external directory: {filename}{externalSuffix}.go
}

// Job represents a job.
type WatchersJob struct {
	Watcher *Watcher
	Message string
}

// Context represents a context of a process.
type WatchersContext struct {
	Wd       string
	Config   *WatchersConfig
	Interval int
}

// Config represents a configuration of a process.
type WatchersConfig struct {
	InitTasks      []*WatchersTask `json:"init_tasks"`
	Watchers       []*Watcher      `json:"watchers"`
	Tasks          []*WatchersTask `json:"tasks"`
	IgnoreCompiled bool            `json:"ignore_compiled"`
}

// Watcher represents a file watcher.
type Watcher struct {
	Name        string `json:"name"`
	IsRecursive bool   `json:"recursive"`
}

// A Task represents a task.
type WatchersTask struct {
	Command string `json:"command"`
	NoWait  bool   `json:"nowait"`
	IsGas   bool   `json:"is_gas"`
}

type Serve struct {
	Port string `json:"port"`
	Dir  string `json:"dir"`
}

type Config struct {
	IgnoreExternal bool `json:"ignore_external"` // build gas files in external libraries and don't use .gaslock
	GoModSupport   bool `json:"go_mod_support"`

	ACSS    ACSS           `json:"acss"`         // for gasx css acss
	Compile Compile        `json:"compile"`      // for gasx compile
	Watch   WatchersConfig `json:"watch"`        // for gasx watch
	Build   Build          `json:"build"`        // for gasx build
	Serve   Serve          `json:"serve"`        // for gasx serve
	Deps    Dependencies   `json:"dependencies"` // web dependencies
}

type Dependencies struct {
	BuildJSOut  string `json:"js_out"`
	BuildCSSOut string `json:"css_out"`
	Deps        []Dep  `json:"deps"`
}

type Dep struct {
	Name          string   `json:"name"`
	Version       string   `json:"version"`
	DefaultFile   string   `json:"file"`
	RequiredFiles []string `json:"files"`
}

func ParseConfig() (*Config, error) {
	configFile, err := ioutil.ReadFile("config.json")
	if err != nil {
		if os.IsExist(err) {
			return nil, err
		} else {
			configFile = []byte(defaultConfig)
		}
	}

	var config Config
	err = json.Unmarshal(configFile, &config)
	if err != nil {
		return nil, err
	}

	if config.Serve.Port == "" {
		config.Serve.Port = ":8080"
	}

	if config.ACSS.Out == "" {
		config.ACSS.Out = "a.css"
	}

	if config.Compile.FilesSuffix == "" {
		config.Compile.FilesSuffix = "_gas"
	}

	if config.Compile.ExternalFilesSuffix == "" {
		config.Compile.ExternalFilesSuffix = "_gas_e"
	}

	return &config, nil
}

func (config *Config) Save() error {
	configJSON, err := json.Marshal(&config)
	if err != nil {
		return err
	}

	configJSON, err = utils.PrettyPrint(configJSON)
	if err != nil {
		return err
	}

	if utils.Exists("config.json") {
		err := os.Remove("config.json")
		if err != nil {
			return err
		}
	}

	configFile, err := os.Create("config.json")
	if err != nil {
		return err
	}

	_, err = configFile.Write(configJSON)
	if err != nil {
		return err
	}

	return nil
}
