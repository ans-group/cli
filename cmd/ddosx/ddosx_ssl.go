package ddosx

import (
	"errors"
	"fmt"

	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/helper"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/ddosx"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func ddosxSSLRootCmd(f factory.ClientFactory, fs afero.Fs) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ssl",
		Short: "sub-commands relating to ssls",
	}

	// Child commands
	cmd.AddCommand(ddosxSSLListCmd(f))
	cmd.AddCommand(ddosxSSLShowCmd(f))
	cmd.AddCommand(ddosxSSLCreateCmd(f, fs))
	cmd.AddCommand(ddosxSSLUpdateCmd(f, fs))
	cmd.AddCommand(ddosxSSLDeleteCmd(f))

	// Child root rommands
	cmd.AddCommand(ddosxSSLContentRootCmd(f))
	cmd.AddCommand(ddosxSSLPrivateKeyRootCmd(f))

	return cmd
}

func ddosxSSLListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists ssls",
		Long:    "This command lists ssls",
		Example: "ans ddosx ssl list",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ddosxSSLList(c.DDoSXService(), cmd, args)
		},
	}
}

func ddosxSSLList(service ddosx.DDoSXService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	ssls, err := service.GetSSLs(params)
	if err != nil {
		return fmt.Errorf("Error retrieving ssls: %s", err)
	}

	return output.CommandOutput(cmd, SSLCollection(ssls))
}

func ddosxSSLShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <ssl: id>...",
		Short:   "Shows a ssl",
		Long:    "This command shows one or more ssls",
		Example: "ans ddosx ssl show 00000000-0000-0000-0000-000000000000",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing ssl")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ddosxSSLShow(c.DDoSXService(), cmd, args)
		},
	}
}

func ddosxSSLShow(service ddosx.DDoSXService, cmd *cobra.Command, args []string) error {
	var ssls []ddosx.SSL
	for _, arg := range args {
		ssl, err := service.GetSSL(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving ssl [%s]: %s", arg, err)
			continue
		}

		ssls = append(ssls, ssl)
	}

	return output.CommandOutput(cmd, SSLCollection(ssls))
}

func ddosxSSLCreateCmd(f factory.ClientFactory, fs afero.Fs) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Creates an ssl",
		Long:    "This command creates an SSL",
		Example: "ans ddosx ssl create",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ddosxSSLCreate(c.DDoSXService(), cmd, fs, args)
		},
	}

	cmd.Flags().String("friendly-name", "", "Friendly name for SSL")
	cmd.MarkFlagRequired("friendly-name")
	cmd.Flags().Int("ans-ssl-id", 0, "Optional ID of ANS SSL to retrieve certificate, key and bundle")
	cmd.Flags().String("key", "", "Key contents for SSL")
	cmd.Flags().String("key-file", "", "Path to file containing key contents for SSL")
	cmd.Flags().String("certificate", "", "Certificate contents for SSL")
	cmd.Flags().String("certificate-file", "", "Path to file containing certificate contents for SSL")
	cmd.Flags().String("ca-bundle", "", "CA bundle contents for SSL")
	cmd.Flags().String("ca-bundle-file", "", "Path to file containing CA bundle contents for SSL")

	return cmd
}

func ddosxSSLCreate(service ddosx.DDoSXService, cmd *cobra.Command, fs afero.Fs, args []string) error {
	createRequest := ddosx.CreateSSLRequest{}
	createRequest.FriendlyName, _ = cmd.Flags().GetString("friendly-name")

	if cmd.Flags().Changed("ans-ssl-id") {
		createRequest.UKFastSSLID, _ = cmd.Flags().GetInt("ans-ssl-id")
	} else {
		var err error
		createRequest.Key, err = helper.GetContentsFromLiteralOrFilePathFlag(cmd, fs, "key", "key-file")
		if err != nil {
			return err
		}
		createRequest.Certificate, err = helper.GetContentsFromLiteralOrFilePathFlag(cmd, fs, "certificate", "certificate-file")
		if err != nil {
			return err
		}
		createRequest.CABundle, err = helper.GetContentsFromLiteralOrFilePathFlag(cmd, fs, "ca-bundle", "ca-bundle-file")
		if err != nil {
			return err
		}
	}

	id, err := service.CreateSSL(createRequest)
	if err != nil {
		return fmt.Errorf("Error creating ssl: %s", err.Error())
	}

	ssl, err := service.GetSSL(id)
	if err != nil {
		return fmt.Errorf("Error retrieving new ssl [%s]: %s", id, err.Error())
	}

	return output.CommandOutput(cmd, SSLCollection([]ddosx.SSL{ssl}))
}

