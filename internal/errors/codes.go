package core_errors

import "errors"

type ErrorCode string

const (
	CodeInternalServerError          ErrorCode = "INTERNAL_SERVER_ERROR"
	CodeAnomaly                      ErrorCode = "ANOMALY" // something unexpected happened in game
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
	CodeRadarAreaTooLarge            ErrorCode = "RADAR_AREA_TOO_LARGE"
	CodeShipNotFound                 ErrorCode = "SHIP_NOT_FOUND"
	CodeInventoryNotFound            ErrorCode = "INVENTORY_NOT_FOUND"
	CodeResourceNotFound             ErrorCode = "RESOURCE_NOT_FOUND"
	CodeItemNotFound                 ErrorCode = "ITEM_NOT_FOUND"
	CodeShipMustBeActive             ErrorCode = "SHIP_MUST_BE_ACTIVE"
	CodeInvalidTransferDirection     ErrorCode = "INVALID_TRANSFER_DIRECTION"
	CodeNotEnoughResources           ErrorCode = "NOT_ENOUGH_RESOURCES"
	CodeInventoryIsFull              ErrorCode = "INVENTORY_IS_FULL"
	CodeItemNotInInventory           ErrorCode = "ITEM_NOT_IN_INVENTORY"
	CodeCharacterInCooldown          ErrorCode = "CHARACTER_IN_COOLDOWN"
	CodeInvalidWarpPath              ErrorCode = "INVALID_WARP_PATH"
	CodeAlreadyAtDestination         ErrorCode = "ALREADY_AT_DESTINATION"
	CodeInvalidCoordinates           ErrorCode = "INVALID_COORDINATES"
	CodeInvalidShipState             ErrorCode = "INVALID_SHIP_STATE"
	CodeShipAlreadyInThisState       ErrorCode = "SHIP_ALREADY_IN_THIS_STATE"
	CodeCannotDock                   ErrorCode = "CANNOT_DOCK_HERE"
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

func IsCode(err error, code ErrorCode) bool {
	if codeErr, ok := errors.AsType[WithCode](err); ok {
		return codeErr.Code == code
	}

	return false
}
