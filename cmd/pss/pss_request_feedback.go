package pss

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/pss"
	"github.com/spf13/cobra"
)

func pssRequestFeedbackRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "feedback",
		Short: "sub-commands relating to request feedback",
	}

	// Child commands
	cmd.AddCommand(pssRequestFeedbackShowCmd(f))
	cmd.AddCommand(pssRequestFeedbackCreateCmd(f))

	return cmd
}

func pssRequestFeedbackShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <request: id>...",
		Short:   "Shows feedback for a request",
		Long:    "This command shows feedback for one or more requests",
		Example: "ans pss request feedback show 123",
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

			return pssRequestFeedbackShow(c.PSSService(), cmd, args)
		},
	}
}

func pssRequestFeedbackShow(service pss.PSSService, cmd *cobra.Command, args []string) error {
	var feedbacks []pss.Feedback
	for _, arg := range args {
		requestID, err := strconv.Atoi(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Invalid request ID [%s]", arg)
			continue
		}

		feedback, err := service.GetRequestFeedback(requestID)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving feedback for request [%s]: %s", arg, err)
			continue
		}

		feedbacks = append(feedbacks, feedback)
	}

	return output.CommandOutput(cmd, OutputPSSFeedbackProvider(feedbacks))
}

func pssRequestFeedbackCreateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create <request: id>",
		Short:   "Creates feedback for a request",
		Long:    "This command creates feedback for a request",
		Example: "ans pss request feedback create 123",
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

			return pssRequestFeedbackCreate(c.PSSService(), cmd, args)
		},
	}

	cmd.Flags().Int("contact", 0, "Specifies contact ID")
	cmd.MarkFlagRequired("contact")
	cmd.Flags().Int("score", 0, "Specifies feedback score")
	cmd.Flags().String("comment", "", "Specifies feedback comment")
	cmd.Flags().Int("speed-resolved", 0, "Specifies feedback speed resolved score")
	cmd.Flags().Int("quality", 0, "Specifies feedback quality")
	cmd.Flags().Int("nps-score", 0, "Specifies feedback NPS score")
	cmd.Flags().Bool("thirdparty-consent", false, "Specifies feedback third party consent")

	return cmd
}

func pssRequestFeedbackCreate(service pss.PSSService, cmd *cobra.Command, args []string) error {
	requestID, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("Invalid request ID [%s]", args[0])
	}

	createRequest := pss.CreateFeedbackRequest{}
	createRequest.ContactID, _ = cmd.Flags().GetInt("contact")
	createRequest.Score, _ = cmd.Flags().GetInt("score")
	createRequest.Comment, _ = cmd.Flags().GetString("comment")
	createRequest.SpeedResolved, _ = cmd.Flags().GetInt("speed-resolved")
	createRequest.Quality, _ = cmd.Flags().GetInt("quality")
	createRequest.NPSScore, _ = cmd.Flags().GetInt("nps-score")
	createRequest.ThirdPartyConsent, _ = cmd.Flags().GetBool("thirdparty-consent")

	_, err = service.CreateRequestFeedback(requestID, createRequest)
	if err != nil {
		return fmt.Errorf("Error creating feedback for request: %s", err)
	}

	feedback, err := service.GetRequestFeedback(requestID)
	if err != nil {
		return fmt.Errorf("Error retrieving new feedback for request: %s", err)
	}

	return output.CommandOutput(cmd, OutputPSSFeedbackProvider([]pss.Feedback{feedback}))
}
