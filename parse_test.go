package gopkg_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thecodedproject/gopkg"
	"github.com/thecodedproject/gopkg/tmpl"
)

func TestParse(t *testing.T) {

	testCases := []struct{
		Name string
		PkgDir string
		ParseOptions []gopkg.ParseOption
		Expected []gopkg.FileContents
	}{
		{
			Name: "no options will auto-detect pkg import dir",
			PkgDir: "test_packages/very_simple",
			Expected: []gopkg.FileContents{
				{
					Filepath: "test_packages/very_simple/very_simple.go",
					PackageName: "very_simple",
					PackageImportPath: "github.com/thecodedproject/gopkg/test_packages/very_simple",
					Vars: []gopkg.DeclVar{
						{
							Name: "MyVar",
							Import: "github.com/thecodedproject/gopkg/test_packages/very_simple",
							Type: gopkg.TypeInt{},
						},
					},
				},
			},
		},
		{
			Name: "all_built_in_golang_types",
			PkgDir: "test_packages/all_built_in_types",
			ParseOptions: []gopkg.ParseOption{
				gopkg.ParseWithPkgImportPath("some/import/all_built_in_types"),
			},
			Expected: []gopkg.FileContents{
				{
					Filepath: "test_packages/all_built_in_types/bools.go",
					PackageName: "all_built_in_types",
					PackageImportPath: "some/import/all_built_in_types",
					Consts: []gopkg.DeclVar{
						{
							Name: "AConstant",
							Import: "some/import/all_built_in_types",
							Type: gopkg.TypeBool{},
							LiteralValue: "false",
						},
					},
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
							ReturnArgs: tmpl.UnnamedReturnArgs(
								gopkg.TypeBool{},
							),
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
					PackageImportPath: "some/import/all_built_in_types",
					Consts: []gopkg.DeclVar{
						{
							Name: "MyConst",
							Import: "some/import/all_built_in_types",
							Type: gopkg.TypeUnnamedLiteral{},
							LiteralValue: "\"some val\"",
						},
						{
							Name: "AnotherConst",
							Import: "some/import/all_built_in_types",
							Type: gopkg.TypeUnnamedLiteral{},
							LiteralValue: "\"other val\"",
						},
						{
							Name: "RealNumberConst",
							Import: "some/import/all_built_in_types",
							Type: gopkg.TypeFloat64{},
							LiteralValue: "1.234",
						},
					},
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
							ReturnArgs: tmpl.UnnamedReturnArgs(
								gopkg.TypeInt{},
								gopkg.TypeInt64{},
								gopkg.TypeInt32{},
							),
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
							ReturnArgs: tmpl.UnnamedReturnArgs(
								gopkg.TypeFloat32{},
								gopkg.TypeFloat64{},
							),
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
							ReturnArgs: tmpl.UnnamedReturnArgs(
								gopkg.TypeString{},
							),
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
			ParseOptions: []gopkg.ParseOption{
				gopkg.ParseWithPkgImportPath("some/import/composite_types"),
			},
			Expected: []gopkg.FileContents{
				{
					Filepath: "test_packages/composite_types/arrays.go",
					PackageName: "composite_types",
					PackageImportPath: "some/import/composite_types",
					Functions: []gopkg.DeclFunc{
						{
							Name: "SomeArrayFunc",
							Import: "some/import/composite_types",
							Args: []gopkg.DeclVar{
								{
									Name: "a",
									Type: gopkg.TypeArray{
										ValueType: gopkg.TypeNamed{
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
							ReturnArgs: tmpl.UnnamedReturnArgs(
								gopkg.TypeArray{
									ValueType: gopkg.TypeNamed{
										Name: "SomeArrayStruct",
										Import: "some/import/composite_types",
									},
								},
							),
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
										ReturnArgs: tmpl.UnnamedReturnArgs(
											gopkg.TypeArray{
												ValueType: gopkg.TypeString{},
											},
										),
									},
								},
							},
						},
					},
				},
				{
					Filepath: "test_packages/composite_types/function_types.go",
					PackageName: "composite_types",
					PackageImportPath: "some/import/composite_types",
					Vars: []gopkg.DeclVar{
						{
							Name: "SomeVar",
							Import: "some/import/composite_types",
							Type: gopkg.TypeFunc{
								Args: tmpl.UnnamedReturnArgs(
									gopkg.TypeInt{},
								),
								ReturnArgs: tmpl.UnnamedReturnArgs(
									gopkg.TypeString{},
								),
							},
						},
					},
					Functions: []gopkg.DeclFunc{
						{
							Name: "SomeFunc",
							Import: "some/import/composite_types",
							Args: []gopkg.DeclVar{
								{
									Name: "f",
									Type: gopkg.TypeFunc{},
								},
							},
							ReturnArgs: tmpl.UnnamedReturnArgs(
								gopkg.TypeFunc{
									ReturnArgs: tmpl.UnnamedReturnArgs(
										gopkg.TypeInt{},
										gopkg.TypeInt{},
									),
								},
							),
						},
					},
					Types: []gopkg.DeclType{
						{
							Name: "SomeType",
							Import: "some/import/composite_types",
							Type: gopkg.TypeFunc{
								ReturnArgs: tmpl.UnnamedReturnArgs(
									gopkg.TypeError{},
								),
							},
						},
						{
							Name: "SomeStruct",
							Import: "some/import/composite_types",
							Type: gopkg.TypeStruct{
								Fields: []gopkg.DeclVar{
									{
										Name: "UnnamedFunc",
										Type: gopkg.TypeFunc{
											Args: tmpl.UnnamedReturnArgs(
												gopkg.TypeInt{},
												gopkg.TypeFloat32{},
											),
											ReturnArgs: tmpl.UnnamedReturnArgs(
												gopkg.TypeString{},
												gopkg.TypeError{},
											),
										},
									},
									{
										Name: "NamedFunc",
										Type: gopkg.TypeFunc{
											Args: []gopkg.DeclVar{
												{
													Name: "a",
													Type: gopkg.TypeInt64{},
												},
												{
													Name: "b",
													Type: gopkg.TypeBool{},
												},
											},
											ReturnArgs: []gopkg.DeclVar{
												{
													Name: "c",
													Type: gopkg.TypeFloat64{},
												},
												{
													Name: "d",
													Type: gopkg.TypeString{},
												},
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
					PackageImportPath: "some/import/composite_types",
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
										ValueType: gopkg.TypeNamed{
											Name: "SomePointerStruct",
											Import: "some/import/composite_types",
										},
									},
								},
							},
							ReturnArgs: tmpl.UnnamedReturnArgs(
								gopkg.TypePointer{
									ValueType: gopkg.TypeString{},
								},
							),
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
										ReturnArgs: tmpl.UnnamedReturnArgs(
											gopkg.TypePointer{
												ValueType: gopkg.TypeInt64{},
											},
										),
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
										ReturnArgs: tmpl.UnnamedReturnArgs(
											gopkg.TypePointer{
												ValueType: gopkg.TypeFloat64{},
											},
										),
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
		{
			Name: "custom_types",
			PkgDir: "test_packages/custom_types",
			ParseOptions: []gopkg.ParseOption{
				gopkg.ParseWithPkgImportPath("some/import/custom_types"),
			},
			Expected: []gopkg.FileContents{
				{
					Filepath: "test_packages/custom_types/embedded_types.go",
					PackageName: "custom_types",
					PackageImportPath: "some/import/custom_types",
					Types: []gopkg.DeclType{
						{
							Name: "SingleEmbed",
							Import: "some/import/custom_types",
							Type: gopkg.TypeStruct{
								Embeds: []gopkg.Type{
									gopkg.TypeNamed{
										Name: "Context",
										Import: "context",
									},
								},
							},
						},
						{
							Name: "InterfaceEmbed",
							Import: "some/import/custom_types",
							Type: gopkg.TypeInterface{
								Embeds: []gopkg.Type{
									gopkg.TypeFloat64{},
								},
								Funcs: []gopkg.DeclFunc{
									{Name: "MyFunc"},
								},
							},
						},
						{
							Name: "ManyEmbeds",
							Import: "some/import/custom_types",
							Type: gopkg.TypeStruct{
								Embeds: []gopkg.Type{
									gopkg.TypeError{},
									gopkg.TypeNamed{
										Name: "Context",
										Import: "context",
									},
									gopkg.TypeInt32{},
								},
								Fields: []gopkg.DeclVar{
									{
										Name: "myVar",
										Type: gopkg.TypeString{},
									},
								},
							},
						},
						{
							Name: "InterfaceManyEmbeds",
							Import: "some/import/custom_types",
							Type: gopkg.TypeInterface{
								Embeds: []gopkg.Type{
									gopkg.TypeNamed{
										Name: "SingleEmbed",
										Import: "some/import/custom_types",
									},
									gopkg.TypeNamed{
										Name: "InterfaceEmbed",
										Import: "some/import/custom_types",
									},
									gopkg.TypeNamed{
										Name: "Context",
										Import: "context",
									},
									gopkg.TypeError{},
								},
								Funcs: []gopkg.DeclFunc{},
							},
						},
					},
				},
			},
		},
		{
			Name: "proto_conversion_package",
			PkgDir: "test_packages/proto_conversion",
			ParseOptions: []gopkg.ParseOption{
				gopkg.ParseWithPkgImportPath("some/import/proto_conversion"),
			},
			Expected: []gopkg.FileContents{
				{
					Filepath: "test_packages/proto_conversion/converters.go",
					PackageName: "proto_conversion",
					PackageImportPath: "some/import/proto_conversion",
					Functions: []gopkg.DeclFunc{
						{
							Name: "IntAsStringFromProto",
							Import: "some/import/proto_conversion",
							Args: []gopkg.DeclVar{
								{
									Name: "v",
									Type: gopkg.TypePointer{
										ValueType: gopkg.TypeNamed{
											Name: "IntAsString",
											Import: "some/import/proto_conversion",
										},
									},
								},
							},
							ReturnArgs: tmpl.UnnamedReturnArgs(
								gopkg.TypeInt{},
								gopkg.TypeError{},
							),
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
							ReturnArgs: tmpl.UnnamedReturnArgs(
								gopkg.TypePointer{
									ValueType: gopkg.TypeNamed{
										Name: "IntAsString",
										Import: "some/import/proto_conversion",
									},
								},
								gopkg.TypeError{},
							),
						},
						{
							Name: "ShopspringDecimalFromProto",
							Import: "some/import/proto_conversion",
							Args: []gopkg.DeclVar{
								{
									Name: "v",
									Type: gopkg.TypePointer{
										ValueType: gopkg.TypeNamed{
											Name: "ShopspringDecimal",
											Import: "some/import/proto_conversion",
										},
									},
								},
							},
							ReturnArgs: tmpl.UnnamedReturnArgs(
								gopkg.TypeNamed{
									Name: "Decimal",
									Import: "github.com/shopspring/decimal",
								},
								gopkg.TypeError{},
							),
						},
						{
							Name: "ShopspringDecimalToProto",
							Import: "some/import/proto_conversion",
							Args: []gopkg.DeclVar{
								{
									Name: "v",
									Type: gopkg.TypeNamed{
										Name: "Decimal",
										Import: "github.com/shopspring/decimal",
									},
								},
							},
							ReturnArgs: tmpl.UnnamedReturnArgs(
								gopkg.TypePointer{
									ValueType: gopkg.TypeNamed{
										Name: "ShopspringDecimal",
										Import: "some/import/proto_conversion",
									},
								},
								gopkg.TypeError{},
							),
						},
					},
				},
				{
					Filepath: "test_packages/proto_conversion/def.pb.go",
					PackageName: "proto_conversion",
					PackageImportPath: "some/import/proto_conversion",
					Consts: []gopkg.DeclVar{
						{
							Name: "_",
							Import: "some/import/proto_conversion",
							Type: gopkg.TypeUnnamedLiteral{},
						},
					},
					Vars: []gopkg.DeclVar{
						{
							Name: "_",
							Import: "some/import/proto_conversion",
							Type: gopkg.TypeUnnamedLiteral{},
						},
						{
							Name: "_",
							Import: "some/import/proto_conversion",
							Type: gopkg.TypeUnnamedLiteral{},
						},
						{
							Name: "_",
							Import: "some/import/proto_conversion",
							Type: gopkg.TypeUnnamedLiteral{},
						},
						{
							Name: "xxx_messageInfo_IntAsString",
							Import: "some/import/proto_conversion",
							Type: gopkg.TypeNamed{
								Name: "InternalMessageInfo",
								Import: "github.com/golang/protobuf/proto",
							},
						},
						{
							Name: "xxx_messageInfo_ShopspringDecimal",
							Import: "some/import/proto_conversion",
							Type: gopkg.TypeNamed{
								Name: "InternalMessageInfo",
								Import: "github.com/golang/protobuf/proto",
							},
						},
						{
							Name: "fileDescriptor_76fb0470a3b910d8",
							Import: "some/import/proto_conversion",
							Type: gopkg.TypeUnnamedLiteral{},
						},
					},
					Functions: protoConversionPackageFuncs(),
					Types: []gopkg.DeclType{
						{
							Name: "IntAsString",
							Import: "some/import/proto_conversion",
							Type: gopkg.TypeStruct{
								Fields: []gopkg.DeclVar{
									{
										Name: "Value",
										Type: gopkg.TypeString{},
										StructTag: "`protobuf:\"bytes,1,opt,name=value,proto3\" json:\"value,omitempty\"`",
									},
									{
										Name: "XXX_NoUnkeyedLiteral",
										Type: gopkg.TypeStruct{},
										StructTag: "`json:\"-\"`",
									},
									{
										Name: "XXX_unrecognized",
										Type: gopkg.TypeArray{
											ValueType: gopkg.TypeByte{},
										},
										StructTag: "`json:\"-\"`",
									},
									{
										Name: "XXX_sizecache",
										Type: gopkg.TypeInt32{},
										StructTag: "`json:\"-\"`",
									},
								},
							},
						},
						{
							Name: "ShopspringDecimal",
							Import: "some/import/proto_conversion",
							Type: gopkg.TypeStruct{
								Fields: []gopkg.DeclVar{
									{
										Name: "Value",
										Type: gopkg.TypeString{},
										StructTag: "`protobuf:\"bytes,1,opt,name=value,proto3\" json:\"value,omitempty\"`",
									},
									{
										Name: "XXX_NoUnkeyedLiteral",
										Type: gopkg.TypeStruct{},
										StructTag: "`json:\"-\"`",
									},
									{
										Name: "XXX_unrecognized",
										Type: gopkg.TypeArray{
											ValueType: gopkg.TypeByte{},
										},
										StructTag: "`json:\"-\"`",
									},
									{
										Name: "XXX_sizecache",
										Type: gopkg.TypeInt32{},
										StructTag: "`json:\"-\"`",
									},
								},
							},
						},
					},
				},
				{
					Filepath: "test_packages/proto_conversion/generate.go",
					PackageName: "proto_conversion",
					PackageImportPath: "some/import/proto_conversion",
				},
			},
		},
		{
			Name: "receiver_funcs",
			PkgDir: "test_packages/receiver_funcs",
			ParseOptions: []gopkg.ParseOption{
				gopkg.ParseWithPkgImportPath("some/import/receiver_funcs"),
			},
			Expected: []gopkg.FileContents{
				{
					Filepath: "test_packages/receiver_funcs/receiver_funcs.go",
					PackageName: "receiver_funcs",
					PackageImportPath: "some/import/receiver_funcs",
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
		{
			Name: "pkg_with_tests",
			PkgDir: "test_packages/pkg_with_tests",
			ParseOptions: []gopkg.ParseOption{
				gopkg.ParseWithPkgImportPath("some/import/pkg_with_tests"),
			},
			Expected: []gopkg.FileContents{
				{
					Filepath: "test_packages/pkg_with_tests/logic.go",
					PackageName: "pkg_with_tests",
					PackageImportPath: "some/import/pkg_with_tests",
					Functions: []gopkg.DeclFunc{
						{
							Name: "MyCoolLogic",
							Import: "some/import/pkg_with_tests",
							Args: []gopkg.DeclVar{
								{
									Name: "i",
									Type: gopkg.TypeInt{},
								},
								{
									Name: "j",
									Type: gopkg.TypeInt{},
								},
							},
							ReturnArgs: tmpl.UnnamedReturnArgs(
								gopkg.TypeInt{},
							),
						},
					},
				},
				{
					Filepath: "test_packages/pkg_with_tests/logic_test.go",
					PackageName: "pkg_with_tests_test",
					PackageImportPath: "some/import/pkg_with_tests",
					Functions: []gopkg.DeclFunc{
						{
							Name: "TestMyCoolLogic",
							Import: "some/import/pkg_with_tests",
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
						},
					},
				},
			},
		},
		{
			Name: "named_return_args",
			PkgDir: "test_packages/named_return_args",
			ParseOptions: []gopkg.ParseOption{
				gopkg.ParseWithPkgImportPath("myimport/named_return_args"),
			},
			Expected: []gopkg.FileContents{
				{
					Filepath: "test_packages/named_return_args/named_return_args.go",
					PackageName: "named_return_args",
					PackageImportPath: "myimport/named_return_args",
					Types: []gopkg.DeclType{
						{
							Name: "SomeInterface",
							Import: "myimport/named_return_args",
							Type: gopkg.TypeInterface{
								Funcs: []gopkg.DeclFunc{
									{
										Name: "SomeFunc",
										ReturnArgs: []gopkg.DeclVar{
											{
												Name: "a",
												Type: gopkg.TypeInt{},
											},
											{
												Name: "b",
												Type: gopkg.TypeError{},
											},
										},
									},
									{
										Name: "OtherFunc",
										ReturnArgs: []gopkg.DeclVar{
											{
												Name: "c",
												Type: gopkg.TypeInt64{},
											},
											{
												Name: "d",
												Type: gopkg.TypeInt64{},
											},
										},
									},
								},
							},
						},
					},
					Functions: []gopkg.DeclFunc{
						{
							Name: "MyMethod",
							Import: "myimport/named_return_args",
							ReturnArgs: []gopkg.DeclVar{
								{
									Name: "e",
									Type: gopkg.TypeInt32{},
								},
								{
									Name: "f",
									Type: gopkg.TypeInt32{},
								},
								{
									Name: "g",
									Type: gopkg.TypeInt32{},
								},
							},
						},
						{
							Name: "MyOtherMethod",
							Import: "myimport/named_return_args",
							ReturnArgs: []gopkg.DeclVar{
								{
									Name: "i",
									Type: gopkg.TypeInt32{},
								},
								{
									Name: "j",
									Type: gopkg.TypeFloat64{},
								},
								{
									Name: "k",
									Type: gopkg.TypeError{},
								},
							},
						},
					},
				},
			},
		},
		{
			Name: "structs_with_tags",
			PkgDir: "test_packages/struct_with_tags",
			ParseOptions: []gopkg.ParseOption{
				gopkg.ParseWithPkgImportPath("myimport/struct_with_tags"),
			},
			Expected: []gopkg.FileContents{
				{
					Filepath: "test_packages/struct_with_tags/acouple_of_structs.go",
					PackageName: "struct_with_tags",
					PackageImportPath: "myimport/struct_with_tags",
					Types: []gopkg.DeclType{
						{
							Name: "AStruct",
							Import: "myimport/struct_with_tags",
							Type: gopkg.TypeStruct{
								Fields: []gopkg.DeclVar{
									{
										Name: "AField",
										Type: gopkg.TypeInt{},
										StructTag: "`AKey:\"some_value\"`",
									},
									{
										Name: "BField",
										Type: gopkg.TypeBool{},
										StructTag: "`BKey:\"some_other_value\"`",
									},
									{
										Name: "privateField",
										Type: gopkg.TypeFloat32{},
										StructTag: "`CKey:\"some_third_value\"`",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {

			pc, err := gopkg.Parse(test.PkgDir, test.ParseOptions...)
			require.NoError(t, err)

			assert.Equal(t, test.Expected, pc)
		})
	}
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
			ReturnArgs: tmpl.UnnamedReturnArgs(
				gopkg.TypeString{},
			),
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
			ReturnArgs: tmpl.UnnamedReturnArgs(
				gopkg.TypeArray{ValueType: gopkg.TypeByte{}},
				gopkg.TypeArray{ValueType: gopkg.TypeInt{}},
			),
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
			ReturnArgs: tmpl.UnnamedReturnArgs(
				gopkg.TypeError{},
			),
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
			ReturnArgs: tmpl.UnnamedReturnArgs(
				gopkg.TypeArray{ValueType: gopkg.TypeByte{}},
				gopkg.TypeError{},
			),
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
					Type: gopkg.TypeNamed{
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
			ReturnArgs: tmpl.UnnamedReturnArgs(
				gopkg.TypeInt{},
			),
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
			ReturnArgs: tmpl.UnnamedReturnArgs(
				gopkg.TypeString{},
			),
		},
	}
}
