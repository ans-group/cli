package draas

import (
	"errors"
	"testing"

	"github.com/ans-group/cli/internal/pkg/clierrors"
	"github.com/ans-group/cli/test/mocks"
	"github.com/ans-group/cli/test/test_output"
	"github.com/ans-group/sdk-go/pkg/service/draas"
	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_draasIOPSTierList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDRaaSService(mockCtrl)

		service.EXPECT().GetIOPSTiers(gomock.Any()).Return([]draas.IOPSTier{}, nil).Times(1)

		draasIOPSTierList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDRaaSService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := draasIOPSTierList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetIOPSTiersError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDRaaSService(mockCtrl)

		service.EXPECT().GetIOPSTiers(gomock.Any()).Return([]draas.IOPSTier{}, errors.New("test error")).Times(1)

		err := draasIOPSTierList(service, &cobra.Command{}, []string{})

		assert.Equal(t, "Error retrieving IOPS tiers: test error", err.Error())
	})
}

func Test_draasIOPSTierShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := draasIOPSTierShowCmd(nil).Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := draasIOPSTierShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing IOPS tier", err.Error())
	})
}

func Test_draasIOPSTierShow(t *testing.T) {
	t.Run("SingleIOPSTier", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDRaaSService(mockCtrl)

		service.EXPECT().GetIOPSTier("00000000-0000-0000-0000-000000000000").Return(draas.IOPSTier{}, nil).Times(1)

		draasIOPSTierShow(service, &cobra.Command{}, []string{"00000000-0000-0000-0000-000000000000"})
	})

	t.Run("GetIOPSTierError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDRaaSService(mockCtrl)

		service.EXPECT().GetIOPSTier("00000000-0000-0000-0000-000000000000").Return(draas.IOPSTier{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving IOPS tier [00000000-0000-0000-0000-000000000000]: test error\n", func() {
			draasIOPSTierShow(service, &cobra.Command{}, []string{"00000000-0000-0000-0000-000000000000"})
		})
	})
}
