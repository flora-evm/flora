package keeper_test

import (
	"cosmossdk.io/math"
	"github.com/rollchains/flora/x/liquidstaking/types"
)

func (suite *KeeperTestSuite) TestTotalLiquidStakedTracking() {
	// Test getting default value (should be zero)
	total := suite.keeper.GetTotalLiquidStaked(suite.ctx)
	suite.True(total.IsZero())

	// Test setting and getting
	amount := math.NewInt(1000000)
	suite.keeper.SetTotalLiquidStaked(suite.ctx, amount)
	total = suite.keeper.GetTotalLiquidStaked(suite.ctx)
	suite.Equal(amount, total)

	// Test increasing
	increaseAmount := math.NewInt(500000)
	suite.keeper.IncreaseTotalLiquidStaked(suite.ctx, increaseAmount)
	total = suite.keeper.GetTotalLiquidStaked(suite.ctx)
	suite.Equal(math.NewInt(1500000), total)

	// Test decreasing
	decreaseAmount := math.NewInt(300000)
	suite.keeper.DecreaseTotalLiquidStaked(suite.ctx, decreaseAmount)
	total = suite.keeper.GetTotalLiquidStaked(suite.ctx)
	suite.Equal(math.NewInt(1200000), total)

	// Test decreasing to zero
	suite.keeper.DecreaseTotalLiquidStaked(suite.ctx, math.NewInt(1200000))
	total = suite.keeper.GetTotalLiquidStaked(suite.ctx)
	suite.True(total.IsZero())

	// Test handling negative (should set to zero)
	suite.keeper.SetTotalLiquidStaked(suite.ctx, math.NewInt(100))
	suite.keeper.DecreaseTotalLiquidStaked(suite.ctx, math.NewInt(200))
	total = suite.keeper.GetTotalLiquidStaked(suite.ctx)
	suite.True(total.IsZero())
}

func (suite *KeeperTestSuite) TestValidatorLiquidStakedTracking() {
	validator1 := "floravaloper1test1"
	validator2 := "floravaloper1test2"

	// Test getting default value (should be zero)
	amount := suite.keeper.GetValidatorLiquidStaked(suite.ctx, validator1)
	suite.True(amount.IsZero())

	// Test setting and getting for different validators
	amount1 := math.NewInt(500000)
	amount2 := math.NewInt(300000)
	suite.keeper.SetValidatorLiquidStaked(suite.ctx, validator1, amount1)
	suite.keeper.SetValidatorLiquidStaked(suite.ctx, validator2, amount2)

	got1 := suite.keeper.GetValidatorLiquidStaked(suite.ctx, validator1)
	got2 := suite.keeper.GetValidatorLiquidStaked(suite.ctx, validator2)
	suite.Equal(amount1, got1)
	suite.Equal(amount2, got2)

	// Test increasing
	suite.keeper.IncreaseValidatorLiquidStaked(suite.ctx, validator1, math.NewInt(100000))
	got1 = suite.keeper.GetValidatorLiquidStaked(suite.ctx, validator1)
	suite.Equal(math.NewInt(600000), got1)

	// Test decreasing
	suite.keeper.DecreaseValidatorLiquidStaked(suite.ctx, validator2, math.NewInt(100000))
	got2 = suite.keeper.GetValidatorLiquidStaked(suite.ctx, validator2)
	suite.Equal(math.NewInt(200000), got2)

	// Test setting to zero (should remove the key)
	suite.keeper.SetValidatorLiquidStaked(suite.ctx, validator1, math.ZeroInt())
	got1 = suite.keeper.GetValidatorLiquidStaked(suite.ctx, validator1)
	suite.True(got1.IsZero())

	// Test GetAllValidatorLiquidStaked
	// Clear existing data first
	all := suite.keeper.GetAllValidatorLiquidStaked(suite.ctx)
	for val := range all {
		suite.keeper.SetValidatorLiquidStaked(suite.ctx, val, math.ZeroInt())
	}

	// Set new data
	validators := map[string]math.Int{
		"floravaloper1aaa": math.NewInt(100000),
		"floravaloper1bbb": math.NewInt(200000),
		"floravaloper1ccc": math.NewInt(300000),
	}

	for val, amt := range validators {
		suite.keeper.SetValidatorLiquidStaked(suite.ctx, val, amt)
	}

	// Get all and verify
	allStaked := suite.keeper.GetAllValidatorLiquidStaked(suite.ctx)
	suite.Len(allStaked, 3)

	for val, expectedAmount := range validators {
		actualAmount, found := allStaked[val]
		suite.True(found)
		suite.Equal(expectedAmount, actualAmount)
	}
}

