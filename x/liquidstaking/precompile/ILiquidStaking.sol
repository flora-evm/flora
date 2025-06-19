// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

/**
 * @title ILiquidStaking
 * @dev Interface for the liquid staking precompile contract
 * @notice This precompile allows EVM contracts to interact with the Cosmos liquid staking module
 */
interface ILiquidStaking {
    // Structs
    
    /**
     * @dev Module parameters
     */
    struct Params {
        bool enabled;
        uint256 minLiquidStakeAmount;
        uint256 globalLiquidStakingCap; // Basis points (10000 = 100%)
        uint256 validatorLiquidCap;     // Basis points (10000 = 100%)
    }
    
    /**
     * @dev Tokenization record information
     */
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
    
    /**
     * @dev Response for tokenizeShares method
     */
    struct TokenizeSharesResponse {
        uint256 recordId;
        string tokensDenom;
        uint256 tokensAmount;
    }
    
    /**
     * @dev Response for redeemTokens method
     */
    struct RedeemTokensResponse {
        string validatorAddress;
        uint256 sharesAmount;
        bool completed;
    }
    
    /**
     * @dev Information about a liquid staking token
     */
    struct LiquidStakingTokenInfo {
        bool isLiquidStakingToken;
        string validatorAddress;
        uint256 recordId;
        address originalDelegator;
        bool active;
    }
    
    // Events
    
    /**
     * @dev Emitted when shares are tokenized
     * @param delegator The address that delegated the shares
     * @param validator The validator address
     * @param recordId The tokenization record ID
     * @param owner The owner of the liquid staking tokens
     * @param sharesAmount The amount of shares tokenized
     * @param tokensAmount The amount of tokens minted
     * @param tokenDenom The denomination of the minted tokens
     */
    event TokenizeSharesEvent(
        address indexed delegator,
        string indexed validator,
        uint256 indexed recordId,
        address owner,
        uint256 sharesAmount,
        uint256 tokensAmount,
        string tokenDenom
    );
    
    /**
     * @dev Emitted when tokens are redeemed
     * @param owner The owner redeeming the tokens
     * @param validator The validator address
     * @param recordId The tokenization record ID
     * @param tokenDenom The denomination of the redeemed tokens
     * @param tokensAmount The amount of tokens redeemed
     * @param sharesAmount The amount of shares restored
     * @param completed Whether the redemption completed the record
     */
    event RedeemTokensEvent(
        address indexed owner,
        string indexed validator,
        uint256 indexed recordId,
        string tokenDenom,
        uint256 tokensAmount,
        uint256 sharesAmount,
        bool completed
    );
    
    // Query Functions
    
    /**
     * @notice Get the module parameters
     * @return params The current module parameters
     */
    function getParams() external view returns (Params memory params);
    
    /**
     * @notice Get a specific tokenization record
     * @param recordId The ID of the record to retrieve
     * @return record The tokenization record
     */
    function getTokenizationRecord(uint256 recordId) external view returns (TokenizationRecord memory record);
    
    /**
     * @notice Get all tokenization records with pagination
     * @param offset The starting index
     * @param limit The maximum number of records to return
     * @return records Array of tokenization records
     * @return total The total number of records
     */
    function getTokenizationRecords(uint256 offset, uint256 limit) 
        external view returns (TokenizationRecord[] memory records, uint256 total);
    
    /**
     * @notice Get tokenization records by owner
     * @param owner The owner address
     * @param offset The starting index
     * @param limit The maximum number of records to return
     * @return records Array of tokenization records
     * @return total The total number of records for this owner
     */
    function getRecordsByOwner(address owner, uint256 offset, uint256 limit) 
        external view returns (TokenizationRecord[] memory records, uint256 total);
    
    /**
     * @notice Get tokenization records by validator
     * @param validatorAddress The validator address
     * @param offset The starting index
     * @param limit The maximum number of records to return
     * @return records Array of tokenization records
     * @return total The total number of records for this validator
     */
    function getRecordsByValidator(string memory validatorAddress, uint256 offset, uint256 limit) 
        external view returns (TokenizationRecord[] memory records, uint256 total);
    
    /**
     * @notice Get the total amount of liquid staked tokens
     * @return amount The total liquid staked amount
     */
    function getTotalLiquidStaked() external view returns (uint256 amount);
    
    /**
     * @notice Get the amount of liquid staked tokens for a validator
     * @param validatorAddress The validator address
     * @return amount The liquid staked amount for the validator
     */
    function getValidatorLiquidStaked(string memory validatorAddress) external view returns (uint256 amount);
    
    /**
     * @notice Get information about a liquid staking token
     * @param tokenDenom The token denomination
     * @return info Information about the token
     */
    function getLiquidStakingTokenInfo(string memory tokenDenom) 
        external view returns (LiquidStakingTokenInfo memory info);
    
    // Transaction Functions
    
    /**
     * @notice Tokenize delegation shares to receive liquid staking tokens
     * @param validatorAddress The validator to tokenize shares from
     * @param amount The amount of shares to tokenize (in shares denomination)
     * @param owner The owner of the liquid staking tokens (use address(0) for msg.sender)
     * @return response The tokenization response containing record ID and token info
     */
    function tokenizeShares(string memory validatorAddress, uint256 amount, address owner) 
        external returns (TokenizeSharesResponse memory response);
    
    /**
     * @notice Redeem liquid staking tokens to restore delegation shares
     * @param tokenDenom The liquid staking token denomination
     * @param amount The amount of tokens to redeem
     * @return response The redemption response containing validator and shares info
     */
    function redeemTokens(string memory tokenDenom, uint256 amount) 
        external returns (RedeemTokensResponse memory response);
}

/**
 * @title LiquidStaking
 * @dev Precompile contract address for liquid staking
 */
library LiquidStaking {
    // The precompile contract address
    ILiquidStaking constant CONTRACT = ILiquidStaking(0x0000000000000000000000000000000000000800);
    
    // Status enum for readability
    enum RecordStatus {
        UNSPECIFIED,
        ACTIVE,
        REDEEMED
    }
}