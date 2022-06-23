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

func Test_ecloudSolutionTagListCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudSolutionTagListCmd(nil).Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudSolutionTagListCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing solution", err.Error())
	})
}

func Test_ecloudSolutionTagList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetSolutionTags(123, gomock.Any()).Return([]ecloud.Tag{}, nil).Times(1)

		ecloudSolutionTagList(service, &cobra.Command{}, []string{"123"})
	})

	t.Run("InvalidSolutionID_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		err := ecloudSolutionTagList(service, &cobra.Command{}, []string{"abc"})

		assert.Equal(t, "Invalid solution ID [abc]", err.Error())
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := ecloudSolutionTagList(service, cmd, []string{"123"})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetSolutionTagsError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetSolutionTags(123, gomock.Any()).Return([]ecloud.Tag{}, errors.New("test error 1")).Times(1)

		err := ecloudSolutionTagList(service, &cobra.Command{}, []string{"123"})

		assert.Equal(t, "Error retrieving solution tags: test error 1", err.Error())
	})
}

func Test_ecloudSolutionTagShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudSolutionTagShowCmd(nil).Args(nil, []string{"123", "testkey1"})

		assert.Nil(t, err)
	})

	t.Run("MissingSolution_Error", func(t *testing.T) {
		err := ecloudSolutionTagShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing solution", err.Error())
	})

	t.Run("MissingTag_Error", func(t *testing.T) {
		err := ecloudSolutionTagShowCmd(nil).Args(nil, []string{"123"})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing tag", err.Error())
	})
}

func Test_ecloudSolutionTagShow(t *testing.T) {
	t.Run("RetrieveSingle", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetSolutionTag(123, "testkey1").Return(ecloud.Tag{}, nil).Times(1)

		ecloudSolutionTagShow(service, &cobra.Command{}, []string{"123", "testkey1"})
	})

	t.Run("RetrieveMultiple", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetSolutionTag(123, "testkey1").Return(ecloud.Tag{}, nil),
			service.EXPECT().GetSolutionTag(123, "testkey2").Return(ecloud.Tag{}, nil),
		)

		ecloudSolutionTagShow(service, &cobra.Command{}, []string{"123", "testkey1", "testkey2"})
	})

	t.Run("InvalidSolutionID_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		err := ecloudSolutionTagShow(service, &cobra.Command{}, []string{"abc", "testkey1"})

		assert.Equal(t, "Invalid solution ID [abc]", err.Error())
	})

	t.Run("GetSolutionTagError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetSolutionTag(123, "testkey1").Return(ecloud.Tag{}, errors.New("test error 1")).Times(1)

		test_output.AssertErrorOutput(t, "Error retrieving solution tag [testkey1]: test error 1\n", func() {
			ecloudSolutionTagShow(service, &cobra.Command{}, []string{"123", "testkey1"})
		})
	})
}

func Test_ecloudSolutionTagCreateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudSolutionTagCreateCmd(nil).Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudSolutionTagCreateCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing solution", err.Error())
	})
}

func Test_ecloudSolutionTagCreate(t *testing.T) {
	t.Run("DefaultCreate", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudSolutionTagCreateCmd(nil)
		cmd.Flags().Set("key", "testkey1")
		cmd.Flags().Set("value", "testvalue1")

		expectedRequest := ecloud.CreateTagRequest{
			Key:   "testkey1",
			Value: "testvalue1",
		}

		gomock.InOrder(
			service.EXPECT().CreateSolutionTag(123, gomock.Eq(expectedRequest)).Return(nil),
			service.EXPECT().GetSolutionTag(123, "testkey1").Return(ecloud.Tag{}, nil),
		)

		ecloudSolutionTagCreate(service, cmd, []string{"123"})
	})

	t.Run("InvalidSolutionID_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		err := ecloudSolutionTagCreate(service, &cobra.Command{}, []string{"abc"})

		assert.Equal(t, "Invalid solution ID [abc]", err.Error())
	})

	t.Run("CreateSolutionTagError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().CreateSolutionTag(123, gomock.Any()).Return(errors.New("test error 1")).Times(1)

		err := ecloudSolutionTagCreate(service, &cobra.Command{}, []string{"123"})

		assert.Equal(t, "Error creating solution tag: test error 1", err.Error())
	})

	t.Run("GetSolutionTagError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudSolutionTagCreateCmd(nil)
		cmd.Flags().Set("key", "testkey1")
		cmd.Flags().Set("value", "testvalue1")

		gomock.InOrder(
			service.EXPECT().CreateSolutionTag(123, gomock.Any()).Return(nil),
			service.EXPECT().GetSolutionTag(123, "testkey1").Return(ecloud.Tag{}, errors.New("test error 1")),
		)

		err := ecloudSolutionTagCreate(service, cmd, []string{"123"})

		assert.Equal(t, "Error retrieving new solution tag: test error 1", err.Error())
	})
}

func Test_ecloudSolutionTagUpdateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudSolutionTagUpdateCmd(nil).Args(nil, []string{"123", "testkey1"})

		assert.Nil(t, err)
	})

	t.Run("MissingSolution_Error", func(t *testing.T) {
		err := ecloudSolutionTagUpdateCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing solution", err.Error())
	})

	t.Run("MissingTag_Error", func(t *testing.T) {
		err := ecloudSolutionTagUpdateCmd(nil).Args(nil, []string{"123"})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing tag", err.Error())
	})
}

func Test_ecloudSolutionTagUpdate(t *testing.T) {
	t.Run("UpdateSingle", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudSolutionTagUpdateCmd(nil)
		cmd.Flags().Set("value", "testvalue1")

		expectedPatch := ecloud.PatchTagRequest{
			Value: "testvalue1",
		}

		gomock.InOrder(
			service.EXPECT().PatchSolutionTag(123, "testkey1", gomock.Eq(expectedPatch)).Return(nil),
			service.EXPECT().GetSolutionTag(123, "testkey1").Return(ecloud.Tag{}, nil),
		)

		ecloudSolutionTagUpdate(service, cmd, []string{"123", "testkey1"})
	})

	t.Run("UpdateMultiple", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudSolutionTagUpdateCmd(nil)
		cmd.Flags().Set("value", "testvalue1")

		expectedPatch := ecloud.PatchTagRequest{
			Value: "testvalue1",
		}

		gomock.InOrder(
			service.EXPECT().PatchSolutionTag(123, "testkey1", gomock.Eq(expectedPatch)).Return(nil),
			service.EXPECT().GetSolutionTag(123, "testkey1").Return(ecloud.Tag{}, nil),
			service.EXPECT().PatchSolutionTag(123, "testkey2", gomock.Eq(expectedPatch)).Return(nil),
			service.EXPECT().GetSolutionTag(123, "testkey2").Return(ecloud.Tag{}, nil),
		)

		ecloudSolutionTagUpdate(service, cmd, []string{"123", "testkey1", "testkey2"})
	})

	t.Run("InvalidSolutionID_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		err := ecloudSolutionTagUpdate(service, &cobra.Command{}, []string{"abc"})

		assert.Equal(t, "Invalid solution ID [abc]", err.Error())
	})

	t.Run("PatchSolutionTag_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().PatchSolutionTag(123, "testkey1", gomock.Any()).Return(errors.New("test error 1")),
		)

		test_output.AssertErrorOutput(t, "Error updating solution tag [testkey1]: test error 1\n", func() {
			ecloudSolutionTagUpdate(service, &cobra.Command{}, []string{"123", "testkey1"})
		})
	})

	t.Run("GetSolutionTagError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudSolutionTagUpdateCmd(nil)
		cmd.Flags().Set("value", "testvalue1")

		gomock.InOrder(
			service.EXPECT().PatchSolutionTag(123, "testkey1", gomock.Any()).Return(nil),
			service.EXPECT().GetSolutionTag(123, "testkey1").Return(ecloud.Tag{}, errors.New("test error 1")),
		)

		test_output.AssertErrorOutput(t, "Error retrieving updated solution tag [testkey1]: test error 1\n", func() {
			ecloudSolutionTagUpdate(service, cmd, []string{"123", "testkey1"})
		})
	})
}

func Test_ecloudSolutionTagDeleteCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudSolutionTagDeleteCmd(nil).Args(nil, []string{"123", "testkey1"})

		assert.Nil(t, err)
	})

	t.Run("MissingSolution_Error", func(t *testing.T) {
		err := ecloudSolutionTagDeleteCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing solution", err.Error())
	})

	t.Run("MissingTag_Error", func(t *testing.T) {
		err := ecloudSolutionTagDeleteCmd(nil).Args(nil, []string{"123"})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing tag", err.Error())
	})
}

func Test_ecloudSolutionTagDelete(t *testing.T) {
	t.Run("RetrieveSingle", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().DeleteSolutionTag(123, "testkey1").Return(nil).Times(1)

		ecloudSolutionTagDelete(service, &cobra.Command{}, []string{"123", "testkey1"})
	})

	t.Run("RetrieveMultiple", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().DeleteSolutionTag(123, "testkey1").Return(nil),
			service.EXPECT().DeleteSolutionTag(123, "testkey2").Return(nil),
		)

		ecloudSolutionTagDelete(service, &cobra.Command{}, []string{"123", "testkey1", "testkey2"})
	})

	t.Run("InvalidSolutionID_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		err := ecloudSolutionTagDelete(service, &cobra.Command{}, []string{"abc", "testkey1"})

		assert.Equal(t, "Invalid solution ID [abc]", err.Error())
	})

	t.Run("DeleteSolutionTagError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().DeleteSolutionTag(123, "testkey1").Return(errors.New("test error 1")).Times(1)

		test_output.AssertErrorOutput(t, "Error removing solution tag [testkey1]: test error 1\n", func() {
			ecloudSolutionTagDelete(service, &cobra.Command{}, []string{"123", "testkey1"})
		})
	})
}
