package gopkg

import (
	"errors"
	"io"
	"text/template"
)

func WriteDeclFunc(
	w io.Writer,
	decl DeclFunc,
	importAliases map[string]string,
) error {

	funcDecl, err := fullFuncDecl(decl, importAliases)
	if err != nil {
		return err
	}

	w.Write([]byte(funcDecl))
	w.Write([]byte(" {\n"))

	if decl.BodyTmpl != "" {

		tmpl := funcBaseTemplate(decl, importAliases)

		funcTmpl, err := tmpl.Parse(decl.BodyTmpl)
		if err != nil {
			return err
		}

		// TODO: Pass in the DeclFunc directly as the template.Execute data
		// There is no need for the extra redirection of `.Func`
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

func fullFuncDecl(
	f DeclFunc,
	importAliases map[string]string,
) (string, error) {

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

	argsAndRets, err := funcArgsAndRetArgs(f, importAliases, true)
	if err != nil {
		return "", err
	}

	decl += f.Name + argsAndRets

	return decl, nil
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
) (string, error) {

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

		argType, err := arg.FullType(importAliases)
		if err != nil {
			return "", err
		}

		decl += arg.Name + " " + argType

		if addNewLinesToArgsList {
			decl += ",\n"
		} else if iArg < len(f.Args)-1 {
			decl += ", "
		}
	}

	decl += ")"

	retArgs, err := funcRetArgs(f, importAliases)
	if err != nil {
		return "", err
	}

	return decl + retArgs, nil
}

func funcRetArgs(
	f DeclFunc,
	importAliases map[string]string,
) (string, error) {

	var retArgs string

	if len(f.ReturnArgs) == 1 {
			retType, err := f.ReturnArgs[0].FullType(importAliases)
			if err != nil {
				return "", err
			}

			if f.ReturnArgs[0].Name == "" {
				retArgs += " " + retType
			} else {
				retArgs += " (" + f.ReturnArgs[0].Name + " " + retType + ")"
			}
	} else if len(f.ReturnArgs) > 1 {

		namedRetArgs := (f.ReturnArgs[0].Name != "")

		retArgs += " ("
		for i, ret := range f.ReturnArgs {

			retType, err := ret.FullType(importAliases)
			if err != nil {
				return "", err
			}

			if (namedRetArgs && ret.Name == "") || (!namedRetArgs && ret.Name != "") {
				return "", errors.New("mix of named and unnamed func return args")
			}

			if namedRetArgs {
				retArgs += ret.Name + " "
			}

			retArgs += retType

			if i < len(f.ReturnArgs) - 1 {
				retArgs += ", "
			}
		}
		retArgs += ")"
	}

	return retArgs, nil
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
