package keeper_test

import (
	"testing"

	corestore "cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/suite"

	"github.com/rollchains/flora/x/liquidstaking/keeper"
	"github.com/rollchains/flora/x/liquidstaking/testutil/mocks"
	"github.com/rollchains/flora/x/liquidstaking/types"
)

func init() {
	// Set bech32 prefixes for testing
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount("flora", "florapub")
	config.SetBech32PrefixForValidator("floravaloper", "floravaloperpub")
	config.SetBech32PrefixForConsensusNode("floravalcons", "floravalconspub")
}

type KeeperTestSuite struct {
	suite.Suite

	ctx          sdk.Context
	keeper       keeper.Keeper
	storeService corestore.KVStoreService
	storeKey     storetypes.StoreKey
	cdc          codec.BinaryCodec
	
	// Mocks
	mockBankKeeper     *mocks.MockBankKeeper
	mockStakingKeeper  *mocks.MockStakingKeeper
	mockAccountKeeper  *mocks.MockAccountKeeper
	mockTransferKeeper *mocks.MockTransferKeeper
	mockChannelKeeper  *mocks.MockChannelKeeper
	mockDistributionKeeper *mocks.MockDistributionKeeper
}

func (suite *KeeperTestSuite) SetupTest() {
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	
	db := dbm.NewMemDB()
	stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	require := suite.Require()
	require.NoError(stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	types.RegisterInterfaces(registry)
	cdc := codec.NewProtoCodec(registry)

	// Create mock keepers
	suite.mockBankKeeper = &mocks.MockBankKeeper{}
	suite.mockStakingKeeper = &mocks.MockStakingKeeper{}
	suite.mockAccountKeeper = &mocks.MockAccountKeeper{
		GetModuleAddressFn: func(moduleName string) sdk.AccAddress {
			return authtypes.NewModuleAddress(moduleName)
		},
	}
	suite.mockTransferKeeper = &mocks.MockTransferKeeper{}
	suite.mockChannelKeeper = &mocks.MockChannelKeeper{}
	suite.mockDistributionKeeper = &mocks.MockDistributionKeeper{}

	// Use gov module address as authority
	auth := authtypes.NewModuleAddress("gov").String()

	storeService := runtime.NewKVStoreService(storeKey)
	suite.keeper = keeper.NewKeeper(
		storeService,
		cdc,
		suite.mockStakingKeeper,
		suite.mockBankKeeper,
		suite.mockAccountKeeper,
		suite.mockTransferKeeper,
		suite.mockChannelKeeper,
		suite.mockDistributionKeeper,
		auth,
	)

	suite.ctx = sdk.NewContext(stateStore, cmtproto.Header{}, false, log.NewNopLogger())
	suite.storeService = storeService
	suite.storeKey = storeKey
	suite.cdc = cdc
	
	// Set default params
	params := types.DefaultParams()
	params.Enabled = true
	suite.keeper.SetParams(suite.ctx, params)
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}