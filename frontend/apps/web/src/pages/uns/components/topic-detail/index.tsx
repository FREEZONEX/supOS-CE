import { useState, useEffect, useRef, useMemo, useCallback, type FC } from 'react';
import { getInstanceInfo } from '@/apis/core-api/uns';
import { Collapse, Flex, theme, Typography } from 'antd';
import { CaretRight, ClipboardList, SendAlt, ChartLine } from '@/components/lucide-icon/carbon';
import { UNS_TOPIC_ICON_COLORS, UnsTopicTypeIconWrap } from '@/pages/uns/components/uns-tree/tree-icons';
import { useTranslate } from '@/hooks';
import type { CSSProperties } from 'react';
import type { CollapseProps } from 'antd';
import Details from './Details';
import TopologyChart from './topology/TopologyChart';
import Definition from './Definition';
import Payload from './Payload';
import Dashboard from './Dashboard';
import RawData from './RawData';
import PayloadViewSegmented, { type PayloadViewMode } from './PayloadViewSegmented';
// import SqlQuery from './SqlQuery';
import DocumentList from '@/pages/uns/components/DocumentList.tsx';
import UploadButton from '@/pages/uns/components/UploadButton.tsx';
import { ButtonPermission } from '@/common-types/button-permission.ts';
import { useMediaSize } from '@/hooks';
import EditDetailButton from '@/pages/uns/components/EditDetailButton';
import type { InitTreeDataFnType, UnsTreeNode, FieldItem } from '@/pages/uns/types';
import { isJsonString } from '@/utils/common';
import { useBaseStore } from '@/stores/base';
// import Subscribe from '@/pages/uns/components/subscribe';
import EditButton from '@/pages/uns/components/EditButton.tsx';
import useSSE from '@/hooks/useSSE.ts';

export interface FileDetailProps {
  currentNode: UnsTreeNode;
  initTreeData: InitTreeDataFnType;
  handleDelete: (node: UnsTreeNode) => void;
  readOnly?: boolean;
}

interface InstanceInfoType {
  [key: string]: any;
}

const { Title } = Typography;

