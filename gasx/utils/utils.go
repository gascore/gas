package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
)

func Must(err error) {
	if err != nil {
		LogError(err)
		os.Exit(1)
	}
}

func LogError(err error) {
	fmt.Println(color.RedString("ERROR") + ": " + err.Error())
}

func Info(msg string) {
	fmt.Println(color.BlueString("Builder:") + " " + msg)
}

func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func PrettyPrint(b []byte) ([]byte, error) {
	var out bytes.Buffer
	err := json.Indent(&out, b, "", "  ")
	return out.Bytes(), err
}

func RemoveContents(dir string) error {
	files, err := filepath.Glob(filepath.Join(dir, "*"))
	if err != nil {
		return err
	}
	for _, file := range files {
		err = os.RemoveAll(file)
		if err != nil {
			return err
		}
	}
	return nil
}
