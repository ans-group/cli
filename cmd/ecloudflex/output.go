package ecloudflex

import (
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/ecloudflex"
)

func OutputECloudFlexProjectsProvider(projects []ecloudflex.Project) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(projects).WithDefaultFields([]string{"id", "name", "created_at"})
}