func ddosxSSLUpdateCmd(f factory.ClientFactory, fs afero.Fs) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <ssl: id>",
		Short:   "Updates an ssl",
		Long:    "This command updates an SSL",
		Example: "ans ddosx ssl update 00000000-0000-0000-0000-000000000000 --friendly-name myssl",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing ssl")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ddosxSSLUpdate(c.DDoSXService(), cmd, fs, args)
		},
	}

	cmd.Flags().String("friendly-name", "", "Friendly name for SSL")
	cmd.Flags().Int("ans-ssl-id", 0, "Optional ID of ANS SSL to retrieve certificate, key and bundle")
	cmd.Flags().String("key", "", "Key contents for SSL")
	cmd.Flags().String("key-file", "", "Path to file containing key contents for SSL")
	cmd.Flags().String("certificate", "", "Certificate contents for SSL")
	cmd.Flags().String("certificate-file", "", "Path to file containing certificate contents for SSL")
	cmd.Flags().String("ca-bundle", "", "CA bundle contents for SSL")
	cmd.Flags().String("ca-bundle-file", "", "Path to file containing CA bundle contents for SSL")

	return cmd
}

func ddosxSSLUpdate(service ddosx.DDoSXService, cmd *cobra.Command, fs afero.Fs, args []string) error {
	patchRequest := ddosx.PatchSSLRequest{}
	patchRequest.FriendlyName, _ = cmd.Flags().GetString("friendly-name")

	if cmd.Flags().Changed("ans-ssl-id") {
		patchRequest.UKFastSSLID, _ = cmd.Flags().GetInt("ans-ssl-id")
	} else {
		var err error
		patchRequest.Key, err = helper.GetContentsFromLiteralOrFilePathFlag(cmd, fs, "key", "key-file")
		if err != nil {
			return err
		}

		patchRequest.Certificate, err = helper.GetContentsFromLiteralOrFilePathFlag(cmd, fs, "certificate", "certificate-file")
		if err != nil {
			return err
		}

		patchRequest.CABundle, err = helper.GetContentsFromLiteralOrFilePathFlag(cmd, fs, "ca-bundle", "ca-bundle-file")
		if err != nil {
			return err
		}
	}

	_, err := service.PatchSSL(args[0], patchRequest)
	if err != nil {
		return fmt.Errorf("Error updating ssl: %s", err.Error())
	}

	ssl, err := service.GetSSL(args[0])
	if err != nil {
		return fmt.Errorf("Error retrieving updated ssl: %s", err.Error())
	}

	return output.CommandOutput(cmd, SSLCollection([]ddosx.SSL{ssl}))
}

func ddosxSSLDeleteCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "delete <ssl: id>...",
		Short:   "Deletes a ssl",
		Long:    "This command deletes one or more ssls",
		Example: "ans ddosx ssl delete 00000000-0000-0000-0000-000000000000",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing ssl")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			ddosxSSLDelete(c.DDoSXService(), cmd, args)
			return nil
		},
	}
}

func ddosxSSLDelete(service ddosx.DDoSXService, cmd *cobra.Command, args []string) {
	for _, arg := range args {
		err := service.DeleteSSL(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error removing ssl [%s]: %s", arg, err)
			continue
		}
	}
}
