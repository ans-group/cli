package cloudflare

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
	"github.com/ukfast/sdk-go/pkg/service/cloudflare"
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
		Example: "ukfast cloudflare account member create",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing account")
			}

			return nil
		},
		RunE: cloudflareCobraRunEFunc(f, cloudflareAccountMemberCreate),
	}

	cmd.Flags().String("email-address", "", "Email address for account member")
	cmd.MarkFlagRequired("email-address")

	return cmd
}

func cloudflareAccountMemberCreate(service cloudflare.CloudflareService, cmd *cobra.Command, args []string) error {
	createRequest := cloudflare.CreateAccountMemberRequest{}
	createRequest.EmailAddress, _ = cmd.Flags().GetString("email-address")
	err := service.CreateAccountMember(args[0], createRequest)
	if err != nil {
		return fmt.Errorf("Error creating account member: %s", err)
	}

	return nil
}
