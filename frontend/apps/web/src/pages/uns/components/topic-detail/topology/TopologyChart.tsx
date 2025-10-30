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
import { debounce } from 'lodash';
import { data, findDate, markupLine } from '@/pages/uns/components/topic-detail/topology/data.ts';
import { useBaseStore } from '@/stores/base';
import { bindDashboardForUns, getTopologyStatus } from '@/apis/inter-api';
import { getSearchParamsString } from '@/utils';
import { bindFlowForUns, createFlow, goFlow } from '@/apis/inter-api/flow.ts';
import { useNavigate } from 'react-router';
import { useTranslate } from '@/hooks';
import { useDeepCompareEffect } from 'ahooks';
import { getRefreshList, getSourceList } from '@/apis/chat2db';
import error from '@/assets/uns/error.svg';

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

const orderList = ['pushOriginalData', 'pushMqtt', 'pullMqttOrDataPersistence']; //自定义循序
const xxIndexMap = orderList.reduce((acc: any, cur: any, index: number) => {
  acc[cur] = index;
  return acc;
}, {});
const compareByXxOrder = (a: any, b: any) => {
  const indexA = xxIndexMap[a.topologyNode];
  const indexB = xxIndexMap[b.topologyNode];
  return indexA - indexB;
};

const TopologyChart = ({ instanceInfo, dashboardInfo, getFileDetail }: any) => {
  const topologyContainerRef = useRef<any>(null);
  const topologyRef = useRef<Graph>(undefined);
  const dashboardType = useBaseStore((state) => state.dashboardType);
  const modeState = useRef<any>([]);
  const [active, setActive] = useState<any>('');
  const [activeInfo, setActiveInfo] = useState<any>();
  const navigate = useNavigate();
  const [datas, setDatas] = useState<any>({});
  const { message } = App.useApp();
  const interls = useRef<any>(null);
  const formatMessage = useTranslate();
  const [errorState, setErrorState] = useState<any>(false);

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
        getFileDetail(instanceInfo.id).then((dashboardInfo: any) => {
          navigate(
            `/dashboards/preview?${getSearchParamsString({ id: dashboardInfo.id, type: dashboardInfo.type, status: 'preview', name: dashboardInfo.name })}`
          );
        });
      });
    } else {
      return bindFlowForUns({
        flowId: item.id,
        unsAlias: instanceInfo.alias,
      }).then(() => {
        message.success(formatMessage('common.optsuccess'));
        getFileDetail(instanceInfo.id);
        fetchNodeRedData(instanceInfo.alias).then((datas: any) => {
          if (datas) {
            navigate(
              `/collection-flow/flow-editor?${getSearchParamsString({
                id: datas.id,
                name: datas.flowName,
                status: datas.flowStatus,
                flowId: datas.flowId,
                from: location.pathname,
              })}`
            );
          }
        });
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
  // 更新状态
  const updateNodeState = (edgeId: string, status = false, marked = true) => {
    if (!topologyRef.current) return;
    const edge: any = topologyRef.current.getCellById(edgeId);
    if (status) {
      edge?.attr('line/stroke', 'red'); // 更新边的颜色
      edge?.attr('line/style/animation', ''); // 清除动画
      edge?.setLabels(markupLine); //追加label
    } else {
      edge?.attr('line/stroke', 'var(--supos-theme-color)'); // 更新边的颜色
      edge?.attr('line/style/animation', 'ant-line 60s infinite linear'); // 清除动画
      edge?.setLabels({}); //追加label
    }
    // 是否去除叉叉
    if (!marked) {
      edge?.setLabels({}); //追加label
    }
  };

  const getTopologyState = async () => {
    getTopologyStatus({
      id: instanceInfo?.id || '',
    }).then((res: any) => {
      const flag = res.filter((item: any) => item.eventCode != '0');
      if (flag.length > 0) {
        setErrorState(true);
      } else {
        setErrorState(false);
      }
      const sortedData = res.sort(compareByXxOrder);
      for (let i = 0; i < sortedData.length; i++) {
        if (sortedData[i].eventCode !== '0') {
          for (let j = i + 1; j < sortedData.length; j++) {
            sortedData[j].marked = true;
          }
          break; // 找到第一个eventCode不等于0的元素并处理后就可以退出循环了
        }
      }
      //合并pullMqtt和dataPersistence的报错信息
      const errNode = sortedData.find((e: any) => e.eventCode !== '0');
      if (errNode && ['pullMqtt', 'dataPersistence'].includes(errNode.topologyNode)) {
        sortedData.push({ ...errNode, topologyNode: 'pullMqttOrDataPersistence' });
      } else {
        sortedData.push({
          topologyNode: 'pullMqttOrDataPersistence',
          eventCode: '0',
          eventMessage: null,
          eventTime: null,
          marked: errNode ? true : false,
        });
      }
      sortedData.push({
        topologyNode: TypeEnum.Apps + '1',
        eventCode: '0',
        eventMessage: null,
        eventTime: null,
        marked: errNode ? true : false,
      });
      modeState.current = sortedData;
      data.edges.map((item: any) => {
        sortedData.map((item2: any) => {
          if (item.id === item2.topologyNode && item2.eventCode != '0' && !item2.marked) {
            updateNodeState(item.id, true, true);
            clearInterval(interls.current);
          } else if (item.id === item2.topologyNode && item2.eventCode == '0' && item2.marked) {
            updateNodeState(item.id, true, false);
          } else if (item.id === item2.topologyNode && item2.eventCode == '0' && !item2.marked) {
            updateNodeState(item.id, false, true);
          }
        });
      });
      return sortedData;
    });
  };

  // 获取拓扑状态
  const getTopologyStateData = async () => {
    try {
      await getTopologyState();
    } catch (error) {
      console.error('Error fetching topology state:', error);
    }
  };
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
      setActiveInfo(xx[0]);
      modeState.current = [];
    };
    // 设置事件监听器
    topologyRef.current!.on('node:click', nodeClickFn);
    topologyRef.current!.on('edge:click', edgeClickFn);
    interls.current = setInterval(() => {
      getTopologyStateData();
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
      {/*  错误状态 */}
      {errorState &&
        activeInfo?.eventMessage &&
        ['pushOriginalData', 'pushMqtt', 'pullMqttOrDataPersistence'].includes(active) && (
          <div className={styles['detailTable']} style={{ alignItems: 'center' }}>
            <div className={styles['error']}>
              <img src={error} />
              {/* <span>Your connection encountered an issue during the modeling phase.</span> */}
              <span>{activeInfo?.eventMessage}</span>
            </div>
          </div>
        )}
    </Flex>
  );
};

export default TopologyChart;
