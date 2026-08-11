import { Close, Search } from '@carbon/icons-react';
import { type CSSProperties, type FC, useEffect, useMemo, useRef } from 'react';
import ComSelect from '../com-select';
import { useMenuNavigate, usePropsValue, useTranslate } from '@/hooks';
import { Space } from 'antd';
import './index.scss';
import { useBaseStore } from '@/stores/base';

interface SearchSelectProps {
  onSearchCallback?: () => void;
  value?: boolean;
  onChange?: (value: boolean) => void;
  selectStyle?: CSSProperties;
}

const SearchSelect: FC<SearchSelectProps> = ({ onSearchCallback, value, onChange, selectStyle }) => {
  const formatMessage = useTranslate();
  const menuGroup = useBaseStore((state) => state.menuGroup?.filter((f) => !f.subMenu));
  const translatedMenuGroup = useMemo(
    () =>
      menuGroup?.map((item) => ({
        ...item,
        translatedShowName: formatMessage(item.showName || '', undefined, item.showName || ''),
      })),
    [formatMessage, menuGroup]
  );
  const selectRef = useRef<any>(null);

  const handleNavigate = useMenuNavigate();
  const [isIcon, setIcon] = usePropsValue({
    value,
    onChange,
    defaultValue: true,
  });
  useEffect(() => {
    if (!isIcon) {
      selectRef?.current?.focus();
    }
  }, [isIcon]);

  return isIcon ? (
    <div
      className="com-header-search-select"
      style={{ display: 'flex', alignItems: 'center', justifyContent: 'center', height: '100%', width: 48 }}
      onClick={() => {
        setIcon(false);
        onSearchCallback?.();
      }}
    >
      <Search size={20} style={{ color: 'var(--ui-text-color)' }} />
    </div>
  ) : (
    <Space.Compact block className="page-search-compact">
      <ComSelect
        defaultOpen
        ref={selectRef}
        variant="filled"
        options={translatedMenuGroup}
        placeholder={formatMessage('common.searchPage')}
        style={{ width: 180, height: '100%', ...selectStyle }}
        onChange={(_: any, options: any) => {
          console.log(options);
          handleNavigate(options);
          setIcon(true);
        }}
        fieldNames={{
          value: 'id',
          label: 'translatedShowName',
        }}
        filterOption={(input, option) =>
          `${(option?.translatedShowName as string) ?? ''} ${(option?.showName as string) ?? ''}`
            .toLowerCase()
            .includes(input.toLowerCase())
        }
        allowClear
        showSearch
      />
      <div
        className="page-search-close"
        onClick={() => {
          setIcon(true);
        }}
        style={{
          justifyContent: 'center',
          alignItems: 'center',
          display: 'flex',
          width: 40,
          background: 'rgba(0, 0, 0, 0.04)',
        }}
      >
        <Close style={{ cursor: 'pointer' }} />
      </div>
    </Space.Compact>
  );
};

export default SearchSelect;
