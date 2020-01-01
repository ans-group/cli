package cmd

import (
	"errors"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/helper"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/storage"
)

func storageHostRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "host",
		Short: "sub-commands relating to hosts",
	}

	// Child commands
	cmd.AddCommand(storageHostListCmd())
	cmd.AddCommand(storageHostShowCmd())

	return cmd
}

func storageHostListCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists hosts",
		Long:    "This command lists hosts",
		Example: "ukfast storage host list",
		Run: func(cmd *cobra.Command, args []string) {
			storageHostList(getClient().StorageService(), cmd, args)
		},
	}
}

func storageHostList(service storage.StorageService, cmd *cobra.Command, args []string) {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		output.Fatal(err.Error())
		return
	}

	hosts, err := service.GetHosts(params)
	if err != nil {
		output.Fatalf("Error retrieving hosts: %s", err)
		return
	}

	outputStorageHosts(hosts)
}

func storageHostShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "show <host: id>...",
		Short:   "Shows a host",
		Long:    "This command shows one or more hosts",
		Example: "ukfast storage host show 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing host")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			storageHostShow(getClient().StorageService(), cmd, args)
		},
	}
}

func storageHostShow(service storage.StorageService, cmd *cobra.Command, args []string) {
	var hosts []storage.Host
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

	outputStorageHosts(hosts)
}
