package precompile_test

import (
	"encoding/hex"
	"math/big"
	"testing"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/rollchains/flora/x/liquidstaking/keeper"
	"github.com/rollchains/flora/x/liquidstaking/precompile"
	"github.com/rollchains/flora/x/liquidstaking/types"
)

type PrecompileTestSuite struct {
	suite.Suite

	ctx            sdk.Context
	keeper         keeper.Keeper
	bankKeeper     *MockBankKeeper
	stakingKeeper  *MockStakingKeeper
	accountKeeper  *MockAccountKeeper
	precompile     *precompile.Contract
	abi            abi.ABI

	// Test addresses
	delegator      sdk.AccAddress
	validator      sdk.ValAddress
	evmDelegator   common.Address
	evmValidator   string
}

func TestPrecompileTestSuite(t *testing.T) {
	suite.Run(t, new(PrecompileTestSuite))
}

func (suite *PrecompileTestSuite) SetupTest() {
	// Initialize test context and keepers
	ctx, keeper, bk, sk, ak := setupTestEnvironment(suite.T())
	
	suite.ctx = ctx
	suite.keeper = keeper
	suite.bankKeeper = bk
	suite.stakingKeeper = sk
	suite.accountKeeper = ak

	// Create precompile
	suite.precompile = precompile.NewContract(keeper)
	
	// Load ABI
	err := precompile.LoadABI()
	suite.Require().NoError(err)
	suite.abi = precompile.ABI

	// Setup test addresses
	suite.delegator = sdk.AccAddress([]byte("delegator"))
	suite.validator = sdk.ValAddress([]byte("validator"))
	suite.evmDelegator = common.BytesToAddress(suite.delegator.Bytes())
	suite.evmValidator = sdk.ValAddress(suite.validator).String()

	// Setup initial state
	suite.setupInitialState()
}

func (suite *PrecompileTestSuite) setupInitialState() {
	// Create validator
	validator := stakingtypes.Validator{
		OperatorAddress: suite.evmValidator,
		Tokens:          math.NewInt(1000000000),
		DelegatorShares: math.LegacyNewDec(1000000000),
		Status:          stakingtypes.Bonded,
	}
	suite.stakingKeeper.validators[suite.evmValidator] = validator

	// Create delegation
	delegation := stakingtypes.Delegation{
		DelegatorAddress: suite.delegator.String(),
		ValidatorAddress: suite.evmValidator,
		Shares:           math.LegacyNewDec(1000000),
	}
	suite.stakingKeeper.delegations[suite.delegator.String()+suite.evmValidator] = delegation

	// Fund delegator account
	suite.bankKeeper.balances[suite.delegator.String()] = sdk.NewCoins(
		sdk.NewCoin("flora", math.NewInt(1000000000)),
		sdk.NewCoin("stake", math.NewInt(1000000000)),
	)
}

// Test query methods

func (suite *PrecompileTestSuite) TestGetParams() {
	// Pack method call
	input, err := suite.abi.Methods["getParams"].Inputs.Pack()
	suite.Require().NoError(err)
	
	methodID := suite.abi.Methods["getParams"].ID
	callData := append(methodID, input...)

	// Create mock EVM and contract
	evm, contract := suite.createMockEVMContract(callData)

	// Execute
	output, err := suite.precompile.Run(evm, contract, true)
	suite.Require().NoError(err)

	// Unpack output
	var params precompile.GetParamsResponse
	err = suite.abi.Methods["getParams"].Outputs.Unpack(&params, output)
	suite.Require().NoError(err)

	// Verify
	suite.True(params.Enabled)
	suite.Equal(big.NewInt(1000000), params.MinLiquidStakeAmount)
	suite.Equal(big.NewInt(2500), params.GlobalLiquidStakingCap) // 25% = 2500 basis points
	suite.Equal(big.NewInt(1000), params.ValidatorLiquidCap)     // 10% = 1000 basis points
}

