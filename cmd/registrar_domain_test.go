package cmd

import (
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/ukfast/cli/test/mocks"
	"github.com/ukfast/cli/test/test_output"
	"github.com/ukfast/sdk-go/pkg/service/registrar"
)

func Test_registrarDomainList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockRegistrarService(mockCtrl)

		service.EXPECT().GetDomains(gomock.Any()).Return([]registrar.Domain{}, nil).Times(1)

		registrarDomainList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_OutputsFatal", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockRegistrarService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		test_output.AssertFatalOutput(t, "Missing value for filtering\n", func() {
			registrarDomainList(service, cmd, []string{})
		})
	})

	t.Run("GetDomainsError_OutputsFatal", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockRegistrarService(mockCtrl)

		service.EXPECT().GetDomains(gomock.Any()).Return([]registrar.Domain{}, errors.New("test error")).Times(1)

		test_output.AssertFatalOutput(t, "Error retrieving domains: test error\n", func() {
			registrarDomainList(service, &cobra.Command{}, []string{})
		})
	})
}

func Test_registrarDomainShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := registrarDomainShowCmd().Args(nil, []string{"testdomain1.co.uk"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := registrarDomainShowCmd().Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing domain", err.Error())
	})
}

func Test_registrarDomainShow(t *testing.T) {
	t.Run("SingleDomain", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockRegistrarService(mockCtrl)

		service.EXPECT().GetDomain("testdomain1.co.uk").Return(registrar.Domain{}, nil).Times(1)

		registrarDomainShow(service, &cobra.Command{}, []string{"testdomain1.co.uk"})
	})

	t.Run("MultipleDomains", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockRegistrarService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetDomain("testdomain1.co.uk").Return(registrar.Domain{}, nil),
			service.EXPECT().GetDomain("testdomain2.co.uk").Return(registrar.Domain{}, nil),
		)

		registrarDomainShow(service, &cobra.Command{}, []string{"testdomain1.co.uk", "testdomain2.co.uk"})
	})

	t.Run("GetDomainError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockRegistrarService(mockCtrl)

		service.EXPECT().GetDomain("testdomain1.co.uk").Return(registrar.Domain{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving domain [testdomain1.co.uk]: test error\n", func() {
			registrarDomainShow(service, &cobra.Command{}, []string{"testdomain1.co.uk"})
		})
	})
}
