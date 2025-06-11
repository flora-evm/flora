package keeper_test

import (
	"context"
	"errors"
	"testing"

	cometbft "github.com/cometbft/cometbft/proto/tendermint/types"
	"cosmossdk.io/log"
	"cosmossdk.io/math"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/rollchains/flora/x/liquidstaking/keeper"
	"github.com/rollchains/flora/x/liquidstaking/testutil/mocks"
	"github.com/rollchains/flora/x/liquidstaking/types"
)

type StakingIntegrationTestSuite struct {
	suite.Suite
	
	ctx            sdk.Context
	keeper         keeper.Keeper
	stakingKeeper  *mocks.MockStakingKeeper
	bankKeeper     *mocks.MockBankKeeper
	accountKeeper  *mocks.MockAccountKeeper
	msgServer      types.MsgServer
}

func TestStakingIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(StakingIntegrationTestSuite))
}

func (suite *StakingIntegrationTestSuite) SetupTest() {
	// Create in-memory store
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey("mem_capability")
	
	db := store.NewCommitMultiStore(dbm.NewMemDB(), log.NewNopLogger(), metrics.NewNoOpMetrics())
	db.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, nil)
	db.MountStoreWithDB(memStoreKey, storetypes.StoreTypeMemory, nil)
	require.NoError(suite.T(), db.LoadLatestVersion())
	
	// Create codec
	registry := codectypes.NewInterfaceRegistry()
	types.RegisterInterfaces(registry)
	cdc := codec.NewProtoCodec(registry)
	
	// Create mocks
	suite.stakingKeeper = &mocks.MockStakingKeeper{}
	suite.bankKeeper = &mocks.MockBankKeeper{}
	suite.accountKeeper = &mocks.MockAccountKeeper{}
	
	// Create keeper
	storeService := runtime.NewKVStoreService(storeKey)
	suite.keeper = keeper.NewKeeper(
		storeService,
		cdc,
		suite.stakingKeeper,
		suite.bankKeeper,
		suite.accountKeeper,
	)
	
	// Create context
	header := cometbft.Header{Height: 1}
	suite.ctx = sdk.NewContext(db, header, false, log.NewNopLogger())
	
	// Set default params
	params := types.DefaultParams()
	params.Enabled = true
	suite.keeper.SetParams(suite.ctx, params)
	
	// Create msg server
	suite.msgServer = keeper.NewMsgServerImpl(suite.keeper)
}

