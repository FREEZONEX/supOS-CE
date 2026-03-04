import { Button, type CheckboxOptionType, Popover, Radio } from 'antd';
import { Filter } from '@carbon/icons-react';
import useTranslate from '../../hooks/useTranslate.ts';
import usePropsValue from '../../hooks/usePropsValue.ts';
import type { FC, Key } from 'react';

export interface ButtonFilterProps {
  value?: Key;
  onChange?: (value: Key) => void;
  options?: CheckboxOptionType[];
  defaultValue?: Key;
}

const ComButtonFilter: FC<ButtonFilterProps> = ({ value, defaultValue, onChange, options }) => {
  const formatMessage = useTranslate();
  const [v, setV] = usePropsValue<Key>({
    value,
    defaultValue,
    onChange,
  });

  const popoverContent = (
    <Radio.Group
      style={{
        display: 'flex',
        flexDirection: 'column',
        gap: 8,
      }}
      value={v}
      onChange={(e) => {
        setV(e.target.value);
      }}
      options={
        options
          ? options
          : [
              { value: 'all', label: formatMessage('common.all'), title: formatMessage('common.all') },
              { value: 'group', label: formatMessage('common.group'), title: formatMessage('common.group') },
              { value: 'file', label: formatMessage('common.flow'), title: formatMessage('common.flow') },
            ]
      }
    />
  );
  return (
    <Popover placement="bottomLeft" title="" content={popoverContent} trigger="hover">
      <Button style={{ padding: 7 }}>
        <Filter />
      </Button>
    </Popover>
  );
};

export default ComButtonFilter;
