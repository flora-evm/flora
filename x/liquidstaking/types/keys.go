package types

import (
	"encoding/binary"
)

const (
	// ModuleName defines the module name
	ModuleName = "liquidstaking"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for the liquid staking module
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName
)

// KVStore key prefixes
var (
	// Primary storage
	ParamsKey                       = []byte{0x01} // prefix for parameters
	TokenizationRecordPrefix        = []byte{0x02} // prefix for tokenization records
	LastTokenizationRecordIDKey     = []byte{0x03} // key for last tokenization record ID
	
	// Indexes for efficient querying
	TokenizationRecordByOwnerPrefix = []byte{0x04} // prefix for tokenization records by owner index
	TokenizationRecordByValidatorPrefix = []byte{0x05} // prefix for tokenization records by validator index
	TokenizationRecordByDenomPrefix = []byte{0x06} // prefix for tokenization record by denom index
	
	// Aggregates for state tracking
	TotalLiquidStakedKey           = []byte{0x07} // key for total liquid staked amount
	ValidatorLiquidStakedPrefix    = []byte{0x08} // prefix for validator liquid staked amounts
)

// GetTokenizationRecordKey returns the key for a tokenization record
func GetTokenizationRecordKey(id uint64) []byte {
	return append(TokenizationRecordPrefix, Uint64ToBytes(id)...)
}

// GetTokenizationRecordByOwnerKey returns the key for tokenization record by owner index
func GetTokenizationRecordByOwnerKey(owner string, id uint64) []byte {
	return append(append(TokenizationRecordByOwnerPrefix, []byte(owner)...), Uint64ToBytes(id)...)
}

// GetTokenizationRecordByValidatorKey returns the key for tokenization record by validator index
func GetTokenizationRecordByValidatorKey(validator string, id uint64) []byte {
	return append(append(TokenizationRecordByValidatorPrefix, []byte(validator)...), Uint64ToBytes(id)...)
}

// GetTokenizationRecordByDenomKey returns the key for tokenization record by denom index
func GetTokenizationRecordByDenomKey(denom string) []byte {
	return append(TokenizationRecordByDenomPrefix, []byte(denom)...)
}

// GetValidatorLiquidStakedKey returns the key for validator liquid staked amount
func GetValidatorLiquidStakedKey(validator string) []byte {
	return append(ValidatorLiquidStakedPrefix, []byte(validator)...)
}

// GetTokenizationRecordByOwnerPrefix returns the prefix for all records by owner
func GetTokenizationRecordByOwnerPrefixKey(owner string) []byte {
	return append(TokenizationRecordByOwnerPrefix, []byte(owner)...)
}

// GetTokenizationRecordByValidatorPrefix returns the prefix for all records by validator
func GetTokenizationRecordByValidatorPrefixKey(validator string) []byte {
	return append(TokenizationRecordByValidatorPrefix, []byte(validator)...)
}

// Uint64ToBytes converts a uint64 to bytes
func Uint64ToBytes(n uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, n)
	return b
}

// BytesToUint64 converts bytes to uint64
func BytesToUint64(b []byte) uint64 {
	return binary.BigEndian.Uint64(b)
}