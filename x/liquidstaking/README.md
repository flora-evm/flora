# Liquid Staking Module

The liquid staking module enables users to tokenize their staked assets while maintaining the security and decentralization of the network. Users receive liquid staking tokens (LST) that represent their staked position and can be freely transferred, traded, or used in DeFi applications.

## Overview

Liquid staking solves the capital efficiency problem in Proof of Stake networks where staked tokens are locked and cannot be used for other purposes. With this module, users can:

- **Tokenize staked assets** - Convert delegated stakes into fungible LST tokens
- **Maintain staking rewards** - Continue earning staking rewards while holding LST
- **Unlock liquidity** - Use LST tokens in DeFi while contributing to network security
- **Redeem anytime** - Convert LST tokens back to staked tokens

## Features

### Core Functionality
- **Tokenization** - Convert delegated shares to liquid staking tokens
- **Redemption** - Burn LST tokens to restore delegated shares
- **Exchange Rates** - Track value of LST tokens relative to native tokens
- **IBC Compatible** - Transfer LST tokens across IBC-enabled chains

### Safety Mechanisms
- **Liquid Staking Caps** - Global and per-validator limits prevent concentration
- **Rate Limiting** - Daily limits on tokenization activity
- **Emergency Controls** - Pause functionality for crisis management
- **Validator Controls** - Whitelist/blacklist for validator participation

### Advanced Features
- **Auto-compound** - Automatic reinvestment of staking rewards
- **Governance Integration** - Parameter updates via governance proposals
- **Hooks System** - Extensible architecture for custom logic
- **Multi-validator Support** - Tokenize stakes from multiple validators

## Architecture

```
User -> MsgTokenizeShares -> Keeper -> Staking Module
                                    -> Bank Module (Mint LST)
                                    -> Store (Record)
                                    
User -> MsgRedeemTokens -> Keeper -> Bank Module (Burn LST)
                                  -> Staking Module
                                  -> Store (Update/Delete Record)
```

## Usage

### CLI Commands

#### Tokenize Shares
```bash
# Tokenize 1000 FLORA from a delegation
florad tx liquidstaking tokenize-shares 1000flora floravaloper1... --from mykey

# Tokenize to a different owner
florad tx liquidstaking tokenize-shares 1000flora floravaloper1... --owner flora1... --from mykey
```

#### Redeem Tokens
```bash
# Redeem liquid staking tokens
florad tx liquidstaking redeem-tokens 1000liquidstake/floravaloper1.../1 --from mykey
```

#### Query Commands
```bash
# Get module parameters
florad query liquidstaking params

# List all tokenization records
florad query liquidstaking tokenization-records

# Get specific record
florad query liquidstaking tokenization-record 1

# Get exchange rate for a validator
florad query liquidstaking exchange-rate floravaloper1...

# Check liquid staked amount
florad query liquidstaking total-liquid-staked
florad query liquidstaking validator-liquid-staked floravaloper1...
```

### Governance

#### Update Parameters
```bash
# Create parameter change proposal
florad tx gov submit-proposal update-params liquidstaking \
  "Update Liquid Staking Params" \
  "Enable auto-compound feature" \
  1000flora \
  params.json \
  --from mykey

# params.json example:
[
  {
    "key": "auto_compound_enabled",
    "value": "true"
  },
  {
    "key": "auto_compound_frequency_blocks",
    "value": "28800"
  }
]
```

#### Emergency Pause
```bash
# Submit emergency pause proposal
florad tx gov submit-proposal emergency-pause \
  "Emergency Pause Liquid Staking" \
  "Critical vulnerability discovered" \
  1000flora \
  true \
  86400 \
  --from mykey
```

### Integration

For developers integrating with the liquid staking module:

```go
import (
    liquidstakingkeeper "github.com/rollchains/flora/x/liquidstaking/keeper"
    liquidstakingtypes "github.com/rollchains/flora/x/liquidstaking/types"
)

// Query exchange rate
rate, found := keeper.GetExchangeRate(ctx, validatorAddr)

// Check if validator is allowed
allowed := keeper.IsValidatorAllowed(ctx, validatorAddr)

// Get liquid staked amount
amount := keeper.GetValidatorLiquidStakedAmount(ctx, validatorAddr)
```

## Parameters

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `enabled` | bool | true | Enable/disable the module |
| `global_liquid_staking_cap` | Dec | 0.25 | Maximum % of total supply that can be liquid staked |
| `validator_liquid_cap` | Dec | 0.50 | Maximum % of validator's stake that can be liquid |
| `min_liquid_stake_amount` | Int | 1000000 | Minimum amount to tokenize |
| `rate_limit_period_hours` | uint32 | 24 | Rate limit window in hours |
| `global_daily_tokenization_percent` | Dec | 0.10 | Daily global tokenization limit (%) |
| `validator_daily_tokenization_percent` | Dec | 0.10 | Daily per-validator limit (%) |
| `global_daily_tokenization_count` | uint64 | 100 | Max daily tokenization transactions globally |
| `validator_daily_tokenization_count` | uint64 | 10 | Max daily transactions per validator |
| `user_daily_tokenization_count` | uint64 | 5 | Max daily transactions per user |
| `warning_threshold_percent` | Dec | 0.90 | Emit warning when cap usage exceeds this |
| `auto_compound_enabled` | bool | false | Enable automatic reward compounding |
| `auto_compound_frequency_blocks` | int64 | 28800 | Blocks between auto-compound runs |
| `max_rate_change_per_update` | Dec | 0.01 | Maximum exchange rate change per update |
| `min_blocks_between_updates` | int64 | 100 | Minimum blocks between rate updates |

## Events

The module emits detailed events for all operations:

- `tokenize_shares` - Emitted when shares are tokenized
- `redeem_tokens` - Emitted when tokens are redeemed  
- `exchange_rate_updated` - Emitted when exchange rate changes
- `liquid_staking_cap_exceeded` - Warning when approaching caps
- `parameter_updated` - Emitted on parameter changes

## Security Considerations

1. **Validator Risk** - LST tokens inherit the risk of the underlying validator
2. **Slashing** - Validator slashing affects LST token value
3. **Exchange Rate** - Rate changes affect redemption value
4. **Concentration** - Caps prevent too much stake becoming liquid

## FAQ

**Q: What happens to my staking rewards?**
A: Staking rewards continue to accrue and are reflected in the exchange rate. When auto-compound is enabled, rewards are automatically restaked.

**Q: Can I vote with LST tokens?**
A: No, governance voting rights remain with the original delegator until tokens are redeemed.

**Q: What if a validator gets slashed?**
A: The LST token value decreases proportionally to the slashing penalty.

**Q: Can I transfer LST tokens via IBC?**
A: Yes, LST tokens are IBC-compatible and can be transferred to other chains.

**Q: How is the exchange rate calculated?**
A: Exchange rate = (Validator's total tokens + rewards) / Total shares

## Contributing

We welcome contributions! Please see our contributing guidelines and submit PRs to our GitHub repository.

## License

This module is licensed under the Apache 2.0 License.