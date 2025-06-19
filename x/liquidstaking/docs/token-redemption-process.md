# Token Redemption and Burning Process

## Overview

The liquid staking module handles redemption of liquid staking tokens (LSTs) by burning them and re-delegating the underlying staked assets back to the validator.

## Redemption Flow

### 1. Validation

Before redemption can proceed, several checks are performed:

```go
// Verify tokenization record exists
recordID, found := k.getTokenizationRecordByDenom(ctx, msg.Amount.Denom)
if !found {
    return nil, ErrTokenizationRecordNotFound
}

// Verify ownership
if record.Owner != msg.OwnerAddress {
    return nil, ErrUnauthorized
}

// Check balance
balance := k.bankKeeper.GetBalance(ctx, ownerAddr, msg.Amount.Denom)
if balance.IsLT(msg.Amount) {
    return nil, ErrInsufficientFunds
}
```

### 2. Token Burning

The burning process involves two steps:

```go
// Step 1: Transfer tokens from user to module account
burnCoins := sdk.NewCoins(msg.Amount)
if err := k.bankKeeper.SendCoinsFromAccountToModule(
    ctx, ownerAddr, types.ModuleName, burnCoins
); err != nil {
    return nil, err
}

// Step 2: Burn tokens from module account
if err := k.bankKeeper.BurnCoins(ctx, types.ModuleName, burnCoins); err != nil {
    return nil, err
}
```

### 3. Re-delegation

After burning, the equivalent shares are re-delegated:

```go
// Calculate shares based on current exchange rate
sharesToRestore, _ := validator.SharesFromTokens(msg.Amount.Amount)

// Re-delegate to the original validator
_, err = k.stakingKeeper.Delegate(
    ctx, ownerAddr, msg.Amount.Amount, 
    stakingtypes.Unbonded, validator, true
)
```

### 4. Record Update

The tokenization record is updated or deleted:

```go
if record.SharesTokenized.Sub(msg.Amount.Amount).IsZero() {
    // Full redemption - delete record and indexes
    k.DeleteTokenizationRecordWithIndexes(ctx, recordID)
    k.deleteTokenizationRecordDenomIndex(ctx, msg.Amount.Denom)
} else {
    // Partial redemption - update record
    record.SharesTokenized = record.SharesTokenized.Sub(msg.Amount.Amount)
    k.SetTokenizationRecordWithIndexes(ctx, record)
}
```

### 5. State Updates

Finally, liquid staking counters are updated:

```go
// Update liquid staked amounts
k.UpdateLiquidStakedAmounts(ctx, record.Validator, msg.Amount.Amount, false)
```

## Key Design Decisions

### Two-Step Burning Process

The burning happens in two steps for security:
1. **Transfer to Module**: Ensures user has the tokens
2. **Burn from Module**: Removes tokens from circulation

This prevents edge cases where tokens might be burned without proper ownership verification.

### Re-delegation vs Direct Staking

Tokens are re-delegated directly to maintain:
- Validator selection continuity
- Proper accounting of liquid vs regular staking
- Simplified user experience

### Partial Redemption Support

Users can redeem any amount of their LSTs:
- Flexibility for users
- Maintains tokenization record for remaining balance
- Updates all indexes appropriately

## Implementation Details

### Complete Redemption Function

```go
func (k msgServer) RedeemTokens(goCtx context.Context, msg *types.MsgRedeemTokens) (*types.MsgRedeemTokensResponse, error) {
    ctx := sdk.UnwrapSDKContext(goCtx)
    
    // Validation
    ownerAddr, _ := sdk.AccAddressFromBech32(msg.OwnerAddress)
    recordID, found := k.getTokenizationRecordByDenom(ctx, msg.Amount.Denom)
    if !found {
        return nil, ErrTokenizationRecordNotFound
    }
    
    record, _ := k.GetTokenizationRecord(ctx, recordID)
    if record.Owner != msg.OwnerAddress {
        return nil, ErrUnauthorized
    }
    
    // Check balance
    balance := k.bankKeeper.GetBalance(ctx, ownerAddr, msg.Amount.Denom)
    if balance.IsLT(msg.Amount) {
        return nil, ErrInsufficientFunds
    }
    
    // Burn tokens
    burnCoins := sdk.NewCoins(msg.Amount)
    k.bankKeeper.SendCoinsFromAccountToModule(ctx, ownerAddr, types.ModuleName, burnCoins)
    k.bankKeeper.BurnCoins(ctx, types.ModuleName, burnCoins)
    
    // Re-delegate
    validator, _ := k.stakingKeeper.GetValidator(ctx, valAddr)
    k.stakingKeeper.Delegate(ctx, ownerAddr, msg.Amount.Amount, stakingtypes.Unbonded, validator, true)
    
    // Update record
    if record.SharesTokenized.Sub(msg.Amount.Amount).IsZero() {
        k.DeleteTokenizationRecordWithIndexes(ctx, recordID)
        k.deleteTokenizationRecordDenomIndex(ctx, msg.Amount.Denom)
    } else {
        record.SharesTokenized = record.SharesTokenized.Sub(msg.Amount.Amount)
        k.SetTokenizationRecordWithIndexes(ctx, record)
    }
    
    // Update counters
    k.UpdateLiquidStakedAmounts(ctx, record.Validator, msg.Amount.Amount, false)
    
    return &types.MsgRedeemTokensResponse{
        Shares:   sharesToRestore,
        RecordId: recordID,
    }, nil
}
```

## Security Considerations

1. **Ownership Verification**: Only the record owner can redeem tokens
2. **Balance Checks**: Prevents redemption of tokens not owned
3. **Atomic Operations**: All steps happen in one transaction
4. **Module Account Control**: Burning only through module account

## Events Emitted

```go
EventTypeRedeemTokens
- owner: Address of the token owner
- validator: Validator address
- denom: LST denom being redeemed  
- amount: Amount of tokens burned
- shares: Amount of shares restored
- record_id: Tokenization record ID
```

## Future Enhancements

1. **Batch Redemption**: Support redeeming multiple LST types in one transaction
2. **Instant Redemption**: Option to receive liquid FLORA instead of re-delegating
3. **Fee Structure**: Optional redemption fees for instant redemption