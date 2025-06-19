package keeper

import (
	"context"
	"fmt"
	
	"cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	
	"github.com/rollchains/flora/x/liquidstaking/types"
)

var _ types.QueryServer = Keeper{}

// Params returns the module parameters
func (k Keeper) Params(c context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	
	return &types.QueryParamsResponse{
		Params: k.GetParams(ctx),
	}, nil
}

// TokenizationRecord returns a specific tokenization record
func (k Keeper) TokenizationRecord(c context.Context, req *types.QueryTokenizationRecordRequest) (*types.QueryTokenizationRecordResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid record ID")
	}
	
	ctx := sdk.UnwrapSDKContext(c)
	
	record, found := k.GetTokenizationRecord(ctx, req.Id)
	if !found {
		return nil, status.Errorf(codes.NotFound, "tokenization record not found: %d", req.Id)
	}
	
	return &types.QueryTokenizationRecordResponse{
		Record: record,
	}, nil
}

// TokenizationRecords returns all tokenization records with pagination
func (k Keeper) TokenizationRecords(c context.Context, req *types.QueryTokenizationRecordsRequest) (*types.QueryTokenizationRecordsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	
	// Get the store
	store := k.storeService.OpenKVStore(ctx)
	
	// Create a prefixed store for tokenization records
	recordStore := prefix.NewStore(runtime.KVStoreAdapter(store), types.TokenizationRecordPrefix)
	
	// Paginate through the records
	var records []types.TokenizationRecord
	pageRes, err := query.Paginate(recordStore, req.Pagination, func(key []byte, value []byte) error {
		var record types.TokenizationRecord
		if err := k.cdc.Unmarshal(value, &record); err != nil {
			return err
		}
		records = append(records, record)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	
	return &types.QueryTokenizationRecordsResponse{
		Records:    records,
		Pagination: pageRes,
	}, nil
}

// TokenizationRecordsByValidator returns all tokenization records for a validator with pagination
func (k Keeper) TokenizationRecordsByValidator(c context.Context, req *types.QueryTokenizationRecordsByValidatorRequest) (*types.QueryTokenizationRecordsByValidatorResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	
	if req.ValidatorAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "empty validator address")
	}
	
	ctx := sdk.UnwrapSDKContext(c)
	
	// Get the store
	store := k.storeService.OpenKVStore(ctx)
	
	// Create a prefixed store for validator index
	validatorStore := prefix.NewStore(runtime.KVStoreAdapter(store), append(types.TokenizationRecordByValidatorPrefix, []byte(req.ValidatorAddress)...))
	
	// Paginate through the records
	var records []types.TokenizationRecord
	pageRes, err := query.Paginate(validatorStore, req.Pagination, func(key []byte, value []byte) error {
		// Value is the record ID
		recordID := sdk.BigEndianToUint64(value)
		record, found := k.GetTokenizationRecord(ctx, recordID)
		if !found {
			// Skip if record not found (shouldn't happen but be safe)
			return nil
		}
		records = append(records, record)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	
	return &types.QueryTokenizationRecordsByValidatorResponse{
		Records:    records,
		Pagination: pageRes,
	}, nil
}

// TokenizationRecordsByOwner returns all tokenization records for an owner with pagination
func (k Keeper) TokenizationRecordsByOwner(c context.Context, req *types.QueryTokenizationRecordsByOwnerRequest) (*types.QueryTokenizationRecordsByOwnerResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	
	if req.OwnerAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "empty owner address")
	}
	
	ctx := sdk.UnwrapSDKContext(c)
	
	// Get the store
	store := k.storeService.OpenKVStore(ctx)
	
	// Create a prefixed store for owner index
	ownerStore := prefix.NewStore(runtime.KVStoreAdapter(store), append(types.TokenizationRecordByOwnerPrefix, []byte(req.OwnerAddress)...))
	
	// Paginate through the records
	var records []types.TokenizationRecord
	pageRes, err := query.Paginate(ownerStore, req.Pagination, func(key []byte, value []byte) error {
		// Value is the record ID
		recordID := sdk.BigEndianToUint64(value)
		record, found := k.GetTokenizationRecord(ctx, recordID)
		if !found {
			// Skip if record not found (shouldn't happen but be safe)
			return nil
		}
		records = append(records, record)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	
	return &types.QueryTokenizationRecordsByOwnerResponse{
		Records:    records,
		Pagination: pageRes,
	}, nil
}

// TotalLiquidStaked returns the total amount of liquid staked tokens
func (k Keeper) TotalLiquidStaked(c context.Context, req *types.QueryTotalLiquidStakedRequest) (*types.QueryTotalLiquidStakedResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	
	total := k.GetTotalLiquidStaked(ctx)
	
	return &types.QueryTotalLiquidStakedResponse{
		TotalLiquidStaked: total,
	}, nil
}

// ValidatorLiquidStaked returns the amount of liquid staked tokens for a validator
func (k Keeper) ValidatorLiquidStaked(c context.Context, req *types.QueryValidatorLiquidStakedRequest) (*types.QueryValidatorLiquidStakedResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	
	if req.ValidatorAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "empty validator address")
	}
	
	ctx := sdk.UnwrapSDKContext(c)
	
	amount := k.GetValidatorLiquidStaked(ctx, req.ValidatorAddress)
	
	return &types.QueryValidatorLiquidStakedResponse{
		LiquidStaked: amount,
	}, nil
}

