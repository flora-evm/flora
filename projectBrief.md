# Flora Blockchain - Project Brief

## Executive Summary

Flora is a next-generation blockchain platform that combines the power of Ethereum Virtual Machine (EVM) compatibility with the robust infrastructure of the Cosmos SDK. This unique fusion enables developers to deploy Ethereum smart contracts while leveraging Inter-Blockchain Communication (IBC) for seamless cross-chain interactions.

## Core Value Proposition

Flora addresses the fragmentation in the blockchain ecosystem by providing:
- **Full EVM Compatibility**: Deploy existing Ethereum dApps without modification
- **Native IBC Support**: Connect to the entire Cosmos ecosystem
- **Unified Token Model**: No wrapped tokens - native assets work seamlessly in both environments
- **Enterprise-Ready Infrastructure**: Built on battle-tested Cosmos SDK with institutional-grade security

## Technical Architecture

### Foundation
- **Cosmos SDK**: v0.50.x (latest stable version)
- **EVM Integration**: cosmos/evm module (Apache 2.0 licensed)
- **Consensus**: CometBFT (formerly Tendermint)
- **IBC Protocol**: v8 with full transfer and interchain accounts support

### Key Specifications
- **Chain ID**: `localchain_9000-1` (EVM-compatible format)
- **Native Token**: `flora` (18 decimals for EVM compatibility)
- **Address Format**: Ethereum-style with `flora` Bech32 prefix support
- **Gas Model**: EIP-1559 dynamic fees with Cosmos gas metering

## Module Architecture

### Core Cosmos Modules
- **auth**: Account management and authentication
- **bank**: Token transfers and balance tracking
- **staking**: Proof-of-Stake consensus participation
- **gov**: On-chain governance proposals
- **distribution**: Reward distribution to validators/delegators
- **mint**: Inflation and token emission
- **slashing**: Validator penalty enforcement

### EVM Modules
- **evm**: Ethereum Virtual Machine execution environment
- **erc20**: Native token <-> ERC20 conversion
- **feemarket**: EIP-1559 dynamic fee adjustment

### IBC Modules
- **transfer**: Cross-chain token transfers
- **interchain-accounts**: Cross-chain account control
- **fee**: Relayer incentivization

### Custom Modules
- **tokenfactory**: Permissionless token creation and management
- **liquidstaking**: (Planned) Liquid staking derivatives

## Implementation Guide

### Step 1: Update Dependencies (go.mod)

```go
require (
    github.com/cosmos/cosmos-sdk v0.50.13
    github.com/ethereum/go-ethereum v1.10.26
    github.com/cosmos/ibc-go/modules/capability v1.0.1
    github.com/cosmos/ibc-go/v8 v8.7.0
)

replace (
    cosmossdk.io/store => github.com/cosmos/cosmos-sdk/store v1.1.2-0.20250108151001
    github.com/cosmos/cosmos-sdk => github.com/cosmos/cosmos-sdk v0.50.13-0.20250109122504
    github.com/ethereum/go-ethereum => github.com/cosmos/go-ethereum v1.10.26-cosmos-1
)
```

### Step 2: Chain Configuration

Update chain configuration for EVM compatibility:

```go
// app/app.go
const ChainID = "localchain_9000-1"

// Chain ID must follow format: {name}_{evm-id}-{version}
// Example: "mychain_9000-1" where 9000 is the EVM chain ID
```

### Step 3: Create EVM Configuration (app/config.go)

