package ssl

import (
	"errors"

	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/ssl"
	"github.com/spf13/cobra"
)

func sslRecommendationsRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "recommendations",
		Short: "sub-commands relating to recommendations",
	}

	// Child commands
	cmd.AddCommand(sslRecommendationsShowCmd(f))

	return cmd
}

func sslRecommendationsShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <recommendations: id>...",
		Short:   "Shows SSL recommendations",
		Long:    "This command shows one or more SSL recommendations",
		Example: "ukfast ssl recommendations show example.com",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing domain")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return sslRecommendationsShow(c.SSLService(), cmd, args)
		},
	}
}

func sslRecommendationsShow(service ssl.SSLService, cmd *cobra.Command, args []string) error {
	var recommendationsSlice []ssl.Recommendations
	for _, arg := range args {
		recommendations, err := service.GetRecommendations(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving SSL recommendations for domain [%s]: %s", arg, err)
			continue
		}

		recommendationsSlice = append(recommendationsSlice, recommendations)
	}

	return output.CommandOutput(cmd, OutputSSLRecommendationsProvider(recommendationsSlice))
}
