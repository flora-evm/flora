package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// RegisterCodec registers the necessary x/liquidstaking interfaces and concrete types
// on the provided LegacyAmino codec. These types are used for Amino JSON serialization.
func RegisterCodec(cdc *codec.LegacyAmino) {
	// Register concrete types here when we have messages in later stages
}

// RegisterInterfaces registers the x/liquidstaking interfaces types with the interface registry
func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	// Register message interfaces here when we have messages in later stages
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)

func init() {
	RegisterCodec(Amino)
	sdk.RegisterLegacyAminoCodec(Amino)
}