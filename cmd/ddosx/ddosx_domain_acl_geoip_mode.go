package ddosx

import (
	"errors"
	"fmt"

	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/ddosx"
	"github.com/spf13/cobra"
)

func ddosxDomainACLGeoIPRulesModeRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mode",
		Short: "sub-commands relating to domain ACL GeoIP rules mode",
	}

	// Child commands
	cmd.AddCommand(ddosxDomainACLGeoIPRulesModeShowCmd(f))
	cmd.AddCommand(ddosxDomainACLGeoIPRulesModeUpdateCmd(f))

	return cmd
}

func ddosxDomainACLGeoIPRulesModeShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <domain: name>...",
		Short:   "Shows a domain",
		Long:    "This command shows one or more domains",
		Example: "ans ddosx domain show example.com",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing domain")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ddosxDomainACLGeoIPRulesModeShow(c.DDoSXService(), cmd, args)
		},
	}
}

func ddosxDomainACLGeoIPRulesModeShow(service ddosx.DDoSXService, cmd *cobra.Command, args []string) error {
	var modes []ddosx.ACLGeoIPRulesMode
	for _, arg := range args {
		mode, err := service.GetDomainACLGeoIPRulesMode(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving domain [%s] ACL GeoIP rules mode: %s", arg, err)
			continue
		}

		modes = append(modes, mode)
	}

	return output.CommandOutput(cmd, OutputDDoSXACLGeoIPRulesModesProvider(modes))
}

func ddosxDomainACLGeoIPRulesModeUpdateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <domain: name>",
		Short:   "Updates a domain ACL GeoIP rule filtering mode",
		Long:    "This command updates a domain ACL GeoIP rule filtering mode",
		Example: "ans ddosx domain acl geoip mode update example.com --mode whitelist",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing domain")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ddosxDomainACLGeoIPRulesModeUpdate(c.DDoSXService(), cmd, args)
		},
	}

	cmd.Flags().String("mode", "", "Filtering mode for GeoIP ACL rules")

	return cmd
}

func ddosxDomainACLGeoIPRulesModeUpdate(service ddosx.DDoSXService, cmd *cobra.Command, args []string) error {
	patchRequest := ddosx.PatchACLGeoIPRulesModeRequest{}

	if cmd.Flags().Changed("mode") {
		mode, _ := cmd.Flags().GetString("mode")
		parsedMode, err := ddosx.ACLGeoIPRulesModeEnum.Parse(mode)
		if err != nil {
			return err
		}
		patchRequest.Mode = parsedMode
	}

	err := service.PatchDomainACLGeoIPRulesMode(args[0], patchRequest)
	if err != nil {
		return fmt.Errorf("Error updating domain ACL GeoIP rule filtering mode: %s", err.Error())
	}

	mode, err := service.GetDomainACLGeoIPRulesMode(args[0])
	if err != nil {
		return fmt.Errorf("Error retrieving updated domain ACL GeoIP rule filtering mode: %s", err)
	}

	return output.CommandOutput(cmd, OutputDDoSXACLGeoIPRulesModesProvider([]ddosx.ACLGeoIPRulesMode{mode}))
}
