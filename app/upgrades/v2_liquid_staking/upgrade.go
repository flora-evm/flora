package v2liquidstaking

import (
	"cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	
	"github.com/rollchains/flora/app/upgrades"
	liquidstakingtypes "github.com/rollchains/flora/x/liquidstaking/types"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: storetypes.StoreUpgrades{
		Added: []string{
			liquidstakingtypes.ModuleName,
		},
	},
}

func CreateUpgradeHandler(
	mm upgrades.ModuleManager,
	configurator module.Configurator,
	keepers *upgrades.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		ctx.Logger().Info("executing liquid staking module upgrade", "name", UpgradeName)
		
		// Since the liquid staking module is already integrated in app.go,
		// we just need to ensure it's properly initialized if this is the first time
		// The module manager will handle the initialization through InitGenesis
		
		// Run module migrations
		vm, err := mm.RunMigrations(ctx, configurator, fromVM)
		if err != nil {
			return nil, err
		}
		
		ctx.Logger().Info("liquid staking module upgrade completed successfully")
		
		return vm, nil
	}
}