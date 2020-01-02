package cmd

import (
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/ukfast/cli/test/mocks"
	"github.com/ukfast/cli/test/test_output"
	"github.com/ukfast/sdk-go/pkg/service/ltaas"
)

func Test_loadtestThresholdList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLTaaSService(mockCtrl)

		service.EXPECT().GetThresholds(gomock.Any()).Return([]ltaas.Threshold{}, nil).Times(1)

		loadtestThresholdList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_OutputsFatal", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLTaaSService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := loadtestThresholdList(service, cmd, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing value for filtering", err.Error())
	})

	t.Run("GetThresholdsError_OutputsFatal", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLTaaSService(mockCtrl)

		service.EXPECT().GetThresholds(gomock.Any()).Return([]ltaas.Threshold{}, errors.New("test error")).Times(1)

		err := loadtestThresholdList(service, &cobra.Command{}, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Error retrieving thresholds: test error", err.Error())
	})
}

func Test_loadtestThresholdShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := loadtestThresholdShowCmd().Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := loadtestThresholdShowCmd().Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing threshold", err.Error())
	})
}

func Test_loadtestThresholdShow(t *testing.T) {
	t.Run("SingleThreshold", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLTaaSService(mockCtrl)

		service.EXPECT().GetThreshold("00000000-0000-0000-0000-000000000000").Return(ltaas.Threshold{}, nil).Times(1)

		loadtestThresholdShow(service, &cobra.Command{}, []string{"00000000-0000-0000-0000-000000000000"})
	})

	t.Run("MultipleThresholds", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLTaaSService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetThreshold("00000000-0000-0000-0000-000000000000").Return(ltaas.Threshold{}, nil),
			service.EXPECT().GetThreshold("00000000-0000-0000-0000-000000000001").Return(ltaas.Threshold{}, nil),
		)

		loadtestThresholdShow(service, &cobra.Command{}, []string{"00000000-0000-0000-0000-000000000000", "00000000-0000-0000-0000-000000000001"})
	})

	t.Run("GetThresholdError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLTaaSService(mockCtrl)

		service.EXPECT().GetThreshold("00000000-0000-0000-0000-000000000000").Return(ltaas.Threshold{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving threshold [00000000-0000-0000-0000-000000000000]: test error\n", func() {
			loadtestThresholdShow(service, &cobra.Command{}, []string{"00000000-0000-0000-0000-000000000000"})
		})
	})
}
