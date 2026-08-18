import { useState, useEffect, type FC, useCallback } from 'react';
import ComEmpty from '@/components/com-empty';
import { Dropdown, Button, Divider, Flex } from 'antd';
import VirtualList from 'rc-virtual-list';
import { Checkmark, ChevronDown } from '@/components/lucide-icon/carbon';
import { useTranslate } from '@/hooks';
import ComEllipsis from '@/components/com-ellipsis';
import styles from './index.module.scss';
import classNames from 'classnames';
import { debounce } from 'lodash-es';
import ComInput from '@/components/com-input';
import usePropsValue from '@/hooks/usePropsValue.ts';
import usePagination from '@/hooks/usePagination.ts';

const CONTAINER_HEIGHT = 200;
const PAGE_SIZE = 20;

const Binding: FC<{
  onBinding?: (item: any) => Promise<void>;
  selectValue?: string;
  setSelectValue?: (value: string) => void;
  api: any;
}> = ({ onBinding, selectValue, setSelectValue, api }) => {
  const [open, setOpen] = useState(false);
  const [selected, setSelected] = usePropsValue({
    value: selectValue,
    onChange: setSelectValue,
  });
  const formatMessage = useTranslate();
  const [searchValue, setSearchValue] = useState('');
  const [bindingId, setBindingId] = useState<string>();
  const { data, pagination, clearData, setSearchParams, hasMore } = usePagination({
    firstNotGetData: true,
    appendData: true,
    fetchApi: api,
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

  const handleBinding = useCallback(
    async (item: any) => {
      const itemId = String(item.id);
      if (!onBinding || bindingId) return;
      setBindingId(itemId);
      try {
        await onBinding(item);
        setSelected(itemId);
        setOpen(false);
      } finally {
        setBindingId(undefined);
      }
    },
    [bindingId, onBinding, setSelected]
  );

  return (
    <Dropdown
      open={open}
      onOpenChange={setOpen}
      trigger={['click']}
      placement="bottomRight"
      getPopupContainer={() => document.body}
      overlayStyle={{
        zIndex: 998,
      }}
      popupRender={() => {
        return (
          <div
            style={{ width: 350, borderRadius: 5, border: '1px solid var(--ui-line-color)', padding: '4px 0' }}
            className={classNames(styles['container'])}
            onPointerDown={(event) => event.stopPropagation()}
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
                      className={classNames(
                        styles['list-item'],
                        String(selected) === String(item.id) && styles.selected
                      )}
                      align="center"
                      role="option"
                      tabIndex={0}
                      aria-selected={String(selected) === String(item.id)}
                      aria-busy={bindingId === String(item.id)}
                      onPointerDown={(event) => {
                        if (event.button !== 0) return;
                        event.preventDefault();
                        event.stopPropagation();
                        void handleBinding(item);
                      }}
                      onKeyDown={(event) => {
                        if (event.key !== 'Enter' && event.key !== ' ') return;
                        event.preventDefault();
                        void handleBinding(item);
                      }}
                      gap={8}
                    >
                      <Flex align="center" style={{ flexShrink: 0, minWidth: 20 }}>
                        {String(selected) === String(item.id) ? <Checkmark /> : <span></span>}
                      </Flex>
                      <ComEllipsis style={{ flex: 1, color: 'var(--ui-text-color)' }}>
                        {item.name || item.flowName}
                      </ComEllipsis>
                    </Flex>
                  );
                }}
              </VirtualList>
            ) : (
              <ComEmpty />
            )}
          </div>
        );
      }}
    >
      <Button
        size="small"
        color="default"
        variant="text"
        style={{ color: 'var(--ui-text-color)', padding: '0 4px' }}
        title={formatMessage('common.changeBinding')}
      >
        <ChevronDown
          size={20}
          aria-hidden
          style={{
            transform: open ? 'rotate(180deg)' : 'rotate(0deg)',
            transition: 'transform 0.2s ease-in-out',
            color: 'var(--ui-text-color)',
          }}
        />
      </Button>
    </Dropdown>
  );
};

export default Binding;
