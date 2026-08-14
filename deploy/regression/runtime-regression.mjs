#!/usr/bin/env node

import { spawnSync } from 'node:child_process';
import { mkdir, readFile, writeFile } from 'node:fs/promises';
import { dirname, resolve } from 'node:path';
import { fileURLToPath } from 'node:url';

import { regressionContract } from './runtime-contract.mjs';

const scriptPath = fileURLToPath(import.meta.url);

const parseArgs = (args) => {
  const options = {};
  for (let index = 0; index < args.length; index += 1) {
    const arg = args[index];
    if (!arg.startsWith('--')) throw new Error(`unexpected argument: ${arg}`);
    const key = arg.slice(2);
    const value = args[index + 1];
    if (!value || value.startsWith('--')) throw new Error(`missing value for --${key}`);
    options[key] = value;
    index += 1;
  }
  return options;
};

export const normalizeBaseURL = (value) => {
  const url = new URL(value);
  if (!['http:', 'https:'].includes(url.protocol)) throw new Error('base URL must use HTTP or HTTPS');
  url.pathname = '/';
  url.search = '';
  url.hash = '';
  return url.toString().replace(/\/$/, '');
};

export const missingResourceKeys = (actual = [], required = regressionContract.requiredResourceKeys) =>
  required.filter((key) => !actual.includes(key));

const parseBrowserObjects = (output) => {
  const objects = [];
  for (const line of output.split('\n')) {
    const trimmed = line.trim();
    if (!trimmed.startsWith('"{\\"')) continue;
    try {
      objects.push(JSON.parse(JSON.parse(trimmed)));
    } catch {
      // Ignore non-JSON browser output that happens to start with a quote.
    }
  }
  return objects;
};

export const parseBrowserState = (output) => {
  const state = parseBrowserObjects(output).find((item) => Object.hasOwn(item, 'path'));
  if (state) return state;
  throw new Error('browser state was not emitted');
};

const forbiddenFlowCredentialKeys = new Set([
  'authtoken',
  'clientid',
  'credentialclientid',
  'credentials',
  'password',
  'username',
]);

const findForbiddenFlowCredentialPaths = (value, path = '$', found = []) => {
  if (Array.isArray(value)) {
    value.forEach((item, index) => findForbiddenFlowCredentialPaths(item, `${path}[${index}]`, found));
    return found;
  }
  if (!value || typeof value !== 'object') return found;
  for (const [key, item] of Object.entries(value)) {
    const itemPath = `${path}.${key}`;
    if (forbiddenFlowCredentialKeys.has(key.toLowerCase())) found.push(itemPath);
    findForbiddenFlowCredentialPaths(item, itemPath, found);
  }
  return found;
};

const wiredTo = (node, targetID) =>
  Array.isArray(node?.wires) && node.wires.flat(Infinity).map(String).includes(String(targetID));

export const inspectMockFlow = (flowData, expectedTopic, expectedFieldName = regressionContract.foundationJourney.field.name) => {
  const nodes = typeof flowData === 'string' ? JSON.parse(flowData) : flowData;
  if (!Array.isArray(nodes)) throw new Error('mock Flow data is not a node array');
  const inject = nodes.find((node) => node?.type === 'inject');
  const functionNode = nodes.find((node) => node?.type === 'function');
  const mqttOut = nodes.find((node) => node?.type === 'mqtt out');
  const broker = nodes.find((node) => node?.type === 'mqtt-broker');
  if (!inject || !functionNode || !mqttOut || !broker) {
    throw new Error('mock Flow is missing inject, function, MQTT output or broker nodes');
  }
  const repeatSeconds = Number(inject.repeat);
  if (inject.once !== true || !Number.isFinite(repeatSeconds) || repeatSeconds <= 0) {
    throw new Error('mock Flow inject node is not configured for immediate and periodic execution');
  }
  if (!wiredTo(inject, functionNode.id) || !wiredTo(functionNode, mqttOut.id)) {
    throw new Error('mock Flow inject/function/MQTT nodes are not connected');
  }
  if (String(mqttOut.broker) !== String(broker.id) || broker.broker !== 'emqx' || String(broker.port) !== '1883') {
    throw new Error('mock Flow is not connected to the retained anonymous EMQX transport');
  }
  if (broker.name !== 'emqx') {
    throw new Error(`mock Flow MQTT configuration node is named ${broker.name || '<empty>'}, expected emqx`);
  }
  if (!String(functionNode.func || '').includes(JSON.stringify(expectedTopic))) {
    throw new Error(`mock Flow does not publish to ${expectedTopic}`);
  }
  if (!String(functionNode.func || '').includes(JSON.stringify(expectedFieldName))) {
    throw new Error(`mock Flow does not construct the ${expectedFieldName} field`);
  }
  const forbiddenCredentialPaths = findForbiddenFlowCredentialPaths(nodes);
  if (forbiddenCredentialPaths.length) {
    throw new Error(`mock Flow contains MQTT credential fields: ${forbiddenCredentialPaths.join(', ')}`);
  }
  return {
    nodeCount: nodes.length,
    repeatSeconds,
    topic: expectedTopic,
    broker: `${broker.broker}:${broker.port}`,
    configurationName: broker.name,
    credentialFree: true,
  };
};

