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

func Test_ecloudRegionList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetRegions(gomock.Any()).Return([]ecloud.Region{}, nil).Times(1)

		ecloudRegionList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := ecloudRegionList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetRegionsError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetRegions(gomock.Any()).Return([]ecloud.Region{}, errors.New("test error")).Times(1)

		err := ecloudRegionList(service, &cobra.Command{}, []string{})

		assert.Equal(t, "error retrieving regions: test error", err.Error())
	})
}

func Test_ecloudRegionShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudRegionShowCmd(nil).Args(nil, []string{"reg-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudRegionShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "missing region", err.Error())
	})
}

func Test_ecloudRegionShow(t *testing.T) {
	t.Run("SingleRegion", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetRegion("reg-abcdef12").Return(ecloud.Region{}, nil).Times(1)

		ecloudRegionShow(service, &cobra.Command{}, []string{"reg-abcdef12"})
	})

	t.Run("MultipleRegions", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetRegion("reg-abcdef12").Return(ecloud.Region{}, nil),
			service.EXPECT().GetRegion("reg-abcdef23").Return(ecloud.Region{}, nil),
		)

		ecloudRegionShow(service, &cobra.Command{}, []string{"reg-abcdef12", "reg-abcdef23"})
	})

	t.Run("GetRegionError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetRegion("reg-abcdef12").Return(ecloud.Region{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving region [reg-abcdef12]: test error\n", func() {
			ecloudRegionShow(service, &cobra.Command{}, []string{"reg-abcdef12"})
		})
	})
}
