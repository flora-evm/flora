package types

import (
	"fmt"
	
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// TokenizeSharesEvent is emitted when shares are tokenized
type TokenizeSharesEvent struct {
	Delegator      string
	Validator      string
	Owner          string
	SharesAmount   string
	TokensMinted   string
	Denom          string
	RecordID       uint64
}

// ToEvent converts the typed event to sdk.Event
func (e TokenizeSharesEvent) ToEvent() sdk.Event {
	return sdk.NewEvent(
		EventTypeTokenizeShares,
		sdk.NewAttribute(AttributeKeyDelegator, e.Delegator),
		sdk.NewAttribute(AttributeKeyValidator, e.Validator),
		sdk.NewAttribute(AttributeKeyOwner, e.Owner),
		sdk.NewAttribute(AttributeKeyShares, e.SharesAmount),
		sdk.NewAttribute(AttributeKeyTokensMinted, e.TokensMinted),
		sdk.NewAttribute(AttributeKeyDenom, e.Denom),
		sdk.NewAttribute(AttributeKeyRecordID, fmt.Sprintf("%d", e.RecordID)),
	)
}

// RedeemTokensEvent is emitted when tokens are redeemed
type RedeemTokensEvent struct {
	Owner          string
	Validator      string
	TokensBurned   string
	SharesRestored string
	Denom          string
	RecordID       uint64
}

// ToEvent converts the typed event to sdk.Event
func (e RedeemTokensEvent) ToEvent() sdk.Event {
	return sdk.NewEvent(
		EventTypeRedeemTokens,
		sdk.NewAttribute(AttributeKeyOwner, e.Owner),
		sdk.NewAttribute(AttributeKeyValidator, e.Validator),
		sdk.NewAttribute(AttributeKeyTokensBurned, e.TokensBurned),
		sdk.NewAttribute(AttributeKeySharesRestored, e.SharesRestored),
		sdk.NewAttribute(AttributeKeyDenom, e.Denom),
		sdk.NewAttribute(AttributeKeyRecordID, fmt.Sprintf("%d", e.RecordID)),
	)
}

// UpdateParamsEvent is emitted when module parameters are updated
type UpdateParamsEvent struct {
	Authority string
	Changes   []EventParamChange
}

// EventParamChange represents a single parameter change in events
type EventParamChange struct {
	Key      string
	OldValue string
	NewValue string
}

// ToEvent converts the typed event to sdk.Event
func (e UpdateParamsEvent) ToEvent() sdk.Event {
	attributes := []sdk.Attribute{
		sdk.NewAttribute(AttributeKeySender, e.Authority),
		sdk.NewAttribute(AttributeKeyAction, AttributeValueActionUpdate),
	}
	
	// Add each parameter change as separate attributes
	for i, change := range e.Changes {
		prefix := fmt.Sprintf("change_%d_", i)
		attributes = append(attributes,
			sdk.NewAttribute(prefix+AttributeKeyParamKey, change.Key),
			sdk.NewAttribute(prefix+AttributeKeyParamOldValue, change.OldValue),
			sdk.NewAttribute(prefix+AttributeKeyParamNewValue, change.NewValue),
		)
	}
	
	return sdk.NewEvent(EventTypeUpdateParams, attributes...)
}

// TokenizationRecordCreatedEvent is emitted when a new record is created
type TokenizationRecordCreatedEvent struct {
	RecordID        uint64
	Validator       string
	Owner           string
	SharesTokenized string
	Denom           string
}

// ToEvent converts the typed event to sdk.Event
func (e TokenizationRecordCreatedEvent) ToEvent() sdk.Event {
	return sdk.NewEvent(
		EventTypeRecordCreated,
		sdk.NewAttribute(AttributeKeyRecordID, fmt.Sprintf("%d", e.RecordID)),
		sdk.NewAttribute(AttributeKeyValidator, e.Validator),
		sdk.NewAttribute(AttributeKeyOwner, e.Owner),
		sdk.NewAttribute(AttributeKeySharesTokenized, e.SharesTokenized),
		sdk.NewAttribute(AttributeKeyDenom, e.Denom),
		sdk.NewAttribute(AttributeKeyAction, AttributeValueActionCreate),
	)
}

// TokenizationRecordUpdatedEvent is emitted when a record is updated
type TokenizationRecordUpdatedEvent struct {
	RecordID           uint64
	OldSharesTokenized string
	NewSharesTokenized string
}

// ToEvent converts the typed event to sdk.Event
func (e TokenizationRecordUpdatedEvent) ToEvent() sdk.Event {
	return sdk.NewEvent(
		EventTypeRecordUpdated,
		sdk.NewAttribute(AttributeKeyRecordID, fmt.Sprintf("%d", e.RecordID)),
		sdk.NewAttribute("old_shares_tokenized", e.OldSharesTokenized),
		sdk.NewAttribute("new_shares_tokenized", e.NewSharesTokenized),
		sdk.NewAttribute(AttributeKeyAction, AttributeValueActionUpdate),
	)
}

// TokenizationRecordDeletedEvent is emitted when a record is deleted
type TokenizationRecordDeletedEvent struct {
	RecordID  uint64
	Validator string
	Owner     string
	Denom     string
}

// ToEvent converts the typed event to sdk.Event
func (e TokenizationRecordDeletedEvent) ToEvent() sdk.Event {
	return sdk.NewEvent(
		EventTypeRecordDeleted,
		sdk.NewAttribute(AttributeKeyRecordID, fmt.Sprintf("%d", e.RecordID)),
		sdk.NewAttribute(AttributeKeyValidator, e.Validator),
		sdk.NewAttribute(AttributeKeyOwner, e.Owner),
		sdk.NewAttribute(AttributeKeyDenom, e.Denom),
		sdk.NewAttribute(AttributeKeyAction, AttributeValueActionDelete),
	)
}

// LiquidStakingCapEvent is emitted when approaching or exceeding caps
type LiquidStakingCapEvent struct {
	CapType        string // "global" or "validator"
	Validator      string // empty for global cap
	CurrentAmount  string
	CapLimit       string
	PercentageUsed string
}

// ToEvent converts the typed event to sdk.Event
func (e LiquidStakingCapEvent) ToEvent() sdk.Event {
	attributes := []sdk.Attribute{
		sdk.NewAttribute(AttributeKeyCapType, e.CapType),
		sdk.NewAttribute(AttributeKeyCurrentAmount, e.CurrentAmount),
		sdk.NewAttribute(AttributeKeyCapLimit, e.CapLimit),
		sdk.NewAttribute(AttributeKeyPercentageUsed, e.PercentageUsed),
	}
	
	if e.Validator != "" {
		attributes = append(attributes, sdk.NewAttribute(AttributeKeyValidator, e.Validator))
	}
	
	return sdk.NewEvent(EventTypeLiquidStakingCap, attributes...)
}

// EmitTokenizeSharesEvent is a helper to emit a typed tokenize shares event
func EmitTokenizeSharesEvent(ctx sdk.Context, event TokenizeSharesEvent) {
	ctx.EventManager().EmitEvent(event.ToEvent())
	
	// Also emit the standard message event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, event.Delegator),
			sdk.NewAttribute(sdk.AttributeKeyAction, AttributeValueActionTokenize),
		),
	)
}

