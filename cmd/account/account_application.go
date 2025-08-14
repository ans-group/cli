package account

import (
	"errors"
	"fmt"
	"net"

	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/helper"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/account"
	"github.com/spf13/cobra"
)

func accountApplicationRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "application",
		Short: "sub-commands relating to applications",
	}

	// Child commands
	cmd.AddCommand(accountApplicationListCmd(f))
	cmd.AddCommand(accountApplicationShowCmd(f))
	cmd.AddCommand(accountApplicationCreateCmd(f))
	cmd.AddCommand(accountApplicationUpdateCmd(f))
	cmd.AddCommand(accountApplicationDeleteCmd(f))
	cmd.AddCommand(accountApplicationRestrictionsRootCmd(f))

	return cmd
}

func accountApplicationListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists applications",
		Long:    "This command lists applications",
		Example: "ans account application list",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return accountApplicationList(c.AccountService(), cmd, args)
		},
	}
}

func accountApplicationList(service account.AccountService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	applications, err := service.GetApplications(params)
	if err != nil {
		return fmt.Errorf("account: error retrieving applications: %s", err)
	}

	return output.CommandOutput(cmd, ApplicationCollection(applications))
}

func accountApplicationShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <application: id>...",
		Short:   "Shows an application",
		Long:    "This command shows one or more applications",
		Example: "ans account application show 550e8400-e29b-41d4-a716-446655440000",
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

			return accountApplicationShow(c.AccountService(), cmd, args)
		},
	}
}

func accountApplicationShow(service account.AccountService, cmd *cobra.Command, args []string) error {
	var applications []account.Application
	for _, arg := range args {
		application, err := service.GetApplication(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving application [%s]: %s", arg, err)
			continue
		}

		applications = append(applications, application)
	}

	return output.CommandOutput(cmd, ApplicationCollection(applications))
}

func accountApplicationCreateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Creates an application",
		Long:    "This command creates an application",
		Example: "ans account application create --name \"My App\" --description \"My application description\"",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return accountApplicationCreate(c.AccountService(), cmd, args)
		},
	}

	// Setup flags
	cmd.Flags().String("name", "", "Name of application")
	cmd.MarkFlagRequired("name")
	cmd.Flags().String("description", "", "Description of application")
	cmd.Flags().StringSlice("allow-ip", []string{}, "IP addresses/ranges to allow (sets allowlist)")
	cmd.Flags().StringSlice("deny-ip", []string{}, "IP addresses/ranges to deny (sets denylist)")

	return cmd
}

func accountApplicationCreate(service account.AccountService, cmd *cobra.Command, args []string) error {
	createRequest := account.CreateApplicationRequest{}
	createRequest.Name, _ = cmd.Flags().GetString("name")
	createRequest.Description, _ = cmd.Flags().GetString("description")

	// Validate IP restriction flags are not both specified
	allowIPs, _ := cmd.Flags().GetStringSlice("allow-ip")
	denyIPs, _ := cmd.Flags().GetStringSlice("deny-ip")

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

	response, err := service.CreateApplication(createRequest)
	if err != nil {
		return fmt.Errorf("account: error creating application: %s", err)
	}

	// Set IP restrictions if specified
	if len(allowIPs) > 0 || len(denyIPs) > 0 {
		restrictionRequest := account.SetRestrictionRequest{}
		if len(allowIPs) > 0 {
			restrictionRequest.IPRestrictionType = "allowlist"
			restrictionRequest.IPRanges = allowIPs
		} else {
			restrictionRequest.IPRestrictionType = "denylist"
			restrictionRequest.IPRanges = denyIPs
		}

		err = service.SetApplicationRestrictions(response.ID, restrictionRequest)
		if err != nil {
			output.OutputWithErrorLevelf("Warning: Application created but failed to set IP restrictions: %s", err)
		}
	}

	// Get the full application details to display
	application, err := service.GetApplication(response.ID)
	if err != nil {
		return fmt.Errorf("account: error retrieving new application: %s", err)
	}

	return output.CommandOutput(cmd, ApplicationCollection([]account.Application{application}))
}

func accountApplicationUpdateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <application: id>...",
		Short:   "Updates an application",
		Long:    "This command updates one or more applications",
		Example: "ans account application update 550e8400-e29b-41d4-a716-446655440000 --name \"Updated App\"",
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

			return accountApplicationUpdate(c.AccountService(), cmd, args)
		},
	}

	// Setup flags
	cmd.Flags().String("name", "", "Name of application")
	cmd.Flags().String("description", "", "Description of application")

	return cmd
}

func accountApplicationUpdate(service account.AccountService, cmd *cobra.Command, args []string) error {
	updateRequest := account.UpdateApplicationRequest{}
	updateRequest.Name, _ = cmd.Flags().GetString("name")
	updateRequest.Description, _ = cmd.Flags().GetString("description")

	var applications []account.Application

	for _, arg := range args {
		err := service.UpdateApplication(arg, updateRequest)
		if err != nil {
			output.OutputWithErrorLevelf("Error updating application [%s]: %s", arg, err.Error())
			continue
		}

		application, err := service.GetApplication(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving updated application [%s]: %s", arg, err.Error())
			continue
		}

		applications = append(applications, application)
	}

	return output.CommandOutput(cmd, ApplicationCollection(applications))
}

func accountApplicationDeleteCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "delete <application: id>...",
		Short:   "Removes an application",
		Long:    "This command removes one or more applications",
		Example: "ans account application delete 550e8400-e29b-41d4-a716-446655440000",
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

			return accountApplicationDelete(c.AccountService(), cmd, args)
		},
	}
}

func accountApplicationDelete(service account.AccountService, cmd *cobra.Command, args []string) error {
	var hasErrors bool
	for _, arg := range args {
		err := service.DeleteApplication(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error removing application [%s]: %s", arg, err)
			hasErrors = true
			continue
		}
	}

	if hasErrors {
		return fmt.Errorf("account: failed to delete one or more applications")
	}
	return nil
}

// validateIPRanges validates that all provided strings are valid IP addresses or CIDR ranges
func validateIPRanges(ipRanges []string) error {
	for _, ipRange := range ipRanges {
		// Try parsing as CIDR first
		_, _, err := net.ParseCIDR(ipRange)
		if err != nil {
			// If not CIDR, try parsing as IP
			ip := net.ParseIP(ipRange)
			if ip == nil {
				return fmt.Errorf("invalid IP address or CIDR range: %s", ipRange)
			}
		}
	}
	return nil
}
