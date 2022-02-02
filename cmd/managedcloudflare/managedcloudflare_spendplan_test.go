package managedcloudflare

import (
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/ukfast/cli/internal/pkg/clierrors"
	"github.com/ukfast/cli/test/mocks"
	"github.com/ukfast/sdk-go/pkg/service/managedcloudflare"
)

func Test_managedcloudflareSpendPlanList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockManagedCloudflareService(mockCtrl)

		service.EXPECT().GetSpendPlans(gomock.Any()).Return([]managedcloudflare.SpendPlan{}, nil).Times(1)

		managedcloudflareSpendPlanList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockManagedCloudflareService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := managedcloudflareSpendPlanList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetSpendPlansError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockManagedCloudflareService(mockCtrl)

		service.EXPECT().GetSpendPlans(gomock.Any()).Return([]managedcloudflare.SpendPlan{}, errors.New("test error")).Times(1)

		err := managedcloudflareSpendPlanList(service, &cobra.Command{}, []string{})

		assert.Equal(t, "Error retrieving spend plans: test error", err.Error())
	})
}
