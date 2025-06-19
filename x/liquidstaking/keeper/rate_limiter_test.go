package keeper_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"

	"github.com/rollchains/flora/x/liquidstaking/keeper"
	"github.com/rollchains/flora/x/liquidstaking/types"
)

func TestGlobalRateLimit(t *testing.T) {
	k, ctx := setupKeeper(t)

	// Set module enabled
	params := types.DefaultParams()
	params.Enabled = true
	require.NoError(t, k.SetParams(ctx, params))

	// Mock total bonded tokens = 10,000,000
	totalBonded := math.NewInt(10000000)
	mockStakingKeeper.TotalBondedTokensFn = func(c context.Context) (math.Int, error) {
		return totalBonded, nil
	}

	// Daily limit should be 5% = 500,000
	// dailyLimit := math.NewInt(500000) // Currently unused

	t.Run("allow tokenization under daily limit", func(t *testing.T) {
		// First tokenization - 200,000
		err := k.CheckGlobalRateLimit(ctx, math.NewInt(200000))
		require.NoError(t, err)

		// Update activity to simulate it was recorded
		k.UpdateTokenizationActivity(ctx, testValAddr1.String(), testAccAddr1.String(), math.NewInt(200000))

		// Second tokenization - 200,000 more (total 400,000 < 500,000)
		err = k.CheckGlobalRateLimit(ctx, math.NewInt(200000))
		require.NoError(t, err)
	})

	t.Run("reject tokenization over daily limit", func(t *testing.T) {
		// Reset activity
		activity := keeper.TokenizationActivity{
			TotalAmount:   math.NewInt(450000),
			LastActivity:  ctx.BlockTime(),
			ActivityCount: 10,
		}
		k.SetGlobalTokenizationActivity(ctx, activity)

		// Try to tokenize 100,000 more (would exceed 500,000 limit)
		err := k.CheckGlobalRateLimit(ctx, math.NewInt(100000))
		require.Error(t, err)
		require.ErrorIs(t, err, types.ErrRateLimitExceeded)
		require.Contains(t, err.Error(), "daily tokenization limit exceeded")
	})

	t.Run("reject tokenization when count limit exceeded", func(t *testing.T) {
		// Reset with high count
		activity := keeper.TokenizationActivity{
			TotalAmount:   math.NewInt(100000),
			LastActivity:  ctx.BlockTime(),
			ActivityCount: 100, // At the limit
		}
		k.SetGlobalTokenizationActivity(ctx, activity)

		// Should reject due to count limit
		err := k.CheckGlobalRateLimit(ctx, math.NewInt(1000))
		require.Error(t, err)
		require.ErrorIs(t, err, types.ErrRateLimitExceeded)
		require.Contains(t, err.Error(), "daily tokenization count limit exceeded")
	})

	t.Run("reset activity after 24 hours", func(t *testing.T) {
		// Set activity near limit
		activity := keeper.TokenizationActivity{
			TotalAmount:   math.NewInt(490000),
			LastActivity:  ctx.BlockTime(),
			ActivityCount: 99,
		}
		k.SetGlobalTokenizationActivity(ctx, activity)

		// Should fail before 24 hours
		err := k.CheckGlobalRateLimit(ctx, math.NewInt(50000))
		require.Error(t, err)

		// Advance time by 25 hours
		newTime := ctx.BlockTime().Add(25 * time.Hour)
		ctx = ctx.WithBlockTime(newTime)

		// Should succeed after reset window (CheckGlobalRateLimit internally resets for checking)
		err = k.CheckGlobalRateLimit(ctx, math.NewInt(50000))
		require.NoError(t, err)

		// Activity in storage should still be the old one (CheckGlobalRateLimit is read-only)
		storedActivity := k.GetGlobalTokenizationActivity(ctx)
		require.Equal(t, math.NewInt(490000), storedActivity.TotalAmount)
		require.Equal(t, uint64(99), storedActivity.ActivityCount)

		// Now actually update the activity (this is when the reset happens in storage)
		k.UpdateTokenizationActivity(ctx, testValAddr1.String(), testAccAddr1.String(), math.NewInt(50000))

		// Now check that activity was reset and updated
		newActivity := k.GetGlobalTokenizationActivity(ctx)
		require.Equal(t, math.NewInt(50000), newActivity.TotalAmount)
		require.Equal(t, uint64(1), newActivity.ActivityCount)
	})
}

