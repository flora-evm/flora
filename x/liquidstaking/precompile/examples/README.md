# Liquid Staking Precompile Examples

This directory contains example Solidity contracts demonstrating how to interact with the Flora liquid staking precompile from EVM smart contracts.

## Examples

### 1. BasicLiquidStaking.sol

A simple contract showing the basic operations:
- Tokenizing delegation shares
- Redeeming liquid staking tokens
- Querying module parameters
- Checking tokenization records
- Verifying if a token is an LST

**Use Case**: Building simple liquid staking interfaces or integrating basic LST functionality into existing contracts.

### 2. LiquidStakingVault.sol

An advanced vault contract that:
- Manages user deposits and liquid staking positions
- Tracks multiple deposits per user
- Provides emergency withdrawal functionality
- Implements access control and pausability
- Offers paginated queries for user positions

**Use Case**: Building yield aggregators, staking pools, or managed liquid staking services.

### 3. LiquidStakingDEX.sol

A DeFi integration example showing:
- Adding liquidity with newly minted LSTs
- Creating LST/WFLORA trading pairs
- Swapping between LSTs and native tokens
- Price discovery for liquid staking tokens
- Managing liquidity positions

**Use Case**: Integrating liquid staking tokens with DEXs, AMMs, or other DeFi protocols.

## Usage

### Importing the Interface

All examples import the precompile interface:

```solidity
import "../ILiquidStaking.sol";
```

### Accessing the Precompile

The precompile is available at a fixed address:

```solidity
ILiquidStaking constant LIQUID_STAKING = ILiquidStaking(0x0000000000000000000000000000000000000800);
```

Or use the convenience library:

```solidity
using LiquidStaking for *;
// Then access as: LiquidStaking.CONTRACT.methodName()
```

### Key Concepts

1. **Tokenization Records**: Each liquid staking position creates a record with a unique ID
2. **LST Denominations**: Liquid staking tokens have specific denominations like `liquidstake/floravaloper1abc/1`
3. **Share Amounts**: Amounts are in validator shares, not native tokens
4. **Status Tracking**: Records can be ACTIVE or REDEEMED

### Gas Considerations

The precompile methods have different gas costs:
- Query methods: ~30,000 gas (base query)
- `getTokenizationRecord`: ~50,000 gas
- `tokenizeShares`: ~100,000 gas
- `redeemTokens`: ~80,000 gas

### Best Practices

1. **Check Module Status**: Always verify liquid staking is enabled before operations
2. **Validate Amounts**: Ensure amounts meet minimum requirements
3. **Handle Failures**: Implement proper error handling for precompile calls
4. **Event Emission**: Emit events for important operations for off-chain tracking
5. **Access Control**: Implement proper access control for admin functions

### Testing

To test these contracts:

1. Deploy to a local Flora testnet with the liquid staking module enabled
2. Ensure the precompile is registered at the correct address
3. Have test accounts with delegated stakes to tokenize
4. Use a testing framework like Hardhat or Foundry

Example test setup:

```javascript
// Hardhat test example
const BasicLiquidStaking = await ethers.getContractFactory("BasicLiquidStaking");
const contract = await BasicLiquidStaking.deploy();

// Tokenize shares (ensure account has delegation first)
const tx = await contract.tokenizeMyShares("floravaloper1...", ethers.utils.parseEther("100"));
await tx.wait();
```

### Security Considerations

1. **Reentrancy**: The precompile handles state changes atomically
2. **Access Control**: Implement proper ownership and permission systems
3. **Validation**: Always validate inputs and check return values
4. **Upgradability**: Consider upgrade patterns for production contracts
5. **Testing**: Thoroughly test all interactions, especially edge cases

## Additional Resources

- [Liquid Staking Module Documentation](../../README.md)
- [Precompile Interface](../ILiquidStaking.sol)
- [Flora EVM Documentation](https://docs.flora.network/evm)