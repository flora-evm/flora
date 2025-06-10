# Stage 2: State Management Implementation Plan

## Overview

Stage 2 extends the basic infrastructure from Stage 1 with sophisticated state management capabilities, queries, and validation logic. This stage prepares the module for actual tokenization operations in Stage 3.

## Timeline: Week 2

## Objectives

1. Extend keeper with tokenization record query operations
2. Implement validation logic for liquid staking caps
3. Add state consistency checks
4. Create comprehensive integration tests

## Implementation Tasks

### 2.1 Extended Keeper Interface

Define and implement additional keeper methods:

```go
// Queries
GetValidatorTokenizationRecords(ctx, validatorAddr) []TokenizationRecord
GetOwnerTokenizationRecords(ctx, ownerAddr) []TokenizationRecord
GetTokenizationRecordByDenom(ctx, denom) (TokenizationRecord, bool)

// Validation helpers
ValidateGlobalLiquidStakingCap(ctx) error
ValidateValidatorLiquidCap(ctx, validatorAddr) error
CanTokenizeShares(ctx, validatorAddr, shares) error

// State helpers
GetTotalLiquidStaked(ctx) math.Int
GetValidatorLiquidStaked(ctx, validatorAddr) math.Int
```

### 2.2 Query Implementation

Implement efficient queries using store indexes:

1. **Validator Index**: Map validator address -> record IDs
2. **Owner Index**: Map owner address -> record IDs
3. **Denom Index**: Map LST denom -> record ID

### 2.3 Validation Logic

Implement cap validation:

```go
// Check if tokenizing more shares would exceed global cap
func (k Keeper) ValidateGlobalLiquidStakingCap(ctx sdk.Context, additionalShares math.Int) error {
    params := k.GetParams(ctx)
    if !params.Enabled {
        return ErrModuleDisabled
    }
    
    totalStaked := k.stakingKeeper.TotalBondedTokens(ctx)
    totalLiquid := k.GetTotalLiquidStaked(ctx).Add(additionalShares)
    
    maxAllowed := params.GlobalLiquidStakingCap.MulInt(totalStaked).TruncateInt()
    if totalLiquid.GT(maxAllowed) {
        return ErrExceedsGlobalCap
    }
    return nil
}
```

### 2.4 Proto Updates

Add new query proto definitions:

```proto
service Query {
  rpc TokenizationRecords(QueryTokenizationRecordsRequest) returns (QueryTokenizationRecordsResponse);
  rpc TokenizationRecord(QueryTokenizationRecordRequest) returns (QueryTokenizationRecordResponse);
  rpc TokenizationRecordsByValidator(QueryTokenizationRecordsByValidatorRequest) returns (QueryTokenizationRecordsByValidatorResponse);
  rpc TokenizationRecordsByOwner(QueryTokenizationRecordsByOwnerRequest) returns (QueryTokenizationRecordsByOwnerResponse);
  rpc TotalLiquidStaked(QueryTotalLiquidStakedRequest) returns (QueryTotalLiquidStakedResponse);
}
```

### 2.5 Store Keys Design

```go
// Keys for efficient querying
var (
    // Primary storage
    TokenizationRecordKey = []byte{0x01} // id -> TokenizationRecord
    
    // Indexes
    ValidatorIndexKey = []byte{0x02} // validator -> []id
    OwnerIndexKey     = []byte{0x03} // owner -> []id
    DenomIndexKey     = []byte{0x04} // denom -> id
    
    // Aggregates
    TotalLiquidStakedKey = []byte{0x05} // -> total liquid staked
)
```

### 2.6 Testing Strategy

1. **Unit Tests**
   - Each new keeper method
   - Validation logic edge cases
   - Index consistency

2. **Integration Tests**
   - Multi-record scenarios
   - Cap validation with mock staking data
   - Query performance with large datasets

3. **Invariant Tests**
   - Total liquid staked consistency
   - Index synchronization
   - Cap compliance

## Dependencies

- Staking keeper interface (read-only for Stage 2)
- Additional proto generation
- Mock staking keeper for tests

## Success Criteria

- [ ] All queries return correct results
- [ ] Validation logic prevents cap violations
- [ ] Indexes remain consistent with primary storage
- [ ] 100% test coverage for new functionality
- [ ] Integration tests pass with mock staking keeper

## Risk Mitigation

1. **Index Corruption**: Implement index rebuild function
2. **Performance**: Add pagination to queries
3. **State Size**: Design efficient key encoding

## Next Steps (Stage 3 Preview)

With state management complete, Stage 3 will implement:
- MsgTokenizeShares transaction
- Liquid staking token minting
- Event emission
- Integration with bank module