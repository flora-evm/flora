package keeper_test

import (
	"context"
	"testing"
	
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"
	
	"github.com/rollchains/flora/x/liquidstaking/keeper"
	"github.com/rollchains/flora/x/liquidstaking/testutil/mocks"
	"github.com/rollchains/flora/x/liquidstaking/types"
)

func TestMsgServer_TokenizeShares(t *testing.T) {
	delegatorAddr := sdk.AccAddress("delegator_address___")
	validatorAddr := sdk.ValAddress("validator_address___")
	ownerAddr := sdk.AccAddress("owner_address_______")
	
	testCases := []struct {
		name          string
		msg           *types.MsgTokenizeShares
		setupMocks    func(*mocks.MockStakingKeeper, *mocks.MockBankKeeper, *mocks.MockAccountKeeper)
		setupKeeper   func(*KeeperTestSuite)
		expectedError error
		validate      func(*KeeperTestSuite, *types.MsgTokenizeSharesResponse)
	}{
		{
			name: "successful tokenization",
			msg: &types.MsgTokenizeShares{
				DelegatorAddress: delegatorAddr.String(),
				ValidatorAddress: validatorAddr.String(),
				Shares:           sdk.NewCoin("stake", math.NewInt(1000000)),
				OwnerAddress:     ownerAddr.String(),
			},
			setupMocks: func(stakingKeeper *mocks.MockStakingKeeper, bankKeeper *mocks.MockBankKeeper, accountKeeper *mocks.MockAccountKeeper) {
				// Mock delegation exists
				stakingKeeper.GetDelegationFn = func(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (stakingtypes.Delegation, error) {
					return stakingtypes.Delegation{
						DelegatorAddress: delAddr.String(),
						ValidatorAddress: valAddr.String(),
						Shares:           math.LegacyNewDec(2000000), // Has 2M shares
					}, nil
				}
				
				// Mock validator exists and is not jailed
				stakingKeeper.GetValidatorFn = func(ctx context.Context, addr sdk.ValAddress) (stakingtypes.Validator, error) {
					val := stakingtypes.Validator{
						OperatorAddress: addr.String(),
						Jailed:          false,
						Tokens:          math.NewInt(10000000), // 10M tokens
						DelegatorShares: math.LegacyNewDec(10000000), // 10M shares
					}
					return val, nil
				}
				
				// Mock unbond succeeds
				stakingKeeper.UnbondFn = func(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, shares math.LegacyDec) (math.Int, error) {
					return math.NewInt(1000000), nil // Returns 1M tokens
				}
				
				// Mock total bonded tokens
				stakingKeeper.TotalBondedTokensFn = func(ctx context.Context) math.Int {
					return math.NewInt(100000000) // 100M total bonded
				}
				
				// Mock mint coins succeeds
				bankKeeper.MintCoinsFn = func(ctx context.Context, moduleName string, amt sdk.Coins) error {
					return nil
				}
				
				// Mock send coins succeeds
				bankKeeper.SendCoinsFromModuleToAccountFn = func(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error {
					return nil
				}
				
				// Mock set denom metadata
				bankKeeper.SetDenomMetaDataFn = func(ctx context.Context, denomMetaData banktypes.Metadata) {
					// Do nothing, just track it was called
				}
			},
			setupKeeper: func(suite *KeeperTestSuite) {
				// Set module params
				params := types.DefaultParams()
				suite.keeper.SetParams(suite.ctx, params)
			},
			expectedError: nil,
			validate: func(suite *KeeperTestSuite, resp *types.MsgTokenizeSharesResponse) {
				require.NotNil(t, resp)
				require.NotEmpty(t, resp.Denom)
				require.True(t, types.IsLiquidStakingTokenDenom(resp.Denom))
				require.Equal(t, resp.Amount.Amount, math.NewInt(1000000))
				require.Equal(t, resp.RecordId, uint64(1))
				
				// Verify tokenization record was created
				record, found := suite.keeper.GetTokenizationRecord(suite.ctx, resp.RecordId)
				require.True(t, found)
				require.Equal(t, record.Validator, validatorAddr.String())
				require.Equal(t, record.Owner, ownerAddr.String())
				require.Equal(t, record.SharesTokenized, math.NewInt(1000000))
				require.Equal(t, record.Denom, resp.Denom)
				
				// Verify indexes were created
				records := suite.keeper.GetTokenizationRecordsByValidator(suite.ctx, validatorAddr.String())
				require.Len(t, records, 1)
				require.Equal(t, records[0].Id, resp.RecordId)
				
				records = suite.keeper.GetTokenizationRecordsByOwner(suite.ctx, ownerAddr.String())
				require.Len(t, records, 1)
				require.Equal(t, records[0].Id, resp.RecordId)
				
				// Verify denom index
				recordByDenom, found := suite.keeper.GetTokenizationRecordByDenom(suite.ctx, resp.Denom)
				require.True(t, found)
				require.Equal(t, recordByDenom.Id, resp.RecordId)
				
				// Verify liquid staked amounts were updated
				totalLiquidStaked := suite.keeper.GetTotalLiquidStaked(suite.ctx)
				require.Equal(t, totalLiquidStaked, math.NewInt(1000000))
				
				validatorLiquidStaked := suite.keeper.GetValidatorLiquidStaked(suite.ctx, validatorAddr.String())
				require.Equal(t, validatorLiquidStaked, math.NewInt(1000000))
			},
		},
		{
			name: "module disabled",
			msg: &types.MsgTokenizeShares{
				DelegatorAddress: delegatorAddr.String(),
				ValidatorAddress: validatorAddr.String(),
				Shares:           sdk.NewCoin("stake", math.NewInt(1000000)),
				OwnerAddress:     "",
			},
			setupMocks: func(stakingKeeper *mocks.MockStakingKeeper, bankKeeper *mocks.MockBankKeeper, accountKeeper *mocks.MockAccountKeeper) {
				// No mocks needed, should fail at module check
			},
			setupKeeper: func(suite *KeeperTestSuite) {
				// Disable module
				params := types.DefaultParams()
				params.Enabled = false
				suite.keeper.SetParams(suite.ctx, params)
			},
			expectedError: types.ErrModuleDisabled,
			validate:      nil,
		},
		{
			name: "delegation not found",
			msg: &types.MsgTokenizeShares{
				DelegatorAddress: delegatorAddr.String(),
				ValidatorAddress: validatorAddr.String(),
				Shares:           sdk.NewCoin("stake", math.NewInt(1000000)),
				OwnerAddress:     "",
			},
			setupMocks: func(stakingKeeper *mocks.MockStakingKeeper, bankKeeper *mocks.MockBankKeeper, accountKeeper *mocks.MockAccountKeeper) {
				// Mock delegation not found
				stakingKeeper.GetDelegationFn = func(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (stakingtypes.Delegation, error) {
					return stakingtypes.Delegation{}, types.ErrDelegationNotFound
				}
			},
			setupKeeper: func(suite *KeeperTestSuite) {
				params := types.DefaultParams()
				suite.keeper.SetParams(suite.ctx, params)
			},
			expectedError: types.ErrDelegationNotFound,
			validate:      nil,
		},
		{
			name: "insufficient shares",
			msg: &types.MsgTokenizeShares{
				DelegatorAddress: delegatorAddr.String(),
				ValidatorAddress: validatorAddr.String(),
				Shares:           sdk.NewCoin("stake", math.NewInt(3000000)), // Wants 3M but has 2M
				OwnerAddress:     "",
			},
			setupMocks: func(stakingKeeper *mocks.MockStakingKeeper, bankKeeper *mocks.MockBankKeeper, accountKeeper *mocks.MockAccountKeeper) {
				// Mock delegation with insufficient shares
				stakingKeeper.GetDelegationFn = func(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (stakingtypes.Delegation, error) {
					return stakingtypes.Delegation{
						DelegatorAddress: delAddr.String(),
						ValidatorAddress: valAddr.String(),
						Shares:           math.LegacyNewDec(2000000), // Has 2M shares
					}, nil
				}
			},
			setupKeeper: func(suite *KeeperTestSuite) {
				params := types.DefaultParams()
				suite.keeper.SetParams(suite.ctx, params)
			},
			expectedError: types.ErrInsufficientShares,
			validate:      nil,
		},
		{
			name: "validator jailed",
			msg: &types.MsgTokenizeShares{
				DelegatorAddress: delegatorAddr.String(),
				ValidatorAddress: validatorAddr.String(),
				Shares:           sdk.NewCoin("stake", math.NewInt(1000000)),
				OwnerAddress:     "",
			},
			setupMocks: func(stakingKeeper *mocks.MockStakingKeeper, bankKeeper *mocks.MockBankKeeper, accountKeeper *mocks.MockAccountKeeper) {
				// Mock delegation exists
				stakingKeeper.GetDelegationFn = func(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (stakingtypes.Delegation, error) {
					return stakingtypes.Delegation{
						DelegatorAddress: delAddr.String(),
						ValidatorAddress: valAddr.String(),
						Shares:           math.LegacyNewDec(2000000),
					}, nil
				}
				
				// Mock validator is jailed
				stakingKeeper.GetValidatorFn = func(ctx context.Context, addr sdk.ValAddress) (stakingtypes.Validator, error) {
					val := stakingtypes.Validator{
						OperatorAddress: addr.String(),
						Jailed:          true, // Jailed!
						Tokens:          math.NewInt(10000000),
						DelegatorShares: math.LegacyNewDec(10000000),
					}
					return val, nil
				}
			},
			setupKeeper: func(suite *KeeperTestSuite) {
				params := types.DefaultParams()
				suite.keeper.SetParams(suite.ctx, params)
			},
			expectedError: types.ErrInvalidValidator,
			validate:      nil,
		},
		{
			name: "exceeds global cap",
			msg: &types.MsgTokenizeShares{
				DelegatorAddress: delegatorAddr.String(),
				ValidatorAddress: validatorAddr.String(),
				Shares:           sdk.NewCoin("stake", math.NewInt(30000000)), // 30M tokens
				OwnerAddress:     "",
			},
			setupMocks: func(stakingKeeper *mocks.MockStakingKeeper, bankKeeper *mocks.MockBankKeeper, accountKeeper *mocks.MockAccountKeeper) {
				// Mock delegation exists
				stakingKeeper.GetDelegationFn = func(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (stakingtypes.Delegation, error) {
					return stakingtypes.Delegation{
						DelegatorAddress: delAddr.String(),
						ValidatorAddress: valAddr.String(),
						Shares:           math.LegacyNewDec(30000000),
					}, nil
				}
				
				// Mock validator
				stakingKeeper.GetValidatorFn = func(ctx context.Context, addr sdk.ValAddress) (stakingtypes.Validator, error) {
					val := stakingtypes.Validator{
						OperatorAddress: addr.String(),
						Jailed:          false,
						Tokens:          math.NewInt(100000000), // 100M tokens
						DelegatorShares: math.LegacyNewDec(100000000),
					}
					return val, nil
				}
				
				// Mock total bonded tokens (100M)
				stakingKeeper.TotalBondedTokensFn = func(ctx context.Context) math.Int {
					return math.NewInt(100000000)
				}
			},
			setupKeeper: func(suite *KeeperTestSuite) {
				params := types.DefaultParams()
				params.GlobalLiquidStakingCap = math.LegacyNewDecWithPrec(25, 2) // 25%
				suite.keeper.SetParams(suite.ctx, params)
				
				// Already have 5M liquid staked
				suite.keeper.SetTotalLiquidStaked(suite.ctx, math.NewInt(5000000))
			},
			expectedError: types.ErrExceedsGlobalCap,
			validate:      nil,
		},
		{
			name: "exceeds validator cap",
			msg: &types.MsgTokenizeShares{
				DelegatorAddress: delegatorAddr.String(),
				ValidatorAddress: validatorAddr.String(),
				Shares:           sdk.NewCoin("stake", math.NewInt(6000000)), // 6M tokens
				OwnerAddress:     "",
			},
			setupMocks: func(stakingKeeper *mocks.MockStakingKeeper, bankKeeper *mocks.MockBankKeeper, accountKeeper *mocks.MockAccountKeeper) {
				// Mock delegation exists
				stakingKeeper.GetDelegationFn = func(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (stakingtypes.Delegation, error) {
					return stakingtypes.Delegation{
						DelegatorAddress: delAddr.String(),
						ValidatorAddress: valAddr.String(),
						Shares:           math.LegacyNewDec(6000000),
					}, nil
				}
				
				// Mock validator
				stakingKeeper.GetValidatorFn = func(ctx context.Context, addr sdk.ValAddress) (stakingtypes.Validator, error) {
					val := stakingtypes.Validator{
						OperatorAddress: addr.String(),
						Jailed:          false,
						Tokens:          math.NewInt(10000000), // 10M tokens
						DelegatorShares: math.LegacyNewDec(10000000),
					}
					return val, nil
				}
				
				// Mock total bonded tokens
				stakingKeeper.TotalBondedTokensFn = func(ctx context.Context) math.Int {
					return math.NewInt(100000000)
				}
			},
			setupKeeper: func(suite *KeeperTestSuite) {
				params := types.DefaultParams()
				params.ValidatorLiquidCap = math.LegacyNewDecWithPrec(50, 2) // 50%
				suite.keeper.SetParams(suite.ctx, params)
				
				// Already have 4M liquid staked for this validator
				suite.keeper.SetValidatorLiquidStaked(suite.ctx, validatorAddr.String(), math.NewInt(4000000))
			},
			expectedError: types.ErrExceedsValidatorCap,
			validate:      nil,
		},
		{
			name: "default owner to delegator",
			msg: &types.MsgTokenizeShares{
				DelegatorAddress: delegatorAddr.String(),
				ValidatorAddress: validatorAddr.String(),
				Shares:           sdk.NewCoin("stake", math.NewInt(1000000)),
				OwnerAddress:     "", // Empty, should default to delegator
			},
			setupMocks: func(stakingKeeper *mocks.MockStakingKeeper, bankKeeper *mocks.MockBankKeeper, accountKeeper *mocks.MockAccountKeeper) {
				// Mock delegation exists
				stakingKeeper.GetDelegationFn = func(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (stakingtypes.Delegation, error) {
					return stakingtypes.Delegation{
						DelegatorAddress: delAddr.String(),
						ValidatorAddress: valAddr.String(),
						Shares:           math.LegacyNewDec(2000000),
					}, nil
				}
				
				// Mock validator
				stakingKeeper.GetValidatorFn = func(ctx context.Context, addr sdk.ValAddress) (stakingtypes.Validator, error) {
					val := stakingtypes.Validator{
						OperatorAddress: addr.String(),
						Jailed:          false,
						Tokens:          math.NewInt(10000000),
						DelegatorShares: math.LegacyNewDec(10000000),
					}
					return val, nil
				}
				
				// Mock unbond
				stakingKeeper.UnbondFn = func(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, shares math.LegacyDec) (math.Int, error) {
					return math.NewInt(1000000), nil
				}
				
				// Mock total bonded
				stakingKeeper.TotalBondedTokensFn = func(ctx context.Context) math.Int {
					return math.NewInt(100000000)
				}
				
				// Mock bank operations
				bankKeeper.MintCoinsFn = func(ctx context.Context, moduleName string, amt sdk.Coins) error {
					return nil
				}
				
				bankKeeper.SendCoinsFromModuleToAccountFn = func(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error {
					// Verify tokens are sent to delegator (since owner is empty)
					require.Equal(t, delegatorAddr, recipientAddr)
					return nil
				}
				
				bankKeeper.SetDenomMetaDataFn = func(ctx context.Context, denomMetaData banktypes.Metadata) {}
			},
			setupKeeper: func(suite *KeeperTestSuite) {
				params := types.DefaultParams()
				suite.keeper.SetParams(suite.ctx, params)
			},
			expectedError: nil,
			validate: func(suite *KeeperTestSuite, resp *types.MsgTokenizeSharesResponse) {
				require.NotNil(t, resp)
				
				// Verify record owner is delegator
				record, found := suite.keeper.GetTokenizationRecord(suite.ctx, resp.RecordId)
				require.True(t, found)
				require.Equal(t, record.Owner, delegatorAddr.String())
			},
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			suite := setupKeeperTestSuite(t)
			
			// Set up mocks
			stakingKeeper := &mocks.MockStakingKeeper{}
			bankKeeper := &mocks.MockBankKeeper{}
			accountKeeper := &mocks.MockAccountKeeper{}
			
			if tc.setupMocks != nil {
				tc.setupMocks(stakingKeeper, bankKeeper, accountKeeper)
			}
			
			// Create keeper with mocks
			suite.keeper = keeper.NewKeeper(
				suite.storeService,
				suite.cdc,
				stakingKeeper,
				bankKeeper,
				accountKeeper,
			)
			
			// Setup keeper state
			if tc.setupKeeper != nil {
				tc.setupKeeper(suite)
			}
			
			// Create message server
			msgServer := keeper.NewMsgServerImpl(suite.keeper)
			
			// Execute
			resp, err := msgServer.TokenizeShares(suite.ctx, tc.msg)
			
			// Validate
			if tc.expectedError != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tc.expectedError)
			} else {
				require.NoError(t, err)
				if tc.validate != nil {
					tc.validate(suite, resp)
				}
			}
		})
	}
}