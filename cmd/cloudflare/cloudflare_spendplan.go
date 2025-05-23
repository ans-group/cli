package cloudflare

import (
	"fmt"

	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/helper"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/cloudflare"
	"github.com/spf13/cobra"
)

func cloudflareSpendPlanRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "spendplan",
		Short: "sub-commands relating to spend plans",
	}

	// Child commands
	cmd.AddCommand(cloudflareSpendPlanListCmd(f))

	return cmd
}

func cloudflareSpendPlanListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists spend plans",
		Long:    "This command lists spend plans",
		Example: "ans cloudflare spendplan list",
		RunE:    cloudflareCobraRunEFunc(f, cloudflareSpendPlanList),
	}
}

func cloudflareSpendPlanList(service cloudflare.CloudflareService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	plans, err := service.GetSpendPlans(params)
	if err != nil {
		return fmt.Errorf("Error retrieving spend plans: %s", err)
	}

	return output.CommandOutput(cmd, SpendPlanCollection(plans))
}
