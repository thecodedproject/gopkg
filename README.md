

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

