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

func Test_registrarDomainNameserverListCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := registrarDomainNameserverListCmd().Args(nil, []string{"testdomain1.co.uk"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := registrarDomainNameserverListCmd().Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing domain", err.Error())
	})
}

func Test_registrarDomainNameserverList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockRegistrarService(mockCtrl)

		service.EXPECT().GetDomainNameservers(gomock.Any()).Return([]registrar.Nameserver{}, nil).Times(1)

		registrarDomainNameserverList(service, &cobra.Command{}, []string{"testdomain1.co.uk"})
	})

	t.Run("GetDomainNameserversError_OutputsFatal", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockRegistrarService(mockCtrl)

		service.EXPECT().GetDomainNameservers(gomock.Any()).Return([]registrar.Nameserver{}, errors.New("test error")).Times(1)

		test_output.AssertFatalOutput(t, "Error retrieving domain nameservers: test error\n", func() {
			registrarDomainNameserverList(service, &cobra.Command{}, []string{"testdomain1.co.uk"})
		})
	})
}
