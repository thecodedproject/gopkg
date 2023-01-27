package gopkg

import (
	"io"
	"text/template"
)

func GenerateFunc(
	w io.Writer,
	decl DeclFunc,
	importAliases map[string]string,
) error {

	w.Write([]byte(decl.FullDecl(importAliases)))
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
