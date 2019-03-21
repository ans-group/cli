package cmd

import (
	"errors"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/helper"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/ddosx"
)

func ddosxDomainPropertyRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "property",
		Short: "sub-commands relating to domain properties",
	}

	// Child commands
	cmd.AddCommand(ddosxDomainPropertyListCmd())
	cmd.AddCommand(ddosxDomainPropertyShowCmd())
	cmd.AddCommand(ddosxDomainPropertyUpdateCmd())

	return cmd
}

func ddosxDomainPropertyListCmd() *cobra.Command {
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
		Run: func(cmd *cobra.Command, args []string) {
			ddosxDomainPropertyList(getClient().DDoSXService(), cmd, args)
		},
	}

	cmd.Flags().String("name", "", "Property name for filtering")

	return cmd
}

func ddosxDomainPropertyList(service ddosx.DDoSXService, cmd *cobra.Command, args []string) {
	params, err := GetAPIRequestParametersFromFlags()
	if err != nil {
		output.Fatal(err.Error())
		return
	}

	if cmd.Flags().Changed("name") {
		filterName, _ := cmd.Flags().GetString("name")
		params.WithFilter(helper.GetFilteringInferOperator("name", filterName))
	}

	properties, err := service.GetDomainProperties(args[0], params)
	if err != nil {
		output.Fatalf("Error retrieving domain properties: %s", err)
		return
	}

	outputDDoSXDomainProperties(properties)
}

func ddosxDomainPropertyShowCmd() *cobra.Command {
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
		Run: func(cmd *cobra.Command, args []string) {
			ddosxDomainPropertyShow(getClient().DDoSXService(), cmd, args)
		},
	}
}

func ddosxDomainPropertyShow(service ddosx.DDoSXService, cmd *cobra.Command, args []string) {

	var properties []ddosx.DomainProperty

	for _, arg := range args[1:] {
		property, err := service.GetDomainProperty(args[0], arg)
		if err != nil {
			OutputWithErrorLevelf("Error retrieving domain property [%s]: %s", arg, err.Error())
			continue
		}

		properties = append(properties, property)
	}

	outputDDoSXDomainProperties(properties)
}

func ddosxDomainPropertyUpdateCmd() *cobra.Command {
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
		Run: func(cmd *cobra.Command, args []string) {
			ddosxDomainPropertyUpdate(getClient().DDoSXService(), cmd, args)
		},
	}

	cmd.Flags().String("value", "", "Property value")

	return cmd
}

func ddosxDomainPropertyUpdate(service ddosx.DDoSXService, cmd *cobra.Command, args []string) {
	var properties []ddosx.DomainProperty

	updateRequest := ddosx.PatchDomainPropertyRequest{}

	if cmd.Flags().Changed("value") {
		value, _ := cmd.Flags().GetString("value")
		updateRequest.Value = helper.InferTypeFromStringFlag(value)
	}

	for _, arg := range args[1:] {
		err := service.PatchDomainProperty(args[0], arg, updateRequest)
		if err != nil {
			OutputWithErrorLevelf("Error updating domain property [%s]: %s", arg, err.Error())
			continue
		}

		property, err := service.GetDomainProperty(args[0], arg)
		if err != nil {
			OutputWithErrorLevelf("Error retrieving updated domain property [%s]: %s", arg, err.Error())
			continue
		}

		properties = append(properties, property)
	}

	outputDDoSXDomainProperties(properties)
}
