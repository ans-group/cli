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

func draasSolutionComputeResourceRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "computeresource",
		Short: "sub-commands relating to solution compute resources",
	}

	// Child commands
	cmd.AddCommand(draasSolutionComputeResourceListCmd(f))
	cmd.AddCommand(draasSolutionComputeResourceShowCmd(f))

	return cmd
}

func draasSolutionComputeResourceListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list <solution: id>",
		Short:   "Lists solution compute resources",
		Long:    "This command lists solution compute resource",
		Example: "ukfast draas solution computeresource list 00000000-0000-0000-0000-000000000000",
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

			return draasSolutionComputeResourceList(c.DRaaSService(), cmd, args)
		},
	}
}

func draasSolutionComputeResourceList(service draas.DRaaSService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	plans, err := service.GetSolutionComputeResources(args[0], params)
	if err != nil {
		return fmt.Errorf("Error retrieving solution compute resources: %s", err)
	}

	return output.CommandOutput(cmd, OutputDRaaSComputeResourcesProvider(plans))
}

func draasSolutionComputeResourceShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <solution: id> <computeresource: id>...",
		Short:   "Shows solution compute resources",
		Long:    "This command shows a solution compute resource",
		Example: "ukfast draas solution computeresource show 00000000-0000-0000-0000-000000000000 00000000-0000-0000-0000-000000000001",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing solution")
			}
			if len(args) < 2 {
				return errors.New("Missing compute resource")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return draasSolutionComputeResourceShow(c.DRaaSService(), cmd, args)
		},
	}
}

func draasSolutionComputeResourceShow(service draas.DRaaSService, cmd *cobra.Command, args []string) error {
	var plans []draas.ComputeResource

	for _, arg := range args[1:] {
		plan, err := service.GetSolutionComputeResource(args[0], arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving solution compute resource [%s]: %s", arg, err.Error())
			continue
		}

		plans = append(plans, plan)
	}

	return output.CommandOutput(cmd, OutputDRaaSComputeResourcesProvider(plans))
}
