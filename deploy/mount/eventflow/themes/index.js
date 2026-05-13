// 解析url设置cookie
(function setCookiesFromUrlParams(params, days) {
  // 从 URL 参数中获取所有指定的参数
  const urlParams = new URLSearchParams(window.location.search);

  // 遍历传入的参数对象
  Object.entries(params).forEach(([paramName, cookieName]) => {
    const paramValue = urlParams.get(paramName);

    // 如果参数存在，设置 cookie
    if (paramValue) {
      let expires = "";
      if (days) {
        const date = new Date();
        date.setTime(date.getTime() + (days * 24 * 60 * 60 * 1000));
        expires = "; expires=" + date.toUTCString();
      }
      document.cookie = cookieName + "=" + (paramValue || "") + expires + "; path=/";
    }
  });
})({ sup_event_flow_id: "sup_event_flow_id", sup_origin_event_flow_id: "sup_origin_event_flow_id" }, 3);

// 监听requestFlows,发送节点数据
(function() {
  // 处理接收到的消息
  function messageHandler(event) {
    // 当监听到请求数据
    if (event.data.type === 'requestEventFlows') {
      // 获取当前编辑器中的完整节点集合
      const completeNodeSet = RED.nodes.createCompleteNodeSet();
      // 将节点集合通过 postMessage 发送给父窗口
      window.parent.postMessage({ type: 'currentEventFlows', data: { flows: completeNodeSet, type: event.data.data } }, '*');
    } else if (event.data.type === 'openEventMenu') {
      event.data.data.id && document.querySelector(`#${event.data.data.id}`).click()
    } else if (event.data.type === 'updateVersion') {
      RED.nodes.version(event.data.data)
    }
  };
  // 通知父节点是否变化
  function flowsChange(params) {
    window.parent.postMessage({ type: 'eventFlowsChange', data: params }, '*');
  }
  // RED内置事件监听
  RED.events.on('flows:change', flowsChange)
  // 监听来自父窗口的消息请求
  window.addEventListener('message', messageHandler);
  const cleanup = () => {
    window.removeEventListener('message', messageHandler);
    RED.events.off('flows:change', flowsChange)
  };
  // 在页面卸载时清理
  window.addEventListener('pagehide', cleanup, { once: true });
})();

// load后的一些操作
$(document).ready(function() {
  const observer = new MutationObserver((mutationsList, observer) => {
    mutationsList.forEach(mutation => {
      // 如果是子元素的添加或移除
      if (mutation.type === 'attributes' && mutation.attributeName === 'aria-labelledby') {
        const element = mutation.target;
        if (element.getAttribute('aria-labelledby') === 'ui-id-3') {
          observer.disconnect();
          // 监听该元素的显示和隐藏
          observeVisibilityChange(element);
        }
      }
    });
  });

  const ariaElementConfig = { childList: true, subtree: true, attributes: true, attributeFilter: ['aria-labelledby'] };
  observer.observe(document.body, ariaElementConfig);

  // 处理导入的value值
  function setValue(value, textarea) {
    if (value && /^\[[\s\S]*\]$/m.test(value)) {
      try {
        const data = JSON.parse(value)
        if (data.some(item => item.type === 'tab')) {
          const filteredData = data.filter(item => item.type !== 'tab')
          textarea.val(JSON.stringify(filteredData, null, 4))
        }
      } catch (e) {}
    }
  }

  function observeVisibilityChange(ariaElement) {
    const visibilityObserver = new MutationObserver(() => {
      const isVisible = getComputedStyle(ariaElement).display !== 'none';  // 判断元素是否可见
      if (isVisible) {
        const textarea = $("#red-ui-clipboard-dialog-import-text")
        // 输入形式
        textarea.on("input", function (event) {
          const value = $(this).val()
          setValue(value, textarea)
        })
        let previousValue = textarea.val();  // 获取初始值
        // 通过导入json文件形式的导入
        $("#red-ui-clipboard-dialog-import-file-upload").on("change", function() {
          const intervalId = setInterval(function () {
            const currentValue = textarea.val();
            if (currentValue !== previousValue) {
              setValue(currentValue, textarea)
              clearInterval(intervalId);
            }
          }, 50);
        })
        // 隐藏新tab元素
        $("#red-ui-clipboard-dialog-import-opt-new").hide();
      }
    });

    // 监听元素的属性变化，尤其是显示隐藏变化
    visibilityObserver.observe(ariaElement, { attributes: true, attributeFilter: ['style'] });
  }
});

