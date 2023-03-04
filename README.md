

## TODO:

* Make linting method:

  * Include:
    * Change Add required imports to not include the import for the current file import path
    * Sort imports
    * Check for TypeUnnamedLiteral (not in funcs/type defs)
    * Sort FileContents

* Make DeclFunc.BodyTmpl execute pass the whole DeclFunc object as the data to execute (not as `.Func`)

  * Add tests for using DeclFunc.BodyData in template exectution

* Add `TypeFunc` for parsing + generating func types

* Allow generating named return args

* Add doc strings

