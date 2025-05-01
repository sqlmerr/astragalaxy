package test

type HTTPTest struct {
	Description   string
	Route         string
	Body          interface{}
	ExpectedError bool
	ExpectedCode  int
	BodyValidator func([]byte)
	BeforeRequest func()
	Method        string
}
