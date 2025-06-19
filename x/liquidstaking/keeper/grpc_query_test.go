package keeper_test

import (
	"context"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/rollchains/flora/x/liquidstaking/types"
)

func (suite *KeeperTestSuite) TestGRPCParams() {
	// Test default params
	resp, err := suite.keeper.Params(context.Background(), &types.QueryParamsRequest{})
	suite.Require().NoError(err)
	suite.Require().NotNil(resp)
	suite.Require().Equal(types.DefaultParams(), resp.Params)

	// Update params
	newParams := types.ModuleParams{
		GlobalLiquidStakingCap: math.LegacyNewDecWithPrec(30, 2), // 30%
		ValidatorLiquidCap:     math.LegacyNewDecWithPrec(60, 2), // 60%
		Enabled:                false,
	}
	suite.keeper.SetParams(suite.ctx, newParams)

	// Query updated params
	resp, err = suite.keeper.Params(context.Background(), &types.QueryParamsRequest{})
	suite.Require().NoError(err)
	suite.Require().Equal(newParams, resp.Params)
}

func (suite *KeeperTestSuite) TestGRPCTokenizationRecord() {
	// Test non-existent record
	_, err := suite.keeper.TokenizationRecord(context.Background(), &types.QueryTokenizationRecordRequest{
		Id: 1,
	})
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "not found")

	// Create a record
	validator := sdk.ValAddress("validator1").String()
	owner := sdk.AccAddress("owner1").String()
	shares := math.NewInt(1000000)
	
	record := types.NewTokenizationRecord(1, validator, owner, shares)
	suite.keeper.SetTokenizationRecord(suite.ctx, record)

	// Query existing record
	resp, err := suite.keeper.TokenizationRecord(context.Background(), &types.QueryTokenizationRecordRequest{
		Id: 1,
	})
	suite.Require().NoError(err)
	suite.Require().Equal(record, resp.Record)

	// Test invalid request
	_, err = suite.keeper.TokenizationRecord(context.Background(), nil)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "invalid request")

	// Test zero ID
	_, err = suite.keeper.TokenizationRecord(context.Background(), &types.QueryTokenizationRecordRequest{
		Id: 0,
	})
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "invalid record ID")
}

func (suite *KeeperTestSuite) TestGRPCTokenizationRecords() {
	// Test empty records
	resp, err := suite.keeper.TokenizationRecords(context.Background(), &types.QueryTokenizationRecordsRequest{})
	suite.Require().NoError(err)
	suite.Require().Empty(resp.Records)

	// Create multiple records
	validator1 := sdk.ValAddress("validator1").String()
	validator2 := sdk.ValAddress("validator2").String()
	owner1 := sdk.AccAddress("owner1").String()
	owner2 := sdk.AccAddress("owner2").String()

	records := []types.TokenizationRecord{
		types.NewTokenizationRecord(1, validator1, owner1, math.NewInt(1000000)),
		types.NewTokenizationRecord(2, validator2, owner1, math.NewInt(2000000)),
		types.NewTokenizationRecord(3, validator1, owner2, math.NewInt(3000000)),
	}

	for _, record := range records {
		suite.keeper.SetTokenizationRecordWithIndexes(suite.ctx, record)
	}

	// Query all records
	resp, err = suite.keeper.TokenizationRecords(context.Background(), &types.QueryTokenizationRecordsRequest{})
	suite.Require().NoError(err)
	suite.Require().Len(resp.Records, 3)
	suite.Require().ElementsMatch(records, resp.Records)

	// Test with pagination
	resp, err = suite.keeper.TokenizationRecords(context.Background(), &types.QueryTokenizationRecordsRequest{
		Pagination: &query.PageRequest{
			Limit: 2,
		},
	})
	suite.Require().NoError(err)
	suite.Require().Len(resp.Records, 2)
	suite.Require().NotNil(resp.Pagination)
	suite.Require().NotEmpty(resp.Pagination.NextKey)

	// Query next page
	resp, err = suite.keeper.TokenizationRecords(context.Background(), &types.QueryTokenizationRecordsRequest{
		Pagination: &query.PageRequest{
			Key: resp.Pagination.NextKey,
		},
	})
	suite.Require().NoError(err)
	suite.Require().Len(resp.Records, 1)
}