func (suite *PrecompileTestSuite) TestGetTokenizationRecord() {
	// Create a tokenization record
	record := types.TokenizationRecord{
		Id:               1,
		ValidatorAddress: suite.evmValidator,
		Owner:            suite.delegator.String(),
		SharesDenomination: "shares/" + suite.evmValidator,
		LiquidStakingTokenDenom: "liquidstake/" + suite.evmValidator + "/1",
		SharesAmount:     math.LegacyNewDec(1000000),
		Status:           types.TokenizationRecord_ACTIVE,
		CreatedAt:        suite.ctx.BlockTime(),
	}
	suite.keeper.SetTokenizationRecord(suite.ctx, record)

	// Pack method call
	input, err := suite.abi.Methods["getTokenizationRecord"].Inputs.Pack(big.NewInt(1))
	suite.Require().NoError(err)
	
	methodID := suite.abi.Methods["getTokenizationRecord"].ID
	callData := append(methodID, input...)

	// Create mock EVM and contract
	evm, contract := suite.createMockEVMContract(callData)

	// Execute
	output, err := suite.precompile.Run(evm, contract, true)
	suite.Require().NoError(err)

	// Unpack output
	var result precompile.TokenizationRecord
	err = suite.abi.Methods["getTokenizationRecord"].Outputs.Unpack(&result, output)
	suite.Require().NoError(err)

	// Verify
	suite.Equal(big.NewInt(1), result.Id)
	suite.Equal(suite.evmValidator, result.ValidatorAddress)
	suite.Equal(suite.evmDelegator, result.Owner)
	suite.Equal(uint8(1), result.Status) // ACTIVE
}

func (suite *PrecompileTestSuite) TestGetTotalLiquidStaked() {
	// Set some liquid staking amounts
	suite.keeper.SetValidatorLiquidStakingAmount(suite.ctx, suite.evmValidator, math.NewInt(5000000))

	// Pack method call
	input, err := suite.abi.Methods["getTotalLiquidStaked"].Inputs.Pack()
	suite.Require().NoError(err)
	
	methodID := suite.abi.Methods["getTotalLiquidStaked"].ID
	callData := append(methodID, input...)

	// Create mock EVM and contract
	evm, contract := suite.createMockEVMContract(callData)

	// Execute
	output, err := suite.precompile.Run(evm, contract, true)
	suite.Require().NoError(err)

	// Unpack output
	var amount *big.Int
	err = suite.abi.Methods["getTotalLiquidStaked"].Outputs.Unpack(&amount, output)
	suite.Require().NoError(err)

	// Verify
	suite.Equal(big.NewInt(5000000), amount)
}

// Test transaction methods

func (suite *PrecompileTestSuite) TestTokenizeShares() {
	// Pack method call
	amount := big.NewInt(100000)
	input, err := suite.abi.Methods["tokenizeShares"].Inputs.Pack(
		suite.evmValidator,
		amount,
		common.Address{}, // Use zero address for self
	)
	suite.Require().NoError(err)
	
	methodID := suite.abi.Methods["tokenizeShares"].ID
	callData := append(methodID, input...)

	// Create mock EVM and contract with delegator as caller
	evm, contract := suite.createMockEVMContractWithCaller(callData, suite.evmDelegator)

	// Execute
	output, err := suite.precompile.Run(evm, contract, false)
	suite.Require().NoError(err)

	// Unpack output
	var response precompile.TokenizeSharesResponse
	err = suite.abi.Methods["tokenizeShares"].Outputs.Unpack(&response, output)
	suite.Require().NoError(err)

	// Verify
	suite.Equal(big.NewInt(1), response.RecordId)
	suite.Equal("liquidstake/"+suite.evmValidator+"/1", response.TokensDenom)
	suite.Equal(amount, response.TokensAmount)

	// Verify record was created
	record, found := suite.keeper.GetTokenizationRecord(suite.ctx, 1)
	suite.True(found)
	suite.Equal(suite.delegator.String(), record.Owner)
	suite.Equal(types.TokenizationRecord_ACTIVE, record.Status)
}

