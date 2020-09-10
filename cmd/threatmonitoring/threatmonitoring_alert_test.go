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
	"github.com/ukfast/sdk-go/pkg/connection"
	"github.com/ukfast/sdk-go/pkg/service/threatmonitoring"
)

func Test_threatmonitoringAlertList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockThreatMonitoringService(mockCtrl)

		service.EXPECT().GetAlertsPaginated(gomock.Any()).Return(&threatmonitoring.PaginatedAlert{PaginatedBase: &connection.PaginatedBase{}}, nil).Times(1)

		threatmonitoringAlertList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockThreatMonitoringService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := threatmonitoringAlertList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetAlertsError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockThreatMonitoringService(mockCtrl)

		service.EXPECT().GetAlertsPaginated(gomock.Any()).Return(&threatmonitoring.PaginatedAlert{PaginatedBase: &connection.PaginatedBase{}}, errors.New("test error")).Times(1)

		err := threatmonitoringAlertList(service, &cobra.Command{}, []string{})

		assert.Equal(t, "Error retrieving alerts: test error", err.Error())
	})
}

func Test_threatmonitoringAlertShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := threatmonitoringAlertShowCmd(nil).Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := threatmonitoringAlertShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing alert", err.Error())
	})
}

func Test_threatmonitoringAlertShow(t *testing.T) {
	t.Run("SingleAlert", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockThreatMonitoringService(mockCtrl)

		service.EXPECT().GetAlert("123").Return(threatmonitoring.Alert{}, nil).Times(1)

		threatmonitoringAlertShow(service, &cobra.Command{}, []string{"123"})
	})

	t.Run("MultipleAlerts", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockThreatMonitoringService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetAlert("123").Return(threatmonitoring.Alert{}, nil),
			service.EXPECT().GetAlert("456").Return(threatmonitoring.Alert{}, nil),
		)

		threatmonitoringAlertShow(service, &cobra.Command{}, []string{"123", "456"})
	})

	t.Run("GetAlertError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockThreatMonitoringService(mockCtrl)

		service.EXPECT().GetAlert("123").Return(threatmonitoring.Alert{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving alert [123]: test error\n", func() {
			threatmonitoringAlertShow(service, &cobra.Command{}, []string{"123"})
		})
	})
}
