package errs

import "errors"

type ErrorCode struct {
	code string
	base error
}

func newCode(code string, base error) ErrorCode {
	return ErrorCode{code: code, base: base}
}

func (c ErrorCode) String() string {
	return c.code
}

func (c ErrorCode) Error() string {
	return c.base.Error()
}

func (c ErrorCode) Unwrap() error {
	return c.base
}

func (c ErrorCode) WithBase(err error) ErrorCode {
	return ErrorCode{
		code: c.code,
		base: err,
	}
}

func (c ErrorCode) Is(target error) bool {
	if t, ok := target.(ErrorCode); ok {
		return c.code == t.code
	}
	return false
}

var (
	CodeInternalServerError          ErrorCode = newCode("INTERNAL_SERVER_ERROR", ErrInternal)
	CodeAnomaly                      ErrorCode = newCode("ANOMALY", ErrInternal) // something unexpected happened in the game
	CodeUnknown                      ErrorCode = newCode("UNKNOWN", ErrInternal)
	CodeDecodeError                  ErrorCode = newCode("DECODE_ERROR", ErrInvalidArgument)
	CodeValidationError              ErrorCode = newCode("VALIDATION_ERROR", ErrInvalidArgument)
	CodeInvalidJWTToken              ErrorCode = newCode("INVALID_JWT_TOKEN", ErrUnauthorized)
	CodeUserUsernameAlreadyOccupied  ErrorCode = newCode("USER_USERNAME_ALREADY_OCCUPIED", ErrConflict)
	CodeUserNotFound                 ErrorCode = newCode("USER_NOT_FOUND", ErrNotFound)
	CodeInvalidCredentials           ErrorCode = newCode("INVALID_CREDENTIALS", ErrUnauthorized)
	CodeAccessDenied                 ErrorCode = newCode("ACCESS_DENIED", ErrAccessDenied)
	CodeAgentNotFound                ErrorCode = newCode("AGENT_NOT_FOUND", ErrNotFound)
	CodeAgentUsernameAlreadyOccupied ErrorCode = newCode("AGENT_USERNAME_ALREADY_OCCUPIED", ErrConflict)
	CodeInvalidAgentToken            ErrorCode = newCode("INVALID_AGENT_TOKEN", ErrUnauthorized)
	CodeAgentLimitExceeded           ErrorCode = newCode("AGENT_LIMIT_EXCEEDED", ErrUnprocessableEntity)
	CodeRadarAreaTooLarge            ErrorCode = newCode("RADAR_AREA_TOO_LARGE", ErrInvalidArgument)
	CodeShipNotFound                 ErrorCode = newCode("SHIP_NOT_FOUND", ErrNotFound)
	CodeInventoryNotFound            ErrorCode = newCode("INVENTORY_NOT_FOUND", ErrNotFound)
	CodeResourceNotFound             ErrorCode = newCode("RESOURCE_NOT_FOUND", ErrNotFound)
	CodeItemNotFound                 ErrorCode = newCode("ITEM_NOT_FOUND", ErrNotFound)
	CodeShipMustBeActive             ErrorCode = newCode("SHIP_MUST_BE_ACTIVE", ErrInternal)
	CodeInvalidTransferDirection     ErrorCode = newCode("INVALID_TRANSFER_DIRECTION", ErrInvalidArgument)
	CodeNotEnoughResources           ErrorCode = newCode("NOT_ENOUGH_RESOURCES", ErrUnprocessableEntity)
	CodeInventoryIsFull              ErrorCode = newCode("INVENTORY_IS_FULL", ErrUnprocessableEntity)
	CodeItemNotInInventory           ErrorCode = newCode("ITEM_NOT_IN_INVENTORY", ErrUnprocessableEntity)
	CodeAgentInCooldown              ErrorCode = newCode("AGENT_IN_COOLDOWN", ErrUnprocessableEntity)
	CodeInvalidWarpPath              ErrorCode = newCode("INVALID_WARP_PATH", ErrInvalidArgument)
	CodeAlreadyAtDestination         ErrorCode = newCode("ALREADY_AT_DESTINATION", ErrNotModified)
	CodeInvalidCoordinates           ErrorCode = newCode("INVALID_COORDINATES", ErrInvalidArgument)
	CodeInvalidShipState             ErrorCode = newCode("INVALID_SHIP_STATE", ErrUnprocessableEntity)
	CodeShipAlreadyInThisState       ErrorCode = newCode("SHIP_ALREADY_IN_THIS_STATE", ErrNotModified)
	CodeCannotDock                   ErrorCode = newCode("CANNOT_DOCK_HERE", ErrUnprocessableEntity)
	CodeInvalidUUID                  ErrorCode = newCode("INVALID_UUID", ErrInvalidArgument)
	CodeInvalidLocation              ErrorCode = newCode("INVALID_LOCATION", ErrUnprocessableEntity) // can be changed
	CodeWaypointNotFound             ErrorCode = newCode("WAYPOINT_NOT_FOUND", ErrNotFound)
	CodeResourceDepositNotFound      ErrorCode = newCode("RESOURCE_DEPOSIT_NOT_FOUND", ErrNotFound)
	CodeRecipeNotFound               ErrorCode = newCode("RECIPE_NOT_FOUND", ErrNotFound)
	CodeFacilityNotFound             ErrorCode = newCode("FACILITY_NOT_FOUND", ErrNotFound)
	CodeShipModuleAlreadyInstalled   ErrorCode = newCode("SHIP_MODULE_ALREADY_INSTALLED", ErrConflict)
	CodeCannotUseItem                ErrorCode = newCode("CANNOT_USE_ITEM", ErrUnprocessableEntity)
	CodeInvalidUsername              ErrorCode = newCode("INVALID_USERNAME", ErrInvalidArgument)
	CodeInvalidPassword              ErrorCode = newCode("INVALID_PASSWORD", ErrInvalidArgument)
	CodeInvalidFacilityType          ErrorCode = newCode("INVALID_FACILITY_TYPE", ErrInvalidArgument)
	CodeInvalidShipModuleType        ErrorCode = newCode("INVALID_SHIP_MODULE_TYPE", ErrInvalidArgument)
	CodeInvalidShipClass             ErrorCode = newCode("INVALID_SHIP_CLASS", ErrInvalidArgument)
	CodeMiningLimitExceeded          ErrorCode = newCode("MINING_LIMIT_EXCEEDED", ErrUnprocessableEntity)
	CodePlanetDiscoveryNotFound      ErrorCode = newCode("PLANET_DISCOVERY_NOT_FOUND", ErrNotFound)
	CodeSystemDiscoveryNotFound      ErrorCode = newCode("SYSTEM_DISCOVERY_NOT_FOUND", ErrNotFound)
	CodeWalletNotFound               ErrorCode = newCode("WALLET_NOT_FOUND", ErrNotFound)
	CodeInvalidWalletAddress         ErrorCode = newCode("INVALID_WALLET_ADDRESS", ErrInvalidArgument)
	CodeInsufficientBalance          ErrorCode = newCode("INSUFFICIENT_BALANCE", ErrUnprocessableEntity)
	CodeInvalidPaymentAmount         ErrorCode = newCode("INVALID_PAYMENT_AMOUNT", ErrInvalidArgument)
	CodeWalletLimitExceeded          ErrorCode = newCode("WALLET_LIMIT_EXCEEDED", ErrUnprocessableEntity)
	CodeCannotRegisterWallet         ErrorCode = newCode("CANNOT_REGISTER_WALLET", ErrUnprocessableEntity)
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
	return errors.Is(err, code)
}

func GetCode(err error) ErrorCode {
	if codeErr, ok := errors.AsType[ErrorCode](err); ok {
		return codeErr
	}

	if oldCodeErr, ok := errors.AsType[WithCode](err); ok {
		return oldCodeErr.Code
	}
	return CodeInternalServerError
}
