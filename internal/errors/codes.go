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
	CodeInternalServerError = newCode("INTERNAL_SERVER_ERROR", ErrInternal)
	CodeAnomaly             = newCode(
		"ANOMALY",
		ErrInternal,
	) // something unexpected happened in the game
	CodeUnknown                      = newCode("UNKNOWN", ErrInternal)
	CodeDecodeError                  = newCode("DECODE_ERROR", ErrInvalidArgument)
	CodeValidationError              = newCode("VALIDATION_ERROR", ErrInvalidArgument)
	CodeInvalidJWTToken              = newCode("INVALID_JWT_TOKEN", ErrUnauthorized)
	CodeUserUsernameAlreadyOccupied  = newCode("USER_USERNAME_ALREADY_OCCUPIED", ErrConflict)
	CodeUserNotFound                 = newCode("USER_NOT_FOUND", ErrNotFound)
	CodeInvalidCredentials           = newCode("INVALID_CREDENTIALS", ErrUnauthorized)
	CodeAccessDenied                 = newCode("ACCESS_DENIED", ErrAccessDenied)
	CodeAgentNotFound                = newCode("AGENT_NOT_FOUND", ErrNotFound)
	CodeAgentUsernameAlreadyOccupied = newCode("AGENT_USERNAME_ALREADY_OCCUPIED", ErrConflict)
	CodeInvalidAgentToken            = newCode("INVALID_AGENT_TOKEN", ErrUnauthorized)
	CodeAgentLimitExceeded           = newCode("AGENT_LIMIT_EXCEEDED", ErrUnprocessableEntity)
	CodeRadarAreaTooLarge            = newCode("RADAR_AREA_TOO_LARGE", ErrInvalidArgument)
	CodeShipNotFound                 = newCode("SHIP_NOT_FOUND", ErrNotFound)
	CodeInventoryNotFound            = newCode("INVENTORY_NOT_FOUND", ErrNotFound)
	CodeResourceNotFound             = newCode("RESOURCE_NOT_FOUND", ErrNotFound)
	CodeItemNotFound                 = newCode("ITEM_NOT_FOUND", ErrNotFound)
	CodeShipMustBeActive             = newCode("SHIP_MUST_BE_ACTIVE", ErrInternal)
	CodeInvalidTransferDirection     = newCode("INVALID_TRANSFER_DIRECTION", ErrInvalidArgument)
	CodeNotEnoughResources           = newCode("NOT_ENOUGH_RESOURCES", ErrUnprocessableEntity)
	CodeInventoryIsFull              = newCode("INVENTORY_IS_FULL", ErrUnprocessableEntity)
	CodeItemNotInInventory           = newCode("ITEM_NOT_IN_INVENTORY", ErrUnprocessableEntity)
	CodeAgentInCooldown              = newCode("AGENT_IN_COOLDOWN", ErrUnprocessableEntity)
	CodeInvalidWarpPath              = newCode("INVALID_WARP_PATH", ErrInvalidArgument)
	CodeAlreadyAtDestination         = newCode("ALREADY_AT_DESTINATION", ErrNotModified)
	CodeInvalidCoordinates           = newCode("INVALID_COORDINATES", ErrInvalidArgument)
	CodeInvalidShipState             = newCode("INVALID_SHIP_STATE", ErrUnprocessableEntity)
	CodeShipAlreadyInThisState       = newCode("SHIP_ALREADY_IN_THIS_STATE", ErrNotModified)
	CodeCannotDock                   = newCode("CANNOT_DOCK_HERE", ErrUnprocessableEntity)
	CodeInvalidUUID                  = newCode("INVALID_UUID", ErrInvalidArgument)
	CodeInvalidLocation              = newCode("INVALID_LOCATION", ErrUnprocessableEntity) // can be changed
	CodeWaypointNotFound             = newCode("WAYPOINT_NOT_FOUND", ErrNotFound)
	CodeResourceDepositNotFound      = newCode("RESOURCE_DEPOSIT_NOT_FOUND", ErrNotFound)
	CodeRecipeNotFound               = newCode("RECIPE_NOT_FOUND", ErrNotFound)
	CodeFacilityNotFound             = newCode("FACILITY_NOT_FOUND", ErrNotFound)
	CodeShipModuleAlreadyInstalled   = newCode("SHIP_MODULE_ALREADY_INSTALLED", ErrConflict)
	CodeCannotUseItem                = newCode("CANNOT_USE_ITEM", ErrUnprocessableEntity)
	CodeInvalidUsername              = newCode("INVALID_USERNAME", ErrInvalidArgument)
	CodeInvalidPassword              = newCode("INVALID_PASSWORD", ErrInvalidArgument)
	CodeInvalidFacilityType          = newCode("INVALID_FACILITY_TYPE", ErrInvalidArgument)
	CodeInvalidShipModuleType        = newCode("INVALID_SHIP_MODULE_TYPE", ErrInvalidArgument)
	CodeInvalidShipClass             = newCode("INVALID_SHIP_CLASS", ErrInvalidArgument)
	CodeMiningLimitExceeded          = newCode("MINING_LIMIT_EXCEEDED", ErrUnprocessableEntity)
	CodePlanetDiscoveryNotFound      = newCode("PLANET_DISCOVERY_NOT_FOUND", ErrNotFound)
	CodeSystemDiscoveryNotFound      = newCode("SYSTEM_DISCOVERY_NOT_FOUND", ErrNotFound)
	CodeWalletNotFound               = newCode("WALLET_NOT_FOUND", ErrNotFound)
	CodeInvalidWalletAddress         = newCode("INVALID_WALLET_ADDRESS", ErrInvalidArgument)
	CodeInsufficientBalance          = newCode("INSUFFICIENT_BALANCE", ErrUnprocessableEntity)
	CodeInvalidPaymentAmount         = newCode("INVALID_PAYMENT_AMOUNT", ErrInvalidArgument)
	CodeWalletLimitExceeded          = newCode("WALLET_LIMIT_EXCEEDED", ErrUnprocessableEntity)
	CodeCannotRegisterWallet         = newCode("CANNOT_REGISTER_WALLET", ErrUnprocessableEntity)
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
