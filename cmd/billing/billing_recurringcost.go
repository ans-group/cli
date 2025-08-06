package billing

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/helper"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/billing"
	"github.com/spf13/cobra"
)

func billingRecurringCostRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "recurringcost",
		Short: "sub-commands relating to recurring costs",
	}

	// Child commands
	cmd.AddCommand(billingRecurringCostListCmd(f))
	cmd.AddCommand(billingRecurringCostShowCmd(f))

	return cmd
}

func billingRecurringCostListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists recurring costs",
		Long:    "This command lists recurring costs",
		Example: "ans billing recurringcost list",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return billingRecurringCostList(c.BillingService(), cmd, args)
		},
	}
}

func billingRecurringCostList(service billing.BillingService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	costs, err := service.GetRecurringCosts(params)
	if err != nil {
		return fmt.Errorf("error retrieving recurring costs: %s", err)
	}

	return output.CommandOutput(cmd, RecurringCostCollection(costs))
}

func billingRecurringCostShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <recurringcost: id>...",
		Short:   "Shows a recurring cost",
		Long:    "This command shows one or more recurring costs",
		Example: "ans billing recurringcost show 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("missing recurring cost")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return billingRecurringCostShow(c.BillingService(), cmd, args)
		},
	}
}

func billingRecurringCostShow(service billing.BillingService, cmd *cobra.Command, args []string) error {
	var costs []billing.RecurringCost
	for _, arg := range args {
		costID, err := strconv.Atoi(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Invalid recurring cost ID [%s]", arg)
			continue
		}

		cost, err := service.GetRecurringCost(costID)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving recurring cost [%s]: %s", arg, err)
			continue
		}

		costs = append(costs, cost)
	}

	return output.CommandOutput(cmd, RecurringCostCollection(costs))
}
