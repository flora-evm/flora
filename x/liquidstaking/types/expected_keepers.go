package types

import (
	"github.com/rollchains/flora/x/liquidstaking/exported"
)

// Import the expected keeper interfaces
type (
	StakingKeeper = exported.StakingKeeper
	BankKeeper    = exported.BankKeeper
	AccountKeeper = exported.AccountKeeper
)