// Test successful tokenization with active validator
func (suite *StakingIntegrationTestSuite) TestTokenizeShares_Success() {
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
	
	suite.stakingKeeper.UnbondFn = func(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, shares math.LegacyDec) (math.Int, error) {
		// Simulate unbonding - return tokens equal to shares (1:1 ratio)
		return shares.TruncateInt(), nil
	}
	
	// Mock total bonded tokens to avoid cap issues
	suite.stakingKeeper.TotalBondedTokensFn = func(ctx context.Context) math.Int {
		return math.NewInt(100000000) // 100M total bonded
	}
	
	suite.bankKeeper.MintCoinsFn = func(ctx context.Context, moduleName string, amt sdk.Coins) error {
		return nil
	}
	
	suite.bankKeeper.SendCoinsFromModuleToAccountFn = func(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error {
		return nil
	}
	
	suite.bankKeeper.SetDenomMetaDataFn = func(ctx context.Context, denomMetaData banktypes.Metadata) {
		// No-op for test
	}
	
	// Execute tokenization
	msg := &types.MsgTokenizeShares{
		DelegatorAddress: delegatorAddr.String(),
		ValidatorAddress: validatorAddr.String(),
		Shares:           sdk.NewCoin("shares", math.NewInt(500000)),
		OwnerAddress:     "",
	}
	
	resp, err := suite.msgServer.TokenizeShares(suite.ctx, msg)
	suite.Require().NoError(err)
	suite.Require().NotNil(resp)
	suite.Require().Equal(uint64(1), resp.RecordId)
	suite.Require().Equal("flora/lstake/"+validatorAddr.String()+"/1", resp.Denom)
	suite.Require().Equal(math.NewInt(500000), resp.Amount.Amount)
	
	// Verify record was created
	record, found := suite.keeper.GetTokenizationRecord(suite.ctx, resp.RecordId)
	suite.Require().True(found)
	suite.Require().Equal(delegatorAddr.String(), record.Owner)
	suite.Require().Equal(validatorAddr.String(), record.Validator)
	suite.Require().Equal(math.NewInt(500000), record.SharesTokenized)
}

// Test tokenization with jailed validator (should fail)
func (suite *StakingIntegrationTestSuite) TestTokenizeShares_JailedValidator() {
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
		Jailed:          true, // Validator is jailed
	}
	
	suite.stakingKeeper.GetDelegationFn = func(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (stakingtypes.Delegation, error) {
		return delegation, nil
	}
	
	suite.stakingKeeper.GetValidatorFn = func(ctx context.Context, addr sdk.ValAddress) (stakingtypes.Validator, error) {
		return validator, nil
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
	suite.Require().Contains(err.Error(), "validator is jailed")
}

// Test tokenization with insufficient delegation
func (suite *StakingIntegrationTestSuite) TestTokenizeShares_InsufficientDelegation() {
	// Setup addresses
	delegatorAddr := sdk.AccAddress([]byte("delegator"))
	validatorAddr := sdk.ValAddress([]byte("validator"))
	
	// Setup mock expectations
	delegation := stakingtypes.Delegation{
		DelegatorAddress: delegatorAddr.String(),
		ValidatorAddress: validatorAddr.String(),
		Shares:           math.LegacyNewDec(100000), // Only 100k shares
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
	
	// Execute tokenization with more shares than available
	msg := &types.MsgTokenizeShares{
		DelegatorAddress: delegatorAddr.String(),
		ValidatorAddress: validatorAddr.String(),
		Shares:           sdk.NewCoin("shares", math.NewInt(500000)), // Request 500k but only have 100k
		OwnerAddress:     "",
	}
	
	_, err := suite.msgServer.TokenizeShares(suite.ctx, msg)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "insufficient delegation shares")
}

// Test tokenization when delegation doesn't exist
func (suite *StakingIntegrationTestSuite) TestTokenizeShares_NoDelegation() {
	// Setup addresses
	delegatorAddr := sdk.AccAddress([]byte("delegator"))
	validatorAddr := sdk.ValAddress([]byte("validator"))
	
	// Setup mock expectations
	suite.stakingKeeper.GetDelegationFn = func(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (stakingtypes.Delegation, error) {
		return stakingtypes.Delegation{}, errors.New("delegation not found")
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
	suite.Require().Contains(err.Error(), "delegation not found")
}

// Test tokenization with validator not found
func (suite *StakingIntegrationTestSuite) TestTokenizeShares_ValidatorNotFound() {
	// Setup addresses
	delegatorAddr := sdk.AccAddress([]byte("delegator"))
	validatorAddr := sdk.ValAddress([]byte("validator"))
	
	// Setup mock expectations
	delegation := stakingtypes.Delegation{
		DelegatorAddress: delegatorAddr.String(),
		ValidatorAddress: validatorAddr.String(),
		Shares:           math.LegacyNewDec(1000000),
	}
	
	suite.stakingKeeper.GetDelegationFn = func(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (stakingtypes.Delegation, error) {
		return delegation, nil
	}
	
	suite.stakingKeeper.GetValidatorFn = func(ctx context.Context, addr sdk.ValAddress) (stakingtypes.Validator, error) {
		return stakingtypes.Validator{}, errors.New("validator not found")
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
	suite.Require().Contains(err.Error(), "validator not found")
}

// Test concurrent tokenization from multiple delegators
func (suite *StakingIntegrationTestSuite) TestTokenizeShares_MultipleDelegators() {
	// Setup addresses
	delegator1Addr := sdk.AccAddress([]byte("delegator1"))
	delegator2Addr := sdk.AccAddress([]byte("delegator2"))
	validatorAddr := sdk.ValAddress([]byte("validator"))
	
	validator := stakingtypes.Validator{
		OperatorAddress: validatorAddr.String(),
		Status:          stakingtypes.Bonded,
		Tokens:          math.NewInt(2000000),
		DelegatorShares: math.LegacyNewDec(2000000),
		Jailed:          false,
	}
	
	// Setup mock expectations for different delegators
	suite.stakingKeeper.GetDelegationFn = func(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (stakingtypes.Delegation, error) {
		if delAddr.Equals(delegator1Addr) {
			return stakingtypes.Delegation{
				DelegatorAddress: delegator1Addr.String(),
				ValidatorAddress: validatorAddr.String(),
				Shares:           math.LegacyNewDec(1000000),
			}, nil
		} else if delAddr.Equals(delegator2Addr) {
			return stakingtypes.Delegation{
				DelegatorAddress: delegator2Addr.String(),
				ValidatorAddress: validatorAddr.String(),
				Shares:           math.LegacyNewDec(1000000),
			}, nil
		}
		return stakingtypes.Delegation{}, errors.New("delegation not found")
	}
	
	suite.stakingKeeper.GetValidatorFn = func(ctx context.Context, addr sdk.ValAddress) (stakingtypes.Validator, error) {
		return validator, nil
	}
	
	suite.stakingKeeper.UnbondFn = func(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, shares math.LegacyDec) (math.Int, error) {
		return shares.TruncateInt(), nil
	}
	
	// Mock total bonded tokens to avoid cap issues
	suite.stakingKeeper.TotalBondedTokensFn = func(ctx context.Context) math.Int {
		return math.NewInt(100000000) // 100M total bonded
	}
	
	suite.bankKeeper.MintCoinsFn = func(ctx context.Context, moduleName string, amt sdk.Coins) error {
		return nil
	}
	
	suite.bankKeeper.SendCoinsFromModuleToAccountFn = func(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error {
		return nil
	}
	
	suite.bankKeeper.SetDenomMetaDataFn = func(ctx context.Context, denomMetaData banktypes.Metadata) {}
	
	// Tokenize from delegator 1
	msg1 := &types.MsgTokenizeShares{
		DelegatorAddress: delegator1Addr.String(),
		ValidatorAddress: validatorAddr.String(),
		Shares:           sdk.NewCoin("shares", math.NewInt(500000)),
		OwnerAddress:     "",
	}
	
	resp1, err := suite.msgServer.TokenizeShares(suite.ctx, msg1)
	suite.Require().NoError(err)
	suite.Require().Equal(uint64(1), resp1.RecordId)
	
	// Tokenize from delegator 2
	msg2 := &types.MsgTokenizeShares{
		DelegatorAddress: delegator2Addr.String(),
		ValidatorAddress: validatorAddr.String(),
		Shares:           sdk.NewCoin("shares", math.NewInt(300000)),
		OwnerAddress:     "",
	}
	
	resp2, err := suite.msgServer.TokenizeShares(suite.ctx, msg2)
	suite.Require().NoError(err)
	suite.Require().Equal(uint64(2), resp2.RecordId)
	
	// Verify both records exist
	record1, found := suite.keeper.GetTokenizationRecord(suite.ctx, 1)
	suite.Require().True(found)
	suite.Require().Equal(delegator1Addr.String(), record1.Owner)
	
	record2, found := suite.keeper.GetTokenizationRecord(suite.ctx, 2)
	suite.Require().True(found)
	suite.Require().Equal(delegator2Addr.String(), record2.Owner)
}