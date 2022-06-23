package ddosx

import (
	"errors"
	"testing"

	"github.com/ans-group/cli/test/mocks"
	"github.com/ans-group/cli/test/test_output"
	"github.com/ans-group/sdk-go/pkg/service/ddosx"
	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_ddosxDomainDNSActivateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ddosxDomainDNSActivateCmd(nil).Args(nil, []string{"testdomain1.co.uk"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ddosxDomainDNSActivateCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing domain", err.Error())
	})
}

func Test_ddosxDomainDNSActivate(t *testing.T) {
	t.Run("SingleDomain", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().ActivateDomainDNSRouting("testdomain1.co.uk").Return(nil),
			service.EXPECT().GetDomain("testdomain1.co.uk").Return(ddosx.Domain{}, nil),
		)

		ddosxDomainDNSActivate(service, &cobra.Command{}, []string{"testdomain1.co.uk"})
	})

	t.Run("MultipleDomains", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().ActivateDomainDNSRouting("testdomain1.co.uk").Return(nil),
			service.EXPECT().GetDomain("testdomain1.co.uk").Return(ddosx.Domain{}, nil),
			service.EXPECT().ActivateDomainDNSRouting("testdomain2.co.uk").Return(nil),
			service.EXPECT().GetDomain("testdomain2.co.uk").Return(ddosx.Domain{}, nil),
		)

		ddosxDomainDNSActivate(service, &cobra.Command{}, []string{"testdomain1.co.uk", "testdomain2.co.uk"})
	})

	t.Run("DNSRoutingActivateDomainError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().ActivateDomainDNSRouting("testdomain1.co.uk").Return(errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error activating DNS routing for domain [testdomain1.co.uk]: test error\n", func() {
			ddosxDomainDNSActivate(service, &cobra.Command{}, []string{"testdomain1.co.uk"})
		})
	})

	t.Run("GetDomainError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().ActivateDomainDNSRouting("testdomain1.co.uk").Return(nil),
			service.EXPECT().GetDomain("testdomain1.co.uk").Return(ddosx.Domain{}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error retrieving domain [testdomain1.co.uk]: test error\n", func() {
			ddosxDomainDNSActivate(service, &cobra.Command{}, []string{"testdomain1.co.uk"})
		})
	})
}
func Test_ddosxDomainDNSDeactivateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ddosxDomainDNSDeactivateCmd(nil).Args(nil, []string{"testdomain1.co.uk"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ddosxDomainDNSDeactivateCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing domain", err.Error())
	})
}

func Test_ddosxDomainDNSDeactivate(t *testing.T) {
	t.Run("SingleDomain", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().DeactivateDomainDNSRouting("testdomain1.co.uk").Return(nil),
			service.EXPECT().GetDomain("testdomain1.co.uk").Return(ddosx.Domain{}, nil),
		)

		ddosxDomainDNSDeactivate(service, &cobra.Command{}, []string{"testdomain1.co.uk"})
	})

	t.Run("MultipleDomains", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().DeactivateDomainDNSRouting("testdomain1.co.uk").Return(nil),
			service.EXPECT().GetDomain("testdomain1.co.uk").Return(ddosx.Domain{}, nil),
			service.EXPECT().DeactivateDomainDNSRouting("testdomain2.co.uk").Return(nil),
			service.EXPECT().GetDomain("testdomain2.co.uk").Return(ddosx.Domain{}, nil),
		)

		ddosxDomainDNSDeactivate(service, &cobra.Command{}, []string{"testdomain1.co.uk", "testdomain2.co.uk"})
	})

	t.Run("DNSRoutingDeactivateDomainError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().DeactivateDomainDNSRouting("testdomain1.co.uk").Return(errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error deactivating DNS routing for domain [testdomain1.co.uk]: test error\n", func() {
			ddosxDomainDNSDeactivate(service, &cobra.Command{}, []string{"testdomain1.co.uk"})
		})
	})

	t.Run("GetDomainError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().DeactivateDomainDNSRouting("testdomain1.co.uk").Return(nil),
			service.EXPECT().GetDomain("testdomain1.co.uk").Return(ddosx.Domain{}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error retrieving domain [testdomain1.co.uk]: test error\n", func() {
			ddosxDomainDNSDeactivate(service, &cobra.Command{}, []string{"testdomain1.co.uk"})
		})
	})
}
