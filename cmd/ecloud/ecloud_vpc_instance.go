package ecloud

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
	"github.com/ukfast/cli/internal/pkg/helper"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func ecloudVPCInstanceRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "instance",
		Short: "sub-commands relating to VPC instances",
	}

	// Child commands
	cmd.AddCommand(ecloudVPCInstanceListCmd(f))

	return cmd
}

func ecloudVPCInstanceListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Lists VPC instances",
		Long:    "This command lists VPC instances",
		Example: "ukfast ecloud vpc instance list net-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing VPC")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudVPCInstanceList),
	}

	cmd.Flags().String("name", "", "Instance name for filtering")

	return cmd
}

func ecloudVPCInstanceList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd, helper.NewStringFilterFlagOption("name", "name"))
	if err != nil {
		return err
	}

	instances, err := service.GetVPCInstances(args[0], params)
	if err != nil {
		return fmt.Errorf("Error retrieving VPC instances: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudInstancesProvider(instances))
}
