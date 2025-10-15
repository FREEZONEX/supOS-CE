import { useState, useEffect, forwardRef, useImperativeHandle } from 'react';
import { pageListUnsByLabel, cancelLabel } from '@/apis/inter-api/uns';
import { useTranslate } from '@/hooks';
import { Flex, App, Tooltip } from 'antd';
import { Document, TrashCan } from '@carbon/icons-react';
import { ButtonPermission } from '@/common-types/button-permission.ts';
import ProTable from '@/components/pro-table';
import { hasPermission } from '@/utils/auth';

interface FileListRefProps {
  getList: (id: string, params?: any) => void;
}
interface FileListProps {
  labelId: string;
}

const initPagination = {
  pageNo: 1,
  pageSize: 10,
};

const FileList = forwardRef<FileListRefProps, FileListProps>(({ labelId }, ref) => {
  const formatMessage = useTranslate();

  const [list, setList] = useState([]);
  const [pagination, setPagination] = useState({
    ...initPagination,
    total: 0,
  });

  const { message } = App.useApp();

  useImperativeHandle(ref, () => ({
    getList,
  }));

  const getList = (labelId: string, params?: any) => {
    pageListUnsByLabel({ labelId, pageNo: pagination.pageNo, pageSize: pagination.pageSize, ...params })
      .then((res) => {
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
      })
      .catch((error) => {
        console.log(error);
      });
  };
  const deleteRow = (unsId: string) => {
    const { pageNo } = pagination;
    cancelLabel(unsId, [labelId]).then(() => {
      message.success(formatMessage('common.deleteSuccessfully'));
      getList(labelId, { pageNo: list.length === 1 && pageNo > 1 ? pageNo - 1 : pageNo });
    });
  };

  useEffect(() => {
    if (labelId) {
      getList(labelId, initPagination);
    }
  }, [labelId]);

  return (
    <>
      <ProTable
        rowKey={'unsId'}
        size="middle"
        bordered
        rowHoverable={false}
        hiddenEmpty
        dataSource={list || []}
        columns={[
          {
            title: formatMessage('common.name'),
            dataIndex: 'name',
            width: '25%',
            render: (text) => (
              <Flex>
                <Flex style={{ height: '24px' }} align="center">
                  <Document style={{ marginRight: '5px' }} />
                </Flex>
                <span style={{ wordBreak: 'break-word' }}>{text}</span>
              </Flex>
            ),
          },
          {
            title: formatMessage('uns.position'),
            dataIndex: 'path',
            width: '70%',
            render: (text: any) => (
              <span style={{ color: 'var(--supos-theme-color)', wordBreak: 'break-word' }}>{text}</span>
            ),
          },
          {
            title: formatMessage('common.operation'),
            dataIndex: 'operation',
            width: '5%',
            align: 'center',
            hidden: !hasPermission(ButtonPermission['uns.labelDetail']),
            render: (_, record: any) => (
              <Tooltip title={formatMessage('common.delete')}>
                <TrashCan
                  onClick={() => {
                    deleteRow(record?.unsId);
                  }}
                  style={{ cursor: 'pointer' }}
                />
              </Tooltip>
            ),
          },
        ]}
        pagination={{
          current: pagination.pageNo,
          pageSize: pagination.pageSize,
          total: pagination?.total || 0,
          showTotal: (total) => `${formatMessage('common.total')} ${total} ${formatMessage('common.items')}`,
          onChange: (pageNo, pageSize) => {
            getList(labelId, { pageNo, pageSize });
          },
          showSizeChanger: true,
        }}
      />
    </>
  );
});
export default FileList;
