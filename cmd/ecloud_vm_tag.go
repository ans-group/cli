package cmd

import (
	"errors"
	"strconv"

	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/ecloud"

	"github.com/spf13/cobra"
)

func ecloudVirtualMachineTagRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tag",
		Short: "sub-commands relating to virtual machine tags",
	}

	// Child commands
	cmd.AddCommand(ecloudVirtualMachineTagListCmd())
	cmd.AddCommand(ecloudVirtualMachineTagShowCmd())
	cmd.AddCommand(ecloudVirtualMachineTagCreateCmd())
	cmd.AddCommand(ecloudVirtualMachineTagUpdateCmd())
	cmd.AddCommand(ecloudVirtualMachineTagDeleteCmd())

	return cmd
}

func ecloudVirtualMachineTagListCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "list <virtualmachine: id>",
		Short:   "lists virtual machine tags",
		Long:    "This command lists virtual machine tags",
		Example: "ukfast ecloud vm tag list 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing virtual machine")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			ecloudVirtualMachineTagList(getClient().ECloudService(), cmd, args)
		},
	}
}

func ecloudVirtualMachineTagList(service ecloud.ECloudService, cmd *cobra.Command, args []string) {
	vmID, err := strconv.Atoi(args[0])
	if err != nil {
		output.Fatalf("Invalid virtual machine ID [%s]", args[0])
		return
	}

	params, err := GetAPIRequestParametersFromFlags()
	if err != nil {
		output.Fatal(err.Error())
		return
	}

	tags, err := service.GetVirtualMachineTags(vmID, params)
	if err != nil {
		output.Fatalf("Error retrieving virtual machine tags: %s", err)
		return
	}

	outputECloudTags(tags)
}

func ecloudVirtualMachineTagShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "show <virtualmachine: id> <tag: key>...",
		Short:   "Shows a virtual machine tag",
		Long:    "This command shows one or more virtual machine tags",
		Example: "ukfast ecloud vm tag show 123 foo",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing virtual machine")
			}
			if len(args) < 2 {
				return errors.New("Missing tag")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			ecloudVirtualMachineTagShow(getClient().ECloudService(), cmd, args)
		},
	}
}

func ecloudVirtualMachineTagShow(service ecloud.ECloudService, cmd *cobra.Command, args []string) {
	vmID, err := strconv.Atoi(args[0])
	if err != nil {
		output.Fatalf("Invalid virtual machine ID [%s]", args[0])
		return
	}

	var tags []ecloud.Tag

	for _, arg := range args[1:] {
		tag, err := service.GetVirtualMachineTag(vmID, arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving virtual machine tag [%s]: %s", arg, err)
			continue
		}

		tags = append(tags, tag)
	}

	outputECloudTags(tags)
}

func ecloudVirtualMachineTagCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create <virtualmachine: id>",
		Short:   "Creates a virtual machine tag",
		Long:    "This command creates a virtual machine tag",
		Example: "ukfast ecloud vm tag create 123 --key foo --value bar",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing virtual machine")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			ecloudVirtualMachineTagCreate(getClient().ECloudService(), cmd, args)
		},
	}

	cmd.Flags().String("key", "", "Key for tag")
	cmd.MarkFlagRequired("key")
	cmd.Flags().String("value", "", "Value for tag")
	cmd.MarkFlagRequired("value")

	return cmd
}

func ecloudVirtualMachineTagCreate(service ecloud.ECloudService, cmd *cobra.Command, args []string) {
	vmID, err := strconv.Atoi(args[0])
	if err != nil {
		output.Fatalf("Invalid virtual machine ID [%s]", args[0])
		return
	}

	key, _ := cmd.Flags().GetString("key")
	value, _ := cmd.Flags().GetString("value")

	createRequest := ecloud.CreateTagRequest{
		Key:   key,
		Value: value,
	}

	err = service.CreateVirtualMachineTag(vmID, createRequest)
	if err != nil {
		output.Fatalf("Error creating virtual machine tag: %s", err)
		return
	}

	tag, err := service.GetVirtualMachineTag(vmID, key)
	if err != nil {
		output.Fatalf("Error retrieving new virtual machine tag: %s", err)
		return
	}

	outputECloudTags([]ecloud.Tag{tag})
}

func ecloudVirtualMachineTagUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <virtualmachine: id> <tag: key>...",
		Short:   "Updates a virtual machine tag",
		Long:    "This command updates one or more virtual machine tags",
		Example: "ukfast ecloud vm tag update 123 foo --value \"new value\"",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing virtual machine")
			}
			if len(args) < 2 {
				return errors.New("Missing tag")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			ecloudVirtualMachineTagUpdate(getClient().ECloudService(), cmd, args)
		},
	}

	cmd.Flags().String("value", "", "Value for tag")

	return cmd
}

func ecloudVirtualMachineTagUpdate(service ecloud.ECloudService, cmd *cobra.Command, args []string) {
	vmID, err := strconv.Atoi(args[0])
	if err != nil {
		output.Fatalf("Invalid virtual machine ID [%s]", args[0])
		return
	}

	patchRequest := ecloud.PatchTagRequest{}

	if cmd.Flags().Changed("value") {
		recordName, _ := cmd.Flags().GetString("value")
		patchRequest.Value = recordName
	}

	var tags []ecloud.Tag

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

	outputECloudTags(tags)
}

func ecloudVirtualMachineTagDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "delete <virtualmachine: id> <tag: key>...",
		Short:   "Removes a virtual machine tag ",
		Long:    "This command removes one or more virtual machine tags",
		Example: "ukfast ecloud vm tag delete 123 foo",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing virtual machine")
			}
			if len(args) < 2 {
				return errors.New("Missing tag")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			ecloudVirtualMachineTagDelete(getClient().ECloudService(), cmd, args)
		},
	}
}

func ecloudVirtualMachineTagDelete(service ecloud.ECloudService, cmd *cobra.Command, args []string) {
	vmID, err := strconv.Atoi(args[0])
	if err != nil {
		output.Fatalf("Invalid virtual machine ID [%s]", args[0])
		return
	}

	for _, arg := range args[1:] {
		err = service.DeleteVirtualMachineTag(vmID, arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error removing virtual machine tag [%s]: %s", arg, err)
		}
	}
}
