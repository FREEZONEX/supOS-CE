import type { CSSProperties, FC } from 'react';
import { Flex, Tag } from 'antd';

const FieldPanel: FC<{
  onClick?: (item: any) => void;
  fieldList?: { label: string; value: string }[];
  style?: CSSProperties;
  tooltip?: React.ReactNode | false;
}> = ({ onClick, fieldList = [], style, tooltip }) => {
  return (
    <Flex wrap gap={'10px 0'} style={style} align="center">
      {fieldList?.map((item) => (
        <Tag
          key={item.value}
          style={{ cursor: 'pointer', userSelect: 'none' }}
          onClick={() => {
            onClick?.(item);
          }}
        >
          {item.label}
        </Tag>
      ))}
      {tooltip}
    </Flex>
  );
};

export default FieldPanel;
