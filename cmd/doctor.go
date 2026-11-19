package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/client"
	"github.com/ans-group/sdk-go/pkg/config"
	"github.com/ans-group/sdk-go/pkg/connection"
	"github.com/ans-group/sdk-go/pkg/service/account"
	"github.com/blang/semver"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Statuses reported by individual diagnostic checks
const (
	doctorStatusOK   = "ok"
	doctorStatusWarn = "warn"
	doctorStatusFail = "fail"
	doctorStatusSkip = "skip"
)

// doctorDetailMaxLength is the length at which unrecognised error strings are truncated,
// keeping the output table readable
const doctorDetailMaxLength = 80

// doctorDefaultConfigFileName is the config file viper falls back to reporting when no
// config file has been read, mirroring the default applied by config.Init() (sdk-go)
const doctorDefaultConfigFileName = ".ans.yml"

// doctorUnexpectedStatusCodeRegex matches the error returned by
// connection.APIResponse.HandleResponse when the response body carries no structured error
var doctorUnexpectedStatusCodeRegex = regexp.MustCompile(`unexpected status code \((\d+)\)`)

// doctorIPNotAllowedRegex matches the message returned by the API when a request is
// rejected because the caller's IP address is not on the key's allow list. The address
// itself is captured so it can be surfaced to the user
var doctorIPNotAllowedRegex = regexp.MustCompile(`(?i)IP address not allowed: ([0-9a-fA-F:.]+)`)

func doctorCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "doctor",
		Short: "Diagnoses CLI configuration and API connectivity",
		Long: `This command performs a read-only set of diagnostic checks against the CLI
configuration and the ANS API, reporting whether the CLI is correctly configured and which
services the configured API key can reach.

Service reachability is established empirically, by issuing a cheap read against each
service. The API provides no endpoint for describing the key currently in use, so the
declared scopes and IP restrictions of an API application can only be reported when the
application ID is provided with --application.

The API key is never printed by this command.`,
		Example: "ans doctor",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return doctor(f, cmd)
		},
	}

	cmd.Flags().Bool("all-contexts", false, "Specifies that checks should be run for all configured contexts")
	cmd.Flags().String("application", "", "Specifies an API application ID for which declared scopes and restrictions should be reported")
	cmd.Flags().Bool("skip-version-check", false, "Specifies that the CLI version check should be skipped")
	cmd.Flags().Bool("skip-services", false, "Specifies that per-service probes should be skipped")

	return cmd
}

func doctor(f factory.ClientFactory, cmd *cobra.Command) error {
	allContexts, _ := cmd.Flags().GetBool("all-contexts")
	skipVersionCheck, _ := cmd.Flags().GetBool("skip-version-check")

	results := doctorConfigChecks(doctorConfigFilePath())

	if allContexts {
		results = append(results, doctorAllContextChecks(f, cmd)...)
	} else {
		results = append(results, doctorContextChecks(f, cmd)...)
	}

	if skipVersionCheck {
		results = append(results, DoctorResult{
			Check:  "cli version",
			Status: doctorStatusSkip,
			Detail: "skipped (--skip-version-check)",
		})
	} else {
		results = append(results, doctorVersionCheck(appVersion))
	}

	err := output.CommandOutput(cmd, DoctorResultCollection(results))
	if err != nil {
		return err
	}

	// A failed check results in a non-zero exit code, so that the command is usable in CI
	// and setup scripts. The failures are already described in the rendered output, so no
	// further output is written. Warnings do not affect the exit code
	if slices.ContainsFunc(results, func(result DoctorResult) bool {
		return result.Status == doctorStatusFail
	}) {
		output.SetErrorLevel(1)
	}

	return nil
}

