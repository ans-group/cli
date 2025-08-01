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

func ecloudVirtualMachineTagRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tag",
		Short: "sub-commands relating to virtual machine tags",
	}

	// Child commands
	cmd.AddCommand(ecloudVirtualMachineTagListCmd(f))
	cmd.AddCommand(ecloudVirtualMachineTagShowCmd(f))
	cmd.AddCommand(ecloudVirtualMachineTagCreateCmd(f))
	cmd.AddCommand(ecloudVirtualMachineTagUpdateCmd(f))
	cmd.AddCommand(ecloudVirtualMachineTagDeleteCmd(f))

	return cmd
}

func ecloudVirtualMachineTagListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list <virtualmachine: id>",
		Short:   "lists virtual machine tags",
		Long:    "This command lists virtual machine tags",
		Example: "ans ecloud vm tag list 123",
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

			return ecloudVirtualMachineTagList(c.ECloudService(), cmd, args)
		},
	}
}

func ecloudVirtualMachineTagList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	vmID, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("Invalid virtual machine ID [%s]", args[0])
	}

	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	tags, err := service.GetVirtualMachineTags(vmID, params)
	if err != nil {
		return fmt.Errorf("Error retrieving virtual machine tags: %s", err)
	}

	return output.CommandOutput(cmd, TagCollection(tags))
}

func ecloudVirtualMachineTagShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <virtualmachine: id> <tag: key>...",
		Short:   "Shows a virtual machine tag",
		Long:    "This command shows one or more virtual machine tags",
		Example: "ans ecloud vm tag show 123 foo",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing virtual machine")
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

			return ecloudVirtualMachineTagShow(c.ECloudService(), cmd, args)
		},
	}
}

func ecloudVirtualMachineTagShow(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	vmID, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("Invalid virtual machine ID [%s]", args[0])
	}

	var tags []ecloud.TagV1

	for _, arg := range args[1:] {
		tag, err := service.GetVirtualMachineTag(vmID, arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving virtual machine tag [%s]: %s", arg, err)
			continue
		}

		tags = append(tags, tag)
	}

	return output.CommandOutput(cmd, TagCollection(tags))
}

func ecloudVirtualMachineTagCreateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create <virtualmachine: id>",
		Short:   "Creates a virtual machine tag",
		Long:    "This command creates a virtual machine tag",
		Example: "ans ecloud vm tag create 123 --key foo --value bar",
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

			return ecloudVirtualMachineTagCreate(c.ECloudService(), cmd, args)
		},
	}

	cmd.Flags().String("key", "", "Key for tag")
	cmd.MarkFlagRequired("key")
	cmd.Flags().String("value", "", "Value for tag")
	cmd.MarkFlagRequired("value")

	return cmd
}

func ecloudVirtualMachineTagCreate(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	vmID, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("Invalid virtual machine ID [%s]", args[0])
	}

	key, _ := cmd.Flags().GetString("key")
	value, _ := cmd.Flags().GetString("value")

	createRequest := ecloud.CreateTagV1Request{
		Key:   key,
		Value: value,
	}

	err = service.CreateVirtualMachineTag(vmID, createRequest)
	if err != nil {
		return fmt.Errorf("Error creating virtual machine tag: %s", err)
	}

	tag, err := service.GetVirtualMachineTag(vmID, key)
	if err != nil {
		return fmt.Errorf("Error retrieving new virtual machine tag: %s", err)
	}

	return output.CommandOutput(cmd, TagCollection([]ecloud.TagV1{tag}))
}

func ecloudVirtualMachineTagUpdateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <virtualmachine: id> <tag: key>...",
		Short:   "Updates a virtual machine tag",
		Long:    "This command updates one or more virtual machine tags",
		Example: "ans ecloud vm tag update 123 foo --value \"new value\"",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing virtual machine")
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

			return ecloudVirtualMachineTagUpdate(c.ECloudService(), cmd, args)
		},
	}

	cmd.Flags().String("value", "", "Value for tag")

	return cmd
}

func ecloudVirtualMachineTagUpdate(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	vmID, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("Invalid virtual machine ID [%s]", args[0])
	}

	patchRequest := ecloud.PatchTagV1Request{}

	if cmd.Flags().Changed("value") {
		recordName, _ := cmd.Flags().GetString("value")
		patchRequest.Value = recordName
	}

	var tags []ecloud.TagV1

	for _, arg := range args[1:] {
		err = service.PatchVirtualMachineTag(vmID, arg, patchRequest)
		if err != nil {
			output.OutputWithErrorLevelf("Error updating virtual machine tag [%s]: %s", arg, err)
			continue
		}

		tag, err := service.GetVirtualMachineTag(vmID, arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving updated virtual machine tag [%s]: %s", arg, err)
			continue
		}

		tags = append(tags, tag)
	}

	return output.CommandOutput(cmd, TagCollection(tags))
}

func ecloudVirtualMachineTagDeleteCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "delete <virtualmachine: id> <tag: key>...",
		Short:   "Removes a virtual machine tag ",
		Long:    "This command removes one or more virtual machine tags",
		Example: "ans ecloud vm tag delete 123 foo",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing virtual machine")
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

			return ecloudVirtualMachineTagDelete(c.ECloudService(), cmd, args)
		},
	}
}

func ecloudVirtualMachineTagDelete(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	vmID, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("Invalid virtual machine ID [%s]", args[0])
	}

	for _, arg := range args[1:] {
		err = service.DeleteVirtualMachineTag(vmID, arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error removing virtual machine tag [%s]: %s", arg, err)
		}
	}

	return nil
}
