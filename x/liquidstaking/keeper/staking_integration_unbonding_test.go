package keeper_test

import (
	"context"
	"errors"
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/suite"

	"github.com/rollchains/flora/x/liquidstaking/types"
)

type UnbondingTestSuite struct {
	StakingIntegrationTestSuite
}

func TestUnbondingTestSuite(t *testing.T) {
	suite.Run(t, new(UnbondingTestSuite))
}

// Helper method to set tokenization record with all indexes including denom
func (suite *UnbondingTestSuite) setTokenizationRecordIndexes(ctx sdk.Context, record types.TokenizationRecord) {
	store := suite.keeper.GetStoreService().OpenKVStore(ctx)
	
	// Set validator index
	validatorKey := types.GetTokenizationRecordByValidatorKey(record.Validator, record.Id)
	err := store.Set(validatorKey, []byte{})
	if err != nil {
		panic(err)
	}
	
	// Set owner index
	ownerKey := types.GetTokenizationRecordByOwnerKey(record.Owner, record.Id)
	err = store.Set(ownerKey, []byte{})
	if err != nil {
		panic(err)
	}
	
	// Set denom index - this is normally done during minting but we need it for tests
	denomKey := types.GetTokenizationRecordByDenomKey(record.Denom)
	err = store.Set(denomKey, types.Uint64ToBytes(record.Id))
	if err != nil {
		panic(err)
	}
}

