package gopkg

import (
)

type DeclFunc struct {
	Name string
	Import string
	Args []DeclVar
	ReturnArgs []Type
	BodyTmpl string
}

type DeclStruct struct {
	Name string
	Import string
	// TODO: Maybe also add struct field descriptors
	Fields []DeclVar//map[string]Type
}

type DeclType struct {
	Name string
	Import string
	Type Type
}

type DeclVar struct {
	Type
	Name string
	Import string
}

func (f DeclFunc) FullDecl(importAliases map[string]string) string {

	return "func " + f.Name + funcArgsAndRetArgs(f, importAliases, true)
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
