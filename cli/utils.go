package cli

import (
	"github.com/olebedev/when"
	"github.com/olebedev/when/rules/common"
	"github.com/olebedev/when/rules/en"
	"os"
)

func getTimeParser() *when.Parser {
	w := when.New(nil)
	w.Add(en.All...)
	w.Add(common.All...)
	return w
}

func isValidFile(path string) (bool, os.FileInfo, error) {
	info, err := os.Stat(path)
	if os.IsNotExist(err) || info.IsDir() {
		return false, info, nil
	}
	if err != nil {
		return false, info, err
	}
	return true, info, nil
}
