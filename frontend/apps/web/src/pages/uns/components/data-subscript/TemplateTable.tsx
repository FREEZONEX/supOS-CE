import ProTable from '@/components/pro-table';
import { usePagination, useTranslate } from '@/hooks';
import { subscribeTemplatePage } from '@/apis/inter-api/uns.ts';
import { useState } from 'react';
import { Button, Flex, Input } from 'antd';
import { SearchOutlined } from '@ant-design/icons';
import styles from './index.module.scss';

const Index = ({
  isSimple = true,
  isFullscreen,
  onNameClick,
}: {
  isSimple?: boolean;
  isFullscreen?: boolean;
  onNameClick: (item: any, type: string) => void;
}) => {
  const formatMessage = useTranslate();
  const { loading, data, setSearchParams, pagination } = usePagination({
    fetchApi: subscribeTemplatePage,
    initPageSize: isSimple ? 5 : 20,
  });
  const [searchValue, setSearchValue] = useState('');

  const columns: any = [
    {
      dataIndex: 'name',
      ellipsis: true,
      title: () => formatMessage('common.name'),
      width: '20%',
      render: (text: string, record: any) => {
        return (
          <Button
            size="small"
            type="link"
            onClick={() => {
              onNameClick(record, 'template');
            }}
          >
            {text}
          </Button>
        );
      },
    },
    {
      dataIndex: 'topic',
      ellipsis: true,
      title: () => formatMessage('uns.topic'),
      width: '40%',
    },
    {
      dataIndex: 'fileCount',
      ellipsis: true,
      title: () => formatMessage('common.fileNum'),
      width: '20%',
    },
    {
      dataIndex: 'subscribeFrequency',
      ellipsis: true,
      title: () => formatMessage('uns.subscriptionFrequency'),
      width: '20%',
    },
  ];
  const onSearch = () => {
    setSearchParams({
      name: searchValue,
    });
  };
  return (
    <>
      {!isSimple && (
        <Flex justify="flex-end" gap={8} align="center" style={{ marginBottom: 16 }}>
          <Input
            value={searchValue}
            onChange={(e) => {
              setSearchValue(e.target.value);
            }}
            onKeyDown={(e) => {
              if (e.key === 'Enter') {
                onSearch();
              }
            }}
            placeholder={formatMessage('common.searchPlaceholderTem')}
            style={{ width: 300 }}
          />
          <Button onClick={onSearch} type="primary" icon={<SearchOutlined />}>
            {formatMessage('common.search')}
          </Button>
        </Flex>
      )}
      <ProTable
        className={styles['custom-table']}
        loading={loading}
        resizeable
        dataSource={data as any}
        columns={columns}
        pagination={
          !isSimple && {
            ...pagination,
            showSizeChanger: true,
          }
        }
        scroll={!isSimple ? { y: isFullscreen ? 'calc(100vh - 280px)' : 400, x: 'max-content' } : undefined}
      />
    </>
  );
};

export default Index;
