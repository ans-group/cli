package cmd

import (
	"errors"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func ecloudPodRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pod",
		Short: "sub-commands relating to pods",
	}

	// Child commands
	cmd.AddCommand(ecloudPodListCmd())
	cmd.AddCommand(ecloudPodShowCmd())

	// Child root commands
	cmd.AddCommand(ecloudPodTemplateRootCmd())
	cmd.AddCommand(ecloudPodApplianceRootCmd())

	return cmd
}

func ecloudPodListCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists pods",
		Long:    "This command lists pods",
		Example: "ukfast ecloud pod list",
		Run: func(cmd *cobra.Command, args []string) {
			ecloudPodList(getClient().ECloudService(), cmd, args)
		},
	}
}

func ecloudPodList(service ecloud.ECloudService, cmd *cobra.Command, args []string) {
	params, err := GetAPIRequestParametersFromFlags()
	if err != nil {
		output.Fatal(err.Error())
		return
	}

	pods, err := service.GetPods(params)
	if err != nil {
		output.Fatalf("Error retrieving pods: %s", err)
		return
	}

	outputECloudPods(pods)
}

func ecloudPodShowCmd() *cobra.Command {
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
		Run: func(cmd *cobra.Command, args []string) {
			ecloudPodShow(getClient().ECloudService(), cmd, args)
		},
	}
}

func ecloudPodShow(service ecloud.ECloudService, cmd *cobra.Command, args []string) {
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

	outputECloudPods(pods)
}
