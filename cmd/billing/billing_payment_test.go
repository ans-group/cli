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

func Test_billingPaymentList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockBillingService(mockCtrl)

		service.EXPECT().GetPayments(gomock.Any()).Return([]billing.Payment{}, nil).Times(1)

		billingPaymentList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockBillingService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := billingPaymentList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetPaymentsError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockBillingService(mockCtrl)

		service.EXPECT().GetPayments(gomock.Any()).Return([]billing.Payment{}, errors.New("test error")).Times(1)

		err := billingPaymentList(service, &cobra.Command{}, []string{})

		assert.Equal(t, "Error retrieving payments: test error", err.Error())
	})
}

func Test_billingPaymentShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := billingPaymentShowCmd(nil).Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := billingPaymentShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing payment", err.Error())
	})
}

func Test_billingPaymentShow(t *testing.T) {
	t.Run("SinglePayment", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockBillingService(mockCtrl)

		service.EXPECT().GetPayment(123).Return(billing.Payment{}, nil).Times(1)

		billingPaymentShow(service, &cobra.Command{}, []string{"123"})
	})

	t.Run("MultiplePayments", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockBillingService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetPayment(123).Return(billing.Payment{}, nil),
			service.EXPECT().GetPayment(456).Return(billing.Payment{}, nil),
		)

		billingPaymentShow(service, &cobra.Command{}, []string{"123", "456"})
	})

	t.Run("GetPaymentID_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockBillingService(mockCtrl)

		test_output.AssertErrorOutput(t, "Invalid payment ID [abc]\n", func() {
			billingPaymentShow(service, &cobra.Command{}, []string{"abc"})
		})
	})

	t.Run("GetPaymentError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockBillingService(mockCtrl)

		service.EXPECT().GetPayment(123).Return(billing.Payment{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving payment [123]: test error\n", func() {
			billingPaymentShow(service, &cobra.Command{}, []string{"123"})
		})
	})
}
