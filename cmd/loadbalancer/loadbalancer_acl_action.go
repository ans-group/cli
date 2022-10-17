package loadbalancer

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/loadbalancer"
	"github.com/spf13/cobra"
)

func loadbalancerACLActionRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "action",
		Short: "sub-commands relating to ACLs",
	}

	// Child commands
	cmd.AddCommand(loadbalancerACLActionListCmd(f))
	cmd.AddCommand(loadbalancerACLActionShowCmd(f))
	cmd.AddCommand(loadbalancerACLActionCreateCmd(f))
	cmd.AddCommand(loadbalancerACLActionUpdateCmd(f))
	cmd.AddCommand(loadbalancerACLActionDeleteCmd(f))

	return cmd
}
func loadbalancerACLActionListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list <acl: id>",
		Short:   "Lists ACL actions",
		Long:    "This command lists ACL actions",
		Example: "ans loadbalancer acl action list 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing ACL")
			}

			return nil
		},
		RunE: loadbalancerCobraRunEFunc(f, loadbalancerACLActionList),
	}
}

func loadbalancerACLActionList(service loadbalancer.LoadBalancerService, cmd *cobra.Command, args []string) error {
	aclID, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("Invalid ACL ID [%s]", args[0])
	}

	acl, err := service.GetACL(aclID)
	if err != nil {
		return fmt.Errorf("Error retrieving ACL [%d]: %s", aclID, err)
	}

	return output.CommandOutput(cmd, OutputLoadBalancerACLActionsProvider(mapACLActions(acl.Actions)))
}

func loadbalancerACLActionShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <acl: id> <action: index>...",
		Short:   "Shows an ACL action",
		Long:    "This command shows one or more ACL actions",
		Example: "ans loadbalancer acl action show 123 0",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing ACL")
			}
			if len(args) < 2 {
				return errors.New("Missing ACL action index")
			}

			return nil
		},
		RunE: loadbalancerCobraRunEFunc(f, loadbalancerACLActionShow),
	}
}

func loadbalancerACLActionShow(service loadbalancer.LoadBalancerService, cmd *cobra.Command, args []string) error {
	var actions []ACLAction
	aclID, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("Invalid ACL ID [%s]", args[0])
	}

	acl, err := service.GetACL(aclID)
	if err != nil {
		return fmt.Errorf("Error retrieving ACL [%d]: %s", aclID, err)
	}

	for _, arg := range args[1:] {
		actionIndex, err := strconv.Atoi(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Invalid ACL action index [%s]", arg)
			continue
		}

		if len(acl.Actions) < actionIndex+1 {
			output.OutputWithErrorLevelf("ACL action index [%s] out of bounds", arg)
			continue
		}

		actions = append(actions, mapACLAction(acl.Actions[actionIndex], actionIndex))
	}

	return output.CommandOutput(cmd, OutputLoadBalancerACLActionsProvider(actions))
}

func loadbalancerACLActionCreateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Creates an ACL action",
		Long:    "This command creates an ACL action",
		Example: "ans loadbalancer acl action create 123 --name \"header_matches\" --argument \"host=ans.co.uk\"",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing ACL")
			}

			return nil
		},
		RunE: loadbalancerCobraRunEFunc(f, loadbalancerACLActionCreate),
	}

	cmd.Flags().String("name", "", "Name of ACL action")
	cmd.MarkFlagRequired("name")
	cmd.Flags().StringArray("argument", []string{}, "ACL action argument. Can be repeated")

	return cmd
}

func loadbalancerACLActionCreate(service loadbalancer.LoadBalancerService, cmd *cobra.Command, args []string) error {
	aclID, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("Invalid ACL ID [%s]", args[0])
	}

	acl, err := service.GetACL(aclID)
	if err != nil {
		return fmt.Errorf("Error retrieving ACL: %s", err)
	}

	action := loadbalancer.ACLAction{}
	action.Name, _ = cmd.Flags().GetString("name")

	if cmd.Flags().Changed("argument") {
		actionArguments, _ := cmd.Flags().GetStringArray("argument")
		actionArgumentsParsed, err := parseACLArguments(actionArguments)
		if err != nil {
			return fmt.Errorf("Failed to parse arguments: %s", err)
		}

		action.Arguments = actionArgumentsParsed
	}

	updateRequest := loadbalancer.PatchACLRequest{
		Actions: append(acl.Actions, action),
	}

	err = service.PatchACL(aclID, updateRequest)
	if err != nil {
		return fmt.Errorf("Error updating ACL: %s", err)
	}

	acl, err = service.GetACL(aclID)
	if err != nil {
		return fmt.Errorf("Error retrieving updated ACL [%d]: %s", aclID, err)
	}

	return output.CommandOutput(cmd, OutputLoadBalancerACLActionsProvider(mapACLActions(acl.Actions)))
}

func loadbalancerACLActionUpdateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <acl: id> <action: index>",
		Short:   "Updates an ACL action",
		Long:    "This command updates an ACL action",
		Example: "ans loadbalancer acl action update 123 0 --name \"header_matches\"",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing ACL")
			}
			if len(args) < 2 {
				return errors.New("Missing ACL action index")
			}

			return nil
		},
		RunE: loadbalancerCobraRunEFunc(f, loadbalancerACLActionUpdate),
	}

	cmd.Flags().String("name", "", "Name of ACL action")
	cmd.Flags().StringArray("argument", []string{}, "ACL action argument. Can be repeated")

	return cmd
}

func loadbalancerACLActionUpdate(service loadbalancer.LoadBalancerService, cmd *cobra.Command, args []string) error {
	aclID, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("Invalid ACL ID [%s]", args[0])
	}

	acl, err := service.GetACL(aclID)
	if err != nil {
		return fmt.Errorf("Error retrieving ACL: %s", err)
	}

	actionIndex, err := strconv.Atoi(args[1])
	if err != nil {
		return fmt.Errorf("Invalid ACL action index [%s]", args[1])
	}

	if len(acl.Actions) < actionIndex+1 {
		return fmt.Errorf("ACL action index [%d] out of bounds", actionIndex)
	}

	changed := false
	if cmd.Flags().Changed("name") {
		acl.Actions[actionIndex].Name, _ = cmd.Flags().GetString("name")
		changed = true
	}

	if cmd.Flags().Changed("argument") {
		actionArguments, _ := cmd.Flags().GetStringArray("argument")
		actionArgumentsParsed, err := parseACLArguments(actionArguments)
		if err != nil {
			return fmt.Errorf("Failed to parse arguments: %s", err)
		}

		acl.Actions[actionIndex].Arguments = actionArgumentsParsed
		changed = true
	}

	if changed {
		updateRequest := loadbalancer.PatchACLRequest{
			Actions: acl.Actions,
		}

		err = service.PatchACL(aclID, updateRequest)
		if err != nil {
			return fmt.Errorf("Error updating ACL: %s", err)
		}
	}

	acl, err = service.GetACL(aclID)
	if err != nil {
		return fmt.Errorf("Error retrieving updated ACL [%d]: %s", aclID, err)
	}

	return output.CommandOutput(cmd, OutputLoadBalancerACLActionsProvider(mapACLActions(acl.Actions)))
}

func loadbalancerACLActionDeleteCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "delete <acl: id> <action: index>",
		Short:   "Removes a acl",
		Long:    "This command removes one or more acls",
		Example: "ans loadbalancer acl delete 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing ACL")
			}
			if len(args) < 2 {
				return errors.New("Missing ACL action index")
			}

			return nil
		},
		RunE: loadbalancerCobraRunEFunc(f, loadbalancerACLActionDelete),
	}
}

func loadbalancerACLActionDelete(service loadbalancer.LoadBalancerService, cmd *cobra.Command, args []string) error {
	aclID, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("Invalid ACL ID [%s]", args[0])
	}

	acl, err := service.GetACL(aclID)
	if err != nil {
		return fmt.Errorf("Error retrieving ACL: %s", err)
	}

	actionIndex, err := strconv.Atoi(args[1])
	if err != nil {
		return fmt.Errorf("Invalid ACL action index [%s]", args[1])
	}

	if len(acl.Actions) < actionIndex+1 {
		return fmt.Errorf("ACL action index [%d] out of bounds", actionIndex)
	}

	actions := make([]loadbalancer.ACLAction, 0)
	actions = append(actions, acl.Actions[:actionIndex]...)
	actions = append(actions, acl.Actions[actionIndex+1:]...)

	updateRequest := loadbalancer.PatchACLRequest{
		Actions: actions,
	}

	err = service.PatchACL(aclID, updateRequest)
	if err != nil {
		return fmt.Errorf("Error updating ACL: %s", err)
	}

	acl, err = service.GetACL(aclID)
	if err != nil {
		return fmt.Errorf("Error retrieving updated ACL [%d]: %s", aclID, err)
	}

	return output.CommandOutput(cmd, OutputLoadBalancerACLActionsProvider(mapACLActions(acl.Actions)))
}

func mapACLAction(action loadbalancer.ACLAction, index int) ACLAction {
	return ACLAction{
		ACLAction: action,
		Index:     index,
	}
}

func mapACLActions(actions []loadbalancer.ACLAction) []ACLAction {
	var aclActions []ACLAction
	for actionIndex, action := range actions {
		aclActions = append(aclActions, mapACLAction(action, actionIndex))
	}

	return aclActions
}
