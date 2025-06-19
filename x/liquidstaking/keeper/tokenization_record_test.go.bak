package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	
	"github.com/rollchains/flora/x/liquidstaking/types"
)

func (suite *KeeperTestSuite) TestTokenizationRecordOperations() {
	validatorAddr := types.TestValidatorAddr
	ownerAddr := types.TestOwnerAddr
	
	// Test SetTokenizationRecordWithIndexes
	record1 := types.NewTokenizationRecord(1, validatorAddr, ownerAddr, math.NewInt(1000))
	suite.keeper.SetTokenizationRecordWithIndexes(suite.ctx, record1)
	
	// Verify record was stored
	gotRecord, found := suite.keeper.GetTokenizationRecord(suite.ctx, 1)
	suite.True(found)
	suite.Equal(record1, gotRecord)
	
	// Add more records for the same validator and owner
	record2 := types.NewTokenizationRecord(2, validatorAddr, ownerAddr, math.NewInt(2000))
	record3 := types.NewTokenizationRecord(3, validatorAddr, "flora1differentowner1234567890abcdefghijklmn", math.NewInt(3000))
	
	suite.keeper.SetTokenizationRecordWithIndexes(suite.ctx, record2)
	suite.keeper.SetTokenizationRecordWithIndexes(suite.ctx, record3)
	
	// Test GetTokenizationRecordsByValidator
	validatorRecords := suite.keeper.GetTokenizationRecordsByValidator(suite.ctx, validatorAddr)
	suite.Len(validatorRecords, 3)
	suite.Contains(validatorRecords, record1)
	suite.Contains(validatorRecords, record2)
	suite.Contains(validatorRecords, record3)
	
	// Test GetTokenizationRecordsByOwner
	ownerRecords := suite.keeper.GetTokenizationRecordsByOwner(suite.ctx, ownerAddr)
	suite.Len(ownerRecords, 2)
	suite.Contains(ownerRecords, record1)
	suite.Contains(ownerRecords, record2)
	
	// Test DeleteTokenizationRecord
	err := suite.keeper.DeleteTokenizationRecord(suite.ctx, 1)
	suite.NoError(err)
	
	// Verify record was deleted
	_, found = suite.keeper.GetTokenizationRecord(suite.ctx, 1)
	suite.False(found)
	
	// Verify indexes were updated
	validatorRecords = suite.keeper.GetTokenizationRecordsByValidator(suite.ctx, validatorAddr)
	suite.Len(validatorRecords, 2)
	suite.NotContains(validatorRecords, record1)
}

func (suite *KeeperTestSuite) TestLiquidStakedAmounts() {
	validatorAddr1 := types.TestValidatorAddr
	// Generate a second test validator address
	addr := sdk.AccAddress([]byte("test2"))
	valAddr := sdk.ValAddress(addr)
	validatorAddr2 := valAddr.String()
	
	// Test initial state
	total := suite.keeper.GetTotalLiquidStaked(suite.ctx)
	suite.Equal(math.ZeroInt(), total)
	
	val1Amount := suite.keeper.GetValidatorLiquidStaked(suite.ctx, validatorAddr1)
	suite.Equal(math.ZeroInt(), val1Amount)
	
	// Test setting amounts
	suite.keeper.SetTotalLiquidStaked(suite.ctx, math.NewInt(5000))
	suite.keeper.SetValidatorLiquidStaked(suite.ctx, validatorAddr1, math.NewInt(3000))
	suite.keeper.SetValidatorLiquidStaked(suite.ctx, validatorAddr2, math.NewInt(2000))
	
	// Verify amounts
	total = suite.keeper.GetTotalLiquidStaked(suite.ctx)
	suite.Equal(math.NewInt(5000), total)
	
	val1Amount = suite.keeper.GetValidatorLiquidStaked(suite.ctx, validatorAddr1)
	suite.Equal(math.NewInt(3000), val1Amount)
	
	val2Amount := suite.keeper.GetValidatorLiquidStaked(suite.ctx, validatorAddr2)
	suite.Equal(math.NewInt(2000), val2Amount)
	
	// Test UpdateLiquidStakedAmounts - increase
	suite.keeper.UpdateLiquidStakedAmounts(suite.ctx, validatorAddr1, math.NewInt(1000), true)
	
	total = suite.keeper.GetTotalLiquidStaked(suite.ctx)
	suite.Equal(math.NewInt(6000), total)
	
	val1Amount = suite.keeper.GetValidatorLiquidStaked(suite.ctx, validatorAddr1)
	suite.Equal(math.NewInt(4000), val1Amount)
	
	// Test UpdateLiquidStakedAmounts - decrease
	suite.keeper.UpdateLiquidStakedAmounts(suite.ctx, validatorAddr2, math.NewInt(500), false)
	
	total = suite.keeper.GetTotalLiquidStaked(suite.ctx)
	suite.Equal(math.NewInt(5500), total)
	
	val2Amount = suite.keeper.GetValidatorLiquidStaked(suite.ctx, validatorAddr2)
	suite.Equal(math.NewInt(1500), val2Amount)
	
	// Test negative protection
	suite.keeper.UpdateLiquidStakedAmounts(suite.ctx, validatorAddr2, math.NewInt(2000), false)
	
	val2Amount = suite.keeper.GetValidatorLiquidStaked(suite.ctx, validatorAddr2)
	suite.Equal(math.ZeroInt(), val2Amount)
}

func (suite *KeeperTestSuite) TestGetTokenizationRecordByDenom() {
	// For Stage 2, we'll test the denom index functionality
	// The actual denom will be set in Stage 3 when minting is implemented
	
	// Test record not found by denom
	_, found := suite.keeper.GetTokenizationRecordByDenom(suite.ctx, "liquidstake/1")
	suite.False(found)
	
	// This functionality will be fully tested in Stage 3
}