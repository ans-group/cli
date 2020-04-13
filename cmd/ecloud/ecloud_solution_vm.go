package ecloud

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
	"github.com/ukfast/cli/internal/pkg/helper"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func ecloudSolutionVirtualMachineRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vm",
		Short: "sub-commands relating to solution virtual machines",
	}

	// Child commands
	cmd.AddCommand(ecloudSolutionVirtualMachineListCmd(f))

	return cmd
}

func ecloudSolutionVirtualMachineListCmd(f factory.ClientFactory) *cobra.Command {
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
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ecloudSolutionVirtualMachineList(c.ECloudService(), cmd, args)
		},
	}
}

func ecloudSolutionVirtualMachineList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	solutionID, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("Invalid solution ID [%s]", args[0])
	}

	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	vms, err := service.GetSolutionVirtualMachines(solutionID, params)
	if err != nil {
		return fmt.Errorf("Error retrieving solution virtual machines: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudVirtualMachinesProvider(vms))
}
