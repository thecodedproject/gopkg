package gopkg

import (
	"errors"
)

type FileContents struct {

	Filepath string

	PackageName string
	Imports []ImportAndAlias

	Types []DeclType

	Functions []DeclFunc

	// TODO Deprecate - should use []DeclTypes instead
	StructTypes []DeclStruct
}

type ImportAndAlias struct {
	Import string
	Alias string
}

type Type interface {
	DefaultInit(importAliases map[string]string) (string, error)
	FullType(importAliases map[string]string) string
}

type TypeArray struct {
	ValueType Type
}

func (t TypeArray) DefaultInit(importAliases map[string]string) (string, error) {
	return "nil", nil
}

func (t TypeArray) FullType(importAliases map[string]string) string {
	return "[]" + t.ValueType.FullType(importAliases)
}

type TypeBool struct {}

func (t TypeBool) DefaultInit(importAliases map[string]string) (string, error) {
	return "false", nil
}

func (t TypeBool) FullType(importAliases map[string]string) string {
	return "bool"
}

type TypeByte struct {}

func (t TypeByte) DefaultInit(importAliases map[string]string) (string, error) {
	return "0", nil
}

func (t TypeByte) FullType(importAliases map[string]string) string {
	return "byte"
}

type TypeError struct {}

func (t TypeError) DefaultInit(importAliases map[string]string) (string, error) {
	return "nil", nil
}

func (t TypeError) FullType(importAliases map[string]string) string {
	return "error"
}

type TypeFloat32 struct {}

func (t TypeFloat32) DefaultInit(importAliases map[string]string) (string, error) {
	return "0", nil
}

func (t TypeFloat32) FullType(importAliases map[string]string) string {
	return "float32"
}

type TypeFloat64 struct {}

func (t TypeFloat64) DefaultInit(importAliases map[string]string) (string, error) {
	return "0", nil
}

func (t TypeFloat64) FullType(importAliases map[string]string) string {
	return "float64"
}

type TypeInt struct {}

func (t TypeInt) DefaultInit(importAliases map[string]string) (string, error) {
	return "0", nil
}

func (t TypeInt) FullType(importAliases map[string]string) string {
	return "int"
}

type TypeInterface struct {
	Funcs []DeclFunc
}

func (t TypeInterface) DefaultInit(importAliases map[string]string) (string, error) {
	return "nil", nil
}

func (t TypeInterface) FullType(importAliases map[string]string) string {

	if len(t.Funcs) == 0 {
		return "interface{}"
	}

	ret := "interface {\n"
	for _, f := range t.Funcs {
		ret += "\t" + f.Name + funcArgsAndRetArgs(f, nil, false) + "\n"
	}
	ret += "}"

	return ret
}

type TypeInt32 struct {}

func (t TypeInt32) DefaultInit(importAliases map[string]string) (string, error) {
	return "0", nil
}

func (t TypeInt32) FullType(importAliases map[string]string) string {
	return "int32"
}

type TypeInt64 struct {}

func (t TypeInt64) DefaultInit(importAliases map[string]string) (string, error) {
	return "0", nil
}

func (t TypeInt64) FullType(importAliases map[string]string) string {
	return "int64"
}

type TypeString struct {}

func (t TypeString) DefaultInit(importAliases map[string]string) (string, error) {
	return "\"\"", nil
}

func (t TypeString) FullType(importAliases map[string]string) string {
	return "string"
}

type TypeStruct struct {
	Fields []DeclVar
}

func (t TypeStruct) DefaultInit(importAliases map[string]string) (string, error) {
	return "{}", nil
}

func (t TypeStruct) FullType(importAliases map[string]string) string {

	ret := "struct {"

	for i, f := range t.Fields {
		if i == 0 {
			ret += "\n"
		}
		ret += "\t" + f.Name + " " + f.FullType(importAliases) + "\n"
	}

	ret += "}"

	return ret
}

// TODO rename to something more approriate - maybe TypeNamed (or TypeAlias)
type TypeUnknownNamed struct {
	Name string
	Import string
	ValueType Type
}

func (t TypeUnknownNamed) DefaultInit(importAliases map[string]string) (string, error) {

	if t.ValueType != nil {

		switch t.ValueType.(type) {
		case TypeStruct:
			return t.FullType(importAliases) + "{}", nil
		default:
			return t.ValueType.DefaultInit(importAliases)
		}
	}

	return "", errors.New("cannot deduce default init for named type with no value type")
}

func (t TypeUnknownNamed) FullType(importAliases map[string]string) string {
	if alias, hasAlias := importAliases[t.Import]; hasAlias {
		return alias + "." + t.Name
	}

	return t.Name
}

type TypeMap struct {
	KeyType Type
	ValueType Type
}

func (t TypeMap) DefaultInit(importAliases map[string]string) (string, error) {
	return "nil", nil
}

func (t TypeMap) FullType(importAliases map[string]string) string {
	return "map[" + t.KeyType.FullType(importAliases) + "]" + t.ValueType.FullType(importAliases)
}

type TypePointer struct {
	ValueType Type
}

func (t TypePointer) DefaultInit(importAliases map[string]string) (string, error) {
	return "nil", nil
}

func (t TypePointer) FullType(importAliases map[string]string) string {
	return "*" + t.ValueType.FullType(importAliases)
}
