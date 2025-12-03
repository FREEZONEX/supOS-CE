import { Close, Search } from '@carbon/icons-react';
import { type CSSProperties, type FC, useEffect, useRef } from 'react';
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
      <Search size={20} style={{ color: 'var(--supos-text-color)' }} />
    </div>
  ) : (
    <Space.Compact block>
      <ComSelect
        defaultOpen
        ref={selectRef}
        variant="filled"
        options={menuGroup}
        placeholder={formatMessage('common.searchPage')}
        style={{ width: 180, height: '100%', ...selectStyle }}
        onChange={(_: any, options: any) => {
          console.log(options);
          handleNavigate(options);
          setIcon(true);
        }}
        fieldNames={{
          value: 'id',
          label: 'showName',
        }}
        filterOption={(input, option) =>
          ((option?.showName as string) ?? '').toLowerCase().includes(input.toLowerCase())
        }
        allowClear
        showSearch
      />
      <div
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
