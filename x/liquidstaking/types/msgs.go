package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgTokenizeShares = "tokenize_shares"
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
const (
	TypeMsgRedeemTokens = "redeem_tokens"
)

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