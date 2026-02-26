package loadbalancer

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/loadbalancer"
	"github.com/spf13/cobra"
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

	// Child root commands
	cmd.AddCommand(loadbalancerACLConditionRootCmd(f))
	cmd.AddCommand(loadbalancerACLActionRootCmd(f))

	return cmd
}

func loadbalancerACLShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <acl: id>...",
		Short:   "Shows an ACL",
		Long:    "This command shows one or more ACLs",
		Example: "ans loadbalancer acl show 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("missing ACL")
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

	return output.CommandOutput(cmd, ACLCollection(acls))
}

func loadbalancerACLCreateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Creates an ACL",
		Long:    "This command creates a ACLs with a single condition/action. Additional conditions/actions can be added with subcommands",
		Example: "ans loadbalancer acl create --name \"test ACL\" --host-group 1 --condition \"header_matches:host=ans.co.uk,accept=application/json\" --action \"redirect:location=developers.ans.co.uk,status=302\"",
		RunE:    loadbalancerCobraRunEFunc(f, loadbalancerACLCreate),
	}

	cmd.Flags().String("name", "", "Name of ACL")
	_ = cmd.MarkFlagRequired("name")
	cmd.Flags().Int("priority", 0, "Priority of ACL")
	cmd.Flags().Int("listener", 0, "ID of listener")
	cmd.Flags().Int("target-group", 0, "ID of target group")
	cmd.Flags().String("condition-name", "", "Name of ACL condition")
	cmd.Flags().StringArray("condition-argument", []string{}, "ACL condition argument. Can be repeated")
	cmd.Flags().Bool("condition-inverted", false, "Specifies ACL condition should be inverted")
	cmd.Flags().String("action-name", "", "Name of ACL action")
	_ = cmd.MarkFlagRequired("action-name")
	cmd.Flags().StringArray("action-argument", []string{}, "ACL action argument. Can be repeated")

	return cmd
}

func loadbalancerACLCreate(service loadbalancer.LoadBalancerService, cmd *cobra.Command, args []string) error {
	createRequest := loadbalancer.CreateACLRequest{}
	createRequest.Name, _ = cmd.Flags().GetString("name")
	createRequest.Priority, _ = cmd.Flags().GetInt("priority")
	createRequest.ListenerID, _ = cmd.Flags().GetInt("listener")
	createRequest.TargetGroupID, _ = cmd.Flags().GetInt("target-group")

	if cmd.Flags().Changed("condition-name") {
		condition := loadbalancer.ACLCondition{}
		condition.Name, _ = cmd.Flags().GetString("condition-name")
		condition.Inverted, _ = cmd.Flags().GetBool("condition-inverted")

		if cmd.Flags().Changed("condition-argument") {
			condition.Arguments = make(map[string]loadbalancer.ACLArgument)
			conditionArguments, _ := cmd.Flags().GetStringArray("condition-argument")
			conditionArgumentsParsed, err := parseACLArguments(conditionArguments)
			if err != nil {
				return fmt.Errorf("failed to parse arguments: %s", err)
			}

			condition.Arguments = conditionArgumentsParsed
		}

		createRequest.Conditions = []loadbalancer.ACLCondition{condition}
	}

	if cmd.Flags().Changed("action-name") {
		action := loadbalancer.ACLAction{}
		action.Name, _ = cmd.Flags().GetString("action-name")

		if cmd.Flags().Changed("action-argument") {
			action.Arguments = make(map[string]loadbalancer.ACLArgument)
			actionArguments, _ := cmd.Flags().GetStringArray("action-argument")

			for _, actionArgument := range actionArguments {
				key, value := parseACLArgumentKV(actionArgument)
				action.Arguments[key] = loadbalancer.ACLArgument{
					Name:  key,
					Value: value,
				}
			}
		}

		createRequest.Actions = []loadbalancer.ACLAction{action}
	}

	aclID, err := service.CreateACL(createRequest)
	if err != nil {
		return fmt.Errorf("error creating ACL: %s", err)
	}

	acl, err := service.GetACL(aclID)
	if err != nil {
		return fmt.Errorf("error retrieving new ACL: %s", err)
	}

	return output.CommandOutput(cmd, ACLCollection([]loadbalancer.ACL{acl}))
}

func loadbalancerACLUpdateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <acl: id>...",
		Short:   "Updates an ACL",
		Long:    "This command updates one or more ACLs",
		Example: "ans loadbalancer acl update 123 --name myacl",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("missing ACL")
			}

			return nil
		},
		RunE: loadbalancerCobraRunEFunc(f, loadbalancerACLUpdate),
	}

	cmd.Flags().String("name", "", "Name of ACL")
	cmd.Flags().Int("priority", 0, "Priority of ACL")

	return cmd
}

func loadbalancerACLUpdate(service loadbalancer.LoadBalancerService, cmd *cobra.Command, args []string) error {
	patchRequest := loadbalancer.PatchACLRequest{}
	patchRequest.Name, _ = cmd.Flags().GetString("name")
	patchRequest.Priority, _ = cmd.Flags().GetInt("priority")

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

	return output.CommandOutput(cmd, ACLCollection(acls))
}

func loadbalancerACLDeleteCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "delete <acl: id>...",
		Short:   "Removes a acl",
		Long:    "This command removes one or more acls",
		Example: "ans loadbalancer acl delete 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("missing ACL")
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

type aclArgument struct {
	Name  string
	Value any
	Array bool
}

func parseACLArguments(args []string) (map[string]loadbalancer.ACLArgument, error) {
	var tmpArguments []*aclArgument
	for _, arg := range args {
		parts := strings.Split(arg, "=")
		if len(parts) != 2 {
			return nil, errors.New("invalid arguments format. Expected format name=value")
		}

		argName := parts[0]
		argValue := parts[1]
		existingArg := false

		if before, ok := strings.CutSuffix(argName, "[]"); ok {
			argNameTrimmed := before

			for _, searchArg := range tmpArguments {
				if searchArg.Name == argNameTrimmed && searchArg.Array {
					existingArg = true
					searchArg.Value = append(searchArg.Value.([]string), argValue)
					break
				}
			}

			if !existingArg {
				tmpArguments = append(tmpArguments, &aclArgument{
					Name:  argNameTrimmed,
					Value: []string{argValue},
					Array: true,
				})
			}

			continue
		}

		tmpArguments = append(tmpArguments, &aclArgument{
			Name:  argName,
			Value: argValue,
		})
	}

	arguments := make(map[string]loadbalancer.ACLArgument)
	for _, tmpArgument := range tmpArguments {
		arguments[tmpArgument.Name] = loadbalancer.ACLArgument{
			Name:  tmpArgument.Name,
			Value: tmpArgument.Value,
		}
	}
	return arguments, nil
}

func parseACLArgumentKV(arg string) (string, string) {
	var value string

	argParts := strings.SplitN(arg, "=", 2)

	if len(argParts) == 2 {
		value = argParts[1]
	}

	return argParts[0], value
}
