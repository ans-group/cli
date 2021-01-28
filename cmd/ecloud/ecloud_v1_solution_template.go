package ecloud

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
	"github.com/ukfast/cli/internal/pkg/helper"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func ecloudSolutionTemplateRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "template",
		Short: "sub-commands relating to solution templates",
	}

	// Child commands
	cmd.AddCommand(ecloudSolutionTemplateListCmd(f))
	cmd.AddCommand(ecloudSolutionTemplateShowCmd(f))
	cmd.AddCommand(ecloudSolutionTemplateUpdateCmd(f))
	cmd.AddCommand(ecloudSolutionTemplateDeleteCmd(f))

	return cmd
}

func ecloudSolutionTemplateListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists solution templates",
		Long:    "This command lists solution templates",
		Example: "ukfast ecloud solution template list 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing solution")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ecloudSolutionTemplateList(c.ECloudService(), cmd, args)
		},
	}
}

func ecloudSolutionTemplateList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	solutionID, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("Invalid solution ID [%s]", args[0])
	}

	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	templates, err := service.GetSolutionTemplates(solutionID, params)
	if err != nil {
		return fmt.Errorf("Error retrieving solution templates: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudTemplatesProvider(templates))
}

func ecloudSolutionTemplateShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <solution: id> <template: name>...",
		Short:   "Shows a solution template",
		Long:    "This command shows one or more solution templates",
		Example: "ukfast ecloud solution template show 123 foo",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing solution")
			}
			if len(args) < 2 {
				return errors.New("Missing template")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ecloudSolutionTemplateShow(c.ECloudService(), cmd, args)
		},
	}
}

func ecloudSolutionTemplateShow(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	solutionID, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("Invalid solution ID [%s]", args[0])
	}

	var templates []ecloud.Template

	for _, arg := range args[1:] {
		template, err := service.GetSolutionTemplate(solutionID, arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving solution template [%s]: %s", arg, err)
			continue
		}

		templates = append(templates, template)
	}

	return output.CommandOutput(cmd, OutputECloudTemplatesProvider(templates))
}

func ecloudSolutionTemplateUpdateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <solution: id> <template: name>...",
		Short:   "Updates a solution template",
		Long:    "This command updates a solution template",
		Example: "ukfast ecloud solution template update 123 foo --name \"bar\"",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing solution")
			}
			if len(args) < 2 {
				return errors.New("Missing template")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ecloudSolutionTemplateUpdate(c.ECloudService(), cmd, args)
		},
	}

	cmd.Flags().String("name", "", "Name for template")

	return cmd
}

func ecloudSolutionTemplateUpdate(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	solutionID, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("Invalid solution ID [%s]", args[0])
	}

	templateName := args[1]

	if cmd.Flags().Changed("name") {
		name, _ := cmd.Flags().GetString("name")
		patchRequest := ecloud.RenameTemplateRequest{
			Destination: name,
		}

		err = service.RenameSolutionTemplate(solutionID, templateName, patchRequest)
		if err != nil {
			return fmt.Errorf("Error updating solution template: %s", err)
		}

		err := helper.WaitForCommand(SolutionTemplateExistsWaitFunc(service, solutionID, name, true))
		if err != nil {
			return fmt.Errorf("Error waiting for solution template update: %s", err)
		}

		templateName = name
	}

	template, err := service.GetSolutionTemplate(solutionID, templateName)
	if err != nil {
		return fmt.Errorf("Error retrieving updated solution template: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudTemplatesProvider([]ecloud.Template{template}))
}

func ecloudSolutionTemplateDeleteCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete <solution: id> <template: name>...",
		Short:   "Removes a solution template ",
		Long:    "This command removes one or more solution templates",
		Example: "ukfast ecloud solution template delete 123 foo",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing solution")
			}
			if len(args) < 2 {
				return errors.New("Missing template")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ecloudSolutionTemplateDelete(c.ECloudService(), cmd, args)
		},
	}

	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the template has been completely created before continuing on")

	return cmd
}

func ecloudSolutionTemplateDelete(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	solutionID, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("Invalid solution ID [%s]", args[0])
	}

	for _, arg := range args[1:] {
		err = service.DeleteSolutionTemplate(solutionID, arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error removing solution template [%s]: %s", arg, err)
			continue
		}

		waitFlag, _ := cmd.Flags().GetBool("wait")
		if waitFlag {
			err := helper.WaitForCommand(SolutionTemplateExistsWaitFunc(service, solutionID, arg, false))
			if err != nil {
				output.OutputWithErrorLevelf("Error removing solution template [%s]: %s", arg, err)
			}
		}
	}

	return nil
}
