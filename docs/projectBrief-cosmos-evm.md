# Project Brief: Building a Cosmos Chain with EVM Support

## Overview
This document provides a comprehensive guide for creating a new Cosmos blockchain with full EVM (Ethereum Virtual Machine) support using the official cosmos/evm module from Interchain Labs. The integration enables your chain to run Ethereum smart contracts while maintaining all Cosmos SDK features like IBC, governance, and staking.

## Prerequisites
- Go 1.23+ installed
- Basic understanding of Cosmos SDK
- Familiarity with Ethereum/EVM concepts
- Development environment with git, make, and Docker

## Architecture Overview

### Dual Transaction Support
The chain will support two types of transactions:
1. **Cosmos Transactions**: Native Cosmos SDK messages (bank transfers, staking, governance)
2. **Ethereum Transactions**: EVM transactions for smart contracts and ERC20 tokens

### Key Modules
- **x/evm**: Core EVM functionality
- **x/erc20**: Bidirectional token conversion between Cosmos and EVM
- **x/feemarket**: EIP-1559 dynamic fee market implementation

## Step 1: Create Base Cosmos Chain

### Option A: Using Ignite CLI
```bash
# Install Ignite CLI
curl https://get.ignite.com/cli | bash

# Scaffold a new chain
ignite scaffold chain evmchain \
  --address-prefix evmc \
  --no-module

# Navigate to the chain directory
cd evmchain
```

### Option B: Manual Setup
```bash
# Create project structure
mkdir evmchain && cd evmchain
git init

# Initialize go module
go mod init github.com/myorg/evmchain

# Create basic directory structure
mkdir -p app cmd/evmchaind x proto scripts
```

## Step 2: Update Dependencies

Edit `go.mod` to include EVM dependencies:

```go
module github.com/myorg/evmchain

go 1.23

require (
    cosmossdk.io/api v0.7.5
    cosmossdk.io/core v0.12.0
    cosmossdk.io/depinject v1.0.0
    cosmossdk.io/errors v1.0.1
    cosmossdk.io/log v1.4.1
    cosmossdk.io/math v1.3.0
    cosmossdk.io/store v1.1.1
    cosmossdk.io/x/tx v0.13.5
    github.com/cosmos/cosmos-db v1.0.3
    github.com/cosmos/cosmos-proto v1.0.0-beta.5
    github.com/cosmos/cosmos-sdk v0.50.10
    github.com/cosmos/evm v1.0.0-rc2  // Official EVM module
    github.com/cosmos/gogoproto v1.7.0
    github.com/cosmos/ibc-go/v8 v8.5.1
    github.com/spf13/cobra v1.8.1
    github.com/spf13/viper v1.19.0
    google.golang.org/grpc v1.67.1
    google.golang.org/protobuf v1.35.1
)

// Required replacements for EVM
replace (
    // Replace go-ethereum with evmos fork
    github.com/ethereum/go-ethereum => github.com/evmos/go-ethereum v1.10.26-evmos-rc4
)
```

## Step 3: Configure Chain ID

The chain ID must follow a specific format for EVM compatibility:

```
Format: {identifier}_{evm-chain-id}-{version}

Example: evmchain_9000-1
- identifier: evmchain
- evm-chain-id: 9000 (used in MetaMask)
- version: 1
```

## Step 4: Create EVM Configuration

Create `app/evm_config.go`:

```go
package app

import (
    "github.com/cosmos/evm/x/evm/types"
    feemarkettypes "github.com/cosmos/evm/x/feemarket/types"
)

// EVMConfig returns the default EVM configuration
func EVMConfig() *types.EVMConfig {
    return &types.EVMConfig{
        // Enable all Ethereum hard forks
        ChainConfig: types.DefaultChainConfig(),
        // Extra EIPs to enable
        ExtraEIPs: []int64{3855}, // PUSH0 opcode
        // Allow unprotected transactions (for development)
        AllowUnprotectedTxs: false,
    }
}

// FeeMarketConfig returns the default fee market configuration
func FeeMarketConfig() *feemarkettypes.Params {
    return &feemarkettypes.Params{
        NoBaseFee:                false,
        BaseFeeChangeDenominator: 8,
        ElasticityMultiplier:     2,
        EnableHeight:             0,
        BaseFee:                  types.DefaultBaseFee,
        MinGasPrice:              types.DefaultMinGasPrice,
        MinGasMultiplier:         types.DefaultMinGasMultiplier,
    }
}
```

## Step 5: Update App Structure

Update `app/app.go` to include EVM modules:

