import { type FC, useState, useEffect } from 'react';
import { Flex, Popover, type PopoverProps, Pagination, Select, Button, Tag, Typography, App, Popconfirm } from 'antd';
import { CaretRight, TrashCan } from '@carbon/icons-react';
import { useTranslate } from '@/hooks';
import { queryNoticeList, noticeToRead, deleteNotice } from '@/apis/inter-api/notify';
import type { TableColumnsType } from 'antd';
import { formatTimestamp } from '@/utils/format';
import ProTable from '../pro-table';

const { Paragraph } = Typography;

interface NoticeItemType {
  id: string;
  msgContent: string;
  createAt: string;
  readStatus: number;
}

export interface NoticePopoverProps extends PopoverProps {
  updateDotStatus: () => void;
}

const NoticePopover: FC<NoticePopoverProps> = ({ children, updateDotStatus, ...restProps }) => {
  const formatMessage = useTranslate();
  const { message } = App.useApp();

  const [open, setOpen] = useState<boolean>(false);
  const [noticeList, setNoticeList] = useState<NoticeItemType[]>([]);
  const [pagination, setPagination] = useState({
    pageNo: 1,
    pageSize: 10,
    total: 0,
  });
  const [readStatus, setReadStatus] = useState<number | undefined>(0);
  const [expandIds, setExpandIds] = useState<string[]>([]);
  const [selectedRowKeys, setSelectedRowKeys] = useState<string[]>([]);

  const getNoticeList = async (params?: any) => {
    const res: any = await queryNoticeList({
      pageNo: pagination.pageNo,
      pageSize: pagination.pageSize,
      // orderField: 'read_status',
      // orderType: 'asc',
      readStatus,
      ...params,
    });
    if (res.code === 0 || res.code === 200) {
      const { data, pageNo, pageSize, total } = res;
      setNoticeList(data || []);
      setPagination({
        pageNo,
        pageSize,
        total,
      });
    }
  };

  const handleDelete = async (data: string[], cb?: () => void) => {
    try {
      await deleteNotice(data);
      getNoticeList({
        pageNo:
          pagination.pageNo > 1 && data.length >= noticeList.length
            ? pagination.pageNo - (Math.ceil((data.length - noticeList.length) / pagination.pageSize) || 1)
            : pagination.pageNo,
      });
      const newExpandIds = expandIds.filter((id) => !data.includes(id));
      const newSelectedRowKeys = selectedRowKeys.filter((id) => !data.includes(id));
      setExpandIds(newExpandIds);
      setSelectedRowKeys(newSelectedRowKeys);
      cb?.();
    } catch (err) {
      console.error(err);
    }
  };

  const handleRead = async (data: string[], cb?: () => void) => {
    try {
      await noticeToRead(data);
      modifyListReadStatus(data);
      updateDotStatus?.();
      cb?.();
    } catch (err) {
      console.error(err);
    }
  };

  const modifyListReadStatus = (data: string[]) => {
    setNoticeList((prev) => {
      return prev.map((item) => {
        if (data.includes(item.id)) {
          item.readStatus = 1;
        }
        return item;
      });
    });
  };

  const resetState = () => {
    setPagination({
      ...pagination,
      pageNo: 1,
    });
    setExpandIds([]);
    setSelectedRowKeys([]);
    setReadStatus(0);
  };

  useEffect(() => {
    // eslint-disable-next-line react-hooks/set-state-in-effect
    if (open) getNoticeList();
  }, [open, readStatus]);

  const columns: TableColumnsType = [
    {
      title: formatMessage('notice.noticeContent'),
      dataIndex: 'msgContent',
      width: '70%',
      render: (text, record) => {
        const { createAt, readStatus, id } = record;
        const isExpand = expandIds.includes(id);
        return (
          <Flex
            style={{ cursor: 'pointer' }}
            onClick={() => {
              const newExpandIds = isExpand ? expandIds.filter((item) => item !== id) : [...expandIds, id];
              setExpandIds(newExpandIds);
              if (readStatus === 0 && !isExpand) handleRead([id]);
            }}
          >
            <CaretRight
              size={20}
              style={{
                flexShrink: 0,
                transform: isExpand ? 'rotate(90deg)' : 'rotate(0deg)',
                transition: '200ms',
                marginTop: '2px',
              }}
            />
            <div>
              {/* <div className="noticeContent">
                {read || text.length <= showLength ? text : `${text.slice(0, showLength)}...`}
              </div> */}
              <Paragraph className="noticeContent" ellipsis={isExpand ? false : { rows: 2 }}>
                {text}
              </Paragraph>
              <div className="noticeSendTime">{formatTimestamp(createAt)}</div>
            </div>
          </Flex>
        );
      },
    },
    {
      title: '',
      dataIndex: 'readStatus',
      width: '20%',
      align: 'center',
      render: (text) => {
        return (
          <Tag color={text === 1 ? 'default' : 'error'} style={{ margin: 0, width: '57px', textAlign: 'center' }}>
            {formatMessage(text === 1 ? 'notice.read' : 'notice.unRead')}
          </Tag>
        );
      },
    },
    {
      title: '',
      dataIndex: 'id',
      width: '10%',
      render: (text) => {
        return (
          <Popconfirm
            title={formatMessage('notice.deleteTip')}
            onConfirm={() => {
              handleDelete([text], () => message.success(formatMessage('common.deleteSuccessfully')));
            }}
          >
            <TrashCan style={{ cursor: 'pointer', marginTop: '2px' }} />
          </Popconfirm>
        );
      },
    },
  ];

  const rowSelection = {
    selectedRowKeys: selectedRowKeys,
    preserveSelectedRowKeys: false,
    columnWidth: '34px',
    onChange: (selectedKeys: any[]) => {
      // setSelectedRowKeys([...new Set([...selectedRowKeys, ...selectedKeys])]);
      setSelectedRowKeys(selectedKeys);
    },
  };

  const noticeContent = (
    <Flex className="noticePopoverContent" vertical>
      <Flex justify="space-between" align="center" style={{ padding: '0 10px', marginBottom: '10px' }}>
        <div className="noticePopoverTitle">{formatMessage('notice.notification')}</div>
        <Flex gap={10} align="center">
          <Popconfirm
            title={formatMessage('notice.deleteMoreTip')}
            onConfirm={() => {
              handleDelete(selectedRowKeys, () => {
                message.success(formatMessage('notice.batchDeleteSuccess'));
                setSelectedRowKeys([]);
              });
            }}
          >
            <Button size="small" disabled={selectedRowKeys.length === 0}>
              {formatMessage('notice.batchDelete')}
            </Button>
          </Popconfirm>
          <Button
            size="small"
            disabled={selectedRowKeys.length === 0}
            type="primary"
            onClick={() =>
              handleRead(selectedRowKeys, () => {
                message.success(formatMessage('notice.batchReadSuccess'));
                setSelectedRowKeys([]);
              })
            }
          >
            {formatMessage('notice.batchRead')}
          </Button>
          <Select
            allowClear
            // prefix={<Filter style={{ verticalAlign: 'text-top' }} />}
            variant="borderless"
            placeholder={formatMessage('common.status')}
            options={[
              { label: formatMessage('notice.unRead'), value: 0 },
              { label: formatMessage('notice.read'), value: 1 },
            ]}
            value={readStatus}
            onChange={(value) => {
              resetState();
              setReadStatus(value);
            }}
          />
        </Flex>
      </Flex>
      <div className="noticeListWrap">
        <ProTable
          className="noticeListTable"
          rowSelection={rowSelection}
          columns={columns}
          dataSource={noticeList}
          rowHoverable={false}
          scroll={{ y: 575 }}
          pagination={false}
        />
      </div>
      {pagination.total > 0 && (
        <div className="noticePagination">
          <Pagination
            size="small"
            align="center"
            current={pagination.pageNo}
            pageSize={pagination.pageSize}
            total={pagination.total}
            showSizeChanger={false}
            onChange={(pageNo) => {
              setExpandIds([]);
              setSelectedRowKeys([]);
              getNoticeList({ pageNo });
            }}
          />
        </div>
      )}
    </Flex>
  );

  return (
    <Popover
      rootClassName="noticePopover"
      placement="bottomRight"
      {...restProps}
      content={noticeContent}
      onOpenChange={(open) => {
        setOpen(open);
        if (open) resetState();
      }}
    >
      {children}
    </Popover>
  );
};

export default NoticePopover;
