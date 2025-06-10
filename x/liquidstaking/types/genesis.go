package types

import (
	"fmt"
)

// NewGenesisState creates a new GenesisState instance
func NewGenesisState(params ModuleParams, records []TokenizationRecord, lastRecordID uint64) *GenesisState {
	return &GenesisState{
		Params:                    params,
		TokenizationRecords:       records,
		LastTokenizationRecordId:  lastRecordID,
	}
}

// DefaultGenesisState returns the default genesis state
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params:                    DefaultParams(),
		TokenizationRecords:       []TokenizationRecord{},
		LastTokenizationRecordId:  0,
	}
}

// Validate performs basic genesis state validation
func (gs GenesisState) Validate() error {
	// Validate params
	if err := gs.Params.Validate(); err != nil {
		return fmt.Errorf("invalid params: %w", err)
	}

	// Track seen IDs to check for duplicates
	seenIDs := make(map[uint64]bool)
	maxID := uint64(0)

	// Validate each tokenization record
	for i, record := range gs.TokenizationRecords {
		if err := record.Validate(); err != nil {
			return fmt.Errorf("invalid tokenization record at index %d: %w", i, err)
		}

		// Check for duplicate IDs
		if seenIDs[record.Id] {
			return fmt.Errorf("duplicate tokenization record ID: %d", record.Id)
		}
		seenIDs[record.Id] = true

		// Track max ID
		if record.Id > maxID {
			maxID = record.Id
		}
	}

	// Ensure last tokenization record ID is consistent
	if maxID > gs.LastTokenizationRecordId {
		return fmt.Errorf("last tokenization record ID (%d) is less than maximum record ID (%d)", 
			gs.LastTokenizationRecordId, maxID)
	}

	return nil
}