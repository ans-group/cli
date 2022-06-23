package cloudflare

import (
	"fmt"

	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/sdk-go/pkg/service/cloudflare"
	"github.com/spf13/cobra"
)

func cloudflareOrchestratorRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "orchestrator",
		Short: "sub-commands relating to orchestrator",
	}

	// Child commands
	cmd.AddCommand(cloudflareOrchestratorCreateCmd(f))

	return cmd
}

func cloudflareOrchestratorCreateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Creates an orchestration",
		Long:    "This command creates an orchestration",
		Example: "ukfast cloudflare orchestrator create --zone-name testzone --zone-subscription c46411f8-13e3-484e-a114-9d0ce5d53502 --account-name testaccount --administrator-email-address test@test.com",
		RunE:    cloudflareCobraRunEFunc(f, cloudflareOrchestratorCreate),
	}

	cmd.Flags().String("zone-name", "", "Name of zone")
	cmd.MarkFlagRequired("zone-name")
	cmd.Flags().String("zone-subscription", "", "ID of zone plan subscription")
	cmd.MarkFlagRequired("zone-subscription")
	cmd.Flags().String("account", "", "ID of account")
	cmd.Flags().String("account-name", "", "Name of account")
	cmd.Flags().String("administrator-email-address", "", "Email address of administrator")

	return cmd
}

func cloudflareOrchestratorCreate(service cloudflare.CloudflareService, cmd *cobra.Command, args []string) error {
	createRequest := cloudflare.CreateOrchestrationRequest{}
	createRequest.ZoneName, _ = cmd.Flags().GetString("zone-name")
	createRequest.ZoneSubscriptionID, _ = cmd.Flags().GetString("zone-subscription")
	createRequest.AccountID, _ = cmd.Flags().GetString("account")
	createRequest.AccountName, _ = cmd.Flags().GetString("account-name")
	createRequest.AdministratorEmailAddress, _ = cmd.Flags().GetString("administrator-email-address")

	err := service.CreateOrchestration(createRequest)
	if err != nil {
		return fmt.Errorf("Error creating orchestration: %s", err)
	}

	return nil
}
