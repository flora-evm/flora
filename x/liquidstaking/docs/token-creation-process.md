# Token Creation Process

## Overview

The liquid staking module creates liquid staking tokens (LSTs) during the tokenization process using the Cosmos SDK Bank module.

## Token Creation Flow

### 1. Denom Generation

When shares are tokenized, a unique denom is generated:

```go
denom := types.GenerateLiquidStakingTokenDenom(validatorAddr, recordID)
// Format: "flora/lstake/{validator_address}/{record_id}"
```

### 2. Token Minting

The actual token creation happens through the Bank module:

```go
// Calculate tokens from shares
tokensToMint := validator.TokensFromShares(sharesToTokenize).TruncateInt()

// Mint the liquid staking tokens
mintCoins := sdk.NewCoins(sdk.NewCoin(denom, tokensToMint))
if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, mintCoins); err != nil {
    return nil, err
}
```

### 3. Token Distribution

After minting, tokens are sent to the owner:

```go
if err := k.bankKeeper.SendCoinsFromModuleToAccount(
    ctx, types.ModuleName, ownerAddr, mintCoins
); err != nil {
    return nil, err
}
```

### 4. Metadata Registration

Token metadata is set using the Bank module:

```go
metadata := types.GenerateLiquidStakingTokenMetadata(validatorAddr, recordID)
k.bankKeeper.SetDenomMetaData(ctx, metadata)
```

## Key Design Decisions

### Why Bank Module Instead of Token Factory?

1. **Simplicity**: The Bank module provides all necessary functionality for LST management
2. **Compatibility**: Standard Cosmos SDK patterns work out of the box
3. **Custom Denoms**: Our LST denoms follow a specific format that doesn't align with Token Factory's `factory/{creator}/{subdenom}` pattern
4. **Direct Control**: We need precise control over minting and burning based on staking shares

### Token Properties

- **1:1 Backing**: Each LST is backed by staked FLORA at the time of minting
- **Unique Denoms**: Each tokenization creates a unique denom tied to the validator and record ID
- **Metadata**: Includes description, display name, symbol, and units
- **Supply Tracking**: Managed through tokenization records and liquid staking counters

## Implementation Details

### Minting Process (in msg_server.go)

```go
func (k msgServer) TokenizeShares(goCtx context.Context, msg *types.MsgTokenizeShares) (*types.MsgTokenizeSharesResponse, error) {
    // ... validation ...
    
    // Unbond shares from delegation
    unbondedTokens, err := k.stakingKeeper.Unbond(ctx, delegatorAddr, validatorAddr, sharesToTokenize)
    
    // Mint liquid staking tokens
    mintCoins := sdk.NewCoins(sdk.NewCoin(denom, unbondedTokens))
    if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, mintCoins); err != nil {
        return nil, err
    }
    
    // Send to owner
    if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, ownerAddr, mintCoins); err != nil {
        return nil, err
    }
    
    // Set metadata
    metadata := types.GenerateLiquidStakingTokenMetadata(msg.ValidatorAddress, recordID)
    k.bankKeeper.SetDenomMetaData(ctx, metadata)
    
    // Update tracking
    k.UpdateLiquidStakedAmounts(ctx, msg.ValidatorAddress, unbondedTokens, true)
}
```

### Metadata Generation

```go
func GenerateLiquidStakingTokenMetadata(validatorAddr string, recordID uint64) banktypes.Metadata {
    denom := GenerateLiquidStakingTokenDenom(validatorAddr, recordID)
    
    return banktypes.Metadata{
        Description: fmt.Sprintf("Liquid staking token for validator %s", validatorAddr),
        DenomUnits: []*banktypes.DenomUnit{
            {
                Denom:    denom,
                Exponent: 0,
            },
            {
                Denom:    fmt.Sprintf("mlst%d", recordID), // milli-LST
                Exponent: 3,
            },
        },
        Base:    denom,
        Display: denom,
        Name:    fmt.Sprintf("Liquid Staked FLORA #%d", recordID),
        Symbol:  fmt.Sprintf("lstFLORA%d", recordID),
    }
}
```

## Security Considerations

1. **Module Account**: Only the liquid staking module can mint LSTs
2. **Validation**: Extensive checks before minting (caps, validator status, delegation existence)
3. **Atomic Operations**: Minting and distribution happen in the same transaction
4. **Supply Tracking**: Total and per-validator liquid staked amounts are tracked

## Future Enhancements

1. **Batch Minting**: Support for tokenizing multiple delegations in one transaction
2. **Custom Metadata**: Allow validators to provide custom metadata for their LSTs
3. **Governance Control**: Parameters for controlling metadata standards