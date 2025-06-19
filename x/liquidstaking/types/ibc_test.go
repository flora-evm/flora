package types_test

import (
	"encoding/json"
	"testing"
	"time"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	"github.com/rollchains/flora/x/liquidstaking/types"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
)

func TestLiquidStakingTokenPacketData(t *testing.T) {
	validMetadata := &types.LiquidStakingMetadata{
		ValidatorAddress: "floravaloper1abc",
		RecordId:         1,
		SharesAmount:     math.LegacyNewDec(1000000),
		SourceChainId:    "flora-1",
		CreatedAt:        time.Now().UTC().Format(time.RFC3339),
	}

	testCases := []struct {
		name      string
		packet    types.LiquidStakingTokenPacketData
		expectErr bool
	}{
		{
			name: "valid packet without metadata",
			packet: types.NewLiquidStakingTokenPacketData(
				"flora",
				"1000000",
				"flora1sender",
				"flora1receiver",
				"test memo",
				nil,
			),
			expectErr: false,
		},
		{
			name: "valid packet with metadata",
			packet: types.NewLiquidStakingTokenPacketData(
				"liquidstake/floravaloper1abc/1",
				"1000000",
				"flora1sender",
				"flora1receiver",
				"",
				validMetadata,
			),
			expectErr: false,
		},
		{
			name: "empty denom",
			packet: types.NewLiquidStakingTokenPacketData(
				"",
				"1000000",
				"flora1sender",
				"flora1receiver",
				"",
				nil,
			),
			expectErr: true,
		},
		{
			name: "invalid amount",
			packet: types.NewLiquidStakingTokenPacketData(
				"flora",
				"-1000000",
				"flora1sender",
				"flora1receiver",
				"",
				nil,
			),
			expectErr: true,
		},
		{
			name: "empty sender",
			packet: types.NewLiquidStakingTokenPacketData(
				"flora",
				"1000000",
				"",
				"flora1receiver",
				"",
				nil,
			),
			expectErr: true,
		},
		{
			name: "empty receiver",
			packet: types.NewLiquidStakingTokenPacketData(
				"flora",
				"1000000",
				"flora1sender",
				"",
				"",
				nil,
			),
			expectErr: true,
		},
		{
			name: "invalid metadata - empty validator",
			packet: types.NewLiquidStakingTokenPacketData(
				"liquidstake/floravaloper1abc/1",
				"1000000",
				"flora1sender",
				"flora1receiver",
				"",
				&types.LiquidStakingMetadata{
					ValidatorAddress: "",
					RecordId:         1,
					SharesAmount:     math.LegacyNewDec(1000000),
					SourceChainId:    "flora-1",
					CreatedAt:        time.Now().UTC().Format(time.RFC3339),
				},
			),
			expectErr: true,
		},
		{
			name: "invalid metadata - zero record ID",
			packet: types.NewLiquidStakingTokenPacketData(
				"liquidstake/floravaloper1abc/1",
				"1000000",
				"flora1sender",
				"flora1receiver",
				"",
				&types.LiquidStakingMetadata{
					ValidatorAddress: "floravaloper1abc",
					RecordId:         0,
					SharesAmount:     math.LegacyNewDec(1000000),
					SourceChainId:    "flora-1",
					CreatedAt:        time.Now().UTC().Format(time.RFC3339),
				},
			),
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.packet.ValidateBasic()
			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestLiquidStakingMetadata_ValidateBasic(t *testing.T) {
	testCases := []struct {
		name      string
		metadata  types.LiquidStakingMetadata
		expectErr bool
	}{
		{
			name: "valid metadata",
			metadata: types.LiquidStakingMetadata{
				ValidatorAddress: "floravaloper1abc",
				RecordId:         1,
				SharesAmount:     math.LegacyNewDec(1000000),
				SourceChainId:    "flora-1",
				CreatedAt:        time.Now().UTC().Format(time.RFC3339),
			},
			expectErr: false,
		},
		{
			name: "empty validator address",
			metadata: types.LiquidStakingMetadata{
				ValidatorAddress: "",
				RecordId:         1,
				SharesAmount:     math.LegacyNewDec(1000000),
				SourceChainId:    "flora-1",
				CreatedAt:        time.Now().UTC().Format(time.RFC3339),
			},
			expectErr: true,
		},
		{
			name: "zero record ID",
			metadata: types.LiquidStakingMetadata{
				ValidatorAddress: "floravaloper1abc",
				RecordId:         0,
				SharesAmount:     math.LegacyNewDec(1000000),
				SourceChainId:    "flora-1",
				CreatedAt:        time.Now().UTC().Format(time.RFC3339),
			},
			expectErr: true,
		},
		{
			name: "negative shares amount",
			metadata: types.LiquidStakingMetadata{
				ValidatorAddress: "floravaloper1abc",
				RecordId:         1,
				SharesAmount:     math.LegacyNewDec(-1000000),
				SourceChainId:    "flora-1",
				CreatedAt:        time.Now().UTC().Format(time.RFC3339),
			},
			expectErr: true,
		},
		{
			name: "empty source chain ID",
			metadata: types.LiquidStakingMetadata{
				ValidatorAddress: "floravaloper1abc",
				RecordId:         1,
				SharesAmount:     math.LegacyNewDec(1000000),
				SourceChainId:    "",
				CreatedAt:        time.Now().UTC().Format(time.RFC3339),
			},
			expectErr: true,
		},
		{
			name: "empty created at",
			metadata: types.LiquidStakingMetadata{
				ValidatorAddress: "floravaloper1abc",
				RecordId:         1,
				SharesAmount:     math.LegacyNewDec(1000000),
				SourceChainId:    "flora-1",
				CreatedAt:        "",
			},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.metadata.ValidateBasic()
			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestPacketDataSerialization(t *testing.T) {
	metadata := &types.LiquidStakingMetadata{
		ValidatorAddress: "floravaloper1abc",
		RecordId:         1,
		SharesAmount:     math.LegacyNewDec(1000000),
		SourceChainId:    "flora-1",
		CreatedAt:        time.Now().UTC().Format(time.RFC3339),
	}

	packet := types.NewLiquidStakingTokenPacketData(
		"liquidstake/floravaloper1abc/1",
		"1000000",
		"flora1sender",
		"flora1receiver",
		"test memo",
		metadata,
	)

	// Test GetBytes
	bz := packet.GetBytes()
	require.NotEmpty(t, bz)

	// Test unmarshal
	unmarshaled := types.MustUnmarshalLiquidStakingTokenPacketData(bz)
	require.Equal(t, packet.Denom, unmarshaled.Denom)
	require.Equal(t, packet.Amount, unmarshaled.Amount)
	require.Equal(t, packet.Sender, unmarshaled.Sender)
	require.Equal(t, packet.Receiver, unmarshaled.Receiver)
	require.Equal(t, packet.Memo, unmarshaled.Memo)
	require.NotNil(t, unmarshaled.LiquidStakingMetadata)
	require.Equal(t, metadata.ValidatorAddress, unmarshaled.LiquidStakingMetadata.ValidatorAddress)
	require.Equal(t, metadata.RecordId, unmarshaled.LiquidStakingMetadata.RecordId)

	// Test UnmarshalLiquidStakingTokenPacketData with error
	_, err := types.UnmarshalLiquidStakingTokenPacketData([]byte("invalid json"))
	require.Error(t, err)
}

func TestConvertToTransferPacket(t *testing.T) {
	metadata := &types.LiquidStakingMetadata{
		ValidatorAddress: "floravaloper1abc",
		RecordId:         1,
		SharesAmount:     math.LegacyNewDec(1000000),
		SourceChainId:    "flora-1",
		CreatedAt:        time.Now().UTC().Format(time.RFC3339),
	}

	testCases := []struct {
		name     string
		packet   types.LiquidStakingTokenPacketData
		checkMemo bool
	}{
		{
			name: "packet without metadata",
			packet: types.NewLiquidStakingTokenPacketData(
				"flora",
				"1000000",
				"flora1sender",
				"flora1receiver",
				"original memo",
				nil,
			),
			checkMemo: false,
		},
		{
			name: "packet with metadata and no memo",
			packet: types.NewLiquidStakingTokenPacketData(
				"liquidstake/floravaloper1abc/1",
				"1000000",
				"flora1sender",
				"flora1receiver",
				"",
				metadata,
			),
			checkMemo: true,
		},
		{
			name: "packet with metadata and existing memo",
			packet: types.NewLiquidStakingTokenPacketData(
				"liquidstake/floravaloper1abc/1",
				"1000000",
				"flora1sender",
				"flora1receiver",
				"original memo",
				metadata,
			),
			checkMemo: true,
		},
		{
			name: "packet with metadata and JSON memo",
			packet: types.NewLiquidStakingTokenPacketData(
				"liquidstake/floravaloper1abc/1",
				"1000000",
				"flora1sender",
				"flora1receiver",
				`{"key": "value"}`,
				metadata,
			),
			checkMemo: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			transferPacket := tc.packet.ConvertToTransferPacket()
			
			require.Equal(t, tc.packet.Denom, transferPacket.Denom)
			require.Equal(t, tc.packet.Amount, transferPacket.Amount)
			require.Equal(t, tc.packet.Sender, transferPacket.Sender)
			require.Equal(t, tc.packet.Receiver, transferPacket.Receiver)
			
			if tc.checkMemo && tc.packet.LiquidStakingMetadata != nil {
				// Verify metadata is in memo
				var memoData map[string]interface{}
				err := json.Unmarshal([]byte(transferPacket.Memo), &memoData)
				require.NoError(t, err)
				require.Contains(t, memoData, "liquid_staking_metadata")
			}
		})
	}
}

func TestExtractFromTransferPacket(t *testing.T) {
	metadata := &types.LiquidStakingMetadata{
		ValidatorAddress: "floravaloper1abc",
		RecordId:         1,
		SharesAmount:     math.LegacyNewDec(1000000),
		SourceChainId:    "flora-1",
		CreatedAt:        time.Now().UTC().Format(time.RFC3339),
	}

	// Create memo with metadata
	memoData := map[string]interface{}{
		"liquid_staking_metadata": metadata,
		"text": "additional info",
	}
	memoJSON, _ := json.Marshal(memoData)

	testCases := []struct {
		name           string
		transferPacket transfertypes.FungibleTokenPacketData
		expectMetadata bool
	}{
		{
			name: "packet with liquid staking metadata",
			transferPacket: transfertypes.FungibleTokenPacketData{
				Denom:    "liquidstake/floravaloper1abc/1",
				Amount:   "1000000",
				Sender:   "flora1sender",
				Receiver: "flora1receiver",
				Memo:     string(memoJSON),
			},
			expectMetadata: true,
		},
		{
			name: "packet without metadata",
			transferPacket: transfertypes.FungibleTokenPacketData{
				Denom:    "flora",
				Amount:   "1000000",
				Sender:   "flora1sender",
				Receiver: "flora1receiver",
				Memo:     "simple text memo",
			},
			expectMetadata: false,
		},
		{
			name: "packet with empty memo",
			transferPacket: transfertypes.FungibleTokenPacketData{
				Denom:    "flora",
				Amount:   "1000000",
				Sender:   "flora1sender",
				Receiver: "flora1receiver",
				Memo:     "",
			},
			expectMetadata: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			lstPacket, err := types.ExtractFromTransferPacket(tc.transferPacket)
			require.NoError(t, err)
			
			require.Equal(t, tc.transferPacket.Denom, lstPacket.Denom)
			require.Equal(t, tc.transferPacket.Amount, lstPacket.Amount)
			require.Equal(t, tc.transferPacket.Sender, lstPacket.Sender)
			require.Equal(t, tc.transferPacket.Receiver, lstPacket.Receiver)
			
			if tc.expectMetadata {
				require.NotNil(t, lstPacket.LiquidStakingMetadata)
				require.Equal(t, metadata.ValidatorAddress, lstPacket.LiquidStakingMetadata.ValidatorAddress)
				require.Equal(t, metadata.RecordId, lstPacket.LiquidStakingMetadata.RecordId)
			} else {
				require.Nil(t, lstPacket.LiquidStakingMetadata)
			}
		})
	}
}

func TestExtractLiquidStakingMetadata(t *testing.T) {
	record := types.TokenizationRecord{
		Id:               1,
		ValidatorAddress: "floravaloper1abc",
		Owner:            "flora1xyz",
		SharesDenomination: "shares/floravaloper1abc",
		LiquidStakingTokenDenom: "liquidstake/floravaloper1abc/1",
		SharesAmount:     math.LegacyNewDec(1000000),
		Status:           types.TokenizationRecord_ACTIVE,
		CreatedAt:        time.Now().UTC(),
	}

	testCases := []struct {
		name           string
		denom          string
		expectMetadata bool
	}{
		{
			name:           "valid liquid staking token",
			denom:          "liquidstake/floravaloper1abc/1",
			expectMetadata: true,
		},
		{
			name:           "regular token",
			denom:          "flora",
			expectMetadata: false,
		},
		{
			name:           "invalid liquid staking format",
			denom:          "liquidstake/invalid",
			expectMetadata: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			metadata := types.ExtractLiquidStakingMetadata(tc.denom, record)
			
			if tc.expectMetadata {
				require.NotNil(t, metadata)
				require.Equal(t, record.ValidatorAddress, metadata.ValidatorAddress)
				require.Equal(t, record.Id, metadata.RecordId)
				require.Equal(t, record.SharesAmount, metadata.SharesAmount)
				require.NotEmpty(t, metadata.CreatedAt)
			} else {
				require.Nil(t, metadata)
			}
		})
	}
}

func TestLiquidStakingAcknowledgement(t *testing.T) {
	testCases := []struct {
		name          string
		result        []byte
		err           error
		expectSuccess bool
	}{
		{
			name:          "successful acknowledgement",
			result:        []byte("success"),
			err:           nil,
			expectSuccess: true,
		},
		{
			name:          "error acknowledgement",
			result:        nil,
			err:           types.ErrInvalidRequest,
			expectSuccess: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ack := types.NewLiquidStakingAcknowledgement(tc.result, tc.err)
			
			require.Equal(t, tc.expectSuccess, ack.Success())
			
			if tc.expectSuccess {
				require.Equal(t, tc.result, ack.Result)
				require.Empty(t, ack.Error)
			} else {
				require.Nil(t, ack.Result)
				require.NotEmpty(t, ack.Error)
			}
			
			// Test serialization
			bz := ack.GetBytes()
			require.NotEmpty(t, bz)
			
			var unmarshaled types.LiquidStakingAcknowledgement
			err := json.Unmarshal(bz, &unmarshaled)
			require.NoError(t, err)
			require.Equal(t, ack.Success(), unmarshaled.Success())
		})
	}
}