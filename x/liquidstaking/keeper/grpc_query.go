package keeper

import (
	"context"
	
	sdk "github.com/cosmos/cosmos-sdk/types"
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
	ctx := sdk.UnwrapSDKContext(c)
	
	record, found := k.GetTokenizationRecord(ctx, req.Id)
	if !found {
		return nil, status.Error(codes.NotFound, "tokenization record not found")
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
	
	// For now, return all records without pagination
	// Pagination will be implemented in later stages
	records := k.GetAllTokenizationRecords(ctx)
	
	return &types.QueryTokenizationRecordsResponse{
		Records: records,
	}, nil
}

// TokenizationRecordsByValidator returns all tokenization records for a validator
func (k Keeper) TokenizationRecordsByValidator(c context.Context, req *types.QueryTokenizationRecordsByValidatorRequest) (*types.QueryTokenizationRecordsByValidatorResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	
	records := k.GetTokenizationRecordsByValidator(ctx, req.ValidatorAddress)
	
	return &types.QueryTokenizationRecordsByValidatorResponse{
		Records: records,
	}, nil
}

// TokenizationRecordsByOwner returns all tokenization records for an owner
func (k Keeper) TokenizationRecordsByOwner(c context.Context, req *types.QueryTokenizationRecordsByOwnerRequest) (*types.QueryTokenizationRecordsByOwnerResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	
	records := k.GetTokenizationRecordsByOwner(ctx, req.OwnerAddress)
	
	return &types.QueryTokenizationRecordsByOwnerResponse{
		Records: records,
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
	ctx := sdk.UnwrapSDKContext(c)
	
	amount := k.GetValidatorLiquidStaked(ctx, req.ValidatorAddress)
	
	return &types.QueryValidatorLiquidStakedResponse{
		LiquidStaked: amount,
	}, nil
}