```go
package app

import (
    // ... standard imports ...
    
    // EVM imports
    evmkeeper "github.com/cosmos/evm/x/evm/keeper"
    evmtypes "github.com/cosmos/evm/x/evm/types"
    erc20keeper "github.com/cosmos/evm/x/erc20/keeper"
    erc20types "github.com/cosmos/evm/x/erc20/types"
    feemarketkeeper "github.com/cosmos/evm/x/feemarket/keeper"
    feemarkettypes "github.com/cosmos/evm/x/feemarket/types"
)

// Add to module account permissions
var maccPerms = map[string][]string{
    // ... existing permissions ...
    evmtypes.ModuleName:      {authtypes.Minter, authtypes.Burner},
    erc20types.ModuleName:    {},
    feemarkettypes.ModuleName: {},
}

type App struct {
    // ... existing fields ...
    
    // EVM keepers
    EvmKeeper       *evmkeeper.Keeper
    Erc20Keeper     erc20keeper.Keeper
    FeeMarketKeeper feemarketkeeper.Keeper
}

// In NewApp function, add keeper initialization:
func NewApp(...) *App {
    // ... existing initialization ...
    
    // Create EVM keepers
    app.FeeMarketKeeper = feemarketkeeper.NewKeeper(
        appCodec,
        runtime.NewKVStoreService(keys[feemarkettypes.StoreKey]),
        app.ConsensusParamsKeeper,
        app.Logger(),
    )
    
    app.EvmKeeper = evmkeeper.NewKeeper(
        appCodec,
        runtime.NewKVStoreService(keys[evmtypes.StoreKey]),
        app.AccountKeeper,
        app.BankKeeper,
        app.StakingKeeper,
        app.FeeMarketKeeper,
        nil, // precompiles added later
        app.Logger(),
    )
    
    app.Erc20Keeper = erc20keeper.NewKeeper(
        runtime.NewKVStoreService(keys[erc20types.StoreKey]),
        appCodec,
        app.AccountKeeper,
        app.BankKeeper,
        app.EvmKeeper,
        app.Logger(),
    )
    
    // ... rest of initialization ...
}
```

## Step 6: Configure Ante Handlers

Create `app/ante/ante.go`:

```go
package ante

import (
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/cosmos/cosmos-sdk/x/auth/ante"
    evmante "github.com/cosmos/evm/x/evm/ante"
    evmtypes "github.com/cosmos/evm/x/evm/types"
)

// NewAnteHandler creates a dual ante handler for both Cosmos and EVM transactions
func NewAnteHandler(options HandlerOptions) sdk.AnteHandler {
    return func(ctx sdk.Context, tx sdk.Tx, sim bool) (sdk.Context, error) {
        // Route to EVM ante handler for Ethereum transactions
        if _, ok := tx.(evmtypes.EvmTxWrapper); ok {
            return NewEVMAnteHandler(options)(ctx, tx, sim)
        }
        
        // Use Cosmos ante handler for standard transactions
        return NewCosmosAnteHandler(options)(ctx, tx, sim)
    }
}

func NewCosmosAnteHandler(options HandlerOptions) sdk.AnteHandler {
    return sdk.ChainAnteDecorators(
        ante.NewSetUpContextDecorator(),
        ante.NewExtensionOptionsDecorator(options.ExtensionOptionChecker),
        ante.NewValidateBasicDecorator(),
        ante.NewTxTimeoutHeightDecorator(),
        ante.NewValidateMemoDecorator(options.AccountKeeper),
        ante.NewConsumeGasForTxSizeDecorator(options.AccountKeeper),
        ante.NewDeductFeeDecorator(options.AccountKeeper, options.BankKeeper, options.FeegrantKeeper, options.TxFeeChecker),
        ante.NewSetPubKeyDecorator(options.AccountKeeper),
        ante.NewValidateSigCountDecorator(options.AccountKeeper),
        ante.NewSigGasConsumeDecorator(options.AccountKeeper, options.SigGasConsumer),
        ante.NewSigVerificationDecorator(options.AccountKeeper, options.SignModeHandler),
        ante.NewIncrementSequenceDecorator(options.AccountKeeper),
    )
}

func NewEVMAnteHandler(options HandlerOptions) sdk.AnteHandler {
    return sdk.ChainAnteDecorators(
        evmante.NewEthAccountVerificationDecorator(options.AccountKeeper, options.EvmKeeper),
        evmante.NewEthGasConsumeDecorator(options.EvmKeeper, options.MaxTxGasWanted),
        evmante.NewEthIncrementSenderSequenceDecorator(options.AccountKeeper),
        evmante.NewGasWantedDecorator(options.EvmKeeper, options.FeeMarketKeeper),
        evmante.NewEthEmitEventDecorator(options.EvmKeeper),
    )
}
```

## Step 7: Configure Genesis

Update `app/genesis.go`:

```go
// DefaultGenesis returns default genesis state
func (app *App) DefaultGenesis() map[string]json.RawMessage {
    genesis := ModuleBasics.DefaultGenesis(app.appCodec)
    
    // Configure EVM genesis
    evmGenesis := evmtypes.DefaultGenesisState()
    evmGenesis.Params.EvmDenom = "aevmc" // smallest unit of your token
    genesis[evmtypes.ModuleName] = app.appCodec.MustMarshalJSON(evmGenesis)
    
    // Configure ERC20 genesis
    erc20Genesis := erc20types.DefaultGenesisState()
    genesis[erc20types.ModuleName] = app.appCodec.MustMarshalJSON(erc20Genesis)
    
    // Configure fee market genesis
    feeMarketGenesis := feemarkettypes.DefaultGenesisState()
    feeMarketGenesis.Params.BaseFee = sdk.NewInt(1000000000) // 1 gwei
    genesis[feemarkettypes.ModuleName] = app.appCodec.MustMarshalJSON(feeMarketGenesis)
    
    return genesis
}
```

