package ecloudflex

import (
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/ukfast/cli/internal/pkg/clierrors"
	"github.com/ukfast/cli/test/mocks"
	"github.com/ukfast/cli/test/test_output"
	"github.com/ukfast/sdk-go/pkg/service/ecloudflex"
)

func Test_ecloudflexProjectList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudFlexService(mockCtrl)

		service.EXPECT().GetProjects(gomock.Any()).Return([]ecloudflex.Project{}, nil).Times(1)

		ecloudflexProjectList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudFlexService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := ecloudflexProjectList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetProjectsError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudFlexService(mockCtrl)

		service.EXPECT().GetProjects(gomock.Any()).Return([]ecloudflex.Project{}, errors.New("test error")).Times(1)

		err := ecloudflexProjectList(service, &cobra.Command{}, []string{})

		assert.Equal(t, "Error retrieving projects: test error", err.Error())
	})
}

func Test_ecloudflexProjectShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudflexProjectShowCmd(nil).Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudflexProjectShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing project", err.Error())
	})
}

func Test_ecloudflexProjectShow(t *testing.T) {
	t.Run("SingleProject", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudFlexService(mockCtrl)

		service.EXPECT().GetProject(123).Return(ecloudflex.Project{}, nil).Times(1)

		ecloudflexProjectShow(service, &cobra.Command{}, []string{"123"})
	})

	t.Run("MultipleProjects", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudFlexService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetProject(123).Return(ecloudflex.Project{}, nil),
			service.EXPECT().GetProject(456).Return(ecloudflex.Project{}, nil),
		)

		ecloudflexProjectShow(service, &cobra.Command{}, []string{"123", "456"})
	})

	t.Run("GetProjectID_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudFlexService(mockCtrl)

		test_output.AssertErrorOutput(t, "Invalid project ID [abc]\n", func() {
			ecloudflexProjectShow(service, &cobra.Command{}, []string{"abc"})
		})
	})

	t.Run("GetProjectError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudFlexService(mockCtrl)

		service.EXPECT().GetProject(123).Return(ecloudflex.Project{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving project [123]: test error\n", func() {
			ecloudflexProjectShow(service, &cobra.Command{}, []string{"123"})
		})
	})
}
