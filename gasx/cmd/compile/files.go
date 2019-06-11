package compile

import (
	"go/build"
	"os"
	"path/filepath"
	"strings"

	"github.com/visualfc/fastmod"
)

func GetGasFiles(currentDir string, buildExternal bool) ([]File, error) {
	var already alreadyList = []string{}
	files, err := getGasFilesBody(currentDir, false) // get files from current dir
	if err != nil {
		return nil, err
	}

	if !buildExternal {
		return files, nil
	}

	pkg, err := fastmod.LoadPackage(currentDir, &build.Default)
	if err != nil {
		return files, nil
	}

	for _, nodeValue := range pkg.NodeMap {
		err = parseModDir(nodeValue.ModDir(), files, already)
		if err != nil {
			return files, err
		}
	}

	return files, nil
}

func parseModDir(root string, files []File, already alreadyList) error {
	if already.in(root) {
		return nil
	}

	newFiles, err := getGasFilesBody(root, true)
	if err != nil {
		return nil
	}

	files = append(files, newFiles...)
	already = append(already, root)

	pkg, err := fastmod.LoadPackage(root, &build.Default)
	if err != nil {
		return err
	}

	for _, nodeValue := range pkg.NodeMap {
		err = parseModDir(nodeValue.ModDir(), files, already)
		if err != nil {
			return err
		}
	}

	return nil
}

func getGasFilesBody(root string, isExternal bool) ([]File, error) {
	var files []File

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".gas") {
			files = append(files, File{Path: path, IsExternal: isExternal})
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return files, err
}


type alreadyList []string

func (al alreadyList) in(a string) bool {
	for _, el := range al {
		if el == a {
			return true
		}
	}
	return false
}
