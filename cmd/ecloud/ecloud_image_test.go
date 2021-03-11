package ecloud

import (
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/ukfast/cli/internal/pkg/clierrors"
	"github.com/ukfast/cli/test/mocks"
	"github.com/ukfast/cli/test/test_output"
	"github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func Test_ecloudImageList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetImages(gomock.Any()).Return([]ecloud.Image{}, nil).Times(1)

		ecloudImageList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := ecloudImageList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetImagesError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetImages(gomock.Any()).Return([]ecloud.Image{}, errors.New("test error")).Times(1)

		err := ecloudImageList(service, &cobra.Command{}, []string{})

		assert.Equal(t, "Error retrieving images: test error", err.Error())
	})
}

func Test_ecloudImageShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudImageShowCmd(nil).Args(nil, []string{"img-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudImageShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing image", err.Error())
	})
}

func Test_ecloudImageShow(t *testing.T) {
	t.Run("SingleImage", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetImage("img-abcdef12").Return(ecloud.Image{}, nil).Times(1)

		ecloudImageShow(service, &cobra.Command{}, []string{"img-abcdef12"})
	})

	t.Run("MultipleImages", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetImage("img-abcdef12").Return(ecloud.Image{}, nil),
			service.EXPECT().GetImage("img-abcdef23").Return(ecloud.Image{}, nil),
		)

		ecloudImageShow(service, &cobra.Command{}, []string{"img-abcdef12", "img-abcdef23"})
	})

	t.Run("GetImageError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetImage("img-abcdef12").Return(ecloud.Image{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving image [img-abcdef12]: test error\n", func() {
			ecloudImageShow(service, &cobra.Command{}, []string{"img-abcdef12"})
		})
	})
}
