package pss

import (
	"strconv"

	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/pss"
)

func OutputPSSRequestsProvider(requests []pss.Request) output.OutputHandlerProvider {
	return output.NewGenericOutputHandlerProvider(
		output.WithData(requests),
		output.WithFieldDataFunc(func() ([]*output.OrderedFields, error) {
			var data []*output.OrderedFields
			for _, request := range requests {
				fields := output.NewOrderedFields()
				fields.Set("id", output.NewFieldValue(strconv.Itoa(request.ID), true))
				fields.Set("author_id", output.NewFieldValue(strconv.Itoa(request.Author.ID), false))
				fields.Set("author_name", output.NewFieldValue(request.Author.Name, true))
				fields.Set("type", output.NewFieldValue(request.Type, true))
				fields.Set("secure", output.NewFieldValue(strconv.FormatBool(request.Secure), false))
				fields.Set("subject", output.NewFieldValue(request.Subject, true))
				fields.Set("created_at", output.NewFieldValue(request.CreatedAt.String(), true))
				fields.Set("priority", output.NewFieldValue(request.Priority.String(), false))
				fields.Set("archived", output.NewFieldValue(strconv.FormatBool(request.Archived), false))
				fields.Set("status", output.NewFieldValue(request.Status.String(), true))
				fields.Set("request_sms", output.NewFieldValue(strconv.FormatBool(request.RequestSMS), false))
				fields.Set("version", output.NewFieldValue(strconv.Itoa(request.Version), false))
				fields.Set("customer_reference", output.NewFieldValue(request.CustomerReference, false))
				fields.Set("last_replied_at", output.NewFieldValue(request.LastRepliedAt.String(), true))
				fields.Set("product_id", output.NewFieldValue(strconv.Itoa(request.Product.ID), false))
				fields.Set("product_name", output.NewFieldValue(request.Product.Name, false))
				fields.Set("product_type", output.NewFieldValue(request.Product.Type, false))

				data = append(data, fields)
			}

			return data, nil
		}),
	)
}

func OutputPSSRepliesProvider(replies []pss.Reply) output.OutputHandlerProvider {
	return output.NewGenericOutputHandlerProvider(
		output.WithData(replies),
		output.WithFieldDataFunc(func() ([]*output.OrderedFields, error) {
			var data []*output.OrderedFields
			for _, reply := range replies {
				fields := output.NewOrderedFields()
				fields.Set("id", output.NewFieldValue(reply.ID, true))
				fields.Set("author_name", output.NewFieldValue(reply.Author.Name, true))
				fields.Set("description", output.NewFieldValue(reply.Description, false))
				fields.Set("created_at", output.NewFieldValue(reply.CreatedAt.String(), true))

				data = append(data, fields)
			}

			return data, nil
		}),
	)
}
