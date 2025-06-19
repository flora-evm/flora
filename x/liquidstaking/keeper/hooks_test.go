package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/rollchains/flora/x/liquidstaking/types"
)

// MockLiquidStakingHooks is a mock implementation of LiquidStakingHooks for testing
type MockLiquidStakingHooks struct {
	PreTokenizeSharesCalled    bool
	PostTokenizeSharesCalled   bool
	PreRedeemTokensCalled      bool
	PostRedeemTokensCalled     bool
	RecordCreatedCalled        bool
	RecordUpdatedCalled        bool
	RecordDeletedCalled        bool
	RateLimitExceededCalled    bool
	CapReachedCalled           bool
	
	// Store the last call parameters for verification
	LastDelegator              sdk.AccAddress
	LastValidator              sdk.ValAddress
	LastOwner                  sdk.AccAddress
	LastShares                 math.LegacyDec
	LastTokens                 math.Int
	LastDenom                  string
	LastRecordID               uint64
	LastLimitType              string
	LastAddress                string
	LastRejectedAmount         math.Int
	
	// Return error from Pre hooks to test error handling
	PreTokenizeSharesError     error
	PreRedeemTokensError       error
}

var _ types.LiquidStakingHooks = &MockLiquidStakingHooks{}

// OnParametersUpdated mock implementation
func (h *MockLiquidStakingHooks) OnParametersUpdated(
	ctx sdk.Context,
	newParams types.ModuleParams,
) {
	// Mock implementation - no action needed for tests
}

// PreTokenizeShares mock implementation
func (h *MockLiquidStakingHooks) PreTokenizeShares(
	ctx sdk.Context,
	delegator sdk.AccAddress,
	validator sdk.ValAddress,
	owner sdk.AccAddress,
	shares math.LegacyDec,
) error {
	h.PreTokenizeSharesCalled = true
	h.LastDelegator = delegator
	h.LastValidator = validator
	h.LastOwner = owner
	h.LastShares = shares
	return h.PreTokenizeSharesError
}

// PostTokenizeShares mock implementation
func (h *MockLiquidStakingHooks) PostTokenizeShares(
	ctx sdk.Context,
	delegator sdk.AccAddress,
	validator sdk.ValAddress,
	owner sdk.AccAddress,
	shares math.LegacyDec,
	tokens math.Int,
	denom string,
	recordID uint64,
) {
	h.PostTokenizeSharesCalled = true
	h.LastDelegator = delegator
	h.LastValidator = validator
	h.LastOwner = owner
	h.LastShares = shares
	h.LastTokens = tokens
	h.LastDenom = denom
	h.LastRecordID = recordID
}

// PreRedeemTokens mock implementation
func (h *MockLiquidStakingHooks) PreRedeemTokens(
	ctx sdk.Context,
	owner sdk.AccAddress,
	tokens sdk.Coin,
	recordID uint64,
) error {
	h.PreRedeemTokensCalled = true
	h.LastOwner = owner
	h.LastTokens = tokens.Amount
	h.LastDenom = tokens.Denom
	h.LastRecordID = recordID
	return h.PreRedeemTokensError
}

// PostRedeemTokens mock implementation
func (h *MockLiquidStakingHooks) PostRedeemTokens(
	ctx sdk.Context,
	owner sdk.AccAddress,
	validator sdk.ValAddress,
	tokens sdk.Coin,
	shares math.LegacyDec,
	recordID uint64,
) {
	h.PostRedeemTokensCalled = true
	h.LastOwner = owner
	h.LastValidator = validator
	h.LastTokens = tokens.Amount
	h.LastDenom = tokens.Denom
	h.LastShares = shares
	h.LastRecordID = recordID
}

// OnTokenizationRecordCreated mock implementation
func (h *MockLiquidStakingHooks) OnTokenizationRecordCreated(
	ctx sdk.Context,
	record types.TokenizationRecord,
) {
	h.RecordCreatedCalled = true
	h.LastRecordID = record.Id
}

// OnTokenizationRecordUpdated mock implementation
func (h *MockLiquidStakingHooks) OnTokenizationRecordUpdated(
	ctx sdk.Context,
	oldRecord types.TokenizationRecord,
	newRecord types.TokenizationRecord,
) {
	h.RecordUpdatedCalled = true
	h.LastRecordID = newRecord.Id
}

// OnTokenizationRecordDeleted mock implementation
func (h *MockLiquidStakingHooks) OnTokenizationRecordDeleted(
	ctx sdk.Context,
	record types.TokenizationRecord,
) {
	h.RecordDeletedCalled = true
	h.LastRecordID = record.Id
}

