package ecloud_v2

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
	"github.com/ukfast/cli/internal/pkg/helper"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func ecloudInstanceRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "instance",
		Short: "sub-commands relating to instances",
	}

	// Child commands
	cmd.AddCommand(ecloudInstanceListCmd(f))
	cmd.AddCommand(ecloudInstanceShowCmd(f))

	return cmd
}

func ecloudInstanceListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Lists instances",
		Long:    "This command lists instances",
		Example: "ukfast ecloud instance list",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ecloudInstanceList(c.ECloudService(), cmd, args)
		},
	}

	cmd.Flags().String("name", "", "Instance name for filtering")

	return cmd
}

func ecloudInstanceList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	helper.HydrateAPIRequestParametersWithStringFilterFlag(&params, cmd, helper.NewStringFilterFlag("name", "name"))

	instances, err := service.GetInstances(params)
	if err != nil {
		return fmt.Errorf("Error retrieving instances: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudInstancesProvider(instances))
}

func ecloudInstanceShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <instance: id>...",
		Short:   "Shows a instance",
		Long:    "This command shows one or more instances",
		Example: "ukfast ecloud instance show i-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing instance")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ecloudInstanceShow(c.ECloudService(), cmd, args)
		},
	}
}

func ecloudInstanceShow(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	var instances []ecloud.Instance
	for _, arg := range args {
		instance, err := service.GetInstance(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving instance [%s]: %s", arg, err)
			continue
		}

		instances = append(instances, instance)
	}

	return output.CommandOutput(cmd, OutputECloudInstancesProvider(instances))
}
