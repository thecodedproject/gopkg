package gopkg_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/thecodedproject/gopkg"
	"github.com/thecodedproject/gopkg/tmpl"
)

func TestTypeDefaultInit(t *testing.T) {

	testCases := []struct{
		Def gopkg.Type
		ImportAliases map[string]string
		Expected string
		ExpectedErr error
	}{
		{
			Def: gopkg.TypeBool{},
			Expected: "false",
		},
		{
			Def: gopkg.TypeByte{},
			Expected: "0",
		},
		{
			Def: gopkg.TypeError{},
			Expected: "nil",
		},
		{
			Def: gopkg.TypeFloat32{},
			Expected: "0",
		},
		{
			Def: gopkg.TypeFloat64{},
			Expected: "0",
		},
		{
			Def: gopkg.TypeInt{},
			Expected: "0",
		},
		{
			Def: gopkg.TypeInt32{},
			Expected: "0",
		},
		{
			Def: gopkg.TypeInt64{},
			Expected: "0",
		},
		{
			Def: gopkg.TypeString{},
			Expected: "\"\"",
		},
		{
			Def: gopkg.TypeInterface{},
			Expected: "nil",
		},
		{
			Def: gopkg.TypeStruct{},
			Expected: "{}",
		},
		{
			Def: gopkg.TypeArray{
				ValueType: gopkg.TypeUnknownNamed{
					Name: "SomeType",
					Import: "some/import",
				},
			},
			Expected: "nil",
		},
		{
			Def: gopkg.TypePointer{
				ValueType: gopkg.TypeUnknownNamed{
					Name: "float32",
				},
			},
			Expected: "nil",
		},
		{
			Def: gopkg.TypeMap{
				KeyType: gopkg.TypeUnknownNamed{
					Name: "MyType",
					Import: "my/import/path",
				},
				ValueType: gopkg.TypePointer{
					ValueType: gopkg.TypeArray{
						ValueType: gopkg.TypeUnknownNamed{
							Name: "MyOtherType",
							Import: "other/import",
						},
					},
				},
			},
			Expected: "nil",
		},
		{
			Def: gopkg.TypeUnknownNamed{
				Name: "MyType",
				Import: "my/import/path",
			},
			ExpectedErr: errors.New("cannot deduce default init for named type with no value type"),
		},
		{
			Def: gopkg.TypeUnnamedLiteral{},
			ExpectedErr: errors.New("no default init for unnamed literal"),
		},
	}

	for _, test := range testCases {
		t.Run(test.Def.FullType(nil), func(t *testing.T) {
			actual, err := test.Def.DefaultInit(test.ImportAliases)

			if test.ExpectedErr != nil {
				require.Equal(t, test.ExpectedErr, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, test.Expected, actual)
		})

		t.Run(
			"TypeUnknownNamed_with_value_type" + test.Def.FullType(nil),
			func(t *testing.T) {

			namedType := gopkg.TypeUnknownNamed{
				ValueType: test.Def,
			}

			actual, err := namedType.DefaultInit(test.ImportAliases)

			if test.ExpectedErr != nil {
				require.Equal(t, test.ExpectedErr, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, test.Expected, actual)
		})
	}
}

func TestTypeFullType(t *testing.T) {

	testCases := []struct{
		Def gopkg.Type
		ImportAliases map[string]string
		Expected string
	}{
		{
			Def: gopkg.TypeBool{},
			Expected: "bool",
		},
		{
			Def: gopkg.TypeByte{},
			Expected: "byte",
		},
		{
			Def: gopkg.TypeError{},
			Expected: "error",
		},
		{
			Def: gopkg.TypeFloat32{},
			Expected: "float32",
		},
		{
			Def: gopkg.TypeFloat64{},
			Expected: "float64",
		},
		{
			Def: gopkg.TypeInt{},
			Expected: "int",
		},
		{
			Def: gopkg.TypeInt32{},
			Expected: "int32",
		},
		{
			Def: gopkg.TypeInt64{},
			Expected: "int64",
		},
		{
			Def: gopkg.TypeString{},
			Expected: "string",
		},
		{
			Def: gopkg.TypeArray{
				ValueType: gopkg.TypeUnknownNamed{
					Name: "SomeType",
					Import: "some/import",
				},
			},
			ImportAliases: map[string]string{
				"some/import": "some_alias",
			},
			Expected: "[]some_alias.SomeType",
		},
		{
			Def: gopkg.TypePointer{
				ValueType: gopkg.TypeUnknownNamed{
					Name: "float32",
				},
			},
			Expected: "*float32",
		},
		{
			Def: gopkg.TypeMap{
				KeyType: gopkg.TypeUnknownNamed{
					Name: "MyType",
					Import: "my/import/path",
				},
				ValueType: gopkg.TypePointer{
					ValueType: gopkg.TypeArray{
						ValueType: gopkg.TypeUnknownNamed{
							Name: "MyOtherType",
							Import: "other/import",
						},
					},
				},
			},
			ImportAliases: map[string]string{
				"my/import/path": "path_alias",
				"other/import": "other_alias",
			},
			Expected: "map[path_alias.MyType]*[]other_alias.MyOtherType",
		},
		{
			Def: gopkg.TypeUnnamedLiteral{},
			Expected: "",
		},
	}

	for _, test := range testCases {
		t.Run(test.Expected, func(t *testing.T) {
			require.Equal(t, test.Expected, test.Def.FullType(test.ImportAliases))
		})
	}
}

// TestTypeStructFullType tests the TypeStruct.FullType functionality
// It is in it's own test fixutre as it's morecomplicated than the other
// types
func TestTypeStructFullType(t *testing.T) {

	testCases := []struct{
		Name string
		Def gopkg.TypeStruct
		ImportAliases map[string]string
		Expected string
	}{
		{
			Name: "empty struct",
			Def: gopkg.TypeStruct{},
			Expected: "struct {}",
		},
		{
			Name: "struct with fields no import aliases",
			Def: gopkg.TypeStruct{
				Fields: []gopkg.DeclVar{
					{
						Name: "SomeValue",
						Type: gopkg.TypeInt{},
					},
					{
						Name: "SomeOtherValue",
						Type: gopkg.TypeString{},
					},
				},
			},
			Expected:
`struct {
	SomeValue int
	SomeOtherValue string
}`,
		},
		{
			Name: "struct with fields with import aliases",
			Def: gopkg.TypeStruct{
				Fields: []gopkg.DeclVar{
					{
						Name: "MyVal",
						Type: gopkg.TypeUnknownNamed{
							Name: "SomeImportedType",
							Import: "github.com/myrepo",
						},
					},
					{
						Name: "MyOtherVal",
						Type: gopkg.TypePointer{
							ValueType: gopkg.TypeUnknownNamed{
								Name: "SomeOtherImportedType",
								Import: "github.com/myotherrepo",
							},
						},
					},
				},
			},
			ImportAliases: map[string]string{
				"github.com/myrepo": "myrepo",
				"github.com/myotherrepo": "myotherrepo",
			},
			Expected:
`struct {
	MyVal myrepo.SomeImportedType
	MyOtherVal *myotherrepo.SomeOtherImportedType
}`,
		},
	}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {
			require.Equal(t, test.Expected, test.Def.FullType(test.ImportAliases))
		})
	}
}

