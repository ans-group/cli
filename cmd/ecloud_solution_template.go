package cmd

import (
	"errors"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func ecloudSolutionTemplateRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "template",
		Short: "sub-commands relating to solution templates",
	}

	// Child commands
	cmd.AddCommand(ecloudSolutionTemplateListCmd())
	cmd.AddCommand(ecloudSolutionTemplateShowCmd())
	cmd.AddCommand(ecloudSolutionTemplateUpdateCmd())
	cmd.AddCommand(ecloudSolutionTemplateDeleteCmd())

	return cmd
}

func ecloudSolutionTemplateListCmd() *cobra.Command {
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
		Run: func(cmd *cobra.Command, args []string) {
			ecloudSolutionTemplateList(getClient().ECloudService(), cmd, args)
		},
	}
}

func ecloudSolutionTemplateList(service ecloud.ECloudService, cmd *cobra.Command, args []string) {
	solutionID, err := strconv.Atoi(args[0])
	if err != nil {
		output.Fatalf("Invalid solution ID [%s]", args[0])
		return
	}

	params, err := GetAPIRequestParametersFromFlags()
	if err != nil {
		output.Fatal(err.Error())
		return
	}

	templates, err := service.GetSolutionTemplates(solutionID, params)
	if err != nil {
		output.Fatalf("Error retrieving solution templates: %s", err)
		return
	}

	outputECloudTemplates(templates)
}

func ecloudSolutionTemplateShowCmd() *cobra.Command {
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
		Run: func(cmd *cobra.Command, args []string) {
			ecloudSolutionTemplateShow(getClient().ECloudService(), cmd, args)
		},
	}
}

func ecloudSolutionTemplateShow(service ecloud.ECloudService, cmd *cobra.Command, args []string) {
	solutionID, err := strconv.Atoi(args[0])
	if err != nil {
		output.Fatalf("Invalid solution ID [%s]", args[0])
		return
	}

	var templates []ecloud.Template

	for _, arg := range args[1:] {
		template, err := service.GetSolutionTemplate(solutionID, arg)
		if err != nil {
			OutputWithErrorLevelf("Error retrieving solution template [%s]: %s", arg, err)
			continue
		}

		templates = append(templates, template)
	}

	outputECloudTemplates(templates)
}

func ecloudSolutionTemplateUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <solution: id> <template: template>...",
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
		Run: func(cmd *cobra.Command, args []string) {
			ecloudSolutionTemplateUpdate(getClient().ECloudService(), cmd, args)
		},
	}

	cmd.Flags().String("name", "", "Name for template")

	return cmd
}

func ecloudSolutionTemplateUpdate(service ecloud.ECloudService, cmd *cobra.Command, args []string) {
	solutionID, err := strconv.Atoi(args[0])
	if err != nil {
		output.Fatalf("Invalid solution ID [%s]", args[0])
		return
	}

	templateName := args[1]

	if cmd.Flags().Changed("name") {
		name, _ := cmd.Flags().GetString("name")
		patchRequest := ecloud.RenameTemplateRequest{
			Destination: name,
		}

		err = service.RenameSolutionTemplate(solutionID, templateName, patchRequest)
		if err != nil {
			output.Fatalf("Error updating solution template: %s", err)
			return
		}

		err := WaitForCommand(SolutionTemplateExistsWaitFunc(service, solutionID, name, true))
		if err != nil {
			output.Fatalf("Error waiting for solution template update: %s", err)
			return
		}

		templateName = name
	}

	template, err := service.GetSolutionTemplate(solutionID, templateName)
	if err != nil {
		output.Fatalf("Error retrieving updated solution template: %s", err)
		return
	}

	outputECloudTemplates([]ecloud.Template{template})
}

func ecloudSolutionTemplateDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "delete <solution: id> <template: template>...",
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
		Run: func(cmd *cobra.Command, args []string) {
			ecloudSolutionTemplateDelete(getClient().ECloudService(), cmd, args)
		},
	}
}

func ecloudSolutionTemplateDelete(service ecloud.ECloudService, cmd *cobra.Command, args []string) {
	solutionID, err := strconv.Atoi(args[0])
	if err != nil {
		output.Fatalf("Invalid solution ID [%s]", args[0])
		return
	}

	for _, arg := range args[1:] {
		err = service.DeleteSolutionTemplate(solutionID, arg)
		if err != nil {
			OutputWithErrorLevelf("Error removing solution template [%s]: %s", arg, err)
		}
	}
}
