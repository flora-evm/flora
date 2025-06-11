package keeper

import (
	"fmt"

	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	
	"github.com/rollchains/flora/x/liquidstaking/types"
)

// Keeper of the liquid staking module
type Keeper struct {
	storeService store.KVStoreService
	cdc          codec.BinaryCodec
	
	stakingKeeper types.StakingKeeper
	bankKeeper    types.BankKeeper
	accountKeeper types.AccountKeeper
}

// NewKeeper creates a new liquid staking Keeper instance
func NewKeeper(
	storeService store.KVStoreService,
	cdc codec.BinaryCodec,
	sk types.StakingKeeper,
	bk types.BankKeeper,
	ak types.AccountKeeper,
) Keeper {
	return Keeper{
		storeService:  storeService,
		cdc:           cdc,
		stakingKeeper: sk,
		bankKeeper:    bk,
		accountKeeper: ak,
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