export const hasPersistedMetricData = (detail, dashboard, fieldName = 'value') => {
  const latestValue = detail?.lastPayload?.data?.[fieldName];
  const history = Array.isArray(dashboard?.list) ? dashboard.list : [];
  return typeof latestValue === 'number' && history.some((item) => typeof item?.payload?.[fieldName] === 'number');
};

const sleep = (milliseconds) => new Promise((resolvePromise) => setTimeout(resolvePromise, milliseconds));

const runCommand = (command, args, options = {}) => {
  const result = spawnSync(command, args, { encoding: 'utf8', ...options });
  const output = `${result.stdout || ''}${result.stderr || ''}`;
  if (result.error) throw result.error;
  if (result.status !== 0) {
    throw new Error(`${command} ${args.join(' ')} failed (${result.status}):\n${output}`);
  }
  return String(result.stdout || '').trim();
};

const runCompose = (deployRoot, args) =>
  runCommand('bash', [resolve(deployRoot, 'bin/compose.sh'), ...args], { cwd: deployRoot });

const waitForComposeServiceHealthy = async (deployRoot, service, timeoutMs = 45_000) => {
  const deadline = Date.now() + timeoutMs;
  let lastStatus = 'missing';
  while (Date.now() < deadline) {
    const containerID = runCompose(deployRoot, ['ps', '-q', service]);
    if (containerID) {
      lastStatus = runCommand('docker', [
        'inspect',
        '--format',
        '{{if .State.Health}}{{.State.Health.Status}}{{else}}{{.State.Status}}{{end}}',
        containerID,
      ]);
      if (lastStatus === 'healthy' || lastStatus === 'running') return lastStatus;
    }
    await sleep(1_000);
  }
  throw new Error(`${service} did not become healthy within ${timeoutMs}ms (last status: ${lastStatus})`);
};

const request = async (baseURL, path, options = {}) => {
  const { signal, ...fetchOptions } = options;
  const response = await fetch(`${baseURL}${path}`, {
    redirect: 'manual',
    ...fetchOptions,
    signal: signal || AbortSignal.timeout(regressionContract.foundationJourney.requestTimeoutMs),
    headers: {
      ...(options.body ? { 'content-type': 'application/json' } : {}),
      ...(options.headers || {}),
    },
  });
  const contentType = response.headers.get('content-type') || '';
  const rawBody = await response.text();
  const body = contentType.includes('application/json') && rawBody ? JSON.parse(rawBody) : rawBody;
  return { status: response.status, body, contentType };
};

const requireEnvelopeSuccess = (name, response) => {
  if (response.status !== 200 || response.body?.code !== 200) {
    throw new Error(`${name} failed: HTTP ${response.status}, code ${response.body?.code ?? 'none'}`);
  }
  return response.body.data;
};

const quoteBatchValue = (value) => JSON.stringify(String(value));
const browserEval = (expression) => `eval ${JSON.stringify(expression)}`;

