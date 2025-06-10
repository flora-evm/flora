package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenerateTestAddresses generates valid test addresses
func GenerateTestAddresses() (validatorAddr, ownerAddr string) {
	// Set bech32 prefixes before generating addresses
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount("flora", "florapub")
	config.SetBech32PrefixForValidator("floravaloper", "floravaloperpub")
	config.SetBech32PrefixForConsensusNode("floravalcons", "floravalconspub")
	
	// Generate a test address from a fixed seed for consistency
	addr := sdk.AccAddress([]byte("test"))
	ownerAddr = addr.String()
	
	// Convert to validator address
	valAddr := sdk.ValAddress(addr)
	validatorAddr = valAddr.String()
	
	return validatorAddr, ownerAddr
}

// Test addresses generated with proper checksums
var TestValidatorAddr, TestOwnerAddr = GenerateTestAddresses()