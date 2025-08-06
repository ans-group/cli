package ecloud

import (
	"errors"
	"testing"

	"github.com/ans-group/cli/internal/pkg/clierrors"
	"github.com/ans-group/cli/test/mocks"
	"github.com/ans-group/cli/test/test_output"
	"github.com/ans-group/sdk-go/pkg/service/ecloud"
	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_ecloudResourceTierList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetResourceTiers(gomock.Any()).Return([]ecloud.ResourceTier{}, nil).Times(1)

		ecloudResourceTierList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := ecloudResourceTierList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetResourceTiersError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetResourceTiers(gomock.Any()).Return([]ecloud.ResourceTier{}, errors.New("test error")).Times(1)

		err := ecloudResourceTierList(service, &cobra.Command{}, []string{})

		assert.Equal(t, "error retrieving resource tiers: test error", err.Error())
	})
}

func Test_ecloudResourceTierShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudResourceTierShowCmd(nil).Args(nil, []string{"rt-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudResourceTierShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "missing resource tier", err.Error())
	})
}

func Test_ecloudResourceTierShow(t *testing.T) {
	t.Run("SingleResourceTier", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetResourceTier("rt-abcdef12").Return(ecloud.ResourceTier{}, nil).Times(1)

		ecloudResourceTierShow(service, &cobra.Command{}, []string{"rt-abcdef12"})
	})

	t.Run("MultipleResourceTiers", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetResourceTier("rt-abcdef12").Return(ecloud.ResourceTier{}, nil),
			service.EXPECT().GetResourceTier("rt-abcdef23").Return(ecloud.ResourceTier{}, nil),
		)

		ecloudResourceTierShow(service, &cobra.Command{}, []string{"rt-abcdef12", "rt-abcdef23"})
	})

	t.Run("GetResourceTierError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetResourceTier("rt-abcdef12").Return(ecloud.ResourceTier{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving resource tier [rt-abcdef12]: test error\n", func() {
			ecloudResourceTierShow(service, &cobra.Command{}, []string{"rt-abcdef12"})
		})
	})
}