// doctorAllContextChecks runs the per-context checks against every configured context.
//
// The SDK config package wraps package-global viper state, so contexts must be switched -
// and therefore processed - serially. The original context is restored before returning
func doctorAllContextChecks(f factory.ClientFactory, cmd *cobra.Command) []DoctorResult {
	restoreContext := doctorCaptureContext()
	defer restoreContext()

	contextNames := config.GetContextNames()
	slices.Sort(contextNames)

	if len(contextNames) < 1 {
		return []DoctorResult{
			{
				Check:  "contexts",
				Status: doctorStatusWarn,
				Detail: "no contexts configured",
			},
		}
	}

	var results []DoctorResult
	for _, contextName := range contextNames {
		err := config.SwitchCurrentContext(contextName)
		if err != nil {
			results = append(results, DoctorResult{
				Context: contextName,
				Check:   "context",
				Status:  doctorStatusFail,
				Detail:  doctorTruncate(err.Error()),
			})
			continue
		}

		contextResults := doctorContextChecks(f, cmd)
		for i := range contextResults {
			contextResults[i].Context = contextName
		}

		results = append(results, contextResults...)
	}

	return results
}

// doctorContextChecks runs the checks which depend on the current context - connectivity,
// authentication, per-service probes and, when requested, application scopes
func doctorContextChecks(f factory.ClientFactory, cmd *cobra.Command) []DoctorResult {
	skipServices, _ := cmd.Flags().GetBool("skip-services")
	applicationID, _ := cmd.Flags().GetString("application")

	// The client is created once per context and shared, as its per-service accessors are
	// cheap constructors over a single connection
	c, err := f.NewClient()
	if err != nil {
		// The client is built from the very config these checks describe, so a failure
		// here is reported as a result rather than returned as an error. Returning would
		// discard the rendered report, which already explains the cause - a missing api
		// key being the common case
		return []DoctorResult{
			{
				Check:  "api client",
				Status: doctorStatusFail,
				Detail: doctorTruncate(fmt.Sprintf("failed to create client: %s", err)),
			},
		}
	}

	results := doctorAuthChecks(c.AccountService())

	if skipServices {
		results = append(results, DoctorResult{
			Check:  "services",
			Status: doctorStatusSkip,
			Detail: "skipped (--skip-services)",
		})
	} else {
		results = append(results, doctorServiceChecks(c)...)
	}

	if len(applicationID) > 0 {
		results = append(results, doctorApplicationChecks(c.AccountService(), applicationID)...)
	}

	return results
}

// doctorConfigFilePath resolves the config file path that was actually read by viper
// during initConfig() (cobra.OnInitialize runs before any RunE, including doctor's).
// This is the file config.Init() (sdk-go) resolved - accounting for --config, the '.ans'
// name and every extension viper supports - so no guessing is needed here.
//
// viper.ConfigFileUsed() is empty when no config file was found, mirroring the fallback
// applied by config.Save() (sdk-go), so the default path is returned in that case, purely
// so the check can report it as absent
func doctorConfigFilePath() string {
	if configFile := viper.ConfigFileUsed(); len(configFile) > 0 {
		return configFile
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	return filepath.Join(home, doctorDefaultConfigFileName)
}

// doctorConfigChecks performs the checks which require no network access
func doctorConfigChecks(configPath string) []DoctorResult {
	return []DoctorResult{
		doctorConfigFileCheck(configPath),
		doctorCurrentContextCheck(),
		doctorAPIKeyCheck(),
		doctorOutputDefaultCheck(),
		doctorContextsCheck(),
	}
}

// doctorConfigFileCheck checks that the config file exists and is readable. An absent
// config file is not a failure, as the CLI can be configured entirely by environment
// variables
func doctorConfigFileCheck(configPath string) DoctorResult {
	result := DoctorResult{Check: "config file"}

	if len(configPath) < 1 {
		result.Status = doctorStatusWarn
		result.Detail = "unable to determine config file path"
		return result
	}

	_, err := os.Stat(configPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			result.Status = doctorStatusWarn
			result.Detail = fmt.Sprintf("%s not found (environment variables may be in use)", configPath)
			return result
		}

		result.Status = doctorStatusFail
		result.Detail = doctorTruncate(err.Error())
		return result
	}

	f, err := os.Open(configPath)
	if err != nil {
		result.Status = doctorStatusFail
		result.Detail = doctorTruncate(err.Error())
		return result
	}
	_ = f.Close()

	result.Status = doctorStatusOK
	result.Detail = configPath
	return result
}

