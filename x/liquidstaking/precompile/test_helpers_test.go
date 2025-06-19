package precompile_test

import (
	"context"
	"testing"
	"time"

	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"

	"github.com/rollchains/flora/x/liquidstaking/keeper"
	"github.com/rollchains/flora/x/liquidstaking/types"
)

// Mock keepers for testing

type MockBankKeeper struct {
	balances map[string]sdk.Coins
	supply   map[string]sdk.Coin
	metadata map[string]banktypes.Metadata
}

func NewMockBankKeeper() *MockBankKeeper {
	return &MockBankKeeper{
		balances: make(map[string]sdk.Coins),
		supply:   make(map[string]sdk.Coin),
		metadata: make(map[string]banktypes.Metadata),
	}
}

func (m *MockBankKeeper) MintCoins(ctx context.Context, moduleName string, amt sdk.Coins) error {
	for _, coin := range amt {
		if supply, ok := m.supply[coin.Denom]; ok {
			m.supply[coin.Denom] = supply.Add(coin)
		} else {
			m.supply[coin.Denom] = coin
		}
	}
	return nil
}

func (m *MockBankKeeper) SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error {
	addr := recipientAddr.String()
	if balance, ok := m.balances[addr]; ok {
		m.balances[addr] = balance.Add(amt...)
	} else {
		m.balances[addr] = amt
	}
	return nil
}

func (m *MockBankKeeper) GetDenomMetaData(ctx context.Context, denom string) (banktypes.Metadata, bool) {
	metadata, found := m.metadata[denom]
	return metadata, found
}

func (m *MockBankKeeper) SetDenomMetaData(ctx context.Context, denomMetaData banktypes.Metadata) {
	m.metadata[denomMetaData.Base] = denomMetaData
}

func (m *MockBankKeeper) GetSupply(ctx context.Context, denom string) sdk.Coin {
	if supply, ok := m.supply[denom]; ok {
		return supply
	}
	return sdk.NewCoin(denom, math.ZeroInt())
}

func (m *MockBankKeeper) GetBalance(ctx context.Context, addr sdk.AccAddress, denom string) sdk.Coin {
	if balance, ok := m.balances[addr.String()]; ok {
		amount := balance.AmountOf(denom)
		return sdk.NewCoin(denom, amount)
	}
	return sdk.NewCoin(denom, math.ZeroInt())
}

func (m *MockBankKeeper) SendCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error {
	addr := senderAddr.String()
	if balance, ok := m.balances[addr]; ok {
		newBalance, neg := balance.SafeSub(amt...)
		if neg {
			return sdk.ErrInsufficientFunds
		}
		m.balances[addr] = newBalance
	} else {
		return sdk.ErrInsufficientFunds
	}
	return nil
}

func (m *MockBankKeeper) BurnCoins(ctx context.Context, moduleName string, amt sdk.Coins) error {
	for _, coin := range amt {
		if supply, ok := m.supply[coin.Denom]; ok {
			m.supply[coin.Denom] = supply.Sub(coin)
		}
	}
	return nil
}

