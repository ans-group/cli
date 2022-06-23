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

func ecloudPodRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pod",
		Short: "sub-commands relating to pods",
	}

	// Child commands
	cmd.AddCommand(ecloudPodListCmd(f))
	cmd.AddCommand(ecloudPodShowCmd(f))

	// Child root commands
	cmd.AddCommand(ecloudPodTemplateRootCmd(f))
	cmd.AddCommand(ecloudPodApplianceRootCmd(f))

	return cmd
}

func ecloudPodListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists pods",
		Long:    "This command lists pods",
		Example: "ukfast ecloud pod list",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ecloudPodList(c.ECloudService(), cmd, args)
		},
	}
}

func ecloudPodList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	pods, err := service.GetPods(params)
	if err != nil {
		return fmt.Errorf("Error retrieving pods: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudPodsProvider(pods))
}

func ecloudPodShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <pod: id>...",
		Short:   "Shows a pod",
		Long:    "This command shows one or more pods",
		Example: "ukfast ecloud vm pod 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing pod")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ecloudPodShow(c.ECloudService(), cmd, args)
		},
	}
}

func ecloudPodShow(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	var pods []ecloud.Pod
	for _, arg := range args {
		podID, err := strconv.Atoi(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Invalid pod ID [%s]", arg)
			continue
		}

		pod, err := service.GetPod(podID)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving pod [%s]: %s", arg, err)
			continue
		}

		pods = append(pods, pod)
	}

	return output.CommandOutput(cmd, OutputECloudPodsProvider(pods))
}
