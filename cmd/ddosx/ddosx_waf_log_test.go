package ddosx

import (
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/ukfast/cli/test/mocks"
	"github.com/ukfast/cli/test/test_output"
	"github.com/ukfast/sdk-go/pkg/connection"
	"github.com/ukfast/sdk-go/pkg/service/ddosx"
)

func Test_ddosxWAFLogList(t *testing.T) {
	t.Run("NoError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().GetWAFLogs(gomock.Any()).Return([]ddosx.WAFLog{}, nil).Times(1)

		ddosxWAFLogList(service, &cobra.Command{}, []string{})
	})

	t.Run("WithDomainFilter_ExpectedFiltering", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		filtering := connection.NewAPIRequestParameters().
			WithFilter(connection.APIRequestFiltering{
				Property: "domain",
				Operator: connection.EQOperator,
				Value:    []string{"example.com"},
			})

		cmd := ddosxWAFLogListCmd(nil)
		cmd.Flags().Set("domain", "example.com")

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().GetWAFLogs(gomock.Eq(*filtering)).Return([]ddosx.WAFLog{}, nil).Times(1)

		ddosxWAFLogList(service, cmd, []string{})
	})

	t.Run("GetWAFLogsError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().GetWAFLogs(gomock.Any()).Return([]ddosx.WAFLog{}, errors.New("test error"))

		err := ddosxWAFLogList(service, &cobra.Command{}, []string{})

		assert.Equal(t, "Error retrieving WAF logs: test error", err.Error())
	})
}

func Test_ddosxWAFLogShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ddosxWAFLogShowCmd(nil).Args(nil, []string{"2d8556677081cecf112b555c359a78c6"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ddosxWAFLogShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing request", err.Error())
	})
}

func Test_ddosxWAFLogShow(t *testing.T) {
	t.Run("SingleLog", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().GetWAFLog("2d8556677081cecf112b555c359a78c6").Return(ddosx.WAFLog{}, nil).Times(1)

		ddosxWAFLogShow(service, &cobra.Command{}, []string{"2d8556677081cecf112b555c359a78c6"})
	})

	t.Run("GetWAFLogError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().GetWAFLog("2d8556677081cecf112b555c359a78c6").Return(ddosx.WAFLog{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving WAF log [2d8556677081cecf112b555c359a78c6]: test error\n", func() {
			ddosxWAFLogShow(service, &cobra.Command{}, []string{"2d8556677081cecf112b555c359a78c6"})
		})
	})
}
