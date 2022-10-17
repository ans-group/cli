package loadbalancer

import (
	"errors"
	"testing"

	"github.com/ans-group/cli/test/mocks"
	"github.com/ans-group/sdk-go/pkg/service/loadbalancer"
	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_loadbalancerACLActionListCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := loadbalancerACLActionListCmd(nil).Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := loadbalancerACLActionListCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing ACL", err.Error())
	})
}

func Test_loadbalancerACLActionList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().GetACL(gomock.Any()).Return(loadbalancer.ACL{}, nil)

		err := loadbalancerACLActionList(service, &cobra.Command{}, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidACLID_Error", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		err := loadbalancerACLActionList(service, &cobra.Command{}, []string{"abc"})

		assert.NotNil(t, err)
		assert.Equal(t, "Invalid ACL ID [abc]", err.Error())
	})

	t.Run("GetACLError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().GetACL(gomock.Any()).Return(loadbalancer.ACL{}, errors.New("test error"))

		err := loadbalancerACLActionList(service, &cobra.Command{}, []string{"123"})

		assert.NotNil(t, err)
		assert.Equal(t, "Error retrieving ACL [123]: test error", err.Error())
	})
}

func Test_loadbalancerACLActionShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := loadbalancerACLActionShowCmd(nil).Args(nil, []string{"123", "0"})

		assert.Nil(t, err)
	})

	t.Run("MissingACL_Error", func(t *testing.T) {
		err := loadbalancerACLActionShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing ACL", err.Error())
	})

	t.Run("MissingIndex_Error", func(t *testing.T) {
		err := loadbalancerACLActionShowCmd(nil).Args(nil, []string{"123"})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing ACL action index", err.Error())
	})
}

func Test_loadbalancerACLActionShow(t *testing.T) {
	t.Run("SingleACL", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().GetACL(123).Return(loadbalancer.ACL{}, nil).Times(1)

		err := loadbalancerACLActionShow(service, &cobra.Command{}, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidACLID_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		err := loadbalancerACLActionShow(service, &cobra.Command{}, []string{"abc", "0"})

		assert.NotNil(t, err)
		assert.Equal(t, "Invalid ACL ID [abc]", err.Error())
	})

	t.Run("GetACLError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().GetACL(123).Return(loadbalancer.ACL{}, errors.New("test error"))

		err := loadbalancerACLActionShow(service, &cobra.Command{}, []string{"123", "0"})

		assert.NotNil(t, err)
		assert.Equal(t, "Error retrieving ACL [123]: test error", err.Error())
	})
}

func Test_loadbalancerACLActionCreateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := loadbalancerACLActionCreateCmd(nil).Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("MissingACL_Error", func(t *testing.T) {
		err := loadbalancerACLActionCreateCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing ACL", err.Error())
	})
}

func Test_loadbalancerACLActionCreate(t *testing.T) {
	t.Run("DefaultCreate_ExpectedPatch", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)
		cmd := loadbalancerACLActionCreateCmd(nil)
		cmd.Flags().Set("name", "header_matches")
		cmd.Flags().Set("argument", "header=host")
		cmd.Flags().Set("argument", "value=test.com")

		req := loadbalancer.PatchACLRequest{
			Actions: []loadbalancer.ACLAction{
				{
					Name: "header_matches",
					Arguments: map[string]loadbalancer.ACLArgument{
						"header": {
							Name:  "header",
							Value: "host",
						},
						"value": {
							Name:  "value",
							Value: "test.com",
						},
					},
				},
			},
		}

		gomock.InOrder(
			service.EXPECT().GetACL(123).Return(loadbalancer.ACL{}, nil),
			service.EXPECT().PatchACL(123, req).Return(nil),
			service.EXPECT().GetACL(123).Return(loadbalancer.ACL{}, nil),
		)

		err := loadbalancerACLActionCreate(service, cmd, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("GetACLError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)
		cmd := loadbalancerACLActionCreateCmd(nil)

		service.EXPECT().GetACL(gomock.Any()).Return(loadbalancer.ACL{}, errors.New("test error"))

		err := loadbalancerACLActionCreate(service, cmd, []string{"123"})

		assert.Equal(t, "Error retrieving ACL: test error", err.Error())
	})

	t.Run("PatchACLError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)
		cmd := loadbalancerACLActionCreateCmd(nil)

		gomock.InOrder(
			service.EXPECT().GetACL(123).Return(loadbalancer.ACL{}, nil),
			service.EXPECT().PatchACL(123, gomock.Any()).Return(errors.New("test error")),
		)

		err := loadbalancerACLActionCreate(service, cmd, []string{"123"})

		assert.Equal(t, "Error updating ACL: test error", err.Error())
	})

	t.Run("UpdatedGetACLError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)
		cmd := loadbalancerACLActionCreateCmd(nil)

		gomock.InOrder(
			service.EXPECT().GetACL(123).Return(loadbalancer.ACL{}, nil),
			service.EXPECT().PatchACL(123, gomock.Any()).Return(nil),
			service.EXPECT().GetACL(123).Return(loadbalancer.ACL{}, errors.New("test error")),
		)

		err := loadbalancerACLActionCreate(service, cmd, []string{"123"})

		assert.Equal(t, "Error retrieving updated ACL [123]: test error", err.Error())
	})
}

func Test_loadbalancerACLActionUpdateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := loadbalancerACLActionUpdateCmd(nil).Args(nil, []string{"123", "0"})

		assert.Nil(t, err)
	})

	t.Run("MissingACL_Error", func(t *testing.T) {
		err := loadbalancerACLActionUpdateCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing ACL", err.Error())
	})

	t.Run("MissingACLActionIndex_Error", func(t *testing.T) {
		err := loadbalancerACLActionUpdateCmd(nil).Args(nil, []string{"123"})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing ACL action index", err.Error())
	})
}

func Test_loadbalancerACLActionUpdate(t *testing.T) {
	t.Run("DefaultUpdate_ExpectedPatch", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)
		cmd := loadbalancerACLActionUpdateCmd(nil)
		cmd.Flags().Set("name", "header_matches")

		req := loadbalancer.PatchACLRequest{
			Actions: []loadbalancer.ACLAction{
				{
					Name: "header_matches",
				},
			},
		}

		gomock.InOrder(
			service.EXPECT().GetACL(123).Return(loadbalancer.ACL{
				Actions: []loadbalancer.ACLAction{
					{
						Name: "redirect",
					},
				},
			}, nil),
			service.EXPECT().PatchACL(123, req).Return(nil),
			service.EXPECT().GetACL(123).Return(loadbalancer.ACL{}, nil),
		)

		err := loadbalancerACLActionUpdate(service, cmd, []string{"123", "0"})
		assert.Nil(t, err)
	})

	t.Run("InvalidACLID_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)
		cmd := loadbalancerACLActionUpdateCmd(nil)

		err := loadbalancerACLActionUpdate(service, cmd, []string{"abc"})

		assert.NotNil(t, err)
		assert.Equal(t, "Invalid ACL ID [abc]", err.Error())
	})

	t.Run("GetACLError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)
		cmd := loadbalancerACLActionUpdateCmd(nil)

		service.EXPECT().GetACL(123).Return(loadbalancer.ACL{}, errors.New("test error"))

		err := loadbalancerACLActionUpdate(service, cmd, []string{"123"})

		assert.NotNil(t, err)
		assert.Equal(t, "Error retrieving ACL: test error", err.Error())
	})

	t.Run("InvalidActionIndex_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)
		cmd := loadbalancerACLActionUpdateCmd(nil)

		service.EXPECT().GetACL(123).Return(loadbalancer.ACL{}, nil)

		err := loadbalancerACLActionUpdate(service, cmd, []string{"123", "abc"})

		assert.NotNil(t, err)
		assert.Equal(t, "Invalid ACL action index [abc]", err.Error())
	})

	t.Run("ActionIndexOutOfBounds_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)
		cmd := loadbalancerACLActionUpdateCmd(nil)

		service.EXPECT().GetACL(123).Return(loadbalancer.ACL{
			Actions: []loadbalancer.ACLAction{
				{
					Name: "redirect",
				},
			},
		}, nil)

		err := loadbalancerACLActionUpdate(service, cmd, []string{"123", "1"})

		assert.NotNil(t, err)
		assert.Equal(t, "ACL action index [1] out of bounds", err.Error())
	})

	t.Run("PatchACLError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)
		cmd := loadbalancerACLActionUpdateCmd(nil)
		cmd.Flags().Set("name", "header_matches")

		gomock.InOrder(
			service.EXPECT().GetACL(123).Return(loadbalancer.ACL{
				Actions: []loadbalancer.ACLAction{
					{
						Name: "redirect",
					},
				},
			}, nil),
			service.EXPECT().PatchACL(123, gomock.Any()).Return(errors.New("test error")),
		)

		err := loadbalancerACLActionUpdate(service, cmd, []string{"123", "0"})

		assert.NotNil(t, err)
		assert.Equal(t, "Error updating ACL: test error", err.Error())
	})

	t.Run("UpdatedGetACLError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)
		cmd := loadbalancerACLActionUpdateCmd(nil)
		cmd.Flags().Set("name", "header_matches")

		gomock.InOrder(
			service.EXPECT().GetACL(123).Return(loadbalancer.ACL{
				Actions: []loadbalancer.ACLAction{
					{
						Name: "redirect",
					},
				},
			}, nil),
			service.EXPECT().PatchACL(123, gomock.Any()).Return(nil),
			service.EXPECT().GetACL(123).Return(loadbalancer.ACL{}, errors.New("test error")),
		)

		err := loadbalancerACLActionUpdate(service, cmd, []string{"123", "0"})

		assert.NotNil(t, err)
		assert.Equal(t, "Error retrieving updated ACL [123]: test error", err.Error())
	})
}

func Test_loadbalancerACLActionDeleteCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := loadbalancerACLActionDeleteCmd(nil).Args(nil, []string{"123", "0"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := loadbalancerACLActionDeleteCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing ACL", err.Error())
	})

	t.Run("MissingACLActionIndex_Error", func(t *testing.T) {
		err := loadbalancerACLActionDeleteCmd(nil).Args(nil, []string{"123"})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing ACL action index", err.Error())
	})
}

func Test_loadbalancerACLActionDelete(t *testing.T) {
	t.Run("DefaultDelete_ExpectedPatch", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)
		cmd := loadbalancerACLActionDeleteCmd(nil)
		cmd.Flags().Set("name", "header_matches")

		req := loadbalancer.PatchACLRequest{
			Actions: make([]loadbalancer.ACLAction, 0),
		}

		gomock.InOrder(
			service.EXPECT().GetACL(123).Return(loadbalancer.ACL{
				Actions: []loadbalancer.ACLAction{
					{
						Name: "redirect",
					},
				},
			}, nil),
			service.EXPECT().PatchACL(123, req).Return(nil),
			service.EXPECT().GetACL(123).Return(loadbalancer.ACL{}, nil),
		)

		err := loadbalancerACLActionDelete(service, cmd, []string{"123", "0"})
		assert.Nil(t, err)
	})

	t.Run("InvalidACLID_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)
		cmd := loadbalancerACLActionDeleteCmd(nil)

		err := loadbalancerACLActionDelete(service, cmd, []string{"abc"})

		assert.NotNil(t, err)
		assert.Equal(t, "Invalid ACL ID [abc]", err.Error())
	})

	t.Run("GetACLError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)
		cmd := loadbalancerACLActionDeleteCmd(nil)

		service.EXPECT().GetACL(123).Return(loadbalancer.ACL{}, errors.New("test error"))

		err := loadbalancerACLActionDelete(service, cmd, []string{"123"})

		assert.NotNil(t, err)
		assert.Equal(t, "Error retrieving ACL: test error", err.Error())
	})

	t.Run("InvalidActionIndex_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)
		cmd := loadbalancerACLActionDeleteCmd(nil)

		service.EXPECT().GetACL(123).Return(loadbalancer.ACL{}, nil)

		err := loadbalancerACLActionDelete(service, cmd, []string{"123", "abc"})

		assert.NotNil(t, err)
		assert.Equal(t, "Invalid ACL action index [abc]", err.Error())
	})

	t.Run("ActionIndexOutOfBounds_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)
		cmd := loadbalancerACLActionDeleteCmd(nil)

		service.EXPECT().GetACL(123).Return(loadbalancer.ACL{
			Actions: []loadbalancer.ACLAction{
				{
					Name: "redirect",
				},
			},
		}, nil)

		err := loadbalancerACLActionDelete(service, cmd, []string{"123", "1"})

		assert.NotNil(t, err)
		assert.Equal(t, "ACL action index [1] out of bounds", err.Error())
	})

	t.Run("PatchACLError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)
		cmd := loadbalancerACLActionDeleteCmd(nil)

		gomock.InOrder(
			service.EXPECT().GetACL(123).Return(loadbalancer.ACL{
				Actions: []loadbalancer.ACLAction{
					{
						Name: "redirect",
					},
				},
			}, nil),
			service.EXPECT().PatchACL(123, gomock.Any()).Return(errors.New("test error")),
		)

		err := loadbalancerACLActionDelete(service, cmd, []string{"123", "0"})

		assert.NotNil(t, err)
		assert.Equal(t, "Error updating ACL: test error", err.Error())
	})

	t.Run("UpdatedGetACLError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)
		cmd := loadbalancerACLActionDeleteCmd(nil)

		gomock.InOrder(
			service.EXPECT().GetACL(123).Return(loadbalancer.ACL{
				Actions: []loadbalancer.ACLAction{
					{
						Name: "redirect",
					},
				},
			}, nil),
			service.EXPECT().PatchACL(123, gomock.Any()).Return(nil),
			service.EXPECT().GetACL(123).Return(loadbalancer.ACL{}, errors.New("test error")),
		)

		err := loadbalancerACLActionDelete(service, cmd, []string{"123", "0"})

		assert.NotNil(t, err)
		assert.Equal(t, "Error retrieving updated ACL [123]: test error", err.Error())
	})
}
