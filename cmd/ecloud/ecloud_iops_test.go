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

func Test_ecloudIOPSTierList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetIOPSTiers(gomock.Any()).Return([]ecloud.IOPSTier{}, nil).Times(1)

		ecloudIOPSTierList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := ecloudIOPSTierList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetIOPSTiersError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetIOPSTiers(gomock.Any()).Return([]ecloud.IOPSTier{}, errors.New("test error")).Times(1)

		err := ecloudIOPSTierList(service, &cobra.Command{}, []string{})

		assert.Equal(t, "Error retrieving IOPS tiers: test error", err.Error())
	})
}

func Test_ecloudIOPSTierShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudIOPSTierShowCmd(nil).Args(nil, []string{"iops-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudIOPSTierShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing IOPS tier", err.Error())
	})
}

func Test_ecloudIOPSTierShow(t *testing.T) {
	t.Run("SingleIOPS", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetIOPSTier("iops-abcdef12").Return(ecloud.IOPSTier{}, nil).Times(1)

		ecloudIOPSTierShow(service, &cobra.Command{}, []string{"iops-abcdef12"})
	})

	t.Run("MultipleIOPSs", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetIOPSTier("iops-abcdef12").Return(ecloud.IOPSTier{}, nil),
			service.EXPECT().GetIOPSTier("iops-abcdef23").Return(ecloud.IOPSTier{}, nil),
		)

		ecloudIOPSTierShow(service, &cobra.Command{}, []string{"iops-abcdef12", "iops-abcdef23"})
	})

	t.Run("GetIOPSError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetIOPSTier("iops-abcdef12").Return(ecloud.IOPSTier{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving IOPS tier [iops-abcdef12]: test error\n", func() {
			ecloudIOPSTierShow(service, &cobra.Command{}, []string{"iops-abcdef12"})
		})
	})
}
