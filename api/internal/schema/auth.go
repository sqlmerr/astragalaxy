package schema

type AuthBody struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

type AuthPayload struct {
	Token string `json:"token"`
}
