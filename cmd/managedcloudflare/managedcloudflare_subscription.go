package managedcloudflare

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
	"github.com/ukfast/cli/internal/pkg/helper"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/managedcloudflare"
)

func managedcloudflareSubscriptionRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "subscription",
		Short: "sub-commands relating to subscriptions",
	}

	// Child commands
	cmd.AddCommand(managedcloudflareSubscriptionListCmd(f))

	return cmd
}

func managedcloudflareSubscriptionListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists subscriptions",
		Long:    "This command lists subscriptions",
		Example: "ukfast managedcloudflare subscription list",
		RunE:    managedcloudflareCobraRunEFunc(f, managedcloudflareSubscriptionList),
	}
}

func managedcloudflareSubscriptionList(service managedcloudflare.ManagedCloudflareService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	subscriptions, err := service.GetSubscriptions(params)
	if err != nil {
		return fmt.Errorf("Error retrieving subscriptions: %s", err)
	}

	return output.CommandOutput(cmd, OutputManagedCloudflareSubscriptionsProvider(subscriptions))
}
