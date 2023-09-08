package account

import (
	"errors"
	"testing"

	"github.com/ans-group/cli/internal/pkg/clierrors"
	"github.com/ans-group/cli/test/mocks"
	"github.com/ans-group/cli/test/test_output"
	"github.com/ans-group/sdk-go/pkg/service/account"
	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_accountClientList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockAccountService(mockCtrl)

		service.EXPECT().GetClients(gomock.Any()).Return([]account.Client{}, nil).Times(1)

		accountClientList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockAccountService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := accountClientList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetClientsError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockAccountService(mockCtrl)

		service.EXPECT().GetClients(gomock.Any()).Return([]account.Client{}, errors.New("test error")).Times(1)

		err := accountClientList(service, &cobra.Command{}, []string{})

		assert.Equal(t, "Error retrieving clients: test error", err.Error())
	})
}

func Test_accountClientShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := accountClientShowCmd(nil).Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := accountClientShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing client", err.Error())
	})
}

func Test_accountClientShow(t *testing.T) {
	t.Run("SingleClient", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockAccountService(mockCtrl)

		service.EXPECT().GetClient(123).Return(account.Client{}, nil).Times(1)

		accountClientShow(service, &cobra.Command{}, []string{"123"})
	})

	t.Run("MultipleClients", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockAccountService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetClient(123).Return(account.Client{}, nil),
			service.EXPECT().GetClient(456).Return(account.Client{}, nil),
		)

		accountClientShow(service, &cobra.Command{}, []string{"123", "456"})
	})

	t.Run("GetClientID_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockAccountService(mockCtrl)

		test_output.AssertErrorOutput(t, "Invalid client ID [abc]\n", func() {
			accountClientShow(service, &cobra.Command{}, []string{"abc"})
		})
	})

	t.Run("GetClientError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockAccountService(mockCtrl)

		service.EXPECT().GetClient(123).Return(account.Client{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving client [123]: test error\n", func() {
			accountClientShow(service, &cobra.Command{}, []string{"123"})
		})
	})
}

func Test_accountClientCreate(t *testing.T) {
	t.Run("DefaultCreate", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockAccountService(mockCtrl)
		cmd := accountClientCreateCmd(nil)
		cmd.Flags().Set("company-name", "test client 1")

		expectedReq := account.CreateClientRequest{
			CompanyName: "test client 1",
		}

		gomock.InOrder(
			service.EXPECT().CreateClient(expectedReq).Return(123, nil),
			service.EXPECT().GetClient(123).Return(account.Client{}, nil),
		)

		accountClientCreate(service, cmd, []string{})
	})

	t.Run("CreateClientError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockAccountService(mockCtrl)

		service.EXPECT().CreateClient(gomock.Any()).Return(0, errors.New("test error")).Times(1)

		err := accountClientCreate(service, &cobra.Command{}, []string{})

		assert.Equal(t, "Error creating client: test error", err.Error())
	})

	t.Run("GetClientError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockAccountService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().CreateClient(gomock.Any()).Return(123, nil),
			service.EXPECT().GetClient(123).Return(account.Client{}, errors.New("test error")),
		)

		err := accountClientCreate(service, &cobra.Command{}, []string{})

		assert.Equal(t, "Error retrieving new client: test error", err.Error())
	})
}

func Test_accountClientUpdateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := accountClientUpdateCmd(nil).Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := accountClientUpdateCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing client", err.Error())
	})
}

func Test_accountClientUpdate(t *testing.T) {
	t.Run("SingleClient_SetsCPU", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockAccountService(mockCtrl)
		cmd := accountClientUpdateCmd(nil)
		cmd.Flags().Set("company-name", "test")

		expectedPatch := account.PatchClientRequest{
			CompanyName: "test",
		}

		gomock.InOrder(
			service.EXPECT().PatchClient(123, expectedPatch).Return(nil),
			service.EXPECT().GetClient(123).Return(account.Client{}, nil),
		)

		accountClientUpdate(service, cmd, []string{"123"})
	})

	t.Run("MultipleClients", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockAccountService(mockCtrl)
		cmd := accountClientUpdateCmd(nil)
		cmd.Flags().Set("company-name", "test")

		expectedPatch := account.PatchClientRequest{
			CompanyName: "test",
		}

		gomock.InOrder(
			service.EXPECT().PatchClient(123, gomock.Eq(expectedPatch)).Return(nil),
			service.EXPECT().GetClient(123).Return(account.Client{}, nil),
			service.EXPECT().PatchClient(456, gomock.Eq(expectedPatch)).Return(nil),
			service.EXPECT().GetClient(456).Return(account.Client{}, nil),
		)

		accountClientUpdate(service, cmd, []string{"123", "456"})
	})

	t.Run("InvalidClientID_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockAccountService(mockCtrl)

		test_output.AssertErrorOutput(t, "Invalid client ID [abc]\n", func() {
			accountClientUpdate(service, &cobra.Command{}, []string{"abc"})
		})
	})

	t.Run("PatchClientError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockAccountService(mockCtrl)

		service.EXPECT().PatchClient(123, gomock.Any()).Return(errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error updating client [123]: test error\n", func() {
			accountClientUpdate(service, &cobra.Command{}, []string{"123"})
		})
	})

	t.Run("GetClientError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockAccountService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().PatchClient(123, gomock.Any()).Return(nil),
			service.EXPECT().GetClient(123).Return(account.Client{}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error retrieving updated client [123]: test error\n", func() {
			accountClientUpdate(service, &cobra.Command{}, []string{"123"})
		})
	})
}

func Test_accountClientDeleteCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := accountClientDeleteCmd(nil).Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := accountClientDeleteCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing client", err.Error())
	})
}

func Test_accountClientDelete(t *testing.T) {
	t.Run("SingleClient", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockAccountService(mockCtrl)

		service.EXPECT().DeleteClient(123).Return(nil).Times(1)

		accountClientDelete(service, &cobra.Command{}, []string{"123"})
	})

	t.Run("MultipleClients", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockAccountService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().DeleteClient(123).Return(nil),
			service.EXPECT().DeleteClient(456).Return(nil),
		)

		accountClientDelete(service, &cobra.Command{}, []string{"123", "456"})
	})

	t.Run("InvalidClientID_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockAccountService(mockCtrl)

		test_output.AssertErrorOutput(t, "Invalid client ID [abc]\n", func() {
			accountClientDelete(service, &cobra.Command{}, []string{"abc"})
		})
	})

	t.Run("GetClientError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockAccountService(mockCtrl)

		service.EXPECT().DeleteClient(123).Return(errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error removing client [123]: test error\n", func() {
			accountClientDelete(service, &cobra.Command{}, []string{"123"})
		})
	})
}
