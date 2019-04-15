package cmd

import (
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/ukfast/cli/test/mocks"
	"github.com/ukfast/cli/test/test_output"
	"github.com/ukfast/sdk-go/pkg/service/pss"
)

func Test_pssRequestList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)

		service.EXPECT().GetRequests(gomock.Any()).Return([]pss.Request{}, nil).Times(1)

		pssRequestList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_OutputsFatal", func(t *testing.T) {
		defer func() { flagFilter = nil }()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)
		flagFilter = []string{"invalidfilter"}

		test_output.AssertFatalOutput(t, "Missing value for filtering\n", func() {
			pssRequestList(service, &cobra.Command{}, []string{})
		})
	})

	t.Run("GetRequestsError_OutputsFatal", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)

		service.EXPECT().GetRequests(gomock.Any()).Return([]pss.Request{}, errors.New("test error")).Times(1)

		test_output.AssertFatalOutput(t, "Error retrieving requests: test error\n", func() {
			pssRequestList(service, &cobra.Command{}, []string{})
		})
	})
}

func Test_pssRequestShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := pssRequestShowCmd().Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := pssRequestShowCmd().Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing request", err.Error())
	})
}

func Test_pssRequestShow(t *testing.T) {
	t.Run("SingleRequest", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)

		service.EXPECT().GetRequest(123).Return(pss.Request{}, nil).Times(1)

		pssRequestShow(service, &cobra.Command{}, []string{"123"})
	})

	t.Run("MultipleRequests", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetRequest(123).Return(pss.Request{}, nil),
			service.EXPECT().GetRequest(456).Return(pss.Request{}, nil),
		)

		pssRequestShow(service, &cobra.Command{}, []string{"123", "456"})
	})

	t.Run("GetRequestID_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)

		test_output.AssertErrorOutput(t, "Invalid request ID [abc]\n", func() {
			pssRequestShow(service, &cobra.Command{}, []string{"abc"})
		})
	})

	t.Run("GetRequestError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)

		service.EXPECT().GetRequest(123).Return(pss.Request{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving request [123]: test error\n", func() {
			pssRequestShow(service, &cobra.Command{}, []string{"123"})
		})
	})
}
