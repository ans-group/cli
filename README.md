# UKFast CLI

[![Build Status](https://travis-ci.org/ukfast/cli.svg?branch=master)](https://travis-ci.org/ukfast/cli)
[![MIT licensed](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

This is the official UKFast command-line client, allowing for querying and controlling
supported UKFast services.

The client utilises UKFast APIs to provide access to most service features. You should refer to the 
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

Pre-compiled binaries are available at [Releases](https://github.com/ukfast/cli/releases)

## Getting started

To get started, we will define a single environment variable to store our API key:

Bash:
> UKF_API_KEY="iqmxgom0kairfnxzcopte5hx"

PowerShell:
> $env:UKF_API_KEY="iqmxgom0kairfnxzcopte5hx"

And away we go!

```
> ukfast safedns zone record show example.co.uk 3337874
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
`$HOME/.ukfast{.extension}` by default (extension being one of the `viper` supported formats such as `yml`, `yaml`, `json`, `toml` etc.). This path can be overridden with the `--config` flag. 
Values defined in the configuration file take precedence over environment variables.

#### Required

* `api_key`: API key for interacting with UKFast APIs

#### Debug

* `api_timeout_seconds`: (int) HTTP timeout for API requests. Default: `90`
* `api_uri`: (string) API URI. Default: `api.ukfast.io`
* `api_insecure`: (bool) Specifies to ignore API certificate validation checks
* `api_debug`: (bool) Specifies for debug messages to be output to stderr
* `api_pagination_perpage` (int) Specifies the per-page for paginated requests

### Environment variables

Environment variables can be used to configure/manipulate the CLI. These variables match the naming of directives in the configuration file 
defined above, however are uppercased and prefixed with `UKF`, such as `UKF_API_KEY`


## Output Formatting

The output of all commands is determined by a single global flag `--output` / `-o`.
In addition to output, there are several output modifier flags which are explained below.

### Table (Default)

The default output format for the CLI is `Table`, which will be used when the value of the `--output` flag
is `table` or unspecified:

```
> ukfast safedns zone record list example.co.uk
+---------+--------------------+------+-----------------------------------------------------------------------+---------------------------+----------+-------+
|   ID    |        NAME        | TYPE |                                CONTENT                                |        UPDATED AT         | PRIORITY |  TTL  |
+---------+--------------------+------+-----------------------------------------------------------------------+---------------------------+----------+-------+
| 3337865 | ns0.ukfast.net     | NS   | 185.226.220.128                                                       | 2019-03-19T16:31:48+00:00 |        0 |     0 |
| 3337868 | ns1.ukfast.net     | NS   | 185.226.221.128                                                       | 2019-03-19T16:31:48+00:00 |        0 |     0 |
| 3337871 | example.co.uk      | SOA  | ns0.ukfast.net support.ukfast.co.uk 2019031901 7200 3600 604800 86400 | 2019-03-19T16:31:48+00:00 |        0 | 86400 |
| 3337874 | test.example.co.uk | A    | 1.2.3.4                                                               | 2019-03-19T16:33:55+00:00 |        0 |     0 |
+---------+--------------------+------+-----------------------------------------------------------------------+---------------------------+----------+-------+
```

The [Property Modifier](#property) is available for this format

### List

Results can be output as a list using the `list` format:

```
> ukfast safedns zone record show example.co.uk 3337874 --output list
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
> ukfast safedns zone record show example.co.uk 3337874 --output json
[{"id":3337874,"template_id":0,"name":"test.example.co.uk","type":"A","content":"1.2.3.4","updated_at":"2019-03-19T16:33:55+00:00","ttl":0,"priority":0}]
```

### Value

Results can be output with a value or set of values using the `value` format:

```
> ukfast safedns zone record show example.co.uk 3337874 --output value
3337874 test.example.co.uk A 1.2.3.4 2019-03-19T16:33:55+00:00 0 0
```

```
> ukfast safedns zone record show example.co.uk 3337874 --output value --property id
3337874
```

The [Property Modifier](#property) is available for this format

### CSV

Results can be output as CSV using the `csv` format:

```
> ukfast safedns zone record show example.co.uk 3337874 --output csv
id,name,type,content,updated_at,priority,ttl
3337874,test.example.co.uk,A,1.2.3.4,2019-03-19T16:33:55+00:00,0,0
```

The [Property Modifier](#property) is available for this format

### Template

Results can be output using a supplied Golang template string using the `template` format
in conjunction with output arguments, e.g.

```
> ukfast safedns zone record list example.co.uk --output template="Record name: {{ .Name }}, Type: {{ .Type }}"
Record name: ns0.ukfast.net, Type: NS
Record name: ns1.ukfast.net, Type: NS
Record name: example.co.uk, Type: SOA
Record name: test.example.co.uk, Type: A
```

### JSON path

Results can be output via JSON Path using the `jsonpath` format
in conjunction with output arguments, e.g.

```
> ukfast safedns zone record list example.co.uk --output jsonpath="{[*].name}"
example.co.uk example.co.uk example.co.uk test.example.co.uk
```


## Output modifiers

### Property

Some output formats support the `--property` output modifier.

Required properties can be specified with the `--property` format modifer flag:

```
> ukfast safedns zone record show example.co.uk 3337874 --output value --property name
test.example.co.uk
```

The property modifier accepts a comma-delimited list of property names, and is also repeatable:

```
> ukfast safedns zone record show example.co.uk 3337874 --output value --property id,name --property content
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
> ukfast update
```

This command in-place updates the CLI, with the old binary moved to `ukfast.old` (`ukfast.old.exe` on Windows), should roll-back be required.

## Shell autocompletions

The CLI supports generating shell completions for the following shells:

* Bash
* Zsh
* PowerShell

The commands at `ukfast completion <shell: bash|zsh|powershell>` provide help for installation on different platforms

## eCloud V2 resources

eCloud V2 resource commands are available by default under the `ecloud` subcommand.

To display only V2 commands, the following environment variable can be set:

```
> export UKF_ECLOUD_V2=true
```

To display only V2 commands, the following environment variable can be set:

```
> export UKF_ECLOUD_V1=true
```