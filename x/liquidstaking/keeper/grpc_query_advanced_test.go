package keeper_test

import (
	"testing"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/rollchains/flora/x/liquidstaking/keeper"
	"github.com/rollchains/flora/x/liquidstaking/types"
)

// TestRateLimitStatus tests the RateLimitStatus query
func TestRateLimitStatus(t *testing.T) {
	// TODO: Uncomment after proto regeneration
	t.Skip("Skipping until proto regeneration")
	
	/*
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
	
	t.Run("query global rate limit status", func(t *testing.T) {
		// Set some global activity
		activity := keeper.TokenizationActivity{
			TotalAmount:   math.NewInt(200000),
			LastActivity:  ctx.BlockTime(),
			ActivityCount: 10,
		}
		k.SetGlobalTokenizationActivity(ctx, activity)
		
		// Query global rate limit status
		resp, err := k.RateLimitStatus(ctx, &types.QueryRateLimitStatusRequest{
			Address: "global",
		})
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Len(t, resp.RateLimits, 1)
		
		globalLimit := resp.RateLimits[0]
		require.Equal(t, "global", globalLimit.LimitType)
		require.Equal(t, math.NewInt(200000), globalLimit.CurrentAmount)
		require.Equal(t, math.NewInt(500000), globalLimit.MaxAmount) // 5% of 10M
		require.Equal(t, uint64(10), globalLimit.CurrentCount)
		require.Equal(t, uint64(100), globalLimit.MaxCount)
	})
	
	t.Run("query validator rate limit status", func(t *testing.T) {
		valAddr := testValAddr1
		validator := createTestValidator(valAddr, math.NewInt(2000000))
		mockStakingKeeper.GetValidatorFn = func(c context.Context, addr sdk.ValAddress) (stakingtypes.Validator, error) {
			if addr.Equals(valAddr) {
				return validator, nil
			}
			return stakingtypes.Validator{}, errors.New("validator not found")
		}
		
		// Set validator activity
		activity := keeper.TokenizationActivity{
			TotalAmount:   math.NewInt(100000),
			LastActivity:  ctx.BlockTime(),
			ActivityCount: 5,
		}
		k.SetValidatorTokenizationActivity(ctx, valAddr.String(), activity)
		
		// Query validator rate limit status
		resp, err := k.RateLimitStatus(ctx, &types.QueryRateLimitStatusRequest{
			Address: valAddr.String(),
		})
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Len(t, resp.RateLimits, 1)
		
		valLimit := resp.RateLimits[0]
		require.Equal(t, "validator", valLimit.LimitType)
		require.Equal(t, math.NewInt(100000), valLimit.CurrentAmount)
		require.Equal(t, math.NewInt(200000), valLimit.MaxAmount) // 10% of 2M
		require.Equal(t, uint64(5), valLimit.CurrentCount)
		require.Equal(t, uint64(20), valLimit.MaxCount)
	})
	
	t.Run("query user rate limit status", func(t *testing.T) {
		userAddr := testAccAddr1.String()
		
		// Set user activity
		activity := keeper.TokenizationActivity{
			TotalAmount:   math.NewInt(50000),
			LastActivity:  ctx.BlockTime(),
			ActivityCount: 3,
		}
		k.SetUserTokenizationActivity(ctx, userAddr, activity)
		
		// Query user rate limit status
		resp, err := k.RateLimitStatus(ctx, &types.QueryRateLimitStatusRequest{
			Address: userAddr,
		})
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Len(t, resp.RateLimits, 1)
		
		userLimit := resp.RateLimits[0]
		require.Equal(t, "user", userLimit.LimitType)
		require.Equal(t, math.NewInt(50000), userLimit.CurrentAmount)
		require.True(t, userLimit.MaxAmount.IsZero()) // Users don't have amount limits
		require.Equal(t, uint64(3), userLimit.CurrentCount)
		require.Equal(t, uint64(5), userLimit.MaxCount)
	})
	
	t.Run("query with expired window", func(t *testing.T) {
		// Set old global activity
		oldTime := ctx.BlockTime().Add(-30 * time.Hour)
		activity := keeper.TokenizationActivity{
			TotalAmount:   math.NewInt(400000),
			LastActivity:  oldTime,
			ActivityCount: 50,
		}
		k.SetGlobalTokenizationActivity(ctx, activity)
		
		// Query should show reset values
		resp, err := k.RateLimitStatus(ctx, &types.QueryRateLimitStatusRequest{
			Address: "global",
		})
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Len(t, resp.RateLimits, 1)
		
		globalLimit := resp.RateLimits[0]
		require.True(t, globalLimit.CurrentAmount.IsZero())
		require.Equal(t, uint64(0), globalLimit.CurrentCount)
		require.Equal(t, ctx.BlockTime(), globalLimit.WindowStart)
		require.Equal(t, ctx.BlockTime().Add(24*time.Hour), globalLimit.WindowEnd)
	})
	
	t.Run("query with invalid address", func(t *testing.T) {
		resp, err := k.RateLimitStatus(ctx, &types.QueryRateLimitStatusRequest{
			Address: "invalid-address",
		})
		require.Error(t, err)
		require.Nil(t, resp)
		require.Contains(t, err.Error(), "invalid address format")
	})
	*/
}

