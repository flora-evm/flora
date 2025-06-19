package keeper_test

import (
	"context"
	"errors"
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"

	"github.com/rollchains/flora/x/liquidstaking/types"
)

func TestGlobalLiquidStakingCap(t *testing.T) {
	keeper, ctx := setupKeeper(t)

	// Set module parameters
	params := types.DefaultParams()
	params.GlobalLiquidStakingCap = math.LegacyNewDecWithPrec(25, 2) // 25%
	// TODO: Uncomment after proto regeneration
	// params.MinLiquidStakeAmount = math.NewInt(1000)
	require.NoError(t, keeper.SetParams(ctx, params))

	// Mock total bonded tokens = 1,000,000
	totalBonded := math.NewInt(1000000)
	mockStakingKeeper.TotalBondedTokensFn = func(c context.Context) (math.Int, error) {
		return totalBonded, nil
	}

	// Cap should be 25% of 1,000,000 = 250,000
	// expectedCap := math.NewInt(250000) // Currently unused

	t.Run("allow tokenization under cap", func(t *testing.T) {
		// Current liquid staked = 100,000
		keeper.SetTotalLiquidStaked(ctx, math.NewInt(100000))

		// Try to tokenize 50,000 more (total would be 150,000 < 250,000)
		err := keeper.CheckGlobalLiquidStakingCap(ctx, math.NewInt(50000))
		require.NoError(t, err)
	})

	t.Run("allow tokenization at exact cap", func(t *testing.T) {
		// Current liquid staked = 200,000
		keeper.SetTotalLiquidStaked(ctx, math.NewInt(200000))

		// Try to tokenize 50,000 more (total would be exactly 250,000)
		err := keeper.CheckGlobalLiquidStakingCap(ctx, math.NewInt(50000))
		require.NoError(t, err)
	})

	t.Run("reject tokenization over cap", func(t *testing.T) {
		// Current liquid staked = 240,000
		keeper.SetTotalLiquidStaked(ctx, math.NewInt(240000))

		// Try to tokenize 20,000 more (total would be 260,000 > 250,000)
		err := keeper.CheckGlobalLiquidStakingCap(ctx, math.NewInt(20000))
		require.Error(t, err)
		require.ErrorIs(t, err, types.ErrGlobalCapExceeded)
		require.Contains(t, err.Error(), "would exceed global liquid staking cap")
	})

	t.Run("handle zero total bonded gracefully", func(t *testing.T) {
		// Mock zero total bonded
		mockStakingKeeper.TotalBondedTokensFn = func(c context.Context) (math.Int, error) {
			return math.ZeroInt(), nil
		}
		keeper.SetTotalLiquidStaked(ctx, math.ZeroInt())

		// Should reject any tokenization when total bonded is zero
		err := keeper.CheckGlobalLiquidStakingCap(ctx, math.NewInt(1000))
		require.Error(t, err)
		require.ErrorIs(t, err, types.ErrGlobalCapExceeded)
	})
}

