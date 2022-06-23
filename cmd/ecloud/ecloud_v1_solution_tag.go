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

func ecloudSolutionTagRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tag",
		Short: "sub-commands relating to solution tags",
	}

	// Child commands
	cmd.AddCommand(ecloudSolutionTagListCmd(f))
	cmd.AddCommand(ecloudSolutionTagShowCmd(f))
	cmd.AddCommand(ecloudSolutionTagCreateCmd(f))
	cmd.AddCommand(ecloudSolutionTagUpdateCmd(f))
	cmd.AddCommand(ecloudSolutionTagDeleteCmd(f))

	return cmd
}

func ecloudSolutionTagListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list <solution: id>",
		Short:   "lists solution tags",
		Long:    "This command lists solution tags",
		Example: "ukfast ecloud solution tag list 123",
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

			return ecloudSolutionTagList(c.ECloudService(), cmd, args)
		},
	}
}

func ecloudSolutionTagList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	solutionID, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("Invalid solution ID [%s]", args[0])
	}

	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	tags, err := service.GetSolutionTags(solutionID, params)
	if err != nil {
		return fmt.Errorf("Error retrieving solution tags: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudTagsProvider(tags))
}

func ecloudSolutionTagShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <solution: id> <tag: key>...",
		Short:   "Shows a solution tag",
		Long:    "This command shows one or more solution tags",
		Example: "ukfast ecloud solution tag show 123 foo",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing solution")
			}
			if len(args) < 2 {
				return errors.New("Missing tag")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ecloudSolutionTagShow(c.ECloudService(), cmd, args)
		},
	}
}

func ecloudSolutionTagShow(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	solutionID, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("Invalid solution ID [%s]", args[0])
	}

	var tags []ecloud.Tag

	for _, arg := range args[1:] {
		tag, err := service.GetSolutionTag(solutionID, arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving solution tag [%s]: %s", arg, err)
			continue
		}

		tags = append(tags, tag)
	}

	return output.CommandOutput(cmd, OutputECloudTagsProvider(tags))
}

func ecloudSolutionTagCreateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create <solution: id>",
		Short:   "Creates a solution tag",
		Long:    "This command creates a solution tag",
		Example: "ukfast ecloud solution tag create 123 --key foo --value bar",
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

			return ecloudSolutionTagCreate(c.ECloudService(), cmd, args)
		},
	}

	cmd.Flags().String("key", "", "Key for tag")
	cmd.MarkFlagRequired("key")
	cmd.Flags().String("value", "", "Value for tag")
	cmd.MarkFlagRequired("value")

	return cmd
}

func ecloudSolutionTagCreate(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	solutionID, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("Invalid solution ID [%s]", args[0])
	}

	key, _ := cmd.Flags().GetString("key")
	value, _ := cmd.Flags().GetString("value")

	createRequest := ecloud.CreateTagRequest{
		Key:   key,
		Value: value,
	}

	err = service.CreateSolutionTag(solutionID, createRequest)
	if err != nil {
		return fmt.Errorf("Error creating solution tag: %s", err)
	}

	tag, err := service.GetSolutionTag(solutionID, key)
	if err != nil {
		return fmt.Errorf("Error retrieving new solution tag: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudTagsProvider([]ecloud.Tag{tag}))
}

func ecloudSolutionTagUpdateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <solution: id> <tag: key>...",
		Short:   "Updates a solution tag",
		Long:    "This command updates one or more solution tags",
		Example: "ukfast ecloud solution tag update 123 foo --value \"new value\"",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing solution")
			}
			if len(args) < 2 {
				return errors.New("Missing tag")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ecloudSolutionTagUpdate(c.ECloudService(), cmd, args)
		},
	}

	cmd.Flags().String("value", "", "Value for tag")

	return cmd
}

func ecloudSolutionTagUpdate(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	solutionID, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("Invalid solution ID [%s]", args[0])
	}

	patchRequest := ecloud.PatchTagRequest{}

	if cmd.Flags().Changed("value") {
		value, _ := cmd.Flags().GetString("value")
		patchRequest.Value = value
	}

	var tags []ecloud.Tag

	for _, arg := range args[1:] {
		err = service.PatchSolutionTag(solutionID, arg, patchRequest)
		if err != nil {
			output.OutputWithErrorLevelf("Error updating solution tag [%s]: %s", arg, err)
			continue
		}

		tag, err := service.GetSolutionTag(solutionID, arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving updated solution tag [%s]: %s", arg, err)
			continue
		}

		tags = append(tags, tag)
	}

	return output.CommandOutput(cmd, OutputECloudTagsProvider(tags))
}

func ecloudSolutionTagDeleteCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "delete <solution: id> <tag: key>...",
		Short:   "Removes a solution tag ",
		Long:    "This command removes one or more solution tags",
		Example: "ukfast ecloud solution tag delete 123 foo",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing solution")
			}
			if len(args) < 2 {
				return errors.New("Missing tag")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ecloudSolutionTagDelete(c.ECloudService(), cmd, args)
		},
	}
}

func ecloudSolutionTagDelete(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	solutionID, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("Invalid solution ID [%s]", args[0])
	}

	for _, arg := range args[1:] {
		err = service.DeleteSolutionTag(solutionID, arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error removing solution tag [%s]: %s", arg, err)
		}
	}

	return nil
}