// TokenizationRecordsByDenom returns the tokenization record for a specific LST denomination
// TokenizationRecordsByDenom returns the tokenization record for a specific LST denomination
func (k Keeper) TokenizationRecordsByDenom(c context.Context, req *types.QueryTokenizationRecordsByDenomRequest) (*types.QueryTokenizationRecordsByDenomResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	
	if req.Denom == "" {
		return nil, status.Error(codes.InvalidArgument, "empty denomination")
	}
	
	ctx := sdk.UnwrapSDKContext(c)
	
	// Since denom is unique per record, we can directly look it up
	record, found := k.GetTokenizationRecordByDenom(ctx, req.Denom)
	if !found {
		// Return empty response, not an error
		return &types.QueryTokenizationRecordsByDenomResponse{
			Records: []types.TokenizationRecord{},
		}, nil
	}
	
	return &types.QueryTokenizationRecordsByDenomResponse{
		Records: []types.TokenizationRecord{record},
	}, nil
}

// RateLimitStatus returns the current rate limit usage for an address
// TODO: Uncomment after proto regeneration
func (k Keeper) RateLimitStatus(c context.Context, req *types.QueryRateLimitStatusRequest) (*types.QueryRateLimitStatusResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	
	if req.Address == "" {
		return nil, status.Error(codes.InvalidArgument, "empty address")
	}
	
	ctx := sdk.UnwrapSDKContext(c)
	
	var rateLimits []types.RateLimitInfo
	
	// Special case for "global" address
	if req.Address == "global" {
		globalActivity := k.GetGlobalTokenizationActivity(ctx)
		totalBonded, err := k.stakingKeeper.TotalBondedTokens(ctx)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		
		params := k.GetParams(ctx)
		maxAmount := totalBonded.Mul(math.NewInt(params.GlobalDailyTokenizationPercent.TruncateInt64())).Quo(math.NewInt(100))
		maxCount := params.GlobalDailyTokenizationCount
		
		rateLimitPeriod := k.GetRateLimitPeriod(ctx)
		windowStart := globalActivity.LastActivity.Add(-rateLimitPeriod)
		windowEnd := globalActivity.LastActivity.Add(rateLimitPeriod)
		if ctx.BlockTime().Sub(globalActivity.LastActivity) > rateLimitPeriod {
			// Window expired, show reset values
			windowStart = ctx.BlockTime()
			windowEnd = ctx.BlockTime().Add(rateLimitPeriod)
			globalActivity.TotalAmount = math.ZeroInt()
			globalActivity.ActivityCount = 0
		}
		
		rateLimits = append(rateLimits, types.RateLimitInfo{
			LimitType:     "global",
			CurrentAmount: globalActivity.TotalAmount,
			MaxAmount:     maxAmount,
			CurrentCount:  globalActivity.ActivityCount,
			MaxCount:      maxCount,
			WindowStart:   windowStart,
			WindowEnd:     windowEnd,
		})
	} else {
		// Check if it's a validator address
		valAddr, valErr := sdk.ValAddressFromBech32(req.Address)
		if valErr == nil {
			// It's a validator address
			validator, err := k.stakingKeeper.GetValidator(ctx, valAddr)
			if err == nil {
				valActivity := k.GetValidatorTokenizationActivity(ctx, req.Address)
				params := k.GetParams(ctx)
				maxAmount := validator.Tokens.Mul(math.NewInt(params.ValidatorDailyTokenizationPercent.TruncateInt64())).Quo(math.NewInt(100))
				maxCount := params.ValidatorDailyTokenizationCount
				
				rateLimitPeriod := k.GetRateLimitPeriod(ctx)
				windowStart := valActivity.LastActivity.Add(-rateLimitPeriod)
				windowEnd := valActivity.LastActivity.Add(rateLimitPeriod)
				if ctx.BlockTime().Sub(valActivity.LastActivity) > rateLimitPeriod {
					// Window expired, show reset values
					windowStart = ctx.BlockTime()
					windowEnd = ctx.BlockTime().Add(rateLimitPeriod)
					valActivity.TotalAmount = math.ZeroInt()
					valActivity.ActivityCount = 0
				}
				
				rateLimits = append(rateLimits, types.RateLimitInfo{
					LimitType:     "validator",
					CurrentAmount: valActivity.TotalAmount,
					MaxAmount:     maxAmount,
					CurrentCount:  valActivity.ActivityCount,
					MaxCount:      maxCount,
					WindowStart:   windowStart,
					WindowEnd:     windowEnd,
				})
			}
		}
		
		// Check if it's a user address
		_, userErr := sdk.AccAddressFromBech32(req.Address)
		if userErr == nil {
			// It's a user address
			userActivity := k.GetUserTokenizationActivity(ctx, req.Address)
			params := k.GetParams(ctx)
			maxCount := params.UserDailyTokenizationCount
			
			rateLimitPeriod := k.GetRateLimitPeriod(ctx)
			windowStart := userActivity.LastActivity.Add(-rateLimitPeriod)
			windowEnd := userActivity.LastActivity.Add(rateLimitPeriod)
			if ctx.BlockTime().Sub(userActivity.LastActivity) > rateLimitPeriod {
				// Window expired, show reset values
				windowStart = ctx.BlockTime()
				windowEnd = ctx.BlockTime().Add(rateLimitPeriod)
				userActivity.TotalAmount = math.ZeroInt()
				userActivity.ActivityCount = 0
			}
			
			rateLimits = append(rateLimits, types.RateLimitInfo{
				LimitType:     "user",
				CurrentAmount: userActivity.TotalAmount,
				MaxAmount:     math.ZeroInt(), // Users don't have amount limits
				CurrentCount:  userActivity.ActivityCount,
				MaxCount:      maxCount,
				WindowStart:   windowStart,
				WindowEnd:     windowEnd,
			})
		}
		
		if len(rateLimits) == 0 {
			return nil, status.Error(codes.InvalidArgument, "invalid address format")
		}
	}
	
	return &types.QueryRateLimitStatusResponse{
		RateLimits: rateLimits,
	}, nil
}

