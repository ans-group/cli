package ecloud

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/ans-group/cli/test/mocks"
	"github.com/ans-group/cli/test/test_output"
	"github.com/ans-group/sdk-go/pkg/connection"
	"github.com/ans-group/sdk-go/pkg/service/ecloud"
	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

// describeTestCmd returns a describe command with the persistent output flag declared
// by the root command, with the given flags parsed
func describeTestCmd(t *testing.T, flags ...string) *cobra.Command {
	cmd := ecloudDescribeCmd(nil)
	cmd.Flags().String("output", "table", "")

	err := cmd.ParseFlags(flags)
	assert.Nil(t, err)

	return cmd
}

// describeTestOptions returns the describe options resolved from the given flags
func describeTestOptions(t *testing.T, flags ...string) describeOptions {
	return describeOptionsFromCommand(describeTestCmd(t, flags...))
}

// describeTestPaginated returns a page of items reporting the given total number of
// items available across all pages
func describeTestPaginated[T any](items []T, total int) *connection.Paginated[T] {
	body := &connection.APIResponseBodyData[[]T]{Data: items}
	body.Metadata.Pagination.Total = total

	return connection.NewPaginated(body, connection.APIRequestParameters{}, nil)
}

// describeTestField returns the value rendered against the named field, or an empty
// string where the field was not rendered
func describeTestField(out string, name string) string {
	for line := range strings.SplitSeq(out, "\n") {
		fields := strings.Fields(line)
		if len(fields) > 0 && fields[0] == name {
			return strings.Join(fields[1:], " ")
		}
	}

	return ""
}

func Test_ecloudDescribeCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudDescribeCmd(nil).Args(nil, []string{"i-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudDescribeCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "missing id", err.Error())
	})
}

func Test_ECloudRootCmd_RegistersDescribe(t *testing.T) {
	var found bool
	for _, cmd := range ECloudRootCmd(nil, nil).Commands() {
		if cmd.Name() == "describe" {
			found = true
		}
	}

	assert.True(t, found, "expected describe to be registered on the eCloud root command")
}

