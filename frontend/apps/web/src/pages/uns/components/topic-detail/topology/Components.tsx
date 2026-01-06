import type { FC, ReactNode } from 'react';
import { Button, Flex, Tooltip } from 'antd';
import Binding from '../binding/DashboardBinding.tsx';
import { flowPage } from '@/apis/inter-api/flow.ts';
import type { NodeDataType } from './types.ts';
import { useTranslate } from '@/hooks';
import { AddLarge, ApplicationWeb, InformationFilled, Launch } from '@carbon/icons-react';
import classNames from 'classnames';
import styles from './TopologyChart.module.scss';
import { useBaseStore } from '@/stores/base';
import postgresql from '@/assets/home-icons/postgresql.svg';
import tdengine from '@/assets/home-icons/tdengine.png';
import timescaleDB from '@/assets/home-icons/timescaleDB.svg';
import { getDashboardList } from '@/apis/inter-api';
import nodeRed from '@/assets/home-icons/node-red.svg';
import ProTable from '@/components/pro-table';
import ComCodeSnippet from '@/components/com-code-snippet';
import ComFormula from '@/components/com-formula';
import c2 from '@/assets/uns/cw.svg';

export const CommonNode = ({
  Icon,
  active,
  title,
  subtitle,
  NavigateBtn,
  bindConfig,
  indicatorConfig,
}: {
  Icon?: ReactNode;
  active?: boolean;
  title: string;
  subtitle?: string;
  NavigateBtn?: ReactNode;
  bindConfig?: {
    api?: any;
    selectValue?: string;
    onBinding?: (item: any) => any;
  };
  indicatorConfig?: {
    statusColor: string;
    title: string;
  };
}) => {
  return (
    <Flex
      className={classNames(styles['common-node'], styles['common-node-hover'], {
        [styles['activeBg']]: active,
      })}
      align="center"
      gap={12}
    >
      {Icon}
      <Flex vertical gap={4} className={styles['common-node-content']}>
        {subtitle && (
          <span className={styles['common-node-subtitle']} title={subtitle}>
            {subtitle}
          </span>
        )}
        <span className={styles['common-node-title']} title={title}>
          {title}
        </span>
      </Flex>
      <Flex>
        {NavigateBtn && (
          <div className={styles['common-node-btn']} data-action="navigate">
            <Button
              size="small"
              color="default"
              variant="text"
              style={{ color: 'var(--supos-text-color)', padding: '0 4px' }}
            >
              {NavigateBtn}
            </Button>
          </div>
        )}
        {bindConfig && (
          <div className={styles['common-node-btn']} data-action="noNavigate">
            <Binding selectValue={bindConfig?.selectValue} api={bindConfig?.api} onBinding={bindConfig?.onBinding} />
          </div>
        )}
      </Flex>
      {indicatorConfig && (
        <div className={styles['status-indicator']}>
          <span className={styles['status-dot']} style={{ background: indicatorConfig?.statusColor }} />
          <span className={styles['status-content']} title={indicatorConfig?.title}>
            {indicatorConfig?.title}
          </span>
        </div>
      )}
    </Flex>
  );
};

const TooltipContent = () => {
  const formatMessage = useTranslate();
  return (
    <Flex align="center" gap={4}>
      <Flex align="center" gap={4}>
        <InformationFilled style={{ color: 'var(--supos-theme-color)', marginRight: 4, width: 30, height: 30 }} />
        <div>
          <span style={{ color: 'var(--supos-text-color)', fontWeight: 600, fontSize: 12 }}>
            {formatMessage('common.nextStep')}:
          </span>
          <span style={{ color: 'var(--supos-text-color)', fontSize: 12 }}>
            {formatMessage('common.clickSourceFlow')}
          </span>
        </div>
      </Flex>
    </Flex>
  );
};

export const NodeRed: FC<NodeDataType> = (data) => {
  const configured = data.node.data.id || data.node.data.flowId || data.node.data.flowName;
  const loading = data.node.data.loading;
  const statusColor = configured ? '#4CAF50' : '#B1973B';
  const formatMessage = useTranslate();
  return (
    <div id={`node-red-container-${data.node.id}`} style={{ position: 'relative', width: '100%', height: '100%' }}>
      <Tooltip
        title={TooltipContent}
        open={data.node && !configured}
        placement="topRight"
        color="var(--supos-bg-color)"
        align={{
          offset: [-33, -10],
        }}
        styles={{
          body: {
            position: 'relative',
            right: -100,
            backgroundColor: '#fff',
            borderRadius: 2,
            padding: 8,
          },
        }}
        getPopupContainer={() => document.getElementById(`node-red-container-${data.node.id}`) as HTMLElement}
      >
        <CommonNode
          Icon={<img src={nodeRed} alt="" width="28px" />}
          title={formatMessage('home.sourceFlow')}
          subtitle={formatMessage('common.nodeRed', 'Node-Red')}
          active={data.node.data.active}
          NavigateBtn={
            loading ? (
              <div className={styles['loading-spinner']} />
            ) : configured ? (
              <Launch size={20} />
            ) : (
              <AddLarge size={20} />
            )
          }
          bindConfig={{
            selectValue: data.node.data.bindId,
            api: flowPage,
            onBinding: (item: any) => {
              return data.node.data.onBindingChange?.('nodeRed1', item);
            },
          }}
          indicatorConfig={{
            statusColor,
            title: formatMessage(configured ? 'common.configured' : 'common.unconfigured'),
          }}
        />
      </Tooltip>
    </div>
  );
};

