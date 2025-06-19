package keeper

import (
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/rollchains/flora/x/liquidstaking/types"
)

// SetExchangeRate stores an exchange rate for a validator
func (k Keeper) SetExchangeRate(ctx sdk.Context, validatorAddr string, rate math.LegacyDec, lastUpdated time.Time) {
	store := k.storeService.OpenKVStore(ctx)
	
	exchangeRate := types.ExchangeRate{
		ValidatorAddress: validatorAddr,
		Denom:           types.GetLSTDenom(validatorAddr),
		Rate:            rate,
		LastUpdated:     lastUpdated.Unix(),
	}
	
	bz := k.cdc.MustMarshal(&exchangeRate)
	key := types.GetExchangeRateKey(validatorAddr)
	
	if err := store.Set(key, bz); err != nil {
		panic(err)
	}
}

// GetExchangeRate retrieves the exchange rate for a validator
func (k Keeper) GetExchangeRate(ctx sdk.Context, validatorAddr string) (types.ExchangeRate, bool) {
	store := k.storeService.OpenKVStore(ctx)
	key := types.GetExchangeRateKey(validatorAddr)
	
	bz, err := store.Get(key)
	if err != nil {
		panic(err)
	}
	if bz == nil {
		return types.ExchangeRate{}, false
	}

	var rate types.ExchangeRate
	k.cdc.MustUnmarshal(bz, &rate)
	return rate, true
}

// GetOrInitExchangeRate gets the exchange rate or initializes it to 1:1 if not exists
func (k Keeper) GetOrInitExchangeRate(ctx sdk.Context, validatorAddr string) types.ExchangeRate {
	rate, found := k.GetExchangeRate(ctx, validatorAddr)
	if !found {
		// Initialize with 1:1 rate
		k.SetExchangeRate(ctx, validatorAddr, math.LegacyOneDec(), ctx.BlockTime())
		rate = types.ExchangeRate{
			ValidatorAddress: validatorAddr,
			Denom:           types.GetLSTDenom(validatorAddr),
			Rate:            math.LegacyOneDec(),
			LastUpdated:     ctx.BlockTime().Unix(),
		}
	}
	return rate
}

// DeleteExchangeRate removes an exchange rate
func (k Keeper) DeleteExchangeRate(ctx sdk.Context, validatorAddr string) {
	store := k.storeService.OpenKVStore(ctx)
	key := types.GetExchangeRateKey(validatorAddr)
	
	if err := store.Delete(key); err != nil {
		panic(err)
	}
}

// IterateExchangeRates iterates over all exchange rates
func (k Keeper) IterateExchangeRates(ctx sdk.Context, cb func(rate types.ExchangeRate) bool) {
	store := k.storeService.OpenKVStore(ctx)
	
	iterator, err := store.Iterator(types.ExchangeRatePrefix, nil)
	if err != nil {
		panic(err)
	}
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var rate types.ExchangeRate
		k.cdc.MustUnmarshal(iterator.Value(), &rate)
		if cb(rate) {
			break
		}
	}
}

// GetAllExchangeRates returns all exchange rates
func (k Keeper) GetAllExchangeRates(ctx sdk.Context) []types.ExchangeRate {
	var rates []types.ExchangeRate
	k.IterateExchangeRates(ctx, func(rate types.ExchangeRate) bool {
		rates = append(rates, rate)
		return false
	})
	return rates
}

// SetGlobalExchangeRate stores the global exchange rate statistics
func (k Keeper) SetGlobalExchangeRate(ctx sdk.Context, rate types.GlobalExchangeRate) {
	store := k.storeService.OpenKVStore(ctx)
	bz := k.cdc.MustMarshal(&rate)
	
	if err := store.Set(types.GlobalExchangeRateKey, bz); err != nil {
		panic(err)
	}
}

// GetGlobalExchangeRate retrieves the global exchange rate statistics
func (k Keeper) GetGlobalExchangeRate(ctx sdk.Context) (types.GlobalExchangeRate, bool) {
	store := k.storeService.OpenKVStore(ctx)
	
	bz, err := store.Get(types.GlobalExchangeRateKey)
	if err != nil {
		panic(err)
	}
	if bz == nil {
		return types.GlobalExchangeRate{}, false
	}

	var rate types.GlobalExchangeRate
	k.cdc.MustUnmarshal(bz, &rate)
	return rate, true
}

// CalculateExchangeRate calculates the exchange rate for a validator
func (k Keeper) CalculateExchangeRate(ctx sdk.Context, validatorAddr string) (math.LegacyDec, error) {
	valAddr, err := sdk.ValAddressFromBech32(validatorAddr)
	if err != nil {
		return math.LegacyZeroDec(), err
	}

	// Get validator's total delegations
	validator, err := k.stakingKeeper.GetValidator(ctx, valAddr)
	if err != nil {
		return math.LegacyZeroDec(), err
	}

	// Get total LST tokens for this validator
	lstDenom := types.GetLSTDenom(validatorAddr)
	totalLST := k.bankKeeper.GetSupply(ctx, lstDenom).Amount

	// If no LST tokens exist, return 1:1 rate
	if totalLST.IsZero() {
		return math.LegacyOneDec(), nil
	}

	// Get bond denomination
	// TODO: Uncomment when distribution keeper is fixed
	// bondDenom, err := k.stakingKeeper.BondDenom(ctx)
	// if err != nil {
	// 	return math.LegacyZeroDec(), err
	// }

	// TODO: Fix distribution keeper interface for rewards
	// Get accumulated rewards for the validator
	// rewards, err := k.distributionKeeper.GetValidatorAccumulatedRewards(ctx, valAddr)
	// if err != nil {
	// 	return math.LegacyZeroDec(), err
	// }

	// Calculate total value (delegations + rewards)
	// Note: validator.TokensFromShares gives us the total tokens delegated to the validator
	totalDelegated := validator.TokensFromShares(validator.DelegatorShares).TruncateInt()
	// rewardAmount := rewards.AmountOf(bondDenom)
	rewardAmount := math.ZeroInt() // TODO: Add rewards when distribution keeper is fixed
	
	// Total value = delegated tokens + accumulated rewards
	totalValue := totalDelegated.Add(rewardAmount)

	// Exchange rate = total value / total LST supply
	// This means 1 LST token can be redeemed for (exchange rate) native tokens
	exchangeRate := math.LegacyNewDecFromInt(totalValue).Quo(math.LegacyNewDecFromInt(totalLST))

	return exchangeRate, nil
}

