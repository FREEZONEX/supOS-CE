'use strict';

// Standalone integration test for the open62541 N-API addon.
// Starts a Python OPC UA fixture server, connects a client, reads/writes/subscribes,
// and asserts values are correct. Does not require Node-RED runtime.

const assert = require('assert');
const { spawn } = require('child_process');
const path = require('path');

let addon;
try {
  addon = require('../build/Release/opcua_open62541.node');
} catch (err) {
  console.log('Native addon not built; skipping integration test.');
  console.log('Run `npm run build` (or `npm run build:dev`) to compile the addon.');
  process.exit(0);
}

function sleep(ms) {
  return new Promise((resolve) => setTimeout(resolve, ms));
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

async function run() {
  const port = 4841;

  console.log('Starting fixture server...');
  const { proc, endpointUrl, nodeId } = await startFixtureServer(port);
  console.log(`Fixture server ready at ${endpointUrl}, nodeId=${nodeId}`);

  try {
    // Connect client.
    console.log('Connecting client...');
    const clientId = await addon.connectAsync({
      endpoint: endpointUrl,
      securityMode: 'None',
      securityPolicy: 'None',
    });
    assert(typeof clientId === 'number' && clientId > 0, 'client should connect and return an id');
    console.log('Client connected:', clientId);

    // Read initial value.
    console.log('Reading initial value...');
    const readResult = await addon.readAsync(clientId, nodeId);
    assert.strictEqual(readResult.statusCode, 'Good', `read should succeed: ${readResult.statusCode}`);
    assert.strictEqual(readResult.dataType, 'Double', 'read should report Double');
    assert.strictEqual(readResult.value, 12.3, 'read should return initial value');

    // Write new value.
    console.log('Writing new value...');
    const writeResult = await addon.writeAsync(clientId, nodeId, 23.0, 'Double');
    console.log('Write result:', writeResult.statusCode);
    assert.strictEqual(writeResult.statusCode, 'Good', `write should succeed: ${writeResult.statusCode}`);

    // Read updated value.
    console.log('Reading updated value...');
    const readResult2 = await addon.readAsync(clientId, nodeId);
    assert.strictEqual(readResult2.value, 23.0, 'read should return updated value');

    // Subscribe and wait for notification.
    console.log('Subscribing...');
    const notifications = [];
    const subResult = await addon.subscribeAsync(clientId, nodeId, 100, (err, data) => {
      if (err) {
        console.error('subscription callback error:', err);
        return;
      }
      notifications.push(data);
    });
    assert.strictEqual(subResult.statusCode, 'Good', `subscribe should succeed: ${subResult.statusCode}`);
    assert(subResult.subscriptionId > 0, 'subscribe should return a subscription id');
    console.log('Subscribed:', subResult.subscriptionId);

    // Write a new value via the client to trigger a data-change notification.
    console.log('Writing value to trigger notification...');
    await addon.writeAsync(clientId, nodeId, 25.0, 'Double');

    // Wait for notification to arrive.
    await sleep(1200);

    assert(notifications.length > 0, 'subscription should receive at least one notification');
    const last = notifications[notifications.length - 1];
    assert.strictEqual(last.value, 25.0, 'last subscription notification should match written value');

    // Unsubscribe.
    const unsubResult = await addon.unsubscribeAsync(clientId, subResult.subscriptionId);
    assert.strictEqual(unsubResult.statusCode, 'Good', `unsubscribe should succeed: ${unsubResult.statusCode}`);

    // Disconnect client.
    await addon.disconnectAsync(clientId);

    console.log('All addon integration tests passed.');
  } finally {
    proc.kill();
    await new Promise((resolve) => proc.on('exit', resolve));
  }
}

run().catch((err) => {
  console.error('Test failed:', err);
  process.exit(1);
});
