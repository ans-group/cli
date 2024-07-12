package pss

import (
	"errors"
	"testing"

	"github.com/ans-group/cli/internal/pkg/clierrors"
	"github.com/ans-group/cli/test/mocks"
	"github.com/ans-group/cli/test/test_output"
	"github.com/ans-group/sdk-go/pkg/service/pss"
	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_pssIncidentList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)

		service.EXPECT().GetIncidentCases(gomock.Any()).Return([]pss.IncidentCase{}, nil).Times(1)

		pssIncidentList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := pssIncidentList(service, cmd, []string{})
		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetIncidentsError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)

		service.EXPECT().GetIncidentCases(gomock.Any()).Return([]pss.IncidentCase{}, errors.New("test error")).Times(1)

		err := pssIncidentList(service, &cobra.Command{}, []string{})
		assert.Equal(t, "test error", err.Error())
	})
}

func Test_pssIncidentShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := pssIncidentShowCmd(nil).Args(nil, []string{"INC123456"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := pssIncidentShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing incident", err.Error())
	})
}

func Test_pssIncidentShow(t *testing.T) {
	t.Run("SingleIncident", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)

		service.EXPECT().GetIncidentCase("INC123456").Return(pss.IncidentCase{}, nil)

		pssIncidentShow(service, &cobra.Command{}, []string{"INC123456"})
	})

	t.Run("MultipleIncidents", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetIncidentCase("INC123").Return(pss.IncidentCase{}, nil),
			service.EXPECT().GetIncidentCase("INC456").Return(pss.IncidentCase{}, nil),
		)

		pssIncidentShow(service, &cobra.Command{}, []string{"INC123", "INC456"})
	})

	t.Run("GetIncidentError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)

		service.EXPECT().GetIncidentCase("INC123456").Return(pss.IncidentCase{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving incident [INC123456]: test error\n", func() {
			pssIncidentShow(service, &cobra.Command{}, []string{"INC123456"})
		})
	})
}

// func Test_pssIncidentCreate(t *testing.T) {
// 	t.Run("DefaultCreate", func(t *testing.T) {
// 		mockCtrl := gomock.NewController(t)
// 		defer mockCtrl.Finish()

// 		service := mocks.NewMockPSSService(mockCtrl)
// 		cmd := pssIncidentCreateCmd(nil)
// 		cmd.Flags().Set("subject", "test subject")
// 		cmd.Flags().Set("product-id", "456")
// 		cmd.Flags().Set("product-name", "testname")
// 		cmd.Flags().Set("product-type", "testtype")

// 		gomock.InOrder(
// 			service.EXPECT().CreateIncident(gomock.Any()).Do(func(req pss.CreateIncidentIncident) {
// 				assert.Equal(t, "test subject", req.Subject)
// 				assert.Equal(t, 456, req.Product.ID)
// 				assert.Equal(t, "testname", req.Product.Name)
// 				assert.Equal(t, "testtype", req.Product.Type)
// 			}).Return(123, nil),
// 			service.EXPECT().GetIncident(123).Return(pss.Incident{}, nil),
// 		)

// 		pssIncidentCreate(service, cmd, []string{})
// 	})

// 	t.Run("InvalidPriority_ReturnsError", func(t *testing.T) {
// 		mockCtrl := gomock.NewController(t)
// 		defer mockCtrl.Finish()

// 		service := mocks.NewMockPSSService(mockCtrl)
// 		cmd := pssIncidentCreateCmd(nil)
// 		cmd.Flags().Set("priority", "invalid")

// 		err := pssIncidentCreate(service, cmd, []string{})
// 		assert.Contains(t, err.Error(), "Invalid pss.IncidentPriority")
// 	})

// 	t.Run("CreateIncidentError_ReturnsError", func(t *testing.T) {
// 		mockCtrl := gomock.NewController(t)
// 		defer mockCtrl.Finish()

// 		service := mocks.NewMockPSSService(mockCtrl)
// 		cmd := pssIncidentCreateCmd(nil)

