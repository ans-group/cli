package ecloud

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
	"github.com/ukfast/cli/internal/pkg/helper"
	"github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func ECloudRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ecloud",
		Short: "Commands relating to eCloud service",
	}

	// Child root commands
	cmd.AddCommand(ecloudVirtualMachineRootCmd(f))
	cmd.AddCommand(ecloudSolutionRootCmd(f))
	cmd.AddCommand(ecloudSiteRootCmd(f))
	cmd.AddCommand(ecloudHostRootCmd(f))
	cmd.AddCommand(ecloudFirewallRootCmd(f))
	cmd.AddCommand(ecloudPodRootCmd(f))
	cmd.AddCommand(ecloudDatastoreRootCmd(f))
	cmd.AddCommand(ecloudApplianceRootCmd(f))
	cmd.AddCommand(ecloudCreditRootCmd(f))

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