const runBrowserRegression = ({ baseURL, username, password, screenshot }) => {
  const session = `edge-open-source-regression-${process.pid}-${Date.now()}`;
  const loginURL = `${baseURL}${regressionContract.loginPath}?redirectUri=${encodeURIComponent(regressionContract.defaultRoute)}`;
  const commands = [
    `open ${loginURL}`,
    'wait --load networkidle',
    `fill input[autocomplete=username] ${quoteBatchValue(username)}`,
    `fill input[type=password] ${quoteBatchValue(password)}`,
    'click button[type=submit]',
    'wait --load networkidle',
    browserEval('JSON.stringify({path:location.pathname,body:document.body.innerText.slice(0,1000)})'),
    `screenshot ${screenshot}`,
    `open ${baseURL}/settings/profile`,
    'wait --load networkidle',
    browserEval("JSON.stringify({settingsPath:location.pathname,settingsItems:document.querySelector('aside').innerText,settingsBody:document.body.innerText.slice(0,1500)})"),
    ...regressionContract.retainedSettingsRoutes.flatMap((path) => [
      `open ${baseURL}${path}`,
      'wait --load networkidle',
      browserEval(`JSON.stringify({retainedSettingsRoute:${JSON.stringify(path)},path:location.pathname,body:document.body.innerText.slice(0,3000)})`),
      browserEval(`(() => { const buttons = [...document.querySelectorAll('button')]; const roleButton = buttons.find((item) => /Role Permission|Role Settings|角色权限/.test(item.innerText || '')); const newUserButton = buttons.find((item) => /New Users?|新增用户|新建用户/.test(item.innerText || '')); const headers = [...document.querySelectorAll('th')].map((item) => (item.innerText || '').trim()).filter(Boolean); if (newUserButton) newUserButton.click(); return JSON.stringify({rolePermissionTrigger:Boolean(roleButton),newUserTrigger:Boolean(newUserButton),userTableHeaders:headers}); })()`),
      'wait --load networkidle',
      browserEval(`(() => { const dialogs = [...document.querySelectorAll('[role="dialog"]')]; const dialog = dialogs.at(-1); return JSON.stringify({newUserDialog:Boolean(dialog),newUserDialogText:(dialog?.innerText || '').slice(0,4000),newUserDialogLabels:[...(dialog?.querySelectorAll('label') || [])].map((item) => (item.innerText || '').trim()).filter(Boolean)}); })()`),
    ]),
    ...regressionContract.disabledSettingsRoutes.flatMap((path) => [
      `open ${baseURL}${path}`,
      'wait --load networkidle',
      browserEval(`JSON.stringify({disabledSettingsRoute:${JSON.stringify(path)},path:location.pathname,body:document.body.innerText.slice(0,1500)})`),
    ]),
    ...regressionContract.disabledFrontendRoutes.flatMap((path) => [
      `open ${baseURL}${path}`,
      'wait --load networkidle',
      browserEval(`JSON.stringify({disabledRoute:${JSON.stringify(path)},path:location.pathname,body:document.body.innerText.slice(0,1000)})`),
    ]),
    `network requests --filter ${regressionContract.forbiddenBrowserRequests[0]}`,
    'errors',
    'close',
  ];
  const browserPackage = process.env.AGENT_BROWSER_PACKAGE || 'agent-browser@0.27.0';
  const result = spawnSync(
    'npx',
    ['-y', browserPackage, '--session', session, 'batch', '--bail', ...commands],
    {
      cwd: dirname(scriptPath),
      encoding: 'utf8',
      timeout: 120_000,
      env: {
        ...process.env,
        AGENT_BROWSER_ARGS: process.env.AGENT_BROWSER_ARGS || '--no-sandbox',
      },
    }
  );
  const output = `${result.stdout || ''}${result.stderr || ''}`;
  if (result.error) throw result.error;
  if (result.status !== 0) throw new Error(`browser regression command failed (${result.status}):\n${output}`);
  const browserObjects = parseBrowserObjects(output);
  const state = parseBrowserState(output);
  if (state.path !== regressionContract.defaultRoute) {
    throw new Error(`browser landed on ${state.path || '<empty>'}, expected ${regressionContract.defaultRoute}`);
  }
  if (regressionContract.forbiddenBrowserPaths.includes(state.path)) {
    throw new Error(`browser landed on forbidden path ${state.path}`);
  }
  for (const path of regressionContract.forbiddenBrowserRequests) {
    if (output.includes(` ${baseURL}${path} `) || output.includes(`${baseURL}${path} (XHR)`)) {
      throw new Error(`browser requested removed API ${path}`);
    }
  }
  if (/Not found\.|No permission|无权访问|Interface does not exist/i.test(state.body || '')) {
    throw new Error('browser rendered a not-found or no-permission state');
  }
  const settingsState = browserObjects.find((item) => Object.hasOwn(item, 'settingsPath'));
  if (!settingsState) throw new Error('settings browser state was not emitted');
  const settingsText = [settingsState.settingsItems || '', settingsState.settingsBody || ''].join('\n');
  const leakedSettingsLabels = regressionContract.disabledSettingsLabels.filter((label) => settingsText.includes(label));
  if (leakedSettingsLabels.length) {
    throw new Error(`disabled settings remain visible: ${leakedSettingsLabels.join(', ')}`);
  }
  const retainedSettingsRouteStates = browserObjects.filter((item) => Object.hasOwn(item, 'retainedSettingsRoute'));
  for (const path of regressionContract.retainedSettingsRoutes) {
    const routeState = retainedSettingsRouteStates.find((item) => item.retainedSettingsRoute === path);
    if (!routeState) throw new Error(`retained settings route state was not emitted: ${path}`);
    const body = String(routeState.body || '');
    if (routeState.path !== path || /Interface does not exist|接口不存在|No data|暂无数据/i.test(body)) {
      throw new Error(`retained user management route is not functional: ${path}`);
    }
    if (!body.includes('tier0') || !/User Management|用户管理/i.test(body)) {
      throw new Error(`retained user management route did not render the tier0 account: ${path}`);
    }
  }
  const userManagementState = browserObjects.find((item) => Object.hasOwn(item, 'rolePermissionTrigger'));
  if (!userManagementState) throw new Error('user management UI state was not emitted');
  if (userManagementState.rolePermissionTrigger) throw new Error('Role Permission trigger remained visible');
  const leakedHeaders = (userManagementState.userTableHeaders || []).filter((label) =>
    regressionContract.hiddenUserTableHeaders.includes(label)
  );
  if (leakedHeaders.length) throw new Error(`Role table column remained visible: ${leakedHeaders.join(', ')}`);
  if (!userManagementState.newUserTrigger) throw new Error('New User trigger was not rendered');
  const newUserDialog = browserObjects.find((item) => Object.hasOwn(item, 'newUserDialog'));
  if (!newUserDialog?.newUserDialog) throw new Error('New User dialog did not render');
  const leakedDialogLabels = (newUserDialog.newUserDialogLabels || []).filter((label) =>
    regressionContract.hiddenUserTableHeaders.includes(label)
  );
  if (leakedDialogLabels.length) {
    throw new Error(`New User dialog still requires a role selection: ${leakedDialogLabels.join(', ')}`);
  }
  const disabledSettingsRouteStates = browserObjects.filter((item) => Object.hasOwn(item, 'disabledSettingsRoute'));
  for (const path of regressionContract.disabledSettingsRoutes) {
    const routeState = disabledSettingsRouteStates.find((item) => item.disabledSettingsRoute === path);
    if (!routeState) throw new Error(`disabled settings route state was not emitted: ${path}`);
    const body = String(routeState.body || '');
    if (/Interface does not exist|接口不存在/i.test(body)) {
      throw new Error(`disabled settings route requested a removed interface: ${path}`);
    }
    const leakedLabels = regressionContract.disabledSettingsLabels.filter((label) => body.includes(label));
    if (leakedLabels.length) {
      throw new Error(`disabled settings route still renders IAM management: ${path} (${leakedLabels.join(', ')})`);
    }
  }
  const disabledRouteStates = browserObjects.filter((item) => Object.hasOwn(item, 'disabledRoute'));
  for (const path of regressionContract.disabledFrontendRoutes) {
    const routeState = disabledRouteStates.find((item) => item.disabledRoute === path);
    if (!routeState) throw new Error(`disabled frontend route state was not emitted: ${path}`);
    const body = String(routeState.body || '');
    if (regressionContract.disabledSettingsLabels.some((label) => body.includes(label))) {
      throw new Error(`disabled frontend route still renders its management page: ${path}`);
    }
    const visiblyNotFound = /not found|404|页面不存在|未找到/i.test(body);
    if (routeState.path === path && !visiblyNotFound) {
      throw new Error(`disabled frontend route remains reachable: ${path}`);
    }
  }
  return {
    path: state.path,
    screenshot,
    forbiddenRequestsObserved: false,
    notFoundOrPermissionState: false,
    disabledSettingsVisible: false,
    userManagementVisible: true,
    roleManagementHidden: true,
    newUserRoleSelectorHidden: true,
    disabledRoutes: regressionContract.disabledFrontendRoutes,
  };
};

