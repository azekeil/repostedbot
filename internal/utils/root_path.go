package utils

import (
	"os"
)

var rootPath string

func GetRootPath() string {
	if rootPath == "" {
		var err error
		rootPath, err = os.Executable()
		if err != nil {
			panic(err)
		}
		// _, b, _, _ := runtime.Caller(0)
		// basepath := filepath.Dir(b)
		// rootPath = filepath.Join(basepath, "../../")
	}
	return rootPath
}