func TestValidatorRateLimit(t *testing.T) {
	k, ctx := setupKeeper(t)

	// Set module enabled
	params := types.DefaultParams()
	params.Enabled = true
	require.NoError(t, k.SetParams(ctx, params))

	validatorAddr := testValAddr1.String()
	valAddr := testValAddr1

	// Mock validator with 2,000,000 tokens
	validator := createTestValidator(valAddr, math.NewInt(2000000))
	mockStakingKeeper.GetValidatorFn = func(c context.Context, addr sdk.ValAddress) (stakingtypes.Validator, error) {
		if addr.Equals(valAddr) {
			return validator, nil
		}
		return stakingtypes.Validator{}, errors.New("validator not found")
	}

	// Daily limit should be 10% = 200,000
	// dailyLimit := math.NewInt(200000) // Currently unused

	t.Run("allow tokenization under validator daily limit", func(t *testing.T) {
		// First tokenization - 80,000
		err := k.CheckValidatorRateLimit(ctx, validatorAddr, math.NewInt(80000))
		require.NoError(t, err)

		// Update activity
		k.UpdateTokenizationActivity(ctx, validatorAddr, testAccAddr1.String(), math.NewInt(80000))

		// Second tokenization - 80,000 more (total 160,000 < 200,000)
		err = k.CheckValidatorRateLimit(ctx, validatorAddr, math.NewInt(80000))
		require.NoError(t, err)
	})

	t.Run("reject tokenization over validator daily limit", func(t *testing.T) {
		// Set activity near limit
		activity := keeper.TokenizationActivity{
			TotalAmount:   math.NewInt(180000),
			LastActivity:  ctx.BlockTime(),
			ActivityCount: 5,
		}
		k.SetValidatorTokenizationActivity(ctx, validatorAddr, activity)

		// Try to tokenize 50,000 more (would exceed 200,000 limit)
		err := k.CheckValidatorRateLimit(ctx, validatorAddr, math.NewInt(50000))
		require.Error(t, err)
		require.ErrorIs(t, err, types.ErrRateLimitExceeded)
		require.Contains(t, err.Error(), "validator daily tokenization limit exceeded")
	})

	t.Run("reject when validator count limit exceeded", func(t *testing.T) {
		// Reset with high count
		activity := keeper.TokenizationActivity{
			TotalAmount:   math.NewInt(50000),
			LastActivity:  ctx.BlockTime(),
			ActivityCount: 20, // At the limit
		}
		k.SetValidatorTokenizationActivity(ctx, validatorAddr, activity)

		// Should reject due to count limit
		err := k.CheckValidatorRateLimit(ctx, validatorAddr, math.NewInt(1000))
		require.Error(t, err)
		require.ErrorIs(t, err, types.ErrRateLimitExceeded)
		require.Contains(t, err.Error(), "validator daily tokenization count limit exceeded")
	})

	t.Run("handle invalid validator address", func(t *testing.T) {
		err := k.CheckValidatorRateLimit(ctx, "invalid", math.NewInt(1000))
		require.Error(t, err)
	})
}

