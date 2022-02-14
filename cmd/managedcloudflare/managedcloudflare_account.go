package managedcloudflare

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
	"github.com/ukfast/cli/internal/pkg/helper"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/managedcloudflare"
)

func managedcloudflareAccountRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "account",
		Short: "sub-commands relating to accounts",
	}

	// Child commands
	cmd.AddCommand(managedcloudflareAccountListCmd(f))
	cmd.AddCommand(managedcloudflareAccountShowCmd(f))
	cmd.AddCommand(managedcloudflareAccountCreateCmd(f))

	// Child root commands
	cmd.AddCommand(managedcloudflareAccountMemberRootCmd(f))

	return cmd
}

func managedcloudflareAccountListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists accounts",
		Long:    "This command lists accounts",
		Example: "ukfast managedcloudflare account list",
		RunE:    managedcloudflareCobraRunEFunc(f, managedcloudflareAccountList),
	}
}

func managedcloudflareAccountList(service managedcloudflare.ManagedCloudflareService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	accounts, err := service.GetAccounts(params)
	if err != nil {
		return fmt.Errorf("Error retrieving accounts: %s", err)
	}

	return output.CommandOutput(cmd, OutputManagedCloudflareAccountsProvider(accounts))
}

func managedcloudflareAccountShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <account: id>...",
		Short:   "Shows a account",
		Long:    "This command shows one or more accounts",
		Example: "ukfast managedcloudflare account show 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing account")
			}

			return nil
		},
		RunE: managedcloudflareCobraRunEFunc(f, managedcloudflareAccountShow),
	}
}

func managedcloudflareAccountShow(service managedcloudflare.ManagedCloudflareService, cmd *cobra.Command, args []string) error {
	var accounts []managedcloudflare.Account
	for _, arg := range args {
		account, err := service.GetAccount(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving account [%s]: %s", arg, err)
			continue
		}

		accounts = append(accounts, account)
	}

	return output.CommandOutput(cmd, OutputManagedCloudflareAccountsProvider(accounts))
}

func managedcloudflareAccountCreateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Creates an account",
		Long:    "This command creates an account",
		Example: "ukfast managedcloudflare account create --name test",
		RunE:    managedcloudflareCobraRunEFunc(f, managedcloudflareAccountCreate),
	}

	cmd.Flags().String("name", "", "Name of account")
	cmd.MarkFlagRequired("name")

	return cmd
}

func managedcloudflareAccountCreate(service managedcloudflare.ManagedCloudflareService, cmd *cobra.Command, args []string) error {
	createRequest := managedcloudflare.CreateAccountRequest{}
	createRequest.Name, _ = cmd.Flags().GetString("name")
	accountID, err := service.CreateAccount(createRequest)
	if err != nil {
		return fmt.Errorf("Error creating account: %s", err)
	}

	account, err := service.GetAccount(accountID)
	if err != nil {
		return fmt.Errorf("Error retrieving new account: %s", err)
	}

	return output.CommandOutput(cmd, OutputManagedCloudflareAccountsProvider([]managedcloudflare.Account{account}))
}
