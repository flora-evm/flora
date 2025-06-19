# Liquid Staking Safety Features

This document describes the comprehensive safety mechanisms implemented in the Flora liquid staking module to protect users, validators, and the network from potential risks.

## Overview

The liquid staking module implements multiple layers of safety controls to ensure secure and stable operation:

1. **Capacity Limits** - Prevent excessive concentration of liquid staked tokens
2. **Rate Limiting** - Control the velocity of tokenization activities
3. **Minimum Amount Requirements** - Prevent dust attacks and ensure economic viability
4. **Module Enable/Disable** - Emergency circuit breaker for the entire module

## Capacity Limits

### Global Liquid Staking Cap

The global cap limits the total amount of tokens that can be liquid staked across the entire network.

**Purpose:**
- Prevents excessive concentration of liquid staked tokens
- Maintains network stability by ensuring sufficient non-liquid staked tokens remain
- Protects against systemic risks from over-tokenization

**Configuration:**
- Parameter: `GlobalLiquidStakingCap` (decimal percentage, e.g., 0.25 = 25%)
- Default: 25% of total bonded tokens
- Enforced in: `keeper/caps.go::CheckGlobalLiquidStakingCap()`

**Example:**
```go
// If total bonded = 1,000,000 FLORA and cap = 25%
// Maximum liquid staked allowed = 250,000 FLORA
```

### Per-Validator Liquid Cap

Each validator has an individual cap limiting how much of their delegated tokens can be liquid staked.

**Purpose:**
- Prevents concentration risk at individual validator level
- Ensures validators maintain sufficient regular delegations
- Protects validator operations from excessive liquid staking

**Configuration:**
- Parameter: `ValidatorLiquidCap` (decimal percentage, e.g., 0.50 = 50%)
- Default: 50% of validator's total tokens
- Must be greater than or equal to `GlobalLiquidStakingCap`
- Enforced in: `keeper/caps.go::CheckValidatorLiquidCap()`

**Example:**
```go
// If validator has 200,000 FLORA delegated and cap = 50%
// Maximum liquid staked for this validator = 100,000 FLORA
```

### Minimum Tokenization Amount

Prevents creation of dust liquid staking tokens that could spam the network.

**Purpose:**
- Prevents dust attacks
- Ensures economic viability of liquid staking operations
- Reduces state bloat from tiny positions

**Configuration:**
- Parameter: `MinLiquidStakeAmount` (integer amount)
- Default: 10,000 uflora (0.01 FLORA)
- Enforced in: `keeper/caps.go::CheckMinimumAmount()`

## Rate Limiting

The module implements sophisticated rate limiting to control the velocity of tokenization activities at multiple levels.

### Global Rate Limit

Limits the total amount and frequency of tokenization across the entire network within a 24-hour window.

**Configuration:**
- Daily amount limit: 5% of total bonded tokens
- Daily transaction count: 100 tokenizations
- Window: 24 hours (rolling)
- Enforced in: `keeper/rate_limiter.go::CheckGlobalRateLimit()`

### Per-Validator Rate Limit

Limits tokenization activity for each individual validator.

**Configuration:**
- Daily amount limit: 10% of validator's tokens
- Daily transaction count: 20 tokenizations per validator
- Window: 24 hours (rolling)
- Enforced in: `keeper/rate_limiter.go::CheckValidatorRateLimit()`

### Per-User Rate Limit

Limits how frequently individual users can perform tokenization.

**Configuration:**
- Daily transaction count: 5 tokenizations per user
- Window: 24 hours (rolling)
- Enforced in: `keeper/rate_limiter.go::CheckUserRateLimit()`

### Rate Limit Implementation

Rate limits use a rolling 24-hour window that automatically resets when the period expires:

```go
type TokenizationActivity struct {
    TotalAmount   math.Int  // Total amount tokenized in current window
    LastActivity  time.Time // Timestamp of last activity
    ActivityCount uint64    // Number of tokenizations in current window
}
```

## Cap Enforcement Flow

All safety checks are enforced before any tokenization occurs:

```go
// In keeper/msg_server.go::TokenizeShares()
1. Check module is enabled
2. Check minimum amount requirement
3. Check global liquid staking cap
4. Check validator liquid cap
5. Check global rate limit
6. Check validator rate limit
7. Check user rate limit
8. Proceed with tokenization
9. Update all tracking metrics
```

## Tracking and Monitoring

The module maintains persistent tracking of all liquid staking metrics:

### Global Metrics
- Total liquid staked amount: `GetTotalLiquidStaked()`
- Global tokenization activity: `GetGlobalTokenizationActivity()`

### Per-Validator Metrics
- Validator liquid staked amounts: `GetValidatorLiquidStaked()`
- Validator tokenization activity: `GetValidatorTokenizationActivity()`

### Per-User Metrics
- User tokenization activity: `GetUserTokenizationActivity()`

## Emergency Controls

### Module Enable/Disable

The entire liquid staking module can be disabled through governance:

```go
// Disable module in emergency
params.Enabled = false
```

When disabled:
- No new tokenizations are allowed
- Existing liquid staking tokens remain valid
- Redemptions continue to function
- All other operations are halted

## Parameter Updates

All safety parameters can be updated through governance proposals:

```go
// Example governance proposal to update caps
{
    "title": "Update Liquid Staking Safety Parameters",
    "description": "Increase global cap to 30% due to increased adoption",
    "changes": [
        {
            "subspace": "liquidstaking",
            "key": "GlobalLiquidStakingCap",
            "value": "0.30"
        }
    ]
}
```

## Best Practices

1. **Conservative Defaults**: Start with conservative limits and increase gradually
2. **Monitor Metrics**: Regularly monitor liquid staking metrics and activity patterns
3. **Gradual Updates**: Make parameter changes incrementally to observe effects
4. **Emergency Planning**: Have governance proposals ready for emergency parameter updates
5. **Cross-Chain Considerations**: Account for IBC liquid staking tokens in capacity planning

## Security Considerations

1. **Parameter Validation**: All parameters are validated to prevent invalid configurations
2. **Atomic Checks**: All safety checks are performed atomically before state changes
3. **Persistent State**: All tracking data is persisted to survive node restarts
4. **No Retroactive Changes**: Cap changes don't affect existing liquid staking positions
5. **Governance Control**: Only governance can modify safety parameters

## Monitoring Queries

Monitor safety metrics using these queries:

```bash
# Check current parameters
florad query liquidstaking params

# Check global liquid staking amount
florad query liquidstaking total-liquid-staked

# Check validator liquid staking amount
florad query liquidstaking validator-liquid-staked [validator-address]

# Check if specific amount would exceed caps
florad query liquidstaking check-caps [validator-address] [amount]
```

## Future Enhancements

Potential future safety improvements:
- Dynamic cap adjustment based on network conditions
- Time-based cap increases for gradual rollout
- Delegator-specific caps for additional granularity
- Integration with monitoring and alerting systems
- Automated cap adjustment through on-chain algorithms