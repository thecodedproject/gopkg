package gopkg

import (
	"errors"
	"fmt"
	"go/parser"
	"go/token"
	"go/ast"
	"path"
	"sort"

	//"reflect"
)

const CURRENT_PKG = "current_pkg_import"

func Parse(pkgDir string, pkgImportPath string) ([]FileContents, error) {

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

	// TODO remove this requirement - need to be able to parse a package which has tests in as well!
	if len(pkgs) != 1 {
		for k := range pkgs {
			fmt.Println(k)
		}

		return nil, fmt.Errorf("more than one package found in dir %s", pkgDir)
	}


	pkgContents := make([]FileContents, 0)

	for _, pkg := range pkgs {


		for filepath, fileNode := range pkg.Files {

			fileContents, err := parseNodeAst(pkgImportPath, fileNode)
			if err != nil {
				return nil, err
			}

			fileContents.PackageName = pkg.Name
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

func parseNodeAst(
	pkgImportPath string,
	p ast.Node,
) (FileContents, error) {

	var pc FileContents

	currentFileImports := make(map[string]string)

	var inspectingErr error
	ast.Inspect(p, func(node ast.Node) bool {

		// If we have encountered an error stop parsing the AST asap (by stopping
		// any more recursion into the ast)
		if inspectingErr != nil {
			return false
		}

		switch n := node.(type) {

			case *ast.ImportSpec:
				addImport(currentFileImports, n)
				return true

			case *ast.File:

				//fmt.Println("File:")

				//for filepath, fileObj := range p.Files {
				//	if n == fileObj {
				//		fmt.Println("File path:", filepath)
				//	}
				//}

				//fmt.Printf("Previous imports: %+v\n", currentFileImports)
				currentFileImports = make(map[string]string)
				currentFileImports[CURRENT_PKG] = pkgImportPath
				return true

			case *ast.FuncDecl:

				receiver, err := getFuncReceiverFromFieldList(n.Recv)
				if err != nil {
					inspectingErr = err
					return false
				}

				args, err := getDeclVarsFromFieldList(currentFileImports, n.Type.Params)
				if err != nil {
					inspectingErr = err
					return false
				}

				retArgs, err := getArgTypeList(currentFileImports, n.Type.Results)
				if err != nil {
					inspectingErr = err
					return false
				}

				f := DeclFunc{
					Name: n.Name.String(),
					Import: pkgImportPath,
					Receiver: receiver,
					Args: args,
					ReturnArgs: retArgs,
				}

				pc.Functions = append(pc.Functions, f)
				return true

			case *ast.GenDecl:

				for _, declSpec := range n.Specs {

					switch s := declSpec.(type) {
					case *ast.TypeSpec:

						fullType, err := getFullType(currentFileImports, s.Type)
						if err != nil {
							inspectingErr = err
							return false
						}

						pc.Types = append(
							pc.Types,
							DeclType{
								Name: s.Name.Name,
								Import: pkgImportPath,
								Type: fullType,
							},
						)
					case *ast.ValueSpec:
						declVars, err := declVarsFromAstValueSpec(pkgImportPath, s)
						if err != nil {
							inspectingErr = err
							return false
						}
						pc.Vars = append(pc.Vars, declVars...)
					}
				}
				return false

			default:
				return true
		}
	})

	if inspectingErr != nil {
		return FileContents{}, inspectingErr
	}

	return pc, nil
}

// getArgTypeList gets an order list of arguments from an `ast.FieldList`
//
// Used to get either the types of the parameters arguments for a function,
// or the return arguments for a function whilst parsing the ast.
func getArgTypeList(
	imports map[string]string,
	fieldList *ast.FieldList,
) ([]Type, error) {

	if fieldList == nil || fieldList.List == nil {
		return nil, nil
	}

	typeList := make([]Type, 0, len(fieldList.List))

	for i := range fieldList.List {
		fieldType, err := getFullType(imports, fieldList.List[i].Type)
		if err != nil {
			return nil, err
		}
		typeList = append(typeList, fieldType)
	}

	return typeList, nil
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

		var name string
		if len(f.Names) > 0 {
			name = f.Names[0].String()
		}

		typeList = append(typeList, DeclVar{
			Name: name,
			Type: fieldType,
		})
	}

	return typeList, nil
}

func getDeclFuncsFromFieldList(
	imports map[string]string,
	fieldList *ast.FieldList,
) ([]DeclFunc, error) {

	funcs := make([]DeclFunc, 0, len(fieldList.List))

	for i, method := range fieldList.List {

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

				retArgs, err := getArgTypeList(imports, funcType.Results)
				if err != nil {
					return nil, err
				}

				funcs = append(funcs, DeclFunc{
					Name: fieldList.List[i].Names[0].String(),
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

		t, ok := p.ValueType.(TypeUnknownNamed)
		if !ok {
			return FuncReceiver{}, errors.New("expected TypeUnknownNamed in pointer receiver but found different type")
		}
		receiver.TypeName = t.Name
	} else {
		t, ok := receiverType.Type.(TypeUnknownNamed)
		if !ok {
			return FuncReceiver{}, errors.New("expected TypeUnknownNamed in receiver but found different type")
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

				//return TypeUnknownNamed{
				//	Name: t.Name,
				//}, nil
			}

			importPath := imports[CURRENT_PKG]
			return TypeUnknownNamed{
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

			return TypeUnknownNamed{
				Name: t.Sel.Name,
				Import: importPath,
			}, nil

		case *ast.StructType:

			structFields, err := getDeclVarsFromFieldList(
				imports,
				t.Fields,
			)
			if err != nil {
				return nil, err
			}

			return TypeStruct{
				Fields: structFields,
			}, nil

		case *ast.InterfaceType:

			interfaceFuncs, err := getDeclFuncsFromFieldList(
				imports,
				t.Methods,
			)
			if err != nil {
				return nil, err
			}

			return TypeInterface{
				Funcs: interfaceFuncs,
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

func addImport(imports map[string]string, n *ast.ImportSpec) {

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
	spec *ast.ValueSpec,
) ([]DeclVar, error) {

	var sType Type
	if spec.Type == nil {
		sType = TypeUnnamedLiteral{}
	} else {
		var err error
		sType, err = getFullType(nil, spec.Type)
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
