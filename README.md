# NimbusDb

A lean, high-performance distributed database that uses object storage (MinIO, S3, Azure Blob Storage) as its storage layer. Built for massive scale at minimal cost.

> Azure blob storage support is upcoming!

## Advantages

- **Cost-efficient at scale**: Leverage object storage economics to store petabytes at a fraction of traditional database costs
- **Inherited durability**: Built-in fault tolerance, redundancy, and geographic distribution from object storage providers
- **Zero infrastructure overhead**: No sharding, replication, or cluster management required—object storage handles it all
- **Global distribution**: Deploy instances worldwide with NATS providing seamless, location-transparent connectivity

## Functionality

### Lightweight & Performant

Built from scratch in Go with minimal resource footprint. Uses NATS as the communication layer for effortless global deployment—hundreds of instances, zero load balancers, no DNS complexity.

### Optimized Access Patterns

Nimbus excels for workloads matching these patterns:

1. **Point lookups by ID**: Fast, direct retrieval of individual objects. e.g. get object with id=some_uuid
2. **Temporal streaming**: Stream objects from collections in insertion timestamp order. e.g. get all activities performed by a user.

**Ideal use cases**: Event sourcing, audit trails, time-series logs, document stores, and any workload where data is append-heavy with lookups by a single Id.

If you have a workload that can fit into these criteria, you can benefit by using Nimbus to achieve massive scale at extremely low cost.

## API

See [API Doc](/devdocs/api.md)

## Development

- For running the project locally, see [Local Setup](/devdocs/local_setup.md)
- If you are interested in contributing, coding standards and style guidelines, see the [Contribution Guide](/devdocs/style_guide.md).
- For Nimbus architecture, please see [Architecture & Dev Docs](/devdocs/README.md)

## Credits

🏗️ Architected by Gaurav Kalele, <br />
💫 Vibe coded with Cursor <br />
👮 Thoroughly Reviewd by: Gaurav (Human), Cursor (self review), Copilot (PR Review)

© Connection Loops
