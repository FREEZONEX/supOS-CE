(function applyEmbeddedTheme() {
  const DARK_CLASS = 'tier0-flow-theme-dark';
  const LIGHT_CLASS = 'tier0-flow-theme-light';
  const THEME_PARAM = 'theme';
  const THEME_MESSAGE_TYPE = 'tier0FlowThemeChange';

  function normalizeTheme(value) {
    const normalized = String(value || '')
      .trim()
      .toLowerCase();
    return normalized === 'dark' || normalized === 'on' || normalized.includes('dark') ? 'dark' : 'light';
  }

  function themeFromQuery() {
    return normalizeTheme(new URLSearchParams(window.location.search).get(THEME_PARAM));
  }

  function applyTheme(value) {
    const theme = normalizeTheme(value || themeFromQuery());
    const isDark = theme === 'dark';
    const root = document.documentElement;

    root.classList.toggle(DARK_CLASS, isDark);
    root.classList.toggle(LIGHT_CLASS, !isDark);
    root.dataset.tier0FlowTheme = theme;
    root.style.colorScheme = theme;

    if (document.body) {
      document.body.classList.toggle(DARK_CLASS, isDark);
      document.body.classList.toggle(LIGHT_CLASS, !isDark);
      document.body.dataset.tier0FlowTheme = theme;
    }
  }

  function handleThemeMessage(event) {
    if (event.source !== window.parent || !event.data || event.data.type !== THEME_MESSAGE_TYPE) return;
    applyTheme(event.data.theme || (event.data.data && event.data.data.theme));
  }

  applyTheme();
  window.addEventListener('message', handleThemeMessage);
  if (document.readyState === 'loading') {
    document.addEventListener(
      'DOMContentLoaded',
      function onReady() {
        applyTheme();
      },
      { once: true }
    );
  }
  window.addEventListener(
    'pagehide',
    function cleanup() {
      window.removeEventListener('message', handleThemeMessage);
    },
    { once: true }
  );
})();

(function patchFlowPayloadContext(params) {
  const urlParams = new URLSearchParams(window.location.search);
  const context = {};

  Object.keys(params).forEach((paramName) => {
    const paramValue = urlParams.get(paramName);
    if (!paramValue) return;
    context[paramName] = paramValue;
  });

  function hasContext() {
    return Object.keys(context).length > 0;
  }

  function isReadMethod(method) {
    return !method || method.toUpperCase() === 'GET' || method.toUpperCase() === 'HEAD';
  }

  function isFlowPayloadUrl(rawUrl) {
    try {
      const url = new URL(rawUrl, window.location.href);
      if (url.origin !== window.location.origin) return false;
      const path = url.pathname.replace(/\/+$/, '');
      if (path.startsWith('/api/') || path.startsWith('/openapi/')) return false;
      return path === '/flows' || path.endsWith('/flows');
    } catch (e) {
      return false;
    }
  }

  function withFlowContext(rawUrl) {
    if (!hasContext() || !isFlowPayloadUrl(rawUrl)) return rawUrl;
    const url = new URL(rawUrl, window.location.href);
    Object.entries(context).forEach(([key, value]) => {
      if (!url.searchParams.has(key)) {
        url.searchParams.set(key, value);
      }
    });
    return url.toString();
  }

  if (window.fetch) {
    const originalFetch = window.fetch;
    window.fetch = function patchedFetch(input, init) {
      const method = (init && init.method) || (input && input.method);
      if (!isReadMethod(method)) {
        return originalFetch.apply(this, arguments);
      }
      if (typeof Request !== 'undefined' && input instanceof Request) {
        const nextUrl = withFlowContext(input.url);
        if (nextUrl !== input.url) {
          return originalFetch.call(this, new Request(nextUrl, input), init);
        }
        return originalFetch.apply(this, arguments);
      }
      return originalFetch.call(this, withFlowContext(input), init);
    };
  }

  if (window.XMLHttpRequest && XMLHttpRequest.prototype.open) {
    const originalOpen = XMLHttpRequest.prototype.open;
    XMLHttpRequest.prototype.open = function patchedOpen(method, url) {
      if (isReadMethod(method)) {
        arguments[1] = withFlowContext(url);
      }
      return originalOpen.apply(this, arguments);
    };
  }
})({ flowId: 'flowId', rootWorkspaceId: 'rootWorkspaceId' });

