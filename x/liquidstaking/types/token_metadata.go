package types

import (
	"fmt"
	"strings"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

// GenerateLiquidStakingTokenMetadata creates bank metadata for a liquid staking token
func GenerateLiquidStakingTokenMetadata(validatorAddr string, recordID uint64) banktypes.Metadata {
	denom := GenerateLiquidStakingTokenDenom(validatorAddr, recordID)
	
	// Extract validator moniker suffix for display
	// e.g., floravaloper1abcd... -> abcd
	monikerSuffix := ""
	if strings.HasPrefix(validatorAddr, "floravaloper1") && len(validatorAddr) > 16 {
		monikerSuffix = validatorAddr[13:17]
	} else if strings.HasPrefix(validatorAddr, "floravaloper1") && len(validatorAddr) > 13 {
		// For shorter addresses, take what's available after "floravaloper1"
		monikerSuffix = validatorAddr[13:]
	}
	
	return banktypes.Metadata{
		Description: fmt.Sprintf("Liquid staking token for validator %s", validatorAddr),
		DenomUnits: []*banktypes.DenomUnit{
			{
				Denom:    denom,
				Exponent: 0,
				Aliases:  []string{},
			},
			{
				Denom:    fmt.Sprintf("LSTFLORA%s", monikerSuffix),
				Exponent: 18, // Same as native flora token
				Aliases:  []string{},
			},
		},
		Base:    denom,
		Display: fmt.Sprintf("LSTFLORA%s", monikerSuffix),
		Name:    fmt.Sprintf("Liquid Staked FLORA %s", monikerSuffix),
		Symbol:  fmt.Sprintf("LSTFLORA%s", monikerSuffix),
		URI:     "", // Can be set later if needed
		URIHash: "", // Can be set later if needed
	}
}

// ValidateLiquidStakingTokenMetadata validates that the metadata is for a liquid staking token
func ValidateLiquidStakingTokenMetadata(metadata banktypes.Metadata) error {
	if !IsLiquidStakingTokenDenom(metadata.Base) {
		return fmt.Errorf("metadata base denom is not a liquid staking token: %s", metadata.Base)
	}
	
	if len(metadata.DenomUnits) < 1 {
		return fmt.Errorf("metadata must have at least one denom unit")
	}
	
	if metadata.DenomUnits[0].Denom != metadata.Base {
		return fmt.Errorf("first denom unit must match base denom")
	}
	
	if metadata.Display == "" {
		return fmt.Errorf("display denom cannot be empty")
	}
	
	if metadata.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	
	if metadata.Symbol == "" {
		return fmt.Errorf("symbol cannot be empty")
	}
	
	return nil
}