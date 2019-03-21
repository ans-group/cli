# UKFast CLI

[![pipeline status](https://github.com/ukfast/cli/badges/master/pipeline.svg)](https://github.com/ukfast/cli/commits/master) [![coverage report](https://github.com/ukfast/cli/badges/master/coverage.svg)](https://github.com/ukfast/cli/commits/master)

This is the official UKFast command-line client, allowing for querying and controlling
supported UKFast services.

The client utilises UKFast APIs to provide access to most service features. You should refer to the 
[Getting Started](https://developers.ukfast.io/getting-started) section of the API documentation before 
proceeding below

## Installation

The CLI is distributed as a single binary, and is available for Windows, Linux and Mac. This binary 
should be downloaded and placed into a directory included in your `PATH`. This would typically 
be `/usr/local/bin` on most Linux distributions

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

### Environment variables

Environment variables can be used to configure/manipulate the CLI

#### Required

* `UKF_API_KEY`: API key for interacting with UKFast APIs

#### Debug

* `UKF_API_TIMEOUT_SECONDS`: (int) HTTP timeout for API requests. Default: `90`
* `UKF_API_URL`: (string) API URL. Default: `api.ukfast.io`
* `UKF_API_INSECURE`: (bool) Specifies to ignore API certificate validation checks
* `UKF_API_DEBUG`: (bool) Specifies for debug messages to be output to stderr
* `UKF_API_PAGINATION_PERPAGE` (int) Specifies the per-page for paginated requests

### Configuration File

An alternative to environment variables is using a configuration file, which is read from
`$HOME/.ukfast.yaml` by default. This path can be overridden with the `--config` flag. 
Values defined in the configuration file take precedence over environment variables.

#### Required

* `api_key`: API key for interacting with UKFast APIs

#### Debug

* `api_timeout_seconds`: (int) HTTP timeout for API requests. Default: `90`
* `api_url`: (string) API URL. Default: `api.ukfast.io`
* `api_insecure`: (bool) Specifies to ignore API certificate validation checks
* `api_debug`: (bool) Specifies for debug messages to be output to stderr
* `api_pagination_perpage` (int) Specifies the per-page for paginated requests


## Output Formatting

The output of all commands is determined by a single global flag `--format` / `-f`.
In additional to format, there are several format modifier flags which are explained below.

### Table (Default)

The default output format for the CLI is `Table`, which will be used when the `--format` flag
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

### JSON

Results can be output in JSON using the `json` format:

```
> ukfast safedns zone record show example.co.uk 3337874 --format json
[{"id":3337874,"template_id":0,"name":"test.example.co.uk","type":"A","content":"1.2.3.4","updated_at":"2019-03-19T16:33:55+00:00","ttl":0,"priority":0}]
```

### Value

Results can be output with a value or set of values using the `value` format:

```
> ukfast safedns zone record show example.co.uk 3337874 --format value
3337874 test.example.co.uk A 1.2.3.4 2019-03-19T16:33:55+00:00 0 0
```

```
> ukfast safedns zone record show example.co.uk 3337874 --format value --property id
3337874
```

The [Property Modifier](#property) is available for this format

### CSV

Results can be output as CSV using the `csv` format:

```
> ukfast safedns zone record show example.co.uk 3337874 --format csv
id,name,type,content,updatedat,priority
3337874,test.example.co.uk,A,1.2.3.4,2019-03-19T16:33:55+00:00,0,0
```

The [Property Modifier](#property) is available for this format

### Template

Results can be output using a supplied Golang template string using the `template` format
in conjunction with the `--outputtemplate` modifier flag:

```
> ukfast safedns zone record list example.co.uk --format template --outputtemplate "Record name: {{ .Name }}, Type: {{ .Type }}"
Record name: ns0.ukfast.net, Type: NS
Record name: ns1.ukfast.net, Type: NS
Record name: example.co.uk, Type: SOA
Record name: test.example.co.uk, Type: A
```


## Output modifiers

### Property

Some output formats support the `--property` output modifier.

Required properties can be specified with the `--property` format modifer flag:

```
> ukfast safedns zone record show example.co.uk 3337874 --format value --property name
test.example.co.uk
```

The property modifier accepts a comma-delimited list of property names, and is also repeatable:

```
> ukfast safedns zone record show example.co.uk 3337874 --format value --property id,name --property content
3337874 test.example.co.uk 1.2.3.4
```

The property modifier also accepts a single wildcard/glob `*`, which denotes all properties
should be returned


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
