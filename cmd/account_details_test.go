package cmd

import (
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/ukfast/cli/test/mocks"
	"github.com/ukfast/cli/test/test_output"
	"github.com/ukfast/sdk-go/pkg/service/account"
)

func Test_accountDetailShow(t *testing.T) {
	t.Run("SingleDetail", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockAccountService(mockCtrl)

		service.EXPECT().GetDetails().Return(account.Details{}, nil).Times(1)

		accountDetailsShow(service, &cobra.Command{}, []string{""})
	})

	t.Run("GetDetailsError_OutputsFatal", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockAccountService(mockCtrl)

		service.EXPECT().GetDetails().Return(account.Details{}, errors.New("test error"))

		test_output.AssertFatalOutput(t, "Error retrieving details: test error\n", func() {
			accountDetailsShow(service, &cobra.Command{}, []string{""})
		})
	})
}
