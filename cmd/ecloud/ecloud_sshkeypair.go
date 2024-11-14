package ecloud

import (
	"errors"
	"fmt"

	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/helper"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func ecloudSSHKeyPairRootCmd(f factory.ClientFactory, fs afero.Fs) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sshkeypair",
		Short: "sub-commands relating to SSH key pairs",
	}

	// Child commands
	cmd.AddCommand(ecloudSSHKeyPairListCmd(f))
	cmd.AddCommand(ecloudSSHKeyPairShowCmd(f))
	cmd.AddCommand(ecloudSSHKeyPairCreateCmd(f, fs))
	cmd.AddCommand(ecloudSSHKeyPairUpdateCmd(f, fs))
	cmd.AddCommand(ecloudSSHKeyPairDeleteCmd(f))

	return cmd
}

func ecloudSSHKeyPairListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Lists SSH key pairs",
		Long:    "This command lists SSH key pairs",
		Example: "ans ecloud sshkeypair list",
		RunE:    ecloudCobraRunEFunc(f, ecloudSSHKeyPairList),
	}

	cmd.Flags().String("name", "", "SSH key pair name for filtering")

	return cmd
}

func ecloudSSHKeyPairList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd,
		helper.NewStringFilterFlagOption("name", "name"),
	)
	if err != nil {
		return err
	}

	keypairs, err := service.GetSSHKeyPairs(params)
	if err != nil {
		return fmt.Errorf("Error retrieving SSH key pairs: %s", err)
	}

	return output.CommandOutput(cmd, SSHKeyPairCollection(keypairs))
}

func ecloudSSHKeyPairShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <keypair: id>...",
		Short:   "Shows a SSH key pair",
		Long:    "This command shows one or more SSH key pairs",
		Example: "ans ecloud sshkeypair show ssh-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing SSH key pair")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudSSHKeyPairShow),
	}
}

func ecloudSSHKeyPairShow(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	var keypairs []ecloud.SSHKeyPair
	for _, arg := range args {
		keypair, err := service.GetSSHKeyPair(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving SSH key pair [%s]: %s", arg, err)
			continue
		}

		keypairs = append(keypairs, keypair)
	}

	return output.CommandOutput(cmd, SSHKeyPairCollection(keypairs))
}

func ecloudSSHKeyPairCreateCmd(f factory.ClientFactory, fs afero.Fs) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Creates a SSH key pair",
		Long:    "This command creates a SSH key pair",
		Example: "ans ecloud sshkeypair create --name test --public-key-file ~/.ssh/id_rsa.pub",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ecloudSSHKeyPairCreate(c.ECloudService(), fs, cmd, args)
		},
	}

	// Setup flags
	cmd.Flags().String("name", "", "Name of SSH key pair")
	cmd.Flags().String("public-key", "", "Public key for SSH key pair")
	cmd.Flags().String("public-key-file", "", "Path to file containing public key for SSH key pair")

	return cmd
}

func ecloudSSHKeyPairCreate(service ecloud.ECloudService, fs afero.Fs, cmd *cobra.Command, args []string) error {
	createRequest := ecloud.CreateSSHKeyPairRequest{}
	if cmd.Flags().Changed("name") {
		createRequest.Name, _ = cmd.Flags().GetString("name")
	}

	var err error
	createRequest.PublicKey, err = helper.GetContentsFromLiteralOrFilePathFlag(cmd, fs, "public-key", "public-key-file")
	if err != nil {
		return err
	}

	keypairID, err := service.CreateSSHKeyPair(createRequest)
	if err != nil {
		return fmt.Errorf("Error creating SSH key pair: %s", err)
	}

	keypair, err := service.GetSSHKeyPair(keypairID)
	if err != nil {
		return fmt.Errorf("Error retrieving new SSH key pair: %s", err)
	}

	return output.CommandOutput(cmd, SSHKeyPairCollection([]ecloud.SSHKeyPair{keypair}))
}

func ecloudSSHKeyPairUpdateCmd(f factory.ClientFactory, fs afero.Fs) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <keypair: id>...",
		Short:   "Updates an SSH key pair",
		Long:    "This command updates one or more SSH key pairs",
		Example: "ans ecloud sshkeypair update ssh-abcdef12 --name \"my keypair\"",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing SSH key pair")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ecloudSSHKeyPairUpdate(c.ECloudService(), fs, cmd, args)
		},
	}

	cmd.Flags().String("name", "", "Name of SSH key pair")
	cmd.Flags().String("public-key", "", "Public key for SSH key pair")
	cmd.Flags().String("public-key-file", "", "Path to file containing public key for SSH key pair")

	return cmd
}

func ecloudSSHKeyPairUpdate(service ecloud.ECloudService, fs afero.Fs, cmd *cobra.Command, args []string) error {
	patchRequest := ecloud.PatchSSHKeyPairRequest{}

	if cmd.Flags().Changed("name") {
		patchRequest.Name, _ = cmd.Flags().GetString("name")
	}

	if cmd.Flags().Changed("public-key") || cmd.Flags().Changed("public-key-file") {

		var err error
		patchRequest.PublicKey, err = helper.GetContentsFromLiteralOrFilePathFlag(cmd, fs, "public-key", "public-key-file")
		if err != nil {
			return err
		}
	}

	var keypairs []ecloud.SSHKeyPair
	for _, arg := range args {
		err := service.PatchSSHKeyPair(arg, patchRequest)
		if err != nil {
			output.OutputWithErrorLevelf("Error updating SSH key pair [%s]: %s", arg, err)
			continue
		}

		keypair, err := service.GetSSHKeyPair(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving updated SSH key pair [%s]: %s", arg, err)
			continue
		}

		keypairs = append(keypairs, keypair)
	}

	return output.CommandOutput(cmd, SSHKeyPairCollection(keypairs))
}

func ecloudSSHKeyPairDeleteCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "delete <keypair: id>...",
		Short:   "Removes an SSH key pair",
		Long:    "This command removes one or more SSH key pairs",
		Example: "ans ecloud sshkeypair delete ssh-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing SSH key pair")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudSSHKeyPairDelete),
	}
}

func ecloudSSHKeyPairDelete(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	for _, arg := range args {
		err := service.DeleteSSHKeyPair(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error removing SSH key pair [%s]: %s", arg, err)
		}
	}
	return nil
}