func (m *MockBankKeeper) SendCoins(ctx context.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error {
	from := fromAddr.String()
	to := toAddr.String()
	
	// Deduct from sender
	if balance, ok := m.balances[from]; ok {
		newBalance, neg := balance.SafeSub(amt...)
		if neg {
			return sdk.ErrInsufficientFunds
		}
		m.balances[from] = newBalance
	} else {
		return sdk.ErrInsufficientFunds
	}
	
	// Add to recipient
	if balance, ok := m.balances[to]; ok {
		m.balances[to] = balance.Add(amt...)
	} else {
		m.balances[to] = amt
	}
	
	return nil
}

type MockStakingKeeper struct {
	validators  map[string]stakingtypes.Validator
	delegations map[string]stakingtypes.Delegation
	params      stakingtypes.Params
}

func NewMockStakingKeeper() *MockStakingKeeper {
	return &MockStakingKeeper{
		validators:  make(map[string]stakingtypes.Validator),
		delegations: make(map[string]stakingtypes.Delegation),
		params: stakingtypes.Params{
			BondDenom: "stake",
		},
	}
}

func (m *MockStakingKeeper) GetDelegation(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (stakingtypes.Delegation, error) {
	key := delAddr.String() + valAddr.String()
	if del, ok := m.delegations[key]; ok {
		return del, nil
	}
	return stakingtypes.Delegation{}, stakingtypes.ErrNoDelegation
}

func (m *MockStakingKeeper) GetValidator(ctx context.Context, addr sdk.ValAddress) (stakingtypes.Validator, error) {
	if val, ok := m.validators[addr.String()]; ok {
		return val, nil
	}
	return stakingtypes.Validator{}, stakingtypes.ErrNoValidatorFound
}

func (m *MockStakingKeeper) Unbond(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, shares math.LegacyDec) (math.Int, error) {
	key := delAddr.String() + valAddr.String()
	if del, ok := m.delegations[key]; ok {
		del.Shares = del.Shares.Sub(shares)
		m.delegations[key] = del
		
		// Return tokens based on validator's shares to tokens ratio
		if val, ok := m.validators[valAddr.String()]; ok {
			return val.TokensFromShares(shares).TruncateInt(), nil
		}
	}
	return math.ZeroInt(), stakingtypes.ErrNoDelegation
}

func (m *MockStakingKeeper) GetParams(ctx context.Context) (stakingtypes.Params, error) {
	return m.params, nil
}

func (m *MockStakingKeeper) BondDenom(ctx context.Context) (string, error) {
	return m.params.BondDenom, nil
}

func (m *MockStakingKeeper) TotalBondedTokens(ctx context.Context) (math.Int, error) {
	total := math.ZeroInt()
	for _, val := range m.validators {
		if val.IsBonded() {
			total = total.Add(val.Tokens)
		}
	}
	return total, nil
}

func (m *MockStakingKeeper) Delegate(ctx context.Context, delAddr sdk.AccAddress, bondAmt math.Int, tokenSrc stakingtypes.BondStatus, validator stakingtypes.Validator, subtractAccount bool) (math.LegacyDec, error) {
	shares := validator.SharesFromTokens(bondAmt)
	
	key := delAddr.String() + validator.OperatorAddress
	if del, ok := m.delegations[key]; ok {
		del.Shares = del.Shares.Add(shares)
		m.delegations[key] = del
	} else {
		m.delegations[key] = stakingtypes.Delegation{
			DelegatorAddress: delAddr.String(),
			ValidatorAddress: validator.OperatorAddress,
			Shares:           shares,
		}
	}
	
	return shares, nil
}

type MockAccountKeeper struct {
	accounts map[string]sdk.AccountI
	modules  map[string]sdk.AccAddress
}

func NewMockAccountKeeper() *MockAccountKeeper {
	return &MockAccountKeeper{
		accounts: make(map[string]sdk.AccountI),
		modules:  make(map[string]sdk.AccAddress),
	}
}

func (m *MockAccountKeeper) GetAccount(ctx context.Context, addr sdk.AccAddress) sdk.AccountI {
	if acc, ok := m.accounts[addr.String()]; ok {
		return acc
	}
	// Return a basic account if not found
	return authtypes.NewBaseAccountWithAddress(addr)
}

func (m *MockAccountKeeper) GetModuleAddress(moduleName string) sdk.AccAddress {
	if addr, ok := m.modules[moduleName]; ok {
		return addr
	}
	// Generate a deterministic address for the module
	return sdk.AccAddress([]byte(moduleName))
}

// Setup functions

func setupTestEnvironment(t *testing.T) (sdk.Context, keeper.Keeper, *MockBankKeeper, *MockStakingKeeper, *MockAccountKeeper) {
	// Create store
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey("mem_" + types.StoreKey)
	
	storeService := runtime.NewKVStoreService(storeKey)
	
	// Create mock keepers
	bankKeeper := NewMockBankKeeper()
	stakingKeeper := NewMockStakingKeeper()
	accountKeeper := NewMockAccountKeeper()
	
	// Create keeper
	k := keeper.NewKeeper(
		log.NewNopLogger(),
		storeService,
		bankKeeper,
		stakingKeeper,
		accountKeeper,
		"flora",
	)
	
	// Create context
	ctx := sdk.NewContext(nil, sdk.BlockHeader{
		Height: 1,
		Time:   time.Now().UTC(),
	}, false, log.NewNopLogger())
	
	// Initialize params
	err := k.SetParams(ctx, types.DefaultParams())
	require.NoError(t, err)
	
	// Set module account
	accountKeeper.modules[types.ModuleName] = authtypes.NewModuleAddress(types.ModuleName)
	
	return ctx, k, bankKeeper, stakingKeeper, accountKeeper
}