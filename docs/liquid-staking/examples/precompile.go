// Package liquidstaking provides EVM precompiles for liquid staking functionality
package liquidstaking

import (
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	
	sdk "github.com/cosmos/cosmos-sdk/types"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	
	"github.com/cosmos/evm/x/erc20/keeper"
	evmkeeper "github.com/cosmos/evm/x/vm/keeper"
	liquidkeeper "github.com/cosmos/gaia/x/liquid/keeper" // Would need to import
	tokenfactorykeeper "github.com/strangelove-ventures/tokenfactory/x/tokenfactory/keeper"
)

const (
	// LiquidStakingAddress defines the liquid staking precompile address
	LiquidStakingAddress = "0x0000000000000000000000000000000000000800"
	
	// Gas costs for operations
	GasTokenizeShares   = 100_000
	GasRedeemTokens     = 80_000
	GasTransferRecord   = 50_000
	GasQueryRecord      = 5_000
	GasQueryStats       = 10_000
)

// Precompile implements the liquid staking precompile
type Precompile struct {
	abi                   abi.ABI
	liquidStakingKeeper   liquidkeeper.Keeper
	stakingKeeper         stakingkeeper.Keeper
	bankKeeper           bankkeeper.Keeper
	tokenFactoryKeeper   tokenfactorykeeper.Keeper
	erc20Keeper          keeper.Keeper
	evmKeeper            *evmkeeper.Keeper
	authzKeeper          authzkeeper.Keeper
	lstTokenContracts    map[string]common.Address // validator -> LST contract
}

// NewPrecompile creates a new liquid staking precompile
func NewPrecompile(
	liquidStakingKeeper liquidkeeper.Keeper,
	stakingKeeper stakingkeeper.Keeper,
	bankKeeper bankkeeper.Keeper,
	tokenFactoryKeeper tokenfactorykeeper.Keeper,
	erc20Keeper keeper.Keeper,
	evmKeeper *evmkeeper.Keeper,
	authzKeeper authzkeeper.Keeper,
) (*Precompile, error) {
	// Parse ABI
	parsedABI, err := abi.JSON(strings.NewReader(LiquidStakingABI))
	if err != nil {
		return nil, err
	}

	return &Precompile{
		abi:                 parsedABI,
		liquidStakingKeeper: liquidStakingKeeper,
		stakingKeeper:       stakingKeeper,
		bankKeeper:          bankKeeper,
		tokenFactoryKeeper:  tokenFactoryKeeper,
		erc20Keeper:         erc20Keeper,
		evmKeeper:           evmKeeper,
		authzKeeper:         authzKeeper,
		lstTokenContracts:   make(map[string]common.Address),
	}, nil
}

// Address returns the precompile address
func (p *Precompile) Address() common.Address {
	return common.HexToAddress(LiquidStakingAddress)
}

// RequiredGas calculates the gas required for the given input
func (p *Precompile) RequiredGas(input []byte) uint64 {
	// Parse method ID from input
	if len(input) < 4 {
		return 0
	}
	
	method, err := p.abi.MethodById(input[:4])
	if err != nil {
		return 0
	}
	
	switch method.Name {
	case "tokenizeShares":
		return GasTokenizeShares
	case "redeemTokens":
		return GasRedeemTokens
	case "transferRecord":
		return GasTransferRecord
	case "getRecord":
		return GasQueryRecord
	case "getGlobalStats", "getValidatorStats":
		return GasQueryStats
	default:
		return 0
	}
}

// Run executes the precompile
func (p *Precompile) Run(evm *vm.EVM, contract *vm.Contract, readOnly bool) ([]byte, error) {
	// Decode input
	method, err := p.abi.MethodById(contract.Input[:4])
	if err != nil {
		return nil, err
	}
	
	args, err := method.Inputs.Unpack(contract.Input[4:])
	if err != nil {
		return nil, err
	}
	
	// Get SDK context
	ctx := evm.StateDB.(*statedb.StateDB).Context()
	
	switch method.Name {
	case "tokenizeShares":
		return p.tokenizeShares(ctx, evm, contract, args)
	case "redeemTokens":
		return p.redeemTokens(ctx, evm, contract, args)
	case "transferRecord":
		return p.transferRecord(ctx, evm, contract, args)
	case "getRecord":
		return p.getRecord(ctx, args)
	case "getGlobalStats":
		return p.getGlobalStats(ctx)
	case "getValidatorStats":
		return p.getValidatorStats(ctx, args)
	default:
		return nil, errors.New("unknown method")
	}
}

// tokenizeShares handles the tokenization of staked shares
func (p *Precompile) tokenizeShares(
	ctx sdk.Context,
	evm *vm.EVM,
	contract *vm.Contract,
	args []interface{},
) ([]byte, error) {
	// Parse arguments
	validator := args[0].(common.Address)
	amount := args[1].(*big.Int)
	
	// Convert addresses
	msgSender := contract.Caller()
	senderCosmosAddr := p.convertEthToCosmosAddr(msgSender)
	validatorCosmosAddr := p.convertEthToCosmosAddr(validator)
	
	// Create tokenize shares message
	msg := &liquidtypes.MsgTokenizeShares{
		DelegatorAddress: senderCosmosAddr.String(),
		ValidatorAddress: validatorCosmosAddr.String(),
		Amount:           sdk.NewCoin("petal", sdk.NewIntFromBigInt(amount)),
		TokenizedShareOwner: senderCosmosAddr.String(),
	}
	
	// Execute tokenization
	msgServer := liquidkeeper.NewMsgServerImpl(p.liquidStakingKeeper)
	res, err := msgServer.TokenizeShares(sdk.WrapSDKContext(ctx), msg)
	if err != nil {
		return nil, err
	}
	
	// Get or create LST token contract for this validator
	lstToken, err := p.getOrCreateLSTToken(ctx, evm, validatorCosmosAddr)
	if err != nil {
		return nil, err
	}
	
	// Mint LST tokens to sender
	mintAmount := res.Amount.Amount.BigInt()
	if err := p.mintLSTTokens(ctx, evm, lstToken, msgSender, mintAmount); err != nil {
		return nil, err
	}
	
	// Pack return values
	return p.abi.Methods["tokenizeShares"].Outputs.Pack(
		res.RecordId,
		lstToken,
	)
}

