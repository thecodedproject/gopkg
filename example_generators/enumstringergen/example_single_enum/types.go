package enum_stringer

//go:generate go run ../main.go -enum MyEnum
type MyEnum int

const (
	MyEnumUnknown  MyEnum = 0
	MyEnumOne      MyEnum = 1
	MyEnumTwo      MyEnum = 2
	MyEnumThree    MyEnum = 3
	MyEnumSentinal MyEnum = 5
)
