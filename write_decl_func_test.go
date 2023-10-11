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

func TestWriteDeclFunc(t *testing.T) {

	testCases := []struct {
		Name          string
		F             gopkg.DeclFunc
		ImportAliases map[string]string
		ExpectedErr   error
	}{
		{
			Name: "no args or return or body",
			F: gopkg.DeclFunc{
				Name: "NoArgsOrReturn",
			},
		},
		{
			Name: "value reciever with no args or returns no body",
			F: gopkg.DeclFunc{
				Name: "ANamedValueRecvFunc",
				Receiver: gopkg.FuncReceiver{
					VarName:  "v",
					TypeName: "ValueType",
				},
			},
		},
		{
			Name: "value reciever unnamed with no args or returns no body",
			F: gopkg.DeclFunc{
				Name: "AUnnamedValueRecvFunc",
				Receiver: gopkg.FuncReceiver{
					TypeName: "ValueType",
				},
			},
		},
		{
			Name: "pointer reciever with no args or returns no body",
			F: gopkg.DeclFunc{
				Name: "ANamedPointerRecvFunc",
				Receiver: gopkg.FuncReceiver{
					VarName:   "p",
					TypeName:  "PointerType",
					IsPointer: true,
				},
			},
		},
		{
			Name: "pointer reciever unnamed with no args or returns no body",
			F: gopkg.DeclFunc{
				Name: "AUnnamedPointerRecvFunc",
				Receiver: gopkg.FuncReceiver{
					TypeName:  "ValueType",
					IsPointer: true,
				},
			},
		},
		{
			Name: "single arg plus return no body",
			F: gopkg.DeclFunc{
				Name: "SingleArgPlusReturn",
				Args: []gopkg.DeclVar{
					{
						Name: "myVar",
						Type: gopkg.TypePointer{
							ValueType: gopkg.TypeNamed{
								Name: "int32",
							},
						},
					},
				},
				ReturnArgs: tmpl.UnnamedReturnArgs(
					gopkg.TypeNamed{Name: "error"},
				),
			},
		},
		{
			Name: "multiple args and single return no body",
			F: gopkg.DeclFunc{
				Name: "MultiArgPlusReturn",
				Args: []gopkg.DeclVar{
					{
						Name: "myVar",
						Type: gopkg.TypeNamed{
							Name: "int32",
						},
					},
					{
						Name: "myOtherVar",
						Type: gopkg.TypeNamed{
							Name: "string",
						},
					},
				},
				ReturnArgs: tmpl.UnnamedReturnArgs(
					gopkg.TypeNamed{Name: "error"},
				),
			},
		},
		{
			Name: "multiple args and multiple returns no body",
			F: gopkg.DeclFunc{
				Name: "MultiArgPlusMultiReturn",
				Args: []gopkg.DeclVar{
					{
						Name: "a",
						Type: gopkg.TypeNamed{
							Name: "int",
						},
					},
					{
						Name: "b",
						Type: gopkg.TypeNamed{
							Name: "int",
						},
					},
				},
				ReturnArgs: tmpl.UnnamedReturnArgs(
					gopkg.TypeNamed{Name: "int"},
					gopkg.TypeNamed{Name: "float32"},
					gopkg.TypeNamed{Name: "error"},
				),
			},
		},
		{
			Name: "single arg and return with imports no body",
			F: gopkg.DeclFunc{
				Name: "SingleArgAndReturnWithImport",
				Args: []gopkg.DeclVar{
					{
						Name: "a",
						Type: gopkg.TypeNamed{
							Name:   "MyType",
							Import: "some/import/path",
						},
					},
				},
				ReturnArgs: tmpl.UnnamedReturnArgs(
					gopkg.TypeNamed{
						Name:   "OtherType",
						Import: "other/import/path",
					},
				),
			},
			ImportAliases: map[string]string{
				"some/import/path":  "some_path",
				"other/import/path": "other_path",
			},
		},
		{
			Name: "empty body no imports",
			F: gopkg.DeclFunc{
				Name: "MyFunction",
				Args: []gopkg.DeclVar{
					{
						Name: "myVar",
						Type: gopkg.TypeInt32{},
					},
				},
				ReturnArgs: tmpl.UnnamedReturnArgs(
					gopkg.TypeInt32{},
					gopkg.TypeError{},
				),
			},
		},
		{
			Name: "empty body with imports",
			F: gopkg.DeclFunc{
				Name: "MyFunction",
				Args: []gopkg.DeclVar{
					{
						Name: "myVar",
						Type: gopkg.TypeNamed{
							Name:   "MyType",
							Import: "/some/path/tomypkg",
						},
					},
				},
				ReturnArgs: tmpl.UnnamedReturnArgs(
					gopkg.TypeNamed{
						Name:   "OtherType",
						Import: "github.com/otherpackage",
					},
				),
			},
			ImportAliases: map[string]string{
				"/some/path/tomypkg":      "my_pkg",
				"github.com/otherpackage": "someotherpkg",
			},
		},
		{
			Name: "single named return type",
			F: gopkg.DeclFunc{
				Name: "SingleNamedReturn",
				ReturnArgs: []gopkg.DeclVar{
					{
						Name: "first",
						Type: gopkg.TypeString{},
					},
				},
			},
		},
		{
			Name: "multiple named return types",
			F: gopkg.DeclFunc{
				Name: "NamedReturns",
				ReturnArgs: []gopkg.DeclVar{
					{
						Name: "one",
						Type: gopkg.TypeInt32{},
					},
					{
						Name: "two",
						Type: gopkg.TypeString{},
					},
					{
						Name: "err",
						Type: gopkg.TypeError{},
					},
				},
			},
		},
		{
			Name: "mix of named and unnamed return types returns error",
			F: gopkg.DeclFunc{
				Name: "NamedReturns",
				ReturnArgs: []gopkg.DeclVar{
					{
						Name: "one",
						Type: gopkg.TypeInt32{},
					},
					{
						Type: gopkg.TypeString{},
					},
					{
						Name: "err",
						Type: gopkg.TypeError{},
					},
				},
			},
			ExpectedErr: errors.New("mix of named and unnamed func args"),
		},
		{
			Name: "mix of unnamed and named return types returns error",
			F: gopkg.DeclFunc{
				Name: "NamedReturns",
				ReturnArgs: []gopkg.DeclVar{
					{
						Type: gopkg.TypeInt32{},
					},
					{
						Type: gopkg.TypeString{},
					},
					{
						Name: "err",
						Type: gopkg.TypeError{},
					},
				},
			},
			ExpectedErr: errors.New("mix of named and unnamed func args"),
		},
		{
			Name: "return default return types",
			F: gopkg.DeclFunc{
				Name: "MyFunction",
				ReturnArgs: tmpl.UnnamedReturnArgs(
					gopkg.TypeInt32{},
					gopkg.TypeFloat64{},
					gopkg.TypeString{},
					gopkg.TypeNamed{
						Name:      "MyStruct",
						ValueType: gopkg.TypeStruct{},
					},
					gopkg.TypePointer{
						ValueType: gopkg.TypeInt32{},
					},
					gopkg.TypeArray{
						ValueType: gopkg.TypeString{},
					},
				),
				BodyTmpl: `
	{{FuncReturnDefaults}}
`,
			},
		},
		{
			Name: "return default return types with import aliases",
			F: gopkg.DeclFunc{
				Name: "MyFunction",
				ReturnArgs: tmpl.UnnamedReturnArgs(
					gopkg.TypeNamed{
						Name:      "MyStruct",
						Import:    "github.com/some/nice_package",
						ValueType: gopkg.TypeStruct{},
					},
					gopkg.TypeNamed{
						Name:      "OtherStruct",
						Import:    "github.com/some/other_package",
						ValueType: gopkg.TypeStruct{},
					},
				),
				BodyTmpl: `
	{{FuncReturnDefaults}}
`,
			},
			ImportAliases: map[string]string{
				"github.com/some/nice_package":  "nice_package",
				"github.com/some/other_package": "other_package",
			},
		},
		{
			Name: "return default return types with error",
			F: gopkg.DeclFunc{
				Name: "MyFunction",
				ReturnArgs: tmpl.UnnamedReturnArgs(
					gopkg.TypeNamed{
						Name:      "AStruct",
						Import:    "github.com/some/nice_package",
						ValueType: gopkg.TypeStruct{},
					},
					gopkg.TypeNamed{
						Name:      "OtherStruct",
						Import:    "github.com/some/other_package",
						ValueType: gopkg.TypeStruct{},
					},
					gopkg.TypeError{},
				),
				BodyTmpl: `
	{{FuncReturnDefaultsWithErr}}
`,
			},
			ImportAliases: map[string]string{
				"github.com/some/nice_package":  "nice_package",
				"github.com/some/other_package": "other_package",
			},
		},
		{
			Name: "return default return types with error when there is no error type just returns defaults",
			F: gopkg.DeclFunc{
				Name: "MyFunction",
				ReturnArgs: tmpl.UnnamedReturnArgs(
					gopkg.TypeNamed{
						Name:      "AStruct",
						Import:    "github.com/some/nice_package",
						ValueType: gopkg.TypeStruct{},
					},
					gopkg.TypeNamed{
						Name:      "OtherStruct",
						Import:    "github.com/some/other_package",
						ValueType: gopkg.TypeStruct{},
					},
					gopkg.TypeInt32{},
				),
				BodyTmpl: `
	{{FuncReturnDefaultsWithErr}}
`,
			},
			ImportAliases: map[string]string{
				"github.com/some/nice_package":  "nice_package",
				"github.com/some/other_package": "other_package",
			},
		},
		{
			Name: "simple wrapper function",
			F: gopkg.DeclFunc{
				Name: "SomeWrapper",
				Args: []gopkg.DeclVar{
					{
						Name: "ctx",
						Type: gopkg.TypeNamed{
							Name:   "Context",
							Import: "context",
						},
					},
					{
						Name: "myArg1",
						Type: gopkg.TypeString{},
					},
					{
						Name: "myArg2",
						Type: gopkg.TypeFloat64{},
					},
				},
				ReturnArgs: tmpl.UnnamedReturnArgs(
					gopkg.TypeInt32{},
					gopkg.TypeNamed{
						Name:   "MyStruct",
						Import: "package/to/wrap",
					},
					gopkg.TypeError{},
				),
				BodyTmpl: `
	return wrap.SomeFuncToWrap(
{{- range .Args}}
		{{.Name}},
{{- end}}
	)
`,
			},
			ImportAliases: map[string]string{
				"context":         "context",
				"package/to/wrap": "wrap",
			},
		},
		{
			Name: "wrapper function which returns error",
			F: gopkg.DeclFunc{
				Name: "SomeWrapper",
				Args: []gopkg.DeclVar{
					{
						Name: "ctx",
						Type: gopkg.TypeNamed{
							Name:   "Context",
							Import: "context",
						},
					},
					{
						Name: "a",
						Type: gopkg.TypeFloat64{},
					},
					{
						Name: "b",
						Type: gopkg.TypeFloat64{},
					},
				},
				ReturnArgs: tmpl.UnnamedReturnArgs(
					gopkg.TypeFloat64{},
					gopkg.TypeError{},
				),
				BodyTmpl: `
	ret, err := SomeFuncToWrap(
{{- range .Args}}
		{{.Name}},
{{- end}}
	)
	if err != nil {
		return 0, errors.Wrap(err, "wrapped error")
	}
	return ret, nil
`,
			},
			ImportAliases: map[string]string{
				"context": "context",
			},
		},
		{
			Name: "using strcase mathods",
			F: gopkg.DeclFunc{
				Name: "StringCasing",
				BodyData: struct {
					Snake      string
					LowerCamel string
					UpperCamel string
				}{
					Snake:      "some_snake_case_str",
					LowerCamel: "someLowerCamelStr",
					UpperCamel: "SomeUpperCamelStr",
				},
				BodyTmpl: `
	{{ToLowerCamel .BodyData.Snake}}
	{{ToCamel .BodyData.LowerCamel}}
	{{ToSnake .BodyData.UpperCamel}}
`,
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {

			buffer := bytes.NewBuffer(nil)

			err := gopkg.WriteDeclFunc(
				buffer,
				test.F,
				test.ImportAliases,
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
