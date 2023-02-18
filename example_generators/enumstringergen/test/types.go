package enum_stringer

//go:generate go run ../main.go -enum MyEnum
type MyEnum int

const (
	MyEnumOne MyEnum = 1
	MyEnumTwo MyEnum = 2
	MyEnumThree MyEnum = 3
	MyEnumSentinal MyEnum = 4
)
