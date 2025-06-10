package keeper_test

import (
	"testing"

	"cosmossdk.io/log"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	dbm "github.com/cosmos/cosmos-db"
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	"github.com/rollchains/flora/x/liquidstaking/keeper"
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

	ctx    sdk.Context
	keeper keeper.Keeper
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

	storeService := runtime.NewKVStoreService(storeKey)
	suite.keeper = keeper.NewKeeper(storeService, cdc)

	ctx := sdk.NewContext(stateStore, cmtproto.Header{}, false, log.NewNopLogger())
	suite.ctx = ctx
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) TestKeeperInstantiation() {
	suite.NotNil(suite.keeper)
	suite.NotNil(suite.keeper.GetStoreService())
}

func (suite *KeeperTestSuite) TestParams() {
	params := types.DefaultParams()

	// Test SetParams
	err := suite.keeper.SetParams(suite.ctx, params)
	suite.NoError(err)

	// Test GetParams
	gotParams := suite.keeper.GetParams(suite.ctx)
	suite.Equal(params, gotParams)

	// Test setting invalid params
	invalidParams := types.NewParams(
		math.LegacyNewDecWithPrec(-10, 2), // negative cap
		math.LegacyNewDecWithPrec(50, 2),
		true,
	)
	err = suite.keeper.SetParams(suite.ctx, invalidParams)
	suite.Error(err)
}

func (suite *KeeperTestSuite) TestTokenizationRecordID() {
	// Test default value
	id := suite.keeper.GetLastTokenizationRecordID(suite.ctx)
	suite.Equal(uint64(0), id)

	// Test setting and getting
	suite.keeper.SetLastTokenizationRecordID(suite.ctx, 42)
	id = suite.keeper.GetLastTokenizationRecordID(suite.ctx)
	suite.Equal(uint64(42), id)

	// Test updating
	suite.keeper.SetLastTokenizationRecordID(suite.ctx, 100)
	id = suite.keeper.GetLastTokenizationRecordID(suite.ctx)
	suite.Equal(uint64(100), id)
}

func (suite *KeeperTestSuite) TestTokenizationRecord() {
	validatorAddr := "floravaloper1tnh2q55v8wyygtt9srz5safamzdengsn9dsd7"
	ownerAddr := "flora1tnh2q55v8wyygtt9srz5safamzdengsnqeycj"
	
	record := types.NewTokenizationRecord(
		1,
		validatorAddr,
		ownerAddr,
		math.NewInt(1000),
	)

	// Test record not found
	_, found := suite.keeper.GetTokenizationRecord(suite.ctx, 1)
	suite.False(found)

	// Test setting record
	suite.keeper.SetTokenizationRecord(suite.ctx, record)

	// Test getting record
	gotRecord, found := suite.keeper.GetTokenizationRecord(suite.ctx, 1)
	suite.True(found)
	suite.Equal(record, gotRecord)

	// Test getting non-existent record
	_, found = suite.keeper.GetTokenizationRecord(suite.ctx, 999)
	suite.False(found)
}

func (suite *KeeperTestSuite) TestGetAllTokenizationRecords() {
	validatorAddr := "floravaloper1tnh2q55v8wyygtt9srz5safamzdengsn9dsd7"
	ownerAddr := "flora1tnh2q55v8wyygtt9srz5safamzdengsnqeycj"

	// Test empty state
	records := suite.keeper.GetAllTokenizationRecords(suite.ctx)
	suite.Empty(records)

	// Add some records
	record1 := types.NewTokenizationRecord(1, validatorAddr, ownerAddr, math.NewInt(1000))
	record2 := types.NewTokenizationRecord(2, validatorAddr, ownerAddr, math.NewInt(2000))
	record3 := types.NewTokenizationRecord(3, validatorAddr, ownerAddr, math.NewInt(3000))

	suite.keeper.SetTokenizationRecord(suite.ctx, record1)
	suite.keeper.SetTokenizationRecord(suite.ctx, record2)
	suite.keeper.SetTokenizationRecord(suite.ctx, record3)

	// Get all records
	records = suite.keeper.GetAllTokenizationRecords(suite.ctx)
	suite.Len(records, 3)
	suite.Contains(records, record1)
	suite.Contains(records, record2)
	suite.Contains(records, record3)
}

func (suite *KeeperTestSuite) TestInitExportGenesis() {
	validatorAddr := "floravaloper1tnh2q55v8wyygtt9srz5safamzdengsn9dsd7"
	ownerAddr := "flora1tnh2q55v8wyygtt9srz5safamzdengsnqeycj"

	// Create genesis state
	params := types.NewParams(
		math.LegacyNewDecWithPrec(30, 2), // 30%
		math.LegacyNewDecWithPrec(60, 2), // 60%
		false,
	)
	records := []types.TokenizationRecord{
		types.NewTokenizationRecord(1, validatorAddr, ownerAddr, math.NewInt(1000)),
		types.NewTokenizationRecord(2, validatorAddr, ownerAddr, math.NewInt(2000)),
	}
	genesis := types.NewGenesisState(params, records, 5)

	// Test InitGenesis
	suite.keeper.InitGenesis(suite.ctx, *genesis)

	// Verify state was set correctly
	gotParams := suite.keeper.GetParams(suite.ctx)
	suite.Equal(params, gotParams)

	gotLastID := suite.keeper.GetLastTokenizationRecordID(suite.ctx)
	suite.Equal(uint64(5), gotLastID)

	gotRecords := suite.keeper.GetAllTokenizationRecords(suite.ctx)
	suite.Len(gotRecords, 2)

	// Test ExportGenesis
	exportedGenesis := suite.keeper.ExportGenesis(suite.ctx)
	suite.Equal(genesis.Params, exportedGenesis.Params)
	suite.Equal(genesis.LastTokenizationRecordId, exportedGenesis.LastTokenizationRecordId)
	suite.ElementsMatch(genesis.TokenizationRecords, exportedGenesis.TokenizationRecords)
}