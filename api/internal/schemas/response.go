package schemas

type OkResponseSchema struct {
	Ok               bool `json:"ok"`
	CustomStatusCode int  `json:"custom_status_code"`
}

type UserTokenResponseSchema struct {
	Token string `json:"token"`
}
