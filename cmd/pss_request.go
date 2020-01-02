package cmd

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/helper"
	"github.com/ukfast/cli/internal/pkg/input"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/pss"
)

func pssRequestRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "request",
		Short: "sub-commands relating to requests",
	}

	// Child commands
	cmd.AddCommand(pssRequestListCmd())
	cmd.AddCommand(pssRequestShowCmd())
	cmd.AddCommand(pssRequestCreateCmd())
	cmd.AddCommand(pssRequestUpdateCmd())

	// Child root commands
	cmd.AddCommand(pssRequestReplyRootCmd())

	return cmd
}

func pssRequestListCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists requests",
		Long:    "This command lists requests",
		Example: "ukfast pss request list",
		RunE: func(cmd *cobra.Command, args []string) error {
			return pssRequestList(getClient().PSSService(), cmd, args)
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

	outputPSSRequests(requests)
	return nil
}

func pssRequestShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "show <request: id>...",
		Short:   "Shows a request",
		Long:    "This command shows one or more requests",
		Example: "ukfast pss request show 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing request")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			pssRequestShow(getClient().PSSService(), cmd, args)
		},
	}
}

func pssRequestShow(service pss.PSSService, cmd *cobra.Command, args []string) {
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

	outputPSSRequests(requests)
}

func pssRequestCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Creates a request",
		Long:    "This command creates a new request",
		Example: "ukfast pss request create --subject 'example ticket' --details 'example' --author 123",
		RunE: func(cmd *cobra.Command, args []string) error {
			return pssRequestCreate(getClient().PSSService(), cmd, args)
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

	outputPSSRequests([]pss.Request{request})
	return nil
}

func pssRequestUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <request: id>...",
		Short:   "Updates requests",
		Long:    "This command updates one or more requests",
		Example: "ukfast pss request update 123 --priority high",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing request")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return pssRequestUpdate(getClient().PSSService(), cmd, args)
		},
	}

	// Setup flags
	cmd.Flags().String("priority", "", "Specifies priority for request")
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

	outputPSSRequests(requests)
	return nil
}
