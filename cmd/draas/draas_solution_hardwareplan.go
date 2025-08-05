package draas

import (
	"errors"
	"fmt"

	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/helper"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/draas"
	"github.com/spf13/cobra"
)

func draasSolutionHardwarePlanRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hardwareplan",
		Short: "sub-commands relating to solution hardware plans",
	}

	// Child commands
	cmd.AddCommand(draasSolutionHardwarePlanListCmd(f))
	cmd.AddCommand(draasSolutionHardwarePlanShowCmd(f))

	// Child root commands
	cmd.AddCommand(draasSolutionHardwarePlanReplicaRootCmd(f))

	return cmd
}

func draasSolutionHardwarePlanListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list <solution: id>",
		Short:   "Lists solution hardware plans",
		Long:    "This command lists solution hardware plan",
		Example: "ans draas solution hardwareplan list 00000000-0000-0000-0000-000000000000",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("missing solution")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return draasSolutionHardwarePlanList(c.DRaaSService(), cmd, args)
		},
	}
}

func draasSolutionHardwarePlanList(service draas.DRaaSService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	plans, err := service.GetSolutionHardwarePlans(args[0], params)
	if err != nil {
		return fmt.Errorf("error retrieving solution hardware plans: %s", err)
	}

	return output.CommandOutput(cmd, HardwarePlanCollection(plans))
}

func draasSolutionHardwarePlanShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <solution: id> <hardwareplan: id>...",
		Short:   "Shows solution hardware plans",
		Long:    "This command shows one or more solution hardware plans",
		Example: "ans draas solution hardwareplan show 00000000-0000-0000-0000-000000000000 00000000-0000-0000-0000-000000000001",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("missing solution")
			}
			if len(args) < 2 {
				return errors.New("missing hardware plan")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return draasSolutionHardwarePlanShow(c.DRaaSService(), cmd, args)
		},
	}
}

func draasSolutionHardwarePlanShow(service draas.DRaaSService, cmd *cobra.Command, args []string) error {
	var plans []draas.HardwarePlan

	for _, arg := range args[1:] {
		plan, err := service.GetSolutionHardwarePlan(args[0], arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving solution hardware plan [%s]: %s", arg, err.Error())
			continue
		}

		plans = append(plans, plan)
	}

	return output.CommandOutput(cmd, HardwarePlanCollection(plans))
}