// 单页签模式下，Node-RED 自带 unknown node 搜索只能搜索当前已加载节点。
// 这里补一个全局修复入口，交给后端扫描所有 tab/global 后定点删除。
(function() {
  const API = "/service-api/supos/proxy/missing-nodes?flowType=event-flow";
  const SEARCH_TEXTS = new Set(["Search for unknown nodes", "搜索未知节点"]);
  const REPAIR_BUTTON_MARK = "suposGlobalMissingNodeRepair";
  const REPAIR_STYLE_ID = "supos-missing-node-repair-style";
  let repairAttachTimer = null;
  const I18N = {
    en: {
      allFlows: "On all flows",
      close: "Close",
      delete: "Delete",
      deleteFailed: "Delete missing node failed: ",
      deleteSkipped: "Missing node was not deleted: ",
      deleteSuccess: "Deleted missing node: ",
      deleteSummary: "The editor will refresh after deletion",
      dialogTitle: "Global Missing Node Repair",
      fetchFailed: "Fetch missing nodes failed: ",
      flow: "Flow",
      flowNode: "Flow node",
      globalConfig: "Global config",
      missingType: "Missing type",
      name: "Name",
      noMissingNodes: "No missing Node-RED nodes found",
      repair: "Global repair",
      scope: "Scope",
      tabConfig: "Tab config",
      title: "Global missing nodes"
    },
    zh: {
      allFlows: "全部流程",
      close: "关闭",
      delete: "删除",
      deleteFailed: "删除未知节点失败: ",
      deleteSkipped: "未知节点未删除: ",
      deleteSuccess: "已删除未知节点: ",
      deleteSummary: "删除后将刷新当前编辑器",
      dialogTitle: "全局未知节点修复",
      fetchFailed: "获取未知节点失败: ",
      flow: "页签",
      flowNode: "普通节点",
      globalConfig: "全局配置",
      missingType: "缺失类型",
      name: "名称",
      noMissingNodes: "未发现未知 Node-RED 节点",
      repair: "全局修复",
      scope: "范围",
      tabConfig: "页签配置",
      title: "全局未知节点"
    }
  };

  function currentLocale() {
    const searchText = $.trim($("button,a").filter(function() {
      return SEARCH_TEXTS.has($.trim($(this).text()));
    }).first().text());
    if (searchText === "Search for unknown nodes") {
      return "en";
    }
    if (searchText === "搜索未知节点") {
      return "zh";
    }
    const lang = (
      (window.RED && RED.settings && RED.settings.lang) ||
      document.documentElement.lang ||
      navigator.language ||
      ""
    ).toLowerCase();
    return lang.startsWith("zh") ? "zh" : "en";
  }

  function t(key) {
    const locale = currentLocale();
    return (I18N[locale] && I18N[locale][key]) || I18N.en[key] || key;
  }

  function setTextIfChanged(element, text) {
    if (!element.length || $.trim(element.text()) === text) {
      return;
    }
    element.text(text);
  }

  function notify(message, type) {
    if (window.RED && RED.notify) {
      RED.notify(message, type || "info");
    } else {
      console.log(message);
    }
  }

  async function fetchMissingNodes() {
    const res = await fetch(API, { credentials: "same-origin" });
    if (!res.ok) {
      throw new Error(await res.text() || ("HTTP " + res.status));
    }
    const data = await res.json();
    return Array.isArray(data.nodes) ? data.nodes : [];
  }

  async function deleteMissingNode(node) {
    const res = await fetch(API, {
      method: "DELETE",
      credentials: "same-origin",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        flowType: "event-flow",
        id: node.id,
        flowId: node.flowId,
        scope: node.scope
      })
    });
    if (!res.ok) {
      throw new Error(await res.text() || ("HTTP " + res.status));
    }
    return res.json();
  }

  function scopeLabel(scope) {
    if (scope === "globalConfig") {
      return t("globalConfig");
    }
    if (scope === "flowConfig") {
      return t("tabConfig");
    }
    return t("flowNode");
  }

  function addCell(row, className, text, title) {
    return $("<div></div>")
      .addClass(className)
      .attr("title", title || text || "")
      .text(text || "")
      .appendTo(row);
  }

  function ensureRepairStyle() {
    if ($("#" + REPAIR_STYLE_ID).length) {
      return;
    }
    $("<style></style>")
      .attr("id", REPAIR_STYLE_ID)
      .text(`
        .supos-missing-node-repair-button {
          margin-left: 12px !important;
          box-shadow: none !important;
          font-weight: normal !important;
        }
        .supos-missing-node-repair-button:hover {
          border-color: #999 !important;
          background: #f7f7f7 !important;
          color: #555 !important;
        }
        .supos-missing-node-repair-button:focus {
          outline: 2px solid rgba(0, 0, 0, 0.12) !important;
          outline-offset: 2px !important;
        }
      `)
      .appendTo("head");
  }

  function openRepairDialog(nodes) {
    if (!nodes.length) {
      notify(t("noMissingNodes"), "success");
      return;
    }

    const dialog = $("<div class='supos-missing-node-dialog'></div>");
    const style = $("<style></style>").text(`
      .supos-missing-node-dialog { padding: 0 2px; }
      .supos-missing-node-summary {
        display: flex;
        align-items: center;
        justify-content: space-between;
        gap: 12px;
        margin: 2px 0 12px;
        color: #555;
        font-size: 13px;
      }
      .supos-missing-node-count {
        display: inline-flex;
        align-items: center;
        justify-content: center;
        min-width: 28px;
        height: 24px;
        padding: 0 8px;
        border-radius: 12px;
        background: #fff3d6;
        color: #9a6500;
        font-weight: 600;
      }
      .supos-missing-node-table {
        border: 1px solid #d7d7d7;
        border-radius: 4px;
        overflow: hidden;
        background: #fff;
      }
      .supos-missing-node-header,
      .supos-missing-node-row {
        display: grid;
        grid-template-columns: minmax(120px, 1.2fr) 92px minmax(140px, 1.4fr) minmax(170px, 1.6fr) 76px;
        align-items: center;
        column-gap: 12px;
        min-height: 42px;
        padding: 0 12px;
      }
      .supos-missing-node-header {
        background: #f6f6f6;
        color: #555;
        font-weight: 600;
        border-bottom: 1px solid #d7d7d7;
      }
      .supos-missing-node-row:not(:last-child) {
        border-bottom: 1px solid #ececec;
      }
      .supos-missing-node-row:hover {
        background: #fafafa;
      }
      .supos-missing-node-flow,
      .supos-missing-node-name,
      .supos-missing-node-type {
        min-width: 0;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
      }
      .supos-missing-node-type {
        color: #b11b2f;
        font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
        font-size: 12px;
      }
      .supos-missing-node-scope-badge {
        justify-self: start;
        padding: 3px 8px;
        border-radius: 10px;
        background: #eef5e2;
        color: #5f8f00;
        font-size: 12px;
        line-height: 1.2;
        white-space: nowrap;
      }
      .supos-missing-node-action {
        justify-self: end;
      }
      .red-ui-button.supos-missing-node-delete-button {
        min-width: 48px !important;
        height: 26px !important;
        padding: 0 10px !important;
        border: 1px solid #b10f24 !important;
        border-radius: 0 !important;
        background: #b10f24 !important;
        color: #fff !important;
        -webkit-text-fill-color: #fff !important;
        font-size: 13px !important;
        line-height: 24px !important;
      }
      .red-ui-button.supos-missing-node-delete-button:hover {
        border-color: #990d1f !important;
        background: #990d1f !important;
        color: #fff !important;
        -webkit-text-fill-color: #fff !important;
      }
      .red-ui-button.supos-missing-node-delete-button:disabled {
        opacity: 0.55;
      }
      .red-ui-button.supos-missing-node-delete-button:focus {
        color: #fff !important;
        -webkit-text-fill-color: #fff !important;
        outline: none !important;
        box-shadow: none !important;
      }
    `).appendTo(dialog);
    const summary = $("<div class='supos-missing-node-summary'></div>").appendTo(dialog);
    $("<div></div>").text(t("title")).append($("<span class='supos-missing-node-count'></span>").text(nodes.length)).appendTo(summary);
    $("<div></div>").text(t("deleteSummary")).appendTo(summary);

    const table = $("<div class='supos-missing-node-table'></div>").appendTo(dialog);
    const header = $("<div class='supos-missing-node-header'></div>").appendTo(table);
    addCell(header, "", t("flow"));
    addCell(header, "", t("scope"));
    addCell(header, "", t("name"));
    addCell(header, "", t("missingType"));
    addCell(header, "", "");

    const list = $("<div></div>").appendTo(table);
    nodes.forEach((node) => {
      const row = $("<div class='supos-missing-node-row'></div>").appendTo(list);
      addCell(row, "supos-missing-node-flow", node.flowId === "global" ? t("allFlows") : (node.flowLabel || node.flowId || "global"));
      $("<div class='supos-missing-node-scope-badge'></div>").text(scopeLabel(node.scope)).appendTo(row);
      addCell(row, "supos-missing-node-name", node.name || node.id, node.id);
      addCell(row, "supos-missing-node-type", node.type);
      const action = $("<div class='supos-missing-node-action'></div>").appendTo(row);
      $("<button class='red-ui-button red-ui-button-small supos-missing-node-delete-button'></button>").text(t("delete")).appendTo(action).on("click", async function() {
        const button = $(this);
        button.prop("disabled", true);
        try {
          const result = await deleteMissingNode(node);
          if (result.deleted > 0) {
            row.fadeOut(120, function() {
              row.remove();
              $(".supos-missing-node-count", dialog).text(list.children().length);
              if (!list.children().length) {
                dialog.dialog("close");
                setTimeout(() => window.location.reload(), 300);
              }
            });
            notify(t("deleteSuccess") + (node.name || node.id), "success");
          } else {
            notify(t("deleteSkipped") + (node.name || node.id), "warning");
            button.prop("disabled", false);
          }
        } catch (err) {
          notify(t("deleteFailed") + err.message, "error");
          button.prop("disabled", false);
        }
      });
    });

    dialog.dialog({
      title: t("dialogTitle"),
      modal: true,
      width: Math.min(920, window.innerWidth - 48),
      close: function() {
        style.remove();
        dialog.dialog("destroy").remove();
      },
      buttons: [
        {
          text: t("close"),
          click: function() {
            dialog.dialog("close");
          }
        }
      ]
    });
  }

  async function openGlobalRepair() {
    try {
      const nodes = await fetchMissingNodes();
      openRepairDialog(nodes);
    } catch (err) {
      notify(t("fetchFailed") + err.message, "error");
    }
  }

  function attachGlobalRepairButtons() {
    $("button,a").each(function() {
      const searchButton = $(this);
      if (!SEARCH_TEXTS.has($.trim(searchButton.text()))) {
        return;
      }
      if (searchButton.data(REPAIR_BUTTON_MARK)) {
        setTextIfChanged(searchButton.next(".supos-missing-node-repair-button"), t("repair"));
        return;
      }
      ensureRepairStyle();
      searchButton.data(REPAIR_BUTTON_MARK, true);
      const repairButton = $("<button type='button'></button>").text(t("repair"))
        .attr("class", $.trim((searchButton.attr("class") || "") + " supos-missing-node-repair-button"))
        .on("click", function(evt) {
          evt.preventDefault();
          evt.stopPropagation();
          openGlobalRepair();
        });
      repairButton.insertAfter(searchButton);
      const height = searchButton.outerHeight();
      if (height) {
        repairButton.css({
          height: height + "px",
          lineHeight: searchButton.css("line-height"),
          paddingTop: searchButton.css("padding-top"),
          paddingRight: searchButton.css("padding-right"),
          paddingBottom: searchButton.css("padding-bottom"),
          paddingLeft: searchButton.css("padding-left"),
          fontSize: searchButton.css("font-size"),
          borderRadius: searchButton.css("border-radius")
        });
      }
    });
  }

  function scheduleAttachGlobalRepairButtons() {
    if (repairAttachTimer) {
      return;
    }
    repairAttachTimer = setTimeout(function() {
      repairAttachTimer = null;
      attachGlobalRepairButtons();
    }, 120);
  }

  $(document).ready(function() {
    attachGlobalRepairButtons();
    const observer = new MutationObserver(scheduleAttachGlobalRepairButtons);
    observer.observe(document.body, { childList: true, subtree: true });
  });
})();
