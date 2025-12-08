// request.js

// 基础 URL（根据你的项目修改）
const BASE_URL = "";

// 默认配置
const defaultOptions = {
  headers: {
    "Content-Type": "application/json",
    // 可以在这里加 token，比如：
    // 'Authorization': 'Bearer ' + localStorage.getItem('token')
  },
};

// 获取国际化文本
const getMessage = (key) => {
  const i18nContainer = document.getElementById("i18n-messages");
  if (i18nContainer && i18nContainer.dataset[key]) {
    return i18nContainer.dataset[key];
  }
  return "";
};

/**
 * 封装 fetch 请求
 * @param {string} url       请求地址
 * @param {object} options   fetch 配置项
 * @param {boolean} isJson   是否自动解析 JSON（默认是）
 * @returns {Promise<any>}   返回处理后的响应数据
 */
export const request = async (url, options = {}, isJson = true) => {
  const fullUrl = url.startsWith("http") ? url : BASE_URL + url;

  const config = {
    ...defaultOptions,
    ...options,
    headers: {
      ...defaultOptions.headers,
      ...(options.headers || {}),
    },
  };

  // 如果 body 是 FormData，则删除 Content-Type 头，让浏览器自动设置
  if (options.body instanceof FormData) {
    delete config.headers["Content-Type"];
  }

  try {
    const response = await fetch(fullUrl, config);

    if (isJson) {
      const res = await response.json();
      
      // 处理标准 API 响应格式 { code, data, msg } 或 { code, data, message }
      if (res && typeof res === 'object') {
        const { code, status, data, msg, message } = res;
        const newCode = status || code;
        if (newCode === 200) {
          return data;
        } else if(newCode !== undefined){
          const errorMsg = `${getMessage('requestfailed')}: ${newCode} ${msg || message || res.error || ''}`;
          const error = new Error(errorMsg);
          error.code = newCode
          error.data = data;
          error.response = res;
          return Promise.reject(error);
        }
      }
      
      // 如果 res 不是对象或为 null，直接返回
      return res;
    } else {
      return await response.text();
    }
  } catch (error) {
    console.error('请求出错:', error);
    // 重新抛出错误，这样调用者可以使用 .catch() 捕获它
    throw error;
  }
};