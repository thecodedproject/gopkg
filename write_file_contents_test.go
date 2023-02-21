package gopkg_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/require"

	"github.com/neurotempest/gopkg"
	"github.com/neurotempest/gopkg/tmpl"
)

func TestWriteFileContents(t *testing.T) {

	testCases := []struct{
		Name string
		C gopkg.FileContents
		ExpectedErr error
	}{
		{
			Name: "empty contents returns error",
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
					{"context", "context"},
					{"some/pkg/path", "pkg_path"},
				},
				Types: []gopkg.DeclType{
					{
						Name: "MyContainer",
						Type: gopkg.TypeStruct{
							Fields: []gopkg.DeclVar{
								{Name: "Num", Type: gopkg.TypeInt64{}},
								{
									Name: "OtherThing",
									Type: gopkg.TypeUnknownNamed{
										Name: "SomeOtherType",
										Import: "some/pkg/path",
									},
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
								Type: gopkg.TypeUnknownNamed{
									Name: "Context",
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
