package mocks

import (
	"context"
	
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// MockStakingKeeper is a mock implementation of the StakingKeeper interface
type MockStakingKeeper struct {
	GetDelegationFn      func(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (stakingtypes.Delegation, error)
	GetValidatorFn       func(ctx context.Context, addr sdk.ValAddress) (stakingtypes.Validator, error)
	UnbondFn             func(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, shares math.LegacyDec) (math.Int, error)
	GetParamsFn          func(ctx context.Context) (stakingtypes.Params, error)
	BondDenomFn          func(ctx context.Context) (string, error)
	TotalBondedTokensFn  func(ctx context.Context) math.Int
}

func (m *MockStakingKeeper) GetDelegation(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (stakingtypes.Delegation, error) {
	if m.GetDelegationFn != nil {
		return m.GetDelegationFn(ctx, delAddr, valAddr)
	}
	return stakingtypes.Delegation{}, nil
}

func (m *MockStakingKeeper) GetValidator(ctx context.Context, addr sdk.ValAddress) (stakingtypes.Validator, error) {
	if m.GetValidatorFn != nil {
		return m.GetValidatorFn(ctx, addr)
	}
	return stakingtypes.Validator{}, nil
}

func (m *MockStakingKeeper) Unbond(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, shares math.LegacyDec) (math.Int, error) {
	if m.UnbondFn != nil {
		return m.UnbondFn(ctx, delAddr, valAddr, shares)
	}
	return math.ZeroInt(), nil
}

func (m *MockStakingKeeper) GetParams(ctx context.Context) (stakingtypes.Params, error) {
	if m.GetParamsFn != nil {
		return m.GetParamsFn(ctx)
	}
	return stakingtypes.Params{}, nil
}

func (m *MockStakingKeeper) BondDenom(ctx context.Context) (string, error) {
	if m.BondDenomFn != nil {
		return m.BondDenomFn(ctx)
	}
	return "stake", nil
}

func (m *MockStakingKeeper) TotalBondedTokens(ctx context.Context) math.Int {
	if m.TotalBondedTokensFn != nil {
		return m.TotalBondedTokensFn(ctx)
	}
	return math.NewInt(1000000)
}

// MockBankKeeper is a mock implementation of the BankKeeper interface
type MockBankKeeper struct {
	MintCoinsFn                      func(ctx context.Context, moduleName string, amt sdk.Coins) error
	SendCoinsFromModuleToAccountFn   func(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	GetDenomMetaDataFn               func(ctx context.Context, denom string) (banktypes.Metadata, bool)
	SetDenomMetaDataFn               func(ctx context.Context, denomMetaData banktypes.Metadata)
	GetSupplyFn                      func(ctx context.Context, denom string) sdk.Coin
}

func (m *MockBankKeeper) MintCoins(ctx context.Context, moduleName string, amt sdk.Coins) error {
	if m.MintCoinsFn != nil {
		return m.MintCoinsFn(ctx, moduleName, amt)
	}
	return nil
}

func (m *MockBankKeeper) SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error {
	if m.SendCoinsFromModuleToAccountFn != nil {
		return m.SendCoinsFromModuleToAccountFn(ctx, senderModule, recipientAddr, amt)
	}
	return nil
}

func (m *MockBankKeeper) GetDenomMetaData(ctx context.Context, denom string) (banktypes.Metadata, bool) {
	if m.GetDenomMetaDataFn != nil {
		return m.GetDenomMetaDataFn(ctx, denom)
	}
	return banktypes.Metadata{}, false
}

func (m *MockBankKeeper) SetDenomMetaData(ctx context.Context, denomMetaData banktypes.Metadata) {
	if m.SetDenomMetaDataFn != nil {
		m.SetDenomMetaDataFn(ctx, denomMetaData)
	}
}

func (m *MockBankKeeper) GetSupply(ctx context.Context, denom string) sdk.Coin {
	if m.GetSupplyFn != nil {
		return m.GetSupplyFn(ctx, denom)
	}
	return sdk.NewCoin(denom, math.ZeroInt())
}

// MockAccountKeeper is a mock implementation of the AccountKeeper interface
type MockAccountKeeper struct {
	GetAccountFn       func(ctx context.Context, addr sdk.AccAddress) sdk.AccountI
	GetModuleAddressFn func(moduleName string) sdk.AccAddress
}

func (m *MockAccountKeeper) GetAccount(ctx context.Context, addr sdk.AccAddress) sdk.AccountI {
	if m.GetAccountFn != nil {
		return m.GetAccountFn(ctx, addr)
	}
	return nil
}

func (m *MockAccountKeeper) GetModuleAddress(moduleName string) sdk.AccAddress {
	if m.GetModuleAddressFn != nil {
		return m.GetModuleAddressFn(moduleName)
	}
	return sdk.AccAddress([]byte(moduleName))
}