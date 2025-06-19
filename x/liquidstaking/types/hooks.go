package types

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// LiquidStakingHooks defines the hooks interface for the liquid staking module
// External modules can implement these hooks to react to liquid staking events
type LiquidStakingHooks interface {
	// PreTokenizeShares is called before shares are tokenized
	// It can be used to validate or reject the tokenization
	// Returning an error will abort the tokenization
	PreTokenizeShares(
		ctx sdk.Context,
		delegator sdk.AccAddress,
		validator sdk.ValAddress,
		owner sdk.AccAddress,
		shares math.LegacyDec,
	) error

	// PostTokenizeShares is called after shares are successfully tokenized
	// It can be used to update external module state or emit additional events
	PostTokenizeShares(
		ctx sdk.Context,
		delegator sdk.AccAddress,
		validator sdk.ValAddress,
		owner sdk.AccAddress,
		shares math.LegacyDec,
		tokens math.Int,
		denom string,
		recordID uint64,
	)

	// PreRedeemTokens is called before tokens are redeemed
	// It can be used to validate or reject the redemption
	// Returning an error will abort the redemption
	PreRedeemTokens(
		ctx sdk.Context,
		owner sdk.AccAddress,
		tokens sdk.Coin,
		recordID uint64,
	) error

	// PostRedeemTokens is called after tokens are successfully redeemed
	// It can be used to update external module state or emit additional events
	PostRedeemTokens(
		ctx sdk.Context,
		owner sdk.AccAddress,
		validator sdk.ValAddress,
		tokens sdk.Coin,
		shares math.LegacyDec,
		recordID uint64,
	)

	// OnTokenizationRecordCreated is called when a new tokenization record is created
	OnTokenizationRecordCreated(
		ctx sdk.Context,
		record TokenizationRecord,
	)

	// OnTokenizationRecordUpdated is called when a tokenization record is updated
	OnTokenizationRecordUpdated(
		ctx sdk.Context,
		oldRecord TokenizationRecord,
		newRecord TokenizationRecord,
	)

	// OnTokenizationRecordDeleted is called when a tokenization record is deleted
	OnTokenizationRecordDeleted(
		ctx sdk.Context,
		record TokenizationRecord,
	)

	// OnRateLimitExceeded is called when a rate limit is exceeded
	OnRateLimitExceeded(
		ctx sdk.Context,
		limitType string,
		address string,
		rejectedAmount math.Int,
	)

	// OnLiquidStakingCapReached is called when approaching or exceeding liquid staking caps
	OnLiquidStakingCapReached(
		ctx sdk.Context,
		capType string,
		validator string,
		currentAmount math.Int,
		capLimit math.Int,
		percentageUsed math.LegacyDec,
	)

	// OnParametersUpdated is called when module parameters are updated via governance
	OnParametersUpdated(
		ctx sdk.Context,
		newParams ModuleParams,
	)
}

// MultiLiquidStakingHooks combines multiple liquid staking hooks into one
type MultiLiquidStakingHooks []LiquidStakingHooks

// NewMultiLiquidStakingHooks creates a new MultiLiquidStakingHooks instance
func NewMultiLiquidStakingHooks(hooks ...LiquidStakingHooks) MultiLiquidStakingHooks {
	return hooks
}

// PreTokenizeShares calls all PreTokenizeShares hooks in order
// Returns on first error
func (h MultiLiquidStakingHooks) PreTokenizeShares(
	ctx sdk.Context,
	delegator sdk.AccAddress,
	validator sdk.ValAddress,
	owner sdk.AccAddress,
	shares math.LegacyDec,
) error {
	for i := range h {
		if err := h[i].PreTokenizeShares(ctx, delegator, validator, owner, shares); err != nil {
			return err
		}
	}
	return nil
}

// PostTokenizeShares calls all PostTokenizeShares hooks
func (h MultiLiquidStakingHooks) PostTokenizeShares(
	ctx sdk.Context,
	delegator sdk.AccAddress,
	validator sdk.ValAddress,
	owner sdk.AccAddress,
	shares math.LegacyDec,
	tokens math.Int,
	denom string,
	recordID uint64,
) {
	for i := range h {
		h[i].PostTokenizeShares(ctx, delegator, validator, owner, shares, tokens, denom, recordID)
	}
}

// PreRedeemTokens calls all PreRedeemTokens hooks in order
// Returns on first error
func (h MultiLiquidStakingHooks) PreRedeemTokens(
	ctx sdk.Context,
	owner sdk.AccAddress,
	tokens sdk.Coin,
	recordID uint64,
) error {
	for i := range h {
		if err := h[i].PreRedeemTokens(ctx, owner, tokens, recordID); err != nil {
			return err
		}
	}
	return nil
}

// PostRedeemTokens calls all PostRedeemTokens hooks
func (h MultiLiquidStakingHooks) PostRedeemTokens(
	ctx sdk.Context,
	owner sdk.AccAddress,
	validator sdk.ValAddress,
	tokens sdk.Coin,
	shares math.LegacyDec,
	recordID uint64,
) {
	for i := range h {
		h[i].PostRedeemTokens(ctx, owner, validator, tokens, shares, recordID)
	}
}

