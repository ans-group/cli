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

func Test_ecloudVPCList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetVPCs(gomock.Any()).Return([]ecloud.VPC{}, nil).Times(1)

		ecloudVPCList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := ecloudVPCList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetVPCsError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetVPCs(gomock.Any()).Return([]ecloud.VPC{}, errors.New("test error")).Times(1)

		err := ecloudVPCList(service, &cobra.Command{}, []string{})

		assert.Equal(t, "Error retrieving vpcs: test error", err.Error())
	})
}

func Test_ecloudVPCShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudVPCShowCmd(nil).Args(nil, []string{"vpc-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudVPCShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing vpc", err.Error())
	})
}

func Test_ecloudVPCShow(t *testing.T) {
	t.Run("SingleVPC", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetVPC("vpc-abcdef12").Return(ecloud.VPC{}, nil).Times(1)

		ecloudVPCShow(service, &cobra.Command{}, []string{"vpc-abcdef12"})
	})

	t.Run("MultipleVPCs", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetVPC("vpc-abcdef12").Return(ecloud.VPC{}, nil),
			service.EXPECT().GetVPC("vpc-abcdef23").Return(ecloud.VPC{}, nil),
		)

		ecloudVPCShow(service, &cobra.Command{}, []string{"vpc-abcdef12", "vpc-abcdef23"})
	})

	t.Run("GetVPCError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetVPC("vpc-abcdef12").Return(ecloud.VPC{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving vpc [vpc-abcdef12]: test error\n", func() {
			ecloudVPCShow(service, &cobra.Command{}, []string{"vpc-abcdef12"})
		})
	})
}