const runDefaultAdminUserJourney = async ({ baseURL, authHeaders }) => {
  const name = `e2e_user_${Date.now().toString(36)}_${process.pid.toString(36)}`;
  let userID = 0;
  let primaryError;
  const evidence = { name, rolePayloadOmitted: true, assignedRole: '', created: false, deleted: false };
  try {
    const created = requireEnvelopeSuccess(
      'create local user without role payload',
      await request(baseURL, '/api/core/iam/users', {
        method: 'POST',
        headers: authHeaders,
        body: JSON.stringify({
          username: name,
          password: `Tier0#${Date.now()}Aa`,
          firstName: name,
          enabled: true,
        }),
      })
    );
    userID = Number(created?.userId ?? created?.id ?? 0);
    if (!userID) throw new Error('user create returned no user identifier');
    evidence.created = true;
    const createdRole = (created?.roleList || [])[0];
    if (String(createdRole?.roleCode || '').toLowerCase() !== 'admin' || Number(createdRole?.roleId) !== 1) {
      throw new Error('new user response was not assigned the built-in Admin role');
    }
    const users = requireEnvelopeSuccess(
      'read default-Admin local user',
      await request(baseURL, '/api/core/iam/users?pageNo=1&pageSize=100', { headers: authHeaders })
    );
    const user = (users?.list || []).find((item) => Number(item?.userId ?? item?.id) === userID);
    const listedRole = (user?.roleList || [])[0];
    if (!user || String(listedRole?.roleCode || '').toLowerCase() !== 'admin' || Number(listedRole?.roleId) !== 1) {
      throw new Error('new user did not persist the built-in Admin role');
    }
    evidence.assignedRole = 'admin';
  } catch (error) {
    primaryError = error;
  } finally {
    if (userID) {
      try {
        requireEnvelopeSuccess(
          `delete local user ${userID}`,
          await request(baseURL, `/api/core/iam/users/${userID}`, { method: 'DELETE', headers: authHeaders })
        );
        evidence.deleted = true;
      } catch (error) {
        if (primaryError) primaryError.message = `${primaryError.message}; user cleanup failed: ${error.message}`;
        else primaryError = error;
      }
    }
  }
  if (primaryError) throw primaryError;
  return evidence;
};

