package keeper_test

import (
	"fmt"
	"testing"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/rollchains/flora/x/liquidstaking/keeper"
	"github.com/rollchains/flora/x/liquidstaking/types"
)

// TestRateLimitEventEmission tests that rate limit events are properly emitted
func TestRateLimitEventEmission(t *testing.T) {
	k, ctx := setupKeeper(t)
	
	// Enable module
	params := types.DefaultParams()
	params.Enabled = true
	require.NoError(t, k.SetParams(ctx, params))
	
	// Mock total bonded tokens
	totalBonded := math.NewInt(10000000)
	mockStakingKeeper.TotalBondedTokensFn = func(c sdk.Context) (math.Int, error) {
		return totalBonded, nil
	}
	
	t.Run("rate limit exceeded event", func(t *testing.T) {
		// Set activity near limit
		activity := keeper.TokenizationActivity{
			TotalAmount:   math.NewInt(499000), // Just under 5% of 10M
			LastActivity:  ctx.BlockTime(),
			ActivityCount: 10,
		}
		k.SetGlobalTokenizationActivity(ctx, activity)
		
		// Try to exceed limit
		err := k.CheckGlobalRateLimit(ctx, math.NewInt(2000))
		require.Error(t, err)
		require.ErrorIs(t, err, types.ErrRateLimitExceeded)
		
		// Check that event was emitted
		events := ctx.EventManager().Events()
		found := false
		for _, event := range events {
			if event.Type == types.EventTypeRateLimitExceeded {
				found = true
				
				// Verify event attributes
				attrs := event.Attributes
				require.Greater(t, len(attrs), 0)
				
				// Check specific attributes
				for _, attr := range attrs {
					switch string(attr.Key) {
					case types.AttributeKeyLimitType:
						require.Equal(t, "global", string(attr.Value))
					case types.AttributeKeyAddress:
						require.Equal(t, "global", string(attr.Value))
					case types.AttributeKeyCurrentUsage:
						require.Equal(t, "499000", string(attr.Value))
					case types.AttributeKeyMaxUsage:
						require.Equal(t, "500000", string(attr.Value))
					case types.AttributeKeyRejectedAmount:
						require.Equal(t, "2000", string(attr.Value))
					}
				}
			}
		}
		require.True(t, found, "rate limit exceeded event not found")
	})
	
	t.Run("rate limit warning event", func(t *testing.T) {
		ctx = ctx.WithEventManager(sdk.NewEventManager()) // Reset events
		
		// Set activity at 85% of limit (above 80% threshold)
		activity := keeper.TokenizationActivity{
			TotalAmount:   math.NewInt(425000), // 85% of 500,000 limit
			LastActivity:  ctx.BlockTime(),
			ActivityCount: 10,
		}
		k.SetGlobalTokenizationActivity(ctx, activity)
		
		// Check limit - should pass but emit warning
		err := k.CheckGlobalRateLimit(ctx, math.NewInt(1000))
		require.NoError(t, err)
		
		// Check that warning event was emitted
		events := ctx.EventManager().Events()
		found := false
		for _, event := range events {
			if event.Type == types.EventTypeRateLimitWarning {
				found = true
				
				// Verify event attributes
				attrs := event.Attributes
				for _, attr := range attrs {
					switch string(attr.Key) {
					case types.AttributeKeyLimitType:
						require.Equal(t, "global", string(attr.Value))
					case types.AttributeKeyPercentageUsed:
						require.Equal(t, "85", string(attr.Value)) // (425000 + 1000) / 500000 * 100
					case types.AttributeKeyLimitThreshold:
						require.Equal(t, "80", string(attr.Value))
					}
				}
			}
		}
		require.True(t, found, "rate limit warning event not found")
	})
	
	t.Run("count limit exceeded event", func(t *testing.T) {
		ctx = ctx.WithEventManager(sdk.NewEventManager()) // Reset events
		
		// Set activity at count limit
		activity := keeper.TokenizationActivity{
			TotalAmount:   math.NewInt(100000),
			LastActivity:  ctx.BlockTime(),
			ActivityCount: 100, // At the limit
		}
		k.SetGlobalTokenizationActivity(ctx, activity)
		
		// Try to exceed count limit
		err := k.CheckGlobalRateLimit(ctx, math.NewInt(1000))
		require.Error(t, err)
		require.ErrorIs(t, err, types.ErrRateLimitExceeded)
		
		// Check that event was emitted
		events := ctx.EventManager().Events()
		found := false
		for _, event := range events {
			if event.Type == types.EventTypeRateLimitExceeded {
				found = true
				
				// Verify event attributes
				attrs := event.Attributes
				for _, attr := range attrs {
					switch string(attr.Key) {
					case types.AttributeKeyCurrentUsage:
						require.Equal(t, "100", string(attr.Value))
					case types.AttributeKeyMaxUsage:
						require.Equal(t, "100", string(attr.Value))
					}
				}
			}
		}
		require.True(t, found, "count limit exceeded event not found")
	})
}

