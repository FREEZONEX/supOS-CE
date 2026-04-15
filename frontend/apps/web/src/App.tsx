import { useEffect, useMemo } from 'react';
import { App as AntApp } from 'antd';
import { BrowserRouter } from 'react-router';
import { getRoutesDom, RoutesElement } from '@/routers';
import CookieContext from '@/CookieContext';
import themeToken from './theme/theme-token.ts';
import 'shepherd.js/dist/css/shepherd.css';
import './App.css';
import { userLogin } from '@/apis/chat2db';
import { UnsTreeMapProvider } from '@/UnsTreeMapContext';
import {
  APP_TITLE,
  LOGIN_URL,
  MENU_TARGET_PATH,
  OMC_MODEL,
  STORAGE_PATH,
  SUPOS_LANG,
} from '@/common-types/constans.ts';
import LanguageProvider from './LanguageProvider.tsx';
import { queryChat2dbCurUser } from '@/utils/chat2db.ts';
import { checkImageExists, isInIframe } from '@/utils/url-util.ts';
import { fetchBaseStore, useBaseStore } from '@/stores/base';
import { setThemeBySystem, ThemeType, useThemeStore } from '@/stores/theme-store.ts';
import { cleanupI18nSubscriptions, defaultLanguage, initI18n } from './stores/i18n-store.ts';
import Cookies from 'js-cookie';
import { useI18nStore } from '@/stores/i18n-store';
import { CookiesProvider } from 'react-cookie';
import { storageOpt } from '@/utils/storage.ts';

function App() {
  const { systemInfo, loading, routesStatus, currentUserInfo, menuGroup } = useBaseStore((state) => ({
    systemInfo: state.systemInfo,
    loading: state.loading,
    routesStatus: state.routesStatus,
    menuGroup: state.menuGroup,
    currentUserInfo: state.currentUserInfo,
  }));
  const _theme = useThemeStore((state) => state._theme);
  const lang = useI18nStore((state) => state.lang);

  useEffect(() => {
    const isOmc = isInIframe([], 'webview');
    const currentPath = window.location.pathname;
    const skipBootstrap = currentPath === LOGIN_URL || currentPath === '/freeLogin';
    const preferredEnvLang = import.meta.env.REACT_APP_LOCAL_LANG || import.meta.env.REACT_APP_OS_LANG;
    if (isOmc) {
      Cookies.set(OMC_MODEL, '1', {
        expires: 365,
      });
    } else {
      Cookies.remove(OMC_MODEL, { path: '/' });
    }
    if (skipBootstrap) {
      const initPublicPage = async () => {
        const fallbackLang = preferredEnvLang || storageOpt.getOrigin(SUPOS_LANG) || defaultLanguage;
        try {
          useBaseStore.setState({
            systemInfo: {
              ...useBaseStore.getState().systemInfo,
              appTitle: APP_TITLE,
            },
          });
          await initI18n(fallbackLang, {}, { skipRemoteMessages: true });
        } catch (error) {
          console.log(error);
          await initI18n(fallbackLang, {}, { skipRemoteMessages: true });
        } finally {
          useBaseStore.setState({
            loading: false,
          });
        }
      };
      initPublicPage();
    } else {
      // 初始化
      fetchBaseStore(true);
    }
    return () => {
      cleanupI18nSubscriptions();
    };
  }, []);

  useEffect(() => {
    if (systemInfo?.containerMap?.chat2db) {
      // chat2db登录逻辑
      try {
        queryChat2dbCurUser?.()?.then(async (res) => {
          if (!res) {
            // 重新登录
            await userLogin?.();
            await queryChat2dbCurUser?.();
          }
        });
      } catch (e) {
        console.log(e);
      }
    }
  }, [systemInfo?.containerMap]);

  useEffect(() => {
    if (!systemInfo.appTitle) return;
    const loadFavicon = async () => {
      // 浏览器标题
      document.title = `${systemInfo.themeConfig?.general?.browserTitle || systemInfo.appTitle}`;
      const baseUrl = `${STORAGE_PATH}${MENU_TARGET_PATH}/${systemInfo?.themeConfig?.general?.browserIco || `logo-ico.svg`}`;
      const themeExists = await checkImageExists(baseUrl);

      // 统一处理文件类型和路径
      const [type, path] = themeExists ? ['image/svg+xml', baseUrl] : ['image/svg+xml', '/logo-ico.svg'];

      // 统一处理时间戳
      const href = `${path}?v=${Date.now()}`;

      // 查找或创建 link 元素
      let link = document.querySelector<HTMLLinkElement>("link[rel~='icon']");
      if (!link) {
        link = document.createElement('link');
        link.rel = 'icon';
        document.head.append(link);
      }

      // 统一设置属性
      Object.assign(link, { type, href });
    };

    loadFavicon();
  }, [systemInfo, lang]);

  useEffect(() => {
    if (_theme === ThemeType.System) {
      const mediaChange = (event: any) => {
        setThemeBySystem(event.matches);
      };
      const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)');
      mediaQuery.addEventListener('change', mediaChange);
      return () => {
        mediaQuery.removeEventListener('change', mediaChange);
      };
    }
  }, [_theme]);

  const routeDom = useMemo(() => {
    return getRoutesDom({ menuGroup, systemInfo, currentUserInfo });
  }, [menuGroup, systemInfo, currentUserInfo, lang]);

  if (loading) {
    return <div>Loading...</div>;
  }

  if (routesStatus === 401) {
    return <div></div>;
  }

  return (
    <CookiesProvider defaultSetOptions={{ path: '/' }}>
      <CookieContext />
      <LanguageProvider config={{ theme: themeToken }}>
        {/*antd组件库的主题*/}
        <UnsTreeMapProvider>
          <AntApp>
            <BrowserRouter>
              <RoutesElement routeDom={routeDom} />
            </BrowserRouter>
          </AntApp>
        </UnsTreeMapProvider>
      </LanguageProvider>
    </CookiesProvider>
  );
}

export default App;
