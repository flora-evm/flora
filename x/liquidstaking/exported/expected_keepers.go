package exported

// StakingKeeper defines the expected staking keeper interface
// This will be expanded in later stages when we integrate with staking module
type StakingKeeper interface {
	// Methods will be added in Stage 5
}

// BankKeeper defines the expected bank keeper interface
// This will be expanded in later stages when we handle token transfers
type BankKeeper interface {
	// Methods will be added in later stages
}

// TokenFactoryKeeper defines the expected token factory keeper interface
// This will be expanded in Stage 6 when we integrate with token factory
type TokenFactoryKeeper interface {
	// Methods will be added in Stage 6
}