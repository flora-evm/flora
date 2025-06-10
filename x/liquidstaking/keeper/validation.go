package keeper

import (
	"fmt"
	
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	
	"github.com/rollchains/flora/x/liquidstaking/types"
)

// ValidateGlobalLiquidStakingCap checks if tokenizing additional shares would exceed the global cap
func (k Keeper) ValidateGlobalLiquidStakingCap(ctx sdk.Context, additionalShares math.Int) error {
	params := k.GetParams(ctx)
	if !params.Enabled {
		return types.ErrModuleDisabled
	}
	
	// For Stage 2, we'll use a placeholder for total bonded tokens
	// In Stage 3+, this will integrate with the staking keeper
	// TODO: Replace with k.stakingKeeper.TotalBondedTokens(ctx) when integrated
	totalBonded := math.NewInt(1_000_000_000) // 1 billion tokens placeholder
	
	totalLiquid := k.GetTotalLiquidStaked(ctx).Add(additionalShares)
	maxAllowed := params.GlobalLiquidStakingCap.MulInt(totalBonded).TruncateInt()
	
	if totalLiquid.GT(maxAllowed) {
		return fmt.Errorf("would exceed global liquid staking cap of %s: current %s + additional %s > max %s",
			params.GlobalLiquidStakingCap.String(),
			k.GetTotalLiquidStaked(ctx).String(),
			additionalShares.String(),
			maxAllowed.String())
	}
	
	return nil
}

// ValidateValidatorLiquidCap checks if tokenizing additional shares would exceed a validator's cap
func (k Keeper) ValidateValidatorLiquidCap(ctx sdk.Context, validatorAddr string, additionalShares math.Int) error {
	params := k.GetParams(ctx)
	if !params.Enabled {
		return types.ErrModuleDisabled
	}
	
	// For Stage 2, we'll use a placeholder for validator's total shares
	// In Stage 3+, this will integrate with the staking keeper
	// TODO: Replace with k.stakingKeeper.GetValidator(ctx, valAddr).GetTokens() when integrated
	validatorTotalShares := math.NewInt(100_000_000) // 100 million tokens placeholder
	
	validatorLiquid := k.GetValidatorLiquidStaked(ctx, validatorAddr).Add(additionalShares)
	maxAllowed := params.ValidatorLiquidCap.MulInt(validatorTotalShares).TruncateInt()
	
	if validatorLiquid.GT(maxAllowed) {
		return fmt.Errorf("would exceed validator liquid cap of %s: current %s + additional %s > max %s",
			params.ValidatorLiquidCap.String(),
			k.GetValidatorLiquidStaked(ctx, validatorAddr).String(),
			additionalShares.String(),
			maxAllowed.String())
	}
	
	return nil
}

// CanTokenizeShares performs all validation checks for tokenizing shares
func (k Keeper) CanTokenizeShares(ctx sdk.Context, validatorAddr string, shares math.Int) error {
	// Check if module is enabled
	params := k.GetParams(ctx)
	if !params.Enabled {
		return types.ErrModuleDisabled
	}
	
	// Validate shares amount
	if !shares.IsPositive() {
		return fmt.Errorf("shares must be positive: %s", shares.String())
	}
	
	// Check global cap
	if err := k.ValidateGlobalLiquidStakingCap(ctx, shares); err != nil {
		return err
	}
	
	// Check validator cap
	if err := k.ValidateValidatorLiquidCap(ctx, validatorAddr, shares); err != nil {
		return err
	}
	
	return nil
}

// UpdateLiquidStakedAmounts updates the total and validator liquid staked amounts
// This should be called after successful tokenization or redemption
func (k Keeper) UpdateLiquidStakedAmounts(ctx sdk.Context, validatorAddr string, delta math.Int, increase bool) {
	// Update total liquid staked
	total := k.GetTotalLiquidStaked(ctx)
	if increase {
		total = total.Add(delta)
	} else {
		total = total.Sub(delta)
		if total.IsNegative() {
			total = math.ZeroInt()
		}
	}
	k.SetTotalLiquidStaked(ctx, total)
	
	// Update validator liquid staked
	validatorAmount := k.GetValidatorLiquidStaked(ctx, validatorAddr)
	if increase {
		validatorAmount = validatorAmount.Add(delta)
	} else {
		validatorAmount = validatorAmount.Sub(delta)
		if validatorAmount.IsNegative() {
			validatorAmount = math.ZeroInt()
		}
	}
	k.SetValidatorLiquidStaked(ctx, validatorAddr, validatorAmount)
}

// ValidateTokenizationRecord validates a tokenization record before storage
func (k Keeper) ValidateTokenizationRecord(ctx sdk.Context, record types.TokenizationRecord) error {
	// Basic validation
	if err := record.Validate(); err != nil {
		return err
	}
	
	// Check if ID already exists
	if _, found := k.GetTokenizationRecord(ctx, record.Id); found {
		return fmt.Errorf("tokenization record with ID %d already exists", record.Id)
	}
	
	// Additional validation can be added here as needed
	
	return nil
}