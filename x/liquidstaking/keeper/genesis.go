package keeper

import (
	"fmt"
	
	sdk "github.com/cosmos/cosmos-sdk/types"
	
	"github.com/rollchains/flora/x/liquidstaking/types"
)

// InitGenesis initializes the liquid staking module's state from a provided genesis state
func (k Keeper) InitGenesis(ctx sdk.Context, genState types.GenesisState) {
	// Set module parameters
	if err := k.SetParams(ctx, genState.Params); err != nil {
		panic(err)
	}
	
	// Set last tokenization record ID
	k.SetLastTokenizationRecordID(ctx, genState.LastTokenizationRecordId)
	
	// Set tokenization records
	for _, record := range genState.TokenizationRecords {
		k.SetTokenizationRecord(ctx, record)
	}
}

// ExportGenesis returns the liquid staking module's exported genesis state
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	return &types.GenesisState{
		Params:                   k.GetParams(ctx),
		TokenizationRecords:      k.GetAllTokenizationRecords(ctx),
		LastTokenizationRecordId: k.GetLastTokenizationRecordID(ctx),
	}
}

// SetLastTokenizationRecordID sets the last tokenization record ID
func (k Keeper) SetLastTokenizationRecordID(ctx sdk.Context, id uint64) {
	store := k.storeService.OpenKVStore(ctx)
	bz := types.Uint64ToBytes(id)
	if err := store.Set(types.LastTokenizationRecordIDKey, bz); err != nil {
		panic(err)
	}
}

// GetLastTokenizationRecordID gets the last tokenization record ID
func (k Keeper) GetLastTokenizationRecordID(ctx sdk.Context) uint64 {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := store.Get(types.LastTokenizationRecordIDKey)
	if err != nil {
		panic(err)
	}
	if bz == nil {
		return 0
	}
	return types.BytesToUint64(bz)
}

// SetTokenizationRecord sets a tokenization record in the store
func (k Keeper) SetTokenizationRecord(ctx sdk.Context, record types.TokenizationRecord) {
	store := k.storeService.OpenKVStore(ctx)
	bz := k.cdc.MustMarshal(&record)
	key := types.GetTokenizationRecordKey(record.Id)
	if err := store.Set(key, bz); err != nil {
		panic(err)
	}
}

// GetTokenizationRecord gets a tokenization record from the store
func (k Keeper) GetTokenizationRecord(ctx sdk.Context, id uint64) (types.TokenizationRecord, bool) {
	store := k.storeService.OpenKVStore(ctx)
	key := types.GetTokenizationRecordKey(id)
	bz, err := store.Get(key)
	if err != nil {
		panic(err)
	}
	if bz == nil {
		return types.TokenizationRecord{}, false
	}
	
	var record types.TokenizationRecord
	k.cdc.MustUnmarshal(bz, &record)
	return record, true
}

// GetAllTokenizationRecords returns all tokenization records
func (k Keeper) GetAllTokenizationRecords(ctx sdk.Context) []types.TokenizationRecord {
	store := k.storeService.OpenKVStore(ctx)
	
	// Create a prefix range for iteration
	startKey := types.TokenizationRecordPrefix
	endKey := append(types.TokenizationRecordPrefix, 0xFF)
	
	iterator, err := store.Iterator(startKey, endKey)
	if err != nil {
		panic(err)
	}
	defer iterator.Close()
	
	var records []types.TokenizationRecord
	for ; iterator.Valid(); iterator.Next() {
		value := iterator.Value()
		if len(value) == 0 {
			continue // Skip empty values
		}
		
		var record types.TokenizationRecord
		if err := k.cdc.Unmarshal(value, &record); err != nil {
			panic(fmt.Errorf("failed to unmarshal tokenization record: %w", err))
		}
		records = append(records, record)
	}
	
	return records
}