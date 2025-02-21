package pss

import (
	"errors"
	"testing"

	"github.com/ans-group/cli/internal/pkg/clierrors"
	"github.com/ans-group/cli/test/mocks"
	"github.com/ans-group/cli/test/test_output"
	"github.com/ans-group/sdk-go/pkg/connection"
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
		service.EXPECT().GetIncidentCasesPaginated(gomock.Any()).Return(connection.NewPaginated(&connection.APIResponseBodyData[[]pss.IncidentCase]{Data: []pss.IncidentCase{}}, *connection.NewAPIRequestParameters(), nil), nil).Times(1)

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

		service.EXPECT().GetIncidentCasesPaginated(gomock.Any()).Return(connection.NewPaginated(&connection.APIResponseBodyData[[]pss.IncidentCase]{Data: []pss.IncidentCase{}}, *connection.NewAPIRequestParameters(), nil), errors.New("test error")).Times(1)

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

func Test_pssIncidentCreate(t *testing.T) {
	setFlags := func(cmd *cobra.Command) {
		cmd.Flags().Set("title", "test incident")
		cmd.Flags().Set("description", "test description")
		cmd.Flags().Set("type", string(pss.IncidentCaseTypeServiceRequest))
		cmd.Flags().Set("category", "04f48547-96ee-4c49-901f-875a72396a60")
		cmd.Flags().Set("supported-service", "5dda44db-dd06-466b-85f6-14669d471bfd")
		cmd.Flags().Set("impact", string(pss.IncidentCaseImpactMinor))
	}

	t.Run("DefaultCreate", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)
		cmd := pssIncidentCreateCmd(nil)
		setFlags(cmd)

		gomock.InOrder(
			service.EXPECT().CreateIncidentCase(gomock.Any()).Do(func(req pss.CreateIncidentCaseRequest) {
				assert.Equal(t, "test incident", req.Title)
				assert.Equal(t, "test description", req.Description)
				assert.Equal(t, pss.IncidentCaseTypeServiceRequest, req.Type)
				assert.Equal(t, "04f48547-96ee-4c49-901f-875a72396a60", req.CategoryID)
				assert.Equal(t, "5dda44db-dd06-466b-85f6-14669d471bfd", req.SupportedServiceID)
				assert.Equal(t, pss.IncidentCaseImpactMinor, req.Impact)
			}).Return("INC123456", nil),
			service.EXPECT().GetIncidentCase("INC123456").Return(pss.IncidentCase{}, nil),
		)

		pssIncidentCreate(service, cmd, []string{})
	})

	t.Run("InvalidType_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)
		cmd := pssIncidentCreateCmd(nil)
		setFlags(cmd)
		cmd.Flags().Set("type", "invalid")

		err := pssIncidentCreate(service, cmd, []string{})
		assert.Contains(t, err.Error(), "Invalid pss.IncidentCaseType")
	})

	t.Run("InvalidImpact_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)
		cmd := pssIncidentCreateCmd(nil)
		setFlags(cmd)
		cmd.Flags().Set("impact", "invalid")

		err := pssIncidentCreate(service, cmd, []string{})
		assert.Contains(t, err.Error(), "Invalid pss.IncidentCaseImpact")
	})

	t.Run("CreateIncidentError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)
		cmd := pssIncidentCreateCmd(nil)
		setFlags(cmd)

		service.EXPECT().CreateIncidentCase(gomock.Any()).Return("", errors.New("test error"))

		err := pssIncidentCreate(service, cmd, []string{})
		assert.Equal(t, "Error creating incident: test error", err.Error())
	})

	t.Run("GetIncidentError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)
		cmd := pssIncidentCreateCmd(nil)
		setFlags(cmd)

		gomock.InOrder(
			service.EXPECT().CreateIncidentCase(gomock.Any()).Return("INC123456", nil),
			service.EXPECT().GetIncidentCase("INC123456").Return(pss.IncidentCase{}, errors.New("test error")),
		)

		err := pssIncidentCreate(service, cmd, []string{})
		assert.Equal(t, "Error retrieving new incident: test error", err.Error())
	})
}

func Test_pssIncidentCloseCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := pssIncidentCloseCmd(nil).Args(nil, []string{"INC123456"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := pssIncidentCloseCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing incident", err.Error())
	})
}

func Test_pssIncidentClose(t *testing.T) {
	t.Run("DefaultClose", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)
		cmd := pssIncidentCloseCmd(nil)
		cmd.Flags().Set("reason", "test reason")
		cmd.Flags().Set("contact", "123")

		gomock.InOrder(
			service.EXPECT().CloseIncidentCase("INC123456", gomock.Any()).Do(func(incidentID string, req pss.CloseIncidentCaseRequest) {
				assert.Equal(t, "test reason", req.Reason)
				assert.Equal(t, 123, req.ContactID)
			}).Return("INC123456", nil),
			service.EXPECT().GetIncidentCase("INC123456").Return(pss.IncidentCase{}, nil),
		)

		pssIncidentClose(service, cmd, []string{"INC123456"})
	})

	t.Run("CloseIncidentCaseError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)
		cmd := pssIncidentCloseCmd(nil)

		service.EXPECT().CloseIncidentCase("INC123456", gomock.Any()).Return("", errors.New("test error"))

		test_output.AssertErrorOutput(t, "Failed to close incident [INC123456]: test error\n", func() {
			pssIncidentClose(service, cmd, []string{"INC123456"})
		})
	})

	t.Run("GetIncidentCaseError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)
		cmd := pssIncidentCloseCmd(nil)

		gomock.InOrder(
			service.EXPECT().CloseIncidentCase("INC123456", gomock.Any()).Return("INC123456", nil),
			service.EXPECT().GetIncidentCase("INC123456").Return(pss.IncidentCase{}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error retrieving closed incident [INC123456]: test error\n", func() {
			pssIncidentClose(service, cmd, []string{"INC123456"})
		})
	})
}
