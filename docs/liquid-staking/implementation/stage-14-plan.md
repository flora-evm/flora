# Stage 14: Exchange Rate Updates - Implementation Plan

## Overview
Stage 14 implements dynamic exchange rates for LST tokens, allowing them to appreciate in value as staking rewards accumulate. This creates a mechanism where 1 LST token represents more than 1 staked token over time due to accumulated rewards.

## Technical Design

### 1. Exchange Rate Calculation
The exchange rate represents how many native tokens can be redeemed per LST token:

```
Exchange Rate = (Total Staked + Total Rewards) / Total LST Supply
```

### 2. Key Components

#### a. Exchange Rate Storage
```protobuf
message ExchangeRate {
  string validator_address = 1;
  string denom = 2; // LST denom
  string rate = 3 [
    (cosmos_proto.scalar) = "cosmos.Dec",
    (gogoproto.customtype) = "cosmossdk.io/math.LegacyDec",
    (gogoproto.nullable) = false
  ];
  int64 last_updated = 4; // Unix timestamp
}
```

#### b. Global Rate Tracking
```protobuf
message GlobalExchangeRate {
  string rate = 1 [
    (cosmos_proto.scalar) = "cosmos.Dec",
    (gogoproto.customtype) = "cosmossdk.io/math.LegacyDec",
    (gogoproto.nullable) = false
  ];
  int64 last_updated = 2;
  string total_staked = 3 [
    (cosmos_proto.scalar) = "cosmos.Int",
    (gogoproto.customtype) = "cosmossdk.io/math.Int",
    (gogoproto.nullable) = false
  ];
  string total_rewards = 4 [
    (cosmos_proto.scalar) = "cosmos.Int",
    (gogoproto.customtype) = "cosmossdk.io/math.Int",
    (gogoproto.nullable) = false
  ];
  string total_lst_supply = 5 [
    (cosmos_proto.scalar) = "cosmos.Int",
    (gogoproto.customtype) = "cosmossdk.io/math.Int",
    (gogoproto.nullable) = false
  ];
}
```

### 3. Update Mechanisms

#### a. Manual Update (Stage 14)
```go
// MsgUpdateExchangeRates - manually trigger rate updates
message MsgUpdateExchangeRates {
  string updater = 1; // Must be authorized updater
  repeated string validators = 2; // Empty means update all
}
```

#### b. Rate Calculation Logic
```go
func (k Keeper) CalculateExchangeRate(ctx sdk.Context, validatorAddr string) (math.LegacyDec, error) {
    // Get validator's total delegations
    validator, found := k.stakingKeeper.GetValidator(ctx, valAddr)
    if !found {
        return math.LegacyZeroDec(), ErrValidatorNotFound
    }
    
    // Get total LST tokens for this validator
    lstDenom := GetLSTDenom(validatorAddr)
    totalLST := k.bankKeeper.GetSupply(ctx, lstDenom)
    
    // Get accumulated rewards
    rewards := k.distributionKeeper.GetValidatorAccumulatedRewards(ctx, valAddr)
    
    // Calculate rate
    totalValue := validator.TokensFromShares(validator.DelegatorShares).Add(rewards.AmountOf(bondDenom))
    
    if totalLST.IsZero() {
        return math.LegacyOneDec(), nil // Initial rate is 1:1
    }
    
    return totalValue.Quo(totalLST), nil
}
```

### 4. Integration Points

#### a. Tokenization Updates
```go
func (k Keeper) TokenizeDelegation(ctx sdk.Context, ...) error {
    // Get current exchange rate
    rate, err := k.GetExchangeRate(ctx, validatorAddr)
    if err != nil {
        return err
    }
    
    // Calculate LST tokens to mint based on rate
    lstAmount := amount.Quo(rate.TruncateInt())
    
    // Mint LST tokens at current rate
    // ...
}
```

#### b. Redemption Updates
```go
func (k Keeper) RedeemTokens(ctx sdk.Context, ...) error {
    // Get current exchange rate
    rate, err := k.GetExchangeRate(ctx, validatorAddr)
    if err != nil {
        return err
    }
    
    // Calculate native tokens to return based on rate
    nativeAmount := lstAmount.Mul(rate.TruncateInt())
    
    // Process redemption at current rate
    // ...
}
```

### 5. Query Updates

#### a. Exchange Rate Query
```protobuf
service Query {
  // ExchangeRate returns the current exchange rate for a validator
  rpc ExchangeRate(QueryExchangeRateRequest) returns (QueryExchangeRateResponse);
  
  // AllExchangeRates returns all exchange rates
  rpc AllExchangeRates(QueryAllExchangeRatesRequest) returns (QueryAllExchangeRatesResponse);
}
```

### 6. Events

```go
const (
    EventTypeExchangeRateUpdated = "exchange_rate_updated"
    
    AttributeKeyValidator = "validator"
    AttributeKeyOldRate = "old_rate"
    AttributeKeyNewRate = "new_rate"
    AttributeKeyTimestamp = "timestamp"
)
```

## Implementation Steps

1. **Add Exchange Rate Types** (Day 1)
   - Define ExchangeRate and GlobalExchangeRate protos
   - Generate protobuf code
   - Add validation logic

2. **Implement Storage** (Day 1-2)
   - Add keeper methods for rate storage/retrieval
   - Implement rate calculation logic
   - Add update authorization checks

3. **Update Tokenization Logic** (Day 2-3)
   - Modify TokenizeDelegation to use exchange rates
   - Update LST minting calculations
   - Add rate validation

4. **Update Redemption Logic** (Day 3-4)
   - Modify RedeemTokens to use exchange rates
   - Update native token calculations
   - Ensure precision handling

5. **Add Manual Update Message** (Day 4)
   - Implement MsgUpdateExchangeRates
   - Add authorization checks
   - Emit update events

6. **Implement Queries** (Day 5)
   - Add exchange rate queries
   - Update existing queries to show rates
   - Add rate history tracking

7. **Testing** (Day 5-6)
   - Unit tests for rate calculations
   - Integration tests for tokenization/redemption
   - Precision and edge case tests

## Testing Strategy

### Unit Tests
- Rate calculation accuracy
- Precision handling (no dust)
- Edge cases (zero supply, first mint)
- Authorization validation

### Integration Tests
- Full tokenization flow with rates
- Full redemption flow with rates
- Rate updates and effects
- Multi-validator scenarios

### Scenarios to Test
1. Initial tokenization (1:1 rate)
2. Tokenization after rewards accumulate
3. Redemption with appreciated rate
4. Multiple rate updates
5. Precision with large numbers
6. Rate calculation with slashing

## Security Considerations

1. **Rate Manipulation**: Only authorized updaters can trigger updates
2. **Precision Loss**: Use sdk.Dec for calculations, careful truncation
3. **Race Conditions**: Lock rates during tokenization/redemption
4. **Oracle Reliability**: Manual updates in Stage 14, automated in Stage 15

## Success Criteria

- [ ] Exchange rates calculate correctly based on validator state
- [ ] Tokenization uses current exchange rate
- [ ] Redemption uses current exchange rate
- [ ] No precision loss in calculations
- [ ] Rate updates emit proper events
- [ ] Queries return accurate rate information
- [ ] All tests pass with 100% coverage

## Next Steps

After Stage 14 completion:
- Stage 15: Auto-compound rewards using exchange rates
- Stage 16: Slashing protection and rate adjustments
- Stage 17: Governance integration for rate parameters
- Stage 18: IBC transfer of LST tokens