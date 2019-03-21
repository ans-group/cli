package cmd

import (
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/cli/test"
	"github.com/ukfast/cli/test/mocks"
	"github.com/ukfast/sdk-go/pkg/service/account"
)

func Test_accountContactList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockAccountService(mockCtrl)

		service.EXPECT().GetContacts(gomock.Any()).Return([]account.Contact{}, nil).Times(1)

		accountContactList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_OutputsFatal", func(t *testing.T) {
		code := 0
		oldOutputExit := output.SetOutputExit(func(c int) {
			code = c
		})
		defer func() { output.SetOutputExit(oldOutputExit) }()
		defer func() { flagFilter = nil }()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockAccountService(mockCtrl)
		flagFilter = []string{"invalidfilter"}

		output := test.CatchStdErr(t, func() {
			accountContactList(service, &cobra.Command{}, []string{})
		})

		assert.Equal(t, 1, code)
		assert.Equal(t, "Missing value for filtering\n", output)
	})

	t.Run("GetContactsError_OutputsFatal", func(t *testing.T) {
		code := 0
		oldOutputExit := output.SetOutputExit(func(c int) {
			code = c
		})
		defer func() { output.SetOutputExit(oldOutputExit) }()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockAccountService(mockCtrl)

		service.EXPECT().GetContacts(gomock.Any()).Return([]account.Contact{}, errors.New("test error")).Times(1)

		output := test.CatchStdErr(t, func() {
			accountContactList(service, &cobra.Command{}, []string{})
		})

		assert.Equal(t, 1, code)
		assert.Equal(t, "Error retrieving contacts: test error\n", output)
	})
}

func Test_accountContactShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := accountContactShowCmd().Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := accountContactShowCmd().Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing contact", err.Error())
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

		output := test.CatchStdErr(t, func() {
			accountContactShow(service, &cobra.Command{}, []string{"abc"})
		})

		assert.Equal(t, "Invalid contact ID [abc]\n", output)
	})

	t.Run("GetContactError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockAccountService(mockCtrl)

		service.EXPECT().GetContact(123).Return(account.Contact{}, errors.New("test error"))

		output := test.CatchStdErr(t, func() {
			accountContactShow(service, &cobra.Command{}, []string{"123"})
		})

		assert.Equal(t, "Error retrieving contact [123]: test error\n", output)
	})
}
