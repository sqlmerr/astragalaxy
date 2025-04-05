package schema

type AuthBody struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

type AuthPayloadToken struct {
	Token string `json:"token"`
}

type AuthPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserTokenResponse struct {
	Token string `json:"token"`
}
