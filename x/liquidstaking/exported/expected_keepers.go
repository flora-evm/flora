package exported

import (
	"context"
	
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"
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
	TotalBondedTokens(ctx context.Context) (math.Int, error)
	
	// Delegate performs a delegation from a delegator to a validator
	Delegate(ctx context.Context, delAddr sdk.AccAddress, bondAmt math.Int, tokenSrc stakingtypes.BondStatus, validator stakingtypes.Validator, subtractAccount bool) (math.LegacyDec, error)
	
	// IterateValidators iterates through all validators
	IterateValidators(ctx context.Context, fn func(index int64, validator stakingtypes.ValidatorI) (stop bool)) error
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
	
	// SendCoins sends coins from one account to another
	SendCoins(ctx context.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
}

// AccountKeeper defines the expected interface for the auth module
type AccountKeeper interface {
	// GetAccount returns an account
	GetAccount(ctx context.Context, addr sdk.AccAddress) sdk.AccountI
	
	// GetModuleAddress returns the address of a module account
	GetModuleAddress(moduleName string) sdk.AccAddress
}

// TransferKeeper defines the expected IBC transfer keeper interface
type TransferKeeper interface {
	// SendTransfer sends IBC tokens from sender to receiver
	// Note: In the SDK, this might be lowercase sendTransfer
	SendTransfer(
		ctx sdk.Context,
		sourcePort,
		sourceChannel string,
		token sdk.Coin,
		sender sdk.AccAddress,
		receiver string,
		timeoutHeight clienttypes.Height,
		timeoutTimestamp uint64,
		memo string,
	) (uint64, error)
}

// ChannelKeeper defines the expected IBC channel keeper interface  
type ChannelKeeper interface {
	// GetChannel returns a channel
	GetChannel(ctx sdk.Context, portID, channelID string) (channeltypes.Channel, bool)
	
	// GetChannelClientState returns the client state for a channel
	GetChannelClientState(ctx sdk.Context, portID, channelID string) (string, ibcexported.ClientState, error)
}

// DistributionKeeper defines the expected distribution keeper interface
type DistributionKeeper interface {
	// TODO: Add methods as needed for liquid staking functionality
	// For now, we don't require any specific distribution keeper methods
}

// TokenFactoryKeeper is currently not used in the liquid staking module
// The module uses the Bank module directly for minting and burning LSTs
// This interface is kept for potential future integration
//
// type TokenFactoryKeeper interface {
// 	// CreateDenom creates a new denom with the given subdenom controlled by the creator
// 	CreateDenom(ctx context.Context, creator string, subdenom string) (string, error)
// 	
// 	// Mint mints tokens of a given denom to an account
// 	Mint(ctx context.Context, mintToAddr string, coin sdk.Coin) error
// 	
// 	// Burn burns tokens from an account
// 	Burn(ctx context.Context, burnFromAddr string, coin sdk.Coin) error
// 	
// 	// SetDenomMetadata sets the metadata for a denom
// 	SetDenomMetadata(ctx context.Context, creator string, metadata banktypes.Metadata) error
// 	
// 	// GetDenomAuthorityMetadata returns the authority metadata for a specific denom
// 	GetDenomAuthorityMetadata(ctx context.Context, denom string) (DenomAuthorityMetadata, error)
// }
//
// // DenomAuthorityMetadata specifies metadata for a denom authority
// type DenomAuthorityMetadata struct {
// 	// Admin is the address that can perform admin operations
// 	Admin string
// }