#!/usr/bin/env python3
"""Minimal OPC UA server used as a test fixture for the client-only addon."""

import sys
import time
from opcua import Server


def main():
    port = 4841
    if len(sys.argv) > 1:
        port = int(sys.argv[1])

    server = Server()
    server.set_endpoint(f"opc.tcp://127.0.0.1:{port}")
    uri = "http://tier0.io/sourceflow-opcua/test"
    idx = server.register_namespace(uri)

    objects = server.get_objects_node()
    node_id = f"ns={idx};s=TestNode"
    var = objects.add_variable(node_id, "TestNode", 12.3)
    var.set_writable()

    server.start()
    print(f"RUNNING opc.tcp://127.0.0.1:{port} {node_id}", flush=True)

    try:
        while True:
            time.sleep(1)
    finally:
        server.stop()


if __name__ == "__main__":
    main()
