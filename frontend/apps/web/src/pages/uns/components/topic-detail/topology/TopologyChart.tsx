import { Flex, App } from 'antd';
import { Fragment, useCallback, useEffect, useMemo, useRef, useState } from 'react';
import type { MouseEvent, PointerEvent as ReactPointerEvent } from 'react';
import styles from './TopologyChart.module.scss';
import { Launch, AddLarge } from '@/components/lucide-icon/carbon';
import { titleIconProps } from '@/components/lucide-icon/icon-props';
import nodeRedIcon from '@/assets/home-icons/node-red.svg';
import postgresql from '@/assets/home-icons/postgresql.svg';
import tdengine from '@/assets/home-icons/tdengine.png';
import timescaleDB from '@/assets/home-icons/timescaleDB.svg';
import { useBaseStore } from '@/stores/base';
import { getSearchParamsString } from '@/utils';
import { bindFlowForUns, createFlow, goFlow } from '@/apis/core-api/flow.ts';
import { useNavigate } from 'react-router';
import { useTranslate } from '@/hooks';
import { getRefreshList, getSourceList } from '@/apis/chat2db';
import Binding from '../binding/Binding.tsx';
import { flowPage } from '@/apis/core-api/flow.ts';
import {
  NodeRedDetail,
  MqttDetail,
  MqttDetail2,
  DataBaseDetail,
  resolveDatabaseSchemaName,
  resolveDatabaseTableName,
} from './Components.tsx';

type TopologyNodeType = 'sourceFlow' | 'eventFlow' | 'mqtt' | 'database';
type TopologyNodeData = {
  active?: boolean;
  id?: string | number;
  flowId?: string | number;
  flowName?: string;
  bindId?: string | number;
  loading?: boolean;
  dataType?: number;
  alias?: string;
  enabled?: boolean;
  onBindingChange?: (_type: string, item: any) => Promise<void>;
};
type TopologyNode = {
  id: TopologyNodeType;
  type: TopologyNodeType;
  data: TopologyNodeData;
};
type TopologyDragState = {
  pointerId: number | null;
  startX: number;
  scrollLeft: number;
};

function SourceFlowNode({ data, readOnly }: { data: TopologyNodeData; readOnly?: boolean }) {
  const formatMessage = useTranslate();
  const configured = data.id || data.flowId || data.flowName;
  const loading = data.loading;
  const statusColor = configured ? '#4CAF50' : '#B1973B';

  return (
    <div className={styles['rf-node-wrap']}>
      <div className={`${styles['rf-node']} ${styles['rf-node-wide']} ${data.active ? styles['rf-node-active'] : ''}`}>
        <img src={nodeRedIcon} alt="" width="28px" />
        <div className={styles['rf-node-text']}>
          <span className={styles['rf-node-subtitle']}>{formatMessage('common.nodeRed', 'Node-Red')}</span>
          <span className={styles['rf-node-title']} title={formatMessage('home.sourceFlow')}>
            {formatMessage('home.sourceFlow')}
          </span>
        </div>
        {!readOnly && (
          <div className={styles['rf-node-actions']}>
            <div data-action="navigate" className={styles['rf-node-btn']}>
              {loading ? (
                <div className={styles['loading-spinner']} />
              ) : configured ? (
                <Launch {...titleIconProps} />
              ) : (
                <AddLarge {...titleIconProps} />
              )}
            </div>
            <div data-action="noNavigate" className={styles['rf-node-btn']}>
              <Binding
                selectValue={data.bindId ? String(data.bindId) : undefined}
                api={flowPage}
                onBinding={(item: any) =>
                  data.onBindingChange ? data.onBindingChange('nodeRed1', item) : Promise.resolve()
                }
              />
            </div>
          </div>
        )}
      </div>
      <div className={styles['rf-indicator']}>
        <span className={styles['rf-indicator-dot']} style={{ background: statusColor }} />
        <span>{formatMessage(configured ? 'common.configured' : 'common.unconfigured')}</span>
      </div>
    </div>
  );
}