// EmitRedeemTokensEvent is a helper to emit a typed redeem tokens event
func EmitRedeemTokensEvent(ctx sdk.Context, event RedeemTokensEvent) {
	ctx.EventManager().EmitEvent(event.ToEvent())
	
	// Also emit the standard message event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, event.Owner),
			sdk.NewAttribute(sdk.AttributeKeyAction, AttributeValueActionRedeem),
		),
	)
}

// EmitUpdateParamsEvent is a helper to emit a typed update params event
func EmitUpdateParamsEvent(ctx sdk.Context, event UpdateParamsEvent) {
	ctx.EventManager().EmitEvent(event.ToEvent())
	
	// Also emit the standard message event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, event.Authority),
			sdk.NewAttribute(sdk.AttributeKeyAction, AttributeValueActionUpdate),
		),
	)
}

// RateLimitExceededEvent is emitted when a rate limit is exceeded
type RateLimitExceededEvent struct {
	LimitType      string // "global", "validator", or "user"
	Address        string // The address that exceeded the limit
	CurrentUsage   string // Current usage amount
	MaxUsage       string // Maximum allowed amount
	RejectedAmount string // Amount that was rejected
	WindowEnd      string // When the current window ends
}

// ToEvent converts the typed event to sdk.Event
func (e RateLimitExceededEvent) ToEvent() sdk.Event {
	return sdk.NewEvent(
		EventTypeRateLimitExceeded,
		sdk.NewAttribute(AttributeKeyLimitType, e.LimitType),
		sdk.NewAttribute(AttributeKeyAddress, e.Address),
		sdk.NewAttribute(AttributeKeyCurrentUsage, e.CurrentUsage),
		sdk.NewAttribute(AttributeKeyMaxUsage, e.MaxUsage),
		sdk.NewAttribute(AttributeKeyRejectedAmount, e.RejectedAmount),
		sdk.NewAttribute(AttributeKeyWindowEnd, e.WindowEnd),
	)
}

