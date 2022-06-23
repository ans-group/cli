package account

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/helper"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/account"
	"github.com/spf13/cobra"
)

func accountContactRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "contact",
		Short: "sub-commands relating to contacts",
	}

	// Child commands
	cmd.AddCommand(accountContactListCmd(f))
	cmd.AddCommand(accountContactShowCmd(f))

	return cmd
}

func accountContactListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists contacts",
		Long:    "This command lists contacts",
		Example: "ukfast account contact list",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return accountContactList(c.AccountService(), cmd, args)
		},
	}
}

func accountContactList(service account.AccountService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	contacts, err := service.GetContacts(params)
	if err != nil {
		return fmt.Errorf("Error retrieving contacts: %s", err)
	}

	return output.CommandOutput(cmd, OutputAccountContactsProvider(contacts))
}

func accountContactShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <contact: id>...",
		Short:   "Shows a contact",
		Long:    "This command shows one or more contacts",
		Example: "ukfast account contact show 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing contact")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return accountContactShow(c.AccountService(), cmd, args)
		},
	}
}

func accountContactShow(service account.AccountService, cmd *cobra.Command, args []string) error {
	var contacts []account.Contact
	for _, arg := range args {
		contactID, err := strconv.Atoi(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Invalid contact ID [%s]", arg)
			continue
		}

		contact, err := service.GetContact(contactID)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving contact [%s]: %s", arg, err)
			continue
		}

		contacts = append(contacts, contact)
	}

	return output.CommandOutput(cmd, OutputAccountContactsProvider(contacts))
}
