package types_test

import (
	"testing"
	
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	
	"github.com/rollchains/flora/x/liquidstaking/types"
)

// Stage 1 Tests: Pure type validation with no external dependencies

func TestTokenizationRecord_Validate(t *testing.T) {
	// Valid addresses for testing
	validValidator := "cosmosvaloper1tnh2q55v8wyygtt9srz5safamzdengsn9dsd7"
	validOwner := "cosmos1tnh2q55v8wyygtt9srz5safamzdengsnqeycj"
	
	testCases := []struct {
		name      string
		record    types.TokenizationRecord
		expectErr bool
		errMsg    string
	}{
		{
			name: "valid record",
			record: types.TokenizationRecord{
				Id:              1,
				Validator:       validValidator,
				Owner:           validOwner,
				SharesTokenized: sdk.NewInt(1000),
				Status:          types.TOKENIZATION_STATUS_TOKENIZED,
			},
			expectErr: false,
		},
		{
			name: "zero id",
			record: types.TokenizationRecord{
				Id:              0,
				Validator:       validValidator,
				Owner:           validOwner,
				SharesTokenized: sdk.NewInt(1000),
				Status:          types.TOKENIZATION_STATUS_TOKENIZED,
			},
			expectErr: true,
			errMsg:    "id cannot be zero",
		},
		{
			name: "invalid validator address",
			record: types.TokenizationRecord{
				Id:              1,
				Validator:       "invalid",
				Owner:           validOwner,
				SharesTokenized: sdk.NewInt(1000),
				Status:          types.TOKENIZATION_STATUS_TOKENIZED,
			},
			expectErr: true,
			errMsg:    "invalid validator address",
		},
		{
			name: "invalid owner address",
			record: types.TokenizationRecord{
				Id:              1,
				Validator:       validValidator,
				Owner:           "invalid",
				SharesTokenized: sdk.NewInt(1000),
				Status:          types.TOKENIZATION_STATUS_TOKENIZED,
			},
			expectErr: true,
			errMsg:    "invalid owner address",
		},
		{
			name: "zero shares",
			record: types.TokenizationRecord{
				Id:              1,
				Validator:       validValidator,
				Owner:           validOwner,
				SharesTokenized: sdk.NewInt(0),
				Status:          types.TOKENIZATION_STATUS_TOKENIZED,
			},
			expectErr: true,
			errMsg:    "shares tokenized must be positive",
		},
		{
			name: "negative shares",
			record: types.TokenizationRecord{
				Id:              1,
				Validator:       validValidator,
				Owner:           validOwner,
				SharesTokenized: sdk.NewInt(-1000),
				Status:          types.TOKENIZATION_STATUS_TOKENIZED,
			},
			expectErr: true,
			errMsg:    "shares tokenized must be positive",
		},
		{
			name: "unspecified status",
			record: types.TokenizationRecord{
				Id:              1,
				Validator:       validValidator,
				Owner:           validOwner,
				SharesTokenized: sdk.NewInt(1000),
				Status:          types.TOKENIZATION_STATUS_UNSPECIFIED,
			},
			expectErr: true,
			errMsg:    "status cannot be unspecified",
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.record.Validate()
			
			if tc.expectErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestModuleParams_Validate(t *testing.T) {
	testCases := []struct {
		name      string
		params    types.ModuleParams
		expectErr bool
		errMsg    string
	}{
		{
			name:      "default params",
			params:    types.DefaultParams(),
			expectErr: false,
		},
		{
			name: "valid custom params",
			params: types.ModuleParams{
				Enabled:                true,
				GlobalLiquidStakingCap: sdk.NewDecWithPrec(30, 2),
				ValidatorLiquidCap:     sdk.NewDecWithPrec(20, 2),
				MinTokenizationAmount:  sdk.NewInt(5000000),
			},
			expectErr: false,
		},
		{
			name: "negative global cap",
			params: types.ModuleParams{
				Enabled:                true,
				GlobalLiquidStakingCap: sdk.NewDecWithPrec(-10, 2),
				ValidatorLiquidCap:     sdk.NewDecWithPrec(20, 2),
				MinTokenizationAmount:  sdk.NewInt(1000000),
			},
			expectErr: true,
			errMsg:    "global liquid staking cap must be between 0 and 1",
		},
		{
			name: "global cap > 1",
			params: types.ModuleParams{
				Enabled:                true,
				GlobalLiquidStakingCap: sdk.NewDecWithPrec(110, 2),
				ValidatorLiquidCap:     sdk.NewDecWithPrec(20, 2),
				MinTokenizationAmount:  sdk.NewInt(1000000),
			},
			expectErr: true,
			errMsg:    "global liquid staking cap must be between 0 and 1",
		},
		{
			name: "validator cap > global cap",
			params: types.ModuleParams{
				Enabled:                true,
				GlobalLiquidStakingCap: sdk.NewDecWithPrec(20, 2),
				ValidatorLiquidCap:     sdk.NewDecWithPrec(30, 2),
				MinTokenizationAmount:  sdk.NewInt(1000000),
			},
			expectErr: true,
			errMsg:    "validator cap cannot exceed global cap",
		},
		{
			name: "zero min amount",
			params: types.ModuleParams{
				Enabled:                true,
				GlobalLiquidStakingCap: sdk.NewDecWithPrec(25, 2),
				ValidatorLiquidCap:     sdk.NewDecWithPrec(20, 2),
				MinTokenizationAmount:  sdk.NewInt(0),
			},
			expectErr: true,
			errMsg:    "minimum tokenization amount must be positive",
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.params.Validate()
			
			if tc.expectErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGenesisState_Validate(t *testing.T) {
	validValidator := "cosmosvaloper1tnh2q55v8wyygtt9srz5safamzdengsn9dsd7"
	validOwner := "cosmos1tnh2q55v8wyygtt9srz5safamzdengsnqeycj"
	
	validRecord := types.TokenizationRecord{
		Id:              1,
		Validator:       validValidator,
		Owner:           validOwner,
		SharesTokenized: sdk.NewInt(1000),
		Status:          types.TOKENIZATION_STATUS_TOKENIZED,
	}
	
	testCases := []struct {
		name      string
		genesis   types.GenesisState
		expectErr bool
		errMsg    string
	}{
		{
			name:      "default genesis",
			genesis:   *types.DefaultGenesisState(),
			expectErr: false,
		},
		{
			name: "valid with records",
			genesis: types.GenesisState{
				Params: types.DefaultParams(),
				TokenizationRecords: []types.TokenizationRecord{
					validRecord,
					{
						Id:              2,
						Validator:       validValidator,
						Owner:           validOwner,
						SharesTokenized: sdk.NewInt(2000),
						Status:          types.TOKENIZATION_STATUS_TOKENIZED,
					},
				},
				LastTokenizationRecordId: 2,
			},
			expectErr: false,
		},
		{
			name: "invalid params",
			genesis: types.GenesisState{
				Params: types.ModuleParams{
					Enabled:                true,
					GlobalLiquidStakingCap: sdk.NewDecWithPrec(-10, 2),
					ValidatorLiquidCap:     sdk.NewDecWithPrec(20, 2),
					MinTokenizationAmount:  sdk.NewInt(1000000),
				},
				TokenizationRecords:      []types.TokenizationRecord{},
				LastTokenizationRecordId: 0,
			},
			expectErr: true,
			errMsg:    "invalid params",
		},
		{
			name: "duplicate record ids",
			genesis: types.GenesisState{
				Params: types.DefaultParams(),
				TokenizationRecords: []types.TokenizationRecord{
					validRecord,
					validRecord, // Same ID
				},
				LastTokenizationRecordId: 1,
			},
			expectErr: true,
			errMsg:    "duplicate tokenization record id: 1",
		},
		{
			name: "inconsistent last id",
			genesis: types.GenesisState{
				Params: types.DefaultParams(),
				TokenizationRecords: []types.TokenizationRecord{
					{
						Id:              5,
						Validator:       validValidator,
						Owner:           validOwner,
						SharesTokenized: sdk.NewInt(1000),
						Status:          types.TOKENIZATION_STATUS_TOKENIZED,
					},
				},
				LastTokenizationRecordId: 3, // Less than max ID (5)
			},
			expectErr: true,
			errMsg:    "last tokenization record id (3) is less than max record id (5)",
		},
		{
			name: "invalid record",
			genesis: types.GenesisState{
				Params: types.DefaultParams(),
				TokenizationRecords: []types.TokenizationRecord{
					{
						Id:              0, // Invalid
						Validator:       validValidator,
						Owner:           validOwner,
						SharesTokenized: sdk.NewInt(1000),
						Status:          types.TOKENIZATION_STATUS_TOKENIZED,
					},
				},
				LastTokenizationRecordId: 0,
			},
			expectErr: true,
			errMsg:    "invalid tokenization record 0",
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.genesis.Validate()
			
			if tc.expectErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// Benchmark tests for type validation
func BenchmarkTokenizationRecord_Validate(b *testing.B) {
	record := types.TokenizationRecord{
		Id:              1,
		Validator:       "cosmosvaloper1tnh2q55v8wyygtt9srz5safamzdengsn9dsd7",
		Owner:           "cosmos1tnh2q55v8wyygtt9srz5safamzdengsnqeycj",
		SharesTokenized: sdk.NewInt(1000000),
		Status:          types.TOKENIZATION_STATUS_TOKENIZED,
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = record.Validate()
	}
}

func BenchmarkGenesisState_Validate(b *testing.B) {
	// Create genesis with 100 records
	records := make([]types.TokenizationRecord, 100)
	for i := 0; i < 100; i++ {
		records[i] = types.TokenizationRecord{
			Id:              uint64(i + 1),
			Validator:       "cosmosvaloper1tnh2q55v8wyygtt9srz5safamzdengsn9dsd7",
			Owner:           "cosmos1tnh2q55v8wyygtt9srz5safamzdengsnqeycj",
			SharesTokenized: sdk.NewInt(1000000),
			Status:          types.TOKENIZATION_STATUS_TOKENIZED,
		}
	}
	
	genesis := types.GenesisState{
		Params:                   types.DefaultParams(),
		TokenizationRecords:      records,
		LastTokenizationRecordId: 100,
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = genesis.Validate()
	}
}