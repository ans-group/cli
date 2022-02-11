package managedcloudflare

import (
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/ukfast/cli/test/mocks"
	"github.com/ukfast/sdk-go/pkg/service/managedcloudflare"
)

func Test_managedcloudflareAccountMemberCreateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := managedcloudflareAccountMemberCreateCmd(nil).Args(nil, []string{"00000000-0000-0000-0000-000000000000"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := managedcloudflareAccountMemberCreateCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing account", err.Error())
	})
}

func Test_managedcloudflareAccountMemberCreate(t *testing.T) {
	t.Run("DefaultCreate", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockManagedCloudflareService(mockCtrl)
		cmd := managedcloudflareAccountMemberCreateCmd(nil)
		cmd.ParseFlags([]string{"--email-address=test@test.com"})

		req := managedcloudflare.CreateAccountMemberRequest{
			EmailAddress: "test@test.com",
		}

		service.EXPECT().CreateAccountMember("00000000-0000-0000-0000-000000000000", req).Return(nil)

		managedcloudflareAccountMemberCreate(service, cmd, []string{"00000000-0000-0000-0000-000000000000"})
	})

	t.Run("CreateAccountMemberError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockManagedCloudflareService(mockCtrl)
		cmd := managedcloudflareAccountMemberCreateCmd(nil)
		cmd.ParseFlags([]string{"--email-address=test@test.com"})

		service.EXPECT().CreateAccountMember("00000000-0000-0000-0000-000000000000", gomock.Any()).Return(errors.New("test error"))

		err := managedcloudflareAccountMemberCreate(service, cmd, []string{"00000000-0000-0000-0000-000000000000"})

		assert.Equal(t, "Error creating account member: test error", err.Error())
	})
}
