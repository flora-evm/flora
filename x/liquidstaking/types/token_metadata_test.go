package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/rollchains/flora/x/liquidstaking/types"
)

func TestGenerateLiquidStakingTokenMetadata(t *testing.T) {
	testCases := []struct {
		name          string
		validatorAddr string
		recordID      uint64
		validate      func(t *testing.T, metadata banktypes.Metadata)
	}{
		{
			name:          "standard validator address",
			validatorAddr: "floravaloper1abcdef123456789",
			recordID:      1,
			validate: func(t *testing.T, metadata banktypes.Metadata) {
				expectedDenom := "flora/lstake/floravaloper1abcdef123456789/1"
				require.Equal(t, expectedDenom, metadata.Base)
				require.Equal(t, "LSTFLORAabcd", metadata.Display)
				require.Equal(t, "Liquid Staked FLORA abcd", metadata.Name)
				require.Equal(t, "LSTFLORAabcd", metadata.Symbol)
				require.Contains(t, metadata.Description, "floravaloper1abcdef123456789")
				
				// Check denom units
				require.Len(t, metadata.DenomUnits, 2)
				require.Equal(t, expectedDenom, metadata.DenomUnits[0].Denom)
				require.Equal(t, uint32(0), metadata.DenomUnits[0].Exponent)
				require.Equal(t, "LSTFLORAabcd", metadata.DenomUnits[1].Denom)
				require.Equal(t, uint32(18), metadata.DenomUnits[1].Exponent)
			},
		},
		{
			name:          "short validator address",
			validatorAddr: "floravaloper1xyz",
			recordID:      100,
			validate: func(t *testing.T, metadata banktypes.Metadata) {
				expectedDenom := "flora/lstake/floravaloper1xyz/100"
				require.Equal(t, expectedDenom, metadata.Base)
				require.Equal(t, "LSTFLORAxyz", metadata.Display) // Takes what's available
				require.Equal(t, "Liquid Staked FLORA xyz", metadata.Name)
				require.Equal(t, "LSTFLORAxyz", metadata.Symbol)
			},
		},
		{
			name:          "non-standard validator prefix",
			validatorAddr: "customvaloper1test123",
			recordID:      42,
			validate: func(t *testing.T, metadata banktypes.Metadata) {
				expectedDenom := "flora/lstake/customvaloper1test123/42"
				require.Equal(t, expectedDenom, metadata.Base)
				require.Equal(t, "LSTFLORA", metadata.Display) // No suffix extracted
				require.Equal(t, "Liquid Staked FLORA ", metadata.Name)
				require.Equal(t, "LSTFLORA", metadata.Symbol)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			metadata := types.GenerateLiquidStakingTokenMetadata(tc.validatorAddr, tc.recordID)
			tc.validate(t, metadata)
		})
	}
}

func TestValidateLiquidStakingTokenMetadata(t *testing.T) {
	validDenom := "flora/lstake/floravaloper1test/1"
	
	testCases := []struct {
		name        string
		metadata    banktypes.Metadata
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid metadata",
			metadata: banktypes.Metadata{
				Base:    validDenom,
				Display: "LSTFLORA",
				Name:    "Liquid Staked FLORA",
				Symbol:  "LSTFLORA",
				DenomUnits: []*banktypes.DenomUnit{
					{Denom: validDenom, Exponent: 0},
					{Denom: "LSTFLORA", Exponent: 18},
				},
			},
			expectError: false,
		},
		{
			name: "invalid base denom",
			metadata: banktypes.Metadata{
				Base:    "flora",
				Display: "LSTFLORA",
				Name:    "Liquid Staked FLORA",
				Symbol:  "LSTFLORA",
				DenomUnits: []*banktypes.DenomUnit{
					{Denom: "flora", Exponent: 0},
				},
			},
			expectError: true,
			errorMsg:    "not a liquid staking token",
		},
		{
			name: "no denom units",
			metadata: banktypes.Metadata{
				Base:       validDenom,
				Display:    "LSTFLORA",
				Name:       "Liquid Staked FLORA",
				Symbol:     "LSTFLORA",
				DenomUnits: []*banktypes.DenomUnit{},
			},
			expectError: true,
			errorMsg:    "must have at least one denom unit",
		},
		{
			name: "first denom unit mismatch",
			metadata: banktypes.Metadata{
				Base:    validDenom,
				Display: "LSTFLORA",
				Name:    "Liquid Staked FLORA",
				Symbol:  "LSTFLORA",
				DenomUnits: []*banktypes.DenomUnit{
					{Denom: "wrong", Exponent: 0},
				},
			},
			expectError: true,
			errorMsg:    "first denom unit must match base denom",
		},
		{
			name: "empty display",
			metadata: banktypes.Metadata{
				Base:    validDenom,
				Display: "",
				Name:    "Liquid Staked FLORA",
				Symbol:  "LSTFLORA",
				DenomUnits: []*banktypes.DenomUnit{
					{Denom: validDenom, Exponent: 0},
				},
			},
			expectError: true,
			errorMsg:    "display denom cannot be empty",
		},
		{
			name: "empty name",
			metadata: banktypes.Metadata{
				Base:    validDenom,
				Display: "LSTFLORA",
				Name:    "",
				Symbol:  "LSTFLORA",
				DenomUnits: []*banktypes.DenomUnit{
					{Denom: validDenom, Exponent: 0},
				},
			},
			expectError: true,
			errorMsg:    "name cannot be empty",
		},
		{
			name: "empty symbol",
			metadata: banktypes.Metadata{
				Base:    validDenom,
				Display: "LSTFLORA",
				Name:    "Liquid Staked FLORA",
				Symbol:  "",
				DenomUnits: []*banktypes.DenomUnit{
					{Denom: validDenom, Exponent: 0},
				},
			},
			expectError: true,
			errorMsg:    "symbol cannot be empty",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := types.ValidateLiquidStakingTokenMetadata(tc.metadata)

			if tc.expectError {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.errorMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestMetadataIntegration(t *testing.T) {
	// Test that generated metadata passes validation
	validatorAddr := "floravaloper1integration123"
	recordID := uint64(999)
	
	metadata := types.GenerateLiquidStakingTokenMetadata(validatorAddr, recordID)
	err := types.ValidateLiquidStakingTokenMetadata(metadata)
	require.NoError(t, err)
	
	// Verify the denom in metadata matches what we expect
	expectedDenom := types.GenerateLiquidStakingTokenDenom(validatorAddr, recordID)
	require.Equal(t, expectedDenom, metadata.Base)
	
	// Verify we can parse the denom back
	parsedValidator, parsedRecordID, err := types.ParseLiquidStakingTokenDenom(metadata.Base)
	require.NoError(t, err)
	require.Equal(t, validatorAddr, parsedValidator)
	require.Equal(t, recordID, parsedRecordID)
}