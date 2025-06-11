# Liquid Staking Module

## Overview

The liquid staking module enables users to tokenize their staked assets while maintaining the security and rewards of the staking system. This module is being implemented in a staged approach over 20 weeks, with each stage building upon the previous functionality.

## Current Status: Stage 5 Complete ✅

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

## Usage

The liquid staking module now supports tokenization and redemption of staked assets.

### Tokenizing Shares

Convert your delegated shares into liquid staking tokens:

```bash
florad tx liquidstaking tokenize-shares [delegator] [validator] [amount] --from [key]

# Example: Tokenize 1000 shares
florad tx liquidstaking tokenize-shares flora1... floravaloper1... 1000shares --from mykey
```

### Redeeming Tokens

Convert liquid staking tokens back to regular delegations:

```bash
florad tx liquidstaking redeem-tokens [amount] --from [key]

# Example: Redeem 500 liquid staking tokens
florad tx liquidstaking redeem-tokens 500flora/lstake/floravaloper1.../1 --from mykey
```

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

- `global_liquid_staking_cap`: Maximum percentage of total staked tokens that can be liquid staked (default: 25%)
- `validator_liquid_cap`: Maximum percentage of a validator's stake that can be liquid staked (default: 50%)
- `enabled`: Module enable/disable flag

## Staking Integration Details

### Tokenization Flow
1. **Validation Phase**
   - Check module is enabled
   - Validate addresses and amounts
   - Verify delegation exists with sufficient shares
   - Ensure validator is not jailed
   - Check liquid staking caps won't be exceeded

2. **Execution Phase**
   - Unbond shares from validator
   - Generate unique LST denomination
   - Mint liquid staking tokens to owner
   - Create and index tokenization record
   - Update liquid staking statistics
   - Emit tokenization events

### Redemption Flow  
1. **Validation Phase**
   - Verify token ownership
   - Check sufficient LST balance
   - Validate tokenization record exists

2. **Execution Phase**
   - Burn liquid staking tokens
   - Re-delegate shares to original validator
   - Update or delete tokenization record
   - Update liquid staking statistics
   - Emit redemption events

### Safety Features
- **Validator Eligibility**: Only non-jailed validators
- **Cap Enforcement**: Global and per-validator limits
- **Address Validation**: Comprehensive bech32 validation
- **Amount Validation**: Positive, non-zero amounts only
- **State Consistency**: Atomic operations with rollback on failure

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

### ✅ Stage 3: Basic Tokenization (Weeks 3-4) - COMPLETED
- Implemented MsgTokenizeShares with full validation
- Integration with staking module for unbonding shares
- Unique liquid staking token denomination generation
- Bank module integration for token minting
- Comprehensive event emission
- Full test coverage with mock keepers

### ✅ Stage 4: Redemption Flow (Week 5) - COMPLETED  
- Implemented MsgRedeemTokens for converting LSTs back to delegations
- Token burning and re-delegation logic
- Partial and full redemption support
- Record lifecycle management (update/delete)
- Event emission for redemption tracking
- Complete test coverage for edge cases

### ✅ Stage 5: Integration with Staking Module (Week 6) - COMPLETED
- Deep integration with Cosmos SDK staking module
- Mock staking keeper for comprehensive testing
- Edge case handling:
  - Jailed validators
  - Insufficient delegations  
  - Unbonding/unbonded validators
  - Validator commission changes
  - Slashed validators
- Liquid staking cap enforcement (global and per-validator)
- Validation helper functions for code reuse
- Enhanced error messages and documentation

### 🚀 Stage 6: Token Factory Integration (Week 7) - NEXT
- Integration with Token Factory module
- Custom denomination metadata
- Enhanced token creation process

### Future Stages (Weeks 8-20)
- Stage 7: EVM Precompiles (Week 8)
- Stage 8: Reward Distribution (Weeks 9-10)
- Stage 9: Slashing Handling (Week 11)
- Stage 10: Unbonding Period Management (Week 12)
- Stage 11: Governance Integration (Week 13)
- Stage 12: Query Improvements (Week 14)
- Stage 13: IBC Compatibility (Week 15)
- Stage 14: CLI Enhancements (Week 16)
- Stage 15: Performance Optimization (Week 17)
- Stage 16: Security Audit Prep (Week 18)
- Stage 17: Documentation (Week 19)
- Stage 18: Mainnet Preparation (Week 20)

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

Run integration tests:
```bash
# Staking integration tests
go test ./x/liquidstaking/keeper -run "TestStakingIntegration" -v

# Validator state tests
go test ./x/liquidstaking/keeper -run "TestValidatorState" -v

# Unbonding tests
go test ./x/liquidstaking/keeper -run "TestUnbonding" -v
```

Current test coverage:
- Types: 100% coverage
- Keeper: 95%+ coverage
- Integration: Comprehensive edge case coverage

## Contributing

This module is under active development. Please refer to the Flora contribution guidelines.

## References

- [Cosmos SDK Staking Module](https://docs.cosmos.network/main/modules/staking)
- [Liquid Staking Research](../docs/liquid-staking/)
- [Implementation Plan](../docs/liquid-staking/implementation/)