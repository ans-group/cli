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

func Test_pssChangeList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)
		service.EXPECT().GetChangeCasesPaginated(gomock.Any()).Return(connection.NewPaginated(&connection.APIResponseBodyData[[]pss.ChangeCase]{Data: []pss.ChangeCase{}}, *connection.NewAPIRequestParameters(), nil), nil).Times(1)

		pssChangeList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := pssChangeList(service, cmd, []string{})
		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetChangesError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)

		service.EXPECT().GetChangeCasesPaginated(gomock.Any()).Return(connection.NewPaginated(&connection.APIResponseBodyData[[]pss.ChangeCase]{Data: []pss.ChangeCase{}}, *connection.NewAPIRequestParameters(), nil), errors.New("test error")).Times(1)

		err := pssChangeList(service, &cobra.Command{}, []string{})
		assert.Equal(t, "test error", err.Error())
	})
}

func Test_pssChangeShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := pssChangeShowCmd(nil).Args(nil, []string{"CHG123456"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := pssChangeShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing change", err.Error())
	})
}

func Test_pssChangeShow(t *testing.T) {
	t.Run("SingleChange", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)

		service.EXPECT().GetChangeCase("CHG123456").Return(pss.ChangeCase{}, nil)

		pssChangeShow(service, &cobra.Command{}, []string{"CHG123456"})
	})

	t.Run("MultipleChanges", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetChangeCase("CHG123").Return(pss.ChangeCase{}, nil),
			service.EXPECT().GetChangeCase("CHG456").Return(pss.ChangeCase{}, nil),
		)

		pssChangeShow(service, &cobra.Command{}, []string{"CHG123", "CHG456"})
	})

	t.Run("GetChangeError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)

		service.EXPECT().GetChangeCase("CHG123456").Return(pss.ChangeCase{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving change [CHG123456]: test error\n", func() {
			pssChangeShow(service, &cobra.Command{}, []string{"CHG123456"})
		})
	})
}

func Test_pssChangeCreate(t *testing.T) {
	setFlags := func(cmd *cobra.Command) {
		cmd.Flags().Set("title", "test change")
		cmd.Flags().Set("description", "test description")
		cmd.Flags().Set("risk", string(pss.ChangeCaseRiskLow))
		cmd.Flags().Set("category", "04f48547-96ee-4c49-901f-875a72396a60")
		cmd.Flags().Set("supported-service", "5dda44db-dd06-466b-85f6-14669d471bfd")
		cmd.Flags().Set("impact", string(pss.ChangeCaseImpactLow))
	}

	t.Run("DefaultCreate", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)
		cmd := pssChangeCreateCmd(nil)
		setFlags(cmd)

		gomock.InOrder(
			service.EXPECT().CreateChangeCase(gomock.Any()).Do(func(req pss.CreateChangeCaseRequest) {
				assert.Equal(t, "test change", req.Title)
				assert.Equal(t, "test description", req.Description)
				assert.Equal(t, pss.ChangeCaseRiskLow, req.Risk)
				assert.Equal(t, "04f48547-96ee-4c49-901f-875a72396a60", req.CategoryID)
				assert.Equal(t, "5dda44db-dd06-466b-85f6-14669d471bfd", req.SupportedServiceID)
				assert.Equal(t, pss.ChangeCaseImpactLow, req.Impact)
			}).Return("CHG123456", nil),
			service.EXPECT().GetChangeCase("CHG123456").Return(pss.ChangeCase{}, nil),
		)

		pssChangeCreate(service, cmd, []string{})
	})

	t.Run("InvalidRisk_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)
		cmd := pssChangeCreateCmd(nil)
		setFlags(cmd)
		cmd.Flags().Set("risk", "invalid")

		err := pssChangeCreate(service, cmd, []string{})
		assert.Contains(t, err.Error(), "Invalid pss.ChangeCaseRisk")
	})

	t.Run("InvalidImpact_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)
		cmd := pssChangeCreateCmd(nil)
		setFlags(cmd)
		cmd.Flags().Set("impact", "invalid")

		err := pssChangeCreate(service, cmd, []string{})
		assert.Contains(t, err.Error(), "Invalid pss.ChangeCaseImpact")
	})

	t.Run("CreateChangeError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)
		cmd := pssChangeCreateCmd(nil)
		setFlags(cmd)

		service.EXPECT().CreateChangeCase(gomock.Any()).Return("", errors.New("test error"))

		err := pssChangeCreate(service, cmd, []string{})
		assert.Equal(t, "Error creating change: test error", err.Error())
	})

	t.Run("GetChangeError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)
		cmd := pssChangeCreateCmd(nil)
		setFlags(cmd)

		gomock.InOrder(
			service.EXPECT().CreateChangeCase(gomock.Any()).Return("CHG123456", nil),
			service.EXPECT().GetChangeCase("CHG123456").Return(pss.ChangeCase{}, errors.New("test error")),
		)

		err := pssChangeCreate(service, cmd, []string{})
		assert.Equal(t, "Error retrieving new change: test error", err.Error())
	})
}

func Test_pssChangeApprove(t *testing.T) {
	t.Run("DefaultApprove", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)
		cmd := pssChangeApproveCmd(nil)
		cmd.Flags().Set("reason", "test reason")
		cmd.Flags().Set("contact", "123")

		gomock.InOrder(
			service.EXPECT().ApproveChangeCase("CHG123456", gomock.Any()).Do(func(changeID string, req pss.ApproveChangeCaseRequest) {
				assert.Equal(t, "test reason", req.Reason)
				assert.Equal(t, 123, req.ContactID)
			}).Return("CHG123456", nil),
			service.EXPECT().GetChangeCase("CHG123456").Return(pss.ChangeCase{}, nil),
		)

		pssChangeApprove(service, cmd, []string{"CHG123456"})
	})

	t.Run("ApproveChangeCaseError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)
		cmd := pssChangeApproveCmd(nil)

		service.EXPECT().ApproveChangeCase("CHG123456", gomock.Any()).Return("", errors.New("test error"))

		test_output.AssertErrorOutput(t, "Failed to approve change [CHG123456]: test error\n", func() {
			pssChangeApprove(service, cmd, []string{"CHG123456"})
		})
	})

	t.Run("GetChangeCaseError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)
		cmd := pssChangeApproveCmd(nil)

		gomock.InOrder(
			service.EXPECT().ApproveChangeCase("CHG123456", gomock.Any()).Return("CHG123456", nil),
			service.EXPECT().GetChangeCase("CHG123456").Return(pss.ChangeCase{}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error retrieving approved change [CHG123456]: test error\n", func() {
			pssChangeApprove(service, cmd, []string{"CHG123456"})
		})
	})
}
