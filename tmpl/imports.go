package tmpl

import (
	"github.com/thecodedproject/gopkg"
)

// UnnamedImports returns a `[]gopkg.ImportAndAlias` from a list of import paths
func UnnamedImports(
	importPaths ...string,
) []gopkg.ImportAndAlias {

	ret := make([]gopkg.ImportAndAlias, 0, len(importPaths))
	for _, i := range importPaths {
		ret = append(ret, gopkg.ImportAndAlias{
			Import: i,
		})
	}
	return ret
}
