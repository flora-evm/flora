package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/rollchains/flora/x/liquidstaking/types"
)

func TestGenerateLiquidStakingTokenDenom(t *testing.T) {
	testCases := []struct {
		name          string
		validatorAddr string
		recordID      uint64
		expectedDenom string
	}{
		{
			name:          "valid generation",
			validatorAddr: "floravaloper1abcdef123456",
			recordID:      1,
			expectedDenom: "flora/lstake/floravaloper1abcdef123456/1",
		},
		{
			name:          "another valid generation",
			validatorAddr: "floravaloper1xyz789",
			recordID:      100,
			expectedDenom: "flora/lstake/floravaloper1xyz789/100",
		},
		{
			name:          "large record ID",
			validatorAddr: "floravaloper1test",
			recordID:      999999999,
			expectedDenom: "flora/lstake/floravaloper1test/999999999",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			denom := types.GenerateLiquidStakingTokenDenom(tc.validatorAddr, tc.recordID)
			require.Equal(t, tc.expectedDenom, denom)
		})
	}
}

func TestParseLiquidStakingTokenDenom(t *testing.T) {
	testCases := []struct {
		name              string
		denom             string
		expectedValidator string
		expectedRecordID  uint64
		expectError       bool
	}{
		{
			name:              "valid parsing",
			denom:             "flora/lstake/floravaloper1abcdef123456/1",
			expectedValidator: "floravaloper1abcdef123456",
			expectedRecordID:  1,
			expectError:       false,
		},
		{
			name:              "valid parsing with large ID",
			denom:             "flora/lstake/floravaloper1xyz789/999999999",
			expectedValidator: "floravaloper1xyz789",
			expectedRecordID:  999999999,
			expectError:       false,
		},
		{
			name:              "invalid prefix",
			denom:             "cosmos/lstake/floravaloper1test/1",
			expectedValidator: "",
			expectedRecordID:  0,
			expectError:       true,
		},
		{
			name:              "missing record ID",
			denom:             "flora/lstake/floravaloper1test/",
			expectedValidator: "",
			expectedRecordID:  0,
			expectError:       true,
		},
		{
			name:              "invalid record ID",
			denom:             "flora/lstake/floravaloper1test/abc",
			expectedValidator: "",
			expectedRecordID:  0,
			expectError:       true,
		},
		{
			name:              "missing slash",
			denom:             "flora/lstake/floravaloper1test123",
			expectedValidator: "",
			expectedRecordID:  0,
			expectError:       true,
		},
		{
			name:              "empty denom",
			denom:             "",
			expectedValidator: "",
			expectedRecordID:  0,
			expectError:       true,
		},
		{
			name:              "regular denom",
			denom:             "flora",
			expectedValidator: "",
			expectedRecordID:  0,
			expectError:       true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			validatorAddr, recordID, err := types.ParseLiquidStakingTokenDenom(tc.denom)

			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedValidator, validatorAddr)
				require.Equal(t, tc.expectedRecordID, recordID)
			}
		})
	}
}

func TestIsLiquidStakingTokenDenom(t *testing.T) {
	testCases := []struct {
		name     string
		denom    string
		expected bool
	}{
		{
			name:     "valid liquid staking token",
			denom:    "flora/lstake/floravaloper1abcdef123456/1",
			expected: true,
		},
		{
			name:     "valid liquid staking token with large ID",
			denom:    "flora/lstake/floravaloper1xyz789/999999999",
			expected: true,
		},
		{
			name:     "regular flora denom",
			denom:    "flora",
			expected: false,
		},
		{
			name:     "ibc denom",
			denom:    "ibc/ABCDEF123456",
			expected: false,
		},
		{
			name:     "factory denom",
			denom:    "factory/flora1abc/token",
			expected: false,
		},
		{
			name:     "empty denom",
			denom:    "",
			expected: false,
		},
		{
			name:     "partial prefix",
			denom:    "flora/lstake",
			expected: false, // Needs complete format to be valid
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := types.IsLiquidStakingTokenDenom(tc.denom)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestGetValidatorFromLiquidStakingTokenDenom(t *testing.T) {
	testCases := []struct {
		name              string
		denom             string
		expectedValidator string
		expectError       bool
	}{
		{
			name:              "valid extraction",
			denom:             "flora/lstake/floravaloper1abcdef123456/1",
			expectedValidator: "floravaloper1abcdef123456",
			expectError:       false,
		},
		{
			name:              "invalid denom",
			denom:             "flora",
			expectedValidator: "",
			expectError:       true,
		},
		{
			name:              "invalid format",
			denom:             "flora/lstake/invalid",
			expectedValidator: "",
			expectError:       true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			validatorAddr, err := types.GetValidatorFromLiquidStakingTokenDenom(tc.denom)

			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedValidator, validatorAddr)
			}
		})
	}
}

func TestGetRecordIDFromLiquidStakingTokenDenom(t *testing.T) {
	testCases := []struct {
		name             string
		denom            string
		expectedRecordID uint64
		expectError      bool
	}{
		{
			name:             "valid extraction",
			denom:            "flora/lstake/floravaloper1abcdef123456/1",
			expectedRecordID: 1,
			expectError:      false,
		},
		{
			name:             "valid extraction large ID",
			denom:            "flora/lstake/floravaloper1xyz789/999999999",
			expectedRecordID: 999999999,
			expectError:      false,
		},
		{
			name:             "invalid denom",
			denom:            "flora",
			expectedRecordID: 0,
			expectError:      true,
		},
		{
			name:             "invalid format",
			denom:            "flora/lstake/floravaloper1test/abc",
			expectedRecordID: 0,
			expectError:      true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			recordID, err := types.GetRecordIDFromLiquidStakingTokenDenom(tc.denom)

			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedRecordID, recordID)
			}
		})
	}
}

func TestRoundTripDenomOperations(t *testing.T) {
	// Test that generate -> parse -> generate produces the same result
	validatorAddr := "floravaloper1roundtrip123"
	recordID := uint64(42)

	// Generate denom
	denom := types.GenerateLiquidStakingTokenDenom(validatorAddr, recordID)
	require.True(t, types.IsLiquidStakingTokenDenom(denom))

	// Parse it back
	parsedValidator, parsedRecordID, err := types.ParseLiquidStakingTokenDenom(denom)
	require.NoError(t, err)
	require.Equal(t, validatorAddr, parsedValidator)
	require.Equal(t, recordID, parsedRecordID)

	// Use helper functions
	validatorFromHelper, err := types.GetValidatorFromLiquidStakingTokenDenom(denom)
	require.NoError(t, err)
	require.Equal(t, validatorAddr, validatorFromHelper)

	recordIDFromHelper, err := types.GetRecordIDFromLiquidStakingTokenDenom(denom)
	require.NoError(t, err)
	require.Equal(t, recordID, recordIDFromHelper)

	// Generate again with parsed values
	regeneratedDenom := types.GenerateLiquidStakingTokenDenom(parsedValidator, parsedRecordID)
	require.Equal(t, denom, regeneratedDenom)
}