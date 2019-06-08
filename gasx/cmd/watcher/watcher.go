package watcher

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/gascore/gas/gasx/cmd/build"
	cfg "github.com/gascore/gas/gasx/cmd/config"
	"github.com/gascore/gas/gasx/cmd/serve"
	"github.com/gascore/gas/gasx/utils"
	"github.com/pkg/errors"
	"github.com/radovskyb/watcher"
	"github.com/spf13/cobra"
)

func Watch() *cobra.Command {
	return &cobra.Command{
		Use:   "watch",
		Short: "Watching for files changes",
		Run: func(cmd *cobra.Command, args []string) {
			allConfig, err := cfg.ParseConfig()
			utils.Must(err)

			config := allConfig.Watch

			w := watcher.New()

			w.SetMaxEvents(1)

			if config.IgnoreCompiled {
				w.AddFilterHook(func(info os.FileInfo, fullPath string) error {
					if strings.HasSuffix(fullPath, "_gas.go") {
						return watcher.ErrSkip
					}
					return nil
				})
			}

			for _, watcherInfo := range config.Watchers {
				if !utils.Exists(watcherInfo.Name) {
					continue
				}

				if watcherInfo.IsRecursive {
					utils.Must(w.AddRecursive(watcherInfo.Name))
				} else {
					utils.Must(w.Add(watcherInfo.Name))
				}
			}

			// utils.Must(w.Add("."))

			go func() {
				for {
					select {
					case event := <-w.Event:
						fmt.Println(color.GreenString("Triggered: "), event.Name())
						err := execTasks(config.Tasks, allConfig)
						if err != nil {
							utils.LogError(err)
							continue
						}
					case err := <-w.Error:
						log.Fatalln(err)
					case <-w.Closed:
						return
					}
				}
			}()

			color.Green("Starting")

			err = execTasks(config.InitTasks, allConfig)
			if err != nil {
				utils.LogError(err)
				return
			}

			time.Sleep(time.Second * 8)

			// Start the watching process - it'll check for changes every 100ms.
			if err := w.Start(time.Millisecond * 1000); err != nil {
				color.Red("Fatal error")
				log.Fatalln(err)
			}
		},
	}
}

func execTasks(tasks []*cfg.WatchersTask, config *cfg.Config) error {
	for _, task := range tasks {
		if task.IsGas {
			switch task.Command {
			case "build":
				err := build.Body(config)
				if err != nil {
					return err
				}
			case "serve":
				go func() {
					utils.LogError(serve.Body(config.Serve))
				}()
			default:
				return errors.New("undefined gas command")
			}
			continue
		}

		if task.NoWait {
			go func() {
				parts := strings.Fields(task.Command)
				_, err := exec.Command(parts[0], parts[1:]...).Output()
				if err != nil {
					if _, ok := err.(*exec.ExitError); ok {
						return
					}

					utils.LogError(err)
				}
			}()
			continue
		}

		parts := strings.Fields(task.Command)
		_, err := exec.Command(parts[0], parts[1:]...).Output()
		if err != nil {
			return err
		}
	}

	return nil
}

func RunAlias() *cobra.Command {
	runCommand := Watch()
	runCommand.Use = "run"
	runCommand.Short = "gasx watch alias"
	return runCommand
}
