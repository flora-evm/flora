package types

import (
	"fmt"
	"strings"
)

const (
	// LiquidStakingTokenPrefix is the prefix for all liquid staking token denoms
	LiquidStakingTokenPrefix = "flora/lstake/"
)

// GenerateLiquidStakingTokenDenom generates a unique denomination for a liquid staking token
// Format: flora/lstake/{validatorAddr}/{recordId}
func GenerateLiquidStakingTokenDenom(validatorAddr string, recordID uint64) string {
	return fmt.Sprintf("%s%s/%d", LiquidStakingTokenPrefix, validatorAddr, recordID)
}

// ParseLiquidStakingTokenDenom parses a liquid staking token denom into its components
// Returns validatorAddr, recordID, and error
func ParseLiquidStakingTokenDenom(denom string) (string, uint64, error) {
	if !IsLiquidStakingTokenDenom(denom) {
		return "", 0, fmt.Errorf("not a liquid staking token denom: %s", denom)
	}
	
	// Remove prefix
	remainder := strings.TrimPrefix(denom, LiquidStakingTokenPrefix)
	
	// Split by last slash to get validator and record ID
	lastSlash := strings.LastIndex(remainder, "/")
	if lastSlash == -1 {
		return "", 0, fmt.Errorf("invalid liquid staking token denom format: %s", denom)
	}
	
	validatorAddr := remainder[:lastSlash]
	recordIDStr := remainder[lastSlash+1:]
	
	var recordID uint64
	_, err := fmt.Sscanf(recordIDStr, "%d", &recordID)
	if err != nil {
		return "", 0, fmt.Errorf("failed to parse record ID from denom %s: %w", denom, err)
	}
	
	return validatorAddr, recordID, nil
}

// IsLiquidStakingTokenDenom checks if a denom is a liquid staking token
func IsLiquidStakingTokenDenom(denom string) bool {
	return strings.HasPrefix(denom, LiquidStakingTokenPrefix)
}

// GetValidatorFromLiquidStakingTokenDenom extracts the validator address from a liquid staking token denom
func GetValidatorFromLiquidStakingTokenDenom(denom string) (string, error) {
	validatorAddr, _, err := ParseLiquidStakingTokenDenom(denom)
	if err != nil {
		return "", err
	}
	return validatorAddr, nil
}

// GetRecordIDFromLiquidStakingTokenDenom extracts the record ID from a liquid staking token denom
func GetRecordIDFromLiquidStakingTokenDenom(denom string) (uint64, error) {
	_, recordID, err := ParseLiquidStakingTokenDenom(denom)
	if err != nil {
		return 0, err
	}
	return recordID, nil
}