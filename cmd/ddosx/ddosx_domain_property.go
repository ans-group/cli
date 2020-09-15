package ddosx

import (
	"errors"
	"fmt"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
	"github.com/ukfast/cli/internal/pkg/helper"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/ddosx"
)

func ddosxDomainPropertyRootCmd(f factory.ClientFactory, fs afero.Fs) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "property",
		Short: "sub-commands relating to domain properties",
	}

	// Child commands
	cmd.AddCommand(ddosxDomainPropertyListCmd(f))
	cmd.AddCommand(ddosxDomainPropertyShowCmd(f))
	cmd.AddCommand(ddosxDomainPropertyUpdateCmd(f, fs))

	return cmd
}

func ddosxDomainPropertyListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list <domain: name>",
		Short:   "Lists domain properties",
		Long:    "This command lists domain properties",
		Example: "ukfast ddosx domain property list example.com",
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

			return ddosxDomainPropertyList(c.DDoSXService(), cmd, args)
		},
	}

	cmd.Flags().String("name", "", "Property name for filtering")

	return cmd
}

func ddosxDomainPropertyList(service ddosx.DDoSXService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	helper.HydrateAPIRequestParametersWithStringFilterFlag(&params, cmd, helper.NewStringFilterFlag("name", "name"))

	properties, err := service.GetDomainProperties(args[0], params)
	if err != nil {
		return fmt.Errorf("Error retrieving domain properties: %s", err)
	}

	return output.CommandOutput(cmd, OutputDDoSXDomainPropertiesProvider(properties))
}

func ddosxDomainPropertyShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <domain: name>",
		Short:   "Shows domain properties",
		Long:    "This command shows a domain property",
		Example: "ukfast ddosx domain property show example.com 00000000-0000-0000-0000-000000000000",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing domain")
			}
			if len(args) < 2 {
				return errors.New("Missing domain property")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ddosxDomainPropertyShow(c.DDoSXService(), cmd, args)
		},
	}
}

func ddosxDomainPropertyShow(service ddosx.DDoSXService, cmd *cobra.Command, args []string) error {

	var properties []ddosx.DomainProperty

	for _, arg := range args[1:] {
		property, err := service.GetDomainProperty(args[0], arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving domain property [%s]: %s", arg, err.Error())
			continue
		}

		properties = append(properties, property)
	}

	return output.CommandOutput(cmd, OutputDDoSXDomainPropertiesProvider(properties))
}

func ddosxDomainPropertyUpdateCmd(f factory.ClientFactory, fs afero.Fs) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <domain: name>...",
		Short:   "Updates domain properties",
		Long:    "This command updates one or more domain properties",
		Example: "ukfast ddosx domain property update example.com 00000000-0000-0000-0000-000000000000 --value false",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing domain")
			}
			if len(args) < 2 {
				return errors.New("Missing domain property")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ddosxDomainPropertyUpdate(c.DDoSXService(), cmd, fs, args)
		},
	}

	cmd.Flags().String("value", "", "Property value")
	cmd.Flags().String("value-file", "", "Property value (from file")

	return cmd
}

func ddosxDomainPropertyUpdate(service ddosx.DDoSXService, cmd *cobra.Command, fs afero.Fs, args []string) error {
	var properties []ddosx.DomainProperty

	updateRequest := ddosx.PatchDomainPropertyRequest{}

	if cmd.Flags().Changed("value") {
		value, _ := cmd.Flags().GetString("value")
		updateRequest.Value = helper.InferTypeFromStringFlagValue(value)
	} else if cmd.Flags().Changed("value-file") {
		var err error
		updateRequest.Value, err = helper.GetContentsFromFilePathFlag(cmd, fs, "value-file")
		if err != nil {
			return err
		}
	}

	for _, arg := range args[1:] {
		err := service.PatchDomainProperty(args[0], arg, updateRequest)
		if err != nil {
			output.OutputWithErrorLevelf("Error updating domain property [%s]: %s", arg, err.Error())
			continue
		}

		property, err := service.GetDomainProperty(args[0], arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving updated domain property [%s]: %s", arg, err.Error())
			continue
		}

		properties = append(properties, property)
	}

	return output.CommandOutput(cmd, OutputDDoSXDomainPropertiesProvider(properties))
}