func TestUserRateLimit(t *testing.T) {
	k, ctx := setupKeeper(t)

	userAddr := testAccAddr1.String()

	t.Run("allow tokenization under user daily limit", func(t *testing.T) {
		// First few tokenizations should succeed
		for i := 0; i < 4; i++ {
			err := k.CheckUserRateLimit(ctx, userAddr)
			require.NoError(t, err)
			
			// Simulate activity update
			activity := k.GetUserTokenizationActivity(ctx, userAddr)
			activity.ActivityCount++
			activity.LastActivity = ctx.BlockTime()
			k.SetUserTokenizationActivity(ctx, userAddr, activity)
		}
	})

	t.Run("reject when user count limit exceeded", func(t *testing.T) {
		// Set activity at limit
		activity := keeper.TokenizationActivity{
			TotalAmount:   math.NewInt(100000),
			LastActivity:  ctx.BlockTime(),
			ActivityCount: 5, // At the limit
		}
		k.SetUserTokenizationActivity(ctx, userAddr, activity)

		// Should reject due to count limit
		err := k.CheckUserRateLimit(ctx, userAddr)
		require.Error(t, err)
		require.ErrorIs(t, err, types.ErrRateLimitExceeded)
		require.Contains(t, err.Error(), "user daily tokenization count limit exceeded")
	})

	t.Run("reset user activity after 24 hours", func(t *testing.T) {
		// Set activity at limit
		activity := keeper.TokenizationActivity{
			TotalAmount:   math.NewInt(100000),
			LastActivity:  ctx.BlockTime(),
			ActivityCount: 5,
		}
		k.SetUserTokenizationActivity(ctx, userAddr, activity)

		// Should fail before 24 hours
		err := k.CheckUserRateLimit(ctx, userAddr)
		require.Error(t, err)

		// Advance time by 25 hours
		newTime := ctx.BlockTime().Add(25 * time.Hour)
		ctx = ctx.WithBlockTime(newTime)

		// Should succeed after reset
		err = k.CheckUserRateLimit(ctx, userAddr)
		require.NoError(t, err)
	})
}

func TestUpdateTokenizationActivity(t *testing.T) {
	k, ctx := setupKeeper(t)

	validatorAddr := testValAddr1.String()
	userAddr := testAccAddr1.String()
	amount := math.NewInt(100000)

	t.Run("update all activity levels", func(t *testing.T) {
		// Initial state - all should be zero
		globalActivity := k.GetGlobalTokenizationActivity(ctx)
		require.True(t, globalActivity.TotalAmount.IsZero())
		require.Equal(t, uint64(0), globalActivity.ActivityCount)

		validatorActivity := k.GetValidatorTokenizationActivity(ctx, validatorAddr)
		require.True(t, validatorActivity.TotalAmount.IsZero())
		require.Equal(t, uint64(0), validatorActivity.ActivityCount)

		userActivity := k.GetUserTokenizationActivity(ctx, userAddr)
		require.True(t, userActivity.TotalAmount.IsZero())
		require.Equal(t, uint64(0), userActivity.ActivityCount)

		// Update activity
		k.UpdateTokenizationActivity(ctx, validatorAddr, userAddr, amount)

		// Check all were updated
		globalActivity = k.GetGlobalTokenizationActivity(ctx)
		require.Equal(t, amount, globalActivity.TotalAmount)
		require.Equal(t, uint64(1), globalActivity.ActivityCount)
		require.Equal(t, ctx.BlockTime(), globalActivity.LastActivity)

		validatorActivity = k.GetValidatorTokenizationActivity(ctx, validatorAddr)
		require.Equal(t, amount, validatorActivity.TotalAmount)
		require.Equal(t, uint64(1), validatorActivity.ActivityCount)
		require.Equal(t, ctx.BlockTime(), validatorActivity.LastActivity)

		userActivity = k.GetUserTokenizationActivity(ctx, userAddr)
		require.Equal(t, amount, userActivity.TotalAmount)
		require.Equal(t, uint64(1), userActivity.ActivityCount)
		require.Equal(t, ctx.BlockTime(), userActivity.LastActivity)

		// Update again
		k.UpdateTokenizationActivity(ctx, validatorAddr, userAddr, amount)

		// Check cumulative updates
		globalActivity = k.GetGlobalTokenizationActivity(ctx)
		require.Equal(t, amount.Mul(math.NewInt(2)), globalActivity.TotalAmount)
		require.Equal(t, uint64(2), globalActivity.ActivityCount)
	})

	t.Run("reset on stale activity", func(t *testing.T) {
		// Set old activity
		oldTime := ctx.BlockTime().Add(-30 * time.Hour)
		oldActivity := keeper.TokenizationActivity{
			TotalAmount:   math.NewInt(500000),
			LastActivity:  oldTime,
			ActivityCount: 50,
		}
		k.SetGlobalTokenizationActivity(ctx, oldActivity)

		// Update should reset
		k.UpdateTokenizationActivity(ctx, validatorAddr, userAddr, amount)

		// Check reset happened
		globalActivity := k.GetGlobalTokenizationActivity(ctx)
		require.Equal(t, amount, globalActivity.TotalAmount) // Reset to just the new amount
		require.Equal(t, uint64(1), globalActivity.ActivityCount) // Reset to 1
		require.Equal(t, ctx.BlockTime(), globalActivity.LastActivity)
	})
}

