package types

import (
	"encoding/json"
	"strings"
	"time"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/types/errors"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
)

// IBC packet data structures for liquid staking tokens

// LiquidStakingTokenPacketData extends the standard IBC transfer packet
// with liquid staking specific metadata
type LiquidStakingTokenPacketData struct {
	// Base IBC transfer fields
	Denom    string `json:"denom"`
	Amount   string `json:"amount"`
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	// Optional memo field for compatibility
	Memo string `json:"memo,omitempty"`
	
	// Liquid staking specific metadata
	LiquidStakingMetadata *LiquidStakingMetadata `json:"liquid_staking_metadata,omitempty"`
}

// LiquidStakingMetadata contains metadata about the liquid staking token
type LiquidStakingMetadata struct {
	// Original validator address on the source chain
	ValidatorAddress string `json:"validator_address"`
	// Original tokenization record ID
	RecordId uint64 `json:"record_id"`
	// Original shares amount that was tokenized
	SharesAmount math.LegacyDec `json:"shares_amount"`
	// Source chain ID where the token was originally created
	SourceChainId string `json:"source_chain_id"`
	// Creation timestamp
	CreatedAt string `json:"created_at"`
}

// NewLiquidStakingTokenPacketData creates a new liquid staking token packet
func NewLiquidStakingTokenPacketData(
	denom, amount, sender, receiver, memo string,
	metadata *LiquidStakingMetadata,
) LiquidStakingTokenPacketData {
	return LiquidStakingTokenPacketData{
		Denom:                 denom,
		Amount:                amount,
		Sender:                sender,
		Receiver:              receiver,
		Memo:                  memo,
		LiquidStakingMetadata: metadata,
	}
}

// ValidateBasic performs basic validation
func (pd LiquidStakingTokenPacketData) ValidateBasic() error {
	if pd.Denom == "" {
		return errors.ErrInvalidRequest.Wrap("denom cannot be empty")
	}
	
	amount, ok := math.NewIntFromString(pd.Amount)
	if !ok || amount.IsNegative() {
		return errors.ErrInvalidRequest.Wrap("invalid amount")
	}
	
	if pd.Sender == "" {
		return errors.ErrInvalidRequest.Wrap("sender cannot be empty")
	}
	
	if pd.Receiver == "" {
		return errors.ErrInvalidRequest.Wrap("receiver cannot be empty")
	}
	
	// Validate liquid staking metadata if present
	if pd.LiquidStakingMetadata != nil {
		if err := pd.LiquidStakingMetadata.ValidateBasic(); err != nil {
			return err
		}
	}
	
	return nil
}

// ValidateBasic performs basic validation on metadata
func (m *LiquidStakingMetadata) ValidateBasic() error {
	if m.ValidatorAddress == "" {
		return errors.ErrInvalidRequest.Wrap("validator address cannot be empty in metadata")
	}
	
	if m.RecordId == 0 {
		return errors.ErrInvalidRequest.Wrap("record ID must be positive")
	}
	
	if m.SharesAmount.IsNil() || m.SharesAmount.IsNegative() {
		return errors.ErrInvalidRequest.Wrap("shares amount must be positive")
	}
	
	if m.SourceChainId == "" {
		return errors.ErrInvalidRequest.Wrap("source chain ID cannot be empty")
	}
	
	if m.CreatedAt == "" {
		return errors.ErrInvalidRequest.Wrap("creation timestamp cannot be empty")
	}
	
	return nil
}

// GetBytes returns the JSON bytes of the packet data
func (pd LiquidStakingTokenPacketData) GetBytes() []byte {
	bz, err := json.Marshal(pd)
	if err != nil {
		panic(err)
	}
	return bz
}

// MustUnmarshalLiquidStakingTokenPacketData unmarshals the packet data
func MustUnmarshalLiquidStakingTokenPacketData(bz []byte) LiquidStakingTokenPacketData {
	var data LiquidStakingTokenPacketData
	if err := json.Unmarshal(bz, &data); err != nil {
		panic(err)
	}
	return data
}

// UnmarshalLiquidStakingTokenPacketData unmarshals the packet data
func UnmarshalLiquidStakingTokenPacketData(bz []byte) (LiquidStakingTokenPacketData, error) {
	var data LiquidStakingTokenPacketData
	if err := json.Unmarshal(bz, &data); err != nil {
		return data, err
	}
	return data, nil
}