func (suite *KeeperTestSuite) TestCheckMinimumAmount() {
	// Set minimum amount in params
	params := types.DefaultParams()
	// TODO: Uncomment after proto regeneration
	// params.MinLiquidStakeAmount = math.NewInt(10000)
	suite.NoError(suite.keeper.SetParams(suite.ctx, params))

	// Test amount equal to minimum
	err := suite.keeper.CheckMinimumAmount(suite.ctx, math.NewInt(10000))
	suite.NoError(err)

	// Test amount above minimum
	err = suite.keeper.CheckMinimumAmount(suite.ctx, math.NewInt(50000))
	suite.NoError(err)

	// Test amount below minimum
	err = suite.keeper.CheckMinimumAmount(suite.ctx, math.NewInt(5000))
	suite.Error(err)
	suite.ErrorIs(err, types.ErrAmountTooSmall)

	// Test zero amount
	err = suite.keeper.CheckMinimumAmount(suite.ctx, math.ZeroInt())
	suite.Error(err)
	suite.ErrorIs(err, types.ErrAmountTooSmall)
}

func (suite *KeeperTestSuite) TestUpdateLiquidStakedAmounts() {
	validator := "floravaloper1update"
	
	// Start fresh
	suite.keeper.SetTotalLiquidStaked(suite.ctx, math.ZeroInt())
	suite.keeper.SetValidatorLiquidStaked(suite.ctx, validator, math.ZeroInt())

	// Test increase
	suite.keeper.UpdateLiquidStakedAmounts(suite.ctx, validator, math.NewInt(50000), true)
	suite.Equal(math.NewInt(50000), suite.keeper.GetTotalLiquidStaked(suite.ctx))
	suite.Equal(math.NewInt(50000), suite.keeper.GetValidatorLiquidStaked(suite.ctx, validator))

	// Test another increase
	suite.keeper.UpdateLiquidStakedAmounts(suite.ctx, validator, math.NewInt(30000), true)
	suite.Equal(math.NewInt(80000), suite.keeper.GetTotalLiquidStaked(suite.ctx))
	suite.Equal(math.NewInt(80000), suite.keeper.GetValidatorLiquidStaked(suite.ctx, validator))

	// Test decrease
	suite.keeper.UpdateLiquidStakedAmounts(suite.ctx, validator, math.NewInt(20000), false)
	suite.Equal(math.NewInt(60000), suite.keeper.GetTotalLiquidStaked(suite.ctx))
	suite.Equal(math.NewInt(60000), suite.keeper.GetValidatorLiquidStaked(suite.ctx, validator))
}

func (suite *KeeperTestSuite) TestValidateModuleEnabled() {
	// Test with module enabled
	params := types.DefaultParams()
	params.Enabled = true
	suite.NoError(suite.keeper.SetParams(suite.ctx, params))
	
	err := suite.keeper.ValidateModuleEnabled(suite.ctx)
	suite.NoError(err)

	// Test with module disabled
	params.Enabled = false
	suite.NoError(suite.keeper.SetParams(suite.ctx, params))
	
	err = suite.keeper.ValidateModuleEnabled(suite.ctx)
	suite.Error(err)
	suite.ErrorIs(err, types.ErrDisabled)
}