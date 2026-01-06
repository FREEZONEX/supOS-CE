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
})({ sup_flow_id: "sup_flow_id", sup_origin_flow_id: "sup_origin_flow_id" }, 3);

// 监听requestFlows,发送节点数据
(function() {
  // 处理接收到的消息
  function messageHandler(event) {
    // 当监听到请求数据
    if (event.data.type === 'requestFlows') {
      // 获取当前编辑器中的完整节点集合
      const completeNodeSet = RED.nodes.createCompleteNodeSet();
      // 将节点集合通过 postMessage 发送给父窗口
      window.parent.postMessage({ type: 'currentFlows', data: { flows: completeNodeSet, type: event.data.data } }, '*');
    } else if (event.data.type === 'openMenu') {
      event.data.data.id && document.querySelector(`#${event.data.data.id}`).click()
    } else if (event.data.type === 'updateVersion') {
      RED.nodes.version(event.data.data)
    }
  };
  // 通知父节点是否变化
  function flowsChange(params) {
    window.parent.postMessage({ type: 'flowsChange', data: params }, '*');
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
  window.addEventListener('unload', cleanup);
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