## Step 8: Add Precompiled Contracts

Create `app/precompiles.go`:

```go
package app

import (
    "github.com/ethereum/go-ethereum/common"
    evmkeeper "github.com/cosmos/evm/x/evm/keeper"
    
    // Import standard precompiles
    "github.com/cosmos/evm/precompiles/bank"
    "github.com/cosmos/evm/precompiles/distribution"
    "github.com/cosmos/evm/precompiles/governance"
    "github.com/cosmos/evm/precompiles/staking"
)

// SetupPrecompiles registers all precompiled contracts
func (app *App) SetupPrecompiles() {
    precompiles := map[common.Address]evmkeeper.PrecompiledContract{
        // Bank precompile at 0x0000000000000000000000000000000000000800
        common.HexToAddress("0x0000000000000000000000000000000000000800"): bank.NewPrecompile(app.BankKeeper),
        
        // Staking precompile at 0x0000000000000000000000000000000000000801
        common.HexToAddress("0x0000000000000000000000000000000000000801"): staking.NewPrecompile(app.StakingKeeper),
        
        // Distribution precompile at 0x0000000000000000000000000000000000000802
        common.HexToAddress("0x0000000000000000000000000000000000000802"): distribution.NewPrecompile(app.DistrKeeper),
        
        // Governance precompile at 0x0000000000000000000000000000000000000803
        common.HexToAddress("0x0000000000000000000000000000000000000803"): governance.NewPrecompile(app.GovKeeper),
    }
    
    app.EvmKeeper.SetPrecompiles(precompiles)
}
```

## Step 9: CLI Integration

Add EVM commands to `cmd/evmchaind/cmd/root.go`:

```go
import (
    evmcli "github.com/cosmos/evm/x/evm/client/cli"
    erc20cli "github.com/cosmos/evm/x/erc20/client/cli"
)

func initRootCmd(rootCmd *cobra.Command, encodingConfig params.EncodingConfig) {
    // ... existing initialization ...
    
    // Add EVM transaction commands
    rootCmd.AddCommand(
        evmcli.GetTxCmd(),
        evmcli.GetQueryCmd(),
        erc20cli.GetTxCmd(),
        erc20cli.GetQueryCmd(),
    )
}
```

## Step 10: Testing Your EVM Chain

### Build and Initialize
```bash
# Build the binary
make install

# Initialize the chain
evmchaind init mynode --chain-id evmchain_9000-1

# Add a validator key
evmchaind keys add validator

# Add genesis account with balance
evmchaind genesis add-genesis-account validator 1000000000000000000aevmc

# Create genesis transaction
evmchaind genesis gentx validator 1000000000000000000aevmc \
  --chain-id evmchain_9000-1

# Collect genesis transactions
evmchaind genesis collect-gentxs

# Start the chain
evmchaind start
```

### Deploy a Smart Contract
```bash
# Deploy a simple ERC20 contract
evmchaind tx evm deploy-contract \
  --contract-bytecode 0x608060... \
  --from validator \
  --gas 5000000 \
  --gas-prices 1000000000aevmc
```

### Connect MetaMask
1. Add custom network in MetaMask:
   - Network Name: EVMChain Local
   - RPC URL: http://localhost:8545
   - Chain ID: 9000
   - Currency Symbol: EVMC

## Common Issues and Solutions

### Issue 1: Account Compatibility
**Problem**: Existing Cosmos accounts can't be used with EVM
**Solution**: Generate new accounts with Ethereum-compatible addresses

### Issue 2: Token Decimals
**Problem**: Cosmos uses 6 decimals, Ethereum uses 18
**Solution**: Handle conversion in your application layer

### Issue 3: Gas Estimation
**Problem**: Gas estimation fails for complex contracts
**Solution**: Use higher gas limits or implement custom gas estimation

## Next Steps

1. **Add Custom Precompiles**: Implement chain-specific functionality
2. **Configure IBC**: Enable cross-chain token transfers
3. **Set Up Monitoring**: Use tools like Prometheus and Grafana
4. **Deploy Explorer**: Set up BlockScout or similar
5. **Security Audit**: Before mainnet launch

## Resources

- [Cosmos EVM Documentation](https://evm.cosmos.network/)
- [Cosmos SDK Documentation](https://docs.cosmos.network/)
- [Ethereum JSON-RPC Spec](https://ethereum.org/en/developers/docs/apis/json-rpc/)
- [EVM Opcodes Reference](https://www.evm.codes/)

## Support

- Discord: [Cosmos Discord](https://discord.gg/cosmos)
- Forum: [Cosmos Forum](https://forum.cosmos.network/)
- GitHub: [cosmos/evm](https://github.com/cosmos/evm)