package ecloud

import (
	"errors"
	"testing"

	"github.com/ans-group/cli/internal/pkg/clierrors"
	"github.com/ans-group/cli/test/mocks"
	"github.com/ans-group/sdk-go/pkg/service/ecloud"
	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_ecloudTagList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetTags(gomock.Any()).Return([]ecloud.Tag{}, nil).Times(1)

		ecloudTagList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := ecloudTagList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetTagsError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetTags(gomock.Any()).Return([]ecloud.Tag{}, errors.New("test error")).Times(1)

		err := ecloudTagList(service, &cobra.Command{}, []string{})

		assert.Equal(t, "ecloud: Error retrieving tags: test error", err.Error())
	})
}

func Test_ecloudTagShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudTagShowCmd(nil).Args(nil, []string{"tag-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudTagShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing tag", err.Error())
	})
}

func Test_ecloudTagShow(t *testing.T) {
	t.Run("SingleTag", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetTag("tag-abcdef12").Return(ecloud.Tag{}, nil).Times(1)

		ecloudTagShow(service, &cobra.Command{}, []string{"tag-abcdef12"})
	})

	t.Run("MultipleTags", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetTag("tag-abcdef12").Return(ecloud.Tag{}, nil),
			service.EXPECT().GetTag("tag-abcdef23").Return(ecloud.Tag{}, nil),
		)

		ecloudTagShow(service, &cobra.Command{}, []string{"tag-abcdef12", "tag-abcdef23"})
	})

	t.Run("GetTagError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetTag("tag-abcdef12").Return(ecloud.Tag{}, errors.New("test error")).Times(1)

		ecloudTagShow(service, &cobra.Command{}, []string{"tag-abcdef12"})
	})
}

func Test_ecloudTagCreate(t *testing.T) {
	t.Run("DefaultCreate", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudTagCreateCmd(nil)
		cmd.Flags().Set("name", "production")
		cmd.Flags().Set("scope", "environment")

		expectedRequest := ecloud.CreateTagRequest{
			Name:  "production",
			Scope: "environment",
		}

		gomock.InOrder(
			service.EXPECT().CreateTag(expectedRequest).Return("tag-abcdef12", nil),
			service.EXPECT().GetTag("tag-abcdef12").Return(ecloud.Tag{}, nil),
		)

		ecloudTagCreate(service, cmd, []string{})
	})

	t.Run("CreateTagError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudTagCreateCmd(nil)
		cmd.Flags().Set("name", "production")
		cmd.Flags().Set("scope", "environment")

		service.EXPECT().CreateTag(gomock.Any()).Return("", errors.New("test error")).Times(1)

		err := ecloudTagCreate(service, cmd, []string{})

		assert.Equal(t, "ecloud: Error creating tag: test error", err.Error())
	})

	t.Run("GetTagError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudTagCreateCmd(nil)
		cmd.Flags().Set("name", "production")
		cmd.Flags().Set("scope", "environment")

		gomock.InOrder(
			service.EXPECT().CreateTag(gomock.Any()).Return("tag-abcdef12", nil),
			service.EXPECT().GetTag("tag-abcdef12").Return(ecloud.Tag{}, errors.New("test error")),
		)

		err := ecloudTagCreate(service, cmd, []string{})

		assert.Equal(t, "ecloud: Error retrieving new tag: test error", err.Error())
	})
}

func Test_ecloudTagUpdateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudTagUpdateCmd(nil).Args(nil, []string{"tag-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudTagUpdateCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing tag", err.Error())
	})
}

func Test_ecloudTagUpdate(t *testing.T) {
	t.Run("SingleTag", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudTagUpdateCmd(nil)
		cmd.Flags().Set("name", "staging")

		expectedRequest := ecloud.PatchTagRequest{
			Name: "staging",
		}

		gomock.InOrder(
			service.EXPECT().PatchTag("tag-abcdef12", expectedRequest).Return(nil),
			service.EXPECT().GetTag("tag-abcdef12").Return(ecloud.Tag{}, nil),
		)

		ecloudTagUpdate(service, cmd, []string{"tag-abcdef12"})
	})

	t.Run("MultipleTags", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudTagUpdateCmd(nil)
		cmd.Flags().Set("name", "staging")

		expectedRequest := ecloud.PatchTagRequest{
			Name: "staging",
		}

		gomock.InOrder(
			service.EXPECT().PatchTag("tag-abcdef12", expectedRequest).Return(nil),
			service.EXPECT().GetTag("tag-abcdef12").Return(ecloud.Tag{}, nil),
			service.EXPECT().PatchTag("tag-abcdef23", expectedRequest).Return(nil),
			service.EXPECT().GetTag("tag-abcdef23").Return(ecloud.Tag{}, nil),
		)

		ecloudTagUpdate(service, cmd, []string{"tag-abcdef12", "tag-abcdef23"})
	})

	t.Run("PatchTagError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudTagUpdateCmd(nil)
		cmd.Flags().Set("name", "staging")

		service.EXPECT().PatchTag("tag-abcdef12", gomock.Any()).Return(errors.New("test error")).Times(1)

		ecloudTagUpdate(service, cmd, []string{"tag-abcdef12"})
	})

	t.Run("GetTagError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudTagUpdateCmd(nil)
		cmd.Flags().Set("name", "staging")

		gomock.InOrder(
			service.EXPECT().PatchTag("tag-abcdef12", gomock.Any()).Return(nil),
			service.EXPECT().GetTag("tag-abcdef12").Return(ecloud.Tag{}, errors.New("test error")),
		)

		ecloudTagUpdate(service, cmd, []string{"tag-abcdef12"})
	})
}

func Test_ecloudTagDeleteCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudTagDeleteCmd(nil).Args(nil, []string{"tag-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudTagDeleteCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing tag", err.Error())
	})
}

func Test_ecloudTagDelete(t *testing.T) {
	t.Run("SingleTag", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().DeleteTag("tag-abcdef12").Return(nil).Times(1)

		ecloudTagDelete(service, &cobra.Command{}, []string{"tag-abcdef12"})
	})

	t.Run("MultipleTags", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().DeleteTag("tag-abcdef12").Return(nil),
			service.EXPECT().DeleteTag("tag-abcdef23").Return(nil),
		)

		ecloudTagDelete(service, &cobra.Command{}, []string{"tag-abcdef12", "tag-abcdef23"})
	})

	t.Run("DeleteTagError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().DeleteTag("tag-abcdef12").Return(errors.New("test error")).Times(1)

		ecloudTagDelete(service, &cobra.Command{}, []string{"tag-abcdef12"})
	})
}
