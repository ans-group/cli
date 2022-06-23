package storage

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/helper"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/storage"
	"github.com/spf13/cobra"
)

func storageSolutionRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "solution",
		Short: "sub-commands relating to solutions",
	}

	// Child commands
	cmd.AddCommand(storageSolutionListCmd(f))
	cmd.AddCommand(storageSolutionShowCmd(f))

	return cmd
}

func storageSolutionListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists solutions",
		Long:    "This command lists solutions",
		Example: "ans storage solution list",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return storageSolutionList(c.StorageService(), cmd, args)
		},
	}
}

func storageSolutionList(service storage.StorageService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	solutions, err := service.GetSolutions(params)
	if err != nil {
		return fmt.Errorf("Error retrieving solutions: %s", err)
	}

	return output.CommandOutput(cmd, OutputStorageSolutionsProvider(solutions))
}

func storageSolutionShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <solution: id>...",
		Short:   "Shows a solution",
		Long:    "This command shows one or more solutions",
		Example: "ans storage solution show 123",
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

			return storageSolutionShow(c.StorageService(), cmd, args)
		},
	}
}

func storageSolutionShow(service storage.StorageService, cmd *cobra.Command, args []string) error {
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

	return output.CommandOutput(cmd, OutputStorageSolutionsProvider(solutions))
}
