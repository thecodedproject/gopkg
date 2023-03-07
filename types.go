package gopkg

import (
	"errors"
)

type FileContents struct {

	Filepath string

	PackageName string

	PackageImportPath string

	Imports []ImportAndAlias

	Consts []DeclVar
	Vars []DeclVar
	Types []DeclType
	Functions []DeclFunc
}

type ImportAndAlias struct {
	Import string
	Alias string
}

type Type interface {
	DefaultInit(importAliases map[string]string) (string, error)
	FullType(importAliases map[string]string) (string, error)
	RequiredImports() map[string]bool
}

type TypeArray struct {
	ValueType Type
}

func (t TypeArray) DefaultInit(importAliases map[string]string) (string, error) {
	return "nil", nil
}

func (t TypeArray) FullType(importAliases map[string]string) (string, error) {

	valueFullType, err := t.ValueType.FullType(importAliases)
	if err != nil {
		return "", err
	}

	return "[]" + valueFullType, nil
}

func (t TypeArray) RequiredImports() map[string]bool {
	return t.ValueType.RequiredImports()
}

type TypeBool struct {}

func (t TypeBool) DefaultInit(importAliases map[string]string) (string, error) {
	return "false", nil
}

func (t TypeBool) FullType(importAliases map[string]string) (string, error) {
	return "bool", nil
}

func (t TypeBool) RequiredImports() map[string]bool {
	return nil
}

type TypeByte struct {}

func (t TypeByte) DefaultInit(importAliases map[string]string) (string, error) {
	return "0", nil
}

func (t TypeByte) FullType(importAliases map[string]string) (string, error) {
	return "byte", nil
}

func (t TypeByte) RequiredImports() map[string]bool {
	return nil
}

type TypeError struct {}

func (t TypeError) DefaultInit(importAliases map[string]string) (string, error) {
	return "nil", nil
}

func (t TypeError) FullType(importAliases map[string]string) (string, error) {
	return "error", nil
}

func (t TypeError) RequiredImports() map[string]bool {
	return nil
}

type TypeFloat32 struct {}

func (t TypeFloat32) DefaultInit(importAliases map[string]string) (string, error) {
	return "0", nil
}

func (t TypeFloat32) FullType(importAliases map[string]string) (string, error) {
	return "float32", nil
}

func (t TypeFloat32) RequiredImports() map[string]bool {
	return nil
}

type TypeFloat64 struct {}

func (t TypeFloat64) DefaultInit(importAliases map[string]string) (string, error) {
	return "0", nil
}

func (t TypeFloat64) FullType(importAliases map[string]string) (string, error) {
	return "float64", nil
}

func (t TypeFloat64) RequiredImports() map[string]bool {
	return nil
}

type TypeInt struct {}

func (t TypeInt) DefaultInit(importAliases map[string]string) (string, error) {
	return "0", nil
}

func (t TypeInt) FullType(importAliases map[string]string) (string, error) {
	return "int", nil
}

func (t TypeInt) RequiredImports() map[string]bool {
	return nil
}

type TypeInterface struct {
	Embeds []Type
	Funcs []DeclFunc
}

func (t TypeInterface) DefaultInit(importAliases map[string]string) (string, error) {
	return "nil", nil
}

func (t TypeInterface) FullType(importAliases map[string]string) (string, error) {

	if len(t.Embeds) == 0 && len(t.Funcs) == 0 {
		return "interface{}", nil
	}

	ret := "interface {\n"

	for _, e := range t.Embeds {
		eFullType, err := e.FullType(importAliases)
		if err != nil {
			return "", err
		}
		ret += "\t" + eFullType + "\n"
	}

	if len(t.Embeds) > 0 && len(t.Funcs) > 0 {
		ret += "\n"
	}

	for _, f := range t.Funcs {

		argsAndRets, err := funcArgsAndRetArgs(
			f.Args,
			f.ReturnArgs,
			importAliases,
			false,
		)
		if err != nil {
			return "", err
		}

		ret += "\t" + f.Name + argsAndRets + "\n"
	}
	ret += "}"

	return ret, nil
}

func (t TypeInterface) RequiredImports() map[string]bool {
	ret := make(map[string]bool)
	for _, e := range t.Embeds {
		ret = union(ret, e.RequiredImports())
	}
	for _, f := range t.Funcs {
		ret = union(ret, f.RequiredImports())
	}
	return ret
}

type TypeInt32 struct {}

func (t TypeInt32) DefaultInit(importAliases map[string]string) (string, error) {
	return "0", nil
}

func (t TypeInt32) FullType(importAliases map[string]string) (string, error) {
	return "int32", nil
}

func (t TypeInt32) RequiredImports() map[string]bool {
	return nil
}

type TypeInt64 struct {}

func (t TypeInt64) DefaultInit(importAliases map[string]string) (string, error) {
	return "0", nil
}

func (t TypeInt64) FullType(importAliases map[string]string) (string, error) {
	return "int64", nil
}

func (t TypeInt64) RequiredImports() map[string]bool {
	return nil
}