func (suite *PrecompileTestSuite) TestRedeemTokens() {
	// First create a tokenization record
	record := types.TokenizationRecord{
		Id:               1,
		ValidatorAddress: suite.evmValidator,
		Owner:            suite.delegator.String(),
		SharesDenomination: "shares/" + suite.evmValidator,
		LiquidStakingTokenDenom: "liquidstake/" + suite.evmValidator + "/1",
		SharesAmount:     math.LegacyNewDec(100000),
		Status:           types.TokenizationRecord_ACTIVE,
		CreatedAt:        suite.ctx.BlockTime(),
	}
	suite.keeper.SetTokenizationRecord(suite.ctx, record)

	// Mint LST tokens to owner
	lstCoin := sdk.NewCoin(record.LiquidStakingTokenDenom, math.NewInt(100000))
	suite.bankKeeper.balances[suite.delegator.String()] = suite.bankKeeper.balances[suite.delegator.String()].Add(lstCoin)

	// Pack method call
	amount := big.NewInt(50000)
	input, err := suite.abi.Methods["redeemTokens"].Inputs.Pack(
		record.LiquidStakingTokenDenom,
		amount,
	)
	suite.Require().NoError(err)
	
	methodID := suite.abi.Methods["redeemTokens"].ID
	callData := append(methodID, input...)

	// Create mock EVM and contract with owner as caller
	evm, contract := suite.createMockEVMContractWithCaller(callData, suite.evmDelegator)

	// Execute
	output, err := suite.precompile.Run(evm, contract, false)
	suite.Require().NoError(err)

	// Unpack output
	var response precompile.RedeemTokensResponse
	err = suite.abi.Methods["redeemTokens"].Outputs.Unpack(&response, output)
	suite.Require().NoError(err)

	// Verify
	suite.Equal(suite.evmValidator, response.ValidatorAddress)
	suite.Equal(big.NewInt(50000), response.SharesAmount)
	suite.False(response.Completed) // Not fully redeemed

	// Verify record was updated
	record, found := suite.keeper.GetTokenizationRecord(suite.ctx, 1)
	suite.True(found)
	suite.Equal(types.TokenizationRecord_ACTIVE, record.Status) // Still active
}

// Test error cases

func (suite *PrecompileTestSuite) TestInvalidMethodID() {
	// Create invalid method ID
	callData := []byte{0x00, 0x00, 0x00, 0x00}

	// Create mock EVM and contract
	evm, contract := suite.createMockEVMContract(callData)

	// Execute
	_, err := suite.precompile.Run(evm, contract, true)
	suite.Require().Error(err)
	suite.Contains(err.Error(), "method not found")
}

func (suite *PrecompileTestSuite) TestReadOnlyRestriction() {
	// Try to call tokenizeShares in read-only mode
	input, err := suite.abi.Methods["tokenizeShares"].Inputs.Pack(
		suite.evmValidator,
		big.NewInt(100000),
		common.Address{},
	)
	suite.Require().NoError(err)
	
	methodID := suite.abi.Methods["tokenizeShares"].ID
	callData := append(methodID, input...)

	// Create mock EVM and contract
	evm, contract := suite.createMockEVMContract(callData)

	// Execute in read-only mode
	_, err = suite.precompile.Run(evm, contract, true)
	suite.Require().Error(err)
	suite.Contains(err.Error(), "cannot call state-changing method in read-only mode")
}

func (suite *PrecompileTestSuite) TestRequiredGas() {
	// Test gas calculation for different methods
	testCases := []struct {
		method      string
		expectedGas uint64
	}{
		{"getParams", precompile.GasBaseQuery},
		{"getTokenizationRecord", precompile.GasGetRecord},
		{"getTotalLiquidStaked", precompile.GasBaseQuery},
		{"tokenizeShares", precompile.GasTokenizeShares},
		{"redeemTokens", precompile.GasRedeemTokens},
	}

	for _, tc := range testCases {
		methodID := suite.abi.Methods[tc.method].ID
		gas := suite.precompile.RequiredGas(methodID)
		suite.Equal(tc.expectedGas, gas, "gas mismatch for method %s", tc.method)
	}

	// Test invalid input
	gas := suite.precompile.RequiredGas([]byte{})
	suite.Equal(uint64(0), gas)
}

