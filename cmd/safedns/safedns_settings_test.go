package safedns

import (
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/ukfast/cli/test/mocks"
	"github.com/ukfast/sdk-go/pkg/service/safedns"
)

func Test_safednsSettingsShow(t *testing.T) {
	t.Run("GetSettingsNoError_ReturnsNil", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)

		service.EXPECT().GetSettings().Return(safedns.Settings{}, nil).Times(1)

		err := safednsSettingsShow(service, &cobra.Command{}, []string{})

		assert.Nil(t, err)
	})

	t.Run("GetSettingsError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)

		service.EXPECT().GetSettings().Return(safedns.Settings{}, errors.New("test error")).Times(1)

		err := safednsSettingsShow(service, &cobra.Command{}, []string{})

		assert.NotNil(t, err)
	})
}
