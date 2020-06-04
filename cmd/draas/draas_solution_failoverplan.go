package draas

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
	"github.com/ukfast/cli/internal/pkg/helper"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/connection"
	"github.com/ukfast/sdk-go/pkg/service/draas"
)

func draasSolutionFailoverPlanRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "failoverplan",
		Short: "sub-commands relating to solution failover plans",
	}

	// Child commands
	cmd.AddCommand(draasSolutionFailoverPlanListCmd(f))
	cmd.AddCommand(draasSolutionFailoverPlanShowCmd(f))
	cmd.AddCommand(draasSolutionFailoverPlanStartCmd(f))
	cmd.AddCommand(draasSolutionFailoverPlanStopCmd(f))

	return cmd
}

func draasSolutionFailoverPlanListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list <solution: id>",
		Short:   "Lists solution failover plans",
		Long:    "This command lists solution failover plan",
		Example: "ukfast draas solution failoverplan list 00000000-0000-0000-0000-000000000000",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing solution")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return draasSolutionFailoverPlanList(c.DRaaSService(), cmd, args)
		},
	}
}

func draasSolutionFailoverPlanList(service draas.DRaaSService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	plans, err := service.GetSolutionFailoverPlans(args[0], params)
	if err != nil {
		return fmt.Errorf("Error retrieving solution failover plans: %s", err)
	}

	return output.CommandOutput(cmd, OutputDRaaSFailoverPlansProvider(plans))
}

func draasSolutionFailoverPlanShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <solution: id> <failoverplan: id>...",
		Short:   "Shows solution failover plans",
		Long:    "This command shows a solution failover plan",
		Example: "ukfast draas solution failoverplan show 00000000-0000-0000-0000-000000000000 00000000-0000-0000-0000-000000000001",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing solution")
			}
			if len(args) < 2 {
				return errors.New("Missing failover plan")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return draasSolutionFailoverPlanShow(c.DRaaSService(), cmd, args)
		},
	}
}

func draasSolutionFailoverPlanShow(service draas.DRaaSService, cmd *cobra.Command, args []string) error {
	var plans []draas.FailoverPlan

	for _, arg := range args[1:] {
		plan, err := service.GetSolutionFailoverPlan(args[0], arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving solution failover plan [%s]: %s", arg, err.Error())
			continue
		}

		plans = append(plans, plan)
	}

	return output.CommandOutput(cmd, OutputDRaaSFailoverPlansProvider(plans))
}

func draasSolutionFailoverPlanStartCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "start <solution: id> <failoverplan: id>...",
		Short:   "Starts solution failover plan",
		Long:    "This command starts one or more solution failover plans",
		Example: "ukfast draas solution failoverplan start 00000000-0000-0000-0000-000000000000 00000000-0000-0000-0000-000000000001",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing solution")
			}
			if len(args) < 2 {
				return errors.New("Missing failover plan")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			draasSolutionFailoverPlanStart(c.DRaaSService(), cmd, args)
			return nil
		},
	}

	cmd.Flags().String("date", "", "Indicates failover plan should be started at specified date/time")

	return cmd
}

func draasSolutionFailoverPlanStart(service draas.DRaaSService, cmd *cobra.Command, args []string) {
	req := draas.StartFailoverPlanRequest{}

	if cmd.Flags().Changed("date") {
		date, _ := cmd.Flags().GetString("date")
		req.StartDate = connection.DateTime(date)
	}

	for _, arg := range args[1:] {
		err := service.StartSolutionFailoverPlan(args[0], arg, req)
		if err != nil {
			output.OutputWithErrorLevelf("Error starting solution failover plan [%s]: %s", arg, err.Error())
			continue
		}
	}
}

func draasSolutionFailoverPlanStopCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "stop <solution: id> <failoverplan: id>...",
		Short:   "Stops solution failover plan",
		Long:    "This command stops one or more solution failover plans",
		Example: "ukfast draas solution failoverplan stop 00000000-0000-0000-0000-000000000000 00000000-0000-0000-0000-000000000001",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing solution")
			}
			if len(args) < 2 {
				return errors.New("Missing failover plan")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			draasSolutionFailoverPlanStop(c.DRaaSService(), cmd, args)
			return nil
		},
	}
}

func draasSolutionFailoverPlanStop(service draas.DRaaSService, cmd *cobra.Command, args []string) {
	for _, arg := range args[1:] {
		err := service.StopSolutionFailoverPlan(args[0], arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error stopping solution failover plan [%s]: %s", arg, err.Error())
			continue
		}
	}
}
