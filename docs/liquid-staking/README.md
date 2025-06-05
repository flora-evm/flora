# Liquid Staking Documentation

## Overview
This directory contains comprehensive documentation for implementing liquid staking on Flora blockchain, combining Cosmos SDK staking with EVM precompiles for seamless DeFi integration.

## Documentation Structure

### 📐 [Architecture](./architecture/)
Technical design and integration patterns for liquid staking.

- [01-overview.md](./architecture/01-overview.md) - Complete architectural overview
  - Cosmos liquid staking module analysis
  - EVM precompile design
  - LST token implementation
  - Integration challenges and solutions

### 🔨 [Implementation](./implementation/)
Step-by-step implementation guide using a staged approach.

- [01-staged-approach.md](./implementation/01-staged-approach.md) - 18-stage implementation plan
  - Detailed breakdown of each stage
  - Dependencies and prerequisites
  - Success criteria
  
- [02-testing-strategy.md](./implementation/02-testing-strategy.md) - Comprehensive testing framework
  - Feature flags for stage control
  - Mock implementations
  - CI/CD integration
  
- [03-summary.md](./implementation/03-summary.md) - Visual overview and metrics
  - Architecture evolution diagram
  - Implementation checkpoints
  - Timeline and deliverables

### 💻 [Examples](./examples/)
Working code examples and reference implementations.

- [precompile.go](./examples/precompile.go) - Liquid staking precompile implementation
- [LiquidStakedToken.sol](./examples/LiquidStakedToken.sol) - Auto-compounding LST token contract
- [stage1-example/](./examples/stage1-example/) - Complete Stage 1 implementation
  - Basic types and validation
  - 100% test coverage
  - Benchmark tests

## Quick Start

### For Developers
1. Read the [architectural overview](./architecture/01-overview.md)
2. Review the [staged implementation plan](./implementation/01-staged-approach.md)
3. Start with [Stage 1 example](./examples/stage1-example/)

### For Architects
1. Review [integration challenges](./architecture/01-overview.md#integration-challenges--solutions)
2. Understand the [testing strategy](./implementation/02-testing-strategy.md)
3. Check [implementation checkpoints](./implementation/03-summary.md#implementation-checkpoints)

### For Project Managers
1. Review the [timeline](./implementation/03-summary.md#timeline-summary)
2. Understand [risk mitigation](./implementation/03-summary.md#risk-mitigation)
3. Track [success metrics](./implementation/03-summary.md#success-metrics)

## Key Features

### 🌉 EVM-Cosmos Bridge
- Precompiles at `0x800` and `0x801`
- Direct access to Cosmos staking from smart contracts
- Seamless address conversion

### 🪙 Liquid Staking Tokens (LST)
- Auto-compounding ERC20 tokens
- Per-validator tokens (stFLORA-{ValidatorID})
- IBC-enabled for cross-chain transfers

### 🔒 Security Features
- Global and per-validator caps
- Slashing protection
- Governance controls
- Rate limiting

### 🧪 Testing Approach
- 18 independently testable stages
- Feature flags for progressive rollout
- Comprehensive mock framework
- Automated CI/CD pipeline

## Implementation Timeline

| Phase | Stages | Duration | Deliverable |
|-------|--------|----------|-------------|
| Foundation | 1-3 | 3 weeks | Basic infrastructure |
| Core Logic | 4-6 | 3 weeks | Cosmos integration |
| EVM Bridge | 7-11 | 5 weeks | Precompile implementation |
| Advanced | 12-15 | 4 weeks | Full features |
| Production | 16-18 | 3 weeks | Enterprise features |
| **Total** | **1-18** | **18 weeks** | **Complete liquid staking** |

## Technical Specifications

### Gas Costs
- Tokenize shares: 100,000 gas
- Redeem tokens: 80,000 gas  
- Transfer record: 50,000 gas
- Query operations: 5,000-10,000 gas

### Risk Parameters
- Global liquid staking cap: 25%
- Per-validator cap: 50%
- Redemption cooldown: 21 days
- Tokenization/redemption fee: 0.1%

### Dependencies
- Cosmos SDK v0.50.13
- cosmos/gaia liquid staking module
- EVM precompile framework
- Token Factory module

## Contributing

When adding new documentation:
1. Follow the existing structure
2. Include code examples
3. Add to appropriate section
4. Update this README

## Resources

- [Cosmos Liquid Staking Module](https://github.com/cosmos/gaia/tree/main/x/liquid)
- [EVM Precompiles Guide](https://docs.evmos.org/develop/smart-contracts/evm-extensions/precompiles)
- [Flora Architecture](../README.md)

## Status

🚧 **Research Phase** - Architecture designed, implementation plan ready

Next steps:
1. Team review and approval
2. Begin Stage 1 implementation
3. Set up testing infrastructure