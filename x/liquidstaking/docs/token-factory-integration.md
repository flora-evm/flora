# Token Factory Integration Design

## Overview

This document outlines the integration between the Liquid Staking module and the Token Factory module for managing liquid staking tokens (LSTs).

## Token Factory Interface Mismatch

The strangelove-ventures tokenfactory module (v0.50.3) uses a message-based approach rather than direct keeper methods. This requires adapting our integration approach.

### Actual Token Factory Interface

```go
// Message server methods (not direct keeper methods)
- CreateDenom(ctx, creator, subdenom) -> through MsgCreateDenom
- Mint(ctx, mintTo, coin) -> through MsgMint
- Burn(ctx, burnFrom, coin) -> through MsgBurn
- SetDenomMetadata(ctx, creator, metadata) -> through MsgSetDenomMetadata

// Direct keeper methods available
- GetAuthorityMetadata(ctx, denom) -> DenomAuthorityMetadata
```

### Solution: Adapter Pattern

We'll create a wrapper interface that provides the expected methods while internally using the message server approach.

## LST Token Metadata Structure

### Denom Format

Liquid staking tokens follow a standardized format:
```
flora/lstake/{validator_address}/{record_id}
```

Example: `flora/lstake/floravaloper1xyz.../1`

### Metadata Fields

```protobuf
message Metadata {
  string description = 1;
  repeated DenomUnit denom_units = 2;
  string base = 3;
  string display = 4;
  string name = 5;
  string symbol = 6;
  string uri = 7;
  string uri_hash = 8;
}
```

### LST Metadata Design

```go
func GenerateLiquidStakingTokenMetadata(validatorAddr string, recordID uint64) banktypes.Metadata {
    denom := types.GenerateLiquidStakingTokenDenom(validatorAddr, recordID)
    
    // Extract validator moniker for display
    // In production, this would query the validator info
    validatorMoniker := "Validator" // placeholder
    
    return banktypes.Metadata{
        Description: fmt.Sprintf("Liquid staking token for delegation to %s", validatorMoniker),
        DenomUnits: []*banktypes.DenomUnit{
            {
                Denom:    denom,
                Exponent: 0,
                Aliases:  []string{},
            },
            {
                Denom:    fmt.Sprintf("LST-%s-%d", validatorAddr[:8], recordID),
                Exponent: 18,
                Aliases:  []string{},
            },
        },
        Base:    denom,
        Display: fmt.Sprintf("LST-%s-%d", validatorAddr[:8], recordID),
        Name:    fmt.Sprintf("Liquid Staked FLORA - %s #%d", validatorMoniker, recordID),
        Symbol:  fmt.Sprintf("lstFLORA-%d", recordID),
        URI:     fmt.Sprintf("https://flora.network/liquid-staking/%s/%d", validatorAddr, recordID),
        URIHash: "", // Can be added for metadata verification
    }
}
```

## Token Factory Adapter Implementation

### Interface Definition

```go
// TokenFactoryAdapter wraps the token factory message server to provide direct methods
type TokenFactoryAdapter interface {
    CreateDenom(ctx context.Context, creator string, subdenom string) (string, error)
    MintTokens(ctx context.Context, mintTo string, coin sdk.Coin) error
    BurnTokens(ctx context.Context, burnFrom string, coin sdk.Coin) error
    SetTokenMetadata(ctx context.Context, creator string, metadata banktypes.Metadata) error
    GetDenomAuthority(ctx context.Context, denom string) (string, error)
}
```

### Implementation Approach

Since the Token Factory module uses a denom creation approach different from our LST format, we'll:

1. Use the existing Bank module's `MintCoins` and `BurnCoins` for LST management
2. Set metadata using Bank module's `SetDenomMetaData`
3. Skip Token Factory's denom creation since our denoms follow a different pattern
4. Maintain compatibility for future Token Factory features

## Integration Points

### During Tokenization (MsgTokenizeShares)

1. Generate LST denom: `flora/lstake/{validator}/{record_id}`
2. Create token metadata with validator and staking information
3. Mint tokens using Bank module
4. Set denom metadata using Bank module
5. Emit appropriate events

### During Redemption (MsgRedeemTokens)

1. Verify token ownership and balance
2. Burn tokens using Bank module
3. Re-delegate to validator
4. Update or remove tokenization record
5. Emit appropriate events

## Benefits of This Approach

1. **Simplicity**: Uses existing Bank module functionality
2. **Compatibility**: Maintains standard Cosmos SDK patterns
3. **Flexibility**: Can integrate Token Factory features later if needed
4. **Consistency**: Follows established LST denom patterns

## Future Enhancements

1. **Token Factory Integration**: If Token Factory adds features we need, we can integrate them
2. **Custom Permissions**: Could use Token Factory's admin features for governance
3. **Advanced Features**: Force transfers, community pool funding, etc.

## Implementation Notes

- The `TokenFactoryKeeper` field in the Keeper struct can be removed or made optional
- Focus on Bank module integration for core functionality
- Maintain clean separation of concerns between modules