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

type BaseResponse[B any] struct {
	Body   B
	Status int
}

type BaseDataResponse[T any] BaseResponse[DataGenericResponse[T]]

type BaseRequest[B any] struct {
	Body B
}
