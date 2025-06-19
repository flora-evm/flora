// +build abi_test

package precompile_test

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/stretchr/testify/require"
)

// TestABIOnly tests just the ABI file without any dependencies
func TestABIOnly(t *testing.T) {
	// Read ABI file directly
	abiData, err := ioutil.ReadFile("abi.json")
	require.NoError(t, err, "Should be able to read abi.json")
	
	// Parse ABI
	var parsedABI abi.ABI
	err = json.Unmarshal(abiData, &parsedABI)
	require.NoError(t, err, "Should be able to parse ABI JSON")
	
	// Check methods exist
	expectedMethods := []string{
		"getParams",
		"getTokenizationRecord",
		"getTokenizationRecords",
		"getRecordsByOwner",
		"getRecordsByValidator",
		"getTotalLiquidStaked",
		"getValidatorLiquidStaked",
		"getLiquidStakingTokenInfo",
		"tokenizeShares",
		"redeemTokens",
	}
	
	for _, method := range expectedMethods {
		_, exists := parsedABI.Methods[method]
		require.True(t, exists, "Method %s should exist in ABI", method)
	}
	
	// Check events exist
	expectedEvents := []string{
		"TokenizeSharesEvent",
		"RedeemTokensEvent",
	}
	
	for _, event := range expectedEvents {
		_, exists := parsedABI.Events[event]
		require.True(t, exists, "Event %s should exist in ABI", event)
	}
	
	t.Log("ABI validation successful!")
}