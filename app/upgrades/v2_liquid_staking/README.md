# V2 Liquid Staking Upgrade

This upgrade adds the liquid staking module to an existing Flora chain.

## Registering the Upgrade

In your app's `RegisterUpgradeHandlers()` function (usually in `app/upgrades.go`), add:

```go
import (
    v2liquidstaking "github.com/rollchains/flora/app/upgrades/v2_liquid_staking"
)

// RegisterUpgradeHandlers registers all upgrade handlers
func (app *ChainApp) RegisterUpgradeHandlers() {
    // ... existing upgrade handlers ...
    
    // Add v2 liquid staking upgrade
    app.UpgradeKeeper.SetUpgradeHandler(
        v2liquidstaking.UpgradeName,
        v2liquidstaking.CreateUpgradeHandler(
            app.ModuleManager,
            app.configurator,
            app.LiquidStakingKeeper,
        ),
    )
    
    // When this upgrade is planned, set the store upgrades
    upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
    if err != nil {
        panic(fmt.Sprintf("failed to read upgrade info from disk %s", err))
    }
    
    if upgradeInfo.Name == v2liquidstaking.UpgradeName && !app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
        storeUpgrades := v2liquidstaking.StoreUpgrades()
        
        // Configure store loader with the store upgrades
        app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, storeUpgrades))
    }
}
```

## Testing the Upgrade

### Local Testing

1. Start a local chain without the liquid staking module
2. Submit an upgrade proposal:
```bash
florad tx gov submit-proposal software-upgrade v2-liquid-staking \
  --title "Add Liquid Staking Module" \
  --description "This upgrade adds the liquid staking module to Flora" \
  --deposit 10000000flora \
  --upgrade-height 100 \
  --from validator
```

3. Vote on the proposal:
```bash
florad tx gov vote 1 yes --from validator
```

4. Wait for the upgrade height
5. The chain will halt at the upgrade height
6. Restart the chain with the new binary containing the liquid staking module

### Testnet Deployment

1. Coordinate with validators to prepare the new binary
2. Submit the upgrade proposal with appropriate height (usually 7-14 days out)
3. Ensure 2/3+ voting power approves the proposal
4. Validators must update their binaries before the upgrade height
5. The chain will automatically resume after the upgrade

## Post-Upgrade Verification

After the upgrade completes:

```bash
# Check module is loaded
florad query upgrade module-versions | grep liquidstaking

# Query module parameters
florad query liquidstaking params

# Verify module is functional
florad tx liquidstaking tokenize-shares 1000000stake [validator-address] --from mykey
```

## Rollback Plan

If issues occur:
1. Validators can coordinate to restart from the last snapshot before upgrade
2. Fix any issues in the upgrade handler
3. Propose a new upgrade with fixes

## Parameters After Upgrade

The upgrade initializes the module with conservative default parameters:
- Module enabled: true
- Global liquid cap: 25%
- Validator liquid cap: 50%
- Daily limits: 10% and max 100 global operations
- Auto-compound: disabled

These can be adjusted via governance after the upgrade.