/**
 * 国际化store
 * @description 国际化相关
 */
/**
 * 国际化store
 * @description 国际化相关
 */
import type { StoreApi } from 'zustand';
import { shallow } from 'zustand/vanilla/shallow';
import { createIntl, createIntlCache, type IntlShape } from 'react-intl';
import { APP_LANG_MESSAGE, APP_LANG } from '@/common-types/constans.ts';
import { type I18nData, loadMessages, loadAntdLocale } from '@/utils/i18ns';
import localZhCN from '@/locale/zh-CN.json';
import localEnUS from '@/locale/en-US.json';
import { storageOpt } from '@/utils/storage';
import { createWithEqualityFn, type UseBoundStoreWithEqualityFn } from 'zustand/traditional';
import { subscribeWithSelector } from 'zustand/middleware';
import 'dayjs/locale/zh-cn';
import dayjs from 'dayjs';

const loadDayjsLocale = (locale: string) => {
  if (dayjs.locale() === locale) return;
  dayjs.locale(locale);
};

export enum I18nEnum {
  EnUS = 'en-US',
  ZhCN = 'zh-CN',
}

export const defaultLanguage = I18nEnum.EnUS; // 默认语言为英文
const fallbackEnUS = localEnUS as I18nData;
const fallbackZhCN = localZhCN as I18nData;
const localMessages: Record<string, I18nData> = {
  [I18nEnum.EnUS]: fallbackEnUS,
  [I18nEnum.ZhCN]: fallbackZhCN,
};

const formatLocalMessage = (template: string, opt?: Record<string, unknown>) => {
  if (!opt) {
    return template;
  }
  return template.replace(/\{(\w+)\}/g, (match, key) => {
    const value = opt[key];
    return value === undefined || value === null ? match : String(value);
  });
};

const getLocalFallbackMessage = (id: string, lang: string, opt?: Record<string, unknown>) => {
  const template = localMessages[lang]?.[id] || fallbackEnUS[id];
  return template ? formatLocalMessage(template, opt) : undefined;
};

interface LanguageProps {
  hasUsed?: boolean;
  id: number;
  languageCode: string;
  languageName: string;
  languageType: number;
  value: string;
  label: string;
}
const intlCache = createIntlCache();

export type TI18nStore = {
  lang: string;
  langMessages: I18nData;
  antMessages: I18nData;
  prefix?: string;
  intl: IntlShape;
  langList: LanguageProps[];
};

type InitI18nOptions = {
  skipRemoteMessages?: boolean;
};

export const useI18nStore: UseBoundStoreWithEqualityFn<StoreApi<TI18nStore>> = createWithEqualityFn(
  subscribeWithSelector(() => {
    const lang = storageOpt.getOrigin(APP_LANG) || defaultLanguage;
    const messages = { ...(localMessages[lang] || localEnUS), ...(storageOpt.get(APP_LANG_MESSAGE) || {}) };
    return {
      lang,
      langMessages: messages,
      antMessages: {},
      intl: createIntl(
        {
          locale: lang,
          messages,
        },
        intlCache
      ),
      langList: [],
    };
  }),
  shallow
);

// 初始化国际化
export const initI18n = async (lang: string = defaultLanguage, pluginLang: any = {}, options?: InitI18nOptions) => {
  loadDayjsLocale(lang === I18nEnum.ZhCN ? 'zh-cn' : 'en');
  const antMessages = await loadAntdLocale(lang);
  return await loadMessages(lang as I18nEnum, options).then((res: I18nData) => {
    const fallbackMessages = localMessages[lang] || fallbackEnUS;
    const existingMessages = storageOpt.get(APP_LANG_MESSAGE) || {};
    const finallyMsg = { ...fallbackMessages, ...res, ...pluginLang };
    const shouldKeepExistingCache =
      !options?.skipRemoteMessages &&
      Object.keys(res).length < 100 &&
      Object.keys(existingMessages).length > Object.keys(finallyMsg).length;
    const messagesForStore = shouldKeepExistingCache ? { ...existingMessages, ...pluginLang } : finallyMsg;
    if (!shouldKeepExistingCache) {
      storageOpt.set(APP_LANG_MESSAGE, messagesForStore);
    }
    // node-red的语言
    storageOpt.setOrigin('editor-language', lang);
    // emq的语言 兼容prida需求 原来是en 和zh
    storageOpt.setOrigin('language', lang === I18nEnum.EnUS ? 'en-us' : 'zh-cn');
    // chat2db语言
    storageOpt.setOrigin('lang', lang === I18nEnum.EnUS ? 'en-us' : 'zh-cn');
    // 平台语言
    storageOpt.setOrigin(APP_LANG, lang);
    useI18nStore.setState({
      langMessages: messagesForStore,
      lang,
      antMessages,
      intl: createIntl(
        {
          locale: lang,
          messages: messagesForStore,
        },
        intlCache
      ),
    });
  });
};

let intl: IntlShape;
const unsubscribe = useI18nStore.subscribe((state) => {
  intl = state.intl;
});

export const cleanupI18nSubscriptions = () => {
  unsubscribe();
};

// FIXME: 在react组件中请使用useTranslate这个hooks，其他js文件中使用getIntl
export const getIntl = (id: string, opt?: any, defaultMessage?: string, description?: string | object) => {
  const state = useI18nStore.getState();
  const formatted = (intl || state?.intl)?.formatMessage(
    {
      id: id,
      defaultMessage: defaultMessage,
      description: description,
    },
    opt
  );
  if (formatted && formatted !== id) {
    return formatted;
  }
  return getLocalFallbackMessage(id, state?.lang || defaultLanguage, opt) || formatted || defaultMessage || id;
};

/**
 * 合并remote的国际化资源
 * @param messages 要添加的国际化消息对象
 */
export const connectI18nMessage = (messages: I18nData) => {
  const finalMessages: I18nData = {
    ...useI18nStore.getState()?.langMessages,
    ...messages,
  };

  useI18nStore.setState({
    langMessages: finalMessages,
    intl: createIntl(
      {
        locale: useI18nStore.getState()?.lang,
        messages: finalMessages,
      },
      intlCache
    ),
  });
  storageOpt.set(APP_LANG_MESSAGE, finalMessages);
};
