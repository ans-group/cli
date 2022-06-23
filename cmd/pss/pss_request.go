package pss

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/helper"
	"github.com/ans-group/cli/internal/pkg/input"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/pss"
	"github.com/spf13/cobra"
)

func pssRequestRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "request",
		Short: "sub-commands relating to requests",
	}

	// Child commands
	cmd.AddCommand(pssRequestListCmd(f))
	cmd.AddCommand(pssRequestShowCmd(f))
	cmd.AddCommand(pssRequestCreateCmd(f))
	cmd.AddCommand(pssRequestUpdateCmd(f))
	cmd.AddCommand(pssRequestCloseCmd(f))

	// Child root commands
	cmd.AddCommand(pssRequestReplyRootCmd(f))
	cmd.AddCommand(pssRequestFeedbackRootCmd(f))

	return cmd
}

func pssRequestListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists requests",
		Long:    "This command lists requests",
		Example: "ans pss request list",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return pssRequestList(c.PSSService(), cmd, args)
		},
	}
}

func pssRequestList(service pss.PSSService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	requests, err := service.GetRequests(params)
	if err != nil {
		return err
	}

	return output.CommandOutput(cmd, OutputPSSRequestsProvider(requests))
}

func pssRequestShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <request: id>...",
		Short:   "Shows a request",
		Long:    "This command shows one or more requests",
		Example: "ans pss request show 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing request")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return pssRequestShow(c.PSSService(), cmd, args)
		},
	}
}

func pssRequestShow(service pss.PSSService, cmd *cobra.Command, args []string) error {
	var requests []pss.Request
	for _, arg := range args {
		requestID, err := strconv.Atoi(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Invalid request ID [%s]", arg)
			continue
		}

		request, err := service.GetRequest(requestID)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving request [%s]: %s", arg, err)
			continue
		}

		requests = append(requests, request)
	}

	return output.CommandOutput(cmd, OutputPSSRequestsProvider(requests))
}

func pssRequestCreateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Creates a request",
		Long:    "This command creates a new request",
		Example: "ans pss request create --subject 'example ticket' --details 'example' --author 123",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return pssRequestCreate(c.PSSService(), cmd, args)
		},
	}

	// Setup flags
	cmd.Flags().String("subject", "", "Specifies subject for request")
	cmd.MarkFlagRequired("subject")
	cmd.Flags().String("details", "", "Specifies details for request")
	cmd.Flags().Int("author", 0, "Specifies author ID for request")
	cmd.MarkFlagRequired("author")
	cmd.Flags().String("priority", "Normal", "Specifies priority for request")
	cmd.Flags().Bool("secure", false, "Specifies whether request is secure")
	cmd.Flags().StringSlice("cc", []string{}, "Specifies CC email addresses for request")
	cmd.Flags().Bool("request-sms", false, "Specifies whether SMS updates are required")
	cmd.Flags().String("customer-reference", "", "Specifies customer reference for request")
	cmd.Flags().Int("product-id", 0, "Specifies product ID for request")
	cmd.Flags().String("product-name", "", "Specifies product name for request")
	cmd.Flags().String("product-type", "", "Specifies product type for request")

	return cmd
}

func pssRequestCreate(service pss.PSSService, cmd *cobra.Command, args []string) error {
	createRequest := pss.CreateRequestRequest{}

	priority, _ := cmd.Flags().GetString("priority")
	parsedPriority, err := pss.ParseRequestPriority(priority)
	if err != nil {
		return err
	}
	createRequest.Priority = parsedPriority

	if cmd.Flags().Changed("product-id") || cmd.Flags().Changed("product-name") || cmd.Flags().Changed("product-type") {
		createRequest.Product = &pss.Product{}
		createRequest.Product.ID, _ = cmd.Flags().GetInt("product-id")
		createRequest.Product.Name, _ = cmd.Flags().GetString("product-name")
		createRequest.Product.Type, _ = cmd.Flags().GetString("product-type")
	}

	createRequest.Subject, _ = cmd.Flags().GetString("subject")
	createRequest.Author.ID, _ = cmd.Flags().GetInt("author")
	createRequest.Secure, _ = cmd.Flags().GetBool("secure")
	createRequest.CC, _ = cmd.Flags().GetStringSlice("cc")
	createRequest.RequestSMS, _ = cmd.Flags().GetBool("request-sms")
	createRequest.CustomerReference, _ = cmd.Flags().GetString("customer-reference")

	if cmd.Flags().Changed("details") {
		createRequest.Details, _ = cmd.Flags().GetString("details")
	} else {
		createRequest.Details, err = input.ReadInput("details")
		if err != nil {
			return err
		}
	}

	requestID, err := service.CreateRequest(createRequest)
	if err != nil {
		return fmt.Errorf("Error creating request: %s", err)
	}

	request, err := service.GetRequest(requestID)
	if err != nil {
		return fmt.Errorf("Error retrieving new request: %s", err)
	}

	return output.CommandOutput(cmd, OutputPSSRequestsProvider([]pss.Request{request}))
}

func pssRequestUpdateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <request: id>...",
		Short:   "Updates requests",
		Long:    "This command updates one or more requests",
		Example: "ans pss request update 123 --priority high",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing request")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return pssRequestUpdate(c.PSSService(), cmd, args)
		},
	}

	// Setup flags
	cmd.Flags().String("priority", "", "Specifies priority for request")
	cmd.Flags().String("status", "", "Specifies status for request")
	cmd.Flags().Bool("secure", false, "Specifies whether request is secure")
	cmd.Flags().Bool("read", false, "Specifies whether request is marked as read")
	cmd.Flags().Bool("request-sms", false, "Specifies whether SMS updates are required")
	cmd.Flags().Bool("archived", false, "Specifies whether request is archived")

	return cmd
}

func pssRequestUpdate(service pss.PSSService, cmd *cobra.Command, args []string) error {
	patchRequest := pss.PatchRequestRequest{}

	if cmd.Flags().Changed("priority") {
		priority, _ := cmd.Flags().GetString("priority")
		parsedPriority, err := pss.ParseRequestPriority(priority)
		if err != nil {
			return err
		}
		patchRequest.Priority = parsedPriority
	}

	if cmd.Flags().Changed("status") {
		status, _ := cmd.Flags().GetString("status")
		parsedStatus, err := pss.ParseRequestStatus(status)
		if err != nil {
			return err
		}
		patchRequest.Status = parsedStatus
	}

	if cmd.Flags().Changed("secure") {
		secure, _ := cmd.Flags().GetBool("secure")
		patchRequest.Secure = &secure
	}
	if cmd.Flags().Changed("read") {
		read, _ := cmd.Flags().GetBool("read")
		patchRequest.Read = &read
	}
	if cmd.Flags().Changed("request-sms") {
		requestSMS, _ := cmd.Flags().GetBool("request-sms")
		patchRequest.RequestSMS = &requestSMS
	}
	if cmd.Flags().Changed("archived") {
		archived, _ := cmd.Flags().GetBool("archived")
		patchRequest.Archived = &archived
	}

	var requests []pss.Request

	for _, arg := range args {
		requestID, err := strconv.Atoi(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Invalid request ID [%s]", arg)
			continue
		}

		err = service.PatchRequest(requestID, patchRequest)
		if err != nil {
			output.OutputWithErrorLevelf("Error updating request [%d]: %s", requestID, err)
			continue
		}

		request, err := service.GetRequest(requestID)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving updated request [%d]: %s", requestID, err)
			continue
		}

		requests = append(requests, request)
	}

	return output.CommandOutput(cmd, OutputPSSRequestsProvider(requests))
}

func pssRequestCloseCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "close <request: id>...",
		Short:   "Closes requests",
		Long:    "This command closes one or more requests",
		Example: "ans pss request close 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing request")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return pssRequestClose(c.PSSService(), cmd, args)
		},
	}
}

func pssRequestClose(service pss.PSSService, cmd *cobra.Command, args []string) error {
	patchRequest := pss.PatchRequestRequest{
		Status: pss.RequestStatusCompleted,
	}

	var requests []pss.Request

	for _, arg := range args {
		requestID, err := strconv.Atoi(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Invalid request ID [%s]", arg)
			continue
		}

		err = service.PatchRequest(requestID, patchRequest)
		if err != nil {
			output.OutputWithErrorLevelf("Error closing request [%d]: %s", requestID, err)
			continue
		}

		request, err := service.GetRequest(requestID)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving updated request [%d]: %s", requestID, err)
			continue
		}

		requests = append(requests, request)
	}

	return output.CommandOutput(cmd, OutputPSSRequestsProvider(requests))
}
