package ecloud

import (
	"errors"
	"testing"

	"github.com/ans-group/cli/internal/pkg/clierrors"
	"github.com/ans-group/cli/test/mocks"
	"github.com/ans-group/cli/test/test_output"
	"github.com/ans-group/sdk-go/pkg/service/ecloud"
	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_ecloudVirtualMachineTagListCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudVirtualMachineTagListCmd(nil).Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudVirtualMachineTagListCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing virtual machine", err.Error())
	})
}

func Test_ecloudVirtualMachineTagList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetVirtualMachineTags(123, gomock.Any()).Return([]ecloud.Tag{}, nil).Times(1)

		ecloudVirtualMachineTagList(service, &cobra.Command{}, []string{"123"})
	})

	t.Run("InvalidVirtualMachineID_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		err := ecloudVirtualMachineTagList(service, &cobra.Command{}, []string{"abc"})

		assert.Equal(t, "Invalid virtual machine ID [abc]", err.Error())
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := ecloudVirtualMachineTagList(service, cmd, []string{"123"})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetVirtualMachineTagsError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetVirtualMachineTags(123, gomock.Any()).Return([]ecloud.Tag{}, errors.New("test error 1")).Times(1)

		err := ecloudVirtualMachineTagList(service, &cobra.Command{}, []string{"123"})

		assert.Equal(t, "Error retrieving virtual machine tags: test error 1", err.Error())
	})
}

func Test_ecloudVirtualMachineTagShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudVirtualMachineTagShowCmd(nil).Args(nil, []string{"123", "testkey1"})

		assert.Nil(t, err)
	})

	t.Run("MissingVirtualMachine_Error", func(t *testing.T) {
		err := ecloudVirtualMachineTagShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing virtual machine", err.Error())
	})

	t.Run("MissingTag_Error", func(t *testing.T) {
		err := ecloudVirtualMachineTagShowCmd(nil).Args(nil, []string{"123"})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing tag", err.Error())
	})
}

func Test_ecloudVirtualMachineTagShow(t *testing.T) {
	t.Run("RetrieveSingle", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetVirtualMachineTag(123, "testkey1").Return(ecloud.Tag{}, nil).Times(1)

		ecloudVirtualMachineTagShow(service, &cobra.Command{}, []string{"123", "testkey1"})
	})

	t.Run("RetrieveMultiple", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetVirtualMachineTag(123, "testkey1").Return(ecloud.Tag{}, nil),
			service.EXPECT().GetVirtualMachineTag(123, "testkey2").Return(ecloud.Tag{}, nil),
		)

		ecloudVirtualMachineTagShow(service, &cobra.Command{}, []string{"123", "testkey1", "testkey2"})
	})

	t.Run("InvalidVirtualMachineID_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		err := ecloudVirtualMachineTagShow(service, &cobra.Command{}, []string{"abc", "testkey1"})

		assert.Equal(t, "Invalid virtual machine ID [abc]", err.Error())
	})

	t.Run("GetVirtualMachineTagError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetVirtualMachineTag(123, "testkey1").Return(ecloud.Tag{}, errors.New("test error 1")).Times(1)

		test_output.AssertErrorOutput(t, "Error retrieving virtual machine tag [testkey1]: test error 1\n", func() {
			ecloudVirtualMachineTagShow(service, &cobra.Command{}, []string{"123", "testkey1"})
		})
	})
}

func Test_ecloudVirtualMachineTagCreateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudVirtualMachineTagCreateCmd(nil).Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudVirtualMachineTagCreateCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing virtual machine", err.Error())
	})
}

func Test_ecloudVirtualMachineTagCreate(t *testing.T) {
	t.Run("DefaultCreate", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudVirtualMachineTagCreateCmd(nil)
		cmd.Flags().Set("key", "testkey1")
		cmd.Flags().Set("value", "testvalue1")

		expectedRequest := ecloud.CreateTagRequest{
			Key:   "testkey1",
			Value: "testvalue1",
		}

		gomock.InOrder(
			service.EXPECT().CreateVirtualMachineTag(123, gomock.Eq(expectedRequest)).Return(nil),
			service.EXPECT().GetVirtualMachineTag(123, "testkey1").Return(ecloud.Tag{}, nil),
		)

		ecloudVirtualMachineTagCreate(service, cmd, []string{"123"})
	})

	t.Run("InvalidVirtualMachineID_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		err := ecloudVirtualMachineTagCreate(service, &cobra.Command{}, []string{"abc"})

		assert.Equal(t, "Invalid virtual machine ID [abc]", err.Error())
	})

	t.Run("CreateVirtualMachineTagError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().CreateVirtualMachineTag(123, gomock.Any()).Return(errors.New("test error 1")).Times(1)

		err := ecloudVirtualMachineTagCreate(service, &cobra.Command{}, []string{"123"})

		assert.Equal(t, "Error creating virtual machine tag: test error 1", err.Error())
	})

	t.Run("GetVirtualMachineTagError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudVirtualMachineTagCreateCmd(nil)
		cmd.Flags().Set("key", "testkey1")
		cmd.Flags().Set("value", "testvalue1")

		gomock.InOrder(
			service.EXPECT().CreateVirtualMachineTag(123, gomock.Any()).Return(nil),
			service.EXPECT().GetVirtualMachineTag(123, "testkey1").Return(ecloud.Tag{}, errors.New("test error 1")),
		)

		err := ecloudVirtualMachineTagCreate(service, cmd, []string{"123"})

		assert.Equal(t, "Error retrieving new virtual machine tag: test error 1", err.Error())
	})
}

func Test_ecloudVirtualMachineTagUpdateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudVirtualMachineTagUpdateCmd(nil).Args(nil, []string{"123", "testkey1"})

		assert.Nil(t, err)
	})

	t.Run("MissingVirtualMachine_Error", func(t *testing.T) {
		err := ecloudVirtualMachineTagUpdateCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing virtual machine", err.Error())
	})

	t.Run("MissingTag_Error", func(t *testing.T) {
		err := ecloudVirtualMachineTagUpdateCmd(nil).Args(nil, []string{"123"})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing tag", err.Error())
	})
}

func Test_ecloudVirtualMachineTagUpdate(t *testing.T) {
	t.Run("UpdateSingle", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudVirtualMachineTagUpdateCmd(nil)
		cmd.Flags().Set("value", "testvalue1")

		expectedPatch := ecloud.PatchTagRequest{
			Value: "testvalue1",
		}

		gomock.InOrder(
			service.EXPECT().PatchVirtualMachineTag(123, "testkey1", gomock.Eq(expectedPatch)).Return(nil),
			service.EXPECT().GetVirtualMachineTag(123, "testkey1").Return(ecloud.Tag{}, nil),
		)

		ecloudVirtualMachineTagUpdate(service, cmd, []string{"123", "testkey1"})
	})

	t.Run("UpdateMultiple", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudVirtualMachineTagUpdateCmd(nil)
		cmd.Flags().Set("value", "testvalue1")

		expectedPatch := ecloud.PatchTagRequest{
			Value: "testvalue1",
		}

		gomock.InOrder(
			service.EXPECT().PatchVirtualMachineTag(123, "testkey1", gomock.Eq(expectedPatch)).Return(nil),
			service.EXPECT().GetVirtualMachineTag(123, "testkey1").Return(ecloud.Tag{}, nil),
			service.EXPECT().PatchVirtualMachineTag(123, "testkey2", gomock.Eq(expectedPatch)).Return(nil),
			service.EXPECT().GetVirtualMachineTag(123, "testkey2").Return(ecloud.Tag{}, nil),
		)

		ecloudVirtualMachineTagUpdate(service, cmd, []string{"123", "testkey1", "testkey2"})
	})

	t.Run("InvalidVirtualMachineID_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		err := ecloudVirtualMachineTagUpdate(service, &cobra.Command{}, []string{"abc"})

		assert.Equal(t, "Invalid virtual machine ID [abc]", err.Error())
	})

	t.Run("PatchVirtualMachineTag_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().PatchVirtualMachineTag(123, "testkey1", gomock.Any()).Return(errors.New("test error 1")),
		)

		test_output.AssertErrorOutput(t, "Error updating virtual machine tag [testkey1]: test error 1\n", func() {
			ecloudVirtualMachineTagUpdate(service, &cobra.Command{}, []string{"123", "testkey1"})
		})
	})

	t.Run("GetVirtualMachineTagError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudVirtualMachineTagUpdateCmd(nil)
		cmd.Flags().Set("value", "testvalue1")

		gomock.InOrder(
			service.EXPECT().PatchVirtualMachineTag(123, "testkey1", gomock.Any()).Return(nil),
			service.EXPECT().GetVirtualMachineTag(123, "testkey1").Return(ecloud.Tag{}, errors.New("test error 1")),
		)

		test_output.AssertErrorOutput(t, "Error retrieving updated virtual machine tag [testkey1]: test error 1\n", func() {
			ecloudVirtualMachineTagUpdate(service, cmd, []string{"123", "testkey1"})
		})
	})
}

func Test_ecloudVirtualMachineTagDeleteCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudVirtualMachineTagDeleteCmd(nil).Args(nil, []string{"123", "testkey1"})

		assert.Nil(t, err)
	})

	t.Run("MissingVirtualMachine_Error", func(t *testing.T) {
		err := ecloudVirtualMachineTagDeleteCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing virtual machine", err.Error())
	})

	t.Run("MissingTag_Error", func(t *testing.T) {
		err := ecloudVirtualMachineTagDeleteCmd(nil).Args(nil, []string{"123"})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing tag", err.Error())
	})
}

func Test_ecloudVirtualMachineTagDelete(t *testing.T) {
	t.Run("RetrieveSingle", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().DeleteVirtualMachineTag(123, "testkey1").Return(nil).Times(1)

		ecloudVirtualMachineTagDelete(service, &cobra.Command{}, []string{"123", "testkey1"})
	})

	t.Run("RetrieveMultiple", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().DeleteVirtualMachineTag(123, "testkey1").Return(nil),
			service.EXPECT().DeleteVirtualMachineTag(123, "testkey2").Return(nil),
		)

		ecloudVirtualMachineTagDelete(service, &cobra.Command{}, []string{"123", "testkey1", "testkey2"})
	})

	t.Run("InvalidVirtualMachineID_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		err := ecloudVirtualMachineTagDelete(service, &cobra.Command{}, []string{"abc", "testkey1"})

		assert.Equal(t, "Invalid virtual machine ID [abc]", err.Error())
	})

	t.Run("DeleteVirtualMachineTagError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().DeleteVirtualMachineTag(123, "testkey1").Return(errors.New("test error 1")).Times(1)

		test_output.AssertErrorOutput(t, "Error removing virtual machine tag [testkey1]: test error 1\n", func() {
			ecloudVirtualMachineTagDelete(service, &cobra.Command{}, []string{"123", "testkey1"})
		})
	})
}