func Test_ecloudDescribeResource(t *testing.T) {
	t.Run("InvalidID_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		_, err := ecloudDescribeResource(newDescribeSession(service), "notanid", describeTestOptions(t))

		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "invalid resource ID")
	})

	t.Run("UnknownPrefix_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		_, err := ecloudDescribeResource(newDescribeSession(service), "zzz-abcdef12", describeTestOptions(t))

		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "unknown resource type prefix [zzz]")
	})

	t.Run("GetInstanceError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		testError := errors.New("test error")
		service.EXPECT().GetInstance("i-abcdef12").Return(ecloud.Instance{}, testError).Times(1)

		_, err := ecloudDescribeResource(newDescribeSession(service), "i-abcdef12", describeTestOptions(t))

		assert.NotNil(t, err)
		assert.Equal(t, "error retrieving instance [i-abcdef12]: test error", err.Error())
		// The SDK error is wrapped, so remains available to errors.Is/errors.As
		assert.ErrorIs(t, err, testError)
	})

	t.Run("Depth0_FetchesResourceOnly", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetInstance("i-abcdef12").Return(ecloud.Instance{ID: "i-abcdef12", VPCID: "vpc-abcdef12"}, nil).Times(1)
		service.EXPECT().GetVPC(gomock.Any()).Times(0)
		expectNoInstanceChildren(service)

		result, err := ecloudDescribeResource(newDescribeSession(service), "i-abcdef12", describeTestOptions(t, "--depth=0"))

		assert.Nil(t, err)
		assert.Equal(t, "Instance", result.Type)
		assert.Len(t, result.Related, 0)
		assert.Len(t, result.Children, 0)
	})

	t.Run("ResolvesParents", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetVPC("vpc-abcdef12").Return(ecloud.VPC{ID: "vpc-abcdef12", Name: "test vpc"}, nil).Times(1)
		expectNoInstanceChildren(service)

		result, err := describeTestInstance(t, service, ecloud.Instance{ID: "i-abcdef12", VPCID: "vpc-abcdef12"}, "--no-children")

		assert.Nil(t, err)
		assert.Len(t, result.Related, 1)
		assert.Equal(t, "vpc_id", result.Related[0].Field)
		assert.Equal(t, "VPC", result.Related[0].Type)
		assert.Equal(t, "vpc-abcdef12", result.Related[0].ID)
		assert.Equal(t, "test vpc", result.Related[0].Name)
		assert.Empty(t, result.Related[0].Error)
	})

	t.Run("UnknownParentPrefix_NotResolved", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		expectNoInstanceChildren(service)

		result, err := describeTestInstance(t, service, ecloud.Instance{ID: "i-abcdef12", VPCID: "zzz-abcdef12"}, "--no-children")

		assert.Nil(t, err)
		assert.Len(t, result.Related, 1)
		assert.Empty(t, result.Related[0].Type)
		assert.Empty(t, result.Related[0].Name)
		assert.Empty(t, result.Related[0].Error)
	})

	t.Run("ParentLookupFailure_DoesNotFailCommand", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetVPC("vpc-abcdef12").Return(ecloud.VPC{}, errors.New("test error")).Times(1)
		expectNoInstanceChildren(service)

		result, err := describeTestInstance(t, service, ecloud.Instance{ID: "i-abcdef12", VPCID: "vpc-abcdef12"}, "--no-children")

		assert.Nil(t, err)
		assert.Len(t, result.Related, 1)
		assert.Equal(t, "error retrieving vpc [vpc-abcdef12]: test error", result.Related[0].Error)
		assert.Empty(t, result.Related[0].Name)
	})

	t.Run("EmptyParentID_NotResolved", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetVPC(gomock.Any()).Times(0)
		expectNoInstanceChildren(service)

		result, err := describeTestInstance(t, service, ecloud.Instance{ID: "i-abcdef12", VPCID: ""}, "--no-children")

		assert.Nil(t, err)
		assert.Len(t, result.Related, 0)
	})

	t.Run("FetchesChildren", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		expectInstanceChildren(service, nil)

		result, err := describeTestInstance(t, service, ecloud.Instance{ID: "i-abcdef12"})

		assert.Nil(t, err)
		assert.Len(t, result.Children, 5)
		assert.Equal(t, "NICs", result.Children[0].Title)
		assert.Equal(t, NICCollection{{ID: "nic-abcdef12"}}, result.Children[0].Items)
		assert.Equal(t, 1, result.Children[0].Total)
		assert.False(t, result.Children[0].Truncated)
	})

	t.Run("FetchesFirewallRulePorts", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetFirewallRule("fwr-abcdef12").Return(ecloud.FirewallRule{ID: "fwr-abcdef12"}, nil).Times(1)
		service.EXPECT().GetFirewallRuleFirewallRulePortsPaginated("fwr-abcdef12", gomock.Any()).
			Return(describeTestPaginated([]ecloud.FirewallRulePort{{ID: "fwrp-abcdef12"}}, 1), nil).Times(1)

		result, err := ecloudDescribeResource(newDescribeSession(service), "fwr-abcdef12", describeTestOptions(t))

		assert.Nil(t, err)
		assert.Len(t, result.Children, 1)
		assert.Equal(t, "Ports", result.Children[0].Title)
		assert.Equal(t, FirewallRulePortCollection{{ID: "fwrp-abcdef12"}}, result.Children[0].Items)
	})

	t.Run("FetchesNetworkRulePorts", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetNetworkRule("nr-abcdef12").Return(ecloud.NetworkRule{ID: "nr-abcdef12"}, nil).Times(1)
		service.EXPECT().GetNetworkRuleNetworkRulePortsPaginated("nr-abcdef12", gomock.Any()).
			Return(describeTestPaginated([]ecloud.NetworkRulePort{{ID: "nrp-abcdef12"}}, 1), nil).Times(1)

		result, err := ecloudDescribeResource(newDescribeSession(service), "nr-abcdef12", describeTestOptions(t))

		assert.Nil(t, err)
		assert.Len(t, result.Children, 1)
		assert.Equal(t, "Ports", result.Children[0].Title)
		assert.Equal(t, NetworkRulePortCollection{{ID: "nrp-abcdef12"}}, result.Children[0].Items)
	})

	t.Run("RedactsCredentialPasswords", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		expectInstanceChildren(service, []ecloud.Credential{{ID: "cred-abcdef12", Username: "root", Password: "somepassword"}})

		result, err := describeTestInstance(t, service, ecloud.Instance{ID: "i-abcdef12"})

		assert.Nil(t, err)
		assert.Equal(t, CredentialCollection{{ID: "cred-abcdef12", Username: "root", Password: "****"}}, result.Children[3].Items)
	})

	t.Run("WithPasswords_DisplaysCredentialPasswords", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		expectInstanceChildren(service, []ecloud.Credential{{ID: "cred-abcdef12", Username: "root", Password: "somepassword"}})

		result, err := describeTestInstance(t, service, ecloud.Instance{ID: "i-abcdef12"}, "--with-passwords")

		assert.Nil(t, err)
		assert.Equal(t, CredentialCollection{{ID: "cred-abcdef12", Username: "root", Password: "somepassword"}}, result.Children[3].Items)
	})

	t.Run("ChildFetchFailure_DoesNotFailCommand", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetInstanceNICsPaginated("i-abcdef12", gomock.Any()).Return(nil, errors.New("test error")).Times(1)
		service.EXPECT().GetInstanceVolumesPaginated("i-abcdef12", gomock.Any()).Return(describeTestPaginated([]ecloud.Volume{}, 0), nil).Times(1)
		service.EXPECT().GetInstanceFloatingIPsPaginated("i-abcdef12", gomock.Any()).Return(describeTestPaginated([]ecloud.FloatingIP{}, 0), nil).Times(1)
		service.EXPECT().GetInstanceCredentialsPaginated("i-abcdef12", gomock.Any()).Return(describeTestPaginated([]ecloud.Credential{}, 0), nil).Times(1)
		service.EXPECT().GetInstanceTasksPaginated("i-abcdef12", gomock.Any()).Return(describeTestPaginated([]ecloud.Task{}, 0), nil).Times(1)

		result, err := describeTestInstance(t, service, ecloud.Instance{ID: "i-abcdef12"})

		assert.Nil(t, err)
		assert.Equal(t, "error retrieving nics: test error", result.Children[0].Error)
	})

	t.Run("NoChildren_SkipsChildFetches", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		expectNoInstanceChildren(service)

		result, err := describeTestInstance(t, service, ecloud.Instance{ID: "i-abcdef12"}, "--no-children")

		assert.Nil(t, err)
		assert.Len(t, result.Children, 0)
	})

	t.Run("ChildLimit_RetrievesSinglePageAndReportsTruncation", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		assertPerPage := func(id string, params connection.APIRequestParameters) {
			assert.Equal(t, 1, params.Pagination.PerPage)
		}

		// One of three NICs is returned, so the collection is truncated. Each getter
		// is expected exactly once: the paginated variants retrieve a single page,
		// whereas their non-paginated counterparts would retrieve every page.
		service.EXPECT().GetInstanceNICsPaginated("i-abcdef12", gomock.Any()).Do(assertPerPage).
			Return(describeTestPaginated([]ecloud.NIC{{ID: "nic-abcdef12"}}, 3), nil).Times(1)
		service.EXPECT().GetInstanceVolumesPaginated("i-abcdef12", gomock.Any()).Do(assertPerPage).
			Return(describeTestPaginated([]ecloud.Volume{}, 0), nil).Times(1)
		service.EXPECT().GetInstanceFloatingIPsPaginated("i-abcdef12", gomock.Any()).Do(assertPerPage).
			Return(describeTestPaginated([]ecloud.FloatingIP{}, 0), nil).Times(1)
		service.EXPECT().GetInstanceCredentialsPaginated("i-abcdef12", gomock.Any()).Do(assertPerPage).
			Return(describeTestPaginated([]ecloud.Credential{}, 0), nil).Times(1)
		service.EXPECT().GetInstanceTasksPaginated("i-abcdef12", gomock.Any()).Do(assertPerPage).
			Return(describeTestPaginated([]ecloud.Task{}, 0), nil).Times(1)

		result, err := describeTestInstance(t, service, ecloud.Instance{ID: "i-abcdef12"}, "--child-limit=1")

		assert.Nil(t, err)
		assert.True(t, result.Children[0].Truncated)
		assert.Equal(t, 3, result.Children[0].Total)
		assert.False(t, result.Children[1].Truncated)
	})

	t.Run("FullPage_NotReportedAsTruncated", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		// A page holding every available item is not truncated, even where it is
		// exactly as large as the child limit
		service.EXPECT().GetInstanceNICsPaginated("i-abcdef12", gomock.Any()).
			Return(describeTestPaginated([]ecloud.NIC{{ID: "nic-abcdef12"}, {ID: "nic-abcdef23"}}, 2), nil).Times(1)
		service.EXPECT().GetInstanceVolumesPaginated("i-abcdef12", gomock.Any()).Return(describeTestPaginated([]ecloud.Volume{}, 0), nil).Times(1)
		service.EXPECT().GetInstanceFloatingIPsPaginated("i-abcdef12", gomock.Any()).Return(describeTestPaginated([]ecloud.FloatingIP{}, 0), nil).Times(1)
		service.EXPECT().GetInstanceCredentialsPaginated("i-abcdef12", gomock.Any()).Return(describeTestPaginated([]ecloud.Credential{}, 0), nil).Times(1)
		service.EXPECT().GetInstanceTasksPaginated("i-abcdef12", gomock.Any()).Return(describeTestPaginated([]ecloud.Task{}, 0), nil).Times(1)

		result, err := describeTestInstance(t, service, ecloud.Instance{ID: "i-abcdef12"}, "--child-limit=2")

		assert.Nil(t, err)
		assert.False(t, result.Children[0].Truncated)
	})
}

