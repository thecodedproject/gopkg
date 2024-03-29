package gopkg_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/thecodedproject/gopkg"
	"github.com/thecodedproject/gopkg/tmpl"
)

func TestLint(t *testing.T) {

	testCases := []struct {
		Name string
		Pkg  []gopkg.FileContents
	}{
		{
			Name: "empty inputs returns no error",
		},
	}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {
			err := gopkg.Lint(test.Pkg)
			require.NoError(t, err)
		})
	}
}

func TestAddRequiredImports(t *testing.T) {

	testCases := []struct {
		Name     string
		Pkg      []gopkg.FileContents
		Expected []gopkg.FileContents
	}{
		{
			Name: "empty inputs returns no error",
		},
		{
			Name: "pkg with all required imports makes no changes",
			Pkg: []gopkg.FileContents{
				{
					Imports: []gopkg.ImportAndAlias{
						{Import: "a"},
						{Import: "b"},
						{Import: "c"},
						{Import: "d"},
					},
					Consts: []gopkg.DeclVar{
						varTypeNamed("A", "a"),
						varTypeNamed("B", "b"),
						varTypeNamed("C", "c"),
						varTypeNamed("C", "d"),
					},
				},
			},
			Expected: []gopkg.FileContents{
				{
					Imports: []gopkg.ImportAndAlias{
						{Import: "a"},
						{Import: "b"},
						{Import: "c"},
						{Import: "d"},
					},
					Consts: []gopkg.DeclVar{
						varTypeNamed("A", "a"),
						varTypeNamed("B", "b"),
						varTypeNamed("C", "c"),
						varTypeNamed("C", "d"),
					},
				},
			},
		},
		{
			Name: "adds required imports from consts and vars",
			Pkg: []gopkg.FileContents{
				{
					Consts: []gopkg.DeclVar{
						varTypeNamed("A", "a"),
						varTypeNamed("B", "b"),
					},
					Vars: []gopkg.DeclVar{
						varTypeNamed("C", "c"),
						varTypeNamed("B", "b"),
						varTypeNamed("Bb", "b"),
					},
				},
			},
			Expected: []gopkg.FileContents{
				{
					Imports: []gopkg.ImportAndAlias{
						{Import: "a"},
						{Import: "b"},
						{Import: "c"},
					},
					Consts: []gopkg.DeclVar{
						varTypeNamed("A", "a"),
						varTypeNamed("B", "b"),
					},
					Vars: []gopkg.DeclVar{
						varTypeNamed("C", "c"),
						varTypeNamed("B", "b"),
						varTypeNamed("Bb", "b"),
					},
				},
			},
		},
		{
			Name: "adds required imports from types",
			Pkg: []gopkg.FileContents{
				{
					Types: []gopkg.DeclType{
						{
							Type: gopkg.TypeStruct{
								Fields: []gopkg.DeclVar{
									varTypeNamed("A", "a/import/path"),
									varTypeNamed("B", "b/import/path"),
									varTypeNamed("C", "c/import/path"),
									varTypeNamed("D", "d/import/path"),
								},
							},
						},
						{
							Type: gopkg.TypeNamed{
								Import: "other/import/path",
							},
						},
					},
				},
			},
			Expected: []gopkg.FileContents{
				{
					Imports: []gopkg.ImportAndAlias{
						{Import: "a/import/path"},
						{Import: "b/import/path"},
						{Import: "c/import/path"},
						{Import: "d/import/path"},
						{Import: "other/import/path"},
					},
					Types: []gopkg.DeclType{
						{
							Type: gopkg.TypeStruct{
								Fields: []gopkg.DeclVar{
									varTypeNamed("A", "a/import/path"),
									varTypeNamed("B", "b/import/path"),
									varTypeNamed("C", "c/import/path"),
									varTypeNamed("D", "d/import/path"),
								},
							},
						},
						{
							Type: gopkg.TypeNamed{
								Import: "other/import/path",
							},
						},
					},
				},
			},
		},
		{
			Name: "adds required imports from functions",
			Pkg: []gopkg.FileContents{
				{
					Functions: []gopkg.DeclFunc{
						{
							Args: []gopkg.DeclVar{
								varTypeNamed("One", "import/one"),
								varTypeNamed("OtherOne", "import/one"),
								varTypeNamed("Two", "import/two"),
							},
						},
						{
							Args: []gopkg.DeclVar{
								varTypeNamed("OtherOne", "import/one"),
								varTypeNamed("Three", "import/three"),
							},
						},
					},
				},
			},
			Expected: []gopkg.FileContents{
				{
					Imports: []gopkg.ImportAndAlias{
						{Import: "import/one"},
						{Import: "import/three"},
						{Import: "import/two"},
					},
					Functions: []gopkg.DeclFunc{
						{
							Args: []gopkg.DeclVar{
								varTypeNamed("One", "import/one"),
								varTypeNamed("OtherOne", "import/one"),
								varTypeNamed("Two", "import/two"),
							},
						},
						{
							Args: []gopkg.DeclVar{
								varTypeNamed("OtherOne", "import/one"),
								varTypeNamed("Three", "import/three"),
							},
						},
					},
				},
			},
		},
		{
			Name: "does not duplicate, edit of remove any existing imports",
			Pkg: []gopkg.FileContents{
				{
					Imports: []gopkg.ImportAndAlias{
						{Import: "a", Alias: "my_alias_a"},
						{Import: "other/import", Alias: "my_other_alias"},
					},
					Consts: []gopkg.DeclVar{
						varTypeNamed("A", "a"),
						varTypeNamed("B", "b"),
					},
					Vars: []gopkg.DeclVar{
						varTypeNamed("C", "c"),
						varTypeNamed("B", "b"),
						varTypeNamed("Bb", "b"),
					},
				},
			},
			Expected: []gopkg.FileContents{
				{
					Imports: []gopkg.ImportAndAlias{
						{Import: "a", Alias: "my_alias_a"},
						{Import: "b"},
						{Import: "c"},
						{Import: "other/import", Alias: "my_other_alias"},
					},
					Consts: []gopkg.DeclVar{
						varTypeNamed("A", "a"),
						varTypeNamed("B", "b"),
					},
					Vars: []gopkg.DeclVar{
						varTypeNamed("C", "c"),
						varTypeNamed("B", "b"),
						varTypeNamed("Bb", "b"),
					},
				},
			},
		},
		{
			Name: "adds imports across multiple files",
			Pkg: []gopkg.FileContents{
				{
					Imports: []gopkg.ImportAndAlias{
						{Import: "a", Alias: "my_alias_a"},
					},
					Consts: []gopkg.DeclVar{
						varTypeNamed("A", "a"),
						varTypeNamed("B", "b"),
					},
				},
				{
					Imports: []gopkg.ImportAndAlias{
						{Import: "one", Alias: "my_alias_one"},
						{Import: "two", Alias: "my_alias_one"},
					},
					Vars: []gopkg.DeclVar{
						varTypeNamed("One", "one"),
						varTypeNamed("Two", "two"),
					},
				},
				{
					Vars: []gopkg.DeclVar{
						varTypeNamed("One", "one"),
						varTypeNamed("B", "b"),
					},
				},
			},
			Expected: []gopkg.FileContents{
				{
					Imports: []gopkg.ImportAndAlias{
						{Import: "a", Alias: "my_alias_a"},
						{Import: "b"},
					},
					Consts: []gopkg.DeclVar{
						varTypeNamed("A", "a"),
						varTypeNamed("B", "b"),
					},
				},
				{
					Imports: []gopkg.ImportAndAlias{
						{Import: "one", Alias: "my_alias_one"},
						{Import: "two", Alias: "my_alias_one"},
					},
					Vars: []gopkg.DeclVar{
						varTypeNamed("One", "one"),
						varTypeNamed("Two", "two"),
					},
				},
				{
					Imports: []gopkg.ImportAndAlias{
						{Import: "b"},
						{Import: "one"},
					},
					Vars: []gopkg.DeclVar{
						varTypeNamed("One", "one"),
						varTypeNamed("B", "b"),
					},
				},
			},
		},
		{
			Name: "doesn't add import if it matches the files package import path",
			Pkg: []gopkg.FileContents{
				{
					PackageImportPath: "path/to/a",
					Consts: []gopkg.DeclVar{
						varTypeNamed("A", "path/to/a"),
						varTypeNamed("B", "path/to/b"),
					},
				},
				{
					PackageImportPath: "path/to/b",
					Vars: []gopkg.DeclVar{
						varTypeNamed("One", "path/to/one"),
						varTypeNamed("A", "path/to/a"),
						varTypeNamed("B", "path/to/b"),
					},
				},
			},
			Expected: []gopkg.FileContents{
				{
					PackageImportPath: "path/to/a",
					Imports: []gopkg.ImportAndAlias{
						{Import: "path/to/b"},
					},
					Consts: []gopkg.DeclVar{
						varTypeNamed("A", "path/to/a"),
						varTypeNamed("B", "path/to/b"),
					},
				},
				{
					PackageImportPath: "path/to/b",
					Imports: []gopkg.ImportAndAlias{
						{Import: "path/to/a"},
						{Import: "path/to/one"},
					},
					Vars: []gopkg.DeclVar{
						varTypeNamed("One", "path/to/one"),
						varTypeNamed("A", "path/to/a"),
						varTypeNamed("B", "path/to/b"),
					},
				},
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {
			err := gopkg.AddRequiredImports(test.Pkg)
			require.NoError(t, err)

			require.Equal(t, test.Expected, test.Pkg)
		})
	}
}

func TestAddAliasToAllImports(t *testing.T) {

	testCases := []struct {
		Name     string
		Pkg      []gopkg.FileContents
		Expected []gopkg.FileContents
	}{
		{
			Name: "empty inputs returns no error",
		},
		{
			Name: "when all imports have alias does nothing",
			Pkg: []gopkg.FileContents{
				{
					Imports: []gopkg.ImportAndAlias{
						importWithAlias("a/path", "path"),
						importWithAlias("other/path", "path2"),
					},
				},
				{
					Imports: []gopkg.ImportAndAlias{
						importWithAlias("my/import", "import_alias"),
						importWithAlias("other/pkg", "pkg"),
					},
				},
			},
			Expected: []gopkg.FileContents{
				{
					Imports: []gopkg.ImportAndAlias{
						importWithAlias("a/path", "path"),
						importWithAlias("other/path", "path2"),
					},
				},
				{
					Imports: []gopkg.ImportAndAlias{
						importWithAlias("my/import", "import_alias"),
						importWithAlias("other/pkg", "pkg"),
					},
				},
			},
		},
		{
			Name: "when imports have no alias adds using end of path",
			Pkg: []gopkg.FileContents{
				{
					Imports: []gopkg.ImportAndAlias{
						importWithAlias("context", ""),
						importWithAlias("a/path", ""),
					},
				},
				{
					Imports: []gopkg.ImportAndAlias{
						importWithAlias("my/import", ""),
						importWithAlias("other/pkg", ""),
					},
				},
			},
			Expected: []gopkg.FileContents{
				{
					Imports: []gopkg.ImportAndAlias{
						importWithAlias("context", "context"),
						importWithAlias("a/path", "path"),
					},
				},
				{
					Imports: []gopkg.ImportAndAlias{
						importWithAlias("my/import", "import"),
						importWithAlias("other/pkg", "pkg"),
					},
				},
			},
		},
		{
			Name: "does not edit existing aliases",
			Pkg: []gopkg.FileContents{
				{
					Imports: []gopkg.ImportAndAlias{
						importWithAlias("context", "context_alias"),
						importWithAlias("a/path", ""),
					},
				},
				{
					Imports: []gopkg.ImportAndAlias{
						importWithAlias("my/import", ""),
						importWithAlias("other/pkg", "my_pkg_alias"),
					},
				},
			},
			Expected: []gopkg.FileContents{
				{
					Imports: []gopkg.ImportAndAlias{
						importWithAlias("context", "context_alias"),
						importWithAlias("a/path", "path"),
					},
				},
				{
					Imports: []gopkg.ImportAndAlias{
						importWithAlias("my/import", "import"),
						importWithAlias("other/pkg", "my_pkg_alias"),
					},
				},
			},
		},
		{
			Name: "adds an integer to the end of duplicate aliases to make them unique",
			Pkg: []gopkg.FileContents{
				{
					Imports: []gopkg.ImportAndAlias{
						importWithAlias("a/path", ""),
						importWithAlias("other/path", ""),
					},
				},
				{
					Imports: []gopkg.ImportAndAlias{
						importWithAlias("third/path", ""),
						importWithAlias("other/pkg", ""),
					},
				},
			},
			Expected: []gopkg.FileContents{
				{
					Imports: []gopkg.ImportAndAlias{
						importWithAlias("a/path", "path"),
						importWithAlias("other/path", "path2"),
					},
				},
				{
					Imports: []gopkg.ImportAndAlias{
						importWithAlias("third/path", "path"),
						importWithAlias("other/pkg", "pkg"),
					},
				},
			},
		},
		{
			Name: "many duplicate aliases",
			Pkg: []gopkg.FileContents{
				{
					Imports: []gopkg.ImportAndAlias{
						importWithAlias("a/path", ""),
						importWithAlias("b/path", ""),
						importWithAlias("c/path", ""),
						importWithAlias("d/path", ""),
						importWithAlias("e/path", ""),
						importWithAlias("f/path", ""),
						importWithAlias("g/path", ""),
						importWithAlias("h/path", ""),
					},
				},
			},
			Expected: []gopkg.FileContents{
				{
					Imports: []gopkg.ImportAndAlias{
						importWithAlias("a/path", "path"),
						importWithAlias("b/path", "path2"),
						importWithAlias("c/path", "path3"),
						importWithAlias("d/path", "path4"),
						importWithAlias("e/path", "path5"),
						importWithAlias("f/path", "path6"),
						importWithAlias("g/path", "path7"),
						importWithAlias("h/path", "path8"),
					},
				},
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {
			err := gopkg.AddAliasToAllImports(test.Pkg)
			require.NoError(t, err)
			require.Equal(t, test.Expected, test.Pkg)
		})
	}
}

func TestGroupStdImportsFirst(t *testing.T) {

	testCases := []struct {
		Name     string
		Pkg      []gopkg.FileContents
		Expected []gopkg.FileContents
	}{
		{
			Name: "empty input returns no error",
		},
		{
			Name: "single file with only std imports not in group does nothing",
			Pkg: []gopkg.FileContents{
				{
					Imports: tmpl.UnnamedImports(
						"crypto",
						"time",
						"io",
					),
				},
			},
			Expected: []gopkg.FileContents{
				{
					Imports: tmpl.UnnamedImports(
						"crypto",
						"time",
						"io",
					),
				},
			},
		},
		{
			Name: "single file with only std imports in group removes grouping",
			Pkg: []gopkg.FileContents{
				{
					Imports: []gopkg.ImportAndAlias{
						iag("crypto", "", 0),
						iag("time", "", 1),
						iag("io", "", 2),
					},
				},
			},
			Expected: []gopkg.FileContents{
				{
					Imports: []gopkg.ImportAndAlias{
						iag("crypto", "", 0),
						iag("time", "", 0),
						iag("io", "", 0),
					},
				},
			},
		},
		{
			Name: "single file with std imports and other imports not in group",
			Pkg: []gopkg.FileContents{
				{
					Imports: tmpl.UnnamedImports(
						"crypto",
						"my/custom/import1",
						"time",
						"my/custom/import2",
						"io",
					),
				},
			},
			Expected: []gopkg.FileContents{
				{
					Imports: []gopkg.ImportAndAlias{
						iag("crypto", "", -1),
						iag("time", "", -1),
						iag("io", "", -1),
						iag("my/custom/import1", "", 0),
						iag("my/custom/import2", "", 0),
					},
				},
			},
		},
		{
			Name: "single file with std imports and other imports in multiple groups",
			Pkg: []gopkg.FileContents{
				{
					Imports: []gopkg.ImportAndAlias{
						iag("my/custom/import1", "", 2),
						iag("my/custom/import2", "", 3),
						iag("my/custom/import3", "", 4),
						iag("crypto", "", 5),
						iag("time", "", 6),
						iag("io", "", 7),
					},
				},
			},
			Expected: []gopkg.FileContents{
				{
					Imports: []gopkg.ImportAndAlias{
						iag("crypto", "", 0),
						iag("time", "", 0),
						iag("io", "", 0),
						iag("my/custom/import1", "", 2),
						iag("my/custom/import2", "", 3),
						iag("my/custom/import3", "", 4),
					},
				},
			},
		},
		{
			Name: "multiple files with multiple imports",
			Pkg: []gopkg.FileContents{
				{
					Imports: []gopkg.ImportAndAlias{
						iag("my/custom/import3", "", 4),
						iag("crypto", "", 5),
					},
				},
				{
					Imports: []gopkg.ImportAndAlias{
						iag("path/filepath", "", 5),
						iag("my/custom/import1", "", 4),
					},
				},
				{
					Imports: []gopkg.ImportAndAlias{
						iag("my/custom/import2", "", 0),
						iag("testing", "", 0),
						iag("my/custom/import1", "", 0),
					},
				},
			},
			Expected: []gopkg.FileContents{
				{
					Imports: []gopkg.ImportAndAlias{
						iag("crypto", "", 0),
						iag("my/custom/import3", "", 4),
					},
				},
				{
					Imports: []gopkg.ImportAndAlias{
						iag("path/filepath", "", 0),
						iag("my/custom/import1", "", 4),
					},
				},
				{
					Imports: []gopkg.ImportAndAlias{
						iag("testing", "", -1),
						iag("my/custom/import2", "", 0),
						iag("my/custom/import1", "", 0),
					},
				},
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {
			err := gopkg.GroupStdImportsFirst(test.Pkg)
			require.NoError(t, err)
			require.Equal(t, test.Expected, test.Pkg)
		})
	}
}

func TestGroupModuleImportsLast(t *testing.T) {

	testCases := []struct {
		Name       string
		ModulePath string
		Pkg        []gopkg.FileContents
		Expected   []gopkg.FileContents
	}{
		{
			Name: "empty input returns no error",
		},
		{
			Name:       "single file with no module imports does nothing",
			ModulePath: "my/module/path",
			Pkg: []gopkg.FileContents{
				{
					Imports: []gopkg.ImportAndAlias{
						iag("crypto", "", 0),
						iag("time", "", 1),
						iag("io", "", 2),
						iag("otherImport", "alias", 2),
					},
				},
			},
			Expected: []gopkg.FileContents{
				{
					Imports: []gopkg.ImportAndAlias{
						iag("crypto", "", 0),
						iag("time", "", 1),
						iag("io", "", 2),
						iag("otherImport", "alias", 2),
					},
				},
			},
		},
		{
			Name:       "single file with module imports and no groups",
			ModulePath: "my/module/path",
			Pkg: []gopkg.FileContents{
				{
					Imports: []gopkg.ImportAndAlias{
						iag("crypto", "", 0),
						iag("my/module/path", "path", 0),
						iag("time", "", 0),
						iag("my/module/path/subpath1", "", 0),
						iag("io", "", 0),
						iag("my/module/path/subpath2", "", 0),
						iag("my/module/path/subpath2/subsubpath", "", 0),
						iag("otherImport", "alias", 0),
					},
				},
			},
			Expected: []gopkg.FileContents{
				{
					Imports: []gopkg.ImportAndAlias{
						iag("crypto", "", 0),
						iag("time", "", 0),
						iag("io", "", 0),
						iag("otherImport", "alias", 0),
						iag("my/module/path", "path", 1),
						iag("my/module/path/subpath1", "", 1),
						iag("my/module/path/subpath2", "", 1),
						iag("my/module/path/subpath2/subsubpath", "", 1),
					},
				},
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {
			lintRule := gopkg.GroupModuleImportsLast(test.ModulePath)
			err := lintRule(test.Pkg)
			require.NoError(t, err)
			require.Equal(t, test.Expected, test.Pkg)
		})
	}
}

func importWithAlias(
	importName string,
	alias string,
) gopkg.ImportAndAlias {

	return gopkg.ImportAndAlias{
		Import: importName,
		Alias:  alias,
	}
}

func iag(
	importName string,
	alias string,
	group int64,
) gopkg.ImportAndAlias {

	return gopkg.ImportAndAlias{
		Import: importName,
		Alias:  alias,
		Group:  group,
	}
}

func varTypeNamed(
	typeName string,
	importName string,
) gopkg.DeclVar {

	return gopkg.DeclVar{
		Name: "My" + typeName,
		Type: gopkg.TypeNamed{
			Name:   typeName,
			Import: importName,
		},
	}
}
