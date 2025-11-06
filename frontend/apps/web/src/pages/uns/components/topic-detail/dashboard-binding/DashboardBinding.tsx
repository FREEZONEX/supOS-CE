import { useState, useEffect, type FC, useCallback } from 'react';
import { Dropdown, Button, Space, Divider, Flex, Empty } from 'antd';
import VirtualList from 'rc-virtual-list';
import ComButton from '@/components/com-button';
import { ChartClusterBar, Checkmark, ChevronDown } from '@carbon/icons-react';
import { useTranslate } from '@/hooks';
import ComEllipsis from '@/components/com-ellipsis';
import styles from './index.module.scss';
import classNames from 'classnames';
import { debounce } from 'lodash-es';
import ComInput from '@/components/com-input';
import { getDashboardList } from '@/apis/inter-api';
import usePropsValue from '@/hooks/usePropsValue.ts';
import usePagination from '@/hooks/usePagination.ts';

const CONTAINER_HEIGHT = 200;
const PAGE_SIZE = 20;

const DashboardBinding: FC<{
  isCreated?: boolean;
  onCreated?: () => Promise<void>;
  onBinding?: (item: any) => Promise<void>;
  selectValue?: string;
  setSelectValue?: (value: string) => void;
}> = ({ isCreated, onCreated, onBinding, selectValue, setSelectValue }) => {
  const [open, setOpen] = useState(false);
  const [selected, setSelected] = usePropsValue({
    value: selectValue,
    onChange: setSelectValue,
  });
  const formatMessage = useTranslate();
  const [searchValue, setSearchValue] = useState('');
  const { data, pagination, clearData, setSearchParams, hasMore } = usePagination({
    firstNotGetData: true,
    appendData: true,
    fetchApi: getDashboardList,
    initPageSize: PAGE_SIZE,
  });

  useEffect(() => {
    setSearchValue('');
    if (open) {
      setSearchParams({});
    } else {
      clearData();
    }
  }, [open]);

  // 防抖搜索
  const debouncedSearch = useCallback(
    debounce((value: any) => {
      setSearchParams({ k: value });
    }, 300),
    [setSearchParams]
  );

  const onScroll = (e: React.UIEvent<HTMLElement, UIEvent>) => {
    if (Math.abs(e.currentTarget.scrollHeight - e.currentTarget.scrollTop - CONTAINER_HEIGHT) <= 1) {
      if (hasMore) {
        pagination?.onChange?.(pagination.page + 1);
      }
    }
  };

  return (
    <Space.Compact>
      <ComButton disabled={isCreated} icon={<ChartClusterBar />} onClick={onCreated}>
        {formatMessage('common.create')}
      </ComButton>
      <Dropdown
        open={open}
        onOpenChange={setOpen}
        placement="bottomRight"
        overlayStyle={{
          zIndex: 998,
        }}
        popupRender={() => {
          return (
            <div
              style={{ width: 350, borderRadius: 5, border: '1px solid #E0E0E0', padding: '4px 0' }}
              className={classNames(styles['container'])}
            >
              <ComInput
                allowClear
                placeholder={formatMessage('common.search')}
                variant="borderless"
                value={searchValue}
                onChange={(e) => {
                  const value = e.target.value;
                  setSearchValue(value);
                  debouncedSearch(value);
                }}
                onKeyDown={(e) => {
                  if (e.key === 'Enter') {
                    debouncedSearch(searchValue);
                  }
                }}
                onClear={() => {
                  setSearchParams({});
                }}
              />
              <Divider
                style={{
                  margin: '4px 0',
                  color: '#E0E0E0',
                }}
              />
              {data?.length > 0 ? (
                <VirtualList data={data} height={CONTAINER_HEIGHT} itemHeight={32} itemKey="id" onScroll={onScroll}>
                  {(item: any) => {
                    return (
                      <Flex
                        style={{ height: 32 }}
                        className={classNames(styles['list-item'], selected === item.id && styles.selected)}
                        align="center"
                        onClick={() => {
                          onBinding?.(item).then(() => {
                            setSelected(item.id);
                          });
                        }}
                        gap={8}
                      >
                        <Flex align="center" style={{ flexShrink: 0, minWidth: 20 }}>
                          {selected === item.id ? <Checkmark /> : <span></span>}
                        </Flex>
                        <ComEllipsis style={{ flex: 1 }}>{item.name}</ComEllipsis>
                      </Flex>
                    );
                  }}
                </VirtualList>
              ) : (
                <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} />
              )}
            </div>
          );
        }}
      >
        <Button style={{ padding: 6 }} title={formatMessage('common.changeBinding')}>
          <Flex
            align="center"
            style={{
              height: '100%',
              transform: open ? 'rotate(180deg)' : 'rotate(0deg)',
              transition: 'transform 0.2s ease-in-out',
              display: 'inline-block',
            }}
          >
            <ChevronDown />
          </Flex>
        </Button>
      </Dropdown>
    </Space.Compact>
  );
};

export default DashboardBinding;
