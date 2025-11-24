# NimbusDb API

Nimbus uses NATS to accept read, write and other requests

## Configuration APIs

### Get Shard Count

**Requester**: NimbusDb Client
**Responder**: Any healthy NimbusDb data node
**Description**:

- Clients are expected to be shard aware
- Meaning clients calculate shard id from partition key, using formula - hash(paritionKey) mod shardCount.
- For this shard aware behaviour clients need to know shard count.

> Clients are expected to cache this value to avoid repeated requests.

```bash
nats req nimbus.config.getShardCount
```

**Server side implementation Notes**

- All available data nodes are listening for request on this subject line through a common queue group called "common_config_qg".
- Data nodes use channel

### 0. Save an object (point write)

**Requester**: NimbusDb Client
**Responder**: Shard Owner.
**Description**:

- This is a point write: meaning only object lives under this path ever. No metadata involvement at all.
  - Note that this doesn't mean the object can't be an array. It just means, its a single file on object store.
- Shard owner is the data node that is assigned by raft metadata cluster to be owner of the shard in question.
- Shard owner has a NATS subject `nimbus.shards.{shardId}.op` subscribed on.
- Client sends the point write request to this subject
  - Headers identify file path and operation type
  - Data is just byte[] -> typically MsgPack value of object(s) getting stored
- Shard owner writes to blob
- Shard owner never parses msg body it just directly writes byte[] to blob.
- Upon success, shard owner returns json with { "errors": "", status: 200 } . If there are errors, error (string) is returned along with appropriate status.

**Example Requests**

```bash
nats req \
  -H "type: 0" \
  -H "bucketName: gk-test" \
  -H "fileName: /ts-id-2/p" \
  -H "overwrite: true" \
  nimbus.shards.12.op \
  "some-random-data"
```

> Please note: In reality the nats message is byte[] corrosponding to msgpack data that client will send.

**Server side implementation Notes**

- Data node starts channel subscriber for all its shard ids.
- Channel is needed here as blob operations are much more expensive compared to nats and not using channel would lead to goroutine explosion.
- Channel provides backpressure in case too many requests start coming in.
- one subscription per shard
- one channel per subscription
- one goroutine to dequeue from channel and process write operation.

### 1. Read an object (point read)

**Requester**: NimbusDb Client
**Responder**: Shard Owner.
**Description**:

- This is a point read: meaning only one object lives under this path ever. No metadata involvement.
  - Note that this doesn't mean the object can't be an array. It just means, its a single file on object store.
- Shard owner is the data node that is assigned by raft metadata cluster to be owner of the shard in question.
- Shard owner has a NATS subject `nimbus.shards.{shardId}.op` subscribed on.
- Client sends the point read request to this subject
  - Headers identify file path and operation type
  - The Data returned is just byte[] -> typically MsgPack value of object(s) being returned
- Shard owner reads from db at exact path. Nothing fancy.
- Shard owner never parses data it just directly writes byte[] back to nats response.
- **Example Requests**

```bash
nats req \
  -H "type: 1" \
  -H "bucketName: gk-test" \
  -H "fileName: /ts-id-2/p" \
  nimbus.shards.12.op
```

**Server side implementation Notes**

- Channel subscription and message processing techniques are same as point write or any other data operation.
