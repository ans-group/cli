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

func loadbalancerACLConditionRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "condition",
		Short: "sub-commands relating to ACLs",
	}

	// Child commands
	cmd.AddCommand(loadbalancerACLConditionListCmd(f))
	cmd.AddCommand(loadbalancerACLConditionShowCmd(f))
	cmd.AddCommand(loadbalancerACLConditionCreateCmd(f))
	cmd.AddCommand(loadbalancerACLConditionUpdateCmd(f))
	cmd.AddCommand(loadbalancerACLConditionDeleteCmd(f))

	return cmd
}
func loadbalancerACLConditionListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list <acl: id>",
		Short:   "Lists ACL conditions",
		Long:    "This command lists ACL conditions",
		Example: "ans loadbalancer acl condition list 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing ACL")
			}

			return nil
		},
		RunE: loadbalancerCobraRunEFunc(f, loadbalancerACLConditionList),
	}
}

func loadbalancerACLConditionList(service loadbalancer.LoadBalancerService, cmd *cobra.Command, args []string) error {
	aclID, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("Invalid ACL ID [%s]", args[0])
	}

	acl, err := service.GetACL(aclID)
	if err != nil {
		return fmt.Errorf("Error retrieving ACL [%d]: %s", aclID, err)
	}

	return output.CommandOutput(cmd, OutputLoadBalancerACLConditionsProvider(mapACLConditions(acl.Conditions)))
}

func loadbalancerACLConditionShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <acl: id> <condition: index>...",
		Short:   "Shows an ACL condition",
		Long:    "This command shows one or more ACL conditions",
		Example: "ans loadbalancer acl condition show 123 0",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing ACL")
			}
			if len(args) < 2 {
				return errors.New("Missing ACL condition index")
			}

			return nil
		},
		RunE: loadbalancerCobraRunEFunc(f, loadbalancerACLConditionShow),
	}
}

func loadbalancerACLConditionShow(service loadbalancer.LoadBalancerService, cmd *cobra.Command, args []string) error {
	var conditions []ACLCondition
	aclID, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("Invalid ACL ID [%s]", args[0])
	}

	acl, err := service.GetACL(aclID)
	if err != nil {
		return fmt.Errorf("Error retrieving ACL [%d]: %s", aclID, err)
	}

	for _, arg := range args[1:] {
		conditionIndex, err := strconv.Atoi(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Invalid ACL condition index [%s]", arg)
			continue
		}

		if len(acl.Conditions) < conditionIndex+1 {
			output.OutputWithErrorLevelf("ACL condition index [%s] out of bounds", arg)
			continue
		}

		conditions = append(conditions, mapACLCondition(acl.Conditions[conditionIndex], conditionIndex))
	}

	return output.CommandOutput(cmd, OutputLoadBalancerACLConditionsProvider(conditions))
}

func loadbalancerACLConditionCreateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Creates an ACL condition",
		Long:    "This command creates an ACL condition",
		Example: "ans loadbalancer acl condition create 123 --name \"header_matches\" --argument \"host=ans.co.uk\"",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing ACL")
			}

			return nil
		},
		RunE: loadbalancerCobraRunEFunc(f, loadbalancerACLConditionCreate),
	}

	cmd.Flags().String("name", "", "Name of ACL condition")
	cmd.MarkFlagRequired("name")
	cmd.Flags().StringArray("argument", []string{}, "ACL condition argument. Can be repeated")
	cmd.Flags().Bool("inverted", false, "Specifies ACL condition should be inverted")

	return cmd
}

func loadbalancerACLConditionCreate(service loadbalancer.LoadBalancerService, cmd *cobra.Command, args []string) error {
	aclID, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("Invalid ACL ID [%s]", args[0])
	}

	acl, err := service.GetACL(aclID)
	if err != nil {
		return fmt.Errorf("Error retrieving ACL: %s", err)
	}

	condition := loadbalancer.ACLCondition{}
	condition.Name, _ = cmd.Flags().GetString("name")
	condition.Inverted, _ = cmd.Flags().GetBool("inverted")

	if cmd.Flags().Changed("argument") {
		conditionArguments, _ := cmd.Flags().GetStringArray("argument")
		conditionArgumentsParsed, err := parseACLArguments(conditionArguments)
		if err != nil {
			return fmt.Errorf("Failed to parse arguments: %s", err)
		}

		condition.Arguments = conditionArgumentsParsed
	}

	updateRequest := loadbalancer.PatchACLRequest{
		Conditions: append(acl.Conditions, condition),
	}

	err = service.PatchACL(aclID, updateRequest)
	if err != nil {
		return fmt.Errorf("Error updating ACL: %s", err)
	}

	acl, err = service.GetACL(aclID)
	if err != nil {
		return fmt.Errorf("Error retrieving updated ACL [%d]: %s", aclID, err)
	}

	return output.CommandOutput(cmd, OutputLoadBalancerACLConditionsProvider(mapACLConditions(acl.Conditions)))
}

func loadbalancerACLConditionUpdateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <acl: id> <condition: index>",
		Short:   "Updates an ACL condition",
		Long:    "This command updates an ACL condition",
		Example: "ans loadbalancer acl condition update 123 0 --name \"header_matches\"",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing ACL")
			}
			if len(args) < 2 {
				return errors.New("Missing ACL condition index")
			}

			return nil
		},
		RunE: loadbalancerCobraRunEFunc(f, loadbalancerACLConditionUpdate),
	}

	cmd.Flags().String("name", "", "Name of ACL condition")
	cmd.Flags().StringArray("argument", []string{}, "ACL condition argument. Can be repeated")
	cmd.Flags().Bool("inverted", false, "Specifies ACL condition should be inverted")

	return cmd
}

func loadbalancerACLConditionUpdate(service loadbalancer.LoadBalancerService, cmd *cobra.Command, args []string) error {
	aclID, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("Invalid ACL ID [%s]", args[0])
	}

	acl, err := service.GetACL(aclID)
	if err != nil {
		return fmt.Errorf("Error retrieving ACL: %s", err)
	}

	conditionIndex, err := strconv.Atoi(args[1])
	if err != nil {
		return fmt.Errorf("Invalid ACL condition index [%s]", args[1])
	}

	if len(acl.Conditions) < conditionIndex+1 {
		return fmt.Errorf("ACL condition index [%d] out of bounds", conditionIndex)
	}

	changed := false
	if cmd.Flags().Changed("name") {
		acl.Conditions[conditionIndex].Name, _ = cmd.Flags().GetString("name")
		changed = true
	}

	if cmd.Flags().Changed("inverted") {
		acl.Conditions[conditionIndex].Inverted, _ = cmd.Flags().GetBool("inverted")
		changed = true
	}

	if cmd.Flags().Changed("argument") {
		conditionArguments, _ := cmd.Flags().GetStringArray("argument")
		conditionArgumentsParsed, err := parseACLArguments(conditionArguments)
		if err != nil {
			return fmt.Errorf("Failed to parse arguments: %s", err)
		}

		acl.Conditions[conditionIndex].Arguments = conditionArgumentsParsed
		changed = true
	}

	if changed {
		updateRequest := loadbalancer.PatchACLRequest{
			Conditions: acl.Conditions,
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

	return output.CommandOutput(cmd, OutputLoadBalancerACLConditionsProvider(mapACLConditions(acl.Conditions)))
}

func loadbalancerACLConditionDeleteCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "delete <acl: id> <condition: index>...",
		Short:   "Removes a acl",
		Long:    "This command removes one or more acls",
		Example: "ans loadbalancer acl delete 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing ACL")
			}
			if len(args) < 2 {
				return errors.New("Missing ACL condition index")
			}

			return nil
		},
		RunE: loadbalancerCobraRunEFunc(f, loadbalancerACLConditionDelete),
	}
}

func loadbalancerACLConditionDelete(service loadbalancer.LoadBalancerService, cmd *cobra.Command, args []string) error {
	aclID, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("Invalid ACL ID [%s]", args[0])
	}

	acl, err := service.GetACL(aclID)
	if err != nil {
		return fmt.Errorf("Error retrieving ACL: %s", err)
	}

	conditionIndex, err := strconv.Atoi(args[1])
	if err != nil {
		return fmt.Errorf("Invalid ACL condition index [%s]", args[1])
	}

	if len(acl.Conditions) < conditionIndex+1 {
		return fmt.Errorf("ACL condition index [%d] out of bounds", conditionIndex)
	}

	conditions := make([]loadbalancer.ACLCondition, 0)
	conditions = append(conditions, acl.Conditions[:conditionIndex]...)
	conditions = append(conditions, acl.Conditions[conditionIndex+1:]...)

	updateRequest := loadbalancer.PatchACLRequest{
		Conditions: conditions,
	}

	err = service.PatchACL(aclID, updateRequest)
	if err != nil {
		return fmt.Errorf("Error updating ACL: %s", err)
	}

	acl, err = service.GetACL(aclID)
	if err != nil {
		return fmt.Errorf("Error retrieving updated ACL [%d]: %s", aclID, err)
	}

	return output.CommandOutput(cmd, OutputLoadBalancerACLConditionsProvider(mapACLConditions(acl.Conditions)))
}

func mapACLCondition(condition loadbalancer.ACLCondition, index int) ACLCondition {
	return ACLCondition{
		ACLCondition: condition,
		Index:        index,
	}
}

func mapACLConditions(conditions []loadbalancer.ACLCondition) []ACLCondition {
	var aclConditions []ACLCondition
	for conditionIndex, condition := range conditions {
		aclConditions = append(aclConditions, mapACLCondition(condition, conditionIndex))
	}

	return aclConditions
}