// OnTokenizationRecordCreated calls all OnTokenizationRecordCreated hooks
func (h MultiLiquidStakingHooks) OnTokenizationRecordCreated(
	ctx sdk.Context,
	record TokenizationRecord,
) {
	for i := range h {
		h[i].OnTokenizationRecordCreated(ctx, record)
	}
}

// OnTokenizationRecordUpdated calls all OnTokenizationRecordUpdated hooks
func (h MultiLiquidStakingHooks) OnTokenizationRecordUpdated(
	ctx sdk.Context,
	oldRecord TokenizationRecord,
	newRecord TokenizationRecord,
) {
	for i := range h {
		h[i].OnTokenizationRecordUpdated(ctx, oldRecord, newRecord)
	}
}

// OnTokenizationRecordDeleted calls all OnTokenizationRecordDeleted hooks
func (h MultiLiquidStakingHooks) OnTokenizationRecordDeleted(
	ctx sdk.Context,
	record TokenizationRecord,
) {
	for i := range h {
		h[i].OnTokenizationRecordDeleted(ctx, record)
	}
}

// OnRateLimitExceeded calls all OnRateLimitExceeded hooks
func (h MultiLiquidStakingHooks) OnRateLimitExceeded(
	ctx sdk.Context,
	limitType string,
	address string,
	rejectedAmount math.Int,
) {
	for i := range h {
		h[i].OnRateLimitExceeded(ctx, limitType, address, rejectedAmount)
	}
}

// OnLiquidStakingCapReached calls all OnLiquidStakingCapReached hooks
func (h MultiLiquidStakingHooks) OnLiquidStakingCapReached(
	ctx sdk.Context,
	capType string,
	validator string,
	currentAmount math.Int,
	capLimit math.Int,
	percentageUsed math.LegacyDec,
) {
	for i := range h {
		h[i].OnLiquidStakingCapReached(ctx, capType, validator, currentAmount, capLimit, percentageUsed)
	}
}

// OnParametersUpdated calls all OnParametersUpdated hooks
func (h MultiLiquidStakingHooks) OnParametersUpdated(
	ctx sdk.Context,
	newParams ModuleParams,
) {
	for i := range h {
		h[i].OnParametersUpdated(ctx, newParams)
	}
}

// NoOpLiquidStakingHooks is a no-op implementation of LiquidStakingHooks
// It can be used as a default when no hooks are set
type NoOpLiquidStakingHooks struct{}

var _ LiquidStakingHooks = NoOpLiquidStakingHooks{}

// PreTokenizeShares no-op implementation
func (NoOpLiquidStakingHooks) PreTokenizeShares(
	ctx sdk.Context,
	delegator sdk.AccAddress,
	validator sdk.ValAddress,
	owner sdk.AccAddress,
	shares math.LegacyDec,
) error {
	return nil
}

// PostTokenizeShares no-op implementation
func (NoOpLiquidStakingHooks) PostTokenizeShares(
	ctx sdk.Context,
	delegator sdk.AccAddress,
	validator sdk.ValAddress,
	owner sdk.AccAddress,
	shares math.LegacyDec,
	tokens math.Int,
	denom string,
	recordID uint64,
) {
}

// PreRedeemTokens no-op implementation
func (NoOpLiquidStakingHooks) PreRedeemTokens(
	ctx sdk.Context,
	owner sdk.AccAddress,
	tokens sdk.Coin,
	recordID uint64,
) error {
	return nil
}

// PostRedeemTokens no-op implementation
func (NoOpLiquidStakingHooks) PostRedeemTokens(
	ctx sdk.Context,
	owner sdk.AccAddress,
	validator sdk.ValAddress,
	tokens sdk.Coin,
	shares math.LegacyDec,
	recordID uint64,
) {
}

// OnTokenizationRecordCreated no-op implementation
func (NoOpLiquidStakingHooks) OnTokenizationRecordCreated(
	ctx sdk.Context,
	record TokenizationRecord,
) {
}

// OnTokenizationRecordUpdated no-op implementation
func (NoOpLiquidStakingHooks) OnTokenizationRecordUpdated(
	ctx sdk.Context,
	oldRecord TokenizationRecord,
	newRecord TokenizationRecord,
) {
}

// OnTokenizationRecordDeleted no-op implementation
func (NoOpLiquidStakingHooks) OnTokenizationRecordDeleted(
	ctx sdk.Context,
	record TokenizationRecord,
) {
}

// OnRateLimitExceeded no-op implementation
func (NoOpLiquidStakingHooks) OnRateLimitExceeded(
	ctx sdk.Context,
	limitType string,
	address string,
	rejectedAmount math.Int,
) {
}

// OnLiquidStakingCapReached no-op implementation
func (NoOpLiquidStakingHooks) OnLiquidStakingCapReached(
	ctx sdk.Context,
	capType string,
	validator string,
	currentAmount math.Int,
	capLimit math.Int,
	percentageUsed math.LegacyDec,
) {
}

// OnParametersUpdated no-op implementation
func (NoOpLiquidStakingHooks) OnParametersUpdated(
	ctx sdk.Context,
	newParams ModuleParams,
) {
}