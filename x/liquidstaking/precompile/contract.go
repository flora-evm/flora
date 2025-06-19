package precompile

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/rollchains/flora/x/liquidstaking/keeper"
	"github.com/rollchains/flora/x/liquidstaking/types"
)

const (
	// PrecompileAddress defines the liquid staking precompile address
	PrecompileAddress = "0x0000000000000000000000000000000000000800"
)

var (
	_ vm.PrecompiledContract = &Contract{}

	//go:embed abi.json
	f embed.FS

	// ABI is the liquid staking precompile ABI
	ABI abi.ABI
)

func init() {
	// Load ABI from embedded file
	if err := LoadABI(); err != nil {
		panic(fmt.Errorf("failed to load liquid staking precompile ABI: %w", err))
	}
}

// LoadABI loads the embedded ABI
func LoadABI() error {
	abiBytes, err := f.ReadFile("abi.json")
	if err != nil {
		return err
	}

	return json.Unmarshal(abiBytes, &ABI)
}

// Contract is the liquid staking precompile contract
type Contract struct {
	keeper keeper.Keeper
}

// NewContract creates a new liquid staking precompile contract
func NewContract(k keeper.Keeper) *Contract {
	return &Contract{
		keeper: k,
	}
}

// Address returns the address of the liquid staking precompile contract
func (c *Contract) Address() common.Address {
	return common.HexToAddress(PrecompileAddress)
}

// RequiredGas calculates the gas required for the given input
func (c *Contract) RequiredGas(input []byte) uint64 {
	// Minimum 4 bytes for method ID
	if len(input) < 4 {
		return 0
	}

	methodID := input[:4]
	method, err := ABI.MethodById(methodID)
	if err != nil {
		return 0
	}

	// Calculate gas based on method
	switch method.Name {
	// Query methods - base query gas
	case MethodGetParams:
		return GasBaseQuery
	case MethodGetTokenizationRecord:
		return GasGetRecord
	case MethodGetTokenizationRecords:
		return GasGetRecords
	case MethodGetRecordsByOwner, MethodGetRecordsByValidator:
		return GasGetRecords
	case MethodGetTotalLiquidStaked, MethodGetValidatorLiquidStaked:
		return GasBaseQuery
	case MethodGetLiquidStakingTokenInfo:
		return GasGetRecord

	// Transaction methods - higher gas requirements
	case MethodTokenizeShares:
		return GasTokenizeShares
	case MethodRedeemTokens:
		return GasRedeemTokens
	}

	return GasBaseQuery
}

// Run executes the precompiled contract
func (c *Contract) Run(evm *vm.EVM, contract *vm.Contract, readOnly bool) ([]byte, error) {
	// Extract method ID (first 4 bytes)
	if len(contract.Input) < 4 {
		return nil, fmt.Errorf("invalid input length")
	}

	methodID := contract.Input[:4]
	method, err := ABI.MethodById(methodID)
	if err != nil {
		return nil, fmt.Errorf("method not found: %v", err)
	}

	// Unpack method arguments
	args, err := method.Inputs.Unpack(contract.Input[4:])
	if err != nil {
		return nil, fmt.Errorf("failed to unpack arguments: %v", err)
	}

	// Get the SDK context from the EVM StateDB
	// In production, this would be obtained from the EVM's StateDB
	// For testing, we use a mock StateDB that provides the context
	var ctx sdk.Context
	if stateDB, ok := evm.StateDB.(interface{ GetContext() sdk.Context }); ok {
		ctx = stateDB.GetContext()
	} else {
		// Fallback for cases where StateDB doesn't implement GetContext
		ctx = sdk.Context{}
	}

	// Handle read-only check for state-changing methods
	switch method.Name {
	case MethodTokenizeShares, MethodRedeemTokens:
		if readOnly {
			return nil, fmt.Errorf("cannot call state-changing method in read-only mode")
		}
	}

	// Execute the method
	switch method.Name {
	// Query methods
	case MethodGetParams:
		return c.getParams(ctx)
	case MethodGetTokenizationRecord:
		return c.getTokenizationRecord(ctx, args)
	case MethodGetTokenizationRecords:
		return c.getTokenizationRecords(ctx, args)
	case MethodGetRecordsByOwner:
		return c.getRecordsByOwner(ctx, args)
	case MethodGetRecordsByValidator:
		return c.getRecordsByValidator(ctx, args)
	case MethodGetTotalLiquidStaked:
		return c.getTotalLiquidStaked(ctx)
	case MethodGetValidatorLiquidStaked:
		return c.getValidatorLiquidStaked(ctx, args)
	case MethodGetLiquidStakingTokenInfo:
		return c.getLiquidStakingTokenInfo(ctx, args)

	// Transaction methods
	case MethodTokenizeShares:
		return c.tokenizeShares(ctx, evm, contract, args)
	case MethodRedeemTokens:
		return c.redeemTokens(ctx, evm, contract, args)

	default:
		return nil, fmt.Errorf("unknown method: %s", method.Name)
	}
}

