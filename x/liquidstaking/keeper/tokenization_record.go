package keeper

import (
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	
	"github.com/rollchains/flora/x/liquidstaking/types"
)

// SetTokenizationRecordWithIndexes stores a tokenization record and updates all indexes
// This extends the basic SetTokenizationRecord from genesis.go with index management
func (k Keeper) SetTokenizationRecordWithIndexes(ctx sdk.Context, record types.TokenizationRecord) {
	// Store the record using the existing method
	k.SetTokenizationRecord(ctx, record)
	
	// Update indexes
	k.setTokenizationRecordIndexes(ctx, record)
}

// GetTokenizationRecordsByValidator returns all tokenization records for a validator
func (k Keeper) GetTokenizationRecordsByValidator(ctx sdk.Context, validatorAddr string) []types.TokenizationRecord {
	store := k.storeService.OpenKVStore(ctx)
	
	var records []types.TokenizationRecord
	prefix := types.GetTokenizationRecordByValidatorPrefixKey(validatorAddr)
	iterator, err := store.Iterator(prefix, storetypes.PrefixEndBytes(prefix))
	if err != nil {
		panic(err)
	}
	defer iterator.Close()
	
	for ; iterator.Valid(); iterator.Next() {
		// Extract record ID from the key
		key := iterator.Key()
		idBytes := key[len(prefix):]
		id := types.BytesToUint64(idBytes)
		
		// Fetch the actual record
		if record, found := k.GetTokenizationRecord(ctx, id); found {
			records = append(records, record)
		}
	}
	
	return records
}

// GetTokenizationRecordsByOwner returns all tokenization records for an owner
func (k Keeper) GetTokenizationRecordsByOwner(ctx sdk.Context, ownerAddr string) []types.TokenizationRecord {
	store := k.storeService.OpenKVStore(ctx)
	
	var records []types.TokenizationRecord
	prefix := types.GetTokenizationRecordByOwnerPrefixKey(ownerAddr)
	iterator, err := store.Iterator(prefix, storetypes.PrefixEndBytes(prefix))
	if err != nil {
		panic(err)
	}
	defer iterator.Close()
	
	for ; iterator.Valid(); iterator.Next() {
		// Extract record ID from the key
		key := iterator.Key()
		idBytes := key[len(prefix):]
		id := types.BytesToUint64(idBytes)
		
		// Fetch the actual record
		if record, found := k.GetTokenizationRecord(ctx, id); found {
			records = append(records, record)
		}
	}
	
	return records
}

// GetTokenizationRecordByDenom returns the tokenization record for a specific denom
func (k Keeper) GetTokenizationRecordByDenom(ctx sdk.Context, denom string) (types.TokenizationRecord, bool) {
	store := k.storeService.OpenKVStore(ctx)
	key := types.GetTokenizationRecordByDenomKey(denom)
	
	bz, err := store.Get(key)
	if err != nil {
		panic(err)
	}
	if bz == nil {
		return types.TokenizationRecord{}, false
	}
	
	// The value stored is the record ID
	id := types.BytesToUint64(bz)
	return k.GetTokenizationRecord(ctx, id)
}

// getTokenizationRecordByDenom returns the tokenization record ID for a specific denom
func (k Keeper) getTokenizationRecordByDenom(ctx sdk.Context, denom string) (uint64, bool) {
	store := k.storeService.OpenKVStore(ctx)
	key := types.GetTokenizationRecordByDenomKey(denom)
	
	bz, err := store.Get(key)
	if err != nil {
		panic(err)
	}
	if bz == nil {
		return 0, false
	}
	
	// The value stored is the record ID
	id := types.BytesToUint64(bz)
	return id, true
}

// setTokenizationRecordDenomIndex sets the denom index for a tokenization record
func (k Keeper) setTokenizationRecordDenomIndex(ctx sdk.Context, denom string, recordID uint64) {
	store := k.storeService.OpenKVStore(ctx)
	key := types.GetTokenizationRecordByDenomKey(denom)
	value := types.Uint64ToBytes(recordID)
	
	err := store.Set(key, value)
	if err != nil {
		panic(err)
	}
}

