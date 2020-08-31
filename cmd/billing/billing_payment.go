package billing

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
	"github.com/ukfast/cli/internal/pkg/helper"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/billing"
)

func billingPaymentRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "payment",
		Short: "sub-commands relating to payments",
	}

	// Child commands
	cmd.AddCommand(billingPaymentListCmd(f))
	cmd.AddCommand(billingPaymentShowCmd(f))

	return cmd
}

func billingPaymentListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists payments",
		Long:    "This command lists payments",
		Example: "ukfast billing payment list",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return billingPaymentList(c.BillingService(), cmd, args)
		},
	}
}

func billingPaymentList(service billing.BillingService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	payments, err := service.GetPayments(params)
	if err != nil {
		return fmt.Errorf("Error retrieving payments: %s", err)
	}

	return output.CommandOutput(cmd, OutputBillingPaymentsProvider(payments))
}

func billingPaymentShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <payment: id>...",
		Short:   "Shows a payment",
		Long:    "This command shows one or more payments",
		Example: "ukfast billing payment show 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing payment")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return billingPaymentShow(c.BillingService(), cmd, args)
		},
	}
}

func billingPaymentShow(service billing.BillingService, cmd *cobra.Command, args []string) error {
	var payments []billing.Payment
	for _, arg := range args {
		paymentID, err := strconv.Atoi(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Invalid payment ID [%s]", arg)
			continue
		}

		payment, err := service.GetPayment(paymentID)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving payment [%s]: %s", arg, err)
			continue
		}

		payments = append(payments, payment)
	}

	return output.CommandOutput(cmd, OutputBillingPaymentsProvider(payments))
}
