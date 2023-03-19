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
  pckContents, err := gopkg.Parse("./path/to/my/package")
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

The [`example_genertors`](https://github.com/thecodedproject/gopkg/tree/main/example_generators) folder contains several toy examples of generator implementations.

See also:

[`servicegen`](https://github.com/thecodedproject/servicegen`) - A generator for creating service interfaces and tests.
[`resourcegen`](https://github.com/thecodedproject/servicegen) - For generating interfaces required for dependency injection.


## TODO:

* Make linting method:

  * Include:
    * Change Add required imports to not include the import for the current file import path
    * Sort imports
    * Check for TypeUnnamedLiteral (not in funcs/type defs)
    * Sort FileContents

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

* Add doc strings

