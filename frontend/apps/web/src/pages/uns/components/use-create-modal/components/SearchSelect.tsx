import { useMemo, useState, type FC, useEffect } from 'react';
import { Spin, Select, Divider, Button } from 'antd';
import type { SelectProps } from 'antd';
import { debounce } from 'lodash-es';
import { searchTreeData } from '@/apis/inter-api/uns';
import { useTranslate } from '@/hooks';

export interface DebounceSelectProps extends Omit<SelectProps, 'onChange'> {
  debounceTimeout?: number;
  selectAll?: (value: string[]) => void;
  onChange?: (e: any) => void;
  apiParams?: { type: number; [key: string]: string | number | boolean };
  disabledIds?: string[];
  onFirstBack?: (data: any) => void;
}

const DebounceSelect: FC<DebounceSelectProps> = ({
  value,
  onChange,
  debounceTimeout = 500,
  selectAll,
  apiParams,
  labelInValue = false,
  mode,
  disabledIds,
  onFirstBack,
  ...rest
}) => {
  const [fetching, setFetching] = useState(false);
  const [options, setOptions] = useState([]);
  const formatMessage = useTranslate();

  const { type = 3 } = apiParams || {};

  const debounceFetcher = useMemo(() => {
    const loadOptions = (value: string) => {
      setOptions([]);
      setFetching(true);
      searchData(value);
    };
    return debounce(loadOptions, debounceTimeout);
  }, [debounceTimeout]);

  const searchData = (key?: string) => {
    const params: any = { pageNo: 1, pageSize: 100, type, ...apiParams };
    if (key) params.k = key;
    return searchTreeData(params)
      .then((res: any) => {
        res.forEach((e: any) => {
          e.disabled = disabledIds?.includes(e.id);
        });
        setOptions(res);
        setFetching(false);
        return res;
      })
      .catch((err) => {
        setFetching(false);
        console.log(err);
      });
  };

  useEffect(() => {
    searchData?.().then((data: any) => {
      onFirstBack?.(data);
    });
  }, []);

  const _onChange = (e: any) => {
    onChange?.(
      labelInValue
        ? mode
          ? (e?.map((item: any) => ({ ...item, option: options.find((i: any) => i.id === item.value) })) ?? [])
          : e
            ? { ...e, option: options.find((i: any) => i.id === e?.value) }
            : undefined
        : e
    );
  };

  return (
    <Select
      allowClear
      showSearch
      filterOption={false}
      fieldNames={{ label: 'path', value: 'id' }}
      notFoundContent={fetching ? <Spin size="small" /> : formatMessage('uns.noData')}
      popupRender={(menu) => (
        <>
          {menu}
          {options.length > 0 && ((selectAll && mode) || options.length > 99) && (
            <>
              <Divider style={{ margin: '4px 0', borderColor: '#c6c6c6' }} />
              {selectAll && mode && (
                <div style={{ textAlign: 'center' }}>
                  <Button
                    color="default"
                    variant="filled"
                    onClick={() => {
                      selectAll(options);
                    }}
                    size="small"
                    style={{ backgroundColor: 'var(--supos-uns-button-color)' }}
                  >
                    {formatMessage('uns.select100Items')}
                  </Button>
                </div>
              )}
              {options.length > 99 && (
                <div style={{ textAlign: 'center' }}>{formatMessage('uns.forMoreInformationPleaseSearch')}</div>
              )}
            </>
          )}
        </>
      )}
      {...rest}
      onSearch={debounceFetcher}
      options={options}
      value={value}
      onChange={_onChange}
      onFocus={() => searchData()}
      labelInValue={labelInValue}
      mode={mode}
    />
  );
};

export default DebounceSelect;
