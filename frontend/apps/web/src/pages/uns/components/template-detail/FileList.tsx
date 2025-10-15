import { type FC, useEffect, useState } from 'react';
import { Document, Folder } from '@carbon/icons-react';
import { Flex } from 'antd';
import { useTranslate } from '@/hooks';
import { pageListUnsByTemplate } from '@/apis/inter-api/uns.ts';
import ProTable from '@/components/pro-table';

interface FileListProps {
  templateId: string;
}

const initPagination = {
  pageNo: 1,
  pageSize: 10,
};
const FileList: FC<FileListProps> = ({ templateId }) => {
  const [list, setList] = useState([]);
  const [pagination, setPagination] = useState({
    ...initPagination,
    total: 0,
  });
  const formatMessage = useTranslate();

  const getList = (templateId: string, params?: any) => {
    if (!templateId) return;
    pageListUnsByTemplate({ templateId, pageNo: pagination.pageNo, pageSize: pagination.pageSize, ...params }).then(
      (res) => {
        const { code, data, pageNo = 1, pageSize = 10, total = 0 } = res;
        if (code === 0 || code === 200) {
          setList(data || []);
          setPagination({
            ...pagination,
            pageNo,
            pageSize,
            total,
          });
        }
      }
    );
  };

  useEffect(() => {
    if (templateId) {
      getList(templateId, initPagination);
    }
  }, [templateId]);

  return (
    <ProTable
      rowHoverable={false}
      bordered
      hiddenEmpty
      columns={[
        {
          title: formatMessage('common.name'),
          dataIndex: 'name',
          width: '25%',
          render: (text, record) => {
            const StartIcon = record.pathType === 0 ? Folder : Document;
            return (
              <Flex>
                <Flex style={{ height: '24px' }} align="center">
                  <StartIcon style={{ marginRight: 5, flexShrink: 0 }} />
                </Flex>
                <span style={{ wordBreak: 'break-word' }}>{text}</span>
              </Flex>
            );
          },
        },
        {
          title: formatMessage('uns.position'),
          dataIndex: 'path',
          width: '75%',
          render: (text) => (
            <span style={{ color: 'var(--supos-table-first-color)', wordBreak: 'break-word' }}>{text}</span>
          ),
        },
      ]}
      dataSource={list || []}
      rowKey="unsId"
      pagination={{
        current: pagination.pageNo,
        pageSize: pagination.pageSize,
        total: pagination?.total || 0,
        showTotal: (total) => `${formatMessage('common.total')} ${total} ${formatMessage('common.items')}`,
        onChange: (pageNo, pageSize) => {
          getList(templateId, { pageNo, pageSize });
        },
        showSizeChanger: true,
      }}
      size="middle"
      showHeader={true}
    />
  );
};

export default FileList;
