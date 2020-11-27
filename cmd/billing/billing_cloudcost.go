package billing

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
	flaghelper "github.com/ukfast/cli/internal/pkg/helper/flag"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/billing"
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
		Example: "ukfast billing cloudcost list",
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
	params, err := flaghelper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	costs, err := service.GetCloudCosts(params)
	if err != nil {
		return fmt.Errorf("Error retrieving cloud costs: %s", err)
	}

	return output.CommandOutput(cmd, OutputBillingCloudCostsProvider(costs))
}
