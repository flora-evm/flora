package keeper

import (
	"fmt"

	"cosmossdk.io/core/store"
	storetypes "cosmossdk.io/store/types"
	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	
	"github.com/rollchains/flora/x/liquidstaking/types"
)

// Keeper of the liquid staking module
type Keeper struct {
	storeService store.KVStoreService
	cdc          codec.BinaryCodec
	storeKey     storetypes.StoreKey
	authority    string
	
	stakingKeeper      types.StakingKeeper
	bankKeeper         types.BankKeeper
	accountKeeper      types.AccountKeeper
	// transferKeeper     types.TransferKeeper // TODO: Fix interface mismatch
	// channelKeeper      types.ChannelKeeper // TODO: Fix interface mismatch
	distributionKeeper types.DistributionKeeper
	// tokenFactoryKeeper removed - using Bank module directly for LST management
	
	hooks types.LiquidStakingHooks
}

// NewKeeper creates a new liquid staking Keeper instance
func NewKeeper(
	storeService store.KVStoreService,
	cdc codec.BinaryCodec,
	sk types.StakingKeeper,
	bk types.BankKeeper,
	ak types.AccountKeeper,
	// tk types.TransferKeeper, // TODO: Fix interface mismatch
	// ck types.ChannelKeeper, // TODO: Fix interface mismatch
	dk types.DistributionKeeper,
	authority string,
) Keeper {
	// For backward compatibility, we derive the store key from the store service
	// This matches the app.go initialization pattern
	return Keeper{
		storeService:       storeService,
		cdc:                cdc,
		storeKey:           nil, // Not needed with KVStoreService
		authority:          authority,
		stakingKeeper:      sk,
		bankKeeper:         bk,
		accountKeeper:      ak,
		// transferKeeper:     tk, // TODO: Fix interface mismatch
		// channelKeeper:      ck, // TODO: Fix interface mismatch
		distributionKeeper: dk,
	}
}

// GetStoreService returns the store service
func (k Keeper) GetStoreService() store.KVStoreService {
	return k.storeService
}

// Logger returns a module-specific logger
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetAuthority returns the module's authority address
func (k Keeper) GetAuthority() string {
	return k.authority
}

// StoreKey returns the store key for migrations
func (k Keeper) StoreKey() storetypes.StoreKey {
	return k.storeKey
}

// Codec returns the codec for migrations
func (k Keeper) Codec() codec.BinaryCodec {
	return k.cdc
}

// ValidateModuleEnabled checks if the liquid staking module is enabled in params.
// Returns ErrModuleDisabled if the module is disabled.
func (k Keeper) ValidateModuleEnabled(ctx sdk.Context) error {
	params := k.GetParams(ctx)
	if !params.Enabled {
		return types.ErrModuleDisabled
	}
	return nil
}

// ParseAndValidateAddress parses a bech32 address string and validates it.
// The addrType parameter is used in error messages for clarity (e.g., "delegator", "owner").
// Returns the parsed address or an error if the address is empty or invalid.
func ParseAndValidateAddress(addrStr string, addrType string) (sdk.AccAddress, error) {
	if addrStr == "" {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("%s address cannot be empty", addrType)
	}
	
	addr, err := sdk.AccAddressFromBech32(addrStr)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid %s address: %s", addrType, err)
	}
	
	return addr, nil
}

// ParseAndValidateValidatorAddress parses a bech32 validator address string and validates it.
// Returns the parsed validator address or an error if the address is empty or invalid.
func ParseAndValidateValidatorAddress(addrStr string) (sdk.ValAddress, error) {
	if addrStr == "" {
		return nil, sdkerrors.ErrInvalidAddress.Wrap("validator address cannot be empty")
	}
	
	addr, err := sdk.ValAddressFromBech32(addrStr)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid validator address: %s", err)
	}
	
	return addr, nil
}

// ValidatePositiveAmount validates that a coin amount is valid and positive.
// Returns an error if the amount is invalid, zero, or negative.
// This is commonly used for validating tokenization and redemption amounts.
func ValidatePositiveAmount(amount sdk.Coin) error {
	if !amount.IsValid() {
		return sdkerrors.ErrInvalidRequest.Wrap("invalid amount")
	}
	
	if amount.IsZero() {
		return sdkerrors.ErrInvalidRequest.Wrap("amount cannot be zero")
	}
	
	if amount.IsNegative() {
		return sdkerrors.ErrInvalidRequest.Wrap("amount cannot be negative")
	}
	
	return nil
}