package cloudflare

import (
	"errors"
	"fmt"

	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/sdk-go/pkg/service/cloudflare"
	"github.com/spf13/cobra"
)

func cloudflareAccountMemberRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "member",
		Short: "sub-commands relating to account members",
	}

	// Child commands
	cmd.AddCommand(cloudflareAccountMemberCreateCmd(f))

	return cmd
}

func cloudflareAccountMemberCreateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Creates account members",
		Long:    "This command creates account members",
		Example: "ans cloudflare account member create e84d6820-870a-4d69-89a4-30e9f1016518",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("missing account")
			}

			return nil
		},
		RunE: cloudflareCobraRunEFunc(f, cloudflareAccountMemberCreate),
	}

	cmd.Flags().String("email-address", "", "Email address for account member")
	_ = cmd.MarkFlagRequired("email-address")

	return cmd
}

func cloudflareAccountMemberCreate(service cloudflare.CloudflareService, cmd *cobra.Command, args []string) error {
	createRequest := cloudflare.CreateAccountMemberRequest{}
	createRequest.EmailAddress, _ = cmd.Flags().GetString("email-address")
	err := service.CreateAccountMember(args[0], createRequest)
	if err != nil {
		return fmt.Errorf("error creating account member: %s", err)
	}

	return nil
}