// TokenizationStatistics returns aggregated tokenization statistics
// TODO: Uncomment after proto regeneration
func (k Keeper) TokenizationStatistics(c context.Context, req *types.QueryTokenizationStatisticsRequest) (*types.QueryTokenizationStatisticsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	
	ctx := sdk.UnwrapSDKContext(c)
	
	// Get total liquid staked amount
	totalLiquidStaked := k.GetTotalLiquidStaked(ctx)
	
	// Get total records and count active records
	totalRecords := k.GetLastTokenizationRecordID(ctx)
	var activeRecords uint64
	var totalTokenized math.Int = math.ZeroInt()
	validatorsWithStake := make(map[string]bool)
	denomsCreated := make(map[string]bool)
	
	// Iterate through all records to gather statistics
	for i := uint64(1); i <= totalRecords; i++ {
		record, found := k.GetTokenizationRecord(ctx, i)
		if found {
			activeRecords++
			totalTokenized = totalTokenized.Add(record.SharesTokenized)
			if record.SharesTokenized.IsPositive() {
				validatorsWithStake[record.Validator] = true
			}
			denomsCreated[record.Denom] = true
		}
	}
	
	// Calculate average record size
	var averageRecordSize math.Int
	if activeRecords > 0 {
		averageRecordSize = totalLiquidStaked.Quo(math.NewInt(int64(activeRecords)))
	} else {
		averageRecordSize = math.ZeroInt()
	}
	
	return &types.QueryTokenizationStatisticsResponse{
		TotalTokenized:           totalTokenized,
		ActiveLiquidStaked:       totalLiquidStaked,
		TotalRecords:             totalRecords,
		ActiveRecords:            activeRecords,
		AverageRecordSize:        averageRecordSize,
		ValidatorsWithLiquidStake: uint64(len(validatorsWithStake)),
		TotalDenomsCreated:       uint64(len(denomsCreated)),
	}, nil
}

