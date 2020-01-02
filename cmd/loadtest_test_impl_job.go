package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ukfast/sdk-go/pkg/connection"
	"github.com/ukfast/sdk-go/pkg/service/ltaas"
)

func loadtestTestJobRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "job",
		Short: "sub-commands relating to jobs",
	}

	// Child commands
	cmd.AddCommand(loadtestTestJobCreateCmd())

	return cmd
}

func loadtestTestJobCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Creates a job for a test",
		Long:    "This command creates a job for a test ",
		Example: "ukfast loadtest test job create",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing test")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return loadtestTestJobCreate(getClient().LTaaSService(), cmd, args)
		},
	}

	cmd.Flags().String("scheduled-timestamp", "", "Timestamp for schedule")
	cmd.Flags().Bool("run-now", false, "Indicates test should be started immediately")

	return cmd
}

func loadtestTestJobCreate(service ltaas.LTaaSService, cmd *cobra.Command, args []string) error {
	createRequest := ltaas.CreateTestJobRequest{}
	createRequest.RunNow, _ = cmd.Flags().GetBool("run-now")

	if cmd.Flags().Changed("scheduled-timestamp") {
		scheduledTimestamp, _ := cmd.Flags().GetString("scheduled-timestamp")
		createRequest.ScheduledTimestamp = connection.DateTime(scheduledTimestamp)
	}

	jobID, err := service.CreateTestJob(args[0], createRequest)
	if err != nil {
		return fmt.Errorf("Error creating job: %s", err)
	}

	job, err := service.GetJob(jobID)
	if err != nil {
		return fmt.Errorf("Error retrieving new job [%s]: %s", jobID, err)
	}

	return outputLoadTestJobs([]ltaas.Job{job})
}
