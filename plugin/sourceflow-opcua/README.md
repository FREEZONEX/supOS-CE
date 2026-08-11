# node-red-contrib-opcua-open62541

OPC UA client nodes for Node-RED/sourceflow powered by open62541, with shared client connection support.

## Design

- `opcua-connection` is a config node that maintains a single shared open62541 client session.
- `opcua-read`, `opcua-write`, `opcua-subscribe` operation nodes reference a connection.
- Multiple nodes and flows can share the same connection, reducing OPC UA Server session count.
- Server functionality is intentionally **not** included in this package.

## Nodes

| Node | Type | Description |
|---|---|---|
| `opcua-connection` | Config | Shared open62541 client configuration |
| `opcua-read` | Operation | Read a variable by NodeId |
| `opcua-write` | Operation | Write a variable by NodeId |
| `opcua-subscribe` | Operation | Subscribe to variable changes |

## Build

Requires open62541 headers and library. By default `binding.gyp` looks in `/usr/local`. For local development inside this repository, the sibling `plugin/opcua/.build/open62541-install` directory is detected automatically.

```bash
cd plugin/sourceflow-opcua
npm install
npm run build
```

To force the bundled open62541 location:

```bash
npm run build:dev
```

The native addon must be built before Node-RED loads the package. If the addon is missing, nodes still appear in the palette but operations fail with an error telling you to run `npm run build`.

## Integration

Add to `deploy/mount/sourceflow/packages.txt`:

```text
/data/offline_modules/node-red-contrib-opcua-open62541
```

Run `deploy/bin/hide-nodered.sh` after deployment. When the new package is installed, legacy `node-red-contrib-opcua` nodes are hidden to avoid palette duplication.

## Example Flow (JSON)

A minimal flow that connects to an existing OPC UA server and reads a variable:

```json
[
  {
    "id": "conn",
    "type": "opcua-connection",
    "endpoint": "opc.tcp://localhost:4840",
    "securityMode": "None",
    "securityPolicy": "None"
  },
  {
    "id": "read",
    "type": "opcua-read",
    "connection": "conn",
    "nodeId": "ns=2;s=Temperature",
    "wires": [["debug"]]
  },
  {
    "id": "debug",
    "type": "debug"
  }
]
```

## Message Conventions

### Input

| Field | Description |
|---|---|
| `msg.topic` | NodeId (overrides static config) |
| `msg.nodeId` | NodeId fallback |
| `msg.payload` | Value to write |
| `msg.dataType` | OPC UA data type, e.g. `String`, `Int32`, `Double`, `Boolean` |

### Output

| Field | Description |
|---|---|
| `msg.payload` | Read value or write result |
| `msg.topic` | NodeId that was operated on |
| `msg.dataType` | OPC UA data type |
| `msg.statusCode` | OPC UA StatusCode string, e.g. `Good` |
| `msg.sourceTimestamp` | Source timestamp (reads/subscriptions) |
| `msg.serverTimestamp` | Server timestamp (reads/subscriptions) |

## Security

`None` security mode is fully supported. `Sign` and `SignAndEncrypt` modes are accepted but certificate loading is a placeholder in this version; production use with certificates requires additional PKI configuration in `src/addon.cc`.

## Testing

Standalone integration tests exercise the addon directly (no Node-RED runtime required):

```bash
npm test
```

The test starts a Python OPC UA fixture server, connects a client, reads/writes/subscribes, then shuts down. If the native addon is not built, the test prints instructions and exits with code 0 so CI does not break.

## License

Apache-2.0
