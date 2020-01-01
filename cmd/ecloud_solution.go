package cmd

import (
	"errors"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/helper"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/ptr"
	"github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func ecloudSolutionRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "solution",
		Short: "sub-commands relating to solutions",
	}

	// Child commands
	cmd.AddCommand(ecloudSolutionListCmd())
	cmd.AddCommand(ecloudSolutionShowCmd())
	cmd.AddCommand(ecloudSolutionUpdateCmd())

	// Child root commands
	cmd.AddCommand(ecloudSolutionVirtualMachineRootCmd())
	cmd.AddCommand(ecloudSolutionTagRootCmd())
	cmd.AddCommand(ecloudSolutionSiteRootCmd())
	cmd.AddCommand(ecloudSolutionNetworkRootCmd())
	cmd.AddCommand(ecloudSolutionHostRootCmd())
	cmd.AddCommand(ecloudSolutionFirewallRootCmd())
	cmd.AddCommand(ecloudSolutionTemplateRootCmd())
	cmd.AddCommand(ecloudSolutionDatastoreRootCmd())

	return cmd
}

func ecloudSolutionListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Lists solutions",
		Long:    "This command lists solutions",
		Example: "ukfast ecloud solution list",
		Run: func(cmd *cobra.Command, args []string) {
			ecloudSolutionList(getClient().ECloudService(), cmd, args)
		},
	}

	cmd.Flags().String("name", "", "Solution name for filtering")

	return cmd
}

func ecloudSolutionList(service ecloud.ECloudService, cmd *cobra.Command, args []string) {
	params, err := GetAPIRequestParametersFromFlags()
	if err != nil {
		output.Fatal(err.Error())
		return
	}

	if cmd.Flags().Changed("name") {
		filterName, _ := cmd.Flags().GetString("name")
		params.WithFilter(helper.GetFilteringInferOperator("name", filterName))
	}

	solutions, err := service.GetSolutions(params)
	if err != nil {
		output.Fatalf("Error retrieving solutions: %s", err)
		return
	}

	outputECloudSolutions(solutions)
}

func ecloudSolutionShowCmd() *cobra.Command {
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
		Run: func(cmd *cobra.Command, args []string) {
			ecloudSolutionShow(getClient().ECloudService(), cmd, args)
		},
	}
}

func ecloudSolutionShow(service ecloud.ECloudService, cmd *cobra.Command, args []string) {
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

	outputECloudSolutions(solutions)
}

func ecloudSolutionUpdateCmd() *cobra.Command {
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
		Run: func(cmd *cobra.Command, args []string) {
			ecloudSolutionUpdate(getClient().ECloudService(), cmd, args)
		},
	}

	cmd.Flags().String("name", "", "Name of solution")

	return cmd
}

func ecloudSolutionUpdate(service ecloud.ECloudService, cmd *cobra.Command, args []string) {
	solutionID, err := strconv.Atoi(args[0])
	if err != nil {
		output.Fatalf("Invalid solution ID [%s]", args[0])
		return
	}

	patchRequest := ecloud.PatchSolutionRequest{}

	if cmd.Flags().Changed("name") {
		solutionName, _ := cmd.Flags().GetString("name")
		patchRequest.Name = ptr.String(solutionName)
	}

	id, err := service.PatchSolution(solutionID, patchRequest)
	if err != nil {
		output.Fatalf("Error updating solution: %s", err)
		return
	}

	solution, err := service.GetSolution(id)
	if err != nil {
		output.Fatalf("Error retrieving updated solution: %s", err)
		return
	}

	outputECloudSolutions([]ecloud.Solution{solution})
}
