# Stage 3: Basic Tokenization Implementation Plan

## Overview

Stage 3 implements the core tokenization functionality, allowing users to convert their staked assets into liquid staking tokens. This stage introduces the first user-facing transaction and integrates with the bank and staking modules.

## Timeline: Weeks 3-4

## Objectives

1. Implement MsgTokenizeShares transaction
2. Create liquid staking token minting logic
3. Integrate with bank and staking modules
4. Add comprehensive event emission
5. Build end-to-end tests

## Implementation Tasks

### 3.1 Message Definition

Create `MsgTokenizeShares` in `tx.proto`:

```proto
message MsgTokenizeShares {
  string delegator_address = 1;
  string validator_address = 2;
  cosmos.base.v1beta1.Coin shares = 3;
  string owner_address = 4; // can be different from delegator
}

message MsgTokenizeSharesResponse {
  string denom = 1; // liquid staking token denom
  cosmos.base.v1beta1.Coin amount = 2; // minted tokens
  uint64 record_id = 3; // tokenization record ID
}
```

### 3.2 Expected Keeper Interfaces

Define interfaces in `expected_keepers.go`:

```go
type StakingKeeper interface {
    GetDelegation(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (stakingtypes.Delegation, bool)
    GetValidator(ctx sdk.Context, addr sdk.ValAddress) (stakingtypes.Validator, bool)
    Unbond(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, shares math.LegacyDec) (math.Int, error)
    GetParams(ctx sdk.Context) stakingtypes.Params
    BondDenom(ctx sdk.Context) string
}

type BankKeeper interface {
    MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
    SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
    GetDenomMetaData(ctx sdk.Context, denom string) (banktypes.Metadata, bool)
    SetDenomMetaData(ctx sdk.Context, denomMetaData banktypes.Metadata)
}

type AccountKeeper interface {
    GetAccount(ctx sdk.Context, addr sdk.AccAddress) authtypes.AccountI
    GetModuleAddress(moduleName string) sdk.AccAddress
}
```

### 3.3 Tokenization Logic

Implement in `keeper/msg_server.go`:

```go
func (k msgServer) TokenizeShares(goCtx context.Context, msg *types.MsgTokenizeShares) (*types.MsgTokenizeSharesResponse, error) {
    ctx := sdk.UnwrapSDKContext(goCtx)
    
    // 1. Validate message
    // 2. Check delegation exists
    // 3. Validate liquid staking caps
    // 4. Create tokenization record
    // 5. Generate LST denom
    // 6. Unbond shares from staking
    // 7. Mint liquid staking tokens
    // 8. Send tokens to owner
    // 9. Update indexes and state
    // 10. Emit events
    
    return response, nil
}
```

### 3.4 Denomination Generation

Liquid staking token denominations follow the pattern:
- Format: `liquidstake/{validator_address}/{record_id}`
- Example: `liquidstake/floravaloper1.../1`

This ensures uniqueness and traceability.

### 3.5 Token Metadata

Set metadata for liquid staking tokens:

```go
metadata := banktypes.Metadata{
    Description: fmt.Sprintf("Liquid staking token for %s", validatorName),
    DenomUnits: []*banktypes.DenomUnit{
        {Denom: denom, Exponent: 0},
    },
    Base:    denom,
    Display: denom,
    Name:    fmt.Sprintf("Liquid Staked %s", bondDenom),
    Symbol:  fmt.Sprintf("ls%s", bondDenom),
}
```

### 3.6 Event Emission

Define events in `types/events.go`:

```go
const (
    EventTypeTokenizeShares = "tokenize_shares"
    
    AttributeKeyDelegator = "delegator"
    AttributeKeyValidator = "validator"
    AttributeKeyOwner = "owner"
    AttributeKeyShares = "shares"
    AttributeKeyDenom = "denom"
    AttributeKeyAmount = "amount"
    AttributeKeyRecordID = "record_id"
)
```

### 3.7 Validation

Comprehensive validation in message handler:
- Delegation exists and has sufficient shares
- Module is enabled
- Caps are not exceeded
- Addresses are valid
- Shares amount is positive
- Validator is not jailed/tombstoned

### 3.8 Testing Strategy

1. **Unit Tests**
   - Message validation
   - Denomination generation
   - State updates

2. **Integration Tests with Mocks**
   - Mock staking keeper
   - Mock bank keeper
   - Full tokenization flow

3. **Invariant Tests**
   - Total supply matches tokenization records
   - Indexes remain consistent

## Security Considerations

1. **Reentrancy Protection**: State updates before external calls
2. **Overflow Protection**: Check math operations
3. **Authorization**: Verify delegator owns the delegation
4. **Slashing**: Consider validator slashing state

## Error Handling

- Clear error messages for user-facing errors
- Proper rollback on any failure
- Comprehensive validation before state changes

## Dependencies

- Cosmos SDK v0.50.13 staking module
- Cosmos SDK bank module
- Proper module account setup

## Success Criteria

- [ ] Users can tokenize their delegations
- [ ] Liquid staking tokens are minted correctly
- [ ] All caps are enforced
- [ ] Events are emitted properly
- [ ] All tests pass
- [ ] No security vulnerabilities

## Next Steps (Stage 4 Preview)

With tokenization complete, Stage 4 will implement:
- MsgRedeemTokensForShares
- Unbonding period handling
- Redemption queue management