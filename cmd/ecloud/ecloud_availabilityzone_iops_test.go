package ecloud

import (
	"errors"
	"testing"

	"github.com/ans-group/cli/internal/pkg/clierrors"
	"github.com/ans-group/cli/test/mocks"
	"github.com/ans-group/sdk-go/pkg/service/ecloud"
	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_ecloudAvailabilityZoneIOPSTiersListCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudAvailabilityZoneIOPSTierListCmd(nil).Args(nil, []string{"fwp-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudAvailabilityZoneIOPSTierListCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing availability zone", err.Error())
	})
}

func Test_ecloudAvailabilityZoneIOPSTiersList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetAvailabilityZoneIOPSTiers("fwp-abcdef12", gomock.Any()).Return([]ecloud.IOPSTier{}, nil).Times(1)

		ecloudAvailabilityZoneIOPSTierList(service, &cobra.Command{}, []string{"fwp-abcdef12"})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := ecloudAvailabilityZoneIOPSTierList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetAvailabilityZonesError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetAvailabilityZoneIOPSTiers("fwp-abcdef12", gomock.Any()).Return([]ecloud.IOPSTier{}, errors.New("test error")).Times(1)

		err := ecloudAvailabilityZoneIOPSTierList(service, &cobra.Command{}, []string{"fwp-abcdef12"})

		assert.Equal(t, "Error retrieving availability zone ports: test error", err.Error())
	})
}
