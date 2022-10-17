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

func Test_loadbalancerACLConditionListCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := loadbalancerACLConditionListCmd(nil).Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := loadbalancerACLConditionListCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing ACL", err.Error())
	})
}

func Test_loadbalancerACLConditionList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().GetACL(gomock.Any()).Return(loadbalancer.ACL{}, nil)

		err := loadbalancerACLConditionList(service, &cobra.Command{}, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidACLID_Error", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		err := loadbalancerACLConditionList(service, &cobra.Command{}, []string{"abc"})

		assert.NotNil(t, err)
		assert.Equal(t, "Invalid ACL ID [abc]", err.Error())
	})

	t.Run("GetACLError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().GetACL(gomock.Any()).Return(loadbalancer.ACL{}, errors.New("test error"))

		err := loadbalancerACLConditionList(service, &cobra.Command{}, []string{"123"})

		assert.NotNil(t, err)
		assert.Equal(t, "Error retrieving ACL [123]: test error", err.Error())
	})
}

func Test_loadbalancerACLConditionShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := loadbalancerACLConditionShowCmd(nil).Args(nil, []string{"123", "0"})

		assert.Nil(t, err)
	})

	t.Run("MissingACL_Error", func(t *testing.T) {
		err := loadbalancerACLConditionShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing ACL", err.Error())
	})

	t.Run("MissingIndex_Error", func(t *testing.T) {
		err := loadbalancerACLConditionShowCmd(nil).Args(nil, []string{"123"})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing ACL condition index", err.Error())
	})
}

func Test_loadbalancerACLConditionShow(t *testing.T) {
	t.Run("SingleACL", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().GetACL(123).Return(loadbalancer.ACL{}, nil).Times(1)

		err := loadbalancerACLConditionShow(service, &cobra.Command{}, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidACLID_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		err := loadbalancerACLConditionShow(service, &cobra.Command{}, []string{"abc", "0"})

		assert.NotNil(t, err)
		assert.Equal(t, "Invalid ACL ID [abc]", err.Error())
	})

	t.Run("GetACLError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().GetACL(123).Return(loadbalancer.ACL{}, errors.New("test error"))

		err := loadbalancerACLConditionShow(service, &cobra.Command{}, []string{"123", "0"})

		assert.NotNil(t, err)
		assert.Equal(t, "Error retrieving ACL [123]: test error", err.Error())
	})
}

func Test_loadbalancerACLConditionCreateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := loadbalancerACLConditionCreateCmd(nil).Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("MissingACL_Error", func(t *testing.T) {
		err := loadbalancerACLConditionCreateCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing ACL", err.Error())
	})
}

func Test_loadbalancerACLConditionCreate(t *testing.T) {
	t.Run("DefaultCreate_ExpectedPatch", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)
		cmd := loadbalancerACLConditionCreateCmd(nil)
		cmd.Flags().Set("name", "header_matches")
		cmd.Flags().Set("argument", "header=host")
		cmd.Flags().Set("argument", "value=test.com")

		req := loadbalancer.PatchACLRequest{
			Conditions: []loadbalancer.ACLCondition{
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

		err := loadbalancerACLConditionCreate(service, cmd, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("GetACLError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)
		cmd := loadbalancerACLConditionCreateCmd(nil)

		service.EXPECT().GetACL(gomock.Any()).Return(loadbalancer.ACL{}, errors.New("test error"))

		err := loadbalancerACLConditionCreate(service, cmd, []string{"123"})

		assert.Equal(t, "Error retrieving ACL: test error", err.Error())
	})

	t.Run("PatchACLError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)
		cmd := loadbalancerACLConditionCreateCmd(nil)

		gomock.InOrder(
			service.EXPECT().GetACL(123).Return(loadbalancer.ACL{}, nil),
			service.EXPECT().PatchACL(123, gomock.Any()).Return(errors.New("test error")),
		)

		err := loadbalancerACLConditionCreate(service, cmd, []string{"123"})

		assert.Equal(t, "Error updating ACL: test error", err.Error())
	})

	t.Run("UpdatedGetACLError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)
		cmd := loadbalancerACLConditionCreateCmd(nil)

		gomock.InOrder(
			service.EXPECT().GetACL(123).Return(loadbalancer.ACL{}, nil),
			service.EXPECT().PatchACL(123, gomock.Any()).Return(nil),
			service.EXPECT().GetACL(123).Return(loadbalancer.ACL{}, errors.New("test error")),
		)

		err := loadbalancerACLConditionCreate(service, cmd, []string{"123"})

		assert.Equal(t, "Error retrieving updated ACL [123]: test error", err.Error())
	})
}

func Test_loadbalancerACLConditionUpdateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := loadbalancerACLConditionUpdateCmd(nil).Args(nil, []string{"123", "0"})

		assert.Nil(t, err)
	})

	t.Run("MissingACL_Error", func(t *testing.T) {
		err := loadbalancerACLConditionUpdateCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing ACL", err.Error())
	})

	t.Run("MissingACLConditionIndex_Error", func(t *testing.T) {
		err := loadbalancerACLConditionUpdateCmd(nil).Args(nil, []string{"123"})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing ACL condition index", err.Error())
	})
}

func Test_loadbalancerACLConditionUpdate(t *testing.T) {
	t.Run("DefaultUpdate_ExpectedPatch", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)
		cmd := loadbalancerACLConditionUpdateCmd(nil)
		cmd.Flags().Set("name", "header_matches")

		req := loadbalancer.PatchACLRequest{
			Conditions: []loadbalancer.ACLCondition{
				{
					Name: "header_matches",
				},
			},
		}

		gomock.InOrder(
			service.EXPECT().GetACL(123).Return(loadbalancer.ACL{
				Conditions: []loadbalancer.ACLCondition{
					{
						Name: "redirect",
					},
				},
			}, nil),
			service.EXPECT().PatchACL(123, req).Return(nil),
			service.EXPECT().GetACL(123).Return(loadbalancer.ACL{}, nil),
		)

		err := loadbalancerACLConditionUpdate(service, cmd, []string{"123", "0"})
		assert.Nil(t, err)
	})

	t.Run("InvalidACLID_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)
		cmd := loadbalancerACLConditionUpdateCmd(nil)

		err := loadbalancerACLConditionUpdate(service, cmd, []string{"abc"})

		assert.NotNil(t, err)
		assert.Equal(t, "Invalid ACL ID [abc]", err.Error())
	})

	t.Run("GetACLError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)
		cmd := loadbalancerACLConditionUpdateCmd(nil)

		service.EXPECT().GetACL(123).Return(loadbalancer.ACL{}, errors.New("test error"))

		err := loadbalancerACLConditionUpdate(service, cmd, []string{"123"})

		assert.NotNil(t, err)
		assert.Equal(t, "Error retrieving ACL: test error", err.Error())
	})

	t.Run("InvalidConditionIndex_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)
		cmd := loadbalancerACLConditionUpdateCmd(nil)

		service.EXPECT().GetACL(123).Return(loadbalancer.ACL{}, nil)

		err := loadbalancerACLConditionUpdate(service, cmd, []string{"123", "abc"})

		assert.NotNil(t, err)
		assert.Equal(t, "Invalid ACL condition index [abc]", err.Error())
	})

	t.Run("ConditionIndexOutOfBounds_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)
		cmd := loadbalancerACLConditionUpdateCmd(nil)

		service.EXPECT().GetACL(123).Return(loadbalancer.ACL{
			Conditions: []loadbalancer.ACLCondition{
				{
					Name: "redirect",
				},
			},
		}, nil)

		err := loadbalancerACLConditionUpdate(service, cmd, []string{"123", "1"})

		assert.NotNil(t, err)
		assert.Equal(t, "ACL condition index [1] out of bounds", err.Error())
	})

	t.Run("PatchACLError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)
		cmd := loadbalancerACLConditionUpdateCmd(nil)
		cmd.Flags().Set("name", "header_matches")

		gomock.InOrder(
			service.EXPECT().GetACL(123).Return(loadbalancer.ACL{
				Conditions: []loadbalancer.ACLCondition{
					{
						Name: "redirect",
					},
				},
			}, nil),
			service.EXPECT().PatchACL(123, gomock.Any()).Return(errors.New("test error")),
		)

		err := loadbalancerACLConditionUpdate(service, cmd, []string{"123", "0"})

		assert.NotNil(t, err)
		assert.Equal(t, "Error updating ACL: test error", err.Error())
	})

	t.Run("UpdatedGetACLError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)
		cmd := loadbalancerACLConditionUpdateCmd(nil)
		cmd.Flags().Set("name", "header_matches")

		gomock.InOrder(
			service.EXPECT().GetACL(123).Return(loadbalancer.ACL{
				Conditions: []loadbalancer.ACLCondition{
					{
						Name: "redirect",
					},
				},
			}, nil),
			service.EXPECT().PatchACL(123, gomock.Any()).Return(nil),
			service.EXPECT().GetACL(123).Return(loadbalancer.ACL{}, errors.New("test error")),
		)

		err := loadbalancerACLConditionUpdate(service, cmd, []string{"123", "0"})

		assert.NotNil(t, err)
		assert.Equal(t, "Error retrieving updated ACL [123]: test error", err.Error())
	})
}

func Test_loadbalancerACLConditionDeleteCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := loadbalancerACLConditionDeleteCmd(nil).Args(nil, []string{"123", "0"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := loadbalancerACLConditionDeleteCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing ACL", err.Error())
	})

	t.Run("MissingACLConditionIndex_Error", func(t *testing.T) {
		err := loadbalancerACLConditionDeleteCmd(nil).Args(nil, []string{"123"})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing ACL condition index", err.Error())
	})
}

func Test_loadbalancerACLConditionDelete(t *testing.T) {
	t.Run("DefaultDelete_ExpectedPatch", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)
		cmd := loadbalancerACLConditionDeleteCmd(nil)
		cmd.Flags().Set("name", "header_matches")

		req := loadbalancer.PatchACLRequest{
			Conditions: make([]loadbalancer.ACLCondition, 0),
		}

		gomock.InOrder(
			service.EXPECT().GetACL(123).Return(loadbalancer.ACL{
				Conditions: []loadbalancer.ACLCondition{
					{
						Name: "redirect",
					},
				},
			}, nil),
			service.EXPECT().PatchACL(123, req).Return(nil),
			service.EXPECT().GetACL(123).Return(loadbalancer.ACL{}, nil),
		)

		err := loadbalancerACLConditionDelete(service, cmd, []string{"123", "0"})
		assert.Nil(t, err)
	})

	t.Run("InvalidACLID_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)
		cmd := loadbalancerACLConditionDeleteCmd(nil)

		err := loadbalancerACLConditionDelete(service, cmd, []string{"abc"})

		assert.NotNil(t, err)
		assert.Equal(t, "Invalid ACL ID [abc]", err.Error())
	})

	t.Run("GetACLError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)
		cmd := loadbalancerACLConditionDeleteCmd(nil)

		service.EXPECT().GetACL(123).Return(loadbalancer.ACL{}, errors.New("test error"))

		err := loadbalancerACLConditionDelete(service, cmd, []string{"123"})

		assert.NotNil(t, err)
		assert.Equal(t, "Error retrieving ACL: test error", err.Error())
	})

	t.Run("InvalidConditionIndex_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)
		cmd := loadbalancerACLConditionDeleteCmd(nil)

		service.EXPECT().GetACL(123).Return(loadbalancer.ACL{}, nil)

		err := loadbalancerACLConditionDelete(service, cmd, []string{"123", "abc"})

		assert.NotNil(t, err)
		assert.Equal(t, "Invalid ACL condition index [abc]", err.Error())
	})

	t.Run("ConditionIndexOutOfBounds_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)
		cmd := loadbalancerACLConditionDeleteCmd(nil)

		service.EXPECT().GetACL(123).Return(loadbalancer.ACL{
			Conditions: []loadbalancer.ACLCondition{
				{
					Name: "redirect",
				},
			},
		}, nil)

		err := loadbalancerACLConditionDelete(service, cmd, []string{"123", "1"})

		assert.NotNil(t, err)
		assert.Equal(t, "ACL condition index [1] out of bounds", err.Error())
	})

	t.Run("PatchACLError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)
		cmd := loadbalancerACLConditionDeleteCmd(nil)

		gomock.InOrder(
			service.EXPECT().GetACL(123).Return(loadbalancer.ACL{
				Conditions: []loadbalancer.ACLCondition{
					{
						Name: "redirect",
					},
				},
			}, nil),
			service.EXPECT().PatchACL(123, gomock.Any()).Return(errors.New("test error")),
		)

		err := loadbalancerACLConditionDelete(service, cmd, []string{"123", "0"})

		assert.NotNil(t, err)
		assert.Equal(t, "Error updating ACL: test error", err.Error())
	})

	t.Run("UpdatedGetACLError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)
		cmd := loadbalancerACLConditionDeleteCmd(nil)

		gomock.InOrder(
			service.EXPECT().GetACL(123).Return(loadbalancer.ACL{
				Conditions: []loadbalancer.ACLCondition{
					{
						Name: "redirect",
					},
				},
			}, nil),
			service.EXPECT().PatchACL(123, gomock.Any()).Return(nil),
			service.EXPECT().GetACL(123).Return(loadbalancer.ACL{}, errors.New("test error")),
		)

		err := loadbalancerACLConditionDelete(service, cmd, []string{"123", "0"})

		assert.NotNil(t, err)
		assert.Equal(t, "Error retrieving updated ACL [123]: test error", err.Error())
	})
}
