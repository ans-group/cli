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
	"github.com/ukfast/sdk-go/pkg/service/ddosx"
)

func Test_ddosxDomainACLGeoIPRulesModeShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ddosxDomainACLGeoIPRulesModeShowCmd().Args(nil, []string{"testdomain1.co.uk"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ddosxDomainACLGeoIPRulesModeShowCmd().Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing domain", err.Error())
	})
}

func Test_ddosxDomainACLGeoIPRulesModeShow(t *testing.T) {
	t.Run("SingleDomain", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().GetDomainACLGeoIPRulesMode("testdomain1.co.uk").Return(ddosx.ACLGeoIPRulesMode(""), nil).Times(1)

		ddosxDomainACLGeoIPRulesModeShow(service, &cobra.Command{}, []string{"testdomain1.co.uk"})
	})

	t.Run("MultipleDomains", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetDomainACLGeoIPRulesMode("testdomain1.co.uk").Return(ddosx.ACLGeoIPRulesMode(""), nil),
			service.EXPECT().GetDomainACLGeoIPRulesMode("testdomain2.co.uk").Return(ddosx.ACLGeoIPRulesMode(""), nil),
		)

		ddosxDomainACLGeoIPRulesModeShow(service, &cobra.Command{}, []string{"testdomain1.co.uk", "testdomain2.co.uk"})
	})

	t.Run("GetDomainACLGeoIPRulesModeError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().GetDomainACLGeoIPRulesMode("testdomain1.co.uk").Return(ddosx.ACLGeoIPRulesMode(""), errors.New("test error"))

		output := test.CatchStdErr(t, func() {
			ddosxDomainACLGeoIPRulesModeShow(service, &cobra.Command{}, []string{"testdomain1.co.uk"})
		})

		assert.Equal(t, "Error retrieving domain [testdomain1.co.uk] ACL GeoIP rules mode: test error\n", output)
	})
}

func Test_ddosxDomainACLGeoIPRulesModeUpdateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		cmd := ddosxDomainACLGeoIPRulesModeUpdateCmd()
		err := cmd.Args(nil, []string{"testdomain1.co.uk"})

		assert.Nil(t, err)
	})

	t.Run("MissingDomain_Error", func(t *testing.T) {
		cmd := ddosxDomainACLGeoIPRulesModeUpdateCmd()
		err := cmd.Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing domain", err.Error())
	})
}

func Test_ddosxDomainACLGeoIPRulesModeUpdate(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)
		cmd := ddosxDomainACLGeoIPRulesModeUpdateCmd()
		cmd.Flags().Set("mode", "whitelist")

		expectedRequest := ddosx.PatchACLGeoIPRulesModeRequest{
			Mode: ddosx.ACLGeoIPRulesModeWhitelist,
		}

		gomock.InOrder(
			service.EXPECT().PatchDomainACLGeoIPRulesMode("testdomain1.co.uk", gomock.Eq(expectedRequest)).Return(nil),
			service.EXPECT().GetDomainACLGeoIPRulesMode("testdomain1.co.uk").Return(ddosx.ACLGeoIPRulesMode(""), nil),
		)

		ddosxDomainACLGeoIPRulesModeUpdate(service, cmd, []string{"testdomain1.co.uk"})
	})

	t.Run("InvalidMode_OutputsFatal", func(t *testing.T) {
		code := 0
		oldOutputExit := output.SetOutputExit(func(c int) {
			code = c
		})
		defer func() { output.SetOutputExit(oldOutputExit) }()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)
		cmd := ddosxDomainACLGeoIPRulesModeUpdateCmd()
		cmd.Flags().Set("mode", "invalidmode")

		output := test.CatchStdErr(t, func() {
			ddosxDomainACLGeoIPRulesModeUpdate(service, cmd, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})
		})

		assert.Equal(t, 1, code)
		assert.Equal(t, "Invalid ACL GeoIP rules filtering mode\n", output)
	})

	t.Run("PatchDomainACLGeoIPRulesMode_OutputsFatal", func(t *testing.T) {
		code := 0
		oldOutputExit := output.SetOutputExit(func(c int) {
			code = c
		})
		defer func() { output.SetOutputExit(oldOutputExit) }()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().PatchDomainACLGeoIPRulesMode("testdomain1.co.uk", gomock.Any()).Return(errors.New("test error"))

		output := test.CatchStdErr(t, func() {
			ddosxDomainACLGeoIPRulesModeUpdate(service, &cobra.Command{}, []string{"testdomain1.co.uk"})
		})

		assert.Equal(t, 1, code)
		assert.Equal(t, "Error updating domain ACL GeoIP rule filtering mode: test error\n", output)
	})

	t.Run("GetDomainACLGeoIPRulesMode_OutputsFatal", func(t *testing.T) {
		code := 0
		oldOutputExit := output.SetOutputExit(func(c int) {
			code = c
		})
		defer func() { output.SetOutputExit(oldOutputExit) }()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().PatchDomainACLGeoIPRulesMode("testdomain1.co.uk", gomock.Any()).Return(nil),
			service.EXPECT().GetDomainACLGeoIPRulesMode("testdomain1.co.uk").Return(ddosx.ACLGeoIPRulesMode(""), errors.New("test error")),
		)

		output := test.CatchStdErr(t, func() {
			ddosxDomainACLGeoIPRulesModeUpdate(service, &cobra.Command{}, []string{"testdomain1.co.uk"})
		})

		assert.Equal(t, 1, code)
		assert.Equal(t, "Error retrieving updated domain ACL GeoIP rule filtering mode: test error\n", output)
	})
}
