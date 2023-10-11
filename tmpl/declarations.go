package tmpl

import (
	"github.com/thecodedproject/gopkg"
)

// FuncWithContextAndError returns a gopkg.DeclFunc with the given name, args
// and returnArgs, but with additionally taking a `context.Context` as the
// first argument and return an `error` as the last return argument.
//
// The `context.Context` argument will always have the name `ctx` and the `error`
// return argument will be either unnamed or named `err`, depending on whether
// returnArgs contains unnamed or named `gopkg.DeclVar`s
func FuncWithContextAndError(
	funcName string,
	args []gopkg.DeclVar,
	returnArgs []gopkg.DeclVar,
) gopkg.DeclFunc {

	f := gopkg.DeclFunc{
		Name: funcName,
	}

	f.Args = append(
		[]gopkg.DeclVar{
			{
				Name: "ctx",
				Type: gopkg.TypeNamed{
					Name:      "Context",
					Import:    "context",
					ValueType: gopkg.TypeInterface{},
				},
			},
		},
		args...,
	)

	errVarName := ""
	if len(returnArgs) > 0 && returnArgs[0].Name != "" {
		errVarName = "err"
	}

	f.ReturnArgs = append(
		returnArgs,
		gopkg.DeclVar{
			Name: errVarName,
			Type: gopkg.TypeError{},
		},
	)

	return f
}

func UnnamedReturnArgs(retArgs ...gopkg.Type) []gopkg.DeclVar {

	ret := make([]gopkg.DeclVar, 0, len(retArgs))

	for _, arg := range retArgs {
		ret = append(ret, gopkg.DeclVar{
			Type: arg,
		})
	}

	return ret
}
