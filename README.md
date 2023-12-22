gopkg
=====

A toolkit for writing code generators in Golang.

## Overview

`gopkg` provides a simple, AST-like structure for representing code declarations - `FileContents`.

It provides two main methods for interacting with this structure:

* `Parse` takes a path to a package and returns the contents of each `.go` file as a slice of `FileContents` objects.

* `Generate` takes a slice of `FileContents` objects and writes each object to `.go` file.

E.g.
```go
  // pckContents has type `[]gopkg.FileContents`
  pkgContents, err := gopkg.Parse("./path/to/my/package")
  // check err

  err = gopkg.Generate(pkgContents)
  // check err
```

### [`FileContents`](https://github.com/thecodedproject/gopkg/blob/main/types.go#L7) structure

Contains all package level declaration contained within a `.go` file, as well as information about the path to the file and it's Golang import path.

```go
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
```

## Examples:

The [`example_generators`](https://github.com/thecodedproject/gopkg/tree/main/example_generators) folder contains several toy examples of generator implementations.

See also:

[`servicegen`](https://github.com/thecodedproject/servicegen`) - A generator for creating service interfaces and tests.
[`resourcegen`](https://github.com/thecodedproject/servicegen) - For generating interfaces required for dependency injection.


## TODO:

* Make linting method to sort FileContents


* Add DeclFunc.AdditionalImports field, which can be used to add imports at the function level

* Add DeclFunc.DocString field

* Add error checks for generation, e.g:

  * Return error when generating a func if:
    * Func name not set
    * Func args unnamed
      * Note; TypeFunc args may be unnamed... should also allow unnamed types on DeclFunc in interfaces
    * Func arg or ret arg is missing a type
  * Return error when generation a file if:
    * Filepath not set
    * PackageName not set
  * Error for TypeDecl when:
    * Name not set
    * Type not set
  * etc...

* Add linter to _sanatize literals_ which will strip leading and trailing whitespace from all literals:
  * Remove leading whitespace from DeclFunc.BodyTmpl
  * Remove leading and trailing whitespace from all strings in `FileContents`
  * Maybe remove newlines from things that shouldnt have new lines? (e.g. Decl names?)

* Consider removing `Import` field from declaration types - it doesn't seem that this is used at all for generating; It seems like a convenience field but I'm not sure there is a scenario where this is useful

