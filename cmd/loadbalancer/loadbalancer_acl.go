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
		Example: "ukfast loadbalancer acl create --name \"test ACL\" --host-group 1 --condition \"header_matches:host=ukfast.co.uk,accept=application/json\" --action \"redirect:location=developers.ukfast.io,status=302\"",
		RunE:    loadbalancerCobraRunEFunc(f, loadbalancerACLCreate),
	}

	cmd.Flags().String("name", "", "Name of ACL")
	cmd.Flags().Int("listener", 0, "ID of listener")
	cmd.Flags().Int("target-group", 0, "ID of target group")
	cmd.Flags().StringArray("condition", []string{}, "Name and arguments of condition. Can be repeated. Example: --condition \"header_matches:host=ukfast.co.uk,accept=application/json\"")
	cmd.Flags().StringArray("action", []string{}, "Name and arguments of action. Can be repeated. Example: --action \"redirect:location=developers.ukfast.io,status=302\"")
	cmd.MarkFlagRequired("action")

	return cmd
}

func loadbalancerACLCreate(service loadbalancer.LoadBalancerService, cmd *cobra.Command, args []string) error {
	createRequest := loadbalancer.CreateACLRequest{}
	createRequest.Name, _ = cmd.Flags().GetString("name")
	createRequest.ListenerID, _ = cmd.Flags().GetInt("listener")
	createRequest.TargetGroupID, _ = cmd.Flags().GetInt("target-group")

	if cmd.Flags().Changed("condition") {
		conditionsFlag, _ := cmd.Flags().GetStringArray("condition")
		conditions, err := parseACLConditionsFromFlag(conditionsFlag)
		if err != nil {
			return err
		}
		createRequest.Conditions = conditions
	}

	if cmd.Flags().Changed("action") {
		actionsFlag, _ := cmd.Flags().GetStringArray("action")
		actions, err := parseACLActionsFromFlag(actionsFlag)
		if err != nil {
			return err
		}
		createRequest.Actions = actions
	}

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
	cmd.Flags().StringArray("condition", []string{}, "Name and arguments of condition. Can be repeated. Example: --condition \"header_matches:host=ukfast.co.uk,accept=application/json\"")
	cmd.Flags().StringArray("action", []string{}, "Name and arguments of action. Can be repeated. Example: --action \"redirect:location=developers.ukfast.io,status=302\"")

	return cmd
}

func loadbalancerACLUpdate(service loadbalancer.LoadBalancerService, cmd *cobra.Command, args []string) error {
	patchRequest := loadbalancer.PatchACLRequest{}
	patchRequest.Name, _ = cmd.Flags().GetString("name")

	if cmd.Flags().Changed("condition") {
		conditionsFlag, _ := cmd.Flags().GetStringArray("condition")
		conditions, err := parseACLConditionsFromFlag(conditionsFlag)
		if err != nil {
			return err
		}
		patchRequest.Conditions = conditions
	}

	if cmd.Flags().Changed("action") {
		actionsFlag, _ := cmd.Flags().GetStringArray("action")
		actions, err := parseACLActionsFromFlag(actionsFlag)
		if err != nil {
			return err
		}
		patchRequest.Actions = actions
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

func parseACLActionsFromFlag(actionFlags []string) ([]loadbalancer.ACLAction, error) {
	var actions []loadbalancer.ACLAction
	for _, actionFlag := range actionFlags {
		name, arguments, err := parseACLStatementFlag(actionFlag)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse ACL action from flag: %s", err)
		}
		actions = append(actions, loadbalancer.ACLAction{
			Name:      name,
			Arguments: arguments,
		})
	}

	return actions, nil
}

func parseACLConditionsFromFlag(conditionFlags []string) ([]loadbalancer.ACLCondition, error) {
	var conditions []loadbalancer.ACLCondition
	for _, conditionFlag := range conditionFlags {
		name, arguments, err := parseACLStatementFlag(conditionFlag)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse ACL condition from flag: %s", err)
		}
		conditions = append(conditions, loadbalancer.ACLCondition{
			Name:      name,
			Arguments: arguments,
		})
	}

	return conditions, nil
}

func parseACLStatementFlag(flag string) (string, map[string]loadbalancer.ACLArgument, error) {
	flagNameSplit := strings.SplitN(flag, ":", 2)
	if len(flagNameSplit) != 2 {
		return "", nil, fmt.Errorf("Invalid flag format. Expected format name:arguments")
	}

	flagArgsSplit := strings.Split(flagNameSplit[1], ",")
	arguments, err := parseACLArguments(flagArgsSplit)
	if err != nil {
		return "", nil, fmt.Errorf("Invalid flag arguments format: %s", err)
	}

	return flagNameSplit[0], arguments, nil
}

func parseACLArguments(args []string) (map[string]loadbalancer.ACLArgument, error) {
	arguments := make(map[string]loadbalancer.ACLArgument)
	for _, arg := range args {
		parts := strings.Split(arg, "=")
		if len(parts) != 2 {
			return nil, fmt.Errorf("Invalid arguments format. Expected format name=value")
		}

		arguments[parts[0]] = loadbalancer.ACLArgument{
			Name:  parts[0],
			Value: parts[1],
		}
	}

	return arguments, nil
}
