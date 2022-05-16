package cloudflare

import (
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/ukfast/cli/internal/pkg/clierrors"
	"github.com/ukfast/cli/test/mocks"
	"github.com/ukfast/sdk-go/pkg/service/cloudflare"
)

func Test_cloudflareSpendPlanList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockCloudflareService(mockCtrl)

		service.EXPECT().GetSpendPlans(gomock.Any()).Return([]cloudflare.SpendPlan{}, nil).Times(1)

		cloudflareSpendPlanList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockCloudflareService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := cloudflareSpendPlanList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetSpendPlansError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockCloudflareService(mockCtrl)

		service.EXPECT().GetSpendPlans(gomock.Any()).Return([]cloudflare.SpendPlan{}, errors.New("test error")).Times(1)

		err := cloudflareSpendPlanList(service, &cobra.Command{}, []string{})

		assert.Equal(t, "Error retrieving spend plans: test error", err.Error())
	})
}
