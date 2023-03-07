package main

import (
	"errors"
	"flag"
	"log"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/thecodedproject/gopkg"
	"github.com/thecodedproject/gopkg/tmpl"
)

var (
	enumTypeName = flag.String("enum", "", "name of the enum type to stringify")
)

func main() {

	flag.Parse()

	if *enumTypeName == "" {
		log.Fatal("enum type name must be set")
	}

	outputDir := "."

	pkgFiles, err := gopkg.Parse(outputDir)
	if err != nil {
		log.Fatal(err.Error())
	}

	pkgName, enumConsts, err := getPkgNameAndEnumConsts(pkgFiles, *enumTypeName)
	if err != nil {
		log.Fatal(err.Error())
	}

	pkgImportPath, err := gopkg.PackageImportPath(outputDir)

	toGenerate := []gopkg.FileContents{
		makeStringerImplFile(pkgName, enumConsts, *enumTypeName),
		makeTestFile(pkgName, pkgImportPath, enumConsts, *enumTypeName),
	}

	err = gopkg.Lint(toGenerate)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = gopkg.Generate(toGenerate)
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
			if declType, ok := constDecl.Type.(gopkg.TypeNamed); ok {
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
				ReturnArgs: tmpl.UnnamedReturnArgs(
					gopkg.TypeString{},
				),
				BodyData: enumConsts,
				BodyTmpl: `
	switch ` + varName + ` {
{{- range .BodyData}}
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

func makeTestFile(
	pkgName string,
	pkgImportPath string,
	enumConsts []string,
	enumName string,
) gopkg.FileContents {

	testCases := make([]interface{}, 0, len(enumConsts))
	for _, c := range enumConsts {
		testCases = append(testCases, struct{
			Name string
			Enum string
			Expected string
		}{
			Name: c,
			Enum: c,
			Expected: c,
		})
	}

	ret := gopkg.FileContents{
		PackageName: pkgName + "_test",
		Filepath: strcase.ToSnake(enumName) + "_string_test.go",
		Imports: []gopkg.ImportAndAlias{
			{
				Import: "github.com/stretchr/testify/require",
				Alias: "require",
			},
			{
				Import: pkgImportPath,
				Alias: pkgName,
			},
		},
		Functions: []gopkg.DeclFunc{
			{
				Name: "Test" + enumName + "_String",
				Args: []gopkg.DeclVar{
					{
						Name: "t",
						Type: gopkg.TypePointer{
							ValueType: gopkg.TypeNamed{
								Name: "T",
								Import: "testing",
							},
						},
					},
				},
				BodyData: testCases,
				BodyTmpl: `
	testCases := []struct{
		Name string
		Enum ` + pkgName + "." + enumName + `
		Expected string
	}{
{{- range .BodyData}}
		{
			Name: "{{.Name}}",
			Enum: ` + pkgName + `.{{.Enum}},
			Expected: "{{.Expected}}",
		},
{{end -}}
	}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {
			require.Equal(t, test.Expected, test.Enum.String())
		})
	}
`,
			},
		},
	}

	return ret
}
