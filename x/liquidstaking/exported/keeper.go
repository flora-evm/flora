package exported

import (
	"cosmossdk.io/log"
	"cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper defines the expected keeper interface for migrations
type Keeper interface {
	// StoreKey returns the store key for the module
	StoreKey() types.StoreKey
	
	// Codec returns the codec for the module
	Codec() codec.BinaryCodec
	
	// Logger returns a module-specific logger
	Logger(ctx sdk.Context) log.Logger
}