package ecloud_v2

import (
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/ukfast/cli/internal/pkg/clierrors"
	"github.com/ukfast/cli/test/mocks"
	"github.com/ukfast/cli/test/test_output"
	"github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func Test_ecloudInstanceList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetInstances(gomock.Any()).Return([]ecloud.Instance{}, nil).Times(1)

		ecloudInstanceList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := ecloudInstanceList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetInstancesError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetInstances(gomock.Any()).Return([]ecloud.Instance{}, errors.New("test error")).Times(1)

		err := ecloudInstanceList(service, &cobra.Command{}, []string{})

		assert.Equal(t, "Error retrieving instances: test error", err.Error())
	})
}

func Test_ecloudInstanceShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudInstanceShowCmd(nil).Args(nil, []string{"i-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudInstanceShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing instance", err.Error())
	})
}

func Test_ecloudInstanceShow(t *testing.T) {
	t.Run("SingleInstance", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetInstance("i-abcdef12").Return(ecloud.Instance{}, nil).Times(1)

		ecloudInstanceShow(service, &cobra.Command{}, []string{"i-abcdef12"})
	})

	t.Run("MultipleInstances", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetInstance("i-abcdef12").Return(ecloud.Instance{}, nil),
			service.EXPECT().GetInstance("i-abcdef23").Return(ecloud.Instance{}, nil),
		)

		ecloudInstanceShow(service, &cobra.Command{}, []string{"i-abcdef12", "i-abcdef23"})
	})

	t.Run("GetInstanceError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetInstance("i-abcdef12").Return(ecloud.Instance{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving instance [i-abcdef12]: test error\n", func() {
			ecloudInstanceShow(service, &cobra.Command{}, []string{"i-abcdef12"})
		})
	})
}
