package cmd

import (
	"errors"
	"strconv"

	"github.com/spf13/cobra"
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
		Run: func(cmd *cobra.Command, args []string) {
			pssRequestList(getClient().PSSService(), cmd, args)
		},
	}
}

func pssRequestList(service pss.PSSService, cmd *cobra.Command, args []string) {
	params, err := GetAPIRequestParametersFromFlags()
	if err != nil {
		output.Fatal(err.Error())
		return
	}

	requests, err := service.GetRequests(params)
	if err != nil {
		output.Fatalf("Error retrieving requests: %s", err)
		return
	}

	outputPSSRequests(requests)
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
			OutputWithErrorLevelf("Invalid request ID [%s]", arg)
			continue
		}

		request, err := service.GetRequest(requestID)
		if err != nil {
			OutputWithErrorLevelf("Error retrieving request [%s]: %s", arg, err)
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
		Example: "ukfast pss request create --subject 'example ticket' --details 'example' --contact 123 --priority normal",
		Run: func(cmd *cobra.Command, args []string) {
			pssRequestCreate(getClient().PSSService(), cmd, args)
		},
	}

	// Setup flags
	cmd.Flags().String("subject", "", "Specifies subject for request")
	cmd.MarkFlagRequired("subject")
	cmd.Flags().String("details", "", "Specifies details for request")
	cmd.MarkFlagRequired("details")
	cmd.Flags().Int("contact", 0, "Contact ID for request")
	cmd.MarkFlagRequired("contact")
	cmd.Flags().String("priority", "", "Specifies priority for request")
	cmd.MarkFlagRequired("priority")
	cmd.Flags().Bool("secure", false, "Specifies whether request is secure")
	cmd.Flags().StringSlice("cc", []string{}, "Specifies CC email addresses for request")
	cmd.Flags().Bool("request-sms", false, "Specifies whether SMS updates are required")
	cmd.Flags().String("customer-reference", "", "Specifies customer reference for request")
	cmd.Flags().Int("product-id", 0, "Specifies product ID for request")
	cmd.Flags().String("product-name", "", "Specifies product name for request")
	cmd.Flags().String("product-type", "", "Specifies product type for request")

	return cmd
}

func pssRequestCreate(service pss.PSSService, cmd *cobra.Command, args []string) {
	createRequest := pss.CreateRequestRequest{}

	priority, _ := cmd.Flags().GetString("priority")
	parsedPriority, err := pss.ParseRequestPriority(priority)
	if err != nil {
		output.Fatal(err.Error())
		return
	}
	createRequest.Priority = parsedPriority

	if cmd.Flags().Changed("product-id") || cmd.Flags().Changed("product-name") || cmd.Flags().Changed("product-type") {
		createRequest.Product = &pss.Product{}
		createRequest.Product.ID, _ = cmd.Flags().GetInt("product-id")
		createRequest.Product.Name, _ = cmd.Flags().GetString("product-name")
		createRequest.Product.Type, _ = cmd.Flags().GetString("product-type")
	}

	createRequest.Subject, _ = cmd.Flags().GetString("subject")
	createRequest.Details, _ = cmd.Flags().GetString("details")
	createRequest.ContactID, _ = cmd.Flags().GetInt("contact")
	createRequest.Secure, _ = cmd.Flags().GetBool("secure")
	createRequest.CC, _ = cmd.Flags().GetStringSlice("cc")
	createRequest.RequestSMS, _ = cmd.Flags().GetBool("request-sms")
	createRequest.CustomerReference, _ = cmd.Flags().GetString("customer-reference")

	requestID, err := service.CreateRequest(createRequest)
	if err != nil {
		output.Fatalf("Error creating request: %s", err)
		return
	}

	request, err := service.GetRequest(requestID)
	if err != nil {
		output.Fatalf("Error retrieving new request: %s", err)
		return
	}

	outputPSSRequests([]pss.Request{request})
}
