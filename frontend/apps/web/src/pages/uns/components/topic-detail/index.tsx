import { useState, useEffect, type FC, useRef } from 'react';
import {
  createDashboard,
  getInstanceInfo,
  modifyModel,
  // updateModelSubscribe,
  checkDashboardIsExist,
  getDashboardByUns,
} from '@/apis/inter-api/uns';
import { Button, Collapse, Flex, theme, Typography, App, Space } from 'antd';
import Icon, { FullscreenOutlined } from '@ant-design/icons';
import { CaretRight, Document, Code, TableSplit, SendAlt, ChartLine } from '@carbon/icons-react';
import { useTranslate } from '@/hooks';
import type { CSSProperties } from 'react';
import type { CollapseProps } from 'antd';
import Details from './Details';
import TopologyChart from './topology/TopologyChart';
import Definition from './Definition';
import Payload from './Payload';
import Dashboard from './Dashboard';
import RawData from './RawData';
// import SqlQuery from './SqlQuery';
import DocumentList from '@/pages/uns/components/DocumentList.tsx';
import UploadButton from '@/pages/uns/components/UploadButton.tsx';
import { ButtonPermission } from '@/common-types/button-permission.ts';
import { useMediaSize } from '@/hooks';
import EditDetailButton from '@/pages/uns/components/EditDetailButton';
import type { InitTreeDataFnType, UnsTreeNode, FieldItem } from '@/pages/uns/types';
import FileEdit from '@/components/svg-components/FileEdit';
import { hasPermission, getToken } from '@/utils/auth';
import { isJsonString } from '@/utils/common';
import { useBaseStore } from '@/stores/base';
// import Subscribe from '@/pages/uns/components/subscribe';
import EditButton from '@/pages/uns/components/EditButton.tsx';
import screenfull from 'screenfull';
import { CustomAxiosConfigEnum } from '@/utils';
import useSSE from '@/hooks/useSSE.ts';

const { Title } = Typography;

export interface FileDetailProps {
  currentNode: UnsTreeNode;
  initTreeData: InitTreeDataFnType;
  handleDelete: (node: UnsTreeNode) => void;
}

interface InstanceInfoType {
  [key: string]: any;
}