(function patchImportDialog() {
  function install() {
    if (!window.$ || !document.body) {
      window.setTimeout(install, 200);
      return;
    }
    const observer = new MutationObserver((mutationsList, rootObserver) => {
      mutationsList.forEach((mutation) => {
        if (mutation.type !== 'attributes' || mutation.attributeName !== 'aria-labelledby') return;
        const element = mutation.target;
        if (element.getAttribute('aria-labelledby') !== 'ui-id-3') return;
        rootObserver.disconnect();
        const visibilityObserver = new MutationObserver(() => {
          if (getComputedStyle(element).display === 'none') return;
          const textarea = $('#red-ui-clipboard-dialog-import-text');
          const normalize = () => {
            const value = textarea.val();
            if (!value || !/^\[[\s\S]*\]$/m.test(value)) return;
            try {
              const data = JSON.parse(value);
              if (data.some((item) => item.type === 'tab')) {
                textarea.val(JSON.stringify(data.filter((item) => item.type !== 'tab'), null, 4));
              }
            } catch (e) {
              // Keep Node-RED default validation message.
            }
          };
          textarea.off('input.flow').on('input.flow', normalize);
          $('#red-ui-clipboard-dialog-import-file-upload')
            .off('change.flow')
            .on('change.flow', () => window.setTimeout(normalize, 100));
          $('#red-ui-clipboard-dialog-import-opt-new').hide();
        });
        visibilityObserver.observe(element, { attributes: true, attributeFilter: ['style'] });
      });
    });
    observer.observe(document.body, { childList: true, subtree: true, attributes: true, attributeFilter: ['aria-labelledby'] });
  }
  install();
})();

