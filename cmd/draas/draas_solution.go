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

func draasSolutionRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "solution",
		Short: "sub-commands relating to solutions",
	}

	// Child commands
	cmd.AddCommand(draasSolutionListCmd(f))
	cmd.AddCommand(draasSolutionShowCmd(f))
	cmd.AddCommand(draasSolutionUpdateCmd(f))

	// Child root commands
	cmd.AddCommand(draasSolutionBackupResourceRootCmd(f))
	cmd.AddCommand(draasSolutionBackupServiceRootCmd(f))
	cmd.AddCommand(draasSolutionFailoverPlanRootCmd(f))
	cmd.AddCommand(draasSolutionComputeResourceRootCmd(f))
	cmd.AddCommand(draasSolutionHardwarePlanRootCmd(f))

	return cmd
}

func draasSolutionListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists solutions",
		Long:    "This command lists solutions",
		Example: "ans draas solution list",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return draasSolutionList(c.DRaaSService(), cmd, args)
		},
	}
}

func draasSolutionList(service draas.DRaaSService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	solutions, err := service.GetSolutions(params)
	if err != nil {
		return fmt.Errorf("Error retrieving solutions: %s", err)
	}

	return output.CommandOutput(cmd, SolutionCollection(solutions))
}

func draasSolutionShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <solution: id>...",
		Short:   "Shows a solution",
		Long:    "This command shows one or more solutions",
		Example: "ans draas solution show 00000000-0000-0000-0000-000000000000",
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

			return draasSolutionShow(c.DRaaSService(), cmd, args)
		},
	}
}

func draasSolutionShow(service draas.DRaaSService, cmd *cobra.Command, args []string) error {
	var solutions []draas.Solution
	for _, arg := range args {
		solution, err := service.GetSolution(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving solution [%s]: %s", arg, err)
			continue
		}

		solutions = append(solutions, solution)
	}

	return output.CommandOutput(cmd, SolutionCollection(solutions))
}

func draasSolutionUpdateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <solution: id>",
		Short:   "Updates a solution",
		Long:    "This command updates a solution",
		Example: "ans draas solution update 00000000-0000-0000-0000-000000000000 --name test",
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

			return draasSolutionUpdate(c.DRaaSService(), cmd, args)
		},
	}

	cmd.Flags().String("name", "", "Name for solution")
	cmd.Flags().String("iops-tier", "", "IOPS tier ID")

	return cmd
}

func draasSolutionUpdate(service draas.DRaaSService, cmd *cobra.Command, args []string) error {
	patchRequest := draas.PatchSolutionRequest{}

	if cmd.Flags().Changed("name") {
		name, _ := cmd.Flags().GetString("name")
		patchRequest.Name = name
	}

	if cmd.Flags().Changed("iops-tier") {
		iopsTierID, _ := cmd.Flags().GetString("iops-tier")
		patchRequest.IOPSTierID = iopsTierID
	}

	err := service.PatchSolution(args[0], patchRequest)
	if err != nil {
		return fmt.Errorf("Error updating solution: %s", err.Error())
	}

	solution, err := service.GetSolution(args[0])
	if err != nil {
		return fmt.Errorf("Error retrieving updated solution: %s", err)
	}

	return output.CommandOutput(cmd, SolutionCollection([]draas.Solution{solution}))
}
