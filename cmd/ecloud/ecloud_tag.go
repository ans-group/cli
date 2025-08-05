package ecloud

import (
	"errors"
	"fmt"

	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/helper"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/spf13/cobra"
)

func ecloudTagRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tag",
		Short: "sub-commands relating to tags",
	}

	// Child commands
	cmd.AddCommand(ecloudTagListCmd(f))
	cmd.AddCommand(ecloudTagShowCmd(f))
	cmd.AddCommand(ecloudTagCreateCmd(f))
	cmd.AddCommand(ecloudTagUpdateCmd(f))
	cmd.AddCommand(ecloudTagDeleteCmd(f))

	return cmd
}

func ecloudTagListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Lists tags",
		Long:    "This command lists tags",
		Example: "ans ecloud tag list",
		RunE:    ecloudCobraRunEFunc(f, ecloudTagList),
	}

	cmd.Flags().String("name", "", "Tag name for filtering")
	cmd.Flags().String("scope", "", "Tag scope for filtering")

	return cmd
}

func ecloudTagList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd,
		helper.NewStringFilterFlagOption("name", "name"),
		helper.NewStringFilterFlagOption("scope", "scope"),
	)
	if err != nil {
		return err
	}

	tags, err := service.GetTags(params)
	if err != nil {
		return fmt.Errorf("ecloud: Error retrieving tags: %s", err)
	}

	return output.CommandOutput(cmd, TagCollection(tags))
}

func ecloudTagShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <tag: id>...",
		Short:   "Shows a tag",
		Long:    "This command shows one or more tags",
		Example: "ans ecloud tag show tag-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("missing tag")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudTagShow),
	}
}

func ecloudTagShow(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	var tags []ecloud.Tag
	for _, arg := range args {
		tag, err := service.GetTag(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving tag [%s]: %s", arg, err)
			continue
		}

		tags = append(tags, tag)
	}

	return output.CommandOutput(cmd, TagCollection(tags))
}

func ecloudTagCreateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Creates a tag",
		Long:    "This command creates a tag",
		Example: "ans ecloud tag create --name \"production\" --scope \"environment\"",
		RunE:    ecloudCobraRunEFunc(f, ecloudTagCreate),
	}

	// Setup flags
	cmd.Flags().String("name", "", "Name of tag")
	_ = cmd.MarkFlagRequired("name")
	cmd.Flags().String("scope", "", "Scope of tag")
	_ = cmd.MarkFlagRequired("scope")

	return cmd
}

func ecloudTagCreate(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	createRequest := ecloud.CreateTagRequest{}
	createRequest.Name, _ = cmd.Flags().GetString("name")
	createRequest.Scope, _ = cmd.Flags().GetString("scope")

	tagID, err := service.CreateTag(createRequest)
	if err != nil {
		return fmt.Errorf("ecloud: Error creating tag: %s", err)
	}

	tag, err := service.GetTag(tagID)
	if err != nil {
		return fmt.Errorf("ecloud: Error retrieving new tag: %s", err)
	}

	return output.CommandOutput(cmd, TagCollection([]ecloud.Tag{tag}))
}

func ecloudTagUpdateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <tag: id>...",
		Short:   "Updates a tag",
		Long:    "This command updates one or more tags",
		Example: "ans ecloud tag update tag-abcdef12 --name \"staging\"",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("missing tag")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudTagUpdate),
	}

	cmd.Flags().String("name", "", "Name of tag")
	cmd.Flags().String("scope", "", "Scope of tag")

	return cmd
}

func ecloudTagUpdate(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	patchRequest := ecloud.PatchTagRequest{}

	if cmd.Flags().Changed("name") {
		patchRequest.Name, _ = cmd.Flags().GetString("name")
	}

	if cmd.Flags().Changed("scope") {
		patchRequest.Scope, _ = cmd.Flags().GetString("scope")
	}

	var tags []ecloud.Tag
	for _, arg := range args {
		err := service.PatchTag(arg, patchRequest)
		if err != nil {
			output.OutputWithErrorLevelf("Error updating tag [%s]: %s", arg, err)
			continue
		}

		tag, err := service.GetTag(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving updated tag [%s]: %s", arg, err)
			continue
		}

		tags = append(tags, tag)
	}

	return output.CommandOutput(cmd, TagCollection(tags))
}

func ecloudTagDeleteCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete <tag: id>...",
		Short:   "Removes a tag",
		Long:    "This command removes one or more tags",
		Example: "ans ecloud tag delete tag-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("missing tag")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudTagDelete),
	}

	return cmd
}

func ecloudTagDelete(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	for _, arg := range args {
		err := service.DeleteTag(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error removing tag [%s]: %s", arg, err)
			continue
		}
	}
	return nil
}
