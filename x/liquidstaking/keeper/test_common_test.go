package keeper_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"cosmossdk.io/log"
	"cosmossdk.io/math"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"

	"github.com/rollchains/flora/x/liquidstaking/keeper"
	"github.com/rollchains/flora/x/liquidstaking/testutil/mocks"
	"github.com/rollchains/flora/x/liquidstaking/types"
)

var (
	mockBankKeeper         *mocks.MockBankKeeper
	mockStakingKeeper      *mocks.MockStakingKeeper
	mockAccountKeeper      *mocks.MockAccountKeeper
	mockDistributionKeeper *mocks.MockDistributionKeeper
	
	// Test addresses - generated properly
	testValAddr1 sdk.ValAddress
	testValAddr2 sdk.ValAddress
	testAccAddr1 sdk.AccAddress
	testAccAddr2 sdk.AccAddress
)

// setupKeeper creates a keeper for testing with mock dependencies
func setupKeeper(t *testing.T) (keeper.Keeper, sdk.Context) {
	// Initialize test addresses
	privKey1 := ed25519.GenPrivKey()
	privKey2 := ed25519.GenPrivKey()
	testAccAddr1 = sdk.AccAddress(privKey1.PubKey().Address())
	testAccAddr2 = sdk.AccAddress(privKey2.PubKey().Address())
	testValAddr1 = sdk.ValAddress(privKey1.PubKey().Address())
	testValAddr2 = sdk.ValAddress(privKey2.PubKey().Address())
	
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	
	db := dbm.NewMemDB()
	stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	require := require.New(t)
	require.NoError(stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	types.RegisterInterfaces(registry)
	cdc := codec.NewProtoCodec(registry)

	// Create mock keepers
	mockBankKeeper = &mocks.MockBankKeeper{}
	mockStakingKeeper = &mocks.MockStakingKeeper{}
	mockAccountKeeper = &mocks.MockAccountKeeper{
		GetModuleAddressFn: func(moduleName string) sdk.AccAddress {
			return authtypes.NewModuleAddress(moduleName)
		},
	}

	// Create mock IBC keepers
	mockTransferKeeper := &mocks.MockTransferKeeper{}
	mockChannelKeeper := &mocks.MockChannelKeeper{}
	mockDistributionKeeper = &mocks.MockDistributionKeeper{}

	// Use gov module address as authority
	auth := authtypes.NewModuleAddress("gov").String()

	storeService := runtime.NewKVStoreService(storeKey)
	k := keeper.NewKeeper(
		storeService,
		cdc,
		mockStakingKeeper,
		mockBankKeeper,
		mockAccountKeeper,
		mockTransferKeeper,
		mockChannelKeeper,
		mockDistributionKeeper,
		auth,
	)

	ctx := sdk.NewContext(stateStore, cmtproto.Header{}, false, log.NewNopLogger())
	
	// Set default params
	params := types.DefaultParams()
	params.Enabled = true
	k.SetParams(ctx, params)
	
	return k, ctx
}

// createTestValidator creates a test validator for testing
func createTestValidator(valAddr sdk.ValAddress, tokens math.Int) stakingtypes.Validator {
	pubKey := ed25519.GenPrivKey().PubKey()
	
	return stakingtypes.Validator{
		OperatorAddress: valAddr.String(),
		ConsensusPubkey: codectypes.UnsafePackAny(pubKey),
		Jailed:          false,
		Status:          stakingtypes.Bonded,
		Tokens:          tokens,
		DelegatorShares: math.LegacyNewDecFromInt(tokens),
		Description: stakingtypes.Description{
			Moniker: "test validator",
		},
		UnbondingHeight: 0,
		UnbondingTime:   time.Time{},
		Commission: stakingtypes.Commission{
			CommissionRates: stakingtypes.CommissionRates{
				Rate:          math.LegacyNewDecWithPrec(1, 1), // 10%
				MaxRate:       math.LegacyNewDecWithPrec(2, 1), // 20%
				MaxChangeRate: math.LegacyNewDecWithPrec(1, 2), // 1%
			},
		},
		MinSelfDelegation: math.OneInt(),
	}
}

// createTestDelegation creates a test delegation
func createTestDelegation(delAddr sdk.AccAddress, valAddr sdk.ValAddress, shares math.LegacyDec) stakingtypes.Delegation {
	return stakingtypes.Delegation{
		DelegatorAddress: delAddr.String(),
		ValidatorAddress: valAddr.String(),
		Shares:           shares,
	}
}

// expectMintCoins sets up expectation for minting coins
func expectMintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) {
	mockBankKeeper.MintCoinsFn = func(c context.Context, mName string, a sdk.Coins) error {
		if mName == moduleName && a.Equal(amt) {
			return nil
		}
		return errors.New("unexpected mint coins call")
	}
}

