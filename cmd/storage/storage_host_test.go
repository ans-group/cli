package storage

import (
	"errors"
	"testing"

	"github.com/ans-group/cli/internal/pkg/clierrors"
	"github.com/ans-group/cli/test/mocks"
	"github.com/ans-group/cli/test/test_output"
	"github.com/ans-group/sdk-go/pkg/service/storage"
	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_storageHostList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockStorageService(mockCtrl)

		service.EXPECT().GetHosts(gomock.Any()).Return([]storage.Host{}, nil).Times(1)

		storageHostList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockStorageService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := storageHostList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetHostsError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockStorageService(mockCtrl)

		service.EXPECT().GetHosts(gomock.Any()).Return([]storage.Host{}, errors.New("test error")).Times(1)

		err := storageHostList(service, &cobra.Command{}, []string{})

		assert.Equal(t, "Error retrieving hosts: test error", err.Error())
	})
}

func Test_storageHostShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := storageHostShowCmd(nil).Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := storageHostShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing host", err.Error())
	})
}

func Test_storageHostShow(t *testing.T) {
	t.Run("SingleHost", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockStorageService(mockCtrl)

		service.EXPECT().GetHost(123).Return(storage.Host{}, nil).Times(1)

		storageHostShow(service, &cobra.Command{}, []string{"123"})
	})

	t.Run("MultipleHosts", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockStorageService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetHost(123).Return(storage.Host{}, nil),
			service.EXPECT().GetHost(456).Return(storage.Host{}, nil),
		)

		storageHostShow(service, &cobra.Command{}, []string{"123", "456"})
	})

	t.Run("GetHostID_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockStorageService(mockCtrl)

		test_output.AssertErrorOutput(t, "Invalid host ID [abc]\n", func() {
			storageHostShow(service, &cobra.Command{}, []string{"abc"})
		})
	})

	t.Run("GetHostError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockStorageService(mockCtrl)

		service.EXPECT().GetHost(123).Return(storage.Host{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving host [123]: test error\n", func() {
			storageHostShow(service, &cobra.Command{}, []string{"123"})
		})
	})
}
