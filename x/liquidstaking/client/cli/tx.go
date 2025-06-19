package cli

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"

	"github.com/rollchains/flora/x/liquidstaking/types"
)

// GetTxCmd returns the transaction commands for the liquid staking module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		Long:                       `Transaction commands for the liquid staking module allowing users to tokenize shares and redeem liquid staking tokens.`,
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		NewTokenizeSharesCmd(),
		NewRedeemTokensCmd(),
	)

	return cmd
}

// NewTokenizeSharesCmd returns a CLI command for tokenizing delegation shares
func NewTokenizeSharesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tokenize-shares [validator-address] [amount] --owner [owner-address]",
		Short: "Tokenize delegation shares to receive liquid staking tokens",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Tokenize delegation shares for a specific validator to receive liquid staking tokens.

The shares amount should be specified in the delegation's share denomination.
If the owner flag is not provided, the liquid staking tokens will be sent to the delegator.

Examples:
$ %s tx %s tokenize-shares floravaloper1... 1000000 --from mykey
$ %s tx %s tokenize-shares floravaloper1... 1000000 --owner flora1... --from mykey
`,
				version.AppName, types.ModuleName,
				version.AppName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			validatorAddr := args[0]
			
			// Parse shares amount
			shares, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return fmt.Errorf("invalid shares amount: %w", err)
			}

			// Get owner address from flag
			ownerStr, err := cmd.Flags().GetString(FlagOwner)
			if err != nil {
				return err
			}

			// Create the message
			msg := &types.MsgTokenizeShares{
				DelegatorAddress: clientCtx.GetFromAddress().String(),
				ValidatorAddress: validatorAddr,
				Shares:           shares,
				OwnerAddress:     ownerStr,
			}

			// Validate the message
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	cmd.Flags().String(FlagOwner, "", "Optional owner address for the liquid staking tokens (defaults to delegator)")

	return cmd
}

// NewRedeemTokensCmd returns a CLI command for redeeming liquid staking tokens
func NewRedeemTokensCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "redeem-tokens [amount]",
		Short: "Redeem liquid staking tokens to restore delegation shares",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Redeem liquid staking tokens to restore the underlying delegation shares.

The amount should include the liquid staking token denomination.

Example:
$ %s tx %s redeem-tokens 1000000liquidstake/floravaloper1.../1 --from mykey
`,
				version.AppName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			// Parse amount with denomination
			amount, err := sdk.ParseCoinNormalized(args[0])
			if err != nil {
				return fmt.Errorf("invalid amount: %w", err)
			}

			// Verify it's a liquid staking token
			if !types.IsLiquidStakingTokenDenom(amount.Denom) {
				return fmt.Errorf("invalid liquid staking token denomination: %s", amount.Denom)
			}

			// Create the message
			msg := &types.MsgRedeemTokens{
				OwnerAddress: clientCtx.GetFromAddress().String(),
				Amount:       amount,
			}

			// Validate the message
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// Command flag constants
const (
	FlagOwner = "owner"
)