// expectSendCoinsFromModuleToAccount sets up expectation for sending coins from module to account
func expectSendCoinsFromModuleToAccount(ctx sdk.Context, moduleName string, recipientAddr sdk.AccAddress, amt sdk.Coins) {
	mockBankKeeper.SendCoinsFromModuleToAccountFn = func(c context.Context, senderModule string, recipient sdk.AccAddress, a sdk.Coins) error {
		if senderModule == moduleName && recipient.Equals(recipientAddr) && a.Equal(amt) {
			return nil
		}
		return errors.New("unexpected send coins call")
	}
}

// expectGetBalance sets up expectation for getting balance
func expectGetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string, amount math.Int) {
	mockBankKeeper.GetBalanceFn = func(c context.Context, a sdk.AccAddress, d string) sdk.Coin {
		if a.Equals(addr) && d == denom {
			return sdk.NewCoin(denom, amount)
		}
		return sdk.NewCoin(denom, math.ZeroInt())
	}
}

// expectBurnCoins sets up expectation for burning coins
func expectBurnCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) {
	mockBankKeeper.BurnCoinsFn = func(c context.Context, mName string, a sdk.Coins) error {
		if mName == moduleName && a.Equal(amt) {
			return nil
		}
		return errors.New("unexpected burn coins call")
	}
}

// expectSendCoinsFromAccountToModule sets up expectation for sending coins from account to module
func expectSendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) {
	mockBankKeeper.SendCoinsFromAccountToModuleFn = func(c context.Context, sender sdk.AccAddress, recipient string, a sdk.Coins) error {
		if sender.Equals(senderAddr) && recipient == recipientModule && a.Equal(amt) {
			return nil
		}
		return errors.New("unexpected send coins to module call")
	}
}

// expectSetDenomMetaData sets up expectation for setting denom metadata
func expectSetDenomMetaData(ctx sdk.Context, denomMetaData banktypes.Metadata) {
	mockBankKeeper.SetDenomMetaDataFn = func(c context.Context, metadata banktypes.Metadata) {
		// Just accept any metadata for now
	}
}

// expectGetValidator sets up expectation for getting validator
func expectGetValidator(ctx sdk.Context, valAddr sdk.ValAddress, validator stakingtypes.Validator) {
	mockStakingKeeper.GetValidatorFn = func(c context.Context, addr sdk.ValAddress) (stakingtypes.Validator, error) {
		if addr.Equals(valAddr) {
			return validator, nil
		}
		return stakingtypes.Validator{}, errors.New("validator not found")
	}
}

// expectGetDelegation sets up expectation for getting delegation
func expectGetDelegation(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, delegation stakingtypes.Delegation) {
	mockStakingKeeper.GetDelegationFn = func(c context.Context, del sdk.AccAddress, val sdk.ValAddress) (stakingtypes.Delegation, error) {
		if del.Equals(delAddr) && val.Equals(valAddr) {
			return delegation, nil
		}
		return stakingtypes.Delegation{}, errors.New("delegation not found")
	}
}

// expectUnbond sets up expectation for unbonding
func expectUnbond(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, shares math.LegacyDec, unbondedTokens math.Int) {
	mockStakingKeeper.UnbondFn = func(c context.Context, del sdk.AccAddress, val sdk.ValAddress, s math.LegacyDec) (math.Int, error) {
		if del.Equals(delAddr) && val.Equals(valAddr) && s.Equal(shares) {
			return unbondedTokens, nil
		}
		return math.ZeroInt(), errors.New("unexpected unbond call")
	}
}

// expectDelegate sets up expectation for delegation
func expectDelegate(ctx sdk.Context, delAddr sdk.AccAddress, bondAmt math.Int, tokenSrc stakingtypes.BondStatus, validator stakingtypes.Validator, subtractAccount bool, shares math.LegacyDec) {
	mockStakingKeeper.DelegateFn = func(c context.Context, del sdk.AccAddress, amt math.Int, src stakingtypes.BondStatus, val stakingtypes.Validator, subtract bool) (math.LegacyDec, error) {
		if del.Equals(delAddr) && amt.Equal(bondAmt) && src == tokenSrc && subtract == subtractAccount {
			return shares, nil
		}
		return math.LegacyZeroDec(), errors.New("unexpected delegate call")
	}
}