const Module: FC<FileDetailProps> = (props) => {
  const {
    currentNode: { id },
    initTreeData,
  } = props;
  const {
    systemInfo: { qualityName = 'quality', timestampName = 'timeStamp', useAliasPathAsTopic, enableAutoCategorization },
  } = useBaseStore((state) => ({
    systemInfo: state.systemInfo,
  }));
  const [createLoading, setCreateLoading] = useState(false);
  const formatMessage = useTranslate();
  const { message } = App.useApp();
  const documentListRef = useRef(null);
  const [instanceInfo, setInstanceInfo] = useState<InstanceInfoType>({});
  const [dashboardInfo, setDashboardInfo] = useState<any>();
  const [activeList, setActiveList] = useState<string[]>([
    'topologyChart',
    'definition',
    'payload',
    'dashboard',
    // 'sqlQuery',
  ]);
  const [showPayloadTable, setShowPayloadTable] = useState<boolean>(true);
  const [websocketData, setWebsocketData] = useState<any>({});
  const [wsTimeStamp, setWsTimeStamp] = useState<number>(0);
  const { token } = theme.useToken();

  const panelStyle: CSSProperties = {
    background: 'val(--supos-bg-color)',
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

  useSSE(
    instanceInfo.id && wsTimeStamp
      ? `/inter-api/supos/uns/newMsg?id=${instanceInfo.id}&t=${wsTimeStamp}&token=${getToken()}`
      : '',
    {
      onMessage: (event) => {
        const dataJson = event.data;
        if (isJsonString(dataJson)) {
          const data = JSON.parse(dataJson);
          console.log(data);
          if (qualityName && data?.data?.[qualityName]) {
            //质量码做特殊处理
            data.data[qualityName] = longToJavaHex(data.data[qualityName]);
          }
          if (!isJsonString(data.payload)) {
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
    }
  );

  useEffect(() => {
    setWebsocketData({});
    if (id) {
      getFileDetail(id as string);
    } else {
      setInstanceInfo({});
      setWsTimeStamp(0);
    }
  }, [id]);

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
          if (JSON.stringify(data.fields) !== JSON.stringify(instanceInfo.fields)) {
            setWsTimeStamp(Date.now());
          }
          // data.extendFieldUsed = data.mount
          //   ? ['unit', 'upperLimit', 'lowerLimit', 'decimal']
          //   : data.extendFieldUsed || [];
          data.extendFieldUsed = [];
          if (data.withDashboard) {
            try {
              const { code } = await checkDashboardIsExist({ alias: data.alias });
              data.dashboardIsExist = code === 200;
            } catch (err) {
              console.error(err);
            }
          }
          setInstanceInfo(data);
          return getDashboardByUns(data?.alias, { [CustomAxiosConfigEnum.NoMessage]: true }).then((data) => {
            setDashboardInfo(data || { type: 0 });
            return data;
          });
        }
      })
      .catch(() => {});
  };

  const getItems: (
    panelStyle: CSSProperties,
    instanceInfo: InstanceInfoType,
    dashboardType: number
  ) => CollapseProps['items'] = (panelStyle, instanceInfo) => {
    const items = [
      {
        key: 'detail',
        label: formatMessage('common.detail'),
        children: (
          <Details instanceInfo={instanceInfo} updateTime={websocketData?.updateTime} websocketData={websocketData} />
        ),
        style: panelStyle,
        extra: (
          <EditDetailButton
            auth={ButtonPermission['uns.fileDetail']}
            modelInfo={instanceInfo}
            getModel={() => getFileDetail(id as string)}
          />
        ),
      },
      {
        key: [1, 2, 3, 6, 7, 8].includes(instanceInfo.dataType) ? 'definition' : '',
        label: formatMessage('uns.definition'),
        children: <Definition instanceInfo={instanceInfo} />,
        style: panelStyle,
        extra:
          [1, 2, 8].includes(instanceInfo.dataType) && !instanceInfo.mount && !instanceInfo.modelId ? (
            <EditButton
              auth={ButtonPermission['uns.fileDetail']}
              modelInfo={instanceInfo}
              getModel={() => getFileDetail(id as string)}
              editType="file"
            />
          ) : null,
      },
      {
        key: [1, 2, 3, 6, 7, 8].includes(instanceInfo.dataType) ? 'payload' : '',
        label: formatMessage('uns.payload'),
        children:
          instanceInfo.dataType === 8 || !showPayloadTable ? (
            <RawData payload={websocketData?.data} />
          ) : (
            <Payload websocketData={websocketData} fields={instanceInfo.fields || []} />
          ),
        style: panelStyle,
        extra:
          instanceInfo.dataType === 8 ? null : (
            <Button
              style={{ border: '1px solid #C6C6C6' }}
              color="default"
              variant="filled"
              icon={showPayloadTable ? <Code /> : <TableSplit />}
              onClick={() => setShowPayloadTable(!showPayloadTable)}
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
              // extra: (
              //   <Space>
              //     <DashboardBinding
              //       key={instanceInfo?.id}
              //       selectValue={dashboardInfo?.id}
              //       isCreated={instanceInfo.withDashboard || instanceInfo.dashboardIsExist}
              //       onCreated={handleCreateDashboard}
              //       onBinding={(item: any) => {
              //         return bindDashboardForUns({
              //           dashboardId: item.id,
              //           unsAlias: instanceInfo.alias,
              //         }).then(() => {
              //           message.success(formatMessage('common.optsuccess'));
              //           getFileDetail(instanceInfo.id).then((dashboardInfo) => {
              //             navigate(
              //               `/dashboards/preview?${getSearchParamsString({ id: dashboardInfo.id, type: dashboardInfo.type, status: 'preview', name: dashboardInfo.name })}`
              //             );
              //           });
              //         });
              //       }}
              //     />
              //     {(instanceInfo.withDashboard || instanceInfo.dashboardIsExist) && (
              //       <Button
              //         title={formatMessage('common.fullScreen')}
              //         icon={<FullscreenOutlined />}
              //         onClick={() => {
              //           if (screenfull.isEnabled) {
              //             const el = document.getElementById('dashboardIframe');
              //             if (el) {
              //               screenfull.request(el);
              //             } else {
              //               message.error('未找到全屏元素');
              //             }
              //           } else {
              //             message.error('该浏览器,不支持全屏功能');
              //           }
              //         }}
              //       />
              //     )}
              //   </Space>
              // ),
              extra: instanceInfo.dataType !== 7 && (
                <Space.Compact block>
                  {(!instanceInfo.withDashboard || !instanceInfo.dashboardIsExist) && (
                    <Button loading={createLoading} onClick={handleCreateDashboard}>
                      {formatMessage('common.create')}
                    </Button>
                  )}
                  {(instanceInfo.withDashboard || instanceInfo.dashboardIsExist) && (
                    <Button
                      title={formatMessage('common.fullScreen')}
                      icon={<FullscreenOutlined />}
                      onClick={() => {
                        if (screenfull.isEnabled) {
                          const el = document.getElementById('dashboardIframe');
                          if (el) {
                            screenfull.request(el);
                          } else {
                            message.error('未找到全屏元素');
                          }
                        } else {
                          message.error('该浏览器,不支持全屏功能');
                        }
                      }}
                    />
                  )}
                </Space.Compact>
              ),
            },
            {
              key: [1, 2, 8].includes(instanceInfo.dataType) ? 'topologyChart' : '',
              label: formatMessage('uns.topology'),
              children: (
                <TopologyChart
                  getFileDetail={getFileDetail}
                  instanceInfo={instanceInfo}
                  dashboardInfo={dashboardInfo}
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
        children: <DocumentList alias={instanceInfo.alias} ref={documentListRef} />,
        style: panelStyle,
        extra: (
          <UploadButton
            setActiveList={setActiveList}
            auth={ButtonPermission['uns.fileDetail']}
            alias={instanceInfo.alias}
            documentListRef={documentListRef}
          />
        ),
      },
    ];
    return items.filter((item: any) => item.key);
  };

  const handleCreateDashboard = () => {
    setCreateLoading(true);
    return createDashboard(instanceInfo.alias)
      .then(() => {
        message.success(formatMessage('common.optsuccess'));
        getFileDetail(instanceInfo.id);
      })
      .finally(() => {
        setCreateLoading(false);
      });
  };
  // const handleChangeSubscribe = async (enable: boolean, frequency?: string) => {
  //   await updateModelSubscribe({ id, enable, frequency });
  //   getFileDetail(id as string);
  //   message.success(enable ? formatMessage('uns.subscribeSuccessful') : formatMessage('uns.unsubscribeSuccessful'));
  // };

  const getFileIcon = () => {
    switch (instanceInfo.parentDataType) {
      case 1:
        return <Document size={20} style={{ color: '#D2A106' }} />;
      case 2:
        return <SendAlt size={20} style={{ color: '#94C518' }} />;
      case 3:
        return <ChartLine size={20} style={{ color: '#1D77FE' }} />;
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
            <Flex
              align="center"
              justify="center"
              style={{ width: 36, height: 36, background: '#f4f4f4', borderRadius: 3 }}
            >
              {getFileIcon()}
            </Flex>
          ) : (
            <Document size={20} />
          )}
          <Title
            level={2}
            style={{ margin: 0, width: '100%', insetInlineStart: 0 }}
            editable={
              hasPermission(ButtonPermission['uns.fileDetail']) && useAliasPathAsTopic
                ? {
                    icon: (
                      <Icon
                        data-button-auth={ButtonPermission['uns.fileDetail']}
                        component={FileEdit}
                        style={{
                          fontSize: 25,
                          color: 'var(--supos-text-color)',
                        }}
                      />
                    ),
                    onChange: (val) => {
                      if (val === instanceInfo.pathName || !val) return;
                      if (val.length > 63) {
                        return message.warning(
                          formatMessage('uns.labelMaxLength', { label: formatMessage('common.name'), length: 63 })
                        );
                      }
                      if (!/^[\u4e00-\u9fa5a-zA-Z0-9_-]+$/.test(val)) {
                        return message.warning(formatMessage('uns.nameFormat'));
                      }
                      modifyModel({ id, name: val }).then(() => {
                        message.success(formatMessage('uns.editSuccessful'));
                        getFileDetail(id as string);
                        initTreeData({ queryType: 'editFileName' });
                      });
                    },
                  }
                : false
            }
          >
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
            items={getItems(panelStyle, instanceInfo, dashboardInfo)}
          />
        </div>
      </div>
    </div>
  );
};
export default Module;
