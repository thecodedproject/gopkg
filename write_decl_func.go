package gopkg

import (
	"errors"
	"io"
	"text/template"

	"github.com/iancoleman/strcase"
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

		err = funcTmpl.Execute(w, decl)
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

	argsAndRets, err := funcArgsAndRetArgs(
		f.Args,
		f.ReturnArgs,
		importAliases,
		true,
	)
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
	args []DeclVar,
	returnArgs []DeclVar,
	importAliases map[string]string,
	addNewLinesToArgsList bool,
) (string, error) {

	// Don't add new lines to arg list if there are 0 args or only 1 arg
	if len(args) < 2 {
		addNewLinesToArgsList = false
	}

	argsList, err := funcArgsWithoutParenthesis(
		args,
		importAliases,
		addNewLinesToArgsList,
	)
	if err != nil {
		return "", err
	}

	decl := "(" + argsList + ")"

	if len(returnArgs) == 0 {
		return decl, nil
	}

	retArgsList, err := funcArgsWithoutParenthesis(returnArgs, importAliases, false)
	if err != nil {
		return "", err
	}

	if len(returnArgs) == 1 && returnArgs[0].Name == "" {
		return decl + " " + retArgsList, nil
	}

	return decl + " (" + retArgsList + ")", nil
}

func funcArgsWithoutParenthesis(
	args []DeclVar,
	importAliases map[string]string,
	newlineDelimitted bool,
) (string, error) {

	if len(args) == 0 {
		return "", nil
	}

	areNamedArgs := (args[0].Name != "")

	var argList string
	if newlineDelimitted {
		argList += "\n\t"
	}
	for i, arg := range args {

		if i > 0 {
			if newlineDelimitted {
				argList += ",\n\t"
			} else {
				argList += ", "
			}
		}

		retType, err := arg.FullType(importAliases)
		if err != nil {
			return "", err
		}

		if (areNamedArgs && arg.Name == "") || (!areNamedArgs && arg.Name != "") {
			return "", errors.New("mix of named and unnamed func args")
		}

		if areNamedArgs {
			argList += arg.Name + " "
		}

		argList += retType
	}

	if newlineDelimitted {
		argList += ",\n"
	}

	return argList, nil
}

func funcBaseTemplate(
	decl DeclFunc,
	importAliases map[string]string,
) *template.Template {

	return template.New("").Funcs(map[string]interface{}{
		"FuncReturnDefaults":        funcReturnDefaults(decl, importAliases, false),
		"FuncReturnDefaultsWithErr": funcReturnDefaults(decl, importAliases, true),
		"ToCamel":                   strcase.ToCamel,
		"ToLowerCamel":              strcase.ToLowerCamel,
		"ToSnake":                   strcase.ToSnake,
	})
}

func funcReturnDefaults(
	decl DeclFunc,
	importAliases map[string]string,
	replaceErrorReturnsWithErr bool,
) func() (string, error) {

	return func() (string, error) {

		statement := "return "

		for i, retArg := range decl.ReturnArgs {
			if i > 0 {
				statement += ", "
			}

			if replaceErrorReturnsWithErr {
				if _, isError := retArg.Type.(TypeError); isError {
					statement += "err"
					continue
				}
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