func TestEnforceTokenizationRateLimits(t *testing.T) {
	k, ctx := setupKeeper(t)

	// Set module enabled
	params := types.DefaultParams()
	params.Enabled = true
	require.NoError(t, k.SetParams(ctx, params))

	validatorAddr := testValAddr1.String()
	userAddr := testAccAddr1.String()
	valAddr := testValAddr1

	// Mock validator with 2,000,000 tokens
	validator := createTestValidator(valAddr, math.NewInt(2000000))
	mockStakingKeeper.GetValidatorFn = func(c context.Context, addr sdk.ValAddress) (stakingtypes.Validator, error) {
		if addr.Equals(valAddr) {
			return validator, nil
		}
		return stakingtypes.Validator{}, errors.New("validator not found")
	}

	// Mock total bonded = 10,000,000
	mockStakingKeeper.TotalBondedTokensFn = func(c context.Context) (math.Int, error) {
		return math.NewInt(10000000), nil
	}

	t.Run("pass all rate limit checks", func(t *testing.T) {
		// Small amount should pass all checks
		err := k.EnforceTokenizationRateLimits(ctx, validatorAddr, userAddr, math.NewInt(10000))
		require.NoError(t, err)
	})

	t.Run("fail global rate limit", func(t *testing.T) {
		// Set global activity near limit
		activity := keeper.TokenizationActivity{
			TotalAmount:   math.NewInt(490000), // Near 500,000 limit
			LastActivity:  ctx.BlockTime(),
			ActivityCount: 10,
		}
		k.SetGlobalTokenizationActivity(ctx, activity)

		// Should fail global check first
		err := k.EnforceTokenizationRateLimits(ctx, validatorAddr, userAddr, math.NewInt(50000))
		require.Error(t, err)
		require.ErrorIs(t, err, types.ErrRateLimitExceeded)
		require.Contains(t, err.Error(), "daily tokenization limit exceeded")
	})

	t.Run("fail validator rate limit", func(t *testing.T) {
		// Reset global activity
		k.SetGlobalTokenizationActivity(ctx, keeper.TokenizationActivity{})

		// Set validator activity near limit
		activity := keeper.TokenizationActivity{
			TotalAmount:   math.NewInt(190000), // Near 200,000 limit for validator
			LastActivity:  ctx.BlockTime(),
			ActivityCount: 5,
		}
		k.SetValidatorTokenizationActivity(ctx, validatorAddr, activity)

		// Should fail validator check
		err := k.EnforceTokenizationRateLimits(ctx, validatorAddr, userAddr, math.NewInt(50000))
		require.Error(t, err)
		require.ErrorIs(t, err, types.ErrRateLimitExceeded)
		require.Contains(t, err.Error(), "validator daily tokenization limit exceeded")
	})

	t.Run("fail user rate limit", func(t *testing.T) {
		// Reset global and validator activity
		k.SetGlobalTokenizationActivity(ctx, keeper.TokenizationActivity{})
		k.SetValidatorTokenizationActivity(ctx, validatorAddr, keeper.TokenizationActivity{})

		// Set user activity at limit
		activity := keeper.TokenizationActivity{
			TotalAmount:   math.NewInt(10000),
			LastActivity:  ctx.BlockTime(),
			ActivityCount: 5, // At limit
		}
		k.SetUserTokenizationActivity(ctx, userAddr, activity)

		// Should fail user check
		err := k.EnforceTokenizationRateLimits(ctx, validatorAddr, userAddr, math.NewInt(1000))
		require.Error(t, err)
		require.ErrorIs(t, err, types.ErrRateLimitExceeded)
		require.Contains(t, err.Error(), "user daily tokenization count limit exceeded")
	})
}

