package types_test

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	"github.com/rollchains/flora/x/liquidstaking/types"
)

func TestGenesisState_Validate(t *testing.T) {
	validatorAddr := types.TestValidatorAddr
	ownerAddr := types.TestOwnerAddr

	tests := []struct {
		name    string
		genesis *types.GenesisState
		wantErr bool
		errMsg  string
	}{
		{
			name:    "default genesis is valid",
			genesis: types.DefaultGenesisState(),
			wantErr: false,
		},
		{
			name: "valid genesis with records",
			genesis: types.NewGenesisState(
				types.DefaultParams(),
				[]types.TokenizationRecord{
					types.NewTokenizationRecord(1, validatorAddr, ownerAddr, math.NewInt(1000)),
					types.NewTokenizationRecord(2, validatorAddr, ownerAddr, math.NewInt(2000)),
				},
				2,
			),
			wantErr: false,
		},
		{
			name: "empty genesis",
			genesis: &types.GenesisState{
				Params:                   types.DefaultParams(),
				TokenizationRecords:      []types.TokenizationRecord{},
				LastTokenizationRecordId: 0,
			},
			wantErr: false,
		},
		{
			name: "invalid params",
			genesis: &types.GenesisState{
				Params: types.NewParams(
					math.LegacyNewDecWithPrec(-10, 2), // negative cap
					math.LegacyNewDecWithPrec(50, 2),
					true,
				),
				TokenizationRecords:      []types.TokenizationRecord{},
				LastTokenizationRecordId: 0,
			},
			wantErr: true,
			errMsg:  "invalid params",
		},
		{
			name: "invalid tokenization record",
			genesis: types.NewGenesisState(
				types.DefaultParams(),
				[]types.TokenizationRecord{
					types.NewTokenizationRecord(0, validatorAddr, ownerAddr, math.NewInt(1000)), // zero ID
				},
				1,
			),
			wantErr: true,
			errMsg:  "invalid tokenization record at index 0",
		},
		{
			name: "duplicate record IDs",
			genesis: types.NewGenesisState(
				types.DefaultParams(),
				[]types.TokenizationRecord{
					types.NewTokenizationRecord(1, validatorAddr, ownerAddr, math.NewInt(1000)),
					types.NewTokenizationRecord(1, validatorAddr, ownerAddr, math.NewInt(2000)), // duplicate ID
				},
				2,
			),
			wantErr: true,
			errMsg:  "duplicate tokenization record ID: 1",
		},
		{
			name: "last record ID less than max",
			genesis: types.NewGenesisState(
				types.DefaultParams(),
				[]types.TokenizationRecord{
					types.NewTokenizationRecord(1, validatorAddr, ownerAddr, math.NewInt(1000)),
					types.NewTokenizationRecord(5, validatorAddr, ownerAddr, math.NewInt(2000)),
				},
				3, // less than max ID (5)
			),
			wantErr: true,
			errMsg:  "last tokenization record ID (3) is less than maximum record ID (5)",
		},
		{
			name: "valid with non-sequential IDs",
			genesis: types.NewGenesisState(
				types.DefaultParams(),
				[]types.TokenizationRecord{
					types.NewTokenizationRecord(1, validatorAddr, ownerAddr, math.NewInt(1000)),
					types.NewTokenizationRecord(3, validatorAddr, ownerAddr, math.NewInt(2000)),
					types.NewTokenizationRecord(7, validatorAddr, ownerAddr, math.NewInt(3000)),
				},
				10, // greater than max ID
			),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.genesis.Validate()
			if tt.wantErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestDefaultGenesisState(t *testing.T) {
	genesis := types.DefaultGenesisState()
	
	// Check default values
	require.Equal(t, types.DefaultParams(), genesis.Params)
	require.Empty(t, genesis.TokenizationRecords)
	require.Equal(t, uint64(0), genesis.LastTokenizationRecordId)
	
	// Default genesis should be valid
	require.NoError(t, genesis.Validate())
}