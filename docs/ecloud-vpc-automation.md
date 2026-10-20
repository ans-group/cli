# Automating eCloud VPC with the ANS CLI

This guide shows how to build and manage eCloud VPC (v2) infrastructure from scripts and pipelines with the `ans` CLI. It covers creating an instance, protecting it with a Backup Gateway, controlling traffic with TCP, UDP and ICMP firewall rules, and checking that everything ended up the way you asked.

Every step uses the CLI first. Each one also has a `curl` equivalent against the eCloud API, which is useful when you need a field the CLI does not expose yet, or when you are writing your own tooling. The [ANS Go SDK](https://github.com/ans-group/sdk-go) wraps the same API if you prefer to work in Go.

## Contents

- [Before you start](#before-you-start)
- [How eCloud VPC resources fit together](#how-ecloud-vpc-resources-fit-together)
- [Asynchronous operations: tasks and sync status](#asynchronous-operations-tasks-and-sync-status)
- [Discovering IDs and specification values](#discovering-ids-and-specification-values)
- [Example 1: build the network](#example-1-build-the-network)
- [Example 2: create an instance](#example-2-create-an-instance)
- [Example 3: enable backups](#example-3-enable-backups)
- [Example 4: manage firewall rules](#example-4-manage-firewall-rules)
- [Verifying configuration](#verifying-configuration)
- [Cleaning up](#cleaning-up)
- [Troubleshooting](#troubleshooting)
- [Complete example script](#complete-example-script)

## Before you start

You need a working `ans` installation and an API key with read/write access to eCloud. The [README](../README.md) explains how to [install the CLI](../README.md#installation), [supply your API key](../README.md#getting-started) and [manage contexts](../README.md#contexts). This guide does not repeat those steps.

The README also documents the global options this guide relies on most:

- [Output formatting](../README.md#output-formatting), in particular `-o value`, `-o json` and `-o jsonpath=...`
- [`--property`](../README.md#property) to select fields
- [Filtering](../README.md#filtering) with `--filter`
- [Sorting](../README.md#sorting) with `--sort`

The `curl` examples use `jq` to parse responses. They assume your API key is in `ANS_API_KEY` and use the following shorthand:

```bash
export ECLOUD="https://api.ukfast.io/ecloud/v2"
H=(-H "Authorization: ${ANS_API_KEY}" -H "Content-Type: application/json" -H "Accept: application/json")

curl -s "${H[@]}" "$ECLOUD/vpcs" | jq
```

The API key goes in the `Authorization` header as is, with no `Bearer` prefix.

All IDs in this guide (`vpc-abcdef12` and so on) are placeholders. Substitute the values returned in your own account.

### Writing scripts that fail safely

The CLI exits with a non-zero status when a request fails or a `--wait` ends in failure or timeout, so `set -euo pipefail` works as expected. To capture an ID from a create command, ask for the `id` property as a plain value:

```bash
VPC_ID=$(ans ecloud vpc create --region reg-abcdef12 --name demo -o value --property id)
```

`-o value` prints one line per resource with no headers. With several properties it separates them with spaces.

## How eCloud VPC resources fit together

Most resources belong to a parent, and you create them top down:

```text
Region
└── VPC
    ├── Router  (in one availability zone)
    │   ├── Network (subnet)
    │   │   └── Instance ── NIC ── IP address ◄── Floating IP
    │   ├── Firewall policy
    │   │   └── Firewall rule
    │   │       └── Firewall rule port
    │   ├── Backup Gateway
    │   └── Monitoring Gateway
    └── Floating IP
```

An instance's availability zone comes from its network's router. Backup Gateways, floating IPs and the instances that use them must all be in the same availability zone.

Each resource type has a fixed ID prefix, which makes it easy to tell what an ID refers to:

| Prefix  | Resource                    | CLI command                         |
|---------|-----------------------------|-------------------------------------|
| `reg-`  | Region                      | `ans ecloud region`                 |
| `az-`   | Availability zone           | `ans ecloud availabilityzone`       |
| `vpc-`  | VPC                         | `ans ecloud vpc`                    |
| `rtr-`  | Router                      | `ans ecloud router`                 |
| `rtp-`  | Router throughput           | `ans ecloud routerthroughput`       |
| `net-`  | Network                     | `ans ecloud network`                |
| `img-`  | Image                       | `ans ecloud image`                  |
| `rt-`   | Resource tier               | `ans ecloud resourcetier`           |
| `i-`    | Instance                    | `ans ecloud instance`               |
| `nic-`  | Network interface           | `ans ecloud nic`                    |
| `ip-`   | IP address                  | `ans ecloud nic ipaddress`          |
| `fip-`  | Floating IP                 | `ans ecloud floatingip`             |
| `ssh-`  | SSH key pair                | `ans ecloud sshkeypair`             |
| `fwp-`  | Firewall policy             | `ans ecloud firewallpolicy`         |
| `fwr-`  | Firewall rule               | `ans ecloud firewallrule`           |
| `fwrp-` | Firewall rule port          | `ans ecloud firewallrule port`      |
| `bgw-`  | Backup Gateway              | `ans ecloud backupgateway`          |
| `bgws-` | Backup Gateway specification| `ans ecloud backupgateway spec`     |
| `task-` | Task                        | `ans ecloud task`                   |

Every resource type supports `list` and `show`, and most support `create`, `update` and `delete`. Run any command with `--help` to see its flags.

## Asynchronous operations: tasks and sync status

Creating, changing or deleting infrastructure takes time, so the API accepts the request straight away and does the work in the background. Your automation has to wait for that work to finish before it moves on to anything that depends on it. A network cannot be created until its router is ready, for example.

Two fields tell you where things stand.

A task records one piece of background work. It has an ID starting `task-`, a `name` that says what it is doing (such as `sync_update` or `sync_delete`), the `resource_id` it acts on, and a `status` of `in-progress`, `complete` or `failed`.

Every asynchronous resource also carries a `sync` block that reflects its most recent task:

```json
"sync": { "status": "complete", "type": "update" },
"task": { "in_progress": false, "status": "complete" }
```

`sync.status` is `in-progress`, `complete` or `failed`. `sync.type` is `update` for creates and changes, and `delete` while the resource is being removed. `task.in_progress` is `true` while any task for the resource is still running.

A resource is ready to use, or to have children created under it, when `sync.status` is `complete`.

### Waiting with `--wait`

The simplest approach is to add `--wait` to commands that support it. The CLI then blocks until the work completes and exits non-zero if it fails:

```bash
ans ecloud router create --vpc vpc-abcdef12 --availability-zone az-abcdef12 --name web --wait
```

By default `--wait` gives up after 1200 seconds and checks every 5 seconds. Change both per context, or per run with environment variables:

```bash
ans config context update mycontext --command-wait-timeout-seconds 3600 --command-wait-sleep-seconds 10

ANS_COMMAND_WAIT_TIMEOUT_SECONDS=3600 ans ecloud instance create ... --wait
```

### Waiting on a task

Commands that start a single piece of work, such as power operations, Backup Gateway changes and firewall changes, wait on that task when you pass `--wait`:

```bash
ans ecloud instance restart i-abcdef12 --wait
```

The CLI does not print task IDs. If you have one, for example from an API response or from a task list, wait for it with `ans ecloud task wait`:

```bash
TASK_ID=$(ans ecloud instance task list i-abcdef12 --filter status=in-progress -o value --property id | head -n1)
ans ecloud task wait "$TASK_ID"
```

`task wait` accepts `--status` if you want to wait for something other than `complete`. To inspect tasks directly:

```bash
ans ecloud task show task-abcdef12
ans ecloud task list --filter resource_id=i-abcdef12 --filter status=in-progress
```

Most resources also have a `task list` subcommand that shows their own history, newest work included:

```bash
ans ecloud instance task list i-abcdef12
ans ecloud vpc task list vpc-abcdef12
ans ecloud router task list rtr-abcdef12
ans ecloud firewallpolicy task list fwp-abcdef12
```

### Polling sync status yourself

Create commands return the new resource rather than a task ID. If you would rather not block on `--wait`, for example because you are creating several resources in parallel, poll `sync.status` yourself. This function works for any resource type that has a `show` command:

```bash
# Usage: wait_for_sync <resource> <id> [timeout-seconds]
wait_for_sync() {
  local resource=$1 id=$2 timeout=${3:-1800} waited=0 status
  while (( waited < timeout )); do
    status=$(ans ecloud "$resource" show "$id" -o value --property sync.status)
    case "$status" in
      complete) return 0 ;;
      failed)   echo "$resource $id failed to sync" >&2; return 1 ;;
    esac
    sleep 10; waited=$(( waited + 10 ))
  done
  echo "timed out waiting for $resource $id" >&2
  return 1
}

wait_for_sync instance i-abcdef12
```

### Polling with curl

Create and update requests return `202 Accepted` with the new resource's ID and the ID of the task doing the work:

```json
{
  "data": { "id": "i-abcdef12", "task_id": "task-abcdef12" },
  "meta": { "location": "https://api.ukfast.io/ecloud/v2/instances/i-abcdef12" }
}
```

Delete requests and actions such as power on/off return `202` with `data.task_id` only. Poll the task until it finishes:

```bash
wait_for_task() {
  local task_id=$1 status
  while :; do
    status=$(curl -s "${H[@]}" "$ECLOUD/tasks/$task_id" | jq -r '.data.status')
    case "$status" in
      complete) return 0 ;;
      failed)   echo "task $task_id failed" >&2; return 1 ;;
    esac
    sleep 10
  done
}
```

You can poll the resource instead with `GET $ECLOUD/instances/i-abcdef12` and read `.data.sync.status`. Every resource also has a `/tasks` sub-path, for example `GET $ECLOUD/instances/i-abcdef12/tasks`.

### Conflicts while work is in progress

Two rules decide when you can make the next change:

- You cannot change a resource while it has a task in progress. The API returns `409 Conflict` with "The specified resource has a task in progress". Wait and retry. Changes to firewall rules and ports count as changes to their firewall policy, so you need to wait for the policy too.
- You cannot create a child under a parent whose `sync.status` is not `complete`. The API returns `422` against the parent's ID field. If the parent's status is `failed`, retrying will not help; contact support.

A failed task does not include a reason. If a task fails and the cause is not obvious from your request, raise it with ANS support and quote the task ID.

## Discovering IDs and specification values

Everything you pass to a create command can be looked up from the CLI. Use `--filter` to narrow lists (see the [README](../README.md#filtering)).

### Regions and availability zones

```bash
ans ecloud region list
ans ecloud availabilityzone list
ans ecloud availabilityzone list --filter region_id=reg-abcdef12
```

`curl -s "${H[@]}" "$ECLOUD/availability-zones?region_id=reg-abcdef12" | jq '.data[] | {id, name}'`

### Router throughput

Throughput options differ by availability zone:

```bash
ans ecloud routerthroughput list --filter availability_zone_id=az-abcdef12
```

If you do not specify one when creating a router, it gets the zone's default (25Mb).

### Resource tiers and host groups

Resource tiers group the hosts your instances run on. Each availability zone has a default tier, which instances use unless you choose another. Dedicated host groups are an alternative to a tier; an instance can have one or the other.

```bash
ans ecloud resourcetier list --filter availability_zone_id=az-abcdef12
ans ecloud hostgroup list --filter vpc_id=vpc-abcdef12
```

### Images

Public images are provided by ANS. Private images are ones you have created from your own instances.

```bash
ans ecloud image list --filter visibility=public
ans ecloud image list --filter visibility=public --filter "name=*Ubuntu*"
ans ecloud image list --filter availability_zone_id=az-abcdef12
ans ecloud image show img-abcdef12
```

Two more commands tell you how to deploy an image. `image metadata list` shows constraints such as the minimum volume size, and `image parameter list` shows the values you can pass at deploy time:

```bash
ans ecloud image metadata list img-abcdef12
ans ecloud image parameter list img-abcdef12
```

Parameters marked `required` must be supplied when you create an instance from that image. See [Passing image parameters and a user script](#passing-image-parameters-and-a-user-script).

You can pass an image ID or its exact name to `ans ecloud instance create --image`. Names are not unique across availability zones, so use IDs in automation.

### Instance sizing

There is no list of instance sizes; you choose CPU, RAM and disk directly:

| Setting | CLI flag | Rule |
|---------|----------|------|
| CPU | `--vcpu-sockets`, `--vcpu-cores-per-socket` | Both default to 1. The default maximum is 10 sockets and 10 vCPUs in total; some resource tiers, host groups and images allow different limits. |
| RAM | `--ram` | In MiB, from 1024 to 65536, and a multiple of 8. |
| Disk | `--volume` | In GB, up to 3000. The minimum comes from the image's `ukfast.spec.volume.min` metadata, otherwise 40 GB for Linux and 80 GB for Windows. |

### Backup Gateway specifications

```bash
ans ecloud backupgateway spec list
```

Each specification describes the number of instances it is designed for, for example "Medium, Recommended for up to 50 Instances". Pick one with room for growth.

### SSH key pairs

Upload a public key once and reference it when creating Linux instances:

```bash
ans ecloud sshkeypair create --name deploy --public-key-file ~/.ssh/id_ed25519.pub
ans ecloud sshkeypair list
```

## Example 1: build the network

This example creates a VPC with one router and one network.

### Create the VPC

```bash
REGION_ID=reg-abcdef12
AZ_ID=az-abcdef12

VPC_ID=$(ans ecloud vpc create --region "$REGION_ID" --name demo-vpc --wait -o value --property id)
```

```bash
curl -s "${H[@]}" -X POST "$ECLOUD/vpcs" \
  -d '{"name": "demo-vpc", "region_id": "reg-abcdef12"}' | jq
```

A new VPC is empty. It has no routers or networks until you create them.

### Create a router

```bash
RTR_ID=$(ans ecloud router create --vpc "$VPC_ID" --availability-zone "$AZ_ID" \
  --name demo-router --wait -o value --property id)
```

```bash
curl -s "${H[@]}" -X POST "$ECLOUD/routers" -d '{
  "name": "demo-router",
  "vpc_id": "vpc-abcdef12",
  "availability_zone_id": "az-abcdef12",
  "router_throughput_id": "rtp-abcdef12"
}' | jq
```

`router_throughput_id` (`--throughput`) is optional.

Each new router gets at least one firewall policy managed by ANS, such as System or Management Firewall Policy, at sequence 0. You can see these policies but you cannot change them. See [Example 4](#example-4-manage-firewall-rules).

### Create a network

```bash
NET_ID=$(ans ecloud network create --router "$RTR_ID" --subnet 10.0.0.0/24 \
  --name demo-net --wait -o value --property id)
```

```bash
curl -s "${H[@]}" -X POST "$ECLOUD/networks" \
  -d '{"name": "demo-net", "router_id": "rtr-abcdef12", "subnet": "10.0.0.0/24"}' | jq
```

The subnet has to meet these rules:

- It is a private (RFC 1918) range written in CIDR notation, such as `10.0.0.0/24`.
- It is a /29 or larger.
- It does not overlap any other network on the same router.
- It is not inside `192.168.0.0/16` or `10.255.0.0/16`, which are reserved.

Some addresses in each subnet are reserved and cannot be given to instances: the network and broadcast addresses, the first usable address (the gateway), and two addresses used for DHCP. Run `ans ecloud network show net-abcdef12` to see the subnet.

### Optional: add the default firewall policies

You cannot add rules to the managed policies. To get a starting set of common rules, the router's `deploydefaults` command adds two editable policies:

| Policy | Sequence | Rules |
|--------|----------|-------|
| Infrastructure | 2 | Outbound ICMP ping, DNS (TCP/UDP 53), NTP (TCP/UDP 123), HTTP/HTTPS (TCP 80, 443) |
| Web Services | 4 | Inbound HTTP/HTTPS (TCP 80, 443) |

```bash
ans ecloud router deploydefaults "$RTR_ID"
```

```bash
curl -s "${H[@]}" -X POST "$ECLOUD/routers/rtr-abcdef12/configure-default-policies"
```

This request returns `202` with no task ID. Check that the policies exist and have finished syncing before you rely on them:

```bash
ans ecloud firewallpolicy list --router "$RTR_ID" -o table --property id,name,sequence,sync.status
```

### Alternative: let eCloud create the router and network

`ans ecloud vpc deploydefaults` creates a router and a network with default settings in the availability zone you choose. It is a quick way to get a working VPC when you do not need to pick the subnet or throughput yourself:

```bash
ans ecloud vpc deploydefaults "$VPC_ID" --availability-zone "$AZ_ID" --wait

RTR_ID=$(ans ecloud router list --filter vpc_id="$VPC_ID" -o value --property id)
ans ecloud network list --filter router_id="$RTR_ID" -o table --property id,name,subnet,sync.status
```

```bash
curl -s "${H[@]}" -X POST "$ECLOUD/vpcs/vpc-abcdef12/deploy-defaults" \
  -d '{"availability_zone_id": "az-abcdef12"}' | jq
```

The API response includes a `task_id` you can poll. In CLI v1.32.0 and earlier, this command does not send the availability zone and always fails validation; use the curl request or create the router and network yourself.

## Example 2: create an instance

### With the CLI

```bash
INSTANCE_ID=$(ans ecloud instance create \
  --vpc "$VPC_ID" \
  --network "$NET_ID" \
  --image img-abcdef12 \
  --name web-01 \
  --vcpu-sockets 1 --vcpu-cores-per-socket 2 \
  --ram 4096 \
  --volume 40 \
  --ssh-key-pair ssh-abcdef12 \
  --ip-address 10.0.0.10 \
  --wait -o value --property id)
```

`--vpc`, `--network`, `--image`, `--ram` and `--volume` are required. The rest are optional:

| Flag | Purpose |
|------|---------|
| `--name` | Display name |
| `--vcpu-sockets`, `--vcpu-cores-per-socket` | CPU layout, default 1 × 1 |
| `--ssh-key-pair` | SSH key pair ID. Repeat for more than one. Applied to Linux images only. |
| `--ip-address` | Fixed private IP. Must be inside the subnet, unused and not reserved. Without it, the next free address is used. |
| `--resource-tier` | Resource tier ID, otherwise the availability zone default |
| `--host-group` | Dedicated host group ID. Cannot be combined with `--resource-tier`. |
| `--enable-vm-backups` | VM-level backups. See [Example 3](#example-3-enable-backups). |
| `--backup-gateway-id` | Attach a Backup Gateway and install the backup agent. See [Example 3](#example-3-enable-backups). |
| `--enable-monitoring`, `--monitoring-gateway-id` | Monitoring. Without a monitoring gateway, the instance needs a floating IP. |
| `--tag` | Tag ID. Repeat for more than one. |
| `--wait` | Wait until the instance has finished deploying |

`--wait` on `instance create` follows the instance's `sync.status` rather than a single task, because deployment involves several steps. Allow plenty of time; Windows images take noticeably longer than Linux.

A VPC can hold up to 200 instances by default.

### With curl

```bash
curl -s "${H[@]}" -X POST "$ECLOUD/instances" -d '{
  "name": "web-01",
  "vpc_id": "vpc-abcdef12",
  "network_id": "net-abcdef12",
  "image_id": "img-abcdef12",
  "vcpu_sockets": 1,
  "vcpu_cores_per_socket": 2,
  "ram_capacity": 4096,
  "volume_capacity": 40,
  "ssh_key_pair_ids": ["ssh-abcdef12"],
  "ip_address": "10.0.0.10"
}' | jq
```

The API also accepts `vcpu_cores` on its own instead of `vcpu_sockets` and `vcpu_cores_per_socket`, but not both. For private images, leave out `volume_capacity`; the volume size comes from the image.

### Passing image parameters and a user script

Some options are available through the API but not yet as CLI flags:

| Field | Purpose |
|-------|---------|
| `image_data` | Values for the image parameters listed by `ans ecloud image parameter list`. Keys are the parameter `key`. Values must match the parameter type. |
| `user_script` | A script that runs on first boot, on Linux and Windows images |
| `requires_floating_ip` | Allocate a new floating IP and assign it to the instance |
| `floating_ip_id` | Assign an existing, unassigned floating IP. Cannot be combined with `requires_floating_ip`. |

For example, to deploy an image that takes a database root password, run a first-boot script, and get a public IP in one request:

```bash
jq -n --arg script "$(cat bootstrap.sh)" '{
  name: "db-01",
  vpc_id: "vpc-abcdef12",
  network_id: "net-abcdef12",
  image_id: "img-abcdef12",
  ram_capacity: 4096,
  volume_capacity: 40,
  vcpu_cores: 2,
  image_data: { mysql_root_password: "use-a-secret-store" },
  user_script: $script,
  requires_floating_ip: true
}' | curl -s "${H[@]}" -X POST "$ECLOUD/instances" -d @- | jq
```

Password parameters are stored as instance credentials, so you can read them back later with `ans ecloud instance credential list`.

### Get the login credentials

Credentials become available once the instance has finished deploying. Before then the API returns `404`.

```bash
ans ecloud instance credential list "$INSTANCE_ID"
ans ecloud instance credential list "$INSTANCE_ID" -o value --property username,password
```

`curl -s "${H[@]}" "$ECLOUD/instances/i-abcdef12/credentials" | jq '.data[] | {name, username, password}'`

The default user is `root` on Linux and `graphite.rack` on Windows. Each credential also has a `port`: ANS Linux images listen for SSH on port 2020, not 22, and Windows images use RDP on 3389. Open the right port in your firewall rules.

### Give the instance a public IP

Skip this if you created the instance with `requires_floating_ip`.

Floating IPs are assigned to the instance's private IP address, not to the instance itself. First find the address. Each instance has a NIC, and each NIC has an IP address of type `dhcp`:

```bash
NIC_ID=$(ans ecloud instance nic list "$INSTANCE_ID" -o value --property id)
IP_ID=$(ans ecloud nic ipaddress list "$NIC_ID" --filter type=dhcp -o value --property id)
```

Then create a floating IP in the same VPC and availability zone, and assign it:

```bash
FIP_ID=$(ans ecloud floatingip create --vpc "$VPC_ID" --availability-zone "$AZ_ID" \
  --name web-01-public --wait -o value --property id)

ans ecloud floatingip assign "$FIP_ID" --resource "$IP_ID" --wait

ans ecloud floatingip show "$FIP_ID" -o value --property ip_address
```

```bash
curl -s "${H[@]}" -X POST "$ECLOUD/floating-ips/fip-abcdef12/assign" \
  -d '{"resource_id": "ip-abcdef12"}' | jq
```

Passing an instance ID (`i-`) to `--resource` is rejected, even though the help text in CLI v1.32.0 and earlier shows it. A NIC ID (`nic-`) is still accepted but is deprecated; use the IP address ID.

`ans ecloud floatingip unassign fip-abcdef12 --wait` releases it again.

### Connect

Inbound traffic is blocked until you allow it in a firewall policy (see [Example 4](#example-4-manage-firewall-rules)). Once SSH is open, `ans ecloud instance ssh` looks up the floating IP and connects as `root` on port 2020. Use `--auto` to log in with the stored password (requires `sshpass`), or `--internal` to connect to the private IP from inside the VPC:

```bash
ans ecloud instance ssh "$INSTANCE_ID"
```

### Day 2 operations

```bash
ans ecloud instance update "$INSTANCE_ID" --ram 8192 --wait
ans ecloud instance update "$INSTANCE_ID" --vcpu-sockets 1 --vcpu-cores-per-socket 4 --wait
ans ecloud instance stop "$INSTANCE_ID"           # add --force for a hard power off
ans ecloud instance start "$INSTANCE_ID"
ans ecloud instance restart "$INSTANCE_ID"
ans ecloud instance lock "$INSTANCE_ID"
ans ecloud instance unlock "$INSTANCE_ID"
```

A locked instance cannot be updated, deleted, powered on or off, or have volumes attached or detached; the API returns `403` until you unlock it. Locking is a good guard for production instances that automation should not touch.

## Example 3: enable backups

eCloud VPC offers two kinds of backup, and an instance can use either or both.

VM-level backup protects the whole virtual machine. Your account must be set up for Secure Backup in the instance's availability zone, otherwise the request fails with a validation error.

Agent-level backup uses a Backup Gateway on your router and an agent installed inside the instance. This lets you restore individual files.

### VM-level backup

Enable it when you create the instance:

```bash
ans ecloud instance create ... --enable-vm-backups --wait
```

In the API this is `"backup_vm_enabled": true` on create or `PATCH /instances/{id}`.

### Agent-level backup with a Backup Gateway

#### 1. Create the Backup Gateway

The gateway must be on the same router (and so in the same VPC and availability zone) as the instances it protects. It does not need a floating IP. The availability zone has to support Secure Backup, and the specification must be one offered there.

```bash
BGW_SPEC=$(ans ecloud backupgateway spec list -o value --property id | head -n1)   # or choose explicitly

BGW_ID=$(ans ecloud backupgateway create --vpc "$VPC_ID" --router "$RTR_ID" \
  --specification "$BGW_SPEC" --name demo-backups --wait -o value --property id)
```

```bash
curl -s "${H[@]}" -X POST "$ECLOUD/backup-gateways" -d '{
  "name": "demo-backups",
  "vpc_id": "vpc-abcdef12",
  "router_id": "rtr-abcdef12",
  "gateway_spec_id": "bgws-abcdef12"
}' | jq
```

#### 2. Attach instances and install the agent

Attaching an instance to the gateway and installing the agent are separate settings: `backup_gateway_id` and `backup_agent_enabled`. You need both.

The CLI sets both when you pass `--backup-gateway-id`. For a new instance:

```bash
ans ecloud instance create --vpc "$VPC_ID" --network "$NET_ID" --image img-abcdef12 \
  --name app-01 --ram 4096 --volume 40 \
  --backup-gateway-id "$BGW_ID" --wait
```

For an existing instance:

```bash
ans ecloud instance update "$INSTANCE_ID" --backup-gateway-id "$BGW_ID" --wait
```

If the instance already has the gateway set but the agent is not enabled, use `--enable-agent-backups` on its own:

```bash
ans ecloud instance update "$INSTANCE_ID" --enable-agent-backups --wait
```

You cannot change an instance's Backup Gateway once it is set.

CLI v1.32.0 and earlier set only `backup_gateway_id`, which attaches the gateway without installing the agent, so no backups run. Instances created that way need `--enable-agent-backups` (or the PATCH request below). If you are on an older CLI, use the API instead.

With curl, include both fields in the create request:

```bash
curl -s "${H[@]}" -X POST "$ECLOUD/instances" -d '{
  "name": "app-01",
  "vpc_id": "vpc-abcdef12",
  "network_id": "net-abcdef12",
  "image_id": "img-abcdef12",
  "vcpu_cores": 2,
  "ram_capacity": 4096,
  "volume_capacity": 40,
  "backup_gateway_id": "bgw-abcdef12",
  "backup_agent_enabled": true
}' | jq
```

Or patch an existing instance:

```bash
curl -s "${H[@]}" -X PATCH "$ECLOUD/instances/i-abcdef12" \
  -d '{"backup_gateway_id": "bgw-abcdef12", "backup_agent_enabled": true}' | jq
```

If the instance already has the gateway set, send `{"backup_agent_enabled": true}` alone.

#### 3. Wait for the agent

The agent is installed after the instance has finished its own deployment or update, as a separate task named `instance_enable_backup`:

```bash
ans ecloud instance task list "$INSTANCE_ID" --filter name=instance_enable_backup
```

The agent install fails if the Backup Gateway is not ready, so create the gateway with `--wait` first. Instances that use a volume group are not supported for agent installation.

#### 4. Check the result

```bash
ans ecloud instance show "$INSTANCE_ID" --with-backup
```

```text
+------------+---------+-------------------+----------------------+------------------------+-------------------+
|     ID     |  NAME   | BACKUP VM ENABLED | BACKUP AGENT ENABLED | BACKUP AGENT INSTALLED | BACKUP GATEWAY ID |
+------------+---------+-------------------+----------------------+------------------------+-------------------+
| i-abcdef12 | app-01  | false             | true                 | true                   | bgw-abcdef12      |
+------------+---------+-------------------+----------------------+------------------------+-------------------+
```

`backup_agent_enabled` shows what you asked for. `backup_agent_installed` turns `true` once the agent is installed and registered. In the API response the second field is called `secure_backup`.

You cannot delete a Backup Gateway while instances are attached to it (the API returns `412`), and you cannot delete a router that still has a Backup Gateway.

## Example 4: manage firewall rules

### How the firewall is organised

Each router has firewall policies. Each policy contains rules, and each rule can contain ports.

Policies are evaluated in ascending `sequence` order, and within a policy rules are also evaluated by `sequence`, lowest first. The first rule that matches decides what happens to the traffic.

ANS creates and manages some policies on each router, such as System or Management Firewall Policy, at sequence 0. They are read-only; editing them or their rules returns `403`. Give your own policies higher sequence numbers.

A rule has:

| Field | CLI flag | Values |
|-------|----------|--------|
| Direction | `--direction` | `IN`, `OUT` or `IN_OUT` |
| Action | `--action` | `ALLOW`, `DROP` or `REJECT` |
| Source, destination | `--source`, `--destination` | `ANY`, or a comma-separated list of IPv4 addresses, CIDR ranges or address ranges (`10.0.0.1-10.0.0.50`), up to 128 entries |
| Sequence | `--sequence` | Evaluation order within the policy |
| Enabled | `--enabled` | Whether the rule is active |

A port entry has a protocol (`TCP`, `UDP` or `ICMPv4`), plus source and destination ports for TCP and UDP. Ports can be a single port (`443`), a range (`8000-8100`), a comma-separated list (`80,443`) or `ANY`, with port numbers from 1 to 65535.

A rule with no ports matches all protocols. Add at least one port entry to restrict it.

`ICMPv4` entries take no ports and allow echo request (ping) only. You cannot select other ICMP types or codes.

Use IPv4 addresses only.

### Create a policy

```bash
FWP_ID=$(ans ecloud firewallpolicy create --router "$RTR_ID" --name app --sequence 10 \
  --wait -o value --property id)
```

```bash
curl -s "${H[@]}" -X POST "$ECLOUD/firewall-policies" \
  -d '{"name": "app", "router_id": "rtr-abcdef12", "sequence": 10}' | jq
```

### Allow inbound TCP

Allow SSH (port 2020 on ANS Linux images) from an office range, and HTTP/HTTPS from anywhere, to `web-01`:

```bash
FWR_SSH=$(ans ecloud firewallrule create --policy "$FWP_ID" --name ssh-from-office \
  --sequence 10 --direction IN --action ALLOW \
  --source 203.0.113.0/24 --destination 10.0.0.10 \
  --enabled --wait -o value --property id)

ans ecloud firewallrule port create --rule "$FWR_SSH" --protocol TCP \
  --source ANY --destination 2020 --name ssh --wait

FWR_WEB=$(ans ecloud firewallrule create --policy "$FWP_ID" --name web \
  --sequence 20 --direction IN --action ALLOW \
  --source ANY --destination 10.0.0.10 \
  --enabled --wait -o value --property id)

ans ecloud firewallrule port create --rule "$FWR_WEB" --protocol TCP \
  --source ANY --destination 80,443 --name http-https --wait
```

`ans ecloud firewallrule create` creates the rule disabled unless you pass `--enabled`. Always include the flag, or enable the rule later with `ans ecloud firewallrule update fwr-abcdef12 --enabled`.

The CLI adds ports one at a time. The API can create a rule and its ports in one request, which saves a round of waiting:

```bash
curl -s "${H[@]}" -X POST "$ECLOUD/firewall-rules" -d '{
  "name": "web",
  "firewall_policy_id": "fwp-abcdef12",
  "sequence": 20,
  "direction": "IN",
  "action": "ALLOW",
  "source": "ANY",
  "destination": "10.0.0.10",
  "enabled": true,
  "ports": [
    { "protocol": "TCP", "source": "ANY", "destination": "80,443" }
  ]
}' | jq
```

### Allow UDP

Allow outbound DNS to two resolvers from the whole network:

```bash
FWR_DNS=$(ans ecloud firewallrule create --policy "$FWP_ID" --name dns-out \
  --sequence 30 --direction OUT --action ALLOW \
  --source 10.0.0.0/24 --destination 198.51.100.53,198.51.100.54 \
  --enabled --wait -o value --property id)

ans ecloud firewallrule port create --rule "$FWR_DNS" --protocol UDP \
  --source ANY --destination 53 --wait
ans ecloud firewallrule port create --rule "$FWR_DNS" --protocol TCP \
  --source ANY --destination 53 --wait
```

### Allow ICMP

Allow ping from your monitoring hosts:

```bash
FWR_PING=$(ans ecloud firewallrule create --policy "$FWP_ID" --name ping-in \
  --sequence 40 --direction IN --action ALLOW \
  --source 203.0.113.10,203.0.113.11 --destination 10.0.0.0/24 \
  --enabled --wait -o value --property id)

ans ecloud firewallrule port create --rule "$FWR_PING" --protocol ICMPv4 --wait
```

```bash
curl -s "${H[@]}" -X POST "$ECLOUD/firewall-rule-ports" \
  -d '{"firewall_rule_id": "fwr-abcdef12", "protocol": "ICMPv4"}' | jq
```

### Block traffic explicitly

Policies with a higher sequence are evaluated later, so a final policy can drop anything your earlier rules did not allow:

```bash
FWP_DENY=$(ans ecloud firewallpolicy create --router "$RTR_ID" --name deny-all --sequence 100 \
  --wait -o value --property id)

ans ecloud firewallrule create --policy "$FWP_DENY" --name drop-inbound \
  --sequence 10 --direction IN --action DROP \
  --source ANY --destination ANY --enabled --wait
```

`DROP` discards packets silently. `REJECT` sends a refusal back, which makes connection failures quicker to diagnose but tells the sender that something is there.

### Change, disable and remove rules

```bash
ans ecloud firewallrule update "$FWR_SSH" --source 203.0.113.0/24,198.51.100.7 --wait
ans ecloud firewallrule update "$FWR_SSH" --enabled=false --wait     # disable
ans ecloud firewallrule update "$FWR_SSH" --enabled --wait           # enable again
ans ecloud firewallrule port update fwrp-abcdef12 --destination 2020,2222 --wait
ans ecloud firewallrule port delete fwrp-abcdef12 --wait
ans ecloud firewallrule delete "$FWR_SSH" --wait
```

`ans ecloud firewallrule update` only changes flags you pass, so `--enabled=false` is the way to switch a rule off.

```bash
curl -s "${H[@]}" -X PATCH "$ECLOUD/firewall-rules/fwr-abcdef12" -d '{"enabled": false}' | jq
curl -s "${H[@]}" -X DELETE "$ECLOUD/firewall-rules/fwr-abcdef12" | jq
```

### When the change is live

Each rule or port change makes the parent policy sync to the router. The `--wait` flag waits for that. Without it, wait for the policy:

```bash
wait_for_sync firewallpolicy "$FWP_ID"
ans ecloud firewallpolicy task list "$FWP_ID"
```

While the policy is syncing, further changes to the same policy return `409`. If you script many changes, either use `--wait` on each one or create rules with their ports inline through the API.

### Network policies

VPCs created with advanced networking can also have network policies, which filter east-west traffic between instances on the same network. They are managed with `ans ecloud networkpolicy`, `ans ecloud networkrule` and `ans ecloud networkrule port`, which work much like the router firewall commands above.

## Verifying configuration

After a run, check that the infrastructure matches what you intended. The commands below read state only and are safe to run at any time.

### Everything has finished syncing

```bash
{
  ans ecloud vpc show "$VPC_ID" -o value --property id,sync.status
  for r in router instance backupgateway floatingip; do
    ans ecloud "$r" list --filter vpc_id="$VPC_ID" -o value --property id,sync.status
  done
  for rtr in $(ans ecloud router list --filter vpc_id="$VPC_ID" -o value --property id); do
    ans ecloud network list --filter router_id="$rtr" -o value --property id,sync.status
    ans ecloud firewallpolicy list --router "$rtr" -o value --property id,sync.status
  done
} | awk '$2 != "complete" { print "not ready:", $0; bad=1 } END { exit bad }'
```

Networks and firewall policies belong to a router rather than directly to a VPC, so the script lists them per router.

With curl:

```bash
curl -s "${H[@]}" "$ECLOUD/instances?vpc_id=vpc-abcdef12&per_page=100" \
  | jq -r '.data[] | [.id, .name, .sync.status] | @tsv'
```

API list responses are paginated. Use `per_page` (up to 100) and `page`, and check `.meta.pagination` to see whether there are more pages.

### Instance state

```bash
ans ecloud instance show "$INSTANCE_ID" --with-all
ans ecloud instance show "$INSTANCE_ID" -o json | jq '.[0] | {id, online, agent_running, sync, vcpu_cores, ram_capacity, volume_capacity}'
ans ecloud instance nic list "$INSTANCE_ID"
ans ecloud instance floatingip list "$INSTANCE_ID"
```

`--with-backup`, `--with-monitoring` and `--with-tags` add those groups of columns to the table; `--with-all` adds all of them.

### Backups

```bash
ans ecloud instance show "$INSTANCE_ID" -o jsonpath='{[0].secure_backup}'
ans ecloud instance show "$INSTANCE_ID" -o json | jq -e '.[0].secure_backup == true'
```

JSON and jsonpath output use the API field name `secure_backup`; only the `--with-backup` table calls it `backup_agent_installed`. The second form exits non-zero if the agent is not installed, which makes it useful as a pipeline gate.

### Firewall

```bash
ans ecloud router firewallpolicy list "$RTR_ID" --sort sequence
ans ecloud firewallpolicy rule list "$FWP_ID" --sort sequence
ans ecloud firewallrule port list "$FWR_WEB"
```

To assert that a specific rule exists and is enabled:

```bash
ans ecloud firewallrule list --policy "$FWP_ID" --filter name=web -o json \
  | jq -e '.[0] | .enabled == true and .action == "ALLOW" and .direction == "IN"'
```

To dump a whole policy, including ports, for review or for diffing against a known good copy:

```bash
for rule in $(ans ecloud firewallpolicy rule list "$FWP_ID" -o value --property id); do
  ans ecloud firewallrule show "$rule" -o json | jq '.[0] | {name, sequence, direction, action, source, destination, enabled}'
  ans ecloud firewallrule port list "$rule" -o json | jq '[.[] | {protocol, source, destination}]'
done
```

### Recent failures

```bash
ans ecloud task list --filter status=failed --sort created_at:desc
```

## Cleaning up

Delete resources in reverse order of creation. Deletes are asynchronous too. A resource stays visible with `sync.type` set to `delete` until the work finishes, after which `show` returns not found.

```bash
ans ecloud floatingip unassign "$FIP_ID" --wait
ans ecloud floatingip delete "$FIP_ID" --wait
ans ecloud instance unlock "$INSTANCE_ID"           # if locked
ans ecloud instance delete "$INSTANCE_ID" --wait
ans ecloud backupgateway delete "$BGW_ID" --wait
ans ecloud firewallpolicy delete "$FWP_ID" --wait
ans ecloud network delete "$NET_ID" --wait
ans ecloud router delete "$RTR_ID" --wait
ans ecloud vpc delete "$VPC_ID" --wait
```

`vpc delete` can also remove everything inside a VPC. Run it with `--dry-run` first to see what it would delete:

```bash
ans ecloud vpc delete "$VPC_ID" --recursive --dry-run
ans ecloud vpc delete "$VPC_ID" --recursive --force
```

## Troubleshooting

| Symptom | Cause | Fix |
|---------|-------|-----|
| `409` "The specified resource has a task in progress" | The resource, or the firewall policy that owns the rule or port, is still busy | Wait for `sync.status` to reach `complete`, then retry |
| `422` on `vpc_id`, `router_id`, `network_id` or similar | The parent's `sync.status` is not `complete` | Wait for the parent. If its status is `failed`, contact support. |
| `403` when changing or deleting an instance | The instance is locked | `ans ecloud instance unlock i-abcdef12` |
| `403` "not editable" on a firewall policy or rule | It belongs to a policy managed by ANS | Put your rules in your own policy |
| `412` when deleting a Backup Gateway | Instances are still attached | Delete or detach those instances first |
| Router delete refused | The router still has a Backup Gateway or Monitoring Gateway | Delete the gateway first |
| Floating IP assign rejected | The resource ID is an instance (`i-`) or in another availability zone | Use the instance's `dhcp` IP address ID (`ip-`) in the same zone |
| Firewall rule exists but has no effect | It was created without `--enabled` | `ans ecloud firewallrule update fwr-abcdef12 --enabled --wait` |
| Firewall rule lets through more than expected | The rule has no ports, so it matches every protocol | Add port entries |
| Credentials return `404` | The instance has not finished deploying | Wait for `sync.status` to reach `complete` |
| Backup Gateway set but `backup_agent_installed` stays `false` | `backup_agent_enabled` was never set (for example, the instance was created with CLI v1.32.0 or earlier), or the agent install task failed | `ans ecloud instance update i-abcdef12 --enable-agent-backups --wait`, then check `ans ecloud instance task list` |
| `ans ecloud vpc deploydefaults` fails validation | CLI v1.32.0 and earlier do not send an availability zone | Upgrade the CLI, or create the router and network yourself |
| `--wait` times out on large instances | The default timeout is 20 minutes | Raise `command_wait_timeout_seconds` |
| `ip_address` rejected | The address is outside the subnet, reserved, or in use | Choose another address, or leave it out |
| SSH key not applied | The image is not Linux | Use the instance credentials instead |

In CLI v1.32.0 and earlier, the sync status column is missing from the default table output. Ask for the property directly with `--property sync.status` or use JSON output.

## Complete example script

This script builds a VPC with one Linux web server, opens SSH, HTTP/HTTPS and ping, gives the server a public IP, protects it with an agent-level backup, and checks the result. Fill in the variables at the top before running it.

```bash
#!/usr/bin/env bash
set -euo pipefail

: "${ANS_API_KEY:?set ANS_API_KEY}"
ECLOUD="https://api.ukfast.io/ecloud/v2"
H=(-H "Authorization: ${ANS_API_KEY}" -H "Content-Type: application/json" -H "Accept: application/json")

REGION_ID=reg-abcdef12
AZ_ID=az-abcdef12
IMAGE_ID=img-abcdef12
SSH_KEY_ID=ssh-abcdef12
ADMIN_CIDR=203.0.113.0/24
NAME=demo

wait_for_sync() {
  local resource=$1 id=$2 timeout=${3:-1800} waited=0 status
  while (( waited < timeout )); do
    status=$(ans ecloud "$resource" show "$id" -o value --property sync.status)
    case "$status" in
      complete) return 0 ;;
      failed)   echo "$resource $id failed" >&2; return 1 ;;
    esac
    sleep 10; waited=$(( waited + 10 ))
  done
  echo "timed out waiting for $resource $id" >&2; return 1
}

echo "Creating VPC, router and network"
VPC_ID=$(ans ecloud vpc create --region "$REGION_ID" --name "$NAME" --wait -o value --property id)
RTR_ID=$(ans ecloud router create --vpc "$VPC_ID" --availability-zone "$AZ_ID" \
  --name "$NAME-rtr" --wait -o value --property id)
NET_ID=$(ans ecloud network create --router "$RTR_ID" --subnet 10.0.0.0/24 \
  --name "$NAME-net" --wait -o value --property id)

echo "Adding default firewall policies"
ans ecloud router deploydefaults "$RTR_ID"

echo "Creating Backup Gateway"
BGW_SPEC=$(ans ecloud backupgateway spec list -o value --property id | head -n1)
BGW_ID=$(ans ecloud backupgateway create --vpc "$VPC_ID" --router "$RTR_ID" \
  --specification "$BGW_SPEC" --name "$NAME-bgw" --wait -o value --property id)

echo "Creating instance"
INSTANCE_ID=$(jq -n \
  --arg vpc "$VPC_ID" --arg net "$NET_ID" --arg img "$IMAGE_ID" \
  --arg ssh "$SSH_KEY_ID" --arg bgw "$BGW_ID" --arg name "$NAME-web-01" '{
    name: $name, vpc_id: $vpc, network_id: $net, image_id: $img,
    vcpu_cores: 2, ram_capacity: 4096, volume_capacity: 40,
    ip_address: "10.0.0.10", ssh_key_pair_ids: [$ssh],
    backup_gateway_id: $bgw, backup_agent_enabled: true
  }' | curl -sf "${H[@]}" -X POST "$ECLOUD/instances" -d @- | jq -r '.data.id')
wait_for_sync instance "$INSTANCE_ID" 3600

echo "Assigning floating IP"
NIC_ID=$(ans ecloud instance nic list "$INSTANCE_ID" -o value --property id)
IP_ID=$(ans ecloud nic ipaddress list "$NIC_ID" --filter type=dhcp -o value --property id)
FIP_ID=$(ans ecloud floatingip create --vpc "$VPC_ID" --availability-zone "$AZ_ID" \
  --name "$NAME-web-01" --wait -o value --property id)
ans ecloud floatingip assign "$FIP_ID" --resource "$IP_ID" --wait

echo "Creating firewall rules"
FWP_ID=$(ans ecloud firewallpolicy create --router "$RTR_ID" --name "$NAME-app" \
  --sequence 10 --wait -o value --property id)

create_rule() {
  # name sequence direction source destination protocol port
  local body
  body=$(jq -n --arg policy "$FWP_ID" --arg name "$1" --argjson seq "$2" \
    --arg dir "$3" --arg src "$4" --arg dst "$5" --arg proto "$6" --arg port "${7:-}" '{
      firewall_policy_id: $policy, name: $name, sequence: $seq,
      direction: $dir, action: "ALLOW", source: $src, destination: $dst, enabled: true,
      ports: [ if $proto == "ICMPv4" then {protocol: $proto}
               else {protocol: $proto, source: "ANY", destination: $port} end ]
    }')
  curl -sf "${H[@]}" -X POST "$ECLOUD/firewall-rules" -d "$body" > /dev/null
  wait_for_sync firewallpolicy "$FWP_ID"
}

create_rule ssh  10 IN "$ADMIN_CIDR" 10.0.0.10 TCP 2020
create_rule web  20 IN ANY           10.0.0.10 TCP 80,443
create_rule ping 30 IN "$ADMIN_CIDR" 10.0.0.10 ICMPv4

echo "Waiting for backup agent"
for _ in $(seq 1 60); do
  installed=$(ans ecloud instance show "$INSTANCE_ID" -o jsonpath='{[0].secure_backup}')
  [[ "$installed" == "true" ]] && break
  sleep 30
done
[[ "$installed" == "true" ]] || { echo "backup agent not installed" >&2; exit 1; }

echo "Verifying"
ans ecloud instance show "$INSTANCE_ID" --with-backup
ans ecloud firewallpolicy rule list "$FWP_ID" --sort sequence
ans ecloud floatingip show "$FIP_ID" -o list --property ip_address,resource_id
```