(function patchSidebarTooltipCleanup() {
  const TOOLTIP_SELECTOR = '.red-ui-popover.red-ui-popover-size-small';
  const HOVER_ROOT_SELECTORS = ['#red-ui-sidebar', '#red-ui-sidebar-separator', '#red-ui-palette'];
  let activePopoverOwner = null;
  let staleCloseTimer = null;

  function getPopoverOwner(target) {
    if (!window.$ || !target) return null;
    const nodes = $(target).parents().addBack().toArray();
    for (let index = 0; index < nodes.length; index += 1) {
      const popover = $(nodes[index]).data('red-ui-popover');
      if (popover && typeof popover.close === 'function') {
        return nodes[index];
      }
    }
    return null;
  }

  function getHoveredPopoverOwner() {
    try {
      const hovered = Array.from(document.querySelectorAll(':hover')).reverse();
      for (let index = 0; index < hovered.length; index += 1) {
        if (!HOVER_ROOT_SELECTORS.some((selector) => hovered[index].closest(selector))) continue;
        const owner = getPopoverOwner(hovered[index]);
        if (owner) return owner;
      }
    } catch (error) {
      return null;
    }
    return null;
  }

  function removeTooltipElements() {
    if (window.$) {
      $(TOOLTIP_SELECTOR).remove();
      return;
    }
    document.querySelectorAll(TOOLTIP_SELECTOR).forEach((element) => element.remove());
  }

  function hasTooltipElements() {
    return !!document.querySelector(TOOLTIP_SELECTOR);
  }

  function closePopoverData(root, exceptOwner) {
    if (!window.$ || !root) return;
    $(root)
      .find('*')
      .addBack()
      .each(function () {
        if (exceptOwner && this === exceptOwner) return;
        const popover = $(this).data('red-ui-popover');
        if (popover && typeof popover.close === 'function') {
          popover.close(true);
        }
      });
  }

  function closeTooltips(exceptOwner) {
    if (window.$) {
      HOVER_ROOT_SELECTORS.forEach((selector) => {
        closePopoverData(document.querySelector(selector), exceptOwner);
      });
      if (!exceptOwner) {
        removeTooltipElements();
      }
      return;
    }
    removeTooltipElements();
  }

  function scheduleCloseTooltips() {
    activePopoverOwner = null;
    closeTooltips();
    window.setTimeout(closeTooltips, 0);
    window.setTimeout(closeTooltips, 80);
    window.setTimeout(closeTooltips, 250);
  }

  function closeStaleTooltips() {
    staleCloseTimer = null;
    if (!hasTooltipElements()) return;
    const owner = getHoveredPopoverOwner();
    if (owner) {
      activePopoverOwner = owner;
      return;
    }
    activePopoverOwner = null;
    closeTooltips();
  }

  function scheduleCloseStaleTooltips() {
    if (!hasTooltipElements()) return;
    if (staleCloseTimer) return;
    staleCloseTimer = window.setTimeout(closeStaleTooltips, 60);
  }

  function handlePointerOver(event) {
    const owner = getPopoverOwner(event.target);
    if (!owner) {
      scheduleCloseStaleTooltips();
      return;
    }
    if (owner === activePopoverOwner) return;
    activePopoverOwner = owner;
    closeTooltips(owner);
    removeTooltipElements();
  }

  function handlePointerOut(event) {
    const owner = getPopoverOwner(event.target);
    if (!owner) return;
    if (event.relatedTarget && owner.contains(event.relatedTarget)) return;
    if (owner === activePopoverOwner) {
      activePopoverOwner = null;
    }
    scheduleCloseTooltips();
  }

  function isToggleSidebarShortcut(event) {
    const key = String(event.key || '').toLowerCase();
    return (event.ctrlKey || event.metaKey) && !event.altKey && (event.code === 'Space' || key === ' ' || key === 'spacebar');
  }

  function install() {
    if (!window.RED || !RED.events || !document.body) {
      window.setTimeout(install, 200);
      return;
    }

    const handleKeyDown = (event) => {
      if (isToggleSidebarShortcut(event)) {
        scheduleCloseTooltips();
      }
    };
    const sidebarSeparator = document.getElementById('red-ui-sidebar-separator');
    const mainContainer = document.getElementById('red-ui-main-container');
    const observer = mainContainer ? new MutationObserver(scheduleCloseTooltips) : null;
    const tooltipObserver = new MutationObserver(scheduleCloseStaleTooltips);

    RED.events.on('sidebar:resize', scheduleCloseTooltips);
    document.addEventListener('keydown', handleKeyDown, true);
    document.addEventListener('pointerover', handlePointerOver, true);
    document.addEventListener('pointerout', handlePointerOut, true);
    document.addEventListener('pointermove', scheduleCloseStaleTooltips, true);
    document.addEventListener('scroll', scheduleCloseStaleTooltips, true);
    window.addEventListener('blur', scheduleCloseTooltips);
    if (sidebarSeparator) {
      sidebarSeparator.addEventListener('click', scheduleCloseTooltips, true);
    }
    if (observer) {
      observer.observe(mainContainer, { attributes: true, attributeFilter: ['class'] });
    }
    tooltipObserver.observe(document.body, { childList: true, subtree: true, attributes: true, attributeFilter: ['style', 'class'] });

    window.addEventListener(
      'pagehide',
      function cleanup() {
        RED.events.off('sidebar:resize', scheduleCloseTooltips);
        document.removeEventListener('keydown', handleKeyDown, true);
        document.removeEventListener('pointerover', handlePointerOver, true);
        document.removeEventListener('pointerout', handlePointerOut, true);
        document.removeEventListener('pointermove', scheduleCloseStaleTooltips, true);
        document.removeEventListener('scroll', scheduleCloseStaleTooltips, true);
        window.removeEventListener('blur', scheduleCloseTooltips);
        if (staleCloseTimer) {
          window.clearTimeout(staleCloseTimer);
        }
        if (sidebarSeparator) {
          sidebarSeparator.removeEventListener('click', scheduleCloseTooltips, true);
        }
        if (observer) {
          observer.disconnect();
        }
        tooltipObserver.disconnect();
      },
      { once: true }
    );
  }

  install();
})();

