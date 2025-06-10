package types_test

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	"github.com/rollchains/flora/x/liquidstaking/types"
)

func TestTokenizationRecord_Validate(t *testing.T) {
	validatorAddr := types.TestValidatorAddr
	ownerAddr := types.TestOwnerAddr

	tests := []struct {
		name    string
		record  types.TokenizationRecord
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid record",
			record: types.NewTokenizationRecord(
				1,
				validatorAddr,
				ownerAddr,
				math.NewInt(1000),
			),
			wantErr: false,
		},
		{
			name: "zero id",
			record: types.NewTokenizationRecord(
				0,
				validatorAddr,
				ownerAddr,
				math.NewInt(1000),
			),
			wantErr: true,
			errMsg:  "tokenization record id cannot be zero",
		},
		{
			name: "invalid validator address",
			record: types.NewTokenizationRecord(
				1,
				"invalid",
				ownerAddr,
				math.NewInt(1000),
			),
			wantErr: true,
			errMsg:  "invalid validator address",
		},
		{
			name: "invalid owner address",
			record: types.NewTokenizationRecord(
				1,
				validatorAddr,
				"invalid",
				math.NewInt(1000),
			),
			wantErr: true,
			errMsg:  "invalid owner address",
		},
		{
			name: "zero shares",
			record: types.NewTokenizationRecord(
				1,
				validatorAddr,
				ownerAddr,
				math.ZeroInt(),
			),
			wantErr: true,
			errMsg:  "shares tokenized must be positive",
		},
		{
			name: "negative shares",
			record: types.NewTokenizationRecord(
				1,
				validatorAddr,
				ownerAddr,
				math.NewInt(-1000),
			),
			wantErr: true,
			errMsg:  "shares tokenized must be positive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.record.Validate()
			if tt.wantErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestModuleParams_Validate(t *testing.T) {
	tests := []struct {
		name    string
		params  types.ModuleParams
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid params",
			params:  types.DefaultParams(),
			wantErr: false,
		},
		{
			name: "valid custom params",
			params: types.NewParams(
				math.LegacyNewDecWithPrec(10, 2), // 10%
				math.LegacyNewDecWithPrec(20, 2), // 20%
				true,
			),
			wantErr: false,
		},
		{
			name: "negative global cap",
			params: types.NewParams(
				math.LegacyNewDecWithPrec(-10, 2),
				math.LegacyNewDecWithPrec(20, 2),
				true,
			),
			wantErr: true,
			errMsg:  "global liquid staking cap cannot be negative",
		},
		{
			name: "global cap over 100%",
			params: types.NewParams(
				math.LegacyNewDecWithPrec(101, 2),
				math.LegacyNewDecWithPrec(20, 2),
				true,
			),
			wantErr: true,
			errMsg:  "global liquid staking cap cannot exceed 100%",
		},
		{
			name: "negative validator cap",
			params: types.NewParams(
				math.LegacyNewDecWithPrec(10, 2),
				math.LegacyNewDecWithPrec(-20, 2),
				true,
			),
			wantErr: true,
			errMsg:  "validator liquid cap cannot be negative",
		},
		{
			name: "validator cap over 100%",
			params: types.NewParams(
				math.LegacyNewDecWithPrec(10, 2),
				math.LegacyNewDecWithPrec(101, 2),
				true,
			),
			wantErr: true,
			errMsg:  "validator liquid cap cannot exceed 100%",
		},
		{
			name: "global cap exceeds validator cap",
			params: types.NewParams(
				math.LegacyNewDecWithPrec(50, 2), // 50%
				math.LegacyNewDecWithPrec(40, 2), // 40%
				true,
			),
			wantErr: true,
			errMsg:  "global liquid staking cap cannot exceed validator liquid cap",
		},
		{
			name: "disabled module with valid caps",
			params: types.NewParams(
				math.LegacyNewDecWithPrec(25, 2),
				math.LegacyNewDecWithPrec(50, 2),
				false,
			),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.params.Validate()
			if tt.wantErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestDefaultParams(t *testing.T) {
	params := types.DefaultParams()
	require.Equal(t, math.LegacyNewDecWithPrec(25, 2), params.GlobalLiquidStakingCap)
	require.Equal(t, math.LegacyNewDecWithPrec(50, 2), params.ValidatorLiquidCap)
	require.True(t, params.Enabled)
	require.NoError(t, params.Validate())
}