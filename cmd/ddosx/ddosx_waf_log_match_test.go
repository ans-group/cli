package ddosx

import (
	"errors"
	"testing"

	"github.com/ans-group/cli/test/mocks"
	"github.com/ans-group/cli/test/test_output"
	"github.com/ans-group/sdk-go/pkg/connection"
	"github.com/ans-group/sdk-go/pkg/service/ddosx"
	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_ddosxWAFLogMatchList(t *testing.T) {
	t.Run("WithoutRequest_ExpectedCalls", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().GetWAFLogMatchesPaginated(gomock.Any()).Return(connection.NewPaginated(&connection.APIResponseBodyData[[]ddosx.WAFLogMatch]{}, connection.APIRequestParameters{}, nil), nil).Times(1)

		ddosxWAFLogMatchList(service, ddosxWAFLogMatchListCmd(nil), []string{})
	})

	t.Run("WithRequest_ExpectedCalls", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		cmd := ddosxWAFLogMatchListCmd(nil)
		cmd.Flags().Set("log", "abcdef")

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().GetWAFLogRequestMatchesPaginated("abcdef", gomock.Any()).Return(connection.NewPaginated(&connection.APIResponseBodyData[[]ddosx.WAFLogMatch]{}, connection.APIRequestParameters{}, nil), nil).Times(1)

		ddosxWAFLogMatchList(service, cmd, []string{})
	})

	t.Run("GetWAFLogRequestMatchesPaginatedError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		cmd := ddosxWAFLogMatchListCmd(nil)
		cmd.Flags().Set("log", "abcdef")

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().GetWAFLogRequestMatchesPaginated("abcdef", gomock.Any()).Return(&connection.Paginated[ddosx.WAFLogMatch]{}, errors.New("test error"))

		err := ddosxWAFLogMatchList(service, cmd, []string{})

		assert.Equal(t, "Error retrieving WAF log matches: test error", err.Error())
	})

	t.Run("GetWAFLogMatchesPaginatedError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().GetWAFLogMatchesPaginated(gomock.Any()).Return(&connection.Paginated[ddosx.WAFLogMatch]{}, errors.New("test error"))

		err := ddosxWAFLogMatchList(service, ddosxWAFLogMatchListCmd(nil), []string{})

		assert.Equal(t, "Error retrieving WAF log matches: test error", err.Error())
	})
}

func Test_ddosxWAFLogMatchShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ddosxWAFLogMatchShowCmd(nil).Args(nil, []string{"2d8556677081cecf112b555c359a78c6", "abcdef"})

		assert.Nil(t, err)
	})

	t.Run("MissingRequest_Error", func(t *testing.T) {
		err := ddosxWAFLogMatchShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing log", err.Error())
	})

	t.Run("MissingMatch_Error", func(t *testing.T) {
		err := ddosxWAFLogMatchShowCmd(nil).Args(nil, []string{"2d8556677081cecf112b555c359a78c6"})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing match", err.Error())
	})
}

func Test_ddosxWAFLogMatchShow(t *testing.T) {
	t.Run("SingleLog", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().GetWAFLogRequestMatch("2d8556677081cecf112b555c359a78c6", "abcdef").Return(ddosx.WAFLogMatch{}, nil).Times(1)

		ddosxWAFLogMatchShow(service, &cobra.Command{}, []string{"2d8556677081cecf112b555c359a78c6", "abcdef"})
	})

	t.Run("GetWAFLogMatchError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().GetWAFLogRequestMatch("2d8556677081cecf112b555c359a78c6", "abcdef").Return(ddosx.WAFLogMatch{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving WAF log matches [abcdef]: test error\n", func() {
			ddosxWAFLogMatchShow(service, &cobra.Command{}, []string{"2d8556677081cecf112b555c359a78c6", "abcdef"})
		})
	})
}
