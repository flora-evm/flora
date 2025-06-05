// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/security/ReentrancyGuard.sol";

/**
 * @title LiquidStakedToken (LST)
 * @notice Auto-compounding liquid staking token for Flora blockchain
 * @dev This token represents staked PETAL with a dynamic exchange rate
 */
contract LiquidStakedToken is ERC20, Ownable, ReentrancyGuard {
    // Constants
    uint256 private constant PRECISION = 1e18;
    address public constant LIQUID_STAKING_PRECOMPILE = 0x0000000000000000000000000000000000000800;
    
    // State variables
    address public immutable validator;
    uint256 public totalShares;        // Total shares minted (internal accounting)
    uint256 public totalPooledPetal;   // Total PETAL value including rewards
    uint256 public lastRewardUpdate;   // Last time rewards were updated
    
    // Exchange rate = totalPooledPetal / totalShares
    // 1 LST token represents a share of the total pooled PETAL
    
    // Events
    event SharesMinted(address indexed account, uint256 shares, uint256 petalAmount);
    event SharesBurned(address indexed account, uint256 shares, uint256 petalAmount);
    event RewardsCompounded(uint256 rewardAmount, uint256 newExchangeRate);
    event SlashingApplied(uint256 slashAmount, uint256 newExchangeRate);
    
    modifier onlyPrecompile() {
        require(msg.sender == LIQUID_STAKING_PRECOMPILE, "Only precompile");
        _;
    }
    
    constructor(
        string memory name,
        string memory symbol,
        address _validator
    ) ERC20(name, symbol) {
        validator = _validator;
        totalShares = 0;
        totalPooledPetal = 0;
        lastRewardUpdate = block.timestamp;
    }
    
    /**
     * @notice Returns the current exchange rate (PETAL per share)
     */
    function getExchangeRate() public view returns (uint256) {
        if (totalShares == 0) {
            return PRECISION; // 1:1 initial rate
        }
        return (totalPooledPetal * PRECISION) / totalShares;
    }
    
    /**
     * @notice Converts PETAL amount to shares
     */
    function getPetalByShares(uint256 sharesAmount) public view returns (uint256) {
        return (sharesAmount * getExchangeRate()) / PRECISION;
    }
    
    /**
     * @notice Converts shares to PETAL amount
     */
    function getSharesByPetal(uint256 petalAmount) public view returns (uint256) {
        return (petalAmount * PRECISION) / getExchangeRate();
    }
    
    /**
     * @notice Mints new LST tokens when PETAL is staked
     * @dev Only callable by the liquid staking precompile
     */
    function mint(address account, uint256 petalAmount) external onlyPrecompile nonReentrant {
        require(account != address(0), "Mint to zero address");
        require(petalAmount > 0, "Amount must be positive");
        
        uint256 sharesToMint = getSharesByPetal(petalAmount);
        totalShares += sharesToMint;
        totalPooledPetal += petalAmount;
        
        _mint(account, sharesToMint);
        
        emit SharesMinted(account, sharesToMint, petalAmount);
    }
    
    /**
     * @notice Burns LST tokens when redeeming for staked PETAL
     * @dev Only callable by the liquid staking precompile
     */
    function burn(address account, uint256 shares) external onlyPrecompile nonReentrant {
        require(account != address(0), "Burn from zero address");
        require(shares > 0, "Amount must be positive");
        require(balanceOf(account) >= shares, "Insufficient balance");
        
        uint256 petalAmount = getPetalByShares(shares);
        totalShares -= shares;
        totalPooledPetal -= petalAmount;
        
        _burn(account, shares);
        
        emit SharesBurned(account, shares, petalAmount);
    }
    
    /**
     * @notice Updates the total pooled PETAL to reflect staking rewards
     * @dev Called periodically by oracle or keeper
     */
    function updateRewards(uint256 newTotalPooledPetal) external onlyOwner {
        require(newTotalPooledPetal >= totalPooledPetal, "Cannot decrease via rewards update");
        
        uint256 rewardAmount = newTotalPooledPetal - totalPooledPetal;
        totalPooledPetal = newTotalPooledPetal;
        lastRewardUpdate = block.timestamp;
        
        emit RewardsCompounded(rewardAmount, getExchangeRate());
    }
    
    /**
     * @notice Applies slashing penalty to the total pooled PETAL
     * @dev Called when validator is slashed
     */
    function applySlashing(uint256 slashAmount) external onlyOwner {
        require(slashAmount <= totalPooledPetal, "Slash amount exceeds total");
        
        totalPooledPetal -= slashAmount;
        
        emit SlashingApplied(slashAmount, getExchangeRate());
    }
    
    /**
     * @notice Override transfer to use shares internally
     */
    function transfer(address to, uint256 amount) public override returns (bool) {
        return super.transfer(to, amount);
    }
    
    /**
     * @notice Override transferFrom to use shares internally
     */
    function transferFrom(address from, address to, uint256 amount) public override returns (bool) {
        return super.transferFrom(from, to, amount);
    }
    
    /**
     * @notice Returns the PETAL value of an account's LST balance
     */
    function balanceOfPetal(address account) external view returns (uint256) {
        return getPetalByShares(balanceOf(account));
    }
    
    /**
     * @notice Returns total PETAL value of all LST tokens
     */
    function totalSupplyPetal() external view returns (uint256) {
        return totalPooledPetal;
    }
}

/**
 * @title ILiquidStaking
 * @notice Interface for the liquid staking precompile
 */
interface ILiquidStaking {
    function tokenizeShares(address validator, uint256 amount) external returns (uint256 recordId, address lstToken);
    function redeemTokens(uint256 recordId, uint256 amount) external returns (bool success);
    function transferRecord(uint256 recordId, address newOwner) external returns (bool success);
    function getRecord(uint256 recordId) external view returns (address owner, address validator, uint256 shares, address lstToken);
}