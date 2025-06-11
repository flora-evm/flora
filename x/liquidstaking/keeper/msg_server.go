package keeper

import (
	"context"
	"fmt"
	
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	
	"github.com/rollchains/flora/x/liquidstaking/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the liquid staking MsgServer interface
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

// TokenizeShares implements types.MsgServer
func (k msgServer) TokenizeShares(goCtx context.Context, msg *types.MsgTokenizeShares) (*types.MsgTokenizeSharesResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	
	// Validate the message
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}
	
	// Check if module is enabled
	params := k.GetParams(ctx)
	if !params.Enabled {
		return nil, types.ErrModuleDisabled
	}
	
	// Parse addresses
	delegatorAddr, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid delegator address: %s", err)
	}
	
	validatorAddr, err := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid validator address: %s", err)
	}
	
	// Determine owner address (defaults to delegator if not specified)
	ownerAddr := delegatorAddr
	if msg.OwnerAddress != "" {
		ownerAddr, err = sdk.AccAddressFromBech32(msg.OwnerAddress)
		if err != nil {
			return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid owner address: %s", err)
		}
	}
	
	// Get the delegation
	delegation, err := k.stakingKeeper.GetDelegation(ctx, delegatorAddr, validatorAddr)
	if err != nil {
		return nil, types.ErrDelegationNotFound
	}
	
	// Parse shares amount
	sharesToTokenize, err := math.LegacyNewDecFromStr(msg.Shares.Amount.String())
	if err != nil {
		return nil, sdkerrors.ErrInvalidRequest.Wrapf("invalid shares amount: %s", err)
	}
	
	// Validate shares amount
	if sharesToTokenize.LTE(math.LegacyZeroDec()) {
		return nil, sdkerrors.ErrInvalidRequest.Wrap("shares must be positive")
	}
	
	if sharesToTokenize.GT(delegation.Shares) {
		return nil, types.ErrInsufficientShares
	}
	
	// Get the validator
	validator, err := k.stakingKeeper.GetValidator(ctx, validatorAddr)
	if err != nil {
		return nil, sdkerrors.ErrNotFound.Wrapf("validator not found")
	}
	
	// Check if validator is valid for liquid staking
	if validator.IsJailed() {
		return nil, types.ErrInvalidValidator.Wrap("validator is jailed")
	}
	
	// Calculate the amount of tokens to mint
	// tokens = shares * (validator tokens / validator shares)
	tokensToMint := validator.TokensFromShares(sharesToTokenize).TruncateInt()
	
	// Validate liquid staking caps
	if err := k.CanTokenizeShares(ctx, msg.ValidatorAddress, tokensToMint); err != nil {
		return nil, err
	}
	
	// Generate the next tokenization record ID
	recordID := k.GetLastTokenizationRecordID(ctx) + 1
	k.SetLastTokenizationRecordID(ctx, recordID)
	
	// Generate the liquid staking token denom
	denom := types.GenerateLiquidStakingTokenDenom(msg.ValidatorAddress, recordID)
	
	// Create the tokenization record
	record := types.NewTokenizationRecordWithDenom(
		recordID,
		msg.ValidatorAddress,
		ownerAddr.String(),
		tokensToMint,
		denom,
	)
	
	// Validate and store the tokenization record
	if err := k.ValidateTokenizationRecord(ctx, record); err != nil {
		return nil, err
	}
	k.SetTokenizationRecordWithIndexes(ctx, record)
	
	// Unbond the shares from the delegation
	unbondedTokens, err := k.stakingKeeper.Unbond(ctx, delegatorAddr, validatorAddr, sharesToTokenize)
	if err != nil {
		return nil, sdkerrors.ErrInvalidRequest.Wrapf("failed to unbond shares: %s", err)
	}
	
	// Mint the liquid staking tokens
	mintCoins := sdk.NewCoins(sdk.NewCoin(denom, unbondedTokens))
	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, mintCoins); err != nil {
		return nil, sdkerrors.ErrInvalidRequest.Wrapf("failed to mint tokens: %s", err)
	}
	
	// Send the minted tokens to the owner
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, ownerAddr, mintCoins); err != nil {
		return nil, sdkerrors.ErrInvalidRequest.Wrapf("failed to send tokens: %s", err)
	}
	
	// Set token metadata
	metadata := types.GenerateLiquidStakingTokenMetadata(msg.ValidatorAddress, recordID)
	k.bankKeeper.SetDenomMetaData(ctx, metadata)
	
	// Update liquid staked amounts
	k.UpdateLiquidStakedAmounts(ctx, msg.ValidatorAddress, unbondedTokens, true)
	
	// Update denom index
	k.setTokenizationRecordDenomIndex(ctx, denom, recordID)
	
	// Emit events
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeTokenizeShares,
			sdk.NewAttribute(types.AttributeKeyDelegator, msg.DelegatorAddress),
			sdk.NewAttribute(types.AttributeKeyValidator, msg.ValidatorAddress),
			sdk.NewAttribute(types.AttributeKeyOwner, ownerAddr.String()),
			sdk.NewAttribute(types.AttributeKeyShares, sharesToTokenize.String()),
			sdk.NewAttribute(types.AttributeKeyDenom, denom),
			sdk.NewAttribute(types.AttributeKeyAmount, unbondedTokens.String()),
			sdk.NewAttribute(types.AttributeKeyRecordID, fmt.Sprintf("%d", recordID)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.DelegatorAddress),
		),
	})
	
	k.Logger(ctx).Info("tokenized shares",
		"delegator", msg.DelegatorAddress,
		"validator", msg.ValidatorAddress,
		"shares", sharesToTokenize,
		"tokens", unbondedTokens,
		"denom", denom,
		"record_id", recordID,
	)
	
	return &types.MsgTokenizeSharesResponse{
		Denom:    denom,
		Amount:   sdk.NewCoin(denom, unbondedTokens),
		RecordId: recordID,
	}, nil
}

