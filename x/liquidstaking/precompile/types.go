package precompile

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
)

// LiquidStakingPrecompile defines the interface for the liquid staking precompile
type LiquidStakingPrecompile interface {
	// Address returns the precompile contract address
	Address() common.Address
	
	// RequiredGas calculates the gas required for a precompile call
	RequiredGas(input []byte) uint64
	
	// Run executes the precompile contract
	Run(evm *vm.EVM, contract *vm.Contract, readOnly bool) ([]byte, error)
	
	// IsTransaction checks if the method modifies state
	IsTransaction(method string) bool
}

// Precompile addresses
var (
	// LiquidStakingPrecompileAddress is the address of the liquid staking precompile
	// Using 0x0000000000000000000000000000000000000800 following Evmos convention
	LiquidStakingPrecompileAddress = common.HexToAddress("0x0000000000000000000000000000000000000800")
)

// Method IDs for the liquid staking precompile
const (
	// Query methods
	MethodGetParams                  = "getParams"
	MethodGetTokenizationRecord      = "getTokenizationRecord"
	MethodGetTokenizationRecords     = "getTokenizationRecords"
	MethodGetRecordsByOwner          = "getRecordsByOwner"
	MethodGetRecordsByValidator      = "getRecordsByValidator"
	MethodGetTotalLiquidStaked       = "getTotalLiquidStaked"
	MethodGetValidatorLiquidStaked   = "getValidatorLiquidStaked"
	MethodGetLiquidStakingTokenInfo  = "getLiquidStakingTokenInfo"
	
	// Transaction methods
	MethodTokenizeShares = "tokenizeShares"
	MethodRedeemTokens   = "redeemTokens"
)

// Gas costs for precompile operations
const (
	// Base costs
	GasBaseQuery       = uint64(3000)
	GasBaseTransaction = uint64(10000)
	
	// Per-item costs for iterations
	GasPerRecord = uint64(100)
	
	// Operation-specific costs
	GasTokenizeShares = uint64(50000)
	GasRedeemTokens   = uint64(30000)
	GasGetRecord      = uint64(5000)
	GasGetRecords     = uint64(10000)
)

// Event signatures
var (
	EventTokenizeShares = common.HexToHash("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef") // Placeholder
	EventRedeemTokens   = common.HexToHash("0xfedcba0987654321fedcba0987654321fedcba0987654321fedcba0987654321") // Placeholder
)

// Structs for precompile method arguments and returns

// GetParamsResponse represents the module parameters
type GetParamsResponse struct {
	Enabled                bool
	MinLiquidStakeAmount   *big.Int
	GlobalLiquidStakingCap *big.Int // Represented as basis points (10000 = 100%)
	ValidatorLiquidCap     *big.Int // Represented as basis points (10000 = 100%)
}

// TokenizationRecord represents a tokenization record in EVM format
type TokenizationRecord struct {
	Id                      *big.Int
	ValidatorAddress        string
	Owner                   common.Address
	SharesDenomination      string
	LiquidStakingTokenDenom string
	SharesAmount            *big.Int
	Status                  uint8 // 0 = UNSPECIFIED, 1 = ACTIVE, 2 = REDEEMED
	CreatedAt               *big.Int // Unix timestamp
	RedeemedAt              *big.Int // Unix timestamp, 0 if not redeemed
}

// TokenizeSharesArgs represents arguments for tokenizeShares method
type TokenizeSharesArgs struct {
	ValidatorAddress string
	Amount           *big.Int
	Owner            common.Address // Optional, defaults to msg.sender
}

// TokenizeSharesResponse represents the response from tokenizeShares
type TokenizeSharesResponse struct {
	RecordId     *big.Int
	TokensDenom  string
	TokensAmount *big.Int
}

// RedeemTokensArgs represents arguments for redeemTokens method
type RedeemTokensArgs struct {
	TokenDenom string
	Amount     *big.Int
}

// RedeemTokensResponse represents the response from redeemTokens
type RedeemTokensResponse struct {
	ValidatorAddress string
	SharesAmount     *big.Int
	Completed        bool
}

// LiquidStakingTokenInfo provides information about a liquid staking token
type LiquidStakingTokenInfo struct {
	IsLiquidStakingToken bool
	ValidatorAddress     string
	RecordId             *big.Int
	OriginalDelegator    common.Address
	Active               bool
}

// Errors
var (
	ErrInvalidMethod           = "invalid method"
	ErrInvalidArguments        = "invalid arguments"
	ErrModuleDisabled          = "liquid staking module is disabled"
	ErrInsufficientDelegation  = "insufficient delegation"
	ErrInsufficientBalance     = "insufficient balance"
	ErrRecordNotFound          = "tokenization record not found"
	ErrInvalidValidator        = "invalid validator address"
	ErrInvalidAmount           = "amount must be positive"
	ErrExceedsGlobalCap        = "exceeds global liquid staking cap"
	ErrExceedsValidatorCap     = "exceeds validator liquid staking cap"
	ErrTokenAlreadyRedeemed    = "tokens already redeemed"
	ErrUnauthorized            = "unauthorized"
)

// ABI definitions for the precompile methods
var (
	// Define ABI for each method
	ABILiquidStaking abi.ABI
	
	// Method signatures
	GetParamsMethod                abi.Method
	GetTokenizationRecordMethod    abi.Method
	GetTokenizationRecordsMethod   abi.Method
	GetRecordsByOwnerMethod        abi.Method
	GetRecordsByValidatorMethod    abi.Method
	GetTotalLiquidStakedMethod     abi.Method
	GetValidatorLiquidStakedMethod abi.Method
	GetLiquidStakingTokenInfoMethod abi.Method
	TokenizeSharesMethod           abi.Method
	RedeemTokensMethod             abi.Method
)

// InitializeABI initializes the ABI for the liquid staking precompile
func InitializeABI() error {
	// This will be implemented to parse the ABI JSON
	// For now, we define the structure
	return nil
}

// PackMethod packs the method call with arguments
func PackMethod(method string, args ...interface{}) ([]byte, error) {
	// Implementation will pack the method and arguments according to ABI
	return nil, nil
}

// UnpackMethod unpacks the method call input
func UnpackMethod(data []byte) (method string, args []interface{}, err error) {
	// Implementation will unpack the method and arguments
	return "", nil, nil
}

// PackOutput packs the output according to the method's return type
func PackOutput(method string, output interface{}) ([]byte, error) {
	// Implementation will pack the output
	return nil, nil
}

// UnpackOutput unpacks the output according to the method's return type
func UnpackOutput(method string, data []byte) (interface{}, error) {
	// Implementation will unpack the output
	return nil, nil
}