// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../ILiquidStaking.sol";

/**
 * @title BasicLiquidStaking
 * @dev Basic example of interacting with the liquid staking precompile
 * @notice This contract demonstrates basic liquid staking operations
 */
contract BasicLiquidStaking {
    using LiquidStaking for *;
    
    // Events
    event SharesTokenized(
        address indexed user,
        string validator,
        uint256 recordId,
        uint256 amount,
        string tokenDenom
    );
    
    event TokensRedeemed(
        address indexed user,
        string tokenDenom,
        uint256 amount,
        bool completed
    );
    
    /**
     * @notice Tokenize shares for a specific validator
     * @param validator The validator address to tokenize shares from
     * @param amount The amount of shares to tokenize
     */
    function tokenizeMyShares(string memory validator, uint256 amount) external {
        // Call the precompile to tokenize shares
        ILiquidStaking.TokenizeSharesResponse memory response = 
            LiquidStaking.CONTRACT.tokenizeShares(validator, amount, msg.sender);
        
        // Emit event for tracking
        emit SharesTokenized(
            msg.sender,
            validator,
            response.recordId,
            response.tokensAmount,
            response.tokensDenom
        );
    }
    
    /**
     * @notice Redeem liquid staking tokens
     * @param tokenDenom The liquid staking token denomination
     * @param amount The amount of tokens to redeem
     */
    function redeemMyTokens(string memory tokenDenom, uint256 amount) external {
        // Call the precompile to redeem tokens
        ILiquidStaking.RedeemTokensResponse memory response = 
            LiquidStaking.CONTRACT.redeemTokens(tokenDenom, amount);
        
        // Emit event for tracking
        emit TokensRedeemed(
            msg.sender,
            tokenDenom,
            amount,
            response.completed
        );
    }
    
    /**
     * @notice Get the current module parameters
     * @return enabled Whether liquid staking is enabled
     * @return minAmount The minimum liquid stake amount
     * @return globalCap The global liquid staking cap (basis points)
     * @return validatorCap The validator liquid cap (basis points)
     */
    function getModuleParams() external view returns (
        bool enabled,
        uint256 minAmount,
        uint256 globalCap,
        uint256 validatorCap
    ) {
        ILiquidStaking.Params memory params = LiquidStaking.CONTRACT.getParams();
        return (
            params.enabled,
            params.minLiquidStakeAmount,
            params.globalLiquidStakingCap,
            params.validatorLiquidCap
        );
    }
    
    /**
     * @notice Get information about a tokenization record
     * @param recordId The record ID to query
     * @return record The tokenization record details
     */
    function getRecord(uint256 recordId) 
        external 
        view 
        returns (ILiquidStaking.TokenizationRecord memory record) 
    {
        return LiquidStaking.CONTRACT.getTokenizationRecord(recordId);
    }
    
    /**
     * @notice Check if a token is a liquid staking token
     * @param tokenDenom The token denomination to check
     * @return isLST Whether it's a liquid staking token
     * @return validator The associated validator
     * @return recordId The associated record ID
     */
    function checkLiquidStakingToken(string memory tokenDenom) 
        external 
        view 
        returns (
            bool isLST,
            string memory validator,
            uint256 recordId
        ) 
    {
        ILiquidStaking.LiquidStakingTokenInfo memory info = 
            LiquidStaking.CONTRACT.getLiquidStakingTokenInfo(tokenDenom);
            
        return (
            info.isLiquidStakingToken,
            info.validatorAddress,
            info.recordId
        );
    }
}