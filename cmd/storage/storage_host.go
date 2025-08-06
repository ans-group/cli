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

func storageHostRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "host",
		Short: "sub-commands relating to hosts",
	}

	// Child commands
	cmd.AddCommand(storageHostListCmd(f))
	cmd.AddCommand(storageHostShowCmd(f))

	return cmd
}

func storageHostListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists hosts",
		Long:    "This command lists hosts",
		Example: "ans storage host list",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return storageHostList(c.StorageService(), cmd, args)
		},
	}
}

func storageHostList(service storage.StorageService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	hosts, err := service.GetHosts(params)
	if err != nil {
		return fmt.Errorf("error retrieving hosts: %s", err)
	}

	return output.CommandOutput(cmd, HostCollection(hosts))
}

func storageHostShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <host: id>...",
		Short:   "Shows a host",
		Long:    "This command shows one or more hosts",
		Example: "ans storage host show 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("missing host")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return storageHostShow(c.StorageService(), cmd, args)
		},
	}
}

func storageHostShow(service storage.StorageService, cmd *cobra.Command, args []string) error {
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

	return output.CommandOutput(cmd, HostCollection(hosts))
}
