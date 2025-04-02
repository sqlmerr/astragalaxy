package test

type HTTPTest struct {
	Description   string
	Route         string
	Body          []byte
	ExpectedError bool
	ExpectedCode  int
	BodyValidator func([]byte)
	Method        string
}
