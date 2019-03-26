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

func Test_registrarWhoisShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := registrarWhoisShowCmd().Args(nil, []string{"testdomain1.co.uk"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := registrarWhoisShowCmd().Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing domain", err.Error())
	})
}

func Test_registrarWhoisShow(t *testing.T) {
	t.Run("SingleWhois", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockRegistrarService(mockCtrl)

		service.EXPECT().GetWhois("testdomain1.co.uk").Return(registrar.Whois{}, nil).Times(1)

		registrarWhoisShow(service, &cobra.Command{}, []string{"testdomain1.co.uk"})
	})

	t.Run("MultipleWhoiss", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockRegistrarService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetWhois("testdomain1.co.uk").Return(registrar.Whois{}, nil),
			service.EXPECT().GetWhois("testdomain2.co.uk").Return(registrar.Whois{}, nil),
		)

		registrarWhoisShow(service, &cobra.Command{}, []string{"testdomain1.co.uk", "testdomain2.co.uk"})
	})

	t.Run("GetWhoisError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockRegistrarService(mockCtrl)

		service.EXPECT().GetWhois("testdomain1.co.uk").Return(registrar.Whois{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving whois for domain [testdomain1.co.uk]: test error\n", func() {
			registrarWhoisShow(service, &cobra.Command{}, []string{"testdomain1.co.uk"})
		})
	})
}

func Test_registrarWhoisShowRaw(t *testing.T) {
	t.Run("Default_Success", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockRegistrarService(mockCtrl)

		service.EXPECT().GetWhoisRaw("testdomain1.co.uk").Return("examplewhois", nil).Times(1)

		registrarWhoisShowRaw(service, &cobra.Command{}, []string{"testdomain1.co.uk"})
	})

	t.Run("GetWhoisRawError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockRegistrarService(mockCtrl)

		service.EXPECT().GetWhoisRaw("testdomain1.co.uk").Return("", errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving raw whois for domain [testdomain1.co.uk]: test error\n", func() {
			registrarWhoisShowRaw(service, &cobra.Command{}, []string{"testdomain1.co.uk"})
		})
	})
}
