package ecloud_v2

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

func Test_ecloudAvailabilityZoneList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetAvailabilityZones(gomock.Any()).Return([]ecloud.AvailabilityZone{}, nil).Times(1)

		ecloudAvailabilityZoneList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := ecloudAvailabilityZoneList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetAvailabilityZonesError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetAvailabilityZones(gomock.Any()).Return([]ecloud.AvailabilityZone{}, errors.New("test error")).Times(1)

		err := ecloudAvailabilityZoneList(service, &cobra.Command{}, []string{})

		assert.Equal(t, "Error retrieving availability zones: test error", err.Error())
	})
}

func Test_ecloudAvailabilityZoneShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudAvailabilityZoneShowCmd(nil).Args(nil, []string{"az-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudAvailabilityZoneShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing availability zone", err.Error())
	})
}

func Test_ecloudAvailabilityZoneShow(t *testing.T) {
	t.Run("SingleAvailabilityZone", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetAvailabilityZone("az-abcdef12").Return(ecloud.AvailabilityZone{}, nil).Times(1)

		ecloudAvailabilityZoneShow(service, &cobra.Command{}, []string{"az-abcdef12"})
	})

	t.Run("MultipleAvailabilityZones", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetAvailabilityZone("az-abcdef12").Return(ecloud.AvailabilityZone{}, nil),
			service.EXPECT().GetAvailabilityZone("az-abcdef23").Return(ecloud.AvailabilityZone{}, nil),
		)

		ecloudAvailabilityZoneShow(service, &cobra.Command{}, []string{"az-abcdef12", "az-abcdef23"})
	})

	t.Run("GetAvailabilityZoneError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetAvailabilityZone("az-abcdef12").Return(ecloud.AvailabilityZone{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving availability zone [az-abcdef12]: test error\n", func() {
			ecloudAvailabilityZoneShow(service, &cobra.Command{}, []string{"az-abcdef12"})
		})
	})
}
