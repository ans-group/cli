package account

import (
	"errors"
	"fmt"

	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/account"
	"github.com/spf13/cobra"
)

func accountApplicationRestrictionsRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "restrictions",
		Short: "sub-commands relating to application IP restrictions",
	}

	cmd.AddCommand(accountApplicationRestrictionsShowCmd(f))
	cmd.AddCommand(accountApplicationRestrictionsSetCmd(f))
	cmd.AddCommand(accountApplicationRestrictionsClearCmd(f))

	return cmd
}

func accountApplicationRestrictionsShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <application: id>...",
		Short:   "Shows application IP restrictions",
		Long:    "This command shows IP restrictions for one or more applications",
		Example: "ans account application restrictions show 550e8400-e29b-41d4-a716-446655440000",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("missing application")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return accountApplicationRestrictionsShow(c.AccountService(), cmd, args)
		},
	}
}

func accountApplicationRestrictionsShow(service account.AccountService, cmd *cobra.Command, args []string) error {
	var restrictions []ApplicationRestrictionWithID
	for _, arg := range args {
		restriction, err := service.GetApplicationRestrictions(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving restrictions for application [%s]: %s", arg, err)
			continue
		}

		// If IPRestrictionType is empty, no restrictions are set
		restrictionType := restriction.IPRestrictionType
		if restrictionType == "" {
			restrictionType = "none"
		}

		restrictions = append(restrictions, ApplicationRestrictionWithID{
			ID:                arg,
			IPRestrictionType: restrictionType,
			IPRanges:          restriction.IPRanges,
		})
	}

	return output.CommandOutput(cmd, ApplicationRestrictionCollection(restrictions))
}

func accountApplicationRestrictionsSetCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "set <application: id>",
		Short:   "Sets IP restrictions for an application",
		Long:    "This command sets IP restrictions for an application",
		Example: "ans account application restrictions set 550e8400-e29b-41d4-a716-446655440000 --allow-ip \"198.51.100.0/24\"",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("missing application")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return accountApplicationRestrictionsSet(c.AccountService(), cmd, args)
		},
	}

	// Setup flags
	cmd.Flags().StringSlice("allow-ip", []string{}, "IP addresses/ranges to allow (sets allowlist)")
	cmd.Flags().StringSlice("deny-ip", []string{}, "IP addresses/ranges to deny (sets denylist)")

	return cmd
}

func accountApplicationRestrictionsSet(service account.AccountService, cmd *cobra.Command, args []string) error {
	appID := args[0]

	// Validate IP restriction flags
	allowIPs, _ := cmd.Flags().GetStringSlice("allow-ip")
	denyIPs, _ := cmd.Flags().GetStringSlice("deny-ip")

	if len(allowIPs) == 0 && len(denyIPs) == 0 {
		return fmt.Errorf("account: must specify either --allow-ip or --deny-ip")
	}

	if len(allowIPs) > 0 && len(denyIPs) > 0 {
		return fmt.Errorf("account: cannot specify both --allow-ip and --deny-ip")
	}

	// Validate IP addresses/ranges
	if err := validateIPRanges(allowIPs); err != nil {
		return fmt.Errorf("account: invalid IP in --allow-ip: %s", err)
	}
	if err := validateIPRanges(denyIPs); err != nil {
		return fmt.Errorf("account: invalid IP in --deny-ip: %s", err)
	}

	// Build restriction request
	restrictionRequest := account.SetRestrictionRequest{}
	if len(allowIPs) > 0 {
		restrictionRequest.IPRestrictionType = "allowlist"
		restrictionRequest.IPRanges = allowIPs
	} else {
		restrictionRequest.IPRestrictionType = "denylist"
		restrictionRequest.IPRanges = denyIPs
	}

	err := service.SetApplicationRestrictions(appID, restrictionRequest)
	if err != nil {
		return fmt.Errorf("account: error setting application restrictions: %s", err)
	}

	// Get and display the updated restrictions
	restriction, err := service.GetApplicationRestrictions(appID)
	if err != nil {
		// If we can't retrieve the updated restrictions, at least show what we tried to set
		return output.CommandOutput(cmd, ApplicationRestrictionCollection([]ApplicationRestrictionWithID{{
			ID:                appID,
			IPRestrictionType: restrictionRequest.IPRestrictionType,
			IPRanges:          restrictionRequest.IPRanges,
		}}))
	}

	return output.CommandOutput(cmd, ApplicationRestrictionCollection([]ApplicationRestrictionWithID{{
		ID:                appID,
		IPRestrictionType: restriction.IPRestrictionType,
		IPRanges:          restriction.IPRanges,
	}}))
}

func accountApplicationRestrictionsClearCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "clear <application: id>...",
		Short:   "Removes IP restrictions from applications",
		Long:    "This command removes IP restrictions from one or more applications",
		Example: "ans account application restrictions clear 550e8400-e29b-41d4-a716-446655440000",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("missing application")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return accountApplicationRestrictionsClear(c.AccountService(), cmd, args)
		},
	}
}

func accountApplicationRestrictionsClear(service account.AccountService, cmd *cobra.Command, args []string) error {
	var hasErrors bool
	for _, arg := range args {
		err := service.DeleteApplicationRestrictions(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error clearing restrictions for application [%s]: %s", arg, err)
			hasErrors = true
			continue
		}
	}

	if hasErrors {
		return fmt.Errorf("account: failed to clear restrictions for one or more applications")
	}
	return nil
}

// ApplicationRestrictionWithID adds the application ID to the restriction for output
type ApplicationRestrictionWithID struct {
	ID                string   `json:"id"`
	IPRestrictionType string   `json:"ip_restriction_type"`
	IPRanges          []string `json:"ip_ranges"`
}
