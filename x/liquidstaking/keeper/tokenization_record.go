package keeper

import (
	"cosmossdk.io/math"
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

// GetTotalLiquidStaked returns the total amount of liquid staked tokens
func (k Keeper) GetTotalLiquidStaked(ctx sdk.Context) math.Int {
	store := k.storeService.OpenKVStore(ctx)
	
	bz, err := store.Get(types.TotalLiquidStakedKey)
	if err != nil {
		panic(err)
	}
	if bz == nil {
		return math.ZeroInt()
	}
	
	var amount math.Int
	err = amount.Unmarshal(bz)
	if err != nil {
		panic(err)
	}
	
	return amount
}

// SetTotalLiquidStaked sets the total amount of liquid staked tokens
func (k Keeper) SetTotalLiquidStaked(ctx sdk.Context, amount math.Int) {
	store := k.storeService.OpenKVStore(ctx)
	
	bz, err := amount.Marshal()
	if err != nil {
		panic(err)
	}
	
	err = store.Set(types.TotalLiquidStakedKey, bz)
	if err != nil {
		panic(err)
	}
}

// GetValidatorLiquidStaked returns the amount of liquid staked tokens for a validator
func (k Keeper) GetValidatorLiquidStaked(ctx sdk.Context, validatorAddr string) math.Int {
	store := k.storeService.OpenKVStore(ctx)
	key := types.GetValidatorLiquidStakedKey(validatorAddr)
	
	bz, err := store.Get(key)
	if err != nil {
		panic(err)
	}
	if bz == nil {
		return math.ZeroInt()
	}
	
	var amount math.Int
	err = amount.Unmarshal(bz)
	if err != nil {
		panic(err)
	}
	
	return amount
}

// SetValidatorLiquidStaked sets the amount of liquid staked tokens for a validator
func (k Keeper) SetValidatorLiquidStaked(ctx sdk.Context, validatorAddr string, amount math.Int) {
	store := k.storeService.OpenKVStore(ctx)
	key := types.GetValidatorLiquidStakedKey(validatorAddr)
	
	bz, err := amount.Marshal()
	if err != nil {
		panic(err)
	}
	
	err = store.Set(key, bz)
	if err != nil {
		panic(err)
	}
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