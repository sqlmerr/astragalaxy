package test

type HTTPTest struct {
	Description      string
	Route            string
	ExpectedError    bool
	ExpectedCode     int
	ExpectedBodyKeys map[string]interface{}
	Method           string
}
