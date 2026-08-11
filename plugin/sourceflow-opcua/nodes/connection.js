'use strict';

const path = require('path');

let addon;
try {
  addon = require('../build/Release/opcua_open62541.node');
} catch (err) {
  // Fallback: allow the module to load so Node-RED can still enumerate nodes
  // when the native addon has not been built yet.
  addon = null;
}

const sharedClients = new Map();

class OpcUaClient {
  constructor(config) {
    this.config = config;
    this.id = `${config.endpoint}|${config.securityMode}|${config.securityPolicy}|${config.user}`;
    this.refCount = 0;
    this.connected = false;
    this.connecting = false;
    this.pendingConnect = null;
    this.listeners = new Set();
    this.reconnectTimer = null;
    this.handle = null;
    this.subscriptions = new Map();
  }

  async _ensureAddon() {
    if (!addon) {
      throw new Error(
        'Native open62541 addon is not available. Run `npm run build` in the package directory.'
      );
    }
  }

  async connect() {
    if (this.connected) {
      return;
    }
    if (this.connecting) {
      return this.pendingConnect;
    }

    this.connecting = true;
    this.pendingConnect = this._doConnect();
    return this.pendingConnect;
  }

  async _doConnect() {
    try {
      await this._ensureAddon();

      const options = {
        endpoint: this.config.endpoint,
        securityMode: this.config.securityMode,
        securityPolicy: this.config.securityPolicy,
        user: this.config.user,
        password: this.config.password,
        certificate: this.config.certificate,
        privateKey: this.config.privateKey,
        requestedSessionTimeout: this.config.requestedSessionTimeout,
      };

      this.handle = await addon.connectAsync(options);
      this.connected = true;
      this.connecting = false;
      this.pendingConnect = null;
      this._notifyListeners('connected');
      this.node && this.node.log(`OPC UA connected: ${this.config.endpoint}`);
    } catch (err) {
      this.connected = false;
      this.connecting = false;
      this.pendingConnect = null;
      this._notifyListeners('error', err);
      this.node && this.node.error(`OPC UA connect failed: ${err.message}`);
      this._scheduleReconnect(true);
      throw err;
    }
  }

  async disconnect() {
    this._clearReconnect();
    if (!this.connected || !this.handle) {
      return;
    }
    try {
      await addon.disconnectAsync(this.handle);
    } catch (err) {
      // Ignore disconnect errors.
    } finally {
      this.connected = false;
      this.handle = null;
      this.subscriptions.clear();
      this._notifyListeners('disconnected');
    }
  }

  async read(nodeId) {
    await this.connect();
    return addon.readAsync(this.handle, nodeId);
  }

  async write(nodeId, value, dataType) {
    await this.connect();
    return addon.writeAsync(this.handle, nodeId, value, dataType);
  }

  async subscribe(nodeId, samplingInterval, callback) {
    await this.connect();
    const result = await addon.subscribeAsync(this.handle, nodeId, samplingInterval || 1000, callback);
    if (result.subscriptionId) {
      this.subscriptions.set(result.subscriptionId, { nodeId, callback });
    }
    return result;
  }

  async unsubscribe(subscriptionId) {
    if (!this.connected || !this.handle) {
      return;
    }
    const result = await addon.unsubscribeAsync(this.handle, subscriptionId);
    this.subscriptions.delete(subscriptionId);
    return result;
  }

  addListener(node) {
    this.listeners.add(node);
    this.refCount++;
  }

  removeListener(node) {
    if (this.listeners.delete(node)) {
      this.refCount--;
      if (this.refCount <= 0) {
        this.disconnect();
        sharedClients.delete(this.id);
      }
    }
  }

  _notifyListeners(event, payload) {
    for (const node of this.listeners) {
      try {
        if (typeof node.onOpcUaEvent === 'function') {
          node.onOpcUaEvent(event, payload);
        }
      } catch (err) {
        // Listener errors should not break the client.
      }
    }
  }

  _scheduleReconnect(afterError) {
    this._clearReconnect();
    if (!afterError && this.connected) {
      return;
    }
    const interval = Number(this.config.reconnectInterval) || 5000;
    this.reconnectTimer = setTimeout(() => {
      if (!this.connected && this.refCount > 0) {
        this.connect().catch(() => {
          // Error already logged in _doConnect.
        });
      }
    }, interval);
  }

  _clearReconnect() {
    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer);
      this.reconnectTimer = null;
    }
  }
}

function getSharedClient(config) {
  const key = `${config.endpoint}|${config.securityMode}|${config.securityPolicy}|${config.user}`;
  if (!sharedClients.has(key)) {
    sharedClients.set(key, new OpcUaClient(config));
  }
  return sharedClients.get(key);
}

module.exports = function (RED) {
  function OpcUaConnectionNode(config) {
    RED.nodes.createNode(this, config);
    const node = this;

    node.endpoint = config.endpoint;
    node.securityMode = config.securityMode;
    node.securityPolicy = config.securityPolicy;
    node.user = config.user;
    node.password = config.password;
    node.certificate = config.certificate;
    node.privateKey = config.privateKey;
    node.requestedSessionTimeout = Number(config.requestedSessionTimeout) || 60000;
    node.reconnectInterval = Number(config.reconnectInterval) || 5000;

    node.client = getSharedClient({
      endpoint: node.endpoint,
      securityMode: node.securityMode,
      securityPolicy: node.securityPolicy,
      user: node.user,
      password: node.password,
      certificate: node.certificate,
      privateKey: node.privateKey,
      requestedSessionTimeout: node.requestedSessionTimeout,
      reconnectInterval: node.reconnectInterval,
    });
    node.client.node = node;
  }

  RED.nodes.registerType('opcua-connection', OpcUaConnectionNode);

  // Expose helper for operation nodes.
  RED.plugins = RED.plugins || {};
  RED.plugins.opcuaOpen62541 = {
    getSharedClient,
  };
};
