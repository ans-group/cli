package ecloud

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
	"github.com/ukfast/cli/internal/pkg/helper"
	"github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func ECloudRootCmd(f factory.ClientFactory, fs afero.Fs) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ecloud",
		Short: "Commands relating to eCloud service",
	}

	// Child root commands
	v1envset := len(os.Getenv("UKF_ECLOUD")) > 0
	vpcEnvSet := len(os.Getenv("UKF_ECLOUD_VPC")) > 0

	if v1envset || !vpcEnvSet {
		cmd.AddCommand(ecloudVirtualMachineRootCmd(f))
		cmd.AddCommand(ecloudSolutionRootCmd(f))
		cmd.AddCommand(ecloudSiteRootCmd(f))
		cmd.AddCommand(ecloudV1HostRootCmd(f))
		cmd.AddCommand(ecloudFirewallRootCmd(f))
		cmd.AddCommand(ecloudPodRootCmd(f))
		cmd.AddCommand(ecloudDatastoreRootCmd(f))
		cmd.AddCommand(ecloudApplianceRootCmd(f))
		cmd.AddCommand(ecloudCreditRootCmd(f))
	}
	// -- eCloud VPC
	if vpcEnvSet || !v1envset {
		cmd.AddCommand(ecloudAvailabilityZoneRootCmd(f))
		cmd.AddCommand(ecloudDHCPRootCmd(f))
		cmd.AddCommand(ecloudFirewallPolicyRootCmd(f))
		cmd.AddCommand(ecloudFirewallRuleRootCmd(f))
		cmd.AddCommand(ecloudFirewallRulePortRootCmd(f))
		cmd.AddCommand(ecloudFloatingIPRootCmd(f))
		cmd.AddCommand(ecloudHostRootCmd(f))
		cmd.AddCommand(ecloudHostGroupRootCmd(f))
		cmd.AddCommand(ecloudHostSpecRootCmd(f))
		cmd.AddCommand(ecloudImageRootCmd(f))
		cmd.AddCommand(ecloudInstanceRootCmd(f))
		cmd.AddCommand(ecloudNetworkRootCmd(f))
		cmd.AddCommand(ecloudNetworkPolicyRootCmd(f))
		cmd.AddCommand(ecloudNetworkRuleRootCmd(f))
		cmd.AddCommand(ecloudNetworkRulePortRootCmd(f))
		cmd.AddCommand(ecloudNICRootCmd(f))
		cmd.AddCommand(ecloudRegionRootCmd(f))
		cmd.AddCommand(ecloudRouterRootCmd(f))
		cmd.AddCommand(ecloudRouterThroughputRootCmd(f))
		cmd.AddCommand(ecloudSSHKeyPairRootCmd(f, fs))
		cmd.AddCommand(ecloudTaskRootCmd(f))
		cmd.AddCommand(ecloudVolumeRootCmd(f))
		cmd.AddCommand(ecloudVPCRootCmd(f))
		cmd.AddCommand(ecloudVPNEndpointRootCmd(f))
		cmd.AddCommand(ecloudVPNProfileGroupRootCmd(f))
		cmd.AddCommand(ecloudVPNServiceRootCmd(f))
		cmd.AddCommand(ecloudVPNSessionRootCmd(f))
		cmd.AddCommand(ecloudVolumeGroupRootCmd(f))
	}

	return cmd
}

// GetCreateTagRequestFromStringArrayFlag returns an array of CreateTagRequest structs from given tag string array flag
func GetCreateTagRequestFromStringArrayFlag(tagsFlag []string) ([]ecloud.CreateTagRequest, error) {
	var tags []ecloud.CreateTagRequest
	for _, tagFlag := range tagsFlag {
		key, value, err := GetKeyValueFromStringFlag(tagFlag)
		if err != nil {
			return tags, err
		}

		tags = append(tags, ecloud.CreateTagRequest{Key: key, Value: value})
	}

	return tags, nil
}

