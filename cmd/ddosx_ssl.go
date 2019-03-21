package cmd

import (
	"errors"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/ddosx"
)

func ddosxSSLRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ssl",
		Short: "sub-commands relating to ssls",
	}

	// Child commands
	cmd.AddCommand(ddosxSSLListCmd())
	cmd.AddCommand(ddosxSSLShowCmd())
	cmd.AddCommand(ddosxSSLCreateCmd())
	cmd.AddCommand(ddosxSSLUpdateCmd())
	cmd.AddCommand(ddosxSSLDeleteCmd())

	// Child root rommands
	cmd.AddCommand(ddosxSSLContentRootCmd())
	cmd.AddCommand(ddosxSSLPrivateKeyRootCmd())

	return cmd
}

func ddosxSSLListCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists ssls",
		Long:    "This command lists ssls",
		Example: "ukfast ddosx ssl list",
		Run: func(cmd *cobra.Command, args []string) {
			ddosxSSLList(getClient().DDoSXService(), cmd, args)
		},
	}
}

func ddosxSSLList(service ddosx.DDoSXService, cmd *cobra.Command, args []string) {
	params, err := GetAPIRequestParametersFromFlags()
	if err != nil {
		output.Fatal(err.Error())
		return
	}

	ssls, err := service.GetSSLs(params)
	if err != nil {
		output.Fatalf("Error retrieving ssls: %s", err)
		return
	}

	outputDDoSXSSLs(ssls)
}

func ddosxSSLShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "show <ssl: id>...",
		Short:   "Shows a ssl",
		Long:    "This command shows one or more ssls",
		Example: "ukfast ddosx ssl show 00000000-0000-0000-0000-000000000000",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing ssl")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			ddosxSSLShow(getClient().DDoSXService(), cmd, args)
		},
	}
}

func ddosxSSLShow(service ddosx.DDoSXService, cmd *cobra.Command, args []string) {
	var ssls []ddosx.SSL
	for _, arg := range args {
		ssl, err := service.GetSSL(arg)
		if err != nil {
			OutputWithErrorLevelf("Error retrieving ssl [%s]: %s", arg, err)
			continue
		}

		ssls = append(ssls, ssl)
	}

	outputDDoSXSSLs(ssls)
}

func ddosxSSLCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Creates an ssl",
		Long:    "This command creates an SSL",
		Example: "ukfast ddosx ssl create",
		Run: func(cmd *cobra.Command, args []string) {
			ddosxSSLCreate(getClient().DDoSXService(), cmd, args)
		},
	}

	cmd.Flags().String("friendly-name", "", "Friendly name for SSL")
	cmd.MarkFlagRequired("friendly-name")
	cmd.Flags().Int("ukfast-ssl-id", 0, "Optional ID of UKFast SSL to retrieve certificate, key and bundle")
	cmd.Flags().String("key", "", "Key for SSL")
	cmd.Flags().String("certificate", "", "Certificate contents for SSL")
	cmd.Flags().String("ca-bundle", "", "CA bundle contents for SSL")

	return cmd
}

func ddosxSSLCreate(service ddosx.DDoSXService, cmd *cobra.Command, args []string) {
	createRequest := ddosx.CreateSSLRequest{}
	createRequest.FriendlyName, _ = cmd.Flags().GetString("friendly-name")

	if cmd.Flags().Changed("ukfast-ssl-id") {
		createRequest.UKFastSSLID, _ = cmd.Flags().GetInt("ukfast-ssl-id")
	} else {
		createRequest.Key, _ = cmd.Flags().GetString("key")
		createRequest.Certificate, _ = cmd.Flags().GetString("certificate")
		createRequest.CABundle, _ = cmd.Flags().GetString("ca-bundle")
	}

	id, err := service.CreateSSL(createRequest)
	if err != nil {
		output.Fatalf("Error creating ssl: %s", err.Error())
		return
	}

	ssl, err := service.GetSSL(id)
	if err != nil {
		output.Fatalf("Error retrieving new ssl [%s]: %s", id, err.Error())
		return
	}

	outputDDoSXSSLs([]ddosx.SSL{ssl})
}

func ddosxSSLUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <ssl: id>",
		Short:   "Updates an ssl",
		Long:    "This command updates an SSL",
		Example: "ukfast ddosx ssl update 00000000-0000-0000-0000-000000000000 --friendly-name myssl",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing ssl")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			ddosxSSLUpdate(getClient().DDoSXService(), cmd, args)
		},
	}

	cmd.Flags().String("friendly-name", "", "Friendly name for SSL")
	cmd.Flags().Int("ukfast-ssl-id", 0, "Optional ID of UKFast SSL to retrieve certificate, key and bundle")
	cmd.Flags().String("key", "", "Key for SSL")
	cmd.Flags().String("certificate", "", "Certificate contents for SSL")
	cmd.Flags().String("ca-bundle", "", "CA bundle contents for SSL")

	return cmd
}

func ddosxSSLUpdate(service ddosx.DDoSXService, cmd *cobra.Command, args []string) {
	patchRequest := ddosx.PatchSSLRequest{}
	patchRequest.FriendlyName, _ = cmd.Flags().GetString("friendly-name")

	if cmd.Flags().Changed("ukfast-ssl-id") {
		patchRequest.UKFastSSLID, _ = cmd.Flags().GetInt("ukfast-ssl-id")
	} else {
		if cmd.Flags().Changed("key") {
			patchRequest.Key, _ = cmd.Flags().GetString("key")
		}
		if cmd.Flags().Changed("certificate") {
			patchRequest.Certificate, _ = cmd.Flags().GetString("certificate")
		}
		if cmd.Flags().Changed("ca-bundle") {
			patchRequest.CABundle, _ = cmd.Flags().GetString("ca-bundle")
		}
	}

	_, err := service.PatchSSL(args[0], patchRequest)
	if err != nil {
		output.Fatalf("Error updating ssl: %s", err.Error())
		return
	}

	ssl, err := service.GetSSL(args[0])
	if err != nil {
		output.Fatalf("Error retrieving updated ssl: %s", err.Error())
		return
	}

	outputDDoSXSSLs([]ddosx.SSL{ssl})
}

func ddosxSSLDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "delete <ssl: id>...",
		Short:   "Deletes a ssl",
		Long:    "This command deletes one or more ssls",
		Example: "ukfast ddosx ssl delete 00000000-0000-0000-0000-000000000000",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing ssl")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			ddosxSSLDelete(getClient().DDoSXService(), cmd, args)
		},
	}
}

func ddosxSSLDelete(service ddosx.DDoSXService, cmd *cobra.Command, args []string) {
	for _, arg := range args {
		err := service.DeleteSSL(arg)
		if err != nil {
			OutputWithErrorLevelf("Error removing ssl [%s]: %s", arg, err)
			continue
		}
	}
}
