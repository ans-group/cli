package cloudflare

import (
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/ukfast/cli/test/mocks"
	"github.com/ukfast/sdk-go/pkg/service/cloudflare"
)

func Test_cloudflareTotalSpendeShow(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockCloudflareService(mockCtrl)

		service.EXPECT().GetTotalSpendMonthToDate().Return(cloudflare.TotalSpend{}, nil)

		cloudflareTotalSpendShow(service, &cobra.Command{}, []string{})
	})

	t.Run("GetTotalSpendMonthToDateError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockCloudflareService(mockCtrl)

		service.EXPECT().GetTotalSpendMonthToDate().Return(cloudflare.TotalSpend{}, errors.New("test error")).Times(1)

		err := cloudflareTotalSpendShow(service, &cobra.Command{}, []string{})

		assert.Equal(t, "Error retrieving total spend: test error", err.Error())
	})
}
