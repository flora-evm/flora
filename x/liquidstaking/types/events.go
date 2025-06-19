package types

// Event types for the liquid staking module
const (
	EventTypeTokenizeShares     = "tokenize_shares"
	EventTypeRedeemTokens       = "redeem_tokens"
	EventTypeUpdateParams       = "update_params"
	EventTypeRecordCreated      = "tokenization_record_created"
	EventTypeRecordUpdated      = "tokenization_record_updated"
	EventTypeRecordDeleted      = "tokenization_record_deleted"
	EventTypeLiquidStakingCap   = "liquid_staking_cap"
	EventTypeRateLimitExceeded  = "rate_limit_exceeded"
	EventTypeRateLimitWarning   = "rate_limit_warning"
	EventTypeActivityTracked    = "activity_tracked"
	
	// IBC event types
	EventTypeLiquidStakingIBCTransfer = "liquid_staking_ibc_transfer"
	EventTypeLiquidStakingIBCReceived = "liquid_staking_ibc_received"
	EventTypeLiquidStakingIBCAck      = "liquid_staking_ibc_ack"
	EventTypeLiquidStakingIBCTimeout  = "liquid_staking_ibc_timeout"
	
	// Governance event types
	EventTypeParameterUpdate          = "parameter_update"
	
	// Exchange rate event types
	EventTypeExchangeRateUpdated      = "exchange_rate_updated"
	
	// Auto-compound event types
	EventTypeRewardsCompounded        = "rewards_compounded"
	EventTypeAutoCompoundStarted      = "auto_compound_started"
	EventTypeAutoCompoundCompleted    = "auto_compound_completed"
	EventTypeAutoCompoundFailed       = "auto_compound_failed"
)

// Event attribute keys
const (
	AttributeKeyDelegator       = "delegator"
	AttributeKeyValidator       = "validator"
	AttributeKeyOwner           = "owner"
	AttributeKeyShares          = "shares"
	AttributeKeySharesTokenized = "shares_tokenized"
	AttributeKeyTokens          = "tokens"
	AttributeKeyTokensMinted    = "tokens_minted"
	AttributeKeyTokensBurned    = "tokens_burned"
	AttributeKeySharesRestored  = "shares_restored"
	AttributeKeyDenom           = "denom"
	AttributeKeyAmount          = "amount"
	AttributeKeyRecordID        = "record_id"
	AttributeKeyRecordIDs       = "record_ids"
	
	// Parameter attributes
	AttributeKeyParamKey        = "param_key"
	AttributeKeyParamOldValue   = "param_old_value"
	AttributeKeyParamNewValue   = "param_new_value"
	
	// Cap attributes
	AttributeKeyCapType         = "cap_type"
	AttributeKeyCurrentAmount   = "current_amount"
	AttributeKeyCapLimit        = "cap_limit"
	AttributeKeyPercentageUsed  = "percentage_used"
	
	// Rate limit attributes
	AttributeKeyLimitType       = "limit_type"
	AttributeKeyLimitThreshold  = "limit_threshold"
	AttributeKeyCurrentUsage    = "current_usage"
	AttributeKeyMaxUsage        = "max_usage"
	AttributeKeyWindowStart     = "window_start"
	AttributeKeyWindowEnd       = "window_end"
	AttributeKeyRejectedAmount  = "rejected_amount"
	AttributeKeyAddress         = "address"
	
	// Generic attributes
	AttributeKeyModule          = "module"
	AttributeKeySender          = "sender"
	AttributeKeyAction          = "action"
	
	// Governance attributes
	AttributeKeyProposalTitle   = "proposal_title"
	AttributeKeyProposalID      = "proposal_id"
	
	// IBC attributes
	AttributeKeyReceiver        = "receiver"
	AttributeKeySourcePort      = "source_port"
	AttributeKeySourceChannel   = "source_channel"
	AttributeKeySourceChainId   = "source_chain_id"
	AttributeKeySuccess         = "success"
	
	// Exchange rate attributes
	AttributeKeyOldRate         = "old_rate"
	AttributeKeyNewRate         = "new_rate"
	AttributeKeyTimestamp       = "timestamp"
	
	// Auto-compound attributes
	AttributeKeyValidatorCount  = "validator_count"
	AttributeKeyTotalCompounded = "total_compounded"
	AttributeKeyBlockHeight     = "block_height"
	AttributeKeyCompoundAmount  = "compound_amount"
	AttributeKeyError           = "error"
	
	AttributeValueCategory      = ModuleName
)

// Event attribute values
const (
	AttributeValueActionTokenize = "tokenize"
	AttributeValueActionRedeem   = "redeem"
	AttributeValueActionUpdate   = "update"
	AttributeValueActionCreate   = "create"
	AttributeValueActionDelete   = "delete"
	
	// Cap types
	AttributeValueCapTypeGlobal    = "global"
	AttributeValueCapTypeValidator = "validator"
)