// Query method implementations

func (c *Contract) getParams(ctx sdk.Context) ([]byte, error) {
	params := c.keeper.GetParams(ctx)

	// Convert to precompile response format
	response := GetParamsResponse{
		Enabled:                params.Enabled,
		MinLiquidStakeAmount:   params.MinLiquidStakeAmount.BigInt(),
		GlobalLiquidStakingCap: decToBasisPoints(params.GlobalLiquidStakingCap),
		ValidatorLiquidCap:     decToBasisPoints(params.ValidatorLiquidCap),
	}

	return ABI.Methods[MethodGetParams].Outputs.Pack(response)
}

func (c *Contract) getTokenizationRecord(ctx sdk.Context, args []interface{}) ([]byte, error) {
	recordID := args[0].(*big.Int)

	record, found := c.keeper.GetTokenizationRecord(ctx, recordID.Uint64())
	if !found {
		return nil, fmt.Errorf("tokenization record not found: %d", recordID)
	}

	// Convert to EVM format
	evmRecord := convertTokenizationRecord(record)
	return ABI.Methods[MethodGetTokenizationRecord].Outputs.Pack(evmRecord)
}

func (c *Contract) getTokenizationRecords(ctx sdk.Context, args []interface{}) ([]byte, error) {
	offset := args[0].(*big.Int)
	limit := args[1].(*big.Int)

	// Get all records
	allRecords := c.keeper.GetAllTokenizationRecords(ctx)
	total := uint64(len(allRecords))

	// Apply pagination
	start := offset.Uint64()
	end := start + limit.Uint64()
	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	// Convert records
	records := make([]TokenizationRecord, 0, end-start)
	for i := start; i < end; i++ {
		records = append(records, convertTokenizationRecord(allRecords[i]))
	}

	return ABI.Methods[MethodGetTokenizationRecords].Outputs.Pack(records, big.NewInt(int64(total)))
}

func (c *Contract) getRecordsByOwner(ctx sdk.Context, args []interface{}) ([]byte, error) {
	owner := args[0].(common.Address)
	offset := args[1].(*big.Int)
	limit := args[2].(*big.Int)

	// Convert Ethereum address to Cosmos address
	ownerAddr := sdk.AccAddress(owner.Bytes())

	// Get records by owner
	allRecords := c.keeper.GetTokenizationRecordsByOwner(ctx, ownerAddr.String())
	total := uint64(len(allRecords))

	// Apply pagination
	start := offset.Uint64()
	end := start + limit.Uint64()
	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	// Convert records
	records := make([]TokenizationRecord, 0, end-start)
	for i := start; i < end; i++ {
		records = append(records, convertTokenizationRecord(allRecords[i]))
	}

	return ABI.Methods[MethodGetRecordsByOwner].Outputs.Pack(records, big.NewInt(int64(total)))
}

