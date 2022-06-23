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

func Test_draasBillingTypeList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDRaaSService(mockCtrl)

		service.EXPECT().GetBillingTypes(gomock.Any()).Return([]draas.BillingType{}, nil).Times(1)

		draasBillingTypeList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDRaaSService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := draasBillingTypeList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetBillingTypesError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDRaaSService(mockCtrl)

		service.EXPECT().GetBillingTypes(gomock.Any()).Return([]draas.BillingType{}, errors.New("test error")).Times(1)

		err := draasBillingTypeList(service, &cobra.Command{}, []string{})

		assert.Equal(t, "Error retrieving billing types: test error", err.Error())
	})
}

func Test_draasBillingTypeShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := draasBillingTypeShowCmd(nil).Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := draasBillingTypeShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing billing type", err.Error())
	})
}

func Test_draasBillingTypeShow(t *testing.T) {
	t.Run("SingleBillingType", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDRaaSService(mockCtrl)

		service.EXPECT().GetBillingType("00000000-0000-0000-0000-000000000000").Return(draas.BillingType{}, nil).Times(1)

		draasBillingTypeShow(service, &cobra.Command{}, []string{"00000000-0000-0000-0000-000000000000"})
	})

	t.Run("GetBillingTypeError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDRaaSService(mockCtrl)

		service.EXPECT().GetBillingType("00000000-0000-0000-0000-000000000000").Return(draas.BillingType{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving billing type [00000000-0000-0000-0000-000000000000]: test error\n", func() {
			draasBillingTypeShow(service, &cobra.Command{}, []string{"00000000-0000-0000-0000-000000000000"})
		})
	})
}
