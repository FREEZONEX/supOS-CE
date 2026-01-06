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
import { SUPOS_LANG_MESSAGE, SUPOS_LANG } from '@/common-types/constans.ts';
import { type I18nData, loadMessages, loadAntdLocale } from '@/utils/i18ns';
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

export const useI18nStore: UseBoundStoreWithEqualityFn<StoreApi<TI18nStore>> = createWithEqualityFn(
  subscribeWithSelector(() => {
    const lang = storageOpt.getOrigin(SUPOS_LANG) || defaultLanguage;
    const messages = storageOpt.get(SUPOS_LANG_MESSAGE) || {};
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
export const initI18n = async (lang: string = defaultLanguage, pluginLang: any = {}) => {
  loadDayjsLocale(lang === I18nEnum.ZhCN ? 'zh-cn' : 'en');
  const antMessages = await loadAntdLocale(lang);
  return await loadMessages(lang as I18nEnum).then((res: I18nData) => {
    const finallyMsg = { ...res, ...pluginLang };
    storageOpt.set(SUPOS_LANG_MESSAGE, finallyMsg);
    // node-red的语言
    storageOpt.setOrigin('editor-language', lang);
    // emq的语言 兼容prida需求 原来是en 和zh
    storageOpt.setOrigin('language', lang === I18nEnum.EnUS ? 'en-us' : 'zh-cn');
    // chat2db语言
    storageOpt.setOrigin('lang', lang === I18nEnum.EnUS ? 'en-us' : 'zh-cn');
    // supos语言
    storageOpt.setOrigin(SUPOS_LANG, lang);
    useI18nStore.setState({
      langMessages: finallyMsg,
      lang,
      antMessages,
      intl: createIntl(
        {
          locale: lang,
          messages: finallyMsg,
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
  return (intl || useI18nStore.getState()?.intl)?.formatMessage(
    {
      id: id,
      defaultMessage: defaultMessage,
      description: description,
    },
    opt
  );
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
  storageOpt.set(SUPOS_LANG_MESSAGE, finalMessages);
};
