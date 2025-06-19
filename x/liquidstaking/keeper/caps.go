package keeper

import (
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	
	"github.com/rollchains/flora/x/liquidstaking/types"
)

// CheckGlobalLiquidStakingCap checks if tokenizing the given amount would exceed the global cap
func (k Keeper) CheckGlobalLiquidStakingCap(ctx sdk.Context, additionalAmount math.Int) error {
	params := k.GetParams(ctx)
	if !params.Enabled {
		return types.ErrDisabled
	}

	// Get total bonded tokens from staking module
	totalBonded, err := k.stakingKeeper.TotalBondedTokens(ctx)
	if err != nil {
		return err
	}

	// Get current total liquid staked
	currentLiquidStaked := k.GetTotalLiquidStaked(ctx)
	
	// Calculate what the new total would be
	newTotalLiquidStaked := currentLiquidStaked.Add(additionalAmount)
	
	// Calculate the cap amount (totalBonded * capPercentage)
	capAmount := math.LegacyNewDecFromInt(totalBonded).Mul(params.GlobalLiquidStakingCap).TruncateInt()
	
	// Check if the new total would exceed the cap
	if newTotalLiquidStaked.GT(capAmount) {
		return types.ErrGlobalCapExceeded.Wrapf(
			"would exceed global liquid staking cap: current %s + additional %s = %s > cap %s (%.2f%% of %s bonded)",
			currentLiquidStaked,
			additionalAmount,
			newTotalLiquidStaked,
			capAmount,
			params.GlobalLiquidStakingCap.MulInt64(100).MustFloat64(),
			totalBonded,
		)
	}

	return nil
}

