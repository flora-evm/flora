package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// CheckEmergencyPause checks if the module is in emergency pause
// TODO: Implement proper emergency pause functionality
func (k Keeper) CheckEmergencyPause(ctx sdk.Context) error {
	return nil
}

// RequireNotPaused checks if the module is paused and returns error if it is
// TODO: Implement proper emergency pause functionality
func (k Keeper) RequireNotPaused(ctx sdk.Context) error {
	return nil
}

// IsValidatorAllowed checks if a validator is allowed for liquid staking
// TODO: Implement proper validator whitelist/blacklist functionality
func (k Keeper) IsValidatorAllowed(ctx sdk.Context, valAddr sdk.ValAddress) bool {
	return true
}