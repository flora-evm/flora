package v2liquidstaking

import (
	"fmt"

	"cosmossdk.io/math"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	
	liquidstakingtypes "github.com/rollchains/flora/x/liquidstaking/types"
)

const (
	// UpgradeName defines the on-chain upgrade name for the liquid staking module addition
	UpgradeName = "v2-liquid-staking"
)

// CreateUpgradeHandler creates an upgrade handler for the v2 liquid staking upgrade
func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	liquidStakingKeeper interface{
		InitGenesis(ctx sdk.Context, genState liquidstakingtypes.GenesisState)
	},
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		ctx.Logger().Info("executing liquid staking module upgrade", "name", UpgradeName)
		
		// Initialize liquid staking module with default genesis state
		liquidStakingGenesis := getLiquidStakingGenesisState()
		liquidStakingKeeper.InitGenesis(ctx, *liquidStakingGenesis)
		
		// Run module migrations
		vm, err := mm.RunMigrations(ctx, configurator, fromVM)
		if err != nil {
			return nil, fmt.Errorf("failed to run module migrations: %w", err)
		}
		
		ctx.Logger().Info("liquid staking module upgrade completed successfully")
		
		return vm, nil
	}
}

// getLiquidStakingGenesisState returns the genesis state for liquid staking module
func getLiquidStakingGenesisState() *liquidstakingtypes.GenesisState {
	return &liquidstakingtypes.GenesisState{
		Params: liquidstakingtypes.ModuleParams{
			// Core functionality - start enabled
			Enabled: true,
			
			// Conservative caps to start
			GlobalLiquidStakingCap: math.LegacyNewDecWithPrec(25, 2), // 25%
			ValidatorLiquidCap:     math.LegacyNewDecWithPrec(50, 2), // 50%
			
			// Minimum stake amount (1 FLORA)
			MinLiquidStakeAmount: math.NewInt(1_000_000_000_000_000_000),
			
			// Rate limiting - 24 hour period
			RateLimitPeriodHours: 24,
			
			// Conservative daily limits
			GlobalDailyTokenizationPercent:    math.LegacyNewDecWithPrec(10, 2), // 10%
			ValidatorDailyTokenizationPercent: math.LegacyNewDecWithPrec(10, 2), // 10%
			GlobalDailyTokenizationCount:      100,
			ValidatorDailyTokenizationCount:   10,
			UserDailyTokenizationCount:        5,
			
			// Warning at 90% of cap
			WarningThresholdPercent: math.LegacyNewDecWithPrec(90, 2),
			
			// Auto-compound disabled by default
			AutoCompoundEnabled:         false,
			AutoCompoundFrequencyBlocks: 28800, // ~24 hours at 3s blocks
			
			// Conservative rate change limits
			MaxRateChangePerUpdate:  math.LegacyNewDecWithPrec(1, 2), // 1%
			MinBlocksBetweenUpdates: 100, // ~5 minutes at 3s blocks
		},
		TokenizationRecords:      []liquidstakingtypes.TokenizationRecord{},
		LastTokenizationRecordId: 0,
		TotalLiquidStaked:        math.ZeroInt(),
	}
}

// StoreUpgrades defines the store upgrades for the v2 liquid staking upgrade
func StoreUpgrades() *upgradetypes.StoreUpgrades {
	return &upgradetypes.StoreUpgrades{
		Added: []string{
			liquidstakingtypes.ModuleName,
		},
	}
}