type TypeString struct {}

func (t TypeString) DefaultInit(importAliases map[string]string) (string, error) {
	return "\"\"", nil
}

func (t TypeString) FullType(importAliases map[string]string) (string, error) {
	return "string", nil
}

func (t TypeString) RequiredImports() map[string]bool {
	return nil
}

type TypeStruct struct {
	Embeds []Type
	Fields []DeclVar
}

func (t TypeStruct) DefaultInit(importAliases map[string]string) (string, error) {
	return "{}", nil
}

func (t TypeStruct) FullType(importAliases map[string]string) (string, error) {

	ret := "struct {"

	for i, e := range t.Embeds {
		if i == 0 {
			ret += "\n"
		}
		eFullType, err := e.FullType(importAliases)
		if err != nil {
			return "", err
		}
		ret += "\t" + eFullType + "\n"
	}

	for i, f := range t.Fields {
		if i == 0 {
			ret += "\n"
		}
		fieldFullType, err := f.FullType(importAliases)
		if err != nil {
			return "", err
		}
		ret += "\t" + f.Name + " " + fieldFullType + "\n"
	}

	ret += "}"

	return ret, nil
}

func (t TypeStruct) RequiredImports() map[string]bool {
	ret := make(map[string]bool)
	for _, e := range t.Embeds {
		ret = union(ret, e.RequiredImports())
	}
	for _, f := range t.Fields {
		ret = union(ret, f.RequiredImports())
	}
	return ret
}

// TODO rename to something more approriate - maybe TypeNamed (or TypeAlias)
type TypeNamed struct {
	Name string
	Import string
	ValueType Type
}

func (t TypeNamed) DefaultInit(importAliases map[string]string) (string, error) {

	if t.ValueType != nil {

		switch t.ValueType.(type) {
		case TypeStruct:

			structFullType, err := t.FullType(importAliases)
			if err != nil {
				return "", err
			}

			return structFullType + "{}", nil
		default:
			return t.ValueType.DefaultInit(importAliases)
		}
	}

	return "", errors.New("cannot deduce default init for named type with no value type")
}

func (t TypeNamed) FullType(importAliases map[string]string) (string, error) {
	if alias, hasAlias := importAliases[t.Import]; hasAlias {
		return alias + "." + t.Name, nil
	}

	return t.Name, nil
}

func (t TypeNamed) RequiredImports() map[string]bool {
	if t.Import != "" {
		return map[string]bool {
			t.Import: true,
		}
	}
	return nil
}

type TypeMap struct {
	KeyType Type
	ValueType Type
}

func (t TypeMap) DefaultInit(importAliases map[string]string) (string, error) {
	return "nil", nil
}

func (t TypeMap) FullType(importAliases map[string]string) (string, error) {

	keyFullType, err := t.KeyType.FullType(importAliases)
	if err != nil {
		return "", err
	}

	valueFullType, err := t.ValueType.FullType(importAliases)
	if err != nil {
		return "", err
	}

	return "map[" + keyFullType + "]" + valueFullType, nil
}

func (t TypeMap) RequiredImports() map[string]bool {
	return union(
		t.KeyType.RequiredImports(),
		t.ValueType.RequiredImports(),
	)
}

type TypePointer struct {
	ValueType Type
}

func (t TypePointer) DefaultInit(importAliases map[string]string) (string, error) {
	return "nil", nil
}

func (t TypePointer) FullType(importAliases map[string]string) (string, error) {

	valueFullType, err := t.ValueType.FullType(importAliases)
	if err != nil {
		return "", err
	}

	return "*" + valueFullType, nil
}

func (t TypePointer) RequiredImports() map[string]bool {
	return t.ValueType.RequiredImports()
}

type TypeUnnamedLiteral struct {}

func (t TypeUnnamedLiteral) DefaultInit(importAliases map[string]string) (string, error) {
	return "", errors.New("no default init for unnamed literal")
}

func (t TypeUnnamedLiteral) FullType(importAliases map[string]string) (string, error) {
	return "", nil
}

func (t TypeUnnamedLiteral) RequiredImports() map[string]bool {
	return nil
}

type TypeFunc struct {
	Args []DeclVar
	ReturnArgs []DeclVar
}

func (t TypeFunc) DefaultInit(importAliases map[string]string) (string, error) {
	return "nil", nil
}

func (t TypeFunc) FullType(importAliases map[string]string) (string, error) {

	fullType := "func"

	argsAndRets, err := funcArgsAndRetArgs(
		t.Args,
		t.ReturnArgs,
		importAliases,
		false,
	)
	if err != nil {
		return "", err
	}

	fullType += argsAndRets

	return fullType, nil
}

func (t TypeFunc) RequiredImports() map[string]bool {

	ret := make(map[string]bool)
	for _, a := range t.Args {
		ret = union(ret, a.RequiredImports())
	}
	for _, a := range t.ReturnArgs {
		ret = union(ret, a.RequiredImports())
	}
	return ret
}
