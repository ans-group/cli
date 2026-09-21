package ecloud

import (
	"errors"
	"fmt"
	"reflect"
	"slices"
	"strings"
	"sync"

	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/connection"
	"github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/spf13/cobra"
)

// describeMaxConcurrency is the maximum number of related/child lookups performed
// in parallel for a single resource.
const describeMaxConcurrency = 8

// describeRedactedValue is displayed in place of a secret field value which has not
// been explicitly requested.
const describeRedactedValue = "****"

// describeSecretFields are the json field names redacted unless --with-passwords is
// specified.
var describeSecretFields = []string{"password"}

// describeOptions holds the resolved command flags for a describe invocation.
type describeOptions struct {
	// depth is the depth of relationships to resolve. Depth 0 fetches the resource
	// only, depth 1 (or greater) additionally resolves parents and children.
	depth int
	// children indicates whether child collections should be fetched.
	children bool
	// childLimit is the maximum number of items fetched per child collection.
	childLimit int
	// withPasswords indicates whether secret field values are displayed rather
	// than redacted.
	withPasswords bool
}

func ecloudDescribeCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "describe <id>...",
		Short: "Describes an eCloud resource and its relationships",
		Long: "This command describes one or more eCloud resources, resolving the parent resources " +
			"they reference and retrieving their child collections. The type of each resource is " +
			"determined from its ID prefix, so no resource type needs to be specified.\n\n" +
			"Passwords held by retrieved resources, such as instance credentials, are redacted " +
			"unless --with-passwords is specified.\n\n" +
			"Note that the SDK used by this command does not support request cancellation, so " +
			"in-flight requests cannot be interrupted. The 'api_timeout_seconds' configuration " +
			"value is the only timeout applied to lookups performed by this command.",
		Example: "ans ecloud describe i-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("missing id")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudDescribe),
	}

	cmd.Flags().Int("depth", 1, "Depth of relationships to resolve. 0 retrieves the resource only, 1 additionally resolves parent references and retrieves child collections. Values greater than 1 are accepted but currently behave as 1")
	cmd.Flags().Bool("no-children", false, "Specifies that child collections should not be retrieved")
	cmd.Flags().Int("child-limit", 50, "Maximum number of items to retrieve per child collection")
	cmd.Flags().Bool("with-passwords", false, "Specifies that passwords should be displayed rather than redacted")

	return cmd
}

// describeOptionsFromCommand returns the describe options given by a command's flags.
func describeOptionsFromCommand(cmd *cobra.Command) describeOptions {
	depth, _ := cmd.Flags().GetInt("depth")
	noChildren, _ := cmd.Flags().GetBool("no-children")
	childLimit, _ := cmd.Flags().GetInt("child-limit")
	withPasswords, _ := cmd.Flags().GetBool("with-passwords")

	return describeOptions{
		depth:         depth,
		children:      !noChildren,
		childLimit:    childLimit,
		withPasswords: withPasswords,
	}
}

func ecloudDescribe(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	opts := describeOptionsFromCommand(cmd)
	if opts.childLimit < 1 {
		return fmt.Errorf("ecloud: invalid child limit [%d], must be greater than 0", opts.childLimit)
	}

	// A single session is shared by every resource described, so that a resource
	// referenced by more than one of them is retrieved only once.
	session := newDescribeSession(service)

	results := make([]DescribeResult, 0, len(args))
	for _, arg := range args {
		result, err := ecloudDescribeResource(session, arg, opts)
		if err != nil {
			output.OutputWithErrorLevelf("%s", err)
			continue
		}

		results = append(results, result)
	}

	return ecloudDescribeOutput(cmd, results)
}

// describeSession holds the state shared by the resources described by a single
// command invocation.
type describeSession struct {
	service ecloud.ECloudService

	mutex    sync.Mutex
	resolved map[string]*describeResolution
}

// describeResolution is the cached outcome of resolving one parent reference.
type describeResolution struct {
	once         sync.Once
	resourceType string
	name         string
	err          string
}

func newDescribeSession(service ecloud.ECloudService) *describeSession {
	return &describeSession{
		service:  service,
		resolved: make(map[string]*describeResolution),
	}
}