func (c *Contract) getRecordsByValidator(ctx sdk.Context, args []interface{}) ([]byte, error) {
	validatorAddr := args[0].(string)
	offset := args[1].(*big.Int)
	limit := args[2].(*big.Int)

	// Get records by validator
	allRecords := c.keeper.GetTokenizationRecordsByValidator(ctx, validatorAddr)
	total := uint64(len(allRecords))

	// Apply pagination
	start := offset.Uint64()
	end := start + limit.Uint64()
	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	// Convert records
	records := make([]TokenizationRecord, 0, end-start)
	for i := start; i < end; i++ {
		records = append(records, convertTokenizationRecord(allRecords[i]))
	}

	return ABI.Methods[MethodGetRecordsByValidator].Outputs.Pack(records, big.NewInt(int64(total)))
}

func (c *Contract) getTotalLiquidStaked(ctx sdk.Context) ([]byte, error) {
	total := c.keeper.GetTotalLiquidStaked(ctx)
	return ABI.Methods[MethodGetTotalLiquidStaked].Outputs.Pack(total.BigInt())
}

func (c *Contract) getValidatorLiquidStaked(ctx sdk.Context, args []interface{}) ([]byte, error) {
	validatorAddr := args[0].(string)
	amount := c.keeper.GetValidatorLiquidStaked(ctx, validatorAddr)
	return ABI.Methods[MethodGetValidatorLiquidStaked].Outputs.Pack(amount.BigInt())
}

func (c *Contract) getLiquidStakingTokenInfo(ctx sdk.Context, args []interface{}) ([]byte, error) {
	tokenDenom := args[0].(string)

	info := LiquidStakingTokenInfo{
		IsLiquidStakingToken: false,
		ValidatorAddress:     "",
		RecordId:             big.NewInt(0),
		OriginalDelegator:    common.Address{},
		Active:               false,
	}

	if types.IsLiquidStakingTokenDenom(tokenDenom) {
		_, recordID, err := types.ParseLiquidStakingTokenDenom(tokenDenom)
		if err == nil {
			record, found := c.keeper.GetTokenizationRecord(ctx, recordID)
			if found {
				info.IsLiquidStakingToken = true
				info.ValidatorAddress = record.ValidatorAddress
				info.RecordId = big.NewInt(int64(record.Id))

				// Convert owner to Ethereum address
				ownerAddr, err := sdk.AccAddressFromBech32(record.Owner)
				if err == nil {
					info.OriginalDelegator = common.BytesToAddress(ownerAddr.Bytes())
				}

				info.Active = true // All existing records are considered active
			}
		}
	}

	return ABI.Methods[MethodGetLiquidStakingTokenInfo].Outputs.Pack(info)
}

// Transaction method implementations

func (c *Contract) tokenizeShares(ctx sdk.Context, evm *vm.EVM, contract *vm.Contract, args []interface{}) ([]byte, error) {
	validatorAddr := args[0].(string)
	amount := args[1].(*big.Int)
	owner := args[2].(common.Address)

	// Get the caller address
	caller := contract.Caller()
	delegator := sdk.AccAddress(caller.Bytes())

	// Determine owner
	ownerAddr := delegator
	if !bytes.Equal(owner.Bytes(), common.Address{}.Bytes()) {
		ownerAddr = sdk.AccAddress(owner.Bytes())
	}

	// Convert amount to shares
	shares := sdk.NewDecCoinFromDec("stake", math.LegacyNewDecFromBigInt(amount))

	// Create message
	msg := &types.MsgTokenizeShares{
		DelegatorAddress: delegator.String(),
		ValidatorAddress: validatorAddr,
		Shares:           shares,
		OwnerAddress:     ownerAddr.String(),
	}

	// Execute the message
	msgServer := keeper.NewMsgServerImpl(c.keeper)
	res, err := msgServer.TokenizeShares(ctx, msg)
	if err != nil {
		return nil, err
	}

	// Parse response
	recordID := res.RecordId
	tokenAmount := res.TokensMinted

	response := TokenizeSharesResponse{
		RecordId:     big.NewInt(int64(recordID)),
		TokensDenom:  res.TokenDenom,
		TokensAmount: tokenAmount.BigInt(),
	}

	// Emit event
	event := ABI.Events["TokenizeSharesEvent"]
	topics := []common.Hash{
		event.ID,
		common.BytesToHash(caller.Bytes()),                    // delegator (indexed)
		common.BytesToHash([]byte(validatorAddr)),            // validator (indexed)
		common.BytesToHash(big.NewInt(int64(recordID)).Bytes()), // recordId (indexed)
	}

	data, err := event.Inputs.NonIndexed().Pack(
		common.BytesToAddress(ownerAddr.Bytes()),
		amount,
		tokenAmount.BigInt(),
		res.TokenDenom,
	)
	if err != nil {
		return nil, err
	}

	evm.StateDB.AddLog(&ethtypes.Log{
		Address: c.Address(),
		Topics:  topics,
		Data:    data,
	})

	return ABI.Methods[MethodTokenizeShares].Outputs.Pack(response)
}

