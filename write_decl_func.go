package gopkg

import (
	"io"
	"text/template"
)

func WriteDeclFunc(
	w io.Writer,
	decl DeclFunc,
	importAliases map[string]string,
) error {

	w.Write([]byte(fullFuncDecl(decl, importAliases)))
	w.Write([]byte(" {\n"))

	if decl.BodyTmpl != "" {

		tmpl := funcBaseTemplate(decl, importAliases)

		funcTmpl, err := tmpl.Parse(decl.BodyTmpl)
		if err != nil {
			return err
		}

		data := struct{
			Func DeclFunc
		}{
			Func: decl,
		}

		err = funcTmpl.Execute(w, data)
		if err != nil {
			return err
		}
	}

	w.Write([]byte("}\n"))

	return nil
}

func fullFuncDecl(f DeclFunc, importAliases map[string]string) string {

	decl := "func "

	if f.Receiver.TypeName != "" {
		decl += "("

		if f.Receiver.VarName != "" {
			decl += f.Receiver.VarName + " "
		}

		if f.Receiver.IsPointer {
			decl += "*"
		}

		decl += f.Receiver.TypeName + ") "
	}

	decl += f.Name + funcArgsAndRetArgs(f, importAliases, true)

	return decl
}

// funcArgsAndRetArgs is a helper function which returns a function signature
// (i.e. argument list and return argment list) without the `func` specifier or
// function name.
//
// `addNewLinesToArgsList` will new-line seperate the arguments list (iff there
// is more than 1 argument in the args list)
func funcArgsAndRetArgs(
	f DeclFunc,
	importAliases map[string]string,
	addNewLinesToArgsList bool,
) string {

	decl := "("

	// Don't add new lines to arg list if there are 0 args or only 1 arg
	if len(f.Args) < 2 {
		addNewLinesToArgsList = false
	}

	if addNewLinesToArgsList {
		decl += "\n"
	}

	for iArg, arg := range f.Args {

		if addNewLinesToArgsList {
			decl += "\t"
		}

		decl += arg.Name + " " + arg.FullType(importAliases)

		if addNewLinesToArgsList {
			decl += ",\n"
		} else if iArg < len(f.Args)-1 {
			decl += ", "
		}
	}

	decl += ")"

	if len(f.ReturnArgs) == 1 {
			decl += " " + f.ReturnArgs[0].FullType(importAliases)
	} else if len(f.ReturnArgs) > 1 {
		decl += " ("
		for i, ret := range f.ReturnArgs {
			decl += ret.FullType(importAliases)

			if i < len(f.ReturnArgs) - 1 {
				decl += ", "
			}
		}
		decl += ")"
	}

	return decl
}

func funcBaseTemplate(
	decl DeclFunc,
	importAliases map[string]string,
) *template.Template {

	return template.New("").Funcs(map[string]interface{}{
		"FuncReturnDefaults": funcReturnDefaults(decl, importAliases),
	})
}

func funcReturnDefaults(
	decl DeclFunc,
	importAliases map[string]string,
) func() (string, error) {

	return func() (string, error) {

		statement := "return "

		for i, retArg := range decl.ReturnArgs {
			if i > 0 {
				statement += ", "
			}
			argInit, err := retArg.DefaultInit(importAliases)
			if err != nil {
				return "", err
			}

			statement += argInit
		}

		return statement, nil
	}
}