// ValidatorStatistics returns detailed statistics for a specific validator
// TODO: Uncomment after proto regeneration
func (k Keeper) ValidatorStatistics(c context.Context, req *types.QueryValidatorStatisticsRequest) (*types.QueryValidatorStatisticsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	
	if req.ValidatorAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "empty validator address")
	}
	
	ctx := sdk.UnwrapSDKContext(c)
	
	// Get validator
	valAddr, err := sdk.ValAddressFromBech32(req.ValidatorAddress)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid validator address")
	}
	
	validator, err := k.stakingKeeper.GetValidator(ctx, valAddr)
	if err != nil {
		return nil, status.Error(codes.NotFound, "validator not found")
	}
	
	// Get liquid staked amount for this validator
	liquidStaked := k.GetValidatorLiquidStaked(ctx, req.ValidatorAddress)
	
	// Calculate liquid staking percentage
	var liquidStakingPercentage math.LegacyDec
	if validator.Tokens.IsPositive() {
		liquidStakingPercentage = math.LegacyNewDecFromInt(liquidStaked).Quo(math.LegacyNewDecFromInt(validator.Tokens)).Mul(math.LegacyNewDec(100))
	} else {
		liquidStakingPercentage = math.LegacyZeroDec()
	}
	
	// Count active records for this validator
	var activeRecords uint64
	var totalRecordsCreated uint64
	
	// Get the store
	store := k.storeService.OpenKVStore(ctx)
	
	// Create a prefixed store for validator index
	validatorStore := prefix.NewStore(runtime.KVStoreAdapter(store), append(types.TokenizationRecordByValidatorPrefix, []byte(req.ValidatorAddress)...))
	
	// Count records
	iterator := validatorStore.Iterator(nil, nil)
	defer iterator.Close()
	
	for ; iterator.Valid(); iterator.Next() {
		totalRecordsCreated++
		// Value is the record ID
		recordID := sdk.BigEndianToUint64(iterator.Value())
		record, found := k.GetTokenizationRecord(ctx, recordID)
		if found && record.SharesTokenized.IsPositive() {
			activeRecords++
		}
	}
	
	// Get rate limit info for validator
	valActivity := k.GetValidatorTokenizationActivity(ctx, req.ValidatorAddress)
	params := k.GetParams(ctx)
	maxAmount := validator.Tokens.Mul(math.NewInt(params.ValidatorDailyTokenizationPercent.TruncateInt64())).Quo(math.NewInt(100))
	maxCount := params.ValidatorDailyTokenizationCount
	
	rateLimitPeriod := k.GetRateLimitPeriod(ctx)
	windowStart := valActivity.LastActivity.Add(-rateLimitPeriod)
	windowEnd := valActivity.LastActivity.Add(rateLimitPeriod)
	if ctx.BlockTime().Sub(valActivity.LastActivity) > rateLimitPeriod {
		// Window expired, show reset values
		windowStart = ctx.BlockTime()
		windowEnd = ctx.BlockTime().Add(rateLimitPeriod)
		valActivity.TotalAmount = math.ZeroInt()
		valActivity.ActivityCount = 0
	}
	
	rateLimitUsage := types.RateLimitInfo{
		LimitType:     "validator",
		CurrentAmount: valActivity.TotalAmount,
		MaxAmount:     maxAmount,
		CurrentCount:  valActivity.ActivityCount,
		MaxCount:      maxCount,
		WindowStart:   windowStart,
		WindowEnd:     windowEnd,
	}
	
	return &types.QueryValidatorStatisticsResponse{
		ValidatorAddress:          req.ValidatorAddress,
		TotalLiquidStaked:         liquidStaked,
		LiquidStakingPercentage:   liquidStakingPercentage,
		ActiveRecords:             activeRecords,
		TotalRecordsCreated:       totalRecordsCreated,
		RateLimitUsage:            &rateLimitUsage,
	}, nil
}

