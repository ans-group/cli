package ecloud

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/helper"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/spf13/cobra"
)

func ecloudV1HostRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "v1host",
		Short: "sub-commands relating to hosts",
	}

	// Child commands
	cmd.AddCommand(ecloudV1HostListCmd(f))
	cmd.AddCommand(ecloudV1HostShowCmd(f))

	return cmd
}

func ecloudV1HostListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists hosts",
		Long:    "This command lists hosts",
		Example: "ukfast ecloud v1host list",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ecloudV1HostList(c.ECloudService(), cmd, args)
		},
	}
}

func ecloudV1HostList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	hosts, err := service.GetV1Hosts(params)
	if err != nil {
		return fmt.Errorf("Error retrieving hosts: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudV1HostsProvider(hosts))
}

func ecloudV1HostShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <host: id>...",
		Short:   "Shows a host",
		Long:    "This command shows one or more hosts",
		Example: "ukfast ecloud v1host 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing host")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ecloudV1HostShow(c.ECloudService(), cmd, args)
		},
	}
}

func ecloudV1HostShow(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	var hosts []ecloud.V1Host
	for _, arg := range args {
		hostID, err := strconv.Atoi(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Invalid host ID [%s]", arg)
			continue
		}

		host, err := service.GetV1Host(hostID)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving host [%s]: %s", arg, err)
			continue
		}

		hosts = append(hosts, host)
	}

	return output.CommandOutput(cmd, OutputECloudV1HostsProvider(hosts))
}
