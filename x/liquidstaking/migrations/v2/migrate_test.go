package v2_test

import (
	"testing"

	"cosmossdk.io/math"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	"cosmossdk.io/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	v2 "github.com/rollchains/flora/x/liquidstaking/migrations/v2"
	"github.com/rollchains/flora/x/liquidstaking/types"
)

func TestMigrateStore(t *testing.T) {
	// Setup store and codec
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	db := dbm.NewMemDB()
	stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	require.NoError(t, stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)
	ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())

	store := ctx.KVStore(storeKey)

	// Setup v1 data
	// Add tokenization records
	records := []types.TokenizationRecord{
		{
			Id:              1,
			Validator:       "floravaloper1validator1",
			Owner:           "flora1owner1",
			SharesTokenized: math.NewInt(1000000),
			Denom:           types.GenerateLiquidStakingTokenDenom("floravaloper1validator1", 1),
		},
		{
			Id:              2,
			Validator:       "floravaloper1validator2",
			Owner:           "flora1owner2",
			SharesTokenized: math.NewInt(2000000),
			Denom:           types.GenerateLiquidStakingTokenDenom("floravaloper1validator2", 2),
		},
	}

	for _, record := range records {
		key := types.GetTokenizationRecordKey(record.Id)
		bz := cdc.MustMarshal(&record)
		store.Set(key, bz)
	}

	// Add parameters
	params := types.DefaultParams()
	bz := cdc.MustMarshal(&params)
	store.Set(types.ParamsKey, bz)

	// Add last tokenization record ID
	store.Set(types.LastTokenizationRecordIDKey, types.Uint64ToBytes(2))

	// Run migration
	err := v2.MigrateStore(ctx, storeKey, cdc)
	require.NoError(t, err)

	// Verify data is still accessible after migration
	// Check records
	for _, expectedRecord := range records {
		key := types.GetTokenizationRecordKey(expectedRecord.Id)
		bz := store.Get(key)
		require.NotNil(t, bz)

		var actualRecord types.TokenizationRecord
		err := cdc.Unmarshal(bz, &actualRecord)
		require.NoError(t, err)
		require.Equal(t, expectedRecord, actualRecord)
	}

	// Check params
	paramsBz := store.Get(types.ParamsKey)
	require.NotNil(t, paramsBz)

	var migratedParams types.ModuleParams
	err = cdc.Unmarshal(paramsBz, &migratedParams)
	require.NoError(t, err)
	require.Equal(t, params, migratedParams)

	// Check last ID
	lastIDBz := store.Get(types.LastTokenizationRecordIDKey)
	require.NotNil(t, lastIDBz)
	require.Equal(t, uint64(2), types.BytesToUint64(lastIDBz))
}

func TestMigrateEmptyStore(t *testing.T) {
	// Setup empty store
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	db := dbm.NewMemDB()
	stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	require.NoError(t, stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)
	ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())

	// Run migration on empty store
	err := v2.MigrateStore(ctx, storeKey, cdc)
	require.NoError(t, err)

	// Verify store is still empty
	store := ctx.KVStore(storeKey)
	iterator := store.Iterator(nil, nil)
	require.False(t, iterator.Valid())
	iterator.Close()
}