# Liquid Staking CLI Usage Guide

## Overview

The liquid staking module provides CLI commands to interact with the liquid staking functionality. This guide covers all available commands for querying state and submitting transactions.

## Query Commands

All query commands are available under:
```bash
florad query liquidstaking [command]
```

### Query Module Parameters

Get the current module parameters including caps and minimum amounts:

```bash
florad query liquidstaking params
```

Example output:
```json
{
  "enabled": true,
  "min_liquid_stake_amount": "1000000",
  "global_liquid_staking_cap": "0.250000000000000000",
  "validator_liquid_cap": "0.100000000000000000"
}
```

### Query Tokenization Record

Query a specific tokenization record by ID:

```bash
florad query liquidstaking record [record-id]
```

Example:
```bash
florad query liquidstaking record 1
```

### Query All Tokenization Records

List all tokenization records with optional pagination:

```bash
florad query liquidstaking records
```

With pagination:
```bash
florad query liquidstaking records --page=2 --limit=20
```

### Query Records by Owner

Find all tokenization records owned by a specific address:

```bash
florad query liquidstaking records-by-owner [owner-address]
```

Example:
```bash
florad query liquidstaking records-by-owner flora1abc...xyz
```

### Query Records by Validator

Find all tokenization records for a specific validator:

```bash
florad query liquidstaking records-by-validator [validator-address]
```

Example:
```bash
florad query liquidstaking records-by-validator floravaloper1abc...xyz
```

### Query Total Liquid Staked

Get the total amount of tokens liquid staked across all validators:

```bash
florad query liquidstaking total-liquid-staked
```

### Query Validator Liquid Staked

Get the amount of tokens liquid staked for a specific validator:

```bash
florad query liquidstaking validator-liquid-staked [validator-address]
```

Example:
```bash
florad query liquidstaking validator-liquid-staked floravaloper1abc...xyz
```

## Transaction Commands

All transaction commands are available under:
```bash
florad tx liquidstaking [command]
```

### Tokenize Shares

Convert delegation shares into liquid staking tokens:

```bash
florad tx liquidstaking tokenize-shares [validator-address] [amount] [flags]
```

Arguments:
- `validator-address`: The validator whose delegation shares to tokenize
- `amount`: The amount of shares to tokenize (in share denomination)

Flags:
- `--owner`: Optional owner address for the liquid staking tokens (defaults to delegator)
- `--from`: Required transaction signer (must be the delegator)

Examples:

1. Tokenize shares with tokens sent to delegator:
```bash
florad tx liquidstaking tokenize-shares floravaloper1abc...xyz 1000000stake \
  --from mykey \
  --gas auto \
  --gas-adjustment 1.5 \
  --gas-prices 0.025flora
```

2. Tokenize shares with tokens sent to a different owner:
```bash
florad tx liquidstaking tokenize-shares floravaloper1abc...xyz 1000000stake \
  --owner flora1def...uvw \
  --from mykey \
  --gas auto \
  --gas-adjustment 1.5 \
  --gas-prices 0.025flora
```

### Redeem Tokens

Burn liquid staking tokens to restore delegation shares:

```bash
florad tx liquidstaking redeem-tokens [amount] [flags]
```

Arguments:
- `amount`: The amount of liquid staking tokens to redeem (must include denomination)

Flags:
- `--from`: Required transaction signer (must be the token owner)

Example:
```bash
florad tx liquidstaking redeem-tokens 1000000liquidstake/floravaloper1abc...xyz/1 \
  --from mykey \
  --gas auto \
  --gas-adjustment 1.5 \
  --gas-prices 0.025flora
```

## Common Usage Patterns

### 1. Check Liquid Staking Availability

Before tokenizing shares, check if the module is enabled and has capacity:

```bash
# Check if module is enabled
florad query liquidstaking params

# Check global cap usage
florad query liquidstaking total-liquid-staked

# Check validator cap usage
florad query liquidstaking validator-liquid-staked floravaloper1abc...xyz
```

### 2. Tokenize Delegation

```bash
# First check your delegation
florad query staking delegation flora1abc...xyz floravaloper1abc...xyz

# Tokenize a portion of your delegation
florad tx liquidstaking tokenize-shares floravaloper1abc...xyz 500000stake \
  --from mykey \
  --gas auto
```

### 3. Transfer and Redeem LST Tokens

```bash
# Check your LST balance
florad query bank balances flora1abc...xyz

# Transfer LST tokens to another account
florad tx bank send flora1abc...xyz flora1def...uvw 100000liquidstake/floravaloper1abc...xyz/1 \
  --from mykey

# Redeem LST tokens
florad tx liquidstaking redeem-tokens 100000liquidstake/floravaloper1abc...xyz/1 \
  --from mykey
```

### 4. Monitor Tokenization Records

```bash
# View all your tokenization records
florad query liquidstaking records-by-owner flora1abc...xyz

# Check details of a specific record
florad query liquidstaking record 1
```

## Important Notes

1. **Liquid Staking Token Format**: LST tokens have the format `liquidstake/{validator}/{record-id}`
   - Example: `liquidstake/floravaloper1abc...xyz/1`

2. **Share vs Token Amounts**: 
   - When tokenizing, specify the amount in shares
   - When redeeming, specify the amount in tokens with full denomination

3. **Ownership**: 
   - Only the delegator can tokenize their shares
   - Only the owner can redeem LST tokens
   - Ownership can be transferred via standard bank send

4. **Caps and Limits**:
   - Check global and per-validator caps before tokenizing
   - Transactions will fail if caps are exceeded

5. **Gas Estimation**:
   - Use `--gas auto` for automatic gas estimation
   - Liquid staking operations may require more gas than standard transfers

## Error Handling

Common errors and their meanings:

- `module is disabled`: The liquid staking module is currently disabled
- `delegation not found`: No delegation exists for the specified validator
- `insufficient shares`: Trying to tokenize more shares than available
- `exceeds global cap`: Would exceed the global liquid staking cap
- `exceeds validator cap`: Would exceed the per-validator liquid staking cap
- `tokenization record not found`: Invalid record ID or denomination
- `unauthorized`: Not the owner of the tokenization record

## Integration with Other Modules

### Bank Module
- Use standard bank commands to transfer LST tokens
- Query LST token balances with bank balance queries
- LST tokens are fungible within the same validator/record

### Staking Module
- Delegations are reduced when tokenizing shares
- New delegations are created when redeeming tokens
- Rewards continue to accrue to the original delegator

### IBC Module
- LST tokens can be transferred via IBC (if enabled)
- Maintain the same denomination format across chains