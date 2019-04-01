package cmd

import (
	"strconv"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/pss"
)

func pssRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pss",
		Short: "Commands relating to PSS service",
	}

	// Child root commands
	cmd.AddCommand(pssRequestRootCmd())

	return cmd
}

// OutputPSSRequests implements OutputDataProvider for outputting an array of Requests
type OutputPSSRequests struct {
	Requests []pss.Request
}

func outputPSSRequests(requests []pss.Request) {
	err := Output(&OutputPSSRequests{Requests: requests})
	if err != nil {
		output.Fatalf("Failed to output requests: %s", err)
	}
}

func (o *OutputPSSRequests) GetData() interface{} {
	return o.Requests
}

func (o *OutputPSSRequests) GetFieldData() ([]*output.OrderedFields, error) {
	var data []*output.OrderedFields
	for _, request := range o.Requests {
		fields := o.getOrderedFields(request)
		data = append(data, fields)
	}

	return data, nil
}

func (o *OutputPSSRequests) getOrderedFields(request pss.Request) *output.OrderedFields {

	var assigneeName string
	if request.Assignee != nil {
		assigneeName = request.Assignee.Name
	}

	fields := output.NewOrderedFields()
	fields.Set("id", output.NewFieldValue(strconv.Itoa(request.ID), true))
	fields.Set("assignee_name", output.NewFieldValue(assigneeName, true))
	fields.Set("author_name", output.NewFieldValue(request.Author.Name, true))
	fields.Set("type", output.NewFieldValue(request.Type, true))
	fields.Set("secure", output.NewFieldValue(strconv.FormatBool(request.Secure), false))
	fields.Set("subject", output.NewFieldValue(request.Subject, true))
	fields.Set("created_at", output.NewFieldValue(request.CreatedAt.String(), true))
	fields.Set("priority", output.NewFieldValue(request.Priority.String(), false))
	fields.Set("archived", output.NewFieldValue(strconv.FormatBool(request.Archived), false))
	fields.Set("status", output.NewFieldValue(request.Status.String(), true))
	fields.Set("request_sms", output.NewFieldValue(strconv.FormatBool(request.RequestSMS), false))

	return fields
}

// OutputPSSReplies implements OutputDataProvider for outputting an array of Replies
type OutputPSSReplies struct {
	Replies []pss.Reply
}

func outputPSSReplies(replies []pss.Reply) {
	err := Output(&OutputPSSReplies{Replies: replies})
	if err != nil {
		output.Fatalf("Failed to output replies: %s", err)
	}
}

func (o *OutputPSSReplies) GetData() interface{} {
	return o.Replies
}

func (o *OutputPSSReplies) GetFieldData() ([]*output.OrderedFields, error) {
	var data []*output.OrderedFields
	for _, reply := range o.Replies {
		fields := o.getOrderedFields(reply)
		data = append(data, fields)
	}

	return data, nil
}

func (o *OutputPSSReplies) getOrderedFields(reply pss.Reply) *output.OrderedFields {
	fields := output.NewOrderedFields()
	fields.Set("author_name", output.NewFieldValue(reply.Author.Name, true))
	fields.Set("description", output.NewFieldValue(reply.Description, false))
	fields.Set("created_at", output.NewFieldValue(reply.CreatedAt.String(), true))

	return fields
}
