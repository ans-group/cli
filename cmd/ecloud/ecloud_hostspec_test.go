package ecloud

import (
	"errors"
	"testing"

	"github.com/ans-group/cli/internal/pkg/clierrors"
	"github.com/ans-group/cli/test/mocks"
	"github.com/ans-group/cli/test/test_output"
	"github.com/ans-group/sdk-go/pkg/service/ecloud"
	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_ecloudHostSpecList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetHostSpecs(gomock.Any()).Return([]ecloud.HostSpec{}, nil).Times(1)

		ecloudHostSpecList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := ecloudHostSpecList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetHostSpecsError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetHostSpecs(gomock.Any()).Return([]ecloud.HostSpec{}, errors.New("test error")).Times(1)

		err := ecloudHostSpecList(service, &cobra.Command{}, []string{})

		assert.Equal(t, "error retrieving host specs: test error", err.Error())
	})
}

func Test_ecloudHostSpecShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudHostSpecShowCmd(nil).Args(nil, []string{"hs-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudHostSpecShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "missing host spec", err.Error())
	})
}

func Test_ecloudHostSpecShow(t *testing.T) {
	t.Run("SingleHostSpec", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetHostSpec("hs-abcdef12").Return(ecloud.HostSpec{}, nil).Times(1)

		ecloudHostSpecShow(service, &cobra.Command{}, []string{"hs-abcdef12"})
	})

	t.Run("MultipleHostSpecs", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetHostSpec("hs-abcdef12").Return(ecloud.HostSpec{}, nil),
			service.EXPECT().GetHostSpec("hs-abcdef23").Return(ecloud.HostSpec{}, nil),
		)

		ecloudHostSpecShow(service, &cobra.Command{}, []string{"hs-abcdef12", "hs-abcdef23"})
	})

	t.Run("GetHostSpecError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetHostSpec("hs-abcdef12").Return(ecloud.HostSpec{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving host spec [hs-abcdef12]: test error\n", func() {
			ecloudHostSpecShow(service, &cobra.Command{}, []string{"hs-abcdef12"})
		})
	})
}
