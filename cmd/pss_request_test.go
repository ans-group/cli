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
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := pssRequestList(service, cmd, []string{})
		assert.Equal(t, "Missing value for filtering", err.Error())
	})

	t.Run("GetRequestsError_OutputsFatal", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)

		service.EXPECT().GetRequests(gomock.Any()).Return([]pss.Request{}, errors.New("test error")).Times(1)

		err := pssRequestList(service, &cobra.Command{}, []string{})
		assert.Equal(t, "test error", err.Error())
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

func Test_pssRequestCreate(t *testing.T) {
	t.Run("DefaultCreate", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)
		cmd := pssRequestCreateCmd()
		cmd.Flags().Set("subject", "test subject")
		cmd.Flags().Set("product-id", "456")
		cmd.Flags().Set("product-name", "testname")
		cmd.Flags().Set("product-type", "testtype")

		gomock.InOrder(
			service.EXPECT().CreateRequest(gomock.Any()).Do(func(req pss.CreateRequestRequest) {
				assert.Equal(t, "test subject", req.Subject)
				assert.Equal(t, 456, req.Product.ID)
				assert.Equal(t, "testname", req.Product.Name)
				assert.Equal(t, "testtype", req.Product.Type)
			}).Return(123, nil),
			service.EXPECT().GetRequest(123).Return(pss.Request{}, nil),
		)

		pssRequestCreate(service, cmd, []string{})
	})

	t.Run("InvalidPriority_OutputsFatal", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)
		cmd := pssRequestCreateCmd()
		cmd.Flags().Set("priority", "invalid")

		err := pssRequestCreate(service, cmd, []string{})
		assert.Contains(t, err.Error(), "Invalid pss.RequestPriority")
	})

	t.Run("CreateRequestError_OutputsFatal", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)
		cmd := pssRequestCreateCmd()

		service.EXPECT().CreateRequest(gomock.Any()).Return(0, errors.New("test error")).Times(1)

		err := pssRequestCreate(service, cmd, []string{})
		assert.Equal(t, "Error creating request: test error", err.Error())
	})

	t.Run("GetRequestError_OutputsFatal", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)
		cmd := pssRequestCreateCmd()

		gomock.InOrder(
			service.EXPECT().CreateRequest(gomock.Any()).Return(123, nil),
			service.EXPECT().GetRequest(123).Return(pss.Request{}, errors.New("test error")),
		)

		err := pssRequestCreate(service, cmd, []string{})
		assert.Equal(t, "Error retrieving new request: test error", err.Error())
	})
}

func Test_pssRequestUpdateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := pssRequestUpdateCmd().Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := pssRequestUpdateCmd().Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing request", err.Error())
	})
}

func Test_pssRequestUpdate(t *testing.T) {
	t.Run("DefaultUpdate", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)
		cmd := pssRequestUpdateCmd()
		cmd.Flags().Set("secure", "true")
		cmd.Flags().Set("read", "true")
		cmd.Flags().Set("request-sms", "true")
		cmd.Flags().Set("archived", "true")
		cmd.Flags().Set("priority", "High")

		gomock.InOrder(
			service.EXPECT().PatchRequest(123, gomock.Any()).Do(func(requestID int, req pss.PatchRequestRequest) {
				assert.Equal(t, 123, requestID)
				assert.Equal(t, true, *req.Secure)
				assert.Equal(t, true, *req.Read)
				assert.Equal(t, true, *req.RequestSMS)
				assert.Equal(t, true, *req.Archived)
				assert.Equal(t, pss.RequestPriorityHigh, req.Priority)
			}).Return(nil),
			service.EXPECT().GetRequest(123).Return(pss.Request{}, nil),
		)

		pssRequestUpdate(service, cmd, []string{"123"})
	})

	t.Run("InvalidPriority_OutputsFatal", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)
		cmd := pssRequestUpdateCmd()
		cmd.Flags().Set("priority", "invalid")

		err := pssRequestUpdate(service, cmd, []string{"123"})
		assert.Contains(t, err.Error(), "Invalid pss.RequestPriority")
	})

	t.Run("InvalidRequestID_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)
		cmd := pssRequestUpdateCmd()

		test_output.AssertErrorOutput(t, "Invalid request ID [abc]\n", func() {
			pssRequestUpdate(service, cmd, []string{"abc"})
		})
	})

	t.Run("PatchRequestError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)
		cmd := pssRequestUpdateCmd()

		service.EXPECT().PatchRequest(123, gomock.Any()).Return(errors.New("test error")).Times(1)

		test_output.AssertErrorOutput(t, "Error updating request [123]: test error\n", func() {
			pssRequestUpdate(service, cmd, []string{"123"})
		})
	})

	t.Run("GetRequestError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)
		cmd := pssRequestUpdateCmd()

		gomock.InOrder(
			service.EXPECT().PatchRequest(123, gomock.Any()).Return(nil),
			service.EXPECT().GetRequest(123).Return(pss.Request{}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error retrieving updated request [123]: test error\n", func() {
			pssRequestUpdate(service, cmd, []string{"123"})
		})
	})
}
