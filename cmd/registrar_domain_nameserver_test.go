package cmd

import (
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/cli/test"
	"github.com/ukfast/cli/test/mocks"
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
		code := 0
		oldOutputExit := output.SetOutputExit(func(c int) {
			code = c
		})
		defer func() { output.SetOutputExit(oldOutputExit) }()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockRegistrarService(mockCtrl)

		service.EXPECT().GetDomainNameservers(gomock.Any()).Return([]registrar.Nameserver{}, errors.New("test error")).Times(1)

		output := test.CatchStdErr(t, func() {
			registrarDomainNameserverList(service, &cobra.Command{}, []string{"testdomain1.co.uk"})
		})

		assert.Equal(t, 1, code)
		assert.Equal(t, "Error retrieving domain nameservers: test error\n", output)
	})
}