// TestTypeInterfaceFullType tests the TypeInterface.FullType functionality
// It is in it's own test fixutre as it's morecomplicated than the other
// types
func TestTypeInterfaceFullType(t *testing.T) {

	testCases := []struct{
		Name string
		Def gopkg.TypeInterface
		ImportAliases map[string]string
		Expected string
	}{
		{
			Name: "empty interface",
			Def: gopkg.TypeInterface{},
			Expected: "interface{}",
		},
		{
			Name: "with functions no import aliases",
			Def: gopkg.TypeInterface{
				Funcs: []gopkg.DeclFunc{
					{
						Name: "FirstMethod",
						Args: []gopkg.DeclVar{
							{
								Name: "a",
								Type: gopkg.TypeInt32{},
							},
							{
								Name: "b",
								Type: gopkg.TypeFloat64{},
							},
						},
						ReturnArgs: tmpl.UnnamedReturnArgs(
							gopkg.TypeString{},
							gopkg.TypeError{},
						),
					},
					{
						Name: "SecondMethod",
						ReturnArgs: tmpl.UnnamedReturnArgs(
							gopkg.TypeUnknownNamed{
								Name: "MyType",
							},
						),
					},
				},
			},
			Expected:
`interface {
	FirstMethod(a int32, b float64) (string, error)
	SecondMethod() MyType
}`,
		},
		{
			Name: "with functions and import aliases",
			Def: gopkg.TypeInterface{
				Funcs: []gopkg.DeclFunc{
					{
						Name: "MyMethod",
						Args: []gopkg.DeclVar{
							{
								Name: "val",
								Type: gopkg.TypeUnknownNamed{
									Name: "AStruct",
									Import: "some/path/toa",
								},
							},
						},
						ReturnArgs: tmpl.UnnamedReturnArgs(
							gopkg.TypeUnknownNamed{
								Name: "BStruct",
								Import: "some/other/path/tob",
							},
							gopkg.TypeError{},
						),
					},
					{
						Name: "OtherMethod",
						ReturnArgs: tmpl.UnnamedReturnArgs(
							gopkg.TypeUnknownNamed{
								Name: "CStruct",
								Import: "some/third/path/toc",
							},
						),
					},
				},
			},
			ImportAliases: map[string]string{
				"some/path/toa": "alias_a",
				"some/other/path/tob": "alias_b",
				"some/third/path/toc": "alias_c",
			},
			Expected:
`interface {
	MyMethod(val alias_a.AStruct) (alias_b.BStruct, error)
	OtherMethod() alias_c.CStruct
}`,
		},
	}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {
			require.Equal(t, test.Expected, test.Def.FullType(test.ImportAliases))
		})
	}
}

