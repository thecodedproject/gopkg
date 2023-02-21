package named_return_args

type SomeInterface interface {
	SomeFunc() (a int, b error)
	OtherFunc() (c, d int64)
}

func MyMethod() (e, f, g int32) {

	return 0, 0, 0
}

func MyOtherMethod() (i int32, j float64, k error) {

	return 0, 0, nil
}
