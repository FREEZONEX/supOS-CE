'use strict';

module.exports = function (RED) {
  function OpcUaSubscribeNode(config) {
    RED.nodes.createNode(this, config);
    const node = this;

    node.connectionConfig = RED.nodes.getNode(config.connection);
    if (!node.connectionConfig) {
      node.error('OPC UA connection not configured');
      return;
    }

    node.staticNodeId = config.nodeId || '';
    node.samplingInterval = Number(config.samplingInterval) || 1000;
    node.subscriptionId = null;
    node.client = node.connectionConfig.client;
    node.client.addListener(node);

    node.status({ fill: 'grey', shape: 'ring', text: 'idle' });

    async function startSubscription() {
      const nodeId = node.staticNodeId;
      if (!nodeId) {
        node.status({ fill: 'red', shape: 'ring', text: 'missing nodeId' });
        return;
      }

      try {
        node.status({ fill: 'blue', shape: 'dot', text: 'subscribing' });
        const result = await node.client.subscribe(nodeId, node.samplingInterval, (err, data) => {
          if (err) {
            node.error(err);
            return;
          }
          const msg = {
            topic: data.nodeId,
            payload: data.value,
            dataType: data.dataType,
            sourceTimestamp: data.sourceTimestamp,
            serverTimestamp: data.serverTimestamp,
            statusCode: data.statusCode,
          };
          node.send(msg);
        });

        node.subscriptionId = result.subscriptionId;
        node.status({ fill: 'green', shape: 'dot', text: 'subscribed' });
      } catch (err) {
        node.status({ fill: 'red', shape: 'ring', text: err.message });
        node.error(err);
      }
    }

    node.on('input', async function (msg, send, done) {
      send = send || function () { node.send.apply(node, arguments); };
      // Trigger subscription start on first input; ignore subsequent inputs.
      if (!node.subscriptionId) {
        await startSubscription();
      }
      if (done) done();
    });

    node.on('close', async function (done) {
      if (node.client) {
        if (node.subscriptionId) {
          try {
            await node.client.unsubscribe(node.subscriptionId);
          } catch (err) {
            // Ignore.
          }
          node.subscriptionId = null;
        }
        node.client.removeListener(node);
      }
      node.status({});
      if (done) done();
    });

    node.onOpcUaEvent = function (event, payload) {
      if (event === 'connected') {
        if (node.subscriptionId) {
          node.subscriptionId = null;
          startSubscription();
        }
        node.status({ fill: 'green', shape: 'dot', text: 'connected' });
      } else if (event === 'disconnected') {
        node.subscriptionId = null;
        node.status({ fill: 'red', shape: 'ring', text: 'disconnected' });
      } else if (event === 'error') {
        node.status({ fill: 'red', shape: 'ring', text: payload.message });
      }
    };
  }

  RED.nodes.registerType('opcua-subscribe', OpcUaSubscribeNode);
};
