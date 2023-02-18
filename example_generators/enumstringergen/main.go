package main

import (
	"errors"
	"flag"
	"log"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/neurotempest/gopkg"

	"fmt"
)

var (
	enumTypeName = flag.String("enum", "", "name of the enum type to stringify")
)

func main() {

	flag.Parse()

	if *enumTypeName == "" {
		log.Fatal("enum type name must be set")
	}

	pkgFiles, err := gopkg.Parse(".", "")
	if err != nil {
		log.Fatal(err.Error())
	}

	pkgName, enumConsts, err := getPkgNameAndEnumConsts(pkgFiles, *enumTypeName)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = gopkg.Generate([]gopkg.FileContents{
		makeStringerImplFile(pkgName, enumConsts, *enumTypeName),
		// TODO implement test file
		//makeTestFile(pkgName, enumConsts, *enumTypeName),
	})
	if err != nil {
		log.Fatal(err.Error())
	}
}

func getPkgNameAndEnumConsts(pkgFiles []gopkg.FileContents, enumName string) (string, []string, error) {

	var pkgName string
	var enumConsts []string
	for _, file := range pkgFiles {
		pkgName = file.PackageName
		for _, constDecl := range file.Consts {
			if declType, ok := constDecl.Type.(gopkg.TypeUnknownNamed); ok {
				if declType.Name == enumName {
					enumConsts = append(enumConsts, constDecl.Name)
				}
			}
		}
	}

	if len(enumConsts) == 0 {
		return "", nil, errors.New("no consts found for enum type name" + enumName)
	}

	return pkgName, enumConsts, nil
}

func makeStringerImplFile(
	pkgName string,
	enumConsts []string,
	enumName string,
) gopkg.FileContents {

	for _, t := range enumConsts {
		fmt.Println("const:", t)
	}

	varName := strings.ToLower(string(enumName[0]))

	ret := gopkg.FileContents{
		PackageName: pkgName,
		Filepath: strcase.ToSnake(enumName) + "_string.go",
		Functions: []gopkg.DeclFunc{
			{
				Name: "String",
				Receiver: gopkg.FuncReceiver{
					VarName: varName,
					TypeName: enumName,
				},
				ReturnArgs: []gopkg.Type{
					gopkg.TypeString{},
				},
				BodyData: enumConsts,
				BodyTmpl: `
	switch ` + varName + ` {
{{- range .Func.BodyData}}
	case {{.}}:
		return "{{.}}"
{{- end}}
	default:
		return "` + enumName + `: Unknown value"
	}
`,
			},
		},
	}

	return ret
}

/*
func makeTestFile(
	pkgName string,
	pkgImportPath string
	enumConsts []string,
	enumName string,
) gopkg.FileContents {


	ret := gopkg.FileContents{
		PackageName: pkgName + "_test",
		Filepath: strcase.ToSnake(typeName) + "_impl_test.go",
	}

	ret.Imports = []gopkg.ImportAndAlias{
		{
			Import: "testing",
			Alias: "testing",
		},
		{
			Import: "github.com/stretchr/testify/require",
			Alias: "require",
		},
		{
			Import: pkgImportPath,
			Alias: pkgName,
		},
	}

	iType, ok := iDecl.Type.(gopkg.TypeInterface)
	if !ok {
		log.Fatal(iDecl.Name, "not an interface declaration")
	}

	for _, iFunc := range iType.Funcs {

		ret.Functions = append(ret.Functions, gopkg.DeclFunc{
			Name: "Test" + typeName + "_" + iFunc.Name,
			Args: []gopkg.DeclVar{
				{
					Name: "t",
					Type: gopkg.TypePointer{
						ValueType: gopkg.TypeUnknownNamed{
							Name: "T",
							Import: "testing",
						},
					},
				},
			},
			BodyTmpl: `
	testCases := []struct{
		Name string
		V ` + pkgName + "." + typeName + `
	}{
		{
			Name: "empty type with empty inputs",
		},
	}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {
			require.True(t, false, "TODO: implement test")
		})
	}
`,
		})
	}

	return ret
}
*/