func TestActivityPersistence(t *testing.T) {
	k, ctx := setupKeeper(t)

	t.Run("global activity persistence", func(t *testing.T) {
		activity := keeper.TokenizationActivity{
			TotalAmount:   math.NewInt(123456),
			LastActivity:  ctx.BlockTime(),
			ActivityCount: 42,
		}
		k.SetGlobalTokenizationActivity(ctx, activity)

		// Retrieve and verify
		retrieved := k.GetGlobalTokenizationActivity(ctx)
		require.Equal(t, activity.TotalAmount, retrieved.TotalAmount)
		require.Equal(t, activity.ActivityCount, retrieved.ActivityCount)
		require.Equal(t, activity.LastActivity.Unix(), retrieved.LastActivity.Unix())
	})

	t.Run("validator activity persistence", func(t *testing.T) {
		validatorAddr := testValAddr1.String()
		activity := keeper.TokenizationActivity{
			TotalAmount:   math.NewInt(654321),
			LastActivity:  ctx.BlockTime(),
			ActivityCount: 24,
		}
		k.SetValidatorTokenizationActivity(ctx, validatorAddr, activity)

		// Retrieve and verify
		retrieved := k.GetValidatorTokenizationActivity(ctx, validatorAddr)
		require.Equal(t, activity.TotalAmount, retrieved.TotalAmount)
		require.Equal(t, activity.ActivityCount, retrieved.ActivityCount)
		require.Equal(t, activity.LastActivity.Unix(), retrieved.LastActivity.Unix())

		// Different validator should have default values
		otherRetrieved := k.GetValidatorTokenizationActivity(ctx, testValAddr2.String())
		require.True(t, otherRetrieved.TotalAmount.IsZero())
		require.Equal(t, uint64(0), otherRetrieved.ActivityCount)
	})

	t.Run("user activity persistence", func(t *testing.T) {
		userAddr := testAccAddr1.String()
		activity := keeper.TokenizationActivity{
			TotalAmount:   math.NewInt(789012),
			LastActivity:  ctx.BlockTime(),
			ActivityCount: 3,
		}
		k.SetUserTokenizationActivity(ctx, userAddr, activity)

		// Retrieve and verify
		retrieved := k.GetUserTokenizationActivity(ctx, userAddr)
		require.Equal(t, activity.TotalAmount, retrieved.TotalAmount)
		require.Equal(t, activity.ActivityCount, retrieved.ActivityCount)
		require.Equal(t, activity.LastActivity.Unix(), retrieved.LastActivity.Unix())

		// Different user should have default values
		otherRetrieved := k.GetUserTokenizationActivity(ctx, "flora1xac0prsfpqrppx9yh2pn8v9gzjxf4lqq4tfmz")
		require.True(t, otherRetrieved.TotalAmount.IsZero())
		require.Equal(t, uint64(0), otherRetrieved.ActivityCount)
	})
}