func (suite *KeeperTestSuite) TestGRPCTokenizationRecordsByValidator() {
	// Test empty validator
	_, err := suite.keeper.TokenizationRecordsByValidator(context.Background(), &types.QueryTokenizationRecordsByValidatorRequest{
		ValidatorAddress: "",
	})
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "empty validator address")

	validator1 := sdk.ValAddress("validator1").String()
	validator2 := sdk.ValAddress("validator2").String()
	owner1 := sdk.AccAddress("owner1").String()
	owner2 := sdk.AccAddress("owner2").String()

	// Create records
	records := []types.TokenizationRecord{
		types.NewTokenizationRecord(1, validator1, owner1, math.NewInt(1000000)),
		types.NewTokenizationRecord(2, validator1, owner2, math.NewInt(2000000)),
		types.NewTokenizationRecord(3, validator2, owner1, math.NewInt(3000000)),
		types.NewTokenizationRecord(4, validator1, owner1, math.NewInt(4000000)),
	}

	for _, record := range records {
		suite.keeper.SetTokenizationRecord(suite.ctx, record)
		suite.keeper.SetTokenizationRecordWithIndexes(suite.ctx, record)
	}

	// Query records for validator1
	resp, err := suite.keeper.TokenizationRecordsByValidator(context.Background(), &types.QueryTokenizationRecordsByValidatorRequest{
		ValidatorAddress: validator1,
	})
	suite.Require().NoError(err)
	suite.Require().Len(resp.Records, 3)

	// Verify correct records returned
	for _, record := range resp.Records {
		suite.Require().Equal(validator1, record.Validator)
	}

	// Test with pagination
	resp, err = suite.keeper.TokenizationRecordsByValidator(context.Background(), &types.QueryTokenizationRecordsByValidatorRequest{
		ValidatorAddress: validator1,
		Pagination: &query.PageRequest{
			Limit: 2,
		},
	})
	suite.Require().NoError(err)
	suite.Require().Len(resp.Records, 2)
	suite.Require().NotNil(resp.Pagination)
	suite.Require().NotEmpty(resp.Pagination.NextKey)
}

func (suite *KeeperTestSuite) TestGRPCTokenizationRecordsByOwner() {
	// Test empty owner
	_, err := suite.keeper.TokenizationRecordsByOwner(context.Background(), &types.QueryTokenizationRecordsByOwnerRequest{
		OwnerAddress: "",
	})
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "empty owner address")

	validator1 := sdk.ValAddress("validator1").String()
	validator2 := sdk.ValAddress("validator2").String()
	owner1 := sdk.AccAddress("owner1").String()
	owner2 := sdk.AccAddress("owner2").String()

	// Create records
	records := []types.TokenizationRecord{
		types.NewTokenizationRecord(1, validator1, owner1, math.NewInt(1000000)),
		types.NewTokenizationRecord(2, validator2, owner1, math.NewInt(2000000)),
		types.NewTokenizationRecord(3, validator1, owner2, math.NewInt(3000000)),
		types.NewTokenizationRecord(4, validator2, owner1, math.NewInt(4000000)),
	}

	for _, record := range records {
		suite.keeper.SetTokenizationRecord(suite.ctx, record)
	}

	// Query records for owner1
	resp, err := suite.keeper.TokenizationRecordsByOwner(context.Background(), &types.QueryTokenizationRecordsByOwnerRequest{
		OwnerAddress: owner1,
	})
	suite.Require().NoError(err)
	suite.Require().Len(resp.Records, 3)

	// Verify correct records returned
	for _, record := range resp.Records {
		suite.Require().Equal(owner1, record.Owner)
	}

	// Test with pagination
	resp, err = suite.keeper.TokenizationRecordsByOwner(context.Background(), &types.QueryTokenizationRecordsByOwnerRequest{
		OwnerAddress: owner1,
		Pagination: &query.PageRequest{
			Limit: 2,
		},
	})
	suite.Require().NoError(err)
	suite.Require().Len(resp.Records, 2)
	suite.Require().NotNil(resp.Pagination)
	suite.Require().NotEmpty(resp.Pagination.NextKey)
}

