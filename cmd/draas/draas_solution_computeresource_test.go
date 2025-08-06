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

func Test_draasSolutionComputeResourceListCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		cmd := draasSolutionComputeResourceListCmd(nil)
		err := cmd.Args(nil, []string{"00000000-0000-0000-0000-000000000000"})

		assert.Nil(t, err)
	})

	t.Run("MissingSolution_Error", func(t *testing.T) {
		cmd := draasSolutionComputeResourceListCmd(nil)
		err := cmd.Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "missing solution", err.Error())
	})
}

func Test_draasSolutionComputeResourceList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDRaaSService(mockCtrl)

		service.EXPECT().GetSolutionComputeResources("00000000-0000-0000-0000-000000000000", gomock.Any()).Return([]draas.ComputeResource{}, nil).Times(1)

		draasSolutionComputeResourceList(service, &cobra.Command{}, []string{"00000000-0000-0000-0000-000000000000"})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDRaaSService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := draasSolutionComputeResourceList(service, cmd, []string{"00000000-0000-0000-0000-000000000000"})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetSolutionComputeResourcesError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDRaaSService(mockCtrl)

		service.EXPECT().GetSolutionComputeResources("00000000-0000-0000-0000-000000000000", gomock.Any()).Return([]draas.ComputeResource{}, errors.New("test error")).Times(1)

		err := draasSolutionComputeResourceList(service, &cobra.Command{}, []string{"00000000-0000-0000-0000-000000000000"})

		assert.Equal(t, "error retrieving solution compute resources: test error", err.Error())
	})
}

func Test_draasSolutionComputeResourceShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		cmd := draasSolutionComputeResourceShowCmd(nil)
		err := cmd.Args(nil, []string{"00000000-0000-0000-0000-000000000000", "00000000-0000-0000-0000-000000000001"})

		assert.Nil(t, err)
	})

	t.Run("MissingSolution_Error", func(t *testing.T) {
		cmd := draasSolutionComputeResourceShowCmd(nil)
		err := cmd.Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "missing solution", err.Error())
	})

	t.Run("MissingComputeResource_Error", func(t *testing.T) {
		cmd := draasSolutionComputeResourceShowCmd(nil)
		err := cmd.Args(nil, []string{"00000000-0000-0000-0000-000000000000"})

		assert.NotNil(t, err)
		assert.Equal(t, "missing compute resource", err.Error())
	})
}

func Test_draasSolutionComputeResourceShow(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDRaaSService(mockCtrl)

		service.EXPECT().GetSolutionComputeResource("00000000-0000-0000-0000-000000000000", "00000000-0000-0000-0000-000000000001").Return(draas.ComputeResource{}, nil)

		draasSolutionComputeResourceShow(service, &cobra.Command{}, []string{"00000000-0000-0000-0000-000000000000", "00000000-0000-0000-0000-000000000001"})
	})

	t.Run("GetSolutionComputeResource_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDRaaSService(mockCtrl)

		service.EXPECT().GetSolutionComputeResource("00000000-0000-0000-0000-000000000000", "00000000-0000-0000-0000-000000000001").Return(draas.ComputeResource{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving solution compute resource [00000000-0000-0000-0000-000000000001]: test error\n", func() {
			draasSolutionComputeResourceShow(service, &cobra.Command{}, []string{"00000000-0000-0000-0000-000000000000", "00000000-0000-0000-0000-000000000001"})
		})
	})
}
