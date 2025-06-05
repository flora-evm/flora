# Liquid Staking Token Naming Convention

## Overview
The Flora blockchain uses a consistent naming convention for liquid staking tokens to ensure clarity and compatibility across the ecosystem.

## Token Naming

### Native Token
- **Symbol**: `FLORA`
- **Denom**: `flora`
- **Decimals**: 18
- **Description**: The native staking and gas token of the Flora blockchain

### Liquid Staking Tokens (LST)
- **Symbol Format**: `stFLORA-{ValidatorID}`
- **Denom Format**: `factory/{creator}/stflora_{validator_suffix}`
- **Decimals**: 18 (matching FLORA)
- **Description**: Liquid staked FLORA tokens representing staked positions with specific validators

## Examples

### Per-Validator LST Tokens
Each validator has their own liquid staking token:

```
Validator: floravaloper1abc...xyz
LST Symbol: stFLORA-1abc
LST Denom: factory/flora1.../stflora_1abc
```

### Token Factory Integration
LST tokens are created using the Token Factory module:

```solidity
// ERC20 representation
contract stFLORA_ValidatorX {
    string public name = "Liquid Staked FLORA - Validator X";
    string public symbol = "stFLORA-X";
    uint8 public decimals = 18;
}
```

## Benefits of This Naming Convention

1. **Clarity**: Users immediately understand that stFLORA represents staked FLORA
2. **Validator Differentiation**: Each validator's LST is uniquely identifiable
3. **Ecosystem Consistency**: Follows established patterns (stETH, stATOM, etc.)
4. **IBC Compatibility**: Denom format works seamlessly with IBC transfers
5. **DeFi Integration**: Clear naming helps with AMM pairs (stFLORA/FLORA)

## Usage in DeFi

### AMM Pools
- Primary pairs: `stFLORA-{ValidatorID}/FLORA`
- Cross-validator pairs: `stFLORA-A/stFLORA-B`

### Lending Markets
- Collateral: stFLORA tokens can be used as collateral
- Naming: Markets clearly show "stFLORA-X" as distinct assets

### Governance
- Voting power: Based on underlying FLORA amount
- Proposal format: "Enable stFLORA-{ValidatorID} as collateral"

## Technical Implementation

### Precompile Address Mapping
```go
// Liquid staking precompile creates tokens
subdenom := fmt.Sprintf("stflora_%s", validator.String()[:8])
denom := fmt.Sprintf("factory/%s/%s", creator, subdenom)
```

### ERC20 Metadata
```solidity
function name() public view returns (string memory) {
    return string(abi.encodePacked("Liquid Staked FLORA - ", validatorMoniker));
}

function symbol() public view returns (string memory) {
    return string(abi.encodePacked("stFLORA-", validatorID));
}
```

## Migration from Previous Naming
- Previous: `petal` → `flora`
- Previous: `sPETAL` → `stFLORA`
- Maintains same structure, just updated base token name