// 		service.EXPECT().CreateIncident(gomock.Any()).Return(0, errors.New("test error")).Times(1)

// 		err := pssIncidentCreate(service, cmd, []string{})
// 		assert.Equal(t, "Error creating incident: test error", err.Error())
// 	})

// 	t.Run("GetIncidentError_ReturnsError", func(t *testing.T) {
// 		mockCtrl := gomock.NewController(t)
// 		defer mockCtrl.Finish()

// 		service := mocks.NewMockPSSService(mockCtrl)
// 		cmd := pssIncidentCreateCmd(nil)

// 		gomock.InOrder(
// 			service.EXPECT().CreateIncident(gomock.Any()).Return(123, nil),
// 			service.EXPECT().GetIncident(123).Return(pss.Incident{}, errors.New("test error")),
// 		)

// 		err := pssIncidentCreate(service, cmd, []string{})
// 		assert.Equal(t, "Error retrieving new incident: test error", err.Error())
// 	})
// }

// func Test_pssIncidentUpdateCmd_Args(t *testing.T) {
// 	t.Run("ValidArgs_NoError", func(t *testing.T) {
// 		err := pssIncidentUpdateCmd(nil).Args(nil, []string{"123"})

// 		assert.Nil(t, err)
// 	})

// 	t.Run("InvalidArgs_Error", func(t *testing.T) {
// 		err := pssIncidentUpdateCmd(nil).Args(nil, []string{})

// 		assert.NotNil(t, err)
// 		assert.Equal(t, "Missing incident", err.Error())
// 	})
// }

// func Test_pssIncidentUpdate(t *testing.T) {
// 	t.Run("DefaultUpdate", func(t *testing.T) {
// 		mockCtrl := gomock.NewController(t)
// 		defer mockCtrl.Finish()

// 		service := mocks.NewMockPSSService(mockCtrl)
// 		cmd := pssIncidentUpdateCmd(nil)
// 		cmd.Flags().Set("secure", "true")
// 		cmd.Flags().Set("read", "true")
// 		cmd.Flags().Set("incident-sms", "true")
// 		cmd.Flags().Set("archived", "true")
// 		cmd.Flags().Set("priority", "High")

// 		gomock.InOrder(
// 			service.EXPECT().PatchIncident(123, gomock.Any()).Do(func(incidentID int, req pss.PatchIncidentIncident) {
// 				assert.Equal(t, 123, incidentID)
// 				assert.Equal(t, true, *req.Secure)
// 				assert.Equal(t, true, *req.Read)
// 				assert.Equal(t, true, *req.IncidentSMS)
// 				assert.Equal(t, true, *req.Archived)
// 				assert.Equal(t, pss.IncidentPriorityHigh, req.Priority)
// 			}).Return(nil),
// 			service.EXPECT().GetIncident(123).Return(pss.Incident{}, nil),
// 		)

// 		pssIncidentUpdate(service, cmd, []string{"123"})
// 	})

// 	t.Run("InvalidPriority_ReturnsError", func(t *testing.T) {
// 		mockCtrl := gomock.NewController(t)
// 		defer mockCtrl.Finish()

// 		service := mocks.NewMockPSSService(mockCtrl)
// 		cmd := pssIncidentUpdateCmd(nil)
// 		cmd.Flags().Set("priority", "invalid")

// 		err := pssIncidentUpdate(service, cmd, []string{"123"})
// 		assert.Contains(t, err.Error(), "Invalid pss.IncidentPriority")
// 	})

// 	t.Run("InvalidIncidentID_OutputsError", func(t *testing.T) {
// 		mockCtrl := gomock.NewController(t)
// 		defer mockCtrl.Finish()

// 		service := mocks.NewMockPSSService(mockCtrl)
// 		cmd := pssIncidentUpdateCmd(nil)

// 		test_output.AssertErrorOutput(t, "Invalid incident ID [abc]\n", func() {
// 			pssIncidentUpdate(service, cmd, []string{"abc"})
// 		})
// 	})

