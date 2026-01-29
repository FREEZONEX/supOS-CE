import { forwardRef, useEffect, useState } from 'react';
import { Select, type SelectProps } from 'antd';
import { debounce } from 'lodash-es';
import { useTranslate } from '@/hooks';

interface ApiFunction {
  (key?: string): Promise<any[]>;
}

interface ComSelectProps extends SelectProps {
  options?: any[];
  api?: ApiFunction;
  debounceTimeout?: number;
  isRequest?: boolean;
  otherOptions?: any[];
}

const ComSelect = forwardRef<any, ComSelectProps>(
  ({ options, api, debounceTimeout = 500, isRequest, otherOptions, ...restProps }, ref) => {
    const formatMessage = useTranslate();
    const [apiOptions, setApiOptions] = useState<any[]>(otherOptions || []);
    const [hasRequested, setHasRequested] = useState(false);

    const searchData = async (key?: string) => {
      if (!api) return;

      try {
        const res = await api(key);
        setApiOptions(otherOptions ? otherOptions.concat(res) : res);
      } catch (error) {
        console.error('ComSelect API请求失败:', error);
        setApiOptions(otherOptions || []);
      }
    };

    // 监听isRequest变化，当Modal打开时触发请求
    useEffect(() => {
      if (api && isRequest && !hasRequested) {
        searchData();
        setHasRequested(true);
      }
    }, [api, isRequest, hasRequested]);

    // 当isRequest变为false时（Modal关闭），重置请求状态
    useEffect(() => {
      if (!isRequest) {
        setHasRequested(false);
      }
    }, [isRequest]);

    const handleSearch = api
      ? debounce((searchValue: string) => {
          searchData(searchValue);
        }, debounceTimeout)
      : restProps?.onSearch;

    // 当Modal打开时，如果还没有数据，自动触发一次请求
    const handleOpenChangeChange = (open: boolean) => {
      if (open && api && !hasRequested) {
        searchData();
        setHasRequested(true);
      }
      restProps?.onOpenChange?.(open);
    };

    const selectProps = {
      placeholder: formatMessage('common.select'),
      ...restProps,
      options: api ? apiOptions : options,
      ref,
      onSearch: handleSearch,
      onOpenChange: handleOpenChangeChange,
    };

    return <Select {...selectProps} />;
  }
);

ComSelect.displayName = 'ComSelect';

export default ComSelect;
