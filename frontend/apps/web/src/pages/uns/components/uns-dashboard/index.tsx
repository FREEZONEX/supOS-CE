import { Divider, Flex } from 'antd';
import Overview from './Overview.tsx';
// import { useUnsContext } from '@/pages/uns/UnsContext.tsx';
import { useDeepCompareEffect } from 'ahooks';
import { Connect, DataConnected, Package } from '@carbon/icons-react';
import { useImmer } from 'use-immer';
import type { OverviewListProps } from './type';
import Icon from '@ant-design/icons';
import PackageTop from '@/components/svg-components/PackageTop.tsx';
import styles from './index.module.scss';
import Functions from './Functions.tsx';
import MQTT from './MQTT.tsx';
import useUnsGlobalWs from '@/pages/uns/useUnsGlobalWs.ts';

const UnsDashboard = () => {
  // const { topologyData } = useUnsContext();
  const { topologyData = {} } = useUnsGlobalWs();
  const [overviewList, setOverviewList] = useImmer<OverviewListProps[]>([
    { key: 'messageInThroughput', label: 'uns.messageIn', icon: <Package size={24} />, value: 0, unit: 'uns.msgUnit' },
    {
      key: 'messageOutThroughput',
      label: 'uns.messageOut',
      icon: <Icon component={PackageTop} style={{ fontSize: 24 }} />,
      value: 0,
      unit: 'uns.msgUnit',
    },
    {
      key: 'allConnections',
      label: 'uns.allConnections',
      icon: <Connect size={24} />,
      value: 0,
    },
    {
      key: 'liveConnections',
      label: 'uns.liveConnections',
      icon: <DataConnected size={24} />,
      value: 0,
    },
  ]);

  useDeepCompareEffect(() => {
    const result: { [key: string]: string } = {};
    Object.keys(topologyData).forEach((key) => {
      result[key.toLowerCase()] = topologyData[key];
    });

    setOverviewList((draft) => {
      return draft.map((item: any) => {
        const key = item.key.split(' ').pop().toLowerCase();
        if (Object.prototype.hasOwnProperty.call(result, key)) {
          return {
            ...item,
            value: result[key],
          };
        }
        return item;
      });
    });
  }, [topologyData]);

  return (
    <div className={styles['unsDashboard']}>
      <Overview overviewList={overviewList} />
      <Divider style={{ background: '#e0e0e0', flexShrink: 0 }} />
      <Flex gap={16}>
        <div className={styles['functions-wrapper']}>
          <Functions />
        </div>
        <div className={styles['mqtt-wrapper']}>
          <MQTT />
        </div>
      </Flex>
    </div>
  );
};
export default UnsDashboard;
