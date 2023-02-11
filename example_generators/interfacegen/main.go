package main

import (
	"flag"
	"log"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/neurotempest/gopkg"
)

var (
	interfaceName = flag.String("interface", "", "name of the interface to generate impl for")
	typeName = flag.String("type", "", "name of the type to generate the impl for")
	importPath = flag.String("import", "", "import path of interface pacakage")
)

func main() {

	flag.Parse()

	if *interfaceName == "" {
		log.Fatal("inteferace name must be set with `--interface`")
	}
	if *typeName == "" {
		log.Fatal("type name must be set with `--type`")
	}
	if *importPath == "" {
		log.Fatal("import path must be set with `--import")
	}

	pkgFiles, err := gopkg.GetContents(".", *importPath)
	if err != nil {
		log.Fatal(err.Error())
	}

	pkgName, iDecl := getInterfaceDecl(pkgFiles, *interfaceName)

	err = gopkg.Generate([]gopkg.FileContents{
		makeImplFile(pkgName, *typeName, iDecl),
		makeTestFile(pkgName, *importPath, *typeName, iDecl),
	})
	if err != nil {
		log.Fatal(err.Error())
	}
}

func getInterfaceDecl(pkgFiles []gopkg.FileContents, name string) (string, gopkg.DeclType) {

	for _, file := range pkgFiles {
		for _, typeDecl := range file.Types {
			if typeDecl.Name == name {
				return file.PackageName, typeDecl
			}
		}
	}

	log.Fatal("no such interface", *interfaceName)
	return "", gopkg.DeclType{}
}

func makeImplFile(
	pkgName string,
	typeName string,
	iDecl gopkg.DeclType,
) gopkg.FileContents {

	ret := gopkg.FileContents{
		PackageName: pkgName,
		Filepath: strcase.ToSnake(typeName) + "_impl.go",
	}

	ret.Imports = []gopkg.ImportAndAlias{
		{
			Import: "context",
			Alias: "context",
		},
	}

	ret.Types = []gopkg.DeclType{
		{
			Name: typeName,
			Type: gopkg.TypeStruct{},
		},
	}

	iType, ok := iDecl.Type.(gopkg.TypeInterface)
	if !ok {
		log.Fatal(iDecl.Name, "not an interface declaration")
	}

	for _, iFunc := range iType.Funcs {

		iFunc.Receiver = gopkg.FuncReceiver{
			VarName: strings.ToLower(string(typeName[0])),
			TypeName: typeName,
		}

		iFunc.BodyTmpl = `
	{{FuncReturnDefaults}}
`

		ret.Functions = append(ret.Functions, iFunc)
	}

	return ret
}

func makeTestFile(
	pkgName string,
	pkgImportPath string,
	typeName string,
	iDecl gopkg.DeclType,
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
