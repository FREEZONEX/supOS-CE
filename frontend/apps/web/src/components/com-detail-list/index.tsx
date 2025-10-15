import type { FC, ReactNode } from 'react';
import { Flex } from 'antd';
import './index.scss';

export interface ComDetailListProps {
  data?: any;
  list?: {
    key: string;
    prefix?: ReactNode;
    hide?: boolean;
    label: ReactNode;
    render?: (item: any, data: any) => ReactNode;
  }[];
}

const ComDetailList: FC<ComDetailListProps> = ({ data, list }) => {
  return (
    <div className="com-detail-list">
      {list
        ?.filter((item) => !item.hide)
        .map?.((l) => {
          return (
            <Flex className="com-detail-list-item" align="center" key={l.key}>
              <Flex className="item-label" align="center" gap={4}>
                {l.prefix && <Flex>{l.prefix}</Flex>}
                <div>{l.label}</div>
              </Flex>
              <div className="item-content">{l.render ? l.render(data?.[l?.key], data) : data?.[l?.key]}</div>
            </Flex>
          );
        })}
    </div>
  );
};

export default ComDetailList;