// RateLimitWarningEvent is emitted when approaching a rate limit threshold
type RateLimitWarningEvent struct {
	LimitType       string // "global", "validator", or "user"
	Address         string // The address approaching limit
	CurrentUsage    string // Current usage amount
	MaxUsage        string // Maximum allowed amount
	PercentageUsed  string // Percentage of limit used
	LimitThreshold  string // Warning threshold percentage
}

// ToEvent converts the typed event to sdk.Event
func (e RateLimitWarningEvent) ToEvent() sdk.Event {
	return sdk.NewEvent(
		EventTypeRateLimitWarning,
		sdk.NewAttribute(AttributeKeyLimitType, e.LimitType),
		sdk.NewAttribute(AttributeKeyAddress, e.Address),
		sdk.NewAttribute(AttributeKeyCurrentUsage, e.CurrentUsage),
		sdk.NewAttribute(AttributeKeyMaxUsage, e.MaxUsage),
		sdk.NewAttribute(AttributeKeyPercentageUsed, e.PercentageUsed),
		sdk.NewAttribute(AttributeKeyLimitThreshold, e.LimitThreshold),
	)
}

// ActivityTrackedEvent is emitted when tokenization activity is tracked
type ActivityTrackedEvent struct {
	LimitType     string // "global", "validator", or "user"
	Address       string // The address whose activity was tracked
	Amount        string // Amount tokenized
	TotalAmount   string // Total amount in current window
	ActivityCount string // Number of operations in current window
	WindowStart   string // When the current window started
	WindowEnd     string // When the current window ends
}

// ToEvent converts the typed event to sdk.Event
func (e ActivityTrackedEvent) ToEvent() sdk.Event {
	return sdk.NewEvent(
		EventTypeActivityTracked,
		sdk.NewAttribute(AttributeKeyLimitType, e.LimitType),
		sdk.NewAttribute(AttributeKeyAddress, e.Address),
		sdk.NewAttribute(AttributeKeyAmount, e.Amount),
		sdk.NewAttribute(AttributeKeyCurrentAmount, e.TotalAmount),
		sdk.NewAttribute("activity_count", e.ActivityCount),
		sdk.NewAttribute(AttributeKeyWindowStart, e.WindowStart),
		sdk.NewAttribute(AttributeKeyWindowEnd, e.WindowEnd),
	)
}