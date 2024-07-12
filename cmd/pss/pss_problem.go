package pss

import (
	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/helper"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/pss"
	"github.com/spf13/cobra"
)

func pssProblemRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "problem",
		Short: "sub-commands relating to problem cases",
	}

	// Child commands
	cmd.AddCommand(pssProblemListCmd(f))
	// cmd.AddCommand(pssProblemShowCmd(f))
	// cmd.AddCommand(pssProblemCreateCmd(f))
	// cmd.AddCommand(pssProblemUpdateCmd(f))
	// cmd.AddCommand(pssProblemCloseCmd(f))

	return cmd
}

func pssProblemListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists problem cases",
		Long:    "This command lists problem cases",
		Example: "ans pss problem list",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return pssProblemList(c.PSSService(), cmd, args)
		},
	}
}

func pssProblemList(service pss.PSSService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	problems, err := service.GetProblemCases(params)
	if err != nil {
		return err
	}

	return output.CommandOutput(cmd, OutputPSSProblemCasesProvider(problems))
}

// func pssProblemShowCmd(f factory.ClientFactory) *cobra.Command {
// 	return &cobra.Command{
// 		Use:     "show <problem: id>...",
// 		Short:   "Shows a problem",
// 		Long:    "This command shows one or more problems",
// 		Example: "ans pss problem show 123",
// 		Args: func(cmd *cobra.Command, args []string) error {
// 			if len(args) < 1 {
// 				return errors.New("Missing problem")
// 			}

// 			return nil
// 		},
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			c, err := f.NewClient()
// 			if err != nil {
// 				return err
// 			}

// 			return pssProblemShow(c.PSSService(), cmd, args)
// 		},
// 	}
// }

// func pssProblemShow(service pss.PSSService, cmd *cobra.Command, args []string) error {
// 	var problems []pss.Problem
// 	for _, arg := range args {
// 		problemID, err := strconv.Atoi(arg)
// 		if err != nil {
// 			output.OutputWithErrorLevelf("Invalid problem ID [%s]", arg)
// 			continue
// 		}

// 		problem, err := service.GetProblem(problemID)
// 		if err != nil {
// 			output.OutputWithErrorLevelf("Error retrieving problem [%s]: %s", arg, err)
// 			continue
// 		}

// 		problems = append(problems, problem)
// 	}

// 	return output.CommandOutput(cmd, OutputPSSProblemsProvider(problems))
// }

// func pssProblemCreateCmd(f factory.ClientFactory) *cobra.Command {
// 	cmd := &cobra.Command{
// 		Use:     "create",
// 		Short:   "Creates a problem",
// 		Long:    "This command creates a new problem",
// 		Example: "ans pss problem create --subject 'example ticket' --details 'example' --author 123",
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			c, err := f.NewClient()
// 			if err != nil {
// 				return err
// 			}

// 			return pssProblemCreate(c.PSSService(), cmd, args)
// 		},
// 	}

// 	// Setup flags
// 	cmd.Flags().String("subject", "", "Specifies subject for problem")
// 	cmd.MarkFlagRequired("subject")
// 	cmd.Flags().String("details", "", "Specifies details for problem")
// 	cmd.Flags().Int("author", 0, "Specifies author ID for problem")
// 	cmd.MarkFlagRequired("author")
// 	cmd.Flags().String("priority", "Normal", "Specifies priority for problem")
// 	cmd.Flags().Bool("secure", false, "Specifies whether problem is secure")
// 	cmd.Flags().StringSlice("cc", []string{}, "Specifies CC email addresses for problem")
// 	cmd.Flags().Bool("problem-sms", false, "Specifies whether SMS updates are required")
// 	cmd.Flags().String("customer-reference", "", "Specifies customer reference for problem")
// 	cmd.Flags().Int("product-id", 0, "Specifies product ID for problem")
// 	cmd.Flags().String("product-name", "", "Specifies product name for problem")
// 	cmd.Flags().String("product-type", "", "Specifies product type for problem")

// 	return cmd
// }

// func pssProblemCreate(service pss.PSSService, cmd *cobra.Command, args []string) error {
// 	createProblem := pss.CreateProblemProblem{}

// 	priority, _ := cmd.Flags().GetString("priority")
// 	parsedPriority, err := pss.ProblemPriorityEnum.Parse(priority)
// 	if err != nil {
// 		return err
// 	}
// 	createProblem.Priority = parsedPriority

// 	if cmd.Flags().Changed("product-id") || cmd.Flags().Changed("product-name") || cmd.Flags().Changed("product-type") {
// 		createProblem.Product = &pss.Product{}
// 		createProblem.Product.ID, _ = cmd.Flags().GetInt("product-id")
// 		createProblem.Product.Name, _ = cmd.Flags().GetString("product-name")
// 		createProblem.Product.Type, _ = cmd.Flags().GetString("product-type")
// 	}

// 	createProblem.Subject, _ = cmd.Flags().GetString("subject")
// 	createProblem.Author.ID, _ = cmd.Flags().GetInt("author")
// 	createProblem.Secure, _ = cmd.Flags().GetBool("secure")
// 	createProblem.CC, _ = cmd.Flags().GetStringSlice("cc")
// 	createProblem.ProblemSMS, _ = cmd.Flags().GetBool("problem-sms")
// 	createProblem.CustomerReference, _ = cmd.Flags().GetString("customer-reference")

// 	if cmd.Flags().Changed("details") {
// 		createProblem.Details, _ = cmd.Flags().GetString("details")
// 	} else {
// 		createProblem.Details, err = input.ReadInput("details")
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	problemID, err := service.CreateProblem(createProblem)
// 	if err != nil {
// 		return fmt.Errorf("Error creating problem: %s", err)
// 	}