// removeTokenizationRecordDenomIndex removes the denom index for a tokenization record
func (k Keeper) removeTokenizationRecordDenomIndex(ctx sdk.Context, denom string) {
	store := k.storeService.OpenKVStore(ctx)
	key := types.GetTokenizationRecordByDenomKey(denom)
	
	err := store.Delete(key)
	if err != nil {
		panic(err)
	}
}

// deleteTokenizationRecordDenomIndex is an alias for removeTokenizationRecordDenomIndex
func (k Keeper) deleteTokenizationRecordDenomIndex(ctx sdk.Context, denom string) {
	k.removeTokenizationRecordDenomIndex(ctx, denom)
}


// setTokenizationRecordIndexes updates all indexes for a tokenization record
func (k Keeper) setTokenizationRecordIndexes(ctx sdk.Context, record types.TokenizationRecord) {
	store := k.storeService.OpenKVStore(ctx)
	
	// Update validator index
	validatorKey := types.GetTokenizationRecordByValidatorKey(record.Validator, record.Id)
	err := store.Set(validatorKey, []byte{})
	if err != nil {
		panic(err)
	}
	
	// Update owner index
	ownerKey := types.GetTokenizationRecordByOwnerKey(record.Owner, record.Id)
	err = store.Set(ownerKey, []byte{})
	if err != nil {
		panic(err)
	}
	
	// Note: Denom index will be set when the liquid staking token is minted in Stage 3
}

// removeTokenizationRecordIndexes removes all indexes for a tokenization record
func (k Keeper) removeTokenizationRecordIndexes(ctx sdk.Context, record types.TokenizationRecord) {
	store := k.storeService.OpenKVStore(ctx)
	
	// Remove validator index
	validatorKey := types.GetTokenizationRecordByValidatorKey(record.Validator, record.Id)
	err := store.Delete(validatorKey)
	if err != nil {
		panic(err)
	}
	
	// Remove owner index
	ownerKey := types.GetTokenizationRecordByOwnerKey(record.Owner, record.Id)
	err = store.Delete(ownerKey)
	if err != nil {
		panic(err)
	}
	
	// Note: Denom index removal will be handled when implemented in Stage 3
}

// DeleteTokenizationRecord removes a tokenization record and its indexes
func (k Keeper) DeleteTokenizationRecord(ctx sdk.Context, id uint64) error {
	record, found := k.GetTokenizationRecord(ctx, id)
	if !found {
		return types.ErrTokenizationRecordNotFound
	}
	
	store := k.storeService.OpenKVStore(ctx)
	
	// Remove the record
	key := types.GetTokenizationRecordKey(id)
	err := store.Delete(key)
	if err != nil {
		return err
	}
	
	// Remove indexes
	k.removeTokenizationRecordIndexes(ctx, record)
	
	return nil
}

// DeleteTokenizationRecordWithIndexes removes a tokenization record and all its indexes including denom index
func (k Keeper) DeleteTokenizationRecordWithIndexes(ctx sdk.Context, id uint64) {
	record, found := k.GetTokenizationRecord(ctx, id)
	if !found {
		return
	}
	
	store := k.storeService.OpenKVStore(ctx)
	
	// Remove the record
	key := types.GetTokenizationRecordKey(id)
	err := store.Delete(key)
	if err != nil {
		panic(err)
	}
	
	// Remove all indexes
	k.removeTokenizationRecordIndexes(ctx, record)
	
	// Remove denom index if it exists
	if record.Denom != "" {
		k.removeTokenizationRecordDenomIndex(ctx, record.Denom)
	}
}

// ValidateTokenizationRecord validates a tokenization record before storing
func (k Keeper) ValidateTokenizationRecord(ctx sdk.Context, record types.TokenizationRecord) error {
	// Basic validation
	if err := record.Validate(); err != nil {
		return err
	}
	
	// Check if record ID already exists
	if _, found := k.GetTokenizationRecord(ctx, record.Id); found {
		return types.ErrTokenizationRecordAlreadyExists
	}
	
	// Check if denom already exists
	if _, found := k.GetTokenizationRecordByDenom(ctx, record.Denom); found {
		return types.ErrDuplicateLiquidStakingToken
	}
	
	// Additional validation can be added here in future stages
	
	return nil
}