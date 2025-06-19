// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../ILiquidStaking.sol";

/**
 * @title LiquidStakingVault
 * @dev Advanced example: A vault that manages liquid staking positions
 * @notice This contract demonstrates how to build DeFi protocols on top of liquid staking
 */
contract LiquidStakingVault {
    using LiquidStaking for *;
    
    // Structs
    struct UserDeposit {
        uint256 recordId;
        string validator;
        uint256 sharesAmount;
        string tokenDenom;
        uint256 depositTimestamp;
    }
    
    struct VaultStats {
        uint256 totalDeposits;
        uint256 totalShares;
        uint256 activeRecords;
    }
    
    // State
    mapping(address => UserDeposit[]) public userDeposits;
    mapping(uint256 => address) public recordOwners;
    mapping(string => uint256) public validatorShares;
    
    VaultStats public vaultStats;
    address public owner;
    bool public paused;
    
    // Events
    event Deposited(
        address indexed user,
        string validator,
        uint256 recordId,
        uint256 shares,
        string tokenDenom
    );
    
    event Withdrawn(
        address indexed user,
        uint256 recordId,
        uint256 shares,
        bool completed
    );
    
    event EmergencyWithdraw(
        address indexed user,
        uint256[] recordIds
    );
    
    // Modifiers
    modifier onlyOwner() {
        require(msg.sender == owner, "Not owner");
        _;
    }
    
    modifier notPaused() {
        require(!paused, "Vault is paused");
        _;
    }
    
    constructor() {
        owner = msg.sender;
    }
    
    /**
     * @notice Deposit and tokenize shares in the vault
     * @param validator The validator to delegate to
     * @param amount The amount of shares to tokenize
     */
    function deposit(string memory validator, uint256 amount) 
        external 
        notPaused 
        returns (uint256 recordId) 
    {
        // Check module is enabled
        ILiquidStaking.Params memory params = LiquidStaking.CONTRACT.getParams();
        require(params.enabled, "Liquid staking disabled");
        require(amount >= params.minLiquidStakeAmount, "Amount too small");
        
        // Tokenize shares through the precompile
        ILiquidStaking.TokenizeSharesResponse memory response = 
            LiquidStaking.CONTRACT.tokenizeShares(validator, amount, address(this));
        
        // Record the deposit
        UserDeposit memory deposit = UserDeposit({
            recordId: response.recordId,
            validator: validator,
            sharesAmount: amount,
            tokenDenom: response.tokensDenom,
            depositTimestamp: block.timestamp
        });
        
        userDeposits[msg.sender].push(deposit);
        recordOwners[response.recordId] = msg.sender;
        validatorShares[validator] += amount;
        
        // Update vault stats
        vaultStats.totalDeposits++;
        vaultStats.totalShares += amount;
        vaultStats.activeRecords++;
        
        emit Deposited(
            msg.sender,
            validator,
            response.recordId,
            amount,
            response.tokensDenom
        );
        
        return response.recordId;
    }
    
    /**
     * @notice Withdraw a specific deposit
     * @param depositIndex The index of the deposit in user's array
     */
    function withdraw(uint256 depositIndex) external notPaused {
        require(depositIndex < userDeposits[msg.sender].length, "Invalid index");
        
        UserDeposit memory deposit = userDeposits[msg.sender][depositIndex];
        require(recordOwners[deposit.recordId] == msg.sender, "Not owner");
        
        // Get the tokenization record to check status
        ILiquidStaking.TokenizationRecord memory record = 
            LiquidStaking.CONTRACT.getTokenizationRecord(deposit.recordId);
            
        require(record.status == uint8(LiquidStaking.RecordStatus.ACTIVE), "Not active");
        
        // Redeem the tokens
        ILiquidStaking.RedeemTokensResponse memory response = 
            LiquidStaking.CONTRACT.redeemTokens(deposit.tokenDenom, deposit.sharesAmount);
        
        // Update state
        if (response.completed) {
            // Remove the deposit
            _removeDeposit(msg.sender, depositIndex);
            delete recordOwners[deposit.recordId];
            validatorShares[deposit.validator] -= deposit.sharesAmount;
            vaultStats.activeRecords--;
        }
        
        vaultStats.totalShares -= response.sharesAmount;
        
        emit Withdrawn(
            msg.sender,
            deposit.recordId,
            response.sharesAmount,
            response.completed
        );
    }
    
    /**
     * @notice Emergency withdraw all deposits
     * @dev Only callable by deposit owner, bypasses normal restrictions
     */
    function emergencyWithdrawAll() external {
        UserDeposit[] memory deposits = userDeposits[msg.sender];
        require(deposits.length > 0, "No deposits");
        
        uint256[] memory recordIds = new uint256[](deposits.length);
        
        for (uint256 i = 0; i < deposits.length; i++) {
            UserDeposit memory deposit = deposits[i];
            
            try LiquidStaking.CONTRACT.redeemTokens(
                deposit.tokenDenom, 
                deposit.sharesAmount
            ) returns (ILiquidStaking.RedeemTokensResponse memory response) {
                if (response.completed) {
                    delete recordOwners[deposit.recordId];
                    validatorShares[deposit.validator] -= deposit.sharesAmount;
                    vaultStats.activeRecords--;
                }
                vaultStats.totalShares -= response.sharesAmount;
            } catch {
                // Continue with other withdrawals even if one fails
            }
            
            recordIds[i] = deposit.recordId;
        }
        
        // Clear all deposits for the user
        delete userDeposits[msg.sender];
        vaultStats.totalDeposits -= deposits.length;
        
        emit EmergencyWithdraw(msg.sender, recordIds);
    }
    
    /**
     * @notice Get user's deposit count
     * @param user The user address
     * @return count Number of deposits
     */
    function getUserDepositCount(address user) external view returns (uint256) {
        return userDeposits[user].length;
    }
    
    /**
     * @notice Get paginated user deposits
     * @param user The user address
     * @param offset Starting index
     * @param limit Maximum number to return
     * @return deposits Array of user deposits
     */
    function getUserDeposits(address user, uint256 offset, uint256 limit) 
        external 
        view 
        returns (UserDeposit[] memory deposits) 
    {
        UserDeposit[] memory allDeposits = userDeposits[user];
        uint256 total = allDeposits.length;
        
        if (offset >= total) {
            return new UserDeposit[](0);
        }
        
        uint256 end = offset + limit;
        if (end > total) {
            end = total;
        }
        
        deposits = new UserDeposit[](end - offset);
        for (uint256 i = 0; i < deposits.length; i++) {
            deposits[i] = allDeposits[offset + i];
        }
    }
    
    /**
     * @notice Calculate total liquid staked amount for a validator
     * @param validator The validator address
     * @return precompileAmount Amount from precompile
     * @return vaultAmount Amount tracked by vault
     */
    function getValidatorStats(string memory validator) 
        external 
        view 
        returns (uint256 precompileAmount, uint256 vaultAmount) 
    {
        precompileAmount = LiquidStaking.CONTRACT.getValidatorLiquidStaked(validator);
        vaultAmount = validatorShares[validator];
    }
    
    // Admin functions
    
    /**
     * @notice Pause the vault
     */
    function pause() external onlyOwner {
        paused = true;
    }
    
    /**
     * @notice Unpause the vault
     */
    function unpause() external onlyOwner {
        paused = false;
    }
    
    /**
     * @notice Transfer ownership
     * @param newOwner The new owner address
     */
    function transferOwnership(address newOwner) external onlyOwner {
        require(newOwner != address(0), "Invalid address");
        owner = newOwner;
    }
    
    // Internal functions
    
    /**
     * @dev Remove a deposit from user's array
     * @param user The user address
     * @param index The deposit index to remove
     */
    function _removeDeposit(address user, uint256 index) internal {
        UserDeposit[] storage deposits = userDeposits[user];
        require(index < deposits.length, "Invalid index");
        
        // Move the last element to the deleted spot
        if (index < deposits.length - 1) {
            deposits[index] = deposits[deposits.length - 1];
        }
        
        // Remove the last element
        deposits.pop();
    }
}