// ExchangeRate queries the current exchange rate for a validator
func (k Keeper) ExchangeRate(
	goCtx context.Context,
	req *types.QueryExchangeRateRequest,
) (*types.QueryExchangeRateResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	
	if req.ValidatorAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "validator address cannot be empty")
	}
	
	ctx := sdk.UnwrapSDKContext(goCtx)
	
	// Validate validator address
	_, err := sdk.ValAddressFromBech32(req.ValidatorAddress)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid validator address: %s", err))
	}
	
	// Get or initialize the exchange rate
	rate := k.GetOrInitExchangeRate(ctx, req.ValidatorAddress)
	
	// Get LST denom
	lstDenom := types.GetLSTDenom(req.ValidatorAddress)
	
	// Calculate native amount per 1 LST token
	nativeAmount := rate.Rate
	
	return &types.QueryExchangeRateResponse{
		ExchangeRate: rate,
		LstDenom:     lstDenom,
		NativeAmount: nativeAmount,
	}, nil
}

// AllExchangeRates queries all exchange rates
func (k Keeper) AllExchangeRates(
	goCtx context.Context,
	req *types.QueryAllExchangeRatesRequest,
) (*types.QueryAllExchangeRatesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	
	ctx := sdk.UnwrapSDKContext(goCtx)
	
	var exchangeRates []types.ExchangeRate
	var pageRes *query.PageResponse
	
	store := ctx.KVStore(k.storeKey)
	exchangeRateStore := prefix.NewStore(store, types.ExchangeRatePrefix)
	
	pageRes, err := query.Paginate(exchangeRateStore, req.Pagination, func(key []byte, value []byte) error {
		var rate types.ExchangeRate
		if err := k.cdc.Unmarshal(value, &rate); err != nil {
			return err
		}
		exchangeRates = append(exchangeRates, rate)
		return nil
	})
	
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	
	// Get global exchange rate if it exists
	globalRate, found := k.GetGlobalExchangeRate(ctx)
	var globalRatePtr *types.GlobalExchangeRate
	if found {
		globalRatePtr = &globalRate
	}
	
	return &types.QueryAllExchangeRatesResponse{
		ExchangeRates:      exchangeRates,
		Pagination:         pageRes,
		GlobalExchangeRate: globalRatePtr,
	}, nil
}