// Helper methods

func (suite *PrecompileTestSuite) createMockEVMContract(input []byte) (*vm.EVM, *vm.Contract) {
	return suite.createMockEVMContractWithCaller(input, common.Address{})
}

func (suite *PrecompileTestSuite) createMockEVMContractWithCaller(input []byte, caller common.Address) (*vm.EVM, *vm.Contract) {
	// Create mock StateDB
	stateDB := &MockStateDB{
		ctx: suite.ctx,
		logs: []*ethtypes.Log{},
	}

	// Create mock EVM
	evm := &vm.EVM{
		StateDB: stateDB,
	}

	// Create contract
	contract := &vm.Contract{
		Input:  input,
		caller: vm.AccountRef(caller),
	}

	return evm, contract
}

// MockStateDB implements the minimal StateDB interface needed for tests
type MockStateDB struct {
	ctx  sdk.Context
	logs []*ethtypes.Log
}

func (m *MockStateDB) GetContext() sdk.Context {
	return m.ctx
}

func (m *MockStateDB) AddLog(log *ethtypes.Log) {
	m.logs = append(m.logs, log)
}

// Minimal StateDB interface implementation for testing
func (m *MockStateDB) GetBalance(addr common.Address) *big.Int { return big.NewInt(0) }
func (m *MockStateDB) GetNonce(addr common.Address) uint64     { return 0 }
func (m *MockStateDB) GetCode(addr common.Address) []byte      { return nil }
func (m *MockStateDB) GetCodeSize(addr common.Address) int     { return 0 }
func (m *MockStateDB) GetCodeHash(common.Address) common.Hash  { return common.Hash{} }
func (m *MockStateDB) GetState(addr common.Address, key common.Hash) common.Hash {
	return common.Hash{}
}
func (m *MockStateDB) SetState(addr common.Address, key, value common.Hash) {}
func (m *MockStateDB) Exist(addr common.Address) bool                       { return true }

// Implement other StateDB methods as needed...

// Test hex encoding/decoding
func (suite *PrecompileTestSuite) TestHexEncoding() {
	// Test method IDs are correctly formatted
	for name, method := range suite.abi.Methods {
		suite.Equal(4, len(method.ID), "method ID should be 4 bytes for %s", name)
		
		// Verify it's valid hex
		hexStr := hex.EncodeToString(method.ID)
		suite.Equal(8, len(hexStr), "hex string should be 8 characters for %s", name)
	}
}

// Test ABI packing/unpacking
func (suite *PrecompileTestSuite) TestABIPackingUnpacking() {
	// Test TokenizationRecord struct packing
	record := precompile.TokenizationRecord{
		Id:                      big.NewInt(1),
		ValidatorAddress:        "floravaloper1abc",
		Owner:                   suite.evmDelegator,
		SharesDenomination:      "shares/floravaloper1abc",
		LiquidStakingTokenDenom: "liquidstake/floravaloper1abc/1",
		SharesAmount:            big.NewInt(1000000),
		Status:                  1, // ACTIVE
		CreatedAt:               big.NewInt(time.Now().Unix()),
		RedeemedAt:              big.NewInt(0),
	}

	// Pack
	packed, err := suite.abi.Methods["getTokenizationRecord"].Outputs.Pack(record)
	suite.Require().NoError(err)
	suite.NotEmpty(packed)

	// Unpack
	var unpacked precompile.TokenizationRecord
	err = suite.abi.Methods["getTokenizationRecord"].Outputs.Unpack(&unpacked, packed)
	suite.Require().NoError(err)
	
	// Verify fields match
	suite.Equal(record.Id, unpacked.Id)
	suite.Equal(record.ValidatorAddress, unpacked.ValidatorAddress)
	suite.Equal(record.Owner, unpacked.Owner)
	suite.Equal(record.Status, unpacked.Status)
}