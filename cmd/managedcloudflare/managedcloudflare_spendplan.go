package managedcloudflare

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
	"github.com/ukfast/cli/internal/pkg/helper"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/managedcloudflare"
)

func managedcloudflareSpendPlanRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "spendplans",
		Short: "sub-commands relating to spend plans",
	}

	// Child commands
	cmd.AddCommand(managedcloudflareSpendPlanListCmd(f))

	return cmd
}

func managedcloudflareSpendPlanListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists spend plans",
		Long:    "This command lists spend plans",
		Example: "ukfast managedcloudflare spendplan list",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return managedcloudflareSpendPlanList(c.ManagedCloudflareService(), cmd, args)
		},
	}
}

func managedcloudflareSpendPlanList(service managedcloudflare.ManagedCloudflareService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	plans, err := service.GetSpendPlans(params)
	if err != nil {
		return fmt.Errorf("Error retrieving spend plans: %s", err)
	}

	return output.CommandOutput(cmd, OutputManagedCloudflareSpendPlansProvider(plans))
}
