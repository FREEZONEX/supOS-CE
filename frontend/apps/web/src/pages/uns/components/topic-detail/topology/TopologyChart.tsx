import { App, Flex } from 'antd';
import { useEffect, useRef, useState } from 'react';
import { register } from '@antv/x6-react-shape';
import styles from './TopologyChart.module.scss';
import { TypeEnum } from './types.ts';
import {
  Apps,
  ButtonError,
  DataBase,
  DataBaseDetail,
  Mqtt,
  MqttDetail,
  MqttDetail2,
  NodeRed,
  NodeRedDetail,
} from './Components.tsx';
import ReactDOM from 'react-dom/client'; // React 18 使用 'react-dom/client'
import { Graph } from '@antv/x6';
import { debounce } from 'lodash-es';
import { findDate } from '@/pages/uns/components/topic-detail/topology/data.ts';
import { useBaseStore } from '@/stores/base';
import { bindDashboardForUns } from '@/apis/inter-api';
import { getSearchParamsString } from '@/utils';
import { bindFlowForUns, createFlow, goFlow } from '@/apis/inter-api/flow.ts';
import { useNavigate } from 'react-router';
import { useTranslate } from '@/hooks';
import { useDeepCompareEffect } from 'ahooks';
import { getRefreshList, getSourceList } from '@/apis/chat2db';

register({
  shape: TypeEnum.NodeRed,
  width: 220,
  height: 52,
  component: NodeRed,
});

register({
  shape: TypeEnum.Mqtt,
  width: 150,
  height: 52,
  component: Mqtt,
});

register({
  shape: TypeEnum.DataBase,
  width: 190,
  height: 52,
  component: DataBase,
});

register({
  shape: TypeEnum.Apps,
  width: 210,
  height: 52,
  component: Apps,
});

