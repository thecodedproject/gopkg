package gopkg_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/require"

	"github.com/thecodedproject/gopkg"
	"github.com/thecodedproject/gopkg/tmpl"
)

func TestWriteFileContents(t *testing.T) {

	testCases := []struct {
		Name        string
		C           gopkg.FileContents
		ExpectedErr error
	}{
		{
			Name:        "empty contents returns error",
			ExpectedErr: errors.New("package name cannot be empty"),
		},
		{
			Name: "empty file with package name",
			C: gopkg.FileContents{
				PackageName: "mypkg",
			},
		},
		{
			Name: "imports with a struct type decl and a function decl",
			C: gopkg.FileContents{
				PackageName: "struct_and_func",
				Imports: []gopkg.ImportAndAlias{
					{Import: "context", Alias: "context"},
					{Import: "some/pkg/path", Alias: "pkg_path"},
				},
				Types: []gopkg.DeclType{
					{
						Name: "MyContainer",
						Type: gopkg.TypeStruct{
							Fields: []gopkg.DeclVar{
								{
									Name: "Num",
									Type: gopkg.TypeInt64{},
								},
								{
									Name: "OtherThing",
									Type: gopkg.TypeNamed{
										Name:   "SomeOtherType",
										Import: "some/pkg/path",
									},
									StructTag: "key:\"some,tag,values\"",
								},
							},
						},
					},
				},
				Functions: []gopkg.DeclFunc{
					{
						Name: "MyMethod",
						Args: []gopkg.DeclVar{
							{
								Name: "ctx",
								Type: gopkg.TypeNamed{
									Name:   "Context",
									Import: "context",
								},
							},
						},
						ReturnArgs: tmpl.UnnamedReturnArgs(
							gopkg.TypeInt32{},
							gopkg.TypeError{},
						),
						BodyTmpl: `
	{{FuncReturnDefaults}}
`,
					},
				},
			},
		},
		{
			Name: "imports with single variable definaition and a type",
			C: gopkg.FileContents{
				PackageName: "global_vars",
				Imports: tmpl.UnnamedImports(
					"flag",
				),
				Vars: []gopkg.DeclVar{
					{
						Name:         "someArg",
						Type:         gopkg.TypeUnnamedLiteral{},
						LiteralValue: `flag.String("some_arg", "", "set an arg")`,
					},
				},
				Types: []gopkg.DeclType{
					{
						Name: "AStruct",
						Type: gopkg.TypeStruct{},
					},
				},
			},
		},
		{
			Name: "imports with some const and variable definitions with a func",
			C: gopkg.FileContents{
				PackageName: "global_vars",
				Imports: tmpl.UnnamedImports(
					"flag",
				),
				Vars: []gopkg.DeclVar{
					{
						Name:         "someArg",
						Type:         gopkg.TypeUnnamedLiteral{},
						LiteralValue: `flag.String("some_arg", "", "set an arg")`,
					},
					{
						Name:         "someOtherArg",
						Type:         gopkg.TypeUnnamedLiteral{},
						LiteralValue: `flag.String("some_other_arg", "", "set an arg")`,
					},
				},
				Consts: []gopkg.DeclVar{
					{
						Name:         "SomeConst",
						LiteralValue: `123`,
					},
					{
						Name:         "SomeOtherConst",
						Type:         gopkg.TypeString{},
						LiteralValue: `"hello world"`,
					},
				},
				Functions: []gopkg.DeclFunc{
					{
						Name: "AStub",
					},
				},
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {

			buffer := bytes.NewBuffer(nil)

			err := gopkg.WriteFileContents(
				buffer,
				test.C,
			)

			if test.ExpectedErr != nil {
				require.Equal(t, test.ExpectedErr, err)
				return
			}

			require.NoError(t, err)

			g := goldie.New(t)
			g.Assert(t, t.Name(), buffer.Bytes())
		})
	}
}
