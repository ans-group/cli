package threatmonitoring

import (
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/ukfast/cli/internal/pkg/clierrors"
	"github.com/ukfast/cli/test/mocks"
	"github.com/ukfast/cli/test/test_output"
	"github.com/ukfast/sdk-go/pkg/service/threatmonitoring"
)

func Test_threatmonitoringAgentList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockThreatMonitoringService(mockCtrl)

		service.EXPECT().GetAgents(gomock.Any()).Return([]threatmonitoring.Agent{}, nil).Times(1)

		threatmonitoringAgentList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockThreatMonitoringService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := threatmonitoringAgentList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetAgentsError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockThreatMonitoringService(mockCtrl)

		service.EXPECT().GetAgents(gomock.Any()).Return([]threatmonitoring.Agent{}, errors.New("test error")).Times(1)

		err := threatmonitoringAgentList(service, &cobra.Command{}, []string{})

		assert.Equal(t, "Error retrieving agents: test error", err.Error())
	})
}

func Test_threatmonitoringAgentShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := threatmonitoringAgentShowCmd(nil).Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := threatmonitoringAgentShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing agent", err.Error())
	})
}

func Test_threatmonitoringAgentShow(t *testing.T) {
	t.Run("SingleAgent", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockThreatMonitoringService(mockCtrl)

		service.EXPECT().GetAgent(123).Return(threatmonitoring.Agent{}, nil).Times(1)

		threatmonitoringAgentShow(service, &cobra.Command{}, []string{"123"})
	})

	t.Run("MultipleAgents", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockThreatMonitoringService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetAgent(123).Return(threatmonitoring.Agent{}, nil),
			service.EXPECT().GetAgent(456).Return(threatmonitoring.Agent{}, nil),
		)

		threatmonitoringAgentShow(service, &cobra.Command{}, []string{"123", "456"})
	})

	t.Run("GetAgentID_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockThreatMonitoringService(mockCtrl)

		test_output.AssertErrorOutput(t, "Invalid agent ID [abc]\n", func() {
			threatmonitoringAgentShow(service, &cobra.Command{}, []string{"abc"})
		})
	})

	t.Run("GetAgentError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockThreatMonitoringService(mockCtrl)

		service.EXPECT().GetAgent(123).Return(threatmonitoring.Agent{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving agent [123]: test error\n", func() {
			threatmonitoringAgentShow(service, &cobra.Command{}, []string{"123"})
		})
	})
}
