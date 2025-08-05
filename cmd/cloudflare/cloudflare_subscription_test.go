package cloudflare

import (
	"errors"
	"testing"

	"github.com/ans-group/cli/internal/pkg/clierrors"
	"github.com/ans-group/cli/test/mocks"
	"github.com/ans-group/sdk-go/pkg/service/cloudflare"
	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_cloudflareSubscriptionList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockCloudflareService(mockCtrl)

		service.EXPECT().GetSubscriptions(gomock.Any()).Return([]cloudflare.Subscription{}, nil).Times(1)

		cloudflareSubscriptionList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockCloudflareService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := cloudflareSubscriptionList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetSubscriptionsError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockCloudflareService(mockCtrl)

		service.EXPECT().GetSubscriptions(gomock.Any()).Return([]cloudflare.Subscription{}, errors.New("test error")).Times(1)

		err := cloudflareSubscriptionList(service, &cobra.Command{}, []string{})

		assert.Equal(t, "error retrieving subscriptions: test error", err.Error())
	})
}