```go
package app

import (
    "fmt"
    sdk "github.com/cosmos/cosmos-sdk/types"
    evmtypes "github.com/cosmos/evm/x/vm/types"
)

func EVMOptionsFn(string) error {
    return nil
}

func NoOpEVMOptions(_ string) error {
    return nil
}

var sealed = false

var ChainsCoinInfo = map[string]evmtypes.EvmCoinInfo{
    ChainID: {
        Denom:          BaseDenom,
        DisplayDenom:   DisplayDenom,
        Decimals:       evmtypes.EighteenDecimals,
    },
}

func EVMAppOptions(chainID string) error {
    if sealed {
        return nil
    }
    
    if chainID == "" {
        chainID = ChainID
    }
    
    id := strings.Split(chainID, "-")[0]
    coinInfo, found := ChainsCoinInfo[id]
    if !found {
        coinInfo, found = ChainsCoinInfo[chainID]
        if !found {
            return fmt.Errorf("unknown chain id: %s, %+v", chainID, ChainsCoinInfo)
        }
    }
    
    if err := setBaseDenom(coinInfo); err != nil {
        return err
    }
    
    baseDenom, err := sdk.GetBaseDenom()
    if err != nil {
        return err
    }
    
    ethCfg := evmtypes.DefaultChainConfig()
    
    err = evmtypes.NewEVMConfigurator().
        WithChainConfig(ethCfg).
        WithEVMCoinInfo(baseDenom, uint8(coinInfo.Decimals)).
        Configure()
    
    if err != nil {
        return err
    }
    
    sealed = true
    return nil
}

func setBaseDenom(ci evmtypes.EvmCoinInfo) error {
    return sdk.RegisterDenom(ci.Denom, math.LegacyOneDec())
}
```

### Step 4: Token Pair Configuration (app/token_pair.go)

```go
package app

import erc20types "github.com/cosmos/evm/x/erc20/types"

const WTokenContractMainnet = "0xD0494964cD82660AaE99bEdC034a0deA8A0bd517"

var ExampleTokenPairs = []erc20types.TokenPair{
    {
        Erc20Address:  WTokenContractMainnet,
        Denom:         BaseDenom,
        Enabled:       true,
        ContractOwner: erc20types.OWNER_MODULE,
    },
}
```

### Step 5: Precompiles Configuration (app/precompiles.go)

```go
package app

import (
    "fmt"
    "maps"
    
    evidencekeeper "cosmossdk.io/x/evidence/keeper"
    authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
    bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
    distributionkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
    govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
    slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
    stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
    erc20Keeper "github.com/cosmos/evm/x/erc20/keeper"
    transferkeeper "github.com/cosmos/ibc-go/v8/modules/apps/transfer/keeper"
    channelkeeper "github.com/cosmos/ibc-go/v8/modules/core/04-channel/keeper"
    "github.com/ethereum/go-ethereum/common"
)

const bech32PrecompileBaseGas = 6_000

func NewAvailableStaticPrecompiles(
    stakingKeeper stakingkeeper.Keeper,
    distributionKeeper distributionkeeper.Keeper,
    bankKeeper bankkeeper.Keeper,
    erc20Keeper erc20Keeper.Keeper,
    authzKeeper authzkeeper.Keeper,
    transferKeeper transferkeeper.Keeper,
    channelKeeper channelkeeper.Keeper,
    evmKeeper *evmkeeper.Keeper,
    govKeeper govkeeper.Keeper,
    slashingKeeper slashingkeeper.Keeper,
    evidenceKeeper evidencekeeper.Keeper,
) map[common.Address]vm.PrecompiledContract {
    precompiles := maps.Clone(vm.PrecompiledContractsBerlin)
    
    p256Precompile := &p256.Precompile{}
    bech32Precompile, err := bech32.NewPrecompile(bech32PrecompileBaseGas)
    if err != nil {
        panic(fmt.Errorf("failed to instantiate bech32 precompile: %w", err))
    }
    
    stakingPrecompile, err := stakingprecompile.NewPrecompile(stakingKeeper, authzKeeper)
    if err != nil {
        panic(fmt.Errorf("failed to instantiate staking precompile: %w", err))
    }
    
    distributionPrecompile, err := distprecompile.NewPrecompile(
        distributionKeeper,
        stakingKeeper,
        authzKeeper,
        evmKeeper,
    )
    if err != nil {
        panic(fmt.Errorf("failed to instantiate distribution precompile: %w", err))
    }
    
    // Add all precompiles
    precompiles[bech32Precompile.Address()] = bech32Precompile
    precompiles[p256Precompile.Address()] = p256Precompile
    precompiles[distributionPrecompile.Address()] = distributionPrecompile
    precompiles[ibcTransferPrecompile.Address()] = ibcTransferPrecompile
    precompiles[bankPrecompile.Address()] = bankPrecompile
    precompiles[govPrecompile.Address()] = govPrecompile
    precompiles[slashingPrecompile.Address()] = slashingPrecompile
    precompiles[evidencePrecompile.Address()] = evidencePrecompile
    
    return precompiles
}
```

