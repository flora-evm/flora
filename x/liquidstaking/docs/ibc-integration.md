# Liquid Staking IBC Integration

## Overview

The liquid staking module provides full IBC (Inter-Blockchain Communication) support, allowing liquid staking tokens (LSTs) to be transferred between different Cosmos chains while preserving their metadata and properties.

## Architecture

### Components

1. **IBC Packet Structures** (`types/ibc.go`)
   - Extended packet data with liquid staking metadata
   - Backward compatible with standard IBC transfers
   - Metadata preserved in packet memo field

2. **IBC Transfer Hooks** (`keeper/ibc_hooks.go`)
   - Intercepts outgoing LST transfers
   - Processes incoming LST transfers
   - Handles acknowledgements and timeouts

3. **IBC Transfer Handler** (`keeper/ibc_transfer_handler.go`)
   - Main logic for cross-chain LST transfers
   - Creates local representations of foreign LSTs
   - Manages metadata storage and retrieval

## How It Works

### Sending LSTs via IBC

When a liquid staking token is sent to another chain:

1. **Validation**: The module validates the transfer request
   - Checks if the module is enabled
   - Verifies the tokenization record exists and is active
   - Ensures the sender has sufficient balance

2. **Metadata Extraction**: The module extracts metadata from the tokenization record
   - Validator address
   - Record ID
   - Original shares amount
   - Creation timestamp
   - Source chain ID

3. **Packet Creation**: The metadata is embedded in the IBC packet memo
   ```json
   {
     "liquid_staking_metadata": {
       "validator_address": "floravaloper1abc...",
       "record_id": 1,
       "shares_amount": "1000000.000000000000000000",
       "source_chain_id": "flora-1",
       "created_at": "2024-01-15T10:30:00Z"
     }
   }
   ```

4. **Transfer Execution**: The standard IBC transfer module handles the actual transfer

### Receiving LSTs via IBC

When receiving a liquid staking token from another chain:

1. **Metadata Detection**: The module checks the packet memo for liquid staking metadata

2. **Local Representation**: If metadata is found, the module creates a local representation
   - Generates appropriate denom metadata
   - Stores the original liquid staking information
   - Creates human-readable descriptions

3. **Token Minting**: The IBC tokens are minted to the receiver's account

4. **Metadata Storage**: The original LST metadata is stored for future reference

## Token Denominations

### Native LST Format
```
liquidstake/{validator}/{record-id}
```
Example: `liquidstake/floravaloper1abc.../1`

### IBC LST Format
After crossing chains, the denom follows standard IBC format:
```
ibc/{hash}
```
Example: `ibc/ABCDEF123456...`

The module maintains metadata that links the IBC denom to the original LST information.

## Events

The module emits specific events for IBC operations:

### Outgoing Transfer
```
Event: liquid_staking_ibc_transfer
Attributes:
  - sender: flora1xyz...
  - receiver: cosmos1abc...
  - denom: liquidstake/floravaloper1abc.../1
  - amount: 1000000
  - source_port: transfer
  - source_channel: channel-0
  - record_id: 1
```

### Incoming Transfer
```
Event: liquid_staking_ibc_received
Attributes:
  - receiver: flora1xyz...
  - denom: ibc/ABCDEF123456...
  - amount: 1000000
  - source_chain_id: cosmos-hub-4
  - validator_address: cosmosvaloper1abc...
  - record_id: 1
```

## CLI Usage

### Send LST via IBC
```bash
# Send liquid staking tokens to another chain
florad tx ibc-transfer transfer transfer channel-0 \
  cosmos1receiver... 1000000liquidstake/floravaloper1abc.../1 \
  --from mykey \
  --packet-timeout-height 0-1000 \
  --packet-timeout-timestamp 0
```

### Query IBC LST Metadata
```bash
# Check if an IBC denom is a liquid staking token
florad query liquidstaking ibc-metadata ibc/ABCDEF123456...
```

## Integration Points

### IBC Transfer Module
The liquid staking module integrates with the IBC transfer module through:
- Transfer hooks for intercepting packets
- Shared escrow accounts for token locking
- Standard acknowledgement handling

### Bank Module
For IBC liquid staking tokens:
- Denom metadata is stored in the bank module
- Standard bank queries work with IBC LSTs
- Transfer restrictions apply normally

## Security Considerations

1. **Metadata Validation**: All incoming metadata is validated before storage
2. **Record Status**: Only active tokenization records can be transferred
3. **Channel Security**: Standard IBC channel security applies
4. **Escrow Safety**: Tokens are properly escrowed during transfers

## Best Practices

### For Chains Receiving LSTs

1. **Metadata Handling**: Preserve liquid staking metadata when forwarding tokens
2. **Display Information**: Use the metadata to show meaningful token information
3. **Redemption Path**: Consider the redemption path back to the source chain

### For Users

1. **Track Source Chain**: Remember which chain originally issued the LST
2. **Redemption**: To redeem, transfer back to the source chain first
3. **Multi-hop Awareness**: Each IBC hop changes the denom

## Technical Details

### Packet Data Structure
```go
type LiquidStakingTokenPacketData struct {
    // Standard IBC transfer fields
    Denom    string
    Amount   string
    Sender   string
    Receiver string
    Memo     string
    
    // Liquid staking metadata
    LiquidStakingMetadata *LiquidStakingMetadata
}
```

### Metadata Structure
```go
type LiquidStakingMetadata struct {
    ValidatorAddress string
    RecordId         uint64
    SharesAmount     math.LegacyDec
    SourceChainId    string
    CreatedAt        string
}
```

### Hooks Interface
```go
type IBCHooks interface {
    OnSendPacket(ctx, sourcePort, sourceChannel, token, sender, receiver, memo, relayer) error
    OnRecvPacket(ctx, packet, ack, relayer) error
    OnAcknowledgementPacket(ctx, packet, acknowledgement, relayer) error
    OnTimeoutPacket(ctx, packet, relayer) error
}
```

## Future Enhancements

1. **Cross-chain Redemption**: Allow redemption on non-source chains
2. **Metadata Compression**: Optimize metadata size in memo field
3. **Whitelisted Channels**: Restrict LST transfers to approved channels
4. **Fee Distribution**: Handle reward distribution for IBC LSTs

## Troubleshooting

### Common Issues

1. **"Liquid staking module is disabled"**
   - Solution: Enable the module through governance

2. **"Cannot transfer redeemed liquid staking tokens"**
   - Solution: Only active LSTs can be transferred

3. **"Channel not found"**
   - Solution: Ensure IBC channel is established and open

4. **Metadata Not Preserved**
   - Check that both chains support liquid staking IBC
   - Verify memo field size limits

### Debugging

Enable debug logging:
```bash
florad start --log_level debug
```

Check IBC packet status:
```bash
florad query ibc-transfer packet-commitments transfer channel-0
```

## References

- [IBC Protocol Specification](https://github.com/cosmos/ibc)
- [IBC Transfer Module](https://ibc.cosmos.network/main/apps/transfer/overview.html)
- [Liquid Staking Module Documentation](./README.md)