const cleanupFoundationJourney = async ({ baseURL, authHeaders, nodeID, flowID, name }) => {
  const cleanup = { flowDeleted: !flowID, unsNodeDeleted: !nodeID, errors: [] };
  if (!nodeID) {
    try {
      const listed = requireEnvelopeSuccess(
        `discover partial UNS artifact ${name}`,
        await request(baseURL, `/api/core/uns/nodes?keyword=${encodeURIComponent(name)}`, { headers: authHeaders })
      );
      const node = (listed?.list || []).find((item) => item?.name === name);
      nodeID = Number(node?.id || 0);
      cleanup.unsNodeDeleted = !nodeID;
    } catch {
      // The create failure remains primary; discovery is only best-effort cleanup.
    }
  }
  if (!flowID) {
    try {
      const listed = requireEnvelopeSuccess(
        `discover partial Flow artifact ${name}`,
        await request(baseURL, `/api/core/flows?flowType=source&keyword=${encodeURIComponent(name)}`, { headers: authHeaders })
      );
      const flow = (listed?.list || []).find((item) =>
        (nodeID && (item?.unsNodeIds || []).map(Number).includes(nodeID)) || String(item?.name || '').endsWith(`/${name}`)
      );
      flowID = Number(flow?.id || 0);
      cleanup.flowDeleted = !flowID;
    } catch {
      // The create failure remains primary; discovery is only best-effort cleanup.
    }
  }
  if (flowID) {
    try {
      requireEnvelopeSuccess(
        `cleanup Flow ${flowID}`,
        await request(baseURL, `/api/core/flows/${flowID}`, { method: 'DELETE', headers: authHeaders })
      );
      cleanup.flowDeleted = true;
    } catch (error) {
      cleanup.errors.push(error.message);
    }
  }
  if (nodeID) {
    try {
      requireEnvelopeSuccess(
        `cleanup UNS node ${nodeID}`,
        await request(baseURL, `/api/core/uns/nodes/${nodeID}`, { method: 'DELETE', headers: authHeaders })
      );
      cleanup.unsNodeDeleted = true;
    } catch (error) {
      cleanup.errors.push(error.message);
    }
  }
  return cleanup;
};

const discoverNamedArtifacts = async ({ baseURL, authHeaders, name }) => {
  const listedNodes = requireEnvelopeSuccess(
    `discover UNS artifact ${name}`,
    await request(baseURL, `/api/core/uns/nodes?keyword=${encodeURIComponent(name)}`, { headers: authHeaders })
  );
  const nodes = (listedNodes?.list || []).filter((item) => item?.name === name);
  const nodeIDs = new Set(nodes.map((item) => Number(item?.id || 0)).filter(Boolean));
  const listedFlows = requireEnvelopeSuccess(
    `discover Flow artifact ${name}`,
    await request(baseURL, `/api/core/flows?flowType=source&keyword=${encodeURIComponent(name)}`, { headers: authHeaders })
  );
  const flows = (listedFlows?.list || []).filter((item) =>
    String(item?.name || '') === name ||
    String(item?.name || '').endsWith(`/${name}`) ||
    (item?.unsNodeIds || []).map(Number).some((id) => nodeIDs.has(id))
  );
  return { nodes, flows };
};

