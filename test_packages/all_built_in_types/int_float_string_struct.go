package all_built_in_types

const (
	MyConst, AnotherConst = "some val", "other val"
	RealNumberConst float64 = 1.234
)

var (
	OneInt, TwoInt int = 1, 2
	SomeFloat float32
	SomeUntyped = "a string"
)

type SomeStruct struct{

	IA int
	IB int32
	IC int64

	FA float32
	FB float64

	S string
}

func SomeInts(a int, b int64, c int32) (int, int64, int32) {

	return a, b, c
}

func SomeFloats(a float32, b float64) (float32, float64) {

	return a, b
}

func SomeStrings(a string) string {

	return ""
}
