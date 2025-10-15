import { Flex } from 'antd';
import styles from './Topology.module.scss';
import useTranslate from '@/hooks/useTranslate.ts';
import cx from 'classnames';
import ComEllipsis from '@/components/com-ellipsis';
import classNames from 'classnames';
import nodeRed from '@/assets/home-icons/node-red.svg';
import { useBaseStore } from '@/stores/base';
import tdengine from '@/assets/home-icons/tdengine.png';
import timescaleDB from '@/assets/home-icons/timescaleDB.svg';
import postgresql from '@/assets/home-icons/postgresql.svg';
import { register } from '@antv/x6-react-shape';
import { useNavigate } from 'react-router';
import { useEffect, useRef } from 'react';
import { useDeepCompareEffect, useMemoizedFn, useSize } from 'ahooks';
import { data } from './data.ts';
import { Graph } from '@antv/x6';
import { ApplicationWeb } from '@carbon/icons-react';
import ComCopy from '@/components/com-copy';

const Modbus = (nodes: any) => {
  return (
    <div
      className={styles['common-node']}
      style={
        nodes.node.data.number
          ? { '--supos-mode-color': 'var(--supos-theme-color)' }
          : { '--supos-mode-color': 'transparent' }
      }
    >
      Modbus
      <div className={styles['node-count']}>
        <span style={nodes.node.data.number ? { '--supos-status-color': 'var(--supos-theme-color)' } : {}}></span>{' '}
        {nodes.node.data.number}
      </div>
    </div>
  );
};
const Api = (nodes: any) => {
  return (
    <div
      className={styles['common-node']}
      style={
        nodes.node.data.number
          ? { '--supos-mode-color': 'var(--supos-theme-color)' }
          : { '--supos-mode-color': 'transparent' }
      }
    >
      API
      <div className={styles['node-count']}>
        <span style={nodes.node.data.number ? { '--supos-status-color': 'var(--supos-theme-color)' } : {}}></span>
        {nodes.node.data.number}
      </div>
    </div>
  );
};
const Opcua = (nodes: any) => {
  return (
    <div
      className={styles['common-node']}
      style={
        nodes.node.data.number
          ? { '--supos-mode-color': 'var(--supos-theme-color)' }
          : { '--supos-mode-color': 'transparent' }
      }
    >
      OPC UA
      <div className={styles['node-count']}>
        <span style={nodes.node.data.number ? { '--supos-status-color': 'var(--supos-theme-color)' } : {}}></span>{' '}
        {nodes.node.data.number}
      </div>
    </div>
  );
};
const Opcda = (nodes: any) => {
  return (
    <div
      className={styles['common-node']}
      style={
        nodes.node.data.number
          ? { '--supos-mode-color': 'var(--supos-theme-color)' }
          : { '--supos-mode-color': 'transparent' }
      }
    >
      OPC DA
      <div className={styles['node-count']}>
        <span style={nodes.node.data.number ? { '--supos-status-color': 'var(--supos-theme-color)' } : {}}></span>{' '}
        {nodes.node.data.number}
      </div>
    </div>
  );
};
const ICMP = (nodes: any) => {
  return (
    <div
      className={styles['common-node']}
      style={
        nodes.node.data.number
          ? { '--supos-mode-color': 'var(--supos-theme-color)' }
          : { '--supos-mode-color': 'transparent' }
      }
    >
      ICMP
      <div className={styles['node-count']}>
        <span style={nodes.node.data.number ? { '--supos-status-color': 'var(--supos-theme-color)' } : {}}></span>{' '}
        {nodes.node.data.number}
      </div>
    </div>
  );
};
const Relation = (nodes: any) => {
  return (
    <div
      className={styles['common-node']}
      style={
        nodes.node.data.number
          ? { '--supos-mode-color': 'var(--supos-theme-color)' }
          : { '--supos-mode-color': 'transparent' }
      }
    >
      Relation
      <div className={styles['node-count']}>
        <span style={nodes.node.data.number ? { '--supos-status-color': 'var(--supos-theme-color)' } : {}}></span>{' '}
        {nodes.node.data.number}
      </div>
    </div>
  );
};
const Unknown = (nodes: any) => {
  return (
    <div
      className={styles['common-node']}
      style={
        nodes.node.data.number
          ? { '--supos-mode-color': 'var(--supos-theme-color)' }
          : { '--supos-mode-color': 'transparent' }
      }
    >
      Unknown
      <div className={styles['node-count']}>
        <span style={nodes.node.data.number ? { '--supos-status-color': 'var(--supos-theme-color)' } : {}}></span>{' '}
        {nodes.node.data.number}
      </div>
    </div>
  );
};
const NodeRed = () => {
  // const formatMessage = useTranslate();
  return (
    <div className={classNames(styles['common-node'], styles['common-node-hover'])}>
      <img src={nodeRed} alt="" width="28px" />
      {/*{formatMessage('common.nodeRed')}*/}
      Node-Red
    </div>
  );
};
const Mqtt = () => {
  return <div className={styles['common-node']}>MQTT Broker</div>;
};
const StreamProcessing = () => {
  return <div className={styles['common-node']}>Stream Calculation</div>;
};
const TDEngine = () => {
  const dataBaseType = useBaseStore((state) => state.dataBaseType);
  return (
    <div className={styles['common-node']}>
      <img src={dataBaseType.includes('tdengine') ? tdengine : timescaleDB} width="28px" />
      {dataBaseType.includes('tdengine') ? 'TDEngine' : 'TimescaleDB'}
    </div>
  );
};
const PostgreSQL = () => {
  return (
    <div className={styles['common-node']}>
      <img src={postgresql} alt="" width="28px" />
      PostgreSQL
    </div>
  );
};
const Grafana = () => {
  const formatMessage = useTranslate();
  return (
    <div className={classNames(styles['common-node'], styles['common-node-hover'])}>
      <ApplicationWeb size={28} />
      {formatMessage('dashboards.dashboard')}
    </div>
  );
};
// const Gui = () => {
//   return (
//     <div className={classNames(styles['common-node'], styles['common-node-hover'])}>
//       <Icon
//         component={Blend}
//         style={{
//           color: 'var(--supos-text-color)',
//           fontSize: 28,
//         }}
//       />
//       {getIntl('common.gui', 'GUI')}
//     </div>
//   );
// };
// const Apps = () => {
//   return (
//     <div className={classNames(styles['common-node'], styles['common-node-hover'])}>
//       <ApplicationWeb size={28} />
//       {getIntl('common.apps', 'APPs')}
//     </div>
//   );
// };