const runAtomicityFaultInjection = async ({ baseURL, authHeaders, deployRoot }) => {
  if (!deployRoot) throw new Error('--deploy-root is required with --fault-injection true');
  const contract = regressionContract.foundationJourney;
  const suffix = `${Date.now().toString(36)}_${process.pid.toString(36)}`;
  const name = `${contract.metricNamePrefix}_retry_${suffix}`;
  const payload = {
    parentId: 0,
    name,
    displayName: `E2E Atomic Metric ${suffix}`,
    type: 'file',
    topicType: contract.topicType,
    schema: JSON.stringify([contract.field]),
    persistence: true,
    withFlow: true,
    addFlow: true,
    mockData: true,
  };
  let sourceFlowStopped = false;
  let nodeID = 0;
  let flowID = 0;
  let evidence;
  let primaryError;

  try {
    runCompose(deployRoot, ['stop', 'sourceflow']);
    sourceFlowStopped = true;
    const failed = await request(baseURL, '/api/core/uns/nodes', {
      method: 'POST',
      headers: authHeaders,
      body: JSON.stringify(payload),
    });
    if (failed.status === 200 && failed.body?.code === 200) {
      throw new Error('Metric create unexpectedly succeeded while SourceFlow was stopped');
    }
    const partial = await discoverNamedArtifacts({ baseURL, authHeaders, name });
    if (partial.nodes.length || partial.flows.length) {
      throw new Error(`failed Metric create left partial artifacts: UNS=${partial.nodes.length}, Flow=${partial.flows.length}`);
    }

    runCompose(deployRoot, ['start', 'sourceflow']);
    sourceFlowStopped = false;
    await waitForComposeServiceHealthy(deployRoot, 'sourceflow');
    const retried = requireEnvelopeSuccess(
      'retry identical Metric create after SourceFlow recovery',
      await request(baseURL, '/api/core/uns/nodes', {
        method: 'POST',
        headers: authHeaders,
        body: JSON.stringify(payload),
      })
    );
    nodeID = Number(retried?.id || 0);
    flowID = Number(retried?.flowId || retried?.flow?.id || 0);
    if (!nodeID || !flowID) throw new Error('identical Metric retry did not return both UNS and Flow identifiers');
    const flow = requireEnvelopeSuccess(
      `read retried generated Flow ${flowID}`,
      await request(baseURL, `/api/core/flows/${flowID}`, { headers: authHeaders })
    );
    if (flow?.status !== 'deployed' || !String(flow?.runtimeFlowId || '')) {
      throw new Error('identical Metric retry did not deploy its Mock Flow');
    }
    const flowRuntime = inspectMockFlow(flow?.flowData, String(retried?.namespace || ''), contract.field.name);
    evidence = {
      status: 'passed',
      firstAttemptFailed: true,
      partialUnsArtifacts: 0,
      partialFlowArtifacts: 0,
      identicalRetrySucceeded: true,
      mqttConfigurationName: flowRuntime.configurationName,
    };
  } catch (error) {
    primaryError = error;
  } finally {
    if (sourceFlowStopped) {
      try {
        runCompose(deployRoot, ['start', 'sourceflow']);
        await waitForComposeServiceHealthy(deployRoot, 'sourceflow');
      } catch (error) {
        primaryError ||= error;
        if (primaryError !== error) primaryError.message = `${primaryError.message}; restore SourceFlow: ${error.message}`;
      }
    }
  }

  const cleanup = await cleanupFoundationJourney({ baseURL, authHeaders, nodeID, flowID, name });
  if (cleanup.errors.length) {
    const cleanupError = new Error(`atomicity regression cleanup failed: ${cleanup.errors.join('; ')}`);
    if (primaryError) primaryError.message = `${primaryError.message}; ${cleanupError.message}`;
    else primaryError = cleanupError;
  }
  if (primaryError) throw primaryError;
  evidence.cleanup = { flowDeleted: cleanup.flowDeleted, unsNodeDeleted: cleanup.unsNodeDeleted };
  return evidence;
};

const runFoundationJourney = async ({ baseURL, authHeaders }) => {
  const contract = regressionContract.foundationJourney;
  const suffix = `${Date.now().toString(36)}_${process.pid.toString(36)}`;
  const name = `${contract.metricNamePrefix}_${suffix}`;
  const startedAt = Date.now();
  let nodeID = 0;
  let flowID = 0;
  let journey;
  let primaryError;

  try {
    const created = requireEnvelopeSuccess(
      'create Metric UNS with Mock Data and history',
      await request(baseURL, '/api/core/uns/nodes', {
        method: 'POST',
        headers: authHeaders,
        body: JSON.stringify({
          parentId: 0,
          name,
          displayName: `E2E Metric ${suffix}`,
          type: 'file',
          topicType: contract.topicType,
          schema: JSON.stringify([contract.field]),
          persistence: true,
          withFlow: true,
          addFlow: true,
          mockData: true,
        }),
      })
    );
    nodeID = Number(created?.id || 0);
    flowID = Number(created?.flowId || created?.flow?.id || 0);
    const topic = String(created?.namespace || '');
    if (!nodeID || !flowID || !topic) {
      throw new Error('created Metric UNS did not return node, Flow and topic identifiers');
    }
    if (String(created?.topicType).toLowerCase() !== contract.topicType.toLowerCase()) {
      throw new Error(`created UNS topic type is ${created?.topicType || '<empty>'}, expected ${contract.topicType}`);
    }
    if (
      created?.persistence !== true ||
      Number(created?.enableHistory) !== 1 ||
      created?.mockData !== true ||
      created?.withFlow !== true
    ) {
      throw new Error('created Metric UNS did not retain history, Mock Data and Flow flags');
    }

    const flow = requireEnvelopeSuccess(
      `read generated Flow ${flowID}`,
      await request(baseURL, `/api/core/flows/${flowID}`, { headers: authHeaders })
    );
    if (flow?.status !== 'deployed' || !String(flow?.runtimeFlowId || '')) {
      throw new Error('generated Mock Flow is not deployed to Node-RED');
    }
    const flowRuntime = inspectMockFlow(flow?.flowData, topic, contract.field.name);
    if (flowRuntime.repeatSeconds !== contract.mockIntervalSeconds) {
      throw new Error(`mock Flow interval is ${flowRuntime.repeatSeconds}s, expected ${contract.mockIntervalSeconds}s`);
    }

    const deadline = Date.now() + contract.persistenceTimeoutMs;
    let attempts = 0;
    let detail;
    let dashboard;
    while (Date.now() < deadline) {
      attempts += 1;
      detail = requireEnvelopeSuccess(
        `read Metric UNS ${nodeID}`,
        await request(baseURL, `/api/core/uns/nodes/${nodeID}`, { headers: authHeaders })
      );
      dashboard = requireEnvelopeSuccess(
        `read Metric history ${nodeID}`,
        await request(baseURL, `/api/core/uns/dashboard?nodeId=${nodeID}&limit=10`, { headers: authHeaders })
      );
      if (hasPersistedMetricData(detail, dashboard, contract.field.name)) break;
      await sleep(contract.pollIntervalMs);
    }
    if (!hasPersistedMetricData(detail, dashboard, contract.field.name)) {
      throw new Error(`Metric data was not persisted within ${contract.persistenceTimeoutMs}ms`);
    }
    journey = {
      status: 'passed',
      firstCreateCompletedAtomically: true,
      topicType: created.topicType,
      topic,
      field: contract.field,
      flow: flowRuntime,
      latestValueObserved: true,
      historyRows: Number(dashboard?.total || dashboard?.list?.length || 0),
      persistenceLatencyMs: Date.now() - startedAt,
      pollAttempts: attempts,
    };
  } catch (error) {
    primaryError = error;
  }

  const cleanup = await cleanupFoundationJourney({ baseURL, authHeaders, nodeID, flowID, name });
  if (cleanup.errors.length) {
    const cleanupError = new Error(`foundation regression cleanup failed: ${cleanup.errors.join('; ')}`);
    if (primaryError) {
      primaryError.message = `${primaryError.message}; ${cleanupError.message}`;
    } else {
      primaryError = cleanupError;
    }
  }
  if (primaryError) throw primaryError;
  journey.cleanup = { flowDeleted: cleanup.flowDeleted, unsNodeDeleted: cleanup.unsNodeDeleted };
  return journey;
};

