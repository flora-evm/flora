# Genesis Handling for Liquid Staking Module

## Overview

The liquid staking module's genesis handling ensures proper initialization and export of state, including tokenization records, liquid staking counters, and module parameters.

## Genesis State Structure

```protobuf
message GenesisState {
  Params params = 1;
  repeated TokenizationRecord tokenization_records = 2;
  uint64 last_tokenization_record_id = 3;
}
```

## Key Features

### 1. Import Process (InitGenesis)

The initialization process performs the following steps:

```go
func (k Keeper) InitGenesis(ctx sdk.Context, genState types.GenesisState) {
    // 1. Set module parameters
    k.SetParams(ctx, genState.Params)
    
    // 2. Set last tokenization record ID
    k.SetLastTokenizationRecordID(ctx, genState.LastTokenizationRecordId)
    
    // 3. Set tokenization records with proper indexes
    for _, record := range genState.TokenizationRecords {
        k.SetTokenizationRecordWithIndexes(ctx, record)
        if record.Denom != "" {
            k.setTokenizationRecordDenomIndex(ctx, record.Denom, record.Id)
        }
    }
    
    // 4. Initialize liquid staking counters
    k.initializeLiquidStakingCounters(ctx, genState.TokenizationRecords)
}
```

### 2. Export Process (ExportGenesis)

```go
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
    return &types.GenesisState{
        Params:                   k.GetParams(ctx),
        TokenizationRecords:      k.GetAllTokenizationRecords(ctx),
        LastTokenizationRecordId: k.GetLastTokenizationRecordID(ctx),
    }
}
```

### 3. Counter Initialization

The module automatically calculates and sets liquid staking counters from tokenization records:

```go
func (k Keeper) initializeLiquidStakingCounters(ctx sdk.Context, records []types.TokenizationRecord) {
    totalLiquidStaked := math.ZeroInt()
    validatorLiquidStaked := make(map[string]math.Int)
    
    for _, record := range records {
        totalLiquidStaked = totalLiquidStaked.Add(record.SharesTokenized)
        validatorLiquidStaked[record.Validator] = 
            validatorLiquidStaked[record.Validator].Add(record.SharesTokenized)
    }
    
    k.SetTotalLiquidStaked(ctx, totalLiquidStaked)
    for validator, amount := range validatorLiquidStaked {
        k.SetValidatorLiquidStaked(ctx, validator, amount)
    }
}
```

## LST Metadata Integration

While LST metadata is stored in the Bank module, the liquid staking module ensures proper handling during genesis:

### Metadata Export (in module.go)
```go
func ExportLSTMetadata(ctx sdk.Context, k keeper.Keeper, bankKeeper types.BankKeeper) []banktypes.Metadata {
    var metadata []banktypes.Metadata
    records := k.GetAllTokenizationRecords(ctx)
    
    for _, record := range records {
        if meta, found := bankKeeper.GetDenomMetaData(ctx, record.Denom); found {
            metadata = append(metadata, meta)
        }
    }
    return metadata
}
```

### Metadata Import (in module.go)
```go
func InitLSTMetadata(ctx sdk.Context, k keeper.Keeper, bankKeeper types.BankKeeper, metadata []banktypes.Metadata) {
    for _, meta := range metadata {
        if types.IsLiquidStakingTokenDenom(meta.Base) {
            bankKeeper.SetDenomMetaData(ctx, meta)
        }
    }
}
```

## Index Management

Genesis import ensures all indexes are properly created:

1. **Validator Index**: Maps validator addresses to their tokenization records
2. **Owner Index**: Maps owner addresses to their tokenization records  
3. **Denom Index**: Maps LST denoms to their tokenization record IDs

## Validation

The module includes validation to ensure genesis state integrity:

```go
func ValidateLSTGenesis(ctx sdk.Context, k keeper.Keeper, bankKeeper types.BankKeeper, records []types.TokenizationRecord) error {
    denomMap := make(map[string]bool)
    
    for _, record := range records {
        // Check for duplicate denoms
        if denomMap[record.Denom] {
            return ErrDuplicateLiquidStakingToken
        }
        
        // Validate denom format
        if !types.IsLiquidStakingTokenDenom(record.Denom) {
            return ErrInvalidLiquidStakingTokenDenom
        }
        
        // Check metadata exists
        if _, found := bankKeeper.GetDenomMetaData(ctx, record.Denom); !found {
            return ErrMetadataNotFound
        }
    }
    return nil
}
```

## Migration Support

The module includes a placeholder for future migrations:

```go
func MigrateLSTDenoms(ctx sdk.Context, k keeper.Keeper, bankKeeper types.BankKeeper) error {
    // Placeholder for future denom format migrations
    return nil
}
```

## Testing

Genesis functionality is thoroughly tested including:
- Import/export round trips
- Empty state handling
- Index creation verification
- Counter initialization
- Metadata handling

See `keeper/genesis_test.go` for comprehensive test coverage.