// TestActivityTrackedEvent tests that activity tracking events are properly emitted
func TestActivityTrackedEvent(t *testing.T) {
	k, ctx := setupKeeper(t)
	
	// Update tokenization activity
	k.UpdateTokenizationActivity(ctx, testValAddr1.String(), testAccAddr1.String(), math.NewInt(100000))
	
	// Check that activity tracked events were emitted (3 total: global, validator, user)
	events := ctx.EventManager().Events()
	activityEvents := 0
	
	for _, event := range events {
		if event.Type == types.EventTypeActivityTracked {
			activityEvents++
			
			// Verify event has required attributes
			attrs := event.Attributes
			hasLimitType := false
			hasAddress := false
			hasAmount := false
			hasTotalAmount := false
			hasActivityCount := false
			hasWindowStart := false
			hasWindowEnd := false
			
			for _, attr := range attrs {
				switch string(attr.Key) {
				case types.AttributeKeyLimitType:
					hasLimitType = true
					// Should be one of: global, validator, user
					limitType := string(attr.Value)
					require.Contains(t, []string{"global", "validator", "user"}, limitType)
				case types.AttributeKeyAddress:
					hasAddress = true
				case types.AttributeKeyAmount:
					hasAmount = true
					require.Equal(t, "100000", string(attr.Value))
				case types.AttributeKeyCurrentAmount:
					hasTotalAmount = true
					require.Equal(t, "100000", string(attr.Value)) // First activity
				case "activity_count":
					hasActivityCount = true
					require.Equal(t, "1", string(attr.Value)) // First activity
				case types.AttributeKeyWindowStart:
					hasWindowStart = true
				case types.AttributeKeyWindowEnd:
					hasWindowEnd = true
				}
			}
			
			// Verify all required attributes are present
			require.True(t, hasLimitType, "missing limit type")
			require.True(t, hasAddress, "missing address")
			require.True(t, hasAmount, "missing amount")
			require.True(t, hasTotalAmount, "missing total amount")
			require.True(t, hasActivityCount, "missing activity count")
			require.True(t, hasWindowStart, "missing window start")
			require.True(t, hasWindowEnd, "missing window end")
		}
	}
	
	// Should have 3 activity events: global, validator, user
	require.Equal(t, 3, activityEvents, "expected 3 activity tracked events")
}

// TestHookIntegrationWithRateLimit tests that hooks are called when rate limits are exceeded
func TestHookIntegrationWithRateLimit(t *testing.T) {
	k, ctx := setupKeeper(t)
	
	// Set up mock hooks
	mockHooks := &MockLiquidStakingHooks{}
	k.SetHooks(mockHooks)
	
	// Enable module
	params := types.DefaultParams()
	params.Enabled = true
	require.NoError(t, k.SetParams(ctx, params))
	
	// Mock total bonded tokens
	totalBonded := math.NewInt(10000000)
	mockStakingKeeper.TotalBondedTokensFn = func(c sdk.Context) (math.Int, error) {
		return totalBonded, nil
	}
	
	// Set activity at limit
	activity := keeper.TokenizationActivity{
		TotalAmount:   math.NewInt(499000),
		LastActivity:  ctx.BlockTime(),
		ActivityCount: 10,
	}
	k.SetGlobalTokenizationActivity(ctx, activity)
	
	// Try to exceed limit
	err := k.CheckGlobalRateLimit(ctx, math.NewInt(2000))
	require.Error(t, err)
	
	// Verify hook was called
	require.True(t, mockHooks.RateLimitExceededCalled)
	require.Equal(t, "global", mockHooks.LastLimitType)
	require.Equal(t, "global", mockHooks.LastAddress)
	require.Equal(t, math.NewInt(2000), mockHooks.LastRejectedAmount)
}