(function bridgeParentMessages() {
  const urlParams = new URLSearchParams(window.location.search);
  let rootWorkspaceId = urlParams.get('rootWorkspaceId') || '';
  let forcingEditableWorkspace = false;

  function postCurrentFlows(event) {
    if (!window.RED || !RED.nodes) return;
    ensureEditableWorkspace();
    const completeNodeSet = RED.nodes.createCompleteNodeSet();
    window.parent.postMessage({ type: 'currentEventFlows', data: { flows: currentWorkspaceNodes(completeNodeSet), type: event.data.data } }, '*');
  }

  function nodeType(node) {
    return node && node.type ? String(node.type).trim() : '';
  }

  function nodeId(node) {
    return node && node.id ? String(node.id).trim() : '';
  }

  function nodeWorkspace(node) {
    return node && node.z ? String(node.z).trim() : '';
  }

  function isContextOnly(node) {
    if (!node) return false;
    return node._contextOnly === true || node._contextOnly === 'true';
  }

  function currentWorkspaceNodes(nodes) {
    const rootId = resolveRootWorkspaceId();
    if (!rootId) return nodes.filter((node) => !isContextOnly(node));
    const cleanNodes = nodes.filter((node) => !isContextOnly(node));
    const ownedWorkspaces = ownedWorkspaceIds(cleanNodes, rootId);
    return cleanNodes.filter((node) => {
      const type = nodeType(node);
      const id = nodeId(node);
      const z = nodeWorkspace(node);
      if (type === 'tab') return id === rootId;
      if (type === 'subflow') return ownedWorkspaces.has(id);
      if (z) return ownedWorkspaces.has(z);
      return true;
    });
  }

  function ownedWorkspaceIds(nodes, rootId) {
    const ownedWorkspaces = new Set();
    if (!rootId) return ownedWorkspaces;
    ownedWorkspaces.add(rootId);
    let changed = true;
    while (changed) {
      changed = false;
      nodes.forEach((node) => {
        if (!ownedWorkspaces.has(nodeWorkspace(node))) return;
        const type = nodeType(node);
        if (!type.startsWith('subflow:')) return;
        const subflowId = type.slice('subflow:'.length).trim();
        if (!subflowId || ownedWorkspaces.has(subflowId)) return;
        ownedWorkspaces.add(subflowId);
        changed = true;
      });
    }
    return ownedWorkspaces;
  }

  function contextWorkspaceIds() {
    if (!window.RED || !RED.nodes) return new Set();
    const ids = new Set();
    const completeNodeSet = RED.nodes.createCompleteNodeSet();
    const rootId = resolveRootWorkspaceId();
    const ownedWorkspaces = ownedWorkspaceIds(completeNodeSet, rootId);
    completeNodeSet.forEach((node) => {
      const type = nodeType(node);
      if (type !== 'tab' && type !== 'subflow') return;
      const id = nodeId(node);
      if (!id) return;
      if (isContextOnly(node) || (rootId && !ownedWorkspaces.has(id))) {
        ids.add(id);
      }
    });
    return ids;
  }

  function isContextWorkspaceId(id) {
    id = String(id || '').trim();
    if (!id || !window.RED || !RED.nodes) return false;
    const rootId = resolveRootWorkspaceId();
    if (!rootId || id === rootId) return false;
    const ownedWorkspaces = ownedWorkspaceIds(RED.nodes.createCompleteNodeSet(), rootId);
    return !ownedWorkspaces.has(id);
  }

  function cssEscape(value) {
    if (window.CSS && CSS.escape) return CSS.escape(value);
    return String(value).replace(/["\\#.:,[\]>+~*^$|=]/g, '\\$&');
  }

  function hideContextWorkspaces() {
    const ids = contextWorkspaceIds();
    if (ids.size === 0) return;
    ids.forEach((id) => {
      const escaped = cssEscape(id);
      const selectors = [
        `#red-ui-workspace-tab-${escaped}`,
        `#red-ui-workspace-tabs [data-id="${escaped}"]`,
        `#red-ui-workspace-tabs [data-workspace="${escaped}"]`,
        `#red-ui-workspace-tabs [data-workspace-id="${escaped}"]`,
        `#red-ui-workspace-tabs a[href="#${escaped}"]`,
      ];
      document.querySelectorAll(selectors.join(',')).forEach((element) => {
        const tab = element.closest('li') || element;
        tab.style.display = 'none';
        tab.setAttribute('data-context-only-workspace', 'true');
      });
    });
  }

  function workspaceName(workspace) {
    if (!workspace) return '';
    return workspace.name || workspace.label || '';
  }

  function resolveRootWorkspaceId() {
    if (!window.RED || !RED.nodes) return rootWorkspaceId;
    if (rootWorkspaceId && RED.nodes.workspace(rootWorkspaceId)) {
      return rootWorkspaceId;
    }
    const completeNodeSet = RED.nodes.createCompleteNodeSet();
    const root = completeNodeSet.find((item) => item && item.type === 'tab');
    if (root && root.id) {
      rootWorkspaceId = root.id;
      return rootWorkspaceId;
    }
    const activeId = RED.workspaces && RED.workspaces.active ? RED.workspaces.active() : '';
    if (activeId && RED.nodes.workspace(activeId)) {
      rootWorkspaceId = activeId;
    }
    return rootWorkspaceId;
  }

  function postWorkspaceState() {
    if (!window.RED || !RED.nodes || !RED.workspaces) return;
    const activeId = RED.workspaces.active();
    const activeSubflow = activeId ? RED.nodes.subflow(activeId) : null;
    const activeWorkspace = activeId ? RED.nodes.workspace(activeId) : null;
    const rootId = resolveRootWorkspaceId();
    const rootWorkspace = rootId ? RED.nodes.workspace(rootId) : null;
    window.parent.postMessage(
      {
        type: 'nodeRedWorkspaceState',
        data: {
          activeWorkspaceId: activeId || '',
          activeWorkspaceName: workspaceName(activeSubflow || activeWorkspace),
          activeWorkspaceType: activeSubflow ? 'subflow' : activeWorkspace ? 'tab' : '',
          rootWorkspaceId: rootId || '',
          rootWorkspaceName: workspaceName(rootWorkspace),
          isSubflow: !!activeSubflow,
        },
      },
      '*'
    );
  }

  function showWorkspace(id) {
    if (!window.RED || !RED.workspaces) return;
    const targetId = id || resolveRootWorkspaceId();
    if (!targetId) return;
    forcingEditableWorkspace = true;
    try {
      RED.workspaces.show(isContextWorkspaceId(targetId) ? resolveRootWorkspaceId() : targetId, true);
    } finally {
      forcingEditableWorkspace = false;
    }
    window.setTimeout(postWorkspaceState, 0);
  }

  function installWorkspaceGuard() {
    if (!window.RED || !RED.workspaces || RED.workspaces.__tier0WorkspaceGuardInstalled) return;
    const originalShow = RED.workspaces.show;
    if (typeof originalShow !== 'function') return;
    RED.workspaces.show = function guardedWorkspaceShow(id) {
      const targetId = String(id || '').trim();
      if (!forcingEditableWorkspace && isContextWorkspaceId(targetId)) {
        const rootId = resolveRootWorkspaceId();
        if (rootId && rootId !== targetId) {
          const args = Array.prototype.slice.call(arguments);
          args[0] = rootId;
          if (args.length < 2) args.push(true);
          return originalShow.apply(this, args);
        }
      }
      return originalShow.apply(this, arguments);
    };
    RED.workspaces.__tier0OriginalShow = originalShow;
    RED.workspaces.__tier0WorkspaceGuardInstalled = true;
  }

  function uninstallWorkspaceGuard() {
    if (!window.RED || !RED.workspaces || !RED.workspaces.__tier0WorkspaceGuardInstalled) return;
    if (RED.workspaces.__tier0OriginalShow) {
      RED.workspaces.show = RED.workspaces.__tier0OriginalShow;
    }
    delete RED.workspaces.__tier0OriginalShow;
    delete RED.workspaces.__tier0WorkspaceGuardInstalled;
  }

  function ensureEditableWorkspace() {
    if (!window.RED || !RED.workspaces || !RED.workspaces.active) return;
    const activeId = RED.workspaces.active();
    if (activeId && isContextWorkspaceId(activeId)) {
      showWorkspace(resolveRootWorkspaceId());
    }
  }

  function messageHandler(event) {
    if (!event.data) return;
    if (event.data.type === 'requestEventFlows') {
      postCurrentFlows(event);
    } else if (event.data.type === 'openEventMenu') {
      const id = event.data.data && event.data.data.id;
      const item = id ? document.querySelector('#' + id) : null;
      if (item) item.click();
    } else if (event.data.type === 'updateVersion' && window.RED && RED.nodes) {
      RED.nodes.version(event.data.data);
    } else if (event.data.type === 'nodeRedShowWorkspace') {
      const id = event.data.data && event.data.data.id;
      showWorkspace(id);
    }
  }

  function flowsChange(params) {
    window.parent.postMessage({ type: 'eventFlowsChange', data: params }, '*');
  }

  function setup() {
    if (!window.RED || !RED.events || !RED.workspaces) {
      window.setTimeout(setup, 200);
      return;
    }
    installWorkspaceGuard();
    RED.events.on('flows:change', flowsChange);
    RED.events.on('flows:loaded', hideContextWorkspaces);
    RED.events.on('flows:loaded', ensureEditableWorkspace);
    RED.events.on('flows:loaded', postWorkspaceState);
    RED.events.on('workspace:change', hideContextWorkspaces);
    RED.events.on('workspace:change', ensureEditableWorkspace);
    RED.events.on('workspace:change', postWorkspaceState);
    RED.events.on('subflows:change', hideContextWorkspaces);
    RED.events.on('subflows:change', ensureEditableWorkspace);
    RED.events.on('subflows:change', postWorkspaceState);
    window.addEventListener('message', messageHandler);
    const observer = new MutationObserver(hideContextWorkspaces);
    observer.observe(document.body, { childList: true, subtree: true });
    window.setTimeout(hideContextWorkspaces, 0);
    window.setTimeout(ensureEditableWorkspace, 0);
    window.setTimeout(postWorkspaceState, 0);
    window.addEventListener(
      'pagehide',
      function cleanup() {
        window.removeEventListener('message', messageHandler);
        observer.disconnect();
        uninstallWorkspaceGuard();
        RED.events.off('flows:change', flowsChange);
        RED.events.off('flows:loaded', hideContextWorkspaces);
        RED.events.off('flows:loaded', ensureEditableWorkspace);
        RED.events.off('flows:loaded', postWorkspaceState);
        RED.events.off('workspace:change', hideContextWorkspaces);
        RED.events.off('workspace:change', ensureEditableWorkspace);
        RED.events.off('workspace:change', postWorkspaceState);
        RED.events.off('subflows:change', hideContextWorkspaces);
        RED.events.off('subflows:change', ensureEditableWorkspace);
        RED.events.off('subflows:change', postWorkspaceState);
      },
      { once: true }
    );
  }

  setup();
})();
