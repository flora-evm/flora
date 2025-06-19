# Comprehensive Upgrade Plan: Flora to Official cosmos/evm

## Executive Summary

After deep analysis of the official cosmos/evm repository, migration is **NOT RECOMMENDED** until v1.0 stable release (expected late Q2 2025). The migration involves significant breaking changes and SDK incompatibilities that would require 10-12 weeks of development with high risk.

## Current State Analysis

### Flora's Current Setup
- **Cosmos SDK**: v0.50.13 (via strangelove-ventures fork)
- **EVM Implementation**: strangelove-ventures/cosmos-evm v0.1.5
- **Key Features**: Token Factory, Liquid Staking, EVM Precompiles
- **Status**: Functional, producing blocks

### Official cosmos/evm Status
- **Repository**: https://github.com/cosmos/evm
- **Latest Version**: v1.0.0-rc2 (Release Candidate)
- **SDK Requirement**: v0.53.2+ (Major version jump from v0.50.x)
- **Audit Status**: Pending (targeting late Q2 2025)
- **Maintainer**: Interchain Labs / Interchain Foundation

## Critical Breaking Changes

### 1. SDK Version Incompatibility (Highest Impact)
```diff
- Cosmos SDK v0.50.13
+ Cosmos SDK v0.53.2+
```

**Impact on Flora**:
- All type imports need updating (sdk.Int → math.Int)
- Store access patterns change (KVStore → StoreService)
- Keeper interfaces incompatible (StakingKeeper, BankKeeper, etc.)
- Parameter management migration required

### 2. Module Rename: x/evm → x/vm
```diff
- import "github.com/strangelove-ventures/cosmos-evm/x/evm"
+ import "github.com/cosmos/evm/x/vm"
```

**State Migration Required**:
- All store keys with "evm" prefix must migrate to "vm"
- Module registration updates in app.go
- Proto namespace changes

### 3. ERC20 Denomination Breaking Change
```diff
- Denom: "erc20/0x1234567890abcdef..."
+ Denom: "erc20{0x1234567890abcdef...}"
```

**Impact**: All existing ERC20 tokens need denomination migration

### 4. Precompile Security Model Change
- **New Requirement**: `msg.sender == requester` (no proxy calls)
- **Impact**: Breaks any smart contracts using delegated calls
- **Flora Impact**: TokenFactory precompile interactions may break

## Migration Challenges Discovered

### 1. Evmos Migration Lessons
- **What Happened**: Evmos deprecated ALL Cosmos transactions in Q3 2024
- **User Impact**: Forced migration from Keplr to MetaMask
- **Community Response**: Loss of Cosmos features, user frustration
- **Lesson**: Migration can fundamentally change user experience

### 2. Critical Security Issues
- **Issue Found**: Precompile bug allowing unlimited token mint
- **Affected Chains**: Evmos, FunctionX/PundiX
- **Root Cause**: Incorrect validation in delegation precompile
- **Lesson**: Precompile changes require extensive security auditing

### 3. No Official Migration Guide
- cosmos/evm repository lacks migration documentation
- No automated tools for state migration
- Community left to figure out migration independently

## Detailed Migration Plan (If Proceeding)

### Phase 1: Preparation (3-4 weeks)
1. **Environment Setup**
   ```bash
   # Create migration branch
   git checkout -b migration/cosmos-evm-v1
   
   # Update go.mod
   go get github.com/cosmos/cosmos-sdk@v0.53.2
   go get github.com/cosmos/evm@v1.0.0
   ```

2. **Dependency Analysis**
   - Map all keeper interface changes
   - Identify custom modifications in strangelove fork
   - Document all precompile customizations

3. **Testing Infrastructure**
   - Set up migration test environment
   - Create state export/import tools
   - Prepare rollback procedures

### Phase 2: Code Migration (4-6 weeks)

#### A. Update Imports
```bash
# Automated script for import updates
find . -name "*.go" -type f -exec sed -i '' \
  -e 's|github.com/strangelove-ventures/cosmos-evm|github.com/cosmos/evm|g' \
  -e 's|github.com/evmos/|github.com/cosmos/|g' \
  -e 's|/x/evm|/x/vm|g' \
  -e 's|sdk\.Int|math.Int|g' \
  -e 's|sdk\.Dec|math.LegacyDec|g' \
  {} +
```