// CheckValidatorLiquidCap checks if tokenizing the given amount would exceed the validator's cap
func (k Keeper) CheckValidatorLiquidCap(ctx sdk.Context, validatorAddr string, additionalAmount math.Int) error {
	params := k.GetParams(ctx)
	if !params.Enabled {
		return types.ErrDisabled
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

	// Get current liquid staked for this validator
	currentValidatorLiquid := k.GetValidatorLiquidStaked(ctx, validatorAddr)
	
	// Calculate what the new total would be
	newValidatorLiquid := currentValidatorLiquid.Add(additionalAmount)
	
	// Calculate the cap amount (validator tokens * capPercentage)
	capAmount := math.LegacyNewDecFromInt(validator.Tokens).Mul(params.ValidatorLiquidCap).TruncateInt()
	
	// Check if the new total would exceed the cap
	if newValidatorLiquid.GT(capAmount) {
		return types.ErrValidatorCapExceeded.Wrapf(
			"would exceed validator liquid cap: current %s + additional %s = %s > cap %s (%.2f%% of %s tokens)",
			currentValidatorLiquid,
			additionalAmount,
			newValidatorLiquid,
			capAmount,
			params.ValidatorLiquidCap.MulInt64(100).MustFloat64(),
			validator.Tokens,
		)
	}

	return nil
}

// CheckMinimumAmount checks if the amount meets the minimum requirement
func (k Keeper) CheckMinimumAmount(ctx sdk.Context, amount math.Int) error {
	params := k.GetParams(ctx)
	if amount.LT(params.MinLiquidStakeAmount) {
		return types.ErrAmountTooSmall.Wrapf(
			"amount %s is less than minimum %s",
			amount,
			params.MinLiquidStakeAmount,
		)
	}
	return nil
}

// EnforceTokenizationCaps performs all cap checks before allowing tokenization
func (k Keeper) EnforceTokenizationCaps(ctx sdk.Context, validatorAddr string, amount math.Int) error {
	// Check minimum amount
	if err := k.CheckMinimumAmount(ctx, amount); err != nil {
		return err
	}

	// Check global cap
	if err := k.CheckGlobalLiquidStakingCap(ctx, amount); err != nil {
		return err
	}

	// Check validator cap
	if err := k.CheckValidatorLiquidCap(ctx, validatorAddr, amount); err != nil {
		return err
	}

	return nil
}

// GetTotalLiquidStaked returns the total amount of liquid staked tokens across all validators
func (k Keeper) GetTotalLiquidStaked(ctx sdk.Context) math.Int {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := store.Get(types.TotalLiquidStakedKey)
	if err != nil {
		k.Logger(ctx).Error("failed to get total liquid staked", "error", err)
		return math.ZeroInt()
	}
	if bz == nil {
		return math.ZeroInt()
	}

	var total math.Int
	if err := total.Unmarshal(bz); err != nil {
		// This should never happen, but return zero if it does
		k.Logger(ctx).Error("failed to unmarshal total liquid staked", "error", err)
		return math.ZeroInt()
	}

	return total
}

// SetTotalLiquidStaked sets the total amount of liquid staked tokens
func (k Keeper) SetTotalLiquidStaked(ctx sdk.Context, amount math.Int) {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := amount.Marshal()
	if err != nil {
		// This should never happen
		panic(fmt.Sprintf("failed to marshal total liquid staked: %v", err))
	}
	if err := store.Set(types.TotalLiquidStakedKey, bz); err != nil {
		panic(fmt.Sprintf("failed to set total liquid staked: %v", err))
	}
}

// IncreaseTotalLiquidStaked increases the total liquid staked amount
func (k Keeper) IncreaseTotalLiquidStaked(ctx sdk.Context, amount math.Int) {
	current := k.GetTotalLiquidStaked(ctx)
	k.SetTotalLiquidStaked(ctx, current.Add(amount))
}

// DecreaseTotalLiquidStaked decreases the total liquid staked amount
func (k Keeper) DecreaseTotalLiquidStaked(ctx sdk.Context, amount math.Int) {
	current := k.GetTotalLiquidStaked(ctx)
	newTotal := current.Sub(amount)
	if newTotal.IsNegative() {
		// This should never happen, but set to zero if it does
		k.Logger(ctx).Error("total liquid staked would be negative", "current", current, "decrease", amount)
		newTotal = math.ZeroInt()
	}
	k.SetTotalLiquidStaked(ctx, newTotal)
}

// GetValidatorLiquidStaked returns the amount of liquid staked tokens for a specific validator
func (k Keeper) GetValidatorLiquidStaked(ctx sdk.Context, validatorAddr string) math.Int {
	store := k.storeService.OpenKVStore(ctx)
	key := types.GetValidatorLiquidStakedKey(validatorAddr)
	bz, err := store.Get(key)
	if err != nil {
		k.Logger(ctx).Error("failed to get validator liquid staked", "validator", validatorAddr, "error", err)
		return math.ZeroInt()
	}
	if bz == nil {
		return math.ZeroInt()
	}

	var amount math.Int
	if err := amount.Unmarshal(bz); err != nil {
		// This should never happen, but return zero if it does
		k.Logger(ctx).Error("failed to unmarshal validator liquid staked", "validator", validatorAddr, "error", err)
		return math.ZeroInt()
	}

	return amount
}

// SetValidatorLiquidStaked sets the amount of liquid staked tokens for a validator
func (k Keeper) SetValidatorLiquidStaked(ctx sdk.Context, validatorAddr string, amount math.Int) {
	store := k.storeService.OpenKVStore(ctx)
	key := types.GetValidatorLiquidStakedKey(validatorAddr)
	
	if amount.IsZero() {
		// Remove the key if amount is zero to save space
		if err := store.Delete(key); err != nil {
			panic(fmt.Sprintf("failed to delete validator liquid staked: %v", err))
		}
		return
	}

	bz, err := amount.Marshal()
	if err != nil {
		// This should never happen
		panic(fmt.Sprintf("failed to marshal validator liquid staked: %v", err))
	}
	if err := store.Set(key, bz); err != nil {
		panic(fmt.Sprintf("failed to set validator liquid staked: %v", err))
	}
}

// IncreaseValidatorLiquidStaked increases the liquid staked amount for a validator
func (k Keeper) IncreaseValidatorLiquidStaked(ctx sdk.Context, validatorAddr string, amount math.Int) {
	current := k.GetValidatorLiquidStaked(ctx, validatorAddr)
	k.SetValidatorLiquidStaked(ctx, validatorAddr, current.Add(amount))
}

// DecreaseValidatorLiquidStaked decreases the liquid staked amount for a validator
func (k Keeper) DecreaseValidatorLiquidStaked(ctx sdk.Context, validatorAddr string, amount math.Int) {
	current := k.GetValidatorLiquidStaked(ctx, validatorAddr)
	newAmount := current.Sub(amount)
	if newAmount.IsNegative() {
		// This should never happen, but set to zero if it does
		k.Logger(ctx).Error("validator liquid staked would be negative", "validator", validatorAddr, "current", current, "decrease", amount)
		newAmount = math.ZeroInt()
	}
	k.SetValidatorLiquidStaked(ctx, validatorAddr, newAmount)
}

// GetAllValidatorLiquidStaked returns all validator liquid staked amounts
func (k Keeper) GetAllValidatorLiquidStaked(ctx sdk.Context) map[string]math.Int {
	store := k.storeService.OpenKVStore(ctx)
	iterator, err := store.Iterator(types.ValidatorLiquidStakedPrefix, nil)
	if err != nil {
		k.Logger(ctx).Error("failed to create iterator", "error", err)
		return make(map[string]math.Int)
	}
	defer iterator.Close()

	result := make(map[string]math.Int)
	for ; iterator.Valid(); iterator.Next() {
		// Extract validator address from key
		validatorAddr := string(iterator.Key()[len(types.ValidatorLiquidStakedPrefix):])
		
		var amount math.Int
		if err := amount.Unmarshal(iterator.Value()); err != nil {
			k.Logger(ctx).Error("failed to unmarshal validator liquid staked", "validator", validatorAddr, "error", err)
			continue
		}
		
		result[validatorAddr] = amount
	}

	return result
}