function EventFlowNode() {
  const formatMessage = useTranslate();

  return (
    <div className={styles['rf-node-wrap']}>
      <div className={`${styles['rf-node']} ${styles['rf-node-static']}`}>
        <img src={nodeRedIcon} alt="" width="28px" />
        <div className={styles['rf-node-text']}>
          <span className={styles['rf-node-subtitle']}>{formatMessage('common.nodeRed', 'Node-Red')}</span>
          <span className={styles['rf-node-title']}>{formatMessage('home.eventFlow')}</span>
        </div>
      </div>
    </div>
  );
}

function MqttNode({ data }: { data: TopologyNodeData }) {
  const mqttBrokeType = useBaseStore((state) => state.mqttBrokeType);

  return (
    <div className={styles['rf-node-wrap']}>
      <div className={`${styles['rf-node']} ${data.active ? styles['rf-node-active'] : ''}`}>
        <div className={styles['rf-node-text']}>
          <span className={styles['rf-node-subtitle']}>{mqttBrokeType?.toUpperCase() || 'EMQX'}</span>
          <span className={styles['rf-node-title']}>MQTT Broker</span>
        </div>
      </div>
    </div>
  );
}

function DataBaseNode({ data, readOnly }: { data: TopologyNodeData; readOnly?: boolean }) {
  const formatMessage = useTranslate();
  const { dataBaseType, systemInfo } = useBaseStore((state) => ({
    dataBaseType: state.dataBaseType,
    systemInfo: state.systemInfo,
  }));
  const isRelational = [2, 8].includes(Number(data.dataType));
  const statusColor = data.enabled ? '#4CAF50' : '#B1973B';

  return (
    <div className={styles['rf-node-wrap']}>
      <div className={`${styles['rf-node']} ${data.active ? styles['rf-node-active'] : ''}`}>
        <img
          src={isRelational ? postgresql : dataBaseType.includes('tdengine') ? tdengine : timescaleDB}
          alt=""
          width="28px"
        />
        <div className={styles['rf-node-text']}>
          <span className={styles['rf-node-subtitle']}>{isRelational ? 'PostgreSQL' : 'TimescaleDB'}</span>
          <span className={styles['rf-node-title']}>
            {isRelational ? 'Relational DB' : dataBaseType.includes('tdengine') ? 'tdengine' : 'Database'}
          </span>
        </div>
        {!readOnly && isRelational && systemInfo?.containerMap?.chat2db && (
          <div data-action="navigate" className={styles['rf-node-btn']}>
            <Launch {...titleIconProps} />
          </div>
        )}
      </div>
      <div className={styles['rf-indicator']}>
        <span className={styles['rf-indicator-dot']} style={{ background: statusColor }} />
        <span>{formatMessage(data.enabled ? 'common.configured' : 'resourceAction.disabled')}</span>
      </div>
    </div>
  );
}

const emptyTopologyData = {};
const getTopologyKey = (id?: string | number, alias?: string) => String(id ?? alias ?? 'empty');

// ─── 主组件 ───

