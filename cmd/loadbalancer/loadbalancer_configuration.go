package loadbalancer

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
	"github.com/ukfast/cli/internal/pkg/helper"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/ptr"
	"github.com/ukfast/sdk-go/pkg/service/loadbalancer"
)

func loadbalancerConfigurationRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "configuration",
		Short: "sub-commands relating to configurations",
	}

	// Child commands
	cmd.AddCommand(loadbalancerConfigurationListCmd(f))
	cmd.AddCommand(loadbalancerConfigurationShowCmd(f))
	cmd.AddCommand(loadbalancerConfigurationCreateCmd(f))
	cmd.AddCommand(loadbalancerConfigurationUpdateCmd(f))
	cmd.AddCommand(loadbalancerConfigurationDeleteCmd(f))

	return cmd
}

func loadbalancerConfigurationListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Lists configurations",
		Long:    "This command lists configurations",
		Example: "ukfast loadbalancer configuration list",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return loadbalancerConfigurationList(c.LoadBalancerService(), cmd, args)
		},
	}

	cmd.Flags().String("name", "", "Configuration name for filtering")

	return cmd
}

func loadbalancerConfigurationList(service loadbalancer.LoadBalancerService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	helper.HydrateAPIRequestParametersWithStringFilterFlag(&params, cmd, helper.NewStringFilterFlag("name", "name"))

	configurations, err := service.GetConfigurations(params)
	if err != nil {
		return fmt.Errorf("Error retrieving configurations: %s", err)
	}

	return output.CommandOutput(cmd, OutputLoadBalancerConfigurationsProvider(configurations))
}

func loadbalancerConfigurationShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <configuration: id>...",
		Short:   "Shows a configuration",
		Long:    "This command shows one or more configurations",
		Example: "ukfast loadbalancer configuration show rtr-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing configuration")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return loadbalancerConfigurationShow(c.LoadBalancerService(), cmd, args)
		},
	}
}

func loadbalancerConfigurationShow(service loadbalancer.LoadBalancerService, cmd *cobra.Command, args []string) error {
	var configurations []loadbalancer.Configuration
	for _, arg := range args {
		configuration, err := service.GetConfiguration(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving configuration [%s]: %s", arg, err)
			continue
		}

		configurations = append(configurations, configuration)
	}

	return output.CommandOutput(cmd, OutputLoadBalancerConfigurationsProvider(configurations))
}

func loadbalancerConfigurationCreateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Creates a configuration",
		Long:    "This command creates a configuration",
		Example: "ukfast loadbalancer configuration create --vpc vpc-abcdef12 --az az-abcdef12",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return loadbalancerConfigurationCreate(c.LoadBalancerService(), cmd, args)
		},
	}

	// Setup flags
	cmd.Flags().String("name", "", "Name of configuration")
	cmd.Flags().String("vpc", "", "ID of VPC")
	cmd.MarkFlagRequired("vpc")
	cmd.Flags().String("az", "", "ID of availability zone")
	cmd.MarkFlagRequired("az")

	return cmd
}

func loadbalancerConfigurationCreate(service loadbalancer.LoadBalancerService, cmd *cobra.Command, args []string) error {
	createRequest := loadbalancer.CreateConfigurationRequest{}
	if cmd.Flags().Changed("name") {
		name, _ := cmd.Flags().GetString("name")
		createRequest.Name = ptr.String(name)
	}
	// createRequest.VPCID, _ = cmd.Flags().GetString("vpc")
	// createRequest.AvailabilityZoneID, _ = cmd.Flags().GetString("az")

	configurationID, err := service.CreateConfiguration(createRequest)
	if err != nil {
		return fmt.Errorf("Error creating configuration: %s", err)
	}

	configuration, err := service.GetConfiguration(configurationID)
	if err != nil {
		return fmt.Errorf("Error retrieving new configuration: %s", err)
	}

	return output.CommandOutput(cmd, OutputLoadBalancerConfigurationsProvider([]loadbalancer.Configuration{configuration}))
}

func loadbalancerConfigurationUpdateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <configuration: id>...",
		Short:   "Updates a configuration",
		Long:    "This command updates one or more configurations",
		Example: "ukfast loadbalancer configuration update rtr-abcdef12 --name \"my configuration\"",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing configuration")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return loadbalancerConfigurationUpdate(c.LoadBalancerService(), cmd, args)
		},
	}

	cmd.Flags().String("name", "", "Name of configuration")

	return cmd
}

func loadbalancerConfigurationUpdate(service loadbalancer.LoadBalancerService, cmd *cobra.Command, args []string) error {
	patchRequest := loadbalancer.PatchConfigurationRequest{}

	if cmd.Flags().Changed("name") {
		name, _ := cmd.Flags().GetString("name")
		patchRequest.Name = ptr.String(name)
	}

	var configurations []loadbalancer.Configuration
	for _, arg := range args {
		err := service.PatchConfiguration(arg, patchRequest)
		if err != nil {
			output.OutputWithErrorLevelf("Error updating configuration [%s]: %s", arg, err)
			continue
		}

		configuration, err := service.GetConfiguration(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving updated configuration [%s]: %s", arg, err)
			continue
		}

		configurations = append(configurations, configuration)
	}

	return output.CommandOutput(cmd, OutputLoadBalancerConfigurationsProvider(configurations))
}

func loadbalancerConfigurationDeleteCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "delete <configuration: id...>",
		Short:   "Removes a configuration",
		Long:    "This command removes one or more configurations",
		Example: "ukfast loadbalancer configuration delete rtr-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing configuration")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			loadbalancerConfigurationDelete(c.LoadBalancerService(), cmd, args)
			return nil
		},
	}
}

func loadbalancerConfigurationDelete(service loadbalancer.LoadBalancerService, cmd *cobra.Command, args []string) {
	for _, arg := range args {
		err := service.DeleteConfiguration(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error removing configuration [%s]: %s", arg, err)
		}
	}
}
