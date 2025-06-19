package types_test

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	"github.com/rollchains/flora/x/liquidstaking/types"
)

func TestGenesisState_Validate(t *testing.T) {
	validatorAddr := types.TestValidatorAddr
	validatorAddr2 := "floravaloper1validator2"
	ownerAddr := types.TestOwnerAddr
	ownerAddr2 := "flora1owner2"

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
			name: "valid genesis with records and denoms",
			genesis: &types.GenesisState{
				Params: types.DefaultParams(),
				TokenizationRecords: []types.TokenizationRecord{
					{
						Id:              1,
						Validator:       validatorAddr,
						Owner:           ownerAddr,
						SharesTokenized: math.NewInt(1000),
						Denom:           types.GenerateLiquidStakingTokenDenom(validatorAddr, 1),
					},
					{
						Id:              2,
						Validator:       validatorAddr2,
						Owner:           ownerAddr2,
						SharesTokenized: math.NewInt(2000),
						Denom:           types.GenerateLiquidStakingTokenDenom(validatorAddr2, 2),
					},
				},
				LastTokenizationRecordId: 2,
			},
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
				Params: types.ModuleParams{
					Enabled:                true,
					MinLiquidStakeAmount:   math.NewInt(-1), // negative amount
					GlobalLiquidStakingCap: math.LegacyNewDec(1),
					ValidatorLiquidCap:     math.LegacyNewDec(1),
				},
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
			name: "duplicate denoms",
			genesis: &types.GenesisState{
				Params: types.DefaultParams(),
				TokenizationRecords: []types.TokenizationRecord{
					{
						Id:              1,
						Validator:       validatorAddr,
						Owner:           ownerAddr,
						SharesTokenized: math.NewInt(1000),
						Denom:           types.GenerateLiquidStakingTokenDenom(validatorAddr, 1),
					},
					{
						Id:              2,
						Validator:       validatorAddr,
						Owner:           ownerAddr2,
						SharesTokenized: math.NewInt(2000),
						Denom:           types.GenerateLiquidStakingTokenDenom(validatorAddr, 1), // duplicate denom
					},
				},
				LastTokenizationRecordId: 2,
			},
			wantErr: true,
			errMsg:  "duplicate denom",
		},
		{
			name: "invalid denom format",
			genesis: &types.GenesisState{
				Params: types.DefaultParams(),
				TokenizationRecords: []types.TokenizationRecord{
					{
						Id:              1,
						Validator:       validatorAddr,
						Owner:           ownerAddr,
						SharesTokenized: math.NewInt(1000),
						Denom:           "invalid-denom-format",
					},
				},
				LastTokenizationRecordId: 1,
			},
			wantErr: true,
			errMsg:  "invalid liquid staking token denom format",
		},
		{
			name: "denom validator mismatch",
			genesis: &types.GenesisState{
				Params: types.DefaultParams(),
				TokenizationRecords: []types.TokenizationRecord{
					{
						Id:              1,
						Validator:       validatorAddr,
						Owner:           ownerAddr,
						SharesTokenized: math.NewInt(1000),
						Denom:           types.GenerateLiquidStakingTokenDenom(validatorAddr2, 1), // wrong validator
					},
				},
				LastTokenizationRecordId: 1,
			},
			wantErr: true,
			errMsg:  "does not match expected format",
		},
		{
			name: "denom record ID mismatch",
			genesis: &types.GenesisState{
				Params: types.DefaultParams(),
				TokenizationRecords: []types.TokenizationRecord{
					{
						Id:              1,
						Validator:       validatorAddr,
						Owner:           ownerAddr,
						SharesTokenized: math.NewInt(1000),
						Denom:           types.GenerateLiquidStakingTokenDenom(validatorAddr, 2), // wrong record ID
					},
				},
				LastTokenizationRecordId: 1,
			},
			wantErr: true,
			errMsg:  "does not match expected format",
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
			name: "non-sequential record IDs",
			genesis: &types.GenesisState{
				Params: types.DefaultParams(),
				TokenizationRecords: []types.TokenizationRecord{
					{
						Id:              1,
						Validator:       validatorAddr,
						Owner:           ownerAddr,
						SharesTokenized: math.NewInt(1000),
						Denom:           types.GenerateLiquidStakingTokenDenom(validatorAddr, 1),
					},
					{
						Id:              3, // skipped ID 2
						Validator:       validatorAddr2,
						Owner:           ownerAddr2,
						SharesTokenized: math.NewInt(2000),
						Denom:           types.GenerateLiquidStakingTokenDenom(validatorAddr2, 3),
					},
				},
				LastTokenizationRecordId: 3,
			},
			wantErr: true,
			errMsg:  "missing tokenization record ID 2",
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