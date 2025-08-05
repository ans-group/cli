package cloudflare

import (
	"errors"
	"testing"

	"github.com/ans-group/cli/test/mocks"
	"github.com/ans-group/sdk-go/pkg/service/cloudflare"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_cloudflareAccountMemberCreateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := cloudflareAccountMemberCreateCmd(nil).Args(nil, []string{"00000000-0000-0000-0000-000000000000"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := cloudflareAccountMemberCreateCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "missing account", err.Error())
	})
}

func Test_cloudflareAccountMemberCreate(t *testing.T) {
	t.Run("DefaultCreate", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockCloudflareService(mockCtrl)
		cmd := cloudflareAccountMemberCreateCmd(nil)
		cmd.ParseFlags([]string{"--email-address=test@test.com"})

		req := cloudflare.CreateAccountMemberRequest{
			EmailAddress: "test@test.com",
		}

		service.EXPECT().CreateAccountMember("00000000-0000-0000-0000-000000000000", req).Return(nil)

		cloudflareAccountMemberCreate(service, cmd, []string{"00000000-0000-0000-0000-000000000000"})
	})

	t.Run("CreateAccountMemberError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockCloudflareService(mockCtrl)
		cmd := cloudflareAccountMemberCreateCmd(nil)
		cmd.ParseFlags([]string{"--email-address=test@test.com"})

		service.EXPECT().CreateAccountMember("00000000-0000-0000-0000-000000000000", gomock.Any()).Return(errors.New("test error"))

		err := cloudflareAccountMemberCreate(service, cmd, []string{"00000000-0000-0000-0000-000000000000"})

		assert.Equal(t, "error creating account member: test error", err.Error())
	})
}
