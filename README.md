

## TODO:

* Make linting method:

  * `Lint([]FileContents, extraLints []func([]FileContents)error) error`

  * Include:
    * Add required imports
    * Sort imports
    * Add import aliases
    * Check for TypeUnnamedLiteral (not in funcs/type defs)
    * Sort FileContents

* Add auto detection of pkg import path (for modules)

* Make DeclFunc.ReturnArgs a []DeclVar - handle named return params

* Make DeclFunc.BodyTmpl execute pass the whole DeclFunc object as the data to execute (not as `.Func`)

  * Add tests for using DeclFunc.BodyData in template exectution
