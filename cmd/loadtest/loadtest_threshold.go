package loadtest

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
	flaghelper "github.com/ukfast/cli/internal/pkg/helper/flag"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/ltaas"
)

func loadtestThresholdRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "threshold",
		Short: "sub-commands relating to thresholds",
	}

	// Child commands
	cmd.AddCommand(loadtestThresholdListCmd(f))
	cmd.AddCommand(loadtestThresholdShowCmd(f))

	return cmd
}

func loadtestThresholdListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists thresholds",
		Long:    "This command lists thresholds",
		Example: "ukfast loadtest threshold list",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return loadtestThresholdList(c.LTaaSService(), cmd, args)
		},
	}
}

func loadtestThresholdList(service ltaas.LTaaSService, cmd *cobra.Command, args []string) error {
	params, err := flaghelper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	thresholds, err := service.GetThresholds(params)
	if err != nil {
		return fmt.Errorf("Error retrieving thresholds: %s", err)
	}

	return output.CommandOutput(cmd, OutputLoadTestThresholdsProvider(thresholds))
}

func loadtestThresholdShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <threshold: id>...",
		Short:   "Shows a threshold",
		Long:    "This command shows one or more thresholds",
		Example: "ukfast loadtest threshold show 00000000-0000-0000-0000-000000000000",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing threshold")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return loadtestThresholdShow(c.LTaaSService(), cmd, args)
		},
	}
}

func loadtestThresholdShow(service ltaas.LTaaSService, cmd *cobra.Command, args []string) error {
	var thresholds []ltaas.Threshold
	for _, arg := range args {
		threshold, err := service.GetThreshold(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving threshold [%s]: %s", arg, err)
			continue
		}

		thresholds = append(thresholds, threshold)
	}

	return output.CommandOutput(cmd, OutputLoadTestThresholdsProvider(thresholds))
}
