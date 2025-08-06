package cloudflare

import (
	"errors"
	"testing"

	"github.com/ans-group/cli/test/mocks"
	"github.com/ans-group/sdk-go/pkg/service/cloudflare"
	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
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

		assert.Equal(t, "error retrieving total spend: test error", err.Error())
	})
}
