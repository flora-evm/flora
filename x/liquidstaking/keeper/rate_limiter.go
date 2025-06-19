package keeper

import (
	"encoding/json"
	"fmt"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	
	"github.com/rollchains/flora/x/liquidstaking/types"
)

// GetRateLimitPeriod returns the rate limit period from module parameters
func (k Keeper) GetRateLimitPeriod(ctx sdk.Context) time.Duration {
	params := k.GetParams(ctx)
	return time.Duration(params.RateLimitPeriodHours) * time.Hour
}

// TokenizationActivity tracks tokenization activity for rate limiting
type TokenizationActivity struct {
	TotalAmount   math.Int  `json:"total_amount"`
	LastActivity  time.Time `json:"last_activity"`
	ActivityCount uint64    `json:"activity_count"`
}

// Marshal custom marshal for TokenizationActivity
func (ta TokenizationActivity) Marshal() ([]byte, error) {
	// Use JSON encoding for simplicity
	return json.Marshal(ta)
}

// Unmarshal custom unmarshal for TokenizationActivity
func (ta *TokenizationActivity) Unmarshal(data []byte) error {
	return json.Unmarshal(data, ta)
}

// CheckGlobalRateLimit checks if the global tokenization rate limit would be exceeded
func (k Keeper) CheckGlobalRateLimit(ctx sdk.Context, additionalAmount math.Int) error {
	params := k.GetParams(ctx)
	if !params.Enabled {
		return types.ErrDisabled
	}

	// Get global tokenization activity
	activity := k.GetGlobalTokenizationActivity(ctx)
	
	// Check if we need to reset the activity window
	rateLimitPeriod := k.GetRateLimitPeriod(ctx)
	if ctx.BlockTime().Sub(activity.LastActivity) > rateLimitPeriod {
		// Reset activity for new period
		activity = TokenizationActivity{
			TotalAmount:   math.ZeroInt(),
			LastActivity:  ctx.BlockTime(),
			ActivityCount: 0,
		}
	}

	// Calculate what the new total would be
	newTotal := activity.TotalAmount.Add(additionalAmount)
	
	// Get total bonded tokens for rate limit calculation
	totalBonded, err := k.stakingKeeper.TotalBondedTokens(ctx)
	if err != nil {
		return err
	}

	// Check if daily tokenization exceeds allowed rate
	dailyLimit := math.LegacyNewDecFromInt(totalBonded).Mul(params.GlobalDailyTokenizationPercent).TruncateInt()
	
	if newTotal.GT(dailyLimit) {
		// Emit rate limit exceeded event
		ctx.EventManager().EmitEvent(types.RateLimitExceededEvent{
			LimitType:      "global",
			Address:        "global",
			CurrentUsage:   activity.TotalAmount.String(),
			MaxUsage:       dailyLimit.String(),
			RejectedAmount: additionalAmount.String(),
			WindowEnd:      activity.LastActivity.Add(rateLimitPeriod).Format(time.RFC3339),
		}.ToEvent())
		
		// Call hook
		k.GetHooks().OnRateLimitExceeded(ctx, "global", "global", additionalAmount)
		
		return types.ErrRateLimitExceeded.Wrapf(
			"daily tokenization limit exceeded: current %s + additional %s = %s > limit %s (%s%% of %s bonded)",
			activity.TotalAmount,
			additionalAmount,
			newTotal,
			dailyLimit,
			params.GlobalDailyTokenizationPercent.Mul(math.LegacyNewDec(100)).TruncateInt(),
			totalBonded,
		)
	}
	
	// Check if approaching limit (configurable threshold) and emit warning
	warningThreshold := math.LegacyNewDecFromInt(dailyLimit).Mul(params.WarningThresholdPercent).TruncateInt()
	if newTotal.GT(warningThreshold) {
		percentageUsed := newTotal.Mul(math.NewInt(100)).Quo(dailyLimit)
		ctx.EventManager().EmitEvent(types.RateLimitWarningEvent{
			LimitType:      "global",
			Address:        "global",
			CurrentUsage:   newTotal.String(),
			MaxUsage:       dailyLimit.String(),
			PercentageUsed: percentageUsed.String(),
			LimitThreshold: params.WarningThresholdPercent.Mul(math.LegacyNewDec(100)).TruncateInt().String(),
		}.ToEvent())
	}

	// Check activity count limit
	if activity.ActivityCount >= params.GlobalDailyTokenizationCount {
		// Emit rate limit exceeded event for count
		ctx.EventManager().EmitEvent(types.RateLimitExceededEvent{
			LimitType:      "global",
			Address:        "global",
			CurrentUsage:   fmt.Sprintf("%d", activity.ActivityCount),
			MaxUsage:       fmt.Sprintf("%d", params.GlobalDailyTokenizationCount),
			RejectedAmount: "1",
			WindowEnd:      activity.LastActivity.Add(rateLimitPeriod).Format(time.RFC3339),
		}.ToEvent())
		
		// Call hook
		k.GetHooks().OnRateLimitExceeded(ctx, "global", "global", math.OneInt())
		
		return types.ErrRateLimitExceeded.Wrapf(
			"daily tokenization count limit exceeded: %d >= %d",
			activity.ActivityCount,
			params.GlobalDailyTokenizationCount,
		)
	}

	return nil
}

