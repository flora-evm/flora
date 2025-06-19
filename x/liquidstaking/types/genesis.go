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

	// Track seen IDs and denoms to check for duplicates
	seenIDs := make(map[uint64]bool)
	seenDenoms := make(map[string]bool)
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

		// Check for duplicate denoms
		if record.Denom != "" {
			if seenDenoms[record.Denom] {
				return fmt.Errorf("duplicate denom %s at record ID %d", record.Denom, record.Id)
			}
			seenDenoms[record.Denom] = true

			// Validate denom format
			if !IsLiquidStakingTokenDenom(record.Denom) {
				return fmt.Errorf("invalid liquid staking token denom format: %s", record.Denom)
			}

			// Extract validator and record ID from denom and verify consistency
			expectedDenom := GenerateLiquidStakingTokenDenom(record.Validator, record.Id)
			if record.Denom != expectedDenom {
				return fmt.Errorf("denom %s does not match expected format for validator %s and record ID %d", 
					record.Denom, record.Validator, record.Id)
			}
		}

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

	// Validate record IDs are sequential (no gaps allowed in genesis)
	if len(gs.TokenizationRecords) > 0 {
		for i := uint64(1); i <= maxID; i++ {
			if !seenIDs[i] {
				return fmt.Errorf("missing tokenization record ID %d (IDs must be sequential in genesis)", i)
			}
		}
	}

	return nil
}