package gopkg

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path"
	"reflect"
	"sort"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/tools/go/packages"

	"fmt"
)

const CURRENT_PKG = "current_pkg_import"

// Parse parses the file or package at `inputPath` and returns its `FileContents` representation
//
// `inputPath` may either be a single golang source file or a directory of golang source files (i.e. a package)
func Parse(inputPath string, opts ...ParseOption) ([]FileContents, error) {

	var parseOptions parseOptions
	for _, opt := range opts {
		parseOptions = opt(parseOptions)
	}

	if parseOptions.pkgImportPath == "" {
		var err error
		parseOptions.pkgImportPath, err = PackageImportPath(inputPath)
		if err != nil {
			return nil, err
		}
	}

	fileInfo, err := os.Stat(inputPath)
	if err != nil {
		return nil, err
	}

	if fileInfo.IsDir() {
		return parseSingleDirectory(inputPath, parseOptions)
	}

	return parseSingleFile(inputPath, parseOptions)
}

func parseSingleDirectory(
	dir string,
	parseOpts parseOptions,
) ([]FileContents, error) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(
		fset,
		dir,
		nil,
		parser.ParseComments,
	)
	if err != nil {
		return nil, err
	}

	pkgContents := make([]FileContents, 0)

	for _, pkg := range pkgs {

		for filepath, fileNode := range pkg.Files {

			fileContents, err := fileContentsFromAstFile(
				parseOpts,
				filepath,
				fileNode,
				fset,
			)
			if err != nil {
				return nil, err
			}

			fileContents.PackageName = pkg.Name
			fileContents.PackageImportPath = parseOpts.pkgImportPath
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

func parseSingleFile(
	filepath string,
	parseOpts parseOptions,
) ([]FileContents, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(
		fset,
		filepath,
		nil,
		parser.ParseComments,
	)
	if err != nil {
		return nil, err
	}

	fileContents, err := fileContentsFromAstFile(
		parseOpts,
		filepath,
		f,
		fset,
	)
	if err != nil {
		return nil, err
	}

	fileContents.PackageName = f.Name.String()
	fileContents.PackageImportPath = parseOpts.pkgImportPath
	fileContents.Filepath = filepath

	return []FileContents{fileContents}, nil
}

func fileContentsFromAstFile(
	parseOpts parseOptions,
	filepath string,
	f *ast.File,
	fileSet *token.FileSet,
) (FileContents, error) {

	fp, err := os.Open(filepath)
	if err != nil {
		return FileContents{}, err
	}
	defer fp.Close()

	var contents FileContents
	if f.Doc != nil {
		var err error
		contents.DocString, err = readFromFileSet(fp, fileSet, f.Doc.Pos(), f.Doc.End())
		if err != nil {
			return FileContents{}, err
		}
	}

	contents.Imports, err = parseImportsFromAstFile(f)
	if err != nil {
		return FileContents{}, err
	}

	fileImports := buildFileAliasesAndImports(parseOpts.pkgImportPath, contents.Imports)

	for _, d := range f.Decls {
		switch decl := d.(type) {
		case *ast.FuncDecl:
			f, err := getDeclFunc(parseOpts, fileSet, fileImports, fp, decl)
			if err != nil {
				return FileContents{}, err
			}
			contents.Functions = append(contents.Functions, f)

		case *ast.GenDecl:

			var docString string
			if decl.Doc != nil {
				docString, err = readFromFileSet(fp, fileSet, decl.Doc.Pos(), decl.Doc.End())
				if err != nil {
					return FileContents{}, err
				}
			}

			for _, declSpec := range decl.Specs {

				switch s := declSpec.(type) {
				case *ast.TypeSpec:

					fullType, err := getFullType(parseOpts, fileImports, s.Type)
					if err != nil {
						return FileContents{}, err
					}

					contents.Types = append(
						contents.Types,
						DeclType{
							Name:      s.Name.Name,
							Import:    parseOpts.pkgImportPath,
							Type:      fullType,
							DocString: docString,
						},
					)
				case *ast.ValueSpec:

					declVars, err := declVarsFromAstValueSpec(
						parseOpts,
						fileImports,
						fp,
						fileSet,
						s,
					)
					if err != nil {
						return FileContents{}, err
					}

					for i := range declVars {
						if declVars[i].DocString == "" {
							declVars[i].DocString = docString
						}
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

func parseImportsFromAstFile(
	fileAst *ast.File,
) ([]ImportAndAlias, error) {

	if len(fileAst.Imports) == 0 {
		return nil, nil
	}

	imports := make([]ImportAndAlias, 0, len(fileAst.Imports))
	for _, importSpec := range fileAst.Imports {
		i, err := parseImportSpec(importSpec)
		if err != nil {
			return nil, err
		}
		imports = append(imports, i)
	}
	return imports, nil
}

// buildFileAliasesAndImports builds a map of _local aliases_ to their imports.
// A local alias is the alias used to reference any declaration from a given
// import within a file - this is always non-empty for every import except
// for the current package.
// i.e. within a golang file, wherever a declarative element is used there is
// an alias to indicate which package this element comes from.
// If the package was imported without an explicit alias, then the local alias
// assumed to be the last element of import path.
// For the current package, a special constant `CURRENT_PKG` is used to indicate
// that this is the current package.
// This mapping is used to assign the correct import path to every parsed
// declaration (the local aliases are not returned in the parsed FileContents
// struct)
func buildFileAliasesAndImports(
	currentPkgImportPath string,
	imports []ImportAndAlias,
) map[string]string {
	fileImports := make(map[string]string)
	fileImports[CURRENT_PKG] = currentPkgImportPath
	for _, i := range imports {
		localAlias := i.Alias
		if localAlias == "" {
			_, localAlias = path.Split(i.Import)
		}

		fileImports[localAlias] = i.Import
	}

	return fileImports
}

func getDeclFunc(
	parseOpts parseOptions,
	fileSet *token.FileSet,
	fileImports map[string]string,
	fp *os.File,
	decl *ast.FuncDecl,
) (DeclFunc, error) {

	receiver, err := getFuncReceiverFromFieldList(parseOpts, decl.Recv)
	if err != nil {
		return DeclFunc{}, err
	}

	variadicLastArg := handleVariadicLastArg(decl.Type.Params)

	args, err := getDeclVarsFromFieldList(parseOpts, fileImports, decl.Type.Params)
	if err != nil {
		return DeclFunc{}, err
	}

	retArgs, err := getDeclVarsFromFieldList(parseOpts, fileImports, decl.Type.Results)
	if err != nil {
		return DeclFunc{}, err
	}

	f := DeclFunc{
		Name:       decl.Name.String(),
		Import:     parseOpts.pkgImportPath,
		Receiver:   receiver,
		Args:       args,
		ReturnArgs: retArgs,
		VariadicLastArg: variadicLastArg,
	}

	if decl.Body != nil {
		body, err := readFromFileSet(fp, fileSet, decl.Body.Lbrace+1, decl.Body.Rbrace)
		if err != nil {
			return DeclFunc{}, err
		}
		if body != "\n" {
			f.BodyTmpl = body
		}
	}

	if decl.Doc != nil {
		docString, err := readFromFileSet(fp, fileSet, decl.Doc.Pos(), decl.Doc.End())
		if err != nil {
			return DeclFunc{}, err
		}
		f.DocString = docString
	}

	return f, nil
}

// readFromFileSet read bytes from the open *os.File, `fp`, from the
// the byte at position `from` in the fileset upto, but not including, the
// byte at `to` in the fileset.
func readFromFileSet(
	fp *os.File,
	fileSet *token.FileSet,
	from token.Pos,
	to token.Pos,
) (string, error) {

	fsFile := fileSet.File(from)
	if fsFile == nil {
		return "", errors.New("position is not in the fileset")
	}

	buf := make([]byte, int64(to-from))
	_, err := fp.ReadAt(buf, int64(from)-int64(fsFile.Base()))
	if err != nil {
		return "", err
	}

	return string(buf), nil
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
	parseOpts parseOptions,
	imports map[string]string,
	fieldList *ast.FieldList,
) ([]DeclVar, error) {

	if fieldList == nil || fieldList.List == nil {
		return nil, nil
	}

	typeList := make([]DeclVar, 0, len(fieldList.List))

	for _, f := range fieldList.List {
		fieldType, err := getFullType(parseOpts, imports, f.Type)
		if err != nil {
			return nil, err
		}

		var tag string
		if f.Tag != nil {
			tag = f.Tag.Value
			tag = strings.TrimPrefix(tag, "`")
			tag = strings.TrimSuffix(tag, "`")
		}

		if len(f.Names) == 0 {
			typeList = append(typeList, DeclVar{
				Type: fieldType,
			})
		} else {
			for _, name := range f.Names {
				typeList = append(typeList, DeclVar{
					Name:      name.String(),
					Type:      fieldType,
					StructTag: reflect.StructTag(tag),
				})
			}
		}
	}

	return typeList, nil
}

func getDeclFuncsFromFieldList(
	parseOpts parseOptions,
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

				variadicLastArg := handleVariadicLastArg(funcType.Params)

				args, err := getDeclVarsFromFieldList(parseOpts, imports, funcType.Params)
				if err != nil {
					return nil, err
				}

				retArgs, err := getDeclVarsFromFieldList(parseOpts, imports, funcType.Results)
				if err != nil {
					return nil, err
				}

				funcs = append(funcs, DeclFunc{
					Name:       name.String(),
					Args:       args,
					ReturnArgs: retArgs,
					VariadicLastArg: variadicLastArg,
				})
			}
		}
	}

	return funcs, nil
}

// handleVariadicLastArg will detect if the last parameter of a func type is variadic
//
// If the last arg is variadic it will return true *and* will _strip_ the `ast.Ellipsis`
// from the type of of the last arg in `funcParams` - this is so it can be passed as if
// it is not variadic (i.e. with the same `getFullType` method)
//
// If the last arg is not variabic it will return false and `funcParams` will remain unchanged
func handleVariadicLastArg(funcParams *ast.FieldList) bool {

	// todo: implement

	if len(funcParams.List) == 0 {
		return false
	}

	p := funcParams.List[len(funcParams.List) - 1]

	e, isVariadic := p.Type.(*ast.Ellipsis)

	if !isVariadic {
		return false
	}

	funcParams.List[len(funcParams.List) - 1].Type = e.Elt

	return true
}

func getFuncReceiverFromFieldList(
	parseOpts parseOptions,
	fieldList *ast.FieldList,
) (FuncReceiver, error) {

	if fieldList == nil {
		return FuncReceiver{}, nil
	}

	types, err := getDeclVarsFromFieldList(parseOpts, nil, fieldList)
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
	parseOpts parseOptions,
	imports map[string]string,
	t ast.Expr,
) (Type, error) {

	//fmt.Println("******", reflect.TypeOf(t))

	switch t := t.(type) {
	case *ast.ArrayType:
		if t.Len != nil {
			return nil, errors.New("[...]T array types not supported")
		}
		fullType, err := getFullType(parseOpts, imports, t.Elt)
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
			Name:   t.Name,
			Import: importPath,
		}, nil

	case *ast.MapType:

		keyType, err := getFullType(parseOpts, imports, t.Key)
		if err != nil {
			return nil, err
		}

		valueType, err := getFullType(parseOpts, imports, t.Value)
		if err != nil {
			return nil, err
		}

		return TypeMap{
			KeyType: keyType,
			ValueType: valueType,
		}, nil

	case *ast.StarExpr:
		fullType, err := getFullType(parseOpts, imports, t.X)
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
			return nil, errors.New("unknown import path '" + imp.Name + "'")
		}

		var valueType Type
		if parseOpts.dependentTypes {

			fmt.Println("******* Found dependent type:", importPath, t.Sel.Name)

			dependentImports, depTypeExpr, err := fetchFileImportsAndAstForDependentType(importPath, t.Sel.Name)
			if err != nil {
				return nil, errors.Wrap(err, "failed to get ast for dependent type")
			}


			// TODO:
			// * Stop parsing nested dependent types (need to change the parseOpts passed here to turn off dep types)
			valueType, err = getFullType(parseOptions{}, dependentImports, depTypeExpr)
			if err != nil {
				return nil, errors.Wrap(err, "failed to get full type for dependent type")
			}

		  //get dependent type:
		 	//	Parse importPath with `packages`
		 	//	look for type with `t.Sel.Name`
		 	//	Parse ast of that type
		}


		//fmt.Println("****** Type:", importPath, importPrefix + "." + t.Sel.Name)

		return TypeNamed{
			Name:   t.Sel.Name,
			Import: importPath,
			ValueType: valueType,
		}, nil

	case *ast.StructType:

		structFieldsAndEmbeds, err := getDeclVarsFromFieldList(
			parseOpts,
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
			parseOpts,
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
			parseOpts,
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

		variadicLastArg := handleVariadicLastArg(t.Params)

		args, err := getDeclVarsFromFieldList(parseOpts, imports, t.Params)
		if err != nil {
			return nil, err
		}
		retArgs, err := getDeclVarsFromFieldList(parseOpts, imports, t.Results)
		if err != nil {
			return nil, err
		}

		return TypeFunc{
			Args:       args,
			ReturnArgs: retArgs,
			VariadicLastArg: variadicLastArg,
		}, nil

	default:
		return nil, errors.New("unknown field type")
	}
}

func fetchFileImportsAndAstForDependentType(
	pkgPath string,
	typeName string,
) (map[string]string, ast.Expr, error) {

	conf := &packages.Config{
		Mode: packages.NeedDeps |
			packages.NeedFiles |
			packages.NeedName |
			packages.NeedSyntax,
	}
	pkgs, err := packages.Load(conf, pkgPath)
	if err != nil {
		return nil, nil, err
	}

	if len(pkgs) == 0 {
		return nil, nil, errors.New("package not found '" + pkgPath + "'")
	}

	for _, fileAst := range pkgs[0].Syntax {
		for _, decl := range fileAst.Decls {

			switch d := decl.(type) {
			case *ast.GenDecl:

				for _, declSpec := range d.Specs {

					switch s := declSpec.(type) {
					case *ast.TypeSpec:

						if s.Name.Name == typeName {

							imports, err := parseImportsFromAstFile(fileAst)
							if err != nil {
								return nil, nil, err
							}

							fileImports := buildFileAliasesAndImports(pkgPath, imports)

							return fileImports, s.Type, nil
						}
					}
				}
			}
		}
	}

	return nil, nil, errors.New("type '" + typeName + "' not found in packge '" + pkgPath + "'")
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

func parseImportSpec(
	n *ast.ImportSpec,
) (ImportAndAlias, error) {

	importPath := removeQuotes(n.Path.Value)
	var alias string
	if n.Name != nil {
		alias = n.Name.String()
	}
	if alias == "." {
		return ImportAndAlias{}, errors.New("'.' imports are not supported")
	}

	return ImportAndAlias{
		Import: importPath,
		Alias: alias,
	}, nil
}

func isBuiltInType(t string) bool {

	builtInTypes := map[string]struct{}{
		"bool":    {},
		"byte":    {},
		"error":   {},
		"float32": {},
		"float64": {},
		"int":     {},
		"int32":   {},
		"int64":   {},
		"string":  {},
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
	parseOpts parseOptions,
	imports map[string]string,
	fp *os.File,
	fileSet *token.FileSet,
	spec *ast.ValueSpec,
) ([]DeclVar, error) {

	var sType Type
	if spec.Type == nil {
		sType = TypeUnnamedLiteral{}
	} else {
		var err error
		sType, err = getFullType(parseOpts, imports, spec.Type)
		if err != nil {
			return nil, err
		}
	}

	var docString string
	if spec.Doc != nil {
		var err error
		docString, err = readFromFileSet(fp, fileSet, spec.Doc.Pos(), spec.Doc.End())
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
				Name:         spec.Names[iDecl].String(),
				Import:       parseOpts.pkgImportPath,
				Type:         sType,
				LiteralValue: literalValue,
				DocString:    docString,
			},
		)
	}

	return declVars, nil
}

type ParseOption func(parseOptions) parseOptions

type parseOptions struct {
	pkgImportPath string
	dependentTypes bool
}

func ParseWithPkgImportPath(importPath string) ParseOption {
	return func(o parseOptions) parseOptions {
		o.pkgImportPath = importPath
		return o
	}
}

func ParseDependentTypes() ParseOption {
	return func(o parseOptions) parseOptions {
		o.dependentTypes = true
		return o
	}
}
