package cmd

import (
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/ukfast/cli/test/mocks"
	"github.com/ukfast/cli/test/test_output"
	"github.com/ukfast/sdk-go/pkg/connection"
	"github.com/ukfast/sdk-go/pkg/service/ltaas"
)

func Test_loadtestAgreementShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := loadtestAgreementShowCmd().Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := loadtestAgreementShowCmd().Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing agreement type", err.Error())
	})
}

func Test_loadtestAgreementShow(t *testing.T) {
	t.Run("SingleAgreement", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLTaaSService(mockCtrl)

		service.EXPECT().GetLatestAgreement(ltaas.AgreementTypeSingle).Return(ltaas.Agreement{}, nil).Times(1)

		loadtestAgreementShow(service, &cobra.Command{}, []string{"single"})
	})

	t.Run("MultipleAgreements", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLTaaSService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetLatestAgreement(ltaas.AgreementTypeSingle).Return(ltaas.Agreement{}, nil),
			service.EXPECT().GetLatestAgreement(ltaas.AgreementTypeRecurring).Return(ltaas.Agreement{}, nil),
		)

		loadtestAgreementShow(service, &cobra.Command{}, []string{"single", "recurring"})
	})

	t.Run("ParseAgreementType_ReturnsError", func(t *testing.T) {
		err := loadtestAgreementShow(nil, &cobra.Command{}, []string{"invalid"})

		assert.IsType(t, &connection.ErrInvalidEnumValue{}, err)
	})

	t.Run("GetAgreementError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLTaaSService(mockCtrl)

		service.EXPECT().GetLatestAgreement(ltaas.AgreementTypeSingle).Return(ltaas.Agreement{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving agreement [single]: test error\n", func() {
			loadtestAgreementShow(service, &cobra.Command{}, []string{"single"})
		})
	})
}
