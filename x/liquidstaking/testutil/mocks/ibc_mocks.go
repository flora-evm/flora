package mocks

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"
)

// MockTransferKeeper is a mock implementation of TransferKeeper for testing
type MockTransferKeeper struct {
	SendTransferFn func(
		ctx sdk.Context,
		sourcePort,
		sourceChannel string,
		token sdk.Coin,
		sender sdk.AccAddress,
		receiver string,
		timeoutHeight clienttypes.Height,
		timeoutTimestamp uint64,
		memo string,
	) (uint64, error)
}

func (m *MockTransferKeeper) SendTransfer(
	ctx sdk.Context,
	sourcePort,
	sourceChannel string,
	token sdk.Coin,
	sender sdk.AccAddress,
	receiver string,
	timeoutHeight clienttypes.Height,
	timeoutTimestamp uint64,
	memo string,
) (uint64, error) {
	if m.SendTransferFn != nil {
		return m.SendTransferFn(ctx, sourcePort, sourceChannel, token, sender, receiver, timeoutHeight, timeoutTimestamp, memo)
	}
	return 0, nil
}

// MockChannelKeeper is a mock implementation of ChannelKeeper for testing
type MockChannelKeeper struct {
	GetChannelFn                  func(ctx sdk.Context, portID, channelID string) (channel channeltypes.Channel, found bool)
	GetNextSequenceSendFn         func(ctx sdk.Context, portID, channelID string) (uint64, bool)
	GetChannelClientStateFn       func(ctx sdk.Context, portID, channelID string) (string, ibcexported.ClientState, error)
}

// GetChannel implements ChannelKeeper
func (m *MockChannelKeeper) GetChannel(ctx sdk.Context, portID, channelID string) (channel channeltypes.Channel, found bool) {
	if m.GetChannelFn != nil {
		return m.GetChannelFn(ctx, portID, channelID)
	}
	return channeltypes.Channel{}, false
}

// GetNextSequenceSend implements ChannelKeeper
func (m *MockChannelKeeper) GetNextSequenceSend(ctx sdk.Context, portID, channelID string) (uint64, bool) {
	if m.GetNextSequenceSendFn != nil {
		return m.GetNextSequenceSendFn(ctx, portID, channelID)
	}
	return 0, false
}

// GetChannelClientState implements ChannelKeeper
func (m *MockChannelKeeper) GetChannelClientState(ctx sdk.Context, portID, channelID string) (string, ibcexported.ClientState, error) {
	if m.GetChannelClientStateFn != nil {
		return m.GetChannelClientStateFn(ctx, portID, channelID)
	}
	return "", nil, nil
}