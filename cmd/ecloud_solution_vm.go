package cmd

import (
	"errors"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func ecloudSolutionVirtualMachineRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vm",
		Short: "sub-commands relating to solution virtual machines",
	}

	// Child commands
	cmd.AddCommand(ecloudSolutionVirtualMachineListCmd())

	return cmd
}

func ecloudSolutionVirtualMachineListCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists solution virtual machines",
		Long:    "This command lists solution virtual machines",
		Example: "ukfast ecloud solution vm list 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing solution")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			ecloudSolutionVirtualMachineList(getClient().ECloudService(), cmd, args)
		},
	}
}

func ecloudSolutionVirtualMachineList(service ecloud.ECloudService, cmd *cobra.Command, args []string) {
	solutionID, err := strconv.Atoi(args[0])
	if err != nil {
		output.Fatalf("Invalid solution ID [%s]", args[0])
		return
	}

	params, err := GetAPIRequestParametersFromFlags()
	if err != nil {
		output.Fatal(err.Error())
		return
	}

	vms, err := service.GetSolutionVirtualMachines(solutionID, params)
	if err != nil {
		output.Fatalf("Error retrieving solution virtual machines: %s", err)
		return
	}

	outputECloudVirtualMachines(vms)
}
