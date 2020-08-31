package billing

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
	"github.com/ukfast/cli/internal/pkg/helper"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/ptr"
	"github.com/ukfast/sdk-go/pkg/service/billing"
)

func billingCardRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "card",
		Short: "sub-commands relating to cards",
	}

	// Child commands
	cmd.AddCommand(billingCardListCmd(f))
	cmd.AddCommand(billingCardShowCmd(f))
	cmd.AddCommand(billingCardCreateCmd(f))
	cmd.AddCommand(billingCardUpdateCmd(f))
	cmd.AddCommand(billingCardDeleteCmd(f))

	return cmd
}

func billingCardListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists cards",
		Long:    "This command lists cards",
		Example: "ukfast billing card list",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return billingCardList(c.BillingService(), cmd, args)
		},
	}
}

func billingCardList(service billing.BillingService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	cards, err := service.GetCards(params)
	if err != nil {
		return fmt.Errorf("Error retrieving cards: %s", err)
	}

	return output.CommandOutput(cmd, OutputBillingCardsProvider(cards))
}

func billingCardShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <card: id>...",
		Short:   "Shows a card",
		Long:    "This command shows one or more cards",
		Example: "ukfast billing card show 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing card")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return billingCardShow(c.BillingService(), cmd, args)
		},
	}
}

func billingCardShow(service billing.BillingService, cmd *cobra.Command, args []string) error {
	var cards []billing.Card
	for _, arg := range args {
		cardID, err := strconv.Atoi(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Invalid card ID [%s]", arg)
			continue
		}

		card, err := service.GetCard(cardID)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving card [%s]: %s", arg, err)
			continue
		}

		cards = append(cards, card)
	}

	return output.CommandOutput(cmd, OutputBillingCardsProvider(cards))
}

func billingCardCreateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Creates a card",
		Long:    "This command creates a card",
		Example: "ukfast billing card create",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return billingCardCreate(c.BillingService(), cmd, args)
		},
	}

	// Setup flags
	cmd.Flags().String("friendly-name", "", "Display name of card")
	cmd.Flags().String("name", "", "Name on card")
	cmd.MarkFlagRequired("name")
	cmd.Flags().String("address", "", "Address of cardholder")
	cmd.MarkFlagRequired("address")
	cmd.Flags().String("postcode", "", "Postcode of cardholder")
	cmd.MarkFlagRequired("postcode")
	cmd.Flags().String("card-number", "", "16 digit card number")
	cmd.MarkFlagRequired("card-number")
	cmd.Flags().String("card-type", "", "Type of card")
	cmd.MarkFlagRequired("card-type")
	cmd.Flags().String("valid-from", "", "Date card is valid from")
	cmd.MarkFlagRequired("valid-from")
	cmd.Flags().String("expiry", "", "Expiry date of card")
	cmd.MarkFlagRequired("expiry")
	cmd.Flags().Int("issue-number", 0, "Issue number of card")
	cmd.Flags().Bool("primary-card", false, "Specifies whether this card should be the primary card")

	return cmd
}

func billingCardCreate(service billing.BillingService, cmd *cobra.Command, args []string) error {
	createRequest := billing.CreateCardRequest{}
	createRequest.FriendlyName, _ = cmd.Flags().GetString("friendly-name")
	createRequest.Name, _ = cmd.Flags().GetString("name")
	createRequest.Address, _ = cmd.Flags().GetString("address")
	createRequest.Postcode, _ = cmd.Flags().GetString("postcode")
	createRequest.CardNumber, _ = cmd.Flags().GetString("card-number")
	createRequest.CardType, _ = cmd.Flags().GetString("card-type")
	createRequest.ValidFrom, _ = cmd.Flags().GetString("valid-from")
	createRequest.Expiry, _ = cmd.Flags().GetString("expiry")
	createRequest.IssueNumber, _ = cmd.Flags().GetInt("issue-number")
	createRequest.PrimaryCard, _ = cmd.Flags().GetBool("primary-card")

	id, err := service.CreateCard(createRequest)
	if err != nil {
		return fmt.Errorf("Error creating card: %s", err)
	}

	card, err := service.GetCard(id)
	if err != nil {
		return fmt.Errorf("Error retrieving new card: %s", err)
	}

	return output.CommandOutput(cmd, OutputBillingCardsProvider([]billing.Card{card}))
}

func billingCardUpdateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <virtualmachine: id>...",
		Short:   "Updates a card",
		Long:    "This command updates one or more cards",
		Example: "ukfast billing card update 123 --name \"test card 1\"",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing card")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return billingCardUpdate(c.BillingService(), cmd, args)
		},
	}

	// Setup flags
	cmd.Flags().String("friendly-name", "", "Display name of card")
	cmd.Flags().String("name", "", "Name on card")
	cmd.Flags().String("address", "", "Address of cardholder")
	cmd.Flags().String("postcode", "", "Postcode of cardholder")
	cmd.Flags().String("card-number", "", "16 digit card number")
	cmd.Flags().String("card-type", "", "Type of card")
	cmd.Flags().String("valid-from", "", "Date card is valid from")
	cmd.Flags().String("expiry", "", "Expiry date of card")
	cmd.Flags().Int("issue-number", 0, "Issue number of card")
	cmd.Flags().Bool("primary-card", false, "Specifies whether this card should be the primary card")

	return cmd
}

func billingCardUpdate(service billing.BillingService, cmd *cobra.Command, args []string) error {

	patchRequest := billing.PatchCardRequest{}
	patchRequest.FriendlyName, _ = cmd.Flags().GetString("friendly-name")
	patchRequest.Name, _ = cmd.Flags().GetString("name")
	patchRequest.Address, _ = cmd.Flags().GetString("address")
	patchRequest.Postcode, _ = cmd.Flags().GetString("postcode")
	patchRequest.CardNumber, _ = cmd.Flags().GetString("card-number")
	patchRequest.CardType, _ = cmd.Flags().GetString("card-type")
	patchRequest.ValidFrom, _ = cmd.Flags().GetString("valid-from")
	patchRequest.Expiry, _ = cmd.Flags().GetString("expiry")
	patchRequest.IssueNumber, _ = cmd.Flags().GetInt("issue-number")

	if cmd.Flags().Changed("primary-card") {
		primaryCard, _ := cmd.Flags().GetBool("primary-card")
		patchRequest.PrimaryCard = ptr.Bool(primaryCard)
	}

	var cards []billing.Card

	for _, arg := range args {
		cardID, err := strconv.Atoi(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Invalid card ID [%s]", arg)
			continue
		}

		err = service.PatchCard(cardID, patchRequest)
		if err != nil {
			output.OutputWithErrorLevelf("Error updating card [%d]: %s", cardID, err.Error())
			continue
		}

		card, err := service.GetCard(cardID)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving updated card [%d]: %s", cardID, err.Error())
			continue
		}

		cards = append(cards, card)
	}

	return output.CommandOutput(cmd, OutputBillingCardsProvider(cards))
}

func billingCardDeleteCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "delete <card: id>...",
		Short:   "Removes a card",
		Long:    "This command removes one or more cards",
		Example: "ukfast billing card delete 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing card")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			billingCardDelete(c.BillingService(), cmd, args)
			return nil
		},
	}
}

func billingCardDelete(service billing.BillingService, cmd *cobra.Command, args []string) {
	for _, arg := range args {
		cardID, err := strconv.Atoi(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Invalid card ID [%s]", arg)
			continue
		}

		err = service.DeleteCard(cardID)
		if err != nil {
			output.OutputWithErrorLevelf("Error removing card [%d]: %s", cardID, err)
			continue
		}
	}
}
