import type { FC, ReactNode } from 'react';
import { Divider, Flex } from 'antd';
import classNames from 'classnames';

const InfoList: FC<{ items: any[]; className?: string; extra?: ReactNode }> = ({ items, className }) => {
  return (
    <div className={classNames(className)} style={{ height: '100%', overflow: 'auto', padding: '20px 1rem' }}>
      {items?.map((item, index) => {
        return (
          <Flex vertical gap={4} key={item.key} style={{ width: '100%' }}>
            <Flex justify={'space-between'} style={{ width: '100%' }}>
              <div>{item.label}</div>
              {item.extra && <div>{item.extra}</div>}
            </Flex>
            <div>{item.children}</div>
            {items?.length !== index + 1 && (
              <Divider
                style={{
                  margin: '20px 0',
                  backgroundColor: '#C6C6C6',
                }}
              />
            )}
          </Flex>
        );
      })}
    </div>
  );
};

export default InfoList;
