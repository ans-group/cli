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

func Test_billingRecurringCostList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockBillingService(mockCtrl)

		service.EXPECT().GetRecurringCosts(gomock.Any()).Return([]billing.RecurringCost{}, nil).Times(1)

		billingRecurringCostList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockBillingService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := billingRecurringCostList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetRecurringCostsError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockBillingService(mockCtrl)

		service.EXPECT().GetRecurringCosts(gomock.Any()).Return([]billing.RecurringCost{}, errors.New("test error")).Times(1)

		err := billingRecurringCostList(service, &cobra.Command{}, []string{})

		assert.Equal(t, "Error retrieving recurring costs: test error", err.Error())
	})
}

func Test_billingRecurringCostShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := billingRecurringCostShowCmd(nil).Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := billingRecurringCostShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing recurring cost", err.Error())
	})
}

func Test_billingRecurringCostShow(t *testing.T) {
	t.Run("SingleRecurringCost", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockBillingService(mockCtrl)

		service.EXPECT().GetRecurringCost(123).Return(billing.RecurringCost{}, nil).Times(1)

		billingRecurringCostShow(service, &cobra.Command{}, []string{"123"})
	})

	t.Run("MultipleRecurringCosts", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockBillingService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetRecurringCost(123).Return(billing.RecurringCost{}, nil),
			service.EXPECT().GetRecurringCost(456).Return(billing.RecurringCost{}, nil),
		)

		billingRecurringCostShow(service, &cobra.Command{}, []string{"123", "456"})
	})

	t.Run("GetRecurringCostID_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockBillingService(mockCtrl)

		test_output.AssertErrorOutput(t, "Invalid recurring cost ID [abc]\n", func() {
			billingRecurringCostShow(service, &cobra.Command{}, []string{"abc"})
		})
	})

	t.Run("GetRecurringCostError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockBillingService(mockCtrl)

		service.EXPECT().GetRecurringCost(123).Return(billing.RecurringCost{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving recurring cost [123]: test error\n", func() {
			billingRecurringCostShow(service, &cobra.Command{}, []string{"123"})
		})
	})
}
