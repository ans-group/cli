package loadbalancer

import (
	"errors"
	"testing"

	"github.com/ans-group/cli/test/mocks"
	"github.com/ans-group/cli/test/test_output"
	"github.com/ans-group/sdk-go/pkg/service/loadbalancer"
	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_loadbalancerACLShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := loadbalancerACLShowCmd(nil).Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := loadbalancerACLShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "missing ACL", err.Error())
	})
}

func Test_loadbalancerACLShow(t *testing.T) {
	t.Run("SingleACL", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().GetACL(123).Return(loadbalancer.ACL{}, nil).Times(1)

		loadbalancerACLShow(service, &cobra.Command{}, []string{"123"})
	})

	t.Run("MultipleACLs", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetACL(123).Return(loadbalancer.ACL{}, nil),
			service.EXPECT().GetACL(456).Return(loadbalancer.ACL{}, nil),
		)

		loadbalancerACLShow(service, &cobra.Command{}, []string{"123", "456"})
	})

	t.Run("GetACLID_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		test_output.AssertErrorOutput(t, "Invalid ACL ID [abc]\n", func() {
			loadbalancerACLShow(service, &cobra.Command{}, []string{"abc"})
		})
	})

	t.Run("GetACLError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().GetACL(123).Return(loadbalancer.ACL{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving ACL [123]: test error\n", func() {
			loadbalancerACLShow(service, &cobra.Command{}, []string{"123"})
		})
	})
}

func Test_loadbalancerACLCreate(t *testing.T) {
	t.Run("DefaultCreate", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)
		cmd := loadbalancerACLCreateCmd(nil)
		cmd.Flags().Set("name", "testacl")
		cmd.Flags().Set("mode", "http")

		req := loadbalancer.CreateACLRequest{
			Name: "testacl",
		}

		gomock.InOrder(
			service.EXPECT().CreateACL(req).Return(123, nil),
			service.EXPECT().GetACL(123).Return(loadbalancer.ACL{}, nil),
		)

		loadbalancerACLCreate(service, cmd, []string{})
	})

	t.Run("CreateACLError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)
		cmd := loadbalancerACLCreateCmd(nil)
		cmd.Flags().Set("name", "testacl")
		cmd.Flags().Set("mode", "http")

		service.EXPECT().CreateACL(gomock.Any()).Return(0, errors.New("test error"))

		err := loadbalancerACLCreate(service, cmd, []string{})

		assert.Equal(t, "error creating ACL: test error", err.Error())
	})

	t.Run("GetACLError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)
		cmd := loadbalancerACLCreateCmd(nil)
		cmd.Flags().Set("name", "testacl")
		cmd.Flags().Set("mode", "http")

		gomock.InOrder(
			service.EXPECT().CreateACL(gomock.Any()).Return(123, nil),
			service.EXPECT().GetACL(123).Return(loadbalancer.ACL{}, errors.New("test error")),
		)

		err := loadbalancerACLCreate(service, cmd, []string{})

		assert.Equal(t, "error retrieving new ACL: test error", err.Error())
	})
}

func Test_loadbalancerACLUpdateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := loadbalancerACLUpdateCmd(nil).Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := loadbalancerACLUpdateCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "missing ACL", err.Error())
	})
}

