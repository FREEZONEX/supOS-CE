import React from 'react';
import { List, Skeleton } from 'antd';
import InfiniteScroll from 'react-infinite-scroll-component';
import styles from './index.module.scss';
import classNames from 'classnames';

interface PropsTypes {
  rowKey?: string;
  rowLabel?: string;
  className?: string;
  dataSource: any[];
  selectedKeys: any[];
  hasMore: boolean;
  renderItem?: (arg0: any) => React.ReactNode;
  extra?: (arg0: any) => React.ReactNode;
  onLoadMoreData: () => void;
  onSelect: (arg0: any) => void;
}

const InfiniteScrollList: React.FC<PropsTypes> = (props) => {
  const {
    rowKey = 'code',
    rowLabel = 'title',
    className,
    dataSource,
    renderItem,
    extra,
    selectedKeys,
    hasMore,
    onLoadMoreData,
    onSelect,
  } = props;

  return (
    <div className={classNames(styles.container, className)}>
      <InfiniteScroll
        dataLength={dataSource.length}
        next={onLoadMoreData}
        hasMore={hasMore}
        loader={<Skeleton avatar paragraph={{ rows: 1 }} active />}
        height="100%"
      >
        <List
          dataSource={dataSource}
          renderItem={(item: any) => (
            <List.Item
              className={selectedKeys.includes(item[rowKey]) ? styles.active : ''}
              key={item[rowKey]}
              onClick={() => onSelect(item)}
            >
              {renderItem?.(item) || (
                <div className={styles.title} title={item[rowLabel]}>
                  {item[rowLabel]}
                </div>
              )}
              {extra && <div className={`${styles.extra} infinite-scroll-list-extra`}>{extra(item)}</div>}
            </List.Item>
          )}
        />
      </InfiniteScroll>
    </div>
  );
};

export default InfiniteScrollList;