func Test_describeSession(t *testing.T) {
	t.Run("ResolvesEachParentOnce", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		// Both instances reference the same VPC, which must be retrieved only once
		service.EXPECT().GetInstance("i-abcdef12").Return(ecloud.Instance{ID: "i-abcdef12", VPCID: "vpc-abcdef12"}, nil).Times(1)
		service.EXPECT().GetInstance("i-abcdef23").Return(ecloud.Instance{ID: "i-abcdef23", VPCID: "vpc-abcdef12"}, nil).Times(1)
		service.EXPECT().GetVPC("vpc-abcdef12").Return(ecloud.VPC{ID: "vpc-abcdef12", Name: "test vpc"}, nil).Times(1)

		cmd := describeTestCmd(t, "--no-children")
		cmd.SetOut(&bytes.Buffer{})

		err := ecloudDescribe(service, cmd, []string{"i-abcdef12", "i-abcdef23"})

		assert.Nil(t, err)
	})

	t.Run("UnknownPrefix_NotRetrieved", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		resourceType, name, resolveError := newDescribeSession(service).resolve("zzz-abcdef12")

		assert.Empty(t, resourceType)
		assert.Empty(t, name)
		assert.Empty(t, resolveError)
	})
}

func Test_ecloudDescribe(t *testing.T) {
	t.Run("MultipleIDs_DescribesEach", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetInstance("i-abcdef12").Return(ecloud.Instance{ID: "i-abcdef12"}, nil).Times(1)
		service.EXPECT().GetVPC("vpc-abcdef12").Return(ecloud.VPC{ID: "vpc-abcdef12"}, nil).Times(1)

		cmd := describeTestCmd(t, "--depth=0")
		buf := &bytes.Buffer{}
		cmd.SetOut(buf)

		err := ecloudDescribe(service, cmd, []string{"i-abcdef12", "vpc-abcdef12"})

		assert.Nil(t, err)
		assert.Contains(t, buf.String(), "Instance  i-abcdef12")
		assert.Contains(t, buf.String(), "VPC  vpc-abcdef12")
	})

	t.Run("GetInstanceError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetInstance("i-abcdef12").Return(ecloud.Instance{}, errors.New("test error")).Times(1)
		service.EXPECT().GetVPC("vpc-abcdef12").Return(ecloud.VPC{ID: "vpc-abcdef12", Name: "test vpc"}, nil).Times(1)

		cmd := describeTestCmd(t, "--depth=0")
		buf := &bytes.Buffer{}
		cmd.SetOut(buf)

		// A failure describing one resource must not prevent the others being described
		test_output.AssertErrorOutput(t, "error retrieving instance [i-abcdef12]: test error\n", func() {
			err := ecloudDescribe(service, cmd, []string{"i-abcdef12", "vpc-abcdef12"})
			assert.Nil(t, err)
		})

		assert.Contains(t, buf.String(), "VPC  vpc-abcdef12")
	})

	t.Run("InvalidChildLimit_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetInstance(gomock.Any()).Times(0)

		err := ecloudDescribe(service, describeTestCmd(t, "--child-limit=0"), []string{"i-abcdef12"})

		assert.NotNil(t, err)
		assert.Equal(t, "ecloud: invalid child limit [0], must be greater than 0", err.Error())
	})

	t.Run("WithPasswords_RendersPassword", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetInstance("i-abcdef12").Return(ecloud.Instance{ID: "i-abcdef12"}, nil).Times(1)
		expectInstanceChildren(service, []ecloud.Credential{{ID: "cred-abcdef12", Username: "root", Password: "somepassword"}})

		cmd := describeTestCmd(t, "--with-passwords")
		buf := &bytes.Buffer{}
		cmd.SetOut(buf)

		err := ecloudDescribe(service, cmd, []string{"i-abcdef12"})

		assert.Nil(t, err)
		assert.Contains(t, buf.String(), "somepassword")
		assert.NotContains(t, buf.String(), "****")
	})

	t.Run("RedactsPasswordInRenderedOutput", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetInstance("i-abcdef12").Return(ecloud.Instance{ID: "i-abcdef12"}, nil).Times(1)
		expectInstanceChildren(service, []ecloud.Credential{{ID: "cred-abcdef12", Username: "root", Password: "somepassword"}})

		cmd := describeTestCmd(t)
		buf := &bytes.Buffer{}
		cmd.SetOut(buf)

		err := ecloudDescribe(service, cmd, []string{"i-abcdef12"})

		assert.Nil(t, err)
		assert.NotContains(t, buf.String(), "somepassword")
		assert.Contains(t, buf.String(), "****")
	})

	t.Run("RedactsPasswordInJSONOutput", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetInstance("i-abcdef12").Return(ecloud.Instance{ID: "i-abcdef12"}, nil).Times(1)
		expectInstanceChildren(service, []ecloud.Credential{{ID: "cred-abcdef12", Username: "root", Password: "somepassword"}})

		cmd := describeTestCmd(t, "--output=json")
		buf := &bytes.Buffer{}
		cmd.SetOut(buf)

		err := ecloudDescribe(service, cmd, []string{"i-abcdef12"})

		assert.Nil(t, err)
		assert.NotContains(t, buf.String(), "somepassword")
		assert.Contains(t, buf.String(), `"password":"****"`)
	})

	t.Run("JSONOutput_IsNested", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetInstance("i-abcdef12").Return(ecloud.Instance{ID: "i-abcdef12", Name: "test instance", VCPUCores: 2}, nil).Times(1)

		cmd := describeTestCmd(t, "--depth=0", "--output=json")
		buf := &bytes.Buffer{}
		cmd.SetOut(buf)

		err := ecloudDescribe(service, cmd, []string{"i-abcdef12"})
		assert.Nil(t, err)

		var results []map[string]any
		err = json.Unmarshal(buf.Bytes(), &results)
		assert.Nil(t, err)

		assert.Len(t, results, 1)
		assert.Equal(t, "Instance", results[0]["type"])

		resource, ok := results[0]["resource"].(map[string]any)
		assert.True(t, ok, "expected nested resource object")
		assert.Equal(t, "test instance", resource["name"])
		assert.Equal(t, float64(2), resource["vcpu_cores"])
	})

	t.Run("JSONPrettyOutput_IsIndented", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetInstance("i-abcdef12").Return(ecloud.Instance{ID: "i-abcdef12", Name: "test instance"}, nil).Times(1)

		cmd := describeTestCmd(t, "--depth=0", "--output=json-pretty")
		buf := &bytes.Buffer{}
		cmd.SetOut(buf)

		err := ecloudDescribe(service, cmd, []string{"i-abcdef12"})
		assert.Nil(t, err)

		assert.Contains(t, buf.String(), "\n    \"type\": \"Instance\",")
	})

	t.Run("YAMLOutput_IsNested", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetInstance("i-abcdef12").Return(ecloud.Instance{ID: "i-abcdef12", Name: "test instance", VCPUCores: 2}, nil).Times(1)

		cmd := describeTestCmd(t, "--depth=0", "--output=yaml")
		buf := &bytes.Buffer{}
		cmd.SetOut(buf)

		err := ecloudDescribe(service, cmd, []string{"i-abcdef12"})
		assert.Nil(t, err)

		var results []map[string]any
		err = yaml.Unmarshal(buf.Bytes(), &results)
		assert.Nil(t, err)

		assert.Len(t, results, 1)
		assert.Equal(t, "Instance", results[0]["type"])
		assert.IsType(t, map[string]any{}, results[0]["resource"])
	})

	t.Run("UnsupportedOutputFormat_OutputsWarning", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetInstance("i-abcdef12").Return(ecloud.Instance{ID: "i-abcdef12"}, nil).Times(1)

		cmd := describeTestCmd(t, "--depth=0", "--output=csv")
		buf := &bytes.Buffer{}
		cmd.SetOut(buf)

		test_output.AssertErrorOutput(t, "describe does not support output format [csv], using default\n", func() {
			err := ecloudDescribe(service, cmd, []string{"i-abcdef12"})
			assert.Nil(t, err)
		})

		assert.Contains(t, buf.String(), "Instance  i-abcdef12")
	})
}

