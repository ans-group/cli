package account

import (
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/ukfast/cli/internal/pkg/clierrors"
	"github.com/ukfast/cli/test/mocks"
	"github.com/ukfast/cli/test/test_output"
	"github.com/ukfast/sdk-go/pkg/service/account"
)

func Test_accountInvoiceQueryList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockAccountService(mockCtrl)

		service.EXPECT().GetInvoiceQueries(gomock.Any()).Return([]account.InvoiceQuery{}, nil).Times(1)

		accountInvoiceQueryList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockAccountService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := accountInvoiceQueryList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetInvoiceQueriesError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockAccountService(mockCtrl)

		service.EXPECT().GetInvoiceQueries(gomock.Any()).Return([]account.InvoiceQuery{}, errors.New("test error")).Times(1)

		err := accountInvoiceQueryList(service, &cobra.Command{}, []string{})

		assert.Equal(t, "Error retrieving invoice queries: test error", err.Error())
	})
}

func Test_accountInvoiceQueryShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := accountInvoiceQueryShowCmd(nil).Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := accountInvoiceQueryShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing invoice query", err.Error())
	})
}

func Test_accountInvoiceQueryShow(t *testing.T) {
	t.Run("SingleInvoiceQuery", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockAccountService(mockCtrl)

		service.EXPECT().GetInvoiceQuery(123).Return(account.InvoiceQuery{}, nil).Times(1)

		accountInvoiceQueryShow(service, &cobra.Command{}, []string{"123"})
	})

	t.Run("MultipleInvoiceQueries", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockAccountService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetInvoiceQuery(123).Return(account.InvoiceQuery{}, nil),
			service.EXPECT().GetInvoiceQuery(456).Return(account.InvoiceQuery{}, nil),
		)

		accountInvoiceQueryShow(service, &cobra.Command{}, []string{"123", "456"})
	})

	t.Run("GetInvoiceQueryID_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockAccountService(mockCtrl)

		test_output.AssertErrorOutput(t, "Invalid invoice query ID [abc]\n", func() {
			accountInvoiceQueryShow(service, &cobra.Command{}, []string{"abc"})
		})
	})

	t.Run("GetInvoiceQueryError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockAccountService(mockCtrl)

		service.EXPECT().GetInvoiceQuery(123).Return(account.InvoiceQuery{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving invoice query [123]: test error\n", func() {
			accountInvoiceQueryShow(service, &cobra.Command{}, []string{"123"})
		})
	})
}

func Test_accountInvoiceQueryCreate(t *testing.T) {
	t.Run("DefaultCreate", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockAccountService(mockCtrl)
		cmd := accountInvoiceQueryCreateCmd(nil)
		cmd.Flags().Set("contact-id", "4")

		gomock.InOrder(
			service.EXPECT().CreateInvoiceQuery(gomock.Any()).Do(func(req account.CreateInvoiceQueryRequest) {
				if req.ContactID != 4 {
					t.Fatalf("Expected ContactID '4', got '%d'", req.ContactID)
				}
			}).Return(123, nil),
			service.EXPECT().GetInvoiceQuery(123).Return(account.InvoiceQuery{}, nil),
		)

		accountInvoiceQueryCreate(service, cmd, []string{})
	})

	t.Run("CreateInvoiceQueryError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockAccountService(mockCtrl)

		service.EXPECT().CreateInvoiceQuery(gomock.Any()).Return(123, errors.New("test error")).Times(1)

		err := accountInvoiceQueryCreate(service, &cobra.Command{}, []string{})

		assert.Equal(t, "Error creating invoice query: test error", err.Error())
	})

	t.Run("GetInvoiceQueryError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockAccountService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().CreateInvoiceQuery(gomock.Any()).Return(123, nil),
			service.EXPECT().GetInvoiceQuery(123).Return(account.InvoiceQuery{}, errors.New("test error")),
		)

		err := accountInvoiceQueryCreate(service, &cobra.Command{}, []string{})

		assert.Equal(t, "Error retrieving new invoice query [123]: test error", err.Error())
	})
}
