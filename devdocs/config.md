# Configuration and State

NimbusDB uses a multi-layered configuration and state management system. This document describes all levels of configuration, their purpose, and every parameter available at each level.

## Overview

NimbusDB has three distinct levels of configuration and state:

1. **Cluster Configuration** - Shared across all nodes in the cluster
2. **Program Arguments** - Node-specific startup parameters
3. **Runtime State** - Dynamic node state that changes during execution

The configuration hierarchy works as follows:

- Cluster Configuration is loaded first (from YAML file or environment variables)
- Program Arguments are parsed from command-line flags
- Runtime State is managed dynamically during node execution

---

## 1. Cluster Configuration

Cluster Configuration (`Config`) is shared across all nodes in a NimbusDB cluster. It defines cluster-wide settings such as shard counts, NATS messaging configuration, and blob storage settings.

### Configuration Sources

Cluster Configuration can be loaded from:

1. **YAML Configuration File** (default: `.config.yml`)
2. **Environment Variables** (overrides YAML values)

Environment variables take precedence over YAML file values, allowing for flexible deployment configurations.

### Configuration Structure

The `Config` struct contains the following sections:

- Root-level cluster settings
- Blob storage configuration (`BlobConfig`)
- NATS messaging configuration (`NATSConfig`)
- Database configuration (`DbConfig`)

### Root-Level Configuration Parameters

| Parameter    | Type     | Environment Variable | YAML Key     | Default | Description                                  | Constraints                                                                           |
| ------------ | -------- | -------------------- | ------------ | ------- | -------------------------------------------- | ------------------------------------------------------------------------------------- |
| `ShardCount` | `uint16` | `SHARD_COUNT`        | `shardCount` | `16`    | Total number of shards in the cluster        | Must be between 1 and 256 (inclusive). **Should be more than total nodes in cluster** |
| `HealthPort` | `int`    | `HEALTH_PORT`        | `healthPort` | `8080`  | Port number for the health check HTTP server | Must be between 1 and 65535 (inclusive)                                               |
| `LogLevel`   | `string` | `LOG_LEVEL`          | `logLevel`   | `info`  | Logging verbosity level                      | Must be one of: `trace`, `debug`, `info`, `warn`, `error`, `fatal`, `panic`           |

### Blob Storage Configuration (`BlobConfig`)

The `BlobConfig` struct contains settings for MinIO blob storage integration.

| Parameter                           | Type            | Environment Variable                          | YAML Key                                 | Default | Description                                                                           | Constraints                           |
| ----------------------------------- | --------------- | --------------------------------------------- | ---------------------------------------- | ------- | ------------------------------------------------------------------------------------- | ------------------------------------- |
| `Endpoint`                          | `string`        | `BLOB_ENDPOINT`                               | `blob.endpoint`                          | -       | MinIO server endpoint URL (e.g., `localhost:9000`)                                    | Required for blob operations          |
| `AccessKeyID`                       | `string`        | `BLOB_ACCESS_KEY_ID`                          | `blob.accessKeyID`                       | -       | MinIO access key ID for authentication                                                | Required for blob operations          |
| `SecretAccessKey`                   | `string`        | `BLOB_SECRET_ACCESS_KEY`                      | `blob.secretAccessKey`                   | -       | MinIO secret access key for authentication                                            | Required for blob operations          |
| `UseSSL`                            | `bool`          | `BLOB_USE_SSL`                                | `blob.useSSL`                            | `false` | Whether to use SSL/TLS for MinIO connections                                          | Boolean (true/false)                  |
| `DeleteMarkerCleanupDelayDays`      | `int`           | `BLOB_DELETE_MARKER_CLEANUP_DELAY_DAYS`       | `blob.deleteMarkerCleanupDelayDays`      | `1`     | Number of days to wait before cleaning up delete markers in blob storage              | Must be between 1 and 365 (inclusive) |
| `NonCurrentVersionCleanupDelayDays` | `int`           | `BLOB_NON_CURRENT_VERSION_CLEANUP_DELAY_DAYS` | `blob.nonCurrentVersionCleanupDelayDays` | `1`     | Number of days to wait before cleaning up non-current object versions in blob storage | Must be between 1 and 365 (inclusive) |
| `BlobOperationTimeout`              | `time.Duration` | `BLOB_OPERATION_TIMEOUT`                      | `blob.blobOperationTimeout`              | `30s`   | Timeout for blob operations                                                           | Must be a valid duration              |

