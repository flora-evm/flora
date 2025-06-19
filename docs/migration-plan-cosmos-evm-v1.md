# Migration Plan: strangelove-ventures/cosmos-evm to cosmos/evm

## Executive Summary

The official `cosmos/evm` repository represents a significant evolution from `strangelove-ventures/cosmos-evm`. This document provides a comprehensive analysis of the migration requirements and breaking changes.

## Key Findings

### 1. Repository Status
- **Repository**: https://github.com/cosmos/evm
- **Created**: March 18, 2025
- **Latest Tag**: v1.0.0-rc2 (June 2025)
- **Status**: Under development, pending audit (targeting late Q2 for v1.0 stable release)
- **Base**: Fork of evmOS, maintained by Interchain Labs and Interchain Foundation

### 2. SDK Version Requirements

**Critical Breaking Change**: cosmos/evm requires Cosmos SDK v0.53.2, while Flora currently uses v0.50.13

```go
// cosmos/evm go.mod
require (
    github.com/cosmos/cosmos-sdk v0.53.2
)
```

This is a major version jump (v0.50.x → v0.53.x) that involves significant breaking changes in the SDK itself.

### 3. Module Structure Changes

#### Module Rename: x/evm → x/vm
The most significant change is the renaming of the main EVM module:
- **Old**: `x/evm` (in strangelove-ventures/cosmos-evm)
- **New**: `x/vm` (in cosmos/evm)

#### Available Modules in cosmos/evm:
- `x/vm` - Core EVM functionality (renamed from x/evm)
- `x/erc20` - ERC20 token management
- `x/feemarket` - EIP-1559 fee market
- `x/precisebank` - New module for decimal precision in EVM
- `x/ibc` - IBC-related functionality

### 4. Import Path Changes

All import paths have changed:
```go
// Old (strangelove-ventures)
import "github.com/strangelove-ventures/cosmos-evm/x/evm"

// New (cosmos)
import "github.com/cosmos/evm/x/vm"
```

Additionally, protobuf paths have changed:
- **Old**: `evmos.*` or custom paths
- **New**: `cosmos.evm.*`

### 5. Breaking Changes (from CHANGELOG)

#### State Breaking Changes:
1. Renamed `x/evm` to `x/vm`
2. Renamed protobuf files from evmos to cosmos organization
3. Removed base fee v1 from x/feemarket (#83)
4. Removed legacy subspaces (#93)
5. Changed native ERC20 denoms prefix from `erc20/` to `erc20` for IBC v2 (#95)
6. Removed x/authz dependency from precompiles (#62)

#### API Breaking Changes:
1. Updated ics20 precompile to use `Denom` instead of `DenomTrace` for IBC v2
2. Evidence precompile now requires explicit submitter address as first argument
3. All precompiles enforce `msg.sender == requester` (no proxy calls allowed)

### 6. New Features

1. **x/precisebank**: New module for decimal precision handling in EVM
2. **Permissionless ERC20 registration**: Allows conversion to Cosmos coins
3. **Enhanced precompiles**: More extensive set including bank, distribution, gov, evidence, etc.

### 7. Precompile Interface Changes

The precompile structure has evolved. Example from bank precompile:
```go
type Precompile struct {
    cmn.Precompile
    bankKeeper  cmn.BankKeeper
    erc20Keeper erc20keeper.Keeper
}
```

Key changes:
- All precompiles now enforce direct EOA calls (no proxy contracts)
- Standardized common interface through `cmn.Precompile`
- Integration with multiple keepers for cross-module functionality

### 8. Migration Challenges

1. **No Official Migration Guide**: Currently, there's no documented migration path from strangelove-ventures to cosmos/evm
2. **SDK Version Gap**: Moving from SDK v0.50.x to v0.53.x requires addressing all SDK breaking changes
3. **State Migration**: Module rename (evm → vm) requires state migration
4. **Active Development**: Repository is still under active development with pending audit

## Migration Strategy Recommendations

### Option 1: Wait for Stable Release
- **Pros**: Avoid breaking changes, get audited code
- **Cons**: Delays implementation (late Q2 target)
- **Recommendation**: Best for production environments

### Option 2: Gradual Migration
1. First upgrade Cosmos SDK from v0.50.13 to v0.53.2
2. Test all existing functionality with new SDK
3. Then migrate from strangelove-ventures to cosmos/evm
4. Implement state migration for module rename

### Option 3: Stay with Current Implementation
- Continue using strangelove-ventures/cosmos-evm
- Monitor cosmos/evm development
- Plan migration after v1.0 stable release with official migration guide

## Technical Migration Steps (If Proceeding)

1. **Update SDK Version**:
   ```go
   // go.mod
   require github.com/cosmos/cosmos-sdk v0.53.2
   ```

2. **Update Import Paths**:
   - Replace all `github.com/strangelove-ventures/cosmos-evm` with `github.com/cosmos/evm`
   - Update module references from `x/evm` to `x/vm`

3. **Update Proto Files**:
   - Regenerate all protobuf files with new cosmos.evm namespace

4. **State Migration**:
   - Implement upgrade handler to migrate x/evm state to x/vm
   - Update all store keys and module names

5. **Update Precompiles**:
   - Review all custom precompiles for interface compatibility
   - Ensure direct EOA call enforcement

## Conclusion

The migration from strangelove-ventures/cosmos-evm to cosmos/evm involves significant breaking changes, primarily due to:
1. SDK version upgrade (v0.50.x → v0.53.x)
2. Module rename (x/evm → x/vm)
3. Complete import path changes
4. Precompile interface updates

Given the current state (pre-audit, rc2 release), and lack of official migration documentation, **I recommend waiting for the stable v1.0 release** before attempting migration, unless there's an urgent need for the new features.

## Resources

- Repository: https://github.com/cosmos/evm
- Documentation: https://evm.cosmos.network/
- Contact: [Interchain Labs team](https://share-eu1.hsforms.com/2g6yO-PVaRoKj50rUgG4Pjg2e2sca)
- Support: Cosmos Tech Slack channel or [Telegram](https://t.me/cosmostechstack)