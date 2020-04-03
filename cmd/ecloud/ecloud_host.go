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

func ecloudHostRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "host",
		Short: "sub-commands relating to hosts",
	}

	// Child commands
	cmd.AddCommand(ecloudHostListCmd(f))
	cmd.AddCommand(ecloudHostShowCmd(f))

	return cmd
}

func ecloudHostListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists hosts",
		Long:    "This command lists hosts",
		Example: "ukfast ecloud host list",
		Run: func(cmd *cobra.Command, args []string) {
			ecloudHostList(f.NewClient().ECloudService(), cmd, args)
		},
	}
}

func ecloudHostList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	hosts, err := service.GetHosts(params)
	if err != nil {
		return fmt.Errorf("Error retrieving hosts: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudHostsProvider(hosts))
}

func ecloudHostShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <host: id>...",
		Short:   "Shows a host",
		Long:    "This command shows one or more hosts",
		Example: "ukfast ecloud vm host 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing host")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return ecloudHostShow(f.NewClient().ECloudService(), cmd, args)
		},
	}
}

func ecloudHostShow(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	var hosts []ecloud.Host
	for _, arg := range args {
		hostID, err := strconv.Atoi(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Invalid host ID [%s]", arg)
			continue
		}

		host, err := service.GetHost(hostID)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving host [%s]: %s", arg, err)
			continue
		}

		hosts = append(hosts, host)
	}

	return output.CommandOutput(cmd, OutputECloudHostsProvider(hosts))
}
