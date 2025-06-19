package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/rollchains/flora/x/liquidstaking/exported"
)

// Migrator is a struct for handling in-place store migrations
type Migrator struct {
	keeper         exported.Keeper
	legacySubspace exported.Subspace
}

// NewMigrator returns a new Migrator
func NewMigrator(keeper exported.Keeper, legacySubspace exported.Subspace) Migrator {
	return Migrator{
		keeper:         keeper,
		legacySubspace: legacySubspace,
	}
}

// Migrate1to2 migrates from version 1 to 2
func (m Migrator) Migrate1to2(ctx sdk.Context) error {
	// This is a placeholder for v1 to v2 migration
	// In a real migration, you would:
	// 1. Update store structure
	// 2. Migrate data formats
	// 3. Add new indexes
	
	// Example: Log the migration
	ctx.Logger().Info("migrating liquid staking module from version 1 to 2")
	
	// The actual migration logic would be imported from the v2 package
	// return v2.MigrateStore(ctx, m.keeper.StoreKey(), m.keeper.Codec())
	
	return nil
}

// GetMigrations returns the list of migrations for the liquid staking module
func GetMigrations(keeper exported.Keeper, legacySubspace exported.Subspace) map[string]module.MigrationHandler {
	m := NewMigrator(keeper, legacySubspace)
	
	return map[string]module.MigrationHandler{
		"1": m.Migrate1to2,
		// Future migrations would be added here:
		// "2": m.Migrate2to3,
	}
}

// CurrentVersion returns the current version of the module for migration purposes
func CurrentVersion() uint64 {
	return 2
}