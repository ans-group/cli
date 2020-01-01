package cmd

import (
	"errors"
	"strconv"

	"github.com/ukfast/cli/internal/pkg/helper"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/ecloud"

	"github.com/spf13/cobra"
)

func ecloudSolutionTagRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tag",
		Short: "sub-commands relating to solution tags",
	}

	// Child commands
	cmd.AddCommand(ecloudSolutionTagListCmd())
	cmd.AddCommand(ecloudSolutionTagShowCmd())
	cmd.AddCommand(ecloudSolutionTagCreateCmd())
	cmd.AddCommand(ecloudSolutionTagUpdateCmd())
	cmd.AddCommand(ecloudSolutionTagDeleteCmd())

	return cmd
}

func ecloudSolutionTagListCmd() *cobra.Command {
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
		Run: func(cmd *cobra.Command, args []string) {
			ecloudSolutionTagList(getClient().ECloudService(), cmd, args)
		},
	}
}

func ecloudSolutionTagList(service ecloud.ECloudService, cmd *cobra.Command, args []string) {
	solutionID, err := strconv.Atoi(args[0])
	if err != nil {
		output.Fatalf("Invalid solution ID [%s]", args[0])
		return
	}

	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		output.Fatal(err.Error())
		return
	}

	tags, err := service.GetSolutionTags(solutionID, params)
	if err != nil {
		output.Fatalf("Error retrieving solution tags: %s", err)
		return
	}

	outputECloudTags(tags)
}

func ecloudSolutionTagShowCmd() *cobra.Command {
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
		Run: func(cmd *cobra.Command, args []string) {
			ecloudSolutionTagShow(getClient().ECloudService(), cmd, args)
		},
	}
}

func ecloudSolutionTagShow(service ecloud.ECloudService, cmd *cobra.Command, args []string) {
	solutionID, err := strconv.Atoi(args[0])
	if err != nil {
		output.Fatalf("Invalid solution ID [%s]", args[0])
		return
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

	outputECloudTags(tags)
}

func ecloudSolutionTagCreateCmd() *cobra.Command {
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
		Run: func(cmd *cobra.Command, args []string) {
			ecloudSolutionTagCreate(getClient().ECloudService(), cmd, args)
		},
	}

	cmd.Flags().String("key", "", "Key for tag")
	cmd.MarkFlagRequired("key")
	cmd.Flags().String("value", "", "Value for tag")
	cmd.MarkFlagRequired("value")

	return cmd
}

func ecloudSolutionTagCreate(service ecloud.ECloudService, cmd *cobra.Command, args []string) {
	solutionID, err := strconv.Atoi(args[0])
	if err != nil {
		output.Fatalf("Invalid solution ID [%s]", args[0])
		return
	}

	key, _ := cmd.Flags().GetString("key")
	value, _ := cmd.Flags().GetString("value")

	createRequest := ecloud.CreateTagRequest{
		Key:   key,
		Value: value,
	}

	err = service.CreateSolutionTag(solutionID, createRequest)
	if err != nil {
		output.Fatalf("Error creating solution tag: %s", err)
		return
	}

	tag, err := service.GetSolutionTag(solutionID, key)
	if err != nil {
		output.Fatalf("Error retrieving new solution tag: %s", err)
		return
	}

	outputECloudTags([]ecloud.Tag{tag})
}

func ecloudSolutionTagUpdateCmd() *cobra.Command {
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
		Run: func(cmd *cobra.Command, args []string) {
			ecloudSolutionTagUpdate(getClient().ECloudService(), cmd, args)
		},
	}

	cmd.Flags().String("value", "", "Value for tag")

	return cmd
}

func ecloudSolutionTagUpdate(service ecloud.ECloudService, cmd *cobra.Command, args []string) {
	solutionID, err := strconv.Atoi(args[0])
	if err != nil {
		output.Fatalf("Invalid solution ID [%s]", args[0])
		return
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

	outputECloudTags(tags)
}

func ecloudSolutionTagDeleteCmd() *cobra.Command {
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
		Run: func(cmd *cobra.Command, args []string) {
			ecloudSolutionTagDelete(getClient().ECloudService(), cmd, args)
		},
	}
}

func ecloudSolutionTagDelete(service ecloud.ECloudService, cmd *cobra.Command, args []string) {
	solutionID, err := strconv.Atoi(args[0])
	if err != nil {
		output.Fatalf("Invalid solution ID [%s]", args[0])
		return
	}

	for _, arg := range args[1:] {
		err = service.DeleteSolutionTag(solutionID, arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error removing solution tag [%s]: %s", arg, err)
		}
	}
}
