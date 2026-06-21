package core_errors

type ErrorCode string

const (
	CodeInternalServerError          ErrorCode = "INTERNAL_SERVER_ERROR"
	CodeUnknown                      ErrorCode = "UNKNOWN"
	CodeDecodeError                  ErrorCode = "DECODE_ERROR"
	CodeValidationError              ErrorCode = "VALIDATION_ERROR"
	CodeInvalidJWTToken              ErrorCode = "INVALID_JWT_TOKEN"
	CodeUserUsernameAlreadyOccupied  ErrorCode = "USER_USERNAME_ALREADY_OCCUPIED"
	CodeUserNotFound                 ErrorCode = "USER_NOT_FOUND"
	CodeInvalidCredentials           ErrorCode = "INVALID_CREDENTIALS"
	CodeAccessDenied                 ErrorCode = "ACCESS_DENIED"
	CodeAgentNotFound                ErrorCode = "AGENT_NOT_FOUND"
	CodeAgentUsernameAlreadyOccupied ErrorCode = "AGENT_USERNAME_ALREADY_OCCUPIED"
	CodeInvalidAgentToken            ErrorCode = "INVALID_AGENT_TOKEN"
	CodeAgentLimitExceeded           ErrorCode = "AGENT_LIMIT_EXCEEDED"
)

type WithCode struct {
	Err  error
	Code ErrorCode
}

func (e WithCode) Error() string {
	return e.Err.Error()
}

func (e WithCode) Unwrap() error {
	return e.Err
}

func NewWithCode(code ErrorCode, err error) WithCode {
	return WithCode{
		Err:  err,
		Code: code,
	}
}
