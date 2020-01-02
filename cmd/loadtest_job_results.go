package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/guptarohit/asciigraph"
	"github.com/spf13/cobra"
	"github.com/ukfast/sdk-go/pkg/service/ltaas"
)

func loadtestJobResultsRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "results",
		Short: "sub-commands relating to job results",
	}

	// Child commands
	cmd.AddCommand(loadtestJobResultsShowCmd())

	return cmd
}

func loadtestJobResultsShowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "show",
		Short:   "Shows job results",
		Long:    "This command shows job results",
		Example: "ukfast loadtest job results show",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing job")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return loadtestJobResultsShow(getClient().LTaaSService(), cmd, args)
		},
	}

	cmd.Flags().Int("graph-width", 100, "Specifies width of graphs")
	cmd.Flags().Int("graph-height", 10, "Specifies height of graphs")
	cmd.Flags().Bool("graph-virtualusers", false, "Specifies output should be a graph of virtual users")
	cmd.Flags().Bool("graph-successfulrequests", false, "Specifies output should be a graph of successful requests")
	cmd.Flags().Bool("graph-failedrequests", false, "Specifies output should be a graph of failed requests")
	cmd.Flags().Bool("graph-latency", false, "Specifies output should be a graph of latest")

	return cmd
}

func loadtestJobResultsShow(service ltaas.LTaaSService, cmd *cobra.Command, args []string) error {
	var allResults []ltaas.JobResults
	for _, arg := range args {
		results, err := service.GetJobResults(arg)
		if err != nil {
			return fmt.Errorf("Error retrieving job results: %s", err)
		}

		graphWidth, _ := cmd.Flags().GetInt("graph-width")
		graphHeight, _ := cmd.Flags().GetInt("graph-height")
		graphVirtualUsers, _ := cmd.Flags().GetBool("graph-virtualusers")
		graphSuccessfulRequests, _ := cmd.Flags().GetBool("graph-successfulrequests")
		graphFailedRequests, _ := cmd.Flags().GetBool("graph-failedrequests")
		graphLatency, _ := cmd.Flags().GetBool("graph-latency")

		if graphVirtualUsers || graphSuccessfulRequests || graphFailedRequests || graphLatency {
			var output []string

			if graphVirtualUsers {
				output = append(output, generateGraph(graphWidth, graphHeight, "# Virtual Users", results.VirtualUsers))
			}

			if graphSuccessfulRequests {
				output = append(output, generateGraph(graphWidth, graphHeight, "# Successful Requests", results.SuccessfulRequests))
			}

			if graphFailedRequests {
				output = append(output, generateGraph(graphWidth, graphHeight, "# Failed Requests", results.FailedRequests))
			}

			if graphLatency {
				output = append(output, generateGraph(graphWidth, graphHeight, "Latency (ms)", results.Latency))
			}

			fmt.Println(strings.Join(output, "\n\n"))
			continue
		}

		allResults = append(allResults, results)
	}

	return outputLoadTestJobResults(allResults)
}

// generateGraph returns an ASCII graph for given parameters
func generateGraph(graphWidth int, graphHeight int, graphCaption string, axisArray []ltaas.JobResultsAxis) string {
	values, valid := getGraphValues(axisArray)
	if !valid {
		return ""
	}

	return asciigraph.Plot(values, asciigraph.Caption(graphCaption), asciigraph.Width(graphWidth), asciigraph.Height(graphHeight))
}

// getGraphValues returns the Y-axis values for provided axisArray, and a boolean indicating
// whether values are valid (contain at least 1 positive value), to workaround bug https://github.com/guptarohit/asciigraph/issues/17
func getGraphValues(axisArray []ltaas.JobResultsAxis) ([]float64, bool) {
	values := []float64{}
	valid := false
	for _, axis := range axisArray {
		if axis.Y > 0 {
			valid = true
		}

		values = append(values, axis.Y)
	}

	return values, valid
}
