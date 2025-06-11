package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// UpdateLiquidStakedAmounts updates the total and per-validator liquid staked amounts
func (k Keeper) UpdateLiquidStakedAmounts(ctx sdk.Context, validatorAddr string, amount math.Int, isIncrease bool) {
	// Update total liquid staked amount
	totalLiquidStaked := k.GetTotalLiquidStaked(ctx)
	if isIncrease {
		totalLiquidStaked = totalLiquidStaked.Add(amount)
	} else {
		totalLiquidStaked = totalLiquidStaked.Sub(amount)
		if totalLiquidStaked.IsNegative() {
			totalLiquidStaked = math.ZeroInt()
		}
	}
	k.SetTotalLiquidStaked(ctx, totalLiquidStaked)
	
	// Update validator liquid staked amount
	validatorLiquidStaked := k.GetValidatorLiquidStaked(ctx, validatorAddr)
	if isIncrease {
		validatorLiquidStaked = validatorLiquidStaked.Add(amount)
	} else {
		validatorLiquidStaked = validatorLiquidStaked.Sub(amount)
		if validatorLiquidStaked.IsNegative() {
			validatorLiquidStaked = math.ZeroInt()
		}
	}
	k.SetValidatorLiquidStaked(ctx, validatorAddr, validatorLiquidStaked)
}