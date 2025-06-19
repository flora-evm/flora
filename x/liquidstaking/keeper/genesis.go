package keeper

import (
	"fmt"
	"time"
	
	"cosmossdk.io/math"
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
	
	// Set tokenization records with proper indexes
	for _, record := range genState.TokenizationRecords {
		// Use SetTokenizationRecordWithIndexes to ensure all indexes are created
		k.SetTokenizationRecordWithIndexes(ctx, record)
		
		// Set the denom index which is normally set during minting
		if record.Denom != "" {
			k.setTokenizationRecordDenomIndex(ctx, record.Denom, record.Id)
		}
	}
	
	// Initialize liquid staking counters from records
	k.initializeLiquidStakingCounters(ctx, genState.TokenizationRecords)
	
	// Import exchange rates
	for _, rate := range genState.ExchangeRates {
		k.SetExchangeRate(ctx, rate.ValidatorAddress, rate.Rate, time.Unix(rate.LastUpdated, 0))
	}
	
	// Import global exchange rate if present
	if genState.GlobalExchangeRate != nil {
		k.SetGlobalExchangeRate(ctx, *genState.GlobalExchangeRate)
	}
}

// ExportGenesis returns the liquid staking module's exported genesis state
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	// Collect all exchange rates
	var exchangeRates []types.ExchangeRate
	k.IterateExchangeRates(ctx, func(rate types.ExchangeRate) bool {
		exchangeRates = append(exchangeRates, rate)
		return false
	})
	
	// Get global exchange rate if it exists
	var globalExchangeRate *types.GlobalExchangeRate
	if globalRate, found := k.GetGlobalExchangeRate(ctx); found {
		globalExchangeRate = &globalRate
	}
	
	return &types.GenesisState{
		Params:                   k.GetParams(ctx),
		TokenizationRecords:      k.GetAllTokenizationRecords(ctx),
		LastTokenizationRecordId: k.GetLastTokenizationRecordID(ctx),
		ExchangeRates:           exchangeRates,
		GlobalExchangeRate:      globalExchangeRate,
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

// initializeLiquidStakingCounters initializes the liquid staking counters from tokenization records
func (k Keeper) initializeLiquidStakingCounters(ctx sdk.Context, records []types.TokenizationRecord) {
	// Reset counters
	totalLiquidStaked := math.ZeroInt()
	validatorLiquidStaked := make(map[string]math.Int)
	validatorsWithLST := make(map[string]bool)
	
	// Calculate totals from records
	for _, record := range records {
		totalLiquidStaked = totalLiquidStaked.Add(record.SharesTokenized)
		
		if current, exists := validatorLiquidStaked[record.Validator]; exists {
			validatorLiquidStaked[record.Validator] = current.Add(record.SharesTokenized)
		} else {
			validatorLiquidStaked[record.Validator] = record.SharesTokenized
		}
		
		// Track validators with LST tokens
		validatorsWithLST[record.Validator] = true
	}
	
	// Set total liquid staked
	k.SetTotalLiquidStaked(ctx, totalLiquidStaked)
	
	// Set per-validator liquid staked amounts and ensure exchange rates
	for validator, amount := range validatorLiquidStaked {
		k.SetValidatorLiquidStaked(ctx, validator, amount)
		
		// Ensure exchange rate is initialized for validators with LST tokens
		if _, found := k.GetExchangeRate(ctx, validator); !found {
			// Initialize to 1:1 if not set
			k.SetExchangeRate(ctx, validator, math.LegacyOneDec(), ctx.BlockTime())
			k.Logger(ctx).Info("initialized exchange rate for validator with LST tokens",
				"validator", validator,
				"rate", math.LegacyOneDec().String(),
			)
		}
	}
}