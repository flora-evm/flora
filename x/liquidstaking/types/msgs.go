package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgTokenizeShares      = "tokenize_shares"
	TypeMsgRedeemTokens        = "redeem_tokens"
	TypeMsgUpdateParams        = "update_params"
	TypeMsgUpdateExchangeRates = "update_exchange_rates"
)

var _ sdk.Msg = &MsgTokenizeShares{}

// NewMsgTokenizeShares creates a new MsgTokenizeShares instance
func NewMsgTokenizeShares(
	delegatorAddress string,
	validatorAddress string,
	shares sdk.Coin,
	ownerAddress string,
) *MsgTokenizeShares {
	return &MsgTokenizeShares{
		DelegatorAddress: delegatorAddress,
		ValidatorAddress: validatorAddress,
		Shares:           shares,
		OwnerAddress:     ownerAddress,
	}
}

// Route implements sdk.Msg
func (msg MsgTokenizeShares) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (msg MsgTokenizeShares) Type() string {
	return TypeMsgTokenizeShares
}

// GetSigners implements sdk.Msg
func (msg MsgTokenizeShares) GetSigners() []sdk.AccAddress {
	delegator, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{delegator}
}

// GetSignBytes implements sdk.Msg
func (msg MsgTokenizeShares) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements sdk.Msg
func (msg MsgTokenizeShares) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid delegator address (%s)", err)
	}

	_, err = sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid validator address (%s)", err)
	}

	if msg.OwnerAddress != "" {
		_, err = sdk.AccAddressFromBech32(msg.OwnerAddress)
		if err != nil {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address (%s)", err)
		}
	}

	if !msg.Shares.IsValid() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidCoins, "invalid shares amount")
	}

	if !msg.Shares.IsPositive() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidCoins, "shares amount must be positive")
	}

	// The denom should be "shares" or the validator-specific shares denom
	// This will be validated more thoroughly in the keeper
	
	return nil
}

// MsgRedeemTokens
var _ sdk.Msg = &MsgRedeemTokens{}

// NewMsgRedeemTokens creates a new MsgRedeemTokens instance
func NewMsgRedeemTokens(
	ownerAddress string,
	amount sdk.Coin,
) *MsgRedeemTokens {
	return &MsgRedeemTokens{
		OwnerAddress: ownerAddress,
		Amount:       amount,
	}
}

// Route implements sdk.Msg
func (msg MsgRedeemTokens) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (msg MsgRedeemTokens) Type() string {
	return TypeMsgRedeemTokens
}

// GetSigners implements sdk.Msg
func (msg MsgRedeemTokens) GetSigners() []sdk.AccAddress {
	owner, err := sdk.AccAddressFromBech32(msg.OwnerAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{owner}
}

// GetSignBytes implements sdk.Msg
func (msg MsgRedeemTokens) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements sdk.Msg
func (msg MsgRedeemTokens) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.OwnerAddress)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address (%s)", err)
	}

	if !msg.Amount.IsValid() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidCoins, "invalid amount")
	}

	if !msg.Amount.IsPositive() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidCoins, "amount must be positive")
	}

	// The denom should be a liquid staking token denom
	// This will be validated more thoroughly in the keeper
	
	return nil
}

// MsgUpdateParams implementation
var _ sdk.Msg = &MsgUpdateParams{}

// NewMsgUpdateParams creates a new MsgUpdateParams instance
func NewMsgUpdateParams(authority string, params ModuleParams) *MsgUpdateParams {
	return &MsgUpdateParams{
		Authority: authority,
		Params:    params,
	}
}

// Route implements sdk.Msg
func (msg MsgUpdateParams) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (msg MsgUpdateParams) Type() string {
	return TypeMsgUpdateParams
}

// GetSigners implements sdk.Msg
func (msg MsgUpdateParams) GetSigners() []sdk.AccAddress {
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{authority}
}

// GetSignBytes implements sdk.Msg
func (msg MsgUpdateParams) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements sdk.Msg
func (msg MsgUpdateParams) ValidateBasic() error {
	// Validate authority address
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
	}

	// Validate the parameters
	if err := msg.Params.Validate(); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid parameters: %s", err)
	}

	return nil
}

// MsgUpdateExchangeRates implementation
var _ sdk.Msg = &MsgUpdateExchangeRates{}

// NewMsgUpdateExchangeRates creates a new MsgUpdateExchangeRates instance
func NewMsgUpdateExchangeRates(updater string, validators []string) *MsgUpdateExchangeRates {
	return &MsgUpdateExchangeRates{
		Updater:    updater,
		Validators: validators,
	}
}

// Route implements sdk.Msg
func (msg MsgUpdateExchangeRates) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (msg MsgUpdateExchangeRates) Type() string {
	return TypeMsgUpdateExchangeRates
}

// GetSigners implements sdk.Msg
func (msg MsgUpdateExchangeRates) GetSigners() []sdk.AccAddress {
	updater, err := sdk.AccAddressFromBech32(msg.Updater)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{updater}
}

// GetSignBytes implements sdk.Msg
func (msg MsgUpdateExchangeRates) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements sdk.Msg
func (msg MsgUpdateExchangeRates) ValidateBasic() error {
	// Validate updater address
	_, err := sdk.AccAddressFromBech32(msg.Updater)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid updater address (%s)", err)
	}

	// Validate validator addresses if specified
	for _, validator := range msg.Validators {
		_, err := sdk.ValAddressFromBech32(validator)
		if err != nil {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid validator address %s: %s", validator, err)
		}
	}

	return nil
}