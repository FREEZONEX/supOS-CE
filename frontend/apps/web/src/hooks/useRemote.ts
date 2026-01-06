import { createElement, useEffect, useRef, useState } from 'react';
import { loadRemote, registerRemotes, registerPlugins } from '@module-federation/enhanced/runtime';
import type { PageProps } from '@/common-types';
import { useI18nStore, connectI18nMessage } from '@/stores/i18n-store';
import { getRemotesInfo, preloadPluginLang } from '@/utils';
import { useBaseStore } from '@/stores/base';

const useRemote = ({
  name,
  moduleName,
  location,
}: {
  name: string;
  moduleName?: string;
  location?: PageProps['location'];
}) => {
  const lang = useI18nStore((state) => state.lang);
  const pluginList = useBaseStore((state) => state.pluginList);
  const initStateRef = useRef<boolean>(false);
  const [Module, setModule] = useState<any>(() => () => {
    return createElement('span');
  });
  const [errorMsg, setErrorMsg] = useState<any>('');
  useEffect(() => {
    if (!location?.pathname) return;
    if (!name) return;
    if (location?.pathname !== (moduleName ? `${name}/${moduleName}` : name)) return;
    loadCrmPlugins().then((module) => {
      if (module) {
        setModule(() => module);
      }
    });
  }, [name, moduleName, location?.pathname]);

  const registerRemote = async () => {
    if (initStateRef.current) return;
    try {
      registerPlugins([
        {
          name: 'custom-plugin',
          beforeInit(args) {
            return args;
          },
          init(args) {
            console.warn('init: ', args);
            return args;
          },
          beforeLoadShare(args) {
            console.warn('beforeLoadShare: ', args);
            return args;
          },
          onLoad(args) {
            console.warn('onLoad: ', args);
            return args;
          },
          errorLoadRemote(args) {
            // if (mfId.current === args.id) {
            //   setErrorMsg(args?.error ? args?.error.toString() : '');
            // }
            console.warn('errorLoadRemote', args);
            return args;
          },
        },
      ]);
      registerRemotes([getRemotesInfo({ name })], { force: true });
      initStateRef.current = true;
      return true;
    } catch (e) {
      setErrorMsg(e ? e.toString() : '');
    }
  };

  const loadCrmPlugins = async () => {
    try {
      await registerRemote();
      const remoteModule: any = await loadRemote(`supos-ce${name}/${moduleName || 'index'}`);
      setErrorMsg(remoteModule?.error ? remoteModule?.error.toString() : '');
      try {
        // 再加载一次国际化
        const newMessages = await preloadPluginLang(
          [{ name, backendName: pluginList?.find((f) => '/' + f.plugInfoYml?.route?.name === name)?.name }],
          lang
        );
        connectI18nMessage(newMessages);
      } catch (e) {
        console.log('国际化 error', e);
      }
      return remoteModule?.default || remoteModule;
    } catch (e) {
      setErrorMsg(e ? e.toString() : '');
    }
  };
  return { Module, reLoadRemote: loadCrmPlugins, errorMsg };
};

export default useRemote;
