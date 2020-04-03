package ecloud

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

func Test_ecloudPodList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetPods(gomock.Any()).Return([]ecloud.Pod{}, nil).Times(1)

		ecloudPodList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := ecloudPodList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetPodsError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetPods(gomock.Any()).Return([]ecloud.Pod{}, errors.New("test error")).Times(1)

		err := ecloudPodList(service, &cobra.Command{}, []string{})

		assert.Equal(t, "Error retrieving pods: test error", err.Error())
	})
}

func Test_ecloudPodShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudPodShowCmd(nil).Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudPodShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing pod", err.Error())
	})
}

func Test_ecloudPodShow(t *testing.T) {
	t.Run("SinglePod", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetPod(123).Return(ecloud.Pod{}, nil).Times(1)

		ecloudPodShow(service, &cobra.Command{}, []string{"123"})
	})

	t.Run("MultiplePods", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetPod(123).Return(ecloud.Pod{}, nil),
			service.EXPECT().GetPod(456).Return(ecloud.Pod{}, nil),
		)

		ecloudPodShow(service, &cobra.Command{}, []string{"123", "456"})
	})

	t.Run("GetPodID_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		test_output.AssertErrorOutput(t, "Invalid pod ID [abc]\n", func() {
			ecloudPodShow(service, &cobra.Command{}, []string{"abc"})
		})
	})

	t.Run("GetPodError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetPod(123).Return(ecloud.Pod{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving pod [123]: test error\n", func() {
			ecloudPodShow(service, &cobra.Command{}, []string{"123"})
		})
	})
}
