package test

type HTTPTest struct {
	Description   string
	Route         string
	ExpectedError bool
	ExpectedCode  int
	BodyValidator func([]byte)
	Method        string
}
