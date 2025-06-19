package precompile_test

import (
	"encoding/hex"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/stretchr/testify/require"

	"github.com/rollchains/flora/x/liquidstaking/precompile"
)

// TestMinimalPrecompile tests basic precompile functionality without full keeper dependencies
func TestMinimalPrecompile(t *testing.T) {
	// Test ABI loading
	err := precompile.LoadABI()
	require.NoError(t, err)
	require.NotNil(t, precompile.ABI)

	// Test method IDs
	methods := []string{
		"getParams",
		"getTokenizationRecord",
		"getTotalLiquidStaked",
		"tokenizeShares",
		"redeemTokens",
	}

	for _, methodName := range methods {
		method, exists := precompile.ABI.Methods[methodName]
		require.True(t, exists, "Method %s should exist", methodName)
		require.Equal(t, 4, len(method.ID), "Method ID should be 4 bytes for %s", methodName)
		
		// Verify it's valid hex
		hexStr := hex.EncodeToString(method.ID)
		require.Equal(t, 8, len(hexStr), "Hex string should be 8 characters for %s", methodName)
	}

	// Test gas requirements
	testCases := []struct {
		method      string
		expectedGas uint64
	}{
		{"getParams", precompile.GasBaseQuery},
		{"getTokenizationRecord", precompile.GasGetRecord},
		{"getTotalLiquidStaked", precompile.GasBaseQuery},
		{"tokenizeShares", precompile.GasTokenizeShares},
		{"redeemTokens", precompile.GasRedeemTokens},
	}

	// Create a minimal contract instance (without keeper)
	contract := &precompile.Contract{}

	for _, tc := range testCases {
		methodID := precompile.ABI.Methods[tc.method].ID
		gas := contract.RequiredGas(methodID)
		require.Equal(t, tc.expectedGas, gas, "Gas mismatch for method %s", tc.method)
	}

	// Test invalid method ID
	invalidID := []byte{0x00, 0x00, 0x00, 0x00}
	gas := contract.RequiredGas(invalidID)
	require.Equal(t, uint64(0), gas, "Invalid method ID should return 0 gas")
}

// TestABIStructures tests the ABI structure definitions
func TestABIStructures(t *testing.T) {
	// Test that the ABI contains expected events
	events := []string{
		"TokenizeSharesEvent",
		"RedeemTokensEvent",
	}

	for _, eventName := range events {
		event, exists := precompile.ABI.Events[eventName]
		require.True(t, exists, "Event %s should exist", eventName)
		require.NotEmpty(t, event.ID, "Event %s should have an ID", eventName)
	}
}