const TopologyChart = ({ instanceInfo, getFileDetail, readOnly = false }: any) => {
  const [active, setActive] = useState<string>('');
  const [isTopologyDragging, setIsTopologyDragging] = useState(false);
  const topologyDragStateRef = useRef<TopologyDragState>({ pointerId: null, startX: 0, scrollLeft: 0 });
  const didTopologyDragRef = useRef(false);
  const pendingTopologyNodeRef = useRef<TopologyNodeType | null>(null);
  const suppressNextTopologyClickRef = useRef(false);
  const topologyKey = useMemo(
    () => getTopologyKey(instanceInfo?.id, instanceInfo?.alias),
    [instanceInfo?.id, instanceInfo?.alias]
  );
  const [datasState, setDatasState] = useState<{ topologyKey: string; data: any }>({
    topologyKey: '',
    data: emptyTopologyData,
  });
  const datas = datasState.topologyKey === topologyKey ? datasState.data : emptyTopologyData;
  const navigate = useNavigate();
  const { message } = App.useApp();
  const formatMessage = useTranslate();

  const refreshNodeRedData = useCallback(async (alias: string, unsId?: string | number) => {
    try {
      const result = await goFlow(alias, unsId);
      setDatasState({ topologyKey: getTopologyKey(unsId, alias), data: result || emptyTopologyData });
      return result;
    } catch (error) {
      console.error('Error fetching topology data:', error);
      return null;
    }
  }, []);

  const onBindingChange = useCallback(
    (_type: string, item: any) => {
      return bindFlowForUns({
        flowId: item.id,
        unsAlias: instanceInfo.alias,
        unsId: instanceInfo.id,
      }).then(() => {
        message.success(formatMessage('common.optsuccess'));
        getFileDetail(instanceInfo.id);
        refreshNodeRedData(instanceInfo.alias, instanceInfo.id);
      });
    },
    [instanceInfo, message, formatMessage, getFileDetail, refreshNodeRedData]
  );

  useEffect(() => {
    let isCurrent = true;
    if (!instanceInfo?.alias || !instanceInfo?.id) {
      return () => {
        isCurrent = false;
      };
    }

    goFlow(instanceInfo.alias, instanceInfo.id)
      .then((result) => {
        if (isCurrent) {
          setDatasState({ topologyKey, data: result || emptyTopologyData });
        }
      })
      .catch((error) => {
        if (isCurrent) {
          console.error('Error fetching topology data:', error);
        }
      });

    return () => {
      isCurrent = false;
    };
  }, [instanceInfo?.alias, instanceInfo?.id, topologyKey]);

  const nodes: TopologyNode[] = useMemo(() => {
    const dataType = instanceInfo?.dataType;
    const result: TopologyNode[] = [];
    const hasSourceNode = dataType !== 3;

    if (hasSourceNode) {
      result.push({
        id: 'sourceFlow',
        type: 'sourceFlow',
        data: {
          active: active === 'sourceFlow',
          id: datas?.id || '',
          flowId: datas?.flowId || '',
          flowName: datas?.flowName || '',
          bindId: datas?.id || '',
          loading: false,
          onBindingChange,
        },
      });
    }

    result.push({
      id: 'mqtt',
      type: 'mqtt',
      data: { active: active === 'mqtt' },
    });

    result.push({
      id: 'database',
      type: 'database',
      data: {
        active: active === 'database',
        dataType: instanceInfo?.dataType,
        alias: instanceInfo?.alias,
        enabled: Boolean(instanceInfo?.persistence || Number(instanceInfo?.enableHistory) === 1),
      },
    });

    result.push({
      id: 'eventFlow',
      type: 'eventFlow',
      data: {},
    });

    return result;
  }, [instanceInfo, datas, active, onBindingChange]);

  // 节点点击处理
  // 节点内按钮点击（通过 data-action 区分）
  const handleNavigateClick = useCallback(
    (nodeId: string, nodeData: any) => {
      if (nodeId === 'sourceFlow') {
        const d = nodeData;
        if (d.id || d.flowId || d.flowName) {
          navigate(
            `/collection-flow/flow-editor?${getSearchParamsString({
              id: d.id,
              name: d.flowName,
              flowId: d.flowId,
              from: location.pathname,
            })}`
          );
        } else if (instanceInfo?.id && instanceInfo?.alias && instanceInfo?.path) {
          createFlow({
            unsAlias: instanceInfo?.alias,
            path: instanceInfo?.path,
            unsId: instanceInfo?.id,
            unsNodeIds: [Number(instanceInfo.id)],
            mockData: true,
            mockTopic: instanceInfo?.topic || instanceInfo?.path || instanceInfo?.namespace || instanceInfo?.alias,
            mockFields: instanceInfo?.fields || instanceInfo?.jsonFields || [],
            mockTriggerMode: 'manual',
          })
            .then((res: any) => {
              if (res) {
                setDatasState({ topologyKey, data: res || emptyTopologyData });
                navigate(
                  `/collection-flow/flow-editor?${getSearchParamsString({
                    id: res.id,
                    name: res.flowName,
                    flowId: res.flowId,
                    from: location.pathname,
                  })}`
                );
              }
            })
            .catch(() => {});
        }
      } else if (nodeId === 'database' && nodeData.dataType === 2) {
        getSourceList().then((data: any) => {
          const sourceData = data?.data?.data?.find((i: any) => i.alias === '@postgresql');
          const loadData = (params: any) => {
            getRefreshList(params).then((res: any) => {
              if (res.hasNextPage) {
                loadData({ dataSourceId: sourceData?.id, pageNo: res.data?.pageNo + 1 });
              } else {
                const schemaName = resolveDatabaseSchemaName(instanceInfo);
                const tableName = resolveDatabaseTableName(instanceInfo);
                navigate(
                  `/SQLEditor?${getSearchParamsString({
                    dataSourceName: '@postgresql',
                    databaseName: 'postgres',
                    databaseType: 'POSTGRESQL',
                    schemaName,
                    tableName,
                  })}`
                );
              }
            });
          };
          loadData({ dataSourceId: sourceData?.id });
        });
      }
    },
    [instanceInfo, navigate, topologyKey]
  );

  const selectTopologyNode = useCallback((node: TopologyNode) => {
    if (node.id === 'eventFlow') {
      return;
    }

    setActive((prev) => (prev === node.id ? '' : node.id));
  }, []);

  const handleNodeClick = useCallback(
    (event: MouseEvent<HTMLDivElement>, node: TopologyNode) => {
      const target = event.target as HTMLElement;

      // 点击绑定按钮区域，不做跳转
      if (target.closest('[data-action="noNavigate"]')) {
        return;
      }

      // 只有点击跳转按钮才触发导航
      if (target.closest('[data-action="navigate"]')) {
        handleNavigateClick(node.id, node.data);
        return;
      }

      // 其他区域：切换选中节点详情
      selectTopologyNode(node);
    },
    [handleNavigateClick, selectTopologyNode]
  );

  const renderNode = useCallback(
    (node: TopologyNode) => {
      const nodeContent = (() => {
        switch (node.type) {
          case 'sourceFlow':
            return <SourceFlowNode data={node.data} readOnly={readOnly} />;
          case 'mqtt':
            return <MqttNode data={node.data} />;
          case 'database':
            return <DataBaseNode data={node.data} readOnly={readOnly} />;
          case 'eventFlow':
            return <EventFlowNode />;
          default:
            return null;
        }
      })();

      return (
        <div
          className={styles['topology-node-item']}
          data-topology-node-id={node.id}
          onClick={(event) => handleNodeClick(event, node)}
        >
          {nodeContent}
        </div>
      );
    },
    [handleNodeClick, readOnly]
  );

  const handleTopologyPointerDown = useCallback((event: ReactPointerEvent<HTMLDivElement>) => {
    if (event.button !== 0) return;

    const target = event.target as HTMLElement;
    pendingTopologyNodeRef.current = null;
    if (target.closest('[data-action]')) return;

    const nodeElement = target.closest('[data-topology-node-id]') as HTMLElement | null;
    pendingTopologyNodeRef.current = (nodeElement?.dataset.topologyNodeId as TopologyNodeType | undefined) || null;

    topologyDragStateRef.current = {
      pointerId: event.pointerId,
      startX: event.clientX,
      scrollLeft: event.currentTarget.scrollLeft,
    };
    didTopologyDragRef.current = false;
    setIsTopologyDragging(true);
    event.currentTarget.setPointerCapture(event.pointerId);
  }, []);

  const handleTopologyPointerMove = useCallback((event: ReactPointerEvent<HTMLDivElement>) => {
    const dragState = topologyDragStateRef.current;
    if (dragState.pointerId !== event.pointerId) return;

    const deltaX = event.clientX - dragState.startX;
    if (Math.abs(deltaX) > 3) {
      didTopologyDragRef.current = true;
    }
    if (didTopologyDragRef.current) {
      event.preventDefault();
      event.currentTarget.scrollLeft = dragState.scrollLeft - deltaX;
    }
  }, []);

  const handleTopologyPointerEnd = useCallback(
    (event: ReactPointerEvent<HTMLDivElement>) => {
      const dragState = topologyDragStateRef.current;
      if (dragState.pointerId !== event.pointerId) return;

      if (!didTopologyDragRef.current && pendingTopologyNodeRef.current) {
        const pendingNode = nodes.find((node) => node.id === pendingTopologyNodeRef.current);
        if (pendingNode) {
          selectTopologyNode(pendingNode);
          suppressNextTopologyClickRef.current = true;
          window.setTimeout(() => {
            suppressNextTopologyClickRef.current = false;
          }, 0);
        }
      }

      pendingTopologyNodeRef.current = null;
      if (event.currentTarget.hasPointerCapture(event.pointerId)) {
        event.currentTarget.releasePointerCapture(event.pointerId);
      }
      topologyDragStateRef.current = { pointerId: null, startX: 0, scrollLeft: 0 };
      setIsTopologyDragging(false);
    },
    [nodes, selectTopologyNode]
  );

  const handleTopologyClickCapture = useCallback((event: MouseEvent<HTMLDivElement>) => {
    if (suppressNextTopologyClickRef.current) {
      event.preventDefault();
      event.stopPropagation();
      suppressNextTopologyClickRef.current = false;
      didTopologyDragRef.current = false;
      return;
    }

    if (!didTopologyDragRef.current) return;

    event.preventDefault();
    event.stopPropagation();
    didTopologyDragRef.current = false;
  }, []);

  return (
    <Flex vertical className={styles['topology-wrap']}>
      <div
        className={`${styles['topology-content']} ${isTopologyDragging ? styles['topology-content-dragging'] : ''}`}
        onClickCapture={handleTopologyClickCapture}
        onPointerDown={handleTopologyPointerDown}
        onPointerMove={handleTopologyPointerMove}
        onPointerUp={handleTopologyPointerEnd}
        onPointerCancel={handleTopologyPointerEnd}
      >
        <div className={styles['topology-flow']}>
          {nodes.map((node, index) => (
            <Fragment key={node.id}>
              {renderNode(node)}
              {index < nodes.length - 1 && (
                <div
                  className={`${styles['topology-link']} ${
                    node.type === 'sourceFlow' ? styles['topology-link-after-wide'] : ''
                  }`}
                  aria-hidden
                />
              )}
            </Fragment>
          ))}
        </div>
      </div>
      {/* 节点详情面板 */}
      {active === 'sourceFlow' && (
        <div className={styles['topology-detail']}>
          <NodeRedDetail flowList={datas} />
        </div>
      )}
      {active === 'mqtt' && instanceInfo.dataType !== 3 && (
        <div className={styles['topology-detail']}>
          <MqttDetail />
        </div>
      )}
      {active === 'mqtt' && instanceInfo.dataType === 3 && (
        <div className={styles['topology-detail']}>
          <MqttDetail2 instanceInfo={instanceInfo} />
        </div>
      )}
      {active === 'database' && (
        <div className={styles['topology-detail']}>
          <DataBaseDetail instanceInfo={instanceInfo} />
        </div>
      )}
    </Flex>
  );
};

export default TopologyChart;
