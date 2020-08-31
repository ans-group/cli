package billing

import (
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/ukfast/cli/test/mocks"
	"github.com/ukfast/sdk-go/pkg/service/billing"
)

func Test_billingDirectDebitShow(t *testing.T) {
	t.Run("GetSuccess_NoError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		service := mocks.NewMockBillingService(mockCtrl)
		service.EXPECT().GetDirectDebit().Return(billing.DirectDebit{}, nil).Times(1)

		billingDirectDebitShow(service, &cobra.Command{}, []string{"123"})
	})

	t.Run("GetDirectDebitError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		service := mocks.NewMockBillingService(mockCtrl)
		service.EXPECT().GetDirectDebit().Return(billing.DirectDebit{}, errors.New("test error"))

		err := billingDirectDebitShow(service, &cobra.Command{}, []string{"123"})

		assert.Equal(t, "Error retrieving direct debit details: test error", err.Error())
	})
}
