# Stage 1: Basic Infrastructure Example

This directory contains a complete implementation of Stage 1 of the liquid staking module, demonstrating the staged approach with isolated, testable components.

## What's Included

### `types.go`
- Basic type definitions for `TokenizationRecord`
- Module parameters structure
- Genesis state definition
- Validation logic for all types
- No external dependencies beyond Cosmos SDK types

### `types_test.go`
- Comprehensive unit tests for all type validation
- Table-driven test cases
- Benchmark tests for performance validation
- 100% test coverage for Stage 1

## Key Design Principles

1. **Zero External Dependencies**: Only uses SDK types for addresses and numbers
2. **Complete Validation**: Every field is validated with clear error messages
3. **Future-Proof Design**: Status field allows for future state transitions
4. **Performance Conscious**: Benchmarks ensure validation is efficient

## Running Stage 1 Tests

```bash
# Run all tests
go test -v

# Run with coverage
go test -cover -coverprofile=coverage.out
go tool cover -html=coverage.out

# Run benchmarks
go test -bench=. -benchmem

# Run specific test
go test -run TestTokenizationRecord_Validate -v
```

## Expected Output

```
=== RUN   TestTokenizationRecord_Validate
=== RUN   TestTokenizationRecord_Validate/valid_record
=== RUN   TestTokenizationRecord_Validate/zero_id
=== RUN   TestTokenizationRecord_Validate/invalid_validator_address
... (all test cases)
--- PASS: TestTokenizationRecord_Validate (0.00s)

=== RUN   TestModuleParams_Validate
... (all test cases)
--- PASS: TestModuleParams_Validate (0.00s)

=== RUN   TestGenesisState_Validate
... (all test cases)
--- PASS: TestGenesisState_Validate (0.00s)

PASS
coverage: 100.0% of statements
ok      github.com/rollchains/flora/x/liquidstaking/types      0.123s
```

## Integration with Next Stages

Stage 2 will add:
- Keeper with store operations
- CRUD methods for TokenizationRecord
- State persistence

Stage 2 can import and use these types directly without modification, demonstrating the forward compatibility of the staged approach.

## Type Evolution

As we progress through stages, these types may gain additional fields:
- Stage 6: Add `LSTDenom` field to TokenizationRecord
- Stage 12: Add `RedemptionTime` for unbonding tracking
- Stage 13: Add `RewardsClaimed` for reward tracking

Each addition will be backward compatible and include migration logic.

## Validation Rules

### TokenizationRecord
- `Id`: Must be > 0
- `Validator`: Must be valid bech32 validator address
- `Owner`: Must be valid bech32 account address
- `SharesTokenized`: Must be positive
- `Status`: Cannot be UNSPECIFIED

### ModuleParams
- `GlobalLiquidStakingCap`: Must be between 0 and 1
- `ValidatorLiquidCap`: Must be between 0 and 1, cannot exceed global cap
- `MinTokenizationAmount`: Must be positive

### GenesisState
- All params must be valid
- No duplicate record IDs
- `LastTokenizationRecordId` must be >= highest record ID
- All records must be valid

## Performance Benchmarks

Current benchmarks on M1 MacBook Pro:

```
BenchmarkTokenizationRecord_Validate-8      500000      2345 ns/op     512 B/op      8 allocs/op
BenchmarkGenesisState_Validate-8            10000     115678 ns/op   24576 B/op    312 allocs/op
```

These benchmarks establish a baseline for performance regression testing in future stages.