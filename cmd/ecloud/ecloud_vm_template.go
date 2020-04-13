package ecloud

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/clierrors"
	"github.com/ukfast/cli/internal/pkg/factory"
	"github.com/ukfast/cli/internal/pkg/helper"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func ecloudVirtualMachineTemplateRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "template",
		Short: "sub-commands relating to virtual machine templates",
	}

	// Child commands
	cmd.AddCommand(ecloudVirtualMachineTemplateCreateCmd(f))

	return cmd
}

func ecloudVirtualMachineTemplateCreateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create <virtualmachine: id>",
		Short:   "Creates a virtual machine template",
		Long:    "This command creates a virtual machine template",
		Example: "ukfast ecloud vm template create 123 --name \"foo\" --type \"solution\"",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing virtual machine")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ecloudVirtualMachineTemplateCreate(c.ECloudService(), cmd, args)
		},
	}

	cmd.Flags().String("name", "", "Name for new template")
	cmd.MarkFlagRequired("name")
	cmd.Flags().String("type", "", "Type of template (pod/solution)")
	cmd.MarkFlagRequired("type")
	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the template has been completely created before continuing on")

	return cmd
}

func ecloudVirtualMachineTemplateCreate(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	vmID, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("Invalid virtual machine ID [%s]", args[0])
	}

	templateType, _ := cmd.Flags().GetString("type")
	parsedTemplateType, err := ecloud.ParseTemplateType(templateType)
	if err != nil {
		return clierrors.NewErrInvalidFlagValue("type", templateType, err)
	}

	templateName, _ := cmd.Flags().GetString("name")
	createRequest := ecloud.CreateVirtualMachineTemplateRequest{
		TemplateName: templateName,
		TemplateType: parsedTemplateType,
	}

	err = service.CreateVirtualMachineTemplate(vmID, createRequest)
	if err != nil {
		return fmt.Errorf("Error creating virtual machine template: %s", err)
	}

	waitFlag, _ := cmd.Flags().GetBool("wait")
	if waitFlag {
		err := helper.WaitForCommand(VirtualMachineStatusWaitFunc(service, vmID, ecloud.VirtualMachineStatusComplete))
		if err != nil {
			return err
		}
	}

	template, err := getVMTemplate(service, vmID, templateName, parsedTemplateType)
	if err != nil {
		if _, ok := err.(*ecloud.TemplateNotFoundError); ok {
			return fmt.Errorf("Error creating virtual machine template (unknown failure)")
		}

		return fmt.Errorf("Error retrieving new virtual machine template: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudTemplatesProvider([]ecloud.Template{template}))
}

func getVMTemplate(service ecloud.ECloudService, vmID int, templateName string, templateType ecloud.TemplateType) (ecloud.Template, error) {
	switch templateType {
	case ecloud.TemplateTypePod:
		return getPodVMTemplate(service, vmID, templateName)
	case ecloud.TemplateTypeSolution:
		return getSolutionVMTemplate(service, vmID, templateName)
	}

	return ecloud.Template{}, errors.New("unknown template type")
}

func getPodVMTemplate(service ecloud.ECloudService, vmID int, templateName string) (ecloud.Template, error) {
	vm, err := service.GetVirtualMachine(vmID)
	if err != nil {
		return ecloud.Template{}, fmt.Errorf("Error retrieving virtual machine: %s", err)
	}

	solution, err := service.GetSolution(vm.SolutionID)
	if err != nil {
		return ecloud.Template{}, fmt.Errorf("Error retrieving solution: %s", err)
	}

	return service.GetPodTemplate(solution.PodID, templateName)
}

func getSolutionVMTemplate(service ecloud.ECloudService, vmID int, templateName string) (ecloud.Template, error) {
	vm, err := service.GetVirtualMachine(vmID)
	if err != nil {
		return ecloud.Template{}, fmt.Errorf("Error retrieving virtual machine: %s", err)
	}

	return service.GetSolutionTemplate(vm.SolutionID, templateName)
}