// TestTokenizationStatistics tests the TokenizationStatistics query
func TestTokenizationStatistics(t *testing.T) {
	// TODO: Uncomment after proto regeneration
	t.Skip("Skipping until proto regeneration")
	
	/*
	k, ctx := setupKeeper(t)
	
	// Create some tokenization records
	records := []types.TokenizationRecord{
		types.NewTokenizationRecordWithDenom(
			1,
			testValAddr1.String(),
			testAccAddr1.String(),
			math.NewInt(1000000),
			"liquidstake/val1/1",
		),
		types.NewTokenizationRecordWithDenom(
			2,
			testValAddr1.String(),
			testAccAddr2.String(),
			math.NewInt(2000000),
			"liquidstake/val1/2",
		),
		types.NewTokenizationRecordWithDenom(
			3,
			testValAddr2.String(),
			testAccAddr1.String(),
			math.NewInt(1500000),
			"liquidstake/val2/3",
		),
		// A record with zero amount (fully redeemed)
		types.NewTokenizationRecordWithDenom(
			4,
			testValAddr2.String(),
			testAccAddr2.String(),
			math.ZeroInt(),
			"liquidstake/val2/4",
		),
	}
	
	// Store the records
	for _, record := range records {
		k.SetTokenizationRecordWithIndexes(ctx, record)
	}
	k.SetLastTokenizationRecordID(ctx, 4)
	
	// Set liquid staked amounts
	k.UpdateLiquidStakedAmounts(ctx, testValAddr1.String(), math.NewInt(3000000), true)
	k.UpdateLiquidStakedAmounts(ctx, testValAddr2.String(), math.NewInt(1500000), true)
	
	// Query statistics
	resp, err := k.TokenizationStatistics(ctx, &types.QueryTokenizationStatisticsRequest{})
	require.NoError(t, err)
	require.NotNil(t, resp)
	
	// Verify response
	require.Equal(t, math.NewInt(4500000), resp.TotalTokenized) // Sum of all active records
	require.Equal(t, math.NewInt(4500000), resp.ActiveLiquidStaked)
	require.Equal(t, uint64(4), resp.TotalRecords)
	require.Equal(t, uint64(3), resp.ActiveRecords) // Excluding the zero-amount record
	require.Equal(t, math.NewInt(1500000), resp.AverageRecordSize) // 4,500,000 / 3
	require.Equal(t, uint64(2), resp.ValidatorsWithLiquidStake)
	require.Equal(t, uint64(4), resp.TotalDenomsCreated)
	
	t.Run("empty state", func(t *testing.T) {
		// Create new keeper with empty state
		k2, ctx2 := setupKeeper(t)
		
		resp, err := k2.TokenizationStatistics(ctx2, &types.QueryTokenizationStatisticsRequest{})
		require.NoError(t, err)
		require.NotNil(t, resp)
		
		require.True(t, resp.TotalTokenized.IsZero())
		require.True(t, resp.ActiveLiquidStaked.IsZero())
		require.Equal(t, uint64(0), resp.TotalRecords)
		require.Equal(t, uint64(0), resp.ActiveRecords)
		require.True(t, resp.AverageRecordSize.IsZero())
		require.Equal(t, uint64(0), resp.ValidatorsWithLiquidStake)
		require.Equal(t, uint64(0), resp.TotalDenomsCreated)
	})
	*/
}

