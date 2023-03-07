package gopkg

import (
	"errors"
	"go/parser"
	"go/token"
	"go/ast"
	"path"
	"sort"
)

const CURRENT_PKG = "current_pkg_import"

func Parse(pkgDir string, opts ...ParseOption) ([]FileContents, error) {

	var parseOptions parseOptions
	for _, opt := range opts {
		parseOptions = opt(parseOptions)
	}

	if parseOptions.pkgImportPath == "" {
		var err error
		parseOptions.pkgImportPath, err = PackageImportPath(pkgDir)
		if err != nil {
			return nil, err
		}
	}

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(
		fset,
		pkgDir,
		nil,
		0,
	)
	if err != nil {
		return nil, err
	}

	pkgContents := make([]FileContents, 0)

	for _, pkg := range pkgs {


		for filepath, fileNode := range pkg.Files {

			fileContents, err := fileContentsFromAstFile(
				parseOptions.pkgImportPath,
				fileNode,
			)
			if err != nil {
				return nil, err
			}

			fileContents.PackageName = pkg.Name
			fileContents.PackageImportPath = parseOptions.pkgImportPath
			fileContents.Filepath = filepath

			pkgContents = append(
				pkgContents,
				fileContents,
			)
		}
	}

	sort.Slice(pkgContents, func(i, j int) bool {
		return pkgContents[i].Filepath < pkgContents[j].Filepath
	})

	return pkgContents, nil
}

func fileContentsFromAstFile(
	pkgImportPath string,
	f *ast.File,
) (FileContents, error) {

	fileImports := make(map[string]string)
	fileImports[CURRENT_PKG] = pkgImportPath
	for _, importSpec := range f.Imports {
		addImport(fileImports, importSpec)
	}

	var contents FileContents
	for _, d := range f.Decls {
		switch decl := d.(type) {
		case *ast.FuncDecl:
			receiver, err := getFuncReceiverFromFieldList(decl.Recv)
			if err != nil {
				return FileContents{}, err
			}

			args, err := getDeclVarsFromFieldList(fileImports, decl.Type.Params)
			if err != nil {
				return FileContents{}, err
			}

			retArgs, err := getDeclVarsFromFieldList(fileImports, decl.Type.Results)
			if err != nil {
				return FileContents{}, err
			}

			f := DeclFunc{
				Name: decl.Name.String(),
				Import: pkgImportPath,
				Receiver: receiver,
				Args: args,
				ReturnArgs: retArgs,
			}

			contents.Functions = append(contents.Functions, f)

		case *ast.GenDecl:

			for _, declSpec := range decl.Specs {

				switch s := declSpec.(type) {
				case *ast.TypeSpec:

					fullType, err := getFullType(fileImports, s.Type)
					if err != nil {
						return FileContents{}, err
					}

					contents.Types = append(
						contents.Types,
						DeclType{
							Name: s.Name.Name,
							Import: pkgImportPath,
							Type: fullType,
						},
					)
				case *ast.ValueSpec:

					declVars, err := declVarsFromAstValueSpec(
						pkgImportPath,
						fileImports,
						s,
					)
					if err != nil {
						return FileContents{}, err
					}
					if decl.Tok == token.VAR {
						contents.Vars = append(contents.Vars, declVars...)
					} else if decl.Tok == token.CONST {
						contents.Consts = append(contents.Consts, declVars...)
					}
				}
			}
		}
	}

	return contents, nil
}

// getDeclVarsFromFieldList returns an ordered list of declared variables
//
// The ast field list might be, for example, the list of arguments passed into
// a function.
// It returns the underlying type (as `Type`) as well as the name of the
// declared variable.
// Note that `DeclVar.Import` will always be blank as the field list will
// only contain vars declared in a local scope (i.e. not at the package level)
func getDeclVarsFromFieldList(
	imports map[string]string,
	fieldList *ast.FieldList,
) ([]DeclVar, error) {

	if fieldList == nil || fieldList.List == nil {
		return nil, nil
	}

	typeList := make([]DeclVar, 0, len(fieldList.List))

	for _, f := range fieldList.List {
		fieldType, err := getFullType(imports, f.Type)
		if err != nil {
			return nil, err
		}

		if len(f.Names) == 0 {
			typeList = append(typeList, DeclVar{
				Type: fieldType,
			})
		} else {
			for _, name := range f.Names {
				typeList = append(typeList, DeclVar{
					Name: name.String(),
					Type: fieldType,
				})
			}
		}
	}

	return typeList, nil
}

