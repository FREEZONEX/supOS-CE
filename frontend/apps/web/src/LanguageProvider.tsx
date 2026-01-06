import { ConfigProvider } from 'antd';
import type { ReactNode } from 'react';
import { IntlProvider } from 'react-intl';
import { useI18nStore } from '@/stores/i18n-store.ts';
import { useTranslate } from '@/hooks';

interface PropsTypes {
  config?: any;
  children: ReactNode;
}

const LanguageProvider = (props: PropsTypes) => {
  const { children, config = {} } = props;
  const { lang, langMessages, antMessages } = useI18nStore();
  const formatMessage = useTranslate();

  const customAntdMessages = {
    global: {
      placeholder: formatMessage('common.select'),
      close: formatMessage('common.close'),
    },
    Empty: {
      description: formatMessage('uns.noData'),
    },
    Modal: {
      okText: formatMessage('common.confirm'),
      cancelText: formatMessage('common.cancel'),
      justOkText: formatMessage('common.ok'),
    },
    Popconfirm: {
      cancelText: formatMessage('common.cancel'),
      okText: formatMessage('common.confirm'),
    },
    Pagination: {
      items_per_page: formatMessage('common.items_per_page'),
      jump_to: formatMessage('common.jump_to'),
      jump_to_confirm: formatMessage('common.jump_to_confirm'),
      page: formatMessage('common.page'),
      prev_page: formatMessage('common.prev_page'),
      next_page: formatMessage('common.next_page'),
      prev_5: formatMessage('common.prev_5'),
      next_5: formatMessage('common.next_5'),
      prev_3: formatMessage('common.prev_3'),
      next_3: formatMessage('common.next_3'),
      page_size: formatMessage('common.page_size'),
    },
  };
  return (
    <IntlProvider messages={langMessages} locale={lang} defaultLocale={'en'} onError={(error) => console.error(error)}>
      <ConfigProvider locale={{ ...antMessages, ...customAntdMessages }} {...config}>
        {children}
      </ConfigProvider>
    </IntlProvider>
  );
};

export default LanguageProvider;