// doctorCurrentContextCheck checks that a current context is set. An unset context is not
// a failure, as the SDK falls back to top-level config keys
func doctorCurrentContextCheck() DoctorResult {
	result := DoctorResult{Check: "current context"}

	contextName := config.GetCurrentContextName()
	if len(contextName) < 1 {
		result.Status = doctorStatusWarn
		result.Detail = "no context set (using top-level config)"
		return result
	}

	result.Status = doctorStatusOK
	result.Detail = contextName
	return result
}

// doctorAPIKeyCheck checks that an API key is configured. The key itself is never
// reported - only its length - as the output of this command is likely to be shared
func doctorAPIKeyCheck() DoctorResult {
	result := DoctorResult{Check: "api key"}

	apiKey := config.GetString("api_key")
	if len(apiKey) < 1 {
		result.Status = doctorStatusFail
		result.Detail = "not set"
		return result
	}

	result.Status = doctorStatusOK
	result.Detail = fmt.Sprintf("set (%d chars)", len(apiKey))
	return result
}

// doctorOutputDefaultCheck checks that any configured default output format is one
// understood by the output handler
func doctorOutputDefaultCheck() DoctorResult {
	result := DoctorResult{Check: "config: output.default"}

	outputDefault := config.GetString("output.default")
	if len(outputDefault) < 1 {
		result.Status = doctorStatusOK
		result.Detail = "not set (defaults to table)"
		return result
	}

	// Formats may carry an argument, e.g. 'jsonpath={.id}'
	format, _ := output.ParseOutputFlag(outputDefault)
	if !output.IsValidFormat(format) {
		result.Status = doctorStatusWarn
		result.Detail = fmt.Sprintf("invalid format '%s' (expected one of: %s)", format, strings.Join(output.ValidFormats, ", "))
		return result
	}

	result.Status = doctorStatusOK
	result.Detail = outputDefault
	return result
}

// doctorContextsCheck checks that every configured context has an API key set
func doctorContextsCheck() DoctorResult {
	result := DoctorResult{Check: "config: contexts"}

	contextNames := config.GetContextNames()
	if len(contextNames) < 1 {
		result.Status = doctorStatusOK
		result.Detail = "no contexts configured"
		return result
	}

	slices.Sort(contextNames)

	var missing []string
	restoreContext := doctorCaptureContext()
	for _, contextName := range contextNames {
		err := config.SwitchCurrentContext(contextName)
		if err != nil {
			continue
		}

		if len(config.GetString("api_key")) < 1 {
			missing = append(missing, contextName)
		}
	}
	restoreContext()

	if len(missing) > 0 {
		result.Status = doctorStatusWarn
		result.Detail = fmt.Sprintf("%d context(s) without an api key: %s", len(missing), strings.Join(missing, ", "))
		return result
	}

	result.Status = doctorStatusOK
	result.Detail = fmt.Sprintf("%d context(s) configured, all with an api key", len(contextNames))
	return result
}

// doctorCaptureContext captures the name of the current context, returning a function
// which restores it. The config package is global state, so any code path which switches
// context must restore it before returning
func doctorCaptureContext() func() {
	originalContextName := config.GetCurrentContextName()

	return func() {
		if len(originalContextName) < 1 {
			// SwitchCurrentContext rejects an empty name, so the key is reset directly.
			// Config is never saved by this command, so this only affects the process
			config.Set("", "current_context", "")
			return
		}

		err := config.SwitchCurrentContext(originalContextName)
		if err != nil {
			output.Errorf("doctor: failed to restore context '%s': %s", originalContextName, err)
		}
	}
}

