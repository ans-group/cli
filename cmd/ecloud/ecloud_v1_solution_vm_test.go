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

func Test_ecloudSolutionVirtualMachineListCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudSolutionVirtualMachineListCmd(nil).Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudSolutionVirtualMachineListCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing solution", err.Error())
	})
}

func Test_ecloudSolutionVirtualMachineList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetSolutionVirtualMachines(123, gomock.Any()).Return([]ecloud.VirtualMachine{}, nil).Times(1)

		ecloudSolutionVirtualMachineList(service, &cobra.Command{}, []string{"123"})
	})

	t.Run("InvalidSolutionID_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		err := ecloudSolutionVirtualMachineList(service, &cobra.Command{}, []string{"abc"})

		assert.Equal(t, "Invalid solution ID [abc]", err.Error())
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := ecloudSolutionVirtualMachineList(service, cmd, []string{"123"})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetVirtualMachinesError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetSolutionVirtualMachines(123, gomock.Any()).Return([]ecloud.VirtualMachine{}, errors.New("test error 1")).Times(1)

		err := ecloudSolutionVirtualMachineList(service, &cobra.Command{}, []string{"123"})

		assert.Equal(t, "Error retrieving solution virtual machines: test error 1", err.Error())
	})
}
