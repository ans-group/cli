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

func Test_accountContactList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockAccountService(mockCtrl)

		service.EXPECT().GetContacts(gomock.Any()).Return([]account.Contact{}, nil).Times(1)

		accountContactList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockAccountService(mockCtrl)
		cmd := accountContactListCmd(nil)
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := accountContactList(service, cmd, []string{})

		assert.NotNil(t, err)
		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetContactsError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockAccountService(mockCtrl)

		service.EXPECT().GetContacts(gomock.Any()).Return([]account.Contact{}, errors.New("test error")).Times(1)

		err := accountContactList(service, &cobra.Command{}, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "error retrieving contacts: test error", err.Error())
	})
}

func Test_accountContactShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := accountContactShowCmd(nil).Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := accountContactShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "missing contact", err.Error())
	})
}

func Test_accountContactShow(t *testing.T) {
	t.Run("SingleContact", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockAccountService(mockCtrl)

		service.EXPECT().GetContact(123).Return(account.Contact{}, nil).Times(1)

		accountContactShow(service, &cobra.Command{}, []string{"123"})
	})

	t.Run("MultipleContacts", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockAccountService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetContact(123).Return(account.Contact{}, nil),
			service.EXPECT().GetContact(456).Return(account.Contact{}, nil),
		)

		accountContactShow(service, &cobra.Command{}, []string{"123", "456"})
	})

	t.Run("GetContactID_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockAccountService(mockCtrl)

		test_output.AssertErrorOutput(t, "Invalid contact ID [abc]\n", func() {
			accountContactShow(service, &cobra.Command{}, []string{"abc"})
		})
	})

	t.Run("GetContactError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockAccountService(mockCtrl)

		service.EXPECT().GetContact(123).Return(account.Contact{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving contact [123]: test error\n", func() {
			accountContactShow(service, &cobra.Command{}, []string{"123"})
		})
	})
}
