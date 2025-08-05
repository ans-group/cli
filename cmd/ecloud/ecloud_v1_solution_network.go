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

func ecloudSolutionNetworkRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "network",
		Short: "sub-commands relating to solution networks",
	}

	// Child commands
	cmd.AddCommand(ecloudSolutionNetworkListCmd(f))

	return cmd
}

func ecloudSolutionNetworkListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists solution networks",
		Long:    "This command lists solution networks",
		Example: "ans ecloud solution network list 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("missing solution")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ecloudSolutionNetworkList(c.ECloudService(), cmd, args)
		},
	}
}

func ecloudSolutionNetworkList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	solutionID, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid solution ID [%s]", args[0])
	}

	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	networks, err := service.GetSolutionNetworks(solutionID, params)
	if err != nil {
		return fmt.Errorf("error retrieving solution networks: %s", err)
	}

	return output.CommandOutput(cmd, V1NetworkCollection(networks))
}
