package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	
	"github.com/rollchains/flora/x/liquidstaking/types"
)

// GetParams returns the module parameters
func (k Keeper) GetParams(ctx sdk.Context) (params types.ModuleParams) {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := store.Get(types.ParamsKey)
	if err != nil {
		panic(err)
	}
	if bz == nil {
		return types.DefaultParams()
	}
	
	k.cdc.MustUnmarshal(bz, &params)
	return params
}

// SetParams sets the module parameters
func (k Keeper) SetParams(ctx sdk.Context, params types.ModuleParams) error {
	if err := params.Validate(); err != nil {
		return err
	}
	
	store := k.storeService.OpenKVStore(ctx)
	bz := k.cdc.MustMarshal(&params)
	return store.Set(types.ParamsKey, bz)
}