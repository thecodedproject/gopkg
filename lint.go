package gopkg

import (
	"sort"
	"strings"
	"strconv"
)

func Lint(
	pkg []FileContents,
	extraLintRules ...func([]FileContents)error,
) error {

	defaultRules := []func([]FileContents)error{
		AddRequiredImports,
		AddAliasToAllImports,
	}

	defaultRules = append(defaultRules, extraLintRules...)

	return LintCustom(pkg, defaultRules...)
}

func LintCustom(
	pkg []FileContents,
	lintRules ...func([]FileContents)error,
) error {

	for _, rule := range lintRules {
		err := rule(pkg)
		if err != nil {
			return err
		}
	}
	return nil
}

func AddAliasToAllImports(pkg []FileContents) error {

	for iF := range pkg {

		existingAliases := make(map[string]bool)

		for iI := range pkg[iF].Imports {
			if pkg[iF].Imports[iI].Alias == "" {
				pathElems := strings.Split(pkg[iF].Imports[iI].Import, "/")

				alias := pathElems[len(pathElems)-1]

				if existingAliases[alias] {
					iAlias := 2
					for {
						potentialAlias := alias + strconv.Itoa(iAlias)
						if existingAliases[potentialAlias] {
							iAlias++
						} else {
							alias = potentialAlias
							break
						}
					}
				}

				existingAliases[alias] = true

				pkg[iF].Imports[iI].Alias = alias
			}
		}
	}
	return nil
}

func AddRequiredImports(pkg []FileContents) error {

	for iF := range pkg {
		existingImportSet := make(map[string]bool)
		for _, i := range pkg[iF].Imports {
			existingImportSet[i.Import] = true
		}

		importsToAdd := complement(
			getFileRequiredTypeImports(pkg[iF]),
			existingImportSet,
		)

		for i := range importsToAdd {
			pkg[iF].Imports = append(
				pkg[iF].Imports,
				ImportAndAlias{
					Import: i,
				},
			)
		}

		sort.Slice(pkg[iF].Imports, func(i, j int) bool {
			return pkg[iF].Imports[i].Import < pkg[iF].Imports[j].Import
		})
	}

	return nil
}

func getFileRequiredTypeImports(f FileContents) map[string]bool {

	requiredTypeImports := make(map[string]bool)
	for _, c := range f.Consts {
		requiredTypeImports = union(
			requiredTypeImports,
			c.RequiredImports(),
		)
	}
	for _, v := range f.Vars {
		requiredTypeImports = union(
			requiredTypeImports,
			v.RequiredImports(),
		)
	}
	for _, t := range f.Types {
		requiredTypeImports = union(
			requiredTypeImports,
			t.Type.RequiredImports(),
		)
	}
	for _, f := range f.Functions {
		requiredTypeImports = union(
			requiredTypeImports,
			f.RequiredImports(),
		)
	}
	return requiredTypeImports
}