const TopologyChart = ({ instanceInfo, dashboardInfo, getFileDetail }: any) => {
  const topologyContainerRef = useRef<any>(null);
  const topologyRef = useRef<Graph>(undefined);
  const dashboardType = useBaseStore((state) => state.dashboardType);
  const modeState = useRef<any>([]);
  const [active, setActive] = useState<any>('');
  const navigate = useNavigate();
  const [datas, setDatas] = useState<any>({});
  const { message } = App.useApp();
  const interls = useRef<any>(null);
  const formatMessage = useTranslate();

  const initTopology = () => {
    if (topologyRef.current) return topologyRef.current;
    const topologyInstance = new Graph({
      container: topologyContainerRef.current,
      background: { color: 'var(--supos-gray-color-10-message)' },
      interacting: false,
      panning: true,
      mousewheel: { enabled: true, modifiers: ['ctrl', 'meta'] },
      scaling: { min: 0.05, max: 12 },
    });

    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-expect-error
    topologyInstance.options.onEdgeLabelRendered = (args: any) => {
      const { selectors, edge } = args; // 获取edge对象
      const content = selectors.foContent as HTMLDivElement;
      if (content) {
        const root = ReactDOM.createRoot(content);
        const nodeStatu = modeState.current?.find((item: any) => item.topologyNode === edge.id);
        root.render(<ButtonError nodeStatu={nodeStatu} />); // 渲染组件
      }
    };

    topologyInstance.positionPoint({ x: 210, y: 0 }, 40, '40%');
    topologyRef.current = topologyInstance;
    return topologyInstance;
  };

  const fetchNodeRedData = async (alias: string) => {
    try {
      const result = await goFlow(alias);
      setDatas(result || {});
      return result;
    } catch (error) {
      console.error('Error fetching topology data:', error);
      return null;
    }
  };

  const onBindingChange = (type: string, item: any) => {
    if (type === 'apps1') {
      return bindDashboardForUns({
        dashboardId: item.id,
        unsAlias: instanceInfo.alias,
      }).then(() => {
        message.success(formatMessage('common.optsuccess'));
        getFileDetail(instanceInfo.id);
      });
    } else {
      return bindFlowForUns({
        flowId: item.id,
        unsAlias: instanceInfo.alias,
      }).then(() => {
        message.success(formatMessage('common.optsuccess'));
        getFileDetail(instanceInfo.id);
        fetchNodeRedData(instanceInfo.alias);
      });
    }
  };

  const updateTopologyData = (instanceInfo: any) => {
    const { nodes, edges } = findDate({ ...instanceInfo, dashboardType }) || { nodes: [], edges: [] };
    modeState.current = [];
    if (topologyRef.current) {
      topologyRef.current.clearCells();
      // 添加节点和边
      nodes.forEach((node) => {
        if (node.id === TypeEnum.NodeRed) {
          node.data.onBindingChange = onBindingChange;
          node.data.active = true;
          setActive(TypeEnum.NodeRed);
        } else if (node.id === TypeEnum.Apps) {
          node.data.onBindingChange = onBindingChange;
        }
        topologyRef.current?.addNode(node);
      });
      edges.forEach((edge) => {
        topologyRef.current?.addEdge(edge);
      });
      // 重新设置edge的位置
      const newEdges = topologyRef.current.getEdges();
      newEdges.forEach((edge: any) => {
        const sourceCell = edge.getSourceCell();
        const targetCell = edge.getTargetCell();

        if (sourceCell && targetCell) {
          const sourceBBox = sourceCell.getBBox();
          const targetBBox = targetCell.getBBox();
          const sourcePoint = {
            x: sourceBBox.x + sourceBBox.width,
            y: sourceBBox.y + sourceBBox.height / 2,
          };
          const targetPoint = {
            x: targetBBox.x,
            y: targetBBox.y + targetBBox.height / 2,
          };
          // 使用绝对坐标设置连接点
          edge.setSource({
            x: sourcePoint.x,
            y: sourcePoint.y,
          });

          edge.setTarget({
            x: targetPoint.x,
            y: targetPoint.y,
          });
        }
      });
    }
  };

  useEffect(() => {
    initTopology();
    const handleResize = debounce(() => {
      if (topologyRef.current) {
        const ww = document.getElementsByClassName('treemapTitle')[0].clientWidth;
        const width = window.innerWidth - ww - 50; // 宽度减去侧边栏宽度
        const height = 200; // 计算容器高度
        topologyRef.current.resize(width, height);
      }
    }, 200); // 防抖 200 毫秒
    window.addEventListener('resize', handleResize);

    // 清理函数
    return () => {
      window.removeEventListener('resize', handleResize);
      if (topologyRef.current) {
        topologyRef.current?.clearCells();
        topologyRef.current = undefined;
      }
    };
  }, []);

  useDeepCompareEffect(() => {
    if (topologyRef.current) {
      topologyRef.current.dispose?.();
      topologyRef.current = undefined;
    }
    // 请求nodered
    initTopology();
    updateTopologyData(instanceInfo);
    const nodeClickFn = ({ cell, e }: any) => {
      const target = e.target as HTMLElement;
      const launchButton = target.closest('[data-action="navigate"]');

      if (target.closest('[data-action="noNavigate"]')) {
        e.stopPropagation();
        return;
      }
      if (cell?.id === TypeEnum.NodeRed && launchButton) {
        if (cell.data.id || cell.data.flowId || cell.data.flowName) {
          navigate(
            `/collection-flow/flow-editor?${getSearchParamsString({
              id: cell.data.id,
              name: cell.data.flowName,
              status: cell.data.flowStatus,
              flowId: cell.data.flowId,
              from: location.pathname,
            })}`
          );
        } else {
          if (cell.data.loading) {
            return;
          }
          if (instanceInfo?.alias && instanceInfo?.path) {
            // 设置节点loading状态
            const node = topologyRef.current!.getCellById(cell.id);
            node.setData({
              ...node.data,
              loading: true,
            });
            createFlow({ unsAlias: instanceInfo?.alias, path: instanceInfo?.path })
              .then((res: any) => {
                if (res) {
                  setDatas(res || {});
                  navigate(
                    `/collection-flow/flow-editor?${getSearchParamsString({
                      id: res.id,
                      name: res.flowName,
                      flowId: res.flowId,
                      from: location.pathname,
                    })}`
                  );
                }
                // 清除节点loading状态
                node.setData({
                  ...node.data,
                  loading: false,
                });
                return res;
              })
              .catch(() => {
                // 清除节点loading状态
                node.setData({
                  ...node.data,
                  loading: false,
                });
              });
          }
        }
        return;
      }

      if (cell?.id === TypeEnum.DataBase && cell.data.dataType === 2 && launchButton) {
        getSourceList().then((data: any) => {
          const sourceData = data?.data?.data?.find((i: any) => i.alias === '@postgresql');
          const loadData = (params: any) => {
            getRefreshList(params).then((res: any) => {
              if (res.hasNextPage) {
                loadData({
                  dataSourceId: sourceData?.id,
                  pageNo: res.data?.pageNo + 1,
                });
              } else {
                navigate(
                  `/SQLEditor?dataSourceName=@postgresql&databaseName=postgres&databaseType=POSTGRESQL&schemaName=public&tableName=${instanceInfo?.alias}`
                );
              }
            });
          };
          loadData({ dataSourceId: sourceData?.id });
        });
        return;
      }
      if (cell?.id === TypeEnum.Apps && dashboardType?.includes('grafana') && launchButton) {
        if (!dashboardInfo) {
          message.error(formatMessage('common.dashboardNotFound'));
          return;
        }
        // navigate('/grafana-design', { state: { url: getAppsLink(dashboardInfo), name: 'GrafanaDesign' } });
        navigate(
          `/dashboards/preview?${getSearchParamsString({ id: dashboardInfo.id, type: dashboardInfo.type, status: 'preview', name: dashboardInfo.name })}`
        );
        return;
      }
      setActive((active: any) => (active === cell.id ? '' : cell.id));
      const node = topologyRef.current!.getCellById(cell.id);
      if (node.data.active) {
        // 清空所有节点的选中状态
        topologyRef.current!.getNodes().forEach((node: any) => {
          node.setData({
            active: false,
          });
        });
      } else {
        // 清空所有节点的选中状态
        topologyRef.current!.getNodes().forEach((node: any) => {
          node.setData({
            active: false,
          });
        });
        // 给节点添加选中样式
        node.setData({
          active: true,
        });
      }
    };
    // 为边绑定点击事件
    const edgeClickFn = ({ cell }: any) => {
      setActive((active: any) => (active === cell.id ? '' : cell.id));
      clearInterval(interls.current);
      // getTopologyState();
      const xx = modeState.current?.filter((item: any) => item.topologyNode == cell.id && item.eventCode != 0) || [];
      console.log('click', modeState.current, cell, xx);
      modeState.current = [];
    };
    // 设置事件监听器
    topologyRef.current!.on('node:click', nodeClickFn);
    topologyRef.current!.on('edge:click', edgeClickFn);
    interls.current = setInterval(() => {
      // getTopologyStateData();
    }, 2000);
    // 清理事件监听器
    return () => {
      clearInterval(interls.current);
      interls.current = null;
      if (topologyRef.current) {
        topologyRef.current.off('node:click', nodeClickFn);
        topologyRef.current.off('edge:click', edgeClickFn);
      }
    };
  }, [instanceInfo, dashboardInfo?.id]);

  useDeepCompareEffect(() => {
    fetchNodeRedData(instanceInfo.alias);
  }, [instanceInfo]);

  useEffect(() => {
    const node = topologyRef.current!.getCellById(TypeEnum.NodeRed);
    if (node) {
      // 将nodered的值设置进去
      node.setData({
        id: datas?.id || '',
        flowStatus: datas?.flowStatus || '',
        flowId: datas?.flowId || '',
        flowName: datas?.flowName || '',
        bindId: datas?.id || '',
      });
    }
    if (topologyRef.current) {
      const node = topologyRef.current!.getCellById(TypeEnum.DataBase);
      if (node) {
        // 将DataBase的值设置进去
        node.setData({
          dataType: instanceInfo?.dataType,
          alias: instanceInfo?.alias,
        });
      }
    }
  }, [datas]);

  useDeepCompareEffect(() => {
    if (topologyRef.current) {
      const node = topologyRef.current!.getCellById(TypeEnum.Apps);
      if (node) {
        // 将apps的值设置进去
        node.setData({
          bindId: dashboardInfo?.id,
          subtitle: dashboardInfo?.type === 2 ? 'fuxa' : 'Grafana',
        });
      }
    }
  }, [dashboardInfo]);

  useDeepCompareEffect(() => {
    // 渲染拓扑图
  }, [datas, dashboardInfo?.id]);
  return (
    <Flex vertical wrap className={styles['topology-wrap']}>
      {/*  拓扑图 */}
      <div ref={topologyContainerRef} className={styles['topology-content']}></div>
      {/*  节点详情 */}
      {[TypeEnum.NodeRed, TypeEnum.DataBase, TypeEnum.Mqtt].includes(active) && (
        <div className={styles['topology-detail']}>
          {active == TypeEnum.NodeRed ? <NodeRedDetail flowList={datas} /> : ''}
          {active == TypeEnum.Mqtt && instanceInfo.dataType !== 3 ? <MqttDetail /> : ''}
          {active == TypeEnum.Mqtt && instanceInfo.dataType === 3 ? <MqttDetail2 instanceInfo={instanceInfo} /> : ''}
          {active == TypeEnum.DataBase ? <DataBaseDetail instanceInfo={instanceInfo} /> : ''}
        </div>
      )}
    </Flex>
  );
};

export default TopologyChart;
