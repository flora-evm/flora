# Liquid Staking Module Integration Guide

This guide provides step-by-step instructions for integrating the liquid staking module into your Flora application.

## Prerequisites

- Flora blockchain with Cosmos SDK v0.50.13
- Completed liquid staking module implementation
- Access to modify app/app.go and related files

## Integration Steps

### 1. Import the Module

Add the liquid staking module import to your `app/app.go`:

```go
import (
    // ... existing imports ...
    
    liquidstaking "github.com/rollchains/flora/x/liquidstaking"
    liquidstakingkeeper "github.com/rollchains/flora/x/liquidstaking/keeper"
    liquidstakingtypes "github.com/rollchains/flora/x/liquidstaking/types"
)
```

### 2. Add Keeper to App Structure

Add the liquid staking keeper to your App struct:

```go
type App struct {
    // ... existing fields ...
    
    // Liquid staking keeper
    LiquidStakingKeeper liquidstakingkeeper.Keeper
}
```

### 3. Add Store Key

Add the store key to your application:

```go
// In the store key definition section
keys := storetypes.NewKVStoreKeys(
    // ... existing store keys ...
    liquidstakingtypes.StoreKey,
)
```

### 4. Initialize the Keeper

In the keeper initialization section of NewApp:

```go
// After staking keeper initialization
app.LiquidStakingKeeper = liquidstakingkeeper.NewKeeper(
    appCodec,
    runtime.NewKVStoreService(keys[liquidstakingtypes.StoreKey]),
    app.AccountKeeper,
    app.BankKeeper,
    app.StakingKeeper,
    app.DistrKeeper,
    authtypes.NewModuleAddress(govtypes.ModuleName).String(), // authority
)
```

### 5. Set Staking Hooks

After creating the liquid staking keeper, set it as a staking hook:

```go
// Set the staking hooks
// NOTE: staking module is required if HistoricalEntries param > 0
// NOTE: capability module's hooks are already set in the capability module's InitGenesis
app.StakingKeeper.SetHooks(
    stakingtypes.NewMultiStakingHooks(
        app.DistrKeeper.Hooks(),
        app.SlashingKeeper.Hooks(),
        app.LiquidStakingKeeper.Hooks(), // Add liquid staking hooks
    ),
)
```

### 6. Create the Module

Create the liquid staking module:

```go
liquidStakingModule := liquidstaking.NewAppModule(
    appCodec,
    app.LiquidStakingKeeper,
    app.AccountKeeper,
    app.BankKeeper,
    app.StakingKeeper,
    app.DistrKeeper,
)
```

### 7. Add to Module Manager

Add the module to the module manager:

```go
app.ModuleManager = module.NewManager(
    // ... existing modules ...
    liquidStakingModule,
)
```

### 8. Set Module Order

Add liquid staking to the module execution order:

```go
// During begin block, run in this order:
// - mint new tokens for the previous block
// - distribute rewards for the previous block
// - slash/jail validators for downtime
// - liquid staking auto-compound and rate updates
app.ModuleManager.SetOrderBeginBlockers(
    // ... existing modules ...
    liquidstakingtypes.ModuleName,
)

// During end block, keep existing order
app.ModuleManager.SetOrderEndBlockers(
    // ... existing modules ...
    liquidstakingtypes.ModuleName,
)

// Sets the order of Genesis - Order matters, genutil is to always come last
app.ModuleManager.SetOrderInitGenesis(
    // ... existing modules ...
    liquidstakingtypes.ModuleName,
    // ... genutil last ...
)
```

### 9. Register Governance Proposal Handlers

In the governance keeper initialization or where proposal handlers are registered:

```go
// Add liquid staking proposal handler
govRouter := govv1beta1.NewRouter()
govRouter.AddRoute(govtypes.RouterKey, govv1beta1.ProposalHandler).
    // ... other routes ...
    AddRoute(liquidstakingtypes.RouterKey, liquidstakingkeeper.NewProposalHandler(app.LiquidStakingKeeper))

app.GovKeeper = govkeeper.NewKeeper(
    // ... existing parameters ...
)
app.GovKeeper.SetLegacyRouter(govRouter)
```

### 10. Register Services

