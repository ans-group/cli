package cloudflare

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
	"github.com/ukfast/cli/internal/pkg/helper"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/cloudflare"
)

func cloudflareSubscriptionRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "subscription",
		Short: "sub-commands relating to subscriptions",
	}

	// Child commands
	cmd.AddCommand(cloudflareSubscriptionListCmd(f))

	return cmd
}

func cloudflareSubscriptionListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists subscriptions",
		Long:    "This command lists subscriptions",
		Example: "ukfast cloudflare subscription list",
		RunE:    cloudflareCobraRunEFunc(f, cloudflareSubscriptionList),
	}
}

func cloudflareSubscriptionList(service cloudflare.CloudflareService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	subscriptions, err := service.GetSubscriptions(params)
	if err != nil {
		return fmt.Errorf("Error retrieving subscriptions: %s", err)
	}

	return output.CommandOutput(cmd, OutputCloudflareSubscriptionsProvider(subscriptions))
}
