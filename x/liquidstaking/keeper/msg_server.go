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
	if err := k.ValidateModuleEnabled(ctx); err != nil {
		return nil, err
	}
	
	// Check if module is paused
	if err := k.RequireNotPaused(ctx); err != nil {
		return nil, err
	}
	
	// Parse addresses
	delegatorAddr, err := ParseAndValidateAddress(msg.DelegatorAddress, "delegator")
	if err != nil {
		return nil, err
	}
	
	validatorAddr, err := ParseAndValidateValidatorAddress(msg.ValidatorAddress)
	if err != nil {
		return nil, err
	}
	
	// Determine owner address (defaults to delegator if not specified)
	ownerAddr := delegatorAddr
	if msg.OwnerAddress != "" {
		ownerAddr, err = ParseAndValidateAddress(msg.OwnerAddress, "owner")
		if err != nil {
			return nil, err
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
	
	// Check if validator is allowed (whitelist/blacklist)
	if !k.IsValidatorAllowed(ctx, validatorAddr) {
		return nil, types.ErrInvalidValidator.Wrap("validator not allowed for liquid staking")
	}
	
	// Calculate the amount of native tokens represented by the shares
	// nativeTokens = shares * (validator tokens / validator shares)
	nativeTokens := validator.TokensFromShares(sharesToTokenize).TruncateInt()
	
	// Apply exchange rate to determine LST tokens to mint
	// LST tokens = native tokens / exchange rate
	lstTokensToMint, err := k.ApplyExchangeRate(ctx, msg.ValidatorAddress, nativeTokens)
	if err != nil {
		return nil, sdkerrors.ErrInvalidRequest.Wrapf("failed to apply exchange rate: %s", err)
	}
	
	// Validate liquid staking caps (using native token amount for cap validation)
	if err := k.CanTokenizeShares(ctx, msg.ValidatorAddress, nativeTokens); err != nil {
		return nil, err
	}
	
	// Enforce rate limits (using native token amount)
	if err := k.EnforceTokenizationRateLimits(ctx, msg.ValidatorAddress, msg.DelegatorAddress, nativeTokens); err != nil {
		return nil, err
	}
	
	// Call pre-tokenization hook
	hooks := k.GetHooks()
	if err := hooks.PreTokenizeShares(ctx, delegatorAddr, validatorAddr, ownerAddr, sharesToTokenize); err != nil {
		return nil, sdkerrors.ErrInvalidRequest.Wrapf("pre-tokenization hook failed: %s", err)
	}
	
	// Generate the next tokenization record ID
	recordID := k.GetLastTokenizationRecordID(ctx) + 1
	k.SetLastTokenizationRecordID(ctx, recordID)
	
	// Generate the liquid staking token denom
	denom := types.GenerateLiquidStakingTokenDenom(msg.ValidatorAddress, recordID)
	
	// Create the tokenization record (storing the native token amount for record keeping)
	record := types.NewTokenizationRecordWithDenom(
		recordID,
		msg.ValidatorAddress,
		ownerAddr.String(),
		nativeTokens,
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
	
	// Mint the liquid staking tokens (using the exchange rate adjusted amount)
	mintCoins := sdk.NewCoins(sdk.NewCoin(denom, lstTokensToMint))
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
	
	// Update tokenization activity for rate limiting
	k.UpdateTokenizationActivity(ctx, msg.ValidatorAddress, msg.DelegatorAddress, unbondedTokens)
	
	// Update denom index
	k.setTokenizationRecordDenomIndex(ctx, denom, recordID)
	
	// Emit typed events
	types.EmitTokenizeSharesEvent(ctx, types.TokenizeSharesEvent{
		Delegator:    msg.DelegatorAddress,
		Validator:    msg.ValidatorAddress,
		Owner:        ownerAddr.String(),
		SharesAmount: sharesToTokenize.String(),
		TokensMinted: lstTokensToMint.String(),
		Denom:        denom,
		RecordID:     recordID,
	})
	
	// Emit record created event
	ctx.EventManager().EmitEvent(types.TokenizationRecordCreatedEvent{
		RecordID:        recordID,
		Validator:       msg.ValidatorAddress,
		Owner:           ownerAddr.String(),
		SharesTokenized: unbondedTokens.String(),
		Denom:           denom,
	}.ToEvent())
	
	// Call hooks for record creation
	hooks.OnTokenizationRecordCreated(ctx, record)
	
	// Call post-tokenization hook (passing LST tokens minted)
	hooks.PostTokenizeShares(ctx, delegatorAddr, validatorAddr, ownerAddr, sharesToTokenize, lstTokensToMint, denom, recordID)
	
	k.Logger(ctx).Info("tokenized shares",
		"delegator", msg.DelegatorAddress,
		"validator", msg.ValidatorAddress,
		"shares", sharesToTokenize,
		"native_tokens", unbondedTokens,
		"lst_tokens", lstTokensToMint,
		"denom", denom,
		"record_id", recordID,
	)
	
	return &types.MsgTokenizeSharesResponse{
		Denom:    denom,
		Amount:   sdk.NewCoin(denom, lstTokensToMint),
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
	if err := k.ValidateModuleEnabled(ctx); err != nil {
		return nil, err
	}
	
	// Check if module is paused
	if err := k.RequireNotPaused(ctx); err != nil {
		return nil, err
	}

	// Parse owner address
	ownerAddr, err := ParseAndValidateAddress(msg.OwnerAddress, "owner")
	if err != nil {
		return nil, err
	}

	// Validate amount
	if err := ValidatePositiveAmount(msg.Amount); err != nil {
		return nil, err
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
	valAddr, err := ParseAndValidateValidatorAddress(record.Validator)
	if err != nil {
		return nil, err
	}

	validator, err := k.stakingKeeper.GetValidator(ctx, valAddr)
	if err != nil {
		return nil, sdkerrors.ErrNotFound.Wrapf("validator not found: %s", err)
	}

	// Apply inverse exchange rate to calculate native tokens from LST tokens
	// Native tokens = LST tokens * exchange rate
	nativeTokens, err := k.ApplyInverseExchangeRate(ctx, record.Validator, msg.Amount.Amount)
	if err != nil {
		return nil, sdkerrors.ErrInvalidRequest.Wrapf("failed to apply exchange rate: %s", err)
	}
	
	// Calculate shares to restore based on the native tokens
	// shares = native tokens / (validator tokens / validator shares)
	sharesToRestore, _ := validator.SharesFromTokens(nativeTokens)

	// Call pre-redemption hook
	hooks := k.GetHooks()
	if err := hooks.PreRedeemTokens(ctx, ownerAddr, msg.Amount, recordID); err != nil {
		return nil, sdkerrors.ErrInvalidRequest.Wrapf("pre-redemption hook failed: %s", err)
	}

	// Burn the liquid staking tokens
	burnCoins := sdk.NewCoins(msg.Amount)
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, ownerAddr, types.ModuleName, burnCoins); err != nil {
		return nil, sdkerrors.ErrInvalidRequest.Wrapf("failed to send tokens to module: %s", err)
	}

	if err := k.bankKeeper.BurnCoins(ctx, types.ModuleName, burnCoins); err != nil {
		return nil, sdkerrors.ErrInvalidRequest.Wrapf("failed to burn tokens: %s", err)
	}

	// Re-delegate the native tokens to the validator
	_, err = k.stakingKeeper.Delegate(ctx, ownerAddr, nativeTokens, stakingtypes.Unbonded, validator, true)
	if err != nil {
		return nil, sdkerrors.ErrInvalidRequest.Wrapf("failed to delegate: %s", err)
	}

	// Update the tokenization record (tracking native tokens)
	oldAmount := record.SharesTokenized
	oldRecord := record // Save for hook
	if record.SharesTokenized.Sub(nativeTokens).IsZero() {
		// If all tokens are redeemed, delete the record
		k.DeleteTokenizationRecordWithIndexes(ctx, recordID)
		k.deleteTokenizationRecordDenomIndex(ctx, msg.Amount.Denom)
		
		// Emit record deleted event
		ctx.EventManager().EmitEvent(types.TokenizationRecordDeletedEvent{
			RecordID:  recordID,
			Validator: record.Validator,
			Owner:     record.Owner,
			Denom:     record.Denom,
		}.ToEvent())
		
		// Call hook for record deletion
		hooks.OnTokenizationRecordDeleted(ctx, record)
	} else {
		// Update the record with reduced amount (native tokens)
		record.SharesTokenized = record.SharesTokenized.Sub(nativeTokens)
		k.SetTokenizationRecordWithIndexes(ctx, record)
		
		// Emit record updated event
		ctx.EventManager().EmitEvent(types.TokenizationRecordUpdatedEvent{
			RecordID:           recordID,
			OldSharesTokenized: oldAmount.String(),
			NewSharesTokenized: record.SharesTokenized.String(),
		}.ToEvent())
		
		// Call hook for record update
		hooks.OnTokenizationRecordUpdated(ctx, oldRecord, record)
	}

	// Update liquid staked amounts (using native tokens)
	k.UpdateLiquidStakedAmounts(ctx, record.Validator, nativeTokens, false)

	// Emit typed redemption event (showing LST tokens burned)
	types.EmitRedeemTokensEvent(ctx, types.RedeemTokensEvent{
		Owner:          msg.OwnerAddress,
		Validator:      record.Validator,
		TokensBurned:   msg.Amount.Amount.String(),
		SharesRestored: sharesToRestore.String(),
		Denom:          msg.Amount.Denom,
		RecordID:       recordID,
	})

	// Call post-redemption hook
	hooks.PostRedeemTokens(ctx, ownerAddr, valAddr, msg.Amount, sharesToRestore, recordID)

	k.Logger(ctx).Info("redeemed tokens",
		"owner", msg.OwnerAddress,
		"validator", record.Validator,
		"lst_tokens", msg.Amount.Amount,
		"native_tokens", nativeTokens,
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

// UpdateParams updates the module parameters via governance
func (k msgServer) UpdateParams(goCtx context.Context, msg *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	
	// Check if the authority is the governance module account
	if msg.Authority != k.authority {
		return nil, sdkerrors.ErrUnauthorized.Wrapf("invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}
	
	// Validate the new parameters
	if err := msg.Params.Validate(); err != nil {
		return nil, sdkerrors.ErrInvalidRequest.Wrapf("invalid parameters: %s", err)
	}
	
	// Set the new parameters
	k.SetParams(ctx, msg.Params)
	
	// Emit event for parameter update
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeParameterUpdate,
			sdk.NewAttribute(types.AttributeKeyAction, "update_params"),
		),
	)
	
	// Call hooks if any module is listening for parameter updates
	k.GetHooks().OnParametersUpdated(ctx, msg.Params)
	
	return &types.MsgUpdateParamsResponse{}, nil
}

// UpdateExchangeRates implements types.MsgServer
func (k msgServer) UpdateExchangeRates(goCtx context.Context, msg *types.MsgUpdateExchangeRates) (*types.MsgUpdateExchangeRatesResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	
	// Validate the message
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}
	
	// Check if module is enabled
	if err := k.ValidateModuleEnabled(ctx); err != nil {
		return nil, err
	}
	
	// Check if module is paused
	if err := k.RequireNotPaused(ctx); err != nil {
		return nil, err
	}
	
	// TODO: Add authorization check - for now only authority can update
	// In Stage 15, this will be automated
	if msg.Updater != k.authority {
		return nil, sdkerrors.ErrUnauthorized.Wrapf("only authority can update exchange rates: expected %s, got %s", k.authority, msg.Updater)
	}
	
	var updatedRates []types.ExchangeRateUpdate
	
	// If no validators specified, update all validators with LST tokens
	if len(msg.Validators) == 0 {
		// Update all exchange rates
		if err := k.UpdateAllExchangeRates(ctx); err != nil {
			return nil, sdkerrors.ErrInvalidRequest.Wrapf("failed to update all exchange rates: %s", err)
		}
		
		// Collect all updated rates for response
		k.IterateExchangeRates(ctx, func(rate types.ExchangeRate) bool {
			updatedRates = append(updatedRates, types.ExchangeRateUpdate{
				ValidatorAddress: rate.ValidatorAddress,
				OldRate:         math.LegacyOneDec(), // We don't track old rate in bulk update
				NewRate:         rate.Rate,
			})
			return false
		})
	} else {
		// Update specific validators
		for _, validatorAddr := range msg.Validators {
			// Get old rate
			oldRate := math.LegacyOneDec()
			if existingRate, found := k.GetExchangeRate(ctx, validatorAddr); found {
				oldRate = existingRate.Rate
			}
			
			// Update the rate
			if err := k.UpdateExchangeRate(ctx, validatorAddr); err != nil {
				// Log error but continue with other validators
				k.Logger(ctx).Error("failed to update exchange rate", "validator", validatorAddr, "error", err)
				continue
			}
			
			// Get new rate
			newRate, found := k.GetExchangeRate(ctx, validatorAddr)
			if found {
				updatedRates = append(updatedRates, types.ExchangeRateUpdate{
					ValidatorAddress: validatorAddr,
					OldRate:         oldRate,
					NewRate:         newRate.Rate,
				})
			}
		}
	}
	
	// Log the update
	k.Logger(ctx).Info("exchange rates updated",
		"updater", msg.Updater,
		"validators_count", len(updatedRates),
	)
	
	return &types.MsgUpdateExchangeRatesResponse{
		UpdatedRates: updatedRates,
	}, nil
}