// CheckValidatorRateLimit checks if the validator-specific rate limit would be exceeded
func (k Keeper) CheckValidatorRateLimit(ctx sdk.Context, validatorAddr string, additionalAmount math.Int) error {
	params := k.GetParams(ctx)
	if !params.Enabled {
		return types.ErrDisabled
	}

	// Get validator tokenization activity
	activity := k.GetValidatorTokenizationActivity(ctx, validatorAddr)
	
	// Check if we need to reset the activity window
	rateLimitPeriod := k.GetRateLimitPeriod(ctx)
	if ctx.BlockTime().Sub(activity.LastActivity) > rateLimitPeriod {
		// Reset activity for new period
		activity = TokenizationActivity{
			TotalAmount:   math.ZeroInt(),
			LastActivity:  ctx.BlockTime(),
			ActivityCount: 0,
		}
	}

	// Get validator info
	valAddr, err := sdk.ValAddressFromBech32(validatorAddr)
	if err != nil {
		return err
	}

	validator, err := k.stakingKeeper.GetValidator(ctx, valAddr)
	if err != nil {
		return err
	}

	// Calculate what the new total would be
	newTotal := activity.TotalAmount.Add(additionalAmount)
	
	// Check if daily validator tokenization exceeds allowed rate
	dailyLimit := math.LegacyNewDecFromInt(validator.Tokens).Mul(params.ValidatorDailyTokenizationPercent).TruncateInt()
	
	if newTotal.GT(dailyLimit) {
		return types.ErrRateLimitExceeded.Wrapf(
			"validator daily tokenization limit exceeded: current %s + additional %s = %s > limit %s (%s%% of %s tokens)",
			activity.TotalAmount,
			additionalAmount,
			newTotal,
			dailyLimit,
			params.ValidatorDailyTokenizationPercent.Mul(math.LegacyNewDec(100)).TruncateInt(),
			validator.Tokens,
		)
	}

	// Check activity count limit
	if activity.ActivityCount >= params.ValidatorDailyTokenizationCount {
		return types.ErrRateLimitExceeded.Wrapf(
			"validator daily tokenization count limit exceeded: %d >= %d",
			activity.ActivityCount,
			params.ValidatorDailyTokenizationCount,
		)
	}

	return nil
}

// CheckUserRateLimit checks if the user-specific rate limit would be exceeded
func (k Keeper) CheckUserRateLimit(ctx sdk.Context, userAddr string) error {
	params := k.GetParams(ctx)
	if !params.Enabled {
		return types.ErrDisabled
	}
	activity := k.GetUserTokenizationActivity(ctx, userAddr)
	
	// Check if we need to reset the activity window
	rateLimitPeriod := k.GetRateLimitPeriod(ctx)
	if ctx.BlockTime().Sub(activity.LastActivity) > rateLimitPeriod {
		// Reset activity for new period
		activity = TokenizationActivity{
			TotalAmount:   math.ZeroInt(),
			LastActivity:  ctx.BlockTime(),
			ActivityCount: 0,
		}
	}

	// Check activity count limit
	if activity.ActivityCount >= params.UserDailyTokenizationCount {
		return types.ErrRateLimitExceeded.Wrapf(
			"user daily tokenization count limit exceeded: %d >= %d",
			activity.ActivityCount,
			params.UserDailyTokenizationCount,
		)
	}

	return nil
}