### NATS Configuration (`NATSConfig`)

The `NATSConfig` struct contains settings for NATS messaging system integration.

| Parameter          | Type            | Environment Variable  | YAML Key                | Default                 | Description                                      | Constraints                                                                                                    |
| ------------------ | --------------- | --------------------- | ----------------------- | ----------------------- | ------------------------------------------------ | -------------------------------------------------------------------------------------------------------------- |
| `URL`              | `string`        | `NATS_URL`            | `nats.url`              | `nats://localhost:4222` | NATS server connection URL                       | Must be a valid NATS URL format                                                                                |
| `Creds`            | `string`        | `NATS_CREDS`          | `nats.creds`            | -                       | Path to NATS credentials file for authentication | Optional, used for NATS authentication                                                                         |
| `SubjectPrefix`    | `string`        | `NATS_SUBJECT_PREFIX` | `nats.subjectPrefix`    | `nimbus`                | Prefix for all NATS subjects used by NimbusDB    | Must be non-empty; can contain alphanumeric characters, dots (.), underscores (\_), dashes (-), and colons (:) |
| `NatsDrainTimeout` | `time.Duration` | `NATS_DRAIN_TIMEOUT`  | `nats.natsDrainTimeout` | `30s`                   | Timeout for NATS drain operation                 | Must be a valid duration                                                                                       |

### Database Configuration (`DbConfig`)

The `DbConfig` struct contains settings for database operations.

| Parameter           | Type  | Environment Variable     | YAML Key               | Default | Description                                 | Constraints                |
| ------------------- | ----- | ------------------------ | ---------------------- | ------- | ------------------------------------------- | -------------------------- |
| `ChannelBufferSize` | `int` | `DB_CHANNEL_BUFFER_SIZE` | `db.channelBufferSize` | `256`   | Buffer size for database operation channels | Must be a positive integer |

### Example YAML Configuration

```yaml
shardCount: 16
healthPort: 8080
logLevel: info
blob:
  endpoint: localhost:9000
  accessKeyID: minioadmin
  secretAccessKey: minioadmin
  useSSL: false
  deleteMarkerCleanupDelayDays: 1
  nonCurrentVersionCleanupDelayDays: 1
  blobOperationTimeout: 30s
nats:
  url: nats://localhost:4222
  subjectPrefix: nimbus
  natsDrainTimeout: 30s
db:
  channelBufferSize: 256
```

### Configuration Loading Order

1. Base YAML file is loaded (if path is provided and file exists)
2. Environment variables override YAML values
3. Default values are applied for unset fields
4. Configuration is validated
5. Final configuration is logged (sensitive fields like credentials are not logged)

---

## 2. Program Arguments

Program Arguments (`ProgramArguments`) are node-specific command-line parameters that control the runtime behavior of each individual node. These arguments are parsed at startup and remain constant throughout the node's execution.

### Command-Line Flags

Program Arguments are specified via command-line flags. Each flag has a long form and a short form (shorthand).

### Program Arguments Parameters

