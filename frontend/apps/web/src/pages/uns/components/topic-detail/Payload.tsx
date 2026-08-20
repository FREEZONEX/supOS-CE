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

const normalizeDateTimeValue = (value: any) => {
  if (value === null || value === undefined || value === '') return '';
  if (typeof value === 'string') {
    const trimmed = value.trim();
    if (/^\d+$/.test(trimmed)) return Number(trimmed);
  }
  return value;
};

const formatDateTimeValue = (value: any) => {
  const normalized = normalizeDateTimeValue(value);
  const formatted = formatTimestamp(normalized);
  return formatted && formatted !== 'Invalid Date' ? formatted : value;
};

const Payload: FC<PayloadProps> = ({ websocketData, fields }) => {
  const formatMessage = useTranslate();
  const { data, dt = {}, msg } = websocketData || {};
  if (msg) {
    return <Alert message={<span style={{ color: '#161616' }}>{msg}</span>} type="error" showIcon />;
  }
  const tableData = Object.keys(data || {}).map((key: string) => ({
    key,
    value:
      fields?.find((e) => e.name === key)?.type?.toLowerCase() === 'datetime'
        ? formatDateTimeValue(data[key])
        : data[key],
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
            const _text = simpleFormat(text);
            return (
              <ComCopyContent
                textToCopy={_text}
                className="payload-copy-content"
                style={{
                  color: 'var(--ui-text-color)',
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
            <span style={{ color: 'var(--ui-theme-color)' }}>{formatTimestamp(record.timestamp)}</span>
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