// UpdateExchangeRate updates the exchange rate for a specific validator
func (k Keeper) UpdateExchangeRate(ctx sdk.Context, validatorAddr string) error {
	// Calculate new rate
	newRate, err := k.CalculateExchangeRate(ctx, validatorAddr)
	if err != nil {
		return err
	}

	// Get current rate
	currentRate, _ := k.GetExchangeRate(ctx, validatorAddr)
	oldRate := currentRate.Rate
	if oldRate.IsNil() {
		oldRate = math.LegacyOneDec()
	}

	// Update rate
	k.SetExchangeRate(ctx, validatorAddr, newRate, ctx.BlockTime())

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeExchangeRateUpdated,
			sdk.NewAttribute(types.AttributeKeyValidator, validatorAddr),
			sdk.NewAttribute(types.AttributeKeyOldRate, oldRate.String()),
			sdk.NewAttribute(types.AttributeKeyNewRate, newRate.String()),
			sdk.NewAttribute(types.AttributeKeyTimestamp, time.Unix(currentRate.LastUpdated, 0).String()),
		),
	)

	return nil
}

// UpdateAllExchangeRates updates exchange rates for all validators with LST tokens
func (k Keeper) UpdateAllExchangeRates(ctx sdk.Context) error {
	// Track global statistics
	totalStaked := math.ZeroInt()
	totalRewards := math.ZeroInt()
	totalLSTSupply := math.ZeroInt()
	validatorCount := 0

	// Update each validator's rate
	k.stakingKeeper.IterateValidators(ctx, func(index int64, validator stakingtypes.ValidatorI) bool {
		validatorAddr := validator.GetOperator()
		lstDenom := types.GetLSTDenom(validatorAddr)
		
		// Check if this validator has any LST tokens
		supply := k.bankKeeper.GetSupply(ctx, lstDenom).Amount
		if supply.IsPositive() {
			// Update the exchange rate
			if err := k.UpdateExchangeRate(ctx, validatorAddr); err != nil {
				// Log error but continue with other validators
				ctx.Logger().Error("failed to update exchange rate", "validator", validatorAddr, "error", err)
			} else {
				// Add to global statistics
				valAddr, _ := sdk.ValAddressFromBech32(validatorAddr)
				val, _ := k.stakingKeeper.GetValidator(ctx, valAddr)
				
				delegated := val.TokensFromShares(val.DelegatorShares).TruncateInt()
				// TODO: Fix distribution keeper interface
				// rewards, _ := k.distributionKeeper.GetValidatorAccumulatedRewards(ctx, valAddr)
				// bondDenom, _ := k.stakingKeeper.BondDenom(ctx)
				// rewardAmount := rewards.AmountOf(bondDenom).TruncateInt()
				rewardAmount := math.ZeroInt() // TODO: Add rewards when distribution keeper is fixed
				
				totalStaked = totalStaked.Add(delegated)
				totalRewards = totalRewards.Add(rewardAmount)
				totalLSTSupply = totalLSTSupply.Add(supply)
				validatorCount++
			}
		}
		return false
	})

	// Update global exchange rate statistics
	if validatorCount > 0 && totalLSTSupply.IsPositive() {
		globalRate := totalStaked.Add(totalRewards).ToLegacyDec().Quo(totalLSTSupply.ToLegacyDec())
		
		k.SetGlobalExchangeRate(ctx, types.GlobalExchangeRate{
			Rate:           globalRate,
			LastUpdated:    ctx.BlockTime().Unix(),
			TotalStaked:    totalStaked,
			TotalRewards:   totalRewards,
			TotalLstSupply: totalLSTSupply,
		})
	}

	return nil
}

// ApplyExchangeRate calculates the LST amount based on native amount and exchange rate
func (k Keeper) ApplyExchangeRate(ctx sdk.Context, validatorAddr string, nativeAmount math.Int) (math.Int, error) {
	rate := k.GetOrInitExchangeRate(ctx, validatorAddr)
	
	// LST amount = native amount / exchange rate
	// If rate is 1.5, then 150 native tokens = 100 LST tokens
	lstAmount := nativeAmount.ToLegacyDec().Quo(rate.Rate).TruncateInt()
	
	return lstAmount, nil
}

// ApplyInverseExchangeRate calculates the native amount based on LST amount and exchange rate
func (k Keeper) ApplyInverseExchangeRate(ctx sdk.Context, validatorAddr string, lstAmount math.Int) (math.Int, error) {
	rate := k.GetOrInitExchangeRate(ctx, validatorAddr)
	
	// Native amount = LST amount * exchange rate
	// If rate is 1.5, then 100 LST tokens = 150 native tokens
	nativeAmount := lstAmount.ToLegacyDec().Mul(rate.Rate).TruncateInt()
	
	return nativeAmount, nil
}