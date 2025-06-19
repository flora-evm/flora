package types

// Query request and response types for admin functionality

// QueryEmergencyStatusRequest is the request type for the Query/EmergencyStatus RPC method
type QueryEmergencyStatusRequest struct{}

// QueryEmergencyStatusResponse is the response type for the Query/EmergencyStatus RPC method
type QueryEmergencyStatusResponse struct {
	Paused    bool   `json:"paused"`
	PausedAt  int64  `json:"paused_at,omitempty"`
	UnpauseAt int64  `json:"unpause_at,omitempty"`
	Authority string `json:"authority,omitempty"`
}

// QueryValidatorWhitelistRequest is the request type for the Query/ValidatorWhitelist RPC method
type QueryValidatorWhitelistRequest struct{}

// QueryValidatorWhitelistResponse is the response type for the Query/ValidatorWhitelist RPC method
type QueryValidatorWhitelistResponse struct {
	Validators []string `json:"validators"`
}

// QueryValidatorBlacklistRequest is the request type for the Query/ValidatorBlacklist RPC method
type QueryValidatorBlacklistRequest struct{}

// QueryValidatorBlacklistResponse is the response type for the Query/ValidatorBlacklist RPC method
type QueryValidatorBlacklistResponse struct {
	Validators []string `json:"validators"`
}

// QueryValidatorStatusRequest is the request type for the Query/ValidatorStatus RPC method
type QueryValidatorStatusRequest struct {
	ValidatorAddress string `json:"validator_address"`
}

// QueryValidatorStatusResponse is the response type for the Query/ValidatorStatus RPC method
type QueryValidatorStatusResponse struct {
	IsAllowed      bool   `json:"is_allowed"`
	IsBlacklisted  bool   `json:"is_blacklisted"`
	IsWhitelisted  bool   `json:"is_whitelisted"`
	HasCustomCap   bool   `json:"has_custom_cap"`
	CustomCap      string `json:"custom_cap,omitempty"`
	LiquidStaked   string `json:"liquid_staked"`
	LiquidPercent  string `json:"liquid_percent"`
	TotalStaked    string `json:"total_staked"`
}