func TestValidatorLiquidCap(t *testing.T) {
	keeper, ctx := setupKeeper(t)

	// Set module parameters
	params := types.DefaultParams()
	params.GlobalLiquidStakingCap = math.LegacyNewDecWithPrec(10, 2) // 10% - must be <= ValidatorLiquidCap
	params.ValidatorLiquidCap = math.LegacyNewDecWithPrec(10, 2) // 10%
	// TODO: Uncomment after proto regeneration
	// params.MinLiquidStakeAmount = math.NewInt(1000)
	require.NoError(t, keeper.SetParams(ctx, params))

	validatorAddr := testValAddr1.String()
	valAddr := testValAddr1

	// Mock validator with 500,000 tokens
	validator := createTestValidator(valAddr, math.NewInt(500000))
	mockStakingKeeper.GetValidatorFn = func(c context.Context, addr sdk.ValAddress) (stakingtypes.Validator, error) {
		if addr.Equals(valAddr) {
			return validator, nil
		}
		return stakingtypes.Validator{}, errors.New("validator not found")
	}

	// Cap should be 10% of 500,000 = 50,000
	// expectedCap := math.NewInt(50000) // Currently unused

	t.Run("allow tokenization under cap", func(t *testing.T) {
		// Current validator liquid staked = 20,000
		keeper.SetValidatorLiquidStaked(ctx, validatorAddr, math.NewInt(20000))

		// Try to tokenize 20,000 more (total would be 40,000 < 50,000)
		err := keeper.CheckValidatorLiquidCap(ctx, validatorAddr, math.NewInt(20000))
		require.NoError(t, err)
	})

	t.Run("allow tokenization at exact cap", func(t *testing.T) {
		// Current validator liquid staked = 30,000
		keeper.SetValidatorLiquidStaked(ctx, validatorAddr, math.NewInt(30000))

		// Try to tokenize 20,000 more (total would be exactly 50,000)
		err := keeper.CheckValidatorLiquidCap(ctx, validatorAddr, math.NewInt(20000))
		require.NoError(t, err)
	})

	t.Run("reject tokenization over cap", func(t *testing.T) {
		// Current validator liquid staked = 45,000
		keeper.SetValidatorLiquidStaked(ctx, validatorAddr, math.NewInt(45000))

		// Try to tokenize 10,000 more (total would be 55,000 > 50,000)
		err := keeper.CheckValidatorLiquidCap(ctx, validatorAddr, math.NewInt(10000))
		require.Error(t, err)
		require.ErrorIs(t, err, types.ErrValidatorCapExceeded)
		require.Contains(t, err.Error(), "would exceed validator liquid cap")
	})

	t.Run("handle invalid validator address", func(t *testing.T) {
		err := keeper.CheckValidatorLiquidCap(ctx, "invalid", math.NewInt(1000))
		require.Error(t, err)
	})
}

func TestMinimumAmount(t *testing.T) {
	keeper, ctx := setupKeeper(t)

	// Set minimum amount
	params := types.DefaultParams()
	// TODO: Uncomment after proto regeneration
	// params.MinLiquidStakeAmount = math.NewInt(10000)
	require.NoError(t, keeper.SetParams(ctx, params))

	t.Run("accept amount equal to minimum", func(t *testing.T) {
		err := keeper.CheckMinimumAmount(ctx, math.NewInt(10000))
		require.NoError(t, err)
	})

	t.Run("accept amount above minimum", func(t *testing.T) {
		err := keeper.CheckMinimumAmount(ctx, math.NewInt(50000))
		require.NoError(t, err)
	})

	t.Run("reject amount below minimum", func(t *testing.T) {
		err := keeper.CheckMinimumAmount(ctx, math.NewInt(5000))
		require.Error(t, err)
		require.ErrorIs(t, err, types.ErrAmountTooSmall)
		require.Contains(t, err.Error(), "amount 5000 is less than minimum 10000")
	})

	t.Run("reject zero amount", func(t *testing.T) {
		err := keeper.CheckMinimumAmount(ctx, math.ZeroInt())
		require.Error(t, err)
		require.ErrorIs(t, err, types.ErrAmountTooSmall)
	})
}

