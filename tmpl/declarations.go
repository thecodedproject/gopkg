package tmpl

import (
	"github.com/neurotempest/gopkg"
)

func UnnamedReturnArgs(retArgs ...gopkg.Type) []gopkg.DeclVar {

	ret := make([]gopkg.DeclVar, 0, len(retArgs))

	for _, arg := range retArgs {
		ret = append(ret, gopkg.DeclVar{
			Type: arg,
		})
	}

	return ret
}