// RedeemTokens implements types.MsgServer
func (k msgServer) RedeemTokens(goCtx context.Context, msg *types.MsgRedeemTokens) (*types.MsgRedeemTokensResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate the message
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	// Check if module is enabled
	params := k.GetParams(ctx)
	if !params.Enabled {
		return nil, types.ErrModuleDisabled
	}

	// Parse owner address
	ownerAddr, err := sdk.AccAddressFromBech32(msg.OwnerAddress)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid owner address: %s", err)
	}

	// Validate amount
	if !msg.Amount.IsValid() || msg.Amount.IsZero() {
		return nil, sdkerrors.ErrInvalidRequest.Wrap("invalid redeem amount")
	}

	// Get the tokenization record by denom
	recordID, found := k.getTokenizationRecordByDenom(ctx, msg.Amount.Denom)
	if !found {
		return nil, types.ErrTokenizationRecordNotFound.Wrapf("no record found for denom %s", msg.Amount.Denom)
	}

	record, found := k.GetTokenizationRecord(ctx, recordID)
	if !found {
		return nil, types.ErrTokenizationRecordNotFound
	}

	// Verify ownership - only the owner can redeem
	if record.Owner != msg.OwnerAddress {
		return nil, sdkerrors.ErrUnauthorized.Wrap("only the owner can redeem tokens")
	}

	// Check balance
	balance := k.bankKeeper.GetBalance(ctx, ownerAddr, msg.Amount.Denom)
	if balance.IsLT(msg.Amount) {
		return nil, sdkerrors.ErrInsufficientFunds.Wrapf("insufficient balance: has %s, needs %s", balance, msg.Amount)
	}

	// Get the validator
	valAddr, err := sdk.ValAddressFromBech32(record.Validator)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid validator address in record: %s", err)
	}

	validator, err := k.stakingKeeper.GetValidator(ctx, valAddr)
	if err != nil {
		return nil, sdkerrors.ErrNotFound.Wrapf("validator not found: %s", err)
	}

	// Calculate shares to restore based on the current exchange rate
	// shares = tokens / (validator tokens / validator shares)
	sharesToRestore, _ := validator.SharesFromTokens(msg.Amount.Amount)

	// Burn the liquid staking tokens
	burnCoins := sdk.NewCoins(msg.Amount)
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, ownerAddr, types.ModuleName, burnCoins); err != nil {
		return nil, sdkerrors.ErrInvalidRequest.Wrapf("failed to send tokens to module: %s", err)
	}

	if err := k.bankKeeper.BurnCoins(ctx, types.ModuleName, burnCoins); err != nil {
		return nil, sdkerrors.ErrInvalidRequest.Wrapf("failed to burn tokens: %s", err)
	}

	// Re-delegate the shares to the validator
	_, err = k.stakingKeeper.Delegate(ctx, ownerAddr, msg.Amount.Amount, stakingtypes.Unbonded, validator, true)
	if err != nil {
		return nil, sdkerrors.ErrInvalidRequest.Wrapf("failed to delegate: %s", err)
	}

	// Update the tokenization record
	if record.SharesTokenized.Sub(msg.Amount.Amount).IsZero() {
		// If all tokens are redeemed, delete the record
		k.DeleteTokenizationRecordWithIndexes(ctx, recordID)
		k.deleteTokenizationRecordDenomIndex(ctx, msg.Amount.Denom)
	} else {
		// Update the record with reduced amount
		record.SharesTokenized = record.SharesTokenized.Sub(msg.Amount.Amount)
		k.SetTokenizationRecordWithIndexes(ctx, record)
	}

	// Update liquid staked amounts
	k.UpdateLiquidStakedAmounts(ctx, record.Validator, msg.Amount.Amount, false)

	// Emit events
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeRedeemTokens,
			sdk.NewAttribute(types.AttributeKeyOwner, msg.OwnerAddress),
			sdk.NewAttribute(types.AttributeKeyValidator, record.Validator),
			sdk.NewAttribute(types.AttributeKeyDenom, msg.Amount.Denom),
			sdk.NewAttribute(types.AttributeKeyAmount, msg.Amount.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyShares, sharesToRestore.String()),
			sdk.NewAttribute(types.AttributeKeyRecordID, fmt.Sprintf("%d", recordID)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.OwnerAddress),
		),
	})

	k.Logger(ctx).Info("redeemed tokens",
		"owner", msg.OwnerAddress,
		"validator", record.Validator,
		"tokens", msg.Amount.Amount,
		"shares", sharesToRestore,
		"denom", msg.Amount.Denom,
		"record_id", recordID,
	)

	return &types.MsgRedeemTokensResponse{
		Shares:   sharesToRestore,
		RecordId: recordID,
	}, nil
}

// GenerateLiquidStakingTokenDenom generates a unique denom for liquid staking tokens
func GenerateLiquidStakingTokenDenom(validatorAddr string, recordID uint64) string {
	return fmt.Sprintf("liquidstake/%s/%d", validatorAddr, recordID)
}