package exported

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Subspace defines the expected subspace interface for migrations
type Subspace interface {
	// GetParamSet retrieves a ParamSet from the Subspace
	GetParamSet(ctx sdk.Context, ps paramtypes.ParamSet)
	
	// HasKeyTable returns true if the Subspace has a KeyTable registered
	HasKeyTable() bool
	
	// WithKeyTable returns a Subspace with a KeyTable registered
	WithKeyTable(table paramtypes.KeyTable) paramtypes.Subspace
}