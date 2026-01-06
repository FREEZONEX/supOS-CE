import type { FC } from 'react';
import { useTranslate } from '@/hooks';
import { Alert } from 'antd';

import type { FieldItem } from '@/pages/uns/types';
import ComCopyContent from '../../../../components/com-copy/ComCopyContent.tsx';
import ProTable from '@/components/pro-table/index.tsx';
import { formatTimestamp, simpleFormat } from '@/utils/format.ts';

interface PayloadProps {
  websocketData: { [key: string]: any };
  fields: FieldItem[];
}

const Payload: FC<PayloadProps> = ({ websocketData, fields }) => {
  const { data, dt = {}, msg } = websocketData || {};
  const formatMessage = useTranslate();
  if (msg) {
    return <Alert message={<span style={{ color: '#161616' }}>{msg}</span>} type="error" showIcon />;
  }
  const tableData = Object.keys(data || {}).map((key: string) => ({
    key,
    value:
      fields?.find((e) => e.name === key)?.type?.toLowerCase() === 'datetime' ? formatTimestamp(data[key]) : data[key],
    timestamp: typeof dt[key] === 'string' ? (dt[key] as any) - 0 : dt[key],
  }));
  return (
    <ProTable
      bordered={true}
      columns={[
        {
          title: formatMessage('uns.attribute'),
          dataIndex: 'key',
          width: '30%',
          ellipsis: true,
          render: (text) => <span className="payloadFirstTd">{text}</span>,
        },
        {
          title: formatMessage('uns.value'),
          dataIndex: 'value',
          width: '30%',
          ellipsis: true,
          render: (text) => {
            console.log(text);
            const _text = simpleFormat(text);
            return (
              <ComCopyContent
                textToCopy={_text}
                style={{
                  color: 'var(--supos-theme-color)',
                  background: 'transparent',
                  padding: 0,
                }}
              />
            );
          },
        },
        {
          title: formatMessage('common.latestUpdate'),
          dataIndex: 'latestUpdate',
          width: '40%',
          ellipsis: true,
          render: (_, record) => (
            <span style={{ color: 'var(--supos-theme-color)' }}>{formatTimestamp(record.timestamp)}</span>
          ),
        },
      ]}
      dataSource={tableData || []}
      rowKey="key"
      pagination={false}
      size="middle"
      hiddenEmpty
      className={'payload-table'}
      rowHoverable={false}
    />
  );
};
export default Payload;
