# ANS CLI

`ans` is a Cobra CLI over [`ans-group/sdk-go`](https://github.com/ans-group/sdk-go). Every API call goes through the SDK; the [eCloud OpenAPI spec](https://developers.ukfast.io/api/documentation/ecloud/2) is reference only.

## Commands

- Build: `go build -o ans` (`make build` adds version ldflags).
- Test: `go test ./cmd/<service>/...` for the package you touched; `go test ./...` before finishing.
- Lint: `golangci-lint run`.
- Format: `goimports -w` on every changed file.
- Mocks: `cd test/mocks && ./generate.sh`. The script post-processes `mockgen` output to strip generic import paths, so it is the only working generator (`generate.go.disabled` is off deliberately). Needs `mockgen` (`github.com/golang/mock/mockgen@v1.6.0`) on `PATH`.

`go.mod` / `mise.toml` define the toolchain

## Layout

- `cmd/root.go` registers each service's root command.
- `cmd/<service>/<service>.go`: service root command, the `<service>CobraRunEFunc` adapter, shared flag parsers and wait funcs.
- `cmd/<service>/<service>_<resource>[_<subresource>].go`, each with a sibling `_test.go`.
- `cmd/<service>/output.go`: `<Resource>Collection` output types.
- `internal/pkg/`: `factory` (client factory), `helper` (flags to API params, `WaitForCommand`), `output` (rendering), `clierrors`, `input`, `resource` (name-to-ID lookup).
- `test/`: `mocks` (generated), `test_output.AssertErrorOutput`, `test_input`.
- eCloud has two command sets: v1 (`ecloud_v1_*`) and VPC (everything else). `ANS_ECLOUD` / `ANS_ECLOUD_VPC` gate which set `ECloudRootCmd` registers; new eCloud resources are VPC and go in the VPC block.
- Files with a `_windows.go` twin carry build tags; change both.

## Adding or changing a command

1. Confirm the SDK exposes the method and fields: `go doc github.com/ans-group/sdk-go/pkg/service/<service>`. If missing, the SDK changes first; then `go get github.com/ans-group/sdk-go@<version>` and regenerate mocks.
2. Copy the nearest **sibling** resource file. The shape every sibling follows:
   - `<x>RootCmd(f)` groups the subcommands and is added to the service root.
   - Each subcommand is a builder `<x>Cmd(f)` (Use, Short, Long, Example, Args, flags) plus a body `<x>(service, cmd, args) error`, wired by `<service>CobraRunEFunc`. Tests call the body directly with a mock service.
   - `list`: filter flags through `helper.GetAPIRequestParametersFromFlags` with `helper.NewStringFilterFlagOption`.
   - `show` / `update` / `delete` take many IDs (`Use: "show <vpc: id>..."`), loop over them, report each failure and `continue`.
   - `update`: set a patch field only when `cmd.Flags().Changed(name)`.
   - Asynchronous resources: a `--wait` flag driving `helper.WaitForCommand` with a sync-status or task wait func.
   - Output: `output.CommandOutput(cmd, <Resource>Collection(items))`, or `output.CommandOutputPaginated` for paginated SDK calls.
3. Add the `<Resource>Collection` to `output.go` with `DefaultColumns()`. Columns default to the SDK struct's JSON tags; implement `Fields`, `HiddenFields`, `FieldNameOverrides` or `FieldValueHandlers` only to change that.
4. Test like the sibling `_test.go`: `Test_<x>Cmd_Args` for arg validation and `Test_<x>` with one `t.Run` subtest per branch, gomock `EXPECT`s, and exact error-string assertions. Packages with no existing tests get table-driven tests.

Done when every SDK call the command makes has a success subtest and an error subtest, `go test ./...` and `go build` pass, and `golangci-lint run` reports nothing in the files you changed (the repo has a few existing issues elsewhere).

## Errors

- Fatal errors are returned: `fmt.Errorf("error retrieving DHCPs: %w", err)`.
- Per-ID errors inside a multi-ID loop are printed and skipped: `output.OutputWithErrorLevelf("Error retrieving DHCP [%s]: %s", id, err)`; this sets a non-zero exit code without aborting the remaining IDs.
- Match the message style the package already uses (most `cmd/` packages use the `error <verb>ing <resource>:` form above; `cmd/account` uses `account:`). A package with no established style prefixes errors with its name, e.g. `foo: this has failed`.
- Invalid flag values return `clierrors.NewErrInvalidFlagValue`.

## Go style

Write Go for the version in `go.mod`: `any`, generics, `errors.Is` / `errors.As`, `%w` wrapping, and `slices` / `strings` / `io` / `os` from the standard library in place of hand-rolled helpers or `ioutil`. Doc-comment exported identifiers. Use the /use-modern-go skill, if available.