register({
  shape: 'modbus',
  width: 150,
  height: 40,
  component: Modbus,
});
register({
  shape: 'api',
  width: 150,
  height: 40,
  component: Api,
});
register({
  shape: 'opcua',
  width: 150,
  height: 40,
  component: Opcua,
});
register({
  shape: 'opcda',
  width: 150,
  height: 40,
  component: Opcda,
});
register({
  shape: 'icmp',
  width: 150,
  height: 40,
  component: ICMP,
});
register({
  shape: 'relation',
  width: 150,
  height: 40,
  component: Relation,
});
register({
  shape: 'unknown',
  width: 150,
  height: 40,
  component: Unknown,
});
register({
  shape: 'nodeRed',
  width: 150,
  height: 40,
  component: NodeRed,
});
register({
  shape: 'mqtt',
  width: 150,
  height: 40,
  component: Mqtt,
});
register({
  shape: 'streamProcessing',
  width: 150,
  height: 40,
  component: StreamProcessing,
});
register({
  shape: 'tdEngine',
  width: 150,
  height: 40,
  component: TDEngine,
});
register({
  shape: 'postgreSQL',
  width: 150,
  height: 40,
  component: PostgreSQL,
});
register({
  shape: 'grafana',
  width: 150,
  height: 40,
  component: Grafana,
});
// register({
//   shape: 'apps',
//   width: 150,
//   height: 40,
//   component: Apps,
// });
// register({
//   shape: 'gui',
//   width: 150,
//   height: 40,
//   component: Gui,
// });

