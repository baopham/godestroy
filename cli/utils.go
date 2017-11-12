package cli

import (
	"github.com/olebedev/when"
	"github.com/olebedev/when/rules/common"
	"github.com/olebedev/when/rules/en"
)

func getTimeParser() *when.Parser {
	w := when.New(nil)
	w.Add(en.All...)
	w.Add(common.All...)
	return w
}