// redeemTokens handles redemption of LST tokens for staked shares
func (p *Precompile) redeemTokens(
	ctx sdk.Context,
	evm *vm.EVM,
	contract *vm.Contract,
	args []interface{},
) ([]byte, error) {
	recordId := args[0].(*big.Int)
	amount := args[1].(*big.Int)
	
	msgSender := contract.Caller()
	senderCosmosAddr := p.convertEthToCosmosAddr(msgSender)
	
	// Get tokenization record
	record, found := p.liquidStakingKeeper.GetTokenizeShareRecord(ctx, recordId.Uint64())
	if !found {
		return nil, errors.New("record not found")
	}
	
	// Verify ownership
	if record.Owner != senderCosmosAddr.String() {
		return nil, errors.New("not record owner")
	}
	
	// Get LST token for validator
	valAddr, _ := sdk.ValAddressFromBech32(record.Validator)
	lstToken := p.lstTokenContracts[valAddr.String()]
	
	// Burn LST tokens from sender
	if err := p.burnLSTTokens(ctx, evm, lstToken, msgSender, amount); err != nil {
		return nil, err
	}
	
	// Create redeem message
	msg := &liquidtypes.MsgRedeemTokensForShares{
		DelegatorAddress: senderCosmosAddr.String(),
		Amount:           sdk.NewCoin(record.GetShareTokenDenom(), sdk.NewIntFromBigInt(amount)),
	}
	
	// Execute redemption
	msgServer := liquidkeeper.NewMsgServerImpl(p.liquidStakingKeeper)
	_, err := msgServer.RedeemTokensForShares(sdk.WrapSDKContext(ctx), msg)
	if err != nil {
		return nil, err
	}
	
	return p.abi.Methods["redeemTokens"].Outputs.Pack(true)
}

// Helper functions

func (p *Precompile) convertEthToCosmosAddr(ethAddr common.Address) sdk.AccAddress {
	// Implementation for address conversion
	return sdk.AccAddress(ethAddr.Bytes())
}

func (p *Precompile) getOrCreateLSTToken(
	ctx sdk.Context,
	evm *vm.EVM,
	validator sdk.ValAddress,
) (common.Address, error) {
	// Check if LST token already exists
	if lstAddr, exists := p.lstTokenContracts[validator.String()]; exists {
		return lstAddr, nil
	}
	
	// Create new LST token using token factory
	valInfo, found := p.stakingKeeper.GetValidator(ctx, validator)
	if !found {
		return common.Address{}, errors.New("validator not found")
	}
	
	// Create denomination through token factory
	subdenom := fmt.Sprintf("spetal_%s", validator.String()[:8])
	creator := p.evmKeeper.GetParams(ctx).EvmDenom // Use module account
	
	_, err := p.tokenFactoryKeeper.CreateDenom(ctx, creator, subdenom)
	if err != nil {
		return common.Address{}, err
	}
	
	// Deploy ERC20 wrapper contract
	// This would deploy a special LST ERC20 that auto-compounds
	lstContract, err := p.deployLSTContract(ctx, evm, validator, valInfo.Description.Moniker)
	if err != nil {
		return common.Address{}, err
	}
	
	// Store mapping
	p.lstTokenContracts[validator.String()] = lstContract
	
	return lstContract, nil
}

// ABI definition
const LiquidStakingABI = `[
	{
		"inputs": [
			{"name": "validator", "type": "address"},
			{"name": "amount", "type": "uint256"}
		],
		"name": "tokenizeShares",
		"outputs": [
			{"name": "recordId", "type": "uint256"},
			{"name": "lstToken", "type": "address"}
		],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [
			{"name": "recordId", "type": "uint256"},
			{"name": "amount", "type": "uint256"}
		],
		"name": "redeemTokens",
		"outputs": [
			{"name": "success", "type": "bool"}
		],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [
			{"name": "recordId", "type": "uint256"},
			{"name": "newOwner", "type": "address"}
		],
		"name": "transferRecord",
		"outputs": [
			{"name": "success", "type": "bool"}
		],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [
			{"name": "recordId", "type": "uint256"}
		],
		"name": "getRecord",
		"outputs": [
			{"name": "owner", "type": "address"},
			{"name": "validator", "type": "address"},
			{"name": "shares", "type": "uint256"},
			{"name": "lstToken", "type": "address"}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "getGlobalStats",
		"outputs": [
			{"name": "totalLiquidStaked", "type": "uint256"},
			{"name": "globalCap", "type": "uint256"},
			{"name": "globalCapUsed", "type": "uint256"}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [
			{"name": "validator", "type": "address"}
		],
		"name": "getValidatorStats",
		"outputs": [
			{"name": "validatorLiquidStaked", "type": "uint256"},
			{"name": "validatorCap", "type": "uint256"},
			{"name": "validatorCapUsed", "type": "uint256"}
		],
		"stateMutability": "view",
		"type": "function"
	}
]`