package managedcloudflare

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
	"github.com/ukfast/sdk-go/pkg/service/managedcloudflare"
)

func managedcloudflareOrchestratorRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "orchestrator",
		Short: "sub-commands relating to orchestrator",
	}

	// Child commands
	cmd.AddCommand(managedcloudflareOrchestratorCreateCmd(f))

	return cmd
}

func managedcloudflareOrchestratorCreateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Creates an orchestration",
		Long:    "This command creates an orchestration",
		Example: "ukfast managedcloudflare orchestrator create --zone-name testzone --zone-subscription-type testtype --account-name testaccount --administrator-email-address test@test.com",
		RunE:    managedcloudflareCobraRunEFunc(f, managedcloudflareOrchestratorCreate),
	}

	cmd.Flags().String("zone-name", "", "Name of zone")
	cmd.MarkFlagRequired("zone-name")
	cmd.Flags().String("zone-subscription-type", "", "Type of zone subscription")
	cmd.MarkFlagRequired("zone-subscription-type")
	cmd.Flags().String("account", "", "ID of account")
	cmd.Flags().String("account-name", "", "Name of account")
	cmd.Flags().String("administrator-email-address", "", "Email address of administrator")

	return cmd
}

func managedcloudflareOrchestratorCreate(service managedcloudflare.ManagedCloudflareService, cmd *cobra.Command, args []string) error {
	createRequest := managedcloudflare.CreateOrchestrationRequest{}
	createRequest.ZoneName, _ = cmd.Flags().GetString("zone-name")
	createRequest.ZoneSubscriptionType, _ = cmd.Flags().GetString("zone-subscription-type")
	createRequest.AccountID, _ = cmd.Flags().GetString("account")
	createRequest.AccountName, _ = cmd.Flags().GetString("account-name")
	createRequest.AdministratorEmailAddress, _ = cmd.Flags().GetString("administrator-email-address")

	err := service.CreateOrchestration(createRequest)
	if err != nil {
		return fmt.Errorf("Error creating orchestration: %s", err)
	}

	return nil
}
