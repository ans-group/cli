package cmd

import (
	"errors"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func ecloudHostRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "host",
		Short: "sub-commands relating to hosts",
	}

	// Child commands
	cmd.AddCommand(ecloudHostListCmd())
	cmd.AddCommand(ecloudHostShowCmd())

	return cmd
}

func ecloudHostListCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists hosts",
		Long:    "This command lists hosts",
		Example: "ukfast ecloud host list",
		Run: func(cmd *cobra.Command, args []string) {
			ecloudHostList(getClient().ECloudService(), cmd, args)
		},
	}
}

func ecloudHostList(service ecloud.ECloudService, cmd *cobra.Command, args []string) {
	params, err := GetAPIRequestParametersFromFlags()
	if err != nil {
		output.Fatal(err.Error())
		return
	}

	hosts, err := service.GetHosts(params)
	if err != nil {
		output.Fatalf("Error retrieving hosts: %s", err)
		return
	}

	outputECloudHosts(hosts)
}

func ecloudHostShowCmd() *cobra.Command {
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
		Run: func(cmd *cobra.Command, args []string) {
			ecloudHostShow(getClient().ECloudService(), cmd, args)
		},
	}
}

func ecloudHostShow(service ecloud.ECloudService, cmd *cobra.Command, args []string) {
	var hosts []ecloud.Host
	for _, arg := range args {
		hostID, err := strconv.Atoi(arg)
		if err != nil {
			OutputWithErrorLevelf("Invalid host ID [%s]", arg)
			continue
		}

		host, err := service.GetHost(hostID)
		if err != nil {
			OutputWithErrorLevelf("Error retrieving host [%s]: %s", arg, err)
			continue
		}

		hosts = append(hosts, host)
	}

	outputECloudHosts(hosts)
}