func Test_describeRender(t *testing.T) {
	t.Run("RendersResourceRelatedAndChildren", func(t *testing.T) {
		buf := &bytes.Buffer{}

		describeRender(buf, []DescribeResult{
			{
				Type: "Instance",
				ID:   "i-abcdef12",
				Resource: ecloud.Instance{
					ID:        "i-abcdef12",
					Name:      "webserver-01",
					VPCID:     "vpc-abcdef12",
					VCPUCores: 2,
					Locked:    false,
					Sync:      ecloud.ResourceSync{Status: ecloud.SyncStatusComplete},
					Tags:      []ecloud.ResourceTag{{Scope: "env", Name: "production"}},
				},
				Related: []DescribeRelated{
					{Field: "vpc_id", Type: "VPC", ID: "vpc-abcdef12", Name: "test vpc"},
				},
				Children: []DescribeChildren{
					{Title: "NICs", Items: NICCollection{{ID: "nic-abcdef12", MACAddress: "00:50:56:a1:b2:c3"}}, Total: 1},
					{Title: "Volumes", Items: VolumeCollection{}},
					{Title: "Tasks", Error: "error retrieving tasks: test error"},
				},
			},
		})

		out := buf.String()

		assert.Contains(t, out, "Instance  i-abcdef12")
		assert.Equal(t, "webserver-01", describeTestField(out, "name"))
		assert.Equal(t, "2", describeTestField(out, "vcpu_cores"))
		// Booleans are retained when false, whereas other zero values are omitted
		assert.Equal(t, "false", describeTestField(out, "locked"))
		assert.Equal(t, "", describeTestField(out, "ram_capacity"))
		// Sync and task are flattened, and tags rendered as scope:name
		assert.Equal(t, "complete", describeTestField(out, "sync"))
		assert.Equal(t, "false", describeTestField(out, "task_in_progress"))
		assert.Equal(t, "env:production", describeTestField(out, "tags"))
		// Parent references are displayed in the related section, not inline
		assert.Equal(t, "", describeTestField(out, "vpc_id"))
		assert.Equal(t, "vpc-abcdef12 (test vpc)", describeTestField(out, "vpc"))
		assert.Contains(t, out, "NICs (1)")
		assert.Contains(t, out, "MAC_ADDRESS")
		assert.Contains(t, out, "00:50:56:a1:b2:c3")
		assert.Contains(t, out, "Volumes (0)")
		assert.Contains(t, out, "error retrieving tasks: test error")
	})

	t.Run("Depth0_RendersParentReferencesInline", func(t *testing.T) {
		buf := &bytes.Buffer{}

		describeRender(buf, []DescribeResult{
			{
				Type:     "Instance",
				ID:       "i-abcdef12",
				Resource: ecloud.Instance{ID: "i-abcdef12", VPCID: "vpc-abcdef12"},
			},
		})

		assert.Equal(t, "vpc-abcdef12", describeTestField(buf.String(), "vpc_id"))
	})

	t.Run("UnresolvedReference_RemainsInResourceBlock", func(t *testing.T) {
		buf := &bytes.Buffer{}

		clientID := 42

		// Only string references are resolved into the related section, so a
		// non-string identifier must remain visible in the resource block
		describeRender(buf, []DescribeResult{
			{
				Type:     "VPC",
				ID:       "vpc-abcdef12",
				Resource: ecloud.VPC{ID: "vpc-abcdef12", Name: "test vpc", RegionID: "reg-abcdef12", ClientID: &clientID},
				Related: []DescribeRelated{
					{Field: "region_id", Type: "Region", ID: "reg-abcdef12", Name: "test region"},
				},
			},
		})

		out := buf.String()

		assert.Equal(t, "42", describeTestField(out, "client_id"))
		assert.Equal(t, "", describeTestField(out, "region_id"))
	})

	t.Run("TruncatedChildren_NotesTruncation", func(t *testing.T) {
		buf := &bytes.Buffer{}

		describeRender(buf, []DescribeResult{
			{
				Type:     "VPC",
				ID:       "vpc-abcdef12",
				Resource: ecloud.VPC{ID: "vpc-abcdef12"},
				Children: []DescribeChildren{
					{Title: "Instances", Items: InstanceCollection{{ID: "i-abcdef12"}}, Total: 12, Truncated: true},
				},
			},
		})

		assert.Contains(t, buf.String(), "Instances (1 of 12, truncated)")
	})

	t.Run("CollectionWithoutDefaultColumns_DerivesColumns", func(t *testing.T) {
		buf := &bytes.Buffer{}

		// VPNs have no collection type declaring default columns, so the columns
		// are derived from the fields of the first item
		describeRender(buf, []DescribeResult{
			{
				Type:     "Router",
				ID:       "rtr-abcdef12",
				Resource: ecloud.Router{ID: "rtr-abcdef12"},
				Children: []DescribeChildren{
					{Title: "VPNs", Items: []ecloud.VPN{{ID: "vpn-abcdef12", RouterID: "rtr-abcdef12"}}, Total: 1},
				},
			},
		})

		out := buf.String()

		assert.Contains(t, out, "VPNs (1)")
		assert.Contains(t, out, "ID")
		assert.Contains(t, out, "ROUTER_ID")
		assert.Contains(t, out, "vpn-abcdef12")
	})
}

