package types

import (
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetLSTDenom returns the liquid staking token denomination for a validator
func GetLSTDenom(validatorAddr string) string {
	return fmt.Sprintf("lst/%s", validatorAddr)
}

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

// NewTokenizationRecordWithDenom creates a new TokenizationRecord instance with denom
func NewTokenizationRecordWithDenom(id uint64, validator, owner string, sharesTokenized math.Int, denom string) TokenizationRecord {
	return TokenizationRecord{
		Id:              id,
		Validator:       validator,
		Owner:           owner,
		SharesTokenized: sharesTokenized,
		Denom:           denom,
	}
}

// DefaultParams returns default module parameters
func DefaultParams() ModuleParams {
	return ModuleParams{
		GlobalLiquidStakingCap: math.LegacyNewDecWithPrec(25, 2), // 25%
		ValidatorLiquidCap:     math.LegacyNewDecWithPrec(50, 2), // 50%
		Enabled:                true,
		MinLiquidStakeAmount:   math.NewInt(10000), // Minimum 10,000 units
		// Rate limiting parameters
		RateLimitPeriodHours:              24,                              // 24 hours
		GlobalDailyTokenizationPercent:    math.LegacyNewDecWithPrec(5, 2), // 5%
		ValidatorDailyTokenizationPercent: math.LegacyNewDecWithPrec(10, 2), // 10%
		GlobalDailyTokenizationCount:      100, // 100 tokenizations per day globally
		ValidatorDailyTokenizationCount:   20,  // 20 tokenizations per validator per day
		UserDailyTokenizationCount:        5,   // 5 tokenizations per user per day
		WarningThresholdPercent:           math.LegacyNewDecWithPrec(80, 2), // 80%
		// Auto-compound parameters
		AutoCompoundEnabled:         false,                             // Disabled by default
		AutoCompoundFrequencyBlocks: 28800,                             // ~24 hours at 3s blocks
		MaxRateChangePerUpdate:      math.LegacyNewDecWithPrec(1, 2),   // 1% max change per update
		MinBlocksBetweenUpdates:     100,                               // ~5 minutes at 3s blocks
	}
}

// Validate performs validation of ModuleParams
func (p ModuleParams) Validate() error {
	// Validate existing caps
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

	// Validate minimum liquid stake amount
	if p.MinLiquidStakeAmount.IsNegative() {
		return fmt.Errorf("minimum liquid stake amount cannot be negative")
	}

	// Validate rate limiting parameters
	if p.RateLimitPeriodHours == 0 {
		return fmt.Errorf("rate limit period hours must be greater than 0")
	}
	if p.RateLimitPeriodHours > 168 { // 7 days max
		return fmt.Errorf("rate limit period hours cannot exceed 168 (7 days)")
	}

	if p.GlobalDailyTokenizationPercent.IsNegative() {
		return fmt.Errorf("global daily tokenization percent cannot be negative")
	}
	if p.GlobalDailyTokenizationPercent.GT(math.LegacyOneDec()) {
		return fmt.Errorf("global daily tokenization percent cannot exceed 100%%")
	}

	if p.ValidatorDailyTokenizationPercent.IsNegative() {
		return fmt.Errorf("validator daily tokenization percent cannot be negative")
	}
	if p.ValidatorDailyTokenizationPercent.GT(math.LegacyOneDec()) {
		return fmt.Errorf("validator daily tokenization percent cannot exceed 100%%")
	}

	if p.GlobalDailyTokenizationCount == 0 {
		return fmt.Errorf("global daily tokenization count must be greater than 0")
	}
	if p.ValidatorDailyTokenizationCount == 0 {
		return fmt.Errorf("validator daily tokenization count must be greater than 0")
	}
	if p.UserDailyTokenizationCount == 0 {
		return fmt.Errorf("user daily tokenization count must be greater than 0")
	}

	if p.WarningThresholdPercent.IsNegative() {
		return fmt.Errorf("warning threshold percent cannot be negative")
	}
	if p.WarningThresholdPercent.GT(math.LegacyOneDec()) {
		return fmt.Errorf("warning threshold percent cannot exceed 100%%")
	}

	// Validate auto-compound parameters
	if p.AutoCompoundEnabled && p.AutoCompoundFrequencyBlocks <= 0 {
		return fmt.Errorf("auto-compound frequency must be positive when enabled")
	}
	
	if p.MaxRateChangePerUpdate.IsNegative() {
		return fmt.Errorf("max rate change per update cannot be negative")
	}
	if p.MaxRateChangePerUpdate.GT(math.LegacyOneDec()) {
		return fmt.Errorf("max rate change per update cannot exceed 100%%")
	}
	
	if p.MinBlocksBetweenUpdates < 0 {
		return fmt.Errorf("min blocks between updates cannot be negative")
	}

	return nil
}

// NewParams creates a new ModuleParams instance
func NewParams(globalCap, validatorCap math.LegacyDec, enabled bool, minAmount math.Int) ModuleParams {
	// Return with default rate limiting parameters
	defaults := DefaultParams()
	return ModuleParams{
		GlobalLiquidStakingCap: globalCap,
		ValidatorLiquidCap:     validatorCap,
		Enabled:                enabled,
		MinLiquidStakeAmount:   minAmount,
		// Use default rate limiting parameters
		RateLimitPeriodHours:              defaults.RateLimitPeriodHours,
		GlobalDailyTokenizationPercent:    defaults.GlobalDailyTokenizationPercent,
		ValidatorDailyTokenizationPercent: defaults.ValidatorDailyTokenizationPercent,
		GlobalDailyTokenizationCount:      defaults.GlobalDailyTokenizationCount,
		ValidatorDailyTokenizationCount:   defaults.ValidatorDailyTokenizationCount,
		UserDailyTokenizationCount:        defaults.UserDailyTokenizationCount,
		WarningThresholdPercent:           defaults.WarningThresholdPercent,
		// Use default auto-compound parameters
		AutoCompoundEnabled:         defaults.AutoCompoundEnabled,
		AutoCompoundFrequencyBlocks: defaults.AutoCompoundFrequencyBlocks,
		MaxRateChangePerUpdate:      defaults.MaxRateChangePerUpdate,
		MinBlocksBetweenUpdates:     defaults.MinBlocksBetweenUpdates,
	}
}