#### B. Update Module Registration
```go
// app/app.go changes
- import evmkeeper "github.com/strangelove-ventures/cosmos-evm/x/evm/keeper"
+ import vmkeeper "github.com/cosmos/evm/x/vm/keeper"

- app.EvmKeeper = evmkeeper.NewKeeper(...)
+ app.VmKeeper = vmkeeper.NewKeeper(...)
```

#### C. State Migration Handler
```go
func MigrateV1toV2(ctx sdk.Context, k Keeper) error {
    // 1. Migrate x/evm store to x/vm
    oldStore := ctx.KVStore(oldEvmKey)
    newStore := ctx.KVStore(newVmKey)
    
    iterator := oldStore.Iterator(nil, nil)
    defer iterator.Close()
    
    for ; iterator.Valid(); iterator.Next() {
        newStore.Set(iterator.Key(), iterator.Value())
    }
    
    // 2. Update ERC20 denominations
    bankKeeper.IterateAllDenomMetaData(ctx, func(metadata banktypes.Metadata) bool {
        if strings.HasPrefix(metadata.Base, "erc20/") {
            newDenom := strings.Replace(metadata.Base, "erc20/", "erc20", 1)
            // Update metadata and balances
        }
        return false
    })
    
    // 3. Migrate liquid staking module
    // Update any EVM-related references
    
    return nil
}
```

#### D. Update Precompiles
```go
// Update all precompiles to enforce direct calls
func (p *Precompile) Run(evm *vm.EVM, caller common.Address, input []byte) ([]byte, error) {
    // New security check
    if evm.Origin != caller {
        return nil, errors.New("precompile requires direct call")
    }
    // ... rest of implementation
}
```

### Phase 3: Testing (2-3 weeks)
1. **Unit Tests**
   - Update all test imports
   - Fix keeper mock interfaces
   - Validate precompile behavior

2. **Integration Tests**
   - Test state migration on testnet
   - Validate ERC20 token functionality
   - Test liquid staking with new EVM

3. **Load Testing**
   - Performance comparison
   - Gas cost analysis
   - Transaction throughput

### Phase 4: Deployment (1-2 weeks)
1. **Testnet Deployment**
   - Export genesis state
   - Run migration
   - Validate all functionality

2. **Mainnet Upgrade**
   - Coordinate validator upgrade
   - Execute at predetermined block height
   - Monitor for issues

## Risk Assessment

### High Risks
1. **State Corruption**: Migration could corrupt state if not properly tested
2. **User Experience**: Similar to Evmos, could lose Cosmos functionality
3. **Security Vulnerabilities**: New precompile model could introduce bugs
4. **Validator Coordination**: Requires 67%+ validators to upgrade simultaneously

### Mitigation Strategies
1. **Extensive Testing**: 6+ weeks of testing before mainnet
2. **Rollback Plan**: Keep ability to revert to v0.50.13
3. **Security Audit**: Audit all precompile changes
4. **Phased Rollout**: Test on canary validators first

## Recommendation: Wait for v1.0 Stable

### Reasons to Wait:
1. **cosmos/evm still in RC stage** - Wait for stable v1.0 release
2. **No official migration guide** - Wait for documentation
3. **Security audit pending** - Wait for audit completion
4. **Major SDK version jump** - High risk of breaking changes
5. **Current setup works** - No urgent need to migrate

### Alternative Strategy:
1. **Continue with current setup** (strangelove-ventures/cosmos-evm v0.1.5)
2. **Complete liquid staking implementation**
3. **Monitor cosmos/evm development**
4. **Plan migration for Q3 2025** after v1.0 stable release

### When to Reconsider:
- cosmos/evm v1.0 stable release available
- Official migration guide published
- Security audit completed and passed
- Clear benefits outweigh migration risks
- Community successfully migrates other chains

## Conclusion

The migration from strangelove-ventures/cosmos-evm to the official cosmos/evm involves significant technical challenges and risks. The lack of stable release, pending audit, and major breaking changes make immediate migration inadvisable.

**Recommended Action**: Continue with current implementation and revisit migration after cosmos/evm v1.0 stable release with official migration documentation.

## Resources
- [cosmos/evm Repository](https://github.com/cosmos/evm)
- [Cosmos SDK v0.53 Migration](https://docs.cosmos.network/v0.53/migrations/v0.50-to-v0.53)
- [EVM Documentation](https://evm.cosmos.network/)
- [Interchain Labs Contact](https://share-eu1.hsforms.com/2g6yO-PVaRoKj50rUgG4Pjg2e2sca)