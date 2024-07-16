package pss

import (
	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/helper"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/pss"
	"github.com/spf13/cobra"
)

func pssCaseCategoryRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "category",
		Short: "sub-commands relating to case categories",
	}

	// Child commands
	cmd.AddCommand(pssCaseCategoryListCmd(f))

	return cmd
}

func pssCaseCategoryListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists case categories",
		Long:    "This command lists case categories",
		Example: "ans pss case category list",
		RunE:    pssCobraRunEFunc(f, pssCaseCategoryList),
	}
}

func pssCaseCategoryList(service pss.PSSService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	categories, err := service.GetCaseCategories(params)
	if err != nil {
		return err
	}

	return output.CommandOutput(cmd, OutputPSSCaseCategoriesProvider(categories))
}
