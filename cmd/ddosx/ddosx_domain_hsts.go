package ddosx

import (
	"errors"

	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/ddosx"
	"github.com/spf13/cobra"
)

func ddosxDomainHSTSRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hsts",
		Short: "sub-commands relating to domain HSTS",
	}

	// Child commands
	cmd.AddCommand(ddosxDomainHSTSShowCmd(f))
	cmd.AddCommand(ddosxDomainHSTSEnableCmd(f))
	cmd.AddCommand(ddosxDomainHSTSDisableCmd(f))

	// Child root commands
	cmd.AddCommand(ddosxDomainHSTSRuleRootCmd(f))

	return cmd
}

func ddosxDomainHSTSShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <domain: name>...",
		Short:   "Shows HSTS for a domain",
		Long:    "This command shows HSTS for one or more domains",
		Example: "ukfast ddosx domain hsts show example.com",
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

			return ddosxDomainHSTSShow(c.DDoSXService(), cmd, args)
		},
	}
}

func ddosxDomainHSTSShow(service ddosx.DDoSXService, cmd *cobra.Command, args []string) error {
	var configurations []ddosx.HSTSConfiguration

	for _, arg := range args {
		configuration, err := service.GetDomainHSTSConfiguration(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving HSTS configuration for domain [%s]: %s", arg, err)
			continue
		}

		configurations = append(configurations, configuration)
	}

	return output.CommandOutput(cmd, OutputDDoSXHSTSConfigurationsProvider(configurations))
}

func ddosxDomainHSTSEnableCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "enable <domain: name>...",
		Short:   "Enables HSTS for a domain",
		Long:    "This command enables HSTS for one or more domains",
		Example: "ukfast ddosx domain hsts enable example.com",
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

			return ddosxDomainHSTSEnable(c.DDoSXService(), cmd, args)
		},
	}
}

func ddosxDomainHSTSEnable(service ddosx.DDoSXService, cmd *cobra.Command, args []string) error {
	var configurations []ddosx.HSTSConfiguration

	for _, arg := range args {
		err := service.AddDomainHSTSConfiguration(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error enabling HSTS for domain [%s]: %s", arg, err.Error())
			continue
		}

		configuration, err := service.GetDomainHSTSConfiguration(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving updated HSTS configuration for domain [%s]: %s", arg, err)
			continue
		}

		configurations = append(configurations, configuration)
	}

	return output.CommandOutput(cmd, OutputDDoSXHSTSConfigurationsProvider(configurations))
}

func ddosxDomainHSTSDisableCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "disable <domain: name>...",
		Short:   "Disables HSTS for a domain",
		Long:    "This command disables HSTS for one or more domains",
		Example: "ukfast ddosx domain hsts disable example.com",
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

			return ddosxDomainHSTSDisable(c.DDoSXService(), cmd, args)
		},
	}
}

func ddosxDomainHSTSDisable(service ddosx.DDoSXService, cmd *cobra.Command, args []string) error {
	var configurations []ddosx.HSTSConfiguration

	for _, arg := range args {
		err := service.DeleteDomainHSTSConfiguration(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error disabling HSTS for domain [%s]: %s", arg, err.Error())
			continue
		}

		configuration, err := service.GetDomainHSTSConfiguration(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving updated HSTS configuration for domain [%s]: %s", arg, err)
			continue
		}

		configurations = append(configurations, configuration)
	}

	return output.CommandOutput(cmd, OutputDDoSXHSTSConfigurationsProvider(configurations))
}