func Test_loadbalancerACLUpdate(t *testing.T) {
	t.Run("SingleGroup", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		cmd := loadbalancerACLUpdateCmd(nil)
		cmd.Flags().Set("name", "testacl")

		req := loadbalancer.PatchACLRequest{
			Name: "testacl",
		}

		gomock.InOrder(
			service.EXPECT().PatchACL(123, req).Return(nil),
			service.EXPECT().GetACL(123).Return(loadbalancer.ACL{}, nil),
		)

		loadbalancerACLUpdate(service, cmd, []string{"123"})
	})

	t.Run("MultipleACLs", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().PatchACL(123, gomock.Any()).Return(nil),
			service.EXPECT().GetACL(123).Return(loadbalancer.ACL{}, nil),
			service.EXPECT().PatchACL(456, gomock.Any()).Return(nil),
			service.EXPECT().GetACL(456).Return(loadbalancer.ACL{}, nil),
		)

		loadbalancerACLUpdate(service, &cobra.Command{}, []string{"123", "456"})
	})

	t.Run("PatchACLError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().PatchACL(123, gomock.Any()).Return(errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error updating ACL [123]: test error\n", func() {
			loadbalancerACLUpdate(service, &cobra.Command{}, []string{"123"})
		})
	})

	t.Run("GetACLError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().PatchACL(123, gomock.Any()).Return(nil),
			service.EXPECT().GetACL(123).Return(loadbalancer.ACL{}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error retrieving updated ACL [123]: test error\n", func() {
			loadbalancerACLUpdate(service, &cobra.Command{}, []string{"123"})
		})
	})
}

func Test_loadbalancerACLDeleteCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := loadbalancerACLDeleteCmd(nil).Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := loadbalancerACLDeleteCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "missing ACL", err.Error())
	})
}

func Test_loadbalancerACLDelete(t *testing.T) {
	t.Run("SingleGroup", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().DeleteACL(123).Return(nil).Times(1)

		loadbalancerACLDelete(service, &cobra.Command{}, []string{"123"})
	})

	t.Run("MultipleACLs", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().DeleteACL(123).Return(nil),
			service.EXPECT().DeleteACL(456).Return(nil),
		)

		loadbalancerACLDelete(service, &cobra.Command{}, []string{"123", "456"})
	})

	t.Run("DeleteACLError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().DeleteACL(123).Return(errors.New("test error")).Times(1)

		test_output.AssertErrorOutput(t, "Error removing ACL [123]: test error\n", func() {
			loadbalancerACLDelete(service, &cobra.Command{}, []string{"123"})
		})
	})
}

func Test_parseACLArguments(t *testing.T) {
	t.Run("SingleArgumentExpectedReturnValue", func(t *testing.T) {
		parsedArgs := []string{
			"arg1=testvalue1",
		}
		parsed, err := parseACLArguments(parsedArgs)

		assert.Nil(t, err)
		assert.Equal(t, "arg1", parsed["arg1"].Name)
		assert.Equal(t, "testvalue1", parsed["arg1"].Value)
	})

	t.Run("MultipleArgumentsExpectedReturnValue", func(t *testing.T) {
		parsedArgs := []string{
			"arg1=testvalue1",
			"arg2=testvalue2",
		}
		parsed, err := parseACLArguments(parsedArgs)

		assert.Nil(t, err)
		assert.Equal(t, "arg1", parsed["arg1"].Name)
		assert.Equal(t, "testvalue1", parsed["arg1"].Value)
		assert.Equal(t, "arg2", parsed["arg2"].Name)
		assert.Equal(t, "testvalue2", parsed["arg2"].Value)
	})

	t.Run("ArrayValueArgumentExpectedReturnValue", func(t *testing.T) {
		parsedArgs := []string{
			"arg1[]=testvalue1",
		}
		parsed, err := parseACLArguments(parsedArgs)

		assert.Nil(t, err)
		assert.Equal(t, "arg1", parsed["arg1"].Name)
		assert.Equal(t, "testvalue1", parsed["arg1"].Value.([]string)[0])
	})

	t.Run("MultipleArrayValueArgumentsExpectedReturnValue", func(t *testing.T) {
		parsedArgs := []string{
			"arg1[]=testvalue1",
			"arg1[]=testvalue2",
		}
		parsed, err := parseACLArguments(parsedArgs)

		assert.Nil(t, err)
		assert.Equal(t, "arg1", parsed["arg1"].Name)
		assert.Equal(t, "testvalue1", parsed["arg1"].Value.([]string)[0])
		assert.Equal(t, "testvalue2", parsed["arg1"].Value.([]string)[1])
	})
}
