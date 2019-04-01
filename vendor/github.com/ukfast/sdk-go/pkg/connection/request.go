package connection

import (
	"errors"
	"strings"

	validator "gopkg.in/go-playground/validator.v9"
)

type Validatable interface {
	Validate() *ValidationError
}

type ValidationError struct {
	Message string
}

func NewValidationError(msg string) *ValidationError {
	return &ValidationError{Message: msg}
}

func (e *ValidationError) Error() string {
	return e.Message
}

type RequestBody interface {
	Validatable
}

// APIRequestBodyDefaultValidator provides convenience method Validate() for
// validating a struct/item, utilising the validator.v9 package
type APIRequestBodyDefaultValidator struct {
}

// Validate performs validation on item v using the validator.v9 package
func (a *APIRequestBodyDefaultValidator) Validate(v interface{}) *ValidationError {
	err := validator.New().Struct(v)
	if err != nil {
		return NewValidationError(err.Error())
	}

	return nil
}

type APIRequest struct {
	Method     string
	Resource   string
	Body       RequestBody
	Parameters APIRequestParameters
}

// APIRequestParameters holds a collection of supported API request parameters
type APIRequestParameters struct {
	Pagination APIRequestPagination
	Sorting    APIRequestSorting
	Filtering  []APIRequestFiltering
}

// WithPagination is a fluent method for adding pagination to request parameters
func (p *APIRequestParameters) WithPagination(pagination APIRequestPagination) *APIRequestParameters {
	p.Pagination = pagination
	return p
}

// WithSorting is a fluent method for adding sorting to request parameters
func (p *APIRequestParameters) WithSorting(sorting APIRequestSorting) *APIRequestParameters {
	p.Sorting = sorting
	return p
}

// WithFilter is a fluent method for adding a filter to request parameters
func (p *APIRequestParameters) WithFilter(filter APIRequestFiltering) *APIRequestParameters {
	p.Filtering = append(p.Filtering, filter)
	return p
}

// WithFilters is a fluent method for adding a set of filters to request parameters
func (p *APIRequestParameters) WithFilters(filters []APIRequestFiltering) *APIRequestParameters {
	for _, filter := range filters {
		p.Filtering = append(p.Filtering, filter)
	}

	return p
}

type APIRequestPagination struct {
	PerPage int
	Page    int
}

type APIRequestSorting struct {
	Property   string
	Descending bool
}

type APIRequestFiltering struct {
	Property string
	Operator APIRequestFilteringOperator
	Value    []string
}

type APIRequestFilteringOperator int

const (
	// EQOperator - equals
	EQOperator APIRequestFilteringOperator = iota
	// LKOperator - like
	LKOperator
	// GTOperator - greater than
	GTOperator
	// LTOperator - less than
	LTOperator
	// INOperator - in set
	INOperator
	// NEQOperator - not equal
	NEQOperator
	// NINOperator - not in set
	NINOperator
	// NLKOperator - not like
	NLKOperator
)

func (o APIRequestFilteringOperator) String() string {
	operators := [...]string{
		"eq",
		"lk",
		"gt",
		"lt",
		"in",
		"neq",
		"nin",
		"nlk",
	}

	if o < EQOperator || o > NLKOperator {
		return "Unknown"
	}

	return operators[o]
}

// ParseOperator attempts to parse an operator from string
func ParseOperator(o string) (APIRequestFilteringOperator, error) {
	switch strings.ToUpper(o) {
	case "EQ":
		return EQOperator, nil
	case "LK":
		return LKOperator, nil
	case "GT":
		return GTOperator, nil
	case "LT":
		return LTOperator, nil
	case "IN":
		return INOperator, nil
	case "NEQ":
		return NEQOperator, nil
	case "NIN":
		return NINOperator, nil
	case "NLK":
		return NLKOperator, nil
	}

	return 0, errors.New("Invalid filtering operator")
}

type RequestAll struct {
	GetNext func(parameters APIRequestParameters) (ResponseBody, error)
}

func (r *RequestAll) Invoke(parameters APIRequestParameters) error {
	totalPages := 1
	for currentPage := 1; currentPage <= totalPages; currentPage++ {
		parameters.Pagination.Page = currentPage
		response, err := r.GetNext(parameters)
		if err != nil {
			return err
		}

		totalPages = response.Pagination().TotalPages
	}

	return nil
}
