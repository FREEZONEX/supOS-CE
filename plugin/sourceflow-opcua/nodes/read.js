'use strict';

module.exports = function (RED) {
  function OpcUaReadNode(config) {
    RED.nodes.createNode(this, config);
    const node = this;

    node.connectionConfig = RED.nodes.getNode(config.connection);
    if (!node.connectionConfig) {
      node.error('OPC UA connection not configured');
      return;
    }

    node.staticNodeId = config.nodeId || '';
    node.client = node.connectionConfig.client;
    node.client.addListener(node);

    node.status({ fill: 'grey', shape: 'ring', text: 'idle' });

    node.on('input', async function (msg, send, done) {
      send = send || function () { node.send.apply(node, arguments); };

      const nodeId = node.staticNodeId || msg.topic || msg.nodeId;
      if (!nodeId) {
        node.status({ fill: 'red', shape: 'ring', text: 'missing nodeId' });
        if (done) done(new Error('Missing NodeId'));
        return;
      }

      try {
        node.status({ fill: 'blue', shape: 'dot', text: 'reading' });
        const result = await node.client.read(nodeId);

        msg.topic = nodeId;
        msg.payload = result.value;
        msg.dataType = result.dataType;
        msg.sourceTimestamp = result.sourceTimestamp;
        msg.serverTimestamp = result.serverTimestamp;
        msg.statusCode = result.statusCode;

        node.status({ fill: 'green', shape: 'dot', text: 'ok' });
        send(msg);
        if (done) done();
      } catch (err) {
        node.status({ fill: 'red', shape: 'ring', text: err.message });
        if (done) done(err);
        else node.error(err, msg);
      }
    });

    node.on('close', function (done) {
      if (node.client) {
        node.client.removeListener(node);
      }
      node.status({});
      if (done) done();
    });

    node.onOpcUaEvent = function (event, payload) {
      if (event === 'connected') {
        node.status({ fill: 'green', shape: 'dot', text: 'connected' });
      } else if (event === 'disconnected') {
        node.status({ fill: 'red', shape: 'ring', text: 'disconnected' });
      } else if (event === 'error') {
        node.status({ fill: 'red', shape: 'ring', text: payload.message });
      }
    };
  }

  RED.nodes.registerType('opcua-read', OpcUaReadNode);
};
