//环境变量-类型提示
interface ImportMetaEnv {
  /** 会使用在copilotkit的地址上 */
  readonly REACT_APP_BASE_URL: string;
  /** 本地开发代理地址 */
  readonly API_PROXY_URL: string;
  /** app的title */
  readonly VITE_APP_TITLE: string;
  readonly REACT_APP_LOCAL_LANG: string;
  //加入更多环境变量...
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}
