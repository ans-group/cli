package cloudflare

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
	"github.com/ukfast/cli/internal/pkg/helper"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/cloudflare"
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

	// Child root commands
	cmd.AddCommand(cloudflareAccountMemberRootCmd(f))

	return cmd
}

func cloudflareAccountListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists accounts",
		Long:    "This command lists accounts",
		Example: "ukfast cloudflare account list",
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
		return fmt.Errorf("Error retrieving accounts: %s", err)
	}

	return output.CommandOutput(cmd, OutputCloudflareAccountsProvider(accounts))
}

func cloudflareAccountShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <account: id>...",
		Short:   "Shows a account",
		Long:    "This command shows one or more accounts",
		Example: "ukfast cloudflare account show 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing account")
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

	return output.CommandOutput(cmd, OutputCloudflareAccountsProvider(accounts))
}

func cloudflareAccountCreateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Creates an account",
		Long:    "This command creates an account",
		Example: "ukfast cloudflare account create --name test",
		RunE:    cloudflareCobraRunEFunc(f, cloudflareAccountCreate),
	}

	cmd.Flags().String("name", "", "Name of account")
	cmd.MarkFlagRequired("name")

	return cmd
}

func cloudflareAccountCreate(service cloudflare.CloudflareService, cmd *cobra.Command, args []string) error {
	createRequest := cloudflare.CreateAccountRequest{}
	createRequest.Name, _ = cmd.Flags().GetString("name")
	accountID, err := service.CreateAccount(createRequest)
	if err != nil {
		return fmt.Errorf("Error creating account: %s", err)
	}

	account, err := service.GetAccount(accountID)
	if err != nil {
		return fmt.Errorf("Error retrieving new account: %s", err)
	}

	return output.CommandOutput(cmd, OutputCloudflareAccountsProvider([]cloudflare.Account{account}))
}
