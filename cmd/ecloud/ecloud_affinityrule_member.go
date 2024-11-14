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

func ecloudAffinityRuleAffinityMemberRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "member",
		Short: "sub-commands relating to affinity rule members",
	}

	// Child commands
	cmd.AddCommand(ecloudAffinityRuleAffinityMemberListCmd(f))

	return cmd
}

func ecloudAffinityRuleAffinityMemberListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Lists members for affinity rule",
		Long:    "This command lists members for an affinity rule",
		Example: "ans ecloud affinityrule member list ar-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing affinity rule")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudAffinityRuleAffinityMemberList),
	}

	return cmd
}

func ecloudAffinityRuleAffinityMemberList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	rules, err := service.GetAffinityRuleMembers(args[0], params)
	if err != nil {
		return fmt.Errorf("Error retrieving affinity rule members: %s", err)
	}

	return output.CommandOutput(cmd, AffinityRuleMemberCollection(rules))
}