// OnRateLimitExceeded mock implementation
func (h *MockLiquidStakingHooks) OnRateLimitExceeded(
	ctx sdk.Context,
	limitType string,
	address string,
	rejectedAmount math.Int,
) {
	h.RateLimitExceededCalled = true
	h.LastLimitType = limitType
	h.LastAddress = address
	h.LastRejectedAmount = rejectedAmount
}

// OnLiquidStakingCapReached mock implementation
func (h *MockLiquidStakingHooks) OnLiquidStakingCapReached(
	ctx sdk.Context,
	capType string,
	validator string,
	currentAmount math.Int,
	capLimit math.Int,
	percentageUsed math.LegacyDec,
) {
	h.CapReachedCalled = true
}

// Reset clears all called flags and stored values
func (h *MockLiquidStakingHooks) Reset() {
	h.PreTokenizeSharesCalled = false
	h.PostTokenizeSharesCalled = false
	h.PreRedeemTokensCalled = false
	h.PostRedeemTokensCalled = false
	h.RecordCreatedCalled = false
	h.RecordUpdatedCalled = false
	h.RecordDeletedCalled = false
	h.RateLimitExceededCalled = false
	h.CapReachedCalled = false
	h.PreTokenizeSharesError = nil
	h.PreRedeemTokensError = nil
}

// TestSetHooks tests setting hooks on the keeper
func TestSetHooks(t *testing.T) {
	k, _ := setupKeeper(t)
	
	// Test setting hooks
	mockHooks := &MockLiquidStakingHooks{}
	k.SetHooks(mockHooks)
	
	// Test that setting hooks twice panics
	require.Panics(t, func() {
		k.SetHooks(mockHooks)
	})
	
	// Test GetHooks returns the set hooks
	hooks := k.GetHooks()
	require.Equal(t, mockHooks, hooks)
}

// TestNoOpHooks tests the no-op hooks implementation
func TestNoOpHooks(t *testing.T) {
	k, ctx := setupKeeper(t)
	
	// When no hooks are set, GetHooks should return NoOpLiquidStakingHooks
	hooks := k.GetHooks()
	_, ok := hooks.(types.NoOpLiquidStakingHooks)
	require.True(t, ok)
	
	// Test that no-op hooks don't error
	err := hooks.PreTokenizeShares(ctx, testAccAddr1, testValAddr1, testAccAddr1, math.LegacyNewDec(100))
	require.NoError(t, err)
	
	err = hooks.PreRedeemTokens(ctx, testAccAddr1, sdk.NewCoin("denom", math.NewInt(100)), 1)
	require.NoError(t, err)
}

// TestMultiHooks tests the multi hooks implementation
func TestMultiHooks(t *testing.T) {
	ctx := sdk.Context{}
	
	// Create multiple mock hooks
	hook1 := &MockLiquidStakingHooks{}
	hook2 := &MockLiquidStakingHooks{}
	hook3 := &MockLiquidStakingHooks{}
	
	// Set hook2 to return an error from PreTokenizeShares
	hook2.PreTokenizeSharesError = types.ErrInvalidShares
	
	// Create multi hooks
	multiHooks := types.NewMultiLiquidStakingHooks(hook1, hook2, hook3)
	
	// Test PreTokenizeShares - should stop at hook2 error
	err := multiHooks.PreTokenizeShares(ctx, testAccAddr1, testValAddr1, testAccAddr1, math.LegacyNewDec(100))
	require.Error(t, err)
	require.Equal(t, types.ErrInvalidShares, err)
	require.True(t, hook1.PreTokenizeSharesCalled)
	require.True(t, hook2.PreTokenizeSharesCalled)
	require.False(t, hook3.PreTokenizeSharesCalled) // Should not be called due to hook2 error
	
	// Reset hooks
	hook1.Reset()
	hook2.Reset()
	hook3.Reset()
	
	// Test PostTokenizeShares - all hooks should be called
	multiHooks.PostTokenizeShares(ctx, testAccAddr1, testValAddr1, testAccAddr1, 
		math.LegacyNewDec(100), math.NewInt(100), "denom", 1)
	require.True(t, hook1.PostTokenizeSharesCalled)
	require.True(t, hook2.PostTokenizeSharesCalled)
	require.True(t, hook3.PostTokenizeSharesCalled)
	
	// Verify all hooks received the same parameters
	require.Equal(t, testAccAddr1, hook1.LastDelegator)
	require.Equal(t, testAccAddr1, hook2.LastDelegator)
	require.Equal(t, testAccAddr1, hook3.LastDelegator)
	require.Equal(t, uint64(1), hook1.LastRecordID)
	require.Equal(t, uint64(1), hook2.LastRecordID)
	require.Equal(t, uint64(1), hook3.LastRecordID)
}