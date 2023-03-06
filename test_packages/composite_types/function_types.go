package composite_types

var SomeVar func(int)string

type SomeType func()error

type SomeStruct struct {
	UnnamedFunc func(int, float32) (string, error)
	NamedFunc func(a int64, b bool) (c float64, d string)
}

func SomeFunc(f func()) func() (int, int) {

	return nil
}
