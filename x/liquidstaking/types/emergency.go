package types

import (
	"time"
)

// EmergencyPauseInfo stores information about an emergency pause
type EmergencyPauseInfo struct {
	Paused     bool      `json:"paused" yaml:"paused"`
	PausedAt   time.Time `json:"paused_at" yaml:"paused_at"`
	UnpauseAt  time.Time `json:"unpause_at" yaml:"unpause_at"`
	Authority  string    `json:"authority" yaml:"authority"`
}

// EmergencyState represents the emergency state of the module
type EmergencyState struct {
	Paused      bool      `json:"paused" yaml:"paused"`
	PausedAt    time.Time `json:"paused_at" yaml:"paused_at"`
	PausedUntil time.Time `json:"paused_until" yaml:"paused_until"`
	Reason      string    `json:"reason" yaml:"reason"`
	Authority   string    `json:"authority" yaml:"authority"`
}

// ValidatorControls defines control settings for validators
type ValidatorControls struct {
	// Whitelist contains validators allowed to participate in liquid staking
	Whitelist []string `json:"whitelist" yaml:"whitelist"`
	
	// Blacklist contains validators banned from liquid staking
	Blacklist []string `json:"blacklist" yaml:"blacklist"`
	
	// CustomCaps contains custom liquid caps for specific validators
	CustomCaps map[string]string `json:"custom_caps" yaml:"custom_caps"`
}

// StringList is a simple wrapper for a list of strings
type StringList struct {
	Values []string `json:"values" yaml:"values"`
}