// doctorAuthChecks performs a single call to the account service, deriving both the
// connectivity and authentication results from it
func doctorAuthChecks(service account.AccountService) []DoctorResult {
	connectivity := DoctorResult{Check: "api connectivity"}
	authentication := DoctorResult{Check: "api authentication"}

	start := time.Now()
	_, err := service.GetDetails()
	elapsed := time.Since(start)

	if err == nil {
		connectivity.Status = doctorStatusOK
		connectivity.Detail = fmt.Sprintf("%dms", elapsed.Milliseconds())
		authentication.Status = doctorStatusOK
		authentication.Detail = "authenticated"

		return []DoctorResult{connectivity, authentication}
	}

	statusCode := httpStatusFromError(err)
	if statusCode < 1 {
		// No HTTP status could be determined, so the API was likely never reached
		connectivity.Status = doctorStatusFail
		connectivity.Detail = doctorTruncate(err.Error())
		authentication.Status = doctorStatusWarn
		authentication.Detail = "undetermined (connectivity failed)"

		return []DoctorResult{connectivity, authentication}
	}

	// A status code implies the API responded, so connectivity itself is fine
	connectivity.Status = doctorStatusOK
	connectivity.Detail = fmt.Sprintf("%dms", elapsed.Milliseconds())

	switch statusCode {
	case 401:
		authentication.Status = doctorStatusFail
		authentication.Detail = "401 unauthorised (api key rejected)"
	case 403:
		authentication.Status = doctorStatusFail
		if blocked, ipAddress := doctorIPNotAllowed(err); blocked {
			authentication.Detail = fmt.Sprintf("403 forbidden (ip address %s not allowed for this api key)", ipAddress)
		} else {
			authentication.Detail = "403 forbidden (api key lacks permission)"
		}
	default:
		authentication.Status = doctorStatusWarn
		authentication.Detail = doctorTruncate(err.Error())
	}

	return []DoctorResult{connectivity, authentication}
}

// doctorServiceProbe pairs a service name with a cheap read against that service
type doctorServiceProbe struct {
	name  string
	probe func(parameters connection.APIRequestParameters) error
}

// doctorServiceProbes returns the read-only probe for each service.
//
// The paginated variant of each list call is used deliberately. The non-paginated
// variants are implemented with connection.InvokeRequestAll, which walks every page -
// with the single-item page size used below, that would issue one request per resource
// held by the account. The paginated variants issue exactly one request.
//
// Note that billing is probed with GetCloudCostsPaginated rather than GetCards - a
// diagnostic command must not retrieve stored payment card data
func doctorServiceProbes(c client.Client) []doctorServiceProbe {
	return []doctorServiceProbe{
		{
			name: "account",
			probe: func(parameters connection.APIRequestParameters) error {
				_, err := c.AccountService().GetClientsPaginated(parameters)
				return err
			},
		},
		{
			name: "billing",
			probe: func(parameters connection.APIRequestParameters) error {
				_, err := c.BillingService().GetCloudCostsPaginated(parameters)
				return err
			},
		},
		{
			name: "cloudflare",
			probe: func(parameters connection.APIRequestParameters) error {
				_, err := c.CloudflareService().GetAccountsPaginated(parameters)
				return err
			},
		},
		{
			name: "ddosx",
			probe: func(parameters connection.APIRequestParameters) error {
				_, err := c.DDoSXService().GetDomainsPaginated(parameters)
				return err
			},
		},
		{
			name: "draas",
			probe: func(parameters connection.APIRequestParameters) error {
				_, err := c.DRaaSService().GetSolutionsPaginated(parameters)
				return err
			},
		},
		{
			name: "ecloud",
			probe: func(parameters connection.APIRequestParameters) error {
				_, err := c.ECloudService().GetVPCsPaginated(parameters)
				return err
			},
		},
		{
			name: "loadbalancer",
			probe: func(parameters connection.APIRequestParameters) error {
				_, err := c.LoadBalancerService().GetClustersPaginated(parameters)
				return err
			},
		},
		{
			name: "pss",
			probe: func(parameters connection.APIRequestParameters) error {
				_, err := c.PSSService().GetRequestsPaginated(parameters)
				return err
			},
		},
		{
			name: "registrar",
			probe: func(parameters connection.APIRequestParameters) error {
				_, err := c.RegistrarService().GetDomainsPaginated(parameters)
				return err
			},
		},
		{
			name: "safedns",
			probe: func(parameters connection.APIRequestParameters) error {
				_, err := c.SafeDNSService().GetZonesPaginated(parameters)
				return err
			},
		},
		{
			name: "ssl",
			probe: func(parameters connection.APIRequestParameters) error {
				_, err := c.SSLService().GetCertificatesPaginated(parameters)
				return err
			},
		},
	}
}