### Step 6: Update app.go Module Integration

Add EVM module imports and initialization:

```go
import (
    // Add EVM imports
    ante "github.com/cosmos/evm/ante"
    evmante "github.com/cosmos/evm/ante/evm"
    evmkeeper "github.com/cosmos/evm/x/vm/keeper"
    evmtypes "github.com/cosmos/evm/x/vm/types"
    erc20keeper "github.com/cosmos/evm/x/erc20/keeper"
    erc20types "github.com/cosmos/evm/x/erc20/types"
    feemarketkeeper "github.com/cosmos/evm/x/feemarket/keeper"
    feemarkettypes "github.com/cosmos/evm/x/feemarket/types"
    evm "github.com/cosmos/evm/x/vm"
    vmtypes "github.com/cosmos/evm/x/vm/core/vm"
    "github.com/cosmos/evm/x/vm/core/tracers/js"
    "github.com/cosmos/evm/x/vm/core/tracers/native"
    evmclient "github.com/cosmos/evm/x/vm/client"
)

// Add to module account permissions
var maccPerms = map[string][]string{
    // ... existing permissions ...
    evmtypes.ModuleName:       {authtypes.Minter, authtypes.Burner},
    erc20types.ModuleName:     {},
    feemarkettypes.ModuleName: {},
}

// Register tracers
func init() {
    js.RegisterNativeTracer("prestateTracer", tracers.NewPrestateTracer)
    js.RegisterNativeTracer("callTracer", tracers.NewCallTracer)
    js.RegisterJsTracer("opcountTracer", opcountTracer)
    
    DefaultPowerReduction = math.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))
}
```

### Step 7: Ante Handler Configuration

Create separate ante handlers for Cosmos and EVM transactions:

```go
// app/ante/ante_cosmos.go
func NewCosmosAnteHandler(options CosmosHandlerOptions) sdk.AnteHandler {
    return sdk.ChainAnteDecorators(
        cosmosante.NewSetUpContextDecorator(),
        cosmosante.NewExtensionOptionsDecorator(options.ExtensionOptionChecker),
        cosmosante.NewValidateBasicDecorator(),
        cosmosante.NewTxTimeoutHeightDecorator(),
        cosmosante.NewValidateMemoDecorator(options.AccountKeeper),
        cosmosante.NewConsumeGasForTxSizeDecorator(options.AccountKeeper),
        cosmosante.NewDeductFeeDecorator(options.AccountKeeper, options.BankKeeper, options.FeegrantKeeper, options.TxFeeChecker),
        cosmosante.NewSetPubKeyDecorator(options.AccountKeeper),
        cosmosante.NewValidateSigCountDecorator(options.AccountKeeper),
        cosmosante.NewSigGasConsumeDecorator(options.AccountKeeper, options.SigGasConsumer),
        cosmosante.NewSigVerificationDecorator(options.AccountKeeper, options.SignModeHandler),
        cosmosante.NewIncrementSequenceDecorator(options.AccountKeeper),
    )
}

// app/ante/ante_evm.go  
func NewEVMAnteHandler(options EVMHandlerOptions) sdk.AnteHandler {
    return sdk.ChainAnteDecorators(
        evmante.NewEthSetUpContextDecorator(options.EvmKeeper),
        evmante.NewEthMempoolFeeDecorator(options.EvmKeeper),
        evmante.NewEthMinGasPriceDecorator(options.FeeMarketKeeper, options.EvmKeeper),
        evmante.NewEthValidateBasicDecorator(options.EvmKeeper),
        evmante.NewEthSigVerificationDecorator(options.EvmKeeper),
        evmante.NewEthAccountVerificationDecorator(options.AccountKeeper, options.EvmKeeper),
        evmante.NewEthGasConsumeDecorator(options.EvmKeeper, options.MaxTxGasWanted),
        evmante.NewEthIncrementSenderSequenceDecorator(options.AccountKeeper),
        evmante.NewGasWantedDecorator(options.EvmKeeper, options.FeeMarketKeeper),
        evmante.NewEthEmitEventDecorator(options.EvmKeeper),
    )
}
```