// resolve returns the type name and display name of the resource with the given ID,
// along with any error encountered retrieving it. References to unknown resource
// types are left unresolved rather than treated as an error, and each ID is retrieved
// at most once for the lifetime of the session.
func (s *describeSession) resolve(id string) (resourceType, name, resolveError string) {
	s.mutex.Lock()
	resolution, ok := s.resolved[id]
	if !ok {
		resolution = &describeResolution{}
		s.resolved[id] = resolution
	}
	s.mutex.Unlock()

	// Concurrent callers for the same ID block until the first has finished, which
	// also establishes the happens-before needed to read the fields below.
	resolution.once.Do(func() {
		prefix, _, found := strings.Cut(id, "-")
		if !found {
			return
		}

		descriptor, ok := ecloudDescribeDescriptorsByPrefix[prefix]
		if !ok {
			return
		}

		resolution.resourceType = descriptor.Name

		resource, err := descriptor.Fetch(s.service, id)
		if err != nil {
			resolution.err = fmt.Sprintf("error retrieving %s [%s]: %s", strings.ToLower(descriptor.Name), id, err)
			return
		}

		resolution.name = describeResourceName(resource)
	})

	return resolution.resourceType, resolution.name, resolution.err
}

// ecloudDescribeResource describes a single resource, resolving its parent references
// and child collections as directed by opts. Failures resolving individual parents or
// child collections are recorded against the relevant entry rather than returned, so
// that a partially-visible resource can still be described.
func ecloudDescribeResource(session *describeSession, id string, opts describeOptions) (DescribeResult, error) {
	descriptor, err := ecloudDescribeDescriptorForID(id)
	if err != nil {
		return DescribeResult{}, err
	}

	resource, err := descriptor.Fetch(session.service, id)
	if err != nil {
		return DescribeResult{}, fmt.Errorf("error retrieving %s [%s]: %w", strings.ToLower(descriptor.Name), id, err)
	}

	if !opts.withPasswords {
		resource = describeRedactResource(resource)
	}

	result := DescribeResult{
		Type:     descriptor.Name,
		ID:       id,
		Resource: resource,
	}

	if opts.depth < 1 {
		return result, nil
	}

	related := describeRelatedReferences(resource)

	var children []DescribeChildren
	if opts.children {
		children = make([]DescribeChildren, len(descriptor.Children))
		for i, child := range descriptor.Children {
			children[i] = DescribeChildren{Title: child.Title}
		}
	}

	// Lookups are independent of one another, so are performed concurrently. Results
	// are written to their pre-allocated index to keep output ordering deterministic.
	var lookups []func()
	for i := range related {
		lookups = append(lookups, func() {
			related[i].Type, related[i].Name, related[i].Error = session.resolve(related[i].ID)
		})
	}
	for i := range children {
		child := descriptor.Children[i]
		lookups = append(lookups, func() {
			describeFetchChildren(session.service, child, id, opts, &children[i])
		})
	}

	describeRunConcurrently(lookups)

	result.Related = related
	result.Children = children

	return result, nil
}

// ecloudDescribeDescriptorForID returns the descriptor for the resource type identified
// by the prefix of the given resource ID.
func ecloudDescribeDescriptorForID(id string) (describeDescriptor, error) {
	prefix, _, found := strings.Cut(id, "-")
	if !found || prefix == "" {
		return describeDescriptor{}, fmt.Errorf("ecloud: invalid resource ID [%s], expected format <prefix>-<id>", id)
	}

	descriptor, ok := ecloudDescribeDescriptorsByPrefix[prefix]
	if !ok {
		return describeDescriptor{}, fmt.Errorf("ecloud: unknown resource type prefix [%s] (supported: %s)", prefix, strings.Join(ecloudDescribePrefixes(), ", "))
	}

	return descriptor, nil
}

