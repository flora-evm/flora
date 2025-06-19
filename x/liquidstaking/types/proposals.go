package types

import (
	"fmt"
	"strings"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

const (
	// ProposalTypeUpdateParams defines the type for parameter update proposals
	ProposalTypeUpdateParams = "UpdateLiquidStakingParams"
	
	// ProposalTypeEmergencyPause defines the type for emergency pause proposals
	ProposalTypeEmergencyPause = "EmergencyPauseLiquidStaking"
	
	// ProposalTypeUpdateValidatorCap defines the type for validator cap update proposals
	ProposalTypeUpdateValidatorCap = "UpdateValidatorLiquidCap"
)

// TODO: Uncomment after adding proto definitions
// var (
// 	_ govtypes.Content = &UpdateParamsProposal{}
// 	_ govtypes.Content = &EmergencyPauseProposal{}
// 	_ govtypes.Content = &UpdateValidatorCapProposal{}
// )

// UpdateParamsProposal defines a proposal to update the liquid staking module parameters
type UpdateParamsProposal struct {
	Title       string          `json:"title" yaml:"title"`
	Description string          `json:"description" yaml:"description"`
	Changes     []ParamChange   `json:"changes" yaml:"changes"`
}

// ParamChange defines a single parameter change
type ParamChange struct {
	Key      string `json:"key" yaml:"key"`
	Value    string `json:"value" yaml:"value"`
}

// NewUpdateParamsProposal creates a new parameter update proposal
func NewUpdateParamsProposal(title, description string, changes []ParamChange) *UpdateParamsProposal {
	return &UpdateParamsProposal{
		Title:       title,
		Description: description,
		Changes:     changes,
	}
}

// GetTitle returns the title of the proposal
func (p *UpdateParamsProposal) GetTitle() string { return p.Title }

// GetDescription returns the description of the proposal
func (p *UpdateParamsProposal) GetDescription() string { return p.Description }

// ProposalRoute returns the routing key of the proposal
func (p *UpdateParamsProposal) ProposalRoute() string { return RouterKey }

// ProposalType returns the type of the proposal
func (p *UpdateParamsProposal) ProposalType() string { return ProposalTypeUpdateParams }

// ValidateBasic performs basic validation of the proposal
func (p *UpdateParamsProposal) ValidateBasic() error {
	if err := govtypes.ValidateAbstract(p); err != nil {
		return err
	}
	
	if len(p.Changes) == 0 {
		return fmt.Errorf("no parameter changes specified")
	}
	
	// Validate each parameter change
	validParams := map[string]bool{
		"enabled":                               true,
		"global_liquid_staking_cap":            true,
		"validator_liquid_cap":                 true,
		"min_liquid_stake_amount":              true,
		"rate_limit_period_hours":              true,
		"global_daily_tokenization_percent":    true,
		"validator_daily_tokenization_percent": true,
		"global_daily_tokenization_count":      true,
		"validator_daily_tokenization_count":   true,
		"user_daily_tokenization_count":        true,
		"warning_threshold_percent":            true,
		"auto_compound_enabled":                true,
		"auto_compound_frequency_blocks":       true,
		"max_rate_change_per_update":          true,
		"min_blocks_between_updates":           true,
	}
	
	for _, change := range p.Changes {
		if !validParams[change.Key] {
			return fmt.Errorf("invalid parameter key: %s", change.Key)
		}
		if change.Value == "" {
			return fmt.Errorf("parameter value cannot be empty for key: %s", change.Key)
		}
	}
	
	return nil
}

// String returns a string representation of the proposal
func (p *UpdateParamsProposal) String() string {
	var changes []string
	for _, c := range p.Changes {
		changes = append(changes, fmt.Sprintf("%s: %s", c.Key, c.Value))
	}
	
	return fmt.Sprintf(`Update Liquid Staking Params Proposal:
  Title:       %s
  Description: %s
  Changes:
    %s
`, p.Title, p.Description, strings.Join(changes, "\n    "))
}

// EmergencyPauseProposal defines a proposal to pause/unpause the liquid staking module
type EmergencyPauseProposal struct {
	Title       string `json:"title" yaml:"title"`
	Description string `json:"description" yaml:"description"`
	Pause       bool   `json:"pause" yaml:"pause"`
	Duration    int64  `json:"duration" yaml:"duration"` // Duration in seconds (0 for permanent)
}

// NewEmergencyPauseProposal creates a new emergency pause proposal
func NewEmergencyPauseProposal(title, description string, pause bool, duration int64) *EmergencyPauseProposal {
	return &EmergencyPauseProposal{
		Title:       title,
		Description: description,
		Pause:       pause,
		Duration:    duration,
	}
}

// GetTitle returns the title of the proposal
func (p *EmergencyPauseProposal) GetTitle() string { return p.Title }

// GetDescription returns the description of the proposal
func (p *EmergencyPauseProposal) GetDescription() string { return p.Description }

// ProposalRoute returns the routing key of the proposal
func (p *EmergencyPauseProposal) ProposalRoute() string { return RouterKey }

// ProposalType returns the type of the proposal
func (p *EmergencyPauseProposal) ProposalType() string { return ProposalTypeEmergencyPause }

// ValidateBasic performs basic validation of the proposal
func (p *EmergencyPauseProposal) ValidateBasic() error {
	if err := govtypes.ValidateAbstract(p); err != nil {
		return err
	}
	
	if p.Duration < 0 {
		return fmt.Errorf("pause duration cannot be negative")
	}
	
	if !p.Pause && p.Duration > 0 {
		return fmt.Errorf("duration should be 0 when unpausing")
	}
	
	return nil
}

// String returns a string representation of the proposal
func (p *EmergencyPauseProposal) String() string {
	action := "Unpause"
	if p.Pause {
		action = "Pause"
	}
	
	durationStr := "permanent"
	if p.Duration > 0 {
		durationStr = fmt.Sprintf("%d seconds", p.Duration)
	}
	
	return fmt.Sprintf(`Emergency %s Liquid Staking Proposal:
  Title:       %s
  Description: %s
  Action:      %s
  Duration:    %s
`, action, p.Title, p.Description, action, durationStr)
}

// UpdateValidatorCapProposal defines a proposal to update a specific validator's liquid cap
type UpdateValidatorCapProposal struct {
	Title            string `json:"title" yaml:"title"`
	Description      string `json:"description" yaml:"description"`
	ValidatorAddress string `json:"validator_address" yaml:"validator_address"`
	LiquidCap        string `json:"liquid_cap" yaml:"liquid_cap"` // Decimal string
}

// NewUpdateValidatorCapProposal creates a new validator cap update proposal
func NewUpdateValidatorCapProposal(title, description, validatorAddress, liquidCap string) *UpdateValidatorCapProposal {
	return &UpdateValidatorCapProposal{
		Title:            title,
		Description:      description,
		ValidatorAddress: validatorAddress,
		LiquidCap:        liquidCap,
	}
}

// GetTitle returns the title of the proposal
func (p *UpdateValidatorCapProposal) GetTitle() string { return p.Title }

// GetDescription returns the description of the proposal
func (p *UpdateValidatorCapProposal) GetDescription() string { return p.Description }

// ProposalRoute returns the routing key of the proposal
func (p *UpdateValidatorCapProposal) ProposalRoute() string { return RouterKey }

// ProposalType returns the type of the proposal
func (p *UpdateValidatorCapProposal) ProposalType() string { return ProposalTypeUpdateValidatorCap }

// ValidateBasic performs basic validation of the proposal
func (p *UpdateValidatorCapProposal) ValidateBasic() error {
	if err := govtypes.ValidateAbstract(p); err != nil {
		return err
	}
	
	// Validate validator address
	_, err := sdk.ValAddressFromBech32(p.ValidatorAddress)
	if err != nil {
		return fmt.Errorf("invalid validator address: %w", err)
	}
	
	// Validate liquid cap
	cap, err := math.LegacyNewDecFromStr(p.LiquidCap)
	if err != nil {
		return fmt.Errorf("invalid liquid cap: %w", err)
	}
	
	if cap.IsNegative() {
		return fmt.Errorf("liquid cap cannot be negative")
	}
	
	if cap.GT(math.LegacyOneDec()) {
		return fmt.Errorf("liquid cap cannot exceed 100%%")
	}
	
	return nil
}

// String returns a string representation of the proposal
func (p *UpdateValidatorCapProposal) String() string {
	return fmt.Sprintf(`Update Validator Liquid Cap Proposal:
  Title:       %s
  Description: %s
  Validator:   %s
  Liquid Cap:  %s
`, p.Title, p.Description, p.ValidatorAddress, p.LiquidCap)
}