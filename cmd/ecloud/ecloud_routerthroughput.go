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

func ecloudRouterThroughputRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "routerthroughput",
		Short: "sub-commands relating to router throughputs",
	}

	// Child commands
	cmd.AddCommand(ecloudRouterThroughputListCmd(f))
	cmd.AddCommand(ecloudRouterThroughputShowCmd(f))

	return cmd
}

func ecloudRouterThroughputListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Lists router throughputs",
		Long:    "This command lists router throughputs",
		Example: "ukfast ecloud routerthroughput list",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ecloudRouterThroughputList(c.ECloudService(), cmd, args)
		},
	}

	cmd.Flags().String("name", "", "RouterThroughput name for filtering")

	return cmd
}

func ecloudRouterThroughputList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd, helper.NewStringFilterFlagOption("name", "name"))
	if err != nil {
		return err
	}

	throughputs, err := service.GetRouterThroughputs(params)
	if err != nil {
		return fmt.Errorf("Error retrieving router throughputs: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudRouterThroughputsProvider(throughputs))
}

func ecloudRouterThroughputShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <routerthroughput: id>...",
		Short:   "Shows a router throughput",
		Long:    "This command shows one or more router throughputs",
		Example: "ukfast ecloud routerthroughput show rtp-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing router throughput")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ecloudRouterThroughputShow(c.ECloudService(), cmd, args)
		},
	}
}

func ecloudRouterThroughputShow(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	var throughputs []ecloud.RouterThroughput
	for _, arg := range args {
		throughput, err := service.GetRouterThroughput(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving router throughput [%s]: %s", arg, err)
			continue
		}

		throughputs = append(throughputs, throughput)
	}

	return output.CommandOutput(cmd, OutputECloudRouterThroughputsProvider(throughputs))
}
