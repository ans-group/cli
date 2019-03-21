package cmd

import (
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/cli/test"
	"github.com/ukfast/cli/test/mocks"
	"github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func Test_ecloudSiteList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetSites(gomock.Any()).Return([]ecloud.Site{}, nil).Times(1)

		ecloudSiteList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_OutputsFatal", func(t *testing.T) {
		code := 0
		oldOutputExit := output.SetOutputExit(func(c int) {
			code = c
		})
		defer func() { output.SetOutputExit(oldOutputExit) }()
		defer func() { flagFilter = nil }()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		flagFilter = []string{"invalidfilter"}

		output := test.CatchStdErr(t, func() {
			ecloudSiteList(service, &cobra.Command{}, []string{})
		})

		assert.Equal(t, 1, code)
		assert.Equal(t, "Missing value for filtering\n", output)
	})

	t.Run("GetSitesError_OutputsFatal", func(t *testing.T) {
		code := 0
		oldOutputExit := output.SetOutputExit(func(c int) {
			code = c
		})
		defer func() { output.SetOutputExit(oldOutputExit) }()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetSites(gomock.Any()).Return([]ecloud.Site{}, errors.New("test error")).Times(1)

		output := test.CatchStdErr(t, func() {
			ecloudSiteList(service, &cobra.Command{}, []string{})
		})

		assert.Equal(t, 1, code)
		assert.Equal(t, "Error retrieving sites: test error\n", output)
	})
}

func Test_ecloudSiteShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudSiteShowCmd().Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudSiteShowCmd().Args(nil, []string{})

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

		output := test.CatchStdErr(t, func() {
			ecloudSiteShow(service, &cobra.Command{}, []string{"abc"})
		})

		assert.Equal(t, "Invalid site ID [abc]\n", output)
	})

	t.Run("GetSiteError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetSite(123).Return(ecloud.Site{}, errors.New("test error"))

		output := test.CatchStdErr(t, func() {
			ecloudSiteShow(service, &cobra.Command{}, []string{"123"})
		})

		assert.Equal(t, "Error retrieving site [123]: test error\n", output)
	})
}