func TestEnforceTokenizationCaps(t *testing.T) {
	keeper, ctx := setupKeeper(t)

	// Set up parameters
	params := types.DefaultParams()
	params.Enabled = true
	// TODO: Uncomment after proto regeneration
	// params.MinLiquidStakeAmount = math.NewInt(5000)
	params.GlobalLiquidStakingCap = math.LegacyNewDecWithPrec(15, 2) // 15% - must be <= ValidatorLiquidCap
	params.ValidatorLiquidCap = math.LegacyNewDecWithPrec(30, 2)     // 30%
	require.NoError(t, keeper.SetParams(ctx, params))

	validatorAddr := testValAddr1.String()
	valAddr := testValAddr1

	// Mock validator with 200,000 tokens
	validator := createTestValidator(valAddr, math.NewInt(200000))
	mockStakingKeeper.GetValidatorFn = func(c context.Context, addr sdk.ValAddress) (stakingtypes.Validator, error) {
		if addr.Equals(valAddr) {
			return validator, nil
		}
		return stakingtypes.Validator{}, errors.New("validator not found")
	}

	// Mock total bonded = 1,000,000
	mockStakingKeeper.TotalBondedTokensFn = func(c context.Context) (math.Int, error) {
		return math.NewInt(1000000), nil
	}

	// Set current amounts
	keeper.SetTotalLiquidStaked(ctx, math.NewInt(100000))       // 10% of total
	keeper.SetValidatorLiquidStaked(ctx, validatorAddr, math.NewInt(10000)) // 5% of validator

	t.Run("pass all checks", func(t *testing.T) {
		// Amount is above minimum, within both caps
		err := keeper.EnforceTokenizationCaps(ctx, validatorAddr, math.NewInt(10000))
		require.NoError(t, err)
	})

	t.Run("fail minimum amount check", func(t *testing.T) {
		err := keeper.EnforceTokenizationCaps(ctx, validatorAddr, math.NewInt(1000))
		require.Error(t, err)
		require.ErrorIs(t, err, types.ErrAmountTooSmall)
	})

	t.Run("fail global cap check", func(t *testing.T) {
		// Try to tokenize 50,000 (would make total 150,000 = 15% which is the global cap)
		// But we need to exceed it, so use 50,001
		err := keeper.EnforceTokenizationCaps(ctx, validatorAddr, math.NewInt(50001))
		require.Error(t, err)
		require.ErrorIs(t, err, types.ErrGlobalCapExceeded)
	})

	t.Run("fail validator cap check", func(t *testing.T) {
		// Validator has 200,000 tokens, 30% cap = 60,000
		// Current validator liquid staked = 10,000
		// Try to tokenize 40,000 (would make validator total 50,000 < 60,000 cap)
		// But would still be within global cap (140,000 < 150,000)
		// We need an amount that passes global cap but fails validator cap
		// Let's tokenize 45,000 which keeps global at 145,000 < 150,000
		// But makes validator total 55,000 which is still < 60,000...
		// Actually we need to set the current validator amount higher
		keeper.SetValidatorLiquidStaked(ctx, validatorAddr, math.NewInt(50000)) // Now at 25% of validator
		
		// Try to tokenize 10,001 (would make validator total 60,001 > 60,000 cap)
		// Global would be 110,001 which is still < 150,000 cap
		err := keeper.EnforceTokenizationCaps(ctx, validatorAddr, math.NewInt(10001))
		require.Error(t, err)
		require.ErrorIs(t, err, types.ErrValidatorCapExceeded)
	})
}

