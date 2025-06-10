# Liquid Staking Module

## Overview

The liquid staking module enables users to tokenize their staked assets while maintaining the security and rewards of the staking system. This module is being implemented in a staged approach over 20 weeks, with each stage building upon the previous functionality.

## Current Status: Stage 2 Complete ✅

### Stage 1: Basic Infrastructure (Week 1) - COMPLETED

**Objective**: Establish the foundational module structure with basic types and minimal keeper functionality.

**Completed Components**:

1. **Module Structure**
   - Created standard Cosmos SDK module directory layout
   - Implemented `AppModuleBasic` and `AppModule` interfaces
   - Integrated with app.go

2. **Core Types**
   - `TokenizationRecord`: Tracks tokenized stake positions
     - ID, Validator, Owner, SharesTokenized
     - Full validation logic
   - `ModuleParams`: Module-wide parameters
     - GlobalLiquidStakingCap (default: 25%)
     - ValidatorLiquidCap (default: 50%)
     - Enabled flag

3. **Keeper Implementation**
   - Basic keeper structure with store service
   - Parameter management (Get/Set)
   - TokenizationRecord CRUD operations
   - Genesis import/export functionality

4. **Protobuf Integration**
   - Generated types from proto definitions
   - Fixed import paths for cosmossdk.io/math types
   - Proper codec registration

5. **Testing**
   - Comprehensive unit tests for all types
   - Keeper functionality tests
   - Genesis validation tests
   - All tests passing (30/30)

### Technical Decisions

1. **Storage Design**
   - Using protobuf for all persisted types
   - KVStore keys: params, tokenization records, last ID counter
   - Efficient iteration patterns for getAllRecords

2. **Type Safety**
   - Using cosmossdk.io/math.Int for numeric values
   - Strict validation on all user inputs
   - Bech32 address validation with proper prefixes

3. **Module Integration**
   - Following tokenfactory module pattern
   - Minimal dependencies for Stage 1
   - Prepared hooks for future staking integration

## Usage (Stage 1)

Currently, the module provides basic infrastructure only. Transaction handling will be added in Stage 3.

### Genesis Configuration

```json
{
  "liquidstaking": {
    "params": {
      "global_liquid_staking_cap": "0.250000000000000000",
      "validator_liquid_cap": "0.500000000000000000",
      "enabled": true
    },
    "tokenization_records": [],
    "last_tokenization_record_id": "0"
  }
}
```

### Parameters

- `global_liquid_staking_cap`: Maximum percentage of total staked tokens that can be liquid staked
- `validator_liquid_cap`: Maximum percentage of a validator's stake that can be liquid staked
- `enabled`: Module enable/disable flag

## Development Roadmap

### ✅ Stage 1: Basic Infrastructure (Week 1) - COMPLETED
- Module structure and basic types
- Minimal keeper implementation
- Genesis handling
- Unit tests

### ✅ Stage 2: State Management (Week 2) - COMPLETED
- Extended keeper with tokenization record operations
- Added indexed queries (by validator, by owner)
- Implemented validation logic for liquid staking caps
- Created comprehensive test suite (38 tests, all passing)
- Added proto query service definitions
- Implemented state aggregation (total and per-validator tracking)

### 🚀 Stage 3: Basic Tokenization (Weeks 3-4) - NEXT
- MsgTokenizeShares implementation
- Basic minting of liquid staking tokens
- Event emission
- E2E tests

### Stage 4: Redemption Mechanism (Weeks 5-6)
- MsgRedeemTokensforShares
- Unbonding period handling
- State transitions

### Future Stages (Weeks 7-20)
- Reward distribution
- Slashing handling
- Governance integration
- IBC compatibility
- Advanced features

## Testing

Run all module tests:
```bash
go test ./x/liquidstaking/...
```

Run specific test suites:
```bash
go test ./x/liquidstaking/types -v
go test ./x/liquidstaking/keeper -v
```

## Contributing

This module is under active development. Please refer to the Flora contribution guidelines.

## References

- [Cosmos SDK Staking Module](https://docs.cosmos.network/main/modules/staking)
- [Liquid Staking Research](../docs/liquid-staking/)
- [Implementation Plan](../docs/liquid-staking/implementation/)