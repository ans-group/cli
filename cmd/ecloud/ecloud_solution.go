package ecloud

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
	"github.com/ukfast/cli/internal/pkg/helper"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/ptr"
	"github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func ecloudSolutionRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "solution",
		Short: "sub-commands relating to solutions",
	}

	// Child commands
	cmd.AddCommand(ecloudSolutionListCmd(f))
	cmd.AddCommand(ecloudSolutionShowCmd(f))
	cmd.AddCommand(ecloudSolutionUpdateCmd(f))

	// Child root commands
	cmd.AddCommand(ecloudSolutionVirtualMachineRootCmd(f))
	cmd.AddCommand(ecloudSolutionTagRootCmd(f))
	cmd.AddCommand(ecloudSolutionSiteRootCmd(f))
	cmd.AddCommand(ecloudSolutionNetworkRootCmd(f))
	cmd.AddCommand(ecloudSolutionHostRootCmd(f))
	cmd.AddCommand(ecloudSolutionFirewallRootCmd(f))
	cmd.AddCommand(ecloudSolutionTemplateRootCmd(f))
	cmd.AddCommand(ecloudSolutionDatastoreRootCmd(f))

	return cmd
}

func ecloudSolutionListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Lists solutions",
		Long:    "This command lists solutions",
		Example: "ukfast ecloud solution list",
		Run: func(cmd *cobra.Command, args []string) {
			ecloudSolutionList(f.NewClient().ECloudService(), cmd, args)
		},
	}

	cmd.Flags().String("name", "", "Solution name for filtering")

	return cmd
}

func ecloudSolutionList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	if cmd.Flags().Changed("name") {
		filterName, _ := cmd.Flags().GetString("name")
		params.WithFilter(helper.GetFilteringInferOperator("name", filterName))
	}

	solutions, err := service.GetSolutions(params)
	if err != nil {
		return fmt.Errorf("Error retrieving solutions: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudSolutionsProvider(solutions))
}

func ecloudSolutionShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <solution: id>...",
		Short:   "Shows a solution",
		Long:    "This command shows one or more solutions",
		Example: "ukfast ecloud solution show 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing solution")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return ecloudSolutionShow(f.NewClient().ECloudService(), cmd, args)
		},
	}
}

func ecloudSolutionShow(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	var solutions []ecloud.Solution
	for _, arg := range args {
		solutionID, err := strconv.Atoi(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Invalid solution ID [%s]", arg)
			continue
		}

		solution, err := service.GetSolution(solutionID)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving solution [%s]: %s", arg, err)
			continue
		}

		solutions = append(solutions, solution)
	}

	return output.CommandOutput(cmd, OutputECloudSolutionsProvider(solutions))
}

func ecloudSolutionUpdateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <solution: id>",
		Short:   "Updates a solution",
		Long:    "This command updates a solution",
		Example: "ukfast ecloud solution update 123 --name \"new name\"",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing solution")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return ecloudSolutionUpdate(f.NewClient().ECloudService(), cmd, args)
		},
	}

	cmd.Flags().String("name", "", "Name of solution")

	return cmd
}

func ecloudSolutionUpdate(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	solutionID, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("Invalid solution ID [%s]", args[0])
	}

	patchRequest := ecloud.PatchSolutionRequest{}

	if cmd.Flags().Changed("name") {
		solutionName, _ := cmd.Flags().GetString("name")
		patchRequest.Name = ptr.String(solutionName)
	}

	id, err := service.PatchSolution(solutionID, patchRequest)
	if err != nil {
		return fmt.Errorf("Error updating solution: %s", err)
	}

	solution, err := service.GetSolution(id)
	if err != nil {
		return fmt.Errorf("Error retrieving updated solution: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudSolutionsProvider([]ecloud.Solution{solution}))
}
