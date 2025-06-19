# Event System Documentation

## Overview

The liquid staking module implements a comprehensive event system to track all significant state changes and operations. These events are crucial for monitoring, indexing, and auditing liquid staking activities.

## Event Categories

### 1. Core Operation Events

#### TokenizeShares Event (`tokenize_shares`)
Emitted when staking shares are converted to liquid staking tokens.

**Attributes:**
- `delegator`: Address of the original delegator
- `validator`: Validator address
- `owner`: Owner of the liquid staking tokens (may differ from delegator)
- `shares`: Amount of shares tokenized
- `tokens_minted`: Amount of LST tokens minted
- `denom`: The LST token denomination
- `record_id`: Unique tokenization record identifier

**Example:**
```json
{
  "type": "tokenize_shares",
  "attributes": [
    {"key": "delegator", "value": "flora1abc..."},
    {"key": "validator", "value": "floravaloper1xyz..."},
    {"key": "owner", "value": "flora1def..."},
    {"key": "shares", "value": "1000000"},
    {"key": "tokens_minted", "value": "1000000"},
    {"key": "denom", "value": "liquidstake/floravaloper1xyz.../1"},
    {"key": "record_id", "value": "1"}
  ]
}
```

#### RedeemTokens Event (`redeem_tokens`)
Emitted when liquid staking tokens are burned to restore staking shares.

**Attributes:**
- `owner`: Address redeeming the tokens
- `validator`: Validator address
- `tokens_burned`: Amount of LST tokens burned
- `shares_restored`: Amount of shares restored
- `denom`: The LST token denomination
- `record_id`: Associated tokenization record identifier

### 2. Record Lifecycle Events

#### Record Created (`tokenization_record_created`)
Emitted when a new tokenization record is created.

**Attributes:**
- `record_id`: Unique record identifier
- `validator`: Associated validator
- `owner`: Record owner
- `shares_tokenized`: Initial tokenized amount
- `denom`: LST token denomination
- `action`: Always "create"

#### Record Updated (`tokenization_record_updated`)
Emitted when a tokenization record is modified (partial redemption).

**Attributes:**
- `record_id`: Record identifier
- `old_shares_tokenized`: Previous amount
- `new_shares_tokenized`: Updated amount
- `action`: Always "update"

#### Record Deleted (`tokenization_record_deleted`)
Emitted when a tokenization record is removed (full redemption).

**Attributes:**
- `record_id`: Record identifier
- `validator`: Associated validator
- `owner`: Record owner
- `denom`: LST token denomination
- `action`: Always "delete"

### 3. Governance Events

#### Update Parameters (`update_params`)
Emitted when module parameters are modified.

**Attributes:**
- `sender`: Authority making the change (usually "governance")
- `action`: Always "update"
- `change_X_param_key`: Parameter name for change X
- `change_X_param_old_value`: Previous value for change X
- `change_X_param_new_value`: New value for change X

**Example with multiple changes:**
```json
{
  "type": "update_params",
  "attributes": [
    {"key": "sender", "value": "governance"},
    {"key": "action", "value": "update"},
    {"key": "change_0_param_key", "value": "enabled"},
    {"key": "change_0_param_old_value", "value": "true"},
    {"key": "change_0_param_new_value", "value": "false"},
    {"key": "change_1_param_key", "value": "min_liquid_stake_amount"},
    {"key": "change_1_param_old_value", "value": "1000"},
    {"key": "change_1_param_new_value", "value": "2000"}
  ]
}
```

### 4. Cap Management Events

#### Liquid Staking Cap (`liquid_staking_cap`)
Emitted when liquid staking approaches or exceeds defined caps.

**Attributes:**
- `cap_type`: Either "global" or "validator"
- `validator`: Validator address (only for validator caps)
- `current_amount`: Current liquid staked amount
- `cap_limit`: Maximum allowed amount
- `percentage_used`: Percentage of cap utilized

## Event Usage Patterns

### 1. Monitoring and Alerting

Events can be used to set up monitoring systems:

```go
// Example: Monitor when caps are approaching limits
if event.Type == "liquid_staking_cap" {
    percentageUsed := getAttributeValue(event, "percentage_used")
    if percentageUsed >= "90" {
        sendAlert("Liquid staking cap approaching limit")
    }
}
```

### 2. Indexing for Queries

Events enable efficient querying of historical data:

```sql
-- Example: Find all tokenizations for a specific validator
SELECT * FROM events 
WHERE type = 'tokenize_shares' 
AND attributes->>'validator' = 'floravaloper1xyz...';

-- Example: Track record lifecycle
SELECT * FROM events 
WHERE attributes->>'record_id' = '42'
ORDER BY block_height;
```

### 3. Audit Trail

Events provide a complete audit trail:

```go
// Track all operations by a specific owner
ownerEvents := []sdk.Event{}
for _, event := range allEvents {
    if hasAttribute(event, "owner", ownerAddress) {
        ownerEvents = append(ownerEvents, event)
    }
}
```

## Implementation Details

### Typed Event System

The module uses typed events for type safety:

```go
type TokenizeSharesEvent struct {
    Delegator      string
    Validator      string
    Owner          string
    SharesAmount   string
    TokensMinted   string
    Denom          string
    RecordID       uint64
}

// Convert to SDK event
func (e TokenizeSharesEvent) ToEvent() sdk.Event {
    return sdk.NewEvent(
        EventTypeTokenizeShares,
        sdk.NewAttribute(AttributeKeyDelegator, e.Delegator),
        // ... other attributes
    )
}
```

### Event Emission

Events are emitted at the appropriate points in message handlers:

```go
// In TokenizeShares handler
types.EmitTokenizeSharesEvent(ctx, types.TokenizeSharesEvent{
    Delegator:    msg.DelegatorAddress,
    Validator:    msg.ValidatorAddress,
    Owner:        ownerAddr.String(),
    SharesAmount: sharesToTokenize.String(),
    TokensMinted: unbondedTokens.String(),
    Denom:        denom,
    RecordID:     recordID,
})
```

### Standard Message Events

In addition to custom events, the module emits standard Cosmos SDK message events:

```go
ctx.EventManager().EmitEvent(
    sdk.NewEvent(
        sdk.EventTypeMessage,
        sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
        sdk.NewAttribute(sdk.AttributeKeySender, sender),
        sdk.NewAttribute(sdk.AttributeKeyAction, action),
    ),
)
```

## Best Practices

1. **Always emit events for state changes**: Any operation that modifies state should emit corresponding events.

2. **Use typed events**: Leverage the typed event system for compile-time safety.

3. **Include all relevant context**: Events should contain enough information to understand the operation without additional queries.

4. **Maintain consistency**: Use consistent attribute names across similar events.

5. **Handle errors properly**: Events should only be emitted when operations succeed.

## Testing Events

The module includes comprehensive event tests:

```go
func TestTokenizeSharesEvents(t *testing.T) {
    // Execute operation
    res, err := msgServer.TokenizeShares(ctx, msg)
    require.NoError(t, err)
    
    // Find specific event
    tokenizeEvent := findEvent(events, types.EventTypeTokenizeShares)
    require.NotNil(t, tokenizeEvent)
    
    // Verify attributes
    assertEventAttribute(t, tokenizeEvent, 
        types.AttributeKeyDelegator, delegator.String())
}
```

## Event Constants Reference

All event types and attribute keys are defined in `types/events.go`:

```go
// Event types
const (
    EventTypeTokenizeShares     = "tokenize_shares"
    EventTypeRedeemTokens       = "redeem_tokens"
    EventTypeUpdateParams       = "update_params"
    EventTypeRecordCreated      = "tokenization_record_created"
    EventTypeRecordUpdated      = "tokenization_record_updated"
    EventTypeRecordDeleted      = "tokenization_record_deleted"
    EventTypeLiquidStakingCap   = "liquid_staking_cap"
)
```

## Integration Examples

### Client-side Event Subscription

```javascript
// Subscribe to tokenization events
const subscription = await client.subscribe({
    query: "tm.event='Tx' AND liquidstaking.action='tokenize'"
});

subscription.on('data', (event) => {
    console.log('New tokenization:', event);
});
```

### Event Processing Service

```go
func ProcessLiquidStakingEvents(events []sdk.Event) {
    for _, event := range events {
        switch event.Type {
        case types.EventTypeTokenizeShares:
            handleTokenization(event)
        case types.EventTypeRedeemTokens:
            handleRedemption(event)
        case types.EventTypeRecordDeleted:
            cleanupRecord(event)
        }
    }
}
```

## Future Enhancements

1. **Event Batching**: For operations affecting multiple records, consider batch events.

2. **Event Versioning**: As the module evolves, implement event versioning for backward compatibility.

3. **Custom Indexing**: Develop specialized indexes for common query patterns.

4. **Real-time Notifications**: Integrate with notification systems for critical events.

5. **Analytics Integration**: Export events to analytics platforms for deeper insights.