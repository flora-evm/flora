package types

import (
	errorsmod "cosmossdk.io/errors"
)

// x/liquidstaking module sentinel errors
var (
	// ErrModuleDisabled is returned when attempting operations while the module is disabled
	ErrModuleDisabled          = errorsmod.Register(ModuleName, 2, "liquid staking module is disabled in params")
	
	// ErrInvalidTokenizationRecord is returned when a tokenization record fails validation
	ErrInvalidTokenizationRecord = errorsmod.Register(ModuleName, 3, "invalid tokenization record")
	
	// ErrTokenizationRecordNotFound is returned when a requested tokenization record doesn't exist
	ErrTokenizationRecordNotFound = errorsmod.Register(ModuleName, 4, "tokenization record not found")
	
	// ErrInvalidShares is returned when shares amount is invalid (zero, negative, or malformed)
	ErrInvalidShares           = errorsmod.Register(ModuleName, 5, "invalid shares amount")
	
	// ErrInvalidValidator is returned when validator is not eligible for liquid staking (e.g., jailed)
	ErrInvalidValidator        = errorsmod.Register(ModuleName, 6, "validator not eligible for liquid staking")
	
	// ErrInvalidDelegator is returned when delegator address is invalid
	ErrInvalidDelegator        = errorsmod.Register(ModuleName, 7, "invalid delegator address")
	
	// ErrExceedsGlobalCap is returned when tokenization would exceed the global liquid staking cap
	ErrExceedsGlobalCap        = errorsmod.Register(ModuleName, 8, "tokenization would exceed global liquid staking cap")
	
	// ErrExceedsValidatorCap is returned when tokenization would exceed the validator's liquid staking cap
	ErrExceedsValidatorCap     = errorsmod.Register(ModuleName, 9, "tokenization would exceed validator liquid staking cap")
	
	// ErrTokenizationRecordAlreadyExists is returned when attempting to create a duplicate record
	ErrTokenizationRecordAlreadyExists = errorsmod.Register(ModuleName, 10, "tokenization record already exists")
	
	// ErrDuplicateLiquidStakingToken is returned when a liquid staking token denom already exists
	ErrDuplicateLiquidStakingToken = errorsmod.Register(ModuleName, 11, "liquid staking token denom already exists")
	
	// ErrInsufficientShares is returned when delegator has insufficient shares to tokenize
	ErrInsufficientShares      = errorsmod.Register(ModuleName, 12, "insufficient delegation shares for tokenization")
	
	// ErrDelegationNotFound is returned when delegation doesn't exist for the delegator/validator pair
	ErrDelegationNotFound      = errorsmod.Register(ModuleName, 13, "delegation not found for given delegator and validator")
)