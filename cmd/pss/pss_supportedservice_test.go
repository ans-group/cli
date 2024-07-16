package pss

import (
	"errors"
	"testing"

	"github.com/ans-group/cli/internal/pkg/clierrors"
	"github.com/ans-group/cli/test/mocks"
	"github.com/ans-group/sdk-go/pkg/service/pss"
	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_pssSupportedServiceList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)

		service.EXPECT().GetSupportedServices(gomock.Any()).Return([]pss.SupportedService{}, nil)

		pssSupportedServiceList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := pssSupportedServiceList(service, cmd, []string{})
		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetSupportedServicesError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)

		service.EXPECT().GetSupportedServices(gomock.Any()).Return([]pss.SupportedService{}, errors.New("test error"))

		err := pssSupportedServiceList(service, &cobra.Command{}, []string{})
		assert.Equal(t, "test error", err.Error())
	})
}
