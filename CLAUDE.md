# Development Guidelines

## Building

* Run `go build -o ans` from the root of the repository to build the binary.

## Libraries

* You should use the official ANS Go SDK for all communication with ANS endpoints. The SDK source code can be found at the following repository: https://github.com/ans-group/sdk-go

## API Documentation

* The eCloud API OpenAPI specification can be found here: https://developers.ukfast.io/api/documentation/ecloud/2. 
* API documentation is provided as a reference - you should use the Go SDK instead of communicating with the APIs directly.

## Code Style

* Use standard Go formatting conventions
* Format with `goimports` once you have finished making changes
* Group imports: standard library first, then third-party
* Use PascalCase for exported types/methods, camelCase for variables
* Add comments for public API and complex logic
* Place related functionality in logically named files
* You should follow existing patterns in the code base

## Error Handling

* Returned errors should be prefix with an appropriate identifier related to the package name. 
* This must remain consistent within that package. e.g., for a package called `foo`, prefix errors with `foo:`, e.g. `foo: this has failed`.
* Use consistent prefixes if another prefix is already being used within that package.

## Testing

* Write table-driven tests with clear input/output extensions. Add comments where necessary to clarify complex test logic.
  * However, where there are existing tests for a package, follow the existing conventions.
* Mocks are created with `mockgen`, however you MUST use the script in `test/mocks/generate.sh` to regenerate the mocks rather than using the `mockgen` command directly.

## Modernisation Notes

* Use errors.Is() and errors.As() for error checking
* Replace interface{} with any type alias
* Replace type assertions with type switches where appropriate
* Use generics for type-safe operations
* Implement context cancellation handling for long operations
* Add proper docstring comments for exported functions and types
* Avoid using deprecated functionality, e.g. do not use the `ioutil` package, instead prefer the implementations from `io` or `os`
