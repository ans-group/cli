package managedcloudflare

import (
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/ukfast/cli/test/mocks"
	"github.com/ukfast/sdk-go/pkg/service/managedcloudflare"
)

func Test_managedcloudflareTotalSpendeShow(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockManagedCloudflareService(mockCtrl)

		service.EXPECT().GetTotalSpendMonthToDate().Return(managedcloudflare.TotalSpend{}, nil)

		managedcloudflareTotalSpendShow(service, &cobra.Command{}, []string{})
	})

	t.Run("GetTotalSpendMonthToDateError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockManagedCloudflareService(mockCtrl)

		service.EXPECT().GetTotalSpendMonthToDate().Return(managedcloudflare.TotalSpend{}, errors.New("test error")).Times(1)

		err := managedcloudflareTotalSpendShow(service, &cobra.Command{}, []string{})

		assert.Equal(t, "Error retrieving total spend: test error", err.Error())
	})
}
