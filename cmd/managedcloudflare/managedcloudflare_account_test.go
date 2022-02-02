package managedcloudflare

import (
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/ukfast/cli/internal/pkg/clierrors"
	"github.com/ukfast/cli/test/mocks"
	"github.com/ukfast/cli/test/test_output"
	"github.com/ukfast/sdk-go/pkg/service/managedcloudflare"
)

func Test_managedcloudflareAccountList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockManagedCloudflareService(mockCtrl)

		service.EXPECT().GetAccounts(gomock.Any()).Return([]managedcloudflare.Account{}, nil).Times(1)

		managedcloudflareAccountList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockManagedCloudflareService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := managedcloudflareAccountList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetAccountsError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockManagedCloudflareService(mockCtrl)

		service.EXPECT().GetAccounts(gomock.Any()).Return([]managedcloudflare.Account{}, errors.New("test error")).Times(1)

		err := managedcloudflareAccountList(service, &cobra.Command{}, []string{})

		assert.Equal(t, "Error retrieving accounts: test error", err.Error())
	})
}

func Test_managedcloudflareAccountShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := managedcloudflareAccountShowCmd(nil).Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := managedcloudflareAccountShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing account", err.Error())
	})
}

func Test_managedcloudflareAccountShow(t *testing.T) {
	t.Run("SingleAccount", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockManagedCloudflareService(mockCtrl)

		service.EXPECT().GetAccount(123).Return(managedcloudflare.Account{}, nil).Times(1)

		managedcloudflareAccountShow(service, &cobra.Command{}, []string{"123"})
	})

	t.Run("MultipleAccounts", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockManagedCloudflareService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetAccount(123).Return(managedcloudflare.Account{}, nil),
			service.EXPECT().GetAccount(456).Return(managedcloudflare.Account{}, nil),
		)

		managedcloudflareAccountShow(service, &cobra.Command{}, []string{"123", "456"})
	})

	t.Run("GetAccountID_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockManagedCloudflareService(mockCtrl)

		test_output.AssertErrorOutput(t, "Invalid account ID [abc]\n", func() {
			managedcloudflareAccountShow(service, &cobra.Command{}, []string{"abc"})
		})
	})

	t.Run("GetAccountError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockManagedCloudflareService(mockCtrl)

		service.EXPECT().GetAccount(123).Return(managedcloudflare.Account{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving account [123]: test error\n", func() {
			managedcloudflareAccountShow(service, &cobra.Command{}, []string{"123"})
		})
	})
}
