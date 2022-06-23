package account

import (
	"errors"
	"testing"

	"github.com/ans-group/cli/test/mocks"
	"github.com/ans-group/sdk-go/pkg/service/account"
	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
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