// describeRelatedReferences returns the parent references held by a resource, being
// any non-empty string field with a json tag ending in '_id'. The resource's own 'id'
// field is excluded, as are duplicate references to the same field/ID pair.
func describeRelatedReferences(resource any) []DescribeRelated {
	reflectedValue := reflect.Indirect(reflect.ValueOf(resource))
	if reflectedValue.Kind() != reflect.Struct {
		return nil
	}

	var related []DescribeRelated
	seen := make(map[DescribeRelated]bool)

	reflectedType := reflectedValue.Type()
	for i := range reflectedType.NumField() {
		field := reflectedType.Field(i)
		value := reflectedValue.Field(i)

		if !value.CanInterface() || value.Kind() != reflect.String {
			continue
		}

		name := describeJSONFieldName(field)
		if name == "id" || !strings.HasSuffix(name, "_id") || value.String() == "" {
			continue
		}

		reference := DescribeRelated{Field: name, ID: value.String()}
		if seen[reference] {
			continue
		}

		seen[reference] = true
		related = append(related, reference)
	}

	return related
}

// describeFetchChildren retrieves the first page of a child collection, recording any
// failure against the collection rather than returning it.
func describeFetchChildren(service ecloud.ECloudService, child describeChild, id string, opts describeOptions, children *DescribeChildren) {
	params := connection.APIRequestParameters{
		Pagination: connection.APIRequestPagination{PerPage: opts.childLimit},
	}

	items, total, err := child.Fetch(service, id, params)
	if err != nil {
		children.Error = fmt.Sprintf("error retrieving %s: %s", strings.ToLower(child.Title), err)
		return
	}

	if !opts.withPasswords {
		describeRedactCollection(items)
	}

	children.Items = items
	children.Total = total

	reflectedItems := reflect.ValueOf(items)
	if reflectedItems.Kind() == reflect.Slice && total > reflectedItems.Len() {
		children.Truncated = true
	}
}

// describeRedactResource returns a copy of a resource with its secret fields redacted.
func describeRedactResource(resource any) any {
	reflectedValue := reflect.ValueOf(resource)
	if reflectedValue.Kind() != reflect.Struct {
		return resource
	}

	redacted := reflect.New(reflectedValue.Type()).Elem()
	redacted.Set(reflectedValue)
	describeRedactStruct(redacted)

	return redacted.Interface()
}

// describeRedactCollection redacts the secret fields held by the items of a child
// collection, in place.
func describeRedactCollection(items any) {
	reflectedItems := reflect.ValueOf(items)
	if reflectedItems.Kind() != reflect.Slice {
		return
	}

	for i := range reflectedItems.Len() {
		describeRedactStruct(reflectedItems.Index(i))
	}
}

// describeRedactStruct replaces the value of any secret field held by an addressable
// struct, so that credentials are not displayed unless explicitly requested.
func describeRedactStruct(value reflect.Value) {
	if value.Kind() != reflect.Struct {
		return
	}

	reflectedType := value.Type()
	for i := range reflectedType.NumField() {
		field := value.Field(i)
		if field.Kind() != reflect.String || !field.CanSet() || field.String() == "" {
			continue
		}

		if slices.Contains(describeSecretFields, describeJSONFieldName(reflectedType.Field(i))) {
			field.SetString(describeRedactedValue)
		}
	}
}

// describeResourceName returns the value of a resource's 'name' field, if it has one.
func describeResourceName(resource any) string {
	reflectedValue := reflect.Indirect(reflect.ValueOf(resource))
	if reflectedValue.Kind() != reflect.Struct {
		return ""
	}

	reflectedType := reflectedValue.Type()
	for i := range reflectedType.NumField() {
		value := reflectedValue.Field(i)
		if value.Kind() != reflect.String || !value.CanInterface() {
			continue
		}

		if describeJSONFieldName(reflectedType.Field(i)) == "name" {
			return value.String()
		}
	}

	return ""
}

// describeJSONFieldName returns the name of a struct field as given by its json tag,
// falling back to the field name where no tag is present.
func describeJSONFieldName(field reflect.StructField) string {
	tag, _, _ := strings.Cut(field.Tag.Get("json"), ",")
	if tag == "" {
		return field.Name
	}

	return tag
}

// describeRunConcurrently runs the given functions concurrently, with a maximum of
// describeMaxConcurrency running at any one time.
func describeRunConcurrently(fns []func()) {
	semaphore := make(chan struct{}, describeMaxConcurrency)

	var wg sync.WaitGroup
	for _, fn := range fns {
		wg.Add(1)
		go func() {
			defer wg.Done()

			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			fn()
		}()
	}

	wg.Wait()
}