// 	problem, err := service.GetProblem(problemID)
// 	if err != nil {
// 		return fmt.Errorf("Error retrieving new problem: %s", err)
// 	}

// 	return output.CommandOutput(cmd, OutputPSSProblemsProvider([]pss.Problem{problem}))
// }

// func pssProblemUpdateCmd(f factory.ClientFactory) *cobra.Command {
// 	cmd := &cobra.Command{
// 		Use:     "update <problem: id>...",
// 		Short:   "Updates problems",
// 		Long:    "This command updates one or more problems",
// 		Example: "ans pss problem update 123 --priority high",
// 		Args: func(cmd *cobra.Command, args []string) error {
// 			if len(args) < 1 {
// 				return errors.New("Missing problem")
// 			}

// 			return nil
// 		},
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			c, err := f.NewClient()
// 			if err != nil {
// 				return err
// 			}

// 			return pssProblemUpdate(c.PSSService(), cmd, args)
// 		},
// 	}

// 	// Setup flags
// 	cmd.Flags().String("priority", "", "Specifies priority for problem")
// 	cmd.Flags().String("status", "", "Specifies status for problem")
// 	cmd.Flags().Bool("secure", false, "Specifies whether problem is secure")
// 	cmd.Flags().Bool("read", false, "Specifies whether problem is marked as read")
// 	cmd.Flags().Bool("problem-sms", false, "Specifies whether SMS updates are required")
// 	cmd.Flags().Bool("archived", false, "Specifies whether problem is archived")

// 	return cmd
// }

// func pssProblemUpdate(service pss.PSSService, cmd *cobra.Command, args []string) error {
// 	patchProblem := pss.PatchProblemProblem{}

// 	if cmd.Flags().Changed("priority") {
// 		priority, _ := cmd.Flags().GetString("priority")
// 		parsedPriority, err := pss.ProblemPriorityEnum.Parse(priority)
// 		if err != nil {
// 			return err
// 		}
// 		patchProblem.Priority = parsedPriority
// 	}

// 	if cmd.Flags().Changed("status") {
// 		status, _ := cmd.Flags().GetString("status")
// 		parsedStatus, err := pss.ProblemStatusEnum.Parse(status)
// 		if err != nil {
// 			return err
// 		}
// 		patchProblem.Status = parsedStatus
// 	}

// 	if cmd.Flags().Changed("secure") {
// 		secure, _ := cmd.Flags().GetBool("secure")
// 		patchProblem.Secure = &secure
// 	}
// 	if cmd.Flags().Changed("read") {
// 		read, _ := cmd.Flags().GetBool("read")
// 		patchProblem.Read = &read
// 	}
// 	if cmd.Flags().Changed("problem-sms") {
// 		problemSMS, _ := cmd.Flags().GetBool("problem-sms")
// 		patchProblem.ProblemSMS = &problemSMS
// 	}
// 	if cmd.Flags().Changed("archived") {
// 		archived, _ := cmd.Flags().GetBool("archived")
// 		patchProblem.Archived = &archived
// 	}

// 	var problems []pss.Problem

// 	for _, arg := range args {
// 		problemID, err := strconv.Atoi(arg)
// 		if err != nil {
// 			output.OutputWithErrorLevelf("Invalid problem ID [%s]", arg)
// 			continue
// 		}

// 		err = service.PatchProblem(problemID, patchProblem)
// 		if err != nil {
// 			output.OutputWithErrorLevelf("Error updating problem [%d]: %s", problemID, err)
// 			continue
// 		}

// 		problem, err := service.GetProblem(problemID)
// 		if err != nil {
// 			output.OutputWithErrorLevelf("Error retrieving updated problem [%d]: %s", problemID, err)
// 			continue
// 		}

// 		problems = append(problems, problem)
// 	}

// 	return output.CommandOutput(cmd, OutputPSSProblemsProvider(problems))
// }

// func pssProblemCloseCmd(f factory.ClientFactory) *cobra.Command {
// 	return &cobra.Command{
// 		Use:     "close <problem: id>...",
// 		Short:   "Closes problems",
// 		Long:    "This command closes one or more problems",
// 		Example: "ans pss problem close 123",
// 		Args: func(cmd *cobra.Command, args []string) error {
// 			if len(args) < 1 {
// 				return errors.New("Missing problem")
// 			}

// 			return nil
// 		},
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			c, err := f.NewClient()
// 			if err != nil {
// 				return err
// 			}

// 			return pssProblemClose(c.PSSService(), cmd, args)
// 		},
// 	}
// }

// func pssProblemClose(service pss.PSSService, cmd *cobra.Command, args []string) error {
// 	patchProblem := pss.PatchProblemProblem{
// 		Status: pss.ProblemStatusCompleted,
// 	}

// 	var problems []pss.Problem

// 	for _, arg := range args {
// 		problemID, err := strconv.Atoi(arg)
// 		if err != nil {
// 			output.OutputWithErrorLevelf("Invalid problem ID [%s]", arg)
// 			continue
// 		}

// 		err = service.PatchProblem(problemID, patchProblem)
// 		if err != nil {
// 			output.OutputWithErrorLevelf("Error closing problem [%d]: %s", problemID, err)
// 			continue
// 		}

// 		problem, err := service.GetProblem(problemID)
// 		if err != nil {
// 			output.OutputWithErrorLevelf("Error retrieving updated problem [%d]: %s", problemID, err)
// 			continue
// 		}

// 		problems = append(problems, problem)
// 	}

// 	return output.CommandOutput(cmd, OutputPSSProblemsProvider(problems))
// }