// TestTypedEventStructures tests the typed event structures
func TestTypedEventStructures(t *testing.T) {
	ctx := sdk.NewContext(nil, sdk.BlockHeight(1), false, nil).
		WithBlockTime(time.Now())
	
	t.Run("RateLimitExceededEvent", func(t *testing.T) {
		event := types.RateLimitExceededEvent{
			LimitType:      "validator",
			Address:        testValAddr1.String(),
			CurrentUsage:   "100000",
			MaxUsage:       "200000",
			RejectedAmount: "50000",
			WindowEnd:      time.Now().Add(24 * time.Hour).Format(time.RFC3339),
		}
		
		sdkEvent := event.ToEvent()
		require.Equal(t, types.EventTypeRateLimitExceeded, sdkEvent.Type)
		require.Len(t, sdkEvent.Attributes, 6)
		
		// Verify attributes
		attrMap := make(map[string]string)
		for _, attr := range sdkEvent.Attributes {
			attrMap[string(attr.Key)] = string(attr.Value)
		}
		
		require.Equal(t, "validator", attrMap[types.AttributeKeyLimitType])
		require.Equal(t, testValAddr1.String(), attrMap[types.AttributeKeyAddress])
		require.Equal(t, "100000", attrMap[types.AttributeKeyCurrentUsage])
		require.Equal(t, "200000", attrMap[types.AttributeKeyMaxUsage])
		require.Equal(t, "50000", attrMap[types.AttributeKeyRejectedAmount])
		require.NotEmpty(t, attrMap[types.AttributeKeyWindowEnd])
	})
	
	t.Run("RateLimitWarningEvent", func(t *testing.T) {
		event := types.RateLimitWarningEvent{
			LimitType:       "user",
			Address:         testAccAddr1.String(),
			CurrentUsage:    "4",
			MaxUsage:        "5",
			PercentageUsed:  "80",
			LimitThreshold:  "80",
		}
		
		sdkEvent := event.ToEvent()
		require.Equal(t, types.EventTypeRateLimitWarning, sdkEvent.Type)
		require.Len(t, sdkEvent.Attributes, 6)
	})
	
	t.Run("ActivityTrackedEvent", func(t *testing.T) {
		event := types.ActivityTrackedEvent{
			LimitType:     "global",
			Address:       "global",
			Amount:        "100000",
			TotalAmount:   "500000",
			ActivityCount: "5",
			WindowStart:   time.Now().Add(-12 * time.Hour).Format(time.RFC3339),
			WindowEnd:     time.Now().Add(12 * time.Hour).Format(time.RFC3339),
		}
		
		sdkEvent := event.ToEvent()
		require.Equal(t, types.EventTypeActivityTracked, sdkEvent.Type)
		require.Len(t, sdkEvent.Attributes, 7)
		
		// Verify attributes
		attrMap := make(map[string]string)
		for _, attr := range sdkEvent.Attributes {
			attrMap[string(attr.Key)] = string(attr.Value)
		}
		
		require.Equal(t, "global", attrMap[types.AttributeKeyLimitType])
		require.Equal(t, "global", attrMap[types.AttributeKeyAddress])
		require.Equal(t, "100000", attrMap[types.AttributeKeyAmount])
		require.Equal(t, "500000", attrMap[types.AttributeKeyCurrentAmount])
		require.Equal(t, "5", attrMap["activity_count"])
		require.NotEmpty(t, attrMap[types.AttributeKeyWindowStart])
		require.NotEmpty(t, attrMap[types.AttributeKeyWindowEnd])
	})
}