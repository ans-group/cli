package ecloudflex

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/helper"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/ecloudflex"
	"github.com/spf13/cobra"
)

func ecloudflexProjectRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "project",
		Short: "sub-commands relating to projects",
	}

	// Child commands
	cmd.AddCommand(ecloudflexProjectListCmd(f))
	cmd.AddCommand(ecloudflexProjectShowCmd(f))

	return cmd
}

func ecloudflexProjectListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists projects",
		Long:    "This command lists projects",
		Example: "ans ecloudflex project list",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ecloudflexProjectList(c.ECloudFlexService(), cmd, args)
		},
	}
}

func ecloudflexProjectList(service ecloudflex.ECloudFlexService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	projects, err := service.GetProjects(params)
	if err != nil {
		return fmt.Errorf("Error retrieving projects: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudFlexProjectsProvider(projects))
}

func ecloudflexProjectShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <project: id>...",
		Short:   "Shows a project",
		Long:    "This command shows one or more projects",
		Example: "ans ecloudflex project show 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing project")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ecloudflexProjectShow(c.ECloudFlexService(), cmd, args)
		},
	}
}

func ecloudflexProjectShow(service ecloudflex.ECloudFlexService, cmd *cobra.Command, args []string) error {
	var projects []ecloudflex.Project
	for _, arg := range args {
		projectID, err := strconv.Atoi(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Invalid project ID [%s]", arg)
			continue
		}

		project, err := service.GetProject(projectID)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving project [%s]: %s", arg, err)
			continue
		}

		projects = append(projects, project)
	}

	return output.CommandOutput(cmd, OutputECloudFlexProjectsProvider(projects))
}
