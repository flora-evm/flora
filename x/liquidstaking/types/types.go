package types

import (
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Validate performs validation of TokenizationRecord
func (r TokenizationRecord) Validate() error {
	if r.Id == 0 {
		return fmt.Errorf("tokenization record id cannot be zero")
	}

	_, err := sdk.ValAddressFromBech32(r.Validator)
	if err != nil {
		return fmt.Errorf("invalid validator address: %w", err)
	}

	_, err = sdk.AccAddressFromBech32(r.Owner)
	if err != nil {
		return fmt.Errorf("invalid owner address: %w", err)
	}

	if !r.SharesTokenized.IsPositive() {
		return fmt.Errorf("shares tokenized must be positive")
	}

	return nil
}

// NewTokenizationRecord creates a new TokenizationRecord instance
func NewTokenizationRecord(id uint64, validator, owner string, sharesTokenized math.Int) TokenizationRecord {
	return TokenizationRecord{
		Id:              id,
		Validator:       validator,
		Owner:           owner,
		SharesTokenized: sharesTokenized,
	}
}

// DefaultParams returns default module parameters
func DefaultParams() ModuleParams {
	return ModuleParams{
		GlobalLiquidStakingCap: math.LegacyNewDecWithPrec(25, 2), // 25%
		ValidatorLiquidCap:     math.LegacyNewDecWithPrec(50, 2), // 50%
		Enabled:                true,
	}
}

// Validate performs validation of ModuleParams
func (p ModuleParams) Validate() error {
	if p.GlobalLiquidStakingCap.IsNegative() {
		return fmt.Errorf("global liquid staking cap cannot be negative")
	}
	if p.GlobalLiquidStakingCap.GT(math.LegacyOneDec()) {
		return fmt.Errorf("global liquid staking cap cannot exceed 100%%")
	}

	if p.ValidatorLiquidCap.IsNegative() {
		return fmt.Errorf("validator liquid cap cannot be negative")
	}
	if p.ValidatorLiquidCap.GT(math.LegacyOneDec()) {
		return fmt.Errorf("validator liquid cap cannot exceed 100%%")
	}

	if p.GlobalLiquidStakingCap.GT(p.ValidatorLiquidCap) {
		return fmt.Errorf("global liquid staking cap cannot exceed validator liquid cap")
	}

	return nil
}

// NewParams creates a new ModuleParams instance
func NewParams(globalCap, validatorCap math.LegacyDec, enabled bool) ModuleParams {
	return ModuleParams{
		GlobalLiquidStakingCap: globalCap,
		ValidatorLiquidCap:     validatorCap,
		Enabled:                enabled,
	}
}