// UpdateTokenizationActivity updates the tokenization activity after a successful tokenization
func (k Keeper) UpdateTokenizationActivity(ctx sdk.Context, validatorAddr, userAddr string, amount math.Int) {
	rateLimitPeriod := k.GetRateLimitPeriod(ctx)
	
	// Update global activity
	globalActivity := k.GetGlobalTokenizationActivity(ctx)
	if ctx.BlockTime().Sub(globalActivity.LastActivity) > rateLimitPeriod {
		globalActivity = TokenizationActivity{
			TotalAmount:   amount,
			LastActivity:  ctx.BlockTime(),
			ActivityCount: 1,
		}
	} else {
		globalActivity.TotalAmount = globalActivity.TotalAmount.Add(amount)
		globalActivity.LastActivity = ctx.BlockTime()
		globalActivity.ActivityCount++
	}
	k.SetGlobalTokenizationActivity(ctx, globalActivity)

	// Emit activity tracked event for global
	ctx.EventManager().EmitEvent(types.ActivityTrackedEvent{
		LimitType:     "global",
		Address:       "global",
		Amount:        amount.String(),
		TotalAmount:   globalActivity.TotalAmount.String(),
		ActivityCount: fmt.Sprintf("%d", globalActivity.ActivityCount),
		WindowStart:   globalActivity.LastActivity.Add(-rateLimitPeriod).Format(time.RFC3339),
		WindowEnd:     globalActivity.LastActivity.Add(rateLimitPeriod).Format(time.RFC3339),
	}.ToEvent())

	// Update validator activity
	validatorActivity := k.GetValidatorTokenizationActivity(ctx, validatorAddr)
	if ctx.BlockTime().Sub(validatorActivity.LastActivity) > rateLimitPeriod {
		validatorActivity = TokenizationActivity{
			TotalAmount:   amount,
			LastActivity:  ctx.BlockTime(),
			ActivityCount: 1,
		}
	} else {
		validatorActivity.TotalAmount = validatorActivity.TotalAmount.Add(amount)
		validatorActivity.LastActivity = ctx.BlockTime()
		validatorActivity.ActivityCount++
	}
	k.SetValidatorTokenizationActivity(ctx, validatorAddr, validatorActivity)

	// Emit activity tracked event for validator
	ctx.EventManager().EmitEvent(types.ActivityTrackedEvent{
		LimitType:     "validator",
		Address:       validatorAddr,
		Amount:        amount.String(),
		TotalAmount:   validatorActivity.TotalAmount.String(),
		ActivityCount: fmt.Sprintf("%d", validatorActivity.ActivityCount),
		WindowStart:   validatorActivity.LastActivity.Add(-rateLimitPeriod).Format(time.RFC3339),
		WindowEnd:     validatorActivity.LastActivity.Add(rateLimitPeriod).Format(time.RFC3339),
	}.ToEvent())

	// Update user activity
	userActivity := k.GetUserTokenizationActivity(ctx, userAddr)
	if ctx.BlockTime().Sub(userActivity.LastActivity) > rateLimitPeriod {
		userActivity = TokenizationActivity{
			TotalAmount:   amount,
			LastActivity:  ctx.BlockTime(),
			ActivityCount: 1,
		}
	} else {
		userActivity.TotalAmount = userActivity.TotalAmount.Add(amount)
		userActivity.LastActivity = ctx.BlockTime()
		userActivity.ActivityCount++
	}
	k.SetUserTokenizationActivity(ctx, userAddr, userActivity)

	// Emit activity tracked event for user
	ctx.EventManager().EmitEvent(types.ActivityTrackedEvent{
		LimitType:     "user",
		Address:       userAddr,
		Amount:        amount.String(),
		TotalAmount:   userActivity.TotalAmount.String(),
		ActivityCount: fmt.Sprintf("%d", userActivity.ActivityCount),
		WindowStart:   userActivity.LastActivity.Add(-rateLimitPeriod).Format(time.RFC3339),
		WindowEnd:     userActivity.LastActivity.Add(rateLimitPeriod).Format(time.RFC3339),
	}.ToEvent())
}

