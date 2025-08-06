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

func Test_ecloudSSHKeyPairList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetSSHKeyPairs(gomock.Any()).Return([]ecloud.SSHKeyPair{}, nil).Times(1)

		ecloudSSHKeyPairList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := ecloudSSHKeyPairList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetSSHKeyPairsError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetSSHKeyPairs(gomock.Any()).Return([]ecloud.SSHKeyPair{}, errors.New("test error")).Times(1)

		err := ecloudSSHKeyPairList(service, &cobra.Command{}, []string{})

		assert.Equal(t, "error retrieving SSH key pairs: test error", err.Error())
	})
}

func Test_ecloudSSHKeyPairShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudSSHKeyPairShowCmd(nil).Args(nil, []string{"ssh-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudSSHKeyPairShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "missing SSH key pair", err.Error())
	})
}

func Test_ecloudSSHKeyPairShow(t *testing.T) {
	t.Run("SingleSSHKeyPair", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetSSHKeyPair("ssh-abcdef12").Return(ecloud.SSHKeyPair{}, nil).Times(1)

		ecloudSSHKeyPairShow(service, &cobra.Command{}, []string{"ssh-abcdef12"})
	})

	t.Run("MultipleSSHKeyPairs", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetSSHKeyPair("ssh-abcdef12").Return(ecloud.SSHKeyPair{}, nil),
			service.EXPECT().GetSSHKeyPair("ssh-abcdef23").Return(ecloud.SSHKeyPair{}, nil),
		)

		ecloudSSHKeyPairShow(service, &cobra.Command{}, []string{"ssh-abcdef12", "ssh-abcdef23"})
	})

	t.Run("GetSSHKeyPairError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetSSHKeyPair("ssh-abcdef12").Return(ecloud.SSHKeyPair{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving SSH key pair [ssh-abcdef12]: test error\n", func() {
			ecloudSSHKeyPairShow(service, &cobra.Command{}, []string{"ssh-abcdef12"})
		})
	})
}

func Test_ecloudSSHKeyPairCreate(t *testing.T) {
	t.Run("DefaultCreate", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudSSHKeyPairCreateCmd(nil, nil)
		cmd.ParseFlags([]string{"--name=testkeypair"})

		req := ecloud.CreateSSHKeyPairRequest{
			Name: "testkeypair",
		}

		gomock.InOrder(
			service.EXPECT().CreateSSHKeyPair(req).Return("ssh-abcdef12", nil),
			service.EXPECT().GetSSHKeyPair("ssh-abcdef12").Return(ecloud.SSHKeyPair{}, nil),
		)

		ecloudSSHKeyPairCreate(service, nil, cmd, []string{})
	})

	t.Run("CreateSSHKeyPairError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudSSHKeyPairCreateCmd(nil, nil)
		cmd.ParseFlags([]string{"--name=testkeypair"})

		service.EXPECT().CreateSSHKeyPair(gomock.Any()).Return("", errors.New("test error")).Times(1)

		err := ecloudSSHKeyPairCreate(service, nil, cmd, []string{})

		assert.Equal(t, "error creating SSH key pair: test error", err.Error())
	})

	t.Run("GetSSHKeyPairError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().CreateSSHKeyPair(gomock.Any()).Return("ssh-abcdef12", nil),
			service.EXPECT().GetSSHKeyPair("ssh-abcdef12").Return(ecloud.SSHKeyPair{}, errors.New("test error")),
		)

		err := ecloudSSHKeyPairCreate(service, nil, ecloudSSHKeyPairCreateCmd(nil, nil), []string{})

		assert.Equal(t, "error retrieving new SSH key pair: test error", err.Error())
	})
}

func Test_ecloudSSHKeyPairUpdateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudSSHKeyPairUpdateCmd(nil, nil).Args(nil, []string{"ssh-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudSSHKeyPairUpdateCmd(nil, nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "missing SSH key pair", err.Error())
	})
}

func Test_ecloudSSHKeyPairUpdate(t *testing.T) {
	t.Run("SingleSSHKeyPair", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudSSHKeyPairUpdateCmd(nil, nil)
		cmd.ParseFlags([]string{"--name=testkeypair"})

		req := ecloud.PatchSSHKeyPairRequest{
			Name: "testkeypair",
		}

		gomock.InOrder(
			service.EXPECT().PatchSSHKeyPair("ssh-abcdef12", req).Return(nil),
			service.EXPECT().GetSSHKeyPair("ssh-abcdef12").Return(ecloud.SSHKeyPair{}, nil),
		)

		ecloudSSHKeyPairUpdate(service, nil, cmd, []string{"ssh-abcdef12"})
	})

	t.Run("MultipleSSHKeyPairs", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().PatchSSHKeyPair("ssh-abcdef12", gomock.Any()).Return(nil),
			service.EXPECT().GetSSHKeyPair("ssh-abcdef12").Return(ecloud.SSHKeyPair{}, nil),
			service.EXPECT().PatchSSHKeyPair("ssh-12abcdef", gomock.Any()).Return(nil),
			service.EXPECT().GetSSHKeyPair("ssh-12abcdef").Return(ecloud.SSHKeyPair{}, nil),
		)

		ecloudSSHKeyPairUpdate(service, nil, &cobra.Command{}, []string{"ssh-abcdef12", "ssh-12abcdef"})
	})

	t.Run("PatchSSHKeyPairError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().PatchSSHKeyPair("ssh-abcdef12", gomock.Any()).Return(errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error updating SSH key pair [ssh-abcdef12]: test error\n", func() {
			ecloudSSHKeyPairUpdate(service, nil, &cobra.Command{}, []string{"ssh-abcdef12"})
		})
	})

	t.Run("GetSSHKeyPairError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().PatchSSHKeyPair("ssh-abcdef12", gomock.Any()).Return(nil),
			service.EXPECT().GetSSHKeyPair("ssh-abcdef12").Return(ecloud.SSHKeyPair{}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error retrieving updated SSH key pair [ssh-abcdef12]: test error\n", func() {
			ecloudSSHKeyPairUpdate(service, nil, &cobra.Command{}, []string{"ssh-abcdef12"})
		})
	})
}

func Test_ecloudSSHKeyPairDeleteCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudSSHKeyPairDeleteCmd(nil).Args(nil, []string{"ssh-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudSSHKeyPairDeleteCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "missing SSH key pair", err.Error())
	})
}

func Test_ecloudSSHKeyPairDelete(t *testing.T) {
	t.Run("SingleSSHKeyPair", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().DeleteSSHKeyPair("ssh-abcdef12").Return(nil).Times(1)

		ecloudSSHKeyPairDelete(service, &cobra.Command{}, []string{"ssh-abcdef12"})
	})

	t.Run("MultipleSSHKeyPairs", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().DeleteSSHKeyPair("ssh-abcdef12").Return(nil),
			service.EXPECT().DeleteSSHKeyPair("ssh-12abcdef").Return(nil),
		)

		ecloudSSHKeyPairDelete(service, &cobra.Command{}, []string{"ssh-abcdef12", "ssh-12abcdef"})
	})

	t.Run("DeleteSSHKeyPairError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().DeleteSSHKeyPair("ssh-abcdef12").Return(errors.New("test error")).Times(1)

		test_output.AssertErrorOutput(t, "Error removing SSH key pair [ssh-abcdef12]: test error\n", func() {
			ecloudSSHKeyPairDelete(service, &cobra.Command{}, []string{"ssh-abcdef12"})
		})
	})
}