// Test tokenization with unbonding validator
func (suite *UnbondingTestSuite) TestTokenizeShares_UnbondingValidator() {
	// Setup addresses
	delegatorAddr := sdk.AccAddress([]byte("delegator"))
	validatorAddr := sdk.ValAddress([]byte("validator"))
	
	// Setup mock expectations
	delegation := stakingtypes.Delegation{
		DelegatorAddress: delegatorAddr.String(),
		ValidatorAddress: validatorAddr.String(),
		Shares:           math.LegacyNewDec(1000000),
	}
	
	// Validator in unbonding state
	validator := stakingtypes.Validator{
		OperatorAddress: validatorAddr.String(),
		Status:          stakingtypes.Unbonding, // Unbonding state
		Tokens:          math.NewInt(1000000),
		DelegatorShares: math.LegacyNewDec(1000000),
		Jailed:          false,
	}
	
	suite.stakingKeeper.GetDelegationFn = func(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (stakingtypes.Delegation, error) {
		return delegation, nil
	}
	
	suite.stakingKeeper.GetValidatorFn = func(ctx context.Context, addr sdk.ValAddress) (stakingtypes.Validator, error) {
		return validator, nil
	}
	
	suite.stakingKeeper.UnbondFn = func(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, shares math.LegacyDec) (math.Int, error) {
		// Should still work for unbonding validators
		return shares.TruncateInt(), nil
	}
	
	suite.bankKeeper.MintCoinsFn = func(ctx context.Context, moduleName string, amt sdk.Coins) error {
		return nil
	}
	
	suite.bankKeeper.SendCoinsFromModuleToAccountFn = func(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error {
		return nil
	}
	
	suite.bankKeeper.SetDenomMetaDataFn = func(ctx context.Context, denomMetaData banktypes.Metadata) {}
	
	// Mock total bonded tokens to avoid cap issues
	suite.stakingKeeper.TotalBondedTokensFn = func(ctx context.Context) math.Int {
		return math.NewInt(100000000) // 100M total bonded
	}
	
	// Execute tokenization
	msg := &types.MsgTokenizeShares{
		DelegatorAddress: delegatorAddr.String(),
		ValidatorAddress: validatorAddr.String(),
		Shares:           sdk.NewCoin("shares", math.NewInt(500000)),
		OwnerAddress:     "",
	}
	
	// Should succeed - unbonding validators can still be tokenized
	resp, err := suite.msgServer.TokenizeShares(suite.ctx, msg)
	suite.Require().NoError(err)
	suite.Require().NotNil(resp)
}

// Test tokenization with unbonded validator
func (suite *UnbondingTestSuite) TestTokenizeShares_UnbondedValidator() {
	// Setup addresses
	delegatorAddr := sdk.AccAddress([]byte("delegator"))
	validatorAddr := sdk.ValAddress([]byte("validator"))
	
	// Setup mock expectations
	delegation := stakingtypes.Delegation{
		DelegatorAddress: delegatorAddr.String(),
		ValidatorAddress: validatorAddr.String(),
		Shares:           math.LegacyNewDec(1000000),
	}
	
	// Validator in unbonded state
	validator := stakingtypes.Validator{
		OperatorAddress: validatorAddr.String(),
		Status:          stakingtypes.Unbonded, // Fully unbonded
		Tokens:          math.NewInt(1000000),
		DelegatorShares: math.LegacyNewDec(1000000),
		Jailed:          false,
	}
	
	suite.stakingKeeper.GetDelegationFn = func(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (stakingtypes.Delegation, error) {
		return delegation, nil
	}
	
	suite.stakingKeeper.GetValidatorFn = func(ctx context.Context, addr sdk.ValAddress) (stakingtypes.Validator, error) {
		return validator, nil
	}
	
	suite.stakingKeeper.UnbondFn = func(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, shares math.LegacyDec) (math.Int, error) {
		// Should still work for unbonded validators if delegation exists
		return shares.TruncateInt(), nil
	}
	
	suite.bankKeeper.MintCoinsFn = func(ctx context.Context, moduleName string, amt sdk.Coins) error {
		return nil
	}
	
	suite.bankKeeper.SendCoinsFromModuleToAccountFn = func(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error {
		return nil
	}
	
	suite.bankKeeper.SetDenomMetaDataFn = func(ctx context.Context, denomMetaData banktypes.Metadata) {}
	
	// Mock total bonded tokens to avoid cap issues
	suite.stakingKeeper.TotalBondedTokensFn = func(ctx context.Context) math.Int {
		return math.NewInt(100000000) // 100M total bonded
	}
	
	// Execute tokenization
	msg := &types.MsgTokenizeShares{
		DelegatorAddress: delegatorAddr.String(),
		ValidatorAddress: validatorAddr.String(),
		Shares:           sdk.NewCoin("shares", math.NewInt(500000)),
		OwnerAddress:     "",
	}
	
	// Should succeed - unbonded validators can still be tokenized if delegation exists
	resp, err := suite.msgServer.TokenizeShares(suite.ctx, msg)
	suite.Require().NoError(err)
	suite.Require().NotNil(resp)
}

// Test unbond failure during tokenization
func (suite *UnbondingTestSuite) TestTokenizeShares_UnbondFailure() {
	// Setup addresses
	delegatorAddr := sdk.AccAddress([]byte("delegator"))
	validatorAddr := sdk.ValAddress([]byte("validator"))
	
	// Setup mock expectations
	delegation := stakingtypes.Delegation{
		DelegatorAddress: delegatorAddr.String(),
		ValidatorAddress: validatorAddr.String(),
		Shares:           math.LegacyNewDec(1000000),
	}
	
	validator := stakingtypes.Validator{
		OperatorAddress: validatorAddr.String(),
		Status:          stakingtypes.Bonded,
		Tokens:          math.NewInt(1000000),
		DelegatorShares: math.LegacyNewDec(1000000),
		Jailed:          false,
	}
	
	suite.stakingKeeper.GetDelegationFn = func(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (stakingtypes.Delegation, error) {
		return delegation, nil
	}
	
	suite.stakingKeeper.GetValidatorFn = func(ctx context.Context, addr sdk.ValAddress) (stakingtypes.Validator, error) {
		return validator, nil
	}
	
	// Simulate unbond failure
	suite.stakingKeeper.UnbondFn = func(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, shares math.LegacyDec) (math.Int, error) {
		return math.ZeroInt(), errors.New("unbonding failed: maximum entries reached")
	}
	
	// Mock total bonded tokens to avoid cap issues
	suite.stakingKeeper.TotalBondedTokensFn = func(ctx context.Context) math.Int {
		return math.NewInt(100000000) // 100M total bonded
	}
	
	// Execute tokenization
	msg := &types.MsgTokenizeShares{
		DelegatorAddress: delegatorAddr.String(),
		ValidatorAddress: validatorAddr.String(),
		Shares:           sdk.NewCoin("shares", math.NewInt(500000)),
		OwnerAddress:     "",
	}
	
	_, err := suite.msgServer.TokenizeShares(suite.ctx, msg)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "failed to unbond shares")
}

// Test redemption with delegation failure
func (suite *UnbondingTestSuite) TestRedeemTokens_DelegationFailure() {
	// First, create a tokenization record
	delegatorAddr := sdk.AccAddress([]byte("delegator"))
	validatorAddr := sdk.ValAddress([]byte("validator"))
	
	record := types.NewTokenizationRecordWithDenom(
		1,
		validatorAddr.String(),
		delegatorAddr.String(),
		math.NewInt(500000),
		"flora/lstake/"+validatorAddr.String()+"/1",
	)
	// For testing, we need to manually create the full record with denom index
	// This simulates what would happen during actual tokenization
	suite.keeper.SetTokenizationRecord(suite.ctx, record)
	suite.setTokenizationRecordIndexes(suite.ctx, record)
	
	// Setup mock expectations
	validator := stakingtypes.Validator{
		OperatorAddress: validatorAddr.String(),
		Status:          stakingtypes.Bonded,
		Tokens:          math.NewInt(1000000),
		DelegatorShares: math.LegacyNewDec(1000000),
		Jailed:          false,
	}
	
	suite.stakingKeeper.GetValidatorFn = func(ctx context.Context, addr sdk.ValAddress) (stakingtypes.Validator, error) {
		return validator, nil
	}
	
	// Mock bank keeper for redemption
	suite.bankKeeper.GetBalanceFn = func(ctx context.Context, addr sdk.AccAddress, denom string) sdk.Coin {
		return sdk.NewCoin(denom, math.NewInt(500000))
	}
	
	suite.bankKeeper.SendCoinsFromAccountToModuleFn = func(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error {
		return nil
	}
	
	suite.bankKeeper.BurnCoinsFn = func(ctx context.Context, moduleName string, amt sdk.Coins) error {
		return nil
	}
	
	// Simulate delegation failure
	suite.stakingKeeper.DelegateFn = func(ctx context.Context, delAddr sdk.AccAddress, bondAmt math.Int, tokenSrc stakingtypes.BondStatus, validator stakingtypes.Validator, subtractAccount bool) (math.LegacyDec, error) {
		return math.LegacyZeroDec(), errors.New("delegation failed: validator at max delegations")
	}
	
	// Execute redemption
	msg := &types.MsgRedeemTokens{
		OwnerAddress: delegatorAddr.String(),
		Amount:       sdk.NewCoin("flora/lstake/"+validatorAddr.String()+"/1", math.NewInt(250000)),
	}
	
	_, err := suite.msgServer.RedeemTokens(suite.ctx, msg)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "failed to delegate")
}

// Test redemption with partial amount
func (suite *UnbondingTestSuite) TestRedeemTokens_PartialRedemption() {
	// First, create a tokenization record
	delegatorAddr := sdk.AccAddress([]byte("delegator"))
	validatorAddr := sdk.ValAddress([]byte("validator"))
	
	record := types.NewTokenizationRecordWithDenom(
		1,
		validatorAddr.String(),
		delegatorAddr.String(),
		math.NewInt(1000000),
		"flora/lstake/"+validatorAddr.String()+"/1",
	)
	// For testing, we need to manually create the full record with denom index
	// This simulates what would happen during actual tokenization
	suite.keeper.SetTokenizationRecord(suite.ctx, record)
	suite.setTokenizationRecordIndexes(suite.ctx, record)
	
	// Setup mock expectations
	validator := stakingtypes.Validator{
		OperatorAddress: validatorAddr.String(),
		Status:          stakingtypes.Bonded,
		Tokens:          math.NewInt(1000000),
		DelegatorShares: math.LegacyNewDec(1000000),
		Jailed:          false,
	}
	
	suite.stakingKeeper.GetValidatorFn = func(ctx context.Context, addr sdk.ValAddress) (stakingtypes.Validator, error) {
		return validator, nil
	}
	
	// Mock bank keeper for redemption
	suite.bankKeeper.GetBalanceFn = func(ctx context.Context, addr sdk.AccAddress, denom string) sdk.Coin {
		return sdk.NewCoin(denom, math.NewInt(1000000))
	}
	
	suite.bankKeeper.SendCoinsFromAccountToModuleFn = func(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error {
		return nil
	}
	
	suite.bankKeeper.BurnCoinsFn = func(ctx context.Context, moduleName string, amt sdk.Coins) error {
		return nil
	}
	
	suite.stakingKeeper.DelegateFn = func(ctx context.Context, delAddr sdk.AccAddress, bondAmt math.Int, tokenSrc stakingtypes.BondStatus, validator stakingtypes.Validator, subtractAccount bool) (math.LegacyDec, error) {
		return math.LegacyNewDecFromInt(bondAmt), nil
	}
	
	// Execute partial redemption
	msg := &types.MsgRedeemTokens{
		OwnerAddress: delegatorAddr.String(),
		Amount:       sdk.NewCoin("flora/lstake/"+validatorAddr.String()+"/1", math.NewInt(250000)), // Redeem only 25%
	}
	
	resp, err := suite.msgServer.RedeemTokens(suite.ctx, msg)
	suite.Require().NoError(err)
	suite.Require().NotNil(resp)
	suite.Require().Equal(uint64(1), resp.RecordId)
	
	// Verify record was updated, not deleted
	updatedRecord, found := suite.keeper.GetTokenizationRecord(suite.ctx, 1)
	suite.Require().True(found)
	suite.Require().Equal(math.NewInt(750000), updatedRecord.SharesTokenized) // 1M - 250k = 750k
}

// Test full redemption (record should be deleted)
func (suite *UnbondingTestSuite) TestRedeemTokens_FullRedemption() {
	// First, create a tokenization record
	delegatorAddr := sdk.AccAddress([]byte("delegator"))
	validatorAddr := sdk.ValAddress([]byte("validator"))
	
	record := types.NewTokenizationRecordWithDenom(
		1,
		validatorAddr.String(),
		delegatorAddr.String(),
		math.NewInt(500000),
		"flora/lstake/"+validatorAddr.String()+"/1",
	)
	// For testing, we need to manually create the full record with denom index
	// This simulates what would happen during actual tokenization
	suite.keeper.SetTokenizationRecord(suite.ctx, record)
	suite.setTokenizationRecordIndexes(suite.ctx, record)
	
	// Setup mock expectations
	validator := stakingtypes.Validator{
		OperatorAddress: validatorAddr.String(),
		Status:          stakingtypes.Bonded,
		Tokens:          math.NewInt(1000000),
		DelegatorShares: math.LegacyNewDec(1000000),
		Jailed:          false,
	}
	
	suite.stakingKeeper.GetValidatorFn = func(ctx context.Context, addr sdk.ValAddress) (stakingtypes.Validator, error) {
		return validator, nil
	}
	
	// Mock bank keeper for redemption
	suite.bankKeeper.GetBalanceFn = func(ctx context.Context, addr sdk.AccAddress, denom string) sdk.Coin {
		return sdk.NewCoin(denom, math.NewInt(500000))
	}
	
	suite.bankKeeper.SendCoinsFromAccountToModuleFn = func(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error {
		return nil
	}
	
	suite.bankKeeper.BurnCoinsFn = func(ctx context.Context, moduleName string, amt sdk.Coins) error {
		return nil
	}
	
	suite.stakingKeeper.DelegateFn = func(ctx context.Context, delAddr sdk.AccAddress, bondAmt math.Int, tokenSrc stakingtypes.BondStatus, validator stakingtypes.Validator, subtractAccount bool) (math.LegacyDec, error) {
		return math.LegacyNewDecFromInt(bondAmt), nil
	}
	
	// Execute full redemption
	msg := &types.MsgRedeemTokens{
		OwnerAddress: delegatorAddr.String(),
		Amount:       sdk.NewCoin("flora/lstake/"+validatorAddr.String()+"/1", math.NewInt(500000)), // Redeem all
	}
	
	resp, err := suite.msgServer.RedeemTokens(suite.ctx, msg)
	suite.Require().NoError(err)
	suite.Require().NotNil(resp)
	suite.Require().Equal(uint64(1), resp.RecordId)
	
	// Verify record was deleted
	_, found := suite.keeper.GetTokenizationRecord(suite.ctx, 1)
	suite.Require().False(found)
	
	// Verify denom index was also deleted
	_, found = suite.keeper.GetTokenizationRecordByDenom(suite.ctx, "flora/lstake/"+validatorAddr.String()+"/1")
	suite.Require().False(found)
}