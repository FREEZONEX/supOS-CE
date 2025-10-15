// src/hooks/useTranslate.ts
/* 
  多语言切换时在组件中使用
  getIntl () 获取当前语言  使用方法  
  import { useTranslate } from '@/hooks';
  const formatMessage = useTranslate();
  formatMessage('common.chatbot')
*/
import { getIntl, useI18nStore } from '@/stores/i18n-store.ts';
import { useCallback } from 'react';

/**
 * 自定义的 useTranslate 钩子
 * @param prefix remote名称 如：OpenData，不传，则用主项目的翻译
 * @returns
 */
const useTranslate = (prefix?: string) => {
  const lang = useI18nStore((state) => state.lang);
  return useCallback(
    (id: string, opt?: any, defaultMessage?: string, description?: string | object) => {
      if (id) {
        return getIntl(prefix ? `${prefix}.${id}` : id, opt, defaultMessage, description);
      } else {
        return '';
      }
    },
    [lang]
  );
};

export default useTranslate;