| Parameter    | Type     | Long Flag   | Short Flag | Default       | Description                                 | Valid Values                               |
| ------------ | -------- | ----------- | ---------- | ------------- | ------------------------------------------- | ------------------------------------------ |
| `Mode`       | `string` | `--mode`    | `-m`       | `single`      | Operation mode of the node                  | `single`, `distributed` (case-insensitive) |
| `ConfigPath` | `string` | `--config`  | `-c`       | `.config.yml` | Path to the cluster configuration YAML file | Any valid file path                        |
| `Help`       | `bool`   | `--help`    | `-h`       | `false`       | Display help message and exit               | Boolean flag (no value)                    |
| `Version`    | `bool`   | `--version` | `-v`       | `false`       | Display version information and exit        | Boolean flag (no value)                    |

### Operation Modes

| Mode          | Description                                                                                      |
| ------------- | ------------------------------------------------------------------------------------------------ |
| `single`      | Single-node operation mode. The node operates independently without cluster coordination.        |
| `distributed` | Distributed operation mode. The node participates in a cluster and coordinates with other nodes. |

### Example Usage

```bash
# Run in single mode with default config
./nimbusdb

# Run in distributed mode with custom config
./nimbusdb --mode distributed --config /path/to/config.yml

# Using short flags
./nimbusdb -m distributed -c /path/to/config.yml

# Display help
./nimbusdb --help

# Display version
./nimbusdb --version
```

### Validation

All Program Arguments are validated:

- Mode must be one of the valid operation modes
- Validation is case-insensitive (values are normalized to lowercase)
- Invalid values result in an error and the application exits

---

## 3. Runtime State

Runtime State (`State`) represents the dynamic, node-specific state that changes during execution. Unlike configuration and program arguments, runtime state is mutable and reflects the current operational status of the node.

### State Characteristics

- **Node-Specific**: Each node maintains its own runtime state
- **Dynamic**: State changes during node execution based on cluster operations
- **Volatile**: State is not persisted and is re-established on node startup
- **Operational**: Reflects the current operational responsibilities of the node

### Runtime State Parameters

| Parameter  | Type       | Description                                            | Mutability                                               |
| ---------- | ---------- | ------------------------------------------------------ | -------------------------------------------------------- |
| `shardIDs` | `[]uint16` | List of shard IDs currently owned/managed by this node | Mutable - changes when shards are assigned or reassigned |

### State Management

The runtime state is managed internally by the node and changes based on:

- Cluster membership changes
- Shard assignment and reassignment operations
- Node role changes (if applicable)
- Cluster rebalancing operations

### State Initialization

Runtime state is typically initialized:

- As an empty state when the node starts
- Populated based on cluster coordination and shard assignment logic
- Updated dynamically as the cluster state changes

### Example State Values

```go
// Initial state (no shards assigned)
State{
    shardIDs: []uint16{}
}

// State after shard assignment
State{
    shardIDs: []uint16{0, 1, 2, 3}
}

// State after shard reassignment
State{
    shardIDs: []uint16{4, 5, 6, 7, 8}
}
```

---

## Configuration Hierarchy Summary

| Level                     | Scope         | Persistence           | Mutability                         | Source                             |
| ------------------------- | ------------- | --------------------- | ---------------------------------- | ---------------------------------- |
| **Cluster Configuration** | Cluster-wide  | Persistent (YAML/env) | Immutable (per node restart)       | YAML file or environment variables |
| **Program Arguments**     | Node-specific | Command-line          | Immutable (per execution)          | Command-line flags                 |
| **Runtime State**         | Node-specific | Volatile (in-memory)  | Mutable (changes during execution) | Managed internally by the node     |

---

## Best Practices

1. **Cluster Configuration**: Store sensitive values (credentials, keys) in a secret store and load into env variables. Use yml only in dev and secured environments.
2. **Lease Privileged**:
   - For NATS use a credential that has access to only `{subjectPrefix}.*`
   - For MinIO use an account that only has access to:
     - list
     - put
     - get
     - create bucket
     - create policy
     - apply policy to bucket
       > For MinIO you can even remove the bucket and policy related permissions if you manually create the bucket. Please see `blob/operations.go` to know what permissions are needed.