### Step 8: Genesis Configuration

Configure EVM genesis state:

```go
// x/evm/genesis.go
func DefaultGenesisState() *GenesisState {
    return &GenesisState{
        Accounts: []GenesisAccount{},
        Params: DefaultParams(),
        ChainConfig: types.DefaultChainConfig(),
        Contracts: []GenesisContract{},
    }
}
```

### Step 9: CLI Commands

Add EVM-specific CLI commands:

```go
// x/evm/client/cli/tx.go
func GetTxCmd() *cobra.Command {
    cmd := &cobra.Command{
        Use:   types.ModuleName,
        Short: "EVM transaction subcommands",
        RunE:  client.ValidateCmd,
    }
    
    cmd.AddCommand(
        GetDeployCmd(),
        GetCallCmd(),
        GetEstimateGasCmd(),
    )
    
    return cmd
}
```

### Step 10: Testing Your EVM Chain

1. **Build the chain**:
   ```bash
   make install
   ```

2. **Initialize a local testnet**:
   ```bash
   florad init mynode --chain-id localchain_9000-1
   florad keys add mykey
   florad genesis add-genesis-account mykey 1000000000000000000flora
   florad genesis gentx mykey 1000000000000000000flora
   florad genesis collect-gentxs
   ```

3. **Start the chain**:
   ```bash
   florad start
   ```

4. **Deploy a smart contract**:
   ```bash
   florad tx evm deploy --from mykey --gas auto --gas-adjustment 1.3 \
     --bytecode 0x608060405234801561001057600080fd5b50610150806100206000396000f3fe...
   ```

## Important Considerations

### Account System Changes
- The conversion to EVM impacts the existing account system
- Address migration required between Cosmos and Ethereum formats
- Token balances need proper decimal conversion (6 → 18)

### Asset Migration
- Existing assets must be initialized in the EVM
- Native tokens become ERC20 compatible automatically
- IBC tokens can be used directly in smart contracts

### Genesis Migration
- Requires careful state migration from standard Cosmos chain
- Account balances need decimal adjustment
- Validator state preservation is critical

## Troubleshooting Guide

### Common Issues

1. **Keeper Initialization Order**
   - EVM keeper must be initialized before ERC20 keeper
   - Check dependency order in app.go

2. **Transaction Format Errors**
   - Verify chain ID format matches EVM requirements
   - Ensure proper encoding with `evmencoding.MakeConfig()`

3. **Gas Estimation**
   - EVM transactions require different gas calculation
   - Use `--gas auto` flag with adjustment factor

4. **IBC Transfer Compatibility**
   - Ensure IBC transfer keeper uses Cosmos EVM version
   - Verify channel capabilities include EVM support

## Next Steps

1. **Development Environment**
   - Set up local testnet
   - Configure MetaMask for local chain
   - Deploy test contracts

2. **Integration Testing**
   - Test EVM transaction processing
   - Verify IBC transfers with EVM tokens
   - Validate precompile functionality

3. **Production Preparation**
   - Security audit of custom modules
   - Performance testing under load
   - Documentation for validators

## Resources

- **Official Documentation**: https://evm.cosmos.network
- **GitHub Repository**: https://github.com/cosmos/evm
- **Community Support**: Cosmos Discord #cosmos-evm channel
- **Example Implementation**: Flora blockchain (this project)

## License

This project uses the Apache 2.0 licensed cosmos/evm module. The integration maintains compatibility with the Cosmos SDK's licensing while providing full EVM functionality.