func getDeclFuncsFromFieldList(
	imports map[string]string,
	fieldList *ast.FieldList,
) ([]DeclFunc, error) {

	funcs := make([]DeclFunc, 0, len(fieldList.List))

	for _, method := range fieldList.List {

		for _, name := range method.Names {

			if name.Obj != nil && name.Obj.Kind == ast.Fun {

				funcDecl, ok := name.Obj.Decl.(*ast.Field)
				if !ok {
					return nil, errors.New("bad func decl")
				}

				funcType, ok := funcDecl.Type.(*ast.FuncType)
				if !ok {
					return nil, errors.New("bad func decl")
				}

				args, err := getDeclVarsFromFieldList(imports, funcType.Params)
				if err != nil {
					return nil, err
				}

				retArgs, err := getDeclVarsFromFieldList(imports, funcType.Results)
				if err != nil {
					return nil, err
				}

				funcs = append(funcs, DeclFunc{
					Name: name.String(),
					Args: args,
					ReturnArgs: retArgs,
				})
			}
		}
	}

	return funcs, nil
}

func getFuncReceiverFromFieldList(
	fieldList *ast.FieldList,
) (FuncReceiver, error) {

	if fieldList == nil {
		return FuncReceiver{}, nil
	}

	types, err := getDeclVarsFromFieldList(nil, fieldList)
	if err != nil {
		return FuncReceiver{}, err
	}

	if len(types) != 1 {
		return FuncReceiver{}, errors.New("More than one receiver in ast for method")
	}

	receiverType := types[0]

	receiver := FuncReceiver{
		VarName: receiverType.Name,
	}
	if p, ok := receiverType.Type.(TypePointer); ok {
		receiver.IsPointer = true

		t, ok := p.ValueType.(TypeNamed)
		if !ok {
			return FuncReceiver{}, errors.New("expected TypeNamed in pointer receiver but found different type")
		}
		receiver.TypeName = t.Name
	} else {
		t, ok := receiverType.Type.(TypeNamed)
		if !ok {
			return FuncReceiver{}, errors.New("expected TypeNamed in receiver but found different type")
		}
		receiver.TypeName = t.Name
	}

	return receiver, nil
}

func getFullType(
	imports map[string]string,
	t ast.Expr,
) (Type, error) {

	//fmt.Println("******", reflect.TypeOf(t))

	switch t := t.(type) {
		case *ast.ArrayType:
			if t.Len != nil {
				return nil, errors.New("[...]T array types not supported")
			}
			fullType, err := getFullType(imports, t.Elt)
			if err != nil {
				return nil, err
			}
			return TypeArray{
				ValueType: fullType,
			}, nil

		case *ast.Ident:
			if isBuiltInType(t.Name) {
				return typeFromString(t.Name), nil
			}

			importPath := imports[CURRENT_PKG]
			return TypeNamed{
				Name: t.Name,
				Import: importPath,
			}, nil

		case *ast.StarExpr:
			fullType, err := getFullType(imports, t.X)
			if err != nil {
				return nil, err
			}
			return TypePointer{
				ValueType: fullType,
			}, nil

		// i.e. an expression selecting something from another package
		//	`some_pkg.SomeType`
		case *ast.SelectorExpr:
			imp, ok := t.X.(*ast.Ident)

			if !ok {
				return nil, errors.New("uknown selector X")
			}

			importPath, ok := imports[imp.Name]
			if !ok {
				return nil, errors.New("unknown import path " + imp.Name)
			}

			//fmt.Println("****** Type:", importPath, importPrefix + "." + t.Sel.Name)

			return TypeNamed{
				Name: t.Sel.Name,
				Import: importPath,
			}, nil

		case *ast.StructType:

			structFieldsAndEmbeds, err := getDeclVarsFromFieldList(
				imports,
				t.Fields,
			)
			if err != nil {
				return nil, err
			}

			var s TypeStruct
			for _, f := range structFieldsAndEmbeds {
				if f.Name == "" {
					s.Embeds = append(s.Embeds, f.Type)
				} else {
					s.Fields = append(s.Fields, f)
				}
			}

			return s, nil

		case *ast.InterfaceType:

			interfaceFuncs, err := getDeclFuncsFromFieldList(
				imports,
				t.Methods,
			)
			if err != nil {
				return nil, err
			}

			i := TypeInterface{
				Funcs: interfaceFuncs,
			}

			// Embedded types or interfaces will appear in the field list
			// without any name
			possibleEmbeds, err := getDeclVarsFromFieldList(
				imports,
				t.Methods,
			)
			if err != nil {
				return nil, err
			}

			for _, f := range possibleEmbeds {
				if f.Name == "" {
					i.Embeds = append(i.Embeds, f.Type)
				}
			}

			return i, nil

		case *ast.FuncType:

			args, err := getDeclVarsFromFieldList(imports, t.Params)
			if err != nil {
				return nil, err
			}
			retArgs, err := getDeclVarsFromFieldList(imports, t.Results)
			if err != nil {
				return nil, err
			}

			return TypeFunc{
				Args: args,
				ReturnArgs: retArgs,
			}, nil

		default:
			return nil, errors.New("unknown field type")
	}
}