func (suite *KeeperTestSuite) TestGRPCTotalLiquidStaked() {
	// Test initial state
	resp, err := suite.keeper.TotalLiquidStaked(context.Background(), &types.QueryTotalLiquidStakedRequest{})
	suite.Require().NoError(err)
	suite.Require().Equal(math.ZeroInt(), resp.TotalLiquidStaked)

	// Set total liquid staked
	total := math.NewInt(5000000)
	suite.keeper.SetTotalLiquidStaked(suite.ctx, total)

	// Query updated total
	resp, err = suite.keeper.TotalLiquidStaked(context.Background(), &types.QueryTotalLiquidStakedRequest{})
	suite.Require().NoError(err)
	suite.Require().Equal(total, resp.TotalLiquidStaked)
}

func (suite *KeeperTestSuite) TestGRPCValidatorLiquidStaked() {
	// Test empty validator
	_, err := suite.keeper.ValidatorLiquidStaked(context.Background(), &types.QueryValidatorLiquidStakedRequest{
		ValidatorAddress: "",
	})
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "empty validator address")

	validator := sdk.ValAddress("validator1").String()

	// Test initial state
	resp, err := suite.keeper.ValidatorLiquidStaked(context.Background(), &types.QueryValidatorLiquidStakedRequest{
		ValidatorAddress: validator,
	})
	suite.Require().NoError(err)
	suite.Require().Equal(math.ZeroInt(), resp.LiquidStaked)

	// Set validator liquid staked
	amount := math.NewInt(3000000)
	suite.keeper.SetValidatorLiquidStaked(suite.ctx, validator, amount)

	// Query updated amount
	resp, err = suite.keeper.ValidatorLiquidStaked(context.Background(), &types.QueryValidatorLiquidStakedRequest{
		ValidatorAddress: validator,
	})
	suite.Require().NoError(err)
	suite.Require().Equal(amount, resp.LiquidStaked)
}

// TODO: Uncomment after proto regeneration
// func (suite *KeeperTestSuite) TestGRPCTokenizationRecordsByDenom() {
// 	// Test empty denom
// 	_, err := suite.keeper.TokenizationRecordsByDenom(context.Background(), &types.QueryTokenizationRecordsByDenomRequest{
// 		Denom: "",
// 	})
// 	suite.Require().Error(err)
// 	suite.Require().Contains(err.Error(), "empty denomination")

// 	// Test non-existent denom
// 	resp, err := suite.keeper.TokenizationRecordsByDenom(context.Background(), &types.QueryTokenizationRecordsByDenomRequest{
// 		Denom: "nonexistent",
// 	})
// 	suite.Require().NoError(err)
// 	suite.Require().Empty(resp.Records)

// 	// Create a record with denom
// 	validator := sdk.ValAddress("validator1").String()
// 	owner := sdk.AccAddress("owner1").String()
// 	shares := math.NewInt(1000000)
// 	denom := "liquidstake/validator1/1"
// 	
// 	record := types.TokenizationRecord{
// 		Id:              1,
// 		Validator:       validator,
// 		Owner:           owner,
// 		SharesTokenized: shares,
// 		Denom:           denom,
// 	}
// 	suite.keeper.SetTokenizationRecord(suite.ctx, record)
// 	suite.keeper.setTokenizationRecordDenomIndex(suite.ctx, denom, record.Id)

// 	// Query by denom
// 	resp, err = suite.keeper.TokenizationRecordsByDenom(context.Background(), &types.QueryTokenizationRecordsByDenomRequest{
// 		Denom: denom,
// 	})
// 	suite.Require().NoError(err)
// 	suite.Require().Len(resp.Records, 1)
// 	suite.Require().Equal(record, resp.Records[0])
// }