func TestTypeUnknownNamedDefaultInit(t *testing.T) {

	testCases := []struct{
		// TODO add names to these tests
		Def gopkg.TypeUnknownNamed
		ImportAliases map[string]string
		Expected string
		ExpectedErr error
	}{
		{
			Def: gopkg.TypeUnknownNamed{
				Name: "SomeName",
				ValueType: gopkg.TypeByte{},
			},
			Expected: "0",
		},
		{
			Def: gopkg.TypeUnknownNamed{
				Name: "SomeName",
				ValueType: gopkg.TypeError{},
			},
			Expected: "nil",
		},
		{
			Def: gopkg.TypeUnknownNamed{
				Name: "SomeName",
				ValueType: gopkg.TypeFloat32{},
			},
			Expected: "0",
		},
		{
			Def: gopkg.TypeUnknownNamed{
				Name: "SomeName",
				ValueType: gopkg.TypeFloat64{},
			},
			Expected: "0",
		},
		{
			Def: gopkg.TypeUnknownNamed{
				Name: "SomeName",
				ValueType: gopkg.TypeInt{},
			},
			Expected: "0",
		},
		{
			Def: gopkg.TypeUnknownNamed{
				Name: "SomeName",
				ValueType: gopkg.TypeInt32{},
			},
			Expected: "0",
		},
		{
			Def: gopkg.TypeUnknownNamed{
				Name: "SomeName",
				ValueType: gopkg.TypeInt64{},
			},
			Expected: "0",
		},
		{
			Def: gopkg.TypeUnknownNamed{
				Name: "SomeName",
				ValueType: gopkg.TypeString{},
			},
			Expected: "\"\"",
		},
		{
			Def: gopkg.TypeUnknownNamed{
				Name: "SomeName",
				ValueType: gopkg.TypeInterface{},
			},
			Expected: "nil",
		},
		{
			Def: gopkg.TypeUnknownNamed{
				Name: "SomeNameWithImport",
				Import: "some/other/import",
				ValueType: gopkg.TypeInterface{},
			},
			ImportAliases: map[string]string{
				"some/other/import": "some_other_alias",
			},
			Expected: "nil",
		},
		{
			Def: gopkg.TypeUnknownNamed{
				Name: "SomeName",
				ValueType: gopkg.TypeStruct{},
			},
			Expected: "SomeName{}",
		},
		{
			Def: gopkg.TypeUnknownNamed{
				Name: "SomeOtherName",
				Import: "github.com/somepkg",
				ValueType: gopkg.TypeStruct{},
			},
			ImportAliases: map[string]string{
				"github.com/somepkg": "somepkg",
			},
			Expected: "somepkg.SomeOtherName{}",
		},
		{
			Def: gopkg.TypeUnknownNamed{
				Name: "SomeType",
				ValueType: gopkg.TypeArray{},
			},
			Expected: "nil",
		},
		{
			Def: gopkg.TypeUnknownNamed{
				Name: "SomeType",
				ValueType: gopkg.TypePointer{},
			},
			Expected: "nil",
		},
		{
			Def: gopkg.TypeUnknownNamed{
				Name: "SomeType",
				ValueType: gopkg.TypeMap{},
			},
			Expected: "nil",
		},
		{
			Def: gopkg.TypeUnknownNamed{
				Name: "MyType",
				Import: "my/import/path",
			},
			ExpectedErr: errors.New("cannot deduce default init for named type with no value type"),
		},
	}

	for _, test := range testCases {
		t.Run(test.Def.FullType(nil), func(t *testing.T) {
			actual, err := test.Def.DefaultInit(test.ImportAliases)

			if test.ExpectedErr != nil {
				require.Equal(t, test.ExpectedErr, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, test.Expected, actual)
		})
	}
}

func TestTypeRequiredImports_SimpleTypes(t *testing.T) {

	simpleTypes := []gopkg.Type{
		gopkg.TypeBool{},
		gopkg.TypeByte{},
		gopkg.TypeError{},
		gopkg.TypeFloat32{},
		gopkg.TypeFloat64{},
		gopkg.TypeInt{},
		gopkg.TypeInt32{},
		gopkg.TypeInt64{},
		gopkg.TypeString{},
	}

	for _, simpleType := range simpleTypes {
		t.Run(simpleType.FullType(nil), func(t *testing.T) {
			require.Equal(
				t,
				map[string]bool(nil),
				simpleType.RequiredImports(),
			)
			return
		})
	}
}

func TestTypeRequiredImports_CompositeTypes(t *testing.T) {

	testCases := []struct{
		Name string
		Def gopkg.Type
		Expected map[string]bool
	}{
		{
			Name: "array of simple type",
			Def: gopkg.TypeArray{
				ValueType: gopkg.TypeString{},
			},
			Expected: map[string]bool(nil),
		},
		{
			Name: "array of composite type",
			Def: gopkg.TypeArray{
				ValueType: gopkg.TypeMap{
					KeyType: gopkg.TypeUnknownNamed{
						Import: "array/other/import",
					},
					ValueType: gopkg.TypeUnknownNamed{
						Import: "array/an/import",
					},
				},
			},
			Expected: map[string]bool{
				"array/an/import": true,
				"array/other/import": true,
			},
		},
		{
			Name: "pointer of simple type",
			Def: gopkg.TypePointer{
				ValueType: gopkg.TypeString{},
			},
			Expected: map[string]bool(nil),
		},
		{
			Name: "pointer of composite type",
			Def: gopkg.TypePointer{
				ValueType: gopkg.TypeMap{
					KeyType: gopkg.TypeUnknownNamed{
						Import: "pointer/other/import",
					},
					ValueType: gopkg.TypeUnknownNamed{
						Import: "pointer/an/import",
					},
				},
			},
			Expected: map[string]bool{
				"pointer/an/import": true,
				"pointer/other/import": true,
			},
		},
		{
			Name: "map with simple types",
			Def: gopkg.TypeMap{
				KeyType: gopkg.TypeInt{},
				ValueType: gopkg.TypeString{},
			},
			Expected: map[string]bool(nil),
		},
		{
			Name: "map with named types",
			Def: gopkg.TypeMap{
				KeyType: gopkg.TypeUnknownNamed{
					Import: "other/import",
				},
				ValueType: gopkg.TypeUnknownNamed{
					Import: "an/import",
				},
			},
			Expected: map[string]bool{
				"an/import": true,
				"other/import": true,
			},
		},
		{
			Name: "named without import",
			Def: gopkg.TypeUnknownNamed{},
			Expected: map[string]bool(nil),
		},
		{
			Name: "named with import",
			Def: gopkg.TypeUnknownNamed{
				Import: "some/import",
			},
			Expected: map[string]bool{
				"some/import": true,
			},
		},
		{
			Name: "struct without fields",
			Def: gopkg.TypeStruct{},
			Expected: map[string]bool(nil),
		},
		{
			Name: "struct with named fields",
			Def: gopkg.TypeStruct{
				Fields: []gopkg.DeclVar{
					{
						Type: gopkg.TypeUnknownNamed{
							Import: "import/a",
						},
					},
					{
						Type: gopkg.TypeString{},
					},
					{
						Type: gopkg.TypeUnknownNamed{
							Import: "import/b",
						},
					},
				},
			},
			Expected: map[string]bool{
				"import/a": true,
				"import/b": true,
			},
		},
		{
			Name: "interface without funcs",
			Def: gopkg.TypeInterface{},
			Expected: map[string]bool(nil),
		},
		{
			Name: "interface with named funcs",
			Def: gopkg.TypeInterface{
				Funcs: []gopkg.DeclFunc{
					{
						Args: []gopkg.DeclVar{
							{
								Type: gopkg.TypeUnknownNamed{
									Import: "import/aa",
								},
							},
							{
								Type: gopkg.TypeUnknownNamed{
									Import: "import/bb",
								},
							},
						},
					},
					{
						ReturnArgs: tmpl.UnnamedReturnArgs(
							gopkg.TypeUnknownNamed{
								Import: "import/bb",
							},
							gopkg.TypeUnknownNamed{
								Import: "import/cc",
							},
						),
					},
				},
			},
			Expected: map[string]bool{
				"import/aa": true,
				"import/bb": true,
				"import/cc": true,
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {

			actual := test.Def.RequiredImports()

			if len(test.Expected) == 0 && len(actual) == 0{
				// Allow actual to be either nil or an empty map
				return
			}

			require.Equal(t, test.Expected, test.Def.RequiredImports())
		})
	}

}
