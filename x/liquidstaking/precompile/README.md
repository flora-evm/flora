# Liquid Staking Precompile

The liquid staking precompile provides an EVM interface to the Cosmos liquid staking module, enabling smart contracts to tokenize staked assets and manage liquid staking positions.

## Overview

The precompile is deployed at address `0x0000000000000000000000000000000000000800` and exposes the full functionality of the liquid staking module to EVM contracts.

### Key Features

- **Tokenize Delegation Shares**: Convert staked tokens into liquid staking tokens (LSTs)
- **Redeem LSTs**: Convert liquid staking tokens back to staked positions
- **Query Functions**: Access module parameters, tokenization records, and statistics
- **Full EVM Integration**: Events, structs, and standard Solidity patterns

## Architecture

```
┌─────────────────┐     ┌──────────────────┐     ┌────────────────────┐
│   EVM Contract  │────▶│   Precompile     │────▶│  Liquid Staking    │
│  (Solidity)     │     │  (0x...0800)     │     │    Module          │
└─────────────────┘     └──────────────────┘     └────────────────────┘
         │                        │                         │
         │                        │                         │
         ▼                        ▼                         ▼
    User Calls              ABI Encoding              Cosmos State
```

## Interface

### Data Structures

```solidity
struct Params {
    bool enabled;
    uint256 minLiquidStakeAmount;
    uint256 globalLiquidStakingCap; // Basis points (10000 = 100%)
    uint256 validatorLiquidCap;     // Basis points (10000 = 100%)
}

struct TokenizationRecord {
    uint256 id;
    string validatorAddress;
    address owner;
    string sharesDenomination;
    string liquidStakingTokenDenom;
    uint256 sharesAmount;
    uint8 status; // 0 = UNSPECIFIED, 1 = ACTIVE, 2 = REDEEMED
    uint256 createdAt;
    uint256 redeemedAt;
}
```

### Query Methods

#### getParams()
Returns the current module parameters.

```solidity
function getParams() external view returns (Params memory params);
```

#### getTokenizationRecord(uint256 recordId)
Retrieves a specific tokenization record by ID.

```solidity
function getTokenizationRecord(uint256 recordId) 
    external view returns (TokenizationRecord memory record);
```

#### getTotalLiquidStaked()
Returns the total amount of liquid staked tokens across all validators.

```solidity
function getTotalLiquidStaked() external view returns (uint256 amount);
```

#### getValidatorLiquidStaked(string validatorAddress)
Returns the liquid staked amount for a specific validator.

```solidity
function getValidatorLiquidStaked(string memory validatorAddress) 
    external view returns (uint256 amount);
```

### Transaction Methods

#### tokenizeShares(string validatorAddress, uint256 amount, address owner)
Tokenizes delegation shares to create liquid staking tokens.

```solidity
function tokenizeShares(
    string memory validatorAddress, 
    uint256 amount, 
    address owner
) external returns (TokenizeSharesResponse memory response);
```

**Parameters:**
- `validatorAddress`: The Bech32 address of the validator
- `amount`: Amount of shares to tokenize (in share units, not tokens)
- `owner`: Address to receive the LSTs (use `address(0)` for msg.sender)

**Returns:**
- `recordId`: The unique ID of the tokenization record
- `tokensDenom`: The denomination of the minted LST
- `tokensAmount`: The amount of LST minted

#### redeemTokens(string tokenDenom, uint256 amount)
Redeems liquid staking tokens back to staked positions.

```solidity
function redeemTokens(string memory tokenDenom, uint256 amount) 
    external returns (RedeemTokensResponse memory response);
```

**Parameters:**
- `tokenDenom`: The LST denomination (e.g., "liquidstake/floravaloper1.../1")
- `amount`: Amount of LST to redeem

**Returns:**
- `validatorAddress`: The validator address
- `sharesAmount`: Amount of shares restored
- `completed`: Whether the record was fully redeemed

## Gas Costs

| Method | Gas Cost |
|--------|----------|
| Query Methods | 30,000 |
| getTokenizationRecord | 50,000 |
| tokenizeShares | 100,000 |
| redeemTokens | 80,000 |

## Events

### TokenizeSharesEvent
Emitted when shares are tokenized.

```solidity
event TokenizeSharesEvent(
    address indexed delegator,
    string indexed validator,
    uint256 indexed recordId,
    address owner,
    uint256 sharesAmount,
    uint256 tokensAmount,
    string tokenDenom
);
```

### RedeemTokensEvent
Emitted when tokens are redeemed.

```solidity
event RedeemTokensEvent(
    address indexed owner,
    string indexed validator,
    uint256 indexed recordId,
    string tokenDenom,
    uint256 tokensAmount,
    uint256 sharesAmount,
    bool completed
);
```

## Usage Examples

### Basic Usage

```solidity
import "./ILiquidStaking.sol";

contract MyContract {
    ILiquidStaking constant LIQUID_STAKING = 
        ILiquidStaking(0x0000000000000000000000000000000000000800);
    
    function stakeAndTokenize(string memory validator, uint256 amount) external {
        // First ensure you have delegated to the validator
        // Then tokenize the shares
        ILiquidStaking.TokenizeSharesResponse memory response = 
            LIQUID_STAKING.tokenizeShares(validator, amount, msg.sender);
        
        // response.recordId - unique record identifier
        // response.tokensDenom - LST denomination
        // response.tokensAmount - amount of LST received
    }
}
```

### Using the Helper Library

```solidity
import "./ILiquidStaking.sol";

contract MyContract {
    using LiquidStaking for *;
    
    function checkModuleStatus() external view returns (bool) {
        ILiquidStaking.Params memory params = LiquidStaking.CONTRACT.getParams();
        return params.enabled;
    }
}
```

## Implementation Details

### ABI Encoding

The precompile uses standard Solidity ABI encoding for all inputs and outputs. Complex types like structs are encoded as tuples.

### Method IDs

Method IDs are computed using the standard Keccak-256 hash of the method signature:

```
getParams() => 0x5e615a6b
tokenizeShares(string,uint256,address) => 0x...
```

### Context Handling

The precompile receives the SDK context from the EVM state database, ensuring proper transaction isolation and state management.

### Error Handling

Errors from the underlying Cosmos module are propagated as EVM reverts with descriptive error messages.

## Security Considerations

1. **Validator Validation**: Always verify validator addresses before tokenizing
2. **Amount Validation**: Ensure amounts meet minimum requirements
3. **Owner Validation**: Be careful with the owner parameter in tokenizeShares
4. **Reentrancy**: The precompile is reentrancy-safe at the protocol level
5. **Access Control**: Implement proper access control in your contracts

## Testing

### Unit Tests

```go
go test -v ./x/liquidstaking/precompile/...
```

### Integration Tests

Deploy test contracts to a local Flora node:

```bash
# Start local node with liquid staking enabled
florad start --enable-liquid-staking

# Deploy and test contracts
npx hardhat test --network flora-local
```

## Troubleshooting

### Common Issues

1. **"Liquid staking disabled"**: Ensure the module is enabled in params
2. **"Amount too small"**: Check minLiquidStakeAmount parameter
3. **"No delegation found"**: Ensure delegation exists before tokenizing
4. **"Invalid validator"**: Verify the validator address format

### Debugging

Enable EVM debug logs:

```bash
florad start --log_level="*:info,evm:debug"
```

## Related Documentation

- [Liquid Staking Module](../README.md)
- [Example Contracts](./examples/)
- [ABI Specification](./abi.json)
- [Solidity Interface](./ILiquidStaking.sol)