// doctorServiceChecks probes each service concurrently, establishing empirically which
// services the configured API key can reach
func doctorServiceChecks(c client.Client) []DoctorResult {
	probes := doctorServiceProbes(c)

	// A single page of one item is requested, as only reachability is of interest
	parameters := connection.APIRequestParameters{
		Pagination: connection.APIRequestPagination{PerPage: 1},
	}

	// Results are written to a pre-sized slice by index, keeping output ordering
	// deterministic regardless of the order in which the probes complete
	results := make([]DoctorResult, len(probes))

	var wg sync.WaitGroup
	for i, serviceProbe := range probes {
		wg.Go(func() {
			status, detail := doctorProbeStatus(serviceProbe.probe(parameters))
			results[i] = DoctorResult{
				Check:  fmt.Sprintf("service: %s", serviceProbe.name),
				Status: status,
				Detail: detail,
			}
		})
	}
	wg.Wait()

	return results
}

// doctorProbeStatus classifies the outcome of a service probe. A successful call is
// reachable even when the returned list is empty, as an empty list means the key was
// accepted and the account simply holds no resources of that type.
//
// A 401 or 403 is reported as a warning rather than a failure. API keys are scoped to
// particular applications and roles, so a key which cannot reach every service is usually
// correctly configured rather than broken. A key which is invalid outright is caught by
// the api authentication check, which probes the account service directly.
//
// A 5xx is reported as a failure, as it indicates the service itself is unavailable
// rather than anything to do with the key's permissions
func doctorProbeStatus(err error) (status string, detail string) {
	if err == nil {
		return doctorStatusOK, "reachable"
	}

	statusCode := httpStatusFromError(err)
	switch {
	case statusCode == 401:
		return doctorStatusWarn, "401 - api key not permitted for this service"
	case statusCode == 403:
		if blocked, ipAddress := doctorIPNotAllowed(err); blocked {
			return doctorStatusWarn, fmt.Sprintf("403 - ip address %s not allowed for this api key", ipAddress)
		}
		return doctorStatusWarn, "403 - api key lacks the required role for this service"
	case statusCode == 404:
		return doctorStatusWarn, "404 not found (service may not be provisioned)"
	case statusCode >= 500 && statusCode <= 599:
		return doctorStatusFail, fmt.Sprintf("%d server error", statusCode)
	default:
		return doctorStatusWarn, doctorTruncate(err.Error())
	}
}

// doctorApplicationChecks reports the declared service scopes and IP restrictions for the
// given API application. These cannot be determined for the key currently in use, as the
// API offers no endpoint for describing the caller
func doctorApplicationChecks(service account.AccountService, applicationID string) []DoctorResult {
	services := DoctorResult{Check: "application services"}
	restrictions := DoctorResult{Check: "application restrictions"}

	mapping, err := service.GetApplicationServices(applicationID)
	if err != nil {
		services.Status = doctorStatusFail
		services.Detail = doctorTruncate(err.Error())
	} else if len(mapping.Scopes) < 1 {
		services.Status = doctorStatusOK
		services.Detail = "no scopes declared"
	} else {
		scopes := make([]string, 0, len(mapping.Scopes))
		for _, scope := range mapping.Scopes {
			scopes = append(scopes, fmt.Sprintf("%s(%s)", scope.Service, strings.Join(scope.Roles, ",")))
		}

		services.Status = doctorStatusOK
		services.Detail = strings.Join(scopes, " ")
	}

	restriction, err := service.GetApplicationRestrictions(applicationID)
	switch {
	case err != nil:
		restrictions.Status = doctorStatusFail
		restrictions.Detail = doctorTruncate(err.Error())
	case len(restriction.IPRestrictionType) < 1 && len(restriction.IPRanges) < 1:
		// The API returns an empty array rather than an object when no restrictions are
		// configured, which the SDK unmarshals to the zero value
		restrictions.Status = doctorStatusOK
		restrictions.Detail = "no restrictions configured"
	default:
		restrictions.Status = doctorStatusWarn
		restrictions.Detail = fmt.Sprintf("%s restriction with %d range(s) - ensure your current IP is covered", restriction.IPRestrictionType, len(restriction.IPRanges))
	}

	return []DoctorResult{services, restrictions}
}

