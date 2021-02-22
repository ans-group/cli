package ecloud

import (
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/ukfast/cli/internal/pkg/clierrors"
	"github.com/ukfast/cli/test/mocks"
	"github.com/ukfast/sdk-go/pkg/service/account"
)

func Test_ecloudCreditList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		service.EXPECT().GetCredits(gomock.Any()).Return([]account.Credit{}, nil).Times(1)

		err := ecloudCreditList(service, &cobra.Command{}, []string{})

		assert.Nil(t, err)
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := ecloudCreditList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetCreditsError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		service.EXPECT().GetCredits(gomock.Any()).Return([]account.Credit{}, errors.New("test error")).Times(1)

		err := ecloudCreditList(service, &cobra.Command{}, []string{})

		assert.Equal(t, "Error retrieving credits: test error", err.Error())
	})
}