const Topology = ({ datas }: any) => {
  const formatMessage = useTranslate();
  const navigate = useNavigate();
  const graphRef = useRef<any>(null);
  const containerRef = useRef<HTMLDivElement>(null);
  const outerContainerRef = useRef<HTMLDivElement>(null);
  const systemInfo = useBaseStore((state) => state.systemInfo);
  const wsPort = systemInfo?.mqttTcpPort ?? window.location.port;
  const datahandle = useMemoizedFn(() => {
    if (!datas) return;
    const edges = graphRef.current.getEdges();
    data.nodes.forEach((node: any) => {
      const nodeId = node.id;
      if (nodeId in (datas as Record<string, any>)) {
        // ws获取过来的数据
        const valueFromDatas = (datas as Record<string, any>)[nodeId];
        const graphNode = graphRef.current.getCellById(nodeId);

        if (graphNode) {
          // 更新node数据
          graphNode.setData({
            number: valueFromDatas,
          });
        }
        const graphEdge = edges.find((edge: any) => edge.source.cell === nodeId);

        if (graphEdge) {
          // 更新edge动画
          graphEdge.setAttrs({
            line: {
              style:
                valueFromDatas > 0
                  ? {
                      animation: 'ant-line 60s infinite linear',
                    }
                  : {},
            },
          });
        }
      }
    });
  });

  useEffect(() => {
    if (!containerRef.current) return;
    graphRef.current = new Graph({
      container: containerRef.current,
      background: {
        color: 'var(--supos-card-bg)',
      },
      interacting: false,
      panning: true,
      mousewheel: {
        enabled: true,
        modifiers: ['ctrl', 'meta'],
      },
      scaling: {
        min: 0.05, // 默认值为 0.01
        max: 12, // 默认值为 16
      },
    });
    graphRef.current.fromJSON(data);
    // 位置设置 靠左，修正了原代码中graph变量未定义的问题，使用graphRef.current
    graphRef.current.positionPoint({ x: 266, y: 0 }, '8%', '12%');
    const clickFn = ({ cell }: any) => {
      switch (cell.id) {
        case 'nodeRed':
          navigate('/collection-flow');
          break;
        case 'grafana':
          window.open('/grafana/home/dashboards/');
          break;
        case 'gui':
          navigate('/app-gui');
          break;
        case 'apps':
          navigate('/app-display');
          break;
        default:
          break;
      }
    };

    // 为节点绑定点击事件
    graphRef.current.on('node:click', clickFn);
    return () => {
      if (!containerRef.current) return;
      if (graphRef.current) {
        graphRef.current?.dispose?.();
        graphRef.current.off('node:click', clickFn);
        graphRef.current = null;
      }
    };
  }, []);

  useDeepCompareEffect(() => {
    datahandle();
  }, [datahandle, datas]);
  const size = useSize(outerContainerRef);

  useEffect(() => {
    if (size && containerRef.current) {
      const height = Math.max(size.height, 400);
      if (graphRef.current) {
        setTimeout(() => {
          graphRef.current.resize(size.width, height);
          graphRef.current.positionPoint({ x: 266, y: 0 }, '8%', '12%');
        }, 0);
      }
    }
  }, [size]);

  const mqttList = [
    {
      key: 'url',
      label: 'uns.MQTTUrl',
      text: `mqtt://${window.location.hostname}`,
    },
    {
      key: 'port',
      label: 'uns.MQTTPort',
      text: wsPort,
    },
    {
      key: 'dependent',
      label: 'uns.dependent',
      text: 'npm install mqtt',
    },
    {
      key: 'payload',
      label: 'uns.payload',
      text: JSON.stringify(
        {
          Name: 'Name',
          CurrentValue: 1,
          InitialValue: 1,
        },
        null,
        2
      ),
    },
  ];

  return (
    <Flex gap={16} className={styles['item-wrapper']} style={{ height: '100%' }} ref={outerContainerRef}>
      <Flex vertical className={cx(styles['item'], styles['item-left'])} gap={16} style={{ overflow: 'hidden' }}>
        <ComEllipsis className={styles['title']}>{formatMessage('uns.topologyMap')}</ComEllipsis>
        <div ref={containerRef} style={{ minHeight: 400, flex: 1 }} />
      </Flex>
      <Flex vertical className={cx(styles['item'], styles['item-right'])} gap={16}>
        <ComEllipsis className={styles['title']}>{formatMessage('uns.mqttAccess')}</ComEllipsis>
        <div style={{ flex: 1, background: 'var(--supos-card-bg)', padding: 16, overflow: 'auto' }}>
          {mqttList?.map((item: any) => {
            return (
              <div key={item.key}>
                <ComEllipsis style={{ fontWeight: 500, fontSize: 16 }}>{formatMessage(item.label)}</ComEllipsis>
                <Flex
                  title={item.text}
                  style={{
                    background: 'var(--supos-bg-color)',
                    padding: '8px 16px',
                    margin: '12px 0',
                    borderRadius: '3px',
                  }}
                  align="center"
                  justify="space-between"
                >
                  <pre>{item.text}</pre>
                  <ComCopy textToCopy={item.text} />
                </Flex>
              </div>
            );
          })}
        </div>
      </Flex>
    </Flex>
  );
};

export default Topology;
