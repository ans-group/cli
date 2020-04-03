package loadtest

import (
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/ukfast/cli/test/mocks"
)

func Test_loadtestAccountCreate(t *testing.T) {
	t.Run("DefaultCreate", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLTaaSService(mockCtrl)
		cmd := loadtestAccountCreateCmd(nil)

		service.EXPECT().CreateAccount().Return("00000000-0000-0000-0000-000000000001", nil)

		loadtestAccountCreate(service, cmd, []string{})
	})

	t.Run("CreateAccountError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLTaaSService(mockCtrl)
		cmd := loadtestAccountCreateCmd(nil)

		service.EXPECT().CreateAccount().Return("", errors.New("test error")).Times(1)

		err := loadtestAccountCreate(service, cmd, []string{})
		assert.Equal(t, "Error creating account: test error", err.Error())
	})
}