// GetGlobalTokenizationActivity retrieves the global tokenization activity
func (k Keeper) GetGlobalTokenizationActivity(ctx sdk.Context) TokenizationActivity {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := store.Get(types.GlobalTokenizationActivityKey)
	if err != nil {
		panic(err)
	}
	if bz == nil {
		return TokenizationActivity{
			TotalAmount:   math.ZeroInt(),
			LastActivity:  time.Time{},
			ActivityCount: 0,
		}
	}

	var activity TokenizationActivity
	if err := activity.Unmarshal(bz); err != nil {
		panic(err)
	}
	return activity
}

// SetGlobalTokenizationActivity sets the global tokenization activity
func (k Keeper) SetGlobalTokenizationActivity(ctx sdk.Context, activity TokenizationActivity) {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := activity.Marshal()
	if err != nil {
		panic(err)
	}
	err = store.Set(types.GlobalTokenizationActivityKey, bz)
	if err != nil {
		panic(err)
	}
}

// GetValidatorTokenizationActivity retrieves the validator's tokenization activity
func (k Keeper) GetValidatorTokenizationActivity(ctx sdk.Context, validatorAddr string) TokenizationActivity {
	store := k.storeService.OpenKVStore(ctx)
	key := types.GetValidatorTokenizationActivityKey(validatorAddr)
	bz, err := store.Get(key)
	if err != nil {
		panic(err)
	}
	if bz == nil {
		return TokenizationActivity{
			TotalAmount:   math.ZeroInt(),
			LastActivity:  time.Time{},
			ActivityCount: 0,
		}
	}

	var activity TokenizationActivity
	if err := activity.Unmarshal(bz); err != nil {
		panic(err)
	}
	return activity
}

// SetValidatorTokenizationActivity sets the validator's tokenization activity
func (k Keeper) SetValidatorTokenizationActivity(ctx sdk.Context, validatorAddr string, activity TokenizationActivity) {
	store := k.storeService.OpenKVStore(ctx)
	key := types.GetValidatorTokenizationActivityKey(validatorAddr)
	bz, err := activity.Marshal()
	if err != nil {
		panic(err)
	}
	err = store.Set(key, bz)
	if err != nil {
		panic(err)
	}
}

// GetUserTokenizationActivity retrieves the user's tokenization activity
func (k Keeper) GetUserTokenizationActivity(ctx sdk.Context, userAddr string) TokenizationActivity {
	store := k.storeService.OpenKVStore(ctx)
	key := types.GetUserTokenizationActivityKey(userAddr)
	bz, err := store.Get(key)
	if err != nil {
		panic(err)
	}
	if bz == nil {
		return TokenizationActivity{
			TotalAmount:   math.ZeroInt(),
			LastActivity:  time.Time{},
			ActivityCount: 0,
		}
	}

	var activity TokenizationActivity
	if err := activity.Unmarshal(bz); err != nil {
		panic(err)
	}
	return activity
}

// SetUserTokenizationActivity sets the user's tokenization activity
func (k Keeper) SetUserTokenizationActivity(ctx sdk.Context, userAddr string, activity TokenizationActivity) {
	store := k.storeService.OpenKVStore(ctx)
	key := types.GetUserTokenizationActivityKey(userAddr)
	bz, err := activity.Marshal()
	if err != nil {
		panic(err)
	}
	err = store.Set(key, bz)
	if err != nil {
		panic(err)
	}
}

// EnforceTokenizationRateLimits performs all rate limit checks before allowing tokenization
func (k Keeper) EnforceTokenizationRateLimits(ctx sdk.Context, validatorAddr, userAddr string, amount math.Int) error {
	// Check global rate limit
	if err := k.CheckGlobalRateLimit(ctx, amount); err != nil {
		return err
	}

	// Check validator rate limit
	if err := k.CheckValidatorRateLimit(ctx, validatorAddr, amount); err != nil {
		return err
	}

	// Check user rate limit
	if err := k.CheckUserRateLimit(ctx, userAddr); err != nil {
		return err
	}

	return nil
}