func (c *Contract) redeemTokens(ctx sdk.Context, evm *vm.EVM, contract *vm.Contract, args []interface{}) ([]byte, error) {
	tokenDenom := args[0].(string)
	amount := args[1].(*big.Int)

	// Get the caller address
	caller := contract.Caller()
	owner := sdk.AccAddress(caller.Bytes())

	// Create coin
	tokenAmount := math.NewIntFromBigInt(amount)
	token := sdk.NewCoin(tokenDenom, tokenAmount)

	// Create message
	msg := &types.MsgRedeemTokens{
		OwnerAddress: owner.String(),
		Amount:       token,
	}

	// Execute the message
	msgServer := keeper.NewMsgServerImpl(c.keeper)
	res, err := msgServer.RedeemTokens(ctx, msg)
	if err != nil {
		return nil, err
	}

	// Parse record ID from denom
	_, recordID, _ := types.ParseLiquidStakingTokenDenom(tokenDenom)

	response := RedeemTokensResponse{
		ValidatorAddress: res.ValidatorAddress,
		SharesAmount:     res.SharesRedeemed.Amount.TruncateInt().BigInt(),
		Completed:        res.RecordCompleted,
	}

	// Emit event
	event := ABI.Events["RedeemTokensEvent"]
	topics := []common.Hash{
		event.ID,
		common.BytesToHash(caller.Bytes()),                       // owner (indexed)
		common.BytesToHash([]byte(res.ValidatorAddress)),        // validator (indexed)
		common.BytesToHash(big.NewInt(int64(recordID)).Bytes()), // recordId (indexed)
	}

	data, err := event.Inputs.NonIndexed().Pack(
		tokenDenom,
		amount,
		res.SharesRedeemed.Amount.TruncateInt().BigInt(),
		res.RecordCompleted,
	)
	if err != nil {
		return nil, err
	}

	evm.StateDB.AddLog(&ethtypes.Log{
		Address: c.Address(),
		Topics:  topics,
		Data:    data,
	})

	return ABI.Methods[MethodRedeemTokens].Outputs.Pack(response)
}

// Helper functions

// convertTokenizationRecord converts a Cosmos tokenization record to EVM format
func convertTokenizationRecord(record types.TokenizationRecord) TokenizationRecord {
	ownerAddr, _ := sdk.AccAddressFromBech32(record.Owner)
	
	// All existing records are considered active (status = 1)
	status := uint8(1)

	return TokenizationRecord{
		Id:                      big.NewInt(int64(record.Id)),
		ValidatorAddress:        record.Validator,
		Owner:                   common.BytesToAddress(ownerAddr.Bytes()),
		SharesDenomination:      fmt.Sprintf("shares/%s", record.Validator),
		LiquidStakingTokenDenom: record.Denom,
		SharesAmount:            record.SharesTokenized.BigInt(),
		Status:                  status,
		CreatedAt:               big.NewInt(time.Now().Unix()), // Use current time as placeholder
		RedeemedAt:              big.NewInt(0),                  // Not redeemed
	}
}

// decToBasisPoints converts a decimal to basis points (10000 = 100%)
func decToBasisPoints(dec math.LegacyDec) *big.Int {
	// Multiply by 10000 to get basis points
	basisPoints := dec.MulInt64(10000).TruncateInt()
	return basisPoints.BigInt()
}