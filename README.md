

## TODO:

* Make linting method:

  * `Lint([]FileContents, extraLints []func([]FileContents)error) error`

  * Include:
    * Change Add required imports to not include the import for the current file import path
    * Sort imports
    * Check for TypeUnnamedLiteral (not in funcs/type defs)
    * Sort FileContents

* Add auto detection of pkg import path (for modules)

* Make DeclFunc.ReturnArgs a []DeclVar - handle named return params

* Make DeclFunc.BodyTmpl execute pass the whole DeclFunc object as the data to execute (not as `.Func`)

  * Add tests for using DeclFunc.BodyData in template exectution

