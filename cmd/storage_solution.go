package cmd

import (
	"errors"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/storage"
)

func storageSolutionRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "solution",
		Short: "sub-commands relating to solutions",
	}

	// Child commands
	cmd.AddCommand(storageSolutionListCmd())
	cmd.AddCommand(storageSolutionShowCmd())

	return cmd
}

func storageSolutionListCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists solutions",
		Long:    "This command lists solutions",
		Example: "ukfast storage solution list",
		Run: func(cmd *cobra.Command, args []string) {
			storageSolutionList(getClient().StorageService(), cmd, args)
		},
	}
}

func storageSolutionList(service storage.StorageService, cmd *cobra.Command, args []string) {
	params, err := GetAPIRequestParametersFromFlags()
	if err != nil {
		output.Fatal(err.Error())
		return
	}

	solutions, err := service.GetSolutions(params)
	if err != nil {
		output.Fatalf("Error retrieving solutions: %s", err)
		return
	}

	outputStorageSolutions(solutions)
}

func storageSolutionShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "show <solution: id>...",
		Short:   "Shows a solution",
		Long:    "This command shows one or more solutions",
		Example: "ukfast storage solution show 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing solution")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			storageSolutionShow(getClient().StorageService(), cmd, args)
		},
	}
}

func storageSolutionShow(service storage.StorageService, cmd *cobra.Command, args []string) {
	var solutions []storage.Solution
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

	outputStorageSolutions(solutions)
}
