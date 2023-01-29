package gopkg_test

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/neurotempest/gopkg"
)

func TestGetPackageContents(t *testing.T) {

	testCases := []struct{
		Name string
		PkgDir string
		PkgImportPath string
		Expected gopkg.Contents
	}{
		{
			Name: "all_built_in_golang_types",
			PkgDir: "test_packages/all_built_in_types",
			PkgImportPath: "some/import/all_built_in_types",
			Expected: gopkg.Contents{
				Functions: []gopkg.DeclFunc{
					{
						Name: "SomeFloats",
						Import: "some/import/all_built_in_types",
						Args: []gopkg.DeclVar{
							{
								Name: "a",
								Type: gopkg.TypeFloat32{},
							},
							{
								Name: "b",
								Type: gopkg.TypeFloat64{},
							},
						},
						ReturnArgs: []gopkg.Type{
							gopkg.TypeFloat32{},
							gopkg.TypeFloat64{},
						},
					},
					{
						Name: "SomeInts",
						Import: "some/import/all_built_in_types",
						Args: []gopkg.DeclVar{
							{
								Name: "a",
								Type: gopkg.TypeInt{},
							},
							{
								Name: "b",
								Type: gopkg.TypeInt64{},
							},
							{
								Name: "c",
								Type: gopkg.TypeInt32{},
							},
						},
						ReturnArgs: []gopkg.Type{
							gopkg.TypeInt{},
							gopkg.TypeInt64{},
							gopkg.TypeInt32{},
						},
					},
					{
						Name: "SomeStrings",
						Import: "some/import/all_built_in_types",
						Args: []gopkg.DeclVar{
							{
								Name: "a",
								Type: gopkg.TypeString{},
							},
						},
						ReturnArgs: []gopkg.Type{
							gopkg.TypeString{},
						},
					},
				},
				Types: []gopkg.DeclType{
					{
						Name: "SomeStruct",
						Import: "some/import/all_built_in_types",
						Type: gopkg.TypeStruct{
							Fields: []gopkg.DeclVar{
								{Name: "IA", Type: gopkg.TypeInt{}},
								{Name: "IB", Type: gopkg.TypeInt32{}},
								{Name: "IC", Type: gopkg.TypeInt64{}},
								{Name: "FA", Type: gopkg.TypeFloat32{}},
								{Name: "FB", Type: gopkg.TypeFloat64{}},
								{Name: "S", Type: gopkg.TypeString{}},
							},
						},
					},
				},
			},
		},
		{
			Name: "composite_types",
			PkgDir: "test_packages/composite_types",
			PkgImportPath: "some/import/composite_types",
			Expected: gopkg.Contents{
				Functions: []gopkg.DeclFunc{
					{
						Name: "SomeArrayFunc",
						Import: "some/import/composite_types",
						Args: []gopkg.DeclVar{
							{
								Name: "a",
								Type: gopkg.TypeArray{
									ValueType: gopkg.TypeUnknownNamed{
										Name: "Decimal",
										Import: "github.com/shopspring/decimal",
									},
								},
							},
							{
								Name: "b",
								Type: gopkg.TypeArray{
									ValueType: gopkg.TypeFloat32{},
								},
							},
						},
						ReturnArgs: []gopkg.Type{
							gopkg.TypeArray{
								ValueType: gopkg.TypeUnknownNamed{
									Name: "SomeArrayStruct",
									Import: "some/import/composite_types",
								},
							},
						},
					},
					{
						Name: "SomePointerFunc",
						Import: "some/import/composite_types",
						Args: []gopkg.DeclVar{
							{
								Name: "a",
								Type: gopkg.TypePointer{
									ValueType: gopkg.TypeFloat32{},
								},
							},
							{
								Name: "b",
								Type: gopkg.TypePointer{
									ValueType: gopkg.TypeUnknownNamed{
										Name: "SomePointerStruct",
										Import: "some/import/composite_types",
									},
								},
							},
						},
						ReturnArgs: []gopkg.Type{
							gopkg.TypePointer{
								ValueType: gopkg.TypeString{},
							},
						},
					},
				},
				Types: []gopkg.DeclType{
					{
						Name: "MyCustomArrayType",
						Import: "some/import/composite_types",
						Type: gopkg.TypeArray{
							ValueType: gopkg.TypeArray{
								ValueType: gopkg.TypeArray{
									ValueType: gopkg.TypeFloat64{},
								},
							},
						},
					},
					{
						Name: "MyCustomPointer",
						Import: "some/import/composite_types",
						Type: gopkg.TypePointer{
							ValueType: gopkg.TypePointer{
								ValueType: gopkg.TypeFloat32{},
							},
						},
					},
					{
						Name: "SomeArrayInterface",
						Import: "some/import/composite_types",
						Type: gopkg.TypeInterface{
							Funcs: []gopkg.DeclFunc{
								{
									Name: "ArrayMaker",
									Args: []gopkg.DeclVar{
										{Name: "n", Type: gopkg.TypeInt64{}},
										{Name: "vals", Type: gopkg.TypeString{}},
									},
									ReturnArgs: []gopkg.Type{
										gopkg.TypeArray{
											ValueType: gopkg.TypeString{},
										},
									},
								},
							},
						},
					},
					{
						Name: "SomeArrayStruct",
						Import: "some/import/composite_types",
						Type: gopkg.TypeStruct{
							Fields: []gopkg.DeclVar{
								{
									Name: "AOfInts",
									Type: gopkg.TypeArray{
										ValueType: gopkg.TypeInt64{},
									},
								},
								{
									Name: "AOfPToStrings",
									Type: gopkg.TypeArray{
										ValueType: gopkg.TypePointer{
											ValueType: gopkg.TypeString{},
										},
									},
								},
							},
						},
					},
					{
						Name: "SomePointerInterface",
						Import: "some/import/composite_types",
						Type: gopkg.TypeInterface{
							Funcs: []gopkg.DeclFunc{
								{
									Name: "Something",
									ReturnArgs: []gopkg.Type{
										gopkg.TypePointer{
											ValueType: gopkg.TypeInt64{},
										},
									},
								},
								{
									Name: "PointerMaker",
									Args: []gopkg.DeclVar{
										{
											Name: "val",
											Type: gopkg.TypePointer{
												ValueType: gopkg.TypeString{},
											},
										},
									},
									ReturnArgs: []gopkg.Type{
										gopkg.TypePointer{
											ValueType: gopkg.TypeFloat64{},
										},
									},
								},
							},
						},
					},
					{
						Name: "SomePointerStruct",
						Import: "some/import/composite_types",
						Type: gopkg.TypeStruct{
							Fields: []gopkg.DeclVar{
								{
									Name: "PToInt",
									Type: gopkg.TypePointer{
										ValueType: gopkg.TypeInt32{},
									},
								},
							},
						},
					},
				},
			},
		},
		// TODO implement test for custom types
		// eg structs/interfaces defined in the package being passed and 'unknonw types' from other packages
		// also typedefs and global vars
		/*{
			Name: "custom_types",
			PkgDir: "test_packages/custom_types",
			PkgImportPath: "some/import/custom_types",
			Expected: gopkg.Contents{
			},
		},*/
		{
			Name: "proto_conversion_package",
			PkgDir: "test_packages/proto_conversion",
			PkgImportPath: "some/import/proto_conversion",
			Expected: gopkg.Contents{
				Functions: []gopkg.DeclFunc{
					{
						Name: "IntAsStringFromProto",
						Import: "some/import/proto_conversion",
						Args: []gopkg.DeclVar{
							{
								Name: "v",
								Type: gopkg.TypePointer{
									ValueType: gopkg.TypeUnknownNamed{
										Name: "IntAsString",
										Import: "some/import/proto_conversion",
									},
								},
							},
						},
						ReturnArgs: []gopkg.Type{
							gopkg.TypeInt{},
							gopkg.TypeError{},
						},
					},
					{
						Name: "IntAsStringToProto",
						Import: "some/import/proto_conversion",
						Args: []gopkg.DeclVar{
							{
								Name: "i",
								Type: gopkg.TypeInt{},
							},
						},
						ReturnArgs: []gopkg.Type{
							gopkg.TypePointer{
								ValueType: gopkg.TypeUnknownNamed{
									Name: "IntAsString",
									Import: "some/import/proto_conversion",
								},
							},
							gopkg.TypeError{},
						},
					},
					{
						Name: "ShopspringDecimalFromProto",
						Import: "some/import/proto_conversion",
						Args: []gopkg.DeclVar{
							{
								Name: "v",
								Type: gopkg.TypePointer{
									ValueType: gopkg.TypeUnknownNamed{
										Name: "ShopspringDecimal",
										Import: "some/import/proto_conversion",
									},
								},
							},
						},
						ReturnArgs: []gopkg.Type{
							gopkg.TypeUnknownNamed{
								Name: "Decimal",
								Import: "github.com/shopspring/decimal",
							},
							gopkg.TypeError{},
						},
					},
					{
						Name: "ShopspringDecimalToProto",
						Import: "some/import/proto_conversion",
						Args: []gopkg.DeclVar{
							{
								Name: "v",
								Type: gopkg.TypeUnknownNamed{
									Name: "Decimal",
									Import: "github.com/shopspring/decimal",
								},
							},
						},
						ReturnArgs: []gopkg.Type{
							gopkg.TypePointer{
								ValueType: gopkg.TypeUnknownNamed{
									Name: "ShopspringDecimal",
									Import: "some/import/proto_conversion",
								},
							},
							gopkg.TypeError{},
						},
					},
					{
						Name: "init",
						Import: "some/import/proto_conversion",
					},
					{
						Name: "init",
						Import: "some/import/proto_conversion",
					},
				},
				Types: []gopkg.DeclType{
					{
						Name: "IntAsString",
						Import: "some/import/proto_conversion",
						Type: gopkg.TypeStruct{
							Fields: []gopkg.DeclVar{
								{Name: "Value", Type: gopkg.TypeString{}},
								{Name: "XXX_NoUnkeyedLiteral", Type: gopkg.TypeStruct{}},
								{
									Name: "XXX_unrecognized",
									Type: gopkg.TypeArray{
										ValueType: gopkg.TypeByte{},
									},
								},
								{Name: "XXX_sizecache", Type: gopkg.TypeInt32{}},
							},
						},
					},
					{
						Name: "ShopspringDecimal",
						Import: "some/import/proto_conversion",
						Type: gopkg.TypeStruct{
							Fields: []gopkg.DeclVar{
								{Name: "Value", Type: gopkg.TypeString{}},
								{
									Name: "XXX_NoUnkeyedLiteral",
									Type: gopkg.TypeStruct{},
								},
								{
									Name: "XXX_unrecognized",
									Type: gopkg.TypeArray{
										ValueType: gopkg.TypeByte{},
									},
								},
								{Name: "XXX_sizecache", Type: gopkg.TypeInt32{}},
							},
						},
					},
				},
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {

			pc, err := gopkg.GetContents(test.PkgDir, test.PkgImportPath)
			require.NoError(t, err)

			// TODO consider sorting as part of parsing (as AST parsing is non-deterministic in order)
			sortDecls(pc)

			assert.Equal(t, test.Expected, pc)
		})
	}
}

func sortDecls(c gopkg.Contents) {

	sort.Slice(c.Types, func(i, j int) bool {
		return c.Types[i].Name < c.Types[j].Name
	})

	sort.Slice(c.Functions, func(i, j int) bool {
		return c.Functions[i].Name < c.Functions[j].Name
	})
}
