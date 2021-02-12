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

func Test_ecloudNICList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetNICs(gomock.Any()).Return([]ecloud.NIC{}, nil).Times(1)

		ecloudNICList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := ecloudNICList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetNICsError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetNICs(gomock.Any()).Return([]ecloud.NIC{}, errors.New("test error")).Times(1)

		err := ecloudNICList(service, &cobra.Command{}, []string{})

		assert.Equal(t, "Error retrieving NICs: test error", err.Error())
	})
}

func Test_ecloudNICShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudNICShowCmd(nil).Args(nil, []string{"nic-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudNICShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing nic", err.Error())
	})
}

func Test_ecloudNICShow(t *testing.T) {
	t.Run("SingleNIC", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetNIC("nic-abcdef12").Return(ecloud.NIC{}, nil).Times(1)

		ecloudNICShow(service, &cobra.Command{}, []string{"nic-abcdef12"})
	})

	t.Run("MultipleNICs", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetNIC("nic-abcdef12").Return(ecloud.NIC{}, nil),
			service.EXPECT().GetNIC("nic-abcdef23").Return(ecloud.NIC{}, nil),
		)

		ecloudNICShow(service, &cobra.Command{}, []string{"nic-abcdef12", "nic-abcdef23"})
	})

	t.Run("GetNICError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetNIC("nic-abcdef12").Return(ecloud.NIC{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving NIC [nic-abcdef12]: test error\n", func() {
			ecloudNICShow(service, &cobra.Command{}, []string{"nic-abcdef12"})
		})
	})
}
