package storage

import (
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/ukfast/cli/internal/pkg/clierrors"
	"github.com/ukfast/cli/test/mocks"
	"github.com/ukfast/cli/test/test_output"
	"github.com/ukfast/sdk-go/pkg/service/storage"
)

func Test_storageVolumeList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockStorageService(mockCtrl)

		service.EXPECT().GetVolumes(gomock.Any()).Return([]storage.Volume{}, nil).Times(1)

		storageVolumeList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockStorageService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := storageVolumeList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetVolumesError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockStorageService(mockCtrl)

		service.EXPECT().GetVolumes(gomock.Any()).Return([]storage.Volume{}, errors.New("test error")).Times(1)

		err := storageVolumeList(service, &cobra.Command{}, []string{})

		assert.Equal(t, "Error retrieving volumes: test error", err.Error())
	})
}

func Test_storageVolumeShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := storageVolumeShowCmd(nil).Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := storageVolumeShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing volume", err.Error())
	})
}

func Test_storageVolumeShow(t *testing.T) {
	t.Run("SingleVolume", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockStorageService(mockCtrl)

		service.EXPECT().GetVolume(123).Return(storage.Volume{}, nil).Times(1)

		storageVolumeShow(service, &cobra.Command{}, []string{"123"})
	})

	t.Run("MultipleVolumes", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockStorageService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetVolume(123).Return(storage.Volume{}, nil),
			service.EXPECT().GetVolume(456).Return(storage.Volume{}, nil),
		)

		storageVolumeShow(service, &cobra.Command{}, []string{"123", "456"})
	})

	t.Run("GetVolumeID_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockStorageService(mockCtrl)

		test_output.AssertErrorOutput(t, "Invalid volume ID [abc]\n", func() {
			storageVolumeShow(service, &cobra.Command{}, []string{"abc"})
		})
	})

	t.Run("GetVolumeError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockStorageService(mockCtrl)

		service.EXPECT().GetVolume(123).Return(storage.Volume{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving volume [123]: test error\n", func() {
			storageVolumeShow(service, &cobra.Command{}, []string{"123"})
		})
	})
}
