package cmd

import (
	"errors"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/ddosx"
)

func ddosxDomainACLGeoIPRulesModeRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mode",
		Short: "sub-commands relating to domain ACL GeoIP rules mode",
	}

	// Child commands
	cmd.AddCommand(ddosxDomainACLGeoIPRulesModeShowCmd())
	cmd.AddCommand(ddosxDomainACLGeoIPRulesModeUpdateCmd())

	return cmd
}

func ddosxDomainACLGeoIPRulesModeShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "show <domain: name>...",
		Short:   "Shows a domain",
		Long:    "This command shows one or more domains",
		Example: "ukfast ddosx domain show example.com",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing domain")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			ddosxDomainACLGeoIPRulesModeShow(getClient().DDoSXService(), cmd, args)
		},
	}
}

func ddosxDomainACLGeoIPRulesModeShow(service ddosx.DDoSXService, cmd *cobra.Command, args []string) {
	var modes []ddosx.ACLGeoIPRulesMode
	for _, arg := range args {
		mode, err := service.GetDomainACLGeoIPRulesMode(arg)
		if err != nil {
			OutputWithErrorLevelf("Error retrieving domain [%s] ACL GeoIP rules mode: %s", arg, err)
			continue
		}

		modes = append(modes, mode)
	}

	outputDDoSXACLGeoIPRulesModes(modes)
}

func ddosxDomainACLGeoIPRulesModeUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <domain: name>",
		Short:   "Updates a domain ACL GeoIP rule filtering mode",
		Long:    "This command updates a domain ACL GeoIP rule filtering mode",
		Example: "ukfast ddosx domain acl geoip mode update example.com --mode whitelist",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing domain")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			ddosxDomainACLGeoIPRulesModeUpdate(getClient().DDoSXService(), cmd, args)
		},
	}

	cmd.Flags().String("mode", "", "Filtering mode for GeoIP ACL rules")

	return cmd
}

func ddosxDomainACLGeoIPRulesModeUpdate(service ddosx.DDoSXService, cmd *cobra.Command, args []string) {
	patchRequest := ddosx.PatchACLGeoIPRulesModeRequest{}

	if cmd.Flags().Changed("mode") {
		mode, _ := cmd.Flags().GetString("mode")
		parsedMode, err := ddosx.ParseACLGeoIPRulesMode(mode)
		if err != nil {
			output.Fatal(err.Error())
			return
		}
		patchRequest.Mode = parsedMode
	}

	err := service.PatchDomainACLGeoIPRulesMode(args[0], patchRequest)
	if err != nil {
		output.Fatalf("Error updating domain ACL GeoIP rule filtering mode: %s", err.Error())
		return
	}

	mode, err := service.GetDomainACLGeoIPRulesMode(args[0])
	if err != nil {
		output.Fatalf("Error retrieving updated domain ACL GeoIP rule filtering mode: %s", err)
		return
	}

	outputDDoSXACLGeoIPRulesModes([]ddosx.ACLGeoIPRulesMode{mode})
}