// TestValidatorStatistics tests the ValidatorStatistics query
func TestValidatorStatistics(t *testing.T) {
	// TODO: Uncomment after proto regeneration
	t.Skip("Skipping until proto regeneration")
	
	/*
	k, ctx := setupKeeper(t)
	
	valAddr := testValAddr1
	validator := createTestValidator(valAddr, math.NewInt(5000000))
	mockStakingKeeper.GetValidatorFn = func(c context.Context, addr sdk.ValAddress) (stakingtypes.Validator, error) {
		if addr.Equals(valAddr) {
			return validator, nil
		}
		return stakingtypes.Validator{}, errors.New("validator not found")
	}
	
	// Create tokenization records for the validator
	records := []types.TokenizationRecord{
		types.NewTokenizationRecordWithDenom(
			1,
			valAddr.String(),
			testAccAddr1.String(),
			math.NewInt(1000000),
			"liquidstake/val1/1",
		),
		types.NewTokenizationRecordWithDenom(
			2,
			valAddr.String(),
			testAccAddr2.String(),
			math.NewInt(500000),
			"liquidstake/val1/2",
		),
		// Fully redeemed record
		types.NewTokenizationRecordWithDenom(
			3,
			valAddr.String(),
			testAccAddr1.String(),
			math.ZeroInt(),
			"liquidstake/val1/3",
		),
	}
	
	// Store the records
	for _, record := range records {
		k.SetTokenizationRecordWithIndexes(ctx, record)
	}
	
	// Set liquid staked amount for validator
	k.UpdateLiquidStakedAmounts(ctx, valAddr.String(), math.NewInt(1500000), true)
	
	// Set rate limit activity
	activity := keeper.TokenizationActivity{
		TotalAmount:   math.NewInt(100000),
		LastActivity:  ctx.BlockTime(),
		ActivityCount: 3,
	}
	k.SetValidatorTokenizationActivity(ctx, valAddr.String(), activity)
	
	// Query validator statistics
	resp, err := k.ValidatorStatistics(ctx, &types.QueryValidatorStatisticsRequest{
		ValidatorAddress: valAddr.String(),
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	
	// Verify response
	require.Equal(t, valAddr.String(), resp.ValidatorAddress)
	require.Equal(t, math.NewInt(1500000), resp.TotalLiquidStaked)
	require.Equal(t, "30.000000000000000000", resp.LiquidStakingPercentage.String()) // 1.5M / 5M * 100 = 30%
	require.Equal(t, uint64(2), resp.ActiveRecords) // Only records with positive amounts
	require.Equal(t, uint64(3), resp.TotalRecordsCreated)
	
	// Verify rate limit usage
	require.NotNil(t, resp.RateLimitUsage)
	require.Equal(t, "validator", resp.RateLimitUsage.LimitType)
	require.Equal(t, math.NewInt(100000), resp.RateLimitUsage.CurrentAmount)
	require.Equal(t, math.NewInt(500000), resp.RateLimitUsage.MaxAmount) // 10% of 5M
	require.Equal(t, uint64(3), resp.RateLimitUsage.CurrentCount)
	require.Equal(t, uint64(20), resp.RateLimitUsage.MaxCount)
	
	t.Run("validator not found", func(t *testing.T) {
		resp, err := k.ValidatorStatistics(ctx, &types.QueryValidatorStatisticsRequest{
			ValidatorAddress: testValAddr2.String(),
		})
		require.Error(t, err)
		require.Nil(t, resp)
		require.Contains(t, err.Error(), "validator not found")
	})
	
	t.Run("invalid validator address", func(t *testing.T) {
		resp, err := k.ValidatorStatistics(ctx, &types.QueryValidatorStatisticsRequest{
			ValidatorAddress: "invalid",
		})
		require.Error(t, err)
		require.Nil(t, resp)
		require.Contains(t, err.Error(), "invalid validator address")
	})
	
	t.Run("validator with no records", func(t *testing.T) {
		// Create a new validator
		valAddr2 := testValAddr2
		validator2 := createTestValidator(valAddr2, math.NewInt(3000000))
		mockStakingKeeper.GetValidatorFn = func(c context.Context, addr sdk.ValAddress) (stakingtypes.Validator, error) {
			if addr.Equals(valAddr2) {
				return validator2, nil
			}
			return stakingtypes.Validator{}, errors.New("validator not found")
		}
		
		resp, err := k.ValidatorStatistics(ctx, &types.QueryValidatorStatisticsRequest{
			ValidatorAddress: valAddr2.String(),
		})
		require.NoError(t, err)
		require.NotNil(t, resp)
		
		require.Equal(t, valAddr2.String(), resp.ValidatorAddress)
		require.True(t, resp.TotalLiquidStaked.IsZero())
		require.Equal(t, "0.000000000000000000", resp.LiquidStakingPercentage.String())
		require.Equal(t, uint64(0), resp.ActiveRecords)
		require.Equal(t, uint64(0), resp.TotalRecordsCreated)
		
		// Rate limit should show zero usage
		require.NotNil(t, resp.RateLimitUsage)
		require.True(t, resp.RateLimitUsage.CurrentAmount.IsZero())
		require.Equal(t, uint64(0), resp.RateLimitUsage.CurrentCount)
	})
	*/
}