const Module: FC<FileDetailProps> = (props) => {
  const {
    currentNode: { id, mount },
    readOnly = false,
  } = props;
  const {
    systemInfo: { qualityName = '_quality', timestampName = '_timestamp', enableAutoCategorization },
  } = useBaseStore((state) => ({
    systemInfo: state.systemInfo,
  }));
  const formatMessage = useTranslate();
  const documentListRef = useRef(null);
  const [instanceInfo, setInstanceInfo] = useState<InstanceInfoType>({});
  const [activeList, setActiveList] = useState<string[]>([
    'topologyChart',
    'definition',
    'payload',
    'dashboard',
    // 'sqlQuery',
  ]);
  const [payloadView, setPayloadView] = useState<PayloadViewMode>('table');
  const [websocketData, setWebsocketData] = useState<any>({});
  const [localSchemaFields, setLocalSchemaFields] = useState<FieldItem[] | null>(null);
  const { token } = theme.useToken();

  const panelStyle: CSSProperties = {
    background: 'val(--ui-bg-color)',
    border: 'none',
  };

  const { isH5 } = useMediaSize();

  const longToJavaHex = (value: string, fullLength = false) => {
    let bigIntValue;
    try {
      bigIntValue = BigInt(value);
    } catch (e) {
      console.log(e);
      bigIntValue = BigInt(JSON.parse(value));
    }

    // 获取对应的无符号 64 位表示（补码兼容）
    const mask64 = 0xffffffffffffffffn;
    const unsigned = bigIntValue < 0n ? ((bigIntValue & mask64) + (1n << 64n)) & mask64 : bigIntValue & mask64;

    let hex = unsigned.toString(16);

    if (fullLength) {
      hex = hex.padStart(16, '0');
    }

    return hex;
  };

  const isZeroOrPositiveInteger = (value: number | string | undefined) => {
    // 将输入转换为数字
    const num = Number(value);

    // 检查：
    // 1. 转换后是有限的数字 (排除 NaN, Infinity)
    // 2. 是整数 (Math.floor(num) === num)
    // 3. 大于等于 0
    return Number.isFinite(num) && Math.floor(num) === num && num >= 0;
  };

  const formatDecimal = (str: string, digits: number) => {
    // 验证输入

    const _str = String(str);
    if (!/^-?\d*\.?\d+$/.test(_str.trim()) || digits < 0 || !str) return str;

    let intPart = _str.trim().split('.')[0];
    const decPart = _str.trim().split('.')[1];
    let dec = decPart.slice(0, digits); // 截取到所需位数

    // 判断是否需要进位 (检查第 digits+1 位)
    if (decPart.length > digits && parseInt(decPart[digits]) >= 5) {
      // 小数部分加1 (处理进位)
      const num = (parseInt(dec, 10) + 1).toString().padStart(digits, '0');
      if (num.length > digits) {
        // 进位到整数部分
        intPart = (BigInt(intPart) + (intPart[0] === '-' ? -1n : 1n)).toString();
        dec = '0'.repeat(digits);
      } else {
        dec = num;
      }
    } else {
      dec = dec.padEnd(digits, '0'); // 不足补零
    }

    return digits === 0 ? intPart : `${intPart}.${dec}`;
  };

  useSSE(instanceInfo.id ? `/api/core/uns/newMsg?id=${instanceInfo.id}` : '', {
    onMessage: (event) => {
      const dataJson = event.data;
      if (isJsonString(dataJson)) {
        const data = JSON.parse(dataJson);
        if (qualityName && data?.data?.[qualityName]) {
          //质量码做特殊处理
          data.data[qualityName] = longToJavaHex(data.data[qualityName]);
        }
        if (typeof data.payload === 'string' && !isJsonString(data.payload)) {
          data.payload = null;
        } else if (data.payload != null && typeof data.payload !== 'string' && typeof data.payload !== 'object') {
          data.payload = null;
        }
        if (instanceInfo?.dataType === 2 && timestampName && data?.data?.[timestampName]) {
          //关系型文件手动隐藏消息体里的时间戳
          delete data.data[timestampName];
        }
        instanceInfo?.fields?.forEach((field: FieldItem) => {
          if (
            ['FLOAT', 'DOUBLE'].includes(field.type) &&
            data?.data?.[field.name] &&
            isZeroOrPositiveInteger(field.decimal)
          ) {
            data.data[field.name] = formatDecimal(data.data[field.name], Number(field.decimal));
          }
        });
        setWebsocketData(data);
      }
    },
    onError: (error) => console.error('WebSocket error:', error),
  });

  useEffect(() => {
    setWebsocketData({});
    setLocalSchemaFields(null);
    if (id) {
      getFileDetail(id as string);
    } else {
      setInstanceInfo({});
    }
  }, [id]);

  const schemaFields = useMemo(() => {
    if (localSchemaFields !== null) {
      return localSchemaFields;
    }
    return Array.isArray(instanceInfo.fields) ? instanceInfo.fields : [];
  }, [instanceInfo, localSchemaFields]);

  const getFileDetail = (id: string) => {
    return getInstanceInfo({ id })
      .then(async (data: any) => {
        if (data?.id) {
          if (data?.dataType === 8) {
            data.fields = data?.jsonFields || [];
          }
          data?.fields?.forEach((field: FieldItem) => {
            //特殊处理挂载文件的异常数据
            if (['STRING', 'BOOLEAN'].includes(field.type) && Number(field.decimal) < 0) {
              field.decimal = undefined;
            }
          });
          // data.extendFieldUsed = data.mount
          //   ? ['unit', 'upperLimit', 'lowerLimit', 'decimal']
          //   : data.extendFieldUsed || [];
          data.extendFieldUsed = [];
          if (data?.lastPayload) {
            setWebsocketData(data.lastPayload);
          } else {
            setWebsocketData({});
          }
          setInstanceInfo(data);
          return data;
        }
      })
      .catch(() => {});
  };

  const refreshSchemaModel = useCallback(
    (savedFields?: FieldItem[]) => {
      if (savedFields !== undefined) {
        setLocalSchemaFields(savedFields);
        void getFileDetail(id as string);
        return Promise.resolve();
      }
      setLocalSchemaFields(null);
      return getFileDetail(id as string);
    },
    [id]
  );

  const getItems: (panelStyle: CSSProperties, instanceInfo: InstanceInfoType) => CollapseProps['items'] = (
    panelStyle,
    instanceInfo
  ) => {
    const isMountedTopic = Boolean(mount || instanceInfo.mount);
    const schemaModelInfo = { ...instanceInfo, fields: schemaFields, mount: isMountedTopic };
    const topicReadOnly = readOnly || isMountedTopic;
    const schemaEditable = !topicReadOnly && [1, 2, 8].includes(instanceInfo.dataType);

    const items = [
      {
        key: 'detail',
        label: formatMessage('common.detail'),
        children: (
          <Details instanceInfo={instanceInfo} updateTime={websocketData?.updateTime} websocketData={websocketData} />
        ),
        style: panelStyle,
        extra: readOnly ? null : (
          <EditDetailButton
            auth={ButtonPermission['uns.fileDetail']}
            modelInfo={{ ...instanceInfo, mount: isMountedTopic }}
            getModel={() => getFileDetail(id as string)}
          />
        ),
      },
      {
        key: [1, 2, 3, 6, 7, 8].includes(instanceInfo.dataType) ? 'definition' : '',
        label: formatMessage('uns.definition'),
        children: (
          <Definition
            instanceInfo={{ ...instanceInfo, fields: schemaFields }}
            modelInfo={schemaModelInfo}
            getModel={refreshSchemaModel}
            auth={ButtonPermission['uns.fileDetail']}
            editable={schemaEditable}
          />
        ),
        style: panelStyle,
        extra: schemaEditable ? (
          <EditButton
            auth={ButtonPermission['uns.fileDetail']}
            modelInfo={schemaModelInfo}
            getModel={refreshSchemaModel}
            editType="file"
            triggerIcon="add"
          />
        ) : null,
      },
      {
        key: [1, 2, 3, 6, 7, 8].includes(instanceInfo.dataType) ? 'payload' : '',
        label: formatMessage('uns.payload'),
        children:
          instanceInfo.dataType === 8 || payloadView === 'code' ? (
            <RawData payload={websocketData?.payload ?? websocketData?.data} className="payload-code-view" />
          ) : (
            <Payload websocketData={websocketData} fields={schemaFields} />
          ),
        style: panelStyle,
        extra:
          instanceInfo.dataType === 8 ? null : (
            <PayloadViewSegmented
              value={payloadView}
              onChange={setPayloadView}
              tableTitle={formatMessage('common.table')}
              codeTitle="JSON"
            />
          ),
      },
      ...(!isH5
        ? [
            {
              key: instanceInfo.dataType !== 7 ? 'dashboard' : '',
              label: formatMessage('uns.dashboard'),
              children: <Dashboard instanceInfo={instanceInfo} />,
              style: panelStyle,
            },
            {
              key: [1, 2, 8].includes(instanceInfo.dataType) ? 'topologyChart' : '',
              label: formatMessage('uns.topology'),
              children: (
                <TopologyChart
                  getFileDetail={getFileDetail}
                  instanceInfo={instanceInfo}
                  readOnly={topicReadOnly}
                  // payload={websocketData?.data}
                  // dt={websocketData?.dt || {}}
                />
              ),
              style: panelStyle,
            },
          ]
        : []),
      ...(!isH5
        ? [
            // {
            //   id: 'sqlQuery',
            //   key: 'sqlQuery',
            //   label: formatMessage('uns.dataOperation'),
            //   children: <SqlQuery instanceInfo={instanceInfo} id={id as string} />,
            //   style: panelStyle,
            // },
          ]
        : []),
      {
        key: 'document',
        label: formatMessage('common.document'),
        children: (
          <DocumentList
            alias={instanceInfo.alias}
            ownerId={instanceInfo.id}
            readOnly={topicReadOnly}
            ref={documentListRef}
          />
        ),
        style: panelStyle,
        extra: topicReadOnly ? null : (
          <UploadButton
            setActiveList={setActiveList}
            auth={ButtonPermission['uns.fileDetail']}
            alias={instanceInfo.alias}
            ownerId={instanceInfo.id}
            documentListRef={documentListRef}
          />
        ),
      },
    ];
    return items.filter((item: any) => item.key);
  };

  // const handleChangeSubscribe = async (enable: boolean, frequency?: string) => {
  //   await updateModelSubscribe({ id, enable, frequency });
  //   getFileDetail(id as string);
  //   message.success(enable ? formatMessage('uns.subscribeSuccessful') : formatMessage('uns.unsubscribeSuccessful'));
  // };

  const getFileIcon = () => {
    switch (instanceInfo.parentDataType) {
      case 1:
        return (
          <ClipboardList size={20} strokeWidth={1.75} aria-hidden style={{ color: UNS_TOPIC_ICON_COLORS.state }} />
        );
      case 2:
        return <SendAlt size={20} strokeWidth={1.75} aria-hidden style={{ color: UNS_TOPIC_ICON_COLORS.action }} />;
      case 3:
        return <ChartLine size={20} strokeWidth={1.75} aria-hidden style={{ color: UNS_TOPIC_ICON_COLORS.metric }} />;
      default:
        return null;
    }
  };

  return (
    <div className="topicDetailWrap">
      <div
        className="topicDetailContent"
        style={{
          paddingLeft: 5,
          paddingRight: 5,
          paddingBottom: '20px',
        }}
      >
        <Flex className="detailTitle" gap={8} align="center">
          {enableAutoCategorization ? (
            <UnsTopicTypeIconWrap topicType={instanceInfo.parentDataType ?? 0}>{getFileIcon()}</UnsTopicTypeIconWrap>
          ) : (
            <ClipboardList size={20} strokeWidth={1.75} aria-hidden />
          )}
          <Title level={2} style={{ margin: 0, width: '100%', insetInlineStart: 0 }} editable={false}>
            {instanceInfo.pathName}
          </Title>
          {/*<Subscribe*/}
          {/*  value={instanceInfo.subscribeEnable}*/}
          {/*  subscribeFrequency={instanceInfo.subscribeFrequency}*/}
          {/*  onChange={handleChangeSubscribe}*/}
          {/*/>*/}
        </Flex>
        <div className="tableWrap">
          <Collapse
            bordered={false}
            collapsible="header"
            activeKey={activeList}
            onChange={(even) => setActiveList(even)}
            expandIcon={({ isActive }) => (
              <CaretRight size={20} style={{ rotate: isActive ? '90deg' : '0deg', transition: '200ms' }} />
            )}
            style={{ background: token.colorBgContainer }}
            items={getItems(panelStyle, instanceInfo)}
          />
        </div>
      </div>
    </div>
  );
};
export default Module;
