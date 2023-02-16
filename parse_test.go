package gopkg_test

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/neurotempest/gopkg"
)

func TestParse(t *testing.T) {

	testCases := []struct{
		Name string
		PkgDir string
		PkgImportPath string
		Expected []gopkg.FileContents
	}{
		{
			Name: "all_built_in_golang_types",
			PkgDir: "test_packages/all_built_in_types",
			PkgImportPath: "some/import/all_built_in_types",
			Expected: []gopkg.FileContents{
				{
					Filepath: "test_packages/all_built_in_types/bools.go",
					PackageName: "all_built_in_types",
					Vars: []gopkg.DeclVar{
						{
							Name: "SomeDefaultVar",
							Import: "some/import/all_built_in_types",
							Type: gopkg.TypeBool{},
							LiteralValue: "true",
						},
					},
					Functions: []gopkg.DeclFunc{
						{
							Name: "MyBoolFunc",
							Import: "some/import/all_built_in_types",
							Args: []gopkg.DeclVar{
								{
									Name: "maybe",
									Type: gopkg.TypeBool{},
								},
							},
							ReturnArgs: []gopkg.Type{
								gopkg.TypeBool{},
							},
						},
					},
					Types: []gopkg.DeclType{
						{
							Name: "MyBoolStruct",
							Import: "some/import/all_built_in_types",
							Type: gopkg.TypeStruct{
								Fields: []gopkg.DeclVar{
									{Name: "Some", Type: gopkg.TypeBool{}},
								},
							},
						},
					},
				},
				{
					Filepath: "test_packages/all_built_in_types/int_float_string_struct.go",
					PackageName: "all_built_in_types",
					Vars: []gopkg.DeclVar{
						{
							Name: "OneInt",
							Import: "some/import/all_built_in_types",
							Type: gopkg.TypeInt{},
							LiteralValue: "1",
						},
						{
							Name: "TwoInt",
							Import: "some/import/all_built_in_types",
							Type: gopkg.TypeInt{},
							LiteralValue: "2",
						},
						{
							Name: "SomeFloat",
							Import: "some/import/all_built_in_types",
							Type: gopkg.TypeFloat32{},
						},
						{
							Name: "SomeUntyped",
							Import: "some/import/all_built_in_types",
							Type: gopkg.TypeUnnamedLiteral{},
							LiteralValue: "\"a string\"",
						},
					},
					Functions: []gopkg.DeclFunc{
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
		},
		{
			Name: "composite_types",
			PkgDir: "test_packages/composite_types",
			PkgImportPath: "some/import/composite_types",
			Expected: []gopkg.FileContents{
				{
					Filepath: "test_packages/composite_types/arrays.go",
					PackageName: "composite_types",
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
					},
					Types: []gopkg.DeclType{
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
					},
				},
				{
					Filepath: "test_packages/composite_types/pointers.go",
					PackageName: "composite_types",
					Functions: []gopkg.DeclFunc{
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
							Name: "MyCustomPointer",
							Import: "some/import/composite_types",
							Type: gopkg.TypePointer{
								ValueType: gopkg.TypePointer{
									ValueType: gopkg.TypeFloat32{},
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
		},
		// TODO implement test for custom types
		// eg structs/interfaces defined in the package being passed and 'unknonw types' from other packages
		// also typedefs and global vars
		/*{
			Name: "custom_types",
			PkgDir: "test_packages/custom_types",
			PkgImportPath: "some/import/custom_types",
			Expected: gopkg.FileContents{
			},
		},*/
		{
			Name: "proto_conversion_package",
			PkgDir: "test_packages/proto_conversion",
			PkgImportPath: "some/import/proto_conversion",
			Expected: []gopkg.FileContents{
				{
					Filepath: "test_packages/proto_conversion/converters.go",
					PackageName: "proto_conversion",
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
					},
				},
				{
					Filepath: "test_packages/proto_conversion/def.pb.go",
					PackageName: "proto_conversion",
					Functions: protoConversionPackageFuncs(),
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
				{
					Filepath: "test_packages/proto_conversion/generate.go",
					PackageName: "proto_conversion",
				},
			},
		},
		{
			Name: "receiver_funcs",
			PkgDir: "test_packages/receiver_funcs",
			PkgImportPath: "some/import/receiver_funcs",
			Expected: []gopkg.FileContents{
				{
					Filepath: "test_packages/receiver_funcs/receiver_funcs.go",
					PackageName: "receiver_funcs",
					Functions: []gopkg.DeclFunc{
						{
							Name: "ValueReceiverFunc",
							Import: "some/import/receiver_funcs",
							Receiver: gopkg.FuncReceiver{
								VarName: "m",
								TypeName: "MyType",
							},
						},
						{
							Name: "PointerRecFunc",
							Import: "some/import/receiver_funcs",
							Receiver: gopkg.FuncReceiver{
								VarName: "m",
								TypeName: "MyType",
								IsPointer: true,
							},
						},
						{
							Name: "OtherPRecFunc",
							Import: "some/import/receiver_funcs",
							Receiver: gopkg.FuncReceiver{
								VarName: "o",
								TypeName: "OtherType",
								IsPointer: true,
							},
						},
						{
							Name: "SomeOtherValRec",
							Import: "some/import/receiver_funcs",
							Receiver: gopkg.FuncReceiver{
								VarName: "o",
								TypeName: "OtherType",
							},
						},
					},
					Types: []gopkg.DeclType{
						{
							Name: "MyType",
							Import: "some/import/receiver_funcs",
							Type: gopkg.TypeStruct{},
						},
						{
							Name: "OtherType",
							Import: "some/import/receiver_funcs",
							Type: gopkg.TypeBool{},
						},
					},
				},
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {

			pc, err := gopkg.Parse(test.PkgDir, test.PkgImportPath)
			require.NoError(t, err)

			//sortDecls(pc[0])

			assert.Equal(t, test.Expected, pc)
		})
	}
}

// TODO Use gopkg full sorting of contents method once implemented
func sortDecls(c gopkg.FileContents) {

	sort.Slice(c.Types, func(i, j int) bool {
		return c.Types[i].Name < c.Types[j].Name
	})

	gopkg.SortFuncs(c.Functions)
}


func protoConversionPackageFuncs() []gopkg.DeclFunc {

	ret := make([]gopkg.DeclFunc, 0)

	ret = append(ret, protoTypeFuncs("IntAsString")...)

	ret = append(ret, protoTypeFuncs("ShopspringDecimal")...)

	ret = append(ret, []gopkg.DeclFunc{
		{
			Name: "init",
			Import: "some/import/proto_conversion",
		},
		{
			Name: "init",
			Import: "some/import/proto_conversion",
		},
	}...)

	return ret
}

func protoTypeFuncs(typeName string) []gopkg.DeclFunc {

	return []gopkg.DeclFunc{
		{
			Name: "Reset",
			Import: "some/import/proto_conversion",
			Receiver: gopkg.FuncReceiver{
				VarName: "m",
				TypeName: typeName,
				IsPointer: true,
			},
		},
		{
			Name: "String",
			Import: "some/import/proto_conversion",
			Receiver: gopkg.FuncReceiver{
				VarName: "m",
				TypeName: typeName,
				IsPointer: true,
			},
			ReturnArgs: []gopkg.Type{
				gopkg.TypeString{},
			},
		},
		{
			Name: "ProtoMessage",
			Import: "some/import/proto_conversion",
			Receiver: gopkg.FuncReceiver{
				TypeName: typeName,
				IsPointer: true,
			},
		},
		{
			Name: "Descriptor",
			Import: "some/import/proto_conversion",
			Receiver: gopkg.FuncReceiver{
				TypeName: typeName,
				IsPointer: true,
			},
			ReturnArgs: []gopkg.Type{
				gopkg.TypeArray{ValueType: gopkg.TypeByte{}},
				gopkg.TypeArray{ValueType: gopkg.TypeInt{}},
			},
		},
		{
			Name: "XXX_Unmarshal",
			Import: "some/import/proto_conversion",
			Receiver: gopkg.FuncReceiver{
				VarName: "m",
				TypeName: typeName,
				IsPointer: true,
			},
			Args: []gopkg.DeclVar{
				{
					Name: "b",
					Type: gopkg.TypeArray{ValueType: gopkg.TypeByte{}},
				},
			},
			ReturnArgs: []gopkg.Type{
				gopkg.TypeError{},
			},
		},
		{
			Name: "XXX_Marshal",
			Import: "some/import/proto_conversion",
			Receiver: gopkg.FuncReceiver{
				VarName: "m",
				TypeName: typeName,
				IsPointer: true,
			},
			Args: []gopkg.DeclVar{
				{
					Name: "b",
					Type: gopkg.TypeArray{ValueType: gopkg.TypeByte{}},
				},
				{
					Name: "deterministic",
					Type: gopkg.TypeBool{},
				},
			},
			ReturnArgs: []gopkg.Type{
				gopkg.TypeArray{ValueType: gopkg.TypeByte{}},
				gopkg.TypeError{},
			},
		},
		{
			Name: "XXX_Merge",
			Import: "some/import/proto_conversion",
			Receiver: gopkg.FuncReceiver{
				VarName: "m",
				TypeName: typeName,
				IsPointer: true,
			},
			Args: []gopkg.DeclVar{
				{
					Name: "src",
					Type: gopkg.TypeUnknownNamed{
						Name: "Message",
						Import: "github.com/golang/protobuf/proto",
					},
				},
			},
		},
		{
			Name: "XXX_Size",
			Import: "some/import/proto_conversion",
			Receiver: gopkg.FuncReceiver{
				VarName: "m",
				TypeName: typeName,
				IsPointer: true,
			},
			ReturnArgs: []gopkg.Type{
				gopkg.TypeInt{},
			},
		},
		{
			Name: "XXX_DiscardUnknown",
			Import: "some/import/proto_conversion",
			Receiver: gopkg.FuncReceiver{
				VarName: "m",
				TypeName: typeName,
				IsPointer: true,
			},
		},
		{
			Name: "GetValue",
			Import: "some/import/proto_conversion",
			Receiver: gopkg.FuncReceiver{
				VarName: "m",
				TypeName: typeName,
				IsPointer: true,
			},
			ReturnArgs: []gopkg.Type{
				gopkg.TypeString{},
			},
		},
	}
}
