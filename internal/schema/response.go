package schema

type OkResponse struct {
	Ok               bool `json:"ok"`
	CustomStatusCode int  `json:"custom_status_code"`
}

type DataResponse struct {
	Data any `json:"data"`
}

type DataGenericResponse[T any] struct {
	Data T `json:"data"`
}
