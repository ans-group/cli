package billing

import (
	"fmt"

	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/helper"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/billing"
	"github.com/spf13/cobra"
)

func billingCloudCostRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cloudcost",
		Short: "sub-commands relating to invoices",
	}

	// Child commands
	cmd.AddCommand(billingCloudCostListCmd(f))

	return cmd
}

func billingCloudCostListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists cloud costs",
		Long:    "This command lists cloud costs",
		Example: "ans billing cloudcost list",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return billingCloudCostList(c.BillingService(), cmd, args)
		},
	}
}

func billingCloudCostList(service billing.BillingService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	costs, err := service.GetCloudCosts(params)
	if err != nil {
		return fmt.Errorf("Error retrieving cloud costs: %s", err)
	}

	return output.CommandOutput(cmd, OutputBillingCloudCostsProvider(costs))
}