const updateMarker = async (markerPath, result) => {
  if (!markerPath) return;
  const marker = JSON.parse(await readFile(markerPath, 'utf8'));
  marker.phases ||= {};
  const runtimePhase = {
    status: 'passed',
    sourceRevision: marker.sourceRevision,
    sourceDirty: marker.sourceDirty,
    manifestDigest: marker.manifestDigest,
    expectedVersion: result.expectedVersion || '',
    baseUrl: result.baseUrl,
    defaultRoute: result.browser.path,
    requiredResourceKeys: regressionContract.requiredResourceKeys,
    defaultAdminUserRole: result.defaultAdminUserJourney.assignedRole,
    completedAt: new Date().toISOString(),
  };
  marker.phases.runtimeSmoke = runtimePhase;
  if (result.scope === 'foundation') {
    marker.phases.runtimeRegression = {
      ...runtimePhase,
      scope: result.scope,
      metricTopicType: result.foundationJourney.topicType,
      mockIntervalSeconds: result.foundationJourney.flow.repeatSeconds,
      mqttConfigurationName: result.foundationJourney.flow.configurationName,
      firstCreateCompletedAtomically: result.foundationJourney.firstCreateCompletedAtomically,
      historyRows: result.foundationJourney.historyRows,
      persistenceLatencyMs: result.foundationJourney.persistenceLatencyMs,
      cleanup: result.foundationJourney.cleanup,
      ...(result.atomicityFaultInjection ? { atomicityFaultInjection: result.atomicityFaultInjection } : {}),
    };
  } else {
    delete marker.phases.runtimeRegression;
  }
  await writeFile(markerPath, `${JSON.stringify(marker, null, 2)}\n`);
};

