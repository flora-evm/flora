package keeper_test

import (
	"cosmossdk.io/math"
	
	"github.com/rollchains/flora/x/liquidstaking/types"
)

func (suite *KeeperTestSuite) TestValidateGlobalLiquidStakingCap() {
	// Set up params with 25% global cap
	params := types.NewParams(
		math.LegacyNewDecWithPrec(25, 2), // 25%
		math.LegacyNewDecWithPrec(50, 2), // 50%
		true,
	)
	err := suite.keeper.SetParams(suite.ctx, params)
	suite.NoError(err)
	
	// With placeholder total bonded of 1 billion, max allowed is 250 million
	
	// Test: no existing liquid staked, within cap
	err = suite.keeper.ValidateGlobalLiquidStakingCap(suite.ctx, math.NewInt(100_000_000))
	suite.NoError(err)
	
	// Test: set existing liquid staked
	suite.keeper.SetTotalLiquidStaked(suite.ctx, math.NewInt(200_000_000))
	
	// Test: would be exactly at cap
	err = suite.keeper.ValidateGlobalLiquidStakingCap(suite.ctx, math.NewInt(50_000_000))
	suite.NoError(err)
	
	// Test: would exceed cap
	err = suite.keeper.ValidateGlobalLiquidStakingCap(suite.ctx, math.NewInt(51_000_000))
	suite.Error(err)
	suite.Contains(err.Error(), "would exceed global liquid staking cap")
	
	// Test: module disabled
	params.Enabled = false
	err = suite.keeper.SetParams(suite.ctx, params)
	suite.NoError(err)
	
	err = suite.keeper.ValidateGlobalLiquidStakingCap(suite.ctx, math.NewInt(1000))
	suite.Error(err)
	suite.Equal(types.ErrModuleDisabled, err)
}

func (suite *KeeperTestSuite) TestValidateValidatorLiquidCap() {
	validatorAddr := types.TestValidatorAddr
	
	// Set up params with 50% validator cap
	params := types.NewParams(
		math.LegacyNewDecWithPrec(25, 2), // 25%
		math.LegacyNewDecWithPrec(50, 2), // 50%
		true,
	)
	err := suite.keeper.SetParams(suite.ctx, params)
	suite.NoError(err)
	
	// With placeholder validator shares of 100 million, max allowed is 50 million
	
	// Test: no existing liquid staked, within cap
	err = suite.keeper.ValidateValidatorLiquidCap(suite.ctx, validatorAddr, math.NewInt(30_000_000))
	suite.NoError(err)
	
	// Test: set existing liquid staked for validator
	suite.keeper.SetValidatorLiquidStaked(suite.ctx, validatorAddr, math.NewInt(40_000_000))
	
	// Test: would be exactly at cap
	err = suite.keeper.ValidateValidatorLiquidCap(suite.ctx, validatorAddr, math.NewInt(10_000_000))
	suite.NoError(err)
	
	// Test: would exceed cap
	err = suite.keeper.ValidateValidatorLiquidCap(suite.ctx, validatorAddr, math.NewInt(11_000_000))
	suite.Error(err)
	suite.Contains(err.Error(), "would exceed validator liquid cap")
	
	// Test: module disabled
	params.Enabled = false
	err = suite.keeper.SetParams(suite.ctx, params)
	suite.NoError(err)
	
	err = suite.keeper.ValidateValidatorLiquidCap(suite.ctx, validatorAddr, math.NewInt(1000))
	suite.Error(err)
	suite.Equal(types.ErrModuleDisabled, err)
}

func (suite *KeeperTestSuite) TestCanTokenizeShares() {
	validatorAddr := types.TestValidatorAddr
	
	// Set up params
	params := types.NewParams(
		math.LegacyNewDecWithPrec(25, 2), // 25%
		math.LegacyNewDecWithPrec(50, 2), // 50%
		true,
	)
	err := suite.keeper.SetParams(suite.ctx, params)
	suite.NoError(err)
	
	// Test: valid tokenization
	err = suite.keeper.CanTokenizeShares(suite.ctx, validatorAddr, math.NewInt(1_000_000))
	suite.NoError(err)
	
	// Test: zero shares
	err = suite.keeper.CanTokenizeShares(suite.ctx, validatorAddr, math.ZeroInt())
	suite.Error(err)
	suite.Contains(err.Error(), "shares must be positive")
	
	// Test: negative shares
	err = suite.keeper.CanTokenizeShares(suite.ctx, validatorAddr, math.NewInt(-1000))
	suite.Error(err)
	suite.Contains(err.Error(), "shares must be positive")
	
	// Test: module disabled
	params.Enabled = false
	err = suite.keeper.SetParams(suite.ctx, params)
	suite.NoError(err)
	
	err = suite.keeper.CanTokenizeShares(suite.ctx, validatorAddr, math.NewInt(1000))
	suite.Error(err)
	suite.Equal(types.ErrModuleDisabled, err)
	
	// Test: exceeds global cap
	params.Enabled = true
	err = suite.keeper.SetParams(suite.ctx, params)
	suite.NoError(err)
	
	suite.keeper.SetTotalLiquidStaked(suite.ctx, math.NewInt(240_000_000))
	err = suite.keeper.CanTokenizeShares(suite.ctx, validatorAddr, math.NewInt(20_000_000))
	suite.Error(err)
	suite.Contains(err.Error(), "would exceed global liquid staking cap")
	
	// Test: exceeds validator cap
	suite.keeper.SetTotalLiquidStaked(suite.ctx, math.ZeroInt())
	suite.keeper.SetValidatorLiquidStaked(suite.ctx, validatorAddr, math.NewInt(45_000_000))
	err = suite.keeper.CanTokenizeShares(suite.ctx, validatorAddr, math.NewInt(10_000_000))
	suite.Error(err)
	suite.Contains(err.Error(), "would exceed validator liquid cap")
}

func (suite *KeeperTestSuite) TestValidateTokenizationRecord() {
	validatorAddr := types.TestValidatorAddr
	ownerAddr := types.TestOwnerAddr
	
	// Test: valid record
	record := types.NewTokenizationRecord(1, validatorAddr, ownerAddr, math.NewInt(1000))
	err := suite.keeper.ValidateTokenizationRecord(suite.ctx, record)
	suite.NoError(err)
	
	// Test: invalid record (zero ID)
	invalidRecord := types.NewTokenizationRecord(0, validatorAddr, ownerAddr, math.NewInt(1000))
	err = suite.keeper.ValidateTokenizationRecord(suite.ctx, invalidRecord)
	suite.Error(err)
	suite.Contains(err.Error(), "tokenization record id cannot be zero")
	
	// Test: duplicate ID
	suite.keeper.SetTokenizationRecordWithIndexes(suite.ctx, record)
	err = suite.keeper.ValidateTokenizationRecord(suite.ctx, record)
	suite.Error(err)
	suite.Contains(err.Error(), "already exists")
}