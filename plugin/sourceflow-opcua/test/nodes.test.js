'use strict';

const assert = require('assert');
const { spawn } = require('child_process');
const path = require('path');

const PORT = 4842;

function delay(ms) {
  return new Promise(resolve => setTimeout(resolve, ms));
}

async function startFixtureServer(port) {
  const serverPath = path.join(__dirname, 'fixture', 'server.py');
  const proc = spawn('python3', [serverPath, String(port)], {
    stdio: ['ignore', 'pipe', 'pipe'],
  });

  let stdout = '';
  let stderr = '';
  proc.stdout.on('data', (chunk) => { stdout += chunk.toString(); });
  proc.stderr.on('data', (chunk) => { stderr += chunk.toString(); });

  const ready = await new Promise((resolve, reject) => {
    const timer = setTimeout(() => {
      proc.kill();
      reject(new Error(`Fixture server did not start in time. stderr: ${stderr}`));
    }, 10000);

    function check() {
      const match = stdout.match(/RUNNING (opc\.tcp:\/\/[^\s]+)(?:\s+(ns=[^\s]+))?/);
      if (match) {
        clearTimeout(timer);
        resolve({ endpointUrl: match[1], nodeId: match[2] || 'ns=1;s=TestNode' });
        return;
      }
      if (!proc.killed && proc.exitCode === null) {
        setTimeout(check, 50);
      } else {
        clearTimeout(timer);
        reject(new Error(`Fixture server exited early. stderr: ${stderr}`));
      }
    }
    check();
  });

  return { proc, endpointUrl: ready.endpointUrl, nodeId: ready.nodeId };
}

function createMockRed() {
  const configNodes = new Map();
  const registeredTypes = new Map();

  const RED = {
    log: { debug: () => {}, info: () => {}, warn: () => {}, error: () => {}, trace: () => {} },
    nodes: {
      createNode(node, config) {
        node.id = config.id || `node-${Math.random().toString(36).slice(2)}`;
        node.name = config.name || '';
        node.type = config.type || '';
        node._status = [];
        node._sent = [];
        node._errors = [];
        node._handlers = {};
        node.status = obj => node._status.push(obj);
        node.send = msg => node._sent.push(msg);
        node.error = (err, msg) => node._errors.push({ err, msg });
        node.log = () => {};
        node.warn = () => {};
        node.on = (event, handler) => {
          node._handlers[event] = node._handlers[event] || [];
          node._handlers[event].push(handler);
        };
        node.close = async done => {
          for (const handler of node._handlers.close || []) {
            await new Promise(resolve => handler(resolve));
          }
          if (done) done();
        };
        node._emitInput = async (msg, send, done) => {
          for (const handler of node._handlers.input || []) {
            await handler(msg, send || node.send, done || (() => {}));
          }
        };
      },
      getNode(id) {
        return configNodes.get(id) || null;
      },
      registerType(type, constructor) {
        registeredTypes.set(type, constructor);
      },
    },
  };

  return { RED, configNodes, registeredTypes };
}

async function createNode(RED, registeredTypes, type, config, configNodes) {
  const Constructor = registeredTypes.get(type);
  assert(Constructor, `Node type ${type} is not registered`);
  const instance = new Constructor(config);
  if (configNodes && config.id) {
    configNodes.set(config.id, instance);
  }
  return instance;
}

async function expectStatus(node, text, timeoutMs = 3000) {
  const deadline = Date.now() + timeoutMs;
  while (Date.now() < deadline) {
    if (node._status.some(s => s.text === text)) return;
    await delay(50);
  }
  throw new Error(`Expected status '${text}' but got ${JSON.stringify(node._status)}`);
}

async function expectSent(node, predicate, timeoutMs = 3000) {
  const deadline = Date.now() + timeoutMs;
  while (Date.now() < deadline) {
    const found = node._sent.find(predicate);
    if (found) return found;
    await delay(50);
  }
  throw new Error(`No matching message sent; sent=${JSON.stringify(node._sent)}`);
}

