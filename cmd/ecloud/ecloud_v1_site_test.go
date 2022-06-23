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

func Test_ecloudSiteList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetSites(gomock.Any()).Return([]ecloud.Site{}, nil).Times(1)

		ecloudSiteList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := ecloudSiteList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetSitesError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetSites(gomock.Any()).Return([]ecloud.Site{}, errors.New("test error")).Times(1)

		err := ecloudSiteList(service, &cobra.Command{}, []string{})

		assert.Equal(t, "Error retrieving sites: test error", err.Error())
	})
}

func Test_ecloudSiteShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudSiteShowCmd(nil).Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudSiteShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing site", err.Error())
	})
}

func Test_ecloudSiteShow(t *testing.T) {
	t.Run("SingleSite", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetSite(123).Return(ecloud.Site{}, nil).Times(1)

		ecloudSiteShow(service, &cobra.Command{}, []string{"123"})
	})

	t.Run("MultipleSites", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetSite(123).Return(ecloud.Site{}, nil),
			service.EXPECT().GetSite(456).Return(ecloud.Site{}, nil),
		)

		ecloudSiteShow(service, &cobra.Command{}, []string{"123", "456"})
	})

	t.Run("GetSiteID_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		test_output.AssertErrorOutput(t, "Invalid site ID [abc]\n", func() {
			ecloudSiteShow(service, &cobra.Command{}, []string{"abc"})
		})
	})

	t.Run("GetSiteError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetSite(123).Return(ecloud.Site{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving site [123]: test error\n", func() {
			ecloudSiteShow(service, &cobra.Command{}, []string{"123"})
		})
	})
}