func removeQuotes(s string) string {

	if s[0] == '"' {
		s = s[1:]
	}
	if s[len(s)-1] == '"' {
		s = s[:len(s)-1]
	}
	return s
}

func addImport(
	imports map[string]string,
	n *ast.ImportSpec,
) {

	importPath := removeQuotes(n.Path.Value)
	var localName string
	if n.Name != nil {
		localName = n.Name.String()
	}
	if localName == "." {
		panic("'.' imports are not supported")
	}
	if localName == "" {
		_, localName = path.Split(importPath)
	}

	imports[localName] = importPath
}

func isBuiltInType(t string) bool {

	builtInTypes := map[string]struct{}{
		"bool": {},
		"byte": {},
		"error": {},
		"float32": {},
		"float64": {},
		"int": {},
		"int32": {},
		"int64": {},
		"string": {},
	}

	_, ok := builtInTypes[t]
	return ok
}

func typeFromString(t string) Type {

	switch t {
	case "bool":
		return TypeBool{}
	case "byte":
		return TypeByte{}
	case "error":
		return TypeError{}
	case "float32":
		return TypeFloat32{}
	case "float64":
		return TypeFloat64{}
	case "int":
		return TypeInt{}
	case "int32":
		return TypeInt32{}
	case "int64":
		return TypeInt64{}
	case "string":
		return TypeString{}
	}
	return nil
}

func declVarsFromAstValueSpec(
	pkgImportPath string,
	imports map[string]string,
	spec *ast.ValueSpec,
) ([]DeclVar, error) {

	var sType Type
	if spec.Type == nil {
		sType = TypeUnnamedLiteral{}
	} else {
		var err error
		sType, err = getFullType(imports, spec.Type)
		if err != nil {
			return nil, err
		}
	}

	var declVars []DeclVar
	hasLiteralValues := (len(spec.Names) == len(spec.Values))
	for iDecl := range spec.Names {

		var literalValue string
		if hasLiteralValues {
			switch litVal := spec.Values[iDecl].(type) {
			case *ast.BasicLit:
				literalValue = litVal.Value
			case *ast.Ident:
				literalValue = litVal.String()
			}
		}

		declVars = append(
			declVars,
			DeclVar{
				Name: spec.Names[iDecl].String(),
				Import: pkgImportPath,
				Type: sType,
				LiteralValue: literalValue,
			},
		)
	}

	return declVars, nil
}

type ParseOption func(parseOptions)parseOptions

type parseOptions struct {
	pkgImportPath string
}

func ParseWithPkgImportPath(importPath string) ParseOption {

	return func(o parseOptions) parseOptions {
		o.pkgImportPath = importPath
		return o
	}
}
