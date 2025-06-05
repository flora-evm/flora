package types

import (
	"fmt"
	
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Stage 1: Basic Types Definition
// No external dependencies, pure data structures

// TokenizationRecord represents a liquid staking position
type TokenizationRecord struct {
	// Unique identifier
	Id uint64 `json:"id"`
	
	// Validator address (bech32)
	Validator string `json:"validator"`
	
	// Owner of the tokenized shares
	Owner string `json:"owner"`
	
	// Amount of shares tokenized
	SharesTokenized sdk.Int `json:"shares_tokenized"`
	
	// Status for future use
	Status TokenizationStatus `json:"status"`
}

// TokenizationStatus represents the state of a tokenization record
type TokenizationStatus int32

const (
	// TOKENIZATION_STATUS_UNSPECIFIED defines an invalid status
	TOKENIZATION_STATUS_UNSPECIFIED TokenizationStatus = 0
	// TOKENIZATION_STATUS_TOKENIZED defines an active tokenization
	TOKENIZATION_STATUS_TOKENIZED TokenizationStatus = 1
	// TOKENIZATION_STATUS_REDEEMED defines a redeemed tokenization
	TOKENIZATION_STATUS_REDEEMED TokenizationStatus = 2
	// TOKENIZATION_STATUS_PENDING defines a pending operation
	TOKENIZATION_STATUS_PENDING TokenizationStatus = 3
)

// Validate performs basic validation on TokenizationRecord
func (tr TokenizationRecord) Validate() error {
	if tr.Id == 0 {
		return fmt.Errorf("id cannot be zero")
	}
	
	// Validate validator address format
	if _, err := sdk.ValAddressFromBech32(tr.Validator); err != nil {
		return fmt.Errorf("invalid validator address: %w", err)
	}
	
	// Validate owner address format
	if _, err := sdk.AccAddressFromBech32(tr.Owner); err != nil {
		return fmt.Errorf("invalid owner address: %w", err)
	}
	
	// Validate shares amount
	if tr.SharesTokenized.IsNil() || !tr.SharesTokenized.IsPositive() {
		return fmt.Errorf("shares tokenized must be positive")
	}
	
	// Validate status
	if tr.Status == TOKENIZATION_STATUS_UNSPECIFIED {
		return fmt.Errorf("status cannot be unspecified")
	}
	
	return nil
}

// ModuleParams defines the parameters for the liquid staking module
type ModuleParams struct {
	// Whether the module is enabled
	Enabled bool `json:"enabled"`
	
	// Global cap on liquid staking as percentage (0-1)
	GlobalLiquidStakingCap sdk.Dec `json:"global_liquid_staking_cap"`
	
	// Per-validator cap as percentage (0-1)
	ValidatorLiquidCap sdk.Dec `json:"validator_liquid_cap"`
	
	// Minimum tokenization amount
	MinTokenizationAmount sdk.Int `json:"min_tokenization_amount"`
}

// DefaultParams returns default module parameters
func DefaultParams() ModuleParams {
	return ModuleParams{
		Enabled:                true,
		GlobalLiquidStakingCap: sdk.NewDecWithPrec(25, 2), // 25%
		ValidatorLiquidCap:     sdk.NewDecWithPrec(50, 2), // 50%
		MinTokenizationAmount:  sdk.NewInt(1000000),       // 1 FLORA
	}
}

// Validate performs validation on ModuleParams
func (p ModuleParams) Validate() error {
	// Validate caps are between 0 and 1
	if p.GlobalLiquidStakingCap.IsNegative() || p.GlobalLiquidStakingCap.GT(sdk.OneDec()) {
		return fmt.Errorf("global liquid staking cap must be between 0 and 1")
	}
	
	if p.ValidatorLiquidCap.IsNegative() || p.ValidatorLiquidCap.GT(sdk.OneDec()) {
		return fmt.Errorf("validator liquid cap must be between 0 and 1")
	}
	
	// Validator cap cannot exceed global cap
	if p.ValidatorLiquidCap.GT(p.GlobalLiquidStakingCap) {
		return fmt.Errorf("validator cap cannot exceed global cap")
	}
	
	// Minimum amount must be positive
	if !p.MinTokenizationAmount.IsPositive() {
		return fmt.Errorf("minimum tokenization amount must be positive")
	}
	
	return nil
}

// GenesisState defines the module's genesis state
type GenesisState struct {
	// Module parameters
	Params ModuleParams `json:"params"`
	
	// List of tokenization records
	TokenizationRecords []TokenizationRecord `json:"tokenization_records"`
	
	// Last tokenization record ID
	LastTokenizationRecordId uint64 `json:"last_tokenization_record_id"`
}

// DefaultGenesisState returns default genesis state
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params:                   DefaultParams(),
		TokenizationRecords:      []TokenizationRecord{},
		LastTokenizationRecordId: 0,
	}
}

// Validate performs validation on GenesisState
func (gs GenesisState) Validate() error {
	// Validate params
	if err := gs.Params.Validate(); err != nil {
		return fmt.Errorf("invalid params: %w", err)
	}
	
	// Track seen IDs
	seenIds := make(map[uint64]bool)
	maxId := uint64(0)
	
	// Validate each record
	for _, record := range gs.TokenizationRecords {
		if err := record.Validate(); err != nil {
			return fmt.Errorf("invalid tokenization record %d: %w", record.Id, err)
		}
		
		// Check for duplicate IDs
		if seenIds[record.Id] {
			return fmt.Errorf("duplicate tokenization record id: %d", record.Id)
		}
		seenIds[record.Id] = true
		
		// Track max ID
		if record.Id > maxId {
			maxId = record.Id
		}
	}
	
	// Validate last ID is consistent
	if gs.LastTokenizationRecordId < maxId {
		return fmt.Errorf("last tokenization record id (%d) is less than max record id (%d)", 
			gs.LastTokenizationRecordId, maxId)
	}
	
	return nil
}