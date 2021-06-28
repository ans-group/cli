package loadbalancer

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/loadbalancer"
)

func loadbalancerACLRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "acl",
		Short: "sub-commands relating to ACLs",
	}

	// Child commands
	cmd.AddCommand(loadbalancerACLShowCmd(f))
	cmd.AddCommand(loadbalancerACLCreateCmd(f))
	cmd.AddCommand(loadbalancerACLUpdateCmd(f))
	cmd.AddCommand(loadbalancerACLDeleteCmd(f))

	return cmd
}

func loadbalancerACLShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <acl: id>...",
		Short:   "Shows an ACL",
		Long:    "This command shows one or more ACLs",
		Example: "ukfast loadbalancer acl show 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing ACL")
			}

			return nil
		},
		RunE: loadbalancerCobraRunEFunc(f, loadbalancerACLShow),
	}
}

func loadbalancerACLShow(service loadbalancer.LoadBalancerService, cmd *cobra.Command, args []string) error {
	var acls []loadbalancer.ACL
	for _, arg := range args {
		aclID, err := strconv.Atoi(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Invalid ACL ID [%s]", arg)
			continue
		}

		acl, err := service.GetACL(aclID)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving ACL [%d]: %s", aclID, err)
			continue
		}

		acls = append(acls, acl)
	}

	return output.CommandOutput(cmd, OutputLoadBalancerACLsProvider(acls))
}

func loadbalancerACLCreateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create <acl: id>...",
		Short:   "Creates an ACL",
		Long:    "This command creates a ACLs with a single condition/action. Additional conditions/actions can be added with subcommands",
		Example: "ukfast loadbalancer acl create --name \"test ACL\" --host-group 1 --condition-name \"header_matches\" --condition-argument \"header=host,value=ukfast.co.uk\" --action-name \"redirect\" --action-argument \"location=developers.ukfast.io,status=302\"",
		RunE:    loadbalancerCobraRunEFunc(f, loadbalancerACLCreate),
	}

	cmd.Flags().String("name", "", "Name of ACL")
	cmd.Flags().Int("listener", 0, "ID of listener")
	cmd.Flags().Int("target-group", 0, "ID of target group")
	cmd.Flags().String("condition-name", "", "Name of condition")
	cmd.Flags().StringSlice("condition-argument", []string{}, "Command-seperated arguments for condition. Can be repeated")
	cmd.Flags().String("action-name", "", "Name of action")
	cmd.MarkFlagRequired("action-name")
	cmd.Flags().StringSlice("action-argument", []string{}, "Command-seperated arguments for action. Can be repeated")

	return cmd
}

func loadbalancerACLCreate(service loadbalancer.LoadBalancerService, cmd *cobra.Command, args []string) error {
	createRequest := loadbalancer.CreateACLRequest{}
	createRequest.Name, _ = cmd.Flags().GetString("name")
	createRequest.ListenerID, _ = cmd.Flags().GetInt("listener")
	createRequest.TargetGroupID, _ = cmd.Flags().GetInt("target-group")

	if cmd.Flags().Changed("condition-name") {
		condition := loadbalancer.ACLCondition{}
		condition.Name, _ = cmd.Flags().GetString("condition-name")
		conditionArgumentsFlag, _ := cmd.Flags().GetStringSlice("condition-argument")
		conditionArguments, err := parseACLArguments(conditionArgumentsFlag)
		if err != nil {
			return err
		}
		condition.Arguments = conditionArguments
		createRequest.Conditions = []loadbalancer.ACLCondition{condition}
	}

	action := loadbalancer.ACLAction{}
	action.Name, _ = cmd.Flags().GetString("action-name")
	actionArgumentsFlag, _ := cmd.Flags().GetStringSlice("action-argument")
	actionArguments, err := parseACLArguments(actionArgumentsFlag)
	if err != nil {
		return err
	}
	action.Arguments = actionArguments
	createRequest.Actions = []loadbalancer.ACLAction{action}

	aclID, err := service.CreateACL(createRequest)
	if err != nil {
		return fmt.Errorf("Error creating ACL: %s", err)
	}

	acl, err := service.GetACL(aclID)
	if err != nil {
		return fmt.Errorf("Error retrieving new ACL: %s", err)
	}

	return output.CommandOutput(cmd, OutputLoadBalancerACLsProvider([]loadbalancer.ACL{acl}))
}

func loadbalancerACLUpdateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <acl: id>...",
		Short:   "Updates an ACL",
		Long:    "This command updates one or more ACLs",
		Example: "ukfast loadbalancer acl update 123 --name myacl",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing ACL")
			}

			return nil
		},
		RunE: loadbalancerCobraRunEFunc(f, loadbalancerACLUpdate),
	}

	cmd.Flags().String("name", "", "Name of ACL")
	cmd.Flags().String("condition-name", "", "Name of condition")
	cmd.Flags().StringSlice("condition-argument", []string{}, "Command-seperated arguments for condition. Can be repeated")
	cmd.Flags().String("action-name", "", "Name of action")
	cmd.Flags().StringSlice("action-argument", []string{}, "Command-seperated arguments for action. Can be repeated")

	return cmd
}

func loadbalancerACLUpdate(service loadbalancer.LoadBalancerService, cmd *cobra.Command, args []string) error {
	patchRequest := loadbalancer.PatchACLRequest{}
	patchRequest.Name, _ = cmd.Flags().GetString("name")

	if cmd.Flags().Changed("condition-name") {
		condition := loadbalancer.ACLCondition{}
		condition.Name, _ = cmd.Flags().GetString("condition-name")
		conditionArgumentsFlag, _ := cmd.Flags().GetStringSlice("condition-argument")
		conditionArguments, err := parseACLArguments(conditionArgumentsFlag)
		if err != nil {
			return err
		}
		condition.Arguments = conditionArguments
		patchRequest.Conditions = []loadbalancer.ACLCondition{condition}
	}

	if cmd.Flags().Changed("action-name") {
		action := loadbalancer.ACLAction{}
		action.Name, _ = cmd.Flags().GetString("action-name")
		actionArgumentsFlag, _ := cmd.Flags().GetStringSlice("action-argument")
		actionArguments, err := parseACLArguments(actionArgumentsFlag)
		if err != nil {
			return err
		}
		action.Arguments = actionArguments
		patchRequest.Actions = []loadbalancer.ACLAction{action}
	}

	var acls []loadbalancer.ACL
	for _, arg := range args {
		aclID, err := strconv.Atoi(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Invalid ACL ID [%s]", arg)
			continue
		}

		err = service.PatchACL(aclID, patchRequest)
		if err != nil {
			output.OutputWithErrorLevelf("Error updating ACL [%d]: %s", aclID, err)
			continue
		}

		acl, err := service.GetACL(aclID)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving updated ACL [%d]: %s", aclID, err)
			continue
		}

		acls = append(acls, acl)
	}

	return output.CommandOutput(cmd, OutputLoadBalancerACLsProvider(acls))
}

func loadbalancerACLDeleteCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "delete <acl: id>...",
		Short:   "Removes a acl",
		Long:    "This command removes one or more acls",
		Example: "ukfast loadbalancer acl delete 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing ACL")
			}

			return nil
		},
		RunE: loadbalancerCobraRunEFunc(f, loadbalancerACLDelete),
	}
}

func loadbalancerACLDelete(service loadbalancer.LoadBalancerService, cmd *cobra.Command, args []string) error {
	for _, arg := range args {
		aclID, err := strconv.Atoi(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Invalid ACL ID [%s]", arg)
			continue
		}

		err = service.DeleteACL(aclID)
		if err != nil {
			output.OutputWithErrorLevelf("Error removing ACL [%d]: %s", aclID, err)
			continue
		}
	}

	return nil
}

func parseACLArguments(args []string) (map[string]loadbalancer.ACLArgument, error) {
	arguments := make(map[string]loadbalancer.ACLArgument)
	for _, arg := range args {
		parts := strings.Split(arg, "=")
		if len(parts) != 2 {
			return nil, fmt.Errorf("Expected 2 parts, got %d", len(parts))
		}

		arguments[parts[0]] = loadbalancer.ACLArgument{
			Name:  parts[0],
			Value: parts[1],
		}
	}

	return arguments, nil
}
