package keeper

import (
	"github.com/rollchains/flora/x/liquidstaking/types"
)

// SetHooks sets the liquid staking hooks
// This should be called only once during app initialization
func (k *Keeper) SetHooks(hooks types.LiquidStakingHooks) {
	if k.hooks != nil {
		panic("liquid staking hooks already set")
	}
	k.hooks = hooks
}

// GetHooks returns the liquid staking hooks
// If no hooks are set, it returns a no-op implementation
func (k Keeper) GetHooks() types.LiquidStakingHooks {
	if k.hooks == nil {
		return types.NoOpLiquidStakingHooks{}
	}
	return k.hooks
}