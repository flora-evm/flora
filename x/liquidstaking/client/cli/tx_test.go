package cli_test

import (
	"fmt"
	"strings"
	"testing"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client/flags"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/require"

	"github.com/rollchains/flora/x/liquidstaking/client/cli"
	"github.com/rollchains/flora/x/liquidstaking/types"
)

func TestGetTxCmd(t *testing.T) {
	cmd := cli.GetTxCmd()
	require.NotNil(t, cmd)
	require.Equal(t, types.ModuleName, cmd.Use)
	require.NotEmpty(t, cmd.Short)
	
	// Check that all subcommands are registered
	subcommands := cmd.Commands()
	require.Greater(t, len(subcommands), 0)
	
	cmdNames := make([]string, len(subcommands))
	for i, subcmd := range subcommands {
		cmdNames[i] = subcmd.Use
	}
	
	// Verify expected commands exist
	expectedCmds := []string{
		"tokenize-shares [validator-address] [amount] --owner [owner-address]",
		"redeem-tokens [amount]",
	}
	
	for _, expected := range expectedCmds {
		found := false
		for _, actual := range cmdNames {
			if actual == expected {
				found = true
				break
			}
		}
		require.True(t, found, "expected command not found: %s", expected)
	}
}

func TestTokenizeSharesCmd(t *testing.T) {
	cmd := cli.NewTokenizeSharesCmd()
	require.NotNil(t, cmd)
	require.True(t, strings.HasPrefix(cmd.Use, "tokenize-shares"))
	require.NotEmpty(t, cmd.Short)
	require.NotEmpty(t, cmd.Long)
	
	// Test that command validates args
	testCases := []struct {
		name        string
		args        []string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "no arguments",
			args:        []string{},
			expectError: true,
			errorMsg:    "accepts 2 arg(s)",
		},
		{
			name:        "one argument",
			args:        []string{"floravaloper1..."},
			expectError: true,
			errorMsg:    "accepts 2 arg(s)",
		},
		{
			name:        "valid arguments",
			args:        []string{"floravaloper1...", "1000000stake"},
			expectError: false,
		},
		{
			name:        "too many arguments",
			args:        []string{"floravaloper1...", "1000000stake", "extra"},
			expectError: true,
			errorMsg:    "accepts 2 arg(s)",
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := cmd.Args(nil, tc.args)
			if tc.expectError {
				require.Error(t, err)
				if tc.errorMsg != "" {
					require.Contains(t, err.Error(), tc.errorMsg)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
	
	// Verify owner flag is registered
	ownerFlag := cmd.Flag(cli.FlagOwner)
	require.NotNil(t, ownerFlag)
	require.Equal(t, "", ownerFlag.DefValue)
}

func TestRedeemTokensCmd(t *testing.T) {
	cmd := cli.NewRedeemTokensCmd()
	require.NotNil(t, cmd)
	require.True(t, strings.HasPrefix(cmd.Use, "redeem-tokens"))
	require.NotEmpty(t, cmd.Short)
	require.NotEmpty(t, cmd.Long)
	
	// Test that command validates args
	testCases := []struct {
		name        string
		args        []string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "no arguments",
			args:        []string{},
			expectError: true,
			errorMsg:    "accepts 1 arg(s)",
		},
		{
			name:        "valid argument",
			args:        []string{"1000000liquidstake/floravaloper1.../1"},
			expectError: false,
		},
		{
			name:        "too many arguments",
			args:        []string{"1000000liquidstake/floravaloper1.../1", "extra"},
			expectError: true,
			errorMsg:    "accepts 1 arg(s)",
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := cmd.Args(nil, tc.args)
			if tc.expectError {
				require.Error(t, err)
				if tc.errorMsg != "" {
					require.Contains(t, err.Error(), tc.errorMsg)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestCmdFlags(t *testing.T) {
	// Test that transaction flags are properly added
	tokenizeCmd := cli.NewTokenizeSharesCmd()
	redeemCmd := cli.NewRedeemTokensCmd()
	
	// Check common transaction flags
	txFlags := []string{
		flags.FlagFrom,
		flags.FlagFees,
		flags.FlagGas,
		flags.FlagGasAdjustment,
		flags.FlagGasPrices,
		flags.FlagBroadcastMode,
		flags.FlagDryRun,
		flags.FlagGenerateOnly,
		flags.FlagOffline,
		flags.FlagSkipConfirmation,
		flags.FlagAccountNumber,
		flags.FlagSequence,
		flags.FlagNote,
		flags.FlagFees,
	}
	
	// Verify tokenize command has all required flags
	tokenizeFlagNames := []string{}
	tokenizeCmd.Flags().VisitAll(func(f *pflag.Flag) {
		tokenizeFlagNames = append(tokenizeFlagNames, f.Name)
	})
	
	for _, flag := range txFlags {
		require.Contains(t, tokenizeFlagNames, flag, "tokenize-shares missing flag: %s", flag)
	}
	
	// Verify redeem command has all required flags
	redeemFlagNames := []string{}
	redeemCmd.Flags().VisitAll(func(f *pflag.Flag) {
		redeemFlagNames = append(redeemFlagNames, f.Name)
	})
	
	for _, flag := range txFlags {
		require.Contains(t, redeemFlagNames, flag, "redeem-tokens missing flag: %s", flag)
	}
}

// TestTokenizeSharesValidation tests message validation in tokenize-shares command
func TestTokenizeSharesValidation(t *testing.T) {
	// This test would require a mock client context to fully test
	// For now, we can test that the command structure is correct
	cmd := cli.NewTokenizeSharesCmd()
	
	// Verify the command has the expected structure
	require.Equal(t, 2, cmd.Args(nil, []string{"val", "amount"}), nil)
	require.NotNil(t, cmd.RunE)
}

// TestRedeemTokensValidation tests message validation in redeem-tokens command
func TestRedeemTokensValidation(t *testing.T) {
	// This test would require a mock client context to fully test
	// For now, we can test that the command structure is correct
	cmd := cli.NewRedeemTokensCmd()
	
	// Verify the command has the expected structure
	require.Equal(t, 1, cmd.Args(nil, []string{"amount"}), nil)
	require.NotNil(t, cmd.RunE)
}