export const runRuntimeRegression = async ({
  baseURL,
  username,
  password,
  expectedVersion,
  screenshot,
  marker,
  scope = 'foundation',
  faultInjection = false,
  deployRoot,
}) => {
  if (!['smoke', 'foundation'].includes(scope)) throw new Error(`unsupported runtime regression scope: ${scope}`);
  const checks = [];
  const ready = await request(baseURL, '/readyz');
  if (ready.status !== 200) throw new Error(`readiness failed: HTTP ${ready.status}`);
  checks.push('readyz');

  const system = requireEnvelopeSuccess('system config', await request(baseURL, '/api/core/system/config'));
  if (expectedVersion && String(system?.version || '') !== expectedVersion) {
    throw new Error(`deployed version ${system?.version || '<empty>'} does not match ${expectedVersion}`);
  }
  if (system?.authEnable !== true || system?.loginPath !== regressionContract.loginPath) {
    throw new Error('system config does not expose the retained local-session login contract');
  }
  checks.push('system-config');

  const login = requireEnvelopeSuccess(
    'login',
    await request(baseURL, '/api/core/auth/login', {
      method: 'POST',
      body: JSON.stringify({ username, password }),
    })
  );
  const token = String(login?.token || '');
  if (!token) throw new Error('login returned no session token');
  const missingLoginKeys = missingResourceKeys(login?.resourceKeys || []);
  if (missingLoginKeys.length) throw new Error(`login is missing resource keys: ${missingLoginKeys.join(', ')}`);
  checks.push('login-resource-keys');

  const authHeaders = { authorization: `Bearer ${token}` };
  for (const check of regressionContract.retainedApiChecks) {
    const data = requireEnvelopeSuccess(
      `${check.method} ${check.path}`,
      await request(baseURL, check.path, { method: check.method, headers: authHeaders })
    );
    if (check.path === '/api/core/auth/me') {
      const missingMeKeys = missingResourceKeys(data?.resourceKeys || []);
      if (missingMeKeys.length) throw new Error(`/auth/me is missing resource keys: ${missingMeKeys.join(', ')}`);
    }
    if (check.path.startsWith('/api/core/iam/users')) {
      if (!(data?.list || []).some((item) => item?.userName === 'tier0')) {
        throw new Error('retained user list does not contain the tier0 administrator');
      }
    }
    if (check.path === '/api/core/iam/roles') {
      const admin = (data?.list || []).find((item) => String(item?.code || item?.roleCode).toLowerCase() === 'admin');
      if (!admin || !(admin?.resourceList || []).some((item) => item?.uri === 'resource:iam.user.view')) {
        throw new Error('retained role list does not expose the admin IAM permission set');
      }
    }
    checks.push(`${check.method} ${check.path}`);
  }

  for (let attempt = 0; attempt < 20; attempt += 1) {
    requireEnvelopeSuccess(
      `repeat role list read ${attempt + 1}`,
      await request(baseURL, '/api/core/iam/roles', { headers: authHeaders })
    );
  }
  checks.push('iam-role-list-repeat-read');

  const unsHTML = await request(baseURL, regressionContract.defaultRoute);
  if (unsHTML.status !== 200 || !unsHTML.contentType.includes('text/html')) {
    throw new Error(`${regressionContract.defaultRoute} did not serve the frontend shell`);
  }
  checks.push('frontend-shell');

  await mkdir(dirname(screenshot), { recursive: true });
  const browser = runBrowserRegression({ baseURL, username, password, screenshot });
  checks.push('browser-login-default-route');

  for (const check of regressionContract.disabledApiChecks) {
    const response = await request(baseURL, check.path, { method: check.method, headers: authHeaders });
    if (response.status !== 404) {
      throw new Error(`disabled API ${check.method} ${check.path} returned HTTP ${response.status}, expected 404`);
    }
    checks.push(`404 ${check.method} ${check.path}`);
  }

  const defaultAdminUserJourney = await runDefaultAdminUserJourney({ baseURL, authHeaders });
  checks.push('user-create-default-admin-cleanup');

  let foundationJourney;
  let atomicityFaultInjection;
  if (scope === 'foundation') {
    if (faultInjection) {
      atomicityFaultInjection = await runAtomicityFaultInjection({ baseURL, authHeaders, deployRoot });
      checks.push('metric-create-failure-compensation-identical-retry');
    }
    foundationJourney = await runFoundationJourney({ baseURL, authHeaders });
    checks.push('metric-mock-flow-mqtt-history-cleanup');
  }

  const result = {
    status: 'passed',
    scope,
    baseUrl: baseURL,
    expectedVersion: expectedVersion || '',
    checks,
    browser,
    defaultAdminUserJourney,
    ...(atomicityFaultInjection ? { atomicityFaultInjection } : {}),
    ...(foundationJourney ? { foundationJourney } : {}),
  };
  await updateMarker(marker, result);
  return result;
};

async function main() {
  const options = parseArgs(process.argv.slice(2));
  const password = options.password || process.env.OPEN_SOURCE_REGRESSION_PASSWORD;
  if (!options['base-url'] || !options.username || !password) {
    throw new Error('usage: runtime-regression.mjs --base-url <url> --username <name> [--password <value> or OPEN_SOURCE_REGRESSION_PASSWORD] [--scope <smoke|foundation>] [--expected-version <version>] [--screenshot <path>] [--marker <path>] [--fault-injection true --deploy-root <path>]');
  }
  const baseURL = normalizeBaseURL(options['base-url']);
  const screenshot = resolve(options.screenshot || `/tmp/tier0-edge-open-source-regression-${Date.now()}.png`);
  const result = await runRuntimeRegression({
    baseURL,
    username: options.username,
    password,
    expectedVersion: options['expected-version'],
    screenshot,
    marker: options.marker ? resolve(options.marker) : undefined,
    scope: options.scope || 'foundation',
    faultInjection: options['fault-injection'] === 'true',
    deployRoot: options['deploy-root'] ? resolve(options['deploy-root']) : undefined,
  });
  console.log(JSON.stringify(result, null, 2));
}

if (resolve(process.argv[1] || '') === scriptPath) {
  main().catch((error) => {
    console.error(error.stack || error.message);
    process.exitCode = 1;
  });
}
