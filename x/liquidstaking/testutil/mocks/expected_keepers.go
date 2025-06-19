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
	TotalBondedTokensFn  func(ctx context.Context) (math.Int, error)
	DelegateFn           func(ctx context.Context, delAddr sdk.AccAddress, bondAmt math.Int, tokenSrc stakingtypes.BondStatus, validator stakingtypes.Validator, subtractAccount bool) (math.LegacyDec, error)
	IterateValidatorsFn  func(ctx context.Context, fn func(index int64, validator stakingtypes.ValidatorI) (stop bool))
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

func (m *MockStakingKeeper) TotalBondedTokens(ctx context.Context) (math.Int, error) {
	if m.TotalBondedTokensFn != nil {
		return m.TotalBondedTokensFn(ctx)
	}
	return math.NewInt(1000000), nil
}

func (m *MockStakingKeeper) Delegate(ctx context.Context, delAddr sdk.AccAddress, bondAmt math.Int, tokenSrc stakingtypes.BondStatus, validator stakingtypes.Validator, subtractAccount bool) (math.LegacyDec, error) {
	if m.DelegateFn != nil {
		return m.DelegateFn(ctx, delAddr, bondAmt, tokenSrc, validator, subtractAccount)
	}
	// Default implementation: return shares equal to tokens (1:1 ratio)
	return math.LegacyNewDecFromInt(bondAmt), nil
}

func (m *MockStakingKeeper) IterateValidators(ctx context.Context, fn func(index int64, validator stakingtypes.ValidatorI) (stop bool)) {
	if m.IterateValidatorsFn != nil {
		m.IterateValidatorsFn(ctx, fn)
	}
}

// MockBankKeeper is a mock implementation of the BankKeeper interface
type MockBankKeeper struct {
	MintCoinsFn                      func(ctx context.Context, moduleName string, amt sdk.Coins) error
	SendCoinsFromModuleToAccountFn   func(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	GetDenomMetaDataFn               func(ctx context.Context, denom string) (banktypes.Metadata, bool)
	SetDenomMetaDataFn               func(ctx context.Context, denomMetaData banktypes.Metadata)
	GetSupplyFn                      func(ctx context.Context, denom string) sdk.Coin
	GetBalanceFn                     func(ctx context.Context, addr sdk.AccAddress, denom string) sdk.Coin
	SendCoinsFromAccountToModuleFn   func(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	BurnCoinsFn                      func(ctx context.Context, moduleName string, amt sdk.Coins) error
	SendCoinsFn                      func(ctx context.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
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

func (m *MockBankKeeper) GetBalance(ctx context.Context, addr sdk.AccAddress, denom string) sdk.Coin {
	if m.GetBalanceFn != nil {
		return m.GetBalanceFn(ctx, addr, denom)
	}
	return sdk.NewCoin(denom, math.ZeroInt())
}

func (m *MockBankKeeper) SendCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error {
	if m.SendCoinsFromAccountToModuleFn != nil {
		return m.SendCoinsFromAccountToModuleFn(ctx, senderAddr, recipientModule, amt)
	}
	return nil
}

func (m *MockBankKeeper) BurnCoins(ctx context.Context, moduleName string, amt sdk.Coins) error {
	if m.BurnCoinsFn != nil {
		return m.BurnCoinsFn(ctx, moduleName, amt)
	}
	return nil
}

func (m *MockBankKeeper) SendCoins(ctx context.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error {
	if m.SendCoinsFn != nil {
		return m.SendCoinsFn(ctx, fromAddr, toAddr, amt)
	}
	return nil
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

// MockTokenFactoryKeeper is no longer needed - using Bank module directly
// The liquid staking module uses the Bank module for minting and burning LSTs
// instead of the Token Factory module

// MockDistributionKeeper is a mock implementation of the DistributionKeeper interface
type MockDistributionKeeper struct {
	GetValidatorAccumulatedRewardsFn func(ctx context.Context, val sdk.ValAddress) (sdk.DecCoins, error)
}

func (m *MockDistributionKeeper) GetValidatorAccumulatedRewards(ctx context.Context, val sdk.ValAddress) (sdk.DecCoins, error) {
	if m.GetValidatorAccumulatedRewardsFn != nil {
		return m.GetValidatorAccumulatedRewardsFn(ctx, val)
	}
	return sdk.DecCoins{}, nil
}