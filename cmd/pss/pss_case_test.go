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

func Test_pssCaseList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)

		service.EXPECT().GetCasesPaginated(gomock.Any()).Return(connection.NewPaginated(&connection.APIResponseBodyData[[]pss.Case]{Data: []pss.Case{}}, *connection.NewAPIRequestParameters(), nil), nil).Times(1)

		pssCaseList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := pssCaseList(service, cmd, []string{})
		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetCasesError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)

		service.EXPECT().GetCasesPaginated(gomock.Any()).Return(connection.NewPaginated(&connection.APIResponseBodyData[[]pss.Case]{Data: []pss.Case{}}, *connection.NewAPIRequestParameters(), nil), errors.New("test error")).Times(1)

		err := pssCaseList(service, &cobra.Command{}, []string{})
		assert.Equal(t, "test error", err.Error())
	})
}

func Test_pssCaseShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := pssCaseShowCmd(nil).Args(nil, []string{"CHG123456"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := pssCaseShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing case", err.Error())
	})
}

func Test_pssCaseShow(t *testing.T) {
	t.Run("SingleCase", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)

		service.EXPECT().GetCase("CHG123456").Return(pss.Case{}, nil)

		pssCaseShow(service, &cobra.Command{}, []string{"CHG123456"})
	})

	t.Run("MultipleCases", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetCase("CHG123").Return(pss.Case{}, nil),
			service.EXPECT().GetCase("CHG456").Return(pss.Case{}, nil),
		)

		pssCaseShow(service, &cobra.Command{}, []string{"CHG123", "CHG456"})
	})

	t.Run("GetCaseError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)

		service.EXPECT().GetCase("CHG123456").Return(pss.Case{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving case [CHG123456]: test error\n", func() {
			pssCaseShow(service, &cobra.Command{}, []string{"CHG123456"})
		})
	})
}