// IsLiquidStakingToken checks if a denom is a liquid staking token
// and if it can be transferred via IBC with metadata
func IsLiquidStakingToken(denom string) bool {
	return IsLiquidStakingTokenDenom(denom)
}

// ExtractLiquidStakingMetadata extracts metadata from a liquid staking token denom
// Returns nil if the denom is not a liquid staking token
func ExtractLiquidStakingMetadata(denom string, record TokenizationRecord) *LiquidStakingMetadata {
	if !IsLiquidStakingToken(denom) {
		return nil
	}
	
	// Parse the denom to extract validator and record ID
	// Format: liquidstake/{validator}/{record-id}
	parts := strings.Split(denom, "/")
	if len(parts) != 3 {
		return nil
	}
	
	return &LiquidStakingMetadata{
		ValidatorAddress: record.Validator,
		RecordId:         record.Id,
		SharesAmount:     math.LegacyNewDecFromInt(record.SharesTokenized),
		SourceChainId:    "", // Will be set by the IBC module
		CreatedAt:        time.Now().Format(time.RFC3339), // Use current time since CreatedAt is not available
	}
}

// ConvertToTransferPacket converts liquid staking packet data to standard IBC transfer packet
// This is used for compatibility with existing IBC transfer infrastructure
func (pd LiquidStakingTokenPacketData) ConvertToTransferPacket() transfertypes.FungibleTokenPacketData {
	// Create memo that includes liquid staking metadata
	memo := pd.Memo
	if pd.LiquidStakingMetadata != nil {
		metadataJSON, _ := json.Marshal(map[string]interface{}{
			"liquid_staking_metadata": pd.LiquidStakingMetadata,
		})
		if memo == "" {
			memo = string(metadataJSON)
		} else {
			// Merge with existing memo
			var memoData map[string]interface{}
			if err := json.Unmarshal([]byte(memo), &memoData); err != nil {
				// If memo is not JSON, create new JSON with both
				memoData = map[string]interface{}{
					"text":                    memo,
					"liquid_staking_metadata": pd.LiquidStakingMetadata,
				}
			} else {
				memoData["liquid_staking_metadata"] = pd.LiquidStakingMetadata
			}
			memoJSON, _ := json.Marshal(memoData)
			memo = string(memoJSON)
		}
	}
	
	return transfertypes.FungibleTokenPacketData{
		Denom:    pd.Denom,
		Amount:   pd.Amount,
		Sender:   pd.Sender,
		Receiver: pd.Receiver,
		Memo:     memo,
	}
}

// ExtractFromTransferPacket extracts liquid staking packet data from standard IBC transfer packet
func ExtractFromTransferPacket(packet transfertypes.FungibleTokenPacketData) (LiquidStakingTokenPacketData, error) {
	lstData := LiquidStakingTokenPacketData{
		Denom:    packet.Denom,
		Amount:   packet.Amount,
		Sender:   packet.Sender,
		Receiver: packet.Receiver,
		Memo:     packet.Memo,
	}
	
	// Try to extract liquid staking metadata from memo
	if packet.Memo != "" {
		var memoData map[string]interface{}
		if err := json.Unmarshal([]byte(packet.Memo), &memoData); err == nil {
			if metadataRaw, ok := memoData["liquid_staking_metadata"]; ok {
				metadataJSON, _ := json.Marshal(metadataRaw)
				var metadata LiquidStakingMetadata
				if err := json.Unmarshal(metadataJSON, &metadata); err == nil {
					lstData.LiquidStakingMetadata = &metadata
				}
			}
		}
	}
	
	return lstData, nil
}

// IBC acknowledgement structures

// LiquidStakingAcknowledgement is the acknowledgement for liquid staking token transfers
type LiquidStakingAcknowledgement struct {
	Result []byte `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}

// Success returns a successful acknowledgement
func (ack LiquidStakingAcknowledgement) Success() bool {
	return ack.Error == ""
}

// NewLiquidStakingAcknowledgement creates a new acknowledgement
func NewLiquidStakingAcknowledgement(result []byte, err error) LiquidStakingAcknowledgement {
	if err != nil {
		return LiquidStakingAcknowledgement{
			Error: err.Error(),
		}
	}
	return LiquidStakingAcknowledgement{
		Result: result,
	}
}

// GetBytes returns the acknowledgement bytes
func (ack LiquidStakingAcknowledgement) GetBytes() []byte {
	bz, err := json.Marshal(ack)
	if err != nil {
		panic(err)
	}
	return bz
}