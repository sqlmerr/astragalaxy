package schema

type OkResponseSchema struct {
	Ok               bool `json:"ok"`
	CustomStatusCode int  `json:"custom_status_code"`
}

type DataResponseSchema struct {
	Data any `json:"data"`
}

type DataGenericResponse[T any] struct {
	Data T `json:"data"`
}