func Test_describeRedactResource(t *testing.T) {
	t.Run("RedactsPassword", func(t *testing.T) {
		credential := ecloud.Credential{ID: "cred-abcdef12", Username: "root", Password: "somepassword"}

		redacted, ok := describeRedactResource(credential).(ecloud.Credential)

		assert.True(t, ok)
		assert.Equal(t, "****", redacted.Password)
		assert.Equal(t, "root", redacted.Username)
		// The original is left untouched
		assert.Equal(t, "somepassword", credential.Password)
	})

	t.Run("EmptyPassword_NotRedacted", func(t *testing.T) {
		redacted, ok := describeRedactResource(ecloud.Credential{ID: "cred-abcdef12"}).(ecloud.Credential)

		assert.True(t, ok)
		assert.Empty(t, redacted.Password)
	})

	t.Run("ResourceWithoutSecrets_Unchanged", func(t *testing.T) {
		instance := ecloud.Instance{ID: "i-abcdef12", Name: "test instance"}

		assert.Equal(t, instance, describeRedactResource(instance))
	})
}

func Test_describeFieldValue(t *testing.T) {
	t.Run("ResolvesNestedFieldByUnderscoredName", func(t *testing.T) {
		fields := describeFlattenFields(ecloud.VPC{ID: "vpc-abcdef12", Sync: ecloud.ResourceSync{Status: ecloud.SyncStatusComplete}})

		assert.Equal(t, "complete", describeFieldValue(fields, "sync_status"))
		assert.Equal(t, "complete", describeFieldValue(fields, "sync.status"))
		assert.Equal(t, "vpc-abcdef12", describeFieldValue(fields, "id"))
		assert.Equal(t, "", describeFieldValue(fields, "notafield"))
	})
}

