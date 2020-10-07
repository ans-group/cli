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
	"github.com/ukfast/sdk-go/pkg/ptr"
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

		assert.Equal(t, "Error retrieving VPCs: test error", err.Error())
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

		test_output.AssertErrorOutput(t, "Error retrieving VPC [vpc-abcdef12]: test error\n", func() {
			ecloudVPCShow(service, &cobra.Command{}, []string{"vpc-abcdef12"})
		})
	})
}

func Test_ecloudVPCCreate(t *testing.T) {
	t.Run("DefaultCreate", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudVPCCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testvpc", "--router=rtr-abcdef12"})

		req := ecloud.CreateVPCRequest{
			Name: ptr.String("testvpc"),
		}

		gomock.InOrder(
			service.EXPECT().CreateVPC(req).Return("vpc-abcdef12", nil),
			service.EXPECT().GetVPC("vpc-abcdef12").Return(ecloud.VPC{}, nil),
		)

		ecloudVPCCreate(service, cmd, []string{})
	})

	t.Run("CreateVPCError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudVPCCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testvpc", "--router=rtr-abcdef12"})

		service.EXPECT().CreateVPC(gomock.Any()).Return("", errors.New("test error")).Times(1)

		err := ecloudVPCCreate(service, cmd, []string{})

		assert.Equal(t, "Error creating VPC: test error", err.Error())
	})

	t.Run("GetVPCError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudVPCCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testvpc", "--router=rtr-abcdef12"})

		gomock.InOrder(
			service.EXPECT().CreateVPC(gomock.Any()).Return("vpc-abcdef12", nil),
			service.EXPECT().GetVPC("vpc-abcdef12").Return(ecloud.VPC{}, errors.New("test error")),
		)

		err := ecloudVPCCreate(service, cmd, []string{})

		assert.Equal(t, "Error retrieving new VPC: test error", err.Error())
	})
}

func Test_ecloudVPCUpdateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudVPCUpdateCmd(nil).Args(nil, []string{"vpc-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudVPCUpdateCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing vpc", err.Error())
	})
}

func Test_ecloudVPCUpdate(t *testing.T) {
	t.Run("SingleVPC", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudVPCCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testvpc"})

		req := ecloud.PatchVPCRequest{
			Name: ptr.String("testvpc"),
		}

		gomock.InOrder(
			service.EXPECT().PatchVPC("vpc-abcdef12", req).Return(nil),
			service.EXPECT().GetVPC("vpc-abcdef12").Return(ecloud.VPC{}, nil),
		)

		ecloudVPCUpdate(service, cmd, []string{"vpc-abcdef12"})
	})

	t.Run("MultipleVPCs", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().PatchVPC("vpc-abcdef12", gomock.Any()).Return(nil),
			service.EXPECT().GetVPC("vpc-abcdef12").Return(ecloud.VPC{}, nil),
			service.EXPECT().PatchVPC("vpc-12abcdef", gomock.Any()).Return(nil),
			service.EXPECT().GetVPC("vpc-12abcdef").Return(ecloud.VPC{}, nil),
		)

		ecloudVPCUpdate(service, &cobra.Command{}, []string{"vpc-abcdef12", "vpc-12abcdef"})
	})

	t.Run("PatchVPCError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().PatchVPC("vpc-abcdef12", gomock.Any()).Return(errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error updating VPC [vpc-abcdef12]: test error\n", func() {
			ecloudVPCUpdate(service, &cobra.Command{}, []string{"vpc-abcdef12"})
		})
	})

	t.Run("GetVPCError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().PatchVPC("vpc-abcdef12", gomock.Any()).Return(nil),
			service.EXPECT().GetVPC("vpc-abcdef12").Return(ecloud.VPC{}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error retrieving updated VPC [vpc-abcdef12]: test error\n", func() {
			ecloudVPCUpdate(service, &cobra.Command{}, []string{"vpc-abcdef12"})
		})
	})
}

func Test_ecloudVPCDeleteCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudVPCDeleteCmd(nil).Args(nil, []string{"vpc-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudVPCDeleteCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing vpc", err.Error())
	})
}

func Test_ecloudVPCDelete(t *testing.T) {
	t.Run("SingleVPC", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().DeleteVPC("vpc-abcdef12").Return(nil).Times(1)

		ecloudVPCDelete(service, &cobra.Command{}, []string{"vpc-abcdef12"})
	})

	t.Run("MultipleVPCs", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().DeleteVPC("vpc-abcdef12").Return(nil),
			service.EXPECT().DeleteVPC("vpc-12abcdef").Return(nil),
		)

		ecloudVPCDelete(service, &cobra.Command{}, []string{"vpc-abcdef12", "vpc-12abcdef"})
	})

	t.Run("DeleteVPCError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().DeleteVPC("vpc-abcdef12").Return(errors.New("test error")).Times(1)

		test_output.AssertErrorOutput(t, "Error removing VPC [vpc-abcdef12]: test error\n", func() {
			ecloudVPCDelete(service, &cobra.Command{}, []string{"vpc-abcdef12"})
		})
	})
}
