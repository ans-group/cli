package ecloud

import (
	"fmt"

	accountcmd "github.com/ans-group/cli/cmd/account"
	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/helper"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/spf13/cobra"
)

func ecloudCreditRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "credit",
		Short: "sub-commands relating to credits",
	}

	// Child commands
	cmd.AddCommand(ecloudCreditListCmd(f))

	return cmd
}

func ecloudCreditListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists credits",
		Long:    "This command lists credits",
		Example: "ans account credit list",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ecloudCreditList(c.ECloudService(), cmd, args)
		},
	}
}

func ecloudCreditList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	credits, err := service.GetCredits(params)
	if err != nil {
		return fmt.Errorf("error retrieving credits: %s", err)
	}

	return output.CommandOutput(cmd, accountcmd.CreditCollection(credits))
}
