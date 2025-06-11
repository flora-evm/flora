package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	
	"github.com/rollchains/flora/x/liquidstaking/types"
)

// CanTokenizeShares checks if the requested tokenization respects both global and validator caps
func (k Keeper) CanTokenizeShares(ctx sdk.Context, validatorAddr string, tokensToAdd math.Int) error {
	params := k.GetParams(ctx)
	
	// Check if module is enabled
	if !params.Enabled {
		return types.ErrModuleDisabled
	}
	
	// Get total bonded tokens
	totalBonded := k.stakingKeeper.TotalBondedTokens(ctx)
	if totalBonded.IsZero() {
		return nil // No cap check needed if nothing is bonded
	}
	
	// Check global cap
	totalLiquidStaked := k.GetTotalLiquidStaked(ctx)
	newTotalLiquidStaked := totalLiquidStaked.Add(tokensToAdd)
	
	// Calculate global liquid staked ratio
	globalRatio := math.LegacyNewDecFromInt(newTotalLiquidStaked).Quo(math.LegacyNewDecFromInt(totalBonded))
	if globalRatio.GT(params.GlobalLiquidStakingCap) {
		return types.ErrExceedsGlobalCap
	}
	
	// Get validator
	valAddr, err := sdk.ValAddressFromBech32(validatorAddr)
	if err != nil {
		return err
	}
	
	validator, err := k.stakingKeeper.GetValidator(ctx, valAddr)
	if err != nil {
		return err
	}
	
	// Check validator cap
	validatorTokens := validator.GetTokens()
	if validatorTokens.IsZero() {
		return nil // No cap check needed if validator has no tokens
	}
	
	validatorLiquidStaked := k.GetValidatorLiquidStaked(ctx, validatorAddr)
	newValidatorLiquidStaked := validatorLiquidStaked.Add(tokensToAdd)
	
	// Calculate validator liquid staked ratio
	validatorRatio := math.LegacyNewDecFromInt(newValidatorLiquidStaked).Quo(math.LegacyNewDecFromInt(validatorTokens))
	if validatorRatio.GT(params.ValidatorLiquidCap) {
		return types.ErrExceedsValidatorCap
	}
	
	return nil
}