// doctorVersionCheck compares the running version against the latest GitHub release. Only
// release detection is performed - the binary is never modified
func doctorVersionCheck(currentVersion string) DoctorResult {
	result := DoctorResult{Check: "cli version"}

	if len(currentVersion) < 1 || currentVersion == "UNKNOWN" {
		result.Status = doctorStatusSkip
		result.Detail = "version unknown (non-release build)"
		return result
	}

	current, err := semver.ParseTolerant(currentVersion)
	if err != nil {
		result.Status = doctorStatusWarn
		result.Detail = fmt.Sprintf("unable to parse version '%s'", currentVersion)
		return result
	}

	latest, found, err := selfupdate.DetectLatest("ans-group/cli")
	if err != nil {
		// A failed lookup is not a failure of the CLI setup, e.g. when offline
		result.Status = doctorStatusWarn
		result.Detail = doctorTruncate(fmt.Sprintf("unable to determine latest version: %s", err))
		return result
	}

	if !found || latest == nil {
		result.Status = doctorStatusWarn
		result.Detail = "no releases found"
		return result
	}

	if current.GTE(latest.Version) {
		result.Status = doctorStatusOK
		result.Detail = fmt.Sprintf("v%s (latest)", current)
		return result
	}

	result.Status = doctorStatusWarn
	result.Detail = fmt.Sprintf("v%s (latest is v%s)", current, latest.Version)
	return result
}

// httpStatusFromError attempts to extract an HTTP status code from an SDK error,
// returning 0 if none could be determined
func httpStatusFromError(err error) int {
	if err == nil {
		return 0
	}

	if bodyErr, ok := errors.AsType[*connection.APIResponseBodyError](err); ok && len(bodyErr.Errors) > 0 {
		return bodyErr.Errors[0].Status
	}

	// Responses carrying no structured error body are reported by the SDK as a plain
	// error containing the status code
	matches := doctorUnexpectedStatusCodeRegex.FindStringSubmatch(err.Error())
	if len(matches) == 2 {
		statusCode, err := strconv.Atoi(matches[1])
		if err == nil {
			return statusCode
		}
	}

	return 0
}

// doctorIPNotAllowed reports whether err represents an API rejection caused by the
// caller's IP address not being on the key's allow list, along with the rejected
// address if one could be extracted. This is distinct from a role/permission failure,
// which also surfaces as a 403, so it must be checked explicitly rather than inferred
// from the status code alone
func doctorIPNotAllowed(err error) (blocked bool, ipAddress string) {
	if err == nil {
		return false, ""
	}

	matches := doctorIPNotAllowedRegex.FindStringSubmatch(err.Error())
	if len(matches) == 2 {
		return true, matches[1]
	}

	return false, ""
}

// doctorTruncate shortens a detail string so that it remains readable in table output.
// The string is measured and cut by rune, as API errors may carry non-ASCII characters
func doctorTruncate(s string) string {
	s = strings.TrimSpace(strings.ReplaceAll(s, "\n", " "))

	runes := []rune(s)
	if len(runes) <= doctorDetailMaxLength {
		return s
	}

	return string(runes[:doctorDetailMaxLength-3]) + "..."
}
