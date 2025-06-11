package exported

import (
	"context"
	
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// StakingKeeper defines the expected staking keeper interface
type StakingKeeper interface {
	// GetDelegation returns a specific delegation
	GetDelegation(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (stakingtypes.Delegation, error)
	
	// GetValidator returns a specific validator
	GetValidator(ctx context.Context, addr sdk.ValAddress) (stakingtypes.Validator, error)
	
	// Unbond a delegation
	Unbond(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, shares math.LegacyDec) (math.Int, error)
	
	// GetParams returns the staking module parameters
	GetParams(ctx context.Context) (stakingtypes.Params, error)
	
	// BondDenom returns the denomination of the staking token
	BondDenom(ctx context.Context) (string, error)
	
	// TotalBondedTokens returns the total amount of bonded tokens
	TotalBondedTokens(ctx context.Context) math.Int
	
	// Delegate performs a delegation from a delegator to a validator
	Delegate(ctx context.Context, delAddr sdk.AccAddress, bondAmt math.Int, tokenSrc stakingtypes.BondStatus, validator stakingtypes.Validator, subtractAccount bool) (math.LegacyDec, error)
}

// BankKeeper defines the expected bank keeper interface
type BankKeeper interface {
	// MintCoins mints new coins
	MintCoins(ctx context.Context, moduleName string, amt sdk.Coins) error
	
	// SendCoinsFromModuleToAccount sends coins from a module account to a user account
	SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	
	// GetDenomMetaData returns the metadata of a denom
	GetDenomMetaData(ctx context.Context, denom string) (banktypes.Metadata, bool)
	
	// SetDenomMetaData sets the metadata of a denom
	SetDenomMetaData(ctx context.Context, denomMetaData banktypes.Metadata)
	
	// GetSupply returns the supply of a denom
	GetSupply(ctx context.Context, denom string) sdk.Coin
	
	// GetBalance returns the balance of a specific denomination for an account
	GetBalance(ctx context.Context, addr sdk.AccAddress, denom string) sdk.Coin
	
	// SendCoinsFromAccountToModule sends coins from a user account to a module account
	SendCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	
	// BurnCoins burns coins from a module account
	BurnCoins(ctx context.Context, moduleName string, amt sdk.Coins) error
}

// AccountKeeper defines the expected interface for the auth module
type AccountKeeper interface {
	// GetAccount returns an account
	GetAccount(ctx context.Context, addr sdk.AccAddress) sdk.AccountI
	
	// GetModuleAddress returns the address of a module account
	GetModuleAddress(moduleName string) sdk.AccAddress
}