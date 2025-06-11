package liquidstaking_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/stretchr/testify/require"

	"github.com/rollchains/flora/x/liquidstaking"
)

func TestModuleBasics(t *testing.T) {
	module := liquidstaking.AppModuleBasic{}
	
	// Test module name
	require.Equal(t, "liquidstaking", module.Name())
	
	// Test that RegisterLegacyAminoCodec doesn't panic
	cdc := codec.NewLegacyAmino()
	require.NotPanics(t, func() {
		module.RegisterLegacyAminoCodec(cdc)
	})
}