The module services are automatically registered when creating the module. No additional registration needed.

### 11. Add to Export/Import

In the ExportAppStateAndValidators function:

```go
// Add liquid staking to the module list that exports genesis
```

### 12. Configure Ante Handlers (Optional)

If you need to add custom ante handlers for liquid staking:

```go
// In NewAnteHandler or where ante handlers are configured
// Add any liquid staking specific ante decorators if needed
```

## Genesis Configuration

### Default Genesis State

Create a default genesis state for the module:

```json
{
  "liquidstaking": {
    "params": {
      "enabled": true,
      "global_liquid_staking_cap": "0.25",
      "validator_liquid_cap": "0.5",
      "min_liquid_stake_amount": "1000000",
      "rate_limit_period_hours": 24,
      "global_daily_tokenization_percent": "0.1",
      "validator_daily_tokenization_percent": "0.1",
      "global_daily_tokenization_count": "100",
      "validator_daily_tokenization_count": "10",
      "user_daily_tokenization_count": "5",
      "warning_threshold_percent": "0.9",
      "auto_compound_enabled": false,
      "auto_compound_frequency_blocks": "28800",
      "max_rate_change_per_update": "0.01",
      "min_blocks_between_updates": "100"
    },
    "tokenization_records": [],
    "last_tokenization_record_id": "0",
    "total_liquid_staked": "0"
  }
}
```

## Testing the Integration

### 1. Build and Run

```bash
# Build the application
make install

# Initialize a test chain
florad init test-chain --chain-id test-1

# Add genesis accounts and validators
florad keys add validator
florad genesis add-genesis-account validator 100000000stake
florad genesis gentx validator 50000000stake

# Start the chain
florad start
```

### 2. Test Basic Operations

```bash
# Check module is loaded
florad query liquidstaking params

# Try tokenizing shares (after delegating)
florad tx staking delegate [validator-address] 10000000stake --from validator
florad tx liquidstaking tokenize-shares 5000000stake [validator-address] --from validator

# Query tokenization records
florad query liquidstaking tokenization-records
```

### 3. Test Governance

```bash
# Submit a parameter update proposal
florad tx gov submit-proposal update-params liquidstaking \
  "Update Params" "Enable auto-compound" 1000stake params.json \
  --from validator

# Vote on the proposal
florad tx gov vote 1 yes --from validator
```

## Troubleshooting

### Common Issues

1. **Module not found errors**
   - Ensure all imports are correct
   - Run `go mod tidy` to update dependencies

2. **Keeper initialization panics**
   - Verify all required keepers are passed
   - Check store key is properly registered

3. **Genesis validation failures**
   - Ensure genesis state follows the correct format
   - Validate parameters are within acceptable ranges

4. **Transaction failures**
   - Check module is enabled in params
   - Verify account has delegations before tokenizing
   - Ensure caps are not exceeded

### Debug Commands

```bash
# Check if module is registered
florad query upgrade module-versions

# Verify store is accessible
florad query liquidstaking params

# Check hooks are working
florad query liquidstaking exchange-rates
```

## Migration from Existing Chain

If adding to an existing chain, create an upgrade handler:

```go
app.UpgradeKeeper.SetUpgradeHandler(
    "v2-liquid-staking",
    func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
        // Initialize liquid staking module with default params
        liquidStakingGenesis := liquidstakingtypes.DefaultGenesis()
        app.LiquidStakingKeeper.InitGenesis(ctx, *liquidStakingGenesis)
        
        return app.ModuleManager.RunMigrations(ctx, app.configurator, fromVM)
    },
)
```

## Security Checklist

Before mainnet deployment:

- [ ] Set appropriate parameter values for your network
- [ ] Configure authority address for emergency controls  
- [ ] Test all governance proposals
- [ ] Verify caps and rate limits
- [ ] Run security audit on integration
- [ ] Test upgrade process on testnet
- [ ] Document emergency procedures

## Next Steps

1. Complete integration following this guide
2. Run comprehensive tests on testnet
3. Perform security audit
4. Plan mainnet upgrade proposal
5. Monitor post-deployment metrics

For additional support or questions, please refer to the module documentation or contact the development team.