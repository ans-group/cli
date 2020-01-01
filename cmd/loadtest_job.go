package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/helper"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/connection"
	"github.com/ukfast/sdk-go/pkg/service/ltaas"
)

func loadtestJobRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "job",
		Short: "sub-commands relating to jobs",
	}

	// Child commands
	cmd.AddCommand(loadtestJobListCmd())
	cmd.AddCommand(loadtestJobShowCmd())
	cmd.AddCommand(loadtestJobCreateCmd())

	// Child root commands
	cmd.AddCommand(loadtestJobResultRootCmd())

	return cmd
}

func loadtestJobListCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists jobs",
		Long:    "This command lists jobs",
		Example: "ukfast loadtest job list",
		Run: func(cmd *cobra.Command, args []string) {
			loadtestJobList(getClient().LTaaSService(), cmd, args)
		},
	}
}

func loadtestJobList(service ltaas.LTaaSService, cmd *cobra.Command, args []string) {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		output.Fatal(err.Error())
		return
	}

	jobs, err := service.GetJobs(params)
	if err != nil {
		output.Fatalf("Error retrieving jobs: %s", err)
		return
	}

	outputLoadTestJobs(jobs)
}

func loadtestJobShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "show <job: id>...",
		Short:   "Shows a job",
		Long:    "This command shows one or more jobs",
		Example: "ukfast loadtest job show 00000000-0000-0000-0000-000000000000",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing job")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			loadtestJobShow(getClient().LTaaSService(), cmd, args)
		},
	}
}

func loadtestJobShow(service ltaas.LTaaSService, cmd *cobra.Command, args []string) {
	var jobs []ltaas.Job
	for _, arg := range args {
		job, err := service.GetJob(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving job [%s]: %s", arg, err)
			continue
		}

		jobs = append(jobs, job)
	}

	outputLoadTestJobs(jobs)
}

func loadtestJobCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Creates a job",
		Long:    "This command creates a job ",
		Example: "ukfast loadtest job create",
		RunE: func(cmd *cobra.Command, args []string) error {
			return loadtestJobCreate(getClient().LTaaSService(), cmd, args)
		},
	}

	cmd.Flags().String("test-id", "", "ID for test")
	cmd.MarkFlagRequired("test-id")
	cmd.Flags().String("scheduled-timestamp", "", "Timestamp for schedule")
	cmd.Flags().Bool("run-now", false, "Indicates job should be started immediately")

	return cmd
}

func loadtestJobCreate(service ltaas.LTaaSService, cmd *cobra.Command, args []string) error {
	createRequest := ltaas.CreateJobRequest{}
	createRequest.TestID, _ = cmd.Flags().GetString("test-id")
	createRequest.RunNow, _ = cmd.Flags().GetBool("run-now")

	if cmd.Flags().Changed("scheduled-timestamp") {
		scheduledTimestamp, _ := cmd.Flags().GetString("scheduled-timestamp")
		createRequest.ScheduledTimestamp = connection.DateTime(scheduledTimestamp)
	}

	jobID, err := service.CreateJob(createRequest)
	if err != nil {
		return fmt.Errorf("Error creating job: %s", err)
	}

	job, err := service.GetJob(jobID)
	if err != nil {
		return fmt.Errorf("Error retrieving new job: %s", err)
	}

	outputLoadTestJobs([]ltaas.Job{job})
	return nil
}

func loadtestJobDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "delete <job: id>...",
		Short:   "Deletes a job",
		Long:    "This command deletes one or more jobs",
		Example: "ukfast loadtest job delete 00000000-0000-0000-0000-000000000000",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing job")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			loadtestJobDelete(getClient().LTaaSService(), cmd, args)
		},
	}
}

func loadtestJobDelete(service ltaas.LTaaSService, cmd *cobra.Command, args []string) {
	var jobs []ltaas.Job
	for _, arg := range args {
		job, err := service.GetJob(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error removing job [%s]: %s", arg, err)
			continue
		}

		jobs = append(jobs, job)
	}

	outputLoadTestJobs(jobs)
}
