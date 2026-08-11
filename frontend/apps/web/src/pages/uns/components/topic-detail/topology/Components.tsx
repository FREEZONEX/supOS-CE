import type { FC } from 'react';
import { useTranslate } from '@/hooks';
import styles from './TopologyChart.module.scss';
import { useBaseStore } from '@/stores/base';
import ProTable from '@/components/pro-table';
import ComCodeSnippet from '@/components/com-code-snippet';
import ComFormula from '@/components/com-formula';

const DEFAULT_MQTT_WEBSOCKET_PORT = '8083';
const DEFAULT_MQTT_TCP_PORT = '1883';

const quoteSqlIdent = (value: string | number) => `"${String(value).replace(/"/g, '""')}"`;

const normalizeSqlIdentValue = (value: unknown) => {
  if (typeof value !== 'string' && typeof value !== 'number') {
    return '';
  }
  const normalized = String(value)
    .trim()
    .replace(/^\/+|\/+$/g, '');
  if (
    !normalized ||
    normalized.includes('[object Object]') ||
    normalized.includes('[object Array]') ||
    normalized === 'undefined' ||
    normalized === 'null'
  ) {
    return '';
  }
  return normalized;
};

const firstSqlIdentValue = (...values: unknown[]) => {
  for (const value of values) {
    const normalized = normalizeSqlIdentValue(value);
    if (normalized) {
      return normalized;
    }
  }
  return '';
};

export const resolveDatabaseSchemaName = (instanceInfo: any) =>
  firstSqlIdentValue(instanceInfo?.dbSchemaName, instanceInfo?.schemaName, instanceInfo?.schema) || 'uns';

export const resolveDatabaseTableName = (instanceInfo: any) =>
  firstSqlIdentValue(instanceInfo?.table, instanceInfo?.tableName, instanceInfo?.alias, instanceInfo?.name);

export const NodeRedDetail: FC<any> = ({ flowList }) => {
  const formatMessage = useTranslate();
  return (
    <div style={{ width: '100%', display: 'contents' }}>
      {flowList && (
        <div style={{ width: '100%' }}>
          <div style={{ width: '100%', marginBottom: 12 }} className={styles['name']}>
            {formatMessage('home.sourceFlow')}
          </div>
          <ProTable
            bordered
            rowHoverable={false}
            className={styles.customTable}
            columns={[
              {
                title: formatMessage('common.detail'),
                dataIndex: 'label',
                key: 'label',
                width: '30%',
                render: (text: any) => <span className={styles.detailLabel}>{text}</span>,
              },
              {
                title: formatMessage('uns.content'),
                dataIndex: 'value',
                key: 'value',
                width: '70%',
                render: (value: any) => value || <span className={styles.empty}>-</span>,
              },
            ]}
            dataSource={[
              {
                key: 'flowName',
                label: formatMessage('uns.CollectionFlowName'),
                value: flowList?.flowName,
              },
              {
                key: 'template',
                label: formatMessage('uns.flowTemplate'),
                value: flowList?.template,
              },
              {
                key: 'description',
                label: formatMessage('uns.description'),
                value: flowList?.description,
              },
            ]}
            pagination={false}
            showHeader={true}
            rowKey="key"
          />
        </div>
      )}
    </div>
  );
};

export const DataBaseDetail: FC<any> = ({ instanceInfo }) => {
  const fields = Array.isArray(instanceInfo?.fields)
    ? instanceInfo.fields.map((field: any) => field?.name || field?.fieldName).filter(Boolean)
    : [];
  const fieldList = fields.length ? fields.map((field: string) => quoteSqlIdent(field)).join(', ') : '*';
  const schemaName = resolveDatabaseSchemaName(instanceInfo);
  const tableName = resolveDatabaseTableName(instanceInfo);
  const whereClause = instanceInfo?.tbFieldName
    ? ` WHERE ${quoteSqlIdent(instanceInfo.tbFieldName)} = ${Number(instanceInfo?.id || 0)}`
    : '';
  const sql = tableName
    ? `SELECT ${fieldList} FROM ${quoteSqlIdent(schemaName)}.${quoteSqlIdent(tableName)}${whereClause} LIMIT 10`
    : '';
  return (
    <>
      <div style={{ width: '100%', marginBottom: 12 }} className={styles['name']}>
        SQL
      </div>
      <ComCodeSnippet
        style={{ border: '1px solid var(--ui-table-tr-color)' }}
        minCollapsedNumberOfRows={4}
        maxCollapsedNumberOfRows={4}
        copyPosition={true}
        copyText={sql}
      >
        {sql || '-'}
      </ComCodeSnippet>
    </>
  );
};

