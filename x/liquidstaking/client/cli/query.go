package cli

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"

	"github.com/rollchains/flora/x/liquidstaking/types"
)

// GetQueryCmd returns the cli query commands for the liquid staking module
func GetQueryCmd(queryRoute string) *cobra.Command {
	// Group liquid staking queries under a subcommand
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		Long:                       `Query commands for the liquid staking module including tokenization records, parameters, and liquid staked amounts.`,
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		GetCmdQueryParams(),
		GetCmdQueryTokenizationRecord(),
		GetCmdQueryTokenizationRecords(),
		GetCmdQueryTokenizationRecordsByOwner(),
		GetCmdQueryTokenizationRecordsByValidator(),
		GetCmdQueryTotalLiquidStaked(),
		GetCmdQueryValidatorLiquidStaked(),
	)

	return cmd
}

// GetCmdQueryParams returns the command to query module parameters
func GetCmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Query the current liquid staking module parameters",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query the current parameters of the liquid staking module.

Example:
$ %s query %s params
`,
				version.AppName, types.ModuleName,
			),
		),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.Params(context.Background(), &types.QueryParamsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(&res.Params)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdQueryTokenizationRecord returns the command to query a specific tokenization record
func GetCmdQueryTokenizationRecord() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "record [record-id]",
		Short: "Query a specific tokenization record by ID",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details about a specific tokenization record.

Example:
$ %s query %s record 1
`,
				version.AppName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			recordID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid record ID: %w", err)
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.TokenizationRecord(
				context.Background(),
				&types.QueryTokenizationRecordRequest{Id: recordID},
			)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(&res.Record)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdQueryTokenizationRecords returns the command to query all tokenization records
func GetCmdQueryTokenizationRecords() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "records",
		Short: "Query all tokenization records",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query all tokenization records with optional pagination.

Example:
$ %s query %s records
$ %s query %s records --page=2 --limit=20
`,
				version.AppName, types.ModuleName,
				version.AppName, types.ModuleName,
			),
		),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.TokenizationRecords(
				context.Background(),
				&types.QueryTokenizationRecordsRequest{
					Pagination: pageReq,
				},
			)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "tokenization records")
	return cmd
}

// GetCmdQueryTokenizationRecordsByOwner returns the command to query records by owner
func GetCmdQueryTokenizationRecordsByOwner() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "records-by-owner [owner-address]",
		Short: "Query all tokenization records for a specific owner",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query all tokenization records owned by a specific address.

Example:
$ %s query %s records-by-owner flora1...
`,
				version.AppName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.TokenizationRecordsByOwner(
				context.Background(),
				&types.QueryTokenizationRecordsByOwnerRequest{
					OwnerAddress: args[0],
					Pagination:   pageReq,
				},
			)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "tokenization records")
	return cmd
}

// GetCmdQueryTokenizationRecordsByValidator returns the command to query records by validator
func GetCmdQueryTokenizationRecordsByValidator() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "records-by-validator [validator-address]",
		Short: "Query all tokenization records for a specific validator",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query all tokenization records associated with a specific validator.

Example:
$ %s query %s records-by-validator floravaloper1...
`,
				version.AppName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.TokenizationRecordsByValidator(
				context.Background(),
				&types.QueryTokenizationRecordsByValidatorRequest{
					ValidatorAddress: args[0],
					Pagination:       pageReq,
				},
			)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "tokenization records")
	return cmd
}

// GetCmdQueryTotalLiquidStaked returns the command to query total liquid staked amount
func GetCmdQueryTotalLiquidStaked() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "total-liquid-staked",
		Short: "Query the total amount of liquid staked tokens",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query the total amount of tokens that have been liquid staked across all validators.

Example:
$ %s query %s total-liquid-staked
`,
				version.AppName, types.ModuleName,
			),
		),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.TotalLiquidStaked(
				context.Background(),
				&types.QueryTotalLiquidStakedRequest{},
			)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdQueryValidatorLiquidStaked returns the command to query liquid staked amount for a validator
func GetCmdQueryValidatorLiquidStaked() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validator-liquid-staked [validator-address]",
		Short: "Query the amount of liquid staked tokens for a specific validator",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query the total amount of tokens that have been liquid staked for a specific validator.

Example:
$ %s query %s validator-liquid-staked floravaloper1...
`,
				version.AppName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.ValidatorLiquidStaked(
				context.Background(),
				&types.QueryValidatorLiquidStakedRequest{
					ValidatorAddress: args[0],
				},
			)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}