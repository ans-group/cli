package helper

import (
	"errors"
	"io"
	"strconv"
	"strings"

	"github.com/ans-group/cli/internal/pkg/clierrors"
	"github.com/ans-group/sdk-go/pkg/config"
	"github.com/ans-group/sdk-go/pkg/connection"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

// InferTypeFromStringFlagValue will return a int, bool or string, based on value of flag
func InferTypeFromStringFlagValue(flag string) any {
	intValue, err := strconv.Atoi(flag)
	if err == nil {
		return intValue
	}

	boolValue, err := strconv.ParseBool(flag)
	if err == nil {
		return boolValue
	}

	return flag
}

// GetFilteringInferOperator returns an APIRequestFiltering struct with the operater inferred from the
// input value. This will return an operator based on the following rules:
// If value contains a comma ',' - returns IN operator
// If value contains an asterisk '*' - returns LK operator
// Otherwise returns EQ operator
func GetFilteringInferOperator(property string, value string) connection.APIRequestFiltering {
	return connection.APIRequestFiltering{
		Property: property,
		Operator: inferOperatorFromValue(value),
		Value:    []string{value},
	}
}

func inferOperatorFromValue(value string) connection.APIRequestFilteringOperator {
	if strings.Contains(value, ",") {
		return connection.INOperator
	}
	if strings.Contains(value, "*") {
		return connection.LKOperator
	}

	return connection.EQOperator
}

// GetFilteringArrayFromStringArrayFlagValue retrieves an array of APIRequestFiltering structs for given
// filtering strings
func GetFilteringArrayFromStringArrayFlagValue(filters []string) ([]connection.APIRequestFiltering, error) {
	var filtering []connection.APIRequestFiltering
	for _, filter := range filters {
		f, err := GetFilteringFromStringFlagValue(filter)
		if err != nil {
			return filtering, clierrors.NewErrInvalidFlagValue("filter", filter, err)
		}

		filtering = append(filtering, f)
	}

	return filtering, nil
}

// GetFilteringFromStringFlagValue retrieves a APIRequestFiltering struct from given filtering
// string. This function expects a string in the following format (with optional :operator): propertyname:operator=value,
// Valid examples:
// name:eq=something
// name=something
func GetFilteringFromStringFlagValue(filter string) (connection.APIRequestFiltering, error) {
	filtering := connection.APIRequestFiltering{}

	if filter == "" {
		return filtering, nil
	}

	// Obtain KV parts from filtering flag string. Example: propertyname:eq=value
	// K at index 0 represents propertyname and optional :operator in format propertyname:operator
	// V at index 1 represents the value
	filteringKVParts := strings.Split(filter, "=")
	if len(filteringKVParts) != 2 || filteringKVParts[1] == "" {
		return filtering, errors.New("missing value for filtering")
	}

	// Obtain PropertyOperator parts from K above. Example: propertyname:operator
	// Property at index 0 represents the property name
	// Operator at index 1 represents the operator
	filteringPropertyOperatorParts := strings.Split(filteringKVParts[0], ":")
	if filteringPropertyOperatorParts[0] == "" {
		return filtering, errors.New("missing property for filtering")
	}

	var operator connection.APIRequestFilteringOperator
	if len(filteringPropertyOperatorParts) == 1 {
		operator = inferOperatorFromValue(filteringKVParts[1])
	} else {
		if len(filteringPropertyOperatorParts) != 2 || filteringPropertyOperatorParts[1] == "" {
			return filtering, errors.New("missing operator for filtering")
		}

		// Parse the operator, returning parse error if any
		parsedOperator, err := connection.APIRequestFilteringOperatorEnum.Parse(filteringPropertyOperatorParts[1])
		if err != nil {
			return filtering, errors.New("invalid filtering operator")
		}

		operator = parsedOperator
	}

	// Sanitize comma-separated value by trimming spaces following split
	var sanitizedValues []string
	values := strings.SplitSeq(filteringKVParts[1], ",")
	for value := range values {
		sanitizedValues = append(sanitizedValues, strings.TrimSpace(value))
	}

	filtering.Property = filteringPropertyOperatorParts[0]
	filtering.Operator = operator
	filtering.Value = sanitizedValues

	return filtering, nil
}

// GetSortingFromStringFlagValue return an APIRequestSorting struct from given sorting string flag
func GetSortingFromStringFlagValue(sort string) connection.APIRequestSorting {
	if sort == "" {
		return connection.APIRequestSorting{}
	}

	var descending bool

	sortingParts := strings.Split(sort, ":")
	if (len(sortingParts)) > 1 && strings.ToLower(sortingParts[1]) == "desc" {
		descending = true
	}

	return connection.APIRequestSorting{
		Property:   sortingParts[0],
		Descending: descending,
	}
}

type APIRequestParametersFromFlagsOption interface {
	Hydrate(params *connection.APIRequestParameters, cmd *cobra.Command)
}

// GetAPIRequestParametersFromFlags returns an APIRequestParameters populated from global flags
func GetAPIRequestParametersFromFlags(cmd *cobra.Command, opts ...APIRequestParametersFromFlagsOption) (connection.APIRequestParameters, error) {
	flagFilter, _ := cmd.Flags().GetStringArray("filter")
	filtering, err := GetFilteringArrayFromStringArrayFlagValue(flagFilter)
	if err != nil {
		return connection.APIRequestParameters{}, err
	}

	flagSort, _ := cmd.Flags().GetString("sort")
	flagPage, _ := cmd.Flags().GetInt("page")

	params := connection.APIRequestParameters{
		Sorting:   GetSortingFromStringFlagValue(flagSort),
		Filtering: filtering,
		Pagination: connection.APIRequestPagination{
			PerPage: config.GetInt("api_pagination_perpage"),
			Page:    flagPage,
		},
	}

	for _, opt := range opts {
		opt.Hydrate(&params, cmd)
	}

	return params, nil
}

type StringFilterFlagOption struct {
	FlagName           string
	FilterPropertyName string
}

func NewStringFilterFlagOption(flagName string, filterPropertyName string) *StringFilterFlagOption {
	return &StringFilterFlagOption{
		FlagName:           flagName,
		FilterPropertyName: filterPropertyName,
	}
}

func (f *StringFilterFlagOption) Hydrate(params *connection.APIRequestParameters, cmd *cobra.Command) {
	if cmd.Flags().Changed(f.FlagName) {
		flagValue, _ := cmd.Flags().GetString(f.FlagName)
		params.WithFilter(GetFilteringInferOperator(f.FilterPropertyName, flagValue))
	}
}

func GetContentsFromFilePathFlag(cmd *cobra.Command, fs afero.Fs, filePathFlag string) (string, error) {
	filePath, _ := cmd.Flags().GetString(filePathFlag)
	file, err := fs.Open(filePath)
	if err != nil {
		return "", err
	}

	contentBytes, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	return string(contentBytes), nil
}

func GetContentsFromLiteralOrFilePathFlag(cmd *cobra.Command, fs afero.Fs, literalFlag, filePathFlag string) (string, error) {
	if cmd.Flags().Changed(literalFlag) {
		return cmd.Flags().GetString(literalFlag)
	}
	if cmd.Flags().Changed(filePathFlag) {
		return GetContentsFromFilePathFlag(cmd, fs, filePathFlag)
	}

	return "", nil
}

func GetStringPtrFlagIfChanged(cmd *cobra.Command, name string) *string {
	if cmd.Flags().Changed(name) {
		value, _ := cmd.Flags().GetString(name)
		return &value
	}
	return nil
}

func GetBoolPtrFlagIfChanged(cmd *cobra.Command, name string) *bool {
	if cmd.Flags().Changed(name) {
		value, _ := cmd.Flags().GetBool(name)
		return &value
	}
	return nil
}

func GetIntPtrFlagIfChanged(cmd *cobra.Command, name string) *int {
	if cmd.Flags().Changed(name) {
		value, _ := cmd.Flags().GetInt(name)
		return &value
	}
	return nil
}
