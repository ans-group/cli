package cmd

import (
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/ukfast/cli/test/mocks"
	"github.com/ukfast/sdk-go/pkg/service/ltaas"
)

func Test_loadtestScenarioList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLTaaSService(mockCtrl)

		service.EXPECT().GetScenarios(gomock.Any()).Return([]ltaas.Scenario{}, nil).Times(1)

		loadtestScenarioList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_OutputsFatal", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLTaaSService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := loadtestScenarioList(service, cmd, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing value for filtering", err.Error())
	})

	t.Run("GetScenariosError_OutputsFatal", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLTaaSService(mockCtrl)

		service.EXPECT().GetScenarios(gomock.Any()).Return([]ltaas.Scenario{}, errors.New("test error")).Times(1)

		err := loadtestScenarioList(service, &cobra.Command{}, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Error retrieving scenarios: test error", err.Error())
	})
}
