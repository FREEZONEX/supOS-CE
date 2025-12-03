import dotenv from 'dotenv';
import colors from 'picocolors';

export interface DevInfo {
  BASE_URL: string;
  API_PROXY_URL: string;
  SINGLE_API_PROXY_URL: string;
  SINGLE_API_PROXY_LIST: string;
  VITE_ASSET_PREFIX: string;
  VITE_REMOTE_PREFIX: string;
  VITE_ENABLE_LOCAL_REMOTE: string;
}

/**
 * @description 开发代理配置
 * */
export const parseConfig = (config: any) => {
  const newConfig: any = {};
  Object.keys(config).forEach((key) => {
    const value = config[key];
    newConfig[key] = value === 'true' ? true : value === 'false' ? false : value;
  });
  return newConfig;
};

export const proxyList = [
  'inter-api',
  'gateway',
  'copilotkit',
  'chat2db/api',
  'minio/inter/supos',
  'files/system/resource',
];

export const getProxy = (
  baseUrl: string = 'http://office.unibutton.com:11488',
  singleList: string,
  singleUrl: string = 'http://office.unibutton.com:11488'
) => {
  const proxyConfig: any = {};

  if (singleList) {
    singleList?.split(',')?.forEach?.((name) => {
      proxyConfig[`/${name}`] = {
        target: singleUrl,
        changeOrigin: true,
        ws: true,
      };
    });
  }

  proxyList.forEach((name) => {
    proxyConfig[`/${name}`] = {
      target: baseUrl,
      changeOrigin: true,
      ws: true,
      // vpn代理
      // agent: new HttpsProxyAgent('http://127.0.0.1:7897'),
    };
  });
  // 给iframe加个代理
  proxyConfig['/iframe'] = {
    target: baseUrl,
    changeOrigin: true,
    rewrite: (path: any) => path.replace(/^\/iframe/, ''),
  };
  // 给chat2db加个代理
  proxyConfig['/chat2db/home/'] = {
    target: baseUrl,
    changeOrigin: true,
  };
  return proxyConfig;
};

// == 开发信息
export const getDevInfo = (): DevInfo => {
  const result = dotenv.config({
    path: ['.env', '.env.local'],
    quiet: true,
  });
  const envConfig = parseConfig(result.parsed || {});
  return envConfig as DevInfo;
};

export const logDevInfo = (info: DevInfo) => {
  const isProdCli = process.env.NODE_ENV === 'production';
  if (isProdCli) return;
  const {
    API_PROXY_URL,
    SINGLE_API_PROXY_URL,
    SINGLE_API_PROXY_LIST,
    VITE_ASSET_PREFIX,
    VITE_REMOTE_PREFIX,
    VITE_ENABLE_LOCAL_REMOTE,
  } = info;
  console.log('---------- 开发信息 ----------');
  console.log(colors.gray('接口代理'), API_PROXY_URL, '\n');
  console.log(colors.gray('特殊接口代理'), SINGLE_API_PROXY_URL, '\n');
  console.log(colors.gray('特殊接口List'), SINGLE_API_PROXY_LIST, '\n');
  console.log(colors.gray('host-dev地址'), VITE_ASSET_PREFIX, '\n');
  console.log(colors.gray('remote-dev地址'), VITE_REMOTE_PREFIX, '\n');
  console.log(colors.gray('是否启用本地模块联邦'), VITE_ENABLE_LOCAL_REMOTE, '\n');

  if (!VITE_ASSET_PREFIX) {
    console.error(colors.red('VITE_ASSET_PREFIX 未配置'));
  }
};

export const logBuildTime = () =>
  new Date().toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
    hour12: true, // 24小时制
  });
