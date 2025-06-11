package types

// Event types for the liquid staking module
const (
	EventTypeTokenizeShares = "tokenize_shares"
	EventTypeRedeemTokens   = "redeem_tokens" // For future use in Stage 4
	
	AttributeKeyDelegator = "delegator"
	AttributeKeyValidator = "validator"
	AttributeKeyOwner     = "owner"
	AttributeKeyShares    = "shares"
	AttributeKeyDenom     = "denom"
	AttributeKeyAmount    = "amount"
	AttributeKeyRecordID  = "record_id"
	
	AttributeValueCategory = ModuleName
)