# ANS CLI

[![Build Status](https://travis-ci.org/ans/cli.svg?branch=master)](https://travis-ci.org/ans/cli)
[![MIT licensed](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

This is the official ANS command-line client, allowing for querying and controlling
supported ANS services.

The client utilises ANS APIs to provide access to most service features. You should refer to the 
[Getting Started](https://developers.ukfast.io/getting-started) section of the API documentation before 
proceeding below

## Building from source

Install the dependencies, golang and make, then build from the Makefile.

```
# make
```

To copy to /usr/local/bin, as root or using sudo

```
# make install
```

## Installation

The CLI is distributed as a single binary, and is available for Windows, Linux and Mac. This binary 
should be downloaded and placed into a directory included in your `PATH`. This would typically 
be `/usr/local/bin` on most Linux distributions

Pre-compiled binaries are available at [Releases](https://github.com/ans-group/cli/releases)

## Getting started

To get started, we will define a single environment variable to store our API key:

Bash:
> export ANS_API_KEY="iqmxgom0kairfnxzcopte5hx"

PowerShell:
> $env:ANS_API_KEY="iqmxgom0kairfnxzcopte5hx"

And away we go!

```
> ans safedns zone record show example.co.uk 3337874
+---------+--------------------+------+---------+---------------------------+----------+-----+
|   ID    |        NAME        | TYPE | CONTENT |        UPDATED AT         | PRIORITY | TTL |
+---------+--------------------+------+---------+---------------------------+----------+-----+
| 3337874 | test.example.co.uk | A    | 1.2.3.4 | 2019-03-19T16:33:55+00:00 |        0 |   0 |
+---------+--------------------+------+---------+---------------------------+----------+-----+
```

## Configuration

There are two available methods for configuring the CLI; Environment variables and configuration file. 
Both of these methods are explained below.

### Configuration File

The configuration file is read from
`$HOME/.ans{.extension}` by default (extension being one of the `viper` supported formats such as `yml`, `yaml`, `json`, `toml` etc.). This path can be overridden with the `--config` flag.

Values defined in the configuration file take precedence over environment variables.

### Schema

* `api_key`: (String) *Required* API key for authenticating with API
* `api_timeout_seconds`: (int) HTTP timeout for API requests. Default: `90`
* `api_uri`: (string) API URI. Default: `api.ukfast.io`
* `api_insecure`: (bool) Specifies to ignore API certificate validation checks
* `api_debug`: (bool) Specifies for debug messages to be output to stderr
* `api_pagination_perpage` (int) Specifies the per-page for paginated requests

### Contexts

Contexts can be defined in the config file to allow for different sets of configuration to be defined:

```yaml
contexts:
  testcontext1:
    api_key: mykey1
  testcontext2:
    api_key: mykey2
current_context: testcontext1
```

The current context can also be overridden with the `--context` flag

### Commands

The configuration file can be manipulated using the `config` subcommand, for example:


```
> ans config context update --current --api-key test1
> ans config context update someothercontext --api-key test1
> ans config context switch someothercontext
```

### Environment variables

Environment variables can be used to configure/manipulate the CLI. These variables match the naming of directives in the configuration file 
defined above, however are uppercased and prefixed with `UKF`, such as `ANS_API_KEY`

## Output Formatting

The output of all commands is determined by a single global flag `--output` / `-o`.
In addition to output, there are several output modifier flags which are explained below.

### Table (Default)

The default output format for the CLI is `Table`, which will be used when the value of the `--output` flag
is `table` or unspecified:

```
> ans safedns zone record list example.co.uk
+---------+--------------------+------+-----------------------------------------------------------------------+---------------------------+----------+-------+
|   ID    |        NAME        | TYPE |                                CONTENT                                |        UPDATED AT         | PRIORITY |  TTL  |
+---------+--------------------+------+-----------------------------------------------------------------------+---------------------------+----------+-------+
| 3337865 | ns0.ans.uk         | NS   | 185.226.220.128                                                       | 2019-03-19T16:31:48+00:00 |        0 |     0 |
| 3337868 | ns1.ans.uk         | NS   | 185.226.221.128                                                       | 2019-03-19T16:31:48+00:00 |        0 |     0 |
| 3337871 | example.co.uk      | SOA  | ns0.ans.uk support.ans.co.uk 2019031901 7200 3600 604800 86400        | 2019-03-19T16:31:48+00:00 |        0 | 86400 |
| 3337874 | test.example.co.uk | A    | 1.2.3.4                                                               | 2019-03-19T16:33:55+00:00 |        0 |     0 |
+---------+--------------------+------+-----------------------------------------------------------------------+---------------------------+----------+-------+
```

The [Property Modifier](#property) is available for this format

### List

Results can be output as a list using the `list` format:

```
> ans safedns zone record show example.co.uk 3337874 --output list
id         : 3337874
name       : test.example.co.uk
type       : A
content    : 1.2.3.4
updated_at : 2019-03-19T16:33:55+00:00
priority   : 0
ttl        : 0
```

The [Property Modifier](#property) is available for this format

### JSON

Results can be output in JSON using the `json` format:

```
> ans safedns zone record show example.co.uk 3337874 --output json
[{"id":3337874,"template_id":0,"name":"test.example.co.uk","type":"A","content":"1.2.3.4","updated_at":"2019-03-19T16:33:55+00:00","ttl":0,"priority":0}]
```

### YAML

Results can be output in YAML using the `yaml` format:

```
> ans safedns zone record show example.co.uk 3337874 --output yaml
- id: 3337874
  templateid: 0
  name: test.example.co.uk
  type: A
  content: 1.2.3.4
  updatedat: "2019-03-19T16:33:55+00:00"
  ttl: 0
  priority: 0
```

### Value

Results can be output with a value or set of values using the `value` format:

```
> ans safedns zone record show example.co.uk 3337874 --output value
3337874 test.example.co.uk A 1.2.3.4 2019-03-19T16:33:55+00:00 0 0
```

```
> ans safedns zone record show example.co.uk 3337874 --output value --property id
3337874
```

The [Property Modifier](#property) is available for this format

### CSV

Results can be output as CSV using the `csv` format:

```
> ans safedns zone record show example.co.uk 3337874 --output csv
id,name,type,content,updated_at,priority,ttl
3337874,test.example.co.uk,A,1.2.3.4,2019-03-19T16:33:55+00:00,0,0
```

The [Property Modifier](#property) is available for this format

### Template

Results can be output using a supplied Golang template string using the `template` format
in conjunction with output arguments, e.g.

```
> ans safedns zone record list example.co.uk --output template="Record name: {{ .Name }}, Type: {{ .Type }}"
Record name: ns0.ans.uk, Type: NS
Record name: ns1.ans.uk, Type: NS
Record name: example.co.uk, Type: SOA
Record name: test.example.co.uk, Type: A
```

### JSON path

Results can be output via JSON Path using the `jsonpath` format
in conjunction with output arguments, e.g.

```
> ans safedns zone record list example.co.uk --output jsonpath="{[*].name}"
example.co.uk example.co.uk example.co.uk test.example.co.uk
```


## Output modifiers

### Property

Some output formats support the `--property` output modifier.

Required properties can be specified with the `--property` format modifer flag:

```
> ans safedns zone record show example.co.uk 3337874 --output value --property name
test.example.co.uk
```

The property modifier accepts a comma-delimited list of property names, and is also repeatable:

```
> ans safedns zone record show example.co.uk 3337874 --output value --property id,name --property content
3337874 test.example.co.uk 1.2.3.4
```

The property modifier also accepts globbing e.g. `*`, `some*`, `*thing`


## Filtering

When using `list` commands, filtering is available via the `--filter` flag. Additionally, optional operators are available.

### Examples

```
--filter id=123
--filter id:lt=10
```

Additionally, the `lk` filter is inferred when a glob `*` is included in the filter value (when operator is omitted)


## Sorting

When using `list` commands, sorting is available via the `--sort` flag.

### Examples

```
--sort id
--sort id:desc
```

## Updates

The CLI has self-update functionality, which can be invoked via the command `update`:

```
> ans update
```

This command in-place updates the CLI, with the old binary moved to `ans.old` (`ans.old.exe` on Windows), should roll-back be required.

### Migrating from the UKFast CLI

If you are upgrading from the old `ukfast` client, you will need to install this client from scratch. You will need to rename `~/.ukfast.yml` to `~/.ans.yml` and ensure any environment variables are updated to use the new `ANS_` prefix, e.g. `UKF_ECLOUD_VPC=true` becomes `ANS_ECLOUD_VPC=true`.

## Shell autocompletions

The CLI supports generating shell completions for the following shells:

* Bash
* Zsh
* PowerShell

The commands at `ans completion <shell: bash|zsh|powershell>` provide help for installation on different platforms

## eCloud VPC resources

eCloud VPC resource commands are available by default under the `ecloud` subcommand.

To display only VPC commands, the following environment variable can be set:

```
> export ANS_ECLOUD_VPC=true
```

To display only original commands, the following environment variable can be set:

```
> export ANS_ECLOUD=true
```