export const Mqtt: FC<NodeDataType> = (data) => {
  const mqttBrokeType = useBaseStore((state) => state.mqttBrokeType);
  return (
    <CommonNode
      title={'MQTT Broker'}
      subtitle={mqttBrokeType?.toUpperCase() || 'EMQX'}
      active={data.node.data.active}
    />
  );
};

export const DataBase: FC<NodeDataType> = (data) => {
  const { dataBaseType, systemInfo } = useBaseStore((state) => ({
    dataBaseType: state.dataBaseType,
    systemInfo: state.systemInfo,
  }));
  const props = [2, 8].includes(data.node.data.dataType)
    ? {
        Icon: <img src={postgresql} alt="" width="28px" />,
        subtitle: 'PostgreSQL',
        title: 'Relational DB',
        NavigateBtn: systemInfo?.containerMap?.chat2db ? <Launch size={20} /> : undefined,
      }
    : {
        Icon: <img src={dataBaseType.includes('tdengine') ? tdengine : timescaleDB} width="28px" />,
        title: dataBaseType.includes('tdengine') ? 'tdengine' : 'Database',
        subtitle: 'TimescaleDB',
      };
  return <CommonNode active={data.node.data.active} {...props} />;
};

export const Apps: FC<NodeDataType> = (data) => {
  return (
    <CommonNode
      Icon={<ApplicationWeb size={28} />}
      title="Dashboard"
      subtitle={data.node.data.subtitle}
      active={data.node.data.active}
      NavigateBtn={<Launch size={20} />}
      bindConfig={{
        selectValue: data.node.data.bindId,
        api: getDashboardList,
        onBinding: (item: any) => {
          return data.node.data.onBindingChange?.('apps1', item);
        },
      }}
    />
  );
};

export const NodeRedDetail: FC<any> = ({ flowList }) => {
  const formatMessage = useTranslate();
  return (
    <div style={{ width: '100%', display: 'contents' }}>
      {flowList && (
        <div style={{ width: '100%' }}>
          <div style={{ width: '100%', marginBottom: 12 }} className={styles['name']}>
            {/*<CautionInverted style={{ marginRight: 8, width: 10, height: 10 }} />*/}
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
                value: flowList?.description, // 注意大小写一致性
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
  const fieldList = instanceInfo?.fields?.map((field: any) => `"${field.name}"`).join(', ') || '"*"';
  const sql = !instanceInfo?.tbFieldName
    ? `SELECT ${fieldList} FROM "public"."${instanceInfo?.table}" LIMIT 10`
    : `SELECT ${fieldList} FROM "public"."${instanceInfo?.table}" WHERE tag=${instanceInfo?.id} LIMIT 10`;
  return (
    <>
      <div style={{ width: '100%', marginBottom: 12 }} className={styles['name']}>
        SQL
      </div>
      <ComCodeSnippet
        style={{ border: '1px solid var(--supos-table-tr-color)' }}
        minCollapsedNumberOfRows={4}
        maxCollapsedNumberOfRows={4}
        copyPosition={true}
        copyText={sql}
      >
        {sql}
      </ComCodeSnippet>
    </>
  );
};

export const MqttDetail: FC = () => {
  const formatMessage = useTranslate();
  const systemInfo = useBaseStore((state) => state.systemInfo);

  const dataSource = [
    {
      key: 'front',
      detail: formatMessage('uns.front'),
      content: `mqtt://${window.location.hostname}:${systemInfo?.mqttTcpPort}/mqtt`,
    },
    {
      key: 'backend',
      detail: formatMessage('uns.backend'),
      content: `tcp://${window.location.hostname}:${systemInfo?.mqttTcpPort}/mqtt`,
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
  // 用 Set 来存储已经处理过的组合（topic + field）
  const seen = new Set();
  // 去重
  const uniqueArr = instanceInfo?.refers?.filter((item: any) => {
    const key = `${item.topic}-${item.field}`; // 创建唯一的 key
    if (seen.has(key)) {
      return false; // 如果 key 已经存在，跳过该项
    } else {
      seen.add(key); // 否则加入 seen 集合
      return true; // 保留该项
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

  // 用来存储最终的替换结果
  let resultStr = instanceInfo.expression;

  // 遍历 newd 数组，并替换对应的值
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

export const ButtonError: FC<any> = () => {
  return (
    <div className={styles['buttonError']}>
      <img src={c2} alt="cw" />
    </div>
  );
};
