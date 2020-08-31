package billing

import (
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/ukfast/cli/internal/pkg/clierrors"
	"github.com/ukfast/cli/test/mocks"
	"github.com/ukfast/cli/test/test_output"
	"github.com/ukfast/sdk-go/pkg/service/billing"
)

func Test_billingInvoiceList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockBillingService(mockCtrl)

		service.EXPECT().GetInvoices(gomock.Any()).Return([]billing.Invoice{}, nil).Times(1)

		billingInvoiceList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockBillingService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := billingInvoiceList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetInvoicesError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockBillingService(mockCtrl)

		service.EXPECT().GetInvoices(gomock.Any()).Return([]billing.Invoice{}, errors.New("test error")).Times(1)

		err := billingInvoiceList(service, &cobra.Command{}, []string{})

		assert.Equal(t, "Error retrieving invoices: test error", err.Error())
	})
}

func Test_billingInvoiceShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := billingInvoiceShowCmd(nil).Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := billingInvoiceShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing invoice", err.Error())
	})
}

func Test_billingInvoiceShow(t *testing.T) {
	t.Run("SingleInvoice", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockBillingService(mockCtrl)

		service.EXPECT().GetInvoice(123).Return(billing.Invoice{}, nil).Times(1)

		billingInvoiceShow(service, &cobra.Command{}, []string{"123"})
	})

	t.Run("MultipleInvoices", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockBillingService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetInvoice(123).Return(billing.Invoice{}, nil),
			service.EXPECT().GetInvoice(456).Return(billing.Invoice{}, nil),
		)

		billingInvoiceShow(service, &cobra.Command{}, []string{"123", "456"})
	})

	t.Run("GetInvoiceID_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockBillingService(mockCtrl)

		test_output.AssertErrorOutput(t, "Invalid invoice ID [abc]\n", func() {
			billingInvoiceShow(service, &cobra.Command{}, []string{"abc"})
		})
	})

	t.Run("GetInvoiceError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockBillingService(mockCtrl)

		service.EXPECT().GetInvoice(123).Return(billing.Invoice{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving invoice [123]: test error\n", func() {
			billingInvoiceShow(service, &cobra.Command{}, []string{"123"})
		})
	})
}
