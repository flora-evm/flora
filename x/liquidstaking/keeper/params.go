package keeper

import (
	"fmt"
	
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
	
	// Get old params for comparison
	oldParams := k.GetParams(ctx)
	
	store := k.storeService.OpenKVStore(ctx)
	bz := k.cdc.MustMarshal(&params)
	if err := store.Set(types.ParamsKey, bz); err != nil {
		return err
	}
	
	// Emit parameter update events
	k.emitParamUpdateEvents(ctx, oldParams, params)
	
	return nil
}

// emitParamUpdateEvents emits events for parameter changes
func (k Keeper) emitParamUpdateEvents(ctx sdk.Context, oldParams, newParams types.ModuleParams) {
	var changes []types.EventParamChange
	
	// Check each parameter for changes
	if oldParams.Enabled != newParams.Enabled {
		changes = append(changes, types.EventParamChange{
			Key:      "enabled",
			OldValue: fmt.Sprintf("%v", oldParams.Enabled),
			NewValue: fmt.Sprintf("%v", newParams.Enabled),
		})
	}
	
	// TODO: Uncomment after proto regeneration
	// if !oldParams.MinLiquidStakeAmount.Equal(newParams.MinLiquidStakeAmount) {
	// 	changes = append(changes, types.ParamChange{
	// 		Key:      "min_liquid_stake_amount",
	// 		OldValue: oldParams.MinLiquidStakeAmount.String(),
	// 		NewValue: newParams.MinLiquidStakeAmount.String(),
	// 	})
	// }
	
	if !oldParams.GlobalLiquidStakingCap.Equal(newParams.GlobalLiquidStakingCap) {
		changes = append(changes, types.EventParamChange{
			Key:      "global_liquid_staking_cap",
			OldValue: oldParams.GlobalLiquidStakingCap.String(),
			NewValue: newParams.GlobalLiquidStakingCap.String(),
		})
	}
	
	if !oldParams.ValidatorLiquidCap.Equal(newParams.ValidatorLiquidCap) {
		changes = append(changes, types.EventParamChange{
			Key:      "validator_liquid_cap",
			OldValue: oldParams.ValidatorLiquidCap.String(),
			NewValue: newParams.ValidatorLiquidCap.String(),
		})
	}
	
	// Only emit event if there were actual changes
	if len(changes) > 0 {
		event := types.UpdateParamsEvent{
			Authority: "governance", // Will be updated when governance integration is added
			Changes:   changes,
		}
		ctx.EventManager().EmitEvent(event.ToEvent())
	}
}