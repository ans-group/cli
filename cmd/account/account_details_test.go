package account

import (
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/ukfast/cli/test/mocks"
	"github.com/ukfast/sdk-go/pkg/service/account"
	"gopkg.in/go-playground/assert.v1"
)

func Test_accountDetailShow(t *testing.T) {
	t.Run("SingleDetail", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockAccountService(mockCtrl)

		service.EXPECT().GetDetails().Return(account.Details{}, nil).Times(1)

		accountDetailsShow(service, &cobra.Command{}, []string{""})
	})

	t.Run("GetDetailsError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockAccountService(mockCtrl)

		service.EXPECT().GetDetails().Return(account.Details{}, errors.New("test error"))

		err := accountDetailsShow(service, &cobra.Command{}, []string{""})

		assert.Equal(t, "Error retrieving details: test error", err.Error())
	})
}
