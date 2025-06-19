package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Response types for admin messages
type MsgEmergencyPauseResponse struct{}
type MsgEmergencyUnpauseResponse struct{}
type MsgSetValidatorWhitelistResponse struct{}
type MsgSetValidatorBlacklistResponse struct{}

const (
	TypeMsgEmergencyPause          = "emergency_pause"
	TypeMsgEmergencyUnpause        = "emergency_unpause"
	TypeMsgSetValidatorWhitelist   = "set_validator_whitelist"
	TypeMsgSetValidatorBlacklist   = "set_validator_blacklist"
)

// TODO: Uncomment after adding proto definitions
// var (
// 	_ sdk.Msg = &MsgEmergencyPause{}
// 	_ sdk.Msg = &MsgEmergencyUnpause{}
// 	_ sdk.Msg = &MsgSetValidatorWhitelist{}
// 	_ sdk.Msg = &MsgSetValidatorBlacklist{}
// )

// MsgEmergencyPause defines a message to pause the module
type MsgEmergencyPause struct {
	Authority string `protobuf:"bytes,1,opt,name=authority,proto3" json:"authority,omitempty"`
	Reason    string `protobuf:"bytes,2,opt,name=reason,proto3" json:"reason,omitempty"`
	Duration  int64  `protobuf:"varint,3,opt,name=duration,proto3" json:"duration,omitempty"`
}

// NewMsgEmergencyPause creates a new MsgEmergencyPause
func NewMsgEmergencyPause(authority, reason string, duration int64) *MsgEmergencyPause {
	return &MsgEmergencyPause{
		Authority: authority,
		Reason:    reason,
		Duration:  duration,
	}
}

// Route returns the message route
func (msg *MsgEmergencyPause) Route() string { return RouterKey }

// Type returns the message type
func (msg *MsgEmergencyPause) Type() string { return TypeMsgEmergencyPause }

// GetSigners returns the expected signers
func (msg *MsgEmergencyPause) GetSigners() []sdk.AccAddress {
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{authority}
}

// GetSignBytes returns the message bytes to sign over
// TODO: Uncomment after adding proto definitions
// func (msg *MsgEmergencyPause) GetSignBytes() []byte {
// 	bz := ModuleCdc.MustMarshalJSON(msg)
// 	return sdk.MustSortJSON(bz)
// }

// ValidateBasic performs basic validation
func (msg *MsgEmergencyPause) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
	}
	if msg.Reason == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "reason cannot be empty")
	}
	if msg.Duration < 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "duration cannot be negative")
	}
	return nil
}

// MsgEmergencyUnpause defines a message to unpause the module
type MsgEmergencyUnpause struct {
	Authority string `protobuf:"bytes,1,opt,name=authority,proto3" json:"authority,omitempty"`
}

// NewMsgEmergencyUnpause creates a new MsgEmergencyUnpause
func NewMsgEmergencyUnpause(authority string) *MsgEmergencyUnpause {
	return &MsgEmergencyUnpause{
		Authority: authority,
	}
}

// Route returns the message route
func (msg *MsgEmergencyUnpause) Route() string { return RouterKey }

// Type returns the message type
func (msg *MsgEmergencyUnpause) Type() string { return TypeMsgEmergencyUnpause }

// GetSigners returns the expected signers
func (msg *MsgEmergencyUnpause) GetSigners() []sdk.AccAddress {
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{authority}
}

// GetSignBytes returns the message bytes to sign over
// TODO: Uncomment after adding proto definitions
// func (msg *MsgEmergencyUnpause) GetSignBytes() []byte {
// 	bz := ModuleCdc.MustMarshalJSON(msg)
// 	return sdk.MustSortJSON(bz)
// }

// ValidateBasic performs basic validation
func (msg *MsgEmergencyUnpause) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
	}
	return nil
}

// MsgSetValidatorWhitelist defines a message to set the validator whitelist
type MsgSetValidatorWhitelist struct {
	Authority  string   `protobuf:"bytes,1,opt,name=authority,proto3" json:"authority,omitempty"`
	Validators []string `protobuf:"bytes,2,rep,name=validators,proto3" json:"validators,omitempty"`
}

// NewMsgSetValidatorWhitelist creates a new MsgSetValidatorWhitelist
func NewMsgSetValidatorWhitelist(authority string, validators []string) *MsgSetValidatorWhitelist {
	return &MsgSetValidatorWhitelist{
		Authority:  authority,
		Validators: validators,
	}
}

// Route returns the message route
func (msg *MsgSetValidatorWhitelist) Route() string { return RouterKey }

// Type returns the message type
func (msg *MsgSetValidatorWhitelist) Type() string { return TypeMsgSetValidatorWhitelist }

// GetSigners returns the expected signers
func (msg *MsgSetValidatorWhitelist) GetSigners() []sdk.AccAddress {
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{authority}
}

// GetSignBytes returns the message bytes to sign over
// TODO: Uncomment after adding proto definitions
// func (msg *MsgSetValidatorWhitelist) GetSignBytes() []byte {
// 	bz := ModuleCdc.MustMarshalJSON(msg)
// 	return sdk.MustSortJSON(bz)
// }

// ValidateBasic performs basic validation
func (msg *MsgSetValidatorWhitelist) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
	}
	
	// Validate all validator addresses
	for _, val := range msg.Validators {
		if _, err := sdk.ValAddressFromBech32(val); err != nil {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid validator address %s: %s", val, err)
		}
	}
	
	return nil
}

// MsgSetValidatorBlacklist defines a message to set the validator blacklist
type MsgSetValidatorBlacklist struct {
	Authority  string   `protobuf:"bytes,1,opt,name=authority,proto3" json:"authority,omitempty"`
	Validators []string `protobuf:"bytes,2,rep,name=validators,proto3" json:"validators,omitempty"`
}

// NewMsgSetValidatorBlacklist creates a new MsgSetValidatorBlacklist
func NewMsgSetValidatorBlacklist(authority string, validators []string) *MsgSetValidatorBlacklist {
	return &MsgSetValidatorBlacklist{
		Authority:  authority,
		Validators: validators,
	}
}

// Route returns the message route
func (msg *MsgSetValidatorBlacklist) Route() string { return RouterKey }

// Type returns the message type
func (msg *MsgSetValidatorBlacklist) Type() string { return TypeMsgSetValidatorBlacklist }

// GetSigners returns the expected signers
func (msg *MsgSetValidatorBlacklist) GetSigners() []sdk.AccAddress {
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{authority}
}

// GetSignBytes returns the message bytes to sign over
// TODO: Uncomment after adding proto definitions
// func (msg *MsgSetValidatorBlacklist) GetSignBytes() []byte {
// 	bz := ModuleCdc.MustMarshalJSON(msg)
// 	return sdk.MustSortJSON(bz)
// }

// ValidateBasic performs basic validation
func (msg *MsgSetValidatorBlacklist) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
	}
	
	// Validate all validator addresses
	for _, val := range msg.Validators {
		if _, err := sdk.ValAddressFromBech32(val); err != nil {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid validator address %s: %s", val, err)
		}
	}
	
	return nil
}