(async () => {
  console.log('Loading node modules...');

  const { proc, endpointUrl, nodeId } = await startFixtureServer(PORT);
  console.log(`Fixture server ready at ${endpointUrl}, nodeId=${nodeId}`);

  try {
    const clientRed = createMockRed();
    require('../nodes/connection')(clientRed.RED);
    require('../nodes/read')(clientRed.RED);
    require('../nodes/write')(clientRed.RED);
    require('../nodes/subscribe')(clientRed.RED);

    assert.ok(clientRed.registeredTypes.has('opcua-connection'));
    assert.ok(clientRed.registeredTypes.has('opcua-read'));
    assert.ok(clientRed.registeredTypes.has('opcua-write'));
    assert.ok(clientRed.registeredTypes.has('opcua-subscribe'));

    const connection = await createNode(
      clientRed.RED,
      clientRed.registeredTypes,
      'opcua-connection',
      { id: 'conn1', type: 'opcua-connection', endpoint: endpointUrl, securityMode: 'None', securityPolicy: 'None' },
      clientRed.configNodes
    );

    const readNode = await createNode(
      clientRed.RED,
      clientRed.registeredTypes,
      'opcua-read',
      { id: 'read1', type: 'opcua-read', connection: 'conn1', nodeId }
    );

    await readNode._emitInput({});
    const readMsg = await expectSent(readNode, m => m.statusCode === 'Good');
    assert.strictEqual(readMsg.payload, 12.3);
    assert.strictEqual(connection.client.connected, true, 'client should be connected after first read');

    const writeNode = await createNode(
      clientRed.RED,
      clientRed.registeredTypes,
      'opcua-write',
      { id: 'write1', type: 'opcua-write', connection: 'conn1', nodeId, dataType: 'Double' }
    );

    await writeNode._emitInput({ payload: 45.6 });
    const writeMsg = await expectSent(writeNode, m => m.statusCode === 'Good');
    assert.strictEqual(writeMsg.statusCode, 'Good');

    readNode._sent.length = 0;
    await readNode._emitInput({});
    const readMsg2 = await expectSent(readNode, m => m.payload === 45.6);
    assert.strictEqual(readMsg2.statusCode, 'Good');

    // Subscribe
    const subscribeNode = await createNode(
      clientRed.RED,
      clientRed.registeredTypes,
      'opcua-subscribe',
      { id: 'sub1', type: 'opcua-subscribe', connection: 'conn1', nodeId, samplingInterval: 100 }
    );

    const notificationPromise = new Promise(resolve => {
      const originalSend = subscribeNode.send;
      let resolved = false;
      subscribeNode.send = msg => {
        originalSend.call(subscribeNode, msg);
        if (!resolved && msg.payload === 78.9) {
          resolved = true;
          subscribeNode.send = originalSend;
          resolve(msg);
        }
      };
    });

    await subscribeNode._emitInput({});
    await delay(300);
    assert.ok(subscribeNode.subscriptionId, 'subscribe node should have a subscription id');

    // Trigger a value change from the client side.
    await writeNode._emitInput({ payload: 78.9 });

    const notified = await Promise.race([
      notificationPromise,
      delay(2000).then(() => null),
    ]);
    assert.ok(notified, 'subscribe node should receive changed value within 2 seconds');
    assert.strictEqual(notified.payload, 78.9);
    assert.strictEqual(notified.statusCode, 'Good');

    // Cleanup
    await subscribeNode.close();
    await writeNode.close();
    await readNode.close();
    await connection.close();

    await delay(500);
    assert.strictEqual(connection.client.connected, false, 'client should be disconnected after close');

    console.log('All Node-RED node tests passed.');
  } finally {
    proc.kill();
    await new Promise(resolve => proc.on('exit', resolve));
  }
})().catch(err => {
  console.error(err);
  process.exit(1);
});