// describeTestInstance describes an instance returned by the mocked service, with the
// given flags applied
func describeTestInstance(t *testing.T, service *mocks.MockECloudService, instance ecloud.Instance, flags ...string) (DescribeResult, error) {
	service.EXPECT().GetInstance(instance.ID).Return(instance, nil).Times(1)

	return ecloudDescribeResource(newDescribeSession(service), instance.ID, describeTestOptions(t, flags...))
}

// expectInstanceChildren expects each instance child collection to be retrieved,
// returning the given credentials, one NIC, and empty collections otherwise
func expectInstanceChildren(service *mocks.MockECloudService, credentials []ecloud.Credential) {
	service.EXPECT().GetInstanceNICsPaginated("i-abcdef12", gomock.Any()).Return(describeTestPaginated([]ecloud.NIC{{ID: "nic-abcdef12"}}, 1), nil).Times(1)
	service.EXPECT().GetInstanceVolumesPaginated("i-abcdef12", gomock.Any()).Return(describeTestPaginated([]ecloud.Volume{}, 0), nil).Times(1)
	service.EXPECT().GetInstanceFloatingIPsPaginated("i-abcdef12", gomock.Any()).Return(describeTestPaginated([]ecloud.FloatingIP{}, 0), nil).Times(1)
	service.EXPECT().GetInstanceCredentialsPaginated("i-abcdef12", gomock.Any()).Return(describeTestPaginated(credentials, len(credentials)), nil).Times(1)
	service.EXPECT().GetInstanceTasksPaginated("i-abcdef12", gomock.Any()).Return(describeTestPaginated([]ecloud.Task{}, 0), nil).Times(1)
}

// expectNoInstanceChildren asserts that no instance child collections are retrieved
func expectNoInstanceChildren(service *mocks.MockECloudService) {
	service.EXPECT().GetInstanceNICsPaginated(gomock.Any(), gomock.Any()).Times(0)
	service.EXPECT().GetInstanceVolumesPaginated(gomock.Any(), gomock.Any()).Times(0)
	service.EXPECT().GetInstanceFloatingIPsPaginated(gomock.Any(), gomock.Any()).Times(0)
	service.EXPECT().GetInstanceCredentialsPaginated(gomock.Any(), gomock.Any()).Times(0)
	service.EXPECT().GetInstanceTasksPaginated(gomock.Any(), gomock.Any()).Times(0)
}
