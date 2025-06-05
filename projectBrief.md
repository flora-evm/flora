# Flora Blockchain Project Brief

## Project Overview
Flora is a Cosmos SDK-based blockchain that uniquely combines Ethereum Virtual Machine (EVM) compatibility with native Cosmos functionality, providing developers with the best of both ecosystems.

## Main Goal
Create a high-performance blockchain that enables:
- Seamless deployment of Ethereum smart contracts
- Native token creation through Token Factory
- Inter-blockchain communication via IBC
- Unified experience for both Cosmos and Ethereum developers

## Key Features

### Core Capabilities
- **EVM Support**: Full Ethereum compatibility for smart contracts and dApps
- **Token Factory**: Permissionless token creation and management
- **IBC Protocol**: Connect to any IBC-enabled blockchain
- **Dual Transaction Support**: Process both Cosmos and Ethereum transaction types
- **EVM Precompiles**: Access Cosmos native functions from smart contracts

### Technical Specifications
- **Chain ID**: localchain_9000-1
- **Native Token**: petal (18 decimals, EVM-compatible)
- **Consensus**: CometBFT (Tendermint)
- **Gas Model**: Supports both Cosmos and EVM gas pricing
- **Account Model**: Ethereum-style accounts with Cosmos addresses

## Target Audience
- **DeFi Developers**: Build cross-chain financial applications
- **dApp Developers**: Deploy Ethereum apps with IBC connectivity
- **Token Creators**: Issue custom tokens without smart contracts
- **Institutional Users**: Leverage enterprise-grade Cosmos infrastructure

## Architecture

### Technology Stack
- **Framework**: Cosmos SDK v0.50.13
- **Consensus**: CometBFT v0.38.17
- **EVM**: Ethermint/Evmos modules
- **Language**: Go
- **Testing**: Interchaintest framework

### Module Architecture
1. **Standard Cosmos Modules**: auth, bank, staking, gov, mint, slashing, distribution
2. **IBC Modules**: transfer, ICA host/controller, fee middleware
3. **EVM Modules**: evm, erc20, feemarket
4. **Custom Modules**: tokenfactory

### Key Design Decisions
- Ethereum account compatibility for seamless wallet support
- Precompiled contracts for gas-efficient Cosmos operations
- Token Factory for non-technical token creation
- EIP-1559 style dynamic fees via fee market module

## Development Status
- **Stage**: Testnet-ready
- **Testing**: Comprehensive unit and integration tests
- **Infrastructure**: Docker support, local development tools
- **Documentation**: In development

## Future Roadmap Considerations
- Liquid staking integration
- Enhanced smart contract templates
- Cross-chain DEX functionality
- Governance improvements
- Performance optimizations

## Getting Started
```bash
# Build the chain
make install

# Run tests
go test ./... -v

# Start local testnet
make sh-testnet
```

For detailed development instructions, see [CLAUDE.md](./CLAUDE.md)