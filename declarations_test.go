package gopkg_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/neurotempest/gopkg"
)

func TestDeclFunc_FullDecl(t *testing.T) {

	testCases := []struct{
		Name string
		F gopkg.DeclFunc
		ImportAliases map[string]string
		Expected string
	}{
		{
			F: gopkg.DeclFunc{
				Name: "NoArgsOrReturn",
			},
			Expected: "func NoArgsOrReturn()",
		},
		{
			F: gopkg.DeclFunc{
				Name: "SingleArgPlusReturn",
				Args: []gopkg.DeclVar{
					{
						Name: "myVar",
						Type: gopkg.TypePointer{
							ValueType: gopkg.TypeUnknownNamed{
								Name: "int32",
							},
						},
					},
				},
				ReturnArgs: []gopkg.Type{
					gopkg.TypeUnknownNamed{Name: "error"},
				},
			},
			Expected: "func SingleArgPlusReturn(myVar *int32) error",
		},
		{
			F: gopkg.DeclFunc{
				Name: "MultiArgPlusReturn",
				Args: []gopkg.DeclVar{
					{
						Name: "myVar",
						Type: gopkg.TypeUnknownNamed{
								Name: "int32",
						},
					},
					{
						Name: "myOtherVar",
						Type: gopkg.TypeUnknownNamed{
								Name: "string",
						},
					},
				},
				ReturnArgs: []gopkg.Type{
					gopkg.TypeUnknownNamed{Name: "error"},
				},
			},
			Expected: "func MultiArgPlusReturn(\n\tmyVar int32,\n\tmyOtherVar string,\n) error",
		},
		{
			F: gopkg.DeclFunc{
				Name: "MultiArgPlusMultiReturn",
				Args: []gopkg.DeclVar{
					{
						Name: "a",
						Type: gopkg.TypeUnknownNamed{
								Name: "int",
						},
					},
					{
						Name: "b",
						Type: gopkg.TypeUnknownNamed{
								Name: "int",
						},
					},
				},
				ReturnArgs: []gopkg.Type{
					gopkg.TypeUnknownNamed{Name: "int"},
					gopkg.TypeUnknownNamed{Name: "float32"},
					gopkg.TypeUnknownNamed{Name: "error"},
				},
			},
			Expected: "func MultiArgPlusMultiReturn(\n\ta int,\n\tb int,\n) (int, float32, error)",
		},
		{
			F: gopkg.DeclFunc{
				Name: "SingleArgAndReturnWithImport",
				Args: []gopkg.DeclVar{
					{
						Name: "a",
						Type: gopkg.TypeUnknownNamed{
								Name: "MyType",
								Import: "some/import/path",
						},
					},
				},
				ReturnArgs: []gopkg.Type{
					gopkg.TypeUnknownNamed{
						Name: "OtherType",
						Import: "other/import/path",
					},
				},
			},
			ImportAliases: map[string]string{
				"some/import/path": "some_path",
				"other/import/path": "other_path",
			},
			Expected: "func SingleArgAndReturnWithImport(a some_path.MyType) other_path.OtherType",
		},
	}

	for _, test := range testCases {
		t.Run(test.Expected, func(t *testing.T){
			require.Equal(
				t,
				test.Expected,
				test.F.FullDecl(test.ImportAliases),
			)
		})
	}
}

func TestDeclVar(t *testing.T) {

	_ = gopkg.DeclVar{
		Name: "a",

		Type: gopkg.TypeUnknownNamed{
			Name: "string",
		},
	}
}
