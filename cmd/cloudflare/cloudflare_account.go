package cloudflare

import (
	"errors"
	"fmt"

	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/helper"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/cloudflare"
	"github.com/spf13/cobra"
)

func cloudflareAccountRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "account",
		Short: "sub-commands relating to accounts",
	}

	// Child commands
	cmd.AddCommand(cloudflareAccountListCmd(f))
	cmd.AddCommand(cloudflareAccountShowCmd(f))
	cmd.AddCommand(cloudflareAccountCreateCmd(f))
	cmd.AddCommand(cloudflareAccountUpdateCmd(f))

	// Child root commands
	cmd.AddCommand(cloudflareAccountMemberRootCmd(f))

	return cmd
}

func cloudflareAccountListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists accounts",
		Long:    "This command lists accounts",
		Example: "ans cloudflare account list",
		RunE:    cloudflareCobraRunEFunc(f, cloudflareAccountList),
	}
}

func cloudflareAccountList(service cloudflare.CloudflareService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	accounts, err := service.GetAccounts(params)
	if err != nil {
		return fmt.Errorf("error retrieving accounts: %s", err)
	}

	return output.CommandOutput(cmd, AccountCollection(accounts))
}

func cloudflareAccountShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <account: id>...",
		Short:   "Shows a account",
		Long:    "This command shows one or more accounts",
		Example: "ans cloudflare account show e3f8baa0-b7c3-4a7a-958d-68e1aca3ea25",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("missing account")
			}

			return nil
		},
		RunE: cloudflareCobraRunEFunc(f, cloudflareAccountShow),
	}
}

func cloudflareAccountShow(service cloudflare.CloudflareService, cmd *cobra.Command, args []string) error {
	var accounts []cloudflare.Account
	for _, arg := range args {
		account, err := service.GetAccount(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving account [%s]: %s", arg, err)
			continue
		}

		accounts = append(accounts, account)
	}

	return output.CommandOutput(cmd, AccountCollection(accounts))
}

func cloudflareAccountCreateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Creates an account",
		Long:    "This command creates an account",
		Example: "ans cloudflare account create --name test",
		RunE:    cloudflareCobraRunEFunc(f, cloudflareAccountCreate),
	}

	cmd.Flags().String("name", "", "Name of account")
	_ = cmd.MarkFlagRequired("name")

	return cmd
}

func cloudflareAccountCreate(service cloudflare.CloudflareService, cmd *cobra.Command, args []string) error {
	createRequest := cloudflare.CreateAccountRequest{}
	createRequest.Name, _ = cmd.Flags().GetString("name")
	accountID, err := service.CreateAccount(createRequest)
	if err != nil {
		return fmt.Errorf("error creating account: %s", err)
	}

	account, err := service.GetAccount(accountID)
	if err != nil {
		return fmt.Errorf("error retrieving new account: %s", err)
	}

	return output.CommandOutput(cmd, AccountCollection([]cloudflare.Account{account}))
}

func cloudflareAccountUpdateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <account: id>...",
		Short:   "Updates an account",
		Long:    "This command updates one or more accounts",
		Example: "ans cloudflare account update e3f8baa0-b7c3-4a7a-958d-68e1aca3ea25",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("missing account")
			}

			return nil
		},
		RunE: cloudflareCobraRunEFunc(f, cloudflareAccountUpdate),
	}

	cmd.Flags().String("name", "", "Name of account")

	return cmd
}

func cloudflareAccountUpdate(service cloudflare.CloudflareService, cmd *cobra.Command, args []string) error {
	req := cloudflare.PatchAccountRequest{}
	req.Name, _ = cmd.Flags().GetString("name")

	for _, arg := range args {
		err := service.PatchAccount(arg, req)
		if err != nil {
			output.OutputWithErrorLevelf("Error updating account [%s]: %s", arg, err)
			continue
		}
	}

	return nil
}
