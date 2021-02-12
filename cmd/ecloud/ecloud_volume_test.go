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

func Test_ecloudVolumeList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetVolumes(gomock.Any()).Return([]ecloud.Volume{}, nil).Times(1)

		ecloudVolumeList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := ecloudVolumeList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetVolumesError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetVolumes(gomock.Any()).Return([]ecloud.Volume{}, errors.New("test error")).Times(1)

		err := ecloudVolumeList(service, &cobra.Command{}, []string{})

		assert.Equal(t, "Error retrieving volumes: test error", err.Error())
	})
}

func Test_ecloudVolumeShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudVolumeShowCmd(nil).Args(nil, []string{"net-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudVolumeShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing volume", err.Error())
	})
}

func Test_ecloudVolumeShow(t *testing.T) {
	t.Run("SingleVolume", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetVolume("net-abcdef12").Return(ecloud.Volume{}, nil).Times(1)

		ecloudVolumeShow(service, &cobra.Command{}, []string{"net-abcdef12"})
	})

	t.Run("MultipleVolumes", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetVolume("net-abcdef12").Return(ecloud.Volume{}, nil),
			service.EXPECT().GetVolume("net-abcdef23").Return(ecloud.Volume{}, nil),
		)

		ecloudVolumeShow(service, &cobra.Command{}, []string{"net-abcdef12", "net-abcdef23"})
	})

	t.Run("GetVolumeError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetVolume("net-abcdef12").Return(ecloud.Volume{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving volume [net-abcdef12]: test error\n", func() {
			ecloudVolumeShow(service, &cobra.Command{}, []string{"net-abcdef12"})
		})
	})
}

// func Test_ecloudVolumeCreate(t *testing.T) {
// 	t.Run("DefaultCreate", func(t *testing.T) {
// 		mockCtrl := gomock.NewController(t)
// 		defer mockCtrl.Finish()

// 		service := mocks.NewMockECloudService(mockCtrl)
// 		cmd := ecloudVolumeCreateCmd(nil)
// 		cmd.ParseFlags([]string{"--name=testvolume"})

// 		req := ecloud.CreateVolumeRequest{
// 			Name: "testvolume",
// 		}

// 		gomock.InOrder(
// 			service.EXPECT().CreateVolume(req).Return("net-abcdef12", nil),
// 			service.EXPECT().GetVolume("net-abcdef12").Return(ecloud.Volume{}, nil),
// 		)

// 		ecloudVolumeCreate(service, cmd, []string{})
// 	})

// 	t.Run("CreateVolumeError_ReturnsError", func(t *testing.T) {
// 		mockCtrl := gomock.NewController(t)
// 		defer mockCtrl.Finish()

// 		service := mocks.NewMockECloudService(mockCtrl)
// 		cmd := ecloudVolumeCreateCmd(nil)
// 		cmd.ParseFlags([]string{"--name=testvolume"})

// 		service.EXPECT().CreateVolume(gomock.Any()).Return("", errors.New("test error")).Times(1)

// 		err := ecloudVolumeCreate(service, cmd, []string{})

// 		assert.Equal(t, "Error creating volume: test error", err.Error())
// 	})

// 	t.Run("GetVolumeError_ReturnsError", func(t *testing.T) {
// 		mockCtrl := gomock.NewController(t)
// 		defer mockCtrl.Finish()

// 		service := mocks.NewMockECloudService(mockCtrl)
// 		cmd := ecloudVolumeCreateCmd(nil)
// 		cmd.ParseFlags([]string{"--name=testvolume"})

// 		gomock.InOrder(
// 			service.EXPECT().CreateVolume(gomock.Any()).Return("net-abcdef12", nil),
// 			service.EXPECT().GetVolume("net-abcdef12").Return(ecloud.Volume{}, errors.New("test error")),
// 		)

// 		err := ecloudVolumeCreate(service, cmd, []string{})

// 		assert.Equal(t, "Error retrieving new volume: test error", err.Error())
// 	})
// }

func Test_ecloudVolumeUpdateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudVolumeUpdateCmd(nil).Args(nil, []string{"net-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudVolumeUpdateCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing volume", err.Error())
	})
}

func Test_ecloudVolumeUpdate(t *testing.T) {
	t.Run("SingleVolume", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudVolumeUpdateCmd(nil)
		cmd.ParseFlags([]string{"--name=testvolume"})

		req := ecloud.PatchVolumeRequest{
			Name: "testvolume",
		}

		gomock.InOrder(
			service.EXPECT().PatchVolume("net-abcdef12", req).Return(nil),
			service.EXPECT().GetVolume("net-abcdef12").Return(ecloud.Volume{}, nil),
		)

		ecloudVolumeUpdate(service, cmd, []string{"net-abcdef12"})
	})

	t.Run("MultipleVolumes", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().PatchVolume("net-abcdef12", gomock.Any()).Return(nil),
			service.EXPECT().GetVolume("net-abcdef12").Return(ecloud.Volume{}, nil),
			service.EXPECT().PatchVolume("net-12abcdef", gomock.Any()).Return(nil),
			service.EXPECT().GetVolume("net-12abcdef").Return(ecloud.Volume{}, nil),
		)

		ecloudVolumeUpdate(service, &cobra.Command{}, []string{"net-abcdef12", "net-12abcdef"})
	})

	t.Run("PatchVolumeError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().PatchVolume("net-abcdef12", gomock.Any()).Return(errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error updating volume [net-abcdef12]: test error\n", func() {
			ecloudVolumeUpdate(service, &cobra.Command{}, []string{"net-abcdef12"})
		})
	})

	t.Run("GetVolumeError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().PatchVolume("net-abcdef12", gomock.Any()).Return(nil),
			service.EXPECT().GetVolume("net-abcdef12").Return(ecloud.Volume{}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error retrieving updated volume [net-abcdef12]: test error\n", func() {
			ecloudVolumeUpdate(service, &cobra.Command{}, []string{"net-abcdef12"})
		})
	})
}

func Test_ecloudVolumeDeleteCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudVolumeDeleteCmd(nil).Args(nil, []string{"net-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudVolumeDeleteCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing volume", err.Error())
	})
}

func Test_ecloudVolumeDelete(t *testing.T) {
	t.Run("SingleVolume", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().DeleteVolume("net-abcdef12").Return(nil).Times(1)

		ecloudVolumeDelete(service, &cobra.Command{}, []string{"net-abcdef12"})
	})

	t.Run("MultipleVolumes", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().DeleteVolume("net-abcdef12").Return(nil),
			service.EXPECT().DeleteVolume("net-12abcdef").Return(nil),
		)

		ecloudVolumeDelete(service, &cobra.Command{}, []string{"net-abcdef12", "net-12abcdef"})
	})

	t.Run("DeleteVolumeError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().DeleteVolume("net-abcdef12").Return(errors.New("test error")).Times(1)

		test_output.AssertErrorOutput(t, "Error removing volume [net-abcdef12]: test error\n", func() {
			ecloudVolumeDelete(service, &cobra.Command{}, []string{"net-abcdef12"})
		})
	})
}
