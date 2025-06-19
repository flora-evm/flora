package cli_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/cosmos/cosmos-sdk/client/flags"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/require"

	"github.com/rollchains/flora/x/liquidstaking/client/cli"
	"github.com/rollchains/flora/x/liquidstaking/types"
)

func TestGetQueryCmd(t *testing.T) {
	cmd := cli.GetQueryCmd(types.StoreKey)
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
	require.Contains(t, cmdNames, "params")
	require.Contains(t, cmdNames, "record [record-id]")
	require.Contains(t, cmdNames, "records")
	require.Contains(t, cmdNames, "records-by-owner [owner-address]")
	require.Contains(t, cmdNames, "records-by-validator [validator-address]")
	require.Contains(t, cmdNames, "total-liquid-staked")
	require.Contains(t, cmdNames, "validator-liquid-staked [validator-address]")
}

func TestQueryParamsCmd(t *testing.T) {
	cmd := cli.GetCmdQueryParams()
	require.NotNil(t, cmd)
	require.Equal(t, "params", cmd.Use)
	require.NotEmpty(t, cmd.Short)
	
	// Test that command validates args
	err := cmd.Args(nil, []string{"unexpected-arg"})
	require.Error(t, err)
	
	err = cmd.Args(nil, []string{})
	require.NoError(t, err)
}

func TestQueryTokenizationRecordCmd(t *testing.T) {
	cmd := cli.GetCmdQueryTokenizationRecord()
	require.NotNil(t, cmd)
	require.True(t, strings.HasPrefix(cmd.Use, "record"))
	require.NotEmpty(t, cmd.Short)
	
	// Test that command validates args
	err := cmd.Args(nil, []string{})
	require.Error(t, err) // Expects exactly 1 argument
	
	err = cmd.Args(nil, []string{"1", "2"})
	require.Error(t, err) // Too many arguments
	
	err = cmd.Args(nil, []string{"1"})
	require.NoError(t, err)
}

func TestQueryTokenizationRecordsCmd(t *testing.T) {
	cmd := cli.GetCmdQueryTokenizationRecords()
	require.NotNil(t, cmd)
	require.Equal(t, "records", cmd.Use)
	require.NotEmpty(t, cmd.Short)
	
	// Verify pagination flags are added
	flagNames := []string{}
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		flagNames = append(flagNames, f.Name)
	})
	
	require.Contains(t, flagNames, flags.FlagPage)
	require.Contains(t, flagNames, flags.FlagLimit)
}

func TestQueryTokenizationRecordsByOwnerCmd(t *testing.T) {
	cmd := cli.GetCmdQueryTokenizationRecordsByOwner()
	require.NotNil(t, cmd)
	require.True(t, strings.HasPrefix(cmd.Use, "records-by-owner"))
	require.NotEmpty(t, cmd.Short)
	
	// Test that command validates args
	err := cmd.Args(nil, []string{})
	require.Error(t, err) // Expects exactly 1 argument
	
	err = cmd.Args(nil, []string{"flora1..."})
	require.NoError(t, err)
}

func TestQueryTokenizationRecordsByValidatorCmd(t *testing.T) {
	cmd := cli.GetCmdQueryTokenizationRecordsByValidator()
	require.NotNil(t, cmd)
	require.True(t, strings.HasPrefix(cmd.Use, "records-by-validator"))
	require.NotEmpty(t, cmd.Short)
	
	// Test that command validates args
	err := cmd.Args(nil, []string{})
	require.Error(t, err) // Expects exactly 1 argument
	
	err = cmd.Args(nil, []string{"floravaloper1..."})
	require.NoError(t, err)
}

func TestQueryTotalLiquidStakedCmd(t *testing.T) {
	cmd := cli.GetCmdQueryTotalLiquidStaked()
	require.NotNil(t, cmd)
	require.Equal(t, "total-liquid-staked", cmd.Use)
	require.NotEmpty(t, cmd.Short)
	
	// Test that command validates args
	err := cmd.Args(nil, []string{"unexpected-arg"})
	require.Error(t, err)
	
	err = cmd.Args(nil, []string{})
	require.NoError(t, err)
}

func TestQueryValidatorLiquidStakedCmd(t *testing.T) {
	cmd := cli.GetCmdQueryValidatorLiquidStaked()
	require.NotNil(t, cmd)
	require.True(t, strings.HasPrefix(cmd.Use, "validator-liquid-staked"))
	require.NotEmpty(t, cmd.Short)
	
	// Test that command validates args
	err := cmd.Args(nil, []string{})
	require.Error(t, err) // Expects exactly 1 argument
	
	err = cmd.Args(nil, []string{"floravaloper1..."})
	require.NoError(t, err)
}

// TestQueryCmdIntegration tests query commands with mock client context
func TestQueryCmdIntegration(t *testing.T) {
	// This is a placeholder for integration tests that would use a mock client context
	// In a real test environment, you would:
	// 1. Create a mock client context
	// 2. Set up expected query responses
	// 3. Execute commands and verify outputs
	
	t.Skip("Integration tests require mock client context setup")
}