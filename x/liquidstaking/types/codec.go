package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

// RegisterCodec registers the necessary x/liquidstaking interfaces and concrete types
// on the provided LegacyAmino codec. These types are used for Amino JSON serialization.
func RegisterCodec(cdc *codec.LegacyAmino) {
	// Messages
	cdc.RegisterConcrete(&MsgTokenizeShares{}, "liquidstaking/TokenizeShares", nil)
	cdc.RegisterConcrete(&MsgRedeemTokens{}, "liquidstaking/RedeemTokens", nil)
	cdc.RegisterConcrete(&MsgUpdateParams{}, "liquidstaking/UpdateParams", nil)
	cdc.RegisterConcrete(&MsgUpdateExchangeRates{}, "liquidstaking/UpdateExchangeRates", nil)
	
	// Admin messages - TODO: Add proto definitions
	// cdc.RegisterConcrete(&MsgEmergencyPause{}, "liquidstaking/EmergencyPause", nil)
	// cdc.RegisterConcrete(&MsgEmergencyUnpause{}, "liquidstaking/EmergencyUnpause", nil)
	// cdc.RegisterConcrete(&MsgSetValidatorWhitelist{}, "liquidstaking/SetValidatorWhitelist", nil)
	// cdc.RegisterConcrete(&MsgSetValidatorBlacklist{}, "liquidstaking/SetValidatorBlacklist", nil)
	
	// Governance proposals - TODO: Add proto definitions
	// cdc.RegisterConcrete(&UpdateParamsProposal{}, "liquidstaking/UpdateParamsProposal", nil)
	cdc.RegisterConcrete(&EmergencyPauseProposal{}, "liquidstaking/EmergencyPauseProposal", nil)
	cdc.RegisterConcrete(&UpdateValidatorCapProposal{}, "liquidstaking/UpdateValidatorCapProposal", nil)
}

// RegisterInterfaces registers the x/liquidstaking interfaces types with the interface registry
func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgTokenizeShares{},
		&MsgRedeemTokens{},
		&MsgUpdateParams{},
		&MsgUpdateExchangeRates{},
		// TODO: Add proto definitions for admin messages
		// &MsgEmergencyPause{},
		// &MsgEmergencyUnpause{},
		// &MsgSetValidatorWhitelist{},
		// &MsgSetValidatorBlacklist{},
	)
	
	registry.RegisterImplementations((*govtypes.Content)(nil),
		// TODO: Add proto definitions for governance proposals
		// &UpdateParamsProposal{},
		// &EmergencyPauseProposal{},
		// &UpdateValidatorCapProposal{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)

func init() {
	RegisterCodec(Amino)
	sdk.RegisterLegacyAminoCodec(Amino)
}