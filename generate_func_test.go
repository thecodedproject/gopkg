package gopkg_test

import (
	"bytes"
	"testing"

	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/require"

	"github.com/neurotempest/gopkg"
)

func TestGenerateFunc(t *testing.T) {

	testCases := []struct{
		Name string
		F gopkg.DeclFunc
		ImportAliases map[string]string
	}{
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
				ReturnArgs: []gopkg.Type{
					gopkg.TypeInt32{},
					gopkg.TypeError{},
				},
			},
		},
		{
			Name: "empty body with imports",
			F: gopkg.DeclFunc{
				Name: "MyFunction",
				Args: []gopkg.DeclVar{
					{
						Name: "myVar",
						Type: gopkg.TypeUnknownNamed{
							Name: "MyType",
							Import: "/some/path/tomypkg",
						},
					},
				},
				ReturnArgs: []gopkg.Type{
					gopkg.TypeUnknownNamed{
						Name: "OtherType",
						Import: "github.com/otherpackage",
					},
				},
			},
			ImportAliases: map[string]string{
				"/some/path/tomypkg": "my_pkg",
				"github.com/otherpackage": "someotherpkg",
			},
		},
		{
			Name: "return default return types",
			F: gopkg.DeclFunc{
				Name: "MyFunction",
				ReturnArgs: []gopkg.Type{
					gopkg.TypeInt32{},
					gopkg.TypeFloat64{},
					gopkg.TypeString{},
					gopkg.TypeUnknownNamed{
						Name: "MyStruct",
						ValueType: gopkg.TypeStruct{},
					},
					gopkg.TypePointer{
						ValueType: gopkg.TypeInt32{},
					},
					gopkg.TypeArray{
						ValueType: gopkg.TypeString{},
					},
				},
				BodyTmpl: `
	{{FuncReturnDefaults}}
`,
			},
		},
		{
			Name: "return default return types with import aliases",
			F: gopkg.DeclFunc{
				Name: "MyFunction",
				ReturnArgs: []gopkg.Type{
					gopkg.TypeUnknownNamed{
						Name: "MyStruct",
						Import: "github.com/some/nice_package",
						ValueType: gopkg.TypeStruct{},
					},
					gopkg.TypeUnknownNamed{
						Name: "OtherStruct",
						Import: "github.com/some/other_package",
						ValueType: gopkg.TypeStruct{},
					},
				},
				BodyTmpl: `
	{{FuncReturnDefaults}}
`,
			},
			ImportAliases: map[string]string{
				"github.com/some/nice_package": "nice_package",
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
						Type: gopkg.TypeUnknownNamed{
							Name: "Context",
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
				ReturnArgs: []gopkg.Type{
					gopkg.TypeInt32{},
					gopkg.TypeUnknownNamed{
						Name: "MyStruct",
						Import: "package/to/wrap",
					},
					gopkg.TypeError{},
				},
				BodyTmpl: `
	return wrap.SomeFuncToWrap(
{{- range .Func.Args}}
		{{.Name}},
{{- end}}
	)
`,
			},
			ImportAliases: map[string]string{
				"context": "context",
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
						Type: gopkg.TypeUnknownNamed{
							Name: "Context",
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
				ReturnArgs: []gopkg.Type{
					gopkg.TypeFloat64{},
					gopkg.TypeError{},
				},
				BodyTmpl: `
	ret, err := SomeFuncToWrap(
{{- range .Func.Args}}
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
	}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {

			buffer := bytes.NewBuffer(nil)

			err := gopkg.GenerateFunc(
				buffer,
				test.F,
				test.ImportAliases,
			)
			require.NoError(t, err)

			g := goldie.New(t)
			g.Assert(t, t.Name(), buffer.Bytes())
		})
	}
}
