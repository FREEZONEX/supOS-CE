import { App as AntApp, ConfigProvider } from 'antd';
import antdEnUS from 'antd/es/locale/en_US';
import antdZhCN from 'antd/es/locale/zh_CN';
import { useEffect, useState, type ReactNode } from 'react';
import type { AnchorContextPayload } from '@internal/anchor-bridge';
import { getHostContext, onHostContext } from '../bridge';
import { APP_LANG_STORAGE_KEY, getLocale, setLocaleOverride } from '../i18n';
import themeToken from './theme-token';

// 平台上下文 Provider：主题/主色/语言。
// 初始从同源 localStorage 读取（主应用写入），宿主 bridge 的 context 消息到达后以宿主为准；
// storage 事件保留为兜底（无宿主直开、旧版主应用等场景）。
// 与主应用 theme-store 一致的存储键（apps/web/src/common-types/constans.ts）
const APP_THEME_KEY = 'APP_THEME_V2';
const APP_PRIMARY_COLOR_KEY = 'APP_PRIMARY_COLOR';

function readTheme() {
  try {
    return {
      theme: window.localStorage.getItem(APP_THEME_KEY) || 'light',
      primaryColor: window.localStorage.getItem(APP_PRIMARY_COLOR_KEY) || 'chartreuse',
    };
  } catch {
    return { theme: 'light', primaryColor: 'chartreuse' };
  }
}

// 与 apps/web theme-store 的 setThemeRoot 保持一致
function applyThemeRoot(theme: string, primaryColor: string) {
  const root = document.documentElement;
  switch (`${theme}-${primaryColor}`) {
    case 'dark-blue':
      root.classList.remove('chartreuse', 'chartreuseDark');
      root.classList.add('dark');
      break;
    case 'light-blue':
      root.classList.remove('dark', 'chartreuse', 'chartreuseDark');
      break;
    case 'dark-chartreuse':
      root.classList.add('chartreuse', 'chartreuseDark', 'dark');
      break;
    case 'light-chartreuse':
      root.classList.remove('dark', 'chartreuseDark');
      root.classList.add('chartreuse');
      break;
    default:
      root.classList.remove('dark', 'chartreuse', 'chartreuseDark');
      break;
  }
}

export function ThemeBridge({ children }: { children: ReactNode }) {
  const [{ theme, primaryColor }, setThemeState] = useState(readTheme);
  const [lang, setLang] = useState(getLocale);

  useEffect(() => {
    applyThemeRoot(theme, primaryColor);
  }, [theme, primaryColor]);

  // 宿主 bridge：context 消息（初始 + 主题/主色/语言变更）以宿主为准实时生效；
  // 消息可能早于本组件挂载到达，先消费缓存的最近一次
  useEffect(() => {
    const applyContext = (context: AnchorContextPayload) => {
      setThemeState({ theme: context.theme, primaryColor: context.primaryColor });
      setLocaleOverride(context.language);
      setLang(getLocale());
    };
    const cached = getHostContext();
    if (cached) applyContext(cached);
    const unsubscribe = onHostContext(applyContext);
    return () => {
      unsubscribe();
    };
  }, []);

  // 兜底：同源 iframe 下主应用写 localStorage 会触发本文档 storage 事件（旧版主应用/独立页签场景）
  useEffect(() => {
    const onStorage = (event: StorageEvent) => {
      if (!event.key || event.key === APP_THEME_KEY || event.key === APP_PRIMARY_COLOR_KEY) {
        setThemeState(readTheme());
      }
      if (!event.key || event.key === APP_LANG_STORAGE_KEY) {
        setLang(getLocale());
      }
    };
    window.addEventListener('storage', onStorage);
    return () => window.removeEventListener('storage', onStorage);
  }, []);

  return (
    <ConfigProvider theme={themeToken} locale={lang === 'zh-CN' ? antdZhCN : antdEnUS}>
      {/* 语言切换重挂子树：业务文案经 t() 渲染，需要整树重渲染才能生效 */}
      <AntApp key={lang} style={{ height: '100%' }}>
        {children}
      </AntApp>
    </ConfigProvider>
  );
}
