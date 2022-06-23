package billing

import (
	"errors"
	"testing"

	"github.com/ans-group/cli/internal/pkg/clierrors"
	"github.com/ans-group/cli/test/mocks"
	"github.com/ans-group/cli/test/test_output"
	"github.com/ans-group/sdk-go/pkg/service/billing"
	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_billingInvoiceQueryList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockBillingService(mockCtrl)

		service.EXPECT().GetInvoiceQueries(gomock.Any()).Return([]billing.InvoiceQuery{}, nil).Times(1)

		billingInvoiceQueryList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockBillingService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := billingInvoiceQueryList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetInvoiceQueriesError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockBillingService(mockCtrl)

		service.EXPECT().GetInvoiceQueries(gomock.Any()).Return([]billing.InvoiceQuery{}, errors.New("test error")).Times(1)

		err := billingInvoiceQueryList(service, &cobra.Command{}, []string{})

		assert.Equal(t, "Error retrieving invoice queries: test error", err.Error())
	})
}

func Test_billingInvoiceQueryShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := billingInvoiceQueryShowCmd(nil).Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := billingInvoiceQueryShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing invoice query", err.Error())
	})
}

func Test_billingInvoiceQueryShow(t *testing.T) {
	t.Run("SingleInvoiceQuery", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockBillingService(mockCtrl)

		service.EXPECT().GetInvoiceQuery(123).Return(billing.InvoiceQuery{}, nil).Times(1)

		billingInvoiceQueryShow(service, &cobra.Command{}, []string{"123"})
	})

	t.Run("MultipleInvoiceQueries", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockBillingService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetInvoiceQuery(123).Return(billing.InvoiceQuery{}, nil),
			service.EXPECT().GetInvoiceQuery(456).Return(billing.InvoiceQuery{}, nil),
		)

		billingInvoiceQueryShow(service, &cobra.Command{}, []string{"123", "456"})
	})

	t.Run("GetInvoiceQueryID_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockBillingService(mockCtrl)

		test_output.AssertErrorOutput(t, "Invalid invoice query ID [abc]\n", func() {
			billingInvoiceQueryShow(service, &cobra.Command{}, []string{"abc"})
		})
	})

	t.Run("GetInvoiceQueryError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockBillingService(mockCtrl)

		service.EXPECT().GetInvoiceQuery(123).Return(billing.InvoiceQuery{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving invoice query [123]: test error\n", func() {
			billingInvoiceQueryShow(service, &cobra.Command{}, []string{"123"})
		})
	})
}

func Test_billingInvoiceQueryCreate(t *testing.T) {
	t.Run("DefaultCreate", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockBillingService(mockCtrl)
		cmd := billingInvoiceQueryCreateCmd(nil)
		cmd.Flags().Set("contact-id", "4")

		gomock.InOrder(
			service.EXPECT().CreateInvoiceQuery(gomock.Any()).Do(func(req billing.CreateInvoiceQueryRequest) {
				if req.ContactID != 4 {
					t.Fatalf("Expected ContactID '4', got '%d'", req.ContactID)
				}
			}).Return(123, nil),
			service.EXPECT().GetInvoiceQuery(123).Return(billing.InvoiceQuery{}, nil),
		)

		billingInvoiceQueryCreate(service, cmd, []string{})
	})

	t.Run("CreateInvoiceQueryError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockBillingService(mockCtrl)

		service.EXPECT().CreateInvoiceQuery(gomock.Any()).Return(123, errors.New("test error")).Times(1)

		err := billingInvoiceQueryCreate(service, &cobra.Command{}, []string{})

		assert.Equal(t, "Error creating invoice query: test error", err.Error())
	})

	t.Run("GetInvoiceQueryError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockBillingService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().CreateInvoiceQuery(gomock.Any()).Return(123, nil),
			service.EXPECT().GetInvoiceQuery(123).Return(billing.InvoiceQuery{}, errors.New("test error")),
		)

		err := billingInvoiceQueryCreate(service, &cobra.Command{}, []string{})

		assert.Equal(t, "Error retrieving new invoice query [123]: test error", err.Error())
	})
}
