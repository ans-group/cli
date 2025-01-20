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

func accountClientRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "client",
		Short: "sub-commands relating to clients",
	}

	// Child commands
	cmd.AddCommand(accountClientListCmd(f))
	cmd.AddCommand(accountClientShowCmd(f))
	cmd.AddCommand(accountClientCreateCmd(f))
	cmd.AddCommand(accountClientUpdateCmd(f))
	cmd.AddCommand(accountClientDeleteCmd(f))

	return cmd
}

func accountClientListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists clients",
		Long:    "This command lists clients",
		Example: "ans account client list",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return accountClientList(c.AccountService(), cmd, args)
		},
	}
}

func accountClientList(service account.AccountService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	clients, err := service.GetClients(params)
	if err != nil {
		return fmt.Errorf("Error retrieving clients: %s", err)
	}

	return output.CommandOutput(cmd, ClientCollection(clients))
}

func accountClientShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <client: id>...",
		Short:   "Shows a client",
		Long:    "This command shows one or more clients",
		Example: "ans account client show 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing client")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return accountClientShow(c.AccountService(), cmd, args)
		},
	}
}

func accountClientShow(service account.AccountService, cmd *cobra.Command, args []string) error {
	var clients []account.Client
	for _, arg := range args {
		clientID, err := strconv.Atoi(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Invalid client ID [%s]", arg)
			continue
		}

		client, err := service.GetClient(clientID)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving client [%s]: %s", arg, err)
			continue
		}

		clients = append(clients, client)
	}

	return output.CommandOutput(cmd, ClientCollection(clients))
}

func accountClientCreateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Creates a client",
		Long:    "This command creates a client",
		Example: "ans account client create",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return accountClientCreate(c.AccountService(), cmd, args)
		},
	}

	// Setup flags
	cmd.Flags().String("company-name", "", "Name of company")
	cmd.Flags().String("first-name", "", "First name of client contact")
	cmd.Flags().String("last-name", "", "Last name of client contact")
	cmd.Flags().String("email-address", "", "Email address of client contact")
	cmd.Flags().String("limited-number", "", "Limited company number")
	cmd.Flags().String("vat-number", "", "VAT number")
	cmd.Flags().String("address", "", "Address of company")
	cmd.Flags().String("address1", "", "Address line 1 of company")
	cmd.Flags().String("city", "", "City of company")
	cmd.Flags().String("county", "", "County of company")
	cmd.Flags().String("country", "", "Country of company")
	cmd.Flags().String("postcode", "", "Post code of company")
	cmd.Flags().String("phone-number", "", "Phone number of company")
	cmd.Flags().String("fax-number", "", "Fax number of company")
	cmd.Flags().String("mobile-number", "", "Mobile number of company")
	cmd.Flags().String("type", "", "Type of client")
	cmd.Flags().String("user-name", "", "Username of client contact")
	cmd.Flags().String("id-reference", "", "ID reference for client")
	cmd.Flags().String("nominet-contact-id", "", "ID of Nominet contact")

	return cmd
}

func accountClientCreate(service account.AccountService, cmd *cobra.Command, args []string) error {
	createRequest := account.CreateClientRequest{}
	createRequest.CompanyName, _ = cmd.Flags().GetString("company-name")
	createRequest.FirstName, _ = cmd.Flags().GetString("first-name")
	createRequest.LastName, _ = cmd.Flags().GetString("last-name")
	createRequest.EmailAddress, _ = cmd.Flags().GetString("email-address")
	createRequest.LimitedNumber, _ = cmd.Flags().GetString("limited-number")
	createRequest.VATNumber, _ = cmd.Flags().GetString("vat-number")
	createRequest.Address, _ = cmd.Flags().GetString("address")
	createRequest.Address1, _ = cmd.Flags().GetString("address1")
	createRequest.City, _ = cmd.Flags().GetString("city")
	createRequest.County, _ = cmd.Flags().GetString("county")
	createRequest.Country, _ = cmd.Flags().GetString("country")
	createRequest.Postcode, _ = cmd.Flags().GetString("postcode")
	createRequest.Phone, _ = cmd.Flags().GetString("phone-number")
	createRequest.Fax, _ = cmd.Flags().GetString("fax-number")
	createRequest.Mobile, _ = cmd.Flags().GetString("mobile-number")
	createRequest.Type, _ = cmd.Flags().GetString("type")
	createRequest.UserName, _ = cmd.Flags().GetString("user-name")
	createRequest.IDReference, _ = cmd.Flags().GetString("id-reference")
	createRequest.NominetContactID, _ = cmd.Flags().GetString("nominet-contact-id")

	id, err := service.CreateClient(createRequest)
	if err != nil {
		return fmt.Errorf("Error creating client: %s", err)
	}

	client, err := service.GetClient(id)
	if err != nil {
		return fmt.Errorf("Error retrieving new client: %s", err)
	}

	return output.CommandOutput(cmd, ClientCollection([]account.Client{client}))
}

func accountClientUpdateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <virtualmachine: id>...",
		Short:   "Updates a client",
		Long:    "This command updates one or more clients",
		Example: "ans account client update 123 --name \"test client 1\"",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing client")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return accountClientUpdate(c.AccountService(), cmd, args)
		},
	}

	// Setup flags
	cmd.Flags().String("company-name", "", "Name of company")
	cmd.Flags().String("first-name", "", "First name of client contact")
	cmd.Flags().String("last-name", "", "Last name of client contact")
	cmd.Flags().String("email-address", "", "Email address of client contact")
	cmd.Flags().String("limited-number", "", "Limited company number")
	cmd.Flags().String("vat-number", "", "VAT number")
	cmd.Flags().String("address", "", "Address of company")
	cmd.Flags().String("address1", "", "Address line 1 of company")
	cmd.Flags().String("city", "", "City of company")
	cmd.Flags().String("county", "", "County of company")
	cmd.Flags().String("country", "", "Country of company")
	cmd.Flags().String("postcode", "", "Post code of company")
	cmd.Flags().String("phone-number", "", "Phone number of company")
	cmd.Flags().String("fax-number", "", "Fax number of company")
	cmd.Flags().String("mobile-number", "", "Mobile number of company")
	cmd.Flags().String("type", "", "Type of client")
	cmd.Flags().String("user-name", "", "Username of client contact")
	cmd.Flags().String("id-reference", "", "ID reference for client")
	cmd.Flags().String("nominet-contact-id", "", "ID of Nominet contact")

	return cmd
}

func accountClientUpdate(service account.AccountService, cmd *cobra.Command, args []string) error {

	patchRequest := account.PatchClientRequest{}
	patchRequest.CompanyName, _ = cmd.Flags().GetString("company-name")
	patchRequest.FirstName, _ = cmd.Flags().GetString("first-name")
	patchRequest.LastName, _ = cmd.Flags().GetString("last-name")
	patchRequest.EmailAddress, _ = cmd.Flags().GetString("email-address")
	patchRequest.LimitedNumber, _ = cmd.Flags().GetString("limited-number")
	patchRequest.VATNumber, _ = cmd.Flags().GetString("vat-number")
	patchRequest.Address, _ = cmd.Flags().GetString("address")
	patchRequest.Address1, _ = cmd.Flags().GetString("address1")
	patchRequest.City, _ = cmd.Flags().GetString("city")
	patchRequest.County, _ = cmd.Flags().GetString("county")
	patchRequest.Country, _ = cmd.Flags().GetString("country")
	patchRequest.Postcode, _ = cmd.Flags().GetString("postcode")
	patchRequest.Phone, _ = cmd.Flags().GetString("phone-number")
	patchRequest.Fax, _ = cmd.Flags().GetString("fax-number")
	patchRequest.Mobile, _ = cmd.Flags().GetString("mobile-number")
	patchRequest.Type, _ = cmd.Flags().GetString("type")
	patchRequest.UserName, _ = cmd.Flags().GetString("user-name")
	patchRequest.IDReference, _ = cmd.Flags().GetString("id-reference")
	patchRequest.NominetContactID, _ = cmd.Flags().GetString("nominet-contact-id")

	var clients []account.Client

	for _, arg := range args {
		clientID, err := strconv.Atoi(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Invalid client ID [%s]", arg)
			continue
		}

		err = service.PatchClient(clientID, patchRequest)
		if err != nil {
			output.OutputWithErrorLevelf("Error updating client [%d]: %s", clientID, err.Error())
			continue
		}

		client, err := service.GetClient(clientID)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving updated client [%d]: %s", clientID, err.Error())
			continue
		}

		clients = append(clients, client)
	}

	return output.CommandOutput(cmd, ClientCollection(clients))
}

func accountClientDeleteCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "delete <client: id>...",
		Short:   "Removes a client",
		Long:    "This command removes one or more clients",
		Example: "ans account client delete 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing client")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			accountClientDelete(c.AccountService(), cmd, args)
			return nil
		},
	}
}

func accountClientDelete(service account.AccountService, cmd *cobra.Command, args []string) {
	for _, arg := range args {
		clientID, err := strconv.Atoi(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Invalid client ID [%s]", arg)
			continue
		}

		err = service.DeleteClient(clientID)
		if err != nil {
			output.OutputWithErrorLevelf("Error removing client [%d]: %s", clientID, err)
			continue
		}
	}
}
