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

func ecloudNetworkRulePortRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "networkruleport",
		Short: "sub-commands relating to network rule ports",
	}

	// Child commands
	cmd.AddCommand(ecloudNetworkRulePortListCmd(f))
	cmd.AddCommand(ecloudNetworkRulePortShowCmd(f))
	cmd.AddCommand(ecloudNetworkRulePortCreateCmd(f))
	cmd.AddCommand(ecloudNetworkRulePortUpdateCmd(f))
	cmd.AddCommand(ecloudNetworkRulePortDeleteCmd(f))

	return cmd
}

func ecloudNetworkRulePortListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Lists network rule ports",
		Long:    "This command lists network rule ports",
		Example: "ans ecloud networkruleport list",
		RunE:    ecloudCobraRunEFunc(f, ecloudNetworkRulePortList),
	}

	cmd.Flags().String("rule", "", "Network rule ID for filtering")

	return cmd
}

func ecloudNetworkRulePortList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd,
		helper.NewStringFilterFlagOption("name", "name"),
		helper.NewStringFilterFlagOption("rule", "network_rule_id"),
	)
	if err != nil {
		return err
	}

	rules, err := service.GetNetworkRulePorts(params)
	if err != nil {
		return fmt.Errorf("Error retrieving network rule ports: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudNetworkRulePortsProvider(rules))
}

func ecloudNetworkRulePortShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <port: id>...",
		Short:   "Shows a network rule port",
		Long:    "This command shows one or more network rule ports",
		Example: "ans ecloud networkruleport show nrp-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing network rule port")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudNetworkRulePortShow),
	}
}

func ecloudNetworkRulePortShow(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	var rules []ecloud.NetworkRulePort
	for _, arg := range args {
		rule, err := service.GetNetworkRulePort(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving network rule port [%s]: %s", arg, err)
			continue
		}

		rules = append(rules, rule)
	}

	return output.CommandOutput(cmd, OutputECloudNetworkRulePortsProvider(rules))
}

func ecloudNetworkRulePortCreateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Creates a network rule port",
		Long:    "This command creates a network rule port",
		Example: "ans ecloud networkruleport create --rule nr-abcdef12",
		RunE:    ecloudCobraRunEFunc(f, ecloudNetworkRulePortCreate),
	}

	// Setup flags
	cmd.Flags().String("rule", "", "ID of network rule")
	cmd.MarkFlagRequired("rule")
	cmd.Flags().String("source", "", "Source port. Single port, port range, or ANY")
	cmd.Flags().String("destination", "", "Destination port. Single port, port range, or ANY")
	cmd.Flags().String("protocol", "", "Protocol of port. One of: TCP/UDP/ICMPv4")
	cmd.MarkFlagRequired("protocol")
	cmd.Flags().String("name", "", "Name of port")
	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the network rule port has been completely created")

	return cmd
}

func ecloudNetworkRulePortCreate(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	createRequest := ecloud.CreateNetworkRulePortRequest{}
	createRequest.NetworkRuleID, _ = cmd.Flags().GetString("rule")
	createRequest.Source, _ = cmd.Flags().GetString("source")
	createRequest.Destination, _ = cmd.Flags().GetString("destination")

	protocol, _ := cmd.Flags().GetString("protocol")
	protocolParsed, err := ecloud.ParseNetworkRulePortProtocol(protocol)
	if err != nil {
		return err
	}
	createRequest.Protocol = protocolParsed

	if cmd.Flags().Changed("name") {
		createRequest.Name, _ = cmd.Flags().GetString("name")
	}

	taskRef, err := service.CreateNetworkRulePort(createRequest)
	if err != nil {
		return fmt.Errorf("Error creating network rule port: %s", err)
	}

	waitFlag, _ := cmd.Flags().GetBool("wait")
	if waitFlag {
		err := helper.WaitForCommand(TaskStatusWaitFunc(service, taskRef.TaskID, ecloud.TaskStatusComplete))
		if err != nil {
			return fmt.Errorf("Error waiting for network rule port task to complete: %s", err)
		}
	}

	rule, err := service.GetNetworkRulePort(taskRef.ResourceID)
	if err != nil {
		return fmt.Errorf("Error retrieving new network rule port: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudNetworkRulePortsProvider([]ecloud.NetworkRulePort{rule}))
}

func ecloudNetworkRulePortUpdateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <port: id>...",
		Short:   "Updates a network rule port",
		Long:    "This command updates one or more network rule ports",
		Example: "ans ecloud networkruleport update nrp-abcdef12 --name \"my port\"",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing network rule port")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudNetworkRulePortUpdate),
	}

	cmd.Flags().String("source", "", "Source port. Single port, port range, or ANY")
	cmd.Flags().String("destination", "", "Destination port. Single port, port range, or ANY")
	cmd.Flags().String("protocol", "", "Protocol of port. One of: TCP/UDP/ICMPv4")
	cmd.Flags().String("name", "", "Name of port")
	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the network rule port has been completely updated")

	return cmd
}

func ecloudNetworkRulePortUpdate(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	patchRequest := ecloud.PatchNetworkRulePortRequest{}

	if cmd.Flags().Changed("source") {
		patchRequest.Source, _ = cmd.Flags().GetString("source")
	}

	if cmd.Flags().Changed("destination") {
		patchRequest.Destination, _ = cmd.Flags().GetString("destination")
	}

	if cmd.Flags().Changed("protocol") {

		protocol, _ := cmd.Flags().GetString("protocol")
		protocolParsed, err := ecloud.ParseNetworkRulePortProtocol(protocol)
		if err != nil {
			return err
		}
		patchRequest.Protocol = protocolParsed
	}

	if cmd.Flags().Changed("name") {
		patchRequest.Name, _ = cmd.Flags().GetString("name")
	}

	var rules []ecloud.NetworkRulePort
	for _, arg := range args {
		task, err := service.PatchNetworkRulePort(arg, patchRequest)
		if err != nil {
			output.OutputWithErrorLevelf("Error updating network rule port [%s]: %s", arg, err)
			continue
		}

		waitFlag, _ := cmd.Flags().GetBool("wait")
		if waitFlag {
			err := helper.WaitForCommand(TaskStatusWaitFunc(service, task.TaskID, ecloud.TaskStatusComplete))
			if err != nil {
				output.OutputWithErrorLevelf("Error waiting for task to complete for network rule port [%s]: %s", arg, err)
				continue
			}
		}

		rule, err := service.GetNetworkRulePort(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving updated network rule port [%s]: %s", arg, err)
			continue
		}

		rules = append(rules, rule)
	}

	return output.CommandOutput(cmd, OutputECloudNetworkRulePortsProvider(rules))
}

func ecloudNetworkRulePortDeleteCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete <port: id>...",
		Short:   "Removes a network rule port",
		Long:    "This command removes one or more network rule ports",
		Example: "ans ecloud networkruleport delete nrp-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing network rule port")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudNetworkRulePortDelete),
	}

	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the network rule port has been completely removed")

	return cmd
}

func ecloudNetworkRulePortDelete(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	for _, arg := range args {
		taskID, err := service.DeleteNetworkRulePort(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error removing network rule port [%s]: %s", arg, err)
			continue
		}

		waitFlag, _ := cmd.Flags().GetBool("wait")
		if waitFlag {
			err := helper.WaitForCommand(TaskStatusWaitFunc(service, taskID, ecloud.TaskStatusComplete))
			if err != nil {
				output.OutputWithErrorLevelf("Error waiting for task to complete for network rule port [%s]: %s", arg, err)
				continue
			}
		}
	}
	return nil
}