func TestLiquidStakedAmountTracking(t *testing.T) {
	keeper, ctx := setupKeeper(t)

	t.Run("increase and decrease total liquid staked", func(t *testing.T) {
		// Start with zero
		require.True(t, keeper.GetTotalLiquidStaked(ctx).IsZero())

		// Increase by 100,000
		keeper.IncreaseTotalLiquidStaked(ctx, math.NewInt(100000))
		require.Equal(t, math.NewInt(100000), keeper.GetTotalLiquidStaked(ctx))

		// Increase by another 50,000
		keeper.IncreaseTotalLiquidStaked(ctx, math.NewInt(50000))
		require.Equal(t, math.NewInt(150000), keeper.GetTotalLiquidStaked(ctx))

		// Decrease by 30,000
		keeper.DecreaseTotalLiquidStaked(ctx, math.NewInt(30000))
		require.Equal(t, math.NewInt(120000), keeper.GetTotalLiquidStaked(ctx))

		// Decrease to zero
		keeper.DecreaseTotalLiquidStaked(ctx, math.NewInt(120000))
		require.True(t, keeper.GetTotalLiquidStaked(ctx).IsZero())
	})

	t.Run("handle negative total gracefully", func(t *testing.T) {
		keeper.SetTotalLiquidStaked(ctx, math.NewInt(10000))
		
		// Try to decrease by more than current amount
		keeper.DecreaseTotalLiquidStaked(ctx, math.NewInt(20000))
		
		// Should be set to zero, not negative
		require.True(t, keeper.GetTotalLiquidStaked(ctx).IsZero())
	})

	t.Run("validator liquid staked tracking", func(t *testing.T) {
		validator1 := testValAddr1.String()
		validator2 := testValAddr2.String()

		// Start with zero
		require.True(t, keeper.GetValidatorLiquidStaked(ctx, validator1).IsZero())
		require.True(t, keeper.GetValidatorLiquidStaked(ctx, validator2).IsZero())

		// Set amounts for different validators
		keeper.SetValidatorLiquidStaked(ctx, validator1, math.NewInt(50000))
		keeper.SetValidatorLiquidStaked(ctx, validator2, math.NewInt(30000))

		require.Equal(t, math.NewInt(50000), keeper.GetValidatorLiquidStaked(ctx, validator1))
		require.Equal(t, math.NewInt(30000), keeper.GetValidatorLiquidStaked(ctx, validator2))

		// Increase validator1
		keeper.IncreaseValidatorLiquidStaked(ctx, validator1, math.NewInt(20000))
		require.Equal(t, math.NewInt(70000), keeper.GetValidatorLiquidStaked(ctx, validator1))

		// Decrease validator2
		keeper.DecreaseValidatorLiquidStaked(ctx, validator2, math.NewInt(10000))
		require.Equal(t, math.NewInt(20000), keeper.GetValidatorLiquidStaked(ctx, validator2))

		// Set to zero should remove the key
		keeper.SetValidatorLiquidStaked(ctx, validator2, math.ZeroInt())
		require.True(t, keeper.GetValidatorLiquidStaked(ctx, validator2).IsZero())
	})

	t.Run("get all validator liquid staked", func(t *testing.T) {
		// Clear any existing data
		all := keeper.GetAllValidatorLiquidStaked(ctx)
		for val := range all {
			keeper.SetValidatorLiquidStaked(ctx, val, math.ZeroInt())
		}

		// Set amounts for multiple validators
		validators := map[string]math.Int{
			testValAddr1.String(): math.NewInt(100000),
			testValAddr2.String(): math.NewInt(200000),
			sdk.ValAddress(ed25519.GenPrivKey().PubKey().Address()).String(): math.NewInt(300000),
		}

		for val, amount := range validators {
			keeper.SetValidatorLiquidStaked(ctx, val, amount)
		}

		// Get all and verify
		allStaked := keeper.GetAllValidatorLiquidStaked(ctx)
		require.Len(t, allStaked, 3)

		for val, expectedAmount := range validators {
			actualAmount, found := allStaked[val]
			require.True(t, found)
			require.Equal(t, expectedAmount, actualAmount)
		}
	})

	t.Run("update liquid staked amounts helper", func(t *testing.T) {
		validator := testValAddr1.String()
		
		// Start fresh
		keeper.SetTotalLiquidStaked(ctx, math.ZeroInt())
		keeper.SetValidatorLiquidStaked(ctx, validator, math.ZeroInt())

		// Increase
		keeper.UpdateLiquidStakedAmounts(ctx, validator, math.NewInt(50000), true)
		require.Equal(t, math.NewInt(50000), keeper.GetTotalLiquidStaked(ctx))
		require.Equal(t, math.NewInt(50000), keeper.GetValidatorLiquidStaked(ctx, validator))

		// Increase again
		keeper.UpdateLiquidStakedAmounts(ctx, validator, math.NewInt(30000), true)
		require.Equal(t, math.NewInt(80000), keeper.GetTotalLiquidStaked(ctx))
		require.Equal(t, math.NewInt(80000), keeper.GetValidatorLiquidStaked(ctx, validator))

		// Decrease
		keeper.UpdateLiquidStakedAmounts(ctx, validator, math.NewInt(20000), false)
		require.Equal(t, math.NewInt(60000), keeper.GetTotalLiquidStaked(ctx))
		require.Equal(t, math.NewInt(60000), keeper.GetValidatorLiquidStaked(ctx, validator))
	})
}