// GetCreateVirtualMachineRequestParameterFromStringArrayFlag returns an array of CreateVirtualMachineRequestParameter structs from given string array flag
func GetCreateVirtualMachineRequestParameterFromStringArrayFlag(parametersFlag []string) ([]ecloud.CreateVirtualMachineRequestParameter, error) {
	var parameters []ecloud.CreateVirtualMachineRequestParameter
	for _, parameterFlag := range parametersFlag {
		key, value, err := GetKeyValueFromStringFlag(parameterFlag)
		if err != nil {
			return parameters, err
		}

		parameters = append(parameters, ecloud.CreateVirtualMachineRequestParameter{Key: key, Value: value})
	}

	return parameters, nil
}

// GetKeyValueFromStringFlag returns a string map from given string flag. Expects format 'key=value'
func GetKeyValueFromStringFlag(flag string) (key, value string, err error) {
	if flag == "" {
		return key, value, errors.New("Missing key/value")
	}

	parts := strings.Split(flag, "=")
	if len(parts) < 2 || len(parts) > 2 {
		return key, value, errors.New("Invalid format, expecting: key=value")
	}
	if parts[0] == "" {
		return key, value, errors.New("Missing key")
	}
	if parts[1] == "" {
		return key, value, errors.New("Missing value")
	}

	return parts[0], parts[1], nil
}

// SolutionTemplateExistsWaitFunc returns WaitFunc for waiting for a template to exist
func SolutionTemplateExistsWaitFunc(service ecloud.ECloudService, solutionID int, templateName string, exists bool) helper.WaitFunc {
	return func() (finished bool, err error) {
		_, err = service.GetSolutionTemplate(solutionID, templateName)
		if err != nil {
			if _, ok := err.(*ecloud.TemplateNotFoundError); ok {
				return (exists == false), nil
			}

			return false, fmt.Errorf("Failed to retrieve solution template [%s]: %s", templateName, err.Error())
		}

		return (exists == true), nil
	}
}

// PodTemplateExistsWaitFunc returns WaitFunc for waiting for a template to exist
func PodTemplateExistsWaitFunc(service ecloud.ECloudService, podID int, templateName string, exists bool) helper.WaitFunc {
	return func() (finished bool, err error) {
		_, err = service.GetPodTemplate(podID, templateName)
		if err != nil {
			if _, ok := err.(*ecloud.TemplateNotFoundError); ok {
				return (exists == false), nil
			}

			return false, fmt.Errorf("Failed to retrieve pod template [%s]: %s", templateName, err.Error())
		}

		return (exists == true), nil
	}
}

type GetResourceSyncStatusFunc func() (ecloud.SyncStatus, error)

func ResourceSyncStatusWaitFunc(fn GetResourceSyncStatusFunc, expectedStatus ecloud.SyncStatus) helper.WaitFunc {
	return func() (finished bool, err error) {
		status, err := fn()
		if err != nil {
			return false, fmt.Errorf("Failed to retrieve status for resource: %s", err)
		}
		if status == ecloud.SyncStatusFailed {
			return false, fmt.Errorf("Resource in [%s] state", ecloud.SyncStatusFailed.String())
		}
		if status == expectedStatus {
			return true, nil
		}

		return false, nil
	}
}

func WaitForTask(service ecloud.ECloudService, taskID string, expectedStatus ecloud.TaskStatus) error {
	return helper.WaitForCommandWithStatus(TaskStatusWaitFunc(service, taskID, expectedStatus), fmt.Sprintf(" Waiting for task [%s] to reach status [%s]", taskID, expectedStatus))
}

func TaskStatusWaitFunc(service ecloud.ECloudService, taskID string, expectedStatus ecloud.TaskStatus) helper.WaitFunc {
	return func() (finished bool, err error) {
		task, err := service.GetTask(taskID)
		if err != nil {
			return false, fmt.Errorf("Failed to retrieve task status: %s", err)
		}
		if task.Status == ecloud.TaskStatusFailed {
			return false, fmt.Errorf("Task in [%s] state", ecloud.TaskStatusFailed.String())
		}
		if task.Status == expectedStatus {
			return true, nil
		}

		return false, nil
	}
}

type ecloudServiceCobraRunEFunc func(service ecloud.ECloudService, cmd *cobra.Command, args []string) error

func ecloudCobraRunEFunc(f factory.ClientFactory, rf ecloudServiceCobraRunEFunc) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		c, err := f.NewClient()
		if err != nil {
			return err
		}

		return rf(c.ECloudService(), cmd, args)
	}
}