// 	t.Run("PatchIncidentError_OutputsError", func(t *testing.T) {
// 		mockCtrl := gomock.NewController(t)
// 		defer mockCtrl.Finish()

// 		service := mocks.NewMockPSSService(mockCtrl)
// 		cmd := pssIncidentUpdateCmd(nil)

// 		service.EXPECT().PatchIncident(123, gomock.Any()).Return(errors.New("test error")).Times(1)

// 		test_output.AssertErrorOutput(t, "Error updating incident [123]: test error\n", func() {
// 			pssIncidentUpdate(service, cmd, []string{"123"})
// 		})
// 	})

// 	t.Run("GetIncidentError_OutputsError", func(t *testing.T) {
// 		mockCtrl := gomock.NewController(t)
// 		defer mockCtrl.Finish()

// 		service := mocks.NewMockPSSService(mockCtrl)
// 		cmd := pssIncidentUpdateCmd(nil)

// 		gomock.InOrder(
// 			service.EXPECT().PatchIncident(123, gomock.Any()).Return(nil),
// 			service.EXPECT().GetIncident(123).Return(pss.Incident{}, errors.New("test error")),
// 		)

// 		test_output.AssertErrorOutput(t, "Error retrieving updated incident [123]: test error\n", func() {
// 			pssIncidentUpdate(service, cmd, []string{"123"})
// 		})
// 	})
// }

// func Test_pssIncidentCloseCmd_Args(t *testing.T) {
// 	t.Run("ValidArgs_NoError", func(t *testing.T) {
// 		err := pssIncidentCloseCmd(nil).Args(nil, []string{"123"})

// 		assert.Nil(t, err)
// 	})

// 	t.Run("InvalidArgs_Error", func(t *testing.T) {
// 		err := pssIncidentCloseCmd(nil).Args(nil, []string{})

// 		assert.NotNil(t, err)
// 		assert.Equal(t, "Missing incident", err.Error())
// 	})
// }

// func Test_pssIncidentClose(t *testing.T) {
// 	t.Run("DefaultClose", func(t *testing.T) {
// 		mockCtrl := gomock.NewController(t)
// 		defer mockCtrl.Finish()

// 		service := mocks.NewMockPSSService(mockCtrl)

// 		gomock.InOrder(
// 			service.EXPECT().PatchIncident(123, gomock.Any()).Return(nil),
// 			service.EXPECT().GetIncident(123).Return(pss.Incident{}, nil),
// 		)

// 		pssIncidentClose(service, pssIncidentCloseCmd(nil), []string{"123"})
// 	})

// 	t.Run("PatchIncidentError_OutputsError", func(t *testing.T) {
// 		mockCtrl := gomock.NewController(t)
// 		defer mockCtrl.Finish()

// 		service := mocks.NewMockPSSService(mockCtrl)
// 		cmd := pssIncidentCloseCmd(nil)

// 		service.EXPECT().PatchIncident(123, gomock.Any()).Return(errors.New("test error")).Times(1)

// 		test_output.AssertErrorOutput(t, "Error closing incident [123]: test error\n", func() {
// 			pssIncidentClose(service, cmd, []string{"123"})
// 		})
// 	})

// 	t.Run("GetIncidentError_OutputsError", func(t *testing.T) {
// 		mockCtrl := gomock.NewController(t)
// 		defer mockCtrl.Finish()

// 		service := mocks.NewMockPSSService(mockCtrl)
// 		cmd := pssIncidentCloseCmd(nil)

// 		gomock.InOrder(
// 			service.EXPECT().PatchIncident(123, gomock.Any()).Return(nil),
// 			service.EXPECT().GetIncident(123).Return(pss.Incident{}, errors.New("test error")),
// 		)

// 		test_output.AssertErrorOutput(t, "Error retrieving updated incident [123]: test error\n", func() {
// 			pssIncidentClose(service, cmd, []string{"123"})
// 		})
// 	})
// }
