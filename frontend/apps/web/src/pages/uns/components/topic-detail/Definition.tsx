import type { FC } from 'react';
import { useTranslate } from '@/hooks';
import Icon from '@ant-design/icons';
import ProTable from '@/components/pro-table';
import MainKey from '@/components/svg-components/MainKey';

interface DefinitionProps {
  instanceInfo: { [key: string]: any };
}

const Definition: FC<DefinitionProps> = ({ instanceInfo }) => {
  const formatMessage = useTranslate();
  const { id, extendFieldUsed = [], fields } = instanceInfo || {};

  return (
    <ProTable
      key={id}
      bordered
      showExpand={extendFieldUsed.length > 0}
      columns={[
        {
          title: formatMessage('uns.attribute'),
          dataIndex: 'name',
          width: '20%',
          render: (text: any, record: any) => (
            <div>
              {record.unique && (
                <Icon
                  style={{
                    color: 'var(--supos-theme-color)',
                    marginRight: '5px',
                    verticalAlign: 'middle',
                  }}
                  title={formatMessage('uns.mainKey')}
                  component={MainKey}
                />
              )}
              {text}
            </div>
          ),
        },
        {
          title: formatMessage('uns.type'),
          dataIndex: 'type',
          width: '20%',
          render: (text: any) => <span style={{ color: 'var(--supos-theme-color)' }}>{text}</span>,
        },
        {
          title: formatMessage('common.length'),
          dataIndex: 'maxLen',
          width: '20%',
          render: (text: any) => <span style={{ color: 'var(--supos-theme-color)' }}>{text}</span>,
        },
        {
          title: formatMessage('uns.displayName'),
          dataIndex: 'displayName',
          width: '20%',
          render: (text: any) => <span style={{ color: 'var(--supos-theme-color)' }}>{text}</span>,
        },
        {
          title: formatMessage('uns.remark'),
          dataIndex: 'remark',
          width: '20%',
          render: (text: any) => <span style={{ color: 'var(--supos-theme-color)' }}>{text}</span>,
        },
        {
          title: formatMessage('uns.unit'),
          dataIndex: extendFieldUsed?.includes('unit') ? 'unit' : undefined,
          width: '20%',
          hidden: true,
          render: (text: any) => <span style={{ color: 'var(--supos-theme-color)' }}>{text}</span>,
        },
        {
          title: formatMessage('uns.upperLimit'),
          dataIndex: extendFieldUsed?.includes('upperLimit') ? 'upperLimit' : undefined,
          width: '20%',
          hidden: true,
          render: (text: any) => <span style={{ color: 'var(--supos-theme-color)' }}>{text}</span>,
        },
        {
          title: formatMessage('uns.lowerLimit'),
          dataIndex: extendFieldUsed?.includes('lowerLimit') ? 'lowerLimit' : undefined,
          width: '20%',
          hidden: true,
          render: (text: any) => <span style={{ color: 'var(--supos-theme-color)' }}>{text}</span>,
        },
        {
          title: formatMessage('uns.decimal'),
          dataIndex: extendFieldUsed?.includes('decimal') ? 'decimal' : undefined,
          width: '20%',
          hidden: true,
          render: (text: any) => <span style={{ color: 'var(--supos-theme-color)' }}>{text}</span>,
        },
      ].filter((item) => item.dataIndex)}
      dataSource={fields || []}
      rowKey="name"
      pagination={false}
      size="middle"
      hiddenEmpty
      rowHoverable={false}
    />
  );
};
export default Definition;
