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
		err := managedcloudflareAccountShowCmd(nil).Args(nil, []string{"00000000-0000-0000-0000-000000000000"})

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

		service.EXPECT().GetAccount("00000000-0000-0000-0000-000000000000").Return(managedcloudflare.Account{}, nil).Times(1)

		managedcloudflareAccountShow(service, &cobra.Command{}, []string{"00000000-0000-0000-0000-000000000000"})
	})

	t.Run("GetAccountError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockManagedCloudflareService(mockCtrl)

		service.EXPECT().GetAccount("00000000-0000-0000-0000-000000000000").Return(managedcloudflare.Account{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving account [00000000-0000-0000-0000-000000000000]: test error\n", func() {
			managedcloudflareAccountShow(service, &cobra.Command{}, []string{"00000000-0000-0000-0000-000000000000"})
		})
	})
}

func Test_managedcloudflareAccountCreate(t *testing.T) {
	t.Run("DefaultCreate", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockManagedCloudflareService(mockCtrl)
		cmd := managedcloudflareAccountCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testaccount"})

		req := managedcloudflare.CreateAccountRequest{
			Name: "testaccount",
		}

		gomock.InOrder(
			service.EXPECT().CreateAccount(req).Return("00000000-0000-0000-0000-000000000000", nil),
			service.EXPECT().GetAccount("00000000-0000-0000-0000-000000000000").Return(managedcloudflare.Account{}, nil),
		)

		managedcloudflareAccountCreate(service, cmd, []string{})
	})

	t.Run("CreateAccountError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockManagedCloudflareService(mockCtrl)
		cmd := managedcloudflareAccountCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testaccount"})

		service.EXPECT().CreateAccount(gomock.Any()).Return("", errors.New("test error"))

		err := managedcloudflareAccountCreate(service, cmd, []string{})

		assert.Equal(t, "Error creating account: test error", err.Error())
	})

	t.Run("GetAccountError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockManagedCloudflareService(mockCtrl)
		cmd := managedcloudflareAccountCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testaccount"})

		gomock.InOrder(
			service.EXPECT().CreateAccount(gomock.Any()).Return("00000000-0000-0000-0000-000000000000", nil),
			service.EXPECT().GetAccount("00000000-0000-0000-0000-000000000000").Return(managedcloudflare.Account{}, errors.New("test error")),
		)

		err := managedcloudflareAccountCreate(service, cmd, []string{})

		assert.Equal(t, "Error retrieving new account: test error", err.Error())
	})
}
