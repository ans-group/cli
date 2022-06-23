package pss

import (
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/pss"
)

func OutputPSSRequestsProvider(requests []pss.Request) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(requests).WithDefaultFields([]string{"id", "author_name", "type", "subject", "status", "last_replied_at", "created_at"})
}

func OutputPSSRepliesProvider(replies []pss.Reply) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(replies).WithDefaultFields([]string{"id", "author_name", "created_at"})
}

func OutputPSSFeedbackProvider(feedback []pss.Feedback) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(feedback).WithDefaultFields([]string{"id", "score", "created_at"})
}
