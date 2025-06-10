package types

import (
	errorsmod "cosmossdk.io/errors"
)

// x/liquidstaking module sentinel errors
var (
	ErrModuleDisabled          = errorsmod.Register(ModuleName, 2, "liquid staking module is disabled")
	ErrInvalidTokenizationRecord = errorsmod.Register(ModuleName, 3, "invalid tokenization record")
	ErrTokenizationRecordNotFound = errorsmod.Register(ModuleName, 4, "tokenization record not found")
	ErrInvalidShares           = errorsmod.Register(ModuleName, 5, "invalid shares amount")
	ErrInvalidValidator        = errorsmod.Register(ModuleName, 6, "invalid validator")
	ErrInvalidDelegator        = errorsmod.Register(ModuleName, 7, "invalid delegator")
	ErrExceedsGlobalCap        = errorsmod.Register(ModuleName, 8, "exceeds global liquid staking cap")
	ErrExceedsValidatorCap     = errorsmod.Register(ModuleName, 9, "exceeds validator liquid staking cap")
)