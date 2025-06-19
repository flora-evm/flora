package types

import (
	"github.com/rollchains/flora/x/liquidstaking/exported"
)

// Import the expected keeper interfaces
type (
	StakingKeeper       = exported.StakingKeeper
	BankKeeper          = exported.BankKeeper
	AccountKeeper       = exported.AccountKeeper
	TransferKeeper      = exported.TransferKeeper
	ChannelKeeper       = exported.ChannelKeeper
	DistributionKeeper  = exported.DistributionKeeper
	// TokenFactoryKeeper removed - using Bank module directly for LST management
	// TokenFactoryKeeper  = exported.TokenFactoryKeeper
	// DenomAuthorityMetadata = exported.DenomAuthorityMetadata
)