export const MqttDetail: FC = () => {
  const formatMessage = useTranslate();
  const systemInfo = useBaseStore((state) => state.systemInfo);
  const mqttWebsocketPort = systemInfo?.mqttWebsocketPort || DEFAULT_MQTT_WEBSOCKET_PORT;
  const mqttTcpPort = systemInfo?.mqttTcpPort || DEFAULT_MQTT_TCP_PORT;

  const dataSource = [
    {
      key: 'front',
      detail: formatMessage('uns.front'),
      content: `mqtt://${window.location.hostname}:${mqttWebsocketPort}/mqtt`,
    },
    {
      key: 'backend',
      detail: formatMessage('uns.backend'),
      content: `tcp://${window.location.hostname}:${mqttTcpPort}/mqtt`,
    },
  ];

  const columns = [
    {
      title: formatMessage('common.detail'),
      dataIndex: 'detail',
      key: 'detail',
      width: '30%',
      render: (text: string) => <td className="payloadFirstTd">{text}</td>,
    },
    {
      title: formatMessage('uns.content'),
      dataIndex: 'content',
      width: '70%',
      key: 'content',
    },
  ];

  return (
    <>
      <div style={{ width: '100%', marginBottom: 12 }} className={styles['name']}>
        MQTT Broker
      </div>
      <ProTable
        className={styles.customTable}
        columns={columns}
        dataSource={dataSource}
        pagination={false}
        showHeader={true}
        rowKey="key"
        rowHoverable={false}
        bordered
        hiddenEmpty
      />
    </>
  );
};

export const MqttDetail2: FC<any> = ({ instanceInfo }) => {
  const formatMessage = useTranslate();
  const newd: any = [];
  const newd2: any = [];
  const seen = new Set();
  const uniqueArr = instanceInfo?.refers?.filter((item: any) => {
    const key = `${item.topic}-${item.field}`;
    if (seen.has(key)) {
      return false;
    } else {
      seen.add(key);
      return true;
    }
  });
  uniqueArr?.forEach((item: any, index: number) => {
    newd.push({
      label: 'Variable' + (index + 1),
      value: `"${item.topic}".${item.field}`,
    });
    newd2.push({
      label: 'Variable' + (index + 1),
      value: `${'Variable' + (index + 1)}`,
    });
  });

  let resultStr = instanceInfo.expression;
  newd.forEach((item: any) => {
    const valueRegex = new RegExp(item.value, 'g');
    resultStr = resultStr.replace(valueRegex, `${item.label}`);
  });
  const columns = [
    {
      title: formatMessage('uns.variable'),
      dataIndex: 'variable',
      width: '30%',
      render: (_: any, __: any, index: number) => formatMessage('uns.variable') + (index + 1),
    },
    {
      title: formatMessage('uns.topic'),
      dataIndex: 'topic',
      width: '40%',
      key: 'topic',
    },
    {
      title: formatMessage('uns.attribute'),
      dataIndex: 'field',
      width: '30%',
      key: 'field',
    },
  ];

  return (
    <div className={styles['Tables']}>
      <ComFormula fieldList={newd2} defaultOpenCalculator={false} value={resultStr} readonly={true} />
      <ProTable
        className={styles.customTable}
        columns={columns}
        dataSource={uniqueArr}
        bordered
        pagination={false}
        rowHoverable={false}
        hiddenEmpty
        rowKey={